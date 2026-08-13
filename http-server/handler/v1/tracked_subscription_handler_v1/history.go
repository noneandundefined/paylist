package tracked_subscription_handler_v1

import (
	"context"
	"time"

	"paylist.server/infra/store/postgres/models"
)

const historyEventDateChanged = "date_changed"

type subscriptionHistoryStore interface {
	Create_SubscriptionHistory(ctx context.Context, entry *models.TrackedSubscriptionHistory) error
}

func isSameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()

	return ay == by && am == bm && ad == bd
}

func writeSubscriptionHistory(ctx context.Context, store subscriptionHistoryStore, sub *models.TrackedSubscription, eventType string, previousDatePay *time.Time) error {
	currency := sub.Currency
	if currency == "" {
		currency = "USD"
	}

	return store.Create_SubscriptionHistory(ctx, &models.TrackedSubscriptionHistory{
		TrackedSubscriptionID: sub.ID,
		UserUUID:              sub.UserUUID,
		EventType:             eventType,
		PreviousDatePay:       previousDatePay,
		NewDatePay:            sub.DatePay,
		Price:                 sub.Price,
		Currency:              currency,
	})
}
