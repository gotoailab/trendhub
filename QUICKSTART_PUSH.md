# 定时推送快速开始指南

## 快速配置（3步开始使用）

### 1. 配置推送时间窗口

编辑 `config.yaml`：

```yaml
notification:
  enable_notification: true
  
  push_window:
    enabled: true              # 启用定时推送
    time_range:
      start: "09:00"           # 每天9点开始
      end: "18:00"             # 18点结束
    once_per_day: true         # 每天只推送一次
    push_record_retention_days: 30  # 保留30天记录
```

### 2. 启动服务

```bash
# 编译
go build -o trendhub ./cmd/main.go

# 启动 Web 模式（会自动启用定时推送）
./trendhub -web
```

### 3. 查看推送记录

打开浏览器访问：http://localhost:8080

点击"推送记录"标签页即可查看所有历史推送记录。

## 常见使用场景

### 场景1: 工作日定时推送

```yaml
push_window:
  enabled: true
  time_range:
    start: "09:00"
    end: "18:00"
  once_per_day: true
```

**效果**: 每个工作日在 9:00-18:00 之间第一次检查时推送一次

### 场景2: 晚间汇总推送

```yaml
push_window:
  enabled: true
  time_range:
    start: "20:00"
    end: "22:00"
  once_per_day: true
```

**效果**: 每天晚上 20:00-22:00 之间推送当日汇总

### 场景3: 多次推送（不推荐）

```yaml
push_window:
  enabled: true
  time_range:
    start: "08:00"
    end: "20:00"
  once_per_day: false  # 关闭每日一次限制
```

**效果**: 在 8:00-20:00 期间，每分钟都会检查并可能推送（慎用）

### 场景4: 跨夜推送

```yaml
push_window:
  enabled: true
  time_range:
    start: "22:00"
    end: "02:00"     # 次日凌晨2点
  once_per_day: true
```

**效果**: 每天晚上 22:00 到次日凌晨 02:00 之间推送

## Web 界面功能

### 推送记录页面
- ✅ 查看所有历史推送记录
- ✅ 显示推送时间、状态、条目数
- ✅ 查看执行耗时和错误信息
- ✅ 分页浏览历史记录
- ✅ 刷新按钮实时更新

### 记录信息
每条记录包含：
- **时间**: 精确到秒的推送时间
- **状态**: 成功（绿色）/ 失败（红色）/ 部分成功（黄色）
- **推送条目**: 本次推送的新闻条目数量
- **耗时**: 任务执行时间
- **错误信息**: 如果失败，显示错误原因

## 常见问题

### Q: 为什么没有自动推送？
**A**: 检查以下几点：
1. `push_window.enabled` 是否为 `true`
2. 当前时间是否在配置的时间窗口内
3. 如果启用了 `once_per_day`，今天是否已经推送过
4. 查看日志确认调度器是否正常启动

### Q: 如何手动触发推送？
**A**: 在 Web 界面的"仪表盘"页面点击"立即运行"按钮

### Q: 推送记录保存在哪里？
**A**: 默认保存在 `data/push_records.db` 文件中

### Q: 可以导出推送记录吗？
**A**: 目前支持通过 API 接口获取：
```bash
curl http://localhost:8080/api/push-records?limit=100
```

### Q: 如何清理旧记录？
**A**: 设置 `push_record_retention_days`，系统会自动清理：
```yaml
push_window:
  push_record_retention_days: 7  # 只保留7天
```

## 测试建议

### 1. 测试时间窗口

设置一个即将到来的时间窗口：
```yaml
time_range:
  start: "14:30"  # 设置为当前时间后几分钟
  end: "14:35"
```

### 2. 观察日志

```bash
# 查看完整日志
./trendhub -web 2>&1 | tee trendhub.log

# 日志会显示：
# - Scheduler started with time window: 14:30 - 14:35
# - Time window matched, executing task...
# - Task completed successfully, pushed 10 items
```

### 3. 验证记录

推送完成后：
1. 打开 Web 界面
2. 切换到"推送记录"标签
3. 查看最新的推送记录

## 命令行参数

```bash
./trendhub -web                           # 使用默认配置
./trendhub -web -addr :9090              # 指定端口
./trendhub -web -db data/my_records.db   # 指定数据库路径
./trendhub -web -config my_config.yaml   # 指定配置文件
```

## 生产环境部署

### 使用 systemd

创建 `/etc/systemd/system/trendhub.service`：

```ini
[Unit]
Description=TrendHub Push Service
After=network.target

[Service]
Type=simple
User=trendhub
WorkingDirectory=/opt/trendhub
ExecStart=/opt/trendhub/trendhub -web -addr :8080
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable trendhub
sudo systemctl start trendhub
sudo systemctl status trendhub
```

### 使用 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o trendhub ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/trendhub .
COPY config.yaml frequency_words.txt ./
RUN mkdir -p data
EXPOSE 8080
CMD ["./trendhub", "-web", "-addr", ":8080"]
```

构建和运行：
```bash
docker build -t trendhub .
docker run -d -p 8080:8080 -v $(pwd)/data:/root/data --name trendhub trendhub
```

## 监控和维护

### 检查服务状态
```bash
curl http://localhost:8080/api/status
```

### 查看最近的推送
```bash
curl http://localhost:8080/api/push-records?limit=5
```

### 备份数据库
```bash
# 定期备份推送记录
cp data/push_records.db data/backup/push_records_$(date +%Y%m%d).db
```

## 技术支持

如有问题，请查看：
- 完整文档：`docs/PUSH_SCHEDULE.md`
- 配置示例：`config.example.yaml`
- 项目 Issues：提交问题反馈

祝使用愉快！🎉

