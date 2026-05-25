-- +goose Up
CREATE TABLE moods (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    timeline_entry_id BIGINT NOT NULL REFERENCES timeline_entries(id) ON DELETE CASCADE,

    label VARCHAR(50) NOT NULL,

    score INTEGER CHECK (score >= 1 AND score <= 10),

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_moods_timeline_entry_id
ON moods(timeline_entry_id);

CREATE INDEX idx_moods_user_id
ON moods(user_id);

-- +goose Down
DROP TABLE moods;