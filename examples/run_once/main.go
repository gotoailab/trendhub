package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/crawler"
	"github.com/gotoailab/trendhub/internal/filter"
	"github.com/gotoailab/trendhub/internal/notifier"
	"github.com/gotoailab/trendhub/internal/rank"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Path to config file")
	keywordPath := flag.String("keywords", "config/frequency_words.txt", "Path to keywords file")
	flag.Parse()

	// 1. 加载配置
	cfg, err := config.LoadConfig(*configPath, *keywordPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Config loaded. Platforms: %d, Keywords Groups: %d\n", len(cfg.Config.Platforms), len(cfg.KeywordGroups))

	// 2. 初始化模块
	c := crawler.NewNewsNowCrawler(cfg.Config)
	f := filter.NewKeywordFilter(cfg.KeywordGroups, cfg.GlobalFilters)
	r := rank.NewWeightedRanker(cfg.Config.Weight)
	n := notifier.NewNotificationManager(cfg.Config)

	// 3. 执行任务 (这里演示一次性执行，如果是守护进程可以加 for loop 或 cron)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("Start crawling...")
	// 3.1 爬取
	data, err := c.Run(ctx)
	if err != nil {
		log.Printf("Crawler failed: %v", err)
		return
	}
	log.Printf("Crawled data from %d platforms", len(data))

	// 3.2 过滤
	filteredData, err := f.Filter(data)
	if err != nil {
		log.Printf("Filter failed: %v", err)
		return
	}
	
	totalItems := 0
	for _, items := range filteredData {
		totalItems += len(items)
	}
	log.Printf("Filtered data: %d items remaining", totalItems)

	// 3.3 排序
	rankedItems := r.Rank(filteredData)

	// 3.4 推送
	if cfg.Config.Notification.EnableNotification {
		log.Println("Sending notifications...")
		n.SendAll(ctx, rankedItems)
	} else {
		log.Println("Notification disabled")
	}
	
	// 打印结果到控制台供检查
	for _, item := range rankedItems {
		fmt.Printf("[%s] %d. %s\n", item.SourceName, item.Ranks[0], item.Title)
	}
	
	log.Println("Done.")
}

