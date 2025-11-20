package pushdb

import (
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	pushBucket = "push_records"
)

// PushRecord 推送记录
type PushRecord struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"` // success, failed, partial
	ItemCount   int       `json:"item_count"`
	Notifiers   []string  `json:"notifiers"`
	SuccessNum  int       `json:"success_num"`
	FailedNum   int       `json:"failed_num"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
	Duration    int64     `json:"duration"` // 毫秒
}

// PushDB 推送记录数据库
type PushDB struct {
	db *bolt.DB
}

// NewPushDB 创建推送记录数据库
func NewPushDB(dbPath string) (*PushDB, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 创建 bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(pushBucket))
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &PushDB{db: db}, nil
}

// Close 关闭数据库
func (pdb *PushDB) Close() error {
	return pdb.db.Close()
}

// SaveRecord 保存推送记录
func (pdb *PushDB) SaveRecord(record *PushRecord) error {
	return pdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pushBucket))
		
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}
		
		return b.Put([]byte(record.ID), data)
	})
}

// GetRecord 获取单条记录
func (pdb *PushDB) GetRecord(id string) (*PushRecord, error) {
	var record PushRecord
	
	err := pdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pushBucket))
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("record not found")
		}
		return json.Unmarshal(data, &record)
	})
	
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetRecords 获取推送记录列表（分页）
func (pdb *PushDB) GetRecords(limit int, offset int) ([]*PushRecord, error) {
	var records []*PushRecord
	
	err := pdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pushBucket))
		c := b.Cursor()
		
		// 反向遍历（最新的在前）
		count := 0
		skipped := 0
		
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if skipped < offset {
				skipped++
				continue
			}
			
			if count >= limit {
				break
			}
			
			var record PushRecord
			if err := json.Unmarshal(v, &record); err != nil {
				continue
			}
			records = append(records, &record)
			count++
		}
		
		return nil
	})
	
	return records, err
}

// GetRecordCount 获取记录总数
func (pdb *PushDB) GetRecordCount() (int, error) {
	count := 0
	
	err := pdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pushBucket))
		count = b.Stats().KeyN
		return nil
	})
	
	return count, err
}

// DeleteOldRecords 删除指定天数之前的记录
func (pdb *PushDB) DeleteOldRecords(days int) (int, error) {
	deleted := 0
	cutoff := time.Now().AddDate(0, 0, -days)
	
	err := pdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pushBucket))
		c := b.Cursor()
		
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var record PushRecord
			if err := json.Unmarshal(v, &record); err != nil {
				continue
			}
			
			if record.Timestamp.Before(cutoff) {
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

// GetLastPushTime 获取最后一次推送时间（用于 once_per_day 检查）
func (pdb *PushDB) GetLastPushTime() (time.Time, error) {
	var lastTime time.Time
	
	err := pdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pushBucket))
		c := b.Cursor()
		
		k, v := c.Last()
		if k == nil {
			return nil
		}
		
		var record PushRecord
		if err := json.Unmarshal(v, &record); err != nil {
			return err
		}
		
		lastTime = record.Timestamp
		return nil
	})
	
	return lastTime, err
}

