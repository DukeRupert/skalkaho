package keyboard

import (
	"context"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/google/uuid"
)

// GetOrgID extracts the organization ID from the request context.
// Returns a valid UUID if the user is authenticated, otherwise returns a null UUID.
func GetOrgID(ctx context.Context) uuid.NullUUID {
	orgID := auth.OrgIDFromContext(ctx)
	if orgID == uuid.Nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: orgID, Valid: true}
}

// MustGetOrgID extracts the organization ID, panicking if not found.
// Use only in handlers that require authentication (protected by RequireAuth middleware).
func MustGetOrgID(ctx context.Context) uuid.NullUUID {
	orgID := auth.OrgIDFromContext(ctx)
	if orgID == uuid.Nil {
		panic("MustGetOrgID called without authenticated session")
	}
	return uuid.NullUUID{UUID: orgID, Valid: true}
}

// GetUserID extracts the user ID from the request context.
// Returns uuid.Nil if not authenticated.
func GetUserID(ctx context.Context) uuid.UUID {
	return auth.UserIDFromContext(ctx)
}

// GetUserEmail extracts the user's email from the request context.
// Returns empty string if not authenticated.
func GetUserEmail(ctx context.Context) string {
	return auth.UserEmailFromContext(ctx)
}

// GetUserName extracts the user's name from the request context.
// Returns empty string if not authenticated.
func GetUserName(ctx context.Context) string {
	return auth.UserNameFromContext(ctx)
}

// GetUserRole extracts the user's role from the request context.
// Returns empty string if not authenticated.
func GetUserRole(ctx context.Context) string {
	return auth.UserRoleFromContext(ctx)
}

// IsAuthenticated checks if the request has an authenticated session.
func IsAuthenticated(ctx context.Context) bool {
	return auth.IsAuthenticated(ctx)
}
