package filter

import (
	"strings"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/model"
)

type Filter interface {
	Filter(items map[string][]*model.NewsItem) (map[string][]*model.NewsItem, error)
}

type KeywordFilter struct {
	groups  []config.KeywordGroup
	filters []string
}

func NewKeywordFilter(groups []config.KeywordGroup, filters []string) *KeywordFilter {
	return &KeywordFilter{
		groups:  groups,
		filters: filters,
	}
}

func (f *KeywordFilter) Filter(allData map[string][]*model.NewsItem) (map[string][]*model.NewsItem, error) {
	// 如果没有关键词组，返回空或者全部？原Python代码逻辑：如果没配置，显示全部。
	// 这里我们假设没配置就返回全部
	if len(f.groups) == 0 {
		return allData, nil
	}

	result := make(map[string][]*model.NewsItem)

	for sourceID, items := range allData {
		var filteredItems []*model.NewsItem
		for _, item := range items {
			matched, score, keywords, groupIndex := f.matchWithScore(item.Title)
			if matched {
				// 设置匹配信息
				item.MatchScore = score
				item.MatchedKeywords = keywords
				item.KeywordGroup = groupIndex
				filteredItems = append(filteredItems, item)
			}
		}
		if len(filteredItems) > 0 {
			result[sourceID] = filteredItems
		}
	}

	return result, nil
}

// matchWithScore 匹配标题并返回评分信息
// 返回值：是否匹配, 匹配分数, 匹配的关键词列表, 关键词组索引
func (f *KeywordFilter) matchWithScore(title string) (bool, float64, []string, int) {
	titleLower := strings.ToLower(title)

	// 1. 全局过滤词检查
	for _, filterWord := range f.filters {
		if strings.Contains(titleLower, strings.ToLower(filterWord)) {
			return false, 0, nil, -1
		}
	}

	maxScore := 0.0
	var bestMatchedKeywords []string
	bestGroupIndex := -1

	// 2. 关键词组匹配 - 遍历所有组，找到得分最高的
	for groupIdx, group := range f.groups {
		score := 0.0
		var matched []string

		// 检查必须词（每个必须词 +20分）
		allRequiredMatched := true
		if len(group.Required) > 0 {
			for _, req := range group.Required {
				if strings.Contains(titleLower, strings.ToLower(req)) {
					score += 20
					matched = append(matched, "+"+req)
				} else {
					allRequiredMatched = false
					break
				}
			}
		}

		if !allRequiredMatched {
			continue
		}

		// 检查普通词（每个普通词 +10分）
		normalMatched := false
		if len(group.Normal) > 0 {
			for _, norm := range group.Normal {
				if strings.Contains(titleLower, strings.ToLower(norm)) {
					score += 10
					matched = append(matched, norm)
					normalMatched = true
				}
			}
		} else {
			// 没有普通词，前面必须词匹配了，那就通过
			normalMatched = true
		}

		if !normalMatched {
			continue
		}

		// 考虑关键词组优先级（1-10，默认5）
		priority := group.Priority
		if priority <= 0 {
			priority = 5 // 默认优先级
		}
		// 优先级作为乘数，范围0.2-2.0
		priorityMultiplier := float64(priority) / 5.0
		score *= priorityMultiplier

		// 记录最高分的匹配
		if score > maxScore {
			maxScore = score
			bestMatchedKeywords = matched
			bestGroupIndex = groupIdx
		}
	}

	if maxScore > 0 {
		return true, maxScore, bestMatchedKeywords, bestGroupIndex
	}

	return false, 0, nil, -1
}


