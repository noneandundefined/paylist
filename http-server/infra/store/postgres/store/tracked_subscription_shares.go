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

func (s *TrackedSubscriptionStore) Get_PendingShareProposal(ctx context.Context, subscriptionID uint64) (*models.TrackedSubscriptionShareProposal, error) {
	query := `
		SELECT id, created_at, tracked_subscription_id, proposed_by_user_uuid, status
		FROM tracked_subscription_share_proposals
		WHERE tracked_subscription_id = $1 AND status = 'pending'
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	proposal, err := pgqx.QueryRowContext[models.TrackedSubscriptionShareProposal](ctx, s.db, query, subscriptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_PendingShareProposal req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return proposal, nil
}

func (s *TrackedSubscriptionStore) Get_ShareProposalByID(ctx context.Context, subscriptionID, proposalID uint64) (*models.TrackedSubscriptionShareProposal, error) {
	query := `
		SELECT id, created_at, tracked_subscription_id, proposed_by_user_uuid, status
		FROM tracked_subscription_share_proposals
		WHERE id = $1 AND tracked_subscription_id = $2
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	proposal, err := pgqx.QueryRowContext[models.TrackedSubscriptionShareProposal](ctx, s.db, query, proposalID, subscriptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		logger.Error("Get_ShareProposalByID req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return proposal, nil
}

func (s *TrackedSubscriptionStore) Get_ShareProposalItems(ctx context.Context, proposalID uint64) ([]models.TrackedSubscriptionShareProposalItem, error) {
	query := `
		SELECT proposal_id, member_id, share_percent
		FROM tracked_subscription_share_proposal_items
		WHERE proposal_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	items, err := pgqx.QueryContext[models.TrackedSubscriptionShareProposalItem](ctx, s.db, query, proposalID)
	if err != nil {
		logger.Error("Get_ShareProposalItems req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return items, nil
}

func (s *TrackedSubscriptionStore) Get_ShareProposalVotes(ctx context.Context, proposalID uint64) ([]models.TrackedSubscriptionShareVote, error) {
	query := `
		SELECT proposal_id, user_uuid, accepted, created_at
		FROM tracked_subscription_share_votes
		WHERE proposal_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	votes, err := pgqx.QueryContext[models.TrackedSubscriptionShareVote](ctx, s.db, query, proposalID)
	if err != nil {
		logger.Error("Get_ShareProposalVotes req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	return votes, nil
}

func (s *TrackedSubscriptionStore) Create_ShareProposal(ctx context.Context, subscriptionID uint64, proposedBy string, items []models.TrackedSubscriptionShareProposalItem) (*models.TrackedSubscriptionShareProposal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Create_ShareProposal req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
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
		logger.Error("Create_ShareProposal req={%s}: Failed to cancel previous proposal: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	proposal := &models.TrackedSubscriptionShareProposal{
		TrackedSubscriptionID: subscriptionID,
		ProposedByUserUUID:    proposedBy,
		Status:                "pending",
	}

	if err := tx.QueryRowContext(
		ctx,
		`
			INSERT INTO tracked_subscription_share_proposals (
				tracked_subscription_id, proposed_by_user_uuid, status
			)
			VALUES ($1, $2, 'pending')
			RETURNING id, created_at
		`,
		subscriptionID,
		proposedBy,
	).Scan(&proposal.ID, &proposal.CreatedAt); err != nil {
		logger.Error("Create_ShareProposal req={%s}: Failed to insert proposal: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	for _, item := range items {
		if _, err := tx.ExecContext(
			ctx,
			`
				INSERT INTO tracked_subscription_share_proposal_items (proposal_id, member_id, share_percent)
				VALUES ($1, $2, $3)
			`,
			proposal.ID,
			item.MemberID,
			item.SharePercent,
		); err != nil {
			logger.Error("Create_ShareProposal req={%s}: Failed to insert proposal item: %s", ctx.Value("XREQID").(string), err.Error())
			return nil, err
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO tracked_subscription_share_votes (proposal_id, user_uuid, accepted)
			VALUES ($1, $2, TRUE)
		`,
		proposal.ID,
		proposedBy,
	); err != nil {
		logger.Error("Create_ShareProposal req={%s}: Failed to insert proposer vote: %s", ctx.Value("XREQID").(string), err.Error())
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return proposal, nil
}

func (s *TrackedSubscriptionStore) Upsert_ShareVote(ctx context.Context, proposalID uint64, userUUID string, accepted bool) error {
	query := `
		INSERT INTO tracked_subscription_share_votes (proposal_id, user_uuid, accepted)
		VALUES ($1, $2, $3)
		ON CONFLICT (proposal_id, user_uuid)
		DO UPDATE SET accepted = EXCLUDED.accepted, created_at = timezone('UTC', now())
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, proposalID, userUUID, accepted); err != nil {
		logger.Error("Upsert_ShareVote req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	return nil
}

func (s *TrackedSubscriptionStore) Reject_ShareProposal(ctx context.Context, proposalID uint64) error {
	query := `
		UPDATE tracked_subscription_share_proposals
		SET status = 'rejected'
		WHERE id = $1 AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	upd, err := s.db.ExecContext(ctx, query, proposalID)
	if err != nil {
		logger.Error("Reject_ShareProposal req={%s}: Failed to exec sql: %s", ctx.Value("XREQID").(string), err.Error())
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

func (s *TrackedSubscriptionStore) Apply_ShareProposal(ctx context.Context, proposalID uint64, items []models.TrackedSubscriptionShareProposalItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Apply_ShareProposal req={%s}: Failed to begin tx: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for _, item := range items {
		if _, err := tx.ExecContext(
			ctx,
			`
				UPDATE tracked_subscription_members
				SET share_percent = $1
				WHERE id = $2 AND status = 'accepted'
			`,
			item.SharePercent,
			item.MemberID,
		); err != nil {
			logger.Error("Apply_ShareProposal req={%s}: Failed to apply share: %s", ctx.Value("XREQID").(string), err.Error())
			return err
		}
	}

	upd, err := tx.ExecContext(
		ctx,
		`
			UPDATE tracked_subscription_share_proposals
			SET status = 'applied'
			WHERE id = $1 AND status = 'pending'
		`,
		proposalID,
	)
	if err != nil {
		logger.Error("Apply_ShareProposal req={%s}: Failed to mark applied: %s", ctx.Value("XREQID").(string), err.Error())
		return err
	}

	affected, err := upd.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return httperr.Err_NotUpdated
	}

	return tx.Commit()
}
