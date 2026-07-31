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
    
    drift_threshold_days INTEGER NOT NULL DEFAULT 30,

    drift_status VARCHAR(50) NOT NULL DEFAULT 'healthy',
    warmth_score INTEGER DEFAULT 100,
    last_reminder_sent_at TIMESTAMPTZ,
    next_contact_at TIMESTAMPTZ,
    
    last_contact_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);