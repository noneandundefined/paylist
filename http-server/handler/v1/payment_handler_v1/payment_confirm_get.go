package payment_handler_v1

import (
	"net/http"
	"strings"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) GetPaymentConfirmHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)
	ref := strings.TrimSpace(r.URL.Query().Get("payment_id"))

	if ref == "" {
		httpx.HttpResponse(w, r, http.StatusOK, PaymentConfirmResponse{Paid: false, Status: ""})
		return nil
	}

	row, err := h.Store.Payments.Get_PaymentHistoryByUserAndRef(ctx, authToken.User.UserUUID, ref)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if row == nil {
		httpx.HttpResponse(w, r, http.StatusOK, PaymentConfirmResponse{Paid: false, Status: ""})
		return nil
	}

	status := row.Status
	if status != models.PaymentStatusSucceeded && h.Yookassa != nil && h.Yookassa.Configured() && strings.TrimSpace(row.YookassaPaymentID) != "" {
		payment, ykErr := h.Yookassa.GetPayment(ctx, row.YookassaPaymentID)
		if ykErr != nil {
			return httperr.InternalServerError(ykErr.Error())
		}

		synced, applyErr := applyYookassaPayment(ctx, h.Store, payment)
		if applyErr != nil {
			return httperr.Db(ctx, applyErr)
		}

		status = synced
	}

	httpx.HttpResponse(w, r, http.StatusOK, PaymentConfirmResponse{
		Paid:   status == models.PaymentStatusSucceeded,
		Status: status,
	})
	return nil
}
