package quote

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

const maxCategoryDepth = 3

// getCategoryDepth calculates how deep a category is (1 = top level)
func (h *Handler) getCategoryDepth(ctx context.Context, categoryID string) (int, error) {
	depth := 1
	currentID := categoryID

	for depth <= maxCategoryDepth {
		cat, err := h.queries.GetCategory(ctx, currentID)
		if err != nil {
			return 0, err
		}
		if !cat.ParentID.Valid {
			return depth, nil
		}
		currentID = cat.ParentID.String
		depth++
	}
	return depth, nil
}

// CreateCategory creates a new category.
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	category, err := h.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:               uuid.New().String(),
		JobID:            jobID,
		ParentID:         sql.NullString{},
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
		w.Header().Set("HX-Trigger", `{"removeEmptyState": true}`)

		data := map[string]interface{}{
			"Category":        category,
			"Job":             job,
			"LineItems":       []repository.LineItem{},
			"ChildCategories": []repository.Category{},
			"Depth":           1,
		}
		if err := h.renderer.RenderPartial(w, "category", data); err != nil {
			logger.Error("failed to render category", "error", err)
		}
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}

// CreateSubcategory creates a new subcategory under a parent.
func (h *Handler) CreateSubcategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	parentID := r.PathValue("parentID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get parent category to find job_id and check depth
	parent, err := h.queries.GetCategory(ctx, parentID)
	if err != nil {
		logger.Error("failed to get parent category", "error", err)
		http.Error(w, "Parent category not found", http.StatusNotFound)
		return
	}

	// Check nesting depth
	parentDepth, err := h.getCategoryDepth(ctx, parentID)
	if err != nil {
		logger.Error("failed to get category depth", "error", err)
		http.Error(w, "Failed to check depth", http.StatusInternalServerError)
		return
	}

	if parentDepth >= maxCategoryDepth {
		http.Error(w, "Maximum nesting depth reached", http.StatusBadRequest)
		return
	}

	category, err := h.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:               uuid.New().String(),
		JobID:            parent.JobID,
		ParentID:         sql.NullString{String: parentID, Valid: true},
		Name:             r.FormValue("name"),
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
	})
	if err != nil {
		logger.Error("failed to create subcategory", "error", err)
		http.Error(w, "Failed to create subcategory", http.StatusInternalServerError)
		return
	}

	// Get job for template context
	job, err := h.queries.GetJob(ctx, parent.JobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Return category partial for HTMX
	if r.Header.Get("HX-Request") == "true" {
		data := map[string]interface{}{
			"Category":        category,
			"Job":             job,
			"LineItems":       []repository.LineItem{},
			"ChildCategories": []repository.Category{},
			"Depth":           parentDepth + 1,
		}
		if err := h.renderer.RenderPartial(w, "category", data); err != nil {
			logger.Error("failed to render category", "error", err)
		}
		return
	}

	http.Redirect(w, r, "/jobs/"+parent.JobID, http.StatusSeeOther)
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
