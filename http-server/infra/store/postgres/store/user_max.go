package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"paylist.server/infra/logger"
)

func (s *UserStore) Upsert_UserMax(ctx context.Context, userUuid string, userID int64, username, language string) error {
	query := `
		INSERT INTO user_settings (user_uuid, max_user_id, max_username, max_language)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_uuid) DO UPDATE
		SET max_user_id = EXCLUDED.max_user_id,
		    max_username = EXCLUDED.max_username,
		    max_language = EXCLUDED.max_language,
		    updated_at = timezone('UTC', now())
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid, userID, username, language); err != nil {
		logger.Error("Upsert_UserMax req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Clear_UserMax(ctx context.Context, userUuid string) error {
	query := `
		UPDATE user_settings
		SET max_user_id = NULL,
		    max_username = NULL,
		    max_language = NULL,
		    updated_at = timezone('UTC', now())
		WHERE user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid); err != nil {
		logger.Error("Clear_UserMax req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Get_UserUuidByMaxUserID(ctx context.Context, userID int64) (string, error) {
	query := `
		SELECT user_uuid
		FROM user_settings
		WHERE max_user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userUuid string

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		logger.Error("Get_UserUuidByMaxUserID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return "", err
	}

	return userUuid, nil
}
