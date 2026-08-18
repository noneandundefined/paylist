package payment_handler_v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
)

/* paylist HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.Handle("/payments/webhook", httpx.ErrorHandler(h.PostWebhookHandler_V1)).Methods(http.MethodPost)

	paymentRouter := router.PathPrefix("/payments").Subrouter()
	paymentRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	paymentRouter.Handle("/checkout", middleware.PowDDos()(
		httpx.ErrorHandler(h.PostCheckoutHandler_V1),
	)).Methods(http.MethodPost)

	paymentRouter.Handle("/confirm", httpx.ErrorHandler(h.GetPaymentConfirmHandler_V1)).Methods(http.MethodGet)

	paymentRouter.Handle("/billing", httpx.ErrorHandler(h.GetBillingHandler_V1)).Methods(http.MethodGet)

	paymentRouter.Handle("/history", middleware.PowDDos()(
		httpx.ErrorHandler(h.GetPaymentHistoryHandler_V1),
	)).Methods(http.MethodGet)

	paymentRouter.Handle("/auto-renew", middleware.PowDDos()(
		httpx.ErrorHandler(h.PatchAutoRenewHandler_V1),
	)).Methods(http.MethodPatch)
}
