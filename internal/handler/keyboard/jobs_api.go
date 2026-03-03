package keyboard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// --- Wire types: clean JSON, no sql.Null* wrappers ---

type QuoteResponse struct {
	Job         QuoteJob        `json:"job"`
	Categories  []QuoteCategory `json:"categories"`
	CustomTypes []QuoteItemType `json:"custom_types"`
	Totals      QuoteTotals     `json:"totals"`
}

type QuoteJob struct {
	ID                        string   `json:"id"`
	Name                      string   `json:"name"`
	Status                    string   `json:"status"`
	SurchargePercent          float64  `json:"surcharge_percent"`
	SurchargeMode             string   `json:"surcharge_mode"`
	MaterialSurchargePercent  *float64 `json:"material_surcharge_percent"`
	LaborSurchargePercent     *float64 `json:"labor_surcharge_percent"`
	EquipmentSurchargePercent *float64 `json:"equipment_surcharge_percent"`
}

type QuoteCategory struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	ParentID         *string         `json:"parent_id"`
	SortOrder        int64           `json:"sort_order"`
	SurchargePercent *float64        `json:"surcharge_percent"`
	Items            []QuoteLineItem `json:"items"`
	Children         []QuoteCategory `json:"children"`
}

type QuoteLineItem struct {
	ID               string   `json:"id"`
	CategoryID       string   `json:"category_id"`
	Type             string   `json:"type"`
	Name             string   `json:"name"`
	Description      *string  `json:"description"`
	Quantity         float64  `json:"quantity"`
	Unit             string   `json:"unit"`
	UnitPrice        float64  `json:"unit_price"`
	SurchargePercent *float64 `json:"surcharge_percent"`
	SortOrder        int64    `json:"sort_order"`
	Tag              *string  `json:"tag"`
}

type QuoteItemType struct {
	Slug             string   `json:"slug"`
	Name             string   `json:"name"`
	Color            string   `json:"color"`
	SurchargePercent *float64 `json:"surcharge_percent"`
}

type QuoteTotals struct {
	Subtotal      float64            `json:"subtotal"`
	SurchargeTotal float64           `json:"surcharge_total"`
	GrandTotal    float64            `json:"grand_total"`
	TypeSubtotals map[string]float64 `json:"type_subtotals"`
}

// --- Conversion helpers ---

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullFloat64ToPtr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

func ptrToNullString(p *string) sql.NullString {
	if p != nil {
		return sql.NullString{String: *p, Valid: true}
	}
	return sql.NullString{}
}

func ptrToNullFloat64(p *float64) sql.NullFloat64 {
	if p != nil {
		return sql.NullFloat64{Float64: *p, Valid: true}
	}
	return sql.NullFloat64{}
}

// convertLineItem converts a repository LineItem to a QuoteLineItem.
func convertLineItem(item repository.LineItem) QuoteLineItem {
	return QuoteLineItem{
		ID:               item.ID,
		CategoryID:       item.CategoryID,
		Type:             item.Type,
		Name:             item.Name,
		Description:      nullStringToPtr(item.Description),
		Quantity:         item.Quantity,
		Unit:             item.Unit,
		UnitPrice:        item.UnitPrice,
		SurchargePercent: nullFloat64ToPtr(item.SurchargePercent),
		SortOrder:        item.SortOrder,
		Tag:              nullStringToPtr(item.Tag),
	}
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeJSONError writes a JSON error response.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// --- Handlers ---

// GetJobJSON returns the full job data as JSON for the Svelte quote editor.
func (h *Handler) GetJobJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "Job not found")
			return
		}
		logger.Error("failed to get job", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load job")
		return
	}

	categories, err := h.queries.ListCategoriesByJob(ctx, repository.ListCategoriesByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load categories")
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, repository.ListLineItemsByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load line items")
		return
	}

	customTypes, err := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list job item types", "error", err)
		customTypes = []repository.JobItemType{}
	}

	totals := h.calculateTotals(job, categories, lineItems, customTypes)

	// Build nested category tree with embedded items
	categoryTree := buildQuoteCategoryTree(categories, lineItems)

	// Convert custom types
	quoteCustomTypes := make([]QuoteItemType, len(customTypes))
	for i, ct := range customTypes {
		quoteCustomTypes[i] = QuoteItemType{
			Slug:             ct.Slug,
			Name:             ct.Name,
			Color:            ct.Color,
			SurchargePercent: nullFloat64ToPtr(ct.SurchargePercent),
		}
	}

	resp := QuoteResponse{
		Job: QuoteJob{
			ID:                        job.ID,
			Name:                      job.Name,
			Status:                    job.Status,
			SurchargePercent:          job.SurchargePercent,
			SurchargeMode:             job.SurchargeMode,
			MaterialSurchargePercent:  nullFloat64ToPtr(job.MaterialSurchargePercent),
			LaborSurchargePercent:     nullFloat64ToPtr(job.LaborSurchargePercent),
			EquipmentSurchargePercent: nullFloat64ToPtr(job.EquipmentSurchargePercent),
		},
		Categories:  categoryTree,
		CustomTypes: quoteCustomTypes,
		Totals: QuoteTotals{
			Subtotal:       totals.Subtotal,
			SurchargeTotal: totals.SurchargeTotal,
			GrandTotal:     totals.GrandTotal,
			TypeSubtotals:  totals.TypeSubtotals,
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

// buildQuoteCategoryTree builds a nested category tree with embedded items.
func buildQuoteCategoryTree(categories []repository.Category, lineItems []repository.LineItem) []QuoteCategory {
	// Group items by category ID
	itemsByCategory := make(map[string][]repository.LineItem)
	for _, item := range lineItems {
		itemsByCategory[item.CategoryID] = append(itemsByCategory[item.CategoryID], item)
	}

	// Group children by parent ID
	childrenByParent := make(map[string][]repository.Category)
	for _, cat := range categories {
		if cat.ParentID.Valid {
			childrenByParent[cat.ParentID.String] = append(childrenByParent[cat.ParentID.String], cat)
		}
	}

	var buildNode func(cat repository.Category) QuoteCategory
	buildNode = func(cat repository.Category) QuoteCategory {
		// Convert items
		rawItems := itemsByCategory[cat.ID]
		items := make([]QuoteLineItem, len(rawItems))
		for i, item := range rawItems {
			items[i] = convertLineItem(item)
		}

		// Build children recursively
		children := make([]QuoteCategory, 0)
		for _, child := range childrenByParent[cat.ID] {
			children = append(children, buildNode(child))
		}

		return QuoteCategory{
			ID:               cat.ID,
			Name:             cat.Name,
			ParentID:         nullStringToPtr(cat.ParentID),
			SortOrder:        cat.SortOrder,
			SurchargePercent: nullFloat64ToPtr(cat.SurchargePercent),
			Items:            items,
			Children:         children,
		}
	}

	// Start from root categories (no parent)
	tree := make([]QuoteCategory, 0)
	for _, cat := range categories {
		if !cat.ParentID.Valid {
			tree = append(tree, buildNode(cat))
		}
	}

	return tree
}

// loadJobTotals is a helper to recalculate totals for a job and return a QuoteTotals.
func (h *Handler) loadJobTotals(r *http.Request, jobID string, orgID uuid.NullUUID) (QuoteTotals, error) {
	ctx := r.Context()
	job, err := h.queries.GetJob(ctx, repository.GetJobParams{ID: jobID, OrgID: orgID})
	if err != nil {
		return QuoteTotals{}, fmt.Errorf("get job: %w", err)
	}
	categories, err := h.queries.ListCategoriesByJob(ctx, repository.ListCategoriesByJobParams{JobID: jobID, OrgID: orgID})
	if err != nil {
		return QuoteTotals{}, fmt.Errorf("list categories: %w", err)
	}
	lineItems, err := h.queries.ListLineItemsByJob(ctx, repository.ListLineItemsByJobParams{JobID: jobID, OrgID: orgID})
	if err != nil {
		return QuoteTotals{}, fmt.Errorf("list line items: %w", err)
	}
	customTypes, _ := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{JobID: jobID, OrgID: orgID})
	totals := h.calculateTotals(job, categories, lineItems, customTypes)
	return QuoteTotals{
		Subtotal:       totals.Subtotal,
		SurchargeTotal: totals.SurchargeTotal,
		GrandTotal:     totals.GrandTotal,
		TypeSubtotals:  totals.TypeSubtotals,
	}, nil
}

// categoryBelongsToJob checks that a category is part of the specified job.
func (h *Handler) categoryBelongsToJob(r *http.Request, categoryID, jobID string, orgID uuid.NullUUID) bool {
	cat, err := h.queries.GetCategory(r.Context(), repository.GetCategoryParams{ID: categoryID, OrgID: orgID})
	if err != nil {
		return false
	}
	return cat.JobID == jobID
}

// PatchLineItemJSON updates fields on an existing line item (partial update).
func (h *Handler) PatchLineItemJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")
	itemID := r.PathValue("itemId")

	// Parse partial JSON body
	var fields map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Fetch existing item
	item, err := h.queries.GetLineItem(ctx, repository.GetLineItemParams{ID: itemID, OrgID: orgID})
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "Item not found")
			return
		}
		logger.Error("failed to get line item", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load item")
		return
	}

	// Verify item's category belongs to this job
	if !h.categoryBelongsToJob(r, item.CategoryID, jobID, orgID) {
		writeJSONError(w, http.StatusForbidden, "Item does not belong to this job")
		return
	}

	// Merge incoming fields onto existing values
	if raw, ok := fields["name"]; ok {
		var v string
		if json.Unmarshal(raw, &v) == nil {
			item.Name = v
		}
	}
	if raw, ok := fields["quantity"]; ok {
		var v float64
		if json.Unmarshal(raw, &v) == nil {
			item.Quantity = v
		}
	}
	if raw, ok := fields["unit"]; ok {
		var v string
		if json.Unmarshal(raw, &v) == nil {
			item.Unit = v
		}
	}
	if raw, ok := fields["unit_price"]; ok {
		var v float64
		if json.Unmarshal(raw, &v) == nil {
			item.UnitPrice = v
		}
	}
	if raw, ok := fields["type"]; ok {
		var v string
		if json.Unmarshal(raw, &v) == nil {
			item.Type = v
		}
	}
	if raw, ok := fields["description"]; ok {
		var v *string
		if json.Unmarshal(raw, &v) == nil {
			item.Description = ptrToNullString(v)
		}
	}
	if raw, ok := fields["surcharge_percent"]; ok {
		var v *float64
		if json.Unmarshal(raw, &v) == nil {
			item.SurchargePercent = ptrToNullFloat64(v)
		}
	}
	if raw, ok := fields["tag"]; ok {
		var v *string
		if json.Unmarshal(raw, &v) == nil {
			item.Tag = ptrToNullString(v)
		}
	}
	if raw, ok := fields["sort_order"]; ok {
		var v int64
		if json.Unmarshal(raw, &v) == nil {
			item.SortOrder = v
		}
	}

	updated, err := h.queries.UpdateLineItem(ctx, repository.UpdateLineItemParams{
		ID:               itemID,
		OrgID:            orgID,
		Type:             item.Type,
		Name:             item.Name,
		Description:      item.Description,
		Quantity:         item.Quantity,
		Unit:             item.Unit,
		UnitPrice:        item.UnitPrice,
		SurchargePercent: item.SurchargePercent,
		SortOrder:        item.SortOrder,
		Tag:              item.Tag,
	})
	if err != nil {
		logger.Error("failed to update line item", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	totals, err := h.loadJobTotals(r, jobID, orgID)
	if err != nil {
		logger.Error("failed to calculate totals", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to calculate totals")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"item":   convertLineItem(updated),
		"totals": totals,
	})
}

// CreateLineItemJSON creates a new line item in a job.
func (h *Handler) CreateLineItemJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	var input struct {
		CategoryID string   `json:"category_id"`
		Type       string   `json:"type"`
		Name       string   `json:"name"`
		Quantity   float64  `json:"quantity"`
		Unit       string   `json:"unit"`
		UnitPrice  float64  `json:"unit_price"`
		Tag        *string  `json:"tag"`
		SortOrder  *int64   `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.CategoryID == "" || input.Name == "" || input.Type == "" {
		writeJSONError(w, http.StatusBadRequest, "category_id, type, and name are required")
		return
	}

	if !h.categoryBelongsToJob(r, input.CategoryID, jobID, orgID) {
		writeJSONError(w, http.StatusForbidden, "Category does not belong to this job")
		return
	}

	if input.Quantity == 0 {
		input.Quantity = 1
	}
	if input.Unit == "" {
		input.Unit = "ea"
	}

	var sortOrder int64
	if input.SortOrder != nil {
		sortOrder = *input.SortOrder
	}

	created, err := h.queries.CreateLineItem(ctx, repository.CreateLineItemParams{
		ID:               uuid.New().String(),
		OrgID:            orgID,
		CategoryID:       input.CategoryID,
		Type:             input.Type,
		Name:             input.Name,
		Description:      sql.NullString{},
		Quantity:         input.Quantity,
		Unit:             input.Unit,
		UnitPrice:        input.UnitPrice,
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        sortOrder,
		Tag:              ptrToNullString(input.Tag),
	})
	if err != nil {
		logger.Error("failed to create line item", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	totals, err := h.loadJobTotals(r, jobID, orgID)
	if err != nil {
		logger.Error("failed to calculate totals", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to calculate totals")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"item":   convertLineItem(created),
		"totals": totals,
	})
}

// DeleteLineItemJSON deletes a line item and returns updated totals.
func (h *Handler) DeleteLineItemJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")
	itemID := r.PathValue("itemId")

	// Verify item exists and belongs to this job
	item, err := h.queries.GetLineItem(ctx, repository.GetLineItemParams{ID: itemID, OrgID: orgID})
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "Item not found")
			return
		}
		logger.Error("failed to get line item", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load item")
		return
	}

	if !h.categoryBelongsToJob(r, item.CategoryID, jobID, orgID) {
		writeJSONError(w, http.StatusForbidden, "Item does not belong to this job")
		return
	}

	if err := h.queries.DeleteLineItem(ctx, repository.DeleteLineItemParams{ID: itemID, OrgID: orgID}); err != nil {
		logger.Error("failed to delete line item", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	totals, err := h.loadJobTotals(r, jobID, orgID)
	if err != nil {
		logger.Error("failed to calculate totals", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to calculate totals")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"totals": totals})
}

// ReorderItemsJSON updates sort_order for a batch of items.
func (h *Handler) ReorderItemsJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	var input struct {
		Items []struct {
			ID        string `json:"id"`
			SortOrder int64  `json:"sort_order"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("failed to begin transaction", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to reorder")
		return
	}
	defer func() { _ = tx.Rollback() }()

	qtx := h.queries.WithTx(tx)
	for _, item := range input.Items {
		if err := qtx.UpdateLineItemSortOrder(ctx, repository.UpdateLineItemSortOrderParams{
			SortOrder: item.SortOrder,
			ID:        item.ID,
			OrgID:     orgID,
		}); err != nil {
			logger.Error("failed to update sort order", "error", err, "item_id", item.ID)
			writeJSONError(w, http.StatusInternalServerError, "Failed to reorder")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to reorder")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// SearchItemsJSON searches item templates and returns JSON.
func (h *Handler) SearchItemsJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	query := r.URL.Query().Get("q")
	itemType := r.URL.Query().Get("type")

	if query == "" {
		writeJSON(w, http.StatusOK, []any{})
		return
	}

	items, err := h.queries.SearchItemTemplatesByType(ctx, repository.SearchItemTemplatesByTypeParams{
		OrgID:   orgID,
		Type:    itemType,
		Column3: sql.NullString{String: query, Valid: true},
	})
	if err != nil {
		logger.Error("failed to search items", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Search failed")
		return
	}

	// Convert to clean JSON
	type TemplateResult struct {
		ID           int64   `json:"id"`
		Type         string  `json:"type"`
		Category     string  `json:"category"`
		Name         string  `json:"name"`
		DefaultUnit  string  `json:"default_unit"`
		DefaultPrice float64 `json:"default_price"`
	}

	results := make([]TemplateResult, len(items))
	for i, item := range items {
		results[i] = TemplateResult{
			ID:           item.ID,
			Type:         item.Type,
			Category:     item.Category,
			Name:         item.Name,
			DefaultUnit:  item.DefaultUnit,
			DefaultPrice: item.DefaultPrice,
		}
	}

	writeJSON(w, http.StatusOK, results)
}

