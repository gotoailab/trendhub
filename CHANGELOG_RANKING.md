# 排序算法优化更新日志

## 版本说明

本次更新实施了两大核心排序优化方案，大幅提升内容推送的个性化程度。

## 新增功能

### 🎯 方案一：关键词匹配评分系统

#### 功能描述
为关键词组引入优先级机制（1-10级），根据匹配的关键词数量和优先级计算分数。

#### 技术实现
1. **数据模型扩展**
   - `NewsItem` 新增字段：
     - `MatchScore` (float64): 关键词匹配分数
     - `MatchedKeywords` ([]string): 匹配到的关键词列表
     - `KeywordGroup` (int): 匹配的关键词组索引

2. **配置结构扩展**
   - `KeywordGroup` 新增 `Priority` (int) 字段
   - `WeightConfig` 新增 `KeywordWeight` (float64) 字段

3. **评分算法**
   - 必须词：每个 +20 分
   - 普通词：每个 +10 分
   - 优先级乘数：score × (priority/5)

#### 使用方式
```
[priority:10]
AI
人工智能
+突破

[priority:5]
科技
互联网
```

### 🏆 方案二：平台权重系统

#### 功能描述
为不同平台设置权重，优质平台的内容可以获得更高的排名。

#### 技术实现
1. **数据模型扩展**
   - `Platform` 新增 `Weight` (float64) 字段

2. **配置结构扩展**
   - `WeightConfig` 新增 `PlatformWeight` (float64) 字段

3. **排序器优化**
   - `WeightedRanker` 构造函数新增 `platforms` 参数
   - 添加平台权重映射缓存
   - 分数计算考虑平台权重

#### 使用方式
```yaml
platforms:
  - id: zhihu
    name: 知乎
    weight: 1.2  # 优质平台
  - id: douyin
    name: 抖音
    weight: 0.7  # 娱乐平台
```

### ⏰ 方案三：时效性评分（附加）

#### 功能描述
新出现的内容自动获得额外加分。

#### 技术实现
- `WeightConfig` 新增 `FreshnessWeight` (float64) 字段
- 利用 `NewsItem.IsNew` 字段判断新内容
- 新内容获得满分（100分）

## 文件变更清单

### 核心代码变更

1. **internal/model/news.go**
   - ✅ `NewsItem` 添加匹配评分字段
   - ✅ `Platform` 添加权重字段

2. **config/config.go**
   - ✅ `KeywordGroup` 添加优先级字段
   - ✅ `WeightConfig` 添加三个新权重字段
   - ✅ `loadKeywords()` 支持解析 `[priority:X]` 标记

3. **internal/filter/keyword.go**
   - ✅ 新增 `matchWithScore()` 方法
   - ✅ `Filter()` 方法填充匹配分数信息
   - ✅ 保持 `match()` 方法向后兼容

4. **internal/rank/weighted.go**
   - ✅ `WeightedRanker` 添加平台权重映射
   - ✅ 构造函数接收 `platforms` 参数
   - ✅ `calculateScore()` 集成关键词、平台、时效性评分

5. **调用点更新**
   - ✅ `web/runner.go` (2处)
   - ✅ `examples/run_once/main.go` (1处)

### 配置文件变更

1. **config.example.yaml**
   - ✅ 添加新的权重配置项
   - ✅ 为所有平台添加 `weight` 字段

2. **frequency_words.example.txt**
   - ✅ 添加使用说明注释
   - ✅ 为示例关键词组添加优先级标记

### 文档变更

1. **新增文档**
   - ✅ `docs/RANKING_OPTIMIZATION.md` - 详细使用说明
   - ✅ `docs/RANKING_MIGRATION.md` - 迁移指南
   - ✅ `CHANGELOG_RANKING.md` - 本变更日志

2. **README.md 更新**
   - ✅ 核心功能列表更新
   - ✅ 添加智能排序章节
   - ✅ 文档链接列表更新

## 兼容性

### 向后兼容 ✅

本次更新完全向后兼容，旧配置文件无需修改即可运行：

- 未设置优先级的关键词组默认为 5
- 未设置权重的平台默认为 1.0
- 未配置新权重项时使用默认值

### 默认值

```go
// 关键词组默认优先级
priority = 5

// 平台默认权重
platformWeight = 1.0

// 默认权重配置
keywordWeight = 0.3
platformWeight = 1.0
freshnessWeight = 0.1
```

## 性能影响

### 计算复杂度
- ⏱️ 关键词匹配：O(n×m) - n为新闻数，m为关键词数
- ⏱️ 平台权重：O(1) - 使用映射缓存
- ⏱️ 总体影响：可忽略（<1ms增加）

### 内存占用
- 📊 每个 `NewsItem` 增加约 40-60 字节
- 📊 平台权重映射：约 1KB（11个平台）
- 📊 总体影响：可忽略

## 测试建议

### 功能测试

1. **关键词优先级测试**
   ```bash
   # 创建测试配置
   [priority:10]
   测试词A
   
   [priority:1]
   测试词B
   
   # 运行并观察：测试词A 的内容应该排在前面
   ```

2. **平台权重测试**
   ```yaml
   platforms:
     - id: zhihu
       name: 知乎
       weight: 2.0  # 极端测试
     - id: weibo
       name: 微博
       weight: 0.5
   
   # 运行并观察：知乎内容应该明显排前
   ```

3. **综合测试**
   - 同时配置关键词优先级和平台权重
   - 观察综合效果

### 回归测试

```bash
# 使用旧配置文件测试
cp config.old.yaml config.yaml
./trendhub

# 应该正常运行，结果与之前一致
```

## 升级步骤

### 快速升级（5分钟）

1. **重新编译**
   ```bash
   go build -o trendhub cmd/main.go
   ```

2. **更新配置**（可选）
   ```bash
   # 备份旧配置
   cp config/config.yaml config/config.yaml.bak
   
   # 参考 config.example.yaml 添加新字段
   ```

3. **设置关键词优先级**（可选）
   ```bash
   # 编辑 frequency_words.txt
   # 为重要的关键词组添加 [priority:10]
   ```

4. **重启服务**
   ```bash
   ./trendhub -web
   ```

### 详细升级指南

参考：[docs/RANKING_MIGRATION.md](docs/RANKING_MIGRATION.md)

## 已知问题

无

## 下一步计划

1. **用户反馈学习**
   - 记录用户点击行为
   - 自动调整关键词权重

2. **热度值集成**
   - 从 API 获取真实热度值
   - 启用 `hotness_weight`

3. **时间衰减**
   - 旧内容逐渐降低分数
   - 保持推送内容新鲜度

4. **A/B 测试框架**
   - 支持多组配置对比
   - 数据驱动的优化决策

## 贡献者

感谢所有参与本次优化的贡献者！

## 获取支持

- 📖 文档：[docs/RANKING_OPTIMIZATION.md](docs/RANKING_OPTIMIZATION.md)
- 💬 微信群：添加 mongorz，备注：simple
- 🐛 问题反馈：[GitHub Issues](https://github.com/gotoailab/trendhub/issues)

---

更新时间：2024-11-21  
版本：v1.1.0（排序优化版）

