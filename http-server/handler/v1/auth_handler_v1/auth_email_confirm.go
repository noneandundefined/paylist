package auth_handler_v1

import (
	"net/http"

	"github.com/go-playground/validator"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
	"paylist.server/util"
)

func (h *Handler) AuthEmailConfirmHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	uuid := r.URL.Query().Get("uuid")
	exp := r.URL.Query().Get("exp")
	sig := r.URL.Query().Get("sig")

	if !util.ValidateEmailConfirmLink(uuid, exp, sig) {
		httpx.HttpResponseWithETag(w, r, http.StatusAccepted, map[string]string{"status": "invalid", "message": tr.TErr("error.mail-confirmation")})
		return nil
	}

	if err := h.Store.Users.Update_UserEmailConfirmedByUid(ctx, uuid, true); err != nil {
		httpx.HttpResponseWithETag(w, r, http.StatusAccepted, map[string]string{"status": "error", "message": tr.TErr("error.mail-confirmation")})
		return nil
	}

	if user, err := h.Store.Users.Get_UserCoreByUserUuid(ctx, uuid); err == nil && user != nil {
		_ = redis.RedisEmailConfirmPendingClear(util.NormalizeEmail(user.Email))
	}

	session := &types.Session{
		UserUuid: uuid,
		Platform: "web",
	}

	sessionId, err := redis.RedisSessionCreate(session)
	if err != nil {
		logger.Error("AuthEmailConfirmHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		httpx.HttpResponseWithETag(w, r, http.StatusAccepted, map[string]string{"status": "error", "message": tr.TErr("error.mail-confirmation")})
		return nil
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "success", "message": sessionId})
	return nil
}

func (h *Handler) AuthEmailConfirmReqHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *ReqEmailConfirmPayload

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

	/* Send confirm to email */
	if err := h.sendConfirmEmail(ctx, email, authToken.User.UserUUID, ctx.Value("XREQID").(string)); err != nil {
		return err
	}

	httpx.HttpResponse(w, r, http.StatusAccepted, tr.T("confirm-link-send-to-email"))
	return nil
}
