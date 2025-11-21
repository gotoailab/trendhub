package model

// NewsItem 代表一条新闻数据
type NewsItem struct {
	Title           string   `json:"title"`
	URL             string   `json:"url"`
	MobileURL       string   `json:"mobileUrl"`
	Ranks           []int    `json:"ranks"`            // 在不同时间点的排名或多次抓取的排名
	SourceID        string   `json:"source_id"`        // 来源平台ID
	SourceName      string   `json:"source_name"`      // 来源平台名称
	FirstSeen       string   `json:"first_seen"`       // 首次发现时间
	LastSeen        string   `json:"last_seen"`        // 最后一次发现时间
	AppearCount     int      `json:"appear_count"`     // 出现次数
	IsNew           bool     `json:"is_new"`           // 是否是新增
	MatchScore      float64  `json:"match_score"`      // 关键词匹配分数
	MatchedKeywords []string `json:"matched_keywords"` // 匹配到的关键词列表
	KeywordGroup    int      `json:"keyword_group"`    // 匹配的关键词组索引
}

// Platform 代表一个监控平台
type Platform struct {
	ID     string  `yaml:"id" json:"id"`
	Name   string  `yaml:"name" json:"name"`
	Weight float64 `yaml:"weight" json:"weight"` // 平台权重，默认1.0，范围0-1
}

// Stats 统计数据
type Stats struct {
	TotalProcessed int
	NewItems       int
}

