# 排序算法优化说明

## 概述

TrendHub 已经实施了两大核心排序优化方案，让推送的内容更符合您的个人偏好：

- **方案一：关键词匹配评分** - 根据匹配的关键词数量和优先级进行评分
- **方案二：平台权重** - 为不同平台设置权重，重点关注优质平台

## 方案一：关键词匹配评分

### 核心原理

不同的关键词组具有不同的重要性，匹配更多关键词的新闻应该获得更高的分数。

### 评分规则

1. **必须词**（`+开头`）：每匹配一个 +20 分
2. **普通词**（无前缀）：每匹配一个 +10 分
3. **优先级乘数**：基础分数 × (优先级/5)
   - 优先级 1：分数 × 0.2
   - 优先级 5：分数 × 1.0（默认）
   - 优先级 10：分数 × 2.0

### 示例计算

#### 示例 1：高优先级关键词组

```
[priority:10]
AI
人工智能
+突破
```

标题："OpenAI 发布重大突破：新AI模型性能提升"

- 匹配词：AI(+10), 人工智能(+10), 突破(+20)
- 基础分：10 + 10 + 20 = 40
- 优先级乘数：10/5 = 2.0
- **最终分数：40 × 2.0 = 80 分**

#### 示例 2：低优先级关键词组

```
[priority:3]
娱乐
八卦
```

标题："娱乐圈八卦：明星动态"

- 匹配词：娱乐(+10), 八卦(+10)
- 基础分：10 + 10 = 20
- 优先级乘数：3/5 = 0.6
- **最终分数：20 × 0.6 = 12 分**

### 配置方法

在 `config/frequency_words.txt` 中使用 `[priority:X]` 标记：

```
# 最高优先级（10）- 你最想看的
[priority:10]
AI
人工智能
+突破

# 高优先级（8）- 很感兴趣
[priority:8]
华为
比亚迪
新能源

# 中等优先级（5）- 默认，可省略
[priority:5]
科技
互联网

# 低优先级（3）- 可选内容
[priority:3]
娱乐
体育
```

### 优先级设置建议

| 优先级 | 适用场景 | 效果 |
|-------|---------|------|
| 10 | 核心关注领域 | 分数翻倍，必定排前面 |
| 8-9 | 重点关注 | 分数提升60%-80% |
| 5-7 | 一般关注 | 标准分数到略微提升 |
| 3-4 | 可选关注 | 分数降低，除非很匹配 |
| 1-2 | 低优先级 | 大幅降低分数 |

## 方案二：平台权重

### 核心原理

不同平台的内容质量不同，为平台设置权重可以提升或降低该平台内容的整体排名。

### 权重范围

- **0.6-0.8**：较低权重，一般内容平台
- **0.9-1.0**：标准权重，主流平台
- **1.0+**：高权重，优质平台（可超过1.0）

### 权重效果

平台权重会影响该平台所有内容的最终分数：

```
最终分数 = 基础分数 × [1 + (平台权重 - 1) × 平台权重系数]
```

### 配置方法

在 `config/config.yaml` 的 `platforms` 部分添加 `weight` 字段：

```yaml
platforms:
  - id: zhihu
    name: 知乎
    weight: 1.2  # 高质量平台，提升20%
    
  - id: wallstreetcn-hot
    name: 华尔街见闻
    weight: 1.1  # 优质财经平台
    
  - id: weibo
    name: 微博
    weight: 0.9  # 标准平台
    
  - id: tieba
    name: 贴吧
    weight: 0.7  # 娱乐为主，降低权重
    
  - id: douyin
    name: 抖音
    weight: 0.6  # 娱乐平台
```

### 示例计算

假设一条新闻的基础分数是 60 分：

| 平台 | 权重 | 计算 | 最终分数 |
|-----|------|------|---------|
| 知乎 | 1.2 | 60 × [1+(1.2-1)×1.0] | **72 分** |
| 微博 | 0.9 | 60 × [1+(0.9-1)×1.0] | **54 分** |
| 贴吧 | 0.7 | 60 × [1+(0.7-1)×1.0] | **42 分** |

## 综合排序算法

### 总分计算公式

```
总分 = 排名分×权重 + 频次分×权重 + 关键词分×权重 + 时效性分×权重
应用平台权重后 = 总分 × [1 + (平台权重-1) × 系数]
```

### 权重配置

在 `config/config.yaml` 的 `weight` 部分：

```yaml
weight:
  rank_weight: 0.3       # 原始排名权重（降低）
  frequency_weight: 0.2  # 出现频次权重
  keyword_weight: 0.4    # 关键词匹配权重（新增，重要！）
  freshness_weight: 0.1  # 时效性权重
  platform_weight: 1.0   # 平台权重影响系数
  hotness_weight: 0.0    # 热度值权重（暂无数据）
```

### 权重调优建议

#### 场景1：只看你关心的关键词

```yaml
weight:
  rank_weight: 0.1       # 降低平台原始排名影响
  frequency_weight: 0.2  
  keyword_weight: 0.6    # 大幅提高关键词权重
  freshness_weight: 0.1
  platform_weight: 1.0
```

#### 场景2：平衡关键词和热度

```yaml
weight:
  rank_weight: 0.3       # 保持一定的热度敏感度
  frequency_weight: 0.2  
  keyword_weight: 0.4    # 适中的关键词权重
  freshness_weight: 0.1
  platform_weight: 1.0
```

#### 场景3：降低平台权重影响

如果你觉得平台权重影响太大：

```yaml
weight:
  rank_weight: 0.3
  frequency_weight: 0.2
  keyword_weight: 0.4
  freshness_weight: 0.1
  platform_weight: 0.5   # 减半平台权重影响
```

## 完整配置示例

### config.yaml

```yaml
weight:
  rank_weight: 0.3
  frequency_weight: 0.2
  keyword_weight: 0.4
  freshness_weight: 0.1
  platform_weight: 1.0
  hotness_weight: 0.0

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
  - id: douyin
    name: 抖音
    weight: 0.7
```

### frequency_words.txt

```
[priority:10]
AI
人工智能
+突破
!广告

[priority:8]
华为
鸿蒙
任正非

[priority:5]
科技
互联网

[priority:3]
娱乐
八卦
```

## 调试技巧

### 1. 查看匹配信息

修改后的 NewsItem 包含匹配信息：

- `match_score`: 关键词匹配分数
- `matched_keywords`: 匹配到的关键词列表
- `keyword_group`: 匹配的关键词组索引

### 2. 逐步调整

建议按以下步骤调整：

1. **第一步**：设置关键词组优先级（重要的设10，一般的设5，次要的设3）
2. **第二步**：设置平台权重（优质平台1.0-1.2，一般平台0.7-0.9）
3. **第三步**：调整权重配置（提高 `keyword_weight`）
4. **第四步**：观察效果，微调参数

### 3. 效果评估

运行一段时间后，观察推送内容：

- ✅ 你关心的关键词排在前面了吗？
- ✅ 优质平台的内容是否更突出？
- ✅ 不太重要的内容是否排在后面？

根据效果持续微调参数。

## 常见问题

### Q1: 设置了优先级但效果不明显？

**A**: 提高 `keyword_weight` 权重，降低 `rank_weight`：

```yaml
weight:
  rank_weight: 0.2       # 降低
  keyword_weight: 0.5    # 提高
```

### Q2: 平台权重设置后没有变化？

**A**: 检查 `platform_weight` 系数是否为 1.0：

```yaml
weight:
  platform_weight: 1.0   # 确保设置为1.0
```

### Q3: 如何让某个关键词组的内容必定排最前面？

**A**: 设置优先级为 10，并使用必须词：

```
[priority:10]
+核心关键词
其他关键词
```

### Q4: 不同关键词组匹配了同一条新闻怎么办？

**A**: 系统会选择得分最高的关键词组作为匹配结果。

## 技术细节

### 数据结构变更

#### NewsItem

```go
type NewsItem struct {
    // ... 原有字段 ...
    MatchScore      float64  // 关键词匹配分数
    MatchedKeywords []string // 匹配到的关键词列表
    KeywordGroup    int      // 匹配的关键词组索引
}
```

#### Platform

```go
type Platform struct {
    ID     string
    Name   string
    Weight float64  // 平台权重
}
```

#### KeywordGroup

```go
type KeywordGroup struct {
    Required []string
    Normal   []string
    Priority int      // 优先级 1-10
    // ...
}
```

## 版本兼容性

- 如果不设置 `priority`，默认为 5
- 如果不设置平台 `weight`，默认为 1.0
- 旧的配置文件无需修改即可运行，只是不会有优化效果

## 总结

通过这两个优化方案，你可以：

1. ✅ **精确控制**：哪些关键词最重要
2. ✅ **平台筛选**：提升优质平台内容
3. ✅ **个性化**：完全根据你的偏好排序
4. ✅ **灵活调整**：随时修改配置立即生效（支持热重载）

开始配置你的个性化排序吧！🚀

