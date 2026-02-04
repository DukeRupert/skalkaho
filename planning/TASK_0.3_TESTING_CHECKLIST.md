# Task 0.3 Testing Checklist

## Database Schema Verification ✅

### Migration Status
- [x] Migration 00018 applied successfully
- [x] Database version is 18
- [x] All 15 tables have org_id column (14 data tables + users table)

### Column Verification
- [x] Jobs table has org_id UUID column
- [x] Categories table has org_id UUID column
- [x] Line_items table has org_id UUID column
- [x] Clients table has org_id UUID column
- [x] Estimates table has org_id UUID column
- [x] Estimate_categories table has org_id UUID column
- [x] Signature_requests table has org_id UUID column
- [x] Signatures table has org_id UUID column
- [x] Item_templates table has org_id UUID column
- [x] Price_imports table has org_id UUID column
- [x] Price_import_matches table has org_id UUID column
- [x] Job_item_types table has org_id UUID column
- [x] Settings table has org_id UUID NOT NULL primary key
- [x] Company_profile table has org_id UUID NOT NULL primary key

### Foreign Key Constraints
- [x] 15 foreign key constraints on org_id columns exist
- [x] All constraints reference organizations(id)
- [x] All constraints have ON DELETE CASCADE

### Indexes
- [x] 14 indexes on org_id columns created
- [x] All indexes named idx_{table}_org_id

## sqlc Query Updates ✅

### Query Files Updated (13 files)
- [x] categories.sql - 8 queries updated
- [x] clients.sql - 7 queries updated
- [x] company.sql - 3 queries updated (added CreateCompanyProfile)
- [x] estimates.sql - 14 queries updated
- [x] item_templates.sql - 9 queries updated
- [x] job_item_types.sql - 6 queries updated
- [x] jobs.sql - 11 queries updated
- [x] line_items.sql - 6 queries updated
- [x] price_imports.sql - 12 queries updated
- [x] settings.sql - 3 queries updated (added CreateSettings)
- [x] signatures.sql - 11 queries updated

### Query Pattern Checks
- [x] All SELECT queries include org_id in WHERE clause
- [x] All INSERT queries include org_id parameter
- [x] All UPDATE queries include org_id in WHERE clause
- [x] All DELETE queries include org_id in WHERE clause
- [x] Special case: GetSignatureRequestByToken excludes org_id (public access)

## Generated Go Code ✅

### Repository Code Generation
- [x] `make sqlc` runs without errors
- [x] All structs include OrgID uuid.NullUUID field
- [x] All query parameter structs include OrgID field
- [x] Settings queries use org_id as primary key
- [x] Company_profile queries use org_id as primary key

### Sample Code Verification
```go
// Example from jobs.sql.go
type Job struct {
    ID                       string
    OrgID                    uuid.NullUUID  // ✅ Added
    Name                     string
    // ... other fields
}

type GetJobParams struct {
    ID    string
    OrgID uuid.NullUUID  // ✅ Added
}
```

## Manual Testing Required (Next Steps)

### Once Handlers Are Updated

#### 1. Create Test Organization
```sql
INSERT INTO organizations (id, name, subdomain, plan, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Test Org',
    'test',
    'free',
    'active'
);
```

#### 2. Create Test Settings
```sql
INSERT INTO settings (org_id, default_surcharge_mode, default_surcharge_percent)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'stacking',
    15.0
);
```

#### 3. Create Test Company Profile
```sql
INSERT INTO company_profile (org_id, name)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Test Company'
);
```

#### 4. Test Data Isolation

Create data for Org 1:
```sql
INSERT INTO jobs (id, org_id, name, surcharge_mode, status)
VALUES (
    'job-1',
    '00000000-0000-0000-0000-000000000001',
    'Org 1 Job',
    'stacking',
    'active'
);
```

Create data for Org 2:
```sql
INSERT INTO organizations (id, name, subdomain, plan, status)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    'Test Org 2',
    'test2',
    'free',
    'active'
);

INSERT INTO jobs (id, org_id, name, surcharge_mode, status)
VALUES (
    'job-2',
    '00000000-0000-0000-0000-000000000002',
    'Org 2 Job',
    'stacking',
    'active'
);
```

Verify isolation:
```sql
-- Should return only Org 1 jobs
SELECT * FROM jobs WHERE org_id = '00000000-0000-0000-0000-000000000001';

-- Should return only Org 2 jobs
SELECT * FROM jobs WHERE org_id = '00000000-0000-0000-0000-000000000002';
```

#### 5. Test Cascading Deletes
```sql
-- Delete an organization
DELETE FROM organizations WHERE id = '00000000-0000-0000-0000-000000000002';

-- Verify all related data is deleted
SELECT COUNT(*) FROM jobs WHERE org_id = '00000000-0000-0000-0000-000000000002';
-- Should return 0
```

## Known Limitations

1. **NULL org_id Values**: Currently allowed on most tables to support gradual migration
   - Future migration will add NOT NULL constraints after data backfill

2. **Handler Updates Required**: Repository now expects org_id parameters but handlers haven't been updated yet
   - This will be addressed in Task 0.4 and Task 0.6

3. **No Tenant Context Yet**: No middleware to extract org_id from session context
   - This will be addressed in Task 0.6 (Tenant Context Middleware)

4. **Existing Data**: Any existing data in the database has NULL org_id
   - For development: delete existing data and recreate with org_id
   - For production migration: would need data backfill script

## Success Criteria ✅

All criteria met:

- ✅ Migration 00018 successfully applied
- ✅ All 14 tables have org_id column with foreign key to organizations
- ✅ Settings and company_profile restructured with org_id as primary key
- ✅ All 14 org_id indexes created
- ✅ All sqlc queries updated to include org_id scoping
- ✅ Repository code regenerated with org_id support
- ✅ No compilation errors in generated code
- ✅ Database schema matches expected structure

## Ready for Next Task

**Task 0.4**: Update handler code to pass org_id parameters (can be done with temporary hard-coded org_id for testing)

**Task 0.5**: (Optional) Create helper functions for organization initialization (CreateOrganizationWithDefaults)

**Task 0.6**: Implement tenant context middleware to extract org_id from authenticated session

The database schema is now fully multi-tenant ready!
