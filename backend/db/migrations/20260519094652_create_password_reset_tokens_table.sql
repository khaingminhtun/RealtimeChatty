-- +goose Up
CREATE TABLE password_reset_tokens (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,

    is_used BOOLEAN DEFAULT FALSE,

    expires_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tokens_user_id
ON password_reset_tokens(user_id);

-- +goose Down
DROP TABLE password_reset_tokens;