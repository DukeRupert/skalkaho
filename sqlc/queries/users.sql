-- name: CreateUser :one
INSERT INTO users (
    org_id,
    email,
    password_hash,
    name,
    role,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE org_id = $1 AND email = $2 LIMIT 1;

-- name: GetUserByResetToken :one
SELECT * FROM users
WHERE reset_token = $1
  AND reset_token_expires_at > NOW()
LIMIT 1;

-- name: GetUserByVerificationToken :one
SELECT * FROM users
WHERE verification_token = $1 LIMIT 1;

-- name: ListUsersByOrg :many
SELECT * FROM users
WHERE org_id = $1
ORDER BY created_at ASC;

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE(sqlc.narg('email'), email),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    name = COALESCE(sqlc.narg('name'), name),
    role = COALESCE(sqlc.narg('role'), role),
    status = COALESCE(sqlc.narg('status'), status),
    email_verified = COALESCE(sqlc.narg('email_verified'), email_verified),
    reset_token = sqlc.narg('reset_token'),
    reset_token_expires_at = sqlc.narg('reset_token_expires_at'),
    verification_token = sqlc.narg('verification_token'),
    last_login_at = COALESCE(sqlc.narg('last_login_at'), last_login_at),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = NOW()
WHERE id = $1;

-- name: SetUserResetToken :exec
UPDATE users
SET
    reset_token = $1,
    reset_token_expires_at = $2,
    updated_at = NOW()
WHERE id = $3;

-- name: ClearUserResetToken :exec
UPDATE users
SET
    reset_token = NULL,
    reset_token_expires_at = NULL,
    updated_at = NOW()
WHERE id = $1;

-- name: SetUserVerificationToken :exec
UPDATE users
SET
    verification_token = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: VerifyUserEmail :exec
UPDATE users
SET
    email_verified = TRUE,
    verification_token = NULL,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsersByOrg :one
SELECT COUNT(*) FROM users
WHERE org_id = $1;
