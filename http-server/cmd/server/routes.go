package main

import (
	"net/http"
	"os"
	"strings"

	"paylist.server/handler"
	"paylist.server/handler/v1/auth_handler_v1"
	"paylist.server/handler/v1/country_handler_v1"
	"paylist.server/handler/v1/currency_handler_v1"
	"paylist.server/handler/v1/device_handler_v1"
	"paylist.server/handler/v1/max_handler_v1"
	"paylist.server/handler/v1/payment_handler_v1"
	"paylist.server/handler/v1/plan_handler_v1"
	"paylist.server/handler/v1/telegram_handler_v1"
	"paylist.server/handler/v1/tracked_subscription_handler_v1"
	"paylist.server/handler/v1/user_handler_v1"
	"paylist.server/middleware"

	"github.com/gorilla/mux"
)

func (s *httpServer) routes() http.Handler {
	router := mux.NewRouter()

	/* Middleware for logging API request */
	router.Use(middleware.NewLogger().LoggerMiddleware)
	/* Middleware for X-Request-Id */
	router.Use(middleware.XRequestIdMiddleware())
	/* Middleware for i18n language */
	router.Use(middleware.LanguageMiddleware())
	/* Middleware for get exception errors */
	router.Use(middleware.RecoveryMiddleware())
	/* Middleware for security API */
	router.Use(middleware.SecurityMiddleware())
	/* Middleware rate limiter */
	router.Use(middleware.RateLimiterMiddleware(6, 10))

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	baseHandler := &handler.BaseHandler{
		Db:    s.db,
		Store: s.store,
	}

	/* Authenticate rotues */
	auth_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Device rotues */
	device_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* User rotues */
	user_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Payment rotues */
	payment_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Plan rotues */
	plan_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Currency rotues */
	currency_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Country rotues */
	country_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Tracked_subscriptions rotues */
	tracked_subscription_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)

	if s.telegram != nil {
		webhookSecret := strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_SECRET"))
		if webhookSecret != "" {
			subrouter.HandleFunc("/telegram/webhook/"+webhookSecret, telegram_handler_v1.WebhookHandler(s.telegram)).Methods(http.MethodPost)
		}
	}

	if s.maxbot != nil {
		webhookSecret := strings.TrimSpace(os.Getenv("MAX_WEBHOOK_SECRET"))
		if webhookSecret != "" {
			subrouter.HandleFunc("/max/webhook/"+webhookSecret, max_handler_v1.WebhookHandler(s.maxbot)).Methods(http.MethodPost)
		}
	}

	// docs
	s.docs(subrouter)

	return s.cors(router)
}
