package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/crawler"
	"github.com/gotoailab/trendhub/internal/filter"
	"github.com/gotoailab/trendhub/internal/notifier"
	"github.com/gotoailab/trendhub/internal/rank"
)

// TaskRunner 负责执行任务
type TaskRunner struct {
	ConfigPath  string
	KeywordPath string
	mu          sync.Mutex
	IsRunning   bool
	LastLog     string
	LastRunTime time.Time
	ExtraWriter io.Writer // 额外的日志输出目标（如 os.Stdout）
}

func NewTaskRunner(configPath, keywordPath string) *TaskRunner {
	return &TaskRunner{
		ConfigPath:  configPath,
		KeywordPath: keywordPath,
	}
}

func (tr *TaskRunner) Run() (string, error) {
	tr.mu.Lock()
	if tr.IsRunning {
		tr.mu.Unlock()
		return "", fmt.Errorf("task is already running")
	}
	tr.IsRunning = true
	tr.mu.Unlock()

	defer func() {
		tr.mu.Lock()
		tr.IsRunning = false
		tr.mu.Unlock()
	}()

	var logBuf bytes.Buffer
	var writer io.Writer = &logBuf
	if tr.ExtraWriter != nil {
		writer = io.MultiWriter(&logBuf, tr.ExtraWriter)
	}
	logger := log.New(writer, "", log.LstdFlags)

	logger.Println("Task started...")
	tr.LastRunTime = time.Now()

	// 1. 加载配置
	cfg, err := config.LoadConfig(tr.ConfigPath, tr.KeywordPath)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load config: %v", err)
		logger.Println(errMsg)
		tr.LastLog = logBuf.String()
		return tr.LastLog, err
	}
	logger.Printf("Config loaded. Platforms: %d, Keywords Groups: %d\n", len(cfg.Config.Platforms), len(cfg.KeywordGroups))

	// 2. 初始化模块
	c := crawler.NewNewsNowCrawler(cfg.Config)
	f := filter.NewKeywordFilter(cfg.KeywordGroups, cfg.GlobalFilters)
	r := rank.NewWeightedRanker(cfg.Config.Weight)
	n := notifier.NewNotificationManager(cfg.Config)

	// 3. 执行任务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger.Println("Start crawling...")
	data, err := c.Run(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("Crawler failed: %v", err)
		logger.Println(errMsg)
		tr.LastLog = logBuf.String()
		return tr.LastLog, err
	}
	logger.Printf("Crawled data from %d platforms", len(data))

	filteredData, err := f.Filter(data)
	if err != nil {
		errMsg := fmt.Sprintf("Filter failed: %v", err)
		logger.Println(errMsg)
		tr.LastLog = logBuf.String()
		return tr.LastLog, err
	}

	totalItems := 0
	for _, items := range filteredData {
		totalItems += len(items)
	}
	logger.Printf("Filtered data: %d items remaining", totalItems)

	rankedItems := r.Rank(filteredData)

	if cfg.Config.Notification.EnableNotification {
		logger.Println("Sending notifications...")
		n.SendAll(ctx, rankedItems)
		logger.Println("Notification sent")
	} else {
		logger.Println("Notification disabled")
	}

	logger.Println("Task completed.")
	tr.LastLog = logBuf.String()
	return tr.LastLog, nil
}
