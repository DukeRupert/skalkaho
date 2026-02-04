package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dukerupert/skalkaho/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateSession verifies session creation.
func TestCreateSession(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "test@example.com", "Test User", "admin")

	// Create SessionManager (assumes queries are generated via sqlc)
	// This will fail until implementation exists
	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "session_cookie", false)

	t.Run("creates_session_successfully", func(t *testing.T) {
		token, session, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Mozilla/5.0",
			IPAddress: "192.168.1.1",
		})

		require.NoError(t, err, "CreateSession should not error")
		assert.NotEmpty(t, token, "token should not be empty")
		assert.NotNil(t, session, "session should not be nil")
		assert.Equal(t, userID, session.UserID, "session should have correct user_id")
		assert.Equal(t, orgID, session.OrgID, "session should have correct org_id")
		assert.Equal(t, "test@example.com", session.Email, "session should include user email")
		assert.Equal(t, "Test User", session.Name, "session should include user name")
		assert.Equal(t, "admin", session.Role, "session should include user role")

		// Verify session exists in database
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "session should exist in database")
	})

	t.Run("stores_user_agent_and_ip", func(t *testing.T) {
		token, session, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Custom User Agent 1.0",
			IPAddress: "10.0.0.1",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify stored values
		var userAgent, ipAddress sql.NullString
		err = db.QueryRowContext(ctx,
			`SELECT user_agent, ip_address FROM sessions WHERE id = $1`,
			session.ID,
		).Scan(&userAgent, &ipAddress)
		require.NoError(t, err)
		assert.True(t, userAgent.Valid)
		assert.Equal(t, "Custom User Agent 1.0", userAgent.String)
		assert.True(t, ipAddress.Valid)
		assert.Equal(t, "10.0.0.1", ipAddress.String)
	})

	t.Run("sets_expiration_correctly", func(t *testing.T) {
		beforeCreate := time.Now()
		token, session, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		afterCreate := time.Now()

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Expiration should be approximately 24 hours from now
		expectedExpiry := beforeCreate.Add(24 * time.Hour)
		assert.True(t, session.ExpiresAt.After(expectedExpiry.Add(-1*time.Minute)))
		assert.True(t, session.ExpiresAt.Before(afterCreate.Add(24*time.Hour+time.Minute)))
	})

	t.Run("token_hash_is_stored_not_plaintext", func(t *testing.T) {
		token, session, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify that the stored token_hash is not the plaintext token
		var storedHash string
		err = db.QueryRowContext(ctx, `SELECT token_hash FROM sessions WHERE id = $1`, session.ID).Scan(&storedHash)
		require.NoError(t, err)
		assert.NotEqual(t, token, storedHash, "plaintext token should not be stored")
		assert.NotEmpty(t, storedHash, "token hash should be stored")
	})

	t.Run("multiple_sessions_for_same_user", func(t *testing.T) {
		// Create first session
		token1, session1, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Browser 1",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, token1)

		// Create second session
		token2, session2, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID:    userID,
			OrgID:     orgID,
			UserAgent: "Browser 2",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, token2)

		// Tokens should be different
		assert.NotEqual(t, token1, token2, "different sessions should have different tokens")
		assert.NotEqual(t, session1.ID, session2.ID, "different sessions should have different IDs")

		// Both should exist in database
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2), "multiple sessions should exist for user")
	})

	t.Run("fails_for_nonexistent_user", func(t *testing.T) {
		nonExistentUserID := uuid.New()

		_, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: nonExistentUserID,
			OrgID:  orgID,
		})

		assert.Error(t, err, "creating session for nonexistent user should fail")
	})
}

// TestValidateSession verifies session validation logic.
func TestValidateSession(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "validate@example.com", "Validate User", "member")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "session_cookie", false)

	t.Run("validates_active_session", func(t *testing.T) {
		// Create session
		token, createdSession, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Validate session
		session, err := sm.ValidateSession(ctx, token)
		require.NoError(t, err, "ValidateSession should not error for valid session")
		assert.NotNil(t, session, "session should not be nil")
		assert.Equal(t, createdSession.ID, session.ID, "session ID should match")
		assert.Equal(t, userID, session.UserID, "user ID should match")
		assert.Equal(t, orgID, session.OrgID, "org ID should match")
		assert.Equal(t, "validate@example.com", session.Email)
		assert.Equal(t, "member", session.Role)
	})

	t.Run("rejects_expired_session", func(t *testing.T) {
		// Create SessionManager with very short expiration
		shortSM := NewSessionManager(db, "test-secret", 1*time.Millisecond, "session_cookie", false)

		token, _, err := shortSM.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Wait for session to expire
		time.Sleep(10 * time.Millisecond)

		// Validation should fail
		session, err := shortSM.ValidateSession(ctx, token)
		assert.Error(t, err, "ValidateSession should error for expired session")
		assert.Nil(t, session, "session should be nil for expired token")
	})

	t.Run("rejects_invalid_token", func(t *testing.T) {
		invalidToken := "invalid-token-12345"

		session, err := sm.ValidateSession(ctx, invalidToken)
		assert.Error(t, err, "ValidateSession should error for invalid token")
		assert.Nil(t, session, "session should be nil for invalid token")
	})

	t.Run("rejects_empty_token", func(t *testing.T) {
		session, err := sm.ValidateSession(ctx, "")
		assert.Error(t, err, "ValidateSession should error for empty token")
		assert.Nil(t, session, "session should be nil for empty token")
	})

	t.Run("updates_last_activity", func(t *testing.T) {
		// Create session
		token, createdSession, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		initialActivity := createdSession.LastActivityAt

		// Wait a moment
		time.Sleep(100 * time.Millisecond)

		// Validate session (should update last_activity_at)
		session, err := sm.ValidateSession(ctx, token)
		require.NoError(t, err)

		// Last activity should be updated
		assert.True(t, session.LastActivityAt.After(initialActivity),
			"last_activity_at should be updated on validation")
	})

	t.Run("rejects_session_for_inactive_user", func(t *testing.T) {
		// Create user and session
		inactiveUserID := uuid.New()
		createTestOrgAndUser(t, db, orgID, inactiveUserID, "inactive@example.com", "Inactive User", "member")

		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: inactiveUserID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Mark user as suspended (inactive is not a valid status - valid statuses are: active, suspended, invited)
		_, err = db.ExecContext(ctx, `UPDATE users SET status = 'suspended' WHERE id = $1`, inactiveUserID)
		require.NoError(t, err)

		// Validation should fail
		session, err := sm.ValidateSession(ctx, token)
		assert.Error(t, err, "ValidateSession should fail for inactive user")
		assert.Nil(t, session, "session should be nil for inactive user")
	})
}

// TestDestroySession verifies session destruction.
func TestDestroySession(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "destroy@example.com", "Destroy User", "admin")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "session_cookie", false)

	t.Run("destroys_session_successfully", func(t *testing.T) {
		// Create session
		token, session, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Verify session exists
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Destroy session
		err = sm.DestroySession(ctx, token)
		require.NoError(t, err, "DestroySession should not error")

		// Verify session is deleted
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, session.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "session should be deleted from database")
	})

	t.Run("validates_after_destroy_fails", func(t *testing.T) {
		// Create session
		token, _, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Destroy session
		err = sm.DestroySession(ctx, token)
		require.NoError(t, err)

		// Attempt to validate destroyed session
		session, err := sm.ValidateSession(ctx, token)
		assert.Error(t, err, "ValidateSession should fail after DestroySession")
		assert.Nil(t, session, "session should be nil after destroy")
	})

	t.Run("destroy_nonexistent_session_does_not_error", func(t *testing.T) {
		invalidToken := "nonexistent-token"

		// Destroying a nonexistent session should not error (idempotent)
		err := sm.DestroySession(ctx, invalidToken)
		assert.NoError(t, err, "DestroySession should be idempotent")
	})

	t.Run("destroy_empty_token_does_not_error", func(t *testing.T) {
		err := sm.DestroySession(ctx, "")
		assert.NoError(t, err, "DestroySession should handle empty token gracefully")
	})
}

// TestDestroyAllUserSessions verifies destroying all sessions for a user.
func TestDestroyAllUserSessions(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and users
	orgID := uuid.New()
	user1ID := uuid.New()
	user2ID := uuid.New()
	createTestOrgAndUser(t, db, orgID, user1ID, "user1@example.com", "User 1", "admin")
	createTestOrgAndUser(t, db, orgID, user2ID, "user2@example.com", "User 2", "member")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "session_cookie", false)

	t.Run("destroys_all_sessions_for_user", func(t *testing.T) {
		// Create multiple sessions for user1
		_, _, err := sm.CreateSession(ctx, CreateSessionParams{UserID: user1ID, OrgID: orgID, UserAgent: "Browser 1"})
		require.NoError(t, err)
		_, _, err = sm.CreateSession(ctx, CreateSessionParams{UserID: user1ID, OrgID: orgID, UserAgent: "Browser 2"})
		require.NoError(t, err)
		_, _, err = sm.CreateSession(ctx, CreateSessionParams{UserID: user1ID, OrgID: orgID, UserAgent: "Browser 3"})
		require.NoError(t, err)

		// Create session for user2 (should not be affected)
		_, _, err = sm.CreateSession(ctx, CreateSessionParams{UserID: user2ID, OrgID: orgID})
		require.NoError(t, err)

		// Verify user1 has 3 sessions
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, user1ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Destroy all sessions for user1
		err = sm.DestroyAllUserSessions(ctx, user1ID)
		require.NoError(t, err)

		// Verify user1 has no sessions
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, user1ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "all sessions for user should be destroyed")

		// Verify user2 still has their session
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, user2ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "other users' sessions should not be affected")
	})
}

// TestCleanupExpiredSessions verifies cleanup of expired sessions.
func TestCleanupExpiredSessions(t *testing.T) {
	db, cleanup := testutil.TestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test organization and user
	orgID := uuid.New()
	userID := uuid.New()
	createTestOrgAndUser(t, db, orgID, userID, "cleanup@example.com", "Cleanup User", "admin")

	sm := NewSessionManager(db, "test-secret", 24*time.Hour, "session_cookie", false)

	t.Run("removes_only_expired_sessions", func(t *testing.T) {
		// Create an active session (24h expiry)
		_, activeSession, err := sm.CreateSession(ctx, CreateSessionParams{
			UserID: userID,
			OrgID:  orgID,
		})
		require.NoError(t, err)

		// Create an expired session by manually inserting with past expiration
		expiredSessionID := uuid.New()
		tokenHash := HashSessionToken("expired-token")
		_, err = db.ExecContext(ctx, `
			INSERT INTO sessions (id, user_id, org_id, token_hash, expires_at, last_activity_at, created_at)
			VALUES ($1, $2, $3, $4, NOW() - INTERVAL '1 hour', NOW(), NOW())
		`, expiredSessionID, userID, orgID, tokenHash)
		require.NoError(t, err)

		// Verify both sessions exist
		var count int64
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(2), count, "should have 2 sessions before cleanup")

		// Cleanup expired sessions
		err = sm.CleanupExpiredSessions(ctx)
		require.NoError(t, err)

		// Verify expired session is deleted
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, expiredSessionID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "expired session should be deleted")

		// Verify active session still exists
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions WHERE id = $1`, activeSession.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "active session should not be deleted")
	})
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
	require.NoError(t, err, "failed to create test organization")

	// Create user with a hashed password
	passwordHash, err := HashPassword("password123")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW())
		ON CONFLICT (id) DO NOTHING
	`, userID, orgID, email, name, passwordHash, role)
	require.NoError(t, err, "failed to create test user")
}
