package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/pkg/telegram"
)

func startTelegramWebhook(client *telegram.Client, secret string) {
	webhookURL := telegramWebhookURL(secret)
	if webhookURL == "" {
		logger.Error("Telegram webhook URL is empty; set TELEGRAM_WEBHOOK_URL or SERVER_BASE_ADDR")
		return
	}

	if strings.HasPrefix(webhookURL, "http://") || strings.Contains(webhookURL, "localhost") {
		logger.Info("Telegram webhook skipped: URL must be public https")
		return
	}

	if err := client.SetWebhook(webhookURL); err != nil {
		logger.Error("Telegram setWebhook failed: %s", err.Error())
		return
	}

	logger.Info("Telegram webhook enabled")
}

func telegramWebhookURL(secret string) string {
	if webhookURL := strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_URL")); webhookURL != "" {
		return webhookURL
	}

	base := strings.TrimRight(strings.TrimSpace(os.Getenv("SERVER_BASE_ADDR")), "/")
	if base == "" {
		return ""
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return ""
	}

	if !strings.Contains(base, "/api/v1") {
		base += "/api/v1"
	}

	return base + "/telegram/webhook/" + secret
}

func startTelegramPolling(notifier *telegram.Notifier, client *telegram.Client) {
	go func() {
		if err := client.DeleteWebhook(); err != nil {
			logger.Error("Telegram deleteWebhook failed: %s", err.Error())
		}

		offset := 0
		logger.Info("Telegram long polling started")

		for {
			updates, err := client.GetUpdates(offset)
			if err != nil {
				var apiErr *telegram.APIError
				if errors.As(err, &apiErr) && apiErr.Code == 409 {
					logger.Error(
						"Telegram polling conflict (409): another process is calling getUpdates. Stop duplicate server instances or set TELEGRAM_POLLING=false",
					)
					time.Sleep(15 * time.Second)
					continue
				}

				logger.Error("Telegram polling error: %s", err.Error())
				time.Sleep(3 * time.Second)
				continue
			}

			for _, update := range updates {
				ctx := context.WithValue(context.Background(), "XREQID", "telegram-poll")
				if err := notifier.HandleUpdate(ctx, update); err != nil {
					if telegram.IsDeliveryRejected(err) {
						logger.Warning("Telegram skip blocked chat: %s", err.Error())
					} else {
						logger.Error("Telegram update handling error: %s", err.Error())
					}
				}

				offset = update.UpdateID + 1
			}
		}
	}()
}

func telegramPollingEnabled() bool {
	if strings.TrimSpace(os.Getenv("TELEGRAM_POLLING")) == "false" {
		return false
	}

	// Webhook mode: polling is opt-in to avoid two receivers fighting for updates.
	if strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_SECRET")) != "" {
		if isDevEnv() {
			return false
		}

		return strings.TrimSpace(os.Getenv("TELEGRAM_POLLING")) == "true"
	}

	return true
}

func isDevEnv() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("GO_ENV")), "DEV")
}
