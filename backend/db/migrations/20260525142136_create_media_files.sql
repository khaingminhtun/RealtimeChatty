-- +goose Up
CREATE TABLE media_files (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    relationship_id BIGINT REFERENCES relationships(id) ON DELETE CASCADE,

    timeline_entry_id BIGINT REFERENCES timeline_entries(id) ON DELETE CASCADE,

    storage_key TEXT NOT NULL,

    category VARCHAR(50),

    mime_type VARCHAR(100),

    size_bytes BIGINT,

    shared BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_media_user_id
ON media_files(user_id);

CREATE INDEX idx_media_relationship_id
ON media_files(relationship_id);

CREATE INDEX idx_media_timeline_entry_id
ON media_files(timeline_entry_id);

-- +goose Down
DROP TABLE media_files;