package user_handler_v1

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"paylist.server/middleware"
	"paylist.server/pkg/country"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/types"
)

func (h *Handler) UserSettingsGetHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	settings, err := h.Store.Users.Get_UserSettingsByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	response := UserSettingsResponse{
		DisplayCurrency:   settings.DisplayCurrency,
		Country:           settings.Country,
		TelegramConnected: settings.TelegramChatID != nil,
		TelegramUsername:  settings.TelegramUsername,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, response)
	return nil
}

func (h *Handler) UserSettingsUpdateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !premium.IsPremiumPlan(authToken) {
		return httperr.Forbidden(tr.TErr("error.premium-required"))
	}

	var payload *UserSettingsUpdatePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	if payload.DisplayCurrency == nil && payload.Country == nil {
		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	if payload.DisplayCurrency != nil {
		currencyCode := strings.ToUpper(strings.TrimSpace(*payload.DisplayCurrency))

		if err := h.Store.Users.Upsert_UserDisplayCurrency(ctx, authToken.User.UserUUID, currencyCode); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	if payload.Country != nil {
		countryCode := strings.ToUpper(strings.TrimSpace(*payload.Country))

		if _, found := country.FindCountry(countryCode); !found {
			return httperr.BadRequest(tr.TErr("error.invalid-country"))
		}

		if err := h.Store.Users.Upsert_UserCountry(ctx, authToken.User.UserUUID, countryCode); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("success.settings-updated"))
	return nil
}
