package quote

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

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

	// Return line item partial for HTMX
	if r.Header.Get("HX-Request") == "true" {
		if err := h.renderer.RenderPartial(w, "line_item", lineItem); err != nil {
			logger.Error("failed to render line item", "error", err)
		}
		return
	}

	// Get category to redirect
	category, _ := h.queries.GetCategory(ctx, categoryID)
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

	// Return updated line item for HTMX
	if r.Header.Get("HX-Request") == "true" {
		if err := h.renderer.RenderPartial(w, "line_item", lineItem); err != nil {
			logger.Error("failed to render line item", "error", err)
		}
		return
	}

	category, _ := h.queries.GetCategory(ctx, lineItem.CategoryID)
	http.Redirect(w, r, "/jobs/"+category.JobID, http.StatusSeeOther)
}

// DeleteLineItem deletes a line item.
func (h *Handler) DeleteLineItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	itemID := r.PathValue("id")

	if err := h.queries.DeleteLineItem(ctx, itemID); err != nil {
		logger.Error("failed to delete line item", "error", err)
		http.Error(w, "Failed to delete line item", http.StatusInternalServerError)
		return
	}

	// Return empty response for HTMX
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
