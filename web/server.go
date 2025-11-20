package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gotoailab/trendhub/config"
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
	http.HandleFunc("/api/status", s.enableCors(s.handleStatus))
	http.HandleFunc("/api/run", s.enableCors(s.handleRun))
	http.HandleFunc("/api/push-records", s.enableCors(s.handlePushRecords))

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

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.Runner.mu.Lock()
	defer s.Runner.mu.Unlock()

	status := map[string]interface{}{
		"isRunning":   s.Runner.IsRunning,
		"lastRunTime": s.Runner.LastRunTime.Format(time.RFC3339),
		"lastLog":     s.Runner.LastLog,
	}
	json.NewEncoder(w).Encode(status)
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
