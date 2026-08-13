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

func (s *TrackedSubscriptionStore) Get_AllSubscriptionCategories(ctx context.Context) ([]models.SubscriptionCategory, error) {
	return s.Get_SubscriptionCategoriesForUser(ctx, "")
}

func (s *TrackedSubscriptionStore) Get_SubscriptionCategoriesForUser(ctx context.Context, userUuid string) ([]models.SubscriptionCategory, error) {
	query := `
		SELECT id, created_at, slug, user_uuid, label
		FROM subscription_categories
		WHERE user_uuid IS NULL OR user_uuid = $1
		ORDER BY (user_uuid IS NULL) DESC, slug ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	categories, err := pgqx.QueryContext[models.SubscriptionCategory](ctx, s.db, query, userUuid)
	if err != nil {
		logger.Error("Get_SubscriptionCategoriesForUser req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return categories, nil
}

func (s *TrackedSubscriptionStore) Create_UserSubscriptionCategory(ctx context.Context, userUuid, slug, label string) (*models.SubscriptionCategory, error) {
	query := `
		INSERT INTO subscription_categories (slug, user_uuid, label)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, slug, user_uuid, label
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var category models.SubscriptionCategory

	err := s.db.QueryRowContext(ctx, query, slug, userUuid, label).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Slug,
		&category.UserUUID,
		&category.Label,
	)
	if err != nil {
		logger.Error("Create_UserSubscriptionCategory req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &category, nil
}

func (s *TrackedSubscriptionStore) Delete_UserSubscriptionCategory(ctx context.Context, userUuid string, categoryID uint64) error {
	query := `
		DELETE FROM subscription_categories
		WHERE id = $1 AND user_uuid = $2
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, categoryID, userUuid)
	if err != nil {
		logger.Error("Delete_UserSubscriptionCategory req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *TrackedSubscriptionStore) Get_SubscriptionCategoryBySlugForUser(ctx context.Context, userUuid, slug string) (*models.SubscriptionCategory, error) {
	query := `
		SELECT id, created_at, slug, user_uuid, label
		FROM subscription_categories
		WHERE slug = $1 AND (user_uuid IS NULL OR user_uuid = $2)
		ORDER BY (user_uuid IS NULL) ASC
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var category models.SubscriptionCategory

	err := s.db.QueryRowContext(ctx, query, slug, userUuid).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Slug,
		&category.UserUUID,
		&category.Label,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		logger.Error("Get_SubscriptionCategoryBySlugForUser req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return &category, nil
}
