# TrendRadar-Go

TrendRadar 的 Golang 重构版本。

## 功能

- **模块化设计**：配置、爬虫、过滤、排序、推送完全解耦。
- **多平台支持**：支持 NewsNow API 下的多个平台。
- **多渠道推送**：支持飞书、钉钉、企业微信、Telegram、Bark、Ntfy、邮件等。
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

## 支持的推送渠道

- 📱 **飞书** - 企业即时通讯
- 📱 **钉钉** - 企业即时通讯
- 📱 **企业微信** - 企业即时通讯
- 📱 **Telegram** - 国际即时通讯 [配置指南](docs/BARK_SETUP.md)
- 📱 **Bark** - iOS 推送通知 [配置指南](docs/BARK_SETUP.md)
- 📱 **Ntfy** - 开源推送服务
- 📧 **邮件** - SMTP 邮件推送

详细配置请参考 [定时推送文档](docs/PUSH_SCHEDULE.md) 和 [快速开始指南](QUICKSTART_PUSH.md)

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

