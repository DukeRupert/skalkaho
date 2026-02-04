-- name: GetCompanyProfile :one
SELECT * FROM company_profile WHERE org_id = $1;

-- name: CreateCompanyProfile :one
INSERT INTO company_profile (org_id, name, email, phone, address, city, state, zip, logo_path)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateCompanyProfile :one
UPDATE company_profile SET
    name = $1,
    email = $2,
    phone = $3,
    address = $4,
    city = $5,
    state = $6,
    zip = $7,
    logo_path = $8,
    updated_at = NOW()
WHERE org_id = $9
RETURNING *;
