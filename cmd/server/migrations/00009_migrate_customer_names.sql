-- +goose Up
-- Create clients from unique customer_names
INSERT INTO clients (id, name, created_at)
SELECT
    gen_random_uuid()::text,
    customer_name,
    NOW()
FROM jobs
WHERE customer_name IS NOT NULL AND customer_name != ''
GROUP BY customer_name;

-- Link jobs to their new client records
UPDATE jobs
SET client_id = (SELECT id FROM clients WHERE clients.name = jobs.customer_name)
WHERE customer_name IS NOT NULL AND customer_name != '';

-- +goose Down
UPDATE jobs SET client_id = NULL;
DELETE FROM clients;
