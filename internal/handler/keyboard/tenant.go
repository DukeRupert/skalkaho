package keyboard

import (
	"context"

	"github.com/google/uuid"
)

// TODO: Replace with actual tenant context extraction once auth is implemented.
// For now, returns a nil UUID which will need to be handled in queries.

// GetOrgID extracts the organization ID from the request context.
// Returns a nil UUID until authentication is implemented.
func GetOrgID(ctx context.Context) uuid.NullUUID {
	// TODO: Extract from session/JWT once auth is implemented
	// For now, return empty/null - queries need to handle this gracefully
	return uuid.NullUUID{}
}

// MustGetOrgID extracts the organization ID, panicking if not found.
// Use only in handlers that require authentication.
func MustGetOrgID(ctx context.Context) uuid.NullUUID {
	// TODO: Implement proper extraction and validation
	return uuid.NullUUID{}
}
