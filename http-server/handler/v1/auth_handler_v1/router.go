package auth_handler_v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
)

/* paylist HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	authProtectedRouter := router.PathPrefix("/auth").Subrouter()
	authProtectedRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	authRouter.Handle("/signin", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthSigninHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/signup", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthSignupHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/check", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthCheckHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/confirm", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthEmailConfirmHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	authRouter.Handle("/confirm/pending", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthConfirmPendingHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	authProtectedRouter.Handle("/confirm/req", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthEmailConfirmReqHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/password/reset/request", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthPasswordResetRequestHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/password/reset/confirm", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthPasswordResetConfirmHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authProtectedRouter.Handle("/signout", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthSignoutHandler_V1),
	)).Methods(http.MethodPost)
}
