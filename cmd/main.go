package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotoailab/trendhub/internal/datacache"
	"github.com/gotoailab/trendhub/internal/pushdb"
	"github.com/gotoailab/trendhub/web"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Path to config file")
	keywordPath := flag.String("keywords", "config/frequency_words.txt", "Path to keywords file")
	webMode := flag.Bool("web", false, "Run in web mode")
	webAddr := flag.String("addr", ":8080", "Web server address")
	pushDBPath := flag.String("pushdb", "data/push_records.db", "Path to push records database")
	cacheDBPath := flag.String("cachedb", "data/data_cache.db", "Path to data cache database")
	flag.Parse()

	// 初始化推送记录数据库
	pushDB, err := pushdb.NewPushDB(*pushDBPath)
	if err != nil {
		log.Fatalf("Failed to initialize push database: %v", err)
	}
	defer pushDB.Close()

	// 初始化数据缓存数据库
	dataCache, err := datacache.NewDataCache(*cacheDBPath)
	if err != nil {
		log.Fatalf("Failed to initialize data cache: %v", err)
	}
	defer dataCache.Close()

	runner := web.NewTaskRunner(*configPath, *keywordPath, pushDB, dataCache)

	if *webMode {
		// Web 模式：启动 Web 服务器和定时调度器
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 启动定时调度器
		if err := runner.StartScheduler(ctx); err != nil {
			log.Printf("Warning: Failed to start scheduler: %v", err)
		}

		// 优雅关闭
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			log.Println("Shutting down...")
			runner.StopScheduler()
			cancel()
			os.Exit(0)
		}()

		server := web.NewServer(runner)
		log.Fatal(server.Run(*webAddr))
	} else {
		// 命令行模式：直接执行任务
		runner.ExtraWriter = os.Stdout
		if _, err := runner.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
