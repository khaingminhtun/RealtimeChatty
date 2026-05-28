CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    session_id BIGINT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,

    is_revoked BOOLEAN DEFAULT FALSE,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW()
);