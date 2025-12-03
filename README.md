<div align="center">

<svg width="120" height="120" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path stroke="#667eea" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
</svg>

# TrendHub

<p align="center">
  <strong>🎯 看你想看的信息，把握趋势算法</strong>
</p>

<p align="center">
  一个强大的热点趋势监控和管理平台<br/>
  支持多平台数据采集 · 智能关键词过滤 · 多渠道推送 · 现代化 Web 界面
</p>

<p align="center">
  <a href="http://trendhub.wowwwow.cn/">在线 Demo</a> ·
  <a href="#-快速开始">快速开始</a> ·
  <a href="#-核心功能">核心功能</a> ·
  <a href="#-文档">文档</a> ·
  <a href="#-社区">社区</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/license/gotoailab/trendhub?style=flat-square" alt="License">
  <img src="https://img.shields.io/github/stars/gotoailab/trendhub?style=flat-square" alt="Stars">
  <img src="https://img.shields.io/github/forks/gotoailab/trendhub?style=flat-square" alt="Forks">
  <img src="https://img.shields.io/github/issues/gotoailab/trendhub?style=flat-square" alt="Issues">
  <img src="https://img.shields.io/github/go-mod/go-version/gotoailab/trendhub?style=flat-square" alt="Go Version">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square" alt="PRs Welcome">
</p>

---

<img src="https://raw.githubusercontent.com/gotoailab/trendhub/refs/heads/master/docs/images/demo.png" width="800" alt="TrendHub Demo">

</div>

## ✨ 核心功能

<table>
<tr>
<td width="50%">

### 🎯 智能监控
- **多平台支持** - 11+ 热门平台数据采集
- **三种工作模式** - daily/current/incremental
- **智能过滤** - 支持必须词、过滤词、普通词
- **智能排序** - 多维度加权排序算法
- **配置热重载** - 修改配置即时生效，无需重启

</td>
<td width="50%">

### 📱 多渠道推送
- **企业通讯** - 飞书、钉钉、企业微信
- **即时通讯** - Telegram、Bark、Ntfy
- **邮件推送** - SMTP 邮件支持
- **定时推送** - 时间窗口控制
- **推送记录** - 完整的推送历史管理

</td>
</tr>
<tr>
<td width="50%">

### 🎨 现代化界面
- **深色模式** - 支持深色/浅色主题切换
- **响应式设计** - 完美适配各种屏幕
- **实时更新** - 自动刷新热点数据
- **可视化配置** - 表单模式/源码模式切换
- **版本检测** - 自动检测新版本并提示

</td>
<td width="50%">

### 🔧 开发友好
- **模块化设计** - 配置、爬虫、过滤、推送完全解耦
- **易于扩展** - 简单添加新平台和推送渠道
- **完整文档** - 详细的使用和开发文档
- **活跃社区** - 微信群、QQ群、Discord
- **开源协议** - MIT 许可证

</td>
</tr>
</table>

---

## 📋 三种工作模式

<table>
<tr>
<td width="33%" align="center">

### 🗓️ 当日汇总
**daily**

持续收集全天数据  
定时推送汇总报告

*适合每日新闻总结*

</td>
<td width="33%" align="center">

### ⚡ 当前榜单
**current**

实时爬取推送  
即时响应热点

*适合突发事件追踪*

</td>
<td width="33%" align="center">

### 📈 增量监控
**incremental**

智能去重推送  
只推送新内容

*适合长期话题跟踪*

</td>
</tr>
</table>

<p align="center">
  <a href="docs/MODES_QUICKSTART.md">📖 模式快速入门</a> ·
  <a href="docs/REPORT_MODES.md">📚 报告模式详解</a>
</p>

---

## 📁 项目结构

<details>
<summary><b>点击展开目录结构</b></summary>

```
.
├── 📂 cmd/
│   └── main.go                     # 🚀 程序入口
│
├── 📂 config/
│   └── config.go                   # ⚙️ 配置管理模块
│
├── 📂 internal/                    # 核心业务逻辑
│   ├── collector/                  # 📊 数据收集器
│   ├── crawler/                    # 🕷️ 爬虫模块
│   ├── datacache/                  # 💾 数据缓存
│   ├── filter/                     # 🔍 关键词过滤模块
│   ├── model/                      # 📋 数据模型
│   ├── notifier/                   # 📱 推送模块
│   ├── pushdb/                     # 🗄️ 推送记录数据库
│   ├── rank/                       # 📈 排序模块
│   └── scheduler/                  # ⏰ 定时调度器
│
├── 📂 web/                         # Web 服务
│   ├── server.go                   # 🌐 Web 服务器
│   ├── runner.go                   # 🏃 任务运行器
│   └── static/
│       └── index.html              # 🎨 Web 界面
│
├── 📂 docs/                        # 📚 文档目录
├── 📂 examples/                    # 💡 示例配置
│
├── config.example.yaml             # ⚙️ 配置文件示例
├── frequency_words.example.txt     # 🔤 关键词文件示例
├── version                         # 🏷️ 版本号文件
└── go.mod                          # 📦 Go 模块依赖

```

**模块化设计，职责清晰，易于扩展！**

</details>

## 🚀 快速开始

### 方式一：一键运行（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/gotoailab/trendhub.git
cd trendhub

# 2. 准备配置文件
cp config.example.yaml config/config.yaml
cp frequency_words.example.txt config/frequency_words.txt

# 3. 编译并运行
go mod tidy
go build -o trendhub cmd/main.go
./trendhub -web
```

**🎉 完成！** 在浏览器打开 `http://localhost:8080` 即可使用

---

### 方式二：详细配置

<details>
<summary><b>点击展开详细步骤</b></summary>

#### 1️⃣ 编译项目

```bash
go mod tidy
go build -o trendhub cmd/main.go
```

#### 2️⃣ Web 模式（推荐）

启动 Web 服务器和定时调度器：

```bash
./trendhub -web
```

指定端口和数据库路径：

```bash
./trendhub -web -addr :8080 -pushdb data/push_records.db -cachedb data/data_cache.db
```

然后在浏览器打开：`http://localhost:8080`

#### 3️⃣ 命令行模式（单次运行）

```bash
./trendhub
```

或者指定配置文件路径：

```bash
./trendhub -config /path/to/config.yaml -keywords /path/to/frequency_words.txt
```

</details>

---

### 🎨 Web 界面

<table>
<tr>
<td width="33%">

#### 🏠 首页
- 实时热点数据展示
- 历史记录查看
- 日期筛选
- 自动刷新

</td>
<td width="33%">

#### ⚙️ 功能设置
- 可视化配置编辑
- 表单/源码模式切换
- 平台管理
- 关键词配置

</td>
<td width="33%">

#### 📋 推送记录
- 历史推送查询
- 执行状态追踪
- 分页浏览
- 详细日志

</td>
</tr>
</table>

**界面特性：**
✅ 深色/浅色模式切换  
✅ 响应式设计  
✅ 实时数据更新  
✅ 版本更新提示  
✅ 详细的规则说明  
✅ 返回顶部快捷键

---

## 📊 工作模式详解

<details>
<summary><b>🗓️ 当日汇总模式 (daily)</b></summary>

后台持续收集一天内的所有匹配新闻，定时推送汇总报告。

**适用场景：**
- 每日新闻总结
- 定时日报推送
- 全面热点回顾

**特点：**
- ✅ 自动去重
- ✅ 完整性高
- ✅ 适合定时推送
- ✅ 不漏重要信息

</details>

<details>
<summary><b>⚡ 当前榜单模式 (current)</b></summary>

实时爬取并推送当前热搜榜单。

**适用场景：**
- 实时热点监控
- 突发事件追踪
- 快速响应需求

**特点：**
- ✅ 实时性强
- ✅ 即时响应
- ✅ 简单直接
- ✅ 零延迟推送

</details>

<details>
<summary><b>📈 增量监控模式 (incremental)</b></summary>

智能记录推送历史，只推送新出现的内容。

**适用场景：**
- 长期跟踪特定话题
- 避免重复打扰
- 持续关注某个领域

**特点：**
- ✅ 智能去重
- ✅ 持续监控
- ✅ 避免重复
- ✅ 精准推送

</details>

<p align="center">
  📖 <a href="docs/MODES_QUICKSTART.md">模式快速入门</a> ·
  📚 <a href="docs/REPORT_MODES.md">报告模式详解</a>
</p>

## 🌐 支持的平台

<table>
<tr>
<td align="center" width="20%">

📱<br/>**微博**<br/>`weibo`

</td>
<td align="center" width="20%">

🎓<br/>**知乎**<br/>`zhihu`

</td>
<td align="center" width="20%">

🔍<br/>**百度热搜**<br/>`baidu`

</td>
<td align="center" width="20%">

📰<br/>**今日头条**<br/>`toutiao`

</td>
<td align="center" width="20%">

📺<br/>**bilibili**<br/>`bilibili-hot-search`

</td>
</tr>
<tr>
<td align="center" width="20%">

🎵<br/>**抖音**<br/>`douyin`

</td>
<td align="center" width="20%">

💬<br/>**贴吧**<br/>`tieba`

</td>
<td align="center" width="20%">

🦅<br/>**凤凰网**<br/>`ifeng`

</td>
<td align="center" width="20%">

💰<br/>**财联社**<br/>`cls-hot`

</td>
<td align="center" width="20%">

📋<br/>**澎湃新闻**<br/>`thepaper`

</td>
</tr>
<tr>
<td align="center" width="20%">

📈<br/>**华尔街见闻**<br/>`wallstreetcn-hot`

</td>
<td align="center" width="20%" colspan="4">

💡 **持续更新中...**

</td>
</tr>
</table>

---

## 📱 推送渠道

<div align="center">

| 渠道 | 类型 | 配置文档 |
|:---:|:---:|:---:|
| 🚀 **飞书** | 企业即时通讯 | [快速开始](docs/QUICKSTART_PUSH.md) |
| 💬 **钉钉** | 企业即时通讯 | [快速开始](docs/QUICKSTART_PUSH.md) |
| 💼 **企业微信** | 企业即时通讯 | [快速开始](docs/QUICKSTART_PUSH.md) |
| ✈️ **Telegram** | 国际即时通讯 | [快速开始](docs/QUICKSTART_PUSH.md) |
| 🍎 **Bark** | iOS 推送通知 | [配置指南](docs/BARK_SETUP.md) |
| 🔔 **Ntfy** | 开源推送服务 | [快速开始](docs/QUICKSTART_PUSH.md) |
| 📧 **邮件** | SMTP 邮件推送 | [定时推送](docs/PUSH_SCHEDULE.md) |

</div>

<p align="center">
  💡 更多详细配置请参考 <a href="docs/PUSH_SCHEDULE.md">定时推送文档</a> 和 <a href="docs/QUICKSTART_PUSH.md">快速开始指南</a>
</p>

---

## 🎯 关键词规则

<div align="center">

| 类型 | 前缀 | 说明 | 示例 |
|:---:|:---:|:---|:---|
| 🔍 **普通词** | 无 | 任意匹配，包含任意一个即可 | `AI`、`人工智能` |
| ✅ **必须词** | `+` | 必须包含，标题必须包含所有必须词 | `+突破`、`+发布` |
| ❌ **过滤词** | `!` | 排除规则，包含该词会被过滤 | `!广告`、`!营销` |

</div>

### 匹配逻辑

```
关键词组之间 = OR（或）关系
关键词组内部 = AND（且）关系
```

### 示例配置

```
[priority:10]
AI
人工智能
+突破
!广告

[priority:8]
华为
鸿蒙
!谣言
```

**解读：** 标题需包含"AI"或"人工智能"，且必须包含"突破"，但不能包含"广告"

<p align="center">
  💡 在 Web 界面的关键词配置页面可以查看更详细的规则说明
</p>

---

## 🎯 智能排序算法

<div align="center">

**让你最想看的内容排在最前面！**

</div>

### 🏆 关键词优先级

为不同的关键词组设置优先级（1-10），重要的话题自动排前面：

<table>
<tr>
<td width="50%">

```
[priority:10]
AI
人工智能
+突破

[priority:8]
华为
鸿蒙

[priority:5]
科技
互联网
```

</td>
<td width="50%">

- **优先级 10** ⭐⭐⭐  
  最关心的核心话题，分数翻倍

- **优先级 5** ⭐⭐  
  默认优先级，标准权重

- **优先级 1** ⭐  
  低优先级内容，分数减半

</td>
</tr>
</table>

### ⚖️ 平台权重

为不同平台设置权重，提升优质平台内容的排名：

```yaml
platforms:
  - id: zhihu
    name: 知乎
    weight: 1.2  # 优质平台，提升 20%
  - id: weibo
    name: 微博
    weight: 0.9  # 标准平台
  - id: douyin
    name: 抖音
    weight: 0.7  # 娱乐平台，降低权重
```

### ⚙️ 权重配置

调整各项因素的影响权重：

```yaml
weight:
  rank_weight: 0.3       # 原始排名权重
  frequency_weight: 0.2  # 出现频次权重
  keyword_weight: 0.4    # 关键词匹配权重（核心！）
  freshness_weight: 0.1  # 时效性权重
  platform_weight: 1.0   # 平台权重影响系数
```

<p align="center">
  📖 <a href="docs/RANKING_OPTIMIZATION.md">排序算法优化文档</a> ·
  🚀 <a href="docs/RANKING_MIGRATION.md">排序优化迁移指南</a>
</p>

---

## ⚡ 配置热重载

<div align="center">

**在 Web 界面修改配置后，系统自动重载，无需重启程序！**

</div>

### 支持热重载的配置

<table>
<tr>
<td width="50%">

#### ✅ 基础配置
- **推送窗口** - 启用/禁用、时间范围
- **工作模式** - daily/current/incremental
- **爬取间隔** - 动态调整收集频率
- **推送渠道** - Webhook URL 和配置

</td>
<td width="50%">

#### ✅ 高级配置
- **关键词** - 实时更新监控关键词
- **平台列表** - 添加/删除监控平台
- **排序权重** - 动态调整排序算法
- **过滤规则** - 实时更新过滤条件

</td>
</tr>
</table>

### 使用方法

```
1️⃣ 在 Web 界面修改配置
        ↓
2️⃣ 点击"保存配置"按钮
        ↓
3️⃣ 系统自动应用新配置 ✓
```

### 自动响应效果

| 操作 | 系统响应 |
|:---|:---|
| 修改推送时间 | 🔄 调度器自动重启 |
| 切换到 daily 模式 | 🚀 自动启动后台收集 |
| 关闭推送窗口 | ⏸️ 自动停止调度器 |
| 修改爬取间隔 | ⚡ 收集器自动更新 |

<p align="center">
  📖 <a href="docs/CONFIG_HOT_RELOAD.md">配置热重载详细文档</a>
</p>

---

## 🔧 配置说明

<details>
<summary><b>📄 配置文件 (config.yaml)</b></summary>

```yaml
app:
  show_version_update: true
  version_check_url: https://raw.githubusercontent.com/gotoailab/trendhub/refs/heads/master/version

crawler:
  enable_crawler: true          # 启用爬虫
  request_interval: 1000        # 请求间隔（毫秒）
  use_proxy: false              # 是否使用代理
  default_proxy: http://127.0.0.1:10086

report:
  mode: daily                   # 工作模式：daily/current/incremental
  rank_threshold: 5             # 排名阈值

notification:
  enable_notification: true
  push_window:
    enabled: true               # 启用推送窗口
    time_range:
      start: "20:00"           # 推送开始时间
      end: "22:00"             # 推送结束时间
    once_per_day: true         # 每天只推送一次
  webhooks:
    feishu_url: ""             # 飞书 Webhook URL
    dingtalk_url: ""           # 钉钉 Webhook URL
    # ... 其他推送渠道配置

weight:
  rank_weight: 0.6             # 排名权重
  frequency_weight: 0.3        # 频次权重
  hotness_weight: 0.1          # 热度权重

platforms:
  - id: weibo
    name: 微博
    weight: 1.0                # 平台权重
  # ... 其他平台
```

</details>

<details>
<summary><b>🔤 关键词文件 (frequency_words.txt)</b></summary>

**格式说明：**
- 使用空行分隔不同的关键词组
- `!开头` 为过滤词（排除包含该词的结果）
- `+开头` 为必须词（必须包含该词）
- 无前缀为普通词（任意匹配即可）
- `[priority:N]` 设置优先级（1-10）

**示例：**

```
[priority:10]
AI
人工智能
+突破
!广告

[priority:8]
华为
鸿蒙
+发布

[priority:5]
科技
互联网
!营销
```

</details>

<details>
<summary><b>🌍 环境变量</b></summary>

支持使用环境变量覆盖 `config.yaml` 中的配置：

| 环境变量 | 说明 |
|:---|:---|
| `FEISHU_WEBHOOK_URL` | 飞书 Webhook URL |
| `DINGTALK_WEBHOOK_URL` | 钉钉 Webhook URL |
| `WEWORK_WEBHOOK_URL` | 企业微信 Webhook URL |
| `TELEGRAM_BOT_TOKEN` | Telegram Bot Token |
| `TELEGRAM_CHAT_ID` | Telegram Chat ID |

**使用示例：**

```bash
export FEISHU_WEBHOOK_URL="https://open.feishu.cn/..."
./trendhub -web
```

</details>

---

## 🔄 版本更新

<div align="center">

**TrendHub 支持自动版本检测和更新提示**

</div>

### 配置步骤

```yaml
app:
  show_version_update: true
  version_check_url: https://raw.githubusercontent.com/gotoailab/trendhub/refs/heads/master/version
```

### 工作流程

```
Web 界面启动
     ↓
自动检查新版本
     ↓
发现新版本 → 显示更新提示
     ↓
用户关闭提示 → 不再重复提醒
```

**版本文件格式：** 纯文本，包含版本号（如 `1.0.0`）

---

## 🛠️ 扩展开发

<table>
<tr>
<td width="50%">

### 添加新的推送渠道

1️⃣ 在 `internal/notifier` 下实现 `Notifier` 接口

```go
type Notifier interface {
    Send(message string) error
}
```

2️⃣ 在 `internal/notifier/manager.go` 中注册

```go
manager.RegisterNotifier(newNotifier)
```

</td>
<td width="50%">

### 添加新的爬虫源

1️⃣ 在 `internal/crawler` 下实现 `Crawler` 接口

```go
type Crawler interface {
    Fetch() ([]Item, error)
}
```

2️⃣ 在相关代码中初始化

```go
crawler := NewMyCrawler()
```

</td>
</tr>
</table>

<p align="center">
  💡 <b>欢迎提交 PR 添加更多平台和推送渠道！</b>
</p>

---

## 📚 文档

<div align="center">

| 分类 | 文档 | 说明 |
|:---:|:---|:---|
| 🚀 **快速开始** | [快速开始指南](docs/QUICKSTART_PUSH.md) | 5 分钟上手教程 |
| 📊 **工作模式** | [模式快速入门](docs/MODES_QUICKSTART.md) | 三种模式快速了解 |
| 📖 **工作模式** | [报告模式详解](docs/REPORT_MODES.md) | 深入理解各种模式 |
| 📱 **推送配置** | [定时推送文档](docs/PUSH_SCHEDULE.md) | 推送功能详细配置 |
| 🍎 **推送配置** | [Bark 配置指南](docs/BARK_SETUP.md) | iOS Bark 推送配置 |
| ⚡ **高级功能** | [配置热重载功能](docs/CONFIG_HOT_RELOAD.md) | 动态配置更新 |
| 🎯 **高级功能** | [排序算法优化说明](docs/RANKING_OPTIMIZATION.md) | 智能排序详解 ⭐ |
| 🚀 **高级功能** | [排序优化迁移指南](docs/RANKING_MIGRATION.md) | 排序算法迁移 ⭐ |

</div>

---

## 🤝 贡献

<div align="center">

**我们欢迎所有形式的贡献！**

</div>

### 贡献方式

- 🐛 **报告 Bug** - 提交 Issue 描述问题
- 💡 **提出建议** - 分享你的想法和改进建议
- 📝 **完善文档** - 帮助改进文档质量
- 🔧 **提交代码** - 发送 Pull Request
- ⭐ **Star 项目** - 给项目一个 Star 支持我们

### 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

<p align="center">
  <a href="https://github.com/gotoailab/trendhub/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=gotoailab/trendhub" />
  </a>
</p>

---

## 💬 社区

<div align="center">

**加入我们的社区，与其他用户交流！**

</div>

<table align="center">
<tr>
<td align="center" width="33%">

### 💬 微信群

添加微信：**mongorz**  
备注：**simple**  
加入微信交流群

</td>
<td align="center" width="33%">

### 🐧 QQ 群

群号：**823136692**  
[点击加入 QQ 群](https://qun.qq.com/universal-share/share?ac=1&authKey=C4MX6tSrUhKA2xX6M8IosY%2Bb2RyV45O15osUlidAptAwXBgA641FCsENb%2BfiVmki&busi_data=eyJncm91cENvZGUiOiI4MjMxMzY2OTIiLCJ0b2tlbiI6IjJGZmlVUzdsTGRCSytTNGNpajJkQ2F5djYyRUxMbVM1dFN5Y1RXYTRTMG1xejV5N0MrOUM5aEY3aGN4MXVCYmsiLCJ1aW4iOiIzMTc1NDI1NDgwIn0%3D&data=oAc646m4Q0_WTFE6cABXelyZYDd71nAaQ6nA91CQ3A1VNn83sDV5Z-2eXghHkbIXzG16UHCHF4szPTq3A2fhOw&svctype=4&tempid=h5_group_info)

</td>
<td align="center" width="33%">

### 💬 Discord

[点击加入 Discord](https://discord.gg/B3FwBSQq)  
国际社区交流

</td>
</tr>
</table>

---

## 📄 许可证

<div align="center">

本项目采用 **MIT** 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

**⭐ 如果这个项目对你有帮助，请给它一个 Star！⭐**

<br/>

Made with ❤️ by [GotoAILab](https://github.com/gotoailab)

<br/>

[![Star History Chart](https://api.star-history.com/svg?repos=gotoailab/trendhub&type=Date)](https://star-history.com/#gotoailab/trendhub&Date)

</div>