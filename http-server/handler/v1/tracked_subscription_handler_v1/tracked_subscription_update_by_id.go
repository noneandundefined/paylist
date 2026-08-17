package tracked_subscription_handler_v1

import (
	"net/http"
	"strconv"

	"paylist.server/infra/constants"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/profanity"
	"paylist.server/types"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func (h *Handler) EditSubscriptionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *TrackedSubscriptionEditPayload

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httperr.BadRequest(tr.TErr("error.tracked-subscription-invalid-id"))
	}

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

	existing, err := h.Store.TrackedSubscriptions.Get_SubscriptionById(ctx, uint64(id), authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if existing == nil {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	payload.AutoRenewal, payload.Notification = normalizePremiumSubscriptionFields(authToken, payload.AutoRenewal, payload.Notification)

	texts := make([]string, 0, 2)
	if existing.IsOwner {
		texts = append(texts, payload.Name)
	}
	if note := normalizeNote(payload.Note); note != nil {
		texts = append(texts, *note)
	}
	if err := profanity.Reject(ctx, tr, "subscription-update", texts...); err != nil {
		return err
	}

	if existing.IsOwner {
		if err := validateTrackedSubscriptionPayload(tr, authToken, payload.Price, payload.DatePay.Time, payload.AutoRenewal, payload.Notification); err != nil {
			return err
		}

		sub := trackedSubscriptionFromEditPayload(existing.UserUUID, payload)
		sub.ID = existing.ID

		if err := h.Store.TrackedSubscriptions.Update_SubscriptionById(ctx, sub, id); err != nil {
			return httperr.Db(ctx, err)
		}

		if !isSameDate(existing.DatePay, sub.DatePay) {
			previousDatePay := existing.DatePay
			if err := writeSubscriptionHistory(ctx, h.Store.TrackedSubscriptions, sub, historyEventDateChanged, &previousDatePay); err != nil {
				return httperr.Db(ctx, err)
			}
		}
	}

	if payload.Categories != nil {
		if err := h.Store.TrackedSubscriptions.Replace_SubscriptionCategories(ctx, uint64(id), authToken.User.UserUUID, payload.Categories); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	if err := h.Store.TrackedSubscriptions.Update_MemberPreferences(ctx, uint64(id), authToken.User.UserUUID, payload.Notification, includeInAnalyticsValue(payload.IncludeInAnalytics), normalizeNote(payload.Note)); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.tracked-subscription-updated"))
	return nil
}
