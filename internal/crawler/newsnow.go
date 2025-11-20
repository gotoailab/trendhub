package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/model"
)

// Crawler 定义爬虫接口
type Crawler interface {
	Run(ctx context.Context) (map[string][]*model.NewsItem, error)
}

// NewsNowCrawler 实现 Crawler 接口
type NewsNowCrawler struct {
	cfg    *config.Config
	client *http.Client
}

func NewNewsNowCrawler(cfg *config.Config) *NewsNowCrawler {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// 如果配置了代理，这里需要设置 Transport (省略具体代理实现细节，仅预留)
	
	return &NewsNowCrawler{
		cfg:    cfg,
		client: client,
	}
}

// NewsNowResponse API 响应结构
type NewsNowResponse struct {
	Status string `json:"status"`
	Items  []struct {
		Title     string `json:"title"`
		URL       string `json:"url"`
		MobileURL string `json:"mobileUrl"`
	} `json:"items"`
}

func (c *NewsNowCrawler) Run(ctx context.Context) (map[string][]*model.NewsItem, error) {
	results := make(map[string][]*model.NewsItem)

	for _, platform := range c.cfg.Platforms {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		items, err := c.fetchPlatform(ctx, platform)
		if err != nil {
			fmt.Printf("Error fetching %s (%s): %v\n", platform.Name, platform.ID, err)
			continue
		}
		results[platform.ID] = items

		// 随机延迟，避免请求过快
		interval := c.cfg.Crawler.RequestInterval
		if interval > 0 {
			delay := time.Duration(interval+rand.Intn(200)) * time.Millisecond
			time.Sleep(delay)
		}
	}

	return results, nil
}

func (c *NewsNowCrawler) fetchPlatform(ctx context.Context, platform model.Platform) ([]*model.NewsItem, error) {
	url := fmt.Sprintf("https://newsnow.busiyi.world/api/s?id=%s&latest", platform.ID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp NewsNowResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	var newsItems []*model.NewsItem
	for i, item := range apiResp.Items {
		newsItems = append(newsItems, &model.NewsItem{
			Title:      item.Title,
			URL:        item.URL,
			MobileURL:  item.MobileURL,
			Ranks:      []int{i + 1}, // 原始排名
			SourceID:   platform.ID,
			SourceName: platform.Name,
			FirstSeen:  time.Now().Format("15:04"), // 简单记录时间
			IsNew:      true, // 初始默认为新，后续由 Filter 模块判断
		})
	}

	return newsItems, nil
}

