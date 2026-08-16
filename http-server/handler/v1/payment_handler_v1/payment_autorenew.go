package payment_handler_v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/pkg/yookassa"
)

func ChargeDueRenewals(ctx context.Context, storage store.Storage, client *yookassa.Client) error {
	if client == nil || !client.Configured() {
		return nil
	}

	due, err := storage.Payments.Get_SubscriptionsDueForAutoRenew(ctx)
	if err != nil {
		return err
	}

	for _, item := range due {
		if err := chargeRenewal(ctx, storage, client, item); err != nil {
			logger.Error("ChargeDueRenewals: user=%s plan=%s: %s", item.UserUUID, item.PlanName, err.Error())
		}
	}

	return nil
}

func chargeRenewal(ctx context.Context, storage store.Storage, client *yookassa.Client, item models.UserSubscriptionRenewalDue) error {
	amountValue := fmt.Sprintf("%.2f", item.Amount)
	idempotenceKey := uuid.NewString()

	ykPayment, err := client.CreatePayment(ctx, idempotenceKey, yookassa.CreatePaymentRequest{
		Amount: yookassa.Amount{
			Value:    amountValue,
			Currency: item.Currency,
		},
		Capture:         true,
		PaymentMethodID: item.YookassaPaymentMethodID,
		Description:     yookassaPremiumDescription(true, item.DurationDays),
		Metadata: map[string]string{
			"user_uuid": item.UserUUID,
			"plan_name": item.PlanName,
		},
	})
	if err != nil {
		return err
	}

	description := ykPayment.Description
	methodID := item.YookassaPaymentMethodID
	_, err = storage.Payments.Create_PaymentHistory(ctx, &models.PaymentHistory{
		UserUUID:                item.UserUUID,
		PlanName:                item.PlanName,
		YookassaPaymentID:       ykPayment.ID,
		YookassaPaymentMethodID: &methodID,
		Amount:                  item.Amount,
		Currency:                item.Currency,
		Status:                  mapYookassaStatus(ykPayment.Status, ykPayment.Paid),
		PaymentKind:             models.PaymentKindRenewal,
		Description:             &description,
	})
	if err != nil {
		return err
	}

	if ykPayment.Status == "succeeded" || ykPayment.Paid {
		return fulfillSucceededPayment(ctx, storage, ykPayment)
	}

	return nil
}
