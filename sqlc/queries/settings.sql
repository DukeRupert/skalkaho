-- name: GetSettings :one
SELECT * FROM settings
WHERE id = 'default';

-- name: UpdateSettings :one
UPDATE settings SET
    default_surcharge_mode = ?,
    default_surcharge_percent = ?
WHERE id = 'default'
RETURNING *;
