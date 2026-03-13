package auth

import (
	"log/slog"
	"net/http"
	"net/url"

	internalAuth "github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/templates"
)

// Handler handles authentication routes (login, logout).
type Handler struct {
	queries        *repository.Queries
	renderer       *templates.Renderer
	sessionManager *internalAuth.SessionManager
	logger         *slog.Logger
}

// NewHandler creates a new auth Handler.
func NewHandler(queries *repository.Queries, renderer *templates.Renderer, sm *internalAuth.SessionManager, logger *slog.Logger) *Handler {
	return &Handler{
		queries:        queries,
		renderer:       renderer,
		sessionManager: sm,
		logger:         logger,
	}
}

// loginData holds template data for the login page.
type loginData struct {
	Error    string
	Email    string
	Redirect string
}

// GetLogin renders the login form.
func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to home
	if internalAuth.IsAuthenticated(r.Context()) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := loginData{
		Redirect: r.URL.Query().Get("redirect"),
	}
	if err := h.renderer.Render(w, "login.html", data); err != nil {
		h.logger.Error("rendering login page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// PostLogin validates credentials, creates a session, and redirects.
func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	redirect := r.FormValue("redirect")

	renderError := func(msg string) {
		data := loginData{
			Error:    msg,
			Email:    email,
			Redirect: redirect,
		}
		w.WriteHeader(http.StatusUnauthorized)
		if err := h.renderer.Render(w, "login.html", data); err != nil {
			h.logger.Error("rendering login error", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	// Look up user
	user, err := h.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		renderError("Invalid email or password")
		return
	}

	// Verify password
	match, err := internalAuth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !match {
		renderError("Invalid email or password")
		return
	}

	// Check user is active
	if user.Status != "active" {
		renderError("Account is not active")
		return
	}

	// Create session
	token, _, err := h.sessionManager.CreateSession(r.Context(), internalAuth.CreateSessionParams{
		UserID:    user.ID,
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
	})
	if err != nil {
		h.logger.Error("creating session", "error", err, "email", email)
		renderError("Unable to sign in. Please try again.")
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

	// Redirect
	target := "/"
	if redirect != "" {
		// Validate redirect is a relative path
		parsed, err := url.Parse(redirect)
		if err == nil && parsed.Host == "" && parsed.Scheme == "" {
			target = redirect
		}
	}
	http.Redirect(w, r, target, http.StatusFound)
}

// Logout destroys the session and redirects to login.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie(h.sessionManager.CookieName())
	if err == nil {
		_ = h.sessionManager.DestroySession(r.Context(), cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.sessionManager.CookieName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.sessionManager.IsSecure(),
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}
