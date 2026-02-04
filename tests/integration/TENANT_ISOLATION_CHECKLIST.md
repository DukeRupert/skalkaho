# Tenant Data Isolation - Manual Verification Checklist

This checklist provides manual verification steps to complement the automated integration tests for tenant data isolation in Skalkaho.

## Prerequisites

Before beginning verification:
- [ ] Multi-tenancy schema changes applied (org_id columns added to all tables)
- [ ] Organization context middleware implemented
- [ ] All SQL queries updated to filter by org_id
- [ ] Automated tests passing

## Database Schema Verification

### Check org_id Columns

```sql
-- Verify org_id column exists on all tenant-scoped tables
PRAGMA table_info(jobs);           -- Should show org_id column
PRAGMA table_info(categories);     -- Should show org_id column
PRAGMA table_info(line_items);     -- Should show org_id column
PRAGMA table_info(settings);       -- Should show org_id column
PRAGMA table_info(clients);        -- Should show org_id column
PRAGMA table_info(estimates);      -- Should show org_id column (through jobs)
PRAGMA table_info(signature_requests); -- Should show org_id column (through estimates)
```

### Check Indexes

```sql
-- Verify indexes for performance
.indexes jobs           -- Should include index on org_id
.indexes categories     -- Should include composite index on (org_id, job_id)
.indexes line_items     -- Should include composite index on (org_id, category_id)
```

### Check Foreign Key Constraints

```sql
-- Verify foreign keys maintain org boundaries
PRAGMA foreign_key_list(categories);    -- Should reference jobs(id)
PRAGMA foreign_key_list(line_items);    -- Should reference categories(id)
```

## Functional Testing

### Test Setup: Create Test Data

Create two test organizations with complete data hierarchies:

**Organization A:**
- 2 jobs
- 3 categories per job
- 5 line items per category
- 1 settings record
- 2 clients

**Organization B:**
- 2 jobs
- 3 categories per job
- 5 line items per category
- 1 settings record
- 2 clients

### Jobs Isolation

- [ ] Log in as Organization A
- [ ] Navigate to jobs list
- [ ] Verify only Organization A's jobs are visible (should see 2 jobs)
- [ ] Note the IDs of visible jobs
- [ ] Log out

- [ ] Log in as Organization B
- [ ] Navigate to jobs list
- [ ] Verify only Organization B's jobs are visible (should see 2 jobs)
- [ ] Attempt to access Organization A's job by URL: `/jobs/{org_a_job_id}`
- [ ] Should receive 404 Not Found error
- [ ] Log out

### Categories Isolation

- [ ] Log in as Organization A
- [ ] Open one of Organization A's jobs
- [ ] Verify categories belong to this job only
- [ ] Count categories (should see 3)
- [ ] Attempt to create category under Organization B's job (requires API call with org B's job_id)
- [ ] Should receive validation error or 403 Forbidden
- [ ] Log out

### Line Items Isolation

- [ ] Log in as Organization A
- [ ] Navigate to a job and expand a category
- [ ] Verify line items are only from this category
- [ ] Count line items (should see 5)
- [ ] Attempt to move line item to Organization B's category (requires API manipulation)
- [ ] Should receive validation error
- [ ] Log out

### Settings Isolation

- [ ] Log in as Organization A
- [ ] Navigate to Settings page
- [ ] Update default surcharge to 15%
- [ ] Save settings
- [ ] Log out

- [ ] Log in as Organization B
- [ ] Navigate to Settings page
- [ ] Verify surcharge is NOT 15% (should be Organization B's original value)
- [ ] Update to 20%
- [ ] Save settings
- [ ] Log out

- [ ] Log in as Organization A
- [ ] Navigate to Settings page
- [ ] Verify surcharge is still 15% (unchanged by Organization B)
- [ ] Log out

### Clients Isolation

- [ ] Log in as Organization A
- [ ] Navigate to Clients list
- [ ] Verify only 2 clients visible (Organization A's)
- [ ] Note a client ID
- [ ] Log out

- [ ] Log in as Organization B
- [ ] Navigate to Clients list
- [ ] Verify only 2 clients visible (Organization B's)
- [ ] Attempt to access Organization A's client by URL: `/clients/{org_a_client_id}`
- [ ] Should receive 404 Not Found error
- [ ] Log out

### Estimates Isolation

- [ ] Log in as Organization A
- [ ] Create estimate for one of Organization A's jobs
- [ ] Note the estimate ID
- [ ] Log out

- [ ] Log in as Organization B
- [ ] Attempt to access Organization A's estimate by URL: `/estimates/{org_a_estimate_id}`
- [ ] Should receive 404 Not Found error
- [ ] Create estimate for Organization B's job
- [ ] Verify estimate is created successfully
- [ ] Log out

## CRUD Operations Testing

### Create Operations

- [ ] Log in as Organization A
- [ ] Create new job "Test Job A"
- [ ] Verify job appears in list
- [ ] Check database: `SELECT org_id FROM jobs WHERE name = 'Test Job A'`
- [ ] Verify org_id matches Organization A
- [ ] Create category under new job
- [ ] Check database: `SELECT org_id FROM categories ORDER BY id DESC LIMIT 1`
- [ ] Verify org_id matches Organization A
- [ ] Log out

### Read Operations

- [ ] Log in as Organization A
- [ ] Execute search for "Test"
- [ ] Verify only Organization A's results returned
- [ ] Check total count matches expected Organization A records
- [ ] Log out

### Update Operations

- [ ] Log in as Organization A
- [ ] Update "Test Job A" to "Modified Job A"
- [ ] Verify update succeeds
- [ ] Log out

- [ ] Log in as Organization B
- [ ] Attempt API call to update "Modified Job A"
- [ ] Should receive 403 Forbidden or 404 Not Found
- [ ] Check database: `SELECT name FROM jobs WHERE name LIKE 'Modified%'`
- [ ] Verify name is still "Modified Job A" (not changed by Organization B)
- [ ] Log out

### Delete Operations

- [ ] Log in as Organization A
- [ ] Delete "Modified Job A"
- [ ] Verify deletion succeeds
- [ ] Verify associated categories and line items also deleted (cascade)
- [ ] Log out

- [ ] Create job for Organization B via API or UI
- [ ] Log in as Organization A
- [ ] Attempt API call to delete Organization B's job
- [ ] Should receive 403 Forbidden or 404 Not Found
- [ ] Check database: `SELECT COUNT(*) FROM jobs WHERE org_id = 'org-b'`
- [ ] Verify count unchanged
- [ ] Log out

## Edge Cases and Security

### SQL Injection Testing

- [ ] Log in as Organization A
- [ ] Attempt to inject SQL in search: `' OR org_id != 'org-a' --`
- [ ] Verify no cross-org data leakage
- [ ] Should see error or no results

### Direct API Manipulation

- [ ] Log in as Organization A
- [ ] Get auth token from browser dev tools
- [ ] Use curl/Postman to send request to update Organization B's job:
  ```bash
  curl -X PUT /api/jobs/{org_b_job_id} \
    -H "Authorization: Bearer {org_a_token}" \
    -d '{"name": "Hacked"}'
  ```
- [ ] Should receive 403 Forbidden or 404 Not Found
- [ ] Verify Organization B's job unchanged in database

### Bulk Operations

- [ ] Log in as Organization A
- [ ] Select multiple jobs (if bulk delete feature exists)
- [ ] Include Organization A's job IDs in request
- [ ] Verify only Organization A's jobs are deleted
- [ ] Organization B's jobs remain untouched

### Cross-Org Foreign Key Attempts

Attempt to create data with foreign keys pointing to other org's records:

- [ ] Log in as Organization A
- [ ] Attempt API call to create category with Organization B's job_id
- [ ] Should receive validation error or 404 Not Found
- [ ] Attempt API call to create line item with Organization B's category_id
- [ ] Should receive validation error or 404 Not Found

## Performance Testing

### Large Dataset Performance

- [ ] Create 1000 jobs for Organization A
- [ ] Create 1000 jobs for Organization B
- [ ] Log in as Organization A
- [ ] Navigate to jobs list
- [ ] Verify page loads in < 2 seconds
- [ ] Check database query plan:
  ```sql
  EXPLAIN QUERY PLAN SELECT * FROM jobs WHERE org_id = 'org-a' ORDER BY created_at DESC LIMIT 50;
  ```
- [ ] Verify index on org_id is being used

### Concurrent Access

- [ ] Open two browser sessions
- [ ] Log in as Organization A in session 1
- [ ] Log in as Organization B in session 2
- [ ] Simultaneously create jobs in both sessions
- [ ] Verify each org only sees their own jobs
- [ ] No cross-contamination

## Compliance Checks

### Audit Trail

- [ ] Enable application logging
- [ ] Log in as Organization A
- [ ] Perform CRUD operations
- [ ] Check logs: verify org_id is logged with every database query
- [ ] Verify no queries missing `WHERE org_id = ?` clause

### Data Export

If data export feature exists:

- [ ] Log in as Organization A
- [ ] Export all jobs as CSV/JSON
- [ ] Verify export contains only Organization A's data
- [ ] Check for any Organization B records (should be zero)

### Backup and Restore

- [ ] Create database backup
- [ ] Restore to test environment
- [ ] Verify org_id values are preserved
- [ ] Verify foreign key relationships maintained
- [ ] Test login and data access for both orgs

## Sign-Off

After completing all checklist items:

- [ ] All automated integration tests passing
- [ ] All manual verification steps completed
- [ ] No cross-org data leakage observed
- [ ] Performance acceptable
- [ ] Security tests passed
- [ ] Documentation updated

**Verified by:** _______________
**Date:** _______________
**Environment:** _______________

## Notes

Use this section to document any issues found during testing:

```
Issue 1: [Description]
- Severity: [High/Medium/Low]
- Steps to reproduce:
- Expected vs Actual:
- Resolution:

Issue 2: [Description]
...
```
