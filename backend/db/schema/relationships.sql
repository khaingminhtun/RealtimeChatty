CREATE TABLE relationships (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    type VARCHAR(100), -- e.g., 'partner', 'friend', 'family'

    how_we_met TEXT,
    birthday DATE,
    location VARCHAR(255),

    avatar_url TEXT,
    tags TEXT[] DEFAULT '{}', -- Native PostgreSQL text array
    
    drift_interval_days VARCHAR(50) DEFAULT NULL,
    reminder_sent BOOLEAN NOT NULL DEFAULT FALSE,
    
    last_contact_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);