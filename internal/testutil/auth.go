package testutil

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// CreateTestOrgAndUser creates a test organization and user for auth testing.
// Returns the organization ID and user ID.
// passwordHash should be pre-computed using auth.HashPassword() in the calling test.
func CreateTestOrgAndUser(t *testing.T, db *sql.DB, email, name, role, passwordHash string) (uuid.UUID, uuid.UUID) {
	t.Helper()

	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	// Create organization
	_, err := db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, subdomain, created_at)
		VALUES ($1, $2, $3, NOW())
	`, orgID, "Test Org", "test-org-"+orgID.String()[:8])
	require.NoError(t, err, "failed to create test organization")

	// Create user with pre-hashed password
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, org_id, email, name, password_hash, role, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW())
	`, userID, orgID, email, name, passwordHash, role)
	require.NoError(t, err, "failed to create test user")

	return orgID, userID
}

// Note: CreateTestSession and ValidateTestSession have been removed from testutil
// to avoid import cycles. Use the auth package directly in your tests:
//
//   sm := auth.NewSessionManager(...)
//   token, session, err := sm.CreateSession(ctx, auth.CreateSessionParams{...})
//   session, err := sm.ValidateSession(ctx, token)

// AssertSessionExists verifies that a session exists in the database.
func AssertSessionExists(t *testing.T, db *sql.DB, sessionID uuid.UUID) {
	t.Helper()

	var count int64
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM sessions WHERE id = $1`, sessionID).Scan(&count)
	require.NoError(t, err)

	if count != 1 {
		t.Errorf("expected session %s to exist, but it doesn't", sessionID)
	}
}

// AssertSessionNotExists verifies that a session does not exist in the database.
func AssertSessionNotExists(t *testing.T, db *sql.DB, sessionID uuid.UUID) {
	t.Helper()

	var count int64
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM sessions WHERE id = $1`, sessionID).Scan(&count)
	require.NoError(t, err)

	if count != 0 {
		t.Errorf("expected session %s to not exist, but it does", sessionID)
	}
}

// AssertUserExists verifies that a user exists in the database.
func AssertUserExists(t *testing.T, db *sql.DB, userID uuid.UUID) {
	t.Helper()

	var count int64
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM users WHERE id = $1`, userID).Scan(&count)
	require.NoError(t, err)

	if count != 1 {
		t.Errorf("expected user %s to exist, but it doesn't", userID)
	}
}

// AssertOrgExists verifies that an organization exists in the database.
func AssertOrgExists(t *testing.T, db *sql.DB, orgID uuid.UUID) {
	t.Helper()

	var count int64
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM organizations WHERE id = $1`, orgID).Scan(&count)
	require.NoError(t, err)

	if count != 1 {
		t.Errorf("expected organization %s to exist, but it doesn't", orgID)
	}
}

// GetUserByEmail retrieves a user by email address for testing.
func GetUserByEmail(t *testing.T, db *sql.DB, email string) (userID uuid.UUID, orgID uuid.UUID, role string) {
	t.Helper()

	ctx := context.Background()

	err := db.QueryRowContext(ctx, `
		SELECT id, org_id, role FROM users WHERE email = $1
	`, email).Scan(&userID, &orgID, &role)
	require.NoError(t, err, "failed to get user by email")

	return userID, orgID, role
}

// CountUserSessions returns the number of active sessions for a user.
func CountUserSessions(t *testing.T, db *sql.DB, userID uuid.UUID) int64 {
	t.Helper()

	var count int64
	err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&count)
	require.NoError(t, err)

	return count
}
