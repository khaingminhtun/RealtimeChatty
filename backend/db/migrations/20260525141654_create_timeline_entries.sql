-- +goose Up
CREATE TABLE timeline_entries (
    id BIGSERIAL PRIMARY KEY,

    relationship_id BIGINT NOT NULL REFERENCES relationships(id) ON DELETE CASCADE,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    entry_type VARCHAR(50) NOT NULL,

    title VARCHAR(255) NOT NULL,

    body TEXT,

    mood_tag VARCHAR(50),

    entry_date DATE NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_timeline_relationship_id
ON timeline_entries(relationship_id);

CREATE INDEX idx_timeline_user_id
ON timeline_entries(user_id);

CREATE INDEX idx_timeline_entry_date
ON timeline_entries(entry_date DESC);

CREATE INDEX idx_timeline_entry_type
ON timeline_entries(entry_type);

-- +goose Down
DROP TABLE timeline_entries;