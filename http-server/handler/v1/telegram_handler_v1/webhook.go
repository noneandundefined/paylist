package telegram_handler_v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"paylist.server/infra/logger"
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

		if err := notifier.HandleUpdate(r.Context(), update); err != nil {
			if telegram.IsDeliveryRejected(err) {
				logger.Warning("Telegram webhook skip blocked chat: %s", err.Error())
			} else {
				logger.Error("Telegram webhook update handling error: %s", err.Error())
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
