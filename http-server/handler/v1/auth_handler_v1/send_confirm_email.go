package auth_handler_v1

import (
	"context"
	"fmt"

	"paylist.server/infra/constants"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/util"
)

func (h *Handler) sendConfirmEmail(ctx context.Context, email string, userUUID string, reqID string) error {
	tr := middleware.TranslatorFromContext(ctx)
	email = util.NormalizeEmail(email)

	allowed, err := redis.RedisEmailConfirmCheckAndIncrement(email, constants.EmailConfirmSendLimit24h)
	if err != nil {
		logger.Error("sendConfirmEmail req={%s}: redis limit check failed: %s", reqID, err.Error())
		return httperr.Redis(ctx, err)
	}

	if !allowed {
		logger.Error("sendConfirmEmail req={%s}: exceeded daily limit (%s)", reqID, email)
		return httperr.Conflict(tr.TErr("error.email-limit-24-hours"))
	}

	if err := redis.RedisEmailConfirmPendingSet(email); err != nil {
		logger.Error("sendConfirmEmail req={%s}: redis pending set failed: %s", reqID, err.Error())
		return httperr.Redis(ctx, err)
	}

	link := util.GenerateVerificationEmailLink(userUUID)

	go func() {
		if errSendEmail := pkg.SendEmail(email, tr.T("confirm-email-title"), fmt.Sprintf(`
				<p>%s</p>

				<div style="border-left:4px solid #0085FF; padding:4px 0 4px 16px; margin:16px 0;">
					<a href="%s" style="color:#0085FF; text-decoration:underline; font-weight:bold;">%s</a>
				</div>

				<p>%s</p>
			`,
			tr.T("email-account-created"),
			link,
			tr.T("email-confirm-cta"),
			tr.T("exp-limit-reset-password"),
		), tr); errSendEmail != nil {
			logger.Error("sendConfirmEmail req={%s}: %s", reqID, errSendEmail.Error())
		}

		logger.Info("sendConfirmEmail req={%s}: Send link for %s: %s", reqID, email, link)
	}()

	return nil
}
