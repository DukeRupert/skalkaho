package quote

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/templates"
	"github.com/google/uuid"
)

// Handler handles quote-related HTTP requests.
type Handler struct {
	queries  *repository.Queries
	renderer *templates.Renderer
	logger   *slog.Logger
}

// NewHandler creates a new quote handler.
func NewHandler(queries *repository.Queries, renderer *templates.Renderer, logger *slog.Logger) *Handler {
	return &Handler{
		queries:  queries,
		renderer: renderer,
		logger:   logger,
	}
}

// ListJobs shows the jobs list page.
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	jobs, err := h.queries.ListJobs(ctx)
	if err != nil {
		logger.Error("failed to list jobs", "error", err)
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Jobs": jobs,
	}

	if err := h.renderer.Render(w, "jobs_list", data); err != nil {
		logger.Error("failed to render jobs list", "error", err)
	}
}

// CreateJob creates a new job.
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	// Get default settings
	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		logger.Error("failed to get settings", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Create job with defaults
	job, err := h.queries.CreateJob(ctx, repository.CreateJobParams{
		ID:               uuid.New().String(),
		Name:             "New Quote",
		CustomerName:     sql.NullString{},
		SurchargePercent: settings.DefaultSurchargePercent,
		SurchargeMode:    settings.DefaultSurchargeMode,
	})
	if err != nil {
		logger.Error("failed to create job", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Redirect to the new job page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/m/jobs/"+job.ID)
		return
	}

	http.Redirect(w, r, "/m/jobs/"+job.ID, http.StatusSeeOther)
}

// GetJob shows a single job with all its details.
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	// Get job
	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Get all categories for this job
	categories, err := h.queries.ListCategoriesByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Get all line items for this job
	lineItems, err := h.queries.ListLineItemsByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load line items", http.StatusInternalServerError)
		return
	}

	// Build data structures for template
	topLevelCategories := make([]repository.Category, 0)
	childCategories := make(map[string][]repository.Category)
	lineItemsByCategory := make(map[string][]repository.LineItem)

	for _, cat := range categories {
		if !cat.ParentID.Valid {
			topLevelCategories = append(topLevelCategories, cat)
		} else {
			childCategories[cat.ParentID.String] = append(childCategories[cat.ParentID.String], cat)
		}
	}

	for _, item := range lineItems {
		lineItemsByCategory[item.CategoryID] = append(lineItemsByCategory[item.CategoryID], item)
	}

	// Calculate totals
	totals := h.calculateTotals(job, categories, lineItems)

	data := map[string]interface{}{
		"Job":                 job,
		"Categories":          topLevelCategories,
		"ChildCategories":     childCategories,
		"LineItemsByCategory": lineItemsByCategory,
		"Totals":              totals,
	}

	if err := h.renderer.Render(w, "job", data); err != nil {
		logger.Error("failed to render job page", "error", err)
	}
}

// UpdateJob updates a job.
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	surchargePercent, _ := strconv.ParseFloat(r.FormValue("surcharge_percent"), 64)

	customerName := sql.NullString{}
	if cn := r.FormValue("customer_name"); cn != "" {
		customerName = sql.NullString{String: cn, Valid: true}
	}

	job, err := h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:               jobID,
		Name:             r.FormValue("name"),
		CustomerName:     customerName,
		SurchargePercent: surchargePercent,
		SurchargeMode:    r.FormValue("surcharge_mode"),
	})
	if err != nil {
		logger.Error("failed to update job", "error", err)
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	// Return updated job header for HTMX
	if r.Header.Get("HX-Request") == "true" {
		if err := h.renderer.RenderPartial(w, "job_header", job); err != nil {
			logger.Error("failed to render job header", "error", err)
		}
		return
	}

	http.Redirect(w, r, "/m/jobs/"+jobID, http.StatusSeeOther)
}

// DeleteJob deletes a job.
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	if err := h.queries.DeleteJob(ctx, jobID); err != nil {
		logger.Error("failed to delete job", "error", err)
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	// Return empty response for HTMX to remove element
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/m/", http.StatusSeeOther)
}

// calculateTotals computes job totals from repository types.
func (h *Handler) calculateTotals(job repository.Job, categories []repository.Category, lineItems []repository.LineItem) domain.JobTotal {
	// Convert to domain types
	domainJob := &domain.Job{
		ID:               job.ID,
		SurchargePercent: job.SurchargePercent,
		SurchargeMode:    domain.SurchargeMode(job.SurchargeMode),
	}

	domainCategories := make([]*domain.Category, len(categories))
	for i, cat := range categories {
		var parentID *string
		if cat.ParentID.Valid {
			parentID = &cat.ParentID.String
		}
		var surcharge *float64
		if cat.SurchargePercent.Valid {
			surcharge = &cat.SurchargePercent.Float64
		}
		domainCategories[i] = &domain.Category{
			ID:               cat.ID,
			JobID:            cat.JobID,
			ParentID:         parentID,
			SurchargePercent: surcharge,
		}
	}

	domainLineItems := make([]*domain.LineItem, len(lineItems))
	for i, item := range lineItems {
		var surcharge *float64
		if item.SurchargePercent.Valid {
			surcharge = &item.SurchargePercent.Float64
		}
		domainLineItems[i] = &domain.LineItem{
			ID:               item.ID,
			CategoryID:       item.CategoryID,
			Type:             domain.LineItemType(item.Type),
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			SurchargePercent: surcharge,
		}
	}

	return domain.CalculateJobTotal(domainJob, domainCategories, domainLineItems)
}
