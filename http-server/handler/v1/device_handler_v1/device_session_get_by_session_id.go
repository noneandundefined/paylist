package device_handler_v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
)

func (h *Handler) DeviceSessionGetBySessionIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	sessionId := mux.Vars(r)["session_id"]
	if sessionId == "" {
		httpx.HttpResponseWithETag(w, r, http.StatusNotFound, map[string]any{"valid": false, "message": tr.TErr("error.session-not-found")})
		return nil
	}

	deviceAuth, err := redis.RedisDeviceAuthGet(sessionId)
	if err != nil {
		logger.Error("DeviceSessionGetBySessionIdHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	if deviceAuth == nil {
		httpx.HttpResponseWithETag(w, r, http.StatusNotFound, map[string]any{"valid": false, "message": tr.TErr("error.session-not-found")})
		return nil
	}

	if deviceAuth.Confirmed || deviceAuth.SessionId != "" {
		httpx.HttpResponseWithETag(w, r, http.StatusConflict, map[string]any{"valid": false, "message": tr.TErr("error.session-already-used")})
		return nil
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]any{"valid": true, "message": deviceAuth})
	return nil
}
