package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"paylist.server/infra/logger"
)

func (s *UserStore) Upsert_UserTelegram(ctx context.Context, userUuid string, chatID int64, username, language string) error {
	query := `
		INSERT INTO user_settings (user_uuid, telegram_chat_id, telegram_username, telegram_language)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_uuid) DO UPDATE
		SET telegram_chat_id = EXCLUDED.telegram_chat_id,
		    telegram_username = EXCLUDED.telegram_username,
		    telegram_language = EXCLUDED.telegram_language,
		    updated_at = timezone('UTC', now())
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid, chatID, nullIfEmpty(username), language); err != nil {
		logger.Error("Upsert_UserTelegram req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Clear_UserTelegram(ctx context.Context, userUuid string) error {
	query := `
		UPDATE user_settings
		SET telegram_chat_id = NULL,
		    telegram_username = NULL,
		    telegram_language = NULL,
		    updated_at = timezone('UTC', now())
		WHERE user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid); err != nil {
		logger.Error("Clear_UserTelegram req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Get_UserUuidByTelegramChatID(ctx context.Context, chatID int64) (string, error) {
	query := `
		SELECT user_uuid
		FROM user_settings
		WHERE telegram_chat_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userUuid string

	err := s.db.QueryRowContext(ctx, query, chatID).Scan(&userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		logger.Error("Get_UserUuidByTelegramChatID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return "", err
	}

	return userUuid, nil
}
