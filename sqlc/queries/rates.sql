-- name: CreateRate :one
INSERT INTO rates (id, name, category_id, supplier, rate, unit, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetRate :one
SELECT * FROM rates WHERE id = $1;

-- name: ListRates :many
SELECT * FROM rates ORDER BY name ASC;

-- name: ListRatesByCategory :many
SELECT * FROM rates WHERE category_id = $1 ORDER BY name ASC;

-- name: SearchRates :many
SELECT * FROM rates
WHERE name ILIKE '%' || @search_term::text || '%'
   OR supplier ILIKE '%' || @search_term::text || '%'
   OR notes ILIKE '%' || @search_term::text || '%'
ORDER BY name ASC;

-- name: SearchRatesByCategory :many
SELECT * FROM rates
WHERE category_id = @category_id
  AND (name ILIKE '%' || @search_term::text || '%'
       OR supplier ILIKE '%' || @search_term::text || '%'
       OR notes ILIKE '%' || @search_term::text || '%')
ORDER BY name ASC;

-- name: UpdateRate :one
UPDATE rates
SET name = $2, category_id = $3, supplier = $4, rate = $5, unit = $6, notes = $7, last_updated = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRate :exec
DELETE FROM rates WHERE id = $1;
