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

// SearchItems searches for item templates by type and name.
func (h *Handler) SearchItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	itemType := r.URL.Query().Get("type")
	query := r.URL.Query().Get("q")

	if query == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return
	}

	items, err := h.queries.SearchItemTemplatesByType(ctx, repository.SearchItemTemplatesByTypeParams{
		Type:   itemType,
		Column2: sql.NullString{String: query, Valid: true},
	})
	if err != nil {
		logger.Error("failed to search items", "error", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "search_results", items); err != nil {
		logger.Error("failed to render search results", "error", err)
		http.Error(w, "Failed to render results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

// GetCategory shows a category with its items and subcategories.
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	category, err := h.queries.GetCategory(ctx, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Failed to load category", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.GetJob(ctx, category.JobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	categories, err := h.queries.ListCategoriesByJob(ctx, job.ID)
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, job.ID)
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load line items", http.StatusInternalServerError)
		return
	}

	// Get direct children (subcategories)
	subcategories := make([]repository.Category, 0)
	for _, cat := range categories {
		if cat.ParentID.Valid && cat.ParentID.String == categoryID {
			subcategories = append(subcategories, cat)
		}
	}

	// Get line items for this category only
	categoryItems := make([]repository.LineItem, 0)
	for _, item := range lineItems {
		if item.CategoryID == categoryID {
			categoryItems = append(categoryItems, item)
		}
	}

	// Calculate depth and breadcrumbs
	depth := h.getCategoryDepth(categories, categoryID)
	breadcrumbs := h.getBreadcrumbs(categories, categoryID, job)

	// Calculate category total
	catTotal := h.calculateCategoryTotal(categoryID, job, categories, lineItems)

	// Calculate totals for subcategories
	type SubcategoryWithTotal struct {
		repository.Category
		Total float64
	}
	subcatsWithTotals := make([]SubcategoryWithTotal, len(subcategories))
	for i, sub := range subcategories {
		subTotal := h.calculateCategoryTotal(sub.ID, job, categories, lineItems)
		subcatsWithTotals[i] = SubcategoryWithTotal{
			Category: sub,
			Total:    subTotal.Total,
		}
	}

	data := map[string]interface{}{
		"Job":              job,
		"Category":         category,
		"Subcategories":    subcatsWithTotals,
		"Items":            categoryItems,
		"Breadcrumbs":      breadcrumbs,
		"Depth":            depth,
		"CanAddSubcategory": canAddSubcategory(depth),
		"CategoryTotal":    catTotal,
		"SelectedIndex":    0,
	}

	if err := h.renderer.Render(w, "category", data); err != nil {
		logger.Error("failed to render category page", "error", err)
	}
}

// CreateCategory creates a new top-level category.
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "New Category"
	}

	category, err := h.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:               uuid.New().String(),
		JobID:            jobID,
		ParentID:         sql.NullString{},
		Name:             name,
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
	})
	if err != nil {
		logger.Error("failed to create category", "error", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/k/categories/"+category.ID)
		return
	}

	http.Redirect(w, r, "/k/categories/"+category.ID, http.StatusSeeOther)
}

// CreateSubcategory creates a subcategory under a parent.
func (h *Handler) CreateSubcategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	parentID := r.PathValue("parentID")

	parent, err := h.queries.GetCategory(ctx, parentID)
	if err != nil {
		logger.Error("failed to get parent category", "error", err)
		http.Error(w, "Parent category not found", http.StatusNotFound)
		return
	}

	// Check depth
	categories, _ := h.queries.ListCategoriesByJob(ctx, parent.JobID)
	depth := h.getCategoryDepth(categories, parentID)
	if depth >= 3 {
		http.Error(w, "Maximum category depth reached", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "New Subcategory"
	}

	category, err := h.queries.CreateCategory(ctx, repository.CreateCategoryParams{
		ID:               uuid.New().String(),
		JobID:            parent.JobID,
		ParentID:         sql.NullString{String: parentID, Valid: true},
		Name:             name,
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
	})
	if err != nil {
		logger.Error("failed to create subcategory", "error", err)
		http.Error(w, "Failed to create subcategory", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/k/categories/"+category.ID)
		return
	}

	http.Redirect(w, r, "/k/categories/"+category.ID, http.StatusSeeOther)
}

// DeleteCategory deletes a category.
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	category, err := h.queries.GetCategory(ctx, categoryID)
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	redirectURL := "/k/jobs/" + category.JobID
	if category.ParentID.Valid {
		redirectURL = "/k/categories/" + category.ParentID.String
	}

	if err := h.queries.DeleteCategory(ctx, categoryID); err != nil {
		logger.Error("failed to delete category", "error", err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// CreateLineItem creates a new line item.
func (h *Handler) CreateLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("categoryID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
	if quantity <= 0 {
		quantity = 1
	}

	unitPrice, _ := strconv.ParseFloat(r.FormValue("unit_price"), 64)

	name := r.FormValue("name")
	if name == "" {
		name = "New Item"
	}

	unit := r.FormValue("unit")
	if unit == "" {
		unit = "ea"
	}

	itemType := r.FormValue("type")
	if itemType == "" {
		itemType = "material"
	}

	_, err := h.queries.CreateLineItem(ctx, repository.CreateLineItemParams{
		ID:               uuid.New().String(),
		CategoryID:       categoryID,
		Type:             itemType,
		Name:             name,
		Description:      sql.NullString{},
		Quantity:         quantity,
		Unit:             unit,
		UnitPrice:        unitPrice,
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
	})
	if err != nil {
		logger.Error("failed to create line item", "error", err)
		http.Error(w, "Failed to create line item", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/k/categories/"+categoryID)
		return
	}

	http.Redirect(w, r, "/k/categories/"+categoryID, http.StatusSeeOther)
}

// DeleteLineItem deletes a line item.
func (h *Handler) DeleteLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")

	item, err := h.queries.GetLineItem(ctx, itemID)
	if err != nil {
		logger.Error("failed to get line item", "error", err)
		http.Error(w, "Line item not found", http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteLineItem(ctx, itemID); err != nil {
		logger.Error("failed to delete line item", "error", err)
		http.Error(w, "Failed to delete line item", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/k/categories/"+item.CategoryID)
		return
	}

	http.Redirect(w, r, "/k/categories/"+item.CategoryID, http.StatusSeeOther)
}

// GetInlineForm returns an inline form for creating items.
func (h *Handler) GetInlineForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("categoryID")
	itemType := r.URL.Query().Get("type")

	if itemType == "" {
		itemType = "material"
	}

	// Default units based on type
	defaultUnit := "ea"
	switch itemType {
	case "labor":
		defaultUnit = "hr"
	case "equipment":
		defaultUnit = "day"
	}

	data := map[string]interface{}{
		"CategoryID":  categoryID,
		"Type":        itemType,
		"DefaultUnit": defaultUnit,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "inline_form", data); err != nil {
		logger.Error("failed to render inline form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

// GetCategoryForm returns an inline form for creating categories.
func (h *Handler) GetCategoryForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	// Determine if this is for a job (top-level) or category (subcategory)
	jobID := r.URL.Query().Get("job_id")
	parentID := r.URL.Query().Get("parent_id")

	var action string
	if parentID != "" {
		action = "/k/categories/" + parentID + "/subcategories"
	} else if jobID != "" {
		action = "/k/jobs/" + jobID + "/categories"
	} else {
		http.Error(w, "Missing job_id or parent_id", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Action": action,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "category_form", data); err != nil {
		logger.Error("failed to render category form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}
