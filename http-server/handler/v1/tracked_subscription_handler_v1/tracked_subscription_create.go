package tracked_subscription_handler_v1

import (
	"net/http"

	"paylist.server/infra/constants"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/profanity"
	"paylist.server/types"

	"github.com/go-playground/validator"
)

func (h *Handler) CreateSubscriptionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *TrackedSubscriptionCreatePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	payload.Tariff = constants.NormalizeTariff(payload.Tariff)

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	if err := validateTrackedSubscriptionPayload(tr, authToken, payload.Price, payload.DatePay.Time, payload.AutoRenewal, payload.Notification); err != nil {
		return err
	}

	if err := profanity.Reject(ctx, tr, "subscription-name", payload.Name); err != nil {
		return err
	}

	payload.AutoRenewal, payload.Notification = normalizePremiumSubscriptionFields(authToken, payload.AutoRenewal, payload.Notification)

	count, err := h.Store.TrackedSubscriptions.Count_SubscriptionsByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if err := checkTrackedSubscriptionLimit(tr, authToken, count); err != nil {
		return err
	}

	sub := trackedSubscriptionFromCreatePayload(authToken.User.UserUUID, payload)

	if err := h.Store.TrackedSubscriptions.Create_Subscription(ctx, sub); err != nil {
		return httperr.Db(ctx, err)
	}

	if payload.Categories != nil {
		if err := h.Store.TrackedSubscriptions.Replace_SubscriptionCategories(ctx, sub.ID, authToken.User.UserUUID, payload.Categories); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	httpx.HttpResponse(w, r, http.StatusCreated, tr.T("success.tracked-subscription-created"))
	return nil
}
