# Integration Tests for Multi-Tenant Data Isolation

This directory contains comprehensive integration tests for verifying tenant data isolation in the Skalkaho construction quoting application.

## Overview

These tests verify that multi-tenancy is correctly implemented across all data entities:
- Jobs (quotes)
- Categories
- Line Items
- Settings
- Clients
- Estimates
- Signature Requests

## Test Files

### `tenant_isolation_test.go`
Tests basic tenant isolation across all entities:
- Data visibility is scoped to organization
- Cross-org read attempts return not found errors
- List operations only return org's own data
- Aggregate queries (counts, totals) are org-scoped
- Hierarchy inheritance (category → job, line item → category)

### `tenant_crud_test.go`
Tests CRUD operations with tenant boundaries:
- Create operations associate data with correct org
- Read operations filter by org_id
- Update operations cannot modify other org's data
- Delete operations cannot delete other org's data
- Cascading deletes respect org boundaries
- Cross-org data manipulation attempts are prevented

### `TENANT_ISOLATION_CHECKLIST.md`
Manual verification checklist for:
- UI-level testing
- Security edge cases
- Performance with large datasets
- Compliance and audit trails

## Current State

**IMPORTANT:** These tests are written to fail initially and will pass once multi-tenancy is fully implemented.

Currently, the tests include TODO comments indicating where organization context should be passed once the auth/session middleware is implemented.

## Running the Tests

### Prerequisites

Install required dependencies:
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/google/uuid
go get github.com/pressly/goose/v3
```

### Run all integration tests
```bash
cd /workspaces/skalkaho
go test ./tests/integration/... -v
```

### Run specific test file
```bash
go test ./tests/integration/tenant_isolation_test.go -v
go test ./tests/integration/tenant_crud_test.go -v
```

### Run specific test function
```bash
go test ./tests/integration -run TestTenantIsolation_Jobs -v
go test ./tests/integration -run TestTenantCRUD_JobLifecycle -v
```

### Run with coverage
```bash
go test ./tests/integration/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Expected Failures

Until multi-tenancy implementation is complete, expect these tests to fail with messages like:
- "no such column: org_id" - org_id columns not yet added to schema
- Cross-org access succeeds when it should fail
- Counts include all orgs instead of just one

These failures are **expected** and serve as validation tests for the implementation.

## Implementation Checklist

To make these tests pass, the following must be implemented:

### 1. Database Schema Changes
- [ ] Add `org_id TEXT NOT NULL` column to: jobs, categories, line_items, settings, clients
- [ ] Add indexes: `CREATE INDEX idx_jobs_org_id ON jobs(org_id)`
- [ ] Add foreign keys or triggers to enforce org_id consistency
- [ ] Add migration script in `migrations/` directory

### 2. Organization Context Middleware
- [ ] Implement middleware to extract org_id from session/JWT
- [ ] Add org_id to request context
- [ ] Handle multi-org users (if applicable)

### 3. Repository Layer Updates
- [ ] Update all SQL queries to include `WHERE org_id = ?`
- [ ] Modify CreateJob to auto-set org_id from context
- [ ] Modify CreateCategory to inherit org_id from job
- [ ] Modify CreateLineItem to inherit org_id from category
- [ ] Update all List/Get/Update/Delete queries to filter by org_id

Example query changes:
```sql
-- Before
SELECT * FROM jobs WHERE id = ?

-- After
SELECT * FROM jobs WHERE id = ? AND org_id = ?
```

### 4. Business Logic Updates
- [ ] Validate org ownership before any UPDATE/DELETE
- [ ] Prevent cross-org foreign key references
- [ ] Ensure cascading deletes respect org boundaries

### 5. Test Helpers
- [ ] Update `testutil.NewFixtures` to support org_id
- [ ] Add helper to create org context for tests
- [ ] Update fixture methods to set org_id when creating test data

## Test Utilities

### `internal/testutil/database.go`
Provides helpers for test database setup:
- `TestDB(t)` - Creates isolated test database with migrations
- `ExecSQL(t, db, query, args)` - Execute SQL statements
- `MustCount(t, db, query, args)` - Get count with failure on error
- `AssertRowExists/AssertRowNotExists` - Verify data existence

### `internal/testutil/fixtures.go`
Provides factories for creating test data:
- `CreateJob(params)` - Create test job
- `CreateCategory(params)` - Create test category
- `CreateLineItem(params)` - Create test line item
- `CreateSettings(params)` - Create test settings
- `CreateClient(params)` - Create test client
- `CreateEstimate(params)` - Create test estimate

Example usage:
```go
db, cleanup := testutil.TestDB(t)
defer cleanup()

fixtures := testutil.NewFixtures(t, db)
jobID := fixtures.CreateJob(testutil.JobParams{
    OrgID: "org-a",
    Name:  "Test Job",
    SurchargePercent: 15.0,
})
```

## Writing New Tests

When adding new tenant-scoped features, follow this pattern:

```go
func TestTenantIsolation_NewFeature(t *testing.T) {
    db, cleanup := testutil.TestDB(t)
    defer cleanup()

    fixtures := testutil.NewFixtures(t, db)
    queries := repository.New(db)
    ctx := context.Background()

    // Create data for org-a
    orgAData := fixtures.CreateSomething(testutil.SomeParams{
        OrgID: "org-a",
        // ...
    })

    // Create data for org-b
    orgBData := fixtures.CreateSomething(testutil.SomeParams{
        OrgID: "org-b",
        // ...
    })

    t.Run("org-a cannot access org-b data", func(t *testing.T) {
        // TODO: With org-a context
        result, err := queries.GetSomething(ctx, orgBData)

        // After multi-tenancy:
        // require.Error(t, err)
        // assert.True(t, errors.Is(err, sql.ErrNoRows))
    })
}
```

## Troubleshooting

### Tests fail with "migrations directory not found"
The `findMigrationsDir` function walks up from the current directory to find `migrations/`. Ensure:
- You're running tests from within the project directory
- The `migrations/` directory exists at the project root

### Tests fail with "table not found"
Migrations may not have run. Check:
- Migration files exist in `migrations/`
- Goose can parse migration files
- Test database is being created correctly

### Race conditions in tests
If tests fail intermittently:
- Ensure proper cleanup with `defer cleanup()`
- Use `t.Parallel()` carefully with shared state
- Check for proper transaction isolation

## Best Practices

1. **Use table-driven tests** for multiple scenarios
2. **Always use test helpers** from `testutil` package
3. **Mark TODOs** where org context is needed
4. **Test both success and failure** cases
5. **Verify side effects** (cascade deletes, etc.)
6. **Use descriptive test names** that explain what's being tested
7. **Include comments** explaining complex test logic
8. **Clean up after tests** with defer statements

## Related Documentation

- `/development/MVP_GUIDE.md` - Complete MVP specification
- `/development/GO_STYLE_GUIDE.md` - Code style guidelines
- `/CLAUDE.md` - Project overview and architecture
- `TENANT_ISOLATION_CHECKLIST.md` - Manual verification steps

## Questions or Issues

If you encounter issues with these tests or need to add new test coverage, consult:
1. This README
2. Existing test files for patterns
3. Test utilities in `internal/testutil/`
4. Project documentation in `/development/`
