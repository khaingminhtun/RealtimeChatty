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
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, owner_id, name, type, how_we_met, birthday, location, tags, last_contact_at, created_at, updated_at, deleted_at;