package auth

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// Renderer interface for template rendering.
type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{}) error
}

// DB interface for transaction support.
type DB interface {
	repository.DBTX
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Handler handles authentication-related HTTP requests.
type Handler struct {
	db             DB
	queries        *repository.Queries
	renderer       Renderer
	sessionManager *auth.SessionManager
}

// NewHandler creates a new auth handler.
func NewHandler(db DB, renderer Renderer, sessionManager *auth.SessionManager) *Handler {
	return &Handler{
		db:             db,
		queries:        repository.New(db),
		renderer:       renderer,
		sessionManager: sessionManager,
	}
}

// LoginData contains data for the login template.
type LoginData struct {
	Error    string
	Email    string
	Redirect string
}

// RegisterData contains data for the register template.
type RegisterData struct {
	Error     string
	OrgName   string
	Subdomain string
	Name      string
	Email     string
}
