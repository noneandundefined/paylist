-- +goose Up
CREATE TABLE IF NOT EXISTS subscription_notification_log (
	id BIGSERIAL PRIMARY KEY,
	tracked_subscription_id BIGINT NOT NULL REFERENCES tracked_subscriptions (id) ON DELETE CASCADE,
	channel VARCHAR(20) NOT NULL DEFAULT 'telegram',
	notify_date DATE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC', now()),
	UNIQUE (tracked_subscription_id, channel, notify_date)
);

CREATE INDEX IF NOT EXISTS idx_subscription_notification_log_notify_date
	ON subscription_notification_log (notify_date);

-- +goose Down
DROP TABLE IF EXISTS subscription_notification_log;
