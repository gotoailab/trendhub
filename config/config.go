package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gotoailab/trendhub/internal/model"
	"gopkg.in/yaml.v3"
)

// AppConfig 应用配置
type AppConfig struct {
	VersionCheckURL   string `yaml:"version_check_url" json:"version_check_url"`
	ShowVersionUpdate bool   `yaml:"show_version_update" json:"show_version_update"`
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	RequestInterval int    `yaml:"request_interval" json:"request_interval"`
	EnableCrawler   bool   `yaml:"enable_crawler" json:"enable_crawler"`
	UseProxy        bool   `yaml:"use_proxy" json:"use_proxy"`
	DefaultProxy    string `yaml:"default_proxy" json:"default_proxy"`
}

// ReportConfig 报告配置
type ReportConfig struct {
	Mode          string `yaml:"mode" json:"mode"`
	RankThreshold int    `yaml:"rank_threshold" json:"rank_threshold"`
}

// PushWindowConfig 推送窗口配置
type PushWindowConfig struct {
	Enabled   bool `yaml:"enabled" json:"enabled"`
	TimeRange struct {
		Start string `yaml:"start" json:"start"`
		End   string `yaml:"end" json:"end"`
	} `yaml:"time_range" json:"time_range"`
	OncePerDay              bool `yaml:"once_per_day" json:"once_per_day"`
	PushRecordRetentionDays int  `yaml:"push_record_retention_days" json:"push_record_retention_days"`
}

// WebhooksConfig Webhook配置
type WebhooksConfig struct {
	FeishuURL        string `yaml:"feishu_url" json:"feishu_url"`
	DingtalkURL      string `yaml:"dingtalk_url" json:"dingtalk_url"`
	WeworkURL        string `yaml:"wework_url" json:"wework_url"`
	TelegramBotToken string `yaml:"telegram_bot_token" json:"telegram_bot_token"`
	TelegramChatID   string `yaml:"telegram_chat_id" json:"telegram_chat_id"`
	EmailFrom        string `yaml:"email_from" json:"email_from"`
	EmailPassword    string `yaml:"email_password" json:"email_password"`
	EmailTo          string `yaml:"email_to" json:"email_to"`
	EmailSMTPServer  string `yaml:"email_smtp_server" json:"email_smtp_server"`
	EmailSMTPPort    string `yaml:"email_smtp_port" json:"email_smtp_port"`
	NtfyServerURL    string `yaml:"ntfy_server_url" json:"ntfy_server_url"`
	NtfyTopic        string `yaml:"ntfy_topic" json:"ntfy_topic"`
	NtfyToken        string `yaml:"ntfy_token" json:"ntfy_token"`
	BarkServerURL    string `yaml:"bark_server_url" json:"bark_server_url"`
	BarkDeviceKey    string `yaml:"bark_device_key" json:"bark_device_key"`
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	EnableNotification     bool             `yaml:"enable_notification" json:"enable_notification"`
	MessageBatchSize       int              `yaml:"message_batch_size" json:"message_batch_size"`
	DingtalkBatchSize      int              `yaml:"dingtalk_batch_size" json:"dingtalk_batch_size"`
	FeishuBatchSize        int              `yaml:"feishu_batch_size" json:"feishu_batch_size"`
	BatchSendInterval      int              `yaml:"batch_send_interval" json:"batch_send_interval"`
	FeishuMessageSeparator string           `yaml:"feishu_message_separator" json:"feishu_message_separator"`
	PushWindow             PushWindowConfig `yaml:"push_window" json:"push_window"`
	Webhooks               WebhooksConfig   `yaml:"webhooks" json:"webhooks"`
}

// WeightConfig 权重配置
type WeightConfig struct {
	RankWeight      float64 `yaml:"rank_weight" json:"rank_weight"`
	FrequencyWeight float64 `yaml:"frequency_weight" json:"frequency_weight"`
	HotnessWeight   float64 `yaml:"hotness_weight" json:"hotness_weight"`
}

// Config 总配置结构
type Config struct {
	App          AppConfig          `yaml:"app" json:"app"`
	Crawler      CrawlerConfig      `yaml:"crawler" json:"crawler"`
	Report       ReportConfig       `yaml:"report" json:"report"`
	Notification NotificationConfig `yaml:"notification" json:"notification"`
	Weight       WeightConfig       `yaml:"weight" json:"weight"`
	Platforms    []model.Platform   `yaml:"platforms" json:"platforms"`
}

// KeywordGroup 关键词组
type KeywordGroup struct {
	Required []string
	Normal   []string
	Filters  []string // 该组特定的过滤词（虽然原始实现是全局的，但这里保留扩展性）
	GroupKey string
}

// GlobalConfig 全局配置管理器
type GlobalConfig struct {
	Config        *Config
	KeywordGroups []KeywordGroup
	GlobalFilters []string
}

// LoadConfig 加载配置
func LoadConfig(configPath string, keywordPath string) (*GlobalConfig, error) {
	// 1. 读取 config.yaml
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file failed: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file failed: %w", err)
	}

	// 环境变量覆盖 (示例: 仅覆盖部分关键配置)
	if envVal := os.Getenv("FEISHU_WEBHOOK_URL"); envVal != "" {
		cfg.Notification.Webhooks.FeishuURL = envVal
	}
	if envVal := os.Getenv("DINGTALK_WEBHOOK_URL"); envVal != "" {
		cfg.Notification.Webhooks.DingtalkURL = envVal
	}
	if envVal := os.Getenv("WEWORK_WEBHOOK_URL"); envVal != "" {
		cfg.Notification.Webhooks.WeworkURL = envVal
	}
	if envVal := os.Getenv("TELEGRAM_BOT_TOKEN"); envVal != "" {
		cfg.Notification.Webhooks.TelegramBotToken = envVal
	}
	if envVal := os.Getenv("TELEGRAM_CHAT_ID"); envVal != "" {
		cfg.Notification.Webhooks.TelegramChatID = envVal
	}
	// ... 其他环境变量覆盖逻辑 ...

	// 2. 读取 frequency_words.txt
	kwGroups, filters, err := loadKeywords(keywordPath)
	if err != nil {
		// 如果文件不存在，可能只是没有配置关键词，不一定是错误，视需求而定
		// 这里假设文件必须存在，或者至少能处理空文件
		fmt.Printf("Warning: loading keywords failed or file missing: %v\n", err)
		kwGroups = []KeywordGroup{}
		filters = []string{}
	}

	return &GlobalConfig{
		Config:        &cfg,
		KeywordGroups: kwGroups,
		GlobalFilters: filters,
	}, nil
}

func loadKeywords(path string) ([]KeywordGroup, []string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var groups []KeywordGroup
	var globalFilters []string

	// 按空行分割组
	rawGroups := strings.Split(string(content), "\n\n")

	for _, rawGroup := range rawGroups {
		rawGroup = strings.TrimSpace(rawGroup)
		if rawGroup == "" {
			continue
		}

		lines := strings.Split(rawGroup, "\n")
		var required []string
		var normal []string
		// 这里 Python 原版逻辑：!开头的是过滤词。
		// 原版逻辑里，filter_words 是全局收集的，group里也会标记。
		// 我们这里遵循原版逻辑：!开头的词会被加入到该组的过滤列表，同时也会被加入到返回值里的 globalFilters (如果是全局行为)
		// 看原版代码：filter_words 列表是收集所有以 ! 开头的词。

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "!") {
				globalFilters = append(globalFilters, strings.TrimPrefix(line, "!"))
			} else if strings.HasPrefix(line, "+") {
				required = append(required, strings.TrimPrefix(line, "+"))
			} else {
				normal = append(normal, line)
			}
		}

		key := ""
		if len(normal) > 0 {
			key = strings.Join(normal, " ")
		} else {
			key = strings.Join(required, " ")
		}

		if len(required) > 0 || len(normal) > 0 {
			groups = append(groups, KeywordGroup{
				Required: required,
				Normal:   normal,
				GroupKey: key,
			})
		}
	}

	return groups, globalFilters, nil
}
