package rank

import (
	"sort"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/model"
)

type Ranker interface {
	Rank(items map[string][]*model.NewsItem) []*model.NewsItem
}

type WeightedRanker struct {
	cfg       config.WeightConfig
	platforms map[string]float64 // 平台ID到权重的映射
}

func NewWeightedRanker(cfg config.WeightConfig, platforms []model.Platform) *WeightedRanker {
	// 构建平台权重映射
	platformWeights := make(map[string]float64)
	for _, p := range platforms {
		weight := p.Weight
		if weight <= 0 {
			weight = 1.0 // 默认权重
		}
		platformWeights[p.ID] = weight
	}
	return &WeightedRanker{
		cfg:       cfg,
		platforms: platformWeights,
	}
}

func (r *WeightedRanker) Rank(data map[string][]*model.NewsItem) []*model.NewsItem {
	var allItems []*model.NewsItem

	for _, items := range data {
		allItems = append(allItems, items...)
	}

	sort.Slice(allItems, func(i, j int) bool {
		scoreI := r.calculateScore(allItems[i])
		scoreJ := r.calculateScore(allItems[j])

		// 分数高的在前
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		// 分数相同，按ID排序保证稳定性
		return allItems[i].SourceID < allItems[j].SourceID
	})

	return allItems
}

func (r *WeightedRanker) calculateScore(item *model.NewsItem) float64 {
	// 1. 排名分 (0-100)
	rank := 1
	if len(item.Ranks) > 0 {
		rank = item.Ranks[0]
	}
	if rank < 1 {
		rank = 1
	}
	rankScore := 1.0 / float64(rank) * 100.0 // 归一化到0-100区间

	// 2. 频次分 (0-100)
	freqScore := float64(item.AppearCount)
	if freqScore == 0 {
		freqScore = 1
	}
	// 假设频次一般不超过24次(一天每小时一次)，简单归一化
	freqScore = freqScore / 24.0 * 100.0
	if freqScore > 100 {
		freqScore = 100
	}

	// 3. 热度分 (0-100) - 暂无数据，设为0
	hotnessScore := 0.0

	// 4. 关键词匹配分 (0-100) - 方案一的核心
	keywordScore := item.MatchScore
	// 归一化：假设最高分是100（2个必须词+1个普通词+优先级10 = (40+10)*2=100）
	if keywordScore > 100 {
		keywordScore = 100
	}

	// 5. 时效性分 (0-100) - 新出现的内容加分
	freshnessScore := 0.0
	if item.IsNew {
		freshnessScore = 100.0 // 新内容满分
	}

	// 获取关键词权重，如果配置为0则使用默认值
	keywordWeight := r.cfg.KeywordWeight
	if keywordWeight == 0 {
		keywordWeight = 0.3 // 默认关键词权重
	}

	freshnessWeight := r.cfg.FreshnessWeight
	if freshnessWeight == 0 {
		freshnessWeight = 0.1 // 默认时效性权重
	}

	// 加权总分
	totalScore := rankScore*r.cfg.RankWeight +
		freqScore*r.cfg.FrequencyWeight +
		hotnessScore*r.cfg.HotnessWeight +
		keywordScore*keywordWeight +
		freshnessScore*freshnessWeight

	// 6. 应用平台权重 (0-1) - 方案二的核心
	platformWeight, exists := r.platforms[item.SourceID]
	if !exists {
		platformWeight = 1.0 // 默认权重
	}

	// 平台权重影响系数
	platformWeightEffect := r.cfg.PlatformWeight
	if platformWeightEffect == 0 {
		platformWeightEffect = 1.0 // 默认完全应用平台权重
	}

	// 应用平台权重：基础分数 * (1 + (platformWeight - 1) * effect)
	// 这样当 platformWeight=1.0 时，分数不变
	// 当 platformWeight=0.8 时，分数会减少
	// 当 platformWeight=1.2 时，分数会增加
	totalScore = totalScore * (1.0 + (platformWeight-1.0)*platformWeightEffect)

	return totalScore
}
