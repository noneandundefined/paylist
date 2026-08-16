package auth_handler_v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"paylist.server/infra/constants"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/security"
	"paylist.server/util"
)

func (h *Handler) AuthPasswordResetRequestHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	reqID := ctx.Value("XREQID").(string)

	var payload *AuthPasswordResetRequestPayload

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

	allowed, err := redis.RedisPasswordResetCheckAndIncrement(email, constants.PasswordResetSendLimit24h)
	if err != nil {
		logger.Error("AuthPasswordResetRequestHandler_V1 req={%s}: redis limit check failed: %s", reqID, err.Error())
		return httperr.Redis(ctx, err)
	}

	if !allowed {
		return httperr.Conflict(tr.TErr("error.password-reset-limit-24-hours"))
	}

	user, err := h.Store.Users.Get_UserCoreByEmail(ctx, email)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user != nil {
		if err := h.sendPasswordResetEmail(ctx, email, user.UserUUID, reqID); err != nil {
			return err
		}
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("password-reset-link-send-to-email"))
	return nil
}

func (h *Handler) AuthPasswordResetConfirmHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	reqID := ctx.Value("XREQID").(string)

	uuid := r.URL.Query().Get("uuid")
	exp := r.URL.Query().Get("exp")
	sig := r.URL.Query().Get("sig")

	if !util.ValidatePasswordResetLink(uuid, exp, sig) {
		return httperr.BadRequest(tr.TErr("error.password-reset-invalid"))
	}

	matched, err := redis.RedisPasswordResetTokenMatch(uuid, sig)
	if err != nil {
		logger.Error("AuthPasswordResetConfirmHandler_V1 req={%s}: redis token check failed: %s", reqID, err.Error())
		return httperr.Redis(ctx, err)
	}

	if !matched {
		return httperr.BadRequest(tr.TErr("error.password-reset-invalid"))
	}

	var payload *AuthPasswordResetConfirmPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	password := strings.TrimSpace(payload.Password)

	if _, chPass := constants.CheckSimplePasswords[strings.ToLower(password)]; chPass {
		return httperr.BadRequest(tr.TErr("error.simple-password"))
	}

	user, err := h.Store.Users.Get_UserCoreByUserUuid(ctx, uuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.BadRequest(tr.TErr("error.password-reset-invalid"))
	}

	if security.PasswordEqualsEmail(password, user.Email) {
		return httperr.BadRequest(tr.TErr("error.password-equals-email"))
	}

	passwordHashed, err := security.HashPassword(password)
	if err != nil {
		logger.Error("AuthPasswordResetConfirmHandler_V1 req={%s}: Failed hash password: %s", reqID, err.Error())
		return httperr.InternalServerError(err.Error())
	}

	if err := h.Store.Users.Update_UserPassword(ctx, uuid, passwordHashed); err != nil {
		return httperr.Db(ctx, err)
	}

	_ = redis.RedisPasswordResetTokenClear(uuid)
	_ = redis.RedisDeviceSessionsDeleteAll(uuid)

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("success-password-reset"))
	return nil
}

func (h *Handler) sendPasswordResetEmail(ctx context.Context, email, userUUID, reqID string) error {
	tr := middleware.TranslatorFromContext(ctx)

	link, signature := util.GeneratePasswordResetLink(userUUID)

	if err := redis.RedisPasswordResetTokenSet(userUUID, signature); err != nil {
		logger.Error("sendPasswordResetEmail req={%s}: redis token set failed: %s", reqID, err.Error())
		return httperr.Redis(ctx, err)
	}

	go func() {
		if errSendEmail := pkg.SendEmail(email, tr.T("password-reset-data-subject"), fmt.Sprintf(`
				<p>%s</p>

				<div style="border-left:4px solid #0085FF; padding:4px 0 4px 16px; margin:16px 0;">
					<a href="%s" style="color:#0085FF; text-decoration:underline; font-weight:bold;">%s</a>
				</div>

				<p>%s</p>
			`,
			tr.T("go-to-link-reset-password"),
			link,
			tr.T("email-reset-password-cta"),
			tr.T("exp-limit-reset-password"),
		), tr); errSendEmail != nil {
			logger.Error("sendPasswordResetEmail req={%s}: %s", reqID, errSendEmail.Error())
		}

		logger.Info("sendPasswordResetEmail req={%s}: Send link for %s: %s", reqID, email, link)
	}()

	return nil
}
