package keyboard

import (
	"bytes"
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// ListJobItemTypes returns the custom item types for a job.
func (h *Handler) ListJobItemTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	itemTypes, err := h.queries.ListJobItemTypes(ctx, jobID)
	if err != nil {
		logger.Error("failed to list job item types", "error", err)
		http.Error(w, "Failed to load item types", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Job":       job,
		"ItemTypes": itemTypes,
		"Colors":    domain.ValidColors,
	}

	if err := h.renderer.Render(w, "job_item_types", data); err != nil {
		logger.Error("failed to render job item types page", "error", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// GetJobItemTypeForm returns the form for creating a new custom item type.
func (h *Handler) GetJobItemTypeForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Job":    job,
		"Colors": domain.ValidColors,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "item_type_form", data); err != nil {
		logger.Error("failed to render item type form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// CreateJobItemType creates a new custom item type for a job.
func (h *Handler) CreateJobItemType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	slug := strings.TrimSpace(r.FormValue("slug"))
	color := strings.TrimSpace(r.FormValue("color"))
	sortOrderStr := r.FormValue("sort_order")

	// Auto-generate slug from name if not provided
	if slug == "" && name != "" {
		slug = generateSlug(name)
	}

	// Default color if not provided
	if color == "" {
		color = "slate"
	}

	sortOrder := 0
	if sortOrderStr != "" {
		if so, err := strconv.Atoi(sortOrderStr); err == nil {
			sortOrder = so
		}
	}

	// Parse surcharge percent
	var surchargePercent *float64
	var surchargeSql sql.NullFloat64
	if v := r.FormValue("surcharge_percent"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			surchargePercent = &f
			surchargeSql = sql.NullFloat64{Float64: f, Valid: true}
		}
	}

	// Validate input
	input := domain.JobItemTypeInput{
		JobID:            jobID,
		Name:             name,
		Slug:             slug,
		Color:            color,
		SortOrder:        sortOrder,
		SurchargePercent: surchargePercent,
	}

	if errors := input.Validate(); len(errors) > 0 {
		http.Error(w, errors[0].Message, http.StatusBadRequest)
		return
	}

	// Check if slug already exists for this job
	_, err := h.queries.GetJobItemTypeBySlug(ctx, repository.GetJobItemTypeBySlugParams{
		JobID: jobID,
		Slug:  slug,
	})
	if err == nil {
		http.Error(w, "A type with this slug already exists", http.StatusConflict)
		return
	}

	_, err = h.queries.CreateJobItemType(ctx, repository.CreateJobItemTypeParams{
		ID:               uuid.New().String(),
		JobID:            jobID,
		Name:             name,
		Slug:             slug,
		Color:            color,
		SortOrder:        int64(sortOrder),
		SurchargePercent: surchargeSql,
	})
	if err != nil {
		logger.Error("failed to create job item type", "error", err)
		http.Error(w, "Failed to create item type", http.StatusInternalServerError)
		return
	}

	redirectURL := "/jobs/" + jobID + "/item-types"

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// UpdateJobItemType updates an existing custom item type.
func (h *Handler) UpdateJobItemType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	typeID := r.PathValue("id")

	itemType, err := h.queries.GetJobItemType(ctx, typeID)
	if err != nil {
		logger.Error("failed to get job item type", "error", err)
		http.Error(w, "Item type not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	slug := strings.TrimSpace(r.FormValue("slug"))
	color := strings.TrimSpace(r.FormValue("color"))
	sortOrderStr := r.FormValue("sort_order")

	if name == "" {
		name = itemType.Name
	}
	if slug == "" {
		slug = itemType.Slug
	}
	if color == "" {
		color = itemType.Color
	}

	sortOrder := int(itemType.SortOrder)
	if sortOrderStr != "" {
		if so, err := strconv.Atoi(sortOrderStr); err == nil {
			sortOrder = so
		}
	}

	// Parse surcharge percent
	var surchargePercent *float64
	var surchargeSql sql.NullFloat64
	surchargeStr := r.FormValue("surcharge_percent")
	if surchargeStr != "" {
		if f, err := strconv.ParseFloat(surchargeStr, 64); err == nil {
			surchargePercent = &f
			surchargeSql = sql.NullFloat64{Float64: f, Valid: true}
		}
	}

	// Validate input
	input := domain.JobItemTypeInput{
		JobID:            itemType.JobID,
		Name:             name,
		Slug:             slug,
		Color:            color,
		SortOrder:        sortOrder,
		SurchargePercent: surchargePercent,
	}

	if errors := input.Validate(); len(errors) > 0 {
		http.Error(w, errors[0].Message, http.StatusBadRequest)
		return
	}

	// Check if slug changed and conflicts with existing
	if slug != itemType.Slug {
		existing, err := h.queries.GetJobItemTypeBySlug(ctx, repository.GetJobItemTypeBySlugParams{
			JobID: itemType.JobID,
			Slug:  slug,
		})
		if err == nil && existing.ID != typeID {
			http.Error(w, "A type with this slug already exists", http.StatusConflict)
			return
		}
	}

	_, err = h.queries.UpdateJobItemType(ctx, repository.UpdateJobItemTypeParams{
		ID:               typeID,
		Name:             name,
		Slug:             slug,
		Color:            color,
		SortOrder:        int64(sortOrder),
		SurchargePercent: surchargeSql,
	})
	if err != nil {
		logger.Error("failed to update job item type", "error", err)
		http.Error(w, "Failed to update item type", http.StatusInternalServerError)
		return
	}

	redirectURL := "/jobs/" + itemType.JobID + "/item-types"

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// DeleteJobItemType deletes a custom item type.
func (h *Handler) DeleteJobItemType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	typeID := r.PathValue("id")

	itemType, err := h.queries.GetJobItemType(ctx, typeID)
	if err != nil {
		logger.Error("failed to get job item type", "error", err)
		http.Error(w, "Item type not found", http.StatusNotFound)
		return
	}

	// Check if any line items use this type
	count, err := h.queries.CountLineItemsByType(ctx, repository.CountLineItemsByTypeParams{
		JobID: itemType.JobID,
		Type:  itemType.Slug,
	})
	if err != nil {
		logger.Error("failed to count line items by type", "error", err)
		http.Error(w, "Failed to check type usage", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Cannot delete item type that is in use", http.StatusConflict)
		return
	}

	if err := h.queries.DeleteJobItemType(ctx, typeID); err != nil {
		logger.Error("failed to delete job item type", "error", err)
		http.Error(w, "Failed to delete item type", http.StatusInternalServerError)
		return
	}

	redirectURL := "/jobs/" + itemType.JobID + "/item-types"

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// generateSlug converts a name to a URL-safe slug.
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	return slug
}
