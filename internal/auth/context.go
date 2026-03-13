package auth

import "context"

type contextKey string

const sessionContextKey contextKey = "session"

// WithSession adds a session to the context.
func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

// SessionFromContext retrieves the session from the context.
func SessionFromContext(ctx context.Context) *Session {
	session, ok := ctx.Value(sessionContextKey).(*Session)
	if !ok {
		return nil
	}
	return session
}

// UserIDFromContext extracts the user ID from the session in the context.
func UserIDFromContext(ctx context.Context) string {
	session := SessionFromContext(ctx)
	if session == nil {
		return ""
	}
	return session.UserID
}

// UserEmailFromContext extracts the user email from the session in the context.
func UserEmailFromContext(ctx context.Context) string {
	session := SessionFromContext(ctx)
	if session == nil {
		return ""
	}
	return session.Email
}

// UserNameFromContext extracts the user name from the session in the context.
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
