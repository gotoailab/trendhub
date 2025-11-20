# 三种模式测试示例

这个目录包含三种工作模式的配置示例，用于快速测试和理解每种模式的特点。

## 文件说明

- `config_daily.yaml` - 当日汇总模式配置
- `config_current.yaml` - 当前榜单模式配置  
- `config_incremental.yaml` - 增量监控模式配置
- `test_all_modes.sh` - 自动测试脚本

## 快速测试

### 1. 测试 daily 模式

```bash
# 启动 Web 模式（必须）
../../trendhub -web -config config_daily.yaml -keywords ../../config/frequency_words.txt

# 观察日志
# 应该看到: "Daily collector started"
# 等待几分钟后: "Collected data: added X new items"
```

**预期行为**:
- 后台收集器自动启动
- 每隔一段时间自动爬取数据
- 缓存不断累积（可在日志中查看）
- 推送时使用缓存的汇总数据

### 2. 测试 current 模式

```bash
# 单次运行即可
../../trendhub -config config_current.yaml -keywords ../../config/frequency_words.txt

# 观察日志
# 应该看到: "Mode: Current ranking - fetching real-time data"
```

**预期行为**:
- 实时爬取当前榜单
- 立即处理和推送
- 无缓存机制

### 3. 测试 incremental 模式

```bash
# 第一次运行
../../trendhub -config config_incremental.yaml -keywords ../../config/frequency_words.txt

# 观察日志
# 应该看到: "Found X new items (total: Y, already pushed: 0)"

# 立即第二次运行
../../trendhub -config config_incremental.yaml -keywords ../../config/frequency_words.txt

# 观察日志
# 应该看到: "Found 0 new items (total: Y, already pushed: Y)"
# 不会推送相同内容
```

**预期行为**:
- 第一次: 推送所有匹配的内容
- 第二次: 不推送（已推送过）
- 等待新内容出现后才会推送

## 自动化测试

运行测试脚本：

```bash
chmod +x test_all_modes.sh
./test_all_modes.sh
```

这个脚本会：
1. 测试三种模式的配置文件
2. 验证 daily 模式的后台收集器
3. 验证 incremental 模式的去重功能
4. 生成测试报告

## 配置说明

### daily 模式关键配置

```yaml
report:
  mode: daily              # 当日汇总

crawler:
  request_interval: 300000 # 5分钟爬取一次

notification:
  push_window:
    once_per_day: true     # 每天只推一次
```

### current 模式关键配置

```yaml
report:
  mode: current            # 当前榜单

crawler:
  request_interval: 1000   # 请求间隔1秒

notification:
  push_window:
    once_per_day: false    # 可以多次推送
```

### incremental 模式关键配置

```yaml
report:
  mode: incremental        # 增量监控

notification:
  push_window:
    once_per_day: false    # 发现新内容就推
```

## 日志观察要点

### daily 模式
```
✅ Daily collector started, collecting data every 5m0s
✅ Collecting data for daily aggregation...
✅ Collected data: added 15 new items, total cached: 128 items
✅ Mode: Daily aggregation - using cached data
✅ Retrieved 128 items from daily cache
```

### current 模式
```
✅ Mode: Current ranking - fetching real-time data
✅ Crawled data from 11 platforms
✅ Filtered data: 23 items remaining
✅ Sending notifications for 23 items...
```

### incremental 模式
```
✅ Mode: Incremental monitoring - fetching and filtering new items
✅ Found 8 new items (total: 150, already pushed: 142)
✅ Marked 8 items as pushed
✅ Cleaned 5 expired push records
```

## 常见问题

### Q: daily 模式没有启动收集器？

**A**: 确保使用 `-web` 参数：
```bash
../../trendhub -web -config config_daily.yaml
```

### Q: incremental 模式重复推送？

**A**: 检查数据库文件：
```bash
ls -lh ../../data/data_cache.db
```

如果文件不存在或损坏，删除后重试：
```bash
rm -f ../../data/data_cache.db
```

### Q: 如何清空测试数据？

**A**: 
```bash
# 清空所有数据
rm -rf ../../data/*.db

# 只清空增量记录
rm -f ../../data/data_cache.db
```

## 模拟场景测试

### 场景 1: 每日新闻汇总

**目标**: 模拟收集一天的新闻，晚上18点推送

1. 启动 daily 模式（Web）
2. 观察后台收集日志
3. 等待推送时间到达
4. 查看推送内容

### 场景 2: 实时热点追踪

**目标**: 每小时查看一次TOP5热搜

1. 修改 config_current.yaml 的 rank_threshold 为 5
2. 配置 crontab 每小时运行
3. 观察每次推送的内容变化

### 场景 3: 长期话题监控

**目标**: 监控特定关键词，只推新内容

1. 配置关键词文件
2. 启动 incremental 模式
3. 第一次运行记录所有匹配内容
4. 后续只推送新出现的内容

## 性能测试

### 内存占用测试

```bash
# 启动 daily 模式
../../trendhub -web -config config_daily.yaml &
PID=$!

# 监控内存使用
watch -n 60 "ps aux | grep $PID | grep -v grep"

# 观察缓存增长对内存的影响
```

### 磁盘占用测试

```bash
# 运行 incremental 模式多次
for i in {1..10}; do
    ../../trendhub -config config_incremental.yaml
    sleep 60
done

# 检查数据库大小
ls -lh ../../data/data_cache.db
```

## 下一步

- 查看 [报告模式详解](../../docs/REPORT_MODES.md)
- 查看 [模式快速入门](../../MODES_QUICKSTART.md)
- 根据实际需求选择合适的模式
- 调整配置参数优化性能

## 反馈

测试中发现问题？请提交 GitHub Issues。

