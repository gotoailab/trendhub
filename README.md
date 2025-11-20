# TrendRadar-Go

TrendRadar 的 Golang 重构版本。

## 功能

- **模块化设计**：配置、爬虫、过滤、排序、推送完全解耦。
- **多平台支持**：支持 NewsNow API 下的多个平台。
- **多渠道推送**：支持飞书、钉钉、Telegram 等。
- **加权排序**：根据排名、频次、热度进行加权排序。
- **配置灵活**：支持 YAML 配置文件和环境变量覆盖。

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

```bash
./trendradar
```

或者指定配置文件路径：

```bash
./trendradar -config /path/to/config.yaml -keywords /path/to/frequency_words.txt
```

## 环境变量

支持使用环境变量覆盖 `config.yaml` 中的配置，例如：

- `FEISHU_WEBHOOK_URL`
- `DINGTALK_WEBHOOK_URL`
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`

## 扩展开发

### 添加新的推送渠道

1. 在 `internal/notifier` 下实现 `Notifier` 接口。
2. 在 `internal/notifier/manager.go` 中注册新的 notifier。

### 添加新的爬虫源

1. 在 `internal/crawler` 下实现 `Crawler` 接口。
2. 在 `cmd/trendradar/main.go` 中初始化新的 crawler。

