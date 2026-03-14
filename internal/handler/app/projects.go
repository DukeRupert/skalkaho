package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// ProjectsPageData extends PageData with project-specific data.
type ProjectsPageData struct {
	PageData
	Projects []repository.Project
	Stats    repository.CountProjectsByStatusRow
	Clients  []repository.Client
	Status   string
	Search   string
}

// ListProjects renders the projects list with optional status filter and search.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	var projects []repository.Project
	var err error

	switch {
	case search != "" && status != "":
		projects, err = h.queries.SearchProjectsByStatus(ctx, repository.SearchProjectsByStatusParams{
			SearchTerm: search,
			Status:     status,
		})
	case search != "":
		projects, err = h.queries.SearchProjects(ctx, search)
	case status != "":
		projects, err = h.queries.ListProjectsByStatus(ctx, status)
	default:
		projects, err = h.queries.ListProjects(ctx)
	}
	if err != nil {
		h.logger.Error("listing projects", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	stats, err := h.queries.CountProjectsByStatus(ctx)
	if err != nil {
		h.logger.Error("counting projects", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	clients, err := h.queries.ListClients(ctx)
	if err != nil {
		h.logger.Error("listing clients", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ProjectsPageData{
		PageData: h.pageData(r, "projects"),
		Projects: projects,
		Stats:    stats,
		Clients:  clients,
		Status:   status,
		Search:   search,
	}

	// HTMX partial: return just the table rows
	if isHTMXPartial(r) {
		if err := h.renderer.RenderPartial(w, "projects.html", "project-rows", data); err != nil {
			h.logger.Error("rendering project rows", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.renderer.Render(w, "projects.html", data); err != nil {
		h.logger.Error("rendering projects", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewProjectModal renders the new project modal partial.
func (h *Handler) NewProjectModal(w http.ResponseWriter, r *http.Request) {
	clients, err := h.queries.ListClients(r.Context())
	if err != nil {
		h.logger.Error("listing clients", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Clients []repository.Client
	}{Clients: clients}

	if err := h.renderer.RenderPartial(w, "projects.html", "new-project-modal", data); err != nil {
		h.logger.Error("rendering new project modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateProject handles POST /projects.
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	clientName := r.FormValue("client_name")
	description := r.FormValue("description")

	// If client_id is set but client_name isn't, look up the client
	if clientID != "" && clientName == "" {
		client, err := h.queries.GetClient(r.Context(), clientID)
		if err == nil {
			clientName = client.CompanyName
		}
	}

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateProject(r.Context(), repository.CreateProjectParams{
		ID:          id,
		Name:        name,
		ClientID:    toNullString(clientID),
		ClientName:  toNullString(clientName),
		Description: toNullString(description),
		Status:      "Draft",
	})
	if err != nil {
		h.logger.Error("creating project", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the new project's overview
	http.Redirect(w, r, fmt.Sprintf("/projects/%s", id), http.StatusSeeOther)
}

// DeleteProject handles DELETE /projects/{id}.
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteProject(r.Context(), id); err != nil {
		h.logger.Error("deleting project", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// HTMX: return empty to remove the row
	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// UpdateProjectStatus handles PATCH /projects/{id}/status.
func (h *Handler) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	validStatuses := map[string]bool{
		"Draft": true, "In Review": true, "Active": true, "Completed": true,
	}
	if !validStatuses[status] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	if err := h.queries.UpdateProjectStatus(r.Context(), repository.UpdateProjectStatusParams{
		ID:     id,
		Status: status,
	}); err != nil {
		h.logger.Error("updating project status", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect back to the referring page (overview or project list)
	redirect := "/"
	if ref := r.Header.Get("HX-Current-URL"); ref != "" {
		redirect = ref
	} else if ref := r.Referer(); ref != "" {
		redirect = ref
	}

	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", redirect)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
