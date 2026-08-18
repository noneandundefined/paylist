package payment_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) PatchAutoRenewHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload AutoRenewPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	billing, err := h.Store.Payments.Get_UserSubscriptionBillingByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	state := models.BillingStateFromSubscription(billing)
	if payload.Enabled {
		return httperr.BadRequest(tr.TErr("error.payment-method-required"))
	}

	if err := h.Store.Payments.Update_UserSubscriptionAutoRenew(ctx, authToken.User.UserUUID, payload.Enabled); err != nil {
		return httperr.Db(ctx, err)
	}

	state.AutoRenewEnabled = payload.Enabled

	httpx.HttpResponseWithETag(w, r, http.StatusOK, state)
	return nil
}
