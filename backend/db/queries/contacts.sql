-- name: CreateContactLog :one
INSERT INTO contacts (relationship_id, user_id, channel, note, contacted_at)
VALUES ($1, $2, $3, $4, COALESCE($5, NOW()))
RETURNING *;

-- name: GetContactsByRelationship :many
SELECT *
FROM contacts
WHERE relationship_id = $1
ORDER BY contacted_at DESC;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = $1 AND user_id = $2;

-- name: SearchRelationships :many
SELECT *
FROM relationships
WHERE owner_id = $1
  AND (
    to_tsvector(
      'english',
      coalesce(name, '') || ' ' ||
      coalesce(type, '') || ' ' ||
      coalesce(how_we_met, '')
    )
    @@ plainto_tsquery('english', $2)
  )
ORDER BY last_contact_at DESC;

-- name: ListRelationshipsForDrift :many
SELECT 
    id, owner_id, name, type,
    drift_threshold_days,
    last_contact_at,
    next_contact_at
FROM relationships
WHERE owner_id = $1;

-- name: GetPendingDriftReminders :many
SELECT 
    id, owner_id, name, type,
    drift_threshold_days,
    last_contact_at,
    next_contact_at
FROM relationships
WHERE next_contact_at IS NOT NULL
  AND next_contact_at <= NOW();

-- name: MarkReminderAsSent :exec
UPDATE relationships
SET 
    last_reminder_sent_at = NOW()
WHERE id = $1;


-- name: GetContactByID :one
SELECT *
FROM contacts
WHERE id = $1
  AND user_id = $2;

-- name: UpdateContact :one
UPDATE contacts
SET
    channel = COALESCE($3, channel),
    note = COALESCE($4, note),
    contacted_at = COALESCE($5, contacted_at),
    updated_at = NOW()
WHERE id = $1
  AND user_id = $2
RETURNING *;