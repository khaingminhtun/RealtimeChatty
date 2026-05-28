-- name: CreateRelationship :one
INSERT INTO relationships (
    owner_id, 
    name, 
    type, 
    how_we_met, 
    birthday, 
    location, 
    tags
) VALUES (
    $1, $2, $3, $4, $5, $6,$7 
)
RETURNING id, owner_id, name, type, how_we_met, birthday, location, avatar_url, tags, last_contact_at, created_at, updated_at;


-- name: GetRelationshipByID :one
SELECT id, owner_id, name, type, how_we_met, birthday, location, avatar_url, tags, last_contact_at, created_at, updated_at
FROM relationships
WHERE id = $1 AND owner_id = $2;

-- name: ListRelationships :many
SELECT id, owner_id, name, type, how_we_met, location, birthday, tags, last_contact_at, created_at
FROM relationships
WHERE owner_id = $1 
  AND ($2::text = '' OR type = $2)
ORDER BY last_contact_at DESC NULLS LAST, created_at DESC;

-- name: UpdateRelationship :one
UPDATE relationships
SET 
    name = COALESCE(sqlc.narg('name'), name),
    type = COALESCE(sqlc.narg('type'), type),
    how_we_met = COALESCE(sqlc.narg('how_we_met'), how_we_met),
    birthday = COALESCE(sqlc.narg('birthday'), birthday),
    location = COALESCE(sqlc.narg('location'), location),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = NOW()
WHERE id = $1 AND owner_id = $2
RETURNING id, owner_id, name, type, how_we_met, birthday, location, avatar_url, tags, last_contact_at, created_at, updated_at;

-- name: DeleteRelationship :exec
DELETE FROM relationships
WHERE id = $1 AND owner_id = $2;