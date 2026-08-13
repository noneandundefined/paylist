package payment_handler_v1

import (
	"context"
	"strconv"
	"strings"
	"time"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/yookassa"
)

func (h *Handler) fulfillSucceededPayment(ctx context.Context, payment *yookassa.Payment) error {
	userUUID := strings.TrimSpace(payment.Metadata["user_uuid"])
	planName := strings.TrimSpace(payment.Metadata["plan_name"])
	if userUUID == "" || planName == "" {
		return nil
	}

	plan, err := h.Store.Subscriptions.Get_SubscriptionByPlanName(ctx, planName)
	if err != nil {
		return err
	}

	if plan == nil {
		return nil
	}

	if err := h.Store.Payments.Update_ActivateUserSubscriptionPlan(ctx, userUUID, plan.PlanName, plan.DurationDays); err != nil {
		return err
	}

	if payment.PaymentMethod.ID != "" && payment.PaymentMethod.Saved {
		title := payment.PaymentMethod.Title
		if title == "" {
			title = payment.PaymentMethod.Type
		}

		if err := h.Store.Payments.Update_YookassaPaymentMethod(ctx, userUUID, payment.PaymentMethod.ID, payment.PaymentMethod.Type, title); err != nil {
			return err
		}
	}

	return nil
}

func mapYookassaStatus(status string, paid bool) string {
	switch status {
	case "succeeded":
		return models.PaymentStatusSucceeded
	case "waiting_for_capture":
		return models.PaymentStatusWaitingForCapture
	case "canceled":
		return models.PaymentStatusCanceled
	case "pending":
		return models.PaymentStatusPending
	default:
		if paid {
			return models.PaymentStatusSucceeded
		}

		return models.PaymentStatusFailed
	}
}

func paidAtFromPayment(payment *yookassa.Payment) *time.Time {
	if !payment.Paid && payment.Status != "succeeded" {
		return nil
	}

	now := time.Now().UTC()
	return &now
}

func amountValueToFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
	if err != nil {
		return 0
	}

	return parsed
}
