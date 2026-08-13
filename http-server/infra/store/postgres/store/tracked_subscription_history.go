package store

import (
	"context"
	"database/sql"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
)

func (s *TrackedSubscriptionStore) Create_SubscriptionHistory(ctx context.Context, entry *models.TrackedSubscriptionHistory) error {
	query := `
		INSERT INTO tracked_subscription_history (
			tracked_subscription_id, user_uuid, event_type,
			previous_date_pay, new_date_pay, price, currency
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	currency := entry.Currency
	if currency == "" {
		currency = "USD"
	}

	var previousDatePay sql.NullTime
	if entry.PreviousDatePay != nil {
		previousDatePay = sql.NullTime{Time: *entry.PreviousDatePay, Valid: true}
	}

	if _, err := s.db.ExecContext(
		ctx,
		query,
		entry.TrackedSubscriptionID,
		entry.UserUUID,
		entry.EventType,
		previousDatePay,
		entry.NewDatePay,
		entry.Price,
		currency,
	); err != nil {
		logger.Error("Create_SubscriptionHistory req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
