-- name: CreateMaterial :one
INSERT INTO materials (id, name, supplier_id, unit_price, unit, supplier_code, price_source)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMaterial :one
SELECT * FROM materials WHERE id = $1;

-- name: ListMaterials :many
SELECT * FROM materials ORDER BY name ASC;

-- name: ListMaterialsBySupplier :many
SELECT * FROM materials WHERE supplier_id = $1 ORDER BY name ASC;

-- name: SearchMaterials :many
SELECT * FROM materials
WHERE name ILIKE '%' || @search_term::text || '%'
   OR supplier_code ILIKE '%' || @search_term::text || '%'
ORDER BY name ASC;

-- name: SearchMaterialsBySupplier :many
SELECT * FROM materials
WHERE supplier_id = @supplier_id
  AND (name ILIKE '%' || @search_term::text || '%'
       OR supplier_code ILIKE '%' || @search_term::text || '%')
ORDER BY name ASC;

-- name: ListMaterialsByPriceSource :many
SELECT * FROM materials WHERE price_source = $1 ORDER BY name ASC;

-- name: ListMaterialsBySupplierAndSource :many
SELECT * FROM materials
WHERE supplier_id = $1 AND price_source = $2
ORDER BY name ASC;

-- name: UpdateMaterial :one
UPDATE materials
SET name = $2, supplier_id = $3, unit_price = $4, unit = $5,
    supplier_code = $6, price_source = $7, last_updated = now()
WHERE id = $1
RETURNING *;

-- name: DeleteMaterial :exec
DELETE FROM materials WHERE id = $1;
