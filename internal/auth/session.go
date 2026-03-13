package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dukerupert/skalkaho/internal/repository"
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
	ID             string
	UserID         string
	Email          string
	Name           string
	ExpiresAt      time.Time
	LastActivityAt time.Time
	CreatedAt      time.Time
}

// CreateSessionParams contains parameters for creating a session.
type CreateSessionParams struct {
	UserID    string
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

	// Get user info
	user, err := sm.queries.GetUser(ctx, params.UserID)
	if err != nil {
		return "", nil, fmt.Errorf("getting user: %w", err)
	}

	// Calculate expiration
	expiresAt := time.Now().Add(sm.duration)

	// Create session in database
	dbSession, err := sm.queries.CreateSession(ctx, repository.CreateSessionParams{
		UserID:    params.UserID,
		TokenHash: tokenHash,
		UserAgent: sql.NullString{String: params.UserAgent, Valid: params.UserAgent != ""},
		IpAddress: sql.NullString{String: params.IPAddress, Valid: params.IPAddress != ""},
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", nil, fmt.Errorf("creating session in database: %w", err)
	}

	session := &Session{
		ID:             dbSession.ID,
		UserID:         dbSession.UserID,
		Email:          user.Email,
		Name:           user.Name,
		ExpiresAt:      dbSession.ExpiresAt,
		LastActivityAt: dbSession.LastActivityAt,
		CreatedAt:      dbSession.CreatedAt,
	}

	return token, session, nil
}

// ValidateSession validates a session token and returns the session details.
func (sm *SessionManager) ValidateSession(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	tokenHash := HashSessionToken(token)

	sessionRow, err := sm.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("retrieving session: %w", err)
	}

	if sessionRow.UserStatus != "active" {
		return nil, fmt.Errorf("user is not active")
	}

	// Update last activity (ignore error - session is still valid)
	now := time.Now()
	_ = sm.queries.UpdateSessionActivity(ctx, sessionRow.ID)

	session := &Session{
		ID:             sessionRow.ID,
		UserID:         sessionRow.UserID,
		Email:          sessionRow.Email,
		Name:           sessionRow.Name,
		ExpiresAt:      sessionRow.ExpiresAt,
		LastActivityAt: now,
		CreatedAt:      sessionRow.CreatedAt,
	}

	return session, nil
}

// DestroySession destroys a session given its token.
func (sm *SessionManager) DestroySession(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}

	tokenHash := HashSessionToken(token)
	_ = sm.queries.DeleteSessionByTokenHash(ctx, tokenHash)
	return nil
}

// DestroyAllUserSessions destroys all sessions for a given user.
func (sm *SessionManager) DestroyAllUserSessions(ctx context.Context, userID string) error {
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
