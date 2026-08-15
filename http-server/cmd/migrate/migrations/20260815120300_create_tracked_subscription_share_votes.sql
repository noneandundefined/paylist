-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tracked_subscription_share_votes (
    proposal_id BIGINT NOT NULL REFERENCES tracked_subscription_share_proposals(id) ON DELETE CASCADE,
    user_uuid VARCHAR(255) NOT NULL REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    accepted BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (timezone('UTC', now())),

    PRIMARY KEY (proposal_id, user_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracked_subscription_share_votes;
-- +goose StatementEnd
