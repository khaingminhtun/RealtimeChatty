-- +goose Up
CREATE TABLE email_verifications (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    email VARCHAR(255) NOT NULL,

    otp_hash TEXT NOT NULL,

    attempts INTEGER DEFAULT 0,

    is_used BOOLEAN DEFAULT FALSE,

    expires_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_email_verifications_user_id
ON email_verifications(user_id);