package collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/crawler"
	"github.com/gotoailab/trendhub/internal/datacache"
)

// DailyCollector 当日汇总数据收集器
type DailyCollector struct {
	crawler   crawler.Crawler
	cache     *datacache.DataCache
	interval  time.Duration
	isRunning bool
	stopChan  chan struct{}
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
		crawler:  c,
		cache:    cache,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 启动持续收集
func (dc *DailyCollector) Start(ctx context.Context) {
	dc.mu.Lock()
	if dc.isRunning {
		dc.mu.Unlock()
		return
	}
	dc.isRunning = true
	dc.mu.Unlock()

	log.Println("Daily collector started, collecting data every", dc.interval)

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

	dc.isRunning = false
	close(dc.stopChan)
	log.Println("Daily collector stopped")
}

// collect 执行一次数据收集
func (dc *DailyCollector) collect(ctx context.Context) {
	log.Println("Collecting data for daily aggregation...")

	// 创建超时上下文
	collectCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// 爬取数据
	data, err := dc.crawler.Run(collectCtx)
	if err != nil {
		log.Printf("Failed to collect data: %v", err)
		return
	}

	// 添加到缓存（自动去重）
	totalAdded := 0
	for _, items := range data {
		added := dc.cache.AddToDailyCache(items)
		totalAdded += added
	}

	cacheCount := dc.cache.GetDailyCacheCount()
	log.Printf("Collected data: added %d new items, total cached: %d items", totalAdded, cacheCount)
}

// IsRunning 检查是否正在运行
func (dc *DailyCollector) IsRunning() bool {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return dc.isRunning
}

