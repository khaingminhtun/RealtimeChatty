-- name: GetPrivateNote :one
SELECT id, user_id, relationship_id, content, created_at, updated_at
FROM private_notes
WHERE relationship_id = $1 AND user_id = $2;

-- name: UpsertPrivateNote :one
INSERT INTO private_notes (user_id, relationship_id, content, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (user_id, relationship_id) 
DO UPDATE SET 
    content = EXCLUDED.content,
    updated_at = NOW()
RETURNING id, user_id, relationship_id, content, created_at, updated_at;

-- name: DeletePrivateNote :exec
DELETE FROM private_notes
WHERE relationship_id = $1 AND user_id = $2;