package auth

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	goauth "github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// Subdomain must:
// - Start with a letter (not number)
// - Contain only lowercase letters, numbers, and hyphens
// - Be 3-63 characters long
// - End with a letter or number (not hyphen)
var subdomainRegex = regexp.MustCompile(`^[a-z][a-z0-9-]{1,61}[a-z0-9]$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// GetRegister renders the registration page.
func (h *Handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to home
	if goauth.IsAuthenticated(r.Context()) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := h.renderer.Render(w, "register", RegisterData{}); err != nil {
		slog.Error("failed to render register template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// PostRegister handles registration form submission.
func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	orgName := strings.TrimSpace(r.FormValue("org_name"))
	subdomain := strings.TrimSpace(strings.ToLower(r.FormValue("subdomain")))
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	data := RegisterData{
		OrgName:   orgName,
		Subdomain: subdomain,
		Name:      name,
		Email:     email,
	}

	// Validate inputs
	if orgName == "" {
		data.Error = "Company name is required"
		h.renderRegister(w, data)
		return
	}

	if subdomain == "" {
		data.Error = "Subdomain is required"
		h.renderRegister(w, data)
		return
	}

	if !subdomainRegex.MatchString(subdomain) {
		data.Error = "Subdomain must start with a letter, be 3-63 characters, and contain only lowercase letters, numbers, and hyphens"
		h.renderRegister(w, data)
		return
	}

	if name == "" {
		data.Error = "Your name is required"
		h.renderRegister(w, data)
		return
	}

	if email == "" {
		data.Error = "Email is required"
		h.renderRegister(w, data)
		return
	}

	if !emailRegex.MatchString(email) {
		data.Error = "Please enter a valid email address"
		h.renderRegister(w, data)
		return
	}

	if password == "" {
		data.Error = "Password is required"
		h.renderRegister(w, data)
		return
	}

	if len(password) < 8 {
		data.Error = "Password must be at least 8 characters"
		h.renderRegister(w, data)
		return
	}

	if password != passwordConfirm {
		data.Error = "Passwords do not match"
		h.renderRegister(w, data)
		return
	}

	// Check for existing email before creating anything
	_, err := h.queries.GetUserByEmailOnly(r.Context(), email)
	if err == nil {
		// User with this email already exists
		data.Error = "An account with this email already exists"
		h.renderRegister(w, data)
		return
	}
	// If error is not "no rows", something else went wrong
	if err.Error() != "sql: no rows in result set" {
		slog.Error("database error checking email", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Hash password
	passwordHash, err := goauth.HashPassword(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func() { _ = tx.Rollback() }()

	qtx := h.queries.WithTx(tx)

	// Create organization
	org, err := qtx.CreateOrganization(r.Context(), repository.CreateOrganizationParams{
		Name:      orgName,
		Subdomain: subdomain,
		Plan:      "free",
		Status:    "active",
	})
	if err != nil {
		// Check for duplicate subdomain
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			data.Error = "This subdomain is already taken"
			h.renderRegister(w, data)
			return
		}
		slog.Error("failed to create organization", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create user (first user is admin)
	user, err := qtx.CreateUser(r.Context(), repository.CreateUserParams{
		OrgID:        org.ID,
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		Role:         "admin",
		Status:       "active",
	})
	if err != nil {
		// Check for duplicate email
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			data.Error = "An account with this email already exists"
			h.renderRegister(w, data)
			return
		}
		slog.Error("failed to create user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create session
	token, _, err := h.sessionManager.CreateSession(r.Context(), goauth.CreateSessionParams{
		UserID:    user.ID,
		OrgID:     org.ID,
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

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
}

// renderRegister renders the registration page with data.
func (h *Handler) renderRegister(w http.ResponseWriter, data RegisterData) {
	if err := h.renderer.Render(w, "register", data); err != nil {
		slog.Error("failed to render register template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
