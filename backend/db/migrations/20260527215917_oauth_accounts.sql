-- +goose Up
CREATE TABLE oauth_accounts (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,

    access_token TEXT,
    refresh_token TEXT,

    token_expires_at TIMESTAMPTZ,

    scopes TEXT,
    id_token TEXT,

    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(provider, provider_user_id),
    UNIQUE(user_id, provider)
);

CREATE INDEX idx_oauth_accounts_user_id
ON oauth_accounts(user_id);

-- +goose Down
DROP TABLE oauth_accounts;