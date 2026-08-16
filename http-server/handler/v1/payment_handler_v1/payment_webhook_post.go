package payment_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/postgres/models"
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

	status := mapYookassaStatus(payment.Status, payment.Paid)
	if err := h.Store.Payments.Update_PaymentHistoryStatus(ctx, payment.ID, status, paidAtFromPayment(payment)); err != nil {
		return httperr.Db(ctx, err)
	}

	switch notification.Event {
	case "payment.succeeded":
		if err := fulfillSucceededPayment(ctx, h.Store, payment); err != nil {
			return httperr.Db(ctx, err)
		}

	case "payment.canceled":
		existing, err := h.Store.Payments.Get_PaymentHistoryByYookassaPaymentID(ctx, payment.ID)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if existing != nil && existing.Status == models.PaymentStatusSucceeded {
			break
		}
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, map[string]string{"status": "ok"})
	return nil
}
