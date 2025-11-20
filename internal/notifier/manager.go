package notifier

import (
	"context"
	"fmt"

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
	
	return manager
}

func (nm *NotificationManager) SendAll(ctx context.Context, items []*model.NewsItem) {
	if len(items) == 0 {
		return
	}
	
	for _, n := range nm.notifiers {
		go func(notifier Notifier) {
			if err := notifier.Send(ctx, items); err != nil {
				fmt.Printf("Failed to send notification via %s: %v\n", notifier.Name(), err)
			} else {
				fmt.Printf("Notification sent via %s\n", notifier.Name())
			}
		}(n)
	}
}
