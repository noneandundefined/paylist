package tracked_subscription_handler_v1

import (
	"net/http"

	"paylist.server/middleware"
	"paylist.server/pkg/httpx"

	"github.com/gorilla/mux"
)

/* paylist HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	trackedSubscriptionRouter := router.PathPrefix("/tracked-subscription").Subrouter()
	trackedSubscriptionFreeRouter := router.PathPrefix("/tracked-subscription").Subrouter()

	trackedSubscriptionRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("", httpx.ErrorHandler(h.GetSubscriptionsHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("", httpx.ErrorHandler(h.CreateSubscriptionHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/summary", httpx.ErrorHandler(h.GetSubscriptionSummaryHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/analytics", httpx.ErrorHandler(h.GetSubscriptionAnalyticsHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/export", httpx.ErrorHandler(h.ExportSubscriptionsHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/categories", httpx.ErrorHandler(h.GetSubscriptionCategoriesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/services", httpx.ErrorHandler(h.GetServicesHandler_V1)).Methods(http.MethodGet)

	/* Access: Premium */
	trackedSubscriptionRouter.Handle("/categories", middleware.PowDDos()(
		httpx.ErrorHandler(h.CreateSubscriptionCategoryHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: Premium */
	trackedSubscriptionRouter.Handle("/categories/{categoryId:[0-9]+}", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeleteSubscriptionCategoryHandler_V1),
	)).Methods(http.MethodDelete)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/invites/accept", httpx.ErrorHandler(h.AcceptSubscriptionInviteHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}/members", httpx.ErrorHandler(h.GetSubscriptionMembersHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}/members", httpx.ErrorHandler(h.InviteSubscriptionMemberHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}/members/me", httpx.ErrorHandler(h.LeaveSubscriptionHandler_V1)).Methods(http.MethodDelete)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}/members/{memberId:[0-9]+}", httpx.ErrorHandler(h.DeleteSubscriptionMemberHandler_V1)).Methods(http.MethodDelete)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}/shares", httpx.ErrorHandler(h.ProposeSubscriptionSharesHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}/shares/{proposalId:[0-9]+}/vote", httpx.ErrorHandler(h.VoteSubscriptionSharesHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.GetSubscriptionByIdHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.EditSubscriptionHandler_V1)).Methods(http.MethodPut)

	/* Access: ALL */
	trackedSubscriptionRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.DeleteSubscriptionHandler_V1)).Methods(http.MethodDelete)

	/* Access: ALL */
	trackedSubscriptionFreeRouter.Handle("/invites", httpx.ErrorHandler(h.GetSubscriptionInviteHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	trackedSubscriptionFreeRouter.Handle("/images/w350", httpx.ErrorHandler(h.GetSubscriptionImageHandler_V1)).Methods(http.MethodGet)
}
