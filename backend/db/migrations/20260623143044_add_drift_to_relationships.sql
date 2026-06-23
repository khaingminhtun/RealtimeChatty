-- +goose Up
-- Add drift_interval_days as a NULLable VARCHAR string with NO default value
ALTER TABLE relationships 
ADD COLUMN drift_interval_days VARCHAR(50) DEFAULT NULL;

-- Add reminder_sent tracking flag for the background worker
ALTER TABLE relationships 
ADD COLUMN reminder_sent BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE relationships DROP COLUMN IF EXISTS drift_interval_days;
ALTER TABLE relationships DROP COLUMN IF EXISTS reminder_sent;