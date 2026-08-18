package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"paylist.server/infra/constants"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/pgqx"
)

type TrackedSubscriptionStore struct {
	db *sql.DB
}

const trackedSubscriptionSelect = `
	SELECT
		tracked_subscriptions.id,
		tracked_subscriptions.created_at,
		tracked_subscriptions.updated_at,
		tracked_subscriptions.user_uuid,
		tracked_subscriptions.name,
		tracked_subscriptions.tariff,
		tracked_subscriptions.price,
		tracked_subscriptions.currency,
		tracked_subscriptions.period,
		tracked_subscriptions.date_pay,
		tracked_subscriptions.auto_renewal,
		tracked_subscription_members.note,
		tracked_subscription_members.notification,
		tracked_subscription_members.include_in_analytics,
		tracked_subscription_members.share_percent,
		ROUND((tracked_subscriptions.price * tracked_subscription_members.share_percent / 100)::numeric, 3) AS share_price,
		(tracked_subscription_members.role = 'owner') AS is_owner
	FROM tracked_subscriptions
	JOIN tracked_subscription_members ON tracked_subscription_members.tracked_subscription_id = tracked_subscriptions.id
`

func (s *TrackedSubscriptionStore) Create_Subscription(ctx context.Context, sub *models.TrackedSubscription) error {
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

	tariff := sub.Tariff
	if tariff == "" {
		tariff = constants.TariffNone
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Create_Subscription req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := tx.QueryRowContext(
		ctx,
		`
			INSERT INTO tracked_subscriptions (
				user_uuid, name, tariff, price, currency, period, date_pay,
				auto_renewal, notification, include_in_analytics
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`,
		sub.UserUUID,
		sub.Name,
		tariff,
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

	var ownerEmail string
	if err := tx.QueryRowContext(ctx, `SELECT email FROM user_cores WHERE user_uuid = $1`, sub.UserUUID).Scan(&ownerEmail); err != nil {
		logger.Error("Create_Subscription req={%s}: Failed to load owner email: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO tracked_subscription_members (
				tracked_subscription_id, user_uuid, email, role, share_percent,
				notification, include_in_analytics, status, note
			)
			VALUES ($1, $2, $3, 'owner', 100, $4, $5, 'accepted', $6)
		`,
		sub.ID,
		sub.UserUUID,
		ownerEmail,
		sub.Notification,
		sub.IncludeInAnalytics,
		sub.Note,
	); err != nil {
		logger.Error("Create_Subscription req={%s}: Failed to insert owner member: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	sub.SharePercent = 100
	sub.SharePrice = sub.Price
	sub.IsOwner = true

	return tx.Commit()
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
	query := trackedSubscriptionSelect + `
		WHERE tracked_subscriptions.id = $1
			AND tracked_subscription_members.user_uuid = $2
			AND tracked_subscription_members.status = 'accepted'
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

	query := trackedSubscriptionSelect + ` WHERE tracked_subscription_members.user_uuid = $1 AND tracked_subscription_members.status = 'accepted'`
	args := []any{uuid}
	paramIndex := 2

	if search != "" {
		query += fmt.Sprintf(" AND tracked_subscriptions.name ILIKE $%d", paramIndex)
		args = append(args, "%"+search+"%")
		paramIndex++
	}

	query += fmt.Sprintf(" ORDER BY tracked_subscriptions.date_pay ASC LIMIT $%d", paramIndex)
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
		SELECT
			tracked_subscriptions.id,
			tracked_subscriptions.created_at,
			tracked_subscriptions.updated_at,
			tracked_subscriptions.user_uuid,
			tracked_subscriptions.name,
			tracked_subscriptions.tariff,
			tracked_subscriptions.price,
			tracked_subscriptions.currency,
			tracked_subscriptions.period,
			tracked_subscriptions.date_pay,
			tracked_subscriptions.auto_renewal,
			tracked_subscription_members.note,
			tracked_subscription_members.notification,
			tracked_subscription_members.include_in_analytics,
			tracked_subscription_members.share_percent,
			ROUND((tracked_subscriptions.price * tracked_subscription_members.share_percent / 100)::numeric, 3) AS share_price,
			(tracked_subscription_members.role = 'owner') AS is_owner,
			tracked_subscription_members.user_uuid AS member_user_uuid,
			user_settings.telegram_chat_id,
			COALESCE(NULLIF(user_settings.telegram_language, ''), 'en') AS telegram_language,
			CASE
				WHEN tracked_subscriptions.date_pay::date = CURRENT_DATE THEN 'today'
				ELSE 'before_3d'
			END AS notify_kind
		FROM tracked_subscriptions
		JOIN tracked_subscription_members ON tracked_subscription_members.tracked_subscription_id = tracked_subscriptions.id
		JOIN user_subscriptions ON user_subscriptions.user_uuid = tracked_subscription_members.user_uuid
		JOIN subscriptions ON LOWER(subscriptions.plan_name) = LOWER(user_subscriptions.plan_name)
		JOIN user_settings ON user_settings.user_uuid = tracked_subscription_members.user_uuid
		WHERE tracked_subscription_members.status = 'accepted'
			AND tracked_subscription_members.user_uuid IS NOT NULL
			AND tracked_subscription_members.notification = TRUE
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

func (s *TrackedSubscriptionStore) Get_SubscriptionsForMaxNotify(ctx context.Context) (*[]models.TrackedSubscriptionNotifyCandidate, error) {
	query := `
		SELECT
			tracked_subscriptions.id,
			tracked_subscriptions.created_at,
			tracked_subscriptions.updated_at,
			tracked_subscriptions.user_uuid,
			tracked_subscriptions.name,
			tracked_subscriptions.tariff,
			tracked_subscriptions.price,
			tracked_subscriptions.currency,
			tracked_subscriptions.period,
			tracked_subscriptions.date_pay,
			tracked_subscriptions.auto_renewal,
			tracked_subscription_members.note,
			tracked_subscription_members.notification,
			tracked_subscription_members.include_in_analytics,
			tracked_subscription_members.share_percent,
			ROUND((tracked_subscriptions.price * tracked_subscription_members.share_percent / 100)::numeric, 3) AS share_price,
			(tracked_subscription_members.role = 'owner') AS is_owner,
			tracked_subscription_members.user_uuid AS member_user_uuid,
			user_settings.max_user_id,
			COALESCE(NULLIF(user_settings.max_language, ''), 'en') AS max_language,
			CASE
				WHEN tracked_subscriptions.date_pay::date = CURRENT_DATE THEN 'today'
				ELSE 'before_3d'
			END AS notify_kind
		FROM tracked_subscriptions
		JOIN tracked_subscription_members ON tracked_subscription_members.tracked_subscription_id = tracked_subscriptions.id
		JOIN user_subscriptions ON user_subscriptions.user_uuid = tracked_subscription_members.user_uuid
		JOIN subscriptions ON LOWER(subscriptions.plan_name) = LOWER(user_subscriptions.plan_name)
		JOIN user_settings ON user_settings.user_uuid = tracked_subscription_members.user_uuid
		WHERE tracked_subscription_members.status = 'accepted'
			AND tracked_subscription_members.user_uuid IS NOT NULL
			AND tracked_subscription_members.notification = TRUE
			AND subscriptions.notification_subscriptions = TRUE
			AND user_subscriptions.is_active = TRUE
			AND user_settings.max_user_id IS NOT NULL
			AND (
				tracked_subscriptions.date_pay::date - INTERVAL '3 days' = CURRENT_DATE
				OR tracked_subscriptions.date_pay::date = CURRENT_DATE
			)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subs, err := pgqx.QueryContext[models.TrackedSubscriptionNotifyCandidate](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_SubscriptionsForMaxNotify req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &subs, nil
}

func (s *TrackedSubscriptionStore) Get_CategorySlugsBySubscriptionID(ctx context.Context, id uint64, userUUID string) ([]string, error) {
	query := `
		SELECT subscription_categories.slug
		FROM tracked_subscription_categories
		JOIN subscription_categories ON subscription_categories.id = tracked_subscription_categories.category_id
		WHERE tracked_subscription_categories.tracked_subscription_id = $1
			AND tracked_subscription_categories.user_uuid = $2
		ORDER BY subscription_categories.slug ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id, userUUID)
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
		JOIN tracked_subscription_members ON tracked_subscription_members.tracked_subscription_id = tracked_subscriptions.id
		WHERE tracked_subscription_members.user_uuid = $1
			AND tracked_subscription_members.status = 'accepted'
			AND tracked_subscription_categories.user_uuid = $1
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
	query := trackedSubscriptionSelect + `
		WHERE tracked_subscription_members.user_uuid = $1 AND tracked_subscription_members.status = 'accepted'
		ORDER BY tracked_subscriptions.date_pay ASC
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

	if _, err := tx.ExecContext(ctx, `DELETE FROM tracked_subscription_categories WHERE tracked_subscription_id = $1 AND user_uuid = $2`, id, userUUID); err != nil {
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
			`INSERT INTO tracked_subscription_categories (tracked_subscription_id, user_uuid, category_id) VALUES ($1, $2, $3)`,
			id,
			userUUID,
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
		    tariff = $2,
		    price = $3,
		    currency = $4,
		    period = $5,
		    date_pay = $6,
		    auto_renewal = $7,
		    notification = $8,
		    include_in_analytics = $9,
		    note = $10
		WHERE id = $11 AND user_uuid = $12
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

	tariff := sub.Tariff
	if tariff == "" {
		tariff = constants.TariffNone
	}

	upd, err := s.db.ExecContext(
		ctx,
		query,
		sub.Name,
		tariff,
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

func (s *TrackedSubscriptionStore) Update_MemberPreferences(ctx context.Context, subscriptionID uint64, userUUID string, notification, includeInAnalytics bool, note *string) error {
	query := `
		UPDATE tracked_subscription_members
		SET notification = $1, include_in_analytics = $2, note = $3
		WHERE tracked_subscription_id = $4
			AND user_uuid = $5
			AND status = 'accepted'
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, notification, includeInAnalytics, note, subscriptionID, userUUID)
	if err != nil {
		logger.Error("Update_MemberPreferences req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	updAffected, err := upd.RowsAffected()
	if err != nil {
		logger.Error("Update_MemberPreferences req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if updAffected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *TrackedSubscriptionStore) Enable_NotificationsForUser(ctx context.Context, userUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(
		ctx,
		`
			UPDATE tracked_subscription_members
			SET notification = TRUE
			WHERE user_uuid = $1 AND status = 'accepted'
		`,
		userUUID,
	); err != nil {
		logger.Error("Enable_NotificationsForUser req={%s}: Failed to update members: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if _, err := s.db.ExecContext(
		ctx,
		`
			UPDATE tracked_subscriptions
			SET notification = TRUE
			WHERE user_uuid = $1
		`,
		userUUID,
	); err != nil {
		logger.Error("Enable_NotificationsForUser req={%s}: Failed to update subscriptions: %s", ctx.Value("XREQID").(string), err.Error())
		return err
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
			WHERE date_pay < CURRENT_DATE
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

func (s *TrackedSubscriptionStore) Get_CrowdSubscriptionPrices(ctx context.Context, excludeUserUUID string) ([]models.CrowdSubscriptionPrice, error) {
	query := `
		SELECT
			tracked_subscriptions.name,
			tracked_subscriptions.tariff,
			tracked_subscriptions.price,
			tracked_subscriptions.currency,
			tracked_subscriptions.period,
			user_settings.country
		FROM tracked_subscriptions
		JOIN tracked_subscription_members ON tracked_subscription_members.tracked_subscription_id = tracked_subscriptions.id
			AND tracked_subscription_members.role = 'owner'
			AND tracked_subscription_members.status = 'accepted'
		LEFT JOIN user_settings ON user_settings.user_uuid = tracked_subscriptions.user_uuid
		WHERE tracked_subscriptions.user_uuid <> $1
			AND tracked_subscription_members.include_in_analytics = TRUE
	`

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	rows, err := pgqx.QueryContext[models.CrowdSubscriptionPrice](ctx, s.db, query, excludeUserUUID)
	if err != nil {
		logger.Error("Get_CrowdSubscriptionPrices req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if rows == nil {
		return []models.CrowdSubscriptionPrice{}, nil
	}

	return rows, nil
}
