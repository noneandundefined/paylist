package tracked_subscription_handler_v1

import (
	"net/http"

	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) GetSubscriptionCategoriesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	categories, err := h.Store.TrackedSubscriptions.Get_SubscriptionCategoriesForUser(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, mapSubscriptionCategoriesResponse(categories))
	return nil
}
