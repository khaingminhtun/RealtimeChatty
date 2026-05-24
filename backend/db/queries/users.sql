-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    display_name,
    is_verified,
    is_active
)
VALUES (
    $1, $2, $3, FALSE, TRUE
)
RETURNING *;


-- name: GetUserByID :one
SELECT id, username, email, display_name, avatar_url, bio, timezone, is_verified, is_active, privacy_setting, notifications_enabled, push_token, last_seen_at, created_at, updated_at
FROM users
WHERE id = $1 AND is_active = TRUE LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: MarkUserAsVerified :exec
UPDATE users
SET 
    is_verified = TRUE, -- Or whatever your column name is (e.g., status = 'verified')
    updated_at = NOW()
WHERE email = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET 
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    timezone = COALESCE(sqlc.narg('timezone'), timezone),
    push_token = COALESCE(sqlc.narg('push_token'), push_token),
    updated_at = NOW()
WHERE id = $1 AND is_active = TRUE
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET 
    is_active = FALSE,
    updated_at = NOW()
WHERE id = $1;