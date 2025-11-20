# TrendHub

看你想看的信息，把握趋势算法。

一个强大的热点趋势监控和管理平台，支持多平台数据采集、智能关键词过滤、多渠道推送，并提供现代化的 Web 管理界面。

## ✨ 核心功能

- **模块化设计**：配置、爬虫、过滤、排序、推送完全解耦，易于扩展和维护
- **多平台支持**：支持微博、知乎、百度、今日头条、bilibili、抖音、贴吧、凤凰网、财联社、澎湃新闻、华尔街见闻等多个热门平台
- **多渠道推送**：支持飞书、钉钉、企业微信、Telegram、Bark、Ntfy、邮件等多种推送方式
- **三种工作模式**：
  - 🗓️ **当日汇总 (daily)**: 持续收集全天数据，定时推送汇总
  - ⚡ **当前榜单 (current)**: 实时爬取推送当前热搜
  - 📈 **增量监控 (incremental)**: 智能去重，只推送新内容
- **配置热重载** ⚡：Web 界面修改配置即时生效，无需重启程序
- **加权排序**：根据排名、频次、热度进行智能加权排序
- **现代化 Web 界面**：
  - 🎨 支持深色/浅色模式切换
  - 📊 实时数据展示和历史记录查看
  - ⚙️ 可视化配置管理（表单模式和源码模式）
  - 📋 推送记录查询和分页
  - 🔍 关键词规则说明和帮助
  - 🔄 版本更新自动检测和提示
- **定时推送**：支持时间窗口控制和推送记录管理
- **版本管理**：自动检测新版本并提示更新

## 📁 目录结构

```
.
├── cmd/
│   └── main.go              # 程序入口
├── config/
│   └── config.go            # 配置管理模块
├── internal/
│   ├── collector/           # 数据收集器
│   ├── crawler/             # 爬虫模块
│   ├── datacache/           # 数据缓存
│   ├── filter/              # 关键词过滤模块
│   ├── model/               # 数据模型
│   ├── notifier/            # 推送模块
│   ├── pushdb/              # 推送记录数据库
│   ├── rank/                # 排序模块
│   └── scheduler/           # 定时调度器
├── web/
│   ├── server.go            # Web 服务器
│   ├── runner.go            # 任务运行器
│   └── static/
│       └── index.html       # Web 界面
├── docs/                    # 文档目录
├── examples/                # 示例配置
├── config.example.yaml      # 配置文件示例
├── frequency_words.example.txt  # 关键词文件示例
├── version                  # 版本号文件
└── go.mod
```

## 🚀 快速开始

### 编译

```bash
go mod tidy
go build -o trendhub cmd/main.go
```

### 运行

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

指定端口和数据库路径：

```bash
./trendhub -web -addr :8080 -pushdb data/push_records.db -cachedb data/data_cache.db
```

然后在浏览器打开：`http://localhost:8080`

### Web 界面功能

- 🏠 **首页 (Trending)**：查看实时和历史热点数据，支持日期筛选和自动刷新
- ⚙️ **功能设置**：
  - **系统配置**：可视化编辑配置文件（表单模式/源码模式）
  - **平台管理**：管理监控的热搜平台
  - **关键词配置**：可视化管理关键词组，支持规则说明
- 📋 **推送记录**：查看历史推送记录和执行状态，支持分页查询

### Web 界面特性

- 🎨 **深色模式**：支持深色/浅色模式切换，自动保存用户偏好
- 📱 **响应式设计**：适配各种屏幕尺寸
- 🔄 **实时更新**：Trending 页面支持自动刷新
- 📖 **规则说明**：关键词配置页面提供详细的规则说明
- 🔝 **返回顶部**：滚动时显示返回顶部按钮
- 🔔 **版本更新提示**：自动检测新版本并提示更新

## 📊 工作模式

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

**快速上手**: 查看 [模式快速入门](docs/MODES_QUICKSTART.md)  
**详细文档**: 查看 [报告模式详解](docs/REPORT_MODES.md)

## 🏢 支持的平台

- 微博 (weibo)
- 知乎 (zhihu)
- 百度热搜 (baidu)
- 今日头条 (toutiao)
- bilibili 热搜 (bilibili-hot-search)
- 抖音 (douyin)
- 贴吧 (tieba)
- 凤凰网 (ifeng)
- 财联社热门 (cls-hot)
- 澎湃新闻 (thepaper)
- 华尔街见闻 (wallstreetcn-hot)

## 📱 支持的推送渠道

- 📱 **飞书** - 企业即时通讯
- 📱 **钉钉** - 企业即时通讯
- 📱 **企业微信** - 企业即时通讯
- 📱 **Telegram** - 国际即时通讯
- 📱 **Bark** - iOS 推送通知 [配置指南](docs/BARK_SETUP.md)
- 📱 **Ntfy** - 开源推送服务
- 📧 **邮件** - SMTP 邮件推送

详细配置请参考 [定时推送文档](docs/PUSH_SCHEDULE.md) 和 [快速开始指南](docs/QUICKSTART_PUSH.md)

## 🎯 关键词规则

关键词配置支持三种类型：

- **普通词**（无前缀）：任意匹配，只要标题包含任意一个普通词即可
- **必须词**（+开头）：必须包含，标题必须包含所有必须词
- **过滤词**（!开头）：排除规则，包含过滤词的新闻会被排除

关键词组之间是 OR 关系，组内规则是 AND 关系。在 Web 界面的关键词配置页面可以查看详细的规则说明。

## ⚡ 配置热重载

在 Web 界面修改配置后，系统会自动重载配置，**无需重启程序**！

### 支持热重载的配置

✅ **推送窗口** - 启用/禁用、时间范围、推送频率  
✅ **工作模式** - daily / current / incremental 切换  
✅ **爬取间隔** - 动态调整数据收集频率  
✅ **推送渠道** - 修改 Webhook URL 和推送配置  
✅ **关键词** - 实时更新监控关键词  
✅ **平台列表** - 添加或删除监控平台

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

## 🔧 配置说明

### 配置文件

主要配置文件：`config/config.yaml`

```yaml
app:
  show_version_update: true
  version_check_url: https://raw.githubusercontent.com/gotoailab/trendhub/refs/heads/master/version
crawler:
  enable_crawler: true
  request_interval: 1000
  use_proxy: false
  default_proxy: http://127.0.0.1:10086
report:
  mode: daily
  rank_threshold: 5
notification:
  enable_notification: true
  push_window:
    enabled: true
    time_range:
      start: "20:00"
      end: "22:00"
    once_per_day: true
  webhooks:
    feishu_url: ""
    dingtalk_url: ""
    # ... 其他推送渠道配置
weight:
  rank_weight: 0.6
  frequency_weight: 0.3
  hotness_weight: 0.1
platforms:
  - id: weibo
    name: 微博
  # ... 其他平台
```

### 关键词文件

关键词配置文件：`config/frequency_words.txt`

格式说明：
- 使用空行分隔不同的关键词组
- `!开头` 为过滤词（排除包含该词的结果）
- `+开头` 为必须词（必须包含该词）
- 无前缀为普通词（任意匹配即可）

示例：
```
AI
人工智能
+突破
!广告
```

### 环境变量

支持使用环境变量覆盖 `config.yaml` 中的配置：

- `FEISHU_WEBHOOK_URL`
- `DINGTALK_WEBHOOK_URL`
- `WEWORK_WEBHOOK_URL`
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`

## 🔄 版本更新

TrendHub 支持自动版本检测和更新提示：

1. 在配置文件中设置 `version_check_url` 指向版本文件 URL
2. 启用 `show_version_update` 选项
3. Web 界面启动时会自动检查新版本
4. 如有新版本，会显示更新提示（用户关闭后不会重复提示）

版本文件格式：纯文本，包含版本号（如 `1.0.0`）

## 🛠️ 扩展开发

### 添加新的推送渠道

1. 在 `internal/notifier` 下实现 `Notifier` 接口
2. 在 `internal/notifier/manager.go` 中注册新的 notifier

### 添加新的爬虫源

1. 在 `internal/crawler` 下实现 `Crawler` 接口
2. 在相关代码中初始化新的 crawler

## 📝 许可证

本项目采用开源许可证，详见 LICENSE 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📚 相关文档

- [模式快速入门](docs/MODES_QUICKSTART.md)
- [报告模式详解](docs/REPORT_MODES.md)
- [定时推送文档](docs/PUSH_SCHEDULE.md)
- [快速开始指南](docs/QUICKSTART_PUSH.md)
- [配置热重载功能](docs/CONFIG_HOT_RELOAD.md)
- [Bark 配置指南](docs/BARK_SETUP.md)
