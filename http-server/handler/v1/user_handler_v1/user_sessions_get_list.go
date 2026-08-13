package user_handler_v1

import (
	"net/http"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) UserSessionsGetListHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	sessionIds, err := redis.RedisDeviceSessionIDs(authToken.User.UserUUID)
	if err != nil {
		logger.Error("UserSessionsGetListHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	sessions := make([]UserSessionResponse, 0, len(sessionIds))

	for _, sessionId := range sessionIds {
		session, err := redis.RedisSessionGet(sessionId)
		if err != nil {
			logger.Error("UserSessionsGetListHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
			continue
		}

		if session == nil {
			_ = redis.RedisDeviceSessionRemove(authToken.User.UserUUID, sessionId)
			continue
		}

		sessions = append(sessions, UserSessionResponse{
			SessionId: sessionId,
			Platform:  session.Platform,
			DeviceId:  session.DeviceId,
			CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
			Current:   sessionId == authToken.SessionId,
		})
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, sessions)
	return nil
}
