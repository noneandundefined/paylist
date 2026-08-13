package user_handler_v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
)

/* paylist HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	userRouter.Handle("/settings", httpx.ErrorHandler(h.UserSettingsGetHandler_V1)).Methods(http.MethodGet)

	/* Access: Premium */
	userRouter.Handle("/telegram/link", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserTelegramLinkHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	userRouter.Handle("/telegram", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserTelegramDisconnectHandler_V1),
	)).Methods(http.MethodDelete)

	/* Access: Premium */
	userRouter.Handle("/settings", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserSettingsUpdateHandler_V1),
	)).Methods(http.MethodPatch)

	/* Access: ALL */
	userRouter.Handle("/login-state", httpx.ErrorHandler(h.UserLoginStateHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/sessions", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserSessionsGetListHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/sessions/disconnect", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserSessionsDisconnectHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	userRouter.Handle("/me", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserProfileUpdateHandler_V1),
	)).Methods(http.MethodPatch)

	/* Access: ALL */
	userRouter.Handle("/me", middleware.PowDDos()(
		httpx.ErrorHandler(h.UserAccountDeleteHandler_V1),
	)).Methods(http.MethodDelete)
}
