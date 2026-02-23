package keyboard

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// TagGroup represents a group of items with the same tag within a type group.
type TagGroup struct {
	Tag      string // Empty string for untagged items
	Items    []repository.LineItem
	Subtotal float64
}

// ItemTypeGroup represents a group of line items of the same type.
type ItemTypeGroup struct {
	Slug      string
	Name      string
	Color     string
	Hotkey    string
	Items     []repository.LineItem
	TagGroups []TagGroup
	Subtotal  float64
}

// groupItemsByType organizes line items into groups by their type.
// Standard types (material, labor, equipment) are always included.
// Custom types are included if they exist for the job.
func groupItemsByType(items []repository.LineItem, customTypes []repository.JobItemType) []ItemTypeGroup {
	// Initialize standard type groups (always shown)
	groups := []ItemTypeGroup{
		{Slug: "material", Name: "Materials", Color: "forest", Hotkey: "m", Items: []repository.LineItem{}, TagGroups: []TagGroup{}, Subtotal: 0},
		{Slug: "labor", Name: "Labor", Color: "copper", Hotkey: "l", Items: []repository.LineItem{}, TagGroups: []TagGroup{}, Subtotal: 0},
		{Slug: "equipment", Name: "Equipment", Color: "slate", Hotkey: "e", Items: []repository.LineItem{}, TagGroups: []TagGroup{}, Subtotal: 0},
	}

	// Add custom type groups
	for i, ct := range customTypes {
		groups = append(groups, ItemTypeGroup{
			Slug:      ct.Slug,
			Name:      ct.Name,
			Color:     ct.Color,
			Hotkey:    fmt.Sprintf("%d", i+1),
			Items:     []repository.LineItem{},
			TagGroups: []TagGroup{},
			Subtotal:  0,
		})
	}

	// Create a map for quick lookup
	groupIndex := make(map[string]int)
	for i, g := range groups {
		groupIndex[g.Slug] = i
	}

	// Distribute items into groups
	for _, item := range items {
		if idx, ok := groupIndex[item.Type]; ok {
			groups[idx].Items = append(groups[idx].Items, item)
			groups[idx].Subtotal += item.Quantity * item.UnitPrice
		}
	}

	// Organize each group's items into tag subgroups
	for i := range groups {
		groups[i].TagGroups = organizeByTag(groups[i].Items)
	}

	return groups
}

// organizeByTag groups items by their tag, with untagged items first.
func organizeByTag(items []repository.LineItem) []TagGroup {
	if len(items) == 0 {
		return nil
	}

	// Map to collect items by tag
	tagMap := make(map[string][]repository.LineItem)
	var tags []string

	for _, item := range items {
		tag := ""
		if item.Tag.Valid {
			tag = item.Tag.String
		}

		if _, exists := tagMap[tag]; !exists {
			tags = append(tags, tag)
		}
		tagMap[tag] = append(tagMap[tag], item)
	}

	// Sort tags: empty string (untagged) first, then alphabetically
	sort.Slice(tags, func(i, j int) bool {
		if tags[i] == "" {
			return true
		}
		if tags[j] == "" {
			return false
		}
		return tags[i] < tags[j]
	})

	// Build tag groups
	var tagGroups []TagGroup
	for _, tag := range tags {
		tagItems := tagMap[tag]
		subtotal := 0.0
		for _, item := range tagItems {
			subtotal += item.Quantity * item.UnitPrice
		}
		tagGroups = append(tagGroups, TagGroup{
			Tag:      tag,
			Items:    tagItems,
			Subtotal: subtotal,
		})
	}

	return tagGroups
}

// GetCategoryMarkupForm returns an inline form for editing category markup.
func (h *Handler) GetCategoryMarkupForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	category, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    categoryID,
		OrgID: GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Category": category,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "category_markup_form", data); err != nil {
		logger.Error("failed to render markup form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// GetCategoryRenameForm returns an inline form for renaming a category.
func (h *Handler) GetCategoryRenameForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	category, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    categoryID,
		OrgID: GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Category": category,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "category_rename_form", data); err != nil {
		logger.Error("failed to render rename form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateCategoryName updates only a category's name.
func (h *Handler) UpdateCategoryName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	category, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    categoryID,
		OrgID: GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = category.Name
	}

	_, err = h.queries.UpdateCategory(ctx, repository.UpdateCategoryParams{
		ID:               categoryID,
		Name:             name,
		SurchargePercent: category.SurchargePercent,
		SortOrder:        category.SortOrder,
	})
	if err != nil {
		logger.Error("failed to update category name", "error", err)
		http.Error(w, "Failed to update name", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/categories/"+categoryID)
		return
	}

	http.Redirect(w, r, "/categories/"+categoryID, http.StatusSeeOther)
}

// UpdateCategoryMarkup updates a category's markup percentage.
func (h *Handler) UpdateCategoryMarkup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	category, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    categoryID,
		OrgID: GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	surchargeStr := r.FormValue("surcharge_percent")
	var surchargePercent sql.NullFloat64
	if surchargeStr != "" {
		val, _ := strconv.ParseFloat(surchargeStr, 64)
		surchargePercent = sql.NullFloat64{Float64: val, Valid: true}
	}

	_, err = h.queries.UpdateCategory(ctx, repository.UpdateCategoryParams{
		ID:               categoryID,
		Name:             category.Name,
		SurchargePercent: surchargePercent,
		SortOrder:        category.SortOrder,
	})
	if err != nil {
		logger.Error("failed to update category markup", "error", err)
		http.Error(w, "Failed to update markup", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/categories/"+categoryID)
		return
	}

	http.Redirect(w, r, "/categories/"+categoryID, http.StatusSeeOther)
}

// GetEditForm returns an inline form for editing an existing line item.
func (h *Handler) GetEditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")

	item, err := h.queries.GetLineItem(ctx, repository.GetLineItemParams{
		ID:    itemID,
		OrgID: GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to get line item", "error", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Fetch existing items to get tags for autocomplete
	items, err := h.queries.ListLineItemsByCategory(ctx, repository.ListLineItemsByCategoryParams{
		CategoryID: item.CategoryID,
		OrgID:      GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to list items for tags", "error", err)
		items = []repository.LineItem{} // Continue with empty list
	}

	// Filter tags for this type only
	existingTags := extractUniqueTagsForType(items, item.Type)

	data := map[string]interface{}{
		"Item":         item,
		"ExistingTags": existingTags,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "edit_form", data); err != nil {
		logger.Error("failed to render edit form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateLineItem updates an existing line item.
func (h *Handler) UpdateLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")

	item, err := h.queries.GetLineItem(ctx, repository.GetLineItemParams{
		ID:    itemID,
		OrgID: GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to get line item", "error", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

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
		name = item.Name
	}

	unit := r.FormValue("unit")
	if unit == "" {
		unit = item.Unit
	}

	tag := r.FormValue("tag")
	var tagParam sql.NullString
	if tag != "" {
		tagParam = sql.NullString{String: tag, Valid: true}
	}

	description := r.FormValue("description")
	var descParam sql.NullString
	if description != "" {
		descParam = sql.NullString{String: description, Valid: true}
	}

	_, err = h.queries.UpdateLineItem(ctx, repository.UpdateLineItemParams{
		ID:               itemID,
		Type:             item.Type,
		Name:             name,
		Description:      descParam,
		Quantity:         quantity,
		Unit:             unit,
		UnitPrice:        unitPrice,
		SurchargePercent: item.SurchargePercent,
		SortOrder:        item.SortOrder,
		Tag:              tagParam,
	})
	if err != nil {
		logger.Error("failed to update line item", "error", err)
		http.Error(w, "Failed to update line item", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/categories/"+item.CategoryID)
		return
	}

	http.Redirect(w, r, "/categories/"+item.CategoryID, http.StatusSeeOther)
}

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
		OrgID:   GetOrgID(ctx),
		Type:    itemType,
		Column3: sql.NullString{String: query, Valid: true},
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
	_, _ = w.Write(buf.Bytes())
}

// GetCategory shows a category with its items and subcategories.
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")

	orgID := GetOrgID(ctx)
	category, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    categoryID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Failed to load category", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    category.JobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	categories, err := h.queries.ListCategoriesByJob(ctx, repository.ListCategoriesByJobParams{
		JobID: job.ID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, repository.ListLineItemsByJobParams{
		JobID: job.ID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load line items", http.StatusInternalServerError)
		return
	}

	// Get custom item types for this job (needed for surcharge calculations)
	customTypes, err := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{
		JobID: job.ID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list job item types", "error", err)
		customTypes = []repository.JobItemType{} // Continue with empty list
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
	catTotal := h.calculateCategoryTotal(categoryID, job, categories, lineItems, customTypes)

	// Calculate totals for subcategories
	type SubcategoryWithTotal struct {
		repository.Category
		Total float64
	}
	subcatsWithTotals := make([]SubcategoryWithTotal, len(subcategories))
	for i, sub := range subcategories {
		subTotal := h.calculateCategoryTotal(sub.ID, job, categories, lineItems, customTypes)
		subcatsWithTotals[i] = SubcategoryWithTotal{
			Category: sub,
			Total:    subTotal.Total,
		}
	}

	// Build category tree for sidebar navigation
	categoryTree := buildCategoryTree(categories)

	// Group items by type for sectioned display
	itemsByType := groupItemsByType(categoryItems, customTypes)

	data := map[string]interface{}{
		"Job":               job,
		"Category":          category,
		"Subcategories":     subcatsWithTotals,
		"Items":             categoryItems,
		"ItemsByType":       itemsByType,
		"Breadcrumbs":       breadcrumbs,
		"Depth":             depth,
		"CanAddSubcategory": canAddSubcategory(depth),
		"CategoryTotal":     catTotal,
		"SelectedIndex":     0,
		"CategoryTree":      categoryTree,
		"CurrentCategoryID": categoryID,
		"CustomTypes":       customTypes,
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
		w.Header().Set("HX-Redirect", "/categories/"+category.ID)
		return
	}

	http.Redirect(w, r, "/categories/"+category.ID, http.StatusSeeOther)
}

// CreateSubcategory creates a subcategory under a parent.
func (h *Handler) CreateSubcategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	parentID := r.PathValue("parentID")
	orgID := GetOrgID(ctx)

	parent, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    parentID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get parent category", "error", err)
		http.Error(w, "Parent category not found", http.StatusNotFound)
		return
	}

	// Check depth
	categories, _ := h.queries.ListCategoriesByJob(ctx, repository.ListCategoriesByJobParams{
		JobID: parent.JobID,
		OrgID: orgID,
	})
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
		w.Header().Set("HX-Redirect", "/categories/"+category.ID)
		return
	}

	http.Redirect(w, r, "/categories/"+category.ID, http.StatusSeeOther)
}

// DeleteCategory deletes a category.
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	category, err := h.queries.GetCategory(ctx, repository.GetCategoryParams{
		ID:    categoryID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	redirectURL := "/jobs/" + category.JobID
	if category.ParentID.Valid {
		redirectURL = "/categories/" + category.ParentID.String
	}

	if err := h.queries.DeleteCategory(ctx, repository.DeleteCategoryParams{
		ID:    categoryID,
		OrgID: orgID,
	}); err != nil {
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

	tag := r.FormValue("tag")
	var tagParam sql.NullString
	if tag != "" {
		tagParam = sql.NullString{String: tag, Valid: true}
	}

	description := r.FormValue("description")
	var descParam sql.NullString
	if description != "" {
		descParam = sql.NullString{String: description, Valid: true}
	}

	_, err := h.queries.CreateLineItem(ctx, repository.CreateLineItemParams{
		ID:               uuid.New().String(),
		CategoryID:       categoryID,
		Type:             itemType,
		Name:             name,
		Description:      descParam,
		Quantity:         quantity,
		Unit:             unit,
		UnitPrice:        unitPrice,
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
		Tag:              tagParam,
	})
	if err != nil {
		logger.Error("failed to create line item", "error", err)
		http.Error(w, "Failed to create line item", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/categories/"+categoryID)
		return
	}

	http.Redirect(w, r, "/categories/"+categoryID, http.StatusSeeOther)
}

// DeleteLineItem deletes a line item.
func (h *Handler) DeleteLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	item, err := h.queries.GetLineItem(ctx, repository.GetLineItemParams{
		ID:    itemID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get line item", "error", err)
		http.Error(w, "Line item not found", http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteLineItem(ctx, repository.DeleteLineItemParams{
		ID:    itemID,
		OrgID: orgID,
	}); err != nil {
		logger.Error("failed to delete line item", "error", err)
		http.Error(w, "Failed to delete line item", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/categories/"+item.CategoryID)
		return
	}

	http.Redirect(w, r, "/categories/"+item.CategoryID, http.StatusSeeOther)
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

	// Fetch existing items to get tags for autocomplete
	items, err := h.queries.ListLineItemsByCategory(ctx, repository.ListLineItemsByCategoryParams{
		CategoryID: categoryID,
		OrgID:      GetOrgID(ctx),
	})
	if err != nil {
		logger.Error("failed to list items for tags", "error", err)
		items = []repository.LineItem{} // Continue with empty list
	}

	// Filter tags for this type only (tags are scoped within types)
	existingTags := extractUniqueTagsForType(items, itemType)

	data := map[string]interface{}{
		"CategoryID":   categoryID,
		"Type":         itemType,
		"DefaultUnit":  defaultUnit,
		"ExistingTags": existingTags,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "inline_form", data); err != nil {
		logger.Error("failed to render inline form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// extractUniqueTagsForType returns unique tags for a specific item type.
func extractUniqueTagsForType(items []repository.LineItem, itemType string) []string {
	tagSet := make(map[string]bool)
	for _, item := range items {
		if item.Type == itemType && item.Tag.Valid && item.Tag.String != "" {
			tagSet[item.Tag.String] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// GetBatchForm returns the batch entry form for adding multiple line items from templates.
func (h *Handler) GetBatchForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("categoryID")
	orgID := GetOrgID(ctx)

	templateCategory := r.URL.Query().Get("template_category")

	// Always fetch template categories for the picker
	categories, err := h.queries.ListItemTemplateCategories(ctx, orgID)
	if err != nil {
		logger.Error("failed to list template categories", "error", err)
		categories = []string{}
	}

	// Fetch existing tags for the datalist
	items, err := h.queries.ListLineItemsByCategory(ctx, repository.ListLineItemsByCategoryParams{
		CategoryID: categoryID,
		OrgID:      orgID,
	})
	if err != nil {
		logger.Error("failed to list items for tags", "error", err)
		items = []repository.LineItem{}
	}
	existingTags := extractUniqueTagsForType(items, "material")

	data := map[string]interface{}{
		"CategoryID":         categoryID,
		"TemplateCategories": categories,
		"SelectedCategory":   templateCategory,
		"ExistingTags":       existingTags,
	}

	// If a template category is selected, load the templates
	if templateCategory != "" {
		templates, err := h.queries.ListItemTemplatesByCategory(ctx, repository.ListItemTemplatesByCategoryParams{
			OrgID:    orgID,
			Category: templateCategory,
		})
		if err != nil {
			logger.Error("failed to list templates by category", "error", err)
			http.Error(w, "Failed to load templates", http.StatusInternalServerError)
			return
		}
		data["Templates"] = templates
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "batch_form", data); err != nil {
		logger.Error("failed to render batch form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// BatchCreateLineItems creates multiple line items from a batch form submission.
func (h *Handler) BatchCreateLineItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("categoryID")
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	names := r.Form["name[]"]
	units := r.Form["unit[]"]
	unitPrices := r.Form["unit_price[]"]
	quantities := r.Form["quantity[]"]
	types := r.Form["type[]"]
	tag := r.FormValue("tag")

	var tagParam sql.NullString
	if tag != "" {
		tagParam = sql.NullString{String: tag, Valid: true}
	}

	// Start transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("failed to begin transaction", "error", err)
		http.Error(w, "Failed to create items", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := h.queries.WithTx(tx)
	created := 0

	for i := range names {
		if i >= len(quantities) || i >= len(unitPrices) || i >= len(units) {
			break
		}

		qty, _ := strconv.ParseFloat(quantities[i], 64)
		if qty <= 0 {
			continue
		}

		unitPrice, _ := strconv.ParseFloat(unitPrices[i], 64)

		itemType := "material"
		if i < len(types) && types[i] != "" {
			itemType = types[i]
		}

		_, err := qtx.CreateLineItem(ctx, repository.CreateLineItemParams{
			ID:               uuid.New().String(),
			OrgID:            orgID,
			CategoryID:       categoryID,
			Type:             itemType,
			Name:             names[i],
			Quantity:         qty,
			Unit:             units[i],
			UnitPrice:        unitPrice,
			SurchargePercent: sql.NullFloat64{},
			SortOrder:        0,
			Tag:              tagParam,
		})
		if err != nil {
			logger.Error("failed to create line item in batch", "error", err, "name", names[i])
			http.Error(w, "Failed to create items", http.StatusInternalServerError)
			return
		}
		created++
	}

	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit batch transaction", "error", err)
		http.Error(w, "Failed to create items", http.StatusInternalServerError)
		return
	}

	logger.Info("batch created line items", "count", created, "categoryID", categoryID)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/categories/"+categoryID)
		return
	}

	http.Redirect(w, r, "/categories/"+categoryID, http.StatusSeeOther)
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
		action = "/categories/" + parentID + "/subcategories"
	} else if jobID != "" {
		action = "/jobs/" + jobID + "/categories"
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
	_, _ = w.Write(buf.Bytes())
}
