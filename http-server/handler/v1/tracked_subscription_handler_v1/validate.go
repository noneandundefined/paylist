package tracked_subscription_handler_v1

import (
	"strings"
	"time"

	"paylist.server/infra/constants"
	"paylist.server/infra/locale"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/types"
)

func normalizePremiumSubscriptionFields(authToken *types.AuthToken, autoRenewal, notification bool) (bool, bool) {
	if !premium.IsPremiumPlan(authToken) || !authToken.NotificationSubscriptions {
		notification = false
	}

	return autoRenewal, notification
}

func validateTrackedSubscriptionPayload(
	tr locale.Translator,
	authToken *types.AuthToken,
	price float64,
	datePay time.Time,
	autoRenewal bool,
	notification bool,
) error {
	autoRenewal, notification = normalizePremiumSubscriptionFields(authToken, autoRenewal, notification)

	if price <= 0 {
		return httperr.BadRequest(tr.TErr("error.tracked-subscription-price-invalid"))
	}

	today := time.Now().Truncate(24 * time.Hour)
	if !datePay.After(today) {
		return httperr.BadRequest(tr.TErr("error.tracked-subscription-date-invalid"))
	}

	return nil
}

func checkTrackedSubscriptionLimit(tr locale.Translator, authToken *types.AuthToken, currentCount int) error {
	if authToken.MaxTotalSubscriptions == nil {
		return nil
	}

	if currentCount >= *authToken.MaxTotalSubscriptions {
		return httperr.BadRequest(tr.TErr("error.tracked-subscription-limit-reached"))
	}

	return nil
}

func includeInAnalyticsValue(value *bool) bool {
	if value == nil {
		return true
	}

	return *value
}

func normalizeCurrency(value string) string {
	if value == "" {
		return "USD"
	}

	return strings.ToUpper(value)
}

func normalizePeriod(value string) string {
	if value == "" {
		return "monthly"
	}

	return value
}

func normalizeTariff(value string) string {
	return constants.NormalizeTariff(value)
}

func normalizeNote(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func trackedSubscriptionFromCreatePayload(userUUID string, payload *TrackedSubscriptionCreatePayload) *models.TrackedSubscription {
	return &models.TrackedSubscription{
		UserUUID:           userUUID,
		Name:               payload.Name,
		Tariff:             normalizeTariff(payload.Tariff),
		Price:              payload.Price,
		Currency:           normalizeCurrency(payload.Currency),
		Period:             normalizePeriod(payload.Period),
		DatePay:            payload.DatePay.Time,
		AutoRenewal:        payload.AutoRenewal,
		Notification:       payload.Notification,
		IncludeInAnalytics: includeInAnalyticsValue(payload.IncludeInAnalytics),
	}
}

func trackedSubscriptionFromEditPayload(userUUID string, payload *TrackedSubscriptionEditPayload) *models.TrackedSubscription {
	return &models.TrackedSubscription{
		UserUUID:           userUUID,
		Name:               payload.Name,
		Tariff:             normalizeTariff(payload.Tariff),
		Price:              payload.Price,
		Currency:           normalizeCurrency(payload.Currency),
		Period:             normalizePeriod(payload.Period),
		DatePay:            payload.DatePay.Time,
		AutoRenewal:        payload.AutoRenewal,
		Notification:       payload.Notification,
		IncludeInAnalytics: includeInAnalyticsValue(payload.IncludeInAnalytics),
		Note:               normalizeNote(payload.Note),
	}
}
