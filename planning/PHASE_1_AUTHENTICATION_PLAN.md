# Authentication Implementation Plan

**Phase:** 1 - Alpha Launch
**Task:** User Authentication
**Status:** Planning Complete

---

## Overview

Add authentication to Skalkaho with session management, login/logout, registration, and protected routes. The users and organizations tables already exist with all required fields.

---

## Implementation Steps

### Step 1: Database - Sessions Table

**File:** `cmd/server/migrations/00019_sessions.sql`

```sql
-- +goose Up
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- +goose Down
DROP TABLE IF EXISTS sessions;
```

### Step 2: Session SQL Queries

**File:** `sqlc/queries/sessions.sql`

```sql
-- name: CreateSession :one
INSERT INTO sessions (user_id, org_id, token_hash, user_agent, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT s.*, u.email, u.name, u.role, u.status as user_status
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token_hash = $1 AND s.expires_at > NOW() LIMIT 1;

-- name: UpdateSessionActivity :exec
UPDATE sessions SET last_activity_at = NOW() WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions WHERE token_hash = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < NOW();
```

Then run: `make sqlc`

### Step 3: Configuration Updates

**File:** `internal/config/config.go`

Add fields:
```go
SessionSecret     string        // SESSION_SECRET env var (required in prod)
SessionDuration   time.Duration // SESSION_DURATION env var (default: 720h = 30 days)
SessionCookieName string        // default: "skalkaho_session"
SecureCookies     bool          // true if ENVIRONMENT=production
```

### Step 4: Auth Package

**Directory:** `internal/auth/`

| File | Purpose |
|------|---------|
| `auth.go` | Password hashing (argon2id), token generation |
| `session.go` | SessionManager: create, validate, destroy sessions |
| `context.go` | Context helpers: WithSession, SessionFromContext, OrgIDFromContext |
| `middleware.go` | Session middleware (loads session) + RequireAuth middleware |

**Key Functions:**
- `HashPassword(password string) (string, error)` - argon2id hashing
- `VerifyPassword(password, hash string) (bool, error)` - constant-time comparison
- `GenerateSessionToken() (string, error)` - 32-byte random token
- `HashSessionToken(token string) string` - SHA-256 for storage

### Step 5: Auth Handlers

**Directory:** `internal/handler/auth/`

| File | Handlers |
|------|----------|
| `handler.go` | Handler struct with queries, renderer, sessionManager |
| `login.go` | `GetLogin`, `PostLogin` |
| `logout.go` | `PostLogout` |
| `register.go` | `GetRegister`, `PostRegister` |
| `password_reset.go` | `GetForgotPassword`, `PostForgotPassword`, `GetResetPassword`, `PostResetPassword` |
| `verify_email.go` | `VerifyEmail` |

### Step 6: Templates

**Directory:** `internal/templates/keyboard/pages/`

| File | Purpose |
|------|---------|
| `login.html` | Login form with subdomain, email, password |
| `register.html` | Registration form with org creation |
| `forgot_password.html` | Password reset request |
| `reset_password.html` | New password form |

Templates follow existing patterns: Tailwind CSS, copper/slate colors, HTMX forms.

### Step 7: Route Updates

**File:** `internal/router/router.go`

**Public Routes (no auth):**
- `GET/POST /login`
- `GET/POST /register`
- `POST /logout`
- `GET/POST /forgot-password`
- `GET/POST /reset-password/{token}`
- `GET /verify-email/{token}`
- `GET/POST /sign/{token}` (existing signature pages)
- `GET /health`

**Protected Routes (require auth):**
- All existing routes (`/`, `/jobs/*`, `/clients/*`, `/estimates/*`, `/settings/*`)

### Step 8: Update tenant.go

**File:** `internal/handler/keyboard/tenant.go`

Update `GetOrgID()` and `MustGetOrgID()` to extract from session context:
```go
func GetOrgID(ctx context.Context) uuid.NullUUID {
    orgID := auth.OrgIDFromContext(ctx)
    if orgID == uuid.Nil {
        return uuid.NullUUID{}
    }
    return uuid.NullUUID{UUID: orgID, Valid: true}
}
```

### Step 9: main.go Integration

**File:** `cmd/server/main.go`

1. Validate SESSION_SECRET in production
2. Create SessionManager
3. Create auth handler
4. Add session middleware to chain (after Logger)
5. Update router.Register() call with auth handler

---

## Files to Create

| Path | Description |
|------|-------------|
| `cmd/server/migrations/00019_sessions.sql` | Sessions table migration |
| `sqlc/queries/sessions.sql` | Session CRUD queries |
| `internal/auth/auth.go` | Password hashing, token generation |
| `internal/auth/session.go` | SessionManager |
| `internal/auth/context.go` | Context helpers |
| `internal/auth/middleware.go` | Session + RequireAuth middleware |
| `internal/handler/auth/handler.go` | Auth handler struct |
| `internal/handler/auth/login.go` | Login handlers |
| `internal/handler/auth/logout.go` | Logout handler |
| `internal/handler/auth/register.go` | Registration handlers |
| `internal/handler/auth/password_reset.go` | Password reset handlers |
| `internal/handler/auth/verify_email.go` | Email verification |
| `internal/templates/keyboard/pages/login.html` | Login page |
| `internal/templates/keyboard/pages/register.html` | Registration page |
| `internal/templates/keyboard/pages/forgot_password.html` | Forgot password page |
| `internal/templates/keyboard/pages/reset_password.html` | Reset password page |

## Files to Modify

| Path | Changes |
|------|---------|
| `internal/config/config.go` | Add session config fields |
| `internal/handler/keyboard/tenant.go` | Extract org_id from session |
| `internal/router/router.go` | Add auth routes, wrap protected routes |
| `cmd/server/main.go` | Wire up session manager, auth handler, middleware |
| `internal/templates/keyboard/layouts/base.html` | Add logout button to header |

---

## Dependencies

Add to `go.mod`:
```
golang.org/x/crypto  # Already available as indirect dep
```

No new external dependencies required. Uses:
- `golang.org/x/crypto/argon2` for password hashing
- `crypto/sha256` for token hashing
- `crypto/rand` for token generation

---

## Security Considerations

1. **Password Storage:** argon2id with OWASP parameters (64MB memory, 1 iteration, 4 threads)
2. **Session Tokens:** 256-bit random, stored as SHA-256 hash
3. **Cookies:** HttpOnly, Secure (prod), SameSite=Lax
4. **Constant-time comparison** for password and token verification

---

## Verification

After implementation, verify:

1. **Registration Flow:**
   - Visit `/register`, create account with new org
   - Verify redirected to `/` (home)
   - Verify org and user created in database

2. **Login Flow:**
   - Visit `/login`, enter credentials
   - Verify session created and redirected to home
   - Verify session cookie set

3. **Protected Routes:**
   - Clear cookies, visit `/jobs`
   - Verify redirected to `/login`
   - After login, verify can access `/jobs`

4. **Logout:**
   - Click logout
   - Verify session destroyed and redirected to `/login`

5. **Multi-tenancy:**
   - Create two orgs with different users
   - Verify each user only sees their org's data

6. **Run Tests:**
   ```bash
   go test ./internal/auth/...
   go test ./internal/handler/auth/...
   ```
