-- +goose Up
-- Add org_id to all tenant-scoped tables for multi-tenancy

-- Jobs table
ALTER TABLE jobs ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_jobs_org_id ON jobs(org_id);

-- Categories table (inherits from jobs, but needs direct org_id for queries)
ALTER TABLE categories ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_categories_org_id ON categories(org_id);

-- Line items table (inherits from categories, but needs direct org_id for queries)
ALTER TABLE line_items ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_line_items_org_id ON line_items(org_id);

-- Clients table
ALTER TABLE clients ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_clients_org_id ON clients(org_id);

-- Estimates table (inherits from jobs, but needs direct org_id)
ALTER TABLE estimates ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_estimates_org_id ON estimates(org_id);

-- Estimate categories table
ALTER TABLE estimate_categories ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_estimate_categories_org_id ON estimate_categories(org_id);

-- Signature requests table (inherits from estimates, but needs direct org_id)
ALTER TABLE signature_requests ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_signature_requests_org_id ON signature_requests(org_id);

-- Signatures table (inherits from signature_requests, but needs direct org_id)
ALTER TABLE signatures ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_signatures_org_id ON signatures(org_id);

-- Item templates table
ALTER TABLE item_templates ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_item_templates_org_id ON item_templates(org_id);

-- Price imports table
ALTER TABLE price_imports ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_price_imports_org_id ON price_imports(org_id);

-- Price import matches table (inherits from price_imports, but needs direct org_id)
ALTER TABLE price_import_matches ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_price_import_matches_org_id ON price_import_matches(org_id);

-- Job item types table (inherits from jobs, but needs direct org_id)
ALTER TABLE job_item_types ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_job_item_types_org_id ON job_item_types(org_id);

-- Settings table: one settings record per organization
-- Delete existing default row (will be recreated per-organization in application code)
DELETE FROM settings WHERE id = 'default';
-- Drop the old default value constraint and recreate with org_id as UUID foreign key
ALTER TABLE settings DROP CONSTRAINT IF EXISTS settings_pkey;
ALTER TABLE settings DROP COLUMN id;
ALTER TABLE settings ADD COLUMN org_id UUID NOT NULL;
ALTER TABLE settings ADD PRIMARY KEY (org_id);
ALTER TABLE settings ADD CONSTRAINT settings_org_id_fkey FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- Company profile table: one company profile per organization
-- Delete existing default row (will be recreated per-organization in application code)
DELETE FROM company_profile WHERE id = 'default';
-- Drop the old default value constraint and recreate with org_id as UUID foreign key
ALTER TABLE company_profile DROP CONSTRAINT IF EXISTS company_profile_pkey;
ALTER TABLE company_profile DROP COLUMN id;
ALTER TABLE company_profile ADD COLUMN org_id UUID NOT NULL;
ALTER TABLE company_profile ADD PRIMARY KEY (org_id);
ALTER TABLE company_profile ADD CONSTRAINT company_profile_org_id_fkey FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- NOTE: Migration requires database to be empty or org_id columns to be populated manually.
-- For fresh installations, org_id will be set during data creation.
-- For existing data, populate org_id columns before applying NOT NULL constraints in a future migration.

-- +goose Down
-- Remove org_id columns and indexes

-- Jobs
DROP INDEX IF EXISTS idx_jobs_org_id;
ALTER TABLE jobs DROP COLUMN IF EXISTS org_id;

-- Categories
DROP INDEX IF EXISTS idx_categories_org_id;
ALTER TABLE categories DROP COLUMN IF EXISTS org_id;

-- Line items
DROP INDEX IF EXISTS idx_line_items_org_id;
ALTER TABLE line_items DROP COLUMN IF EXISTS org_id;

-- Clients
DROP INDEX IF EXISTS idx_clients_org_id;
ALTER TABLE clients DROP COLUMN IF EXISTS org_id;

-- Estimates
DROP INDEX IF EXISTS idx_estimates_org_id;
ALTER TABLE estimates DROP COLUMN IF EXISTS org_id;

-- Estimate categories
DROP INDEX IF EXISTS idx_estimate_categories_org_id;
ALTER TABLE estimate_categories DROP COLUMN IF EXISTS org_id;

-- Signature requests
DROP INDEX IF EXISTS idx_signature_requests_org_id;
ALTER TABLE signature_requests DROP COLUMN IF EXISTS org_id;

-- Signatures
DROP INDEX IF EXISTS idx_signatures_org_id;
ALTER TABLE signatures DROP COLUMN IF EXISTS org_id;

-- Item templates
DROP INDEX IF EXISTS idx_item_templates_org_id;
ALTER TABLE item_templates DROP COLUMN IF EXISTS org_id;

-- Price imports
DROP INDEX IF EXISTS idx_price_imports_org_id;
ALTER TABLE price_imports DROP COLUMN IF EXISTS org_id;

-- Price import matches
DROP INDEX IF EXISTS idx_price_import_matches_org_id;
ALTER TABLE price_import_matches DROP COLUMN IF EXISTS org_id;

-- Job item types
DROP INDEX IF EXISTS idx_job_item_types_org_id;
ALTER TABLE job_item_types DROP COLUMN IF EXISTS org_id;

-- Settings: revert to id column
ALTER TABLE settings DROP CONSTRAINT IF EXISTS settings_org_id_fkey;
ALTER TABLE settings DROP CONSTRAINT IF EXISTS settings_pkey;
ALTER TABLE settings DROP COLUMN IF EXISTS org_id;
ALTER TABLE settings ADD COLUMN id TEXT PRIMARY KEY DEFAULT 'default';

-- Company profile: revert to id column
ALTER TABLE company_profile DROP CONSTRAINT IF EXISTS company_profile_org_id_fkey;
ALTER TABLE company_profile DROP CONSTRAINT IF EXISTS company_profile_pkey;
ALTER TABLE company_profile DROP COLUMN IF EXISTS org_id;
ALTER TABLE company_profile ADD COLUMN id TEXT PRIMARY KEY DEFAULT 'default';
