package device_handler_v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
)

/* paylist HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	deviceRouter := router.PathPrefix("/device").Subrouter()

	deviceProtectedRouter := router.PathPrefix("/device").Subrouter()
	deviceProtectedRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	deviceRouter.Handle("/create-session", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeviceCreateAuthSessionHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	deviceProtectedRouter.Handle("/session/{session_id}", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeviceSessionGetBySessionIdHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceProtectedRouter.Handle("/confirm", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeviceConfirmHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	deviceRouter.Handle("/status/{session_id}", middleware.PowDDos().Middleware(
		httpx.ErrorHandler(h.DeviceStatusConfirmHandler_V1),
	)).Methods(http.MethodGet)
}
