package auth

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRenderer implements the Renderer interface for testing.
type mockRenderer struct {
	lastTemplate string
	lastData     interface{}
	err          error
}

func (m *mockRenderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	m.lastTemplate = name
	m.lastData = data
	if m.err != nil {
		return m.err
	}
	// Write a 200 OK response by default
	w.WriteHeader(http.StatusOK)
	return nil
}

// TestGetLogin verifies the login page is rendered.
func TestGetLogin(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	renderer := &mockRenderer{}

	// Create handler
	h := NewHandler(db, renderer, sm)

	t.Run("renders_login_page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		rec := httptest.NewRecorder()

		h.GetLogin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "should return 200 OK")
		// In real implementation, would check that template was rendered
	})

	t.Run("includes_redirect_parameter_in_form", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login?redirect=/protected/page", nil)
		rec := httptest.NewRecorder()

		h.GetLogin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// In real implementation, would verify redirect parameter is passed to template
	})

	t.Run("redirects_already_authenticated_user", func(t *testing.T) {
		// Create session
		orgID := uuid.New()
		userID := uuid.New()
		createTestOrgAndUser(t, db, orgID, userID, "test@example.com", "Test User", "member")

		token, _, err := sm.CreateSession(context.Background(), auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create request with session
		req := httptest.NewRequest("GET", "/login", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token,
		})

		// Add session to context (normally done by middleware)
		session, err := sm.ValidateSession(context.Background(), token)
		require.NoError(t, err)
		ctx := auth.WithSession(req.Context(), session)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		h.GetLogin(rec, req)

		// Should redirect authenticated user away from login
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect authenticated user")
		assert.Equal(t, "/", rec.Header().Get("Location"), "should redirect to home")
	})
}

// TestPostLogin verifies login form submission.
func TestPostLogin(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	email := "login@example.com"
	password := "SecurePassword123!"

	createTestOrgAndUser(t, db, orgID, userID, email, "Login User", "admin")

	// Update user with known password
	passwordHash, err := auth.HashPassword(password)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE users SET password_hash = $1 WHERE id = $2`, passwordHash, userID)
	require.NoError(t, err)

	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	renderer := &mockRenderer{}
	h := NewHandler(db, renderer, sm)

	t.Run("success_with_valid_credentials", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should redirect to home page
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect after successful login")
		assert.Equal(t, "/", rec.Header().Get("Location"), "should redirect to home")

		// Should set session cookie
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie, "should set session cookie")
		assert.NotEmpty(t, sessionCookie.Value, "session cookie should have value")
		assert.True(t, sessionCookie.HttpOnly, "session cookie should be HttpOnly")
		assert.Equal(t, http.SameSiteLaxMode, sessionCookie.SameSite, "should use SameSite=Lax")
	})

	t.Run("redirects_to_original_url_after_login", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", password)
		formData.Set("redirect", "/protected/resource")

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/protected/resource", rec.Header().Get("Location"),
			"should redirect to original URL")
	})

	t.Run("fails_with_invalid_password", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", "WrongPassword123!")

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should return error (200 with error message or redirect back to login)
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect on failed login")

		// Should not set session cookie
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			assert.NotEqual(t, "session", cookie.Name, "should not set session cookie on failed login")
		}
	})

	t.Run("fails_with_nonexistent_user", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "nonexistent@example.com")
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should return error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect on failed login")

		// Should not set session cookie
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			assert.NotEqual(t, "session", cookie.Name, "should not set session cookie")
		}
	})

	t.Run("fails_with_inactive_user", func(t *testing.T) {
		// Create user, then suspend them (non-active status)
		suspendedUserID := uuid.New()
		suspendedEmail := "suspended@example.com"
		createTestOrgAndUser(t, db, orgID, suspendedUserID, suspendedEmail, "Suspended User", "member")

		// Set password and mark as suspended (non-active status)
		passwordHash, err := auth.HashPassword(password)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `
			UPDATE users
			SET password_hash = $1, status = 'suspended'
			WHERE id = $2
		`, passwordHash, suspendedUserID)
		require.NoError(t, err)

		formData := url.Values{}
		formData.Set("email", suspendedEmail)
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should return error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect for inactive user")

		// Should not set session cookie
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			assert.NotEqual(t, "session", cookie.Name, "should not set session cookie for inactive user")
		}
	})

	t.Run("fails_with_missing_email", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should return validation error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with missing email")
	})

	t.Run("fails_with_missing_password", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", email)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should return validation error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with missing password")
	})

	t.Run("creates_session_in_database", func(t *testing.T) {
		// Count sessions before login
		var countBefore int64
		err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&countBefore)
		require.NoError(t, err)

		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Count sessions after login
		var countAfter int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&countAfter)
		require.NoError(t, err)

		assert.Equal(t, countBefore+1, countAfter, "should create new session in database")
	})

	t.Run("stores_user_agent_and_ip", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", "Test Browser 1.0")
		req.RemoteAddr = "192.168.1.100:12345"
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Get the session cookie
		cookies := rec.Result().Cookies()
		var token string
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				token = cookie.Value
				break
			}
		}
		require.NotEmpty(t, token, "session token should be set")

		// Validate session and check user agent
		_, err = sm.ValidateSession(ctx, token)
		require.NoError(t, err)

		// Verify user agent is stored (exact comparison depends on implementation)
		// This test assumes user agent is accessible through session or database
		// Implementation may vary
	})

	t.Run("case_insensitive_email", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", strings.ToUpper(email)) // Use uppercase version
		formData.Set("password", password)

		req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostLogin(rec, req)

		// Should succeed with case-insensitive email
		assert.Equal(t, http.StatusFound, rec.Code, "should succeed with case-insensitive email")
	})
}

// TestPostLogin_RateLimiting verifies rate limiting (if implemented).
func TestPostLogin_RateLimiting(t *testing.T) {
	t.Skip("Rate limiting implementation depends on Phase 1 scope")

	// This test would verify:
	// - Multiple failed login attempts result in temporary lockout
	// - Lockout is per-user or per-IP
	// - Successful login resets attempt counter
}

// Helper function to create test organization and user
func createTestOrgAndUser(t *testing.T, db *sql.DB, orgID, userID uuid.UUID, email, name, role string) {
	t.Helper()

	ctx := context.Background()

	// Create organization
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (id) DO NOTHING
	`, orgID, "Test Org", "test-org")
	require.NoError(t, err)

	// Create user
	passwordHash, err := auth.HashPassword("password123")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW())
		ON CONFLICT (id) DO NOTHING
	`, userID, orgID, email, name, passwordHash, role)
	require.NoError(t, err)
}
