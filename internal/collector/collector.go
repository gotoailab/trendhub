package collector

import (
	"context"
	"sync"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/crawler"
	"github.com/gotoailab/trendhub/internal/datacache"
	"github.com/gotoailab/trendhub/internal/logger"
)

// DailyCollector 当日汇总数据收集器
type DailyCollector struct {
	cfg       *config.Config
	crawler   crawler.Crawler
	cache     *datacache.DataCache
	interval  time.Duration
	isRunning bool
	stopChan  chan struct{}
	ctx       context.Context
	mu        sync.Mutex
}

// NewDailyCollector 创建数据收集器
func NewDailyCollector(cfg *config.Config, cache *datacache.DataCache) *DailyCollector {
	c := crawler.NewNewsNowCrawler(cfg)
	interval := time.Duration(cfg.Crawler.RequestInterval) * time.Millisecond
	if interval < time.Minute {
		interval = 5 * time.Minute // 默认至少5分钟间隔
	}

	return &DailyCollector{
		cfg:      cfg,
		crawler:  c,
		cache:    cache,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动持续收集
func (dc *DailyCollector) Start(ctx context.Context) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if dc.isRunning {
		logger.Info("Daily collector is already running")
		return
	}

	dc.ctx = ctx
	dc.isRunning = true
	dc.stopChan = make(chan struct{})

	logger.Infof("Daily collector started, collecting data every %v", dc.interval)

	go func() {
		// 立即执行一次
		dc.collect(ctx)

		ticker := time.NewTicker(dc.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				dc.Stop()
				return
			case <-dc.stopChan:
				return
			case <-ticker.C:
				dc.collect(ctx)
			}
		}
	}()
}

// Stop 停止收集
func (dc *DailyCollector) Stop() {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.isRunning {
		return
	}

	close(dc.stopChan)
	dc.isRunning = false
	logger.Info("Daily collector stopped")
}

// ReloadConfig 重新加载配置
func (dc *DailyCollector) ReloadConfig(newCfg *config.Config) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	logger.Info("Reloading daily collector configuration...")

	// 计算新的间隔时间
	newInterval := time.Duration(newCfg.Crawler.RequestInterval) * time.Millisecond
	if newInterval < time.Minute {
		newInterval = 5 * time.Minute
	}

	// 保存旧配置
	oldInterval := dc.interval
	oldMode := dc.cfg.Report.Mode

	// 检测配置变化
	configChanged := oldInterval != newInterval

	if oldMode != newCfg.Report.Mode {
		logger.Infof("Report mode changed: %s -> %s", oldMode, newCfg.Report.Mode)
		configChanged = true
	}

	// 更新配置
	dc.cfg = newCfg
	dc.interval = newInterval
	dc.crawler = crawler.NewNewsNowCrawler(newCfg)

	if !configChanged {
		logger.Info("Daily collector configuration unchanged")
		return nil
	}

	logger.Infof("Daily collector configuration changed (interval: %v -> %v)", oldInterval, newInterval)

	// 如果正在运行，需要重启
	if dc.isRunning {
		logger.Info("Restarting daily collector with new configuration...")
		
		// 停止当前收集器
		close(dc.stopChan)
		dc.isRunning = false

		// 如果新模式是 daily，重新启动
		if newCfg.Report.Mode == "daily" {
			dc.stopChan = make(chan struct{})
			dc.isRunning = true

			go func() {
				dc.collect(dc.ctx)

				ticker := time.NewTicker(dc.interval)
				defer ticker.Stop()

				for {
					select {
					case <-dc.ctx.Done():
						dc.Stop()
						return
					case <-dc.stopChan:
						return
					case <-ticker.C:
						dc.collect(dc.ctx)
					}
				}
			}()

			logger.Info("Daily collector restarted successfully")
		} else {
			logger.Info("Daily collector stopped (mode is not daily)")
		}
	}

	return nil
}

// collect 执行一次数据收集
func (dc *DailyCollector) collect(ctx context.Context) {
	logger.Info("Collecting data for daily aggregation...")

	// 创建超时上下文
	collectCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// 爬取数据
	data, err := dc.crawler.Run(collectCtx)
	if err != nil {
		logger.Infof("Failed to collect data: %v", err)
		return
	}

	// 保存到历史记录（覆盖今天的记录，保持最新）
	if err := dc.cache.SaveCrawlHistory(data); err != nil {
		logger.Infof("Warning: Failed to save crawl history: %v", err)
	}

	// 添加到缓存（自动去重）
	totalAdded := 0
	for _, items := range data {
		added := dc.cache.AddToDailyCache(items)
		totalAdded += added
	}

	cacheCount := dc.cache.GetDailyCacheCount()
	logger.Infof("Collected data: added %d new items, total cached: %d items", totalAdded, cacheCount)
}

// IsRunning 检查是否正在运行
func (dc *DailyCollector) IsRunning() bool {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return dc.isRunning
}

