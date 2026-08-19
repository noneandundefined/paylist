package store

import (
	"context"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/pgqx"
)

func (s *UserStore) List_AdminMessageRecipients(ctx context.Context) ([]models.AdminMessageRecipient, error) {
	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_cores.first_name,
			user_cores.last_name,
			(user_settings.telegram_chat_id IS NOT NULL) AS telegram_connected,
			(user_settings.max_user_id IS NOT NULL) AS max_connected
		FROM user_cores
		LEFT JOIN user_settings ON user_settings.user_uuid = user_cores.user_uuid
		WHERE user_cores.email_confirmed = TRUE
		ORDER BY user_cores.email ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	users, err := pgqx.QueryContext[models.AdminMessageRecipient](ctx, s.db, query)
	if err != nil {
		logger.Error("List_AdminMessageRecipients req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if users == nil {
		return []models.AdminMessageRecipient{}, nil
	}

	return users, nil
}

func (s *UserStore) List_AdminMessageTargets(ctx context.Context, userUUID *string) ([]models.AdminMessageTarget, error) {
	query := `
		SELECT
			user_cores.user_uuid,
			user_cores.email,
			user_settings.telegram_chat_id,
			user_settings.max_user_id
		FROM user_cores
		LEFT JOIN user_settings ON user_settings.user_uuid = user_cores.user_uuid
		WHERE user_cores.email_confirmed = TRUE
			AND ($1::text IS NULL OR user_cores.user_uuid = $1)
		ORDER BY user_cores.email ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var uuid any
	if userUUID != nil && *userUUID != "" {
		uuid = *userUUID
	}

	users, err := pgqx.QueryContext[models.AdminMessageTarget](ctx, s.db, query, uuid)
	if err != nil {
		logger.Error("List_AdminMessageTargets req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if users == nil {
		return []models.AdminMessageTarget{}, nil
	}

	return users, nil
}
