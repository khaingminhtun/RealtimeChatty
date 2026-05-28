CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,

    relationship_id BIGINT NOT NULL REFERENCES relationships(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    channel VARCHAR(50) NOT NULL, -- e.g., 'call', 'message', 'email', 'in_person'
    note TEXT,
    
    contacted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
