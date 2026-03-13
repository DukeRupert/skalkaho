package app

import (
	"log/slog"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/service/email"
	"github.com/dukerupert/skalkaho/internal/templates"
)

// PageData holds common data passed to all page templates.
type PageData struct {
	ActivePage string
	UserName   string
	Project    *ProjectStub
}

// ProjectStub holds minimal project info for sidebar rendering.
type ProjectStub struct {
	ID   string
	Name string
}

// Handler serves the main application pages.
type Handler struct {
	queries     *repository.Queries
	renderer    *templates.Renderer
	logger      *slog.Logger
	emailClient *email.Client
}

// NewHandler creates a new app handler.
func NewHandler(queries *repository.Queries, renderer *templates.Renderer, logger *slog.Logger) *Handler {
	return &Handler{
		queries:  queries,
		renderer: renderer,
		logger:   logger,
	}
}

func (h *Handler) pageData(r *http.Request, activePage string) PageData {
	name := auth.UserNameFromContext(r.Context())
	if name == "" {
		name = auth.UserEmailFromContext(r.Context())
	}
	return PageData{
		ActivePage: activePage,
		UserName:   name,
	}
}

// Projects redirects to ListProjects (kept for backward compat during phase transition).
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	h.ListProjects(w, r)
}

// Clients delegates to ListClients for the full client management page.
func (h *Handler) Clients(w http.ResponseWriter, r *http.Request) {
	h.ListClients(w, r)
}

// Materials delegates to ListMaterials for the full materials management page.
func (h *Handler) Materials(w http.ResponseWriter, r *http.Request) {
	h.ListMaterials(w, r)
}

// Rates delegates to ListRates for the full rates management page.
func (h *Handler) Rates(w http.ResponseWriter, r *http.Request) {
	h.ListRates(w, r)
}

// loadProject fetches a project and returns a ProjectStub for sidebar rendering.
func (h *Handler) loadProject(r *http.Request) *ProjectStub {
	projectID := r.PathValue("id")
	project, err := h.queries.GetProject(r.Context(), projectID)
	if err != nil {
		return &ProjectStub{ID: projectID, Name: "Unknown Project"}
	}
	return &ProjectStub{ID: project.ID, Name: project.Name}
}

// ProjectOverview delegates to GetProjectOverview for the full overview page.
func (h *Handler) ProjectOverview(w http.ResponseWriter, r *http.Request) {
	h.GetProjectOverview(w, r)
}

// EstimateBuilder renders the estimate builder page.
func (h *Handler) EstimateBuilder(w http.ResponseWriter, r *http.Request) {
	data := h.pageData(r, "estimate")
	data.Project = h.loadProject(r)
	if err := h.renderer.Render(w, "estimate_builder.html", data); err != nil {
		h.logger.Error("rendering estimate builder", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
