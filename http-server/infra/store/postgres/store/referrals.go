package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/pkg/referral"
)

type ReferralStore struct {
	db *sql.DB
}

func (s *ReferralStore) Ensure_ReferralCode(ctx context.Context, userUuid string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var existing sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT referral_code FROM user_cores WHERE user_uuid = $1`, userUuid).Scan(&existing)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		logger.Error("Ensure_ReferralCode req={%s}: Failed to load code: %s", ctx.Value("XREQID"), err.Error())
		return "", err
	}

	if code := strings.TrimSpace(existing.String); code != "" {
		return code, nil
	}

	for attempt := 0; attempt < 8; attempt++ {
		code, err := referral.GenerateCode()
		if err != nil {
			return "", err
		}

		result, err := s.db.ExecContext(ctx, `
			UPDATE user_cores
			SET referral_code = $2, updated_at = timezone('UTC', now())
			WHERE user_uuid = $1 AND (referral_code IS NULL OR referral_code = '')
		`, userUuid, code)
		if err != nil {
			if strings.Contains(err.Error(), "idx_user_cores_referral_code") {
				continue
			}

			logger.Error("Ensure_ReferralCode req={%s}: Failed to save code: %s", ctx.Value("XREQID"), err.Error())
			return "", err
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			err := s.db.QueryRowContext(ctx, `SELECT referral_code FROM user_cores WHERE user_uuid = $1`, userUuid).Scan(&existing)
			if err != nil {
				return "", err
			}

			if saved := strings.TrimSpace(existing.String); saved != "" {
				return saved, nil
			}

			continue
		}

		return code, nil
	}

	return "", errors.New("failed to allocate referral code")
}

func (s *ReferralStore) Attach_Referral(ctx context.Context, code, referredUuid string) error {
	code = referral.SanitizeCode(code)
	if code == "" || referredUuid == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var referrerUuid string
	err := s.db.QueryRowContext(ctx, `SELECT user_uuid FROM user_cores WHERE referral_code = $1`, code).Scan(&referrerUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		logger.Error("Attach_Referral req={%s}: Failed to find referrer: %s", ctx.Value("XREQID"), err.Error())
		return err
	}

	if referrerUuid == referredUuid {
		return nil
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO user_referrals (referrer_uuid, referred_uuid)
		VALUES ($1, $2)
		ON CONFLICT (referred_uuid) DO NOTHING
	`, referrerUuid, referredUuid)
	if err != nil {
		logger.Error("Attach_Referral req={%s}: Failed to insert referral: %s", ctx.Value("XREQID"), err.Error())
		return err
	}

	return nil
}

func (s *ReferralStore) Count_ConvertedReferrals(ctx context.Context, referrerUuid string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM user_referrals
		WHERE referrer_uuid = $1 AND converted_at IS NOT NULL
	`, referrerUuid).Scan(&count)
	if err != nil {
		logger.Error("Count_ConvertedReferrals req={%s}: Failed to count: %s", ctx.Value("XREQID"), err.Error())
		return 0, err
	}

	return count, nil
}

func (s *ReferralStore) Apply_PaymentConversion(ctx context.Context, referredUuid string) (string, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var referrerUuid string
	err := s.db.QueryRowContext(ctx, `
		UPDATE user_referrals
		SET converted_at = COALESCE(converted_at, timezone('UTC', now()))
		WHERE referred_uuid = $1
		RETURNING referrer_uuid
	`, referredUuid).Scan(&referrerUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 0, nil
		}

		logger.Error("Apply_PaymentConversion req={%s}: Failed to convert: %s", ctx.Value("XREQID"), err.Error())
		return "", 0, err
	}

	daysGranted := 0

	err = WithTx(s.db, ctx, func(tx *sql.Tx) error {
		var count int
		if err := tx.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM user_referrals
			WHERE referrer_uuid = $1 AND converted_at IS NOT NULL
		`, referrerUuid).Scan(&count); err != nil {
			return err
		}

		for _, rank := range referral.Ranks {
			if rank.RewardDays <= 0 || count < rank.MinCount {
				continue
			}

			result, err := tx.ExecContext(ctx, `
				INSERT INTO user_referral_rewards (user_uuid, rank, days_granted)
				VALUES ($1, $2, $3)
				ON CONFLICT (user_uuid, rank) DO NOTHING
			`, referrerUuid, rank.Level, rank.RewardDays)
			if err != nil {
				return err
			}

			rows, _ := result.RowsAffected()
			if rows == 0 {
				continue
			}

			if err := extendPremiumTx(ctx, tx, referrerUuid, rank.RewardDays); err != nil {
				return err
			}

			daysGranted += rank.RewardDays
		}

		return nil
	})
	if err != nil {
		logger.Error("Apply_PaymentConversion req={%s}: Failed to grant rewards: %s", ctx.Value("XREQID"), err.Error())
		return referrerUuid, 0, err
	}

	return referrerUuid, daysGranted, nil
}

func extendPremiumTx(ctx context.Context, tx *sql.Tx, userUuid string, days int) error {
	if days <= 0 {
		return nil
	}

	_, err := tx.ExecContext(ctx, `
		UPDATE user_subscriptions
		SET
			plan_name = $2,
			valid_from = CASE
				WHEN plan_name = $2 AND valid_to IS NOT NULL AND valid_to > timezone('UTC', now()) THEN valid_from
				ELSE timezone('UTC', now())
			END,
			valid_to = GREATEST(COALESCE(valid_to, timezone('UTC', now())), timezone('UTC', now())) + ($3::text || ' days')::interval,
			is_active = TRUE,
			updated_at = timezone('UTC', now())
		WHERE user_uuid = $1 AND is_active = TRUE
	`, userUuid, referral.PremiumPlan, days)

	return err
}
