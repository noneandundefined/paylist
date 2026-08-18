package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
)

func (s *UserStore) Get_UserSettingsByUserUuid(ctx context.Context, userUuid string) (*models.UserSettings, error) {
	query := `
		SELECT user_uuid, created_at, updated_at, display_currency, country, telegram_chat_id, telegram_username, telegram_language, max_user_id, max_username, max_language
		FROM user_settings
		WHERE user_uuid = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var settings models.UserSettings

	err := s.db.QueryRowContext(ctx, query, userUuid).Scan(
		&settings.UserUUID,
		&settings.CreatedAt,
		&settings.UpdatedAt,
		&settings.DisplayCurrency,
		&settings.Country,
		&settings.TelegramChatID,
		&settings.TelegramUsername,
		&settings.TelegramLanguage,
		&settings.MaxUserID,
		&settings.MaxUsername,
		&settings.MaxLanguage,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.UserSettings{UserUUID: userUuid}, nil
		}

		logger.Error("Get_UserSettingsByUserUuid req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &settings, nil
}

func (s *UserStore) Upsert_UserDisplayCurrency(ctx context.Context, userUuid, currency string) error {
	query := `
		INSERT INTO user_settings (user_uuid, display_currency)
		VALUES ($1, $2)
		ON CONFLICT (user_uuid) DO UPDATE
		SET display_currency = EXCLUDED.display_currency,
		    updated_at = timezone('UTC', now())
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid, currency); err != nil {
		logger.Error("Upsert_UserDisplayCurrency req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *UserStore) Upsert_UserCountry(ctx context.Context, userUuid, country string) error {
	query := `
		INSERT INTO user_settings (user_uuid, country)
		VALUES ($1, $2)
		ON CONFLICT (user_uuid) DO UPDATE
		SET country = EXCLUDED.country,
		    updated_at = timezone('UTC', now())
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, userUuid, country); err != nil {
		logger.Error("Upsert_UserCountry req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return value
}
