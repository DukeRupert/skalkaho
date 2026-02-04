-- name: GetSettings :one
SELECT * FROM settings
WHERE id = 'default';

-- name: UpdateSettings :one
UPDATE settings SET
    default_surcharge_mode = $1,
    default_surcharge_percent = $2
WHERE id = 'default'
RETURNING *;
