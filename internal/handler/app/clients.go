package app

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// ClientsPageData extends PageData with client-specific data.
type ClientsPageData struct {
	PageData
	Clients       []repository.Client
	TotalClients  int64
	TotalProjects int64
	Search        string
	// EditClient is set when rendering the edit modal
	EditClient *repository.Client
}

// ListClients renders the clients list with optional search.
func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")

	var clients []repository.Client
	var err error

	if search != "" {
		clients, err = h.queries.SearchClients(ctx, search)
	} else {
		clients, err = h.queries.ListClients(ctx)
	}
	if err != nil {
		h.logger.Error("listing clients", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalClients, err := h.queries.CountClients(ctx)
	if err != nil {
		h.logger.Error("counting clients", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalProjects, err := h.queries.CountTotalClientProjects(ctx)
	if err != nil {
		h.logger.Error("counting client projects", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ClientsPageData{
		PageData:      h.pageData(r, "clients"),
		Clients:       clients,
		TotalClients:  totalClients,
		TotalProjects: totalProjects,
		Search:        search,
	}

	if isHTMXPartial(r) {
		if err := h.renderer.RenderPartial(w, "clients.html", "client-rows", data); err != nil {
			h.logger.Error("rendering client rows", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.renderer.Render(w, "clients.html", data); err != nil {
		h.logger.Error("rendering clients", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateClient handles POST /clients.
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	companyName := r.FormValue("company_name")
	if companyName == "" {
		http.Error(w, "Company name is required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateClient(r.Context(), repository.CreateClientParams{
		ID:          id,
		CompanyName: companyName,
		ContactName: toNullString(r.FormValue("contact_name")),
		Email:       toNullString(r.FormValue("email")),
		Phone:       toNullString(r.FormValue("phone")),
		Address:     toNullString(r.FormValue("address")),
		Notes:       toNullString(r.FormValue("notes")),
	})
	if err != nil {
		h.logger.Error("creating client", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

// GetClientEditForm returns the edit form partial for a client.
func (h *Handler) GetClientEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	client, err := h.queries.GetClient(r.Context(), id)
	if err != nil {
		h.logger.Error("getting client", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := ClientsPageData{
		PageData:   h.pageData(r, "clients"),
		EditClient: &client,
	}

	if err := h.renderer.RenderPartial(w, "clients.html", "edit-modal", data); err != nil {
		h.logger.Error("rendering edit form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateClient handles PUT /clients/{id}.
func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	companyName := r.FormValue("company_name")
	if companyName == "" {
		http.Error(w, "Company name is required", http.StatusBadRequest)
		return
	}

	_, err := h.queries.UpdateClient(r.Context(), repository.UpdateClientParams{
		ID:          id,
		CompanyName: companyName,
		ContactName: toNullString(r.FormValue("contact_name")),
		Email:       toNullString(r.FormValue("email")),
		Phone:       toNullString(r.FormValue("phone")),
		Address:     toNullString(r.FormValue("address")),
		Notes:       toNullString(r.FormValue("notes")),
	})
	if err != nil {
		h.logger.Error("updating client", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Also update client_name on any linked projects
	// (denormalized field for display purposes)
	// This is best-effort; projects table has client_name for quick listing
	if err := h.updateProjectClientNames(r, id, companyName); err != nil {
		h.logger.Warn("updating project client names", "error", err)
	}

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

// updateProjectClientNames updates the denormalized client_name on projects
// when a client's company name changes.
func (h *Handler) updateProjectClientNames(r *http.Request, clientID, newName string) error {
	// We don't have a specific query for this yet, so skip for now.
	// The denormalized name is set at project creation time.
	// TODO: Add UpdateProjectClientName query if needed.
	return nil
}

// DeleteClient handles DELETE /clients/{id}.
func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteClient(r.Context(), id); err != nil {
		h.logger.Error("deleting client", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", "/clients")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

// nullStr safely extracts a sql.NullString for template display.
func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
