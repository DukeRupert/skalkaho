-- name: CreateClient :one
INSERT INTO clients (id, name, company, email, phone, address, city, state, zip, tax_id, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetClient :one
SELECT * FROM clients WHERE id = $1;

-- name: GetClientByName :one
SELECT * FROM clients WHERE name = $1;

-- name: ListClients :many
SELECT * FROM clients ORDER BY name ASC;

-- name: ListClientsPaginated :many
SELECT * FROM clients
WHERE (sqlc.arg('search') = '' OR name LIKE '%' || sqlc.arg('search') || '%' OR company LIKE '%' || sqlc.arg('search') || '%')
ORDER BY name ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountClients :one
SELECT COUNT(*) FROM clients
WHERE (sqlc.arg('search') = '' OR name LIKE '%' || sqlc.arg('search') || '%' OR company LIKE '%' || sqlc.arg('search') || '%');

-- name: UpdateClient :one
UPDATE clients SET
    name = $1,
    company = $2,
    email = $3,
    phone = $4,
    address = $5,
    city = $6,
    state = $7,
    zip = $8,
    tax_id = $9,
    notes = $10
WHERE id = $11
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM clients WHERE id = $1;

-- name: ClientHasJobs :one
SELECT COUNT(*) > 0 FROM jobs WHERE client_id = $1;
