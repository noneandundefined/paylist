package payment_handler_v1

import (
	"net/http"
	"strconv"
	"time"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) GetPaymentHistoryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	rows, err := h.Store.Payments.Get_PaymentHistoryListByUserUuid(ctx, authToken.User.UserUUID, limit)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	resp := make([]PaymentHistoryResponse, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, paymentHistoryToResponse(row))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, resp)
	return nil
}

func paymentHistoryToResponse(row models.PaymentHistory) PaymentHistoryResponse {
	item := PaymentHistoryResponse{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt.UTC().Format(time.RFC3339),
		PlanName:    row.PlanName,
		Amount:      row.Amount,
		Currency:    row.Currency,
		Status:      row.Status,
		PaymentKind: row.PaymentKind,
		Description: row.Description,
	}

	if row.PaidAt != nil {
		paidAt := row.PaidAt.UTC().Format(time.RFC3339)
		item.PaidAt = &paidAt
	}

	return item
}
