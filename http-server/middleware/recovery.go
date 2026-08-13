package middleware

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"paylist.server/infra/logger"
	"paylist.server/pkg/httpx"
)

/* Перехват паник (panic) во время обработки HTTP-запросов */
func RecoveryMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := TranslatorFromContext(r.Context())

			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("%v", rec)
					httpx.HttpResponseError(w, r, errors.New(tr.TErr("error.oops-something-went-wrong")))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
