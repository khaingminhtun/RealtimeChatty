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