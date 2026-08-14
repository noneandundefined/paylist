package middleware

import (
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"paylist.server/infra/logger"
	"paylist.server/pkg/httpx"
)

const (
	MaxBodySize          = 1 << 20
	MaxMultipartBodySize = 3 << 20
	MaxViolationsPerIP   = 3
	IpBlockDuration      = 10 * time.Minute
)

var (
	BlockedIPs   = make(map[string]time.Time)
	IpViolations = make(map[string]int)
	Mutex        sync.Mutex
)

/* Проверока данных от пользователей (XSS, SQLInjection, ...) */
func SecurityMiddleware() mux.MiddlewareFunc { //nolint
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := TranslatorFromContext(r.Context())
			ip := getIP(r)

			// logger.Ip(ip)

			if isBlocked(ip) {
				logger.Warning("SecurityMiddleware ip={%s}: Ip has been blocked", ip)
				httpx.HttpResponse(w, r, http.StatusForbidden, tr.TErr("error.access-denied"))
				return
			}

			/* Limit size */
			maxBody := int64(MaxBodySize)
			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				maxBody = MaxMultipartBodySize
				r.Body = http.MaxBytesReader(w, r.Body, maxBody)
				next.ServeHTTP(w, r)
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, maxBody)

			var bodyContent string
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				data, err := io.ReadAll(r.Body)
				if err != nil {
					httpx.HttpResponse(w, r, http.StatusRequestEntityTooLarge, tr.TErr("error.request-text-toobig"))
					registerViolation(ip)
					return
				}

				bodyContent = string(data)
				r.Body = io.NopCloser(strings.NewReader(bodyContent))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIP(r *http.Request) string {
	/* CF-Connecting-IP (Cloudflare) */
	if cf := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); cf != "" {
		if ip := net.ParseIP(cf); ip != nil {
			return cf
		}
	}

	/* X-Forwarded-For */
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")

		for _, p := range parts {
			ipStr := strings.TrimSpace(p)
			if ip := net.ParseIP(ipStr); ip != nil && !isInternalIP(ip) {
				return ipStr
			}
		}

		first := strings.TrimSpace(parts[0])
		if net.ParseIP(first) != nil {
			return first
		}
	}

	/* X-Real-IP (nginx) */
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		if net.ParseIP(xr) != nil {
			return xr
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func isInternalIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	privateCIDRs := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateCIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func isBlocked(ip string) bool {
	if isInternalIP(net.ParseIP(ip)) {
		return false
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	expiry, exists := BlockedIPs[ip]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(BlockedIPs, ip)
		delete(IpViolations, ip)

		return false
	}

	return true
}

func registerViolation(ip string) {
	Mutex.Lock()
	defer Mutex.Unlock()

	IpViolations[ip]++
	if IpViolations[ip] >= MaxViolationsPerIP {
		BlockedIPs[ip] = time.Now().Add(IpBlockDuration)
	}
}
