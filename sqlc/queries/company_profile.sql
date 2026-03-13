-- name: GetCompanyProfile :one
SELECT * FROM company_profile WHERE id = 'default';

-- name: UpsertCompanyProfile :one
INSERT INTO company_profile (id, name, email, phone, address, city, state, zip)
VALUES ('default', $1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    address = EXCLUDED.address,
    city = EXCLUDED.city,
    state = EXCLUDED.state,
    zip = EXCLUDED.zip,
    updated_at = now()
RETURNING *;
