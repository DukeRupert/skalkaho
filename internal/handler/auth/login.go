package auth

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	goauth "github.com/dukerupert/skalkaho/internal/auth"
)

// GetLogin renders the login page.
func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to home
	if goauth.IsAuthenticated(r.Context()) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := LoginData{
		Redirect: r.URL.Query().Get("redirect"),
	}

	if err := h.renderer.Render(w, "login", data); err != nil {
		slog.Error("failed to render login template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// PostLogin handles login form submission.
func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	password := r.FormValue("password")
	redirect := r.FormValue("redirect")

	// Validate inputs
	if email == "" || password == "" {
		h.renderLoginError(w, "Email and password are required", email, redirect)
		return
	}

	// Get user by email (across all orgs - email must be unique globally)
	user, err := h.queries.GetUserByEmailOnly(r.Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			h.renderLoginError(w, "Invalid email or password", email, redirect)
			return
		}
		slog.Error("database error during login", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check user status
	if user.Status != "active" {
		h.renderLoginError(w, "Your account is not active. Please contact support.", email, redirect)
		return
	}

	// Verify password
	match, err := goauth.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		slog.Error("password verification error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !match {
		h.renderLoginError(w, "Invalid email or password", email, redirect)
		return
	}

	// Create session
	token, _, err := h.sessionManager.CreateSession(r.Context(), goauth.CreateSessionParams{
		UserID:    user.ID,
		OrgID:     user.OrgID,
		UserAgent: r.Header.Get("User-Agent"),
		IPAddress: strings.Split(r.RemoteAddr, ":")[0],
	})
	if err != nil {
		slog.Error("failed to create session", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.sessionManager.CookieName(),
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.sessionManager.Duration().Seconds()),
		HttpOnly: true,
		Secure:   h.sessionManager.IsSecure(),
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to original URL or home
	redirectURL := "/"
	if redirect != "" {
		redirectURL = redirect
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// renderLoginError renders the login page with an error message.
func (h *Handler) renderLoginError(w http.ResponseWriter, errorMsg, email, redirect string) {
	data := LoginData{
		Error:    errorMsg,
		Email:    email,
		Redirect: redirect,
	}
	if err := h.renderer.Render(w, "login", data); err != nil {
		slog.Error("failed to render login template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
