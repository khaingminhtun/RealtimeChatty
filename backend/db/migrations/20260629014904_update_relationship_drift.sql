-- +goose Up
-- =========================================================
-- RELATIONSHIPS DRIFT SYSTEM UPGRADE (CLEAN + CONSISTENT)
-- =========================================================

-- =========================================================
-- 1. FIX + NORMALIZE DRIFT FIELD
-- =========================================================

-- . Convert type first

ALTER TABLE relationships
ALTER COLUMN drift_interval_days DROP DEFAULT;

ALTER TABLE relationships
ALTER COLUMN drift_interval_days TYPE INTEGER
USING NULLIF(drift_interval_days, '')::INTEGER;

-- 2. Rename column
ALTER TABLE relationships
RENAME COLUMN drift_interval_days TO drift_threshold_days;

-- 3. Set NOT NULL + DEFAULT safely
ALTER TABLE relationships
ALTER COLUMN drift_threshold_days SET DEFAULT 30;

-- Optional but recommended (if you want strict consistency)
ALTER TABLE relationships
ALTER COLUMN drift_threshold_days SET NOT NULL;

-- =========================================================
-- 2. REMINDER SYSTEM UPGRADE
-- =========================================================

ALTER TABLE relationships
ADD COLUMN last_reminder_sent_at TIMESTAMPTZ;

ALTER TABLE relationships
DROP COLUMN IF EXISTS reminder_sent;

-- =========================================================
-- 3. DRIFT + INTELLIGENCE FIELDS
-- =========================================================

ALTER TABLE relationships
ADD COLUMN drift_status VARCHAR(50) NOT NULL DEFAULT 'healthy';

ALTER TABLE relationships
ADD COLUMN warmth_score INTEGER DEFAULT 100;

-- =========================================================
-- 4. INDEXES (PERFORMANCE LAYER)
-- =========================================================

CREATE INDEX IF NOT EXISTS idx_relationships_last_contact_at
ON relationships(last_contact_at);

CREATE INDEX IF NOT EXISTS idx_relationships_drift_status
ON relationships(drift_status);

CREATE INDEX IF NOT EXISTS idx_relationships_warmth_score
ON relationships(warmth_score);

-- =========================================================
-- 5. TRIGGER: SYNC LAST CONTACT ON INSERT CONTACT
-- =========================================================

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION sync_relationship_last_contact()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE relationships
    SET
        last_contact_at = NEW.contacted_at,
        updated_at = NOW()
    WHERE id = NEW.relationship_id
      AND (
          last_contact_at IS NULL
          OR NEW.contacted_at > last_contact_at
      );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trigger_sync_relationship_last_contact ON contacts;

CREATE TRIGGER trigger_sync_relationship_last_contact
AFTER INSERT ON contacts
FOR EACH ROW
EXECUTE FUNCTION sync_relationship_last_contact();

-- =========================================================
-- 6. TRIGGER: RECALCULATE ON DELETE CONTACT
-- =========================================================
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION recalculate_relationship_last_contact()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE relationships
    SET
        last_contact_at = (
            SELECT MAX(contacted_at)
            FROM contacts
            WHERE relationship_id = OLD.relationship_id
        ),
        updated_at = NOW()
    WHERE id = OLD.relationship_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trigger_recalculate_relationship_last_contact ON contacts;

CREATE TRIGGER trigger_recalculate_relationship_last_contact
AFTER DELETE ON contacts
FOR EACH ROW
EXECUTE FUNCTION recalculate_relationship_last_contact();


-- +goose Down
-- =========================================================
-- FULL ROLLBACK (SAFE ORDER)
-- =========================================================

-- =========================================================
-- 1. DROP TRIGGERS
-- =========================================================

DROP TRIGGER IF EXISTS trigger_sync_relationship_last_contact ON contacts;
DROP TRIGGER IF EXISTS trigger_recalculate_relationship_last_contact ON contacts;

-- =========================================================
-- 2. DROP FUNCTIONS
-- =========================================================

DROP FUNCTION IF EXISTS sync_relationship_last_contact();
DROP FUNCTION IF EXISTS recalculate_relationship_last_contact();

-- =========================================================
-- 3. DROP INDEXES
-- =========================================================

DROP INDEX IF EXISTS idx_relationships_last_contact_at;
DROP INDEX IF EXISTS idx_relationships_drift_status;
DROP INDEX IF EXISTS idx_relationships_warmth_score;

-- =========================================================
-- 4. DROP NEW COLUMNS
-- =========================================================

ALTER TABLE relationships DROP COLUMN IF EXISTS drift_status;
ALTER TABLE relationships DROP COLUMN IF EXISTS warmth_score;
ALTER TABLE relationships DROP COLUMN IF EXISTS last_reminder_sent_at;

-- =========================================================
-- 5. RESTORE OLD SCHEMA
-- =========================================================

ALTER TABLE relationships
ADD COLUMN reminder_sent BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE relationships
RENAME COLUMN drift_threshold_days TO drift_interval_days;

ALTER TABLE relationships
ALTER COLUMN drift_interval_days TYPE VARCHAR(50)
USING drift_interval_days::TEXT;