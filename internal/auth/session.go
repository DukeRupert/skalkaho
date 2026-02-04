package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// SessionManager manages user sessions.
type SessionManager struct {
	queries    *repository.Queries
	secret     string
	duration   time.Duration
	cookieName string
	secure     bool
}

// Session represents a user session with associated user information.
type Session struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrgID          uuid.UUID
	Email          string
	Name           string
	Role           string
	ExpiresAt      time.Time
	LastActivityAt time.Time
	CreatedAt      time.Time
}

// CreateSessionParams contains parameters for creating a session.
type CreateSessionParams struct {
	UserID    uuid.UUID
	OrgID     uuid.UUID
	UserAgent string
	IPAddress string
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager(db repository.DBTX, secret string, duration time.Duration, cookieName string, secure bool) *SessionManager {
	return &SessionManager{
		queries:    repository.New(db),
		secret:     secret,
		duration:   duration,
		cookieName: cookieName,
		secure:     secure,
	}
}

// CreateSession creates a new session for a user.
// Returns the plaintext token (to be set in cookie) and the session details.
func (sm *SessionManager) CreateSession(ctx context.Context, params CreateSessionParams) (string, *Session, error) {
	// Generate session token
	token, err := GenerateSessionToken()
	if err != nil {
		return "", nil, fmt.Errorf("generating session token: %w", err)
	}

	// Hash token for storage
	tokenHash := HashSessionToken(token)

	// Get user info first (before creating session)
	user, err := sm.queries.GetUser(ctx, params.UserID)
	if err != nil {
		return "", nil, fmt.Errorf("getting user: %w", err)
	}

	// Validate user belongs to the specified organization
	if user.OrgID != params.OrgID {
		return "", nil, fmt.Errorf("user does not belong to specified organization")
	}

	// Calculate expiration
	expiresAt := time.Now().Add(sm.duration)

	// Create session in database
	dbSession, err := sm.queries.CreateSession(ctx, repository.CreateSessionParams{
		UserID:    params.UserID,
		OrgID:     params.OrgID,
		TokenHash: tokenHash,
		UserAgent: sql.NullString{String: params.UserAgent, Valid: params.UserAgent != ""},
		IpAddress: sql.NullString{String: params.IPAddress, Valid: params.IPAddress != ""},
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", nil, fmt.Errorf("creating session in database: %w", err)
	}

	// Build session with user info
	session := &Session{
		ID:             dbSession.ID,
		UserID:         dbSession.UserID,
		OrgID:          dbSession.OrgID,
		Email:          user.Email,
		Name:           user.Name,
		Role:           user.Role,
		ExpiresAt:      dbSession.ExpiresAt,
		LastActivityAt: dbSession.LastActivityAt,
		CreatedAt:      dbSession.CreatedAt,
	}

	return token, session, nil
}

// ValidateSession validates a session token and returns the session details.
// Returns an error if the token is invalid, expired, or the user is inactive.
func (sm *SessionManager) ValidateSession(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	// Hash token to look up in database
	tokenHash := HashSessionToken(token)

	// Retrieve session (query already checks expiration)
	sessionRow, err := sm.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("retrieving session: %w", err)
	}

	// Check user status
	if sessionRow.UserStatus != "active" {
		return nil, fmt.Errorf("user is not active")
	}

	// Update last activity
	now := time.Now()
	if err := sm.queries.UpdateSessionActivity(ctx, sessionRow.ID); err != nil {
		// Log but don't fail - session is still valid
		// In production, you'd use a proper logger here
	}

	session := &Session{
		ID:             sessionRow.ID,
		UserID:         sessionRow.UserID,
		OrgID:          sessionRow.OrgID,
		Email:          sessionRow.Email,
		Name:           sessionRow.Name,
		Role:           sessionRow.Role,
		ExpiresAt:      sessionRow.ExpiresAt,
		LastActivityAt: now, // Use the current time since we just updated it
		CreatedAt:      sessionRow.CreatedAt,
	}

	return session, nil
}

// DestroySession destroys a session given its token.
// This operation is idempotent - destroying a non-existent session does not error.
func (sm *SessionManager) DestroySession(ctx context.Context, token string) error {
	if token == "" {
		return nil // Idempotent
	}

	tokenHash := HashSessionToken(token)
	if err := sm.queries.DeleteSessionByTokenHash(ctx, tokenHash); err != nil {
		// Even if session doesn't exist, don't error (idempotent)
		return nil
	}

	return nil
}

// DestroyAllUserSessions destroys all sessions for a given user.
func (sm *SessionManager) DestroyAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	if err := sm.queries.DeleteUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("deleting user sessions: %w", err)
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions from the database.
func (sm *SessionManager) CleanupExpiredSessions(ctx context.Context) error {
	if err := sm.queries.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("cleaning up expired sessions: %w", err)
	}
	return nil
}

// CookieName returns the session cookie name.
func (sm *SessionManager) CookieName() string {
	return sm.cookieName
}

// Duration returns the session duration.
func (sm *SessionManager) Duration() time.Duration {
	return sm.duration
}

// IsSecure returns whether cookies should be marked as Secure.
func (sm *SessionManager) IsSecure() bool {
	return sm.secure
}
