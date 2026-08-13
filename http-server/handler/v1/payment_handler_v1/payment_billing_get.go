package payment_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) GetBillingHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	billing, err := h.Store.Payments.Get_UserSubscriptionBillingByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	state := models.BillingStateFromSubscription(billing)

	httpx.HttpResponseWithETag(w, r, http.StatusOK, state)
	return nil
}
