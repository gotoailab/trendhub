package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/model"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Runner *TaskRunner
}

func NewServer(runner *TaskRunner) *Server {
	return &Server{Runner: runner}
}

func (s *Server) Run(addr string) error {
	http.HandleFunc("/api/config", s.enableCors(s.handleConfig))
	http.HandleFunc("/api/keywords", s.enableCors(s.handleKeywords))
	http.HandleFunc("/api/run", s.enableCors(s.handleRun))
	http.HandleFunc("/api/push-records", s.enableCors(s.handlePushRecords))
	http.HandleFunc("/api/crawl-history", s.enableCors(s.handleCrawlHistory))
	http.HandleFunc("/api/crawl-history/recent", s.enableCors(s.handleRecentHistory))
	http.HandleFunc("/api/version", s.enableCors(s.handleVersion))

	// 静态文件服务
	// 假设 web/static 在运行目录的相对路径下
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fs)

	fmt.Printf("Web server listening on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

type ConfigRequest struct {
	Content string `json:"content"` // 暂时保留字符串兼容，前端如果传json对象，则用 interface{}
	// 为了支持结构化编辑，我们也许可以扩展API，但最简单的是前端解析后发回yaml字符串
	// 或者后端提供结构化API。
	// 用户要求"表单来配"，前端解析yaml为json生成表单比较好，保存时前端转回yaml或者发json给后端转yaml。
	// 这里我们让后端支持接收 JSON 配置并转回 YAML 保存
	JsonConfig interface{} `json:"jsonConfig"`
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	path := s.Runner.ConfigPath
	if r.Method == "GET" {
		content, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// 尝试解析为 JSON 返回给前端方便生成表单
		var cfg config.Config
		if err := yaml.Unmarshal(content, &cfg); err == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"yaml": string(content),
				"json": cfg,
			})
		} else {
			// 解析失败降级为纯文本
			json.NewEncoder(w).Encode(map[string]string{"yaml": string(content)})
		}
	} else if r.Method == "POST" {
		var req ConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var yamlBytes []byte
		var err error

		if req.JsonConfig != nil {
			// 如果前端传了结构化数据，转为 YAML
			yamlBytes, err = yaml.Marshal(req.JsonConfig)
			if err != nil {
				http.Error(w, "Failed to marshal json to yaml: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			yamlBytes = []byte(req.Content)
		}

		// 简单的备份
		os.WriteFile(path+".bak", yamlBytes, 0644) // ignore err for simplicity
		if err := os.WriteFile(path, yamlBytes, 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 重新加载配置并更新调度器和收集器
		if err := s.Runner.ReloadConfig(r.Context()); err != nil {
			log.Printf("Warning: Failed to reload configuration: %v", err)
			// 仍然返回成功，因为文件已经保存
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleKeywords(w http.ResponseWriter, r *http.Request) {
	path := s.Runner.KeywordPath
	if r.Method == "GET" {
		content, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				json.NewEncoder(w).Encode(map[string]string{"content": ""})
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"content": string(content)})
	} else if r.Method == "POST" {
		var req ConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 关键词仍然是文本格式比较方便，或者前端解析后拼装回文本
		if err := os.WriteFile(path, []byte(req.Content), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCrawlHistory 获取指定日期的爬取历史数据（已过滤）
func (s *Server) handleCrawlHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.Runner.DataCache == nil {
		http.Error(w, "Data cache not initialized", http.StatusInternalServerError)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		// 默认返回今天的数据
		date = time.Now().Format("2006-01-02")
	}

	// 获取原始历史数据
	history, err := s.Runner.DataCache.GetCrawlHistory(date)
	if err != nil {
		// 如果没有历史记录，返回空数据
		json.NewEncoder(w).Encode(map[string]interface{}{
			"date":       date,
			"timestamp":  time.Now(),
			"items":      []interface{}{},
			"item_count": 0,
		})
		return
	}

	// 应用过滤和排序
	filteredItems, err := s.Runner.FilterAndRankData(history.Data)
	if err != nil {
		log.Printf("Failed to filter data: %v", err)
		// 如果过滤失败，返回未过滤的数据
		var allItems []*model.NewsItem
		for _, items := range history.Data {
			allItems = append(allItems, items...)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"date":       history.Date,
			"timestamp":  history.Timestamp,
			"items":      allItems,
			"item_count": len(allItems),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"date":       history.Date,
		"timestamp":  history.Timestamp,
		"items":      filteredItems,
		"item_count": len(filteredItems),
	})
}

// handleRecentHistory 获取最近7天的抓取历史摘要
func (s *Server) handleRecentHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.Runner.DataCache == nil {
		http.Error(w, "Data cache not initialized", http.StatusInternalServerError)
		return
	}

	days := 7
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		fmt.Sscanf(daysStr, "%d", &days)
	}

	histories, err := s.Runner.DataCache.GetRecentCrawlHistory(days)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get recent history: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"histories": histories,
		"total":     len(histories),
	})
}


func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 异步运行
	go func() {
		_, err := s.Runner.Run()
		if err != nil {
			fmt.Printf("Manual run failed: %v\n", err)
		}
	}()

	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func (s *Server) handlePushRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.Runner.PushDB == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"records": []interface{}{},
			"total":   0,
		})
		return
	}

	// 获取分页参数
	limit := 20
	offset := 0
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}

	records, err := s.Runner.PushDB.GetRecords(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := s.Runner.PushDB.GetRecordCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"records": records,
		"total":   total,
	})
}

// handleVersion 处理版本检查请求
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 加载配置以获取版本检查URL
	cfg, err := config.LoadConfig(s.Runner.ConfigPath, s.Runner.KeywordPath)
	if err != nil {
		// 如果加载配置失败，仍然返回当前版本
		currentVersion, _ := config.GetCurrentVersion()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"latest_version":  "",
			"has_update":      false,
			"error":           "Failed to load config",
		})
		return
	}

	// 如果未启用版本更新提示，直接返回当前版本
	if !cfg.Config.App.ShowVersionUpdate {
		currentVersion, _ := config.GetCurrentVersion()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"latest_version":  "",
			"has_update":      false,
		})
		return
	}

	// 检查版本更新
	versionInfo, err := config.CheckVersionUpdate(cfg.Config.App.VersionCheckURL)
	if err != nil {
		// 如果检查失败，返回当前版本信息
		currentVersion, _ := config.GetCurrentVersion()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current_version": currentVersion,
			"latest_version":  "",
			"has_update":      false,
			"error":           err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(versionInfo)
}
