package auth_handler_v1

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/security"
	"paylist.server/types"
	"paylist.server/util"
)

/* Neosync HTTPx V1 */
/* Handler: вход в аккаунт пользователя */

func (h *Handler) AuthSigninHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	// env := os.Getenv("GO_ENV") == "DEV"
	tr := middleware.TranslatorFromContext(ctx)

	var payload *AuthSigninPayload

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
		httpx.HttpResponse(w, r, http.StatusOK, map[string]string{"status": "password", "message": tr.TErr("error.user-not-found")})
		return nil
	}

	if !security.CheckPassword(payload.Password, user.Password) {
		time.Sleep(500 * time.Millisecond)
		return httperr.NotFound(tr.TErr("error.user-not-found"))
	}

	if !user.EmailConfirmed {
		/* Send confirm to email */
		if err := h.sendConfirmEmail(ctx, email, user.UserUUID, ctx.Value("XREQID").(string)); err != nil {
			return err
		}

		httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{
			"status":  "sent",
			"message": tr.T("confirm-link-send-to-email"),
		})
		return nil
	}

	session := &types.Session{
		UserUuid: user.UserUUID,
		Platform: "web",
	}

	sessionId, err := redis.RedisSessionCreate(session)
	if err != nil {
		logger.Error("AuthSigninHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{
		"status":  "signed_in",
		"message": sessionId,
	})
	return nil
}
