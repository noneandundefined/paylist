package user_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/pkg/telegram"
	"paylist.server/types"
)

type TelegramLinkResponse struct {
	BotURL string `json:"bot_url"`
}

func (h *Handler) UserTelegramLinkHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !premium.IsPremiumPlan(authToken) {
		return httperr.Forbidden(tr.TErr("error.premium-required"))
	}

	permissions, err := h.Store.Users.Get_UserPermissionsByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if !permissions.NotificationSubscriptions {
		return httperr.Forbidden(tr.TErr("error.premium-required"))
	}

	client, err := telegram.NewFromEnv()
	if err != nil {
		return httperr.InternalServerError(tr.TErr("error.server-error"))
	}

	token, err := redis.RedisTelegramLinkCreate(authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, TelegramLinkResponse{
		BotURL: client.BotURL(token),
	})

	return nil
}

func (h *Handler) UserTelegramDisconnectHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if err := h.Store.Users.Clear_UserTelegram(ctx, authToken.User.UserUUID); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("success.telegram-disconnected"))
	return nil
}
