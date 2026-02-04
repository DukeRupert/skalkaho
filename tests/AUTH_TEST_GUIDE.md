# Phase 1 Authentication Test Suite

This document describes the comprehensive test suite for the Phase 1 Authentication implementation.

## Test Coverage

The test suite validates all components of the authentication system:

### 1. Core Authentication (`internal/auth/`)

**`auth_test.go`** - Password hashing and token generation
- Password hashing with argon2id
- Password verification with constant-time comparison
- Session token generation (cryptographically random)
- Session token hashing (SHA-256)
- Round-trip password hash/verify workflows

**`session_test.go`** - Session management
- Session creation with user and organization context
- Session validation (active, expired, invalid)
- Session destruction (logout)
- Multiple sessions per user
- Session cleanup (expired sessions)
- User agent and IP address tracking

**`context_test.go`** - Context helpers
- Storing and retrieving sessions from context
- Extracting user ID, org ID, email, name, role from session
- Authentication status checks
- Multiple context layers

**`middleware_test.go`** - HTTP middleware
- SessionMiddleware: Loads session from cookie into context
- RequireAuth: Redirects unauthenticated users to login
- RequireRole: Role-based access control
- Middleware chaining
- Cookie handling (valid, invalid, expired, missing)

### 2. Authentication Handlers (`internal/handler/auth/`)

**`login_test.go`** - Login flow
- GET /login: Render login page
- POST /login: Authenticate user
  - Success with valid credentials
  - Failure with invalid password
  - Failure with nonexistent user
  - Failure with inactive user
  - Session creation on successful login
  - Cookie setting with HttpOnly, SameSite
  - Redirect to original URL after login
  - Case-insensitive email handling

**`register_test.go`** - Registration flow
- GET /register: Render registration page
- POST /register: Create organization, user, and session
  - Success creates all three entities
  - Duplicate email rejection
  - Duplicate subdomain rejection
  - Password confirmation validation
  - Weak password rejection
  - Email format validation
  - Subdomain format validation
  - Email normalization (lowercase)
  - Transaction rollback on partial failure
  - Password hashing (never plaintext)
  - First user is admin role

**`logout_test.go`** - Logout flow
- POST /logout: Destroy session
  - Session destruction in database
  - Cookie clearing (MaxAge=-1)
  - Redirect to login page
  - Idempotent (multiple logouts don't error)
  - Only destroys current session, not all user sessions
  - Handles missing or invalid session gracefully

### 3. Integration Tests (`tests/integration/`)

**`auth_flow_test.go`** - End-to-end authentication flows
- Complete user journey (register → login → access → logout → login again)
- Multiple users in same organization
- Multi-tenant isolation (different orgs can't access each other's data)
- Session expiration behavior
- HTTP cookie handling through middleware
- Password verification in realistic scenarios

### 4. Test Utilities (`internal/testutil/`)

**`auth.go`** - Helper functions for auth testing
- CreateTestOrgAndUser: Sets up test data
- CreateTestSession: Creates authenticated session
- ValidateTestSession: Validates session token
- AssertSessionExists/NotExists: Verifies session state
- AssertUserExists, AssertOrgExists: Verifies entity creation
- GetUserByEmail: Retrieves user for testing
- CountUserSessions: Counts active sessions

## Running the Tests

### Prerequisites

- Go 1.24+
- Docker (for testcontainers PostgreSQL)
- All dependencies installed (`go mod download`)

### Run All Auth Tests

```bash
# Run all unit tests in internal/auth/
go test ./internal/auth/... -v

# Run all handler tests
go test ./internal/handler/auth/... -v

# Run integration tests
go test ./tests/integration/... -v -tags=integration

# Run everything together
go test ./internal/auth/... ./internal/handler/auth/... ./tests/integration/... -v -tags=integration
```

### Run Specific Test Files

```bash
# Password hashing tests only
go test ./internal/auth/ -run TestHashPassword -v
go test ./internal/auth/ -run TestVerifyPassword -v

# Session management tests only
go test ./internal/auth/ -run TestCreateSession -v
go test ./internal/auth/ -run TestValidateSession -v
go test ./internal/auth/ -run TestDestroySession -v

# Middleware tests only
go test ./internal/auth/ -run TestSessionMiddleware -v
go test ./internal/auth/ -run TestRequireAuthMiddleware -v

# Login flow tests only
go test ./internal/handler/auth/ -run TestPostLogin -v

# Registration tests only
go test ./internal/handler/auth/ -run TestPostRegister -v

# Integration flow tests
go test ./tests/integration/ -run TestAuthFlow -v -tags=integration
```

### Run Tests with Coverage

```bash
# Generate coverage report
go test ./internal/auth/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# View coverage in terminal
go test ./internal/auth/... -cover

# Detailed coverage by function
go test ./internal/auth/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Run Tests in Parallel

```bash
# Run tests in parallel (faster)
go test ./internal/auth/... -v -parallel=4

# Integration tests may need sequential execution
go test ./tests/integration/... -v -tags=integration -parallel=1
```

## Test Implementation Status

These tests are written **before implementation** (TDD approach) and will **initially fail**. This is expected and intentional.

### Expected Test Failures

Until the Phase 1 implementation is complete, tests will fail with errors like:
- "undefined: HashPassword"
- "undefined: NewSessionManager"
- "undefined: SessionMiddleware"
- "undefined: RequireAuth"
- Handler types not defined

### Implementation Checklist

Use these tests to drive implementation. Tests should pass after completing each component:

- [ ] `internal/auth/auth.go` - Password hashing and token functions
- [ ] `internal/auth/session.go` - SessionManager implementation
- [ ] `internal/auth/context.go` - Context helper functions
- [ ] `internal/auth/middleware.go` - Session and auth middleware
- [ ] `internal/handler/auth/handler.go` - Handler struct
- [ ] `internal/handler/auth/login.go` - Login handlers
- [ ] `internal/handler/auth/register.go` - Registration handlers
- [ ] `internal/handler/auth/logout.go` - Logout handler
- [ ] `cmd/server/migrations/00019_sessions.sql` - Sessions table
- [ ] `sqlc/queries/sessions.sql` - Session queries (run `make sqlc`)
- [ ] `internal/config/config.go` - Session config fields

## Manual Testing Checklist

After all automated tests pass, perform these manual verifications:

### Registration Flow
1. Visit `/register`
2. Fill out form with new organization and user details
3. Submit form
4. Verify redirected to `/` (home page)
5. Verify session cookie is set in browser
6. Verify organization created in database
7. Verify user created with admin role
8. Verify session created in database

### Login Flow
1. Visit `/login`
2. Enter valid credentials
3. Verify redirected to home page
4. Verify session cookie is set
5. Try with invalid password → should show error
6. Try with nonexistent email → should show error

### Protected Routes
1. Clear cookies
2. Visit `/jobs` or other protected route
3. Verify redirected to `/login?redirect=/jobs`
4. Login successfully
5. Verify redirected back to `/jobs`

### Logout Flow
1. While logged in, click logout button
2. Verify redirected to `/login`
3. Verify session cookie is cleared
4. Try to visit protected route → should redirect to login

### Multi-User Testing
1. Register first organization and user
2. Logout
3. Register second organization and user
4. Login as second user
5. Verify can only see second org's data (not first org's)
6. Logout, login as first user
7. Verify can only see first org's data

### Session Expiration
1. Set `SESSION_DURATION=1m` in environment
2. Login
3. Wait 2 minutes
4. Try to access protected route
5. Verify redirected to login (session expired)

## Test Database

Tests use `testcontainers-go` to spin up isolated PostgreSQL containers. Each test gets a fresh database with migrations applied.

### Test Database Lifecycle

1. Test starts → PostgreSQL container launches
2. Migrations run automatically
3. Test executes with isolated database
4. Test ends → Container and database destroyed

This ensures:
- Tests don't interfere with each other
- No cleanup required between tests
- Real PostgreSQL behavior (not mocks)
- Schema matches production

## Debugging Tests

### Verbose Output

```bash
go test ./internal/auth/... -v
```

### Run Single Test

```bash
go test ./internal/auth/ -run TestHashPassword/simple_password -v
```

### Show Test Names

```bash
go test ./internal/auth/... -list .
```

### Skip Integration Tests

```bash
go test ./... -v -short  # Skips tests marked with t.Skip
```

### Debug with Delve

```bash
dlv test ./internal/auth -- -test.run TestHashPassword
```

## Continuous Integration

Add to CI pipeline:

```yaml
# .github/workflows/test.yml
- name: Run auth tests
  run: |
    go test ./internal/auth/... -v -race
    go test ./internal/handler/auth/... -v -race
    go test ./tests/integration/... -v -tags=integration
```

## Test Maintenance

### Adding New Tests

When adding features, add corresponding tests:
1. Add test case to relevant test file
2. Ensure test fails initially (red)
3. Implement feature (green)
4. Refactor if needed

### Updating Tests

When changing auth behavior:
1. Update tests to match new requirements
2. Run tests to see failures
3. Update implementation
4. Verify all tests pass

## Security Testing

These tests verify security best practices:
- ✅ Passwords hashed with argon2id (not bcrypt/sha256)
- ✅ Constant-time password comparison
- ✅ Session tokens are cryptographically random (32 bytes)
- ✅ Session tokens stored as hashes (SHA-256)
- ✅ Cookies marked HttpOnly, SameSite=Lax
- ✅ Secure cookies in production (HTTPS)
- ✅ No plaintext passwords in database
- ✅ Tenant isolation (org_id scoping)
- ✅ Inactive users cannot login
- ✅ Expired sessions rejected

## Performance Considerations

### Slow Tests

Password hashing (argon2id) is intentionally slow (security feature). If tests are too slow:
- Reduce iteration count in test environment (add test-specific config)
- Run tests in parallel
- Use `-short` flag to skip slow tests

### Test Parallelization

```go
func TestSomething(t *testing.T) {
    t.Parallel()  // Allows test to run in parallel
    // ... test code
}
```

## Troubleshooting

### "testcontainers: Docker not available"
- Ensure Docker is running
- Verify Docker socket is accessible

### "database migration failed"
- Check migrations are in `migrations/` directory
- Verify migration SQL is valid PostgreSQL

### "undefined: HashPassword"
- Implementation not complete yet (expected for TDD)
- Verify file paths match expected structure

### Tests hang indefinitely
- Check for missing `cleanup()` defer calls
- Verify testcontainer terminates properly
- Look for context.Background() instead of test context

## Success Criteria

All tests should pass before considering Phase 1 complete:

```bash
# All tests pass
go test ./internal/auth/... ./internal/handler/auth/... ./tests/integration/... -tags=integration

# No race conditions
go test ./internal/auth/... -race

# Good coverage (>80%)
go test ./internal/auth/... -cover
```

## Next Steps

After Phase 1 tests pass:
1. Run manual testing checklist
2. Test in staging environment
3. Proceed to Phase 2 (Email Verification, Password Reset, etc.)
