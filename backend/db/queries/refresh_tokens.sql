-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, 
    session_id, 
    token_hash, 
    expires_at
) VALUES (
    $1, $2, $3, $4
) 
RETURNING id, user_id, session_id, token_hash, is_revoked, expires_at, created_at;

-- name: GetActiveTokenByHash :one
-- Retrieves a token only if it hasn't been explicitly revoked and hasn't expired yet
SELECT id, user_id, session_id, token_hash, is_revoked, expires_at, created_at
FROM refresh_tokens
WHERE token_hash = $1 
  AND is_revoked = FALSE 
  AND expires_at > NOW()
LIMIT 1;

-- name: RevokeTokenByHash :exec
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE token_hash = $1;

-- name: RevokeAllSessionsForUser :exec
-- Useful for "log out of all devices" or global password resets
UPDATE refresh_tokens    
SET is_revoked = TRUE
WHERE user_id = $1 
  AND is_revoked = FALSE;

-- name: RevokeSpecificSession :exec
-- Useful for logging out of one specific device/browser session
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE session_id = $1 
  AND is_revoked = FALSE;

-- name: DeleteExpiredTokens :exec
-- Run this on a cron job or background worker to keep your DB lean
DELETE FROM refresh_tokens
WHERE expires_at <= NOW();