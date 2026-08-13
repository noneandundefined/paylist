package user_handler_v1

import (
	"net/http"

	"github.com/go-playground/validator"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) UserSessionsDisconnectHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *UserSessionDisconnectPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	session, err := redis.RedisSessionGet(payload.SessionId)
	if err != nil {
		logger.Error("UserSessionsDisconnectHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	if session == nil || session.UserUuid != authToken.User.UserUUID {
		return httperr.NotFound(tr.TErr("error.device-not-found"))
	}

	if err := redis.RedisSessionDelete(payload.SessionId); err != nil {
		logger.Error("UserSessionsDisconnectHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("success.device-disconnected"))
	return nil
}
