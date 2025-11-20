# Bark 推送集成更新日志

## 版本：v2.1.0
**发布日期**: 2025-11-20

## 🎉 新功能

### Bark 推送支持

新增对 [Bark](https://github.com/finb/bark) iOS 推送通知的完整支持！

#### 核心功能
- ✅ 支持官方 Bark 服务器 (`https://api.day.app`)
- ✅ 支持自建 Bark 服务器
- ✅ 自动格式化推送内容
- ✅ 支持点击通知跳转链接
- ✅ 自动分组管理
- ✅ 自定义提示音

#### 集成方式
1. **配置文件支持** (`config.yaml`)
2. **Web 界面配置** (可视化配置)
3. **环境变量覆盖** (可选)

## 📝 新增文件

```
trendhub/
├── internal/notifier/
│   └── bark.go                    # Bark 推送实现
└── docs/
    └── BARK_SETUP.md              # Bark 完整配置指南
```

## 🔧 修改文件

### 1. `config/config.go`
**变更**: 在 `WebhooksConfig` 结构体中添加 Bark 配置字段

```go
type WebhooksConfig struct {
    // ... 其他字段
    BarkServerURL string `yaml:"bark_server_url" json:"bark_server_url"`
    BarkDeviceKey string `yaml:"bark_device_key" json:"bark_device_key"`
}
```

### 2. `internal/notifier/manager.go`
**变更**: 在 `NewNotificationManager` 中注册 Bark notifier

```go
if cfg.Notification.Webhooks.BarkDeviceKey != "" {
    manager.notifiers = append(manager.notifiers, 
        NewBarkNotifier(cfg.Notification.Webhooks.BarkServerURL, 
                       cfg.Notification.Webhooks.BarkDeviceKey))
}
```

### 3. `config.example.yaml`
**变更**: 添加 Bark 配置示例

```yaml
webhooks:
    bark_server_url: "https://api.day.app"
    bark_device_key: ""
```

### 4. `web/static/index.html`
**变更**: Web 界面添加 Bark 配置项

在"即时通讯工具 Webhooks"部分添加：
- Bark Server URL 输入框
- Bark Device Key 输入框
- 配置说明提示

### 5. `README.md`
**变更**: 更新功能列表和推送渠道说明

## 📖 使用方法

### 快速开始（3 步）

#### 1. 安装 Bark App
在 App Store 搜索并安装 Bark

#### 2. 获取设备密钥
打开 Bark App，复制显示的设备密钥

#### 3. 配置 TrendHub

**方式一：Web 界面**
```bash
./trendhub -web
# 打开 http://localhost:8080
# 进入"系统配置" -> "通知推送" -> 填写 Bark 配置
```

**方式二：配置文件**
```yaml
notification:
  enable_notification: true
  webhooks:
    bark_device_key: "your_device_key_here"
```

### 推送内容

推送消息包含：
- 📱 标题: "TrendHub 热点监控"
- 📝 内容: 各平台热搜新闻（最多 10 条）
- 🔗 跳转: 点击可跳转到第一条新闻
- 📁 分组: 自动归类到 "TrendHub"
- 🔔 提示音: calypso

## 🎨 Web 界面更新

### 新增配置项
在"系统配置" -> "通知推送" -> "即时通讯工具 Webhooks"部分：

```
┌─────────────────────────────────────────┐
│ Bark Server URL                         │
│ ┌─────────────────────────────────────┐ │
│ │ https://api.day.app                 │ │
│ └─────────────────────────────────────┘ │
│ 💡 留空使用官方服务器，或填写自建地址    │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│ Bark Device Key                         │
│ ┌─────────────────────────────────────┐ │
│ │ your_device_key                     │ │
│ └─────────────────────────────────────┘ │
│ 💡 在 Bark App 中查看                   │
└─────────────────────────────────────────┘
```

## 🔍 技术细节

### API 实现
```go
// Bark API 格式
GET {serverURL}/{deviceKey}/{title}/{body}?url={url}&group={group}&sound={sound}

// TrendHub 使用的参数
- title: "TrendHub 热点监控"
- body: 热搜新闻内容
- url: 第一条新闻链接
- group: "TrendHub"
- sound: "calypso"
```

### 错误处理
- ✅ 设备密钥验证
- ✅ 网络请求超时（10秒）
- ✅ HTTP 状态码检查
- ✅ 详细错误日志

### 内容优化
- 最多显示 10 条新闻，避免推送内容过长
- 按平台分组显示
- 超过限制显示"还有 X 条"
- URL 编码处理，确保特殊字符正常传输

## 📊 兼容性

- ✅ iOS 14+
- ✅ iPadOS 14+
- ✅ 官方 Bark 服务器
- ✅ 自建 Bark 服务器
- ✅ 与其他推送渠道同时使用

## 🔒 安全性

- 推送内容使用 HTTPS 加密传输
- 设备密钥仅保存在本地配置文件
- 官方服务器不保存推送内容
- 支持自建服务器，完全掌控数据

## 📚 相关文档

- [Bark 完整配置指南](docs/BARK_SETUP.md)
- [定时推送功能文档](docs/PUSH_SCHEDULE.md)
- [快速开始指南](QUICKSTART_PUSH.md)
- [Bark 官方文档](https://github.com/finb/bark)

## 🐛 已知问题

无

## 🔮 未来计划

- [ ] 支持自定义提示音选择
- [ ] 支持图片推送
- [ ] 支持多设备配置
- [ ] 推送失败自动重试

## 💡 使用建议

1. **首次使用**: 先手动运行测试推送是否正常
2. **定时推送**: 配合时间窗口功能，避免打扰
3. **分组管理**: 在 Bark App 中为 TrendHub 设置单独分组
4. **免打扰**: 设置合理的免打扰时段
5. **备份密钥**: 将设备密钥保存在安全的地方

## 🙏 致谢

- [Bark](https://github.com/finb/bark) - 优秀的 iOS 推送工具
- 所有贡献者和用户的支持

## 📞 技术支持

如有问题，欢迎：
- 查看 [Bark 配置指南](docs/BARK_SETUP.md)
- 提交 GitHub Issues
- 查阅 Bark 官方文档

---

**更新内容**: 添加 Bark 推送支持  
**影响范围**: 新增功能，不影响现有功能  
**升级建议**: 直接替换二进制文件即可使用新功能  
**配置迁移**: 无需迁移，仅需添加 Bark 配置即可

