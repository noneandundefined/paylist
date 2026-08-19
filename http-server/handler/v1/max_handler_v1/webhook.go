package max_handler_v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"paylist.server/infra/logger"
	"paylist.server/pkg/maxbot"
)

func WebhookHandler(notifier *maxbot.Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := strings.TrimSpace(os.Getenv("MAX_WEBHOOK_SECRET"))
		if secret == "" || notifier == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if headerSecret := strings.TrimSpace(r.Header.Get("X-Max-Bot-Api-Secret")); headerSecret != "" && headerSecret != secret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var update maxbot.Update

		if err := json.Unmarshal(body, &update); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := notifier.HandleUpdate(r.Context(), update); err != nil {
			logger.Error("MAX webhook update handling error: %s", err.Error())
		}
		w.WriteHeader(http.StatusOK)
	}
}
