-- name: CreateSession :one
INSERT INTO sessions (user_id, token_hash, user_agent, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT
    s.id, s.user_id, s.token_hash, s.user_agent, s.ip_address,
    s.expires_at, s.last_activity_at, s.created_at,
    u.email, u.name, u.status AS user_status
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = $1 AND s.expires_at > now();

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions WHERE token_hash = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at <= now();

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: UpdateSessionActivity :exec
UPDATE sessions SET last_activity_at = now() WHERE id = $1;
