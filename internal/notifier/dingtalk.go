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

type DingtalkNotifier struct {
	webhookURL string
}

func NewDingtalkNotifier(url string) *DingtalkNotifier {
	return &DingtalkNotifier{webhookURL: url}
}

func (n *DingtalkNotifier) Name() string {
	return "Dingtalk"
}

type DingtalkMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

func (n *DingtalkNotifier) Send(ctx context.Context, items []*model.NewsItem) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# TrendHub 热点监控报告 (%s)\n\n", time.Now().Format("15:04")))

	currentSource := ""
	for _, item := range items {
		if item.SourceName != currentSource {
			sb.WriteString(fmt.Sprintf("\n## %s\n", item.SourceName))
			currentSource = item.SourceName
		}
		title := item.Title
		if item.URL != "" {
			title = fmt.Sprintf("[%s](%s)", item.Title, item.URL)
		}
		sb.WriteString(fmt.Sprintf("- **%d.** %s\n", item.Ranks[0], title))
	}

	msg := DingtalkMessage{
		MsgType: "markdown",
	}
	msg.Markdown.Title = "TrendHub 热点监控"
	msg.Markdown.Text = sb.String()

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
		return fmt.Errorf("dingtalk api status code: %d", resp.StatusCode)
	}

	return nil
}

