-- +goose Up

-- =========================================================
-- ADD NEXT CONTACT FIELD (COMPUTED BY BACKEND)
-- =========================================================

ALTER TABLE relationships
ADD COLUMN next_contact_at TIMESTAMPTZ;

-- Index for reminder worker queries
CREATE INDEX IF NOT EXISTS idx_relationships_next_contact_at
ON relationships(next_contact_at);


-- +goose Down

-- =========================================================
-- REMOVE NEXT CONTACT FIELD
-- =========================================================

DROP INDEX IF EXISTS idx_relationships_next_contact_at;

ALTER TABLE relationships
DROP COLUMN IF EXISTS next_contact_at;