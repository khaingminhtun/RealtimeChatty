-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,

    display_name VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,

    status VARCHAR(20) DEFAULT 'offline',

    last_seen_at TIMESTAMP,

    privacy_setting VARCHAR(20) DEFAULT 'public',

    notifications_enabled BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;