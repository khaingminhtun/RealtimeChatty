-- name: CreateContactLog :one
INSERT INTO contacts (relationship_id, user_id, channel, note, contacted_at)
VALUES ($1, $2, $3, $4, COALESCE($5, NOW()))
RETURNING *;

-- name: UpdateRelationshipLastContact :exec
UPDATE relationships
SET 
    last_contact_at = $1,
    reminder_sent = FALSE, -- Reset the worker alert state automatically on new contact
    updated_at = NOW()
WHERE id = $2 AND owner_id = $3;

-- name: ListRelationshipsForDrift :many
SELECT id, owner_id, name, type, drift_interval_days, last_contact_at 
FROM relationships
WHERE owner_id = $1;

-- name: GetPendingDriftReminders :many
SELECT id, owner_id, name, type, drift_interval_days, last_contact_at
FROM relationships
WHERE last_contact_at IS NOT NULL 
  AND reminder_sent = FALSE;

-- name: MarkReminderAsSent :exec
UPDATE relationships
SET reminder_sent = TRUE
WHERE id = $1;

-- name: SearchRelationships :many
SELECT * FROM relationships
WHERE owner_id = $1
  AND (
    to_tsvector('english', coalesce(name, '') || ' ' || coalesce(type, '') || ' ' || coalesce(how_we_met, '')) 
    @@ plainto_tsquery('english', $2)
  )
ORDER BY last_contact_at DESC;