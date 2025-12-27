package quote

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// writeTotalsOOB writes the totals partial with OOB swap to the response
func (h *Handler) writeTotalsOOB(w http.ResponseWriter, r *http.Request, jobID string) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job for totals", "error", err)
		return
	}

	categories, err := h.queries.ListCategoriesByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list categories for totals", "error", err)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list line items for totals", "error", err)
		return
	}

	totals := h.calculateTotals(job, categories, lineItems)

	var buf bytes.Buffer
	if err := h.renderer.RenderToWriter(&buf, "totals_oob", totals); err != nil {
		logger.Error("failed to render totals OOB", "error", err)
		return
	}
	_, _ = w.Write(buf.Bytes())
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
	unitPrice, _ := strconv.ParseFloat(r.FormValue("unit_price"), 64)

	description := sql.NullString{}
	if desc := r.FormValue("description"); desc != "" {
		description = sql.NullString{String: desc, Valid: true}
	}

	lineItem, err := h.queries.CreateLineItem(ctx, repository.CreateLineItemParams{
		ID:               uuid.New().String(),
		CategoryID:       categoryID,
		Type:             r.FormValue("type"),
		Name:             r.FormValue("name"),
		Description:      description,
		Quantity:         quantity,
		Unit:             r.FormValue("unit"),
		UnitPrice:        unitPrice,
		SurchargePercent: sql.NullFloat64{},
		SortOrder:        0,
	})
	if err != nil {
		logger.Error("failed to create line item", "error", err)
		http.Error(w, "Failed to create line item", http.StatusInternalServerError)
		return
	}

	// Get category for jobID
	category, err := h.queries.GetCategory(ctx, categoryID)
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// Return line item partial for HTMX
	if r.Header.Get("HX-Request") == "true" {
		if err := h.renderer.RenderPartial(w, "line_item", lineItem); err != nil {
			logger.Error("failed to render line item", "error", err)
		}
		// Append OOB totals update
		h.writeTotalsOOB(w, r, category.JobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+category.JobID, http.StatusSeeOther)
}

// UpdateLineItem updates a line item.
func (h *Handler) UpdateLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get existing item
	existingItem, err := h.queries.GetLineItem(ctx, itemID)
	if err != nil {
		logger.Error("failed to get line item", "error", err)
		http.Error(w, "Line item not found", http.StatusNotFound)
		return
	}

	// Parse form values with defaults from existing item
	itemType := r.FormValue("type")
	if itemType == "" {
		itemType = existingItem.Type
	}

	name := r.FormValue("name")
	if name == "" {
		name = existingItem.Name
	}

	quantity := existingItem.Quantity
	if q := r.FormValue("quantity"); q != "" {
		if val, err := strconv.ParseFloat(q, 64); err == nil {
			quantity = val
		}
	}

	unit := r.FormValue("unit")
	if unit == "" {
		unit = existingItem.Unit
	}

	unitPrice := existingItem.UnitPrice
	if up := r.FormValue("unit_price"); up != "" {
		if val, err := strconv.ParseFloat(up, 64); err == nil {
			unitPrice = val
		}
	}

	surchargePercent := existingItem.SurchargePercent
	if sp := r.FormValue("surcharge_percent"); sp != "" {
		if val, err := strconv.ParseFloat(sp, 64); err == nil {
			surchargePercent = sql.NullFloat64{Float64: val, Valid: true}
		}
	} else if r.Form.Has("surcharge_percent") {
		// Form field exists but is empty - clear the surcharge
		surchargePercent = sql.NullFloat64{}
	}

	description := existingItem.Description
	if desc := r.FormValue("description"); desc != "" {
		description = sql.NullString{String: desc, Valid: true}
	}

	lineItem, err := h.queries.UpdateLineItem(ctx, repository.UpdateLineItemParams{
		ID:               itemID,
		Type:             itemType,
		Name:             name,
		Description:      description,
		Quantity:         quantity,
		Unit:             unit,
		UnitPrice:        unitPrice,
		SurchargePercent: surchargePercent,
		SortOrder:        existingItem.SortOrder,
	})
	if err != nil {
		logger.Error("failed to update line item", "error", err)
		http.Error(w, "Failed to update line item", http.StatusInternalServerError)
		return
	}

	// Get category for jobID
	category, err := h.queries.GetCategory(ctx, lineItem.CategoryID)
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// Return updated line item for HTMX
	if r.Header.Get("HX-Request") == "true" {
		if err := h.renderer.RenderPartial(w, "line_item", lineItem); err != nil {
			logger.Error("failed to render line item", "error", err)
		}
		// Append OOB totals update
		h.writeTotalsOOB(w, r, category.JobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+category.JobID, http.StatusSeeOther)
}

// DeleteLineItem deletes a line item.
func (h *Handler) DeleteLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")

	// Get item first to find jobID for OOB update
	item, err := h.queries.GetLineItem(ctx, itemID)
	if err != nil {
		logger.Error("failed to get line item", "error", err)
		http.Error(w, "Line item not found", http.StatusNotFound)
		return
	}

	category, err := h.queries.GetCategory(ctx, item.CategoryID)
	if err != nil {
		logger.Error("failed to get category", "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteLineItem(ctx, itemID); err != nil {
		logger.Error("failed to delete line item", "error", err)
		http.Error(w, "Failed to delete line item", http.StatusInternalServerError)
		return
	}

	// Return empty response for HTMX with OOB totals update
	if r.Header.Get("HX-Request") == "true" {
		h.writeTotalsOOB(w, r, category.JobID)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
