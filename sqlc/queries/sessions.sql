-- name: CreateSession :one
INSERT INTO sessions (user_id, org_id, token_hash, user_agent, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT s.*, u.email, u.name, u.role, u.status as user_status
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token_hash = $1 AND s.expires_at > NOW() LIMIT 1;

-- name: UpdateSessionActivity :exec
UPDATE sessions SET last_activity_at = NOW() WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions WHERE token_hash = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < NOW();
