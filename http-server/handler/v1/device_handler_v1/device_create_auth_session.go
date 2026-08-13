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

func (h *Handler) DeviceCreateAuthSessionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	// env := os.Getenv("GO_ENV") == "DEV"
	tr := middleware.TranslatorFromContext(ctx)

	var payload *DeviceCreateAuthSessionPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	session := &types.DeviceAuthSession{
		DeviceId:  payload.DeviceId,
		SessionId: "",
		Confirmed: false,
	}

	sessionId, err := redis.RedisDeviceAuthCreate(session)
	if err != nil {
		logger.Error("DeviceCreateAuthSessionHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusCreated, sessionId)
	return nil
}
