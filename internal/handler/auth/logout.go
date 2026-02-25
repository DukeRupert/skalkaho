package auth

import (
	"net/http"
)

// PostLogout handles logout requests.
func (h *Handler) PostLogout(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	cookie, err := r.Cookie(h.sessionManager.CookieName())
	if err == nil && cookie.Value != "" {
		// Destroy session in database (ignore error - we're logging out anyway)
		_ = h.sessionManager.DestroySession(r.Context(), cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.sessionManager.CookieName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.sessionManager.IsSecure(),
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusFound)
}
