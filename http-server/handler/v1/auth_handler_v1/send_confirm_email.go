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
				<body>
					<p>%s <b>%s</b></p>
					<p>%s:</p>

					<a href="%s">
						%s
					</a>

					<p>%s</p>
					<p>P.S. %s</p>
				</body>
			`,
			tr.T("request-email-mail"), email,
			tr.T("go-to-link-confirm-email"),
			link,
			link,
			tr.T("exp-limit-reset-password"),
			tr.T("email-created-automatic"),
		), tr); errSendEmail != nil {
			logger.Error("sendConfirmEmail req={%s}: %s", reqID, errSendEmail.Error())
		}

		logger.Info("sendConfirmEmail req={%s}: Send link for %s: %s", reqID, email, link)
	}()

	return nil
}
