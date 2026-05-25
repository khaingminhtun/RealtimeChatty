CREATE TABLE relationships (
    id BIGSERIAL PRIMARY KEY,

    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name VARCHAR(120) NOT NULL,

    type VARCHAR(50),

    how_we_met TEXT,

    birthday DATE,

    location VARCHAR(255),

    tags TEXT[] DEFAULT '{}',

    last_contact_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_relationships_owner_id
ON relationships(owner_id);

CREATE INDEX idx_relationships_type
ON relationships(type);

CREATE INDEX idx_relationships_last_contact
ON relationships(last_contact_at);