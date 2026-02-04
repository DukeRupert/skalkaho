//go:build integration

package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTenantIsolation_OrgIDColumnsExist verifies that org_id columns exist on all tenant-scoped tables.
// This is the foundational test - if this fails, multi-tenancy isn't implemented.
func TestTenantIsolation_OrgIDColumnsExist(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	tables := []string{"jobs", "categories", "line_items", "clients", "estimates"}

	for _, table := range tables {
		t.Run(table+"_has_org_id", func(t *testing.T) {
			// Query table info to check for org_id column (PostgreSQL)
			var count int
			err := db.QueryRow(`
				SELECT COUNT(*)
				FROM information_schema.columns
				WHERE table_name = $1 AND column_name = 'org_id'
			`, table).Scan(&count)
			require.NoError(t, err)
			assert.Equal(t, 1, count, "table %s should have org_id column", table)
		})
	}

	// Settings table uses org_id as primary key
	t.Run("settings_has_org_id", func(t *testing.T) {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*)
			FROM information_schema.columns
			WHERE table_name = 'settings' AND column_name = 'org_id'
		`).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "settings table should have org_id column")
	})
}

// TestTenantIsolation_Jobs verifies that jobs are completely isolated by organization.
func TestTenantIsolation_Jobs(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-jobs', NOW()),
			($2, 'Organization B', 'org-b-jobs', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Create jobs for two different organizations using direct SQL
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('job-a1', $1, 'Org A - Kitchen Remodel', 15.0, 'stacking', 'draft', NOW()),
			('job-a2', $1, 'Org A - Bathroom Addition', 20.0, 'stacking', 'draft', NOW()),
			('job-b1', $2, 'Org B - Deck Construction', 10.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to insert jobs with org_id")

	t.Run("ListJobs_filtered_by_org_id", func(t *testing.T) {
		// Query jobs for org-a only
		rows, err := db.QueryContext(ctx, `SELECT id, name FROM jobs WHERE org_id = $1`, orgA)
		require.NoError(t, err)
		defer rows.Close()

		var jobs []string
		for rows.Next() {
			var id, name string
			require.NoError(t, rows.Scan(&id, &name))
			jobs = append(jobs, id)
		}
		require.NoError(t, rows.Err())

		assert.Len(t, jobs, 2, "org-a should have exactly 2 jobs")
		assert.Contains(t, jobs, "job-a1")
		assert.Contains(t, jobs, "job-a2")
		assert.NotContains(t, jobs, "job-b1", "org-b's job should not appear in org-a query")
	})

	t.Run("GetJob_respects_org_boundary", func(t *testing.T) {
		// Org A trying to get Org B's job should return no rows
		var id string
		err := db.QueryRowContext(ctx, `
			SELECT id FROM jobs WHERE id = $1 AND org_id = $2
		`, "job-b1", orgA).Scan(&id)

		assert.ErrorIs(t, err, sql.ErrNoRows, "org-a should not be able to access org-b's job")
	})

	t.Run("CountJobs_scoped_by_org", func(t *testing.T) {
		var countA, countB int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE org_id = $1`, orgA).Scan(&countA)
		require.NoError(t, err)
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE org_id = $1`, orgB).Scan(&countB)
		require.NoError(t, err)

		assert.Equal(t, int64(2), countA, "org-a should have 2 jobs")
		assert.Equal(t, int64(1), countB, "org-b should have 1 job")
	})

	t.Run("UpdateJob_respects_org_boundary", func(t *testing.T) {
		// Org A trying to update Org B's job should affect 0 rows
		result, err := db.ExecContext(ctx, `
			UPDATE jobs SET name = 'HACKED' WHERE id = $1 AND org_id = $2
		`, "job-b1", orgA)
		require.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		require.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected, "update should affect 0 rows for other org's job")

		// Verify job wasn't modified
		var name string
		err = db.QueryRowContext(ctx, `SELECT name FROM jobs WHERE id = $1`, "job-b1").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Org B - Deck Construction", name, "org-b's job should not be modified")
	})

	t.Run("DeleteJob_respects_org_boundary", func(t *testing.T) {
		// Org A trying to delete Org B's job should affect 0 rows
		result, err := db.ExecContext(ctx, `
			DELETE FROM jobs WHERE id = $1 AND org_id = $2
		`, "job-b1", orgA)
		require.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		require.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected, "delete should affect 0 rows for other org's job")

		// Verify job still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE id = $1`, "job-b1").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-b's job should still exist")
	})
}

// TestTenantIsolation_Categories verifies categories are isolated by organization.
func TestTenantIsolation_Categories(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-cat', NOW()),
			($2, 'Organization B', 'org-b-cat', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Create jobs and categories for two orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('job-a', $1, 'Org A Job', 15.0, 'stacking', 'draft', NOW()),
			('job-b', $2, 'Org B Job', 10.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO categories (id, org_id, job_id, name, sort_order)
		VALUES
			('cat-a1', $1, 'job-a', 'Org A - Foundation', 0),
			('cat-a2', $1, 'job-a', 'Org A - Framing', 1),
			('cat-b1', $2, 'job-b', 'Org B - Plumbing', 0)
	`, orgA, orgB)
	require.NoError(t, err, "should be able to insert categories with org_id")

	t.Run("ListCategories_filtered_by_org", func(t *testing.T) {
		var count int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(2), count, "org-a should have 2 categories")
	})

	t.Run("GetCategory_respects_org_boundary", func(t *testing.T) {
		var id string
		err := db.QueryRowContext(ctx, `
			SELECT id FROM categories WHERE id = $1 AND org_id = $2
		`, "cat-b1", orgA).Scan(&id)
		assert.ErrorIs(t, err, sql.ErrNoRows, "org-a should not access org-b's category")
	})

	t.Run("Category_inherits_org_from_job", func(t *testing.T) {
		// Verify that categories have the same org_id as their parent job
		var catOrgID, jobOrgID uuid.UUID
		err := db.QueryRowContext(ctx, `
			SELECT c.org_id, j.org_id
			FROM categories c
			JOIN jobs j ON c.job_id = j.id
			WHERE c.id = $1
		`, "cat-a1").Scan(&catOrgID, &jobOrgID)
		require.NoError(t, err)
		assert.Equal(t, jobOrgID, catOrgID, "category org_id should match job org_id")
	})
}

// TestTenantIsolation_LineItems verifies line items are isolated by organization.
func TestTenantIsolation_LineItems(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-li', NOW()),
			($2, 'Organization B', 'org-b-li', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Create complete hierarchy for two orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('job-a', $1, 'Org A Job', 15.0, 'stacking', 'draft', NOW()),
			('job-b', $2, 'Org B Job', 10.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO categories (id, org_id, job_id, name, sort_order)
		VALUES
			('cat-a', $1, 'job-a', 'Org A Category', 0),
			('cat-b', $2, 'job-b', 'Org B Category', 0)
	`, orgA, orgB)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
		VALUES
			('item-a1', $1, 'cat-a', 'material', 'Org A - 2x4 Lumber', 100, 'ea', 5.50, 0),
			('item-b1', $2, 'cat-b', 'material', 'Org B - Copper Pipe', 50, 'ea', 12.00, 0)
	`, orgA, orgB)
	require.NoError(t, err, "should be able to insert line_items with org_id")

	t.Run("ListLineItems_filtered_by_org", func(t *testing.T) {
		var count int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-a should have 1 line item")
	})

	t.Run("GetLineItem_respects_org_boundary", func(t *testing.T) {
		var id string
		err := db.QueryRowContext(ctx, `
			SELECT id FROM line_items WHERE id = $1 AND org_id = $2
		`, "item-b1", orgA).Scan(&id)
		assert.ErrorIs(t, err, sql.ErrNoRows, "org-a should not access org-b's line item")
	})

	t.Run("LineItem_inherits_org_through_hierarchy", func(t *testing.T) {
		// Verify org_id consistency through the hierarchy
		var itemOrg, catOrg, jobOrg uuid.UUID
		err := db.QueryRowContext(ctx, `
			SELECT li.org_id, c.org_id, j.org_id
			FROM line_items li
			JOIN categories c ON li.category_id = c.id
			JOIN jobs j ON c.job_id = j.id
			WHERE li.id = $1
		`, "item-a1").Scan(&itemOrg, &catOrg, &jobOrg)
		require.NoError(t, err)
		assert.Equal(t, jobOrg, catOrg, "category org should match job org")
		assert.Equal(t, catOrg, itemOrg, "line item org should match category org")
	})
}

// TestTenantIsolation_Settings verifies settings are isolated by organization.
func TestTenantIsolation_Settings(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-set', NOW()),
			($2, 'Organization B', 'org-b-set', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Create settings for two orgs (org_id is the primary key)
	_, err = db.ExecContext(ctx, `
		INSERT INTO settings (org_id, default_surcharge_mode, default_surcharge_percent)
		VALUES
			($1, 'stacking', 15.0),
			($2, 'override', 20.0)
	`, orgA, orgB)
	require.NoError(t, err, "should be able to insert settings with org_id")

	t.Run("GetSettings_filtered_by_org", func(t *testing.T) {
		var mode string
		var percent float64
		err := db.QueryRowContext(ctx, `
			SELECT default_surcharge_mode, default_surcharge_percent
			FROM settings WHERE org_id = $1
		`, orgA).Scan(&mode, &percent)
		require.NoError(t, err)
		assert.Equal(t, "stacking", mode)
		assert.Equal(t, 15.0, percent)
	})

	t.Run("Settings_isolated_between_orgs", func(t *testing.T) {
		// Each org should only see their own settings
		var countA, countB int64
		db.QueryRowContext(ctx, `SELECT COUNT(*) FROM settings WHERE org_id = $1`, orgA).Scan(&countA)
		db.QueryRowContext(ctx, `SELECT COUNT(*) FROM settings WHERE org_id = $1`, orgB).Scan(&countB)

		assert.Equal(t, int64(1), countA)
		assert.Equal(t, int64(1), countB)
	})
}

// TestTenantIsolation_Clients verifies clients are isolated by organization.
func TestTenantIsolation_Clients(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-cli', NOW()),
			($2, 'Organization B', 'org-b-cli', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Create clients for two orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO clients (id, org_id, name, email, created_at)
		VALUES
			('client-a1', $1, 'John Doe', 'john@example.com', NOW()),
			('client-b1', $2, 'Jane Smith', 'jane@example.com', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to insert clients with org_id")

	t.Run("ListClients_filtered_by_org", func(t *testing.T) {
		var count int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM clients WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-a should have 1 client")
	})

	t.Run("GetClient_respects_org_boundary", func(t *testing.T) {
		var id string
		err := db.QueryRowContext(ctx, `
			SELECT id FROM clients WHERE id = $1 AND org_id = $2
		`, "client-b1", orgA).Scan(&id)
		assert.ErrorIs(t, err, sql.ErrNoRows, "org-a should not access org-b's client")
	})
}

// TestTenantIsolation_Estimates verifies estimates are isolated by organization.
func TestTenantIsolation_Estimates(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-est', NOW()),
			($2, 'Organization B', 'org-b-est', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Create jobs and estimates for two orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('job-a', $1, 'Org A Job', 15.0, 'stacking', 'draft', NOW()),
			('job-b', $2, 'Org B Job', 10.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
		VALUES
			('est-a1', $1, 'job-a', 1, 'draft', 5000.00, NOW()),
			('est-b1', $2, 'job-b', 1, 'draft', 8000.00, NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to insert estimates with org_id")

	t.Run("ListEstimates_filtered_by_org", func(t *testing.T) {
		var count int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM estimates WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-a should have 1 estimate")
	})

	t.Run("GetEstimate_respects_org_boundary", func(t *testing.T) {
		var id string
		err := db.QueryRowContext(ctx, `
			SELECT id FROM estimates WHERE id = $1 AND org_id = $2
		`, "est-b1", orgA).Scan(&id)
		assert.ErrorIs(t, err, sql.ErrNoRows, "org-a should not access org-b's estimate")
	})

	t.Run("Estimate_inherits_org_from_job", func(t *testing.T) {
		var estOrg, jobOrg uuid.UUID
		err := db.QueryRowContext(ctx, `
			SELECT e.org_id, j.org_id
			FROM estimates e
			JOIN jobs j ON e.job_id = j.id
			WHERE e.id = $1
		`, "est-a1").Scan(&estOrg, &jobOrg)
		require.NoError(t, err)
		assert.Equal(t, jobOrg, estOrg, "estimate org_id should match job org_id")
	})
}

// TestTenantIsolation_CrossOrgPrevention verifies that cross-org data manipulation is prevented.
// NOTE: These tests require database constraints (CHECK or triggers) that aren't implemented yet.
// They are skipped until cross-org prevention constraints are added to the schema.
func TestTenantIsolation_CrossOrgPrevention(t *testing.T) {
	t.Skip("Cross-org prevention constraints not yet implemented in database schema. See planning/saas-roadmap.md for implementation details.")

	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Generate UUIDs for test organizations
	orgA := uuid.New()
	orgB := uuid.New()

	// Create the organizations first (required by foreign key)
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES
			($1, 'Organization A', 'org-a-cross', NOW()),
			($2, 'Organization B', 'org-b-cross', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	// Setup: Create data for both orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('job-a', $1, 'Org A Job', 15.0, 'stacking', 'draft', NOW()),
			('job-b', $2, 'Org B Job', 10.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO categories (id, org_id, job_id, name, sort_order)
		VALUES
			('cat-a', $1, 'job-a', 'Org A Category', 0),
			('cat-b', $2, 'job-b', 'Org B Category', 0)
	`, orgA, orgB)
	require.NoError(t, err)

	t.Run("Cannot_attach_category_to_other_orgs_job", func(t *testing.T) {
		// Try to create a category for org-a but attach to org-b's job
		// This should be prevented by a CHECK constraint or trigger
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES ('malicious-cat', $1, 'job-b', 'Malicious Category', 0)
		`, orgA)
		// This should fail due to org_id mismatch constraint
		assert.Error(t, err, "should not be able to attach category to job with different org_id")
	})

	t.Run("Cannot_attach_line_item_to_other_orgs_category", func(t *testing.T) {
		// Try to create a line item for org-a but attach to org-b's category
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES ('malicious-item', $1, 'cat-b', 'material', 'Malicious Item', 1, 'ea', 100, 0)
		`, orgA)
		// This should fail due to org_id mismatch constraint
		assert.Error(t, err, "should not be able to attach line item to category with different org_id")
	})

	t.Run("Cannot_link_job_to_other_orgs_client", func(t *testing.T) {
		// Create clients for each org
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES ('client-b', $1, 'Org B Client', NOW())
		`, orgB)
		require.NoError(t, err)

		// Try to link org-a's job to org-b's client
		_, err = db.ExecContext(ctx, `
			UPDATE jobs SET client_id = 'client-b' WHERE id = 'job-a' AND org_id = $1
		`, orgA)
		// This should fail due to cross-org foreign key constraint
		assert.Error(t, err, "should not be able to link job to client with different org_id")
	})
}
