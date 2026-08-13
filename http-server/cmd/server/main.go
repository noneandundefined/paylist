package main

import (
	"database/sql"
	"os"

	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/infra/store/redis"
	"paylist.server/pkg/telegram"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

type httpServer struct {
	db       *sql.DB
	cron     *cron.Cron
	store    store.Storage
	telegram *telegram.Notifier
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
		}
	} else if os.Getenv("TELEGRAM_BOT_TOKEN") != "" {
		logger.Error("Telegram bot is configured but failed to initialize: %s", err.Error())
	}

	server.cron = cron.New(cron.WithSeconds())

	server.startCronJobs()

	server.cron.Start()

	/* Started HTTPx server */
	if err := server.httpStart(); err != nil {
		logger.Error("Failed start TCP server: %s", err.Error())
	}
}
