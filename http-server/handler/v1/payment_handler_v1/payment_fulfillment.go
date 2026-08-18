package payment_handler_v1

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/infra/store/postgres/store"
	"paylist.server/pkg/yookassa"
)

func applyYookassaPayment(ctx context.Context, storage store.Storage, payment *yookassa.Payment) (string, error) {
	status := mapYookassaStatus(payment.Status, payment.Paid)

	existing, err := storage.Payments.Get_PaymentHistoryByYookassaPaymentID(ctx, payment.ID)
	if err != nil {
		return status, err
	}

	if existing != nil && existing.Status == models.PaymentStatusSucceeded && status != models.PaymentStatusSucceeded {
		return existing.Status, nil
	}

	if err := storage.Payments.Update_PaymentHistoryStatus(ctx, payment.ID, status, paidAtFromPayment(payment)); err != nil {
		return status, err
	}

	if status == models.PaymentStatusSucceeded {
		if err := fulfillSucceededPayment(ctx, storage, payment); err != nil {
			return status, err
		}
	}

	return status, nil
}

func fulfillSucceededPayment(ctx context.Context, storage store.Storage, payment *yookassa.Payment) error {
	userUUID := strings.TrimSpace(payment.Metadata["user_uuid"])
	planName := strings.TrimSpace(payment.Metadata["plan_name"])
	if userUUID == "" || planName == "" {
		return nil
	}

	plan, err := storage.Subscriptions.Get_SubscriptionByPlanName(ctx, planName)
	if err != nil {
		return err
	}

	if plan == nil {
		return nil
	}

	if err := storage.Payments.Update_ActivateUserSubscriptionPlan(ctx, userUUID, plan.PlanName, plan.DurationDays); err != nil {
		return err
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

func yookassaPremiumDescription(renewal bool, durationDays int) string {
	action := "Оплата"
	if renewal {
		action = "Продление"
	}

	return fmt.Sprintf(`%s подписки "Paylist" на %s`, action, formatRussianMonths(monthsFromDurationDays(durationDays)))
}

func monthsFromDurationDays(durationDays int) int {
	if durationDays <= 0 {
		return 1
	}

	months := durationDays / 30
	if months < 1 {
		return 1
	}

	return months
}

func formatRussianMonths(months int) string {
	abs := months % 100
	if abs >= 11 && abs <= 14 {
		return fmt.Sprintf("%d месяцев", months)
	}

	switch months % 10 {
	case 1:
		return fmt.Sprintf("%d месяц", months)
	case 2, 3, 4:
		return fmt.Sprintf("%d месяца", months)
	default:
		return fmt.Sprintf("%d месяцев", months)
	}
}

func paymentMethodDisplayTitle(method yookassa.PaymentMethod) string {
	if method.Card != nil && strings.TrimSpace(method.Card.Last4) != "" {
		return fmt.Sprintf("%s **** %s", formatCardBrand(method.Card.CardType), strings.TrimSpace(method.Card.Last4))
	}

	title := strings.TrimSpace(method.Title)
	if last4 := lastFourDigits(title); last4 != "" {
		if paymentMethodTitleLooksFormatted(title) {
			return title
		}

		return fmt.Sprintf("%s **** %s", brandFromPaymentMethod(method.Type, title), last4)
	}

	if title != "" {
		return title
	}

	return method.Type
}

func paymentMethodTitleLooksFormatted(title string) bool {
	return strings.Contains(strings.TrimSpace(title), " **** ")
}

func lastFourDigits(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 4 {
		return ""
	}

	last := value[len(value)-4:]
	for i := 0; i < 4; i++ {
		if last[i] < '0' || last[i] > '9' {
			return ""
		}
	}

	return last
}

func brandFromPaymentMethod(methodType, title string) string {
	haystack := strings.ToLower(methodType + " " + title)

	switch {
	case strings.Contains(haystack, "mir"):
		return "MIR"
	case strings.Contains(haystack, "visa"):
		return "Visa"
	case strings.Contains(haystack, "master"):
		return "Mastercard"
	case strings.Contains(haystack, "union"):
		return "UnionPay"
	case methodType == "sbp" || strings.Contains(haystack, "sbp"):
		return "СБП"
	case methodType == "yoo_money" || strings.Contains(haystack, "yoo"):
		return "ЮMoney"
	default:
		return formatCardBrand("")
	}
}

func formatCardBrand(cardType string) string {
	switch strings.ToLower(strings.TrimSpace(cardType)) {
	case "mir":
		return "MIR"
	case "visa":
		return "Visa"
	case "mastercard", "master card":
		return "Mastercard"
	case "unionpay":
		return "UnionPay"
	case "americanexpress", "american express":
		return "American Express"
	case "":
		return "Card"
	default:
		return cardType
	}
}

func amountValueToFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
	if err != nil {
		return 0
	}

	return parsed
}
