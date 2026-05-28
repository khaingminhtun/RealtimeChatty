-- +goose Up
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
    last_contact_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_relationships_owner_id ON relationships(owner_id);

-- Attach the auto-update trigger
CREATE TRIGGER trigger_update_relationships_updated_at
    BEFORE UPDATE ON relationships
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_relationships_updated_at ON relationships;
DROP TABLE relationships;