package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gotoailab/trendhub/internal/model"
)

type WPSNotifier struct {
	webhookURL string
}

func NewWPSNotifier(url string) *WPSNotifier {
	return &WPSNotifier{
		webhookURL: url,
	}
}

func (n *WPSNotifier) Name() string {
	return "WPS"
}

type WPSMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func (n *WPSNotifier) Send(ctx context.Context, items []*model.NewsItem) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TrendRadar 热点监控报告 (%s)\n\n", time.Now().Format("2006-01-02 15:04")))

	// 简单的文本格式化
	currentSource := ""
	for _, item := range items {
		if item.SourceName != currentSource {
			sb.WriteString(fmt.Sprintf("\n【%s】\n", item.SourceName))
			currentSource = item.SourceName
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", item.Ranks[0], item.Title))
		if item.URL != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", item.URL))
		}
	}

	msg := WPSMessage{
		MsgType: "text",
	}
	msg.Text.Content = sb.String()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wps api status code: %d", resp.StatusCode)
	}

	return nil
}

