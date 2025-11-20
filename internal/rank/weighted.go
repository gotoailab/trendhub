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
	cfg config.WeightConfig
}

func NewWeightedRanker(cfg config.WeightConfig) *WeightedRanker {
	return &WeightedRanker{cfg: cfg}
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
	// 基础分：倒数排名 (排名越小分数越高，例如第1名得1分，第10名得0.1分)
	// 防止除以0，rank至少为1
	rank := 1
	if len(item.Ranks) > 0 {
		rank = item.Ranks[0]
	}
	if rank < 1 {
		rank = 1
	}
	
	rankScore := 1.0 / float64(rank) * 100.0 // 归一化到0-100区间

	// 频次分：目前没有历史记录，默认就是1次
	freqScore := float64(item.AppearCount)
	if freqScore == 0 {
		freqScore = 1
	}
	// 假设频次一般不超过24次(一天每小时一次)，简单归一化
	freqScore = freqScore / 24.0 * 100.0

	// 热度分：暂无数据，设为0
	hotnessScore := 0.0

	// 加权总分
	totalScore := rankScore*r.cfg.RankWeight + freqScore*r.cfg.FrequencyWeight + hotnessScore*r.cfg.HotnessWeight
	
	return totalScore
}
