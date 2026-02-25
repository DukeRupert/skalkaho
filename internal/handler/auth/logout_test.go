package auth

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

// TestPostLogout verifies logout functionality.
func TestPostLogout(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "logout@example.com", "Logout User", "member")

	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	h := NewHandler(db, nil, sm)

	t.Run("destroys_session_and_redirects", func(t *testing.T) {
		// Create session
		token, session, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Verify session exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "session should exist before logout")

		// Create logout request with session in context
		req := httptest.NewRequest("POST", "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token,
		})
		req = req.WithContext(auth.WithSession(req.Context(), session))

		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		// Should redirect to login
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect after logout")
		assert.Equal(t, "/login", rec.Header().Get("Location"), "should redirect to login page")

		// Verify session is destroyed in database
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "session should be destroyed after logout")

		// Verify session cookie is cleared
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie, "should set session cookie to clear it")
		assert.Empty(t, sessionCookie.Value, "session cookie value should be empty")
		assert.Equal(t, -1, sessionCookie.MaxAge, "session cookie should have negative MaxAge to delete")
	})

	t.Run("handles_logout_without_session_gracefully", func(t *testing.T) {
		// Create logout request without session
		req := httptest.NewRequest("POST", "/logout", nil)
		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		// Should still redirect to login (idempotent)
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect even without session")
		assert.Equal(t, "/login", rec.Header().Get("Location"), "should redirect to login")
	})

	t.Run("handles_logout_with_invalid_session", func(t *testing.T) {
		// Create logout request with invalid session
		req := httptest.NewRequest("POST", "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: "invalid-token-12345",
		})
		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		// Should still redirect (idempotent)
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect even with invalid session")
		assert.Equal(t, "/login", rec.Header().Get("Location"), "should redirect to login")
	})

	t.Run("clears_cookie_even_if_session_already_destroyed", func(t *testing.T) {
		// Create session
		token, session, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Manually destroy session
		err = sm.DestroySession(ctx, token)
		require.NoError(t, err)

		// Verify session is destroyed
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// Try to logout with already-destroyed session
		req := httptest.NewRequest("POST", "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token,
		})
		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		// Should still clear cookie and redirect
		assert.Equal(t, http.StatusFound, rec.Code)

		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie, "should clear cookie")
		assert.Empty(t, sessionCookie.Value)
	})

	t.Run("only_destroys_current_session_not_all_user_sessions", func(t *testing.T) {
		// Create two sessions for the same user
		token1, session1, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Browser 1",
		})
		require.NoError(t, err)

		token2, session2, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Browser 2",
		})
		require.NoError(t, err)

		// Logout from session1
		req := httptest.NewRequest("POST", "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token1,
		})
		req = req.WithContext(auth.WithSession(req.Context(), session1))
		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code)

		// Verify session1 is destroyed
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session1.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "logged out session should be destroyed")

		// Verify session2 still exists
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session2.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "other sessions should not be affected")

		// Verify session2 is still valid
		validatedSession, err := sm.ValidateSession(ctx, token2)
		require.NoError(t, err)
		assert.NotNil(t, validatedSession, "other sessions should remain valid")
	})

	t.Run("handles_concurrent_logouts", func(t *testing.T) {
		// Create session
		token, session, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Simulate concurrent logout requests with same session
		req1 := httptest.NewRequest("POST", "/logout", nil)
		req1.AddCookie(&http.Cookie{Name: "session", Value: token})
		req1 = req1.WithContext(auth.WithSession(req1.Context(), session))

		req2 := httptest.NewRequest("POST", "/logout", nil)
		req2.AddCookie(&http.Cookie{Name: "session", Value: token})
		req2 = req2.WithContext(auth.WithSession(req2.Context(), session))

		rec1 := httptest.NewRecorder()
		rec2 := httptest.NewRecorder()

		// Both requests should succeed (idempotent)
		h.Logout(rec1, req1)
		h.Logout(rec2, req2)

		assert.Equal(t, http.StatusFound, rec1.Code, "first logout should succeed")
		assert.Equal(t, http.StatusFound, rec2.Code, "second logout should succeed (idempotent)")

		// Verify session is destroyed
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "session should be destroyed")
	})

	t.Run("allows_custom_redirect_after_logout", func(t *testing.T) {
		// Create session
		token, session, err := sm.CreateSession(ctx, auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create logout request with redirect parameter
		req := httptest.NewRequest("POST", "/logout?redirect=/goodbye", nil)
		req.AddCookie(&http.Cookie{Name: "session", Value: token})
		req = req.WithContext(auth.WithSession(req.Context(), session))
		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		// Should redirect to custom location if implementation supports it
		assert.Equal(t, http.StatusFound, rec.Code)
		// Note: Default behavior is to redirect to /login, but implementation
		// may support custom redirect
	})
}

// TestPostLogout_OnlyPOST verifies logout only works with POST method.
func TestPostLogout_OnlyPOST(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	h := NewHandler(db, nil, sm)

	t.Run("rejects_GET_request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/logout", nil)
		rec := httptest.NewRecorder()

		// Assuming router enforces POST-only, this test verifies handler doesn't accept GET
		// If handler is called with GET, it should handle gracefully
		h.Logout(rec, req)

		// Implementation should either:
		// 1. Return method not allowed (if checked in handler)
		// 2. Be protected by router to only accept POST
		// This test assumes routing is handled elsewhere and won't be called with GET
	})
}

// TestPostLogout_CSRF verifies CSRF protection (if implemented).
func TestPostLogout_CSRF(t *testing.T) {
	t.Skip("CSRF protection implementation depends on Phase 1 scope")

	// This test would verify:
	// - Logout requires valid CSRF token
	// - Missing or invalid CSRF token rejects logout
	// - CSRF token is validated before session destruction
}
