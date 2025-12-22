package quote

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// CreateCategory creates a new category.
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	parentID := sql.NullString{}
	if pid := r.FormValue("parent_id"); pid != "" {
		parentID = sql.NullString{String: pid, Valid: true}
	}

	// TODO: Check nesting depth

	category, err := h.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:               uuid.New().String(),
		JobID:            jobID,
		ParentID:         parentID,
		Name:             r.FormValue("name"),
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
	})
	if err != nil {
		logger.Error("failed to create category", "error", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	// Get job for template context
	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Return category partial for HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Remove empty state if it exists
		w.Header().Set("HX-Trigger", `{"removeEmptyState": true}`)

		data := map[string]interface{}{
			"Category":        category,
			"Job":             job,
			"LineItems":       []repository.LineItem{},
			"ChildCategories": []repository.Category{},
		}
		if err := h.renderer.RenderPartial(w, "category", data); err != nil {
			logger.Error("failed to render category", "error", err)
		}
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}

// UpdateCategory updates a category.
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get existing category to preserve job_id
	existingCat, err := h.queries.GetCategory(ctx, categoryID)
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	surchargePercent := sql.NullFloat64{}
	if sp := r.FormValue("surcharge_percent"); sp != "" {
		if val, err := strconv.ParseFloat(sp, 64); err == nil {
			surchargePercent = sql.NullFloat64{Float64: val, Valid: true}
		}
	}

	sortOrder := existingCat.SortOrder
	if so := r.FormValue("sort_order"); so != "" {
		if val, err := strconv.ParseInt(so, 10, 64); err == nil {
			sortOrder = val
		}
	}

	category, err := h.queries.UpdateCategory(ctx, repository.UpdateCategoryParams{
		ID:               categoryID,
		Name:             r.FormValue("name"),
		SurchargePercent: surchargePercent,
		SortOrder:        sortOrder,
	})
	if err != nil {
		logger.Error("failed to update category", "error", err)
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	// Return updated category for HTMX
	if r.Header.Get("HX-Request") == "true" {
		job, _ := h.queries.GetJob(ctx, category.JobID)
		lineItems, _ := h.queries.ListLineItemsByCategory(ctx, categoryID)
		childCats, _ := h.queries.ListChildCategories(ctx, sql.NullString{String: categoryID, Valid: true})

		data := map[string]interface{}{
			"Category":        category,
			"Job":             job,
			"LineItems":       lineItems,
			"ChildCategories": childCats,
		}
		if err := h.renderer.RenderPartial(w, "category", data); err != nil {
			logger.Error("failed to render category", "error", err)
		}
		return
	}

	http.Redirect(w, r, "/jobs/"+category.JobID, http.StatusSeeOther)
}

// DeleteCategory deletes a category.
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	if err := h.queries.DeleteCategory(ctx, categoryID); err != nil {
		logger.Error("failed to delete category", "error", err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	// Return empty response for HTMX
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
