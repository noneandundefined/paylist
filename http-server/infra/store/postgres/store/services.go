package store

import (
	"context"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/pgqx"
)

func (s *TrackedSubscriptionStore) Get_Services(ctx context.Context) ([]models.Service, error) {
	query := `
		SELECT id, created_at, slug, name, category, aliases
		FROM services
		ORDER BY name ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	services, err := pgqx.QueryContext[models.Service](ctx, s.db, query)
	if err != nil {
		logger.Error("Get_Services req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if services == nil {
		return []models.Service{}, nil
	}

	return services, nil
}
