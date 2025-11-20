# Bark 推送配置指南

## 什么是 Bark？

Bark 是一个开源的 iOS 推送通知工具，可以让你通过简单的 HTTP API 向你的 iPhone/iPad 发送通知。

- GitHub：https://github.com/finb/bark
- 官方服务器：https://api.day.app

## 功能特点

✅ 简单易用 - 只需一个设备密钥即可使用  
✅ 免费服务 - 官方提供免费的推送服务  
✅ 自建支持 - 可以部署自己的 Bark 服务器  
✅ 丰富选项 - 支持自定义声音、分组、跳转链接等  
✅ 隐私保护 - 推送内容不会在服务器保存  

## 快速开始

### 1. 安装 Bark App

在 App Store 搜索并安装 **Bark** 应用。

### 2. 获取设备密钥

打开 Bark App，你会看到类似这样的内容：

```
https://api.day.app/your_device_key_here
```

其中 `your_device_key_here` 就是你的**设备密钥**。

### 3. 配置 TrendHub

#### 方法一：通过 Web 界面配置（推荐）

1. 启动 TrendHub Web 模式：
   ```bash
   ./trendhub -web
   ```

2. 打开浏览器访问：`http://localhost:8080`

3. 进入"系统配置"标签页

4. 找到"通知推送"部分，启用通知

5. 在"即时通讯工具 Webhooks"部分找到 Bark 配置：
   - **Bark Server URL**: `https://api.day.app` (使用官方服务器，可留空)
   - **Bark Device Key**: 填入你的设备密钥（如 `your_device_key_here`）

6. 点击"保存配置"

#### 方法二：编辑配置文件

编辑 `config.yaml`：

```yaml
notification:
  enable_notification: true
  
  webhooks:
    bark_server_url: "https://api.day.app"  # 可选，留空使用官方服务器
    bark_device_key: "your_device_key_here"  # 替换为你的设备密钥
```

### 4. 测试推送

在 Web 界面的"仪表盘"页面点击"立即运行"按钮，或者运行命令：

```bash
./trendhub
```

如果配置正确，你的 iPhone 将收到来自 TrendHub 的推送通知！📱

## 高级配置

### 使用自建 Bark 服务器

如果你部署了自己的 Bark 服务器（比如 `https://bark.example.com`），可以这样配置：

```yaml
webhooks:
  bark_server_url: "https://bark.example.com"
  bark_device_key: "your_device_key"
```

### Bark 推送特性

TrendHub 的 Bark 推送包含以下特性：

- 📱 **标题**: "TrendHub 热点监控"
- 📝 **内容**: 包含各平台热搜新闻，最多显示 10 条
- 🔗 **跳转链接**: 点击通知可跳转到第一条新闻
- 📁 **分组**: 自动归类到"TrendHub"分组
- 🔔 **提示音**: 使用 "calypso" 提示音

## 推送内容示例

```
TrendHub 热点监控

热点监控 (14:30)

【微博】
1. 某某事件登上热搜
2. 某某新闻引发关注

【知乎】
3. 某某话题热度飙升
4. 某某问题获得高赞

...还有 6 条
```

## 常见问题

### Q: 为什么没有收到推送？

**A**: 检查以下几点：
1. ✅ Bark App 是否已安装并允许通知
2. ✅ 设备密钥是否正确填写
3. ✅ iPhone 是否联网
4. ✅ 通知设置中是否启用了 TrendHub 分组
5. ✅ 查看 TrendHub 日志是否有错误信息

### Q: 推送内容太长会被截断吗？

**A**: TrendHub 会自动限制推送内容：
- 最多显示 10 条热搜
- 超过 10 条会显示"还有 X 条"
- Bark 本身支持较长内容，一般不会截断

### Q: 可以自定义推送的声音吗？

**A**: 当前版本使用固定的 "calypso" 提示音。如需自定义，可以修改 `internal/notifier/bark.go` 中的 `sound` 参数。

Bark 支持的声音列表：
- `alarm` - 警告声
- `anticipate` - 期待
- `bell` - 铃声
- `birdsong` - 鸟鸣
- `bloom` - 绽放
- `calypso` - 卡里普索（默认）
- `chime` - 钟鸣
- `choo` - 火车
- `descent` - 下降
- `electronic` - 电子音
- `fanfare` - 号角
- `glass` - 玻璃
- `gotosleep` - 入睡
- `healthnotification` - 健康通知
- `horn` - 喇叭
- `ladder` - 梯子
- `mailsent` - 邮件发送
- `minuet` - 小步舞曲
- `multiwayinvitation` - 多方邀请
- `newmail` - 新邮件
- `newsflash` - 新闻快报
- `noir` - 黑色电影
- `paymentsuccess` - 支付成功
- `shake` - 摇动
- `sherwoodforest` - 舍伍德森林
- `silence` - 静音
- `spell` - 咒语
- `suspense` - 悬疑
- `telegraph` - 电报
- `tiptoes` - 踮脚
- `typewriters` - 打字机
- `update` - 更新

### Q: 推送频率是怎样的？

**A**: 推送频率取决于你的配置：
- 如果启用了定时推送窗口，则在窗口期内每天推送一次
- 如果手动运行任务，则立即推送
- 不会重复推送相同的内容

### Q: Bark 服务是否安全？

**A**: 
- ✅ 官方 Bark 服务器采用 HTTPS 加密传输
- ✅ 推送内容不会在服务器端保存
- ✅ 设备密钥类似于密码，请妥善保管
- ✅ 如果担心隐私，可以自建 Bark 服务器

### Q: 如何自建 Bark 服务器？

**A**: 参考 Bark 官方文档：
```bash
# 使用 Docker 部署
docker run -d \
  --name bark-server \
  -p 8080:8080 \
  finab/bark-server:latest
```

然后在配置中使用自己的服务器地址：
```yaml
bark_server_url: "http://your-server-ip:8080"
```

### Q: 可以同时配置多个推送渠道吗？

**A**: 可以！TrendHub 支持同时配置多种推送方式：
- 飞书
- 钉钉
- 企业微信
- Telegram
- Bark
- Ntfy
- 邮件

所有配置的推送渠道都会同时接收通知。

## Bark App 设置建议

### 推荐设置

1. **通知分组**: 在 Bark App 中为 TrendHub 创建单独分组
2. **免打扰时段**: 设置晚上 23:00 - 早上 7:00 免打扰
3. **提示音**: 可以为 TrendHub 分组设置独特的提示音
4. **通知样式**: 建议使用"横幅"样式，方便快速查看

### 批量管理

如果你有多个设备（iPhone + iPad），可以：
1. 在每个设备上安装 Bark
2. 记录每个设备的密钥
3. 在不同的 TrendHub 实例中配置不同的设备密钥

## API 参考

Bark 推送 API 格式：

```
GET {serverURL}/{deviceKey}/{title}/{body}?url={url}&group={group}&sound={sound}
```

TrendHub 使用的参数：
- `title`: "TrendHub 热点监控"
- `body`: 热搜新闻内容
- `url`: 第一条新闻的链接（点击通知跳转）
- `group`: "TrendHub"
- `sound`: "calypso"

更多 Bark API 参数，请参考：https://github.com/finb/bark

## 故障排查

### 推送失败

如果日志显示推送失败，可能的原因：

1. **网络问题**
   ```bash
   # 测试网络连接
   curl https://api.day.app
   ```

2. **设备密钥错误**
   - 重新检查 Bark App 中的密钥
   - 确保复制时没有多余的空格

3. **服务器不可达**
   - 检查自建服务器是否在线
   - 确认防火墙是否开放端口

### 测试推送

使用 curl 测试 Bark 推送：

```bash
curl "https://api.day.app/your_device_key/测试标题/测试内容"
```

如果能收到推送，说明 Bark 配置正确，问题可能在 TrendHub 配置上。

## 最佳实践

1. **定期检查**: 定期查看推送记录，确保推送正常
2. **合理分组**: 在 Bark App 中为不同的推送源设置分组
3. **备份密钥**: 将设备密钥保存在安全的地方
4. **测试环境**: 在生产使用前先在测试环境验证
5. **关注日志**: 出现问题时检查 TrendHub 日志获取详细信息

## 相关资源

- [Bark GitHub](https://github.com/finb/bark)
- [Bark 文档](https://bark.day.app/)
- [TrendHub 定时推送文档](PUSH_SCHEDULE.md)
- [TrendHub 快速开始](../QUICKSTART_PUSH.md)

## 技术支持

如有问题，欢迎：
- 查看项目 GitHub Issues
- 提交问题反馈
- 参考 Bark 官方文档

祝推送愉快！📱✨

