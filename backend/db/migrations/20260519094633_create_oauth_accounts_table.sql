-- +goose Up
CREATE TABLE oauth_accounts (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    provider VARCHAR(50) NOT NULL,

    provider_user_id VARCHAR(255) NOT NULL,

    access_token TEXT,
    refresh_token TEXT,

    token_expires_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_oauth_accounts_user_id
ON oauth_accounts(user_id);

-- +goose Down
DROP TABLE oauth_accounts;