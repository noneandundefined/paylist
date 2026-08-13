package auth_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/redis"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/util"
)

func (h *Handler) AuthConfirmPendingHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	email := util.NormalizeEmail(r.URL.Query().Get("email"))
	if email == "" {
		return httperr.BadRequest("email is required")
	}

	pending, err := redis.RedisEmailConfirmPendingCheck(email)
	if err != nil {
		return httperr.Redis(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]bool{"pending": pending})
	return nil
}
