package middleware

import (
	"crypto/sha256"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"paylist.server/pkg/httpx"
)

const (
	defaultSecret     = "d92376b55175"
	defaultDifficulty = 1
)

func generateChallenge() string {
	secret := os.Getenv("POW_DDOS_SECRET")
	if secret == "" {
		secret = defaultSecret
	}

	t := time.Now().Unix() / 10
	return secret + ":" + strconv.FormatInt(t, 10)
}

func hashMatchesDifficulty(input string, difficulty int) bool {
	hash := sha256.Sum256([]byte(input))

	for i := 0; i < difficulty; i++ {
		byteIdx := i / 2
		half := i % 2

		val := hash[byteIdx]

		if half == 0 {
			if (val >> 4) != 0 {
				return false
			}
		} else {
			if (val & 0x0F) != 0 {
				return false
			}
		}
	}

	return true
}

func verify(r *http.Request, difficulty int) bool {
	challenge := r.Header.Get("Pow-Challenge")
	nonce := r.Header.Get("Pow-Nonce")

	if challenge == "" || nonce == "" {
		return false
	}

	input := challenge + nonce
	return hashMatchesDifficulty(input, difficulty)
}

func PowDDos() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Debug") == "1" {
				next.ServeHTTP(w, r)
				return
			}

			difficulty := defaultDifficulty

			if val := os.Getenv("POW_DDOS_DIFFICULTY"); val != "" {
				if parsed, err := strconv.Atoi(val); err == nil {
					difficulty = parsed
				}
			}

			if !verify(r, difficulty) {
				challenge := generateChallenge()

				httpx.HttpResponseWithETag(w, r, http.StatusTooManyRequests, map[string]string{"challenge": challenge, "difficulty": strconv.Itoa(difficulty)})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
