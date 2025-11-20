# 报告模式快速入门

## 三种模式简介

TrendHub 提供三种不同的工作模式，选择适合你的使用场景：

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| 🗓️ **daily** | 全天持续收集，推送时汇总 | 每日新闻总结 |
| ⚡ **current** | 实时爬取当前榜单 | 实时热点监控 |
| 📈 **incremental** | 只推送新内容，避免重复 | 长期跟踪特定话题 |

## 快速配置

### 1. 每日汇总模式 (推荐)

**场景**: 每天晚上收到全天的新闻汇总

编辑 `config.yaml`:

```yaml
report:
  mode: daily                    # 当日汇总模式
  rank_threshold: 10             # 只包含排名前10的

crawler:
  request_interval: 600000       # 后台每10分钟爬取一次 (毫秒)

notification:
  enable_notification: true
  push_window:
    enabled: true
    time_range:
      start: "18:00"             # 每天18点推送
      end: "18:10"
    once_per_day: true           # 每天只推一次
```

**启动**:
```bash
./trendhub -web
```

**效果**: 
- ✅ 后台自动每10分钟爬取数据
- ✅ 18:00 推送全天收集的所有匹配新闻
- ✅ 自动去重，同一条新闻只记录一次
- ✅ 次日0点自动清空，开始新的一天

### 2. 实时榜单模式

**场景**: 每小时查看一次当前热搜 TOP 5

编辑 `config.yaml`:

```yaml
report:
  mode: current                  # 当前榜单模式
  rank_threshold: 5              # 只推送前5名

notification:
  push_window:
    enabled: true
    time_range:
      start: "09:00"
      end: "18:00"
    once_per_day: false          # 可以多次推送
```

**配置 crontab** (每小时运行):
```bash
0 * * * * cd /path/to/trendhub && ./trendhub
```

**效果**:
- ✅ 每小时推送当前 TOP 5 热搜
- ✅ 实时获取最新排行
- ✅ 快速响应热点变化

### 3. 增量监控模式

**场景**: 持续监控某个话题，只推送新内容

编辑 `config.yaml`:

```yaml
report:
  mode: incremental              # 增量模式
  rank_threshold: 20             # 监控前20名

crawler:
  request_interval: 1000         # 请求延迟1秒

notification:
  push_window:
    enabled: true
    time_range:
      start: "08:00"
      end: "22:00"
    once_per_day: false
```

**配置 crontab** (每15分钟检查):
```bash
*/15 * * * * cd /path/to/trendhub && ./trendhub
```

**效果**:
- ✅ 每15分钟检查一次
- ✅ 只推送新出现的匹配内容
- ✅ 避免重复推送相同新闻
- ✅ 自动记录推送历史（保留7天）

## 模式切换

直接修改配置文件中的 `report.mode` 即可：

```yaml
report:
  mode: daily        # 改为 daily、current 或 incremental
```

重启程序使配置生效：

```bash
# 如果是 Web 模式
Ctrl+C  # 停止
./trendhub -web  # 重新启动

# 如果是 cron 任务
crontab -e  # 编辑定时任务
```

## 日志查看

### 查看当前模式
```bash
# 日志中会显示当前模式
tail -f trendhub.log | grep "Mode:"
```

输出示例：
```
Mode: Daily aggregation - using cached data
Mode: Current ranking - fetching real-time data  
Mode: Incremental monitoring - fetching and filtering new items
```

### daily 模式日志
```
Daily collector started, collecting data every 10m0s
Collecting data for daily aggregation...
Collected data: added 15 new items, total cached: 128 items
```

### incremental 模式日志
```
Found 8 new items (total: 150, already pushed: 142)
Marked 8 items as pushed
```

## 数据文件

程序会创建以下数据文件：

```bash
data/
├── push_records.db      # 推送记录（所有模式）
└── data_cache.db        # 数据缓存（daily & incremental 模式）
```

- **push_records.db**: 记录推送历史，用于 Web 界面查看
- **data_cache.db**: 
  - daily 模式: 不使用磁盘缓存（内存中）
  - incremental 模式: 记录已推送内容（持久化）

## 常见问题

### Q: daily 模式下推送内容为空？

**A**: 确保使用 Web 模式运行，才会启动后台收集器：
```bash
./trendhub -web  # 正确 ✅
./trendhub       # 错误 ❌ (单次运行无法持续收集)
```

### Q: incremental 模式重复推送？

**A**: 检查数据库文件是否正常：
```bash
ls -lh data/data_cache.db
```

如果文件不存在或损坏，删除后重建：
```bash
rm data/data_cache.db
./trendhub -web
```

### Q: 如何清空历史记录？

**A**: 
```bash
# 清空所有记录
rm data/*.db

# 只清空增量模式的记录
rm data/data_cache.db
```

### Q: 三种模式可以同时用吗？

**A**: 不可以，同一时间只能使用一种模式。但可以：
- 部署多个实例，每个实例使用不同模式
- 根据需求切换模式

## 性能建议

### daily 模式
- ✅ 爬取间隔：5-30 分钟
- ✅ 推送时间：下班时间（如18:00）
- ✅ 内存占用：中等（一天数据量决定）

### current 模式
- ✅ 推送频率：1-4 小时一次
- ✅ 适合：实时监控、快速响应
- ✅ 资源占用：低

### incremental 模式
- ✅ 检查频率：10-30 分钟
- ✅ 适合：长期跟踪、避免重复
- ✅ 磁盘占用：中等

## 最佳实践

1. **新用户**: 先用 current 模式测试，确认关键词配置正确
2. **日常使用**: 使用 daily 模式，每天收到汇总报告
3. **重要监控**: 使用 incremental 模式，避免遗漏新内容
4. **实时追踪**: 使用 current 模式，快速响应突发事件

## 完整示例

### 场景：科技媒体日报

**需求**: 每天 19:00 推送科技领域的全天新闻

**配置**:
```yaml
report:
  mode: daily
  rank_threshold: 15

crawler:
  request_interval: 600000  # 10分钟

notification:
  enable_notification: true
  push_window:
    enabled: true
    time_range:
      start: "19:00"
      end: "19:05"
    once_per_day: true
  webhooks:
    bark_device_key: "your_device_key"  # 推送到手机
```

**关键词** (`frequency_words.txt`):
```
AI
人工智能
ChatGPT
大模型

芯片
半导体

iPhone
华为
```

**运行**:
```bash
./trendhub -web
```

**效果**: 每天 19:00 收到包含所有科技相关新闻的推送。

## 进阶功能

### 组合使用

可以在不同服务器上部署多个实例：

**服务器 A**: daily 模式（每日汇总）
```yaml
report:
  mode: daily
notification:
  webhooks:
    bark_device_key: "device_key_1"  # 推送到个人手机
```

**服务器 B**: incremental 模式（紧急监控）
```yaml
report:
  mode: incremental
notification:
  webhooks:
    feishu_url: "webhook_url"  # 推送到工作群
```

### 定时切换模式

使用脚本在不同时段切换模式：

```bash
#!/bin/bash
# auto-switch-mode.sh

hour=$(date +%H)

if [ $hour -ge 9 ] && [ $hour -lt 18 ]; then
    # 工作时间：增量模式
    sed -i 's/mode: .*/mode: incremental/' config.yaml
else
    # 非工作时间：当日汇总模式
    sed -i 's/mode: .*/mode: daily/' config.yaml
fi

./trendhub
```

## 技术支持

详细文档：[报告模式详解](docs/REPORT_MODES.md)

如有问题：
- 查看日志文件
- 检查配置文件
- 提交 GitHub Issues

---

**快速上手**: 选择一个模式 → 修改配置 → 启动程序 → 查看推送 🎉

