package datacache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gotoailab/trendhub/internal/logger"
	"github.com/gotoailab/trendhub/internal/model"
	bolt "go.etcd.io/bbolt"
)

const (
	dailyBucket       = "daily_cache"
	incrementalBucket = "incremental_pushed"
	historyBucket     = "crawl_history" // 存储每天的抓取历史
)

// DataCache 数据缓存管理器
type DataCache struct {
	db              *bolt.DB
	mu              sync.RWMutex
	dailyCache      map[string]*model.NewsItem // 当日汇总缓存（内存）
	lastResetTime   time.Time                  // 上次重置时间
	incrementalMode bool                       // 是否启用增量模式
}

// NewDataCache 创建数据缓存
func NewDataCache(dbPath string) (*DataCache, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open cache database: %w", err)
	}

	// 创建 buckets
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(dailyBucket)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(incrementalBucket)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(historyBucket)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create buckets: %w", err)
	}

	cache := &DataCache{
		db:            db,
		dailyCache:    make(map[string]*model.NewsItem),
		lastResetTime: time.Now(),
	}

	return cache, nil
}

// Close 关闭数据库
func (dc *DataCache) Close() error {
	return dc.db.Close()
}

// generateHash 生成新闻的唯一标识
func generateHash(item *model.NewsItem) string {
	// 使用 标题 + 平台ID 生成哈希
	data := item.Title + "|" + item.SourceID
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// AddToDailyCache 添加到当日缓存（去重）
func (dc *DataCache) AddToDailyCache(items []*model.NewsItem) int {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	// 检查是否需要重置缓存（每天0点重置）
	now := time.Now()
	if !isSameDay(dc.lastResetTime, now) {
		dc.dailyCache = make(map[string]*model.NewsItem)
		dc.lastResetTime = now
	}

	addedCount := 0
	for _, item := range items {
		hash := generateHash(item)
		if _, exists := dc.dailyCache[hash]; !exists {
			// 新内容，添加到缓存
			dc.dailyCache[hash] = item
			addedCount++
		} else {
			// 已存在，更新排名信息（保留最好的排名）
			existing := dc.dailyCache[hash]
			if len(item.Ranks) > 0 && (len(existing.Ranks) == 0 || item.Ranks[0] < existing.Ranks[0]) {
				existing.Ranks = item.Ranks
			}
		}
	}

	return addedCount
}

// GetDailyCache 获取当日缓存的所有数据
func (dc *DataCache) GetDailyCache() []*model.NewsItem {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	items := make([]*model.NewsItem, 0, len(dc.dailyCache))
	for _, item := range dc.dailyCache {
		items = append(items, item)
	}

	return items
}

// GetDailyCacheCount 获取当日缓存数量
func (dc *DataCache) GetDailyCacheCount() int {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return len(dc.dailyCache)
}

// ClearDailyCache 清空当日缓存
func (dc *DataCache) ClearDailyCache() {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.dailyCache = make(map[string]*model.NewsItem)
	dc.lastResetTime = time.Now()
}

// MarkAsPushed 标记内容已推送（用于增量模式）
func (dc *DataCache) MarkAsPushed(items []*model.NewsItem) error {
	return dc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incrementalBucket))
		now := time.Now()

		for _, item := range items {
			hash := generateHash(item)
			record := map[string]interface{}{
				"hash":       hash,
				"title":      item.Title,
				"source":     item.SourceName,
				"pushed_at":  now,
				"expires_at": now.Add(7 * 24 * time.Hour), // 7天后过期
			}

			data, err := json.Marshal(record)
			if err != nil {
				continue
			}

			if err := b.Put([]byte(hash), data); err != nil {
				return err
			}
		}

		return nil
	})
}

// IsPushed 检查内容是否已推送
func (dc *DataCache) IsPushed(item *model.NewsItem) bool {
	hash := generateHash(item)
	isPushed := false

	dc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incrementalBucket))
		data := b.Get([]byte(hash))

		if data != nil {
			var record map[string]interface{}
			if err := json.Unmarshal(data, &record); err == nil {
				// 检查是否过期
				if expiresStr, ok := record["expires_at"].(string); ok {
					expires, err := time.Parse(time.RFC3339, expiresStr)
					if err == nil && time.Now().Before(expires) {
						isPushed = true
					}
				}
			}
		}

		return nil
	})

	return isPushed
}

// FilterUnpushed 过滤出未推送的内容（增量模式）
func (dc *DataCache) FilterUnpushed(items []*model.NewsItem) []*model.NewsItem {
	unpushed := make([]*model.NewsItem, 0)

	for _, item := range items {
		if !dc.IsPushed(item) {
			unpushed = append(unpushed, item)
		}
	}

	return unpushed
}

// CleanExpiredRecords 清理过期的推送记录
func (dc *DataCache) CleanExpiredRecords() (int, error) {
	deleted := 0

	err := dc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incrementalBucket))
		c := b.Cursor()
		now := time.Now()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var record map[string]interface{}
			if err := json.Unmarshal(v, &record); err != nil {
				continue
			}

			if expiresStr, ok := record["expires_at"].(string); ok {
				expires, err := time.Parse(time.RFC3339, expiresStr)
				if err == nil && now.After(expires) {
					if err := b.Delete(k); err != nil {
						return err
					}
					deleted++
				}
			}
		}

		return nil
	})

	return deleted, err
}

// GetPushedCount 获取已推送记录数量
func (dc *DataCache) GetPushedCount() (int, error) {
	count := 0

	err := dc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incrementalBucket))
		count = b.Stats().KeyN
		return nil
	})

	return count, err
}

// isSameDay 检查两个时间是否是同一天
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// CrawlHistoryRecord 抓取历史记录
type CrawlHistoryRecord struct {
	Date      string                       `json:"date"`       // YYYY-MM-DD
	Timestamp time.Time                    `json:"timestamp"`  // 抓取时间
	Data      map[string][]*model.NewsItem `json:"data"`       // 按平台分组的数据
	ItemCount int                          `json:"item_count"` // 总条目数
}

// SaveCrawlHistory 保存抓取历史
func (dc *DataCache) SaveCrawlHistory(data map[string][]*model.NewsItem) error {
	date := time.Now().Format("2006-01-02")
	
	// 计算总条目数
	totalItems := 0
	for _, items := range data {
		totalItems += len(items)
	}

	record := CrawlHistoryRecord{
		Date:      date,
		Timestamp: time.Now(),
		Data:      data,
		ItemCount: totalItems,
	}

	return dc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", historyBucket)
		}

		// 使用日期作为key
		key := []byte(date)
		
		jsonData, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal history: %w", err)
		}

		return b.Put(key, jsonData)
	})
}

// GetCrawlHistory 获取指定日期的抓取历史
func (dc *DataCache) GetCrawlHistory(date string) (*CrawlHistoryRecord, error) {
	var record CrawlHistoryRecord

	err := dc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", historyBucket)
		}

		data := b.Get([]byte(date))
		if data == nil {
			return fmt.Errorf("no history found for date %s", date)
		}

		return json.Unmarshal(data, &record)
	})

	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetRecentCrawlHistory 获取最近N天的抓取历史列表（仅摘要）
func (dc *DataCache) GetRecentCrawlHistory(days int) ([]map[string]interface{}, error) {
	var histories []map[string]interface{}

	err := dc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		// 从最新的开始倒序遍历
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if len(histories) >= days {
				break
			}

			var record CrawlHistoryRecord
			if err := json.Unmarshal(v, &record); err != nil {
				logger.Errorf("Error unmarshalling history %s: %v", string(k), err)
				continue
			}

			// 只返回摘要信息，不包含完整数据
			summary := map[string]interface{}{
				"date":       record.Date,
				"timestamp":  record.Timestamp,
				"item_count": record.ItemCount,
			}

			// 统计各平台数量
			platformCounts := make(map[string]int)
			for platform, items := range record.Data {
				platformCounts[platform] = len(items)
			}
			summary["platforms"] = platformCounts

			histories = append(histories, summary)
		}

		return nil
	})

	return histories, err
}

// CleanOldHistory 清理N天前的历史记录
func (dc *DataCache) CleanOldHistory(retentionDays int) (int, error) {
	if retentionDays <= 0 {
		return 0, nil
	}

	threshold := time.Now().AddDate(0, 0, -retentionDays).Format("2006-01-02")
	deleted := 0

	err := dc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			dateStr := string(k)
			if dateStr < threshold {
				if err := b.Delete(k); err != nil {
					return err
				}
				deleted++
			}
		}

		return nil
	})

	return deleted, err
}

