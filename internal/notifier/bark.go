package notifier

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gotoailab/trendhub/internal/model"
)

type BarkNotifier struct {
	serverURL string
	deviceKey string
}

// NewBarkNotifier 创建 Bark 推送通知器
// serverURL: Bark 服务器地址，如 https://api.day.app
// deviceKey: 设备密钥
func NewBarkNotifier(serverURL, deviceKey string) *BarkNotifier {
	// 如果没有指定服务器，使用官方服务器
	if serverURL == "" {
		serverURL = "https://api.day.app"
	}
	// 移除末尾的斜杠
	serverURL = strings.TrimRight(serverURL, "/")
	
	return &BarkNotifier{
		serverURL: serverURL,
		deviceKey: deviceKey,
	}
}

func (n *BarkNotifier) Name() string {
	return "Bark"
}

func (n *BarkNotifier) Send(ctx context.Context, items []*model.NewsItem) error {
	if n.deviceKey == "" {
		return fmt.Errorf("bark device key is empty")
	}

	// 构建消息内容
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("热点监控 (%s)\n\n", time.Now().Format("15:04")))

	// 最多显示前 10 条
	maxItems := 10
	if len(items) > maxItems {
		items = items[:maxItems]
	}

	currentSource := ""
	for _, item := range items {
		if item.SourceName != currentSource {
			sb.WriteString(fmt.Sprintf("\n【%s】\n", item.SourceName))
			currentSource = item.SourceName
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", item.Ranks[0], item.Title))
	}

	if len(items) >= maxItems {
		sb.WriteString(fmt.Sprintf("\n...还有 %d 条", len(items)-maxItems))
	}

	title := "TrendHub 热点监控"
	body := sb.String()

	// 构建 URL
	// Bark API: GET {serverURL}/{deviceKey}/{title}/{body}?url={url}&group={group}&sound={sound}
	requestURL := fmt.Sprintf("%s/%s/%s/%s",
		n.serverURL,
		url.PathEscape(n.deviceKey),
		url.PathEscape(title),
		url.PathEscape(body),
	)

	// 添加查询参数
	params := url.Values{}
	params.Add("group", "TrendHub")
	params.Add("sound", "calypso") // 使用 Bark 的默认提示音
	if len(items) > 0 && items[0].URL != "" {
		params.Add("url", items[0].URL) // 点击通知跳转到第一条新闻
	}
	
	requestURL = requestURL + "?" + params.Encode()

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("create bark request failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send bark notification failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark api returned status code: %d", resp.StatusCode)
	}

	return nil
}

