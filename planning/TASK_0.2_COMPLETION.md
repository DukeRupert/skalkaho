# Task 0.2: Organizations and Users Tables - COMPLETED

## Summary
Successfully created multi-tenancy foundation with organizations and users tables for PostgreSQL.

## What Was Implemented

### 1. Database Migrations

**Migration 00016_organizations.sql**
- UUID primary key with `gen_random_uuid()` default
- Organization metadata: name, subdomain (unique)
- Billing fields: stripe_customer_id (unique), plan, status
- Subscription lifecycle: trial_ends_at, subscription_ends_at
- Timestamps: created_at, updated_at (both TIMESTAMPTZ)
- Indexes: subdomain, stripe_customer_id (partial), status
- Check constraints: plan (free/pro/enterprise), status (active/suspended/cancelled)

**Migration 00017_users.sql**
- UUID primary key with `gen_random_uuid()` default
- Foreign key to organizations with CASCADE delete
- Authentication: email, password_hash
- Authorization: role (owner/admin/member)
- Profile: name
- Account fields: email_verified, status
- Password reset: reset_token, reset_token_expires_at
- Email verification: verification_token
- Audit: last_login_at
- Timestamps: created_at, updated_at
- Unique constraint: (org_id, email) - ensures email uniqueness per organization
- Indexes: org_id, email, reset_token (partial), verification_token (partial)

### 2. sqlc Queries

**organizations.sql** - Full CRUD operations:
- CreateOrganization
- GetOrganization
- GetOrganizationBySubdomain
- GetOrganizationByStripeCustomerID
- UpdateOrganization (with COALESCE for partial updates)
- DeleteOrganization
- ListOrganizations
- ListActiveOrganizations

**users.sql** - Complete user management:
- CreateUser
- GetUser
- GetUserByEmail
- GetUserByResetToken
- GetUserByVerificationToken
- ListUsersByOrg
- UpdateUser (with COALESCE for partial updates)
- UpdateUserLastLogin
- SetUserResetToken
- ClearUserResetToken
- SetUserVerificationToken
- VerifyUserEmail
- DeleteUser
- CountUsersByOrg

### 3. PostgreSQL Migration Fixes

Fixed SQLite-specific syntax across all migrations:
- Replaced `AUTOINCREMENT` with `SERIAL`
- Replaced `datetime('now')` with `NOW()`
- Replaced SQLite UUID generation with `gen_random_uuid()`
- Converted query placeholders from `?` to `$1, $2, $3...`
- Fixed named parameters from `@name` to `sqlc.arg('name')`

### 4. sqlc Configuration Updates

Updated `/workspaces/skalkaho/sqlc.yaml` with type overrides:
- `pg_catalog.float4` → `float64` (for REAL columns)
- `pg_catalog.int4` → `int64` (for INTEGER columns)
- `serial` → `int64` (for SERIAL columns)

This ensures compatibility with existing Go code that uses int64/float64.

### 5. Code Fixes

Fixed type conversion issues in handlers:
- `/workspaces/skalkaho/internal/handler/keyboard/clients.go` - Cast offset/limit to int32
- `/workspaces/skalkaho/internal/handler/keyboard/jobs.go` - Cast offset/limit to int32
- `/workspaces/skalkaho/internal/handler/keyboard/price_import.go` - Convert NullInt64 to NullInt32

### 6. Migration Management

- Synced `/workspaces/skalkaho/migrations/` with `/workspaces/skalkaho/cmd/server/migrations/`
- Both directories now contain identical PostgreSQL-compatible migrations

## Testing Results

All migrations ran successfully:
```
✓ 00001-00015: Existing tables migrated
✓ 00016: organizations table created
✓ 00017: users table created
```

Repository functions tested and verified:
- Organization CRUD operations
- User CRUD operations
- Foreign key constraints working
- Unique constraints working
- Indexes created properly

## Files Modified

**Created:**
- `/workspaces/skalkaho/migrations/00016_organizations.sql`
- `/workspaces/skalkaho/migrations/00017_users.sql`
- `/workspaces/skalkaho/sqlc/queries/organizations.sql`
- `/workspaces/skalkaho/sqlc/queries/users.sql`
- `/workspaces/skalkaho/docker-compose.dev.yml`
- `/workspaces/skalkaho/.env`

**Modified:**
- `/workspaces/skalkaho/sqlc.yaml` - Added type overrides
- `/workspaces/skalkaho/migrations/00005_item_templates.sql` - AUTOINCREMENT → SERIAL
- `/workspaces/skalkaho/migrations/00009_migrate_customer_names.sql` - UUID generation
- `/workspaces/skalkaho/migrations/00010_price_imports.sql` - AUTOINCREMENT → SERIAL
- All migrations - datetime('now') → NOW()
- All query files - ? → $N placeholders, @param → sqlc.arg('param')
- Handler files - Type conversions for int32/int64

**Generated:**
- `/workspaces/skalkaho/internal/repository/organizations.sql.go`
- `/workspaces/skalkaho/internal/repository/users.sql.go`
- Updated `/workspaces/skalkaho/internal/repository/models.go` with Organization and User structs

## Next Steps

Ready to proceed to **Task 0.3: Add org_id to Existing Tables**

This will involve:
1. Adding `org_id UUID REFERENCES organizations(id)` to existing tables
2. Creating a default organization for development
3. Migrating existing data to the default organization
4. Adding NOT NULL constraint after data migration
5. Updating all indexes to include org_id
