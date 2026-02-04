//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTenantCRUD_JobLifecycle tests complete job CRUD operations with tenant isolation.
func TestTenantCRUD_JobLifecycle(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("Create_job_with_org_id", func(t *testing.T) {
		// Create a job with org_id
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES ('job-1', 'org-a', 'Kitchen Remodel', 15.0, 'stacking', 'draft', datetime('now'))
		`)
		require.NoError(t, err, "should be able to create job with org_id")

		// Verify org_id was stored
		var orgID string
		err = db.QueryRowContext(ctx, `SELECT org_id FROM jobs WHERE id = ?`, "job-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, "org-a", orgID, "job should have org_id set")
	})

	t.Run("Read_job_requires_org_id_filter", func(t *testing.T) {
		// Setup: create jobs for two orgs
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES
				('job-read-a', 'org-a', 'Org A Job', 10.0, 'stacking', 'draft', datetime('now')),
				('job-read-b', 'org-b', 'Org B Job', 20.0, 'stacking', 'draft', datetime('now'))
		`)
		require.NoError(t, err)

		// Reading with org filter should only return that org's job
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE org_id = ?`, "org-a").Scan(&count)
		require.NoError(t, err)
		// Note: count includes job-1 from previous test
		assert.GreaterOrEqual(t, count, int64(1), "org-a should have at least 1 job")

		// Verify we can't see org-b's job when filtering by org-a
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM jobs WHERE id = ? AND org_id = ?`, "job-read-b", "org-a").Scan(&id)
		assert.Error(t, err, "org-a should not be able to read org-b's job")
	})

	t.Run("Update_job_scoped_by_org_id", func(t *testing.T) {
		// Setup
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES
				('job-update-a', 'org-a', 'Original A', 10.0, 'stacking', 'draft', datetime('now')),
				('job-update-b', 'org-b', 'Original B', 20.0, 'stacking', 'draft', datetime('now'))
		`)
		require.NoError(t, err)

		// Update own org's job - should succeed
		result, err := db.ExecContext(ctx, `
			UPDATE jobs SET name = 'Updated A' WHERE id = ? AND org_id = ?
		`, "job-update-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected, "should update own org's job")

		// Attempt to update other org's job - should affect 0 rows
		result, err = db.ExecContext(ctx, `
			UPDATE jobs SET name = 'Hacked B' WHERE id = ? AND org_id = ?
		`, "job-update-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not update other org's job")

		// Verify org-b's job was not modified
		var name string
		err = db.QueryRowContext(ctx, `SELECT name FROM jobs WHERE id = ?`, "job-update-b").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Original B", name, "org-b job should not be modified")
	})

	t.Run("Delete_job_scoped_by_org_id", func(t *testing.T) {
		// Setup
		_, err := db.ExecContext(ctx, `
			INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
			VALUES
				('job-delete-a', 'org-a', 'To Delete A', 10.0, 'stacking', 'draft', datetime('now')),
				('job-delete-b', 'org-b', 'To Delete B', 20.0, 'stacking', 'draft', datetime('now'))
		`)
		require.NoError(t, err)

		// Delete own org's job - should succeed
		result, err := db.ExecContext(ctx, `
			DELETE FROM jobs WHERE id = ? AND org_id = ?
		`, "job-delete-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected, "should delete own org's job")

		// Verify it's deleted
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE id = ?`, "job-delete-a").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "job should be deleted")

		// Attempt to delete other org's job - should affect 0 rows
		result, err = db.ExecContext(ctx, `
			DELETE FROM jobs WHERE id = ? AND org_id = ?
		`, "job-delete-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not delete other org's job")

		// Verify org-b's job still exists
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE id = ?`, "job-delete-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-b job should still exist")
	})
}

// TestTenantCRUD_CategoryOperations tests category CRUD with org isolation through job hierarchy.
func TestTenantCRUD_CategoryOperations(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Setup: Create jobs for both orgs
	_, err := db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('cat-test-job-a', 'org-a', 'Org A Job', 10.0, 'stacking', 'draft', datetime('now')),
			('cat-test-job-b', 'org-b', 'Org B Job', 20.0, 'stacking', 'draft', datetime('now'))
	`)
	require.NoError(t, err)

	t.Run("Create_category_inherits_org_from_job", func(t *testing.T) {
		// Create category with matching org_id
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES ('cat-inherit-1', 'org-a', 'cat-test-job-a', 'Foundation', 0)
		`)
		require.NoError(t, err)

		// Verify org_id consistency
		var catOrg, jobOrg string
		err = db.QueryRowContext(ctx, `
			SELECT c.org_id, j.org_id
			FROM categories c
			JOIN jobs j ON c.job_id = j.id
			WHERE c.id = ?
		`, "cat-inherit-1").Scan(&catOrg, &jobOrg)
		require.NoError(t, err)
		assert.Equal(t, jobOrg, catOrg, "category org_id should match job org_id")
	})

	t.Run("Create_nested_categories_all_same_org", func(t *testing.T) {
		// Create 3-level nested categories
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, parent_id, name, sort_order)
			VALUES
				('cat-l1', 'org-a', 'cat-test-job-a', NULL, 'Level 1', 0),
				('cat-l2', 'org-a', 'cat-test-job-a', 'cat-l1', 'Level 2', 0),
				('cat-l3', 'org-a', 'cat-test-job-a', 'cat-l2', 'Level 3', 0)
		`)
		require.NoError(t, err)

		// Verify all have same org_id
		rows, err := db.QueryContext(ctx, `SELECT org_id FROM categories WHERE id IN ('cat-l1', 'cat-l2', 'cat-l3')`)
		require.NoError(t, err)
		defer rows.Close()

		for rows.Next() {
			var orgID string
			require.NoError(t, rows.Scan(&orgID))
			assert.Equal(t, "org-a", orgID, "all nested categories should have same org_id")
		}
	})

	t.Run("Update_category_scoped_by_org", func(t *testing.T) {
		// Setup
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES
				('cat-upd-a', 'org-a', 'cat-test-job-a', 'Original A', 0),
				('cat-upd-b', 'org-b', 'cat-test-job-b', 'Original B', 0)
		`)
		require.NoError(t, err)

		// Update own org's category
		result, err := db.ExecContext(ctx, `
			UPDATE categories SET name = 'Updated A' WHERE id = ? AND org_id = ?
		`, "cat-upd-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's category
		result, err = db.ExecContext(ctx, `
			UPDATE categories SET name = 'Hacked B' WHERE id = ? AND org_id = ?
		`, "cat-upd-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not update other org's category")
	})

	t.Run("Delete_category_cascades_within_org", func(t *testing.T) {
		// Setup: category with line items
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (id, org_id, job_id, name, sort_order)
			VALUES ('cat-cascade', 'org-a', 'cat-test-job-a', 'To Cascade', 0)
		`)
		require.NoError(t, err)

		_, err = db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('item-cascade-1', 'org-a', 'cat-cascade', 'material', 'Item 1', 1, 'ea', 100, 0),
				('item-cascade-2', 'org-a', 'cat-cascade', 'labor', 'Item 2', 2, 'hr', 50, 1)
		`)
		require.NoError(t, err)

		// Delete category - should cascade to line items
		_, err = db.ExecContext(ctx, `DELETE FROM categories WHERE id = ? AND org_id = ?`, "cat-cascade", "org-a")
		require.NoError(t, err)

		// Verify line items were cascaded
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE category_id = ?`, "cat-cascade").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "line items should be cascade deleted")
	})
}

// TestTenantCRUD_LineItemOperations tests line item CRUD with org isolation.
func TestTenantCRUD_LineItemOperations(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Setup: Create hierarchy for both orgs
	_, err := db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('li-test-job-a', 'org-a', 'Org A Job', 10.0, 'stacking', 'draft', datetime('now')),
			('li-test-job-b', 'org-b', 'Org B Job', 20.0, 'stacking', 'draft', datetime('now'))
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO categories (id, org_id, job_id, name, sort_order)
		VALUES
			('li-test-cat-a', 'org-a', 'li-test-job-a', 'Org A Category', 0),
			('li-test-cat-b', 'org-b', 'li-test-job-b', 'Org B Category', 0)
	`)
	require.NoError(t, err)

	t.Run("Create_line_item_with_org_id", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES ('li-create-1', 'org-a', 'li-test-cat-a', 'material', '2x4 Lumber', 100, 'ea', 5.50, 0)
		`)
		require.NoError(t, err)

		// Verify org_id was set
		var orgID string
		err = db.QueryRowContext(ctx, `SELECT org_id FROM line_items WHERE id = ?`, "li-create-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, "org-a", orgID)
	})

	t.Run("Read_line_items_filtered_by_org", func(t *testing.T) {
		// Create items for both orgs
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('li-read-a', 'org-a', 'li-test-cat-a', 'material', 'Org A Item', 1, 'ea', 100, 0),
				('li-read-b', 'org-b', 'li-test-cat-b', 'material', 'Org B Item', 1, 'ea', 200, 0)
		`)
		require.NoError(t, err)

		// Query with org filter
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE org_id = ?`, "org-a").Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1), "org-a should have line items")

		// Verify can't access other org's item
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM line_items WHERE id = ? AND org_id = ?`, "li-read-b", "org-a").Scan(&id)
		assert.Error(t, err, "org-a should not access org-b's line item")
	})

	t.Run("Update_line_item_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('li-upd-a', 'org-a', 'li-test-cat-a', 'material', 'Original A', 1, 'ea', 100, 0),
				('li-upd-b', 'org-b', 'li-test-cat-b', 'material', 'Original B', 1, 'ea', 200, 0)
		`)
		require.NoError(t, err)

		// Update own org's item
		result, err := db.ExecContext(ctx, `
			UPDATE line_items SET name = 'Updated A', unit_price = 150 WHERE id = ? AND org_id = ?
		`, "li-upd-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's item
		result, err = db.ExecContext(ctx, `
			UPDATE line_items SET name = 'Hacked B' WHERE id = ? AND org_id = ?
		`, "li-upd-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not update other org's line item")

		// Verify original value preserved
		var name string
		err = db.QueryRowContext(ctx, `SELECT name FROM line_items WHERE id = ?`, "li-upd-b").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Original B", name)
	})

	t.Run("Delete_line_item_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO line_items (id, org_id, category_id, type, name, quantity, unit, unit_price, sort_order)
			VALUES
				('li-del-a', 'org-a', 'li-test-cat-a', 'material', 'To Delete A', 1, 'ea', 100, 0),
				('li-del-b', 'org-b', 'li-test-cat-b', 'material', 'To Delete B', 1, 'ea', 200, 0)
		`)
		require.NoError(t, err)

		// Delete own org's item
		result, err := db.ExecContext(ctx, `
			DELETE FROM line_items WHERE id = ? AND org_id = ?
		`, "li-del-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to delete other org's item
		result, err = db.ExecContext(ctx, `
			DELETE FROM line_items WHERE id = ? AND org_id = ?
		`, "li-del-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected, "should not delete other org's line item")

		// Verify org-b's item still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM line_items WHERE id = ?`, "li-del-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "org-b item should still exist")
	})
}

// TestTenantCRUD_ClientOperations tests client CRUD with org isolation.
func TestTenantCRUD_ClientOperations(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("Create_client_with_org_id", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, email, created_at)
			VALUES ('client-create-1', 'org-a', 'John Doe', 'john@example.com', datetime('now'))
		`)
		require.NoError(t, err)

		var orgID string
		err = db.QueryRowContext(ctx, `SELECT org_id FROM clients WHERE id = ?`, "client-create-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, "org-a", orgID)
	})

	t.Run("List_clients_filtered_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES
				('client-list-a', 'org-a', 'Org A Client', datetime('now')),
				('client-list-b', 'org-b', 'Org B Client', datetime('now'))
		`)
		require.NoError(t, err)

		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM clients WHERE org_id = ?`, "org-a").Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))

		// Verify isolation
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM clients WHERE id = ? AND org_id = ?`, "client-list-b", "org-a").Scan(&id)
		assert.Error(t, err, "org-a should not access org-b's client")
	})

	t.Run("Update_client_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES
				('client-upd-a', 'org-a', 'Original A', datetime('now')),
				('client-upd-b', 'org-b', 'Original B', datetime('now'))
		`)
		require.NoError(t, err)

		// Update own client
		result, err := db.ExecContext(ctx, `UPDATE clients SET name = 'Updated A' WHERE id = ? AND org_id = ?`, "client-upd-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's client
		result, err = db.ExecContext(ctx, `UPDATE clients SET name = 'Hacked B' WHERE id = ? AND org_id = ?`, "client-upd-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)
	})

	t.Run("Delete_client_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO clients (id, org_id, name, created_at)
			VALUES
				('client-del-a', 'org-a', 'To Delete A', datetime('now')),
				('client-del-b', 'org-b', 'To Delete B', datetime('now'))
		`)
		require.NoError(t, err)

		// Delete own client
		result, err := db.ExecContext(ctx, `DELETE FROM clients WHERE id = ? AND org_id = ?`, "client-del-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to delete other org's client
		result, err = db.ExecContext(ctx, `DELETE FROM clients WHERE id = ? AND org_id = ?`, "client-del-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)

		// Verify org-b's client still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM clients WHERE id = ?`, "client-del-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

// TestTenantCRUD_EstimateOperations tests estimate CRUD with org isolation.
func TestTenantCRUD_EstimateOperations(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Setup: Create jobs for both orgs
	_, err := db.ExecContext(ctx, `
		INSERT INTO jobs (id, org_id, name, surcharge_percent, surcharge_mode, status, created_at)
		VALUES
			('est-test-job-a', 'org-a', 'Org A Job', 10.0, 'stacking', 'draft', datetime('now')),
			('est-test-job-b', 'org-b', 'Org B Job', 20.0, 'stacking', 'draft', datetime('now'))
	`)
	require.NoError(t, err)

	t.Run("Create_estimate_with_org_id", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES ('est-create-1', 'org-a', 'est-test-job-a', 1, 'draft', 5000.00, datetime('now'))
		`)
		require.NoError(t, err)

		var orgID string
		err = db.QueryRowContext(ctx, `SELECT org_id FROM estimates WHERE id = ?`, "est-create-1").Scan(&orgID)
		require.NoError(t, err)
		assert.Equal(t, "org-a", orgID)
	})

	t.Run("List_estimates_by_job_filtered_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES
				('est-list-a', 'org-a', 'est-test-job-a', 2, 'draft', 6000.00, datetime('now')),
				('est-list-b', 'org-b', 'est-test-job-b', 1, 'draft', 8000.00, datetime('now'))
		`)
		require.NoError(t, err)

		// Query estimates with org filter
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM estimates WHERE org_id = ?`, "org-a").Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))

		// Verify isolation
		var id string
		err = db.QueryRowContext(ctx, `SELECT id FROM estimates WHERE id = ? AND org_id = ?`, "est-list-b", "org-a").Scan(&id)
		assert.Error(t, err, "org-a should not access org-b's estimate")
	})

	t.Run("Update_estimate_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES
				('est-upd-a', 'org-a', 'est-test-job-a', 3, 'draft', 7000.00, datetime('now')),
				('est-upd-b', 'org-b', 'est-test-job-b', 2, 'draft', 9000.00, datetime('now'))
		`)
		require.NoError(t, err)

		// Update own estimate
		result, err := db.ExecContext(ctx, `
			UPDATE estimates SET status = 'sent', grand_total = 7500.00 WHERE id = ? AND org_id = ?
		`, "est-upd-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to update other org's estimate
		result, err = db.ExecContext(ctx, `
			UPDATE estimates SET status = 'accepted' WHERE id = ? AND org_id = ?
		`, "est-upd-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)

		// Verify original status preserved
		var status string
		err = db.QueryRowContext(ctx, `SELECT status FROM estimates WHERE id = ?`, "est-upd-b").Scan(&status)
		require.NoError(t, err)
		assert.Equal(t, "draft", status)
	})

	t.Run("Delete_estimate_scoped_by_org", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `
			INSERT INTO estimates (id, org_id, job_id, version, status, grand_total, created_at)
			VALUES
				('est-del-a', 'org-a', 'est-test-job-a', 4, 'draft', 1000.00, datetime('now')),
				('est-del-b', 'org-b', 'est-test-job-b', 3, 'draft', 2000.00, datetime('now'))
		`)
		require.NoError(t, err)

		// Delete own estimate
		result, err := db.ExecContext(ctx, `DELETE FROM estimates WHERE id = ? AND org_id = ?`, "est-del-a", "org-a")
		require.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)

		// Attempt to delete other org's estimate
		result, err = db.ExecContext(ctx, `DELETE FROM estimates WHERE id = ? AND org_id = ?`, "est-del-b", "org-a")
		require.NoError(t, err)
		rowsAffected, _ = result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)

		// Verify org-b's estimate still exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM estimates WHERE id = ?`, "est-del-b").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}
