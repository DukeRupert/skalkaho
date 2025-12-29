-- name: CreateClient :one
INSERT INTO clients (id, name, company, email, phone, address, city, state, zip, tax_id, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetClient :one
SELECT * FROM clients WHERE id = ?;

-- name: GetClientByName :one
SELECT * FROM clients WHERE name = ?;

-- name: ListClients :many
SELECT * FROM clients ORDER BY name ASC;

-- name: ListClientsPaginated :many
SELECT * FROM clients
WHERE (@search = '' OR name LIKE '%' || @search || '%' OR company LIKE '%' || @search || '%')
ORDER BY name ASC
LIMIT @limit OFFSET @offset;

-- name: CountClients :one
SELECT COUNT(*) FROM clients
WHERE (@search = '' OR name LIKE '%' || @search || '%' OR company LIKE '%' || @search || '%');

-- name: UpdateClient :one
UPDATE clients SET
    name = ?,
    company = ?,
    email = ?,
    phone = ?,
    address = ?,
    city = ?,
    state = ?,
    zip = ?,
    tax_id = ?,
    notes = ?
WHERE id = ?
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM clients WHERE id = ?;

-- name: ClientHasJobs :one
SELECT COUNT(*) > 0 FROM jobs WHERE client_id = ?;
