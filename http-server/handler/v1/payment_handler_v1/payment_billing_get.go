package payment_handler_v1

import (
	"context"
	"net/http"
	"strings"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/yookassa"
	"paylist.server/types"
)

func (h *Handler) GetBillingHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	billing, err := h.Store.Payments.Get_UserSubscriptionBillingByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	h.refreshStoredPaymentMethodTitle(ctx, authToken.User.UserUUID, billing)

	httpx.HttpResponseWithETag(w, r, http.StatusOK, models.BillingStateFromSubscription(billing))
	return nil
}

func (h *Handler) refreshStoredPaymentMethodTitle(ctx context.Context, userUUID string, billing *models.UserSubscriptionBilling) {
	if billing == nil || billing.YookassaPaymentMethodID == nil {
		return
	}

	methodID := strings.TrimSpace(*billing.YookassaPaymentMethodID)
	if methodID == "" {
		return
	}

	currentTitle := ""
	if billing.PaymentMethodTitle != nil {
		currentTitle = strings.TrimSpace(*billing.PaymentMethodTitle)
	}

	if paymentMethodTitleLooksFormatted(currentTitle) {
		return
	}

	methodType := ""
	if billing.PaymentMethodType != nil {
		methodType = strings.TrimSpace(*billing.PaymentMethodType)
	}

	title := currentTitle
	if h.Yookassa != nil && h.Yookassa.Configured() {
		if method := h.loadYookassaPaymentMethod(ctx, userUUID, methodID); method != nil {
			title = paymentMethodDisplayTitle(*method)
			if strings.TrimSpace(method.Type) != "" {
				methodType = method.Type
			}
		}
	}

	if !paymentMethodTitleLooksFormatted(title) {
		title = paymentMethodDisplayTitle(yookassa.PaymentMethod{
			Type:  methodType,
			Title: title,
		})
	}

	if title == "" {
		return
	}

	if title != currentTitle {
		billing.PaymentMethodTitle = &title
		if methodType != "" {
			billing.PaymentMethodType = &methodType
		}
	}

	if !paymentMethodTitleLooksFormatted(title) || title == currentTitle {
		return
	}

	if err := h.Store.Payments.Update_YookassaPaymentMethod(ctx, userUUID, methodID, methodType, title); err != nil {
		logger.Error("Save formatted payment method title failed: %s", err.Error())
	}
}

func (h *Handler) loadYookassaPaymentMethod(ctx context.Context, userUUID, methodID string) *yookassa.PaymentMethod {
	method, err := h.Yookassa.GetPaymentMethod(ctx, methodID)
	if err != nil {
		logger.Error("YooKassa get payment method failed: %s", err.Error())
	} else if method != nil && method.Card != nil {
		return method
	}

	payments, err := h.Store.Payments.Get_PaymentHistoryListByUserUuid(ctx, userUUID, 10)
	if err != nil {
		return method
	}

	for i := range payments {
		payment := payments[i]
		if payment.Status != models.PaymentStatusSucceeded || strings.TrimSpace(payment.YookassaPaymentID) == "" {
			continue
		}

		ykPayment, err := h.Yookassa.GetPayment(ctx, payment.YookassaPaymentID)
		if err != nil {
			logger.Error("YooKassa get payment failed: %s", err.Error())
			continue
		}

		if ykPayment != nil && (ykPayment.PaymentMethod.Card != nil || strings.TrimSpace(ykPayment.PaymentMethod.Title) != "") {
			pm := ykPayment.PaymentMethod
			return &pm
		}
	}

	return method
}
