package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"paylist.server/pkg/httpx"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string, rps float64, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if l, exists := limiters[ip]; exists {
		return l
	}

	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	limiters[ip] = limiter
	return limiter
}

/* Ограничивает кол-во запросов от клиентов */
func RateLimiterMiddleware(rps float64, burst int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isWebhookPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			tr := TranslatorFromContext(r.Context())
			ip := r.RemoteAddr

			limiter := getLimiter(ip, rps, burst)

			if !limiter.Allow() {
				httpx.HttpResponse(w, r, http.StatusTooManyRequests, tr.TErr("rate-limiter-exceeded"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isWebhookPath(path string) bool {
	return strings.Contains(path, "/telegram/webhook/") || strings.Contains(path, "/max/webhook/")
}
