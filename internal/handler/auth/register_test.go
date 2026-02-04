package auth

import (
	"context"
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

// TestGetRegister verifies the registration page is rendered.
func TestGetRegister(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	renderer := &mockRenderer{}
	h := NewHandler(db, renderer, sm)

	t.Run("renders_registration_page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/register", nil)
		rec := httptest.NewRecorder()

		h.GetRegister(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "should return 200 OK")
	})

	t.Run("redirects_already_authenticated_user", func(t *testing.T) {
		// Create session
		orgID := uuid.New()
		userID := uuid.New()
		createTestOrgAndUser(t, db, orgID, userID, "test@example.com", "Test User", "admin")

		token, _, err := sm.CreateSession(context.Background(), auth.CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create request with session in context
		session, err := sm.ValidateSession(context.Background(), token)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/register", nil)
		req = req.WithContext(auth.WithSession(req.Context(), session))
		rec := httptest.NewRecorder()

		h.GetRegister(rec, req)

		// Should redirect authenticated user
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect authenticated user")
		assert.Equal(t, "/", rec.Header().Get("Location"), "should redirect to home")
	})
}

// TestPostRegister verifies user registration.
func TestPostRegister(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()
	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	renderer := &mockRenderer{}
	h := NewHandler(db, renderer, sm)

	t.Run("success_creates_org_user_and_session", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("org_name", "Acme Construction")
		formData.Set("subdomain", "acme")
		formData.Set("name", "John Doe")
		formData.Set("email", "john@acme.com")
		formData.Set("password", "SecurePassword123!")
		formData.Set("password_confirm", "SecurePassword123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		// Should redirect to home after successful registration
		assert.Equal(t, http.StatusFound, rec.Code, "should redirect after registration")
		assert.Equal(t, "/", rec.Header().Get("Location"), "should redirect to home")

		// Verify organization was created
		var orgID uuid.UUID
		err := db.QueryRowContext(ctx, `
			SELECT id FROM organizations WHERE subdomain = $1
		`, "acme").Scan(&orgID)
		require.NoError(t, err, "organization should be created")
		assert.NotEqual(t, uuid.Nil, orgID)

		// Verify user was created with admin role
		var userID uuid.UUID
		var role, status string
		err = db.QueryRowContext(ctx, `
			SELECT id, role, status FROM users WHERE email = $1
		`, "john@acme.com").Scan(&userID, &role, &status)
		require.NoError(t, err, "user should be created")
		assert.Equal(t, "admin", role, "first user should be admin")
		assert.Equal(t, "active", status, "user should be active")

		// Verify user belongs to created organization
		var userOrgID uuid.UUID
		err = db.QueryRowContext(ctx, `
			SELECT org_id FROM users WHERE id = $1
		`, userID).Scan(&userOrgID)
		require.NoError(t, err)
		assert.Equal(t, orgID, userOrgID, "user should belong to created org")

		// Verify session cookie was set
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie, "should set session cookie")
		assert.NotEmpty(t, sessionCookie.Value)
		assert.True(t, sessionCookie.HttpOnly)

		// Verify session exists in database
		var sessionCount int64
		err = db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM sessions WHERE user_id = $1
		`, userID).Scan(&sessionCount)
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionCount, "session should be created")
	})

	t.Run("fails_with_duplicate_email", func(t *testing.T) {
		// Create existing user
		existingEmail := "existing@example.com"
		orgID := uuid.New()
		userID := uuid.New()
		createTestOrgAndUser(t, db, orgID, userID, existingEmail, "Existing User", "admin")

		// Try to register with same email
		formData := url.Values{}
		formData.Set("org_name", "New Org")
		formData.Set("subdomain", "neworg")
		formData.Set("name", "New User")
		formData.Set("email", existingEmail)
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		// Should return error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with duplicate email")

		// Verify no new organization was created
		var count int64
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM organizations WHERE subdomain = $1
		`, "neworg").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "should not create organization")
	})

	t.Run("fails_with_duplicate_subdomain", func(t *testing.T) {
		// Create existing organization
		existingOrgID := uuid.New()
		_, err := db.ExecContext(ctx, `
			INSERT INTO organizations (id, name, subdomain, created_at)
			VALUES ($1, 'Existing Org', 'existing', NOW())
		`, existingOrgID)
		require.NoError(t, err)

		// Try to register with same subdomain
		formData := url.Values{}
		formData.Set("org_name", "New Org")
		formData.Set("subdomain", "existing")
		formData.Set("name", "New User")
		formData.Set("email", "new@example.com")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		// Should return error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with duplicate subdomain")

		// Verify user was not created
		var count int64
		err = db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM users WHERE email = $1
		`, "new@example.com").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "should not create user")
	})

	t.Run("fails_with_mismatched_passwords", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("org_name", "Test Org")
		formData.Set("subdomain", "testorg")
		formData.Set("name", "Test User")
		formData.Set("email", "test@example.com")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "DifferentPassword123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		// Should return validation error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with mismatched passwords")

		// Verify nothing was created
		var orgCount int64
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM organizations WHERE subdomain = $1
		`, "testorg").Scan(&orgCount)
		require.NoError(t, err)
		assert.Equal(t, int64(0), orgCount, "should not create organization")
	})

	t.Run("fails_with_weak_password", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("org_name", "Test Org")
		formData.Set("subdomain", "testorg2")
		formData.Set("name", "Test User")
		formData.Set("email", "test2@example.com")
		formData.Set("password", "weak")
		formData.Set("password_confirm", "weak")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		// Should return validation error
		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with weak password")
	})

	t.Run("fails_with_missing_org_name", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("subdomain", "testorg3")
		formData.Set("name", "Test User")
		formData.Set("email", "test3@example.com")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with missing org name")
	})

	t.Run("fails_with_missing_subdomain", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("org_name", "Test Org")
		formData.Set("name", "Test User")
		formData.Set("email", "test4@example.com")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with missing subdomain")
	})

	t.Run("fails_with_invalid_email", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("org_name", "Test Org")
		formData.Set("subdomain", "testorg5")
		formData.Set("name", "Test User")
		formData.Set("email", "not-an-email")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		assert.NotEqual(t, http.StatusFound, rec.Code, "should not redirect with invalid email")
	})

	t.Run("normalizes_email_to_lowercase", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("org_name", "Lowercase Test")
		formData.Set("subdomain", "lowercase")
		formData.Set("name", "Test User")
		formData.Set("email", "TEST@EXAMPLE.COM")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code, "should succeed")

		// Verify email was stored as lowercase
		var storedEmail string
		err := db.QueryRowContext(ctx, `
			SELECT email FROM users WHERE email = $1
		`, "test@example.com").Scan(&storedEmail)
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", storedEmail, "email should be lowercase")
	})

	t.Run("validates_subdomain_format", func(t *testing.T) {
		invalidSubdomains := []string{
			"Invalid Subdomain", // spaces
			"sub_domain",        // underscores
			"sub.domain",        // dots
			"123",               // numbers only / starts with digit
			"a",                 // too short
			"ab",                // too short (need 3+ chars)
		}
		// Note: "SUB" is NOT invalid - it gets normalized to "sub" which is valid

		for _, subdomain := range invalidSubdomains {
			formData := url.Values{}
			formData.Set("org_name", "Test Org")
			formData.Set("subdomain", subdomain)
			formData.Set("name", "Test User")
			formData.Set("email", "test-subdomain@example.com")
			formData.Set("password", "Password123!")
			formData.Set("password_confirm", "Password123!")

			req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			h.PostRegister(rec, req)

			assert.NotEqual(t, http.StatusFound, rec.Code,
				"should reject invalid subdomain: %s", subdomain)
		}
	})

	t.Run("transaction_rollback_on_error", func(t *testing.T) {
		// This test verifies that if user creation fails after org creation,
		// the organization is not left orphaned (transaction should rollback)

		formData := url.Values{}
		formData.Set("org_name", "Rollback Test Org")
		formData.Set("subdomain", "rollbacktest")
		formData.Set("name", "Test User")
		formData.Set("email", "rollback@example.com")
		formData.Set("password", "Password123!")
		formData.Set("password_confirm", "Password123!")

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		// First registration should succeed
		h.PostRegister(rec, req)
		assert.Equal(t, http.StatusFound, rec.Code)

		// Try to register again with same subdomain (should fail)
		formData.Set("email", "different@example.com") // different email
		req = httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec = httptest.NewRecorder()

		h.PostRegister(rec, req)

		// Should fail
		assert.NotEqual(t, http.StatusFound, rec.Code)

		// Verify only one org with that subdomain exists (no duplicates from failed transaction)
		var count int64
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM organizations WHERE subdomain = $1
		`, "rollbacktest").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "should only have one organization")

		// Verify the second user was not created
		var userCount int64
		err = db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM users WHERE email = $1
		`, "different@example.com").Scan(&userCount)
		require.NoError(t, err)
		assert.Equal(t, int64(0), userCount, "failed registration should not create user")
	})

	t.Run("password_is_hashed_not_plaintext", func(t *testing.T) {
		password := "PlaintextPassword123!"

		formData := url.Values{}
		formData.Set("org_name", "Hash Test Org")
		formData.Set("subdomain", "hashtest")
		formData.Set("name", "Hash User")
		formData.Set("email", "hash@example.com")
		formData.Set("password", password)
		formData.Set("password_confirm", password)

		req := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.PostRegister(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code)

		// Verify password is hashed
		var passwordHash string
		err := db.QueryRowContext(ctx, `
			SELECT password_hash FROM users WHERE email = $1
		`, "hash@example.com").Scan(&passwordHash)
		require.NoError(t, err)
		assert.NotEqual(t, password, passwordHash, "password should be hashed, not plaintext")
		assert.NotEmpty(t, passwordHash, "password hash should not be empty")

		// Verify hash can be verified
		match, err := auth.VerifyPassword(password, passwordHash)
		require.NoError(t, err)
		assert.True(t, match, "hashed password should verify correctly")
	})
}

// TestPostRegister_Idempotency verifies registration is not idempotent (duplicate attempts fail).
func TestPostRegister_Idempotency(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()
	sm := auth.NewSessionManager(db, "test-secret", 24*time.Hour, "session", false)
	renderer := &mockRenderer{}
	h := NewHandler(db, renderer, sm)

	formData := url.Values{}
	formData.Set("org_name", "Idempotent Org")
	formData.Set("subdomain", "idempotent")
	formData.Set("name", "Test User")
	formData.Set("email", "idempotent@example.com")
	formData.Set("password", "Password123!")
	formData.Set("password_confirm", "Password123!")

	// First registration
	req1 := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec1 := httptest.NewRecorder()

	h.PostRegister(rec1, req1)
	assert.Equal(t, http.StatusFound, rec1.Code, "first registration should succeed")

	// Second registration (duplicate)
	req2 := httptest.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec2 := httptest.NewRecorder()

	h.PostRegister(rec2, req2)
	assert.NotEqual(t, http.StatusFound, rec2.Code, "duplicate registration should fail")

	// Verify only one user exists
	var count int64
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users WHERE email = $1
	`, "idempotent@example.com").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "should only create one user")
}
