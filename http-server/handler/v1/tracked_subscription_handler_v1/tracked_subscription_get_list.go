package tracked_subscription_handler_v1

import (
	"net/http"

	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) GetSubscriptionsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	search := r.URL.Query().Get("search")

	subs, err := h.Store.TrackedSubscriptions.Get_SubscriptionsByUuid(ctx, authToken.User.UserUUID, search)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	categoryMap, err := h.Store.TrackedSubscriptions.Get_CategorySlugsMapByUserUUID(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	response := make([]TrackedSubscriptionDetailResponse, 0, len(*subs))
	for _, sub := range *subs {
		categories := categoryMap[sub.ID]
		if categories == nil {
			categories = []string{}
		}

		response = append(response, TrackedSubscriptionDetailResponse{
			TrackedSubscription: sub,
			Categories:          categories,
		})
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, response)
	return nil
}
