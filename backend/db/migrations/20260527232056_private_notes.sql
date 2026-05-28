-- +goose Up
CREATE TABLE private_notes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    relationship_id BIGINT NOT NULL REFERENCES relationships(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_private_notes_user_relationship ON private_notes(user_id, relationship_id);

-- Attach the auto-update trigger
CREATE TRIGGER trigger_update_private_notes_updated_at
    BEFORE UPDATE ON private_notes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_private_notes_updated_at ON private_notes;
DROP TABLE private_notes;