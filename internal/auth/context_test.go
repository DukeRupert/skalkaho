package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestWithSession_SessionFromContext verifies context storage and retrieval of sessions.
func TestWithSession_SessionFromContext(t *testing.T) {
	t.Run("stores_and_retrieves_session", func(t *testing.T) {
		ctx := context.Background()

		orgID := uuid.New()
		userID := uuid.New()
		sessionID := uuid.New()

		// Create a test session
		session := &Session{
			ID:     sessionID,
			UserID: userID,
			OrgID:  orgID,
			Email:  "test@example.com",
			Name:   "Test User",
			Role:   "admin",
		}

		// Store session in context
		ctxWithSession := WithSession(ctx, session)

		// Retrieve session from context
		retrieved := SessionFromContext(ctxWithSession)
		assert.NotNil(t, retrieved, "session should be retrievable from context")
		assert.Equal(t, session.ID, retrieved.ID)
		assert.Equal(t, session.UserID, retrieved.UserID)
		assert.Equal(t, session.OrgID, retrieved.OrgID)
		assert.Equal(t, session.Email, retrieved.Email)
		assert.Equal(t, session.Name, retrieved.Name)
		assert.Equal(t, session.Role, retrieved.Role)
	})

	t.Run("returns_nil_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		session := SessionFromContext(ctx)
		assert.Nil(t, session, "should return nil when no session in context")
	})

	t.Run("handles_nil_session", func(t *testing.T) {
		ctx := context.Background()

		// Store nil session
		ctxWithNil := WithSession(ctx, nil)

		// Retrieve should return nil
		session := SessionFromContext(ctxWithNil)
		assert.Nil(t, session, "should handle nil session gracefully")
	})

	t.Run("multiple_context_layers", func(t *testing.T) {
		ctx := context.Background()

		session1 := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  uuid.New(),
			Email:  "user1@example.com",
			Name:   "User 1",
		}

		session2 := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  uuid.New(),
			Email:  "user2@example.com",
			Name:   "User 2",
		}

		// Create nested contexts
		ctx1 := WithSession(ctx, session1)
		ctx2 := WithSession(ctx1, session2)

		// Most recent session should be retrieved
		retrieved := SessionFromContext(ctx2)
		assert.NotNil(t, retrieved)
		assert.Equal(t, session2.Email, retrieved.Email, "should retrieve most recent session")

		// Original context should still have session1
		retrieved1 := SessionFromContext(ctx1)
		assert.NotNil(t, retrieved1)
		assert.Equal(t, session1.Email, retrieved1.Email)
	})
}

// TestOrgIDFromContext verifies extraction of org ID from session.
func TestOrgIDFromContext(t *testing.T) {
	t.Run("extracts_org_id_from_session", func(t *testing.T) {
		ctx := context.Background()
		orgID := uuid.New()

		session := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  orgID,
			Email:  "test@example.com",
		}

		ctxWithSession := WithSession(ctx, session)

		extractedOrgID := OrgIDFromContext(ctxWithSession)
		assert.Equal(t, orgID, extractedOrgID, "should extract correct org ID")
	})

	t.Run("returns_nil_uuid_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		orgID := OrgIDFromContext(ctx)
		assert.Equal(t, uuid.Nil, orgID, "should return nil UUID when no session")
	})

	t.Run("returns_nil_uuid_for_nil_session", func(t *testing.T) {
		ctx := WithSession(context.Background(), nil)

		orgID := OrgIDFromContext(ctx)
		assert.Equal(t, uuid.Nil, orgID, "should return nil UUID for nil session")
	})
}

// TestUserIDFromContext verifies extraction of user ID from session.
func TestUserIDFromContext(t *testing.T) {
	t.Run("extracts_user_id_from_session", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		session := &Session{
			ID:     uuid.New(),
			UserID: userID,
			OrgID:  uuid.New(),
			Email:  "test@example.com",
		}

		ctxWithSession := WithSession(ctx, session)

		extractedUserID := UserIDFromContext(ctxWithSession)
		assert.Equal(t, userID, extractedUserID, "should extract correct user ID")
	})

	t.Run("returns_nil_uuid_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		userID := UserIDFromContext(ctx)
		assert.Equal(t, uuid.Nil, userID, "should return nil UUID when no session")
	})
}

// TestUserRoleFromContext verifies extraction of user role from session.
func TestUserRoleFromContext(t *testing.T) {
	t.Run("extracts_role_from_session", func(t *testing.T) {
		ctx := context.Background()

		session := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  uuid.New(),
			Email:  "admin@example.com",
			Role:   "admin",
		}

		ctxWithSession := WithSession(ctx, session)

		role := UserRoleFromContext(ctxWithSession)
		assert.Equal(t, "admin", role, "should extract correct role")
	})

	t.Run("returns_empty_string_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		role := UserRoleFromContext(ctx)
		assert.Equal(t, "", role, "should return empty string when no session")
	})
}

// TestIsAuthenticated verifies authentication check.
func TestIsAuthenticated(t *testing.T) {
	t.Run("returns_true_when_session_exists", func(t *testing.T) {
		ctx := context.Background()

		session := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  uuid.New(),
			Email:  "test@example.com",
		}

		ctxWithSession := WithSession(ctx, session)

		assert.True(t, IsAuthenticated(ctxWithSession), "should return true when session exists")
	})

	t.Run("returns_false_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		assert.False(t, IsAuthenticated(ctx), "should return false when no session")
	})

	t.Run("returns_false_for_nil_session", func(t *testing.T) {
		ctx := WithSession(context.Background(), nil)

		assert.False(t, IsAuthenticated(ctx), "should return false for nil session")
	})
}

// TestUserEmailFromContext verifies extraction of user email from session.
func TestUserEmailFromContext(t *testing.T) {
	t.Run("extracts_email_from_session", func(t *testing.T) {
		ctx := context.Background()

		session := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  uuid.New(),
			Email:  "user@example.com",
		}

		ctxWithSession := WithSession(ctx, session)

		email := UserEmailFromContext(ctxWithSession)
		assert.Equal(t, "user@example.com", email, "should extract correct email")
	})

	t.Run("returns_empty_string_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		email := UserEmailFromContext(ctx)
		assert.Equal(t, "", email, "should return empty string when no session")
	})
}

// TestUserNameFromContext verifies extraction of user name from session.
func TestUserNameFromContext(t *testing.T) {
	t.Run("extracts_name_from_session", func(t *testing.T) {
		ctx := context.Background()

		session := &Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
			OrgID:  uuid.New(),
			Email:  "user@example.com",
			Name:   "John Doe",
		}

		ctxWithSession := WithSession(ctx, session)

		name := UserNameFromContext(ctxWithSession)
		assert.Equal(t, "John Doe", name, "should extract correct name")
	})

	t.Run("returns_empty_string_when_no_session", func(t *testing.T) {
		ctx := context.Background()

		name := UserNameFromContext(ctx)
		assert.Equal(t, "", name, "should return empty string when no session")
	})
}

// TestContextHelpersIntegration tests the context helpers work together.
func TestContextHelpersIntegration(t *testing.T) {
	ctx := context.Background()

	orgID := uuid.New()
	userID := uuid.New()

	session := &Session{
		ID:     uuid.New(),
		UserID: userID,
		OrgID:  orgID,
		Email:  "integration@example.com",
		Name:   "Integration Test User",
		Role:   "member",
	}

	// Store session
	ctxWithSession := WithSession(ctx, session)

	// Verify all extractors work
	assert.True(t, IsAuthenticated(ctxWithSession))
	assert.Equal(t, orgID, OrgIDFromContext(ctxWithSession))
	assert.Equal(t, userID, UserIDFromContext(ctxWithSession))
	assert.Equal(t, "integration@example.com", UserEmailFromContext(ctxWithSession))
	assert.Equal(t, "Integration Test User", UserNameFromContext(ctxWithSession))
	assert.Equal(t, "member", UserRoleFromContext(ctxWithSession))

	// Verify SessionFromContext returns the same session
	retrievedSession := SessionFromContext(ctxWithSession)
	assert.NotNil(t, retrievedSession)
	assert.Equal(t, session.ID, retrievedSession.ID)
	assert.Equal(t, session.UserID, retrievedSession.UserID)
	assert.Equal(t, session.OrgID, retrievedSession.OrgID)
	assert.Equal(t, session.Email, retrievedSession.Email)
	assert.Equal(t, session.Name, retrievedSession.Name)
	assert.Equal(t, session.Role, retrievedSession.Role)
}
