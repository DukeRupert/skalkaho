package keyboard

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// ListJobs shows the keyboard-centric jobs list.
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	jobs, err := h.queries.ListJobs(ctx)
	if err != nil {
		logger.Error("failed to list jobs", "error", err)
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	// Calculate totals for each job
	type JobWithTotal struct {
		repository.Job
		GrandTotal float64
	}

	jobsWithTotals := make([]JobWithTotal, len(jobs))
	for i, job := range jobs {
		categories, _ := h.queries.ListCategoriesByJob(ctx, job.ID)
		lineItems, _ := h.queries.ListLineItemsByJob(ctx, job.ID)
		totals := h.calculateTotals(job, categories, lineItems)
		jobsWithTotals[i] = JobWithTotal{
			Job:        job,
			GrandTotal: totals.GrandTotal,
		}
	}

	data := map[string]interface{}{
		"Jobs":          jobsWithTotals,
		"SelectedIndex": 0,
	}

	if err := h.renderer.Render(w, "jobs_list", data); err != nil {
		logger.Error("failed to render jobs list", "error", err)
	}
}

// GetJob shows a single job with its categories.
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

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

	categories, err := h.queries.ListCategoriesByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load line items", http.StatusInternalServerError)
		return
	}

	// Get only top-level categories
	topLevelCategories := make([]repository.Category, 0)
	for _, cat := range categories {
		if !cat.ParentID.Valid {
			topLevelCategories = append(topLevelCategories, cat)
		}
	}

	// Calculate totals for each category
	type CategoryWithTotal struct {
		repository.Category
		Total float64
	}

	categoriesWithTotals := make([]CategoryWithTotal, len(topLevelCategories))
	for i, cat := range topLevelCategories {
		catTotal := h.calculateCategoryTotal(cat.ID, job, categories, lineItems)
		categoriesWithTotals[i] = CategoryWithTotal{
			Category: cat,
			Total:    catTotal.Total,
		}
	}

	totals := h.calculateTotals(job, categories, lineItems)

	// Build category tree for sidebar navigation
	categoryTree := buildCategoryTree(categories)

	data := map[string]interface{}{
		"Job":               job,
		"Categories":        categoriesWithTotals,
		"Totals":            totals,
		"SelectedIndex":     0,
		"CategoryTree":      categoryTree,
		"CurrentCategoryID": "",
	}

	if err := h.renderer.Render(w, "job", data); err != nil {
		logger.Error("failed to render job page", "error", err)
	}
}

// CreateJob creates a new job and redirects to it.
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "New Quote"
	}

	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		logger.Error("failed to get settings", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.CreateJob(ctx, repository.CreateJobParams{
		ID:               uuid.New().String(),
		Name:             name,
		CustomerName:     sql.NullString{},
		SurchargePercent: settings.DefaultSurchargePercent,
		SurchargeMode:    settings.DefaultSurchargeMode,
	})
	if err != nil {
		logger.Error("failed to create job", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+job.ID)
		return
	}

	http.Redirect(w, r, "/jobs/"+job.ID, http.StatusSeeOther)
}

// UpdateJob updates a job's details.
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

	_, err := h.queries.UpdateJob(ctx, repository.UpdateJobParams{
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

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+jobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
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

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetJobForm returns an inline form for creating jobs.
func (h *Handler) GetJobForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "job_form", nil); err != nil {
		logger.Error("failed to render job form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// GetMarkupForm returns an inline form for editing job markup.
func (h *Handler) GetMarkupForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Job": job,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "markup_form", data); err != nil {
		logger.Error("failed to render markup form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateMarkup updates a job's markup percentage.
func (h *Handler) UpdateMarkup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	surchargePercent, _ := strconv.ParseFloat(r.FormValue("surcharge_percent"), 64)

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:               jobID,
		Name:             job.Name,
		CustomerName:     job.CustomerName,
		SurchargePercent: surchargePercent,
		SurchargeMode:    job.SurchargeMode,
	})
	if err != nil {
		logger.Error("failed to update job markup", "error", err)
		http.Error(w, "Failed to update markup", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+jobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}
