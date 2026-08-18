package payment_handler_v1

import (
	"net/http"

	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/yookassa"
)

func (h *Handler) PostWebhookHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	if h.Yookassa == nil || !h.Yookassa.Configured() {
		httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "ignored"})
		return nil
	}

	var notification yookassa.WebhookNotification
	if err := httpx.HttpParse(r, &notification); err != nil {
		return httperr.BadRequest(err.Error())
	}

	paymentID := notification.Object.ID
	if paymentID == "" {
		httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "ignored"})
		return nil
	}

	payment, err := h.Yookassa.GetPayment(ctx, paymentID)
	if err != nil {
		return httperr.InternalServerError(err.Error())
	}

	if _, err := applyYookassaPayment(ctx, h.Store, payment); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "ok"})
	return nil
}
