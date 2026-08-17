package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/pgqx"
)

func (s *TrackedSubscriptionStore) Get_MembersBySubscriptionID(ctx context.Context, subscriptionID uint64) ([]models.TrackedSubscriptionMember, error) {
	query := `
		SELECT
			tracked_subscription_members.id,
			tracked_subscription_members.created_at,
			tracked_subscription_members.updated_at,
			tracked_subscription_members.tracked_subscription_id,
			tracked_subscription_members.user_uuid,
			tracked_subscription_members.email,
			tracked_subscription_members.role,
			tracked_subscription_members.share_percent,
			tracked_subscription_members.notification,
			tracked_subscription_members.include_in_analytics,
			tracked_subscription_members.status,
			tracked_subscription_members.invite_token,
			tracked_subscription_members.invite_expires_at,
			user_cores.first_name,
			user_cores.last_name,
			user_cores.avatars
		FROM tracked_subscription_members
		LEFT JOIN user_cores ON user_cores.user_uuid = tracked_subscription_members.user_uuid
			OR (
				tracked_subscription_members.user_uuid IS NULL
				AND user_cores.email = tracked_subscription_members.email
			)
		WHERE tracked_subscription_members.tracked_subscription_id = $1
		ORDER BY (tracked_subscription_members.role = 'owner') DESC, (tracked_subscription_members.role = 'observer') ASC, tracked_subscription_members.created_at ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	members, err := pgqx.QueryContext[models.TrackedSubscriptionMember](ctx, s.db, query, subscriptionID)
	if err != nil {
		logger.Error("Get_MembersBySubscriptionID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return members, nil
}

func (s *TrackedSubscriptionStore) Get_AcceptedMember(ctx context.Context, subscriptionID uint64, userUUID string) (*models.TrackedSubscriptionMember, error) {
	query := `
		SELECT
			tracked_subscription_members.id,
			tracked_subscription_members.created_at,
			tracked_subscription_members.updated_at,
			tracked_subscription_members.tracked_subscription_id,
			tracked_subscription_members.user_uuid,
			tracked_subscription_members.email,
			tracked_subscription_members.role,
			tracked_subscription_members.share_percent,
			tracked_subscription_members.notification,
			tracked_subscription_members.include_in_analytics,
			tracked_subscription_members.status,
			tracked_subscription_members.invite_token,
			tracked_subscription_members.invite_expires_at,
			user_cores.first_name,
			user_cores.last_name,
			user_cores.avatars
		FROM tracked_subscription_members
		LEFT JOIN user_cores ON user_cores.user_uuid = tracked_subscription_members.user_uuid
			OR (
				tracked_subscription_members.user_uuid IS NULL
				AND user_cores.email = tracked_subscription_members.email
			)
		WHERE tracked_subscription_members.tracked_subscription_id = $1
			AND tracked_subscription_members.user_uuid = $2
			AND tracked_subscription_members.status = 'accepted'
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	member, err := pgqx.QueryRowContext[models.TrackedSubscriptionMember](ctx, s.db, query, subscriptionID, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_AcceptedMember req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return member, nil
}

func (s *TrackedSubscriptionStore) Get_MemberByID(ctx context.Context, subscriptionID, memberID uint64) (*models.TrackedSubscriptionMember, error) {
	query := `
		SELECT
			id, created_at, updated_at, tracked_subscription_id,
			user_uuid, email, role, share_percent, notification,
			include_in_analytics, status, invite_token, invite_expires_at
		FROM tracked_subscription_members
		WHERE id = $1 AND tracked_subscription_id = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	member, err := pgqx.QueryRowContext[models.TrackedSubscriptionMember](ctx, s.db, query, memberID, subscriptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_MemberByID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return member, nil
}

func (s *TrackedSubscriptionStore) Get_MemberByEmail(ctx context.Context, subscriptionID uint64, email string) (*models.TrackedSubscriptionMember, error) {
	query := `
		SELECT
			id, created_at, updated_at, tracked_subscription_id,
			user_uuid, email, role, share_percent, notification,
			include_in_analytics, status, invite_token, invite_expires_at
		FROM tracked_subscription_members
		WHERE tracked_subscription_id = $1 AND email = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	member, err := pgqx.QueryRowContext[models.TrackedSubscriptionMember](ctx, s.db, query, subscriptionID, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_MemberByEmail req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return member, nil
}

func (s *TrackedSubscriptionStore) Count_ActiveMembers(ctx context.Context, subscriptionID uint64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM tracked_subscription_members
		WHERE tracked_subscription_id = $1
			AND status IN ('pending', 'accepted')
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int
	if err := s.db.QueryRowContext(ctx, query, subscriptionID).Scan(&count); err != nil {
		logger.Error("Count_ActiveMembers req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return 0, err
	}

	return count, nil
}

func (s *TrackedSubscriptionStore) Create_MemberInvite(ctx context.Context, member *models.TrackedSubscriptionMember) error {
	if member.Role == "" {
		member.Role = "member"
	}

	query := `
		INSERT INTO tracked_subscription_members (
			tracked_subscription_id, email, role, share_percent,
			notification, include_in_analytics, status, invite_token, invite_expires_at
		)
		VALUES ($1, $2, $3, $4, FALSE, $5, 'pending', $6, $7)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Create_MemberInvite req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if member.SharePercent != 0 {
		reserve, err := tx.ExecContext(
			ctx,
			`
				UPDATE tracked_subscription_members
				SET share_percent = share_percent - $1
				WHERE tracked_subscription_id = $2
					AND role = 'owner'
					AND status = 'accepted'
					AND share_percent - $1 >= 0
			`,
			member.SharePercent,
			member.TrackedSubscriptionID,
		)
		if err != nil {
			logger.Error("Create_MemberInvite req={%s}: Failed to reserve owner share: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		reserved, err := reserve.RowsAffected()
		if err != nil {
			return err
		}

		if reserved == 0 {
			return httperr.Err_NotUpdated
		}
	}

	if err := tx.QueryRowContext(
		ctx,
		query,
		member.TrackedSubscriptionID,
		member.Email,
		member.Role,
		member.SharePercent,
		member.IncludeInAnalytics,
		member.InviteToken,
		member.InviteExpiresAt,
	).Scan(&member.ID); err != nil {
		logger.Error("Create_MemberInvite req={%s}: Failed to insert invite: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return tx.Commit()
}

func (s *TrackedSubscriptionStore) Refresh_MemberInvite(ctx context.Context, member *models.TrackedSubscriptionMember, newShare float64) error {
	if member.Role == "" {
		member.Role = "member"
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Refresh_MemberInvite req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	delta := newShare - member.SharePercent
	if delta != 0 {
		result, err := tx.ExecContext(
			ctx,
			`
				UPDATE tracked_subscription_members
				SET share_percent = share_percent - $1
				WHERE tracked_subscription_id = $2
					AND role = 'owner'
					AND status = 'accepted'
					AND share_percent - $1 >= 0
			`,
			delta,
			member.TrackedSubscriptionID,
		)
		if err != nil {
			logger.Error("Refresh_MemberInvite req={%s}: Failed to adjust owner share: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			return httperr.Err_NotUpdated
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			UPDATE tracked_subscription_members
			SET share_percent = $1,
				role = $2,
				include_in_analytics = $3,
				status = 'pending',
				invite_token = $4,
				invite_expires_at = $5
			WHERE id = $6
		`,
		newShare,
		member.Role,
		member.Role != "observer",
		member.InviteToken,
		member.InviteExpiresAt,
		member.ID,
	); err != nil {
		logger.Error("Refresh_MemberInvite req={%s}: Failed to refresh invite: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	member.SharePercent = newShare
	member.Status = "pending"
	member.IncludeInAnalytics = member.Role != "observer"

	return tx.Commit()
}

func (s *TrackedSubscriptionStore) Get_InviteByToken(ctx context.Context, token string) (*models.TrackedSubscriptionInvitePreview, *models.TrackedSubscriptionMember, error) {
	query := `
		SELECT
			tracked_subscriptions.id AS subscription_id,
			tracked_subscriptions.name AS subscription_name,
			COALESCE(
				NULLIF(TRIM(CONCAT(COALESCE(user_cores.first_name, ''), ' ', COALESCE(user_cores.last_name, ''))), ''),
				tracked_subscription_owner_members.email
			) AS owner_name,
			tracked_subscription_members.email,
			tracked_subscription_members.share_percent,
			tracked_subscription_members.status,
			tracked_subscription_members.invite_expires_at,
			tracked_subscription_members.id,
			tracked_subscription_members.created_at,
			tracked_subscription_members.updated_at,
			tracked_subscription_members.tracked_subscription_id,
			tracked_subscription_members.user_uuid,
			tracked_subscription_members.role,
			tracked_subscription_members.notification,
			tracked_subscription_members.include_in_analytics,
			tracked_subscription_members.invite_token
		FROM tracked_subscription_members
		JOIN tracked_subscriptions ON tracked_subscriptions.id = tracked_subscription_members.tracked_subscription_id
		JOIN tracked_subscription_members AS tracked_subscription_owner_members
			ON tracked_subscription_owner_members.tracked_subscription_id = tracked_subscriptions.id
			AND tracked_subscription_owner_members.role = 'owner'
		LEFT JOIN user_cores ON user_cores.user_uuid = tracked_subscription_owner_members.user_uuid
		WHERE tracked_subscription_members.invite_token = $1
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, token)

	var preview models.TrackedSubscriptionInvitePreview
	var member models.TrackedSubscriptionMember

	if err := row.Scan(
		&preview.SubscriptionID,
		&preview.SubscriptionName,
		&preview.OwnerName,
		&preview.Email,
		&preview.SharePercent,
		&preview.Status,
		&preview.InviteExpiresAt,
		&member.ID,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.TrackedSubscriptionID,
		&member.UserUUID,
		&member.Role,
		&member.Notification,
		&member.IncludeInAnalytics,
		&member.InviteToken,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}

		logger.Error("Get_InviteByToken req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, nil, err
	}

	member.Email = preview.Email
	member.SharePercent = preview.SharePercent
	member.Status = preview.Status
	member.InviteExpiresAt = preview.InviteExpiresAt
	preview.Role = member.Role

	return &preview, &member, nil
}

func (s *TrackedSubscriptionStore) Accept_MemberInvite(ctx context.Context, memberID uint64, userUUID string) error {
	query := `
		UPDATE tracked_subscription_members
		SET user_uuid = $1,
			status = 'accepted',
			invite_token = NULL,
			invite_expires_at = NULL
		WHERE id = $2
			AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, userUUID, memberID)
	if err != nil {
		logger.Error("Accept_MemberInvite req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	affected, err := upd.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return httperr.Err_NotUpdated
	}

	return nil
}

func (s *TrackedSubscriptionStore) Delete_MemberAndReturnShare(ctx context.Context, subscriptionID, memberID uint64, sharePercent float64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Delete_MemberAndReturnShare req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(
		ctx,
		`
			UPDATE tracked_subscription_share_proposals
			SET status = 'cancelled'
			WHERE tracked_subscription_id = $1 AND status = 'pending'
		`,
		subscriptionID,
	); err != nil {
		logger.Error("Delete_MemberAndReturnShare req={%s}: Failed to cancel proposals: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			UPDATE tracked_subscription_members
			SET share_percent = share_percent + $1
			WHERE tracked_subscription_id = $2
				AND role = 'owner'
				AND status = 'accepted'
		`,
		sharePercent,
		subscriptionID,
	); err != nil {
		logger.Error("Delete_MemberAndReturnShare req={%s}: Failed to return share: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	del, err := tx.ExecContext(
		ctx,
		`DELETE FROM tracked_subscription_members WHERE id = $1 AND tracked_subscription_id = $2 AND role <> 'owner'`,
		memberID,
		subscriptionID,
	)
	if err != nil {
		logger.Error("Delete_MemberAndReturnShare req={%s}: Failed to delete member: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	affected, err := del.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return httperr.Err_NotDeleted
	}

	return tx.Commit()
}
