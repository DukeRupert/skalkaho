//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTenantCRUD_JobLifecycle tests complete job CRUD operations with tenant isolation.
func TestTenantCRUD_JobLifecycle(t *testing.T) {
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
			($1, 'Organization A', 'org-a', NOW()),
			($2, 'Organization B', 'org-b', NOW())
	`, orgA, orgB)
	require.NoError(t, err, "should be able to create organizations")

	t.Run("Create_job_with_org_id", func(t *testing.T) {
		// Create a job with org_id
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES ('job-1', $1, 'Kitchen Remodel', 15.0, 'stacking', 'draft', NOW())
		`, orgA)
		require.NoError(t, err, "should be able to create job with org_id")

		// Verify org_id was stored
		var orgID uuid.UUID
		err = db.QueryRowContext(ctx, `SELECT org_id FROM jobs WHERE id = $1`, "job-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, orgA, orgID, "job should have org_id set")
	})

	t.Run("Read_job_requires_org_id_filter", func(t *testing.T) {
		// Setup: create jobs for two orgs
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES
				('job-read-a', $1, 'Org A Job', 10.0, 'stacking', 'draft', NOW()),
				('job-read-b', $2, 'Org B Job', 20.0, 'stacking', 'draft', NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Reading with org filter should only return that org's job
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		// Note: count includes job-1 from previous test
		assert.GreaterOrEqual(t, count, int64(1), "org-a should have at least 1 job")

		// Verify we can't see org-b's job when filtering by org-a
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM jobs WHERE id = $1 AND org_id = $2`, "job-read-b", orgA).Scan(&id)
		assert.Error(t, err, "org-a should not be able to read org-b's job")
	})

	t.Run("Update_job_scoped_by_org_id", func(t *testing.T) {
		// Setup
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES
				('job-update-a', $1, 'Original A', 10.0, 'stacking', 'draft', NOW()),
				('job-update-b', $2, 'Original B', 20.0, 'stacking', 'draft', NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Update own org's job - should succeed
		result, err := db.ExecContext(ctx, `
			UPDATE jobs SET name = 'Updated A' WHERE id = $1 AND org_id = $2
		`, "job-update-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected, "should update own org's job")

		// Attempt to update other org's job - should affect 0 rows
		result, err = db.ExecContext(ctx, `
			UPDATE jobs SET name = 'Hacked B' WHERE id = $1 AND org_id = $2
		`, "job-update-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not update other org's job")

		// Verify org-b's job was not modified
		var name string
		err = db.QueryRowContext(ctx, `SELECT name FROM jobs WHERE id = $1`, "job-update-b").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Original B", name, "org-b job should not be modified")
	})

	t.Run("Delete_job_scoped_by_org_id", func(t *testing.T) {
		// Setup
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES
				('job-delete-a', $1, 'To Delete A', 10.0, 'stacking', 'draft', NOW()),
				('job-delete-b', $2, 'To Delete B', 20.0, 'stacking', 'draft', NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Delete own org's job - should succeed
		result, err := db.ExecContext(ctx, `
			DELETE FROM jobs WHERE id = $1 AND org_id = $2
		`, "job-delete-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected, "should delete own org's job")

		// Verify it's deleted
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE id = $1`, "job-delete-a").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "job should be deleted")

		// Attempt to delete other org's job - should affect 0 rows
		result, err = db.ExecContext(ctx, `
			DELETE FROM jobs WHERE id = $1 AND org_id = $2
		`, "job-delete-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not delete other org's job")

		// Verify org-b's job still exists
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE id = $1`, "job-delete-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-b job should still exist")
	})
}

// TestTenantCRUD_CategoryOperations tests category CRUD with org isolation through job hierarchy.
func TestTenantCRUD_CategoryOperations(t *testing.T) {
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

	// Setup: Create jobs for both orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('cat-test-job-a', $1, 'Org A Job', 10.0, 'stacking', 'draft', NOW()),
			('cat-test-job-b', $2, 'Org B Job', 20.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	t.Run("Create_category_inherits_org_from_job", func(t *testing.T) {
		// Create category with matching org_id
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES ('cat-inherit-1', $1, 'cat-test-job-a', 'Foundation', 0)
		`, orgA)
		require.NoError(t, err)

		// Verify org_id consistency
		var catOrg, jobOrg string
		err = db.QueryRowContext(ctx, `
			SELECT c.org_id, j.org_id
			FROM categories c
			JOIN jobs j ON c.job_id = j.id
			WHERE c.id = $1
		`, "cat-inherit-1").Scan(&catOrg, &jobOrg)
		require.NoError(t, err)
		assert.Equal(t, jobOrg, catOrg, "category org_id should match job org_id")
	})

	t.Run("Create_nested_categories_all_same_org", func(t *testing.T) {
		// Create 3-level nested categories
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, parent_id, name, sort_order)
			VALUES
				('cat-l1', $1, 'cat-test-job-a', NULL, 'Level 1', 0),
				('cat-l2', $1, 'cat-test-job-a', 'cat-l1', 'Level 2', 0),
				('cat-l3', $1, 'cat-test-job-a', 'cat-l2', 'Level 3', 0)
		`, orgA)
		require.NoError(t, err)

		// Verify all have same org_id
		rows, err := db.QueryContext(ctx, `SELECT org_id FROM categories WHERE id IN ('cat-l1', 'cat-l2', 'cat-l3')`)
		require.NoError(t, err)
		defer rows.Close()

		for rows.Next() {
			var orgID uuid.UUID
			require.NoError(t, rows.Scan(&orgID))
			assert.Equal(t, orgA, orgID, "all nested categories should have same org_id")
		}
	})

	t.Run("Update_category_scoped_by_org", func(t *testing.T) {
		// Setup
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES
				('cat-upd-a', $1, 'cat-test-job-a', 'Original A', 0),
				('cat-upd-b', $2, 'cat-test-job-b', 'Original B', 0)
		`, orgA, orgB)
		require.NoError(t, err)

		// Update own org's category
		result, err := db.ExecContext(ctx, `
			UPDATE categories SET name = 'Updated A' WHERE id = $1 AND org_id = $2
		`, "cat-upd-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's category
		result, err = db.ExecContext(ctx, `
			UPDATE categories SET name = 'Hacked B' WHERE id = $1 AND org_id = $2
		`, "cat-upd-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not update other org's category")
	})

	t.Run("Delete_category_cascades_within_org", func(t *testing.T) {
		// Setup: category with line items
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES ('cat-cascade', $1, 'cat-test-job-a', 'To Cascade', 0)
		`, orgA)
		require.NoError(t, err)

		_, err = db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('item-cascade-1', $1, 'cat-cascade', 'material', 'Item 1', 1, 'ea', 100, 0),
				('item-cascade-2', $1, 'cat-cascade', 'labor', 'Item 2', 2, 'hr', 50, 1)
		`, orgA)
		require.NoError(t, err)

		// Delete category - should cascade to line items
		_, err = db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1 AND org_id = $2`, "cat-cascade", orgA)
		require.NoError(t, err)

		// Verify line items were cascaded
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE category_id = $1`, "cat-cascade").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "line items should be cascade deleted")
	})
}

// TestTenantCRUD_LineItemOperations tests line item CRUD with org isolation.
func TestTenantCRUD_LineItemOperations(t *testing.T) {
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

	// Setup: Create hierarchy for both orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('li-test-job-a', $1, 'Org A Job', 10.0, 'stacking', 'draft', NOW()),
			('li-test-job-b', $2, 'Org B Job', 20.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO categories (id, org_id, job_id, name, sort_order)
		VALUES
			('li-test-cat-a', $1, 'li-test-job-a', 'Org A Category', 0),
			('li-test-cat-b', $2, 'li-test-job-b', 'Org B Category', 0)
	`, orgA, orgB)
	require.NoError(t, err)

	t.Run("Create_line_item_with_org_id", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES ('li-create-1', $1, 'li-test-cat-a', 'material', '2x4 Lumber', 100, 'ea', 5.50, 0)
		`, orgA)
		require.NoError(t, err)

		// Verify org_id was set
		var orgID uuid.UUID
		err = db.QueryRowContext(ctx, `SELECT org_id FROM line_items WHERE id = $1`, "li-create-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, orgA, orgID)
	})

	t.Run("Read_line_items_filtered_by_org", func(t *testing.T) {
		// Create items for both orgs
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('li-read-a', $1, 'li-test-cat-a', 'material', 'Org A Item', 1, 'ea', 100, 0),
				('li-read-b', $2, 'li-test-cat-b', 'material', 'Org B Item', 1, 'ea', 200, 0)
		`, orgA, orgB)
		require.NoError(t, err)

		// Query with org filter
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1), "org-a should have line items")

		// Verify can't access other org's item
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM line_items WHERE id = $1 AND org_id = $2`, "li-read-b", orgA).Scan(&id)
		assert.Error(t, err, "org-a should not access org-b's line item")
	})

	t.Run("Update_line_item_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('li-upd-a', $1, 'li-test-cat-a', 'material', 'Original A', 1, 'ea', 100, 0),
				('li-upd-b', $2, 'li-test-cat-b', 'material', 'Original B', 1, 'ea', 200, 0)
		`, orgA, orgB)
		require.NoError(t, err)

		// Update own org's item
		result, err := db.ExecContext(ctx, `
			UPDATE line_items SET name = 'Updated A', unit_price = 150 WHERE id = $1 AND org_id = $2
		`, "li-upd-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's item
		result, err = db.ExecContext(ctx, `
			UPDATE line_items SET name = 'Hacked B' WHERE id = $1 AND org_id = $2
		`, "li-upd-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not update other org's line item")

		// Verify original value preserved
		var name string
		err = db.QueryRowContext(ctx, `SELECT name FROM line_items WHERE id = $1`, "li-upd-b").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Original B", name)
	})

	t.Run("Delete_line_item_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('li-del-a', $1, 'li-test-cat-a', 'material', 'To Delete A', 1, 'ea', 100, 0),
				('li-del-b', $2, 'li-test-cat-b', 'material', 'To Delete B', 1, 'ea', 200, 0)
		`, orgA, orgB)
		require.NoError(t, err)

		// Delete own org's item
		result, err := db.ExecContext(ctx, `
			DELETE FROM line_items WHERE id = $1 AND org_id = $2
		`, "li-del-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to delete other org's item
		result, err = db.ExecContext(ctx, `
			DELETE FROM line_items WHERE id = $1 AND org_id = $2
		`, "li-del-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not delete other org's line item")

		// Verify org-b's item still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE id = $1`, "li-del-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-b item should still exist")
	})
}

// TestTenantCRUD_ClientOperations tests client CRUD with org isolation.
func TestTenantCRUD_ClientOperations(t *testing.T) {
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

	t.Run("Create_client_with_org_id", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, email, created_at)
			VALUES ('client-create-1', $1, 'John Doe', 'john@example.com', NOW())
		`, orgA)
		require.NoError(t, err)

		var orgID uuid.UUID
		err = db.QueryRowContext(ctx, `SELECT org_id FROM clients WHERE id = $1`, "client-create-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, orgA, orgID)
	})

	t.Run("List_clients_filtered_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES
				('client-list-a', $1, 'Org A Client', NOW()),
				('client-list-b', $2, 'Org B Client', NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM clients WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))

		// Verify isolation
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM clients WHERE id = $1 AND org_id = $2`, "client-list-b", orgA).Scan(&id)
		assert.Error(t, err, "org-a should not access org-b's client")
	})

	t.Run("Update_client_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES
				('client-upd-a', $1, 'Original A', NOW()),
				('client-upd-b', $2, 'Original B', NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Update own client
		result, err := db.ExecContext(ctx, `UPDATE clients SET name = 'Updated A' WHERE id = $1 AND org_id = $2`, "client-upd-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's client
		result, err = db.ExecContext(ctx, `UPDATE clients SET name = 'Hacked B' WHERE id = $1 AND org_id = $2`, "client-upd-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)
	})

	t.Run("Delete_client_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES
				('client-del-a', $1, 'To Delete A', NOW()),
				('client-del-b', $2, 'To Delete B', NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Delete own client
		result, err := db.ExecContext(ctx, `DELETE FROM clients WHERE id = $1 AND org_id = $2`, "client-del-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to delete other org's client
		result, err = db.ExecContext(ctx, `DELETE FROM clients WHERE id = $1 AND org_id = $2`, "client-del-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)

		// Verify org-b's client still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM clients WHERE id = $1`, "client-del-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

// TestTenantCRUD_EstimateOperations tests estimate CRUD with org isolation.
func TestTenantCRUD_EstimateOperations(t *testing.T) {
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

	// Setup: Create jobs for both orgs
	_, err = db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('est-test-job-a', $1, 'Org A Job', 10.0, 'stacking', 'draft', NOW()),
			('est-test-job-b', $2, 'Org B Job', 20.0, 'stacking', 'draft', NOW())
	`, orgA, orgB)
	require.NoError(t, err)

	t.Run("Create_estimate_with_org_id", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES ('est-create-1', $1, 'est-test-job-a', 1, 'draft', 5000.00, NOW())
		`, orgA)
		require.NoError(t, err)

		var orgID uuid.UUID
		err = db.QueryRowContext(ctx, `SELECT org_id FROM estimates WHERE id = $1`, "est-create-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, orgA, orgID)
	})

	t.Run("List_estimates_by_job_filtered_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES
				('est-list-a', $1, 'est-test-job-a', 2, 'draft', 6000.00, NOW()),
				('est-list-b', $2, 'est-test-job-b', 1, 'draft', 8000.00, NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Query estimates with org filter
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM estimates WHERE org_id = $1`, orgA).Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))

		// Verify isolation
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM estimates WHERE id = $1 AND org_id = $2`, "est-list-b", orgA).Scan(&id)
		assert.Error(t, err, "org-a should not access org-b's estimate")
	})

	t.Run("Update_estimate_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES
				('est-upd-a', $1, 'est-test-job-a', 3, 'draft', 7000.00, NOW()),
				('est-upd-b', $2, 'est-test-job-b', 2, 'draft', 9000.00, NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Update own estimate
		result, err := db.ExecContext(ctx, `
			UPDATE estimates SET status = 'sent', grand_total = 7500.00 WHERE id = $1 AND org_id = $2
		`, "est-upd-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's estimate
		result, err = db.ExecContext(ctx, `
			UPDATE estimates SET status = 'accepted' WHERE id = $1 AND org_id = $2
		`, "est-upd-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)

		// Verify original status preserved
		var status string
		err = db.QueryRowContext(ctx, `SELECT status FROM estimates WHERE id = $1`, "est-upd-b").Scan(&status)
		require.NoError(t, err)
		assert.Equal(t, "draft", status)
	})

	t.Run("Delete_estimate_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES
				('est-del-a', $1, 'est-test-job-a', 4, 'draft', 1000.00, NOW()),
				('est-del-b', $2, 'est-test-job-b', 3, 'draft', 2000.00, NOW())
		`, orgA, orgB)
		require.NoError(t, err)

		// Delete own estimate
		result, err := db.ExecContext(ctx, `DELETE FROM estimates WHERE id = $1 AND org_id = $2`, "est-del-a", orgA)
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to delete other org's estimate
		result, err = db.ExecContext(ctx, `DELETE FROM estimates WHERE id = $1 AND org_id = $2`, "est-del-b", orgA)
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)

		// Verify org-b's estimate still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM estimates WHERE id = $1`, "est-del-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}
