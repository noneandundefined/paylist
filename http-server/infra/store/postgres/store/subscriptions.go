package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
)

type SubscriptionStore struct {
	db *sql.DB
}

func scanSubscription(rows *sql.Rows) (*models.Subscription, error) {
	var sub models.Subscription
	var descriptionRaw []byte
	var featuresRaw []byte
	var maxTotal sql.NullInt32

	if err := rows.Scan(
		&sub.ID,
		&sub.CreatedAt,
		&sub.PlanName,
		&sub.Amount,
		&sub.Currency,
		&sub.DurationDays,
		&maxTotal,
		&sub.NotificationSubscriptions,
		&sub.AutoFindSubscriptions,
		&descriptionRaw,
		&featuresRaw,
	); err != nil {
		return nil, err
	}

	if maxTotal.Valid {
		v := int(maxTotal.Int32)
		sub.MaxTotalSubscriptions = &v
	}

	if len(descriptionRaw) > 0 {
		if err := json.Unmarshal(descriptionRaw, &sub.Description); err != nil {
			return nil, err
		}
	}

	if len(featuresRaw) > 0 {
		if err := json.Unmarshal(featuresRaw, &sub.Features); err != nil {
			return nil, err
		}
	}

	return &sub, nil
}

func (s *SubscriptionStore) Get_SubscriptionByPlanName(ctx context.Context, planName string) (*models.Subscription, error) {
	query := `
		SELECT * FROM subscriptions
		WHERE LOWER(plan_name) = LOWER($1)
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, planName)
	if err != nil {
		logger.Error("Get_SubscriptionByPlanName req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	sub, err := scanSubscription(rows)
	if err != nil {
		logger.Error("Get_SubscriptionByPlanName req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return sub, nil
}

func (s *SubscriptionStore) Get_Subscriptions(ctx context.Context) ([]models.Subscription, error) {
	query := `SELECT * FROM subscriptions`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("Get_Subscriptions req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription

	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			logger.Error("Get_Subscriptions req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}

		subs = append(subs, *sub)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Get_Subscriptions req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return subs, nil
}
