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
			if f.match(item.Title) {
				filteredItems = append(filteredItems, item)
			}
		}
		if len(filteredItems) > 0 {
			result[sourceID] = filteredItems
		}
	}

	return result, nil
}

func (f *KeywordFilter) match(title string) bool {
	titleLower := strings.ToLower(title)

	// 1. 全局过滤词检查
	for _, filterWord := range f.filters {
		if strings.Contains(titleLower, strings.ToLower(filterWord)) {
			return false
		}
	}

	// 2. 关键词组匹配
	// 只要匹配任意一个组即可
	for _, group := range f.groups {
		// 检查必须词
		matchRequired := true
		if len(group.Required) > 0 {
			for _, req := range group.Required {
				if !strings.Contains(titleLower, strings.ToLower(req)) {
					matchRequired = false
					break
				}
			}
		}
		if !matchRequired {
			continue
		}

		// 检查普通词 (如果有普通词，必须包含至少一个；如果没有普通词但有必须词，则只要必须词匹配即可)
		matchNormal := false
		if len(group.Normal) > 0 {
			for _, norm := range group.Normal {
				if strings.Contains(titleLower, strings.ToLower(norm)) {
					matchNormal = true
					break
				}
			}
		} else {
			// 没有普通词，前面必须词匹配了，那就通过
			matchNormal = true
		}

		if matchRequired && matchNormal {
			return true
		}
	}

	return false
}

