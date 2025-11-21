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
	"github.com/gotoailab/trendhub/internal/collector"
	"github.com/gotoailab/trendhub/internal/crawler"
	"github.com/gotoailab/trendhub/internal/datacache"
	"github.com/gotoailab/trendhub/internal/filter"
	"github.com/gotoailab/trendhub/internal/model"
	"github.com/gotoailab/trendhub/internal/notifier"
	"github.com/gotoailab/trendhub/internal/pushdb"
	"github.com/gotoailab/trendhub/internal/rank"
	"github.com/gotoailab/trendhub/internal/scheduler"
)

// TaskRunner 负责执行任务
type TaskRunner struct {
	ConfigPath     string
	KeywordPath    string
	PushDB         *pushdb.PushDB
	DataCache      *datacache.DataCache
	Scheduler      *scheduler.Scheduler
	DailyCollector *collector.DailyCollector
	mu             sync.Mutex
	IsRunning      bool
	LastLog        string
	LastRunTime    time.Time
	ExtraWriter    io.Writer // 额外的日志输出目标（如 os.Stdout）
	logFilePath    string    // 日志文件路径
}

func NewTaskRunner(configPath, keywordPath string, pushDB *pushdb.PushDB, dataCache *datacache.DataCache) *TaskRunner {
	return &TaskRunner{
		ConfigPath:  configPath,
		KeywordPath: keywordPath,
		PushDB:      pushDB,
		DataCache:   dataCache,
	}
}

// SetLogFilePath 设置日志文件路径
func (tr *TaskRunner) SetLogFilePath(path string) {
	tr.logFilePath = path
}

// GetLogFilePath 获取日志文件路径
func (tr *TaskRunner) GetLogFilePath() string {
	return tr.logFilePath
}

// StartScheduler 启动定时调度器
func (tr *TaskRunner) StartScheduler(ctx context.Context) error {
	cfg, err := config.LoadConfig(tr.ConfigPath, tr.KeywordPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 如果是 daily 模式，启动持续收集器
	if cfg.Config.Report.Mode == "daily" {
		tr.DailyCollector = collector.NewDailyCollector(cfg.Config, tr.DataCache)
		tr.DailyCollector.Start(ctx)
		log.Println("Daily mode: continuous data collection started")
	}

	if !cfg.Config.Notification.PushWindow.Enabled {
		log.Println("Push window disabled, scheduler will not start")
		// 即使不启用定时推送，daily 收集器也会运行
		return nil
	}

	// 创建任务函数
	taskFunc := func() (int, error) {
		itemCount, err := tr.RunWithCount()
		if err != nil {
			return 0, err
		}
		return itemCount, nil
	}

	tr.Scheduler = scheduler.NewScheduler(&cfg.Config.Notification, tr.PushDB, taskFunc)
	return tr.Scheduler.Start(ctx)
}

// StopScheduler 停止定时调度器
func (tr *TaskRunner) StopScheduler() {
	if tr.Scheduler != nil {
		tr.Scheduler.Stop()
	}
	if tr.DailyCollector != nil {
		tr.DailyCollector.Stop()
	}
}

// ReloadConfig 重新加载配置并更新调度器和收集器
func (tr *TaskRunner) ReloadConfig(ctx context.Context) error {
	log.Println("Reloading configuration for TaskRunner...")

	// 加载新配置
	cfg, err := config.LoadConfig(tr.ConfigPath, tr.KeywordPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	// 1. 处理 Daily Collector
	if cfg.Config.Report.Mode == "daily" {
		if tr.DailyCollector == nil {
			// 之前没有收集器，创建并启动
			log.Println("Creating and starting daily collector (mode changed to daily)...")
			tr.DailyCollector = collector.NewDailyCollector(cfg.Config, tr.DataCache)
			tr.DailyCollector.Start(ctx)
		} else {
			// 已有收集器，重载配置
			if err := tr.DailyCollector.ReloadConfig(cfg.Config); err != nil {
				log.Printf("Failed to reload daily collector config: %v", err)
			}
		}
	} else {
		// 不是 daily 模式，停止收集器
		if tr.DailyCollector != nil {
			log.Println("Stopping daily collector (mode changed from daily)...")
			tr.DailyCollector.Stop()
			tr.DailyCollector = nil
		}
	}

	// 2. 处理 Scheduler
	if tr.Scheduler != nil {
		if err := tr.Scheduler.ReloadConfig(&cfg.Config.Notification); err != nil {
			log.Printf("Failed to reload scheduler config: %v", err)
		}
	} else if cfg.Config.Notification.PushWindow.Enabled {
		// 之前没有调度器，但现在启用了，需要创建并启动
		log.Println("Creating and starting scheduler (push window enabled)...")
		taskFunc := func() (int, error) {
			itemCount, err := tr.RunWithCount()
			return itemCount, err
		}
		tr.Scheduler = scheduler.NewScheduler(&cfg.Config.Notification, tr.PushDB, taskFunc)
		if err := tr.Scheduler.Start(ctx); err != nil {
			log.Printf("Failed to start scheduler: %v", err)
		}
	}

	log.Println("Configuration reloaded successfully")
	return nil
}

// RunWithCount 执行任务并返回推送数量
func (tr *TaskRunner) RunWithCount() (int, error) {
	logOutput, err := tr.Run()
	if err != nil {
		return 0, err
	}
	
	// 这里简化处理，实际可以从日志中解析数量
	_ = logOutput
	return 0, nil
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
	logger.Printf("Config loaded. Mode: %s, Platforms: %d, Keywords Groups: %d\n", 
		cfg.Config.Report.Mode, len(cfg.Config.Platforms), len(cfg.KeywordGroups))

	// 2. 初始化模块
	f := filter.NewKeywordFilter(cfg.KeywordGroups, cfg.GlobalFilters)
	r := rank.NewWeightedRanker(cfg.Config.Weight)
	n := notifier.NewNotificationManager(cfg.Config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var rawData map[string][]*model.NewsItem

	// 3. 根据模式获取数据
	switch cfg.Config.Report.Mode {
	case "daily":
		// 当日汇总模式：使用缓存的数据
		logger.Println("Mode: Daily aggregation - using cached data")
		if tr.DataCache == nil {
			errMsg := "Data cache not initialized for daily mode"
			logger.Println(errMsg)
			tr.LastLog = logBuf.String()
			return tr.LastLog, fmt.Errorf(errMsg)
		}

		cachedItems := tr.DataCache.GetDailyCache()
		logger.Printf("Retrieved %d items from daily cache", len(cachedItems))

		// 将缓存数据转换为按平台分组的格式
		rawData = make(map[string][]*model.NewsItem)
		for _, item := range cachedItems {
			rawData[item.SourceID] = append(rawData[item.SourceID], item)
		}

	case "incremental":
		// 增量监控模式：爬取当前数据，过滤已推送的
		logger.Println("Mode: Incremental monitoring - fetching and filtering new items")
		if tr.DataCache == nil {
			errMsg := "Data cache not initialized for incremental mode"
			logger.Println(errMsg)
			tr.LastLog = logBuf.String()
			return tr.LastLog, fmt.Errorf(errMsg)
		}

		c := crawler.NewNewsNowCrawler(cfg.Config)
		data, err := c.Run(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("Crawler failed: %v", err)
			logger.Println(errMsg)
			tr.LastLog = logBuf.String()
			return tr.LastLog, err
		}
		logger.Printf("Crawled data from %d platforms", len(data))

		// 保存原始爬取数据到历史记录
		if tr.DataCache != nil {
			if err := tr.DataCache.SaveCrawlHistory(data); err != nil {
				logger.Printf("Warning: Failed to save crawl history: %v", err)
			}
		}

		// 过滤出未推送的内容
		rawData = make(map[string][]*model.NewsItem)
		totalItems := 0
		newItems := 0
		for platform, items := range data {
			totalItems += len(items)
			unpushed := tr.DataCache.FilterUnpushed(items)
			if len(unpushed) > 0 {
				rawData[platform] = unpushed
				newItems += len(unpushed)
			}
		}
		logger.Printf("Found %d new items (total: %d, already pushed: %d)", 
			newItems, totalItems, totalItems-newItems)

	default: // "current" 或其他
		// 当前榜单模式（默认）：实时爬取
		logger.Println("Mode: Current ranking - fetching real-time data")
		c := crawler.NewNewsNowCrawler(cfg.Config)
		data, err := c.Run(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("Crawler failed: %v", err)
			logger.Println(errMsg)
			tr.LastLog = logBuf.String()
			return tr.LastLog, err
		}
		logger.Printf("Crawled data from %d platforms", len(data))

		// 保存原始爬取数据到历史记录
		if tr.DataCache != nil {
			if err := tr.DataCache.SaveCrawlHistory(data); err != nil {
				logger.Printf("Warning: Failed to save crawl history: %v", err)
			}
		}

		rawData = data
	}

	// 4. 关键词过滤
	filteredData, err := f.Filter(rawData)
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

	if totalItems == 0 {
		logger.Println("No matching items found, skipping notification")
		tr.LastLog = logBuf.String()
		return tr.LastLog, nil
	}

	// 5. 排序
	rankedItems := r.Rank(filteredData)

	// 6. 推送通知
	if cfg.Config.Notification.EnableNotification {
		logger.Printf("Sending notifications for %d items...", len(rankedItems))
		n.SendAll(ctx, rankedItems)
		logger.Println("Notification sent")

		// 7. 增量模式下标记已推送
		if cfg.Config.Report.Mode == "incremental" && tr.DataCache != nil {
			if err := tr.DataCache.MarkAsPushed(rankedItems); err != nil {
				logger.Printf("Warning: Failed to mark items as pushed: %v", err)
			} else {
				logger.Printf("Marked %d items as pushed", len(rankedItems))
			}

			// 清理过期记录
			if deleted, err := tr.DataCache.CleanExpiredRecords(); err == nil && deleted > 0 {
				logger.Printf("Cleaned %d expired push records", deleted)
			}
		}
	} else {
		logger.Println("Notification disabled")
	}

	logger.Println("Task completed.")
	tr.LastLog = logBuf.String()
	return tr.LastLog, nil
}

// FilterAndRankData 对原始数据进行过滤和排序
func (tr *TaskRunner) FilterAndRankData(rawData map[string][]*model.NewsItem) ([]*model.NewsItem, error) {
	// 加载配置
	cfg, err := config.LoadConfig(tr.ConfigPath, tr.KeywordPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化过滤器和排序器
	f := filter.NewKeywordFilter(cfg.KeywordGroups, cfg.GlobalFilters)
	r := rank.NewWeightedRanker(cfg.Config.Weight)

	// 过滤数据
	filteredData, err := f.Filter(rawData)
	if err != nil {
		return nil, fmt.Errorf("filter failed: %w", err)
	}

	// 排序
	rankedItems := r.Rank(filteredData)

	return rankedItems, nil
}
