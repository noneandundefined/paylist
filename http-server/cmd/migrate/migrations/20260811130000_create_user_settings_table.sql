-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_settings (
    user_uuid VARCHAR(255) PRIMARY KEY REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (timezone('UTC', now())),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (timezone('UTC', now())),
    display_currency CHAR(3) NULL,
    country CHAR(2) NULL,

    telegram_chat_id BIGINT,
    telegram_username VARCHAR(64),
    telegram_language CHAR(2),

    max_user_id BIGINT,
    max_username VARCHAR(64),
    max_language CHAR(2),

    CONSTRAINT chk_user_settings_display_currency CHECK (display_currency IS NULL OR display_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT chk_user_settings_country CHECK (country IS NULL OR country ~ '^[A-Z]{2}$')
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_settings_telegram_chat_id
	ON user_settings (telegram_chat_id)
	WHERE telegram_chat_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_settings_max_user_id
	ON user_settings (max_user_id)
	WHERE max_user_id IS NOT NULL;

CREATE OR REPLACE FUNCTION set_user_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = timezone('UTC', now());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS set_user_settings_updated_at_trigger ON user_settings;
CREATE TRIGGER set_user_settings_updated_at_trigger
    BEFORE UPDATE ON user_settings
    FOR EACH ROW
    EXECUTE FUNCTION set_user_settings_updated_at();

ALTER TABLE subscription_categories
    ADD COLUMN IF NOT EXISTS user_uuid VARCHAR(255) NULL REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS label TEXT NULL;

ALTER TABLE subscription_categories DROP CONSTRAINT IF EXISTS uq_subscription_categories_slug;

CREATE UNIQUE INDEX IF NOT EXISTS uq_subscription_categories_global_slug
    ON subscription_categories (slug)
    WHERE user_uuid IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_subscription_categories_user_slug
    ON subscription_categories (user_uuid, slug)
    WHERE user_uuid IS NOT NULL;

ALTER TABLE subscription_categories
    DROP CONSTRAINT IF EXISTS chk_subscription_categories_user_label;

ALTER TABLE subscription_categories
    ADD CONSTRAINT chk_subscription_categories_user_label
    CHECK (user_uuid IS NULL OR label IS NOT NULL);

COMMENT ON TABLE user_settings IS 'Per-user preferences (Premium display currency, etc.)';
COMMENT ON COLUMN user_settings.display_currency IS 'Preferred currency for spending summary (Premium)';
COMMENT ON COLUMN user_settings.country IS 'User country ISO 3166-1 alpha-2 for inflation analytics (Premium)';
COMMENT ON COLUMN subscription_categories.user_uuid IS 'Owner for custom categories; NULL for global catalog';
COMMENT ON COLUMN subscription_categories.label IS 'Display name for custom user categories';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE subscription_categories DROP CONSTRAINT IF EXISTS chk_subscription_categories_user_label;
DROP INDEX IF EXISTS uq_subscription_categories_user_slug;
DROP INDEX IF EXISTS uq_subscription_categories_global_slug;

ALTER TABLE subscription_categories
    DROP COLUMN IF EXISTS label,
    DROP COLUMN IF EXISTS user_uuid;

ALTER TABLE subscription_categories
    ADD CONSTRAINT uq_subscription_categories_slug UNIQUE (slug);

DROP TRIGGER IF EXISTS set_user_settings_updated_at_trigger ON user_settings;
DROP FUNCTION IF EXISTS set_user_settings_updated_at();
DROP TABLE IF EXISTS user_settings;
-- +goose StatementEnd
