CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,

    display_name VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,

    timezone VARCHAR(100) DEFAULT 'UTC',

    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    privacy_setting VARCHAR(20) DEFAULT 'public',

    notifications_enabled BOOLEAN DEFAULT TRUE,

    push_token TEXT,

    last_seen_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);