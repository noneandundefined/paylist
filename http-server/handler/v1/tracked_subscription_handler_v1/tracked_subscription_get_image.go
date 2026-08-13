package tracked_subscription_handler_v1

import (
	"net/http"
	"os"

	"paylist.server/middleware"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/ml"
)

func (h *Handler) GetSubscriptionImageHandler_V1(w http.ResponseWriter, r *http.Request) error {
	tr := middleware.TranslatorFromContext(r.Context())
	
	nlp := ml.NewNLPBuilder()

	name := r.URL.Query().Get("name")
	if name == "" {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-image-not-found"))
	}

	img := nlp.GetSubscriptionImage(name)

	if img == "" {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-image-not-found"))
	}

	filePath := "media/subscriptions/" + img

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-image-not-found"))
	}

	http.ServeFile(w, r, filePath)
	return nil
}
