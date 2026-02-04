-- name: CreateOrganization :one
INSERT INTO organizations (
    name,
    subdomain,
    stripe_customer_id,
    plan,
    status,
    trial_ends_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1 LIMIT 1;

-- name: GetOrganizationBySubdomain :one
SELECT * FROM organizations
WHERE subdomain = $1 LIMIT 1;

-- name: GetOrganizationByStripeCustomerID :one
SELECT * FROM organizations
WHERE stripe_customer_id = $1 LIMIT 1;

-- name: UpdateOrganization :one
UPDATE organizations
SET
    name = COALESCE(sqlc.narg('name'), name),
    subdomain = COALESCE(sqlc.narg('subdomain'), subdomain),
    stripe_customer_id = COALESCE(sqlc.narg('stripe_customer_id'), stripe_customer_id),
    plan = COALESCE(sqlc.narg('plan'), plan),
    status = COALESCE(sqlc.narg('status'), status),
    trial_ends_at = COALESCE(sqlc.narg('trial_ends_at'), trial_ends_at),
    subscription_ends_at = COALESCE(sqlc.narg('subscription_ends_at'), subscription_ends_at),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: ListOrganizations :many
SELECT * FROM organizations
ORDER BY created_at DESC;

-- name: ListActiveOrganizations :many
SELECT * FROM organizations
WHERE status = 'active'
ORDER BY created_at DESC;
