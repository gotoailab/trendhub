# TrendRadar-Go

TrendRadar 的 Golang 重构版本。

## 功能

- **模块化设计**：配置、爬虫、过滤、排序、推送完全解耦。
- **多平台支持**：支持 NewsNow API 下的多个平台。
- **多渠道推送**：支持飞书、钉钉、企业微信、Telegram、Bark、Ntfy、邮件等。
- **三种工作模式**：
  - 🗓️ **当日汇总 (daily)**: 持续收集全天数据，定时推送汇总
  - ⚡ **当前榜单 (current)**: 实时爬取推送当前热搜
  - 📈 **增量监控 (incremental)**: 智能去重，只推送新内容
- **配置热重载** ⚡：Web 界面修改配置即时生效，无需重启程序
- **加权排序**：根据排名、频次、热度进行加权排序。
- **配置灵活**：支持 YAML 配置文件和环境变量覆盖。
- **Web 管理界面**：可视化配置管理、任务控制、推送记录查看。
- **定时推送**：支持时间窗口控制和推送记录管理。

## 目录结构

```
.
├── cmd/
│   └── trendradar/      # 入口文件
├── config/              # 配置模块
├── internal/
│   ├── crawler/         # 爬取模块
│   ├── filter/          # 过滤模块
│   ├── rank/            # 排序模块
│   ├── notifier/        # 推送模块
│   └── model/           # 数据模型
└── go.mod
```

## 编译与运行

### 编译

```bash
go mod tidy
go build -o trendradar cmd/trendradar/main.go
```

### 运行

确保 `config/config.yaml` 和 `config/frequency_words.txt` 存在。

#### 命令行模式（单次运行）

```bash
./trendhub
```

或者指定配置文件路径：

```bash
./trendhub -config /path/to/config.yaml -keywords /path/to/frequency_words.txt
```

#### Web 模式（推荐）

启动 Web 服务器和定时调度器：

```bash
./trendhub -web
```

指定端口：

```bash
./trendhub -web -addr :8080
```

然后在浏览器打开：`http://localhost:8080`

Web 界面功能：
- 📊 **仪表盘**：查看运行状态和日志
- ⚙️ **系统配置**：可视化编辑配置文件
- 🏢 **平台管理**：管理监控的热搜平台
- 🏷️ **关键词配置**：可视化管理关键词组
- 📋 **推送记录**：查看历史推送记录

## 工作模式

TrendHub 支持三种工作模式，满足不同的使用场景：

### 🗓️ 当日汇总模式 (daily)
后台持续收集一天内的所有匹配新闻，定时推送汇总报告。
- **适用**: 每日新闻总结、定时日报
- **特点**: 自动去重、完整性高、适合定时推送

### ⚡ 当前榜单模式 (current)  
实时爬取并推送当前热搜榜单。
- **适用**: 实时热点监控、突发事件追踪
- **特点**: 实时性强、即时响应、简单直接

### 📈 增量监控模式 (incremental)
智能记录推送历史，只推送新出现的内容。
- **适用**: 长期跟踪特定话题、避免重复打扰
- **特点**: 智能去重、持续监控、避免重复

**快速上手**: 查看 [模式快速入门](MODES_QUICKSTART.md)  
**详细文档**: 查看 [报告模式详解](docs/REPORT_MODES.md)

## 支持的推送渠道

- 📱 **飞书** - 企业即时通讯
- 📱 **钉钉** - 企业即时通讯
- 📱 **企业微信** - 企业即时通讯
- 📱 **Telegram** - 国际即时通讯 [配置指南](docs/BARK_SETUP.md)
- 📱 **Bark** - iOS 推送通知 [配置指南](docs/BARK_SETUP.md)
- 📱 **Ntfy** - 开源推送服务
- 📧 **邮件** - SMTP 邮件推送

详细配置请参考 [定时推送文档](docs/PUSH_SCHEDULE.md) 和 [快速开始指南](QUICKSTART_PUSH.md)

## 配置热重载 ⚡

在 Web 界面修改配置后，系统会自动重载配置，**无需重启程序**！

### 支持热重载的配置

✅ **推送窗口** - 启用/禁用、时间范围、推送频率  
✅ **工作模式** - daily / current / incremental 切换  
✅ **爬取间隔** - 动态调整数据收集频率  
✅ **推送渠道** - 修改 Webhook URL 和推送配置  
✅ **关键词** - 实时更新监控关键词  

### 使用方法

1. 在 Web 界面修改配置
2. 点击"保存配置"
3. 系统自动应用新配置 ✓

**效果**:
- 修改推送时间 → 调度器自动重启
- 切换到 daily 模式 → 自动启动后台收集
- 关闭推送窗口 → 自动停止调度器
- 修改爬取间隔 → 收集器自动更新

**详细文档**: [配置热重载功能](docs/CONFIG_HOT_RELOAD.md)

## 环境变量

支持使用环境变量覆盖 `config.yaml` 中的配置，例如：

- `FEISHU_WEBHOOK_URL`
- `DINGTALK_WEBHOOK_URL`
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`
- `WEWORK_WEBHOOK_URL`

## 扩展开发

### 添加新的推送渠道

1. 在 `internal/notifier` 下实现 `Notifier` 接口。
2. 在 `internal/notifier/manager.go` 中注册新的 notifier。

### 添加新的爬虫源

1. 在 `internal/crawler` 下实现 `Crawler` 接口。
2. 在 `cmd/trendradar/main.go` 中初始化新的 crawler。

