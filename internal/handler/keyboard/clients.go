package keyboard

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

const clientsPageSize = 20

// ListClients shows the clients management page with search and pagination.
func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	// Parse query params
	search := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := int64((page - 1) * clientsPageSize)

	// Get total count for pagination
	totalCount, err := h.queries.CountClients(ctx, search)
	if err != nil {
		logger.Error("failed to count clients", "error", err)
		http.Error(w, "Failed to load clients", http.StatusInternalServerError)
		return
	}

	totalPages := int(totalCount+int64(clientsPageSize)-1) / clientsPageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Get paginated clients
	clients, err := h.queries.ListClientsPaginated(ctx, repository.ListClientsPaginatedParams{
		Search: search,
		Offset: offset,
		Limit:  int64(clientsPageSize),
	})
	if err != nil {
		logger.Error("failed to list clients", "error", err)
		http.Error(w, "Failed to load clients", http.StatusInternalServerError)
		return
	}

	pagination := PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalCount,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}

	data := map[string]interface{}{
		"Clients":    clients,
		"Search":     search,
		"Pagination": pagination,
	}

	if err := h.renderer.Render(w, "clients_list", data); err != nil {
		logger.Error("failed to render clients page", "error", err)
	}
}

// GetClient shows the client detail/edit page.
func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	client, err := h.queries.GetClient(ctx, id)
	if err != nil {
		logger.Error("failed to get client", "error", err, "id", id)
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Get jobs associated with this client
	jobs, err := h.queries.ListJobs(ctx)
	if err != nil {
		logger.Error("failed to list jobs", "error", err)
	}

	// Filter jobs for this client
	var clientJobs []repository.Job
	for _, job := range jobs {
		if job.ClientID.Valid && job.ClientID.String == id {
			clientJobs = append(clientJobs, job)
		}
	}

	// Check if client can be deleted
	hasJobs, _ := h.queries.ClientHasJobs(ctx, sql.NullString{String: id, Valid: true})

	data := map[string]interface{}{
		"Client":  client,
		"Jobs":    clientJobs,
		"HasJobs": hasJobs,
	}

	if err := h.renderer.Render(w, "client", data); err != nil {
		logger.Error("failed to render client page", "error", err)
	}
}

// GetClientForm returns the inline form for creating a new client.
func (h *Handler) GetClientForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "client_form", nil); err != nil {
		logger.Error("failed to render client form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// CreateClient creates a new client.
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Check for duplicate name
	_, err := h.queries.GetClientByName(ctx, name)
	if err == nil {
		http.Error(w, "A client with this name already exists", http.StatusConflict)
		return
	}

	client, err := h.queries.CreateClient(ctx, repository.CreateClientParams{
		ID:      uuid.New().String(),
		Name:    name,
		Company: toNullString(r.FormValue("company")),
		Email:   toNullString(r.FormValue("email")),
		Phone:   toNullString(r.FormValue("phone")),
		Address: toNullString(r.FormValue("address")),
		City:    toNullString(r.FormValue("city")),
		State:   toNullString(r.FormValue("state")),
		Zip:     toNullString(r.FormValue("zip")),
		TaxID:   toNullString(r.FormValue("tax_id")),
		Notes:   toNullString(r.FormValue("notes")),
	})
	if err != nil {
		logger.Error("failed to create client", "error", err)
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	// Redirect to client detail page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/clients/"+client.ID)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/clients/"+client.ID, http.StatusSeeOther)
}

// GetClientEditForm returns the inline form for editing a client.
func (h *Handler) GetClientEditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	client, err := h.queries.GetClient(ctx, id)
	if err != nil {
		logger.Error("failed to get client", "error", err, "id", id)
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Client": client,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "client_edit_form", data); err != nil {
		logger.Error("failed to render client edit form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateClient updates an existing client.
func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Check for duplicate name (excluding current client)
	existing, err := h.queries.GetClientByName(ctx, name)
	if err == nil && existing.ID != id {
		http.Error(w, "A client with this name already exists", http.StatusConflict)
		return
	}

	_, err = h.queries.UpdateClient(ctx, repository.UpdateClientParams{
		ID:      id,
		Name:    name,
		Company: toNullString(r.FormValue("company")),
		Email:   toNullString(r.FormValue("email")),
		Phone:   toNullString(r.FormValue("phone")),
		Address: toNullString(r.FormValue("address")),
		City:    toNullString(r.FormValue("city")),
		State:   toNullString(r.FormValue("state")),
		Zip:     toNullString(r.FormValue("zip")),
		TaxID:   toNullString(r.FormValue("tax_id")),
		Notes:   toNullString(r.FormValue("notes")),
	})
	if err != nil {
		logger.Error("failed to update client", "error", err)
		http.Error(w, "Failed to update client", http.StatusInternalServerError)
		return
	}

	// Redirect back to client detail
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/clients/"+id)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/clients/"+id, http.StatusSeeOther)
}

// DeleteClient deletes a client.
func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	// Check if client has jobs
	hasJobs, err := h.queries.ClientHasJobs(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		logger.Error("failed to check client jobs", "error", err)
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}

	if hasJobs {
		http.Error(w, "Cannot delete client with associated quotes", http.StatusConflict)
		return
	}

	if err := h.queries.DeleteClient(ctx, id); err != nil {
		logger.Error("failed to delete client", "error", err)
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}

	// Redirect to clients list
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/clients")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

// toNullString converts a string to sql.NullString.
func toNullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
