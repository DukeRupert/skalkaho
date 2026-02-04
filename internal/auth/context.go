package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const sessionContextKey contextKey = "session"

// WithSession adds a session to the context.
func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

// SessionFromContext retrieves the session from the context.
// Returns nil if no session is present.
func SessionFromContext(ctx context.Context) *Session {
	session, ok := ctx.Value(sessionContextKey).(*Session)
	if !ok {
		return nil
	}
	return session
}

// OrgIDFromContext extracts the organization ID from the session in the context.
// Returns uuid.Nil if no session is present.
func OrgIDFromContext(ctx context.Context) uuid.UUID {
	session := SessionFromContext(ctx)
	if session == nil {
		return uuid.Nil
	}
	return session.OrgID
}

// UserIDFromContext extracts the user ID from the session in the context.
// Returns uuid.Nil if no session is present.
func UserIDFromContext(ctx context.Context) uuid.UUID {
	session := SessionFromContext(ctx)
	if session == nil {
		return uuid.Nil
	}
	return session.UserID
}

// UserRoleFromContext extracts the user role from the session in the context.
// Returns empty string if no session is present.
func UserRoleFromContext(ctx context.Context) string {
	session := SessionFromContext(ctx)
	if session == nil {
		return ""
	}
	return session.Role
}

// UserEmailFromContext extracts the user email from the session in the context.
// Returns empty string if no session is present.
func UserEmailFromContext(ctx context.Context) string {
	session := SessionFromContext(ctx)
	if session == nil {
		return ""
	}
	return session.Email
}

// UserNameFromContext extracts the user name from the session in the context.
// Returns empty string if no session is present.
func UserNameFromContext(ctx context.Context) string {
	session := SessionFromContext(ctx)
	if session == nil {
		return ""
	}
	return session.Name
}

// IsAuthenticated checks if the context contains a valid session.
func IsAuthenticated(ctx context.Context) bool {
	return SessionFromContext(ctx) != nil
}
