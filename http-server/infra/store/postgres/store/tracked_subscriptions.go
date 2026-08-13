package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/pgqx"
)

type TrackedSubscriptionStore struct {
	db *sql.DB
}

func (s *TrackedSubscriptionStore) Create_Subscription(ctx context.Context, sub *models.TrackedSubscription) error {
	query := `
		INSERT INTO tracked_subscriptions (
			user_uuid, name, price, currency, period, date_pay,
			auto_renewal, notification, include_in_analytics
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	currency := sub.Currency
	if currency == "" {
		currency = "USD"
	}

	period := sub.Period
	if period == "" {
		period = "monthly"
	}

	if err := s.db.QueryRowContext(
		ctx,
		query,
		sub.UserUUID,
		sub.Name,
		sub.Price,
		currency,
		period,
		sub.DatePay,
		sub.AutoRenewal,
		sub.Notification,
		sub.IncludeInAnalytics,
	).Scan(&sub.ID); err != nil {
		logger.Error("Create_Subscription req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *TrackedSubscriptionStore) Count_SubscriptionsByUuid(ctx context.Context, uuid string) (int, error) {
	query := `SELECT COUNT(*) FROM tracked_subscriptions WHERE user_uuid = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int
	if err := s.db.QueryRowContext(ctx, query, uuid).Scan(&count); err != nil {
		logger.Error("Count_SubscriptionsByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return 0, err
	}

	return count, nil
}

func (s *TrackedSubscriptionStore) Get_SubscriptionById(ctx context.Context, id uint64, uuid string) (*models.TrackedSubscription, error) {
	query := `
		SELECT * FROM tracked_subscriptions
		WHERE id = $1 AND user_uuid = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sub, err := pgqx.QueryRowContext[models.TrackedSubscription](ctx, s.db, query, id, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_SubscriptionById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return sub, nil
}

func (s *TrackedSubscriptionStore) Get_SubscriptionsByUuid(ctx context.Context, uuid, search string) (*[]models.TrackedSubscription, error) {
	var limit sql.NullInt32
	if err := s.db.QueryRowContext(ctx, `
		SELECT subscriptions.max_total_subscriptions
		FROM user_subscriptions
		JOIN subscriptions ON LOWER(user_subscriptions.plan_name) = LOWER(subscriptions.plan_name)
		WHERE user_subscriptions.user_uuid = $1
		LIMIT 1
	`, uuid).Scan(&limit); err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("Get_SubscriptionsByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	maxLimit := 1000000
	if limit.Valid {
		maxLimit = int(limit.Int32)
		if maxLimit == 0 {
			maxLimit = 10
		}
	}

	query := `SELECT * FROM tracked_subscriptions WHERE user_uuid = $1`
	args := []any{uuid}
	paramIndex := 2

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", paramIndex)
		args = append(args, "%"+search+"%")
		paramIndex++
	}

	query += fmt.Sprintf(" ORDER BY date_pay ASC LIMIT $%d", paramIndex)
	args = append(args, maxLimit)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subs, err := pgqx.QueryContext[models.TrackedSubscription](ctx, s.db, query, args...)
	if err != nil {
		logger.Error("Get_SubscriptionsByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &subs, nil
}

func (s *TrackedSubscriptionStore) Get_SubscriptionsForTelegramNotify(ctx context.Context) (*[]models.TrackedSubscriptionNotifyCandidate, error) {
	query := `
		SELECT tracked_subscriptions.*,
			user_settings.telegram_chat_id,
			COALESCE(NULLIF(user_settings.telegram_language, ''), 'en') AS telegram_language,
			CASE
				WHEN tracked_subscriptions.date_pay::date = CURRENT_DATE THEN 'today'
				ELSE 'before_3d'
			END AS notify_kind
		FROM tracked_subscriptions
		JOIN user_subscriptions ON user_subscriptions.user_uuid = tracked_subscriptions.user_uuid
		JOIN subscriptions ON LOWER(subscriptions.plan_name) = LOWER(user_subscriptions.plan_name)
		JOIN user_settings ON user_settings.user_uuid = tracked_subscriptions.user_uuid
		WHERE tracked_subscriptions.notification = TRUE
			AND subscriptions.notification_subscriptions = TRUE
			AND user_subscriptions.is_active = TRUE
			AND user_settings.telegram_chat_id IS NOT NULL
			AND (
				tracked_subscriptions.date_pay::date - INTERVAL '3 days' = CURRENT_DATE
				OR tracked_subscriptions.date_pay::date = CURRENT_DATE
			)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subs, err := pgqx.QueryContext[models.TrackedSubscriptionNotifyCandidate](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_SubscriptionsForTelegramNotify req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &subs, nil
}

func (s *TrackedSubscriptionStore) Get_CategorySlugsBySubscriptionID(ctx context.Context, id uint64) ([]string, error) {
	query := `
		SELECT subscription_categories.slug
		FROM tracked_subscription_categories
		JOIN subscription_categories ON subscription_categories.id = tracked_subscription_categories.category_id
		WHERE tracked_subscription_categories.tracked_subscription_id = $1
		ORDER BY subscription_categories.slug ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		logger.Error("Get_CategorySlugsBySubscriptionID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	slugs := make([]string, 0)
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			logger.Error("Get_CategorySlugsBySubscriptionID req={%s}: Failed to scan row: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}

		slugs = append(slugs, slug)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Get_CategorySlugsBySubscriptionID req={%s}: Failed to iterate rows: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return slugs, nil
}

func (s *TrackedSubscriptionStore) Get_CategorySlugsMapByUserUUID(ctx context.Context, userUUID string) (map[uint64][]string, error) {
	query := `
		SELECT tracked_subscription_categories.tracked_subscription_id, subscription_categories.slug
		FROM tracked_subscription_categories
		JOIN subscription_categories ON subscription_categories.id = tracked_subscription_categories.category_id
		JOIN tracked_subscriptions ON tracked_subscriptions.id = tracked_subscription_categories.tracked_subscription_id
		WHERE tracked_subscriptions.user_uuid = $1
		ORDER BY tracked_subscription_categories.tracked_subscription_id ASC, subscription_categories.slug ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userUUID)
	if err != nil {
		logger.Error("Get_CategorySlugsMapByUserUUID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	result := make(map[uint64][]string)
	for rows.Next() {
		var subscriptionID uint64
		var slug string
		if err := rows.Scan(&subscriptionID, &slug); err != nil {
			logger.Error("Get_CategorySlugsMapByUserUUID req={%s}: Failed to scan row: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}

		result[subscriptionID] = append(result[subscriptionID], slug)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Get_CategorySlugsMapByUserUUID req={%s}: Failed to iterate rows: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return result, nil
}

func (s *TrackedSubscriptionStore) Get_AllSubscriptionsByUuid(ctx context.Context, uuid string) ([]models.TrackedSubscription, error) {
	query := `
		SELECT * FROM tracked_subscriptions
		WHERE user_uuid = $1
		ORDER BY date_pay ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subs, err := pgqx.QueryContext[models.TrackedSubscription](ctx, s.db, query, uuid)
	if err != nil {
		logger.Error("Get_AllSubscriptionsByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return subs, nil
}

func (s *TrackedSubscriptionStore) Replace_SubscriptionCategories(ctx context.Context, id uint64, userUUID string, slugs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Replace_SubscriptionCategories req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, `DELETE FROM tracked_subscription_categories WHERE tracked_subscription_id = $1`, id); err != nil {
		logger.Error("Replace_SubscriptionCategories req={%s}: Failed to delete categories: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if len(slugs) == 0 {
		return tx.Commit()
	}

	for _, slug := range slugs {
		var categoryID uint64
		if err := tx.QueryRowContext(
			ctx,
			`SELECT id FROM subscription_categories WHERE slug = $1 AND (user_uuid IS NULL OR user_uuid = $2) ORDER BY (user_uuid IS NULL) ASC LIMIT 1`,
			slug,
			userUUID,
		).Scan(&categoryID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("unknown category slug: %s", slug)
			}

			logger.Error("Replace_SubscriptionCategories req={%s}: Failed to find category: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO tracked_subscription_categories (tracked_subscription_id, category_id) VALUES ($1, $2)`,
			id,
			categoryID,
		); err != nil {
			logger.Error("Replace_SubscriptionCategories req={%s}: Failed to insert category link: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}
	}

	return tx.Commit()
}

func (s *TrackedSubscriptionStore) Update_SubscriptionById(ctx context.Context, sub *models.TrackedSubscription, id int) error {
	query := `
		UPDATE tracked_subscriptions
		SET
		    name = $1,
		    price = $2,
		    currency = $3,
		    period = $4,
		    date_pay = $5,
		    auto_renewal = $6,
		    notification = $7,
		    include_in_analytics = $8,
		    note = $9
		WHERE id = $10 AND user_uuid = $11
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	currency := sub.Currency
	if currency == "" {
		currency = "USD"
	}

	period := sub.Period
	if period == "" {
		period = "monthly"
	}

	upd, err := s.db.ExecContext(
		ctx,
		query,
		sub.Name,
		sub.Price,
		currency,
		period,
		sub.DatePay,
		sub.AutoRenewal,
		sub.Notification,
		sub.IncludeInAnalytics,
		sub.Note,
		id,
		sub.UserUUID,
	)
	if err != nil {
		logger.Error("Update_SubscriptionById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		logger.Error("Update_SubscriptionById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *TrackedSubscriptionStore) Update_SubscriptionsMounth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	query := `
		WITH renewed AS (
			UPDATE tracked_subscriptions
			SET date_pay = date_pay + CASE
				WHEN period = 'yearly' THEN INTERVAL '1 year'
				ELSE INTERVAL '1 month'
			END
			WHERE date_pay <= CURRENT_DATE
				AND auto_renewal = TRUE
			RETURNING
				id,
				user_uuid,
				period,
				(date_pay - CASE
					WHEN period = 'yearly' THEN INTERVAL '1 year'
					ELSE INTERVAL '1 month'
				END)::date AS previous_date_pay,
				date_pay::date AS new_date_pay,
				price,
				currency
		)
		INSERT INTO tracked_subscription_history (
			tracked_subscription_id, user_uuid, event_type,
			previous_date_pay, new_date_pay, price, currency
		)
		SELECT id, user_uuid, 'renewed', previous_date_pay, new_date_pay, price, currency
		FROM renewed
	`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		logger.Error("Update_SubscriptionsMounth req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *TrackedSubscriptionStore) Delete_SubscriptionById(ctx context.Context, id int, uuid string) error {
	query := `DELETE FROM tracked_subscriptions WHERE id = $1 AND user_uuid = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	del, err := s.db.ExecContext(ctx, query, id, uuid)
	if err != nil {
		logger.Error("Delete_SubscriptionById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	delAffected, err := del.RowsAffected()
	if err != nil {
		logger.Error("Delete_SubscriptionById req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if delAffected == 0 {
		return httperr.Err_NotDeleted
	}

	return nil
}
