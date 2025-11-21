# 排序优化迁移指南

## 快速开始

如果你已经在使用 TrendHub，按照以下步骤升级到新的排序算法。

## 第一步：更新配置文件

### 1.1 更新 config.yaml

在 `weight` 部分添加新的权重配置：

```yaml
weight:
  rank_weight: 0.3       # 降低原始排名权重
  frequency_weight: 0.2  
  hotness_weight: 0.0
  keyword_weight: 0.4    # 新增：关键词匹配权重
  platform_weight: 1.0   # 新增：平台权重影响系数
  freshness_weight: 0.1  # 新增：时效性权重
```

### 1.2 为平台添加权重

在 `platforms` 部分为每个平台添加 `weight` 字段：

```yaml
platforms:
  - id: zhihu
    name: 知乎
    weight: 1.0  # 添加这一行
    
  - id: weibo
    name: 微博
    weight: 0.9  # 添加这一行
    
  # ... 其他平台
```

**建议权重值**：
- 优质平台（知乎、华尔街见闻等）：1.0-1.2
- 主流平台（微博、今日头条等）：0.8-1.0
- 娱乐平台（抖音、贴吧等）：0.6-0.8

## 第二步：设置关键词优先级

### 2.1 打开 frequency_words.txt

在你现有的关键词组前添加 `[priority:X]` 标记。

### 2.2 示例转换

**旧格式**（仍然可用）：
```
AI
人工智能
+突破

华为
鸿蒙
```

**新格式**（推荐）：
```
[priority:10]
AI
人工智能
+突破

[priority:8]
华为
鸿蒙
```

### 2.3 优先级分配原则

根据你的兴趣程度分配优先级：

| 优先级 | 说明 | 示例 |
|-------|------|------|
| 10 | 最关心的核心话题 | AI、核心技术突破 |
| 8-9 | 非常感兴趣 | 关注的公司、重要行业 |
| 5-7 | 一般感兴趣 | 普通新闻、行业动态 |
| 3-4 | 可选内容 | 次要话题 |
| 1-2 | 低优先级 | 娱乐八卦等 |

## 第三步：重启或热重载

### 方式1：热重载（推荐）

如果你在使用 Web 模式：

1. 访问 Web 界面
2. 进入"功能设置"页面
3. 修改配置后点击"保存配置"
4. 系统自动应用新配置 ✓

### 方式2：重启

```bash
# 停止当前运行的 trendhub
pkill trendhub

# 重新启动
./trendhub -web
```

## 第四步：验证效果

### 4.1 观察推送结果

运行一段时间后，检查推送的内容：

- ✅ 高优先级关键词的新闻排在前面了吗？
- ✅ 优质平台的内容是否更突出？
- ✅ 整体排序是否更符合你的期望？

### 4.2 微调参数

根据效果调整配置：

**如果高优先级关键词效果不明显**：
```yaml
weight:
  rank_weight: 0.2       # 进一步降低
  keyword_weight: 0.5    # 进一步提高
```

**如果平台权重效果太强**：
```yaml
weight:
  platform_weight: 0.5   # 减半影响
```

**如果新出现的内容不够突出**：
```yaml
weight:
  freshness_weight: 0.2  # 提高时效性权重
```

## 兼容性说明

### 向后兼容

- ✅ 不添加 `[priority:]` 标记，默认优先级为 5
- ✅ 不添加平台 `weight`，默认权重为 1.0
- ✅ 不添加新的 weight 配置项，使用默认值
- ✅ 旧配置文件可以直接运行

### 建议升级

虽然旧配置仍可用，但**强烈建议**完成迁移以获得最佳效果。

## 完整示例

### config.yaml 完整示例

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

weight:
  rank_weight: 0.3
  frequency_weight: 0.2
  hotness_weight: 0.0
  keyword_weight: 0.4
  platform_weight: 1.0
  freshness_weight: 0.1

platforms:
  - id: zhihu
    name: 知乎
    weight: 1.2
  - id: wallstreetcn-hot
    name: 华尔街见闻
    weight: 1.1
  - id: weibo
    name: 微博
    weight: 0.9
  - id: baidu
    name: 百度热搜
    weight: 0.8
  - id: douyin
    name: 抖音
    weight: 0.7
  - id: tieba
    name: 贴吧
    weight: 0.6
```

### frequency_words.txt 完整示例

```
# 核心关注（优先级10）
[priority:10]
AI
人工智能
+突破
!广告

[priority:10]
DeepSeek
ChatGPT
Claude

# 重点关注（优先级8-9）
[priority:9]
华为
鸿蒙
任正非

[priority:8]
比亚迪
王传福

[priority:8]
芯片
光刻机
+国产

# 一般关注（优先级5-7）
[priority:7]
科技
互联网

[priority:6]
新能源
电动车

[priority:5]
手机
iPhone

# 可选关注（优先级3-4）
[priority:4]
游戏
电竞

[priority:3]
娱乐
明星
```

## 常见问题

### Q: 需要重新编译吗？

A: 是的，需要重新编译：
```bash
go build -o trendhub cmd/main.go
```

### Q: 现有的推送记录会受影响吗？

A: 不会。排序优化只影响新的爬取和排序，历史记录保持不变。

### Q: 可以只使用一个方案吗？

A: 可以。
- 只用关键词优先级：设置 `platform_weight: 0`
- 只用平台权重：不设置 `[priority:]`，保持关键词默认优先级

### Q: 如何恢复到旧的排序算法？

A: 设置以下权重即可：
```yaml
weight:
  rank_weight: 0.6
  frequency_weight: 0.3
  hotness_weight: 0.1
  keyword_weight: 0.0
  platform_weight: 0.0
  freshness_weight: 0.0
```

## 获取帮助

- 📖 详细说明：查看 [RANKING_OPTIMIZATION.md](./RANKING_OPTIMIZATION.md)
- 💬 交流群：添加微信 mongorz，备注：simple
- 🐛 问题反馈：[GitHub Issues](https://github.com/gotoailab/trendhub/issues)

## 下一步

完成迁移后，建议阅读：

1. [排序算法优化详细说明](./RANKING_OPTIMIZATION.md)
2. [配置热重载功能](./CONFIG_HOT_RELOAD.md)
3. [报告模式详解](./REPORT_MODES.md)

祝你使用愉快！🎉

