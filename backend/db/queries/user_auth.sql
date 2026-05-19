-- name: CreateUserAuth :one
INSERT INTO user_auth (
    user_id,
    password_hash,
    mfa_enabled,
    failed_attempts
)
VALUES (
    $1, $2, FALSE, 0
)
RETURNING *;

-- name: GetUserAuthByUserID :one
SELECT *
FROM user_auth
WHERE user_id = $1
LIMIT 1;

