package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/pushdb"
)

// TaskFunc 任务执行函数
type TaskFunc func() (int, error) // 返回推送的条目数量和错误

// Scheduler 定时调度器
type Scheduler struct {
	cfg       *config.NotificationConfig
	pushDB    *pushdb.PushDB
	taskFunc  TaskFunc
	ticker    *time.Ticker
	stopChan  chan struct{}
	isRunning bool
	mu        sync.RWMutex
	ctx       context.Context
}

// NewScheduler 创建调度器
func NewScheduler(cfg *config.NotificationConfig, db *pushdb.PushDB, taskFunc TaskFunc) *Scheduler {
	return &Scheduler{
		cfg:      cfg,
		pushDB:   db,
		taskFunc: taskFunc,
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时调度
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		log.Println("Scheduler is already running")
		return nil
	}

	if !s.cfg.PushWindow.Enabled {
		log.Println("Push window is disabled, scheduler will not start")
		return nil
	}

	s.ctx = ctx
	s.isRunning = true
	s.stopChan = make(chan struct{})
	
	log.Printf("Scheduler started with time window: %s - %s\n", 
		s.cfg.PushWindow.TimeRange.Start, 
		s.cfg.PushWindow.TimeRange.End)

	// 每分钟检查一次是否在推送时间窗口内
	s.ticker = time.NewTicker(1 * time.Minute)

	go func() {
		// 立即检查一次
		s.checkAndRun()

		for {
			select {
			case <-ctx.Done():
				s.Stop()
				return
			case <-s.stopChan:
				return
			case <-s.ticker.C:
				s.checkAndRun()
			}
		}
	}()

	return nil
}

// Stop 停止调度
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
	
	close(s.stopChan)
	s.isRunning = false
	log.Println("Scheduler stopped")
}

// ReloadConfig 重新加载配置
func (s *Scheduler) ReloadConfig(newCfg *config.NotificationConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Reloading scheduler configuration...")

	// 保存旧配置
	oldEnabled := s.cfg.PushWindow.Enabled
	oldStart := s.cfg.PushWindow.TimeRange.Start
	oldEnd := s.cfg.PushWindow.TimeRange.End
	oldOncePerDay := s.cfg.PushWindow.OncePerDay

	// 更新配置
	s.cfg = newCfg

	// 检测配置变化
	configChanged := oldEnabled != newCfg.PushWindow.Enabled ||
		oldStart != newCfg.PushWindow.TimeRange.Start ||
		oldEnd != newCfg.PushWindow.TimeRange.End ||
		oldOncePerDay != newCfg.PushWindow.OncePerDay

	if !configChanged {
		log.Println("Scheduler configuration unchanged, no restart needed")
		return nil
	}

	log.Printf("Scheduler configuration changed (enabled: %v -> %v, time: %s-%s -> %s-%s)",
		oldEnabled, newCfg.PushWindow.Enabled,
		oldStart, oldEnd,
		newCfg.PushWindow.TimeRange.Start, newCfg.PushWindow.TimeRange.End)

	// 如果正在运行，需要重启
	if s.isRunning {
		log.Println("Restarting scheduler with new configuration...")
		
		// 停止当前调度器（不加锁，因为已经加锁了）
		if s.ticker != nil {
			s.ticker.Stop()
			s.ticker = nil
		}
		close(s.stopChan)
		s.isRunning = false

		// 如果新配置启用了推送窗口，重新启动
		if newCfg.PushWindow.Enabled {
			s.stopChan = make(chan struct{})
			s.isRunning = true
			s.ticker = time.NewTicker(1 * time.Minute)

			go func() {
				s.checkAndRun()

				for {
					select {
					case <-s.ctx.Done():
						s.Stop()
						return
					case <-s.stopChan:
						return
					case <-s.ticker.C:
						s.checkAndRun()
					}
				}
			}()

			log.Println("Scheduler restarted successfully")
		} else {
			log.Println("Scheduler stopped (push window disabled)")
		}
	} else {
		// 如果之前没运行，但现在启用了，启动调度器
		if newCfg.PushWindow.Enabled && s.ctx != nil {
			log.Println("Starting scheduler (push window enabled)...")
			s.stopChan = make(chan struct{})
			s.isRunning = true
			s.ticker = time.NewTicker(1 * time.Minute)

			go func() {
				s.checkAndRun()

				for {
					select {
					case <-s.ctx.Done():
						s.Stop()
						return
					case <-s.stopChan:
						return
					case <-s.ticker.C:
						s.checkAndRun()
					}
				}
			}()

			log.Println("Scheduler started successfully")
		}
	}

	return nil
}

// IsRunning 返回调度器是否正在运行
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// checkAndRun 检查是否应该运行任务
func (s *Scheduler) checkAndRun() {
	now := time.Now()
	
	// 检查是否在时间窗口内
	if !s.isInTimeWindow(now) {
		return
	}

	// 如果启用了每日一次，检查今天是否已经推送过
	if s.cfg.PushWindow.OncePerDay {
		lastPushTime, err := s.pushDB.GetLastPushTime()
		if err == nil && !lastPushTime.IsZero() {
			// 检查最后一次推送是否是今天
			if isSameDay(lastPushTime, now) {
				log.Println("Already pushed today, skipping...")
				return
			}
		}
	}

	log.Println("Time window matched, executing task...")
	s.executeTask()
}

// isInTimeWindow 检查当前时间是否在推送窗口内
func (s *Scheduler) isInTimeWindow(now time.Time) bool {
	startTime, err := parseTimeOfDay(s.cfg.PushWindow.TimeRange.Start)
	if err != nil {
		log.Printf("Failed to parse start time: %v\n", err)
		return false
	}

	endTime, err := parseTimeOfDay(s.cfg.PushWindow.TimeRange.End)
	if err != nil {
		log.Printf("Failed to parse end time: %v\n", err)
		return false
	}

	currentMinutes := now.Hour()*60 + now.Minute()
	startMinutes := startTime.Hour()*60 + startTime.Minute()
	endMinutes := endTime.Hour()*60 + endTime.Minute()

	// 处理跨日的情况
	if endMinutes < startMinutes {
		// 例如 22:00 到 02:00
		return currentMinutes >= startMinutes || currentMinutes <= endMinutes
	}

	return currentMinutes >= startMinutes && currentMinutes <= endMinutes
}

// executeTask 执行任务并记录结果
func (s *Scheduler) executeTask() {
	recordID := fmt.Sprintf("%d", time.Now().UnixNano())
	startTime := time.Now()

	itemCount, err := s.taskFunc()
	duration := time.Since(startTime).Milliseconds()

	record := &pushdb.PushRecord{
		ID:        recordID,
		Timestamp: startTime,
		ItemCount: itemCount,
		Duration:  duration,
	}

	if err != nil {
		record.Status = "failed"
		record.ErrorMsg = err.Error()
		record.FailedNum = 1
		log.Printf("Task failed: %v\n", err)
	} else {
		record.Status = "success"
		record.SuccessNum = 1
		log.Printf("Task completed successfully, pushed %d items\n", itemCount)
	}

	// 保存记录
	if err := s.pushDB.SaveRecord(record); err != nil {
		log.Printf("Failed to save push record: %v\n", err)
	}

	// 清理旧记录
	if s.cfg.PushWindow.PushRecordRetentionDays > 0 {
		deleted, err := s.pushDB.DeleteOldRecords(s.cfg.PushWindow.PushRecordRetentionDays)
		if err != nil {
			log.Printf("Failed to delete old records: %v\n", err)
		} else if deleted > 0 {
			log.Printf("Deleted %d old push records\n", deleted)
		}
	}
}

// parseTimeOfDay 解析时间字符串 HH:MM
func parseTimeOfDay(timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location()), nil
}

// isSameDay 检查两个时间是否是同一天
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

