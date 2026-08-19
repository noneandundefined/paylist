package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/pgqx"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create_UserCore(ctx context.Context, tx *sql.Tx, user *models.UserCore) error {
	query := `
		INSERT INTO user_cores (user_uuid, email, first_name, last_name, password)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, user.UserUUID, user.Email, user.FirstName, user.LastName, user.Password); err != nil {
		if strings.Contains(err.Error(), `user_cores_email_key`) {
			return httperr.Err_DuplicateEmail
		}

		logger.Error("Create_UserCore req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Create_UserSubscription(ctx context.Context, tx *sql.Tx, userUuid string) error {
	query := `
		INSERT INTO user_subscriptions (user_uuid)
		VALUES ($1)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, userUuid); err != nil {
		logger.Error("Create_UserSubscription req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Create_UserTransaction(ctx context.Context, tx *sql.Tx, user *models.UserTransaction) error {
	query := `
		INSERT INTO user_transactions (user_uuid, subscription_id, amount, currency)
		VALUES ($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, user.UserUUID, user.SubscriptionID, user.Amount, user.Currency); err != nil {
		logger.Error("Create_UserTransaction req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error) {
	query := `
		SELECT * FROM user_cores WHERE email = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserCore](ctx, s.db, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserCoreByEmail req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserCoreByUserUuid(ctx context.Context, userUuid string) (*models.UserCore, error) {
	query := `
		SELECT * FROM user_cores WHERE user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserCore](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserCoreByEmail req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Search_UserPublicProfilesByEmail(ctx context.Context, emailQuery, excludeUserUUID string, limit int) ([]models.UserPublicProfile, error) {
	query := `
		SELECT email, first_name, last_name, avatars
		FROM user_cores
		WHERE email = $1
			AND user_uuid <> $2
			AND email_confirmed = TRUE
		ORDER BY email ASC
		LIMIT $3
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	users, err := pgqx.QueryContext[models.UserPublicProfile](ctx, s.db, query, emailQuery, excludeUserUUID, limit)
	if err != nil {
		logger.Error("Search_UserPublicProfilesByEmail req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return users, nil
}

func (s *UserStore) Get_UserSubscriptionsByUserUuid(ctx context.Context, userUuid string) (*models.UserSubscription, error) {
	query := `
		SELECT * FROM user_subscriptions WHERE user_uuid = $1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserSubscription](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserSubscriptionsByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserLoginStateByUserUuid(ctx context.Context, userUuid string) (*models.UserLoginState, error) {
	query := `
		SELECT
			user_cores.id,
			user_cores.created_at,
			user_cores.email,
			user_cores.email_confirmed,
			user_cores.first_name,
			user_cores.last_name,
			user_cores.avatars,
			user_subscriptions.plan_name,
			user_subscriptions.valid_to,
			COALESCE(subscriptions.amount, 0) AS amount,
			COALESCE(subscriptions.currency, 'RUB') AS currency,
			COALESCE(subscriptions.notification_subscriptions, FALSE) AS notification_subscriptions,
			subscriptions.max_total_subscriptions,
			COALESCE(subscriptions.auto_find_subscriptions, FALSE) AS auto_find_subscriptions,
			COALESCE(user_cores.is_admin, FALSE) AS is_admin
		FROM user_cores
		LEFT JOIN user_subscriptions ON user_subscriptions.user_uuid = user_cores.user_uuid
		LEFT JOIN subscriptions ON LOWER(subscriptions.plan_name) = LOWER(user_subscriptions.plan_name)
		WHERE user_cores.user_uuid = $1 AND user_subscriptions.is_active = TRUE
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := pgqx.QueryRowContext[models.UserLoginState](ctx, s.db, query, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_UserLoginStateByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Get_UserPermissionsByUserUuid(ctx context.Context, userUuid string) (*models.UserPlanPermissions, error) {
	query := `
		SELECT
			user_subscriptions.plan_name,
			COALESCE(subscriptions.notification_subscriptions, FALSE),
			subscriptions.max_total_subscriptions,
			COALESCE(subscriptions.auto_find_subscriptions, FALSE)
		FROM user_subscriptions
		JOIN subscriptions ON LOWER(subscriptions.plan_name) = LOWER(user_subscriptions.plan_name)
		WHERE user_subscriptions.user_uuid = $1 AND user_subscriptions.is_active = TRUE
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var perms models.UserPlanPermissions
	var maxTotal sql.NullInt32

	err := s.db.QueryRowContext(ctx, query, userUuid).Scan(
		&perms.PlanName,
		&perms.NotificationSubscriptions,
		&maxTotal,
		&perms.AutoFindSubscriptions,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			freeLimit := 10
			return &models.UserPlanPermissions{
				PlanName:              "Free",
				MaxTotalSubscriptions: &freeLimit,
			}, nil
		}

		logger.Error("Get_UserPermissionsByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if maxTotal.Valid {
		v := int(maxTotal.Int32)
		perms.MaxTotalSubscriptions = &v
	}

	return &perms, nil
}

func (s *UserStore) Update_UserSubscriptionResetExpired(ctx context.Context) error {
	query := `
		UPDATE user_subscriptions
		SET
			plan_name = 'Free',
			valid_to = NULL,
			updated_at = NOW()
		WHERE is_active = TRUE
			AND valid_to IS NOT NULL
			AND valid_to <= NOW()
			AND LOWER(plan_name) <> 'free'
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		logger.Error("Update_UserSubscriptionResetExpired req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserEmailConfirmedByUid(ctx context.Context, userUuid string, confirmed bool) error {
	query := `
		UPDATE user_cores SET email_confirmed = $1 WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, confirmed, userUuid); err != nil {
		logger.Error("Update_UserEmailConfirmedByUid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserProfile(ctx context.Context, userUuid string, firstName, lastName *string) error {
	query := `
		UPDATE user_cores SET first_name = $1, last_name = $2 WHERE user_uuid = $3
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, firstName, lastName, userUuid)
	if err != nil {
		logger.Error("Update_UserProfile req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserAvatar(ctx context.Context, userUuid string, avatarURL string) error {
	query := `
		UPDATE user_cores SET avatars = $1 WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, avatarURL, userUuid); err != nil {
		logger.Error("Update_UserAvatar req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserEmail(ctx context.Context, userUuid, email string) error {
	query := `
		UPDATE user_cores
		SET email = $1, email_confirmed = TRUE
		WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, email, userUuid); err != nil {
		if strings.Contains(err.Error(), `user_cores_email_key`) {
			return httperr.Err_DuplicateEmail
		}

		logger.Error("Update_UserEmail req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Update_UserPassword(ctx context.Context, userUuid, passwordHash string) error {
	query := `
		UPDATE user_cores
		SET password = $1, email_confirmed = TRUE
		WHERE user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, passwordHash, userUuid); err != nil {
		logger.Error("Update_UserPassword req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Delete_UserByUuid(ctx context.Context, userUuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `DELETE FROM user_cores WHERE user_uuid = $1`

	if _, err := s.db.ExecContext(ctx, query, userUuid); err != nil {
		logger.Error("Delete_UserByUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}
