-- name: GetCompanyProfile :one
SELECT * FROM company_profile WHERE id = 'default';

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
    updated_at = datetime('now')
WHERE id = 'default'
RETURNING *;
