package keyboard

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// ListItemTemplates shows the item templates management page with search and filters.
func (h *Handler) ListItemTemplates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	query := r.URL.Query().Get("q")
	typeFilter := r.URL.Query().Get("type")
	categoryFilter := r.URL.Query().Get("category")

	var items []repository.ItemTemplate
	var err error

	// Get all items for the categories dropdown and filtering
	allItems, err := h.queries.ListItemTemplates(ctx, orgID)
	if err != nil {
		logger.Error("failed to list item templates", "error", err)
		http.Error(w, "Failed to load item templates", http.StatusInternalServerError)
		return
	}

	// Extract unique categories and subcategories
	categorySet := make(map[string]bool)
	subcategorySet := make(map[string]bool)
	for _, item := range allItems {
		categorySet[item.Category] = true
		if item.Subcategory.Valid && item.Subcategory.String != "" {
			subcategorySet[item.Subcategory.String] = true
		}
	}
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	subcategories := make([]string, 0, len(subcategorySet))
	for sub := range subcategorySet {
		subcategories = append(subcategories, sub)
	}

	// Apply filters
	if query != "" || typeFilter != "" || categoryFilter != "" {
		items = filterItems(allItems, query, typeFilter, categoryFilter)
	} else {
		items = allItems
	}

	data := map[string]interface{}{
		"Items":          items,
		"Categories":     categories,
		"Subcategories":  subcategories,
		"Query":          query,
		"TypeFilter":     typeFilter,
		"CategoryFilter": categoryFilter,
	}

	// For HTMX partial requests, return just the items list
	if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "true" {
		var buf bytes.Buffer
		if err := h.renderer.RenderPartial(&buf, "item_templates_list", data); err != nil {
			logger.Error("failed to render item templates list", "error", err)
			http.Error(w, "Failed to render list", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
		return
	}

	if err := h.renderer.Render(w, "item_templates", data); err != nil {
		logger.Error("failed to render item templates page", "error", err)
	}
}

// filterItems filters items based on query, type, and category.
func filterItems(items []repository.ItemTemplate, query, typeFilter, categoryFilter string) []repository.ItemTemplate {
	var result []repository.ItemTemplate
	for _, item := range items {
		// Type filter
		if typeFilter != "" && item.Type != typeFilter {
			continue
		}
		// Category filter
		if categoryFilter != "" && item.Category != categoryFilter {
			continue
		}
		// Name search (case-insensitive contains)
		if query != "" {
			queryLower := stringToLower(query)
			nameLower := stringToLower(item.Name)
			if !stringContains(nameLower, queryLower) {
				continue
			}
		}
		result = append(result, item)
	}
	return result
}

// stringToLower converts a string to lowercase.
func stringToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		} else {
			b[i] = c
		}
	}
	return string(b)
}

// stringContains checks if haystack contains needle.
func stringContains(haystack, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	if len(needle) > len(haystack) {
		return false
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// GetItemTemplateForm returns the inline form for creating a new item template.
func (h *Handler) GetItemTemplateForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	// Get categories for autocomplete
	items, err := h.queries.ListItemTemplates(ctx, orgID)
	if err != nil {
		logger.Error("failed to list item templates", "error", err)
	}

	categorySet := make(map[string]bool)
	subcategorySet := make(map[string]bool)
	for _, item := range items {
		categorySet[item.Category] = true
		if item.Subcategory.Valid && item.Subcategory.String != "" {
			subcategorySet[item.Subcategory.String] = true
		}
	}
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	subcategories := make([]string, 0, len(subcategorySet))
	for sub := range subcategorySet {
		subcategories = append(subcategories, sub)
	}

	data := map[string]interface{}{
		"Categories":    categories,
		"Subcategories": subcategories,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "item_template_form", data); err != nil {
		logger.Error("failed to render item template form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// CreateItemTemplate creates a new item template.
func (h *Handler) CreateItemTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	itemType := r.FormValue("type")
	if itemType == "" {
		itemType = "material"
	}

	category := r.FormValue("category")
	if category == "" {
		category = "Uncategorized"
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	defaultUnit := r.FormValue("default_unit")
	if defaultUnit == "" {
		defaultUnit = "ea"
	}

	defaultPrice, _ := strconv.ParseFloat(r.FormValue("default_price"), 64)

	subcategory := r.FormValue("subcategory")
	if subcategory == "" {
		subcategory = name
	}

	_, err := h.queries.CreateItemTemplate(ctx, repository.CreateItemTemplateParams{
		OrgID:        orgID,
		Type:         itemType,
		Category:     category,
		Name:         name,
		DefaultUnit:  defaultUnit,
		DefaultPrice: defaultPrice,
		Subcategory:  sql.NullString{String: subcategory, Valid: true},
	})
	if err != nil {
		logger.Error("failed to create item template", "error", err)
		http.Error(w, "Failed to create item template", http.StatusInternalServerError)
		return
	}

	// Redirect back to the items page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/items")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

// GetItemTemplateEditForm returns the inline form for editing an item template.
func (h *Handler) GetItemTemplateEditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := h.queries.GetItemTemplate(ctx, repository.GetItemTemplateParams{
		ID:    id,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get item template", "error", err)
		http.Error(w, "Item template not found", http.StatusNotFound)
		return
	}

	// Get categories for autocomplete
	items, err := h.queries.ListItemTemplates(ctx, orgID)
	if err != nil {
		logger.Error("failed to list item templates", "error", err)
	}

	categorySet := make(map[string]bool)
	subcategorySet := make(map[string]bool)
	for _, it := range items {
		categorySet[it.Category] = true
		if it.Subcategory.Valid && it.Subcategory.String != "" {
			subcategorySet[it.Subcategory.String] = true
		}
	}
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	subcategories := make([]string, 0, len(subcategorySet))
	for sub := range subcategorySet {
		subcategories = append(subcategories, sub)
	}

	data := map[string]interface{}{
		"Item":          item,
		"Categories":    categories,
		"Subcategories": subcategories,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "item_template_edit_form", data); err != nil {
		logger.Error("failed to render item template edit form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateItemTemplate updates an existing item template.
func (h *Handler) UpdateItemTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	itemType := r.FormValue("type")
	if itemType == "" {
		itemType = "material"
	}

	category := r.FormValue("category")
	if category == "" {
		category = "Uncategorized"
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	defaultUnit := r.FormValue("default_unit")
	if defaultUnit == "" {
		defaultUnit = "ea"
	}

	defaultPrice, _ := strconv.ParseFloat(r.FormValue("default_price"), 64)

	subcategory := r.FormValue("subcategory")
	if subcategory == "" {
		subcategory = name
	}

	_, err = h.queries.UpdateItemTemplate(ctx, repository.UpdateItemTemplateParams{
		ID:           id,
		OrgID:        orgID,
		Type:         itemType,
		Category:     category,
		Name:         name,
		DefaultUnit:  defaultUnit,
		DefaultPrice: defaultPrice,
		Subcategory:  sql.NullString{String: subcategory, Valid: true},
	})
	if err != nil {
		logger.Error("failed to update item template", "error", err)
		http.Error(w, "Failed to update item template", http.StatusInternalServerError)
		return
	}

	// Redirect back to the items page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/items")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

// BulkUpdateFamily assigns a product family (subcategory) to multiple item templates.
func (h *Handler) BulkUpdateFamily(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	subcategory := r.FormValue("subcategory")
	if subcategory == "" {
		http.Error(w, "Product family is required", http.StatusBadRequest)
		return
	}

	idStrs := r.Form["ids[]"]
	if len(idStrs) == 0 {
		http.Error(w, "No items selected", http.StatusBadRequest)
		return
	}

	ids := make([]int64, 0, len(idStrs))
	for _, s := range idStrs {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	err := h.queries.BulkUpdateItemTemplateSubcategory(ctx, repository.BulkUpdateItemTemplateSubcategoryParams{
		Subcategory: sql.NullString{String: subcategory, Valid: true},
		OrgID:       orgID,
		Ids:         ids,
	})
	if err != nil {
		logger.Error("failed to bulk update item template families", "error", err)
		http.Error(w, "Failed to update families", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/items")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

// DeleteItemTemplate deletes an item template.
func (h *Handler) DeleteItemTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteItemTemplate(ctx, repository.DeleteItemTemplateParams{
		ID:    id,
		OrgID: orgID,
	}); err != nil {
		logger.Error("failed to delete item template", "error", err)
		http.Error(w, "Failed to delete item template", http.StatusInternalServerError)
		return
	}

	// Redirect back to the items page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/items")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}
