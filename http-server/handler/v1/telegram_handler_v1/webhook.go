package telegram_handler_v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"paylist.server/pkg/telegram"
)

func WebhookHandler(notifier *telegram.Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_SECRET"))
		if secret == "" || notifier == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var update telegram.Update

		if err := json.Unmarshal(body, &update); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_ = notifier.HandleUpdate(r.Context(), update)
		w.WriteHeader(http.StatusOK)
	}
}
