# TrendHub 报告模式详解

## 概述

TrendHub 支持三种不同的报告模式，每种模式针对不同的使用场景和需求：

1. **当日汇总模式 (daily)** - 持续收集一天内的所有匹配新闻
2. **当前榜单模式 (current)** - 推送实时热搜榜单
3. **增量监控模式 (incremental)** - 只推送新出现的内容

## 模式详解

### 1. 当日汇总模式 (daily)

#### 适用场景
- 需要了解一天内所有相关热点
- 希望获得完整的新闻汇总
- 定时推送每日总结报告

#### 工作原理

```
启动程序
    ↓
启动后台数据收集器
    ↓
每隔一段时间自动爬取 (根据 request_interval)
    ↓
去重后加入当日缓存
    ↓
推送时使用缓存的所有数据
    ↓
次日 0 点自动清空缓存
```

#### 特点
✅ **完整性** - 不会遗漏任何匹配的新闻  
✅ **去重** - 同一条新闻只记录一次  
✅ **排名优化** - 保留每条新闻的最好排名  
✅ **自动重置** - 每天 0 点自动清空，开始新的一天  

#### 配置示例

```yaml
report:
  mode: daily              # 设置为 daily 模式
  rank_threshold: 10       # 只包含排名前10的新闻

crawler:
  request_interval: 300000 # 5分钟爬取一次 (毫秒)
  
notification:
  enable_notification: true
  push_window:
    enabled: true
    time_range:
      start: "18:00"       # 每天18:00推送当日汇总
      end: "18:30"
    once_per_day: true
```

#### 使用建议
- **推送时间**: 建议设置在下班时间（如 18:00），获取全天汇总
- **爬取间隔**: 建议 5-30 分钟，平衡时效性和服务器负担
- **排名阈值**: 设置合理的 `rank_threshold`，避免数据过多

### 2. 当前榜单模式 (current)

#### 适用场景
- 需要实时热搜信息
- 关注当下最热门的话题
- 快速响应突发事件

#### 工作原理

```
触发推送
    ↓
实时爬取当前榜单
    ↓
关键词过滤
    ↓
立即推送
```

#### 特点
✅ **实时性** - 推送最新的热搜榜单  
✅ **简单直接** - 无需维护缓存  
✅ **即时响应** - 适合突发新闻监控  

#### 配置示例

```yaml
report:
  mode: current           # 设置为 current 模式
  rank_threshold: 5       # 只推送前5名的新闻

notification:
  push_window:
    enabled: true
    time_range:
      start: "09:00"
      end: "22:00"
    once_per_day: false   # 可以多次推送
```

#### 使用建议
- **推送频率**: 可以设置多次推送，实时跟踪热点变化
- **排名阈值**: 设置较小的值，聚焦TOP热点
- **适合**: 新闻媒体、舆情监控等需要实时响应的场景

### 3. 增量监控模式 (incremental)

#### 适用场景
- 避免重复推送相同内容
- 只关注新出现的热点
- 长期持续监控特定关键词

#### 工作原理

```
触发推送
    ↓
爬取当前数据
    ↓
与历史推送记录对比
    ↓
过滤出未推送的新内容
    ↓
推送新内容
    ↓
记录已推送 (保留7天)
```

#### 特点
✅ **避免重复** - 同一条新闻只推送一次  
✅ **持续监控** - 适合长期跟踪  
✅ **智能过滤** - 自动记录推送历史  
✅ **自动清理** - 7天后自动清理旧记录  

#### 配置示例

```yaml
report:
  mode: incremental       # 设置为增量模式
  rank_threshold: 20      # 监控排名前20的新闻

crawler:
  request_interval: 1000  # 爬取间隔（用于请求延迟）

notification:
  push_window:
    enabled: true
    time_range:
      start: "08:00"
      end: "22:00"
    once_per_day: false   # 发现新内容立即推送
```

#### 使用建议
- **推送频率**: 可以频繁检查，只有新内容才会推送
- **记录保留**: 默认保留 7 天，可根据需求调整
- **适合**: 长期跟踪特定话题、竞品监控等场景

## 模式对比

| 特性 | daily | current | incremental |
|------|-------|---------|-------------|
| **实时性** | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **完整性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **去重** | ✅ 自动 | ❌ 无 | ✅ 智能 |
| **内存占用** | 中 | 低 | 中 |
| **磁盘占用** | 低 | 无 | 中 |
| **推送频率** | 每日一次 | 可多次 | 发现新内容 |
| **适用场景** | 日报汇总 | 实时监控 | 长期跟踪 |

## 数据存储

### daily 模式
- **内存缓存**: 当日数据存储在内存中
- **自动重置**: 每天 0 点清空
- **持久化**: 不持久化到磁盘（重启后丢失当日数据）

### incremental 模式
- **持久化数据库**: `data/data_cache.db` (BoltDB)
- **记录内容**: 已推送新闻的哈希值
- **保留时间**: 7 天后自动清理
- **去重依据**: 标题 + 平台ID

## 配置参数详解

### report.mode
```yaml
report:
  mode: "daily"           # daily / current / incremental
  rank_threshold: 10      # 排名阈值（适用于所有模式）
```

### crawler.request_interval
```yaml
crawler:
  request_interval: 5000  # 毫秒
  # daily 模式: 控制后台爬取间隔（建议5-30分钟）
  # current 模式: 控制请求之间的延迟（建议1-2秒）
  # incremental 模式: 控制请求之间的延迟（建议1-2秒）
```

## 使用示例

### 示例 1: 每日新闻汇总

**需求**: 每天晚上 18:00 收到当天所有相关新闻的汇总

```yaml
report:
  mode: daily
  rank_threshold: 10

crawler:
  request_interval: 600000  # 10分钟爬取一次

notification:
  enable_notification: true
  push_window:
    enabled: true
    time_range:
      start: "18:00"
      end: "18:10"
    once_per_day: true
```

**效果**: 程序在后台每 10 分钟爬取一次数据，到 18:00 时推送当天收集的所有匹配新闻。

### 示例 2: 实时热点监控

**需求**: 工作时间内每小时推送一次当前热搜

```yaml
report:
  mode: current
  rank_threshold: 5

notification:
  push_window:
    enabled: true
    time_range:
      start: "09:00"
      end: "18:00"
    once_per_day: false
    
# 配置定时任务（cron）每小时运行一次
```

**效果**: 每小时爬取并推送当前 TOP 5 热搜。

### 示例 3: 持续增量监控

**需求**: 发现新的匹配内容立即推送，避免重复

```yaml
report:
  mode: incremental
  rank_threshold: 20

notification:
  push_window:
    enabled: true
    time_range:
      start: "08:00"
      end: "22:00"
    once_per_day: false

# 配置定时任务（cron）每15分钟运行一次
```

**效果**: 每 15 分钟检查一次，只推送新出现的匹配新闻。

## 命令行使用

### Web 模式（推荐）

```bash
# 启动 Web 模式（自动根据配置启动相应模式）
./trendhub -web

# daily 模式会自动启动后台收集器
# 日志会显示: "Daily mode: continuous data collection started"
```

### 单次运行

```bash
# 直接运行（根据配置的 mode 执行）
./trendhub

# daily 模式: 推送当前缓存的数据（如果缓存为空则只爬取一次）
# current 模式: 爬取并推送当前榜单
# incremental 模式: 爬取并推送新内容
```

## 日志输出

### daily 模式
```
Daily collector started, collecting data every 10m0s
Collecting data for daily aggregation...
Collected data: added 15 new items, total cached: 128 items
Mode: Daily aggregation - using cached data
Retrieved 128 items from daily cache
```

### current 模式
```
Mode: Current ranking - fetching real-time data
Crawled data from 11 platforms
Filtered data: 23 items remaining
```

### incremental 模式
```
Mode: Incremental monitoring - fetching and filtering new items
Crawled data from 11 platforms
Found 8 new items (total: 150, already pushed: 142)
Marked 8 items as pushed
Cleaned 5 expired push records
```

## 故障排查

### daily 模式不工作

**症状**: 推送时没有数据

**检查**:
1. 确认后台收集器是否启动：查找日志 "Daily collector started"
2. 检查缓存数量：查找 "total cached: X items"
3. 确认 Web 模式是否正在运行

**解决**:
```bash
# 必须使用 Web 模式才能启用后台收集
./trendhub -web
```

### incremental 模式重复推送

**症状**: 相同内容被多次推送

**检查**:
1. 检查数据库文件是否正常：`ls -lh data/data_cache.db`
2. 查看日志中的 "Marked X items as pushed"

**解决**:
```bash
# 清空增量记录数据库（谨慎操作）
rm data/data_cache.db
# 重新启动程序
./trendhub -web
```

### 数据库文件过大

**症状**: `data_cache.db` 文件很大

**原因**: incremental 模式会记录所有推送历史

**解决**:
- 默认 7 天后自动清理
- 手动清理：删除数据库文件后重启

## 性能优化

### daily 模式
- **爬取间隔**: 不要设置过短（建议 ≥ 5 分钟）
- **推送时间**: 避开高峰时段
- **内存使用**: 一天内缓存数据较多时会占用较多内存

### incremental 模式  
- **定期清理**: 自动清理 7 天前的记录
- **数据库大小**: 长期运行建议定期备份和重建
- **推送频率**: 可以频繁检查，只有新内容才推送

## 最佳实践

1. **daily 模式**
   - ✅ 配合定时推送使用
   - ✅ 设置在下班时间推送日报
   - ✅ 合理设置爬取间隔

2. **current 模式**
   - ✅ 用于实时监控场景
   - ✅ 可以多次推送
   - ✅ 设置较小的 rank_threshold

3. **incremental 模式**
   - ✅ 用于长期持续监控
   - ✅ 定期检查数据库大小
   - ✅ 避免重复打扰

## 技术细节

### 去重算法

**daily 模式**:
```go
hash = MD5(title + "|" + platform_id)
```

**incremental 模式**:
```go
hash = MD5(title + "|" + platform_id)
record_expiry = 7 days
```

### 数据结构

**daily 缓存** (内存):
```go
map[hash]*NewsItem
```

**incremental 记录** (BoltDB):
```json
{
  "hash": "abc123...",
  "title": "新闻标题",
  "source": "微博",
  "pushed_at": "2025-11-20T10:30:00Z",
  "expires_at": "2025-11-27T10:30:00Z"
}
```

## 版本信息

- 功能版本: v2.2.0
- 添加时间: 2025-11-20
- 依赖: go.etcd.io/bbolt v1.4.3+

## 相关文档

- [定时推送文档](PUSH_SCHEDULE.md)
- [快速开始指南](../QUICKSTART_PUSH.md)
- [配置文件示例](../config.example.yaml)

---

如有问题，欢迎提交 Issues 或查看项目文档。

