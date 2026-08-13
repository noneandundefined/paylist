package currency_handler_v1

import (
	"net/http"

	"paylist.server/pkg/httpx"

	"github.com/gorilla/mux"
)

/* paylist HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	currencyRouter := router.PathPrefix("/currency").Subrouter()

	/* Access: ALL */
	currencyRouter.Handle("/convert", httpx.ErrorHandler(h.ConvertCurrencyHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	currencyRouter.Handle("/rates", httpx.ErrorHandler(h.GetCurrencyRatesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	currencyRouter.Handle("/currencies", httpx.ErrorHandler(h.GetCurrenciesHandler_V1)).Methods(http.MethodGet)
}
