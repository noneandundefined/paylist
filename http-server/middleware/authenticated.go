// nolint
package middleware

import (
	"context"
	"net/http"

	"strings"

	"github.com/gorilla/mux"
	"paylist.server/handler"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/redis"
	"paylist.server/pkg/httpx"
	"paylist.server/types"
)

/* Проверка на аутентифицированного пользователя */
func IsAuthenticatedMiddleware(h *handler.BaseHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := TranslatorFromContext(r.Context())

			var sessionId string

			/* From header */
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					sessionId = parts[1]
				}
			}

			/* From cookie */
			if sessionId == "" {
				cookie, err := r.Cookie("auth-token")
				if err == nil && cookie != nil {
					sessionId = cookie.Value
				}
			}

			if sessionId == "" {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("error.unauthorized"))
				return
			}

			session, err := redis.RedisSessionGet(sessionId)
			if err != nil {
				logger.Error("IsAuthenticatedMiddleware: %v", err.Error())

				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("error.check-token"))
				return
			}

			if session == nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("error.unauthorized"))
				return
			}

			/* Get user core */
			user, err := h.Store.Users.Get_UserCoreByUserUuid(r.Context(), session.UserUuid)
			if err != nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("error.unauthorized"))
				return
			}

			if user == nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("error.unauthorized"))
				return
			}

			permissions, err := h.Store.Users.Get_UserPermissionsByUserUuid(r.Context(), session.UserUuid)
			if err != nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("error.unauthorized"))
				return
			}

			authTokenModel := &types.AuthToken{
				User:                      *user,
				SessionId:                 sessionId,
				PlanName:                  permissions.PlanName,
				NotificationSubscriptions: permissions.NotificationSubscriptions,
				MaxTotalSubscriptions:     permissions.MaxTotalSubscriptions,
				AutoFindSubscriptions:     permissions.AutoFindSubscriptions,
				IsAdmin:                   user.IsAdmin,
			}

			// session REFRESH
			redis.RedisSessionRefresh(sessionId)

			//nolint
			ctx := context.WithValue(r.Context(), "identity", authTokenModel)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
