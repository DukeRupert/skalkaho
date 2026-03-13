-- name: CreateClient :one
INSERT INTO clients (id, company_name, contact_name, email, phone, address, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetClient :one
SELECT * FROM clients WHERE id = $1;

-- name: ListClients :many
SELECT * FROM clients ORDER BY company_name ASC;

-- name: UpdateClient :one
UPDATE clients
SET company_name = $2, contact_name = $3, email = $4, phone = $5, address = $6, notes = $7
WHERE id = $1
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM clients WHERE id = $1;
