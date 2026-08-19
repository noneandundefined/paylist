package main

import (
	"database/sql"
	"os"
	"strings"
	"time"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/infra/store/redis"
	"paylist.server/pkg/maxbot"
	"paylist.server/pkg/telegram"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

type httpServer struct {
	db       *sql.DB
	cron     *cron.Cron
	store    store.Storage
	telegram *telegram.Notifier
	maxbot   *maxbot.Notifier
}

func main() {
	/* local: .env.docker | deploy replaces via scripts/pre_deploy */
	_ = godotenv.Load(".env.docker")

	/* Inital logger */
	logger.InitLogger()

	/* Localization en/es/ru */
	locale.InitI18n()

	/* Initial connect db */
	db, err := postgres.New(os.Getenv("DB_ADDR"), 150, 25, "7m")
	if err != nil {
		logger.Error("Failed connect to database: %s", err.Error())
		return
	}
	defer db.Close()

	/* Store for postgres */
	store := store.NewStorage(db)

	/* Initial connect redis */
	if err := redis.NewRedisDb(); err != nil {
		logger.Error("Failed connect to redis: %s", err.Error())
		return
	}

	server := &httpServer{
		db:    db,
		store: store,
	}

	if tgClient, err := telegram.NewFromEnv(); err == nil {
		server.telegram = telegram.NewNotifier(store, tgClient)

		if telegramPollingEnabled() {
			startTelegramPolling(server.telegram, tgClient)
		} else if webhookSecret := strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_SECRET")); webhookSecret != "" {
			startTelegramWebhook(tgClient, webhookSecret)
		}
	} else if os.Getenv("TELEGRAM_BOT_TOKEN") != "" {
		logger.Error("Telegram bot is configured but failed to initialize: %s", err.Error())
	}

	if maxClient, err := maxbot.NewFromEnv(); err == nil {
		server.maxbot = maxbot.NewNotifier(store, maxClient)

		if maxPollingEnabled() {
			startMaxPolling(server.maxbot, maxClient)
		} else if webhookSecret := strings.TrimSpace(os.Getenv("MAX_WEBHOOK_SECRET")); webhookSecret != "" {
			if webhookURL := maxWebhookURL(webhookSecret); webhookURL != "" {
				if err := maxClient.SubscribeWebhook(webhookURL, webhookSecret); err != nil {
					logger.Error("MAX subscribe webhook failed: %s", err.Error())
				} else {
					logger.Info("MAX webhook subscribed")
				}
			}
		}
	} else if os.Getenv("MAX_BOT_TOKEN") != "" {
		logger.Error("MAX bot is configured but failed to initialize: %s", err.Error())
	}

	server.cron = cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))

	server.startCronJobs()

	server.cron.Start()
	server.runCronJob(cronPremiumPlanPrice, server.refreshPremiumPlanPrice)

	/* Started HTTPx server */
	if err := server.httpStart(); err != nil {
		logger.Error("Failed start TCP server: %s", err.Error())
	}
}
