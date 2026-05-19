-- name: CreateEmailVerification :one
INSERT INTO email_verifications (
    user_id,
    email,
    otp_hash,
    expires_at
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetLatestEmailVerification :one
SELECT *
FROM email_verifications
WHERE email = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkEmailVerificationUsed :exec
UPDATE email_verifications
SET is_used = TRUE
WHERE id = $1;