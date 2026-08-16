package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/pgqx"
)

type PaymentStore struct {
	db *sql.DB
}

func (s *PaymentStore) Create_PaymentHistory(ctx context.Context, payment *models.PaymentHistory) (*models.PaymentHistory, error) {
	query := `
		INSERT INTO payment_history (
			user_uuid,
			plan_name,
			yookassa_payment_id,
			yookassa_payment_method_id,
			amount,
			currency,
			status,
			payment_kind,
			description,
			paid_at,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
		RETURNING *
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var res models.PaymentHistory

	err := s.db.QueryRowContext(
		ctx,
		query,
		payment.UserUUID,
		payment.PlanName,
		payment.YookassaPaymentID,
		payment.YookassaPaymentMethodID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentKind,
		payment.Description,
		payment.PaidAt,
		jsonbArg(payment.Metadata),
	).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.UserUUID,
		&res.PlanName,
		&res.YookassaPaymentID,
		&res.YookassaPaymentMethodID,
		&res.Amount,
		&res.Currency,
		&res.Status,
		&res.PaymentKind,
		&res.Description,
		&res.PaidAt,
		&res.Metadata,
	)

	if err != nil {
		logger.Error("Create_PaymentHistory req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &res, nil
}

func (s *PaymentStore) Get_PaymentHistoryListByUserUuid(ctx context.Context, userUuid string, limit int) ([]models.PaymentHistory, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT * FROM payment_history WHERE user_uuid = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := pgqx.QueryContext[models.PaymentHistory](ctx, s.db, query, userUuid, limit)
	if err != nil {
		logger.Error("Get_PaymentHistoryListByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return rows, nil
}

func (s *PaymentStore) Get_PaymentHistoryByYookassaPaymentID(ctx context.Context, paymentID string) (*models.PaymentHistory, error) {
	query := `
		SELECT * FROM payment_history
		WHERE yookassa_payment_id = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	payment, err := pgqx.QueryRowContext[models.PaymentHistory](ctx, s.db, query, paymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_PaymentHistoryByYookassaPaymentID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return payment, nil
}

func (s *PaymentStore) Get_UserSubscriptionBillingByUserUuid(ctx context.Context, userUuid string) (*models.UserSubscriptionBilling, error) {
	query := `
		SELECT
			COALESCE(auto_renew_enabled, FALSE) AS auto_renew_enabled,
			yookassa_payment_method_id,
			payment_method_type,
			payment_method_title,
			payment_method_saved_at
		FROM user_subscriptions
		WHERE user_uuid = $1 AND is_active = TRUE
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	billing, err := pgqx.QueryRowContext[models.UserSubscriptionBilling](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserSubscriptionBillingByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return billing, nil
}

func (s *PaymentStore) Get_SubscriptionsDueForAutoRenew(ctx context.Context) ([]models.UserSubscriptionRenewalDue, error) {
	query := `
		SELECT
			user_subscriptions.user_uuid,
			user_subscriptions.plan_name,
			user_subscriptions.yookassa_payment_method_id,
			subscriptions.amount,
			subscriptions.currency,
			subscriptions.duration_days
		FROM user_subscriptions
		JOIN subscriptions ON LOWER(subscriptions.plan_name) = LOWER(user_subscriptions.plan_name)
		WHERE user_subscriptions.is_active = TRUE
			AND user_subscriptions.auto_renew_enabled = TRUE
			AND user_subscriptions.yookassa_payment_method_id IS NOT NULL
			AND user_subscriptions.valid_to IS NOT NULL
			AND user_subscriptions.valid_to::date <= CURRENT_DATE
			AND LOWER(user_subscriptions.plan_name) <> 'free'
			AND NOT EXISTS (
				SELECT 1
				FROM payment_history
				WHERE payment_history.user_uuid = user_subscriptions.user_uuid
					AND payment_history.payment_kind = 'renewal'
					AND payment_history.status IN ('pending', 'waiting_for_capture', 'succeeded')
					AND payment_history.created_at::date = CURRENT_DATE
			)
	`

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	items, err := pgqx.QueryContext[models.UserSubscriptionRenewalDue](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_SubscriptionsDueForAutoRenew req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return items, nil
}

func (s *PaymentStore) Get_PaymentActiveCount(ctx context.Context, userUuid string) (uint32, error) {
	query := `
		SELECT COUNT(*) FROM payment_history
		WHERE user_uuid = $1
  			AND status IN ('pending', 'waiting_for_capture')
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count uint32

	err := s.db.QueryRowContext(ctx, query, userUuid).Scan(&count)
	if err != nil {
		logger.Error("Get_PaymentActiveCount req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return 0, err
	}

	return count, nil
}

func (s *PaymentStore) Update_PaymentHistoryStatus(ctx context.Context, yookassaPaymentID, status string, paidAt *time.Time) error {
	query := `
		UPDATE payment_history
		SET
			status = $2,
			paid_at = COALESCE($3, paid_at)
		WHERE yookassa_payment_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, yookassaPaymentID, status, paidAt)
	if err != nil {
		logger.Error("Update_PaymentHistoryStatus req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *PaymentStore) Update_UserSubscriptionAutoRenew(ctx context.Context, userUuid string, enabled bool) error {
	query := `
		UPDATE user_subscriptions
		SET auto_renew_enabled = $2
		WHERE user_uuid = $1 AND is_active = TRUE
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userUuid, enabled)
	if err != nil {
		logger.Error("Update_UserSubscriptionAutoRenew req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *PaymentStore) Update_YookassaPaymentMethod(ctx context.Context, userUuid string, paymentMethodID string, paymentMethodType string, paymentMethodTitle string) error {
	query := `
		UPDATE user_subscriptions
		SET
			yookassa_payment_method_id = $2,
			payment_method_type = $3,
			payment_method_title = $4,
			payment_method_saved_at = timezone('UTC', now()),
			auto_renew_enabled = TRUE
		WHERE user_uuid = $1 AND is_active = TRUE
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userUuid, paymentMethodID, paymentMethodType, paymentMethodTitle)
	if err != nil {
		logger.Error("SaveYookassaPaymentMethod req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *PaymentStore) Update_ActivateUserSubscriptionPlan(ctx context.Context, userUuid, planName string, durationDays int) error {
	if durationDays <= 0 {
		durationDays = 30
	}

	query := `
		UPDATE user_subscriptions
		SET
			plan_name = $2,
			valid_from = timezone('UTC', now()),
			valid_to = timezone('UTC', now()) + ($3::text || ' days')::interval,
			is_active = TRUE,
			updated_at = timezone('UTC', now())
		WHERE user_uuid = $1 AND is_active = TRUE
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userUuid, planName, durationDays)
	if err != nil {
		logger.Error("ActivateUserSubscriptionPlan req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *PaymentStore) Update_ClearYookassaPaymentMethod(ctx context.Context, userUuid string) error {
	query := `
		UPDATE user_subscriptions
		SET
			yookassa_payment_method_id = NULL,
			payment_method_type = NULL,
			payment_method_title = NULL,
			payment_method_saved_at = NULL,
			auto_renew_enabled = FALSE
		WHERE user_uuid = $1 AND is_active = TRUE
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userUuid)
	if err != nil {
		logger.Error("ClearYookassaPaymentMethod req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *PaymentStore) Delete_PaymentWithStatusPending(ctx context.Context) error {
	query := `
		DELETE FROM payment_history
		WHERE status IN ('pending', 'waiting_for_capture')
			AND created_at < now() - INTERVAL '10 minutes';
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		logger.Error("Delete_PaymentWithStatusPending req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func jsonbArg(raw []byte) interface{} {
	if len(raw) == 0 {
		return nil
	}

	return string(raw)
}
