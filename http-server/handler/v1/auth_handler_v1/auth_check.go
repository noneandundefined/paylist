package auth_handler_v1

import (
	"net/http"

	"github.com/go-playground/validator"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/util"
)

/* Neosync HTTPx V1 */
/* Handler: определение вход/регистрация для пользователя */

func (h *Handler) AuthCheckHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	var payload *AuthCheckPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	email := util.NormalizeEmail(payload.Email)

	user, err := h.Store.Users.Get_UserCoreByEmail(ctx, email)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "signup", "message": tr.TErr("error.user-not-found")})
		return nil
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "signin", "message": ""})
	return nil
}
