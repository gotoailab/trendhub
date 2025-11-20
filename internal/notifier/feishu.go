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

type FeishuNotifier struct {
	webhookURL string
}

func NewFeishuNotifier(url string) *FeishuNotifier {
	return &FeishuNotifier{webhookURL: url}
}

func (n *FeishuNotifier) Name() string {
	return "Feishu"
}

type FeishuMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (n *FeishuNotifier) Send(ctx context.Context, items []*model.NewsItem) error {
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

	msg := FeishuMessage{
		MsgType: "text",
	}
	msg.Content.Text = sb.String()

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
		return fmt.Errorf("feishu api status code: %d", resp.StatusCode)
	}

	return nil
}

