-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tracked_subscription_share_proposal_items (
    proposal_id BIGINT NOT NULL REFERENCES tracked_subscription_share_proposals(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES tracked_subscription_members(id) ON DELETE CASCADE,
    share_percent NUMERIC(6, 3) NOT NULL,

    PRIMARY KEY (proposal_id, member_id),
    CONSTRAINT chk_tracked_subscription_share_proposal_items_share CHECK (share_percent > 0 AND share_percent <= 100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracked_subscription_share_proposal_items;
-- +goose StatementEnd
