-- name: GetCompanyProfile :one
SELECT * FROM company_profile WHERE id = 'default';

-- name: UpdateCompanyProfile :one
UPDATE company_profile SET
    name = ?,
    email = ?,
    phone = ?,
    address = ?,
    city = ?,
    state = ?,
    zip = ?,
    logo_path = ?,
    updated_at = datetime('now')
WHERE id = 'default'
RETURNING *;
