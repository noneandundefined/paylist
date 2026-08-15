package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"paylist.server/infra/logger"
)

func (s *TrackedSubscriptionStore) Create_SubscriptionNotificationLog(ctx context.Context, subscriptionID uint64, userUUID, channel string, notifyDate time.Time) error {
	query := `
		INSERT INTO subscription_notification_log (tracked_subscription_id, user_uuid, channel, notify_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, subscriptionID, userUUID, channel, notifyDate.Format("2006-01-02")); err != nil {
		logger.Error("Create_SubscriptionNotificationLog req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *TrackedSubscriptionStore) Has_SubscriptionNotificationLog(ctx context.Context, subscriptionID uint64, userUUID, channel string, notifyDate time.Time) (bool, error) {
	query := `
		SELECT 1
		FROM subscription_notification_log
		WHERE tracked_subscription_id = $1
			AND user_uuid = $2
			AND channel = $3
			AND notify_date = $4
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists int

	err := s.db.QueryRowContext(ctx, query, subscriptionID, userUUID, channel, notifyDate.Format("2006-01-02")).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		logger.Error("Has_SubscriptionNotificationLog req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return false, err
	}

	return true, nil
}

func (s *TrackedSubscriptionStore) Delete_OldSubscriptionNotificationLogs(ctx context.Context, olderThanDays int) error {
	query := `
		DELETE FROM subscription_notification_log
		WHERE notify_date < CURRENT_DATE - make_interval(days => $1)
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, olderThanDays); err != nil {
		logger.Error("Delete_OldSubscriptionNotificationLogs req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
