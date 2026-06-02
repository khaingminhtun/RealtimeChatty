-- +goose Up
-- 1. Drop the old non-unique index to prevent duplicate index overhead
DROP INDEX IF EXISTS idx_private_notes_user_relationship;

-- 2. Add the unique constraint to enforce the 1-to-1 rule and enable SQLC ON CONFLICT upserts
ALTER TABLE private_notes 
ADD CONSTRAINT unique_user_relationship_note UNIQUE (user_id, relationship_id);


-- +goose Down
-- 1. Drop the unique constraint if rolling back
ALTER TABLE private_notes 
DROP CONSTRAINT IF EXISTS unique_user_relationship_note;

-- 2. Re-create the original non-unique index to restore original state
CREATE INDEX idx_private_notes_user_relationship ON private_notes(user_id, relationship_id);