package plan_handler_v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	planRouter := router.PathPrefix("/plans").Subrouter()

	planRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access ALL */
	planRouter.Handle("", httpx.ErrorHandler(h.GetPlansHandler_V1)).Methods(http.MethodGet)
}
