//go:build integration

package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthFlow_CompleteRegistrationAndLogin tests the complete registration and login flow.
func TestAuthFlow_CompleteRegistrationAndLogin(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()
	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)

	t.Run("complete_user_journey", func(t *testing.T) {
		email := "journey@example.com"
		password := "SecurePassword123!"
		subdomain := "journeytest"

		// Step 1: User doesn't exist yet
		var userCount int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&userCount)
		require.NoError(t, err)
		assert.Equal(t, int64(0), userCount, "user should not exist initially")

		// Step 2: Register new user and organization
		// (This would be done through handler, but we test the underlying logic)
		orgID := uuid.New()
		_, err = db.ExecContext(ctx, `
			INSERT INTO organizations (id, name, subdomain, created_at)
			VALUES ($1, 'Journey Org', $2, NOW())
		`, orgID, subdomain)
		require.NoError(t, err)

		passwordHash, err := auth.HashPassword(password)
		require.NoError(t, err)

		userID := uuid.New()
		_, err = db.ExecContext(ctx, `
			INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
			VALUES ($1, $2, $3, 'Journey User', $4, 'admin', 'active', NOW())
		`, userID, orgID, email, passwordHash)
		require.NoError(t, err)

		// Step 3: Create session (simulating login)
		token, session, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Test Browser",
			IPAddress: "127.0.0.1",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, email, session.Email)
		assert.Equal(t, "admin", session.Role)

		// Step 4: Validate session
		validatedSession, err := sm.ValidateSession(ctx, token)
		require.NoError(t, err)
		assert.Equal(t, session.ID, validatedSession.ID)
		assert.Equal(t, userID, validatedSession.UserID)
		assert.Equal(t, orgID, validatedSession.OrgID)

		// Step 5: Access context with session
		testCtx := auth.WithSession(context.Background(), validatedSession)
		assert.True(t, auth.IsAuthenticated(testCtx))
		assert.Equal(t, orgID, auth.OrgIDFromContext(testCtx))
		assert.Equal(t, userID, auth.UserIDFromContext(testCtx))
		assert.Equal(t, "admin", auth.UserRoleFromContext(testCtx))

		// Step 6: Logout (destroy session)
		err = sm.DestroySession(ctx, token)
		require.NoError(t, err)

		// Step 7: Verify session is destroyed
		_, err = sm.ValidateSession(ctx, token)
		assert.Error(t, err, "session should be invalid after logout")

		// Step 8: Login again (create new session)
		token2, session2, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)
		assert.NotEqual(t, token, token2, "new session should have different token")
		assert.NotEqual(t, session.ID, session2.ID, "new session should have different ID")

		// Step 9: Verify new session is valid
		validatedSession2, err := sm.ValidateSession(ctx, token2)
		require.NoError(t, err)
		assert.Equal(t, session2.ID, validatedSession2.ID)
	})
}

// TestAuthFlow_MultipleUsersSameOrg tests multiple users in the same organization.
func TestAuthFlow_MultipleUsersSameOrg(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()
	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)

	// Create organization
	orgID := uuid.New()
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES ($1, 'Multi User Org', 'multiuser', NOW())
	`, orgID)
	require.NoError(t, err)

	// Create multiple users in same org
	adminID := uuid.New()
	memberID := uuid.New()

	users := []struct {
		id    uuid.UUID
		email string
		role  string
	}{
		{adminID, "admin@multiuser.com", "admin"},
		{memberID, "member@multiuser.com", "member"},
	}

	for _, user := range users {
		passwordHash, err := auth.HashPassword("password123")
		require.NoError(t, err)

		_, err = db.ExecContext(ctx, `
			INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW())
		`, user.id, orgID, user.email, user.email, passwordHash, user.role)
		require.NoError(t, err)
	}

	t.Run("both_users_can_login_simultaneously", func(t *testing.T) {
		// Create sessions for both users
		adminToken, adminSession, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: adminID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		memberToken, memberSession, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: memberID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Verify both sessions are valid
		assert.NotEqual(t, adminToken, memberToken, "tokens should be different")

		// Validate admin session
		validatedAdmin, err := sm.ValidateSession(ctx, adminToken)
		require.NoError(t, err)
		assert.Equal(t, "admin", validatedAdmin.Role)
		assert.Equal(t, adminSession.ID, validatedAdmin.ID)

		// Validate member session
		validatedMember, err := sm.ValidateSession(ctx, memberToken)
		require.NoError(t, err)
		assert.Equal(t, "member", validatedMember.Role)
		assert.Equal(t, memberSession.ID, validatedMember.ID)

		// Both should have same org ID
		assert.Equal(t, orgID, validatedAdmin.OrgID)
		assert.Equal(t, orgID, validatedMember.OrgID)
	})

	t.Run("destroying_one_session_does_not_affect_other", func(t *testing.T) {
		// Create sessions for both users
		adminToken, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: adminID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		memberToken, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: memberID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Destroy admin session
		err = sm.DestroySession(ctx, adminToken)
		require.NoError(t, err)

		// Admin session should be invalid
		_, err = sm.ValidateSession(ctx, adminToken)
		assert.Error(t, err, "admin session should be destroyed")

		// Member session should still be valid
		memberSession, err := sm.ValidateSession(ctx, memberToken)
		require.NoError(t, err)
		assert.NotNil(t, memberSession, "member session should still be valid")
	})
}

// TestAuthFlow_MultiTenantIsolation tests tenant isolation in authentication.
func TestAuthFlow_MultiTenantIsolation(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()
	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)

	// Create two separate organizations
	orgA := uuid.New()
	orgB := uuid.New()

	for i, org := range []uuid.UUID{orgA, orgB} {
		_, err := db.ExecContext(ctx, `
			INSERT INTO organizations (id, name, subdomain, created_at)
			VALUES ($1, $2, $3, NOW())
		`, org, "Org "+string(rune('A'+i)), "org-"+string(rune('a'+i)))
		require.NoError(t, err)
	}

	// Create users in different organizations
	userAID := uuid.New()
	userBID := uuid.New()

	passwordHash, err := auth.HashPassword("password123")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
		VALUES
			($1, $2, 'usera@orga.com', 'User A', $3, 'admin', 'active', NOW()),
			($4, $5, 'userb@orgb.com', 'User B', $6, 'admin', 'active', NOW())
	`, userAID, orgA, passwordHash, userBID, orgB, passwordHash)
	require.NoError(t, err)

	t.Run("sessions_are_scoped_to_org", func(t *testing.T) {
		// Create sessions for both users
		tokenA, sessionA, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userAID,
			OrgID:  orgA,
		})
		require.NoError(t, err)

		tokenB, sessionB, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userBID,
			OrgID:  orgB,
		})
		require.NoError(t, err)

		// Validate sessions return correct org IDs
		assert.Equal(t, orgA, sessionA.OrgID)
		assert.Equal(t, orgB, sessionB.OrgID)

		// Verify org IDs are different
		assert.NotEqual(t, sessionA.OrgID, sessionB.OrgID)

		// Verify tokens are different
		assert.NotEqual(t, tokenA, tokenB)
	})

	t.Run("context_helpers_extract_correct_org", func(t *testing.T) {
		// Create session for user A
		tokenA, sessionA, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userAID,
			OrgID:  orgA,
		})
		require.NoError(t, err)

		// Validate and add to context
		validatedSession, err := sm.ValidateSession(ctx, tokenA)
		require.NoError(t, err)

		testCtx := auth.WithSession(context.Background(), validatedSession)

		// Extract org ID from context
		extractedOrgID := auth.OrgIDFromContext(testCtx)
		assert.Equal(t, orgA, extractedOrgID, "should extract correct org ID")
		assert.Equal(t, sessionA.OrgID, extractedOrgID)
	})

	t.Run("cannot_create_session_with_mismatched_user_and_org", func(t *testing.T) {
		// Attempt to create session with user from org A but org B
		_, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userAID,
			OrgID:  orgB, // Mismatched!
		})

		// This should fail due to foreign key constraint or validation
		assert.Error(t, err, "should not allow mismatched user and org")
	})
}

// TestAuthFlow_SessionExpiration tests session expiration behavior.
func TestAuthFlow_SessionExpiration(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create organization and user
	orgID, userID := testutil.CreateTestOrgAndUser(t, db, "expire@example.com", "Expire User", "member", "password123")

	t.Run("session_expires_after_duration", func(t *testing.T) {
		// Create session manager with very short duration
		sm := auth.NewSessionManager(db, "test-secret", 50*time.Millisecond, "session", false)

		// Create session
		token, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Session should be valid immediately
		session, err := sm.ValidateSession(ctx, token)
		require.NoError(t, err)
		assert.NotNil(t, session)

		// Wait for session to expire
		time.Sleep(100 * time.Millisecond)

		// Session should now be invalid
		_, err = sm.ValidateSession(ctx, token)
		assert.Error(t, err, "session should be expired")
	})

	t.Run("long_duration_session_remains_valid", func(t *testing.T) {
		// Create session manager with long duration
		sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)

		// Create session
		token, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Wait a moment
		time.Sleep(50 * time.Millisecond)

		// Session should still be valid
		session, err := sm.ValidateSession(ctx, token)
		require.NoError(t, err)
		assert.NotNil(t, session)
	})
}

// TestAuthFlow_SessionCookieHandling tests cookie handling in HTTP flow.
func TestAuthFlow_SessionCookieHandling(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()
	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "test_session", false)

	// Create organization and user
	orgID, userID := testutil.CreateTestOrgAndUser(t, db, "cookie@example.com", "Cookie User", "member", "password123")

	t.Run("middleware_loads_session_from_cookie", func(t *testing.T) {
		// Create session
		token, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler that checks for session
		var capturedSession *auth.Session
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSession = auth.SessionFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session middleware
		wrappedHandler := sm.SessionMiddleware(handler)

		// Create request with session cookie
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify session was loaded
		assert.NotNil(t, capturedSession, "session should be loaded from cookie")
		assert.Equal(t, userID, capturedSession.UserID)
		assert.Equal(t, orgID, capturedSession.OrgID)
	})

	t.Run("require_auth_middleware_protects_routes", func(t *testing.T) {
		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session and auth middleware
		wrappedHandler := sm.SessionMiddleware(auth.RequireAuth(handler))

		// Request without session
		req := httptest.NewRequest("GET", "/protected", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		// Should redirect to login
		assert.False(t, handlerCalled, "handler should not be called without auth")
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/login", rec.Header().Get("Location"))
	})

	t.Run("authenticated_request_accesses_protected_route", func(t *testing.T) {
		// Create session
		token, _, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session and auth middleware
		wrappedHandler := sm.SessionMiddleware(auth.RequireAuth(handler))

		// Request with valid session
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Should access protected route
		assert.True(t, handlerCalled, "handler should be called with valid auth")
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// TestAuthFlow_PasswordVerification tests password verification in authentication.
func TestAuthFlow_PasswordVerification(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create organization and user with known password
	password := "TestPassword123!"
	email := "pwtest@example.com"

	orgID := uuid.New()
	userID := uuid.New()

	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES ($1, 'PW Test Org', 'pwtest', NOW())
	`, orgID)
	require.NoError(t, err)

	passwordHash, err := auth.HashPassword(password)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
		VALUES ($1, $2, $3, 'PW User', $4, 'member', 'active', NOW())
	`, userID, orgID, email, passwordHash)
	require.NoError(t, err)

	t.Run("correct_password_verifies", func(t *testing.T) {
		// Retrieve user
		var storedHash string
		err := db.QueryRowContext(ctx, `
			SELECT password_hash FROM users WHERE email = $1
		`, email).Scan(&storedHash)
		require.NoError(t, err)

		// Verify password
		match, err := auth.VerifyPassword(password, storedHash)
		require.NoError(t, err)
		assert.True(t, match, "correct password should verify")
	})

	t.Run("incorrect_password_fails", func(t *testing.T) {
		// Retrieve user
		var storedHash string
		err := db.QueryRowContext(ctx, `
			SELECT password_hash FROM users WHERE email = $1
		`, email).Scan(&storedHash)
		require.NoError(t, err)

		// Verify wrong password
		match, err := auth.VerifyPassword("WrongPassword123!", storedHash)
		require.NoError(t, err)
		assert.False(t, match, "incorrect password should not verify")
	})
}
