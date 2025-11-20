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

type TelegramNotifier struct {
	botToken string
	chatID   string
}

func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: token,
		chatID:   chatID,
	}
}

func (n *TelegramNotifier) Name() string {
	return "Telegram"
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (n *TelegramNotifier) Send(ctx context.Context, items []*model.NewsItem) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>TrendRadar 热点监控报告</b> (%s)\n", time.Now().Format("15:04")))

	currentSource := ""
	for _, item := range items {
		if item.SourceName != currentSource {
			sb.WriteString(fmt.Sprintf("\n<b>%s</b>\n", item.SourceName))
			currentSource = item.SourceName
		}
		
		title := item.Title
		// Telegram HTML mode 特殊字符转义需注意，这里简化处理
		title = strings.ReplaceAll(title, "<", "&lt;")
		title = strings.ReplaceAll(title, ">", "&gt;")

		if item.URL != "" {
			sb.WriteString(fmt.Sprintf("%d. <a href=\"%s\">%s</a>\n", item.Ranks[0], item.URL, title))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s\n", item.Ranks[0], title))
		}
	}

	msg := TelegramMessage{
		ChatID:    n.chatID,
		Text:      sb.String(),
		ParseMode: "HTML",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
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
		return fmt.Errorf("telegram api status code: %d", resp.StatusCode)
	}

	return nil
}

