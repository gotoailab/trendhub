package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/gotoailab/trendhub/config"
	"github.com/gotoailab/trendhub/internal/model"
)

type Notifier interface {
	Send(ctx context.Context, items []*model.NewsItem) error
	Name() string
}

type NotificationManager struct {
	notifiers []Notifier
}

func NewNotificationManager(cfg *config.Config) *NotificationManager {
	manager := &NotificationManager{}

	if cfg.Notification.Webhooks.FeishuURL != "" {
		manager.notifiers = append(manager.notifiers, NewFeishuNotifier(cfg.Notification.Webhooks.FeishuURL))
	}

	if cfg.Notification.Webhooks.DingtalkURL != "" {
		manager.notifiers = append(manager.notifiers, NewDingtalkNotifier(cfg.Notification.Webhooks.DingtalkURL))
	}

	if cfg.Notification.Webhooks.TelegramBotToken != "" && cfg.Notification.Webhooks.TelegramChatID != "" {
		manager.notifiers = append(manager.notifiers, NewTelegramNotifier(cfg.Notification.Webhooks.TelegramBotToken, cfg.Notification.Webhooks.TelegramChatID))
	}

	if cfg.Notification.Webhooks.BarkDeviceKey != "" {
		manager.notifiers = append(manager.notifiers, NewBarkNotifier(cfg.Notification.Webhooks.BarkServerURL, cfg.Notification.Webhooks.BarkDeviceKey))
	}

	if cfg.Notification.Webhooks.WPSWebhookURL != "" {
		manager.notifiers = append(manager.notifiers, NewWPSNotifier(cfg.Notification.Webhooks.WPSWebhookURL))
	}

	return manager
}

func (nm *NotificationManager) SendAll(ctx context.Context, items []*model.NewsItem) {
	if len(items) == 0 {
		return
	}

	for _, n := range nm.notifiers {
		go func(notifier Notifier) {
			// 为通知发送创建独立的 context，避免因主 context 取消导致通知失败
			// 设置 30 秒超时，足够发送通知
			notifyCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := notifier.Send(notifyCtx, items); err != nil {
				fmt.Printf("Failed to send notification via %s: %v\n", notifier.Name(), err)
			} else {
				fmt.Printf("Notification sent via %s\n", notifier.Name())
			}
		}(n)
	}
}
