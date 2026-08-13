package device_handler_v1

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

func (h *Handler) DeviceConfirmHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *DeviceConfirmPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	deviceAuth, err := redis.RedisDeviceAuthGet(payload.SessionId)
	if err != nil {
		logger.Error("DeviceConfirmHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	session := &types.Session{
		UserUuid: authToken.User.UserUUID,
		Platform: "device",
		DeviceId: deviceAuth.DeviceId,
	}

	sessionId, err := redis.RedisSessionCreate(session)
	if err != nil {
		logger.Error("DeviceConfirmHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	/* Confirm tmp desktop auth */
	deviceAuth.Confirmed = true
	deviceAuth.SessionId = sessionId

	if err := redis.RedisDeviceAuthUpdate(payload.SessionId, deviceAuth); err != nil {
		logger.Error("DeviceConfirmHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusCreated, deviceAuth)
	return nil
}
