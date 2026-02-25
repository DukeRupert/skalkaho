package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSessionMiddleware verifies that the session middleware loads sessions from cookies.
func TestSessionMiddleware(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "middleware@example.com", "Middleware User", "admin")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "test_session", false)

	t.Run("loads_session_from_valid_cookie", func(t *testing.T) {
		// Create a session
		token, createdSession, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler that checks for session in context
		var capturedSession *Session
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSession = SessionFromContext(r.Context())
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

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify session was loaded into context
		assert.NotNil(t, capturedSession, "session should be loaded into context")
		assert.Equal(t, createdSession.ID, capturedSession.ID)
		assert.Equal(t, userID, capturedSession.UserID)
		assert.Equal(t, orgID, capturedSession.OrgID)
		assert.Equal(t, "middleware@example.com", capturedSession.Email)
	})

	t.Run("continues_without_session_when_no_cookie", func(t *testing.T) {
		// Create test handler that checks for session in context
		var capturedSession *Session
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSession = SessionFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session middleware
		wrappedHandler := sm.SessionMiddleware(handler)

		// Create request without session cookie
		req := httptest.NewRequest("GET", "/test", nil)

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify no session in context, but request continues
		assert.Nil(t, capturedSession, "session should be nil when no cookie")
		assert.Equal(t, http.StatusOK, rec.Code, "request should succeed without session")
	})

	t.Run("continues_without_session_when_invalid_cookie", func(t *testing.T) {
		// Create test handler
		var capturedSession *Session
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSession = SessionFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session middleware
		wrappedHandler := sm.SessionMiddleware(handler)

		// Create request with invalid session cookie
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: "invalid-token-12345",
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify no session in context, but request continues
		assert.Nil(t, capturedSession, "session should be nil for invalid token")
		assert.Equal(t, http.StatusOK, rec.Code, "request should succeed with invalid token")
	})

	t.Run("does_not_load_expired_session", func(t *testing.T) {
		// Create SessionManager with very short expiration
		shortSM := NewSessionManager(db, "test-secret", 1*time.Millisecond, "test_session", false)

		// Create session
		token, _, err := shortSM.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Wait for session to expire
		time.Sleep(10 * time.Millisecond)

		// Create test handler
		var capturedSession *Session
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSession = SessionFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session middleware
		wrappedHandler := shortSM.SessionMiddleware(handler)

		// Create request with expired session cookie
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify session is not loaded
		assert.Nil(t, capturedSession, "expired session should not be loaded")
		assert.Equal(t, http.StatusOK, rec.Code, "request should continue without session")
	})

	t.Run("ignores_wrong_cookie_name", func(t *testing.T) {
		// Create session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler
		var capturedSession *Session
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedSession = SessionFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with session middleware
		wrappedHandler := sm.SessionMiddleware(handler)

		// Create request with wrong cookie name
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "wrong_session_name",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify session is not loaded
		assert.Nil(t, capturedSession, "session should not be loaded from wrong cookie name")
	})
}

// TestRequireAuthMiddleware verifies authentication requirement middleware.
func TestRequireAuthMiddleware(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "auth@example.com", "Auth User", "member")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "test_session", false)

	t.Run("allows_authenticated_request", func(t *testing.T) {
		// Create session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		})

		// Wrap with both middleware (session, then require auth)
		wrappedHandler := sm.SessionMiddleware(RequireAuth(handler))

		// Create request with valid session cookie
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was called and request succeeded
		assert.True(t, handlerCalled, "handler should be called for authenticated request")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})

	t.Run("redirects_unauthenticated_request", func(t *testing.T) {
		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with both middleware (session, then require auth)
		wrappedHandler := sm.SessionMiddleware(RequireAuth(handler))

		// Create request without session cookie
		req := httptest.NewRequest("GET", "/protected", nil)

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was NOT called and request was redirected
		assert.False(t, handlerCalled, "handler should not be called for unauthenticated request")
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect unauthenticated user")

		// Verify redirect location
		location := rec.Header().Get("Location")
		assert.Equal(t, "/login", location, "should redirect to login page")
	})

	t.Run("preserves_original_url_in_redirect", func(t *testing.T) {
		// Create test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with middleware
		wrappedHandler := sm.SessionMiddleware(RequireAuth(handler))

		// Create request to a specific protected path
		req := httptest.NewRequest("GET", "/protected/resource?id=123", nil)

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify redirect includes original URL
		location := rec.Header().Get("Location")
		assert.Contains(t, location, "redirect=", "should include redirect parameter")
		assert.Contains(t, location, "%2Fprotected%2Fresource", "should include encoded original path")
	})

	t.Run("rejects_expired_session", func(t *testing.T) {
		// Create SessionManager with very short expiration
		shortSM := NewSessionManager(db, "test-secret", 1*time.Millisecond, "test_session", false)

		// Create session
		token, _, err := shortSM.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Wait for session to expire
		time.Sleep(10 * time.Millisecond)

		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with middleware
		wrappedHandler := shortSM.SessionMiddleware(RequireAuth(handler))

		// Create request with expired session cookie
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was NOT called and request was redirected
		assert.False(t, handlerCalled, "handler should not be called for expired session")
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect for expired session")
	})

	t.Run("rejects_invalid_token", func(t *testing.T) {
		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with middleware
		wrappedHandler := sm.SessionMiddleware(RequireAuth(handler))

		// Create request with invalid token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: "invalid-token-abc123",
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was NOT called and request was redirected
		assert.False(t, handlerCalled, "handler should not be called for invalid token")
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect for invalid token")
	})
}

// TestRequireAuthMiddleware_WithRoles tests role-based authorization.
func TestRequireAuthMiddleware_WithRoles(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and users with different roles
	orgID := uuid.New()
	adminUserID := uuid.New()
	memberUserID := uuid.New()

	createTestOrgAndUser(t, db, orgID, adminUserID, "admin@example.com", "Admin User", "admin")
	createTestOrgAndUser(t, db, orgID, memberUserID, "member@example.com", "Member User", "member")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "test_session", false)

	t.Run("allows_admin_access_to_admin_only_resource", func(t *testing.T) {
		// Create admin session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: adminUserID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with middleware that requires admin role
		wrappedHandler := sm.SessionMiddleware(RequireRole("admin")(handler))

		// Create request with admin session
		req := httptest.NewRequest("GET", "/admin", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was called
		assert.True(t, handlerCalled, "admin should access admin-only resource")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("denies_member_access_to_admin_only_resource", func(t *testing.T) {
		// Create member session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: memberUserID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with middleware that requires admin role
		wrappedHandler := sm.SessionMiddleware(RequireRole("admin")(handler))

		// Create request with member session
		req := httptest.NewRequest("GET", "/admin", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was NOT called
		assert.False(t, handlerCalled, "member should not access admin-only resource")
		assert.Equal(t, http.StatusForbidden, rec.Code, "should return forbidden for insufficient role")
	})

	t.Run("allows_any_authenticated_user_when_no_role_required", func(t *testing.T) {
		// Create member session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: memberUserID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap with RequireAuth only (no role requirement)
		wrappedHandler := sm.SessionMiddleware(RequireAuth(handler))

		// Create request with member session
		req := httptest.NewRequest("GET", "/dashboard", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify handler was called
		assert.True(t, handlerCalled, "any authenticated user should access resource")
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// TestMiddlewareChaining verifies middleware work correctly when chained.
func TestMiddlewareChaining(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "chain@example.com", "Chain User", "admin")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "test_session", false)

	t.Run("session_and_auth_middleware_chain", func(t *testing.T) {
		// Create session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create test handler that accesses session
		var capturedOrgID uuid.UUID
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedOrgID = OrgIDFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Chain middleware: session -> auth -> handler
		wrappedHandler := sm.SessionMiddleware(RequireAuth(handler))

		// Create authenticated request
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "test_session",
			Value: token,
		})

		// Execute request
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Verify middleware chain worked
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, orgID, capturedOrgID, "org ID should be accessible through context")
	})
}
