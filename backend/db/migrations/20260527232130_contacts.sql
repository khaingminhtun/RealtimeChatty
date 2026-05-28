-- +goose Up
CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    relationship_id BIGINT NOT NULL REFERENCES relationships(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL, -- e.g., 'call', 'message', 'email', 'in_person'
    note TEXT,
    contacted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for lookup filtering and sorting by recency
CREATE INDEX idx_contacts_relationship_recency 
ON contacts(relationship_id, contacted_at DESC);

CREATE INDEX idx_contacts_user_id ON contacts(user_id);

-- +goose Down
DROP TABLE contacts;