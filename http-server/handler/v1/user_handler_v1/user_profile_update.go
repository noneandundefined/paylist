package user_handler_v1

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/profanity"
	"paylist.server/types"
)

func (h *Handler) UserProfileUpdateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *UserProfileUpdatePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	valid := regexp.MustCompile(`^[\p{L}0-9_]+$`)

	if payload.FirstName != nil {
		trimmed := strings.TrimSpace(*payload.FirstName)
		payload.FirstName = &trimmed

		if !valid.MatchString(trimmed) {
			return httperr.BadRequest(tr.TErr("error.invalid-characters-firstname"))
		}
	}

	if payload.LastName != nil {
		trimmed := strings.TrimSpace(*payload.LastName)
		payload.LastName = &trimmed

		if !valid.MatchString(trimmed) {
			return httperr.BadRequest(tr.TErr("error.invalid-characters-lastname"))
		}
	}

	if err := profanity.Reject(ctx, tr, "profile-name", profanity.Pointers(payload.FirstName, payload.LastName)...); err != nil {
		return err
	}

	if err := h.Store.Users.Update_UserProfile(ctx, authToken.User.UserUUID, payload.FirstName, payload.LastName); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("profile-updated"))
	return nil
}
