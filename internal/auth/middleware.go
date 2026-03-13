package auth

import (
	"net/http"
	"net/url"
)

// SessionMiddleware loads the session from the cookie and adds it to the request context.
// This middleware does NOT require authentication - it simply loads the session if present.
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sm.cookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := sm.ValidateSession(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := WithSession(r.Context(), session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is middleware that requires authentication.
// If the user is not authenticated, they are redirected to the login page.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r.Context()) {
			redirectURL := "/login"
			if r.URL.RawQuery != "" || (r.URL.Path != "/" && r.URL.Path != "/login" && hasMultiplePathSegments(r.URL.Path)) {
				redirectURL = "/login?redirect=" + url.QueryEscape(r.URL.RequestURI())
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func hasMultiplePathSegments(path string) bool {
	segments := 0
	for _, char := range path {
		if char == '/' {
			segments++
		}
	}
	return segments > 1
}
