package main

import (
	"context"
	"os"
	"strings"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/pkg/maxbot"
)

func startMaxPolling(notifier *maxbot.Notifier, client *maxbot.Client) {
	go func() {
		if err := client.DeleteSubscriptions(); err != nil {
			logger.Error("MAX deleteSubscriptions failed: %s", err.Error())
		}

		var marker *int64
		logger.Info("MAX long polling started")

		for {
			updates, nextMarker, err := client.GetUpdates(marker)
			if err != nil {
				logger.Error("MAX polling error: %s", err.Error())
				time.Sleep(3 * time.Second)
				continue
			}

			if nextMarker != nil {
				marker = nextMarker
			}

			for _, update := range updates {
				ctx := context.WithValue(context.Background(), "XREQID", "max-poll")
				if err := notifier.HandleUpdate(ctx, update); err != nil {
					logger.Error("MAX update handling error: %s", err.Error())
				}
			}
		}
	}()
}

func maxPollingEnabled() bool {
	if strings.TrimSpace(os.Getenv("MAX_POLLING")) == "false" {
		return false
	}

	if strings.TrimSpace(os.Getenv("MAX_WEBHOOK_SECRET")) != "" {
		return strings.TrimSpace(os.Getenv("MAX_POLLING")) == "true"
	}

	return true
}

func maxWebhookURL(secret string) string {
	if webhookURL := strings.TrimSpace(os.Getenv("MAX_WEBHOOK_URL")); webhookURL != "" {
		return webhookURL
	}

	base := strings.TrimRight(strings.TrimSpace(os.Getenv("SERVER_BASE_ADDR")), "/")
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return ""
	}

	return base + "/max/webhook/" + secret
}
