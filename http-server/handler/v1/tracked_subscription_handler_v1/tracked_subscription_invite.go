package tracked_subscription_handler_v1

import (
	"net/http"
	"time"

	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
	"paylist.server/util"

	"github.com/go-playground/validator"
)

func (h *Handler) GetSubscriptionInviteHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	token := r.URL.Query().Get("token")
	if token == "" {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-invite-invalid"))
	}

	preview, member, err := h.Store.TrackedSubscriptions.Get_InviteByToken(ctx, token)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if preview == nil || member == nil || preview.Status != "pending" {
		return httperr.NotFound(tr.TErr("error.shared-subscription-invite-invalid"))
	}

	if preview.InviteExpiresAt != nil && preview.InviteExpiresAt.Before(time.Now().UTC()) {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-invite-expired"))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, preview)
	return nil
}

func (h *Handler) AcceptSubscriptionInviteHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload TrackedSubscriptionInviteAcceptPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	preview, member, err := h.Store.TrackedSubscriptions.Get_InviteByToken(ctx, payload.Token)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if preview == nil || member == nil || preview.Status != "pending" {
		return httperr.NotFound(tr.TErr("error.shared-subscription-invite-invalid"))
	}

	if preview.InviteExpiresAt != nil && preview.InviteExpiresAt.Before(time.Now().UTC()) {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-invite-expired"))
	}

	if util.NormalizeEmail(preview.Email) != util.NormalizeEmail(authToken.User.Email) {
		return httperr.Forbidden(tr.TErr("error.shared-subscription-invite-email-mismatch"))
	}

	if err := h.Store.TrackedSubscriptions.Accept_MemberInvite(ctx, member.ID, authToken.User.UserUUID); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, TrackedSubscriptionInviteAcceptResponse{
		Message:        tr.T("success.shared-subscription-accepted"),
		SubscriptionID: preview.SubscriptionID,
	})
	return nil
}
