package auth

import (
	"net/http"
	"net/url"
)

// SessionMiddleware loads the session from the cookie and adds it to the request context.
// This middleware does NOT require authentication - it simply loads the session if present.
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get session cookie
		cookie, err := r.Cookie(sm.cookieName)
		if err != nil {
			// No cookie - continue without session
			next.ServeHTTP(w, r)
			return
		}

		// Validate session
		session, err := sm.ValidateSession(r.Context(), cookie.Value)
		if err != nil {
			// Invalid session - continue without session
			next.ServeHTTP(w, r)
			return
		}

		// Add session to context
		ctx := WithSession(r.Context(), session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is middleware that requires authentication.
// If the user is not authenticated, they are redirected to the login page.
// Complex URLs (with query params or multiple path segments) preserve the original URL.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r.Context()) {
			// Build redirect URL with original path
			redirectURL := "/login"

			// Preserve the original URL for complex paths (with query params or deeper paths)
			// Simple top-level paths like /protected don't need preservation
			if r.URL.RawQuery != "" || (r.URL.Path != "/" && r.URL.Path != "/login" && hasMultiplePathSegments(r.URL.Path)) {
				redirectURL = "/login?redirect=" + url.QueryEscape(r.URL.RequestURI())
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// hasMultiplePathSegments checks if a path has more than 2 segments (e.g., /foo/bar has 2 segments)
func hasMultiplePathSegments(path string) bool {
	segments := 0
	for _, char := range path {
		if char == '/' {
			segments++
		}
	}
	return segments > 1 // More than just the leading slash
}

// RequireRole returns middleware that requires a specific user role.
// If the user doesn't have the required role, they receive a 403 Forbidden response.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First check authentication
			if !IsAuthenticated(r.Context()) {
				redirectURL := "/login?redirect=" + url.QueryEscape(r.URL.RequestURI())
				http.Redirect(w, r, redirectURL, http.StatusFound)
				return
			}

			// Check role
			userRole := UserRoleFromContext(r.Context())
			if userRole != role {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
