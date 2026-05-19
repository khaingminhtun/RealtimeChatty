-- +goose Up
CREATE TABLE user_auth (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    password_hash TEXT NOT NULL,

    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret TEXT,

    failed_attempts INTEGER DEFAULT 0,

    locked_until TIMESTAMP,

    password_changed_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE user_auth;