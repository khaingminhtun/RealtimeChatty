-- name: CreateSession :one
INSERT INTO sessions (
    user_id, 
    token_hash, 
    device_name, 
    device_type, 
    ip_address, 
    user_agent, 
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) 
RETURNING id, user_id, token_hash, device_name, device_type, ip_address, user_agent, is_active, expires_at, created_at;

-- name: GetActiveSessionByHash :one
-- Fetches the session details only if it's active and hasn't expired
SELECT id, user_id, token_hash, device_name, device_type, ip_address, user_agent, is_active, expires_at, created_at
FROM sessions
WHERE token_hash = $1 
  AND is_active = TRUE 
  AND expires_at > NOW()
LIMIT 1;

-- name: GetActiveSessionsByUser :many
-- Useful for showing a user a list of their "Logged-in Devices"
SELECT id, device_name, device_type, ip_address, user_agent, expires_at, created_at
FROM sessions
WHERE user_id = $1 
  AND is_active = TRUE 
  AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: UpdateSessionActivity :exec
-- Call this if you use a rolling session mechanism to extend expiration on use
UPDATE sessions
SET expires_at = $2,
    ip_address = $3,
    user_agent = $4
WHERE id = $1 
  AND is_active = TRUE;

-- name: DeactivateSession :exec
-- Standard single device logout
UPDATE sessions
SET is_active = FALSE
WHERE id = $1;

-- name: DeactivateAllUserSessions :exec
-- Great for a "Log out of all other devices" feature
UPDATE sessions
SET is_active = FALSE
WHERE user_id = $1 
  AND id <> $2
  AND is_active = TRUE;

-- name: DeleteExpiredSessions :exec
-- Housekeeping query to clean up dead rows over time
DELETE FROM sessions
WHERE expires_at <= NOW() 
   OR is_active = FALSE;