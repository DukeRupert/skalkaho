-- name: GetSettings :one
SELECT * FROM settings
WHERE org_id = $1;

-- name: CreateSettings :one
INSERT INTO settings (org_id, default_surcharge_mode, default_surcharge_percent)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSettings :one
UPDATE settings SET
    default_surcharge_mode = $1,
    default_surcharge_percent = $2
WHERE org_id = $3
RETURNING *;
