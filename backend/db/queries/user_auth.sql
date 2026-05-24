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

-- name: UpdateUserPassword :exec
UPDATE user_auth
SET password_hash = $2,
    updated_at = NOW()
WHERE user_id = $1;
