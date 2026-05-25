-- +goose Up
CREATE TABLE private_notes (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    relationship_id BIGINT NOT NULL REFERENCES relationships(id) ON DELETE CASCADE,

    content TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_private_notes_user_relationship
ON private_notes(user_id, relationship_id);

-- +goose Down
DROP TABLE private_notes;