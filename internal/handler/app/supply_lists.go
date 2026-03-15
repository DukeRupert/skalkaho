package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// ── View models ────────────────────────────────────────────────

// SupplyListIndexData holds data for the supply list index page.
type SupplyListIndexData struct {
	PageData
	Quote      repository.Quote
	Lists      []SupplyListView
	Suppliers  []repository.ListSuppliersForProjectRow
	ProjectID  string
}

// SupplyListView is a template-friendly view of a supply list.
type SupplyListView struct {
	ID                string
	Label             string
	ItemCount         int64
	TotalValue        float64
	GeneratedAt       string
	EstimatedDelivery string
	ActualDelivery    string
	HasActualDelivery bool
	Notes             string
}

// SupplyListShowData holds data for the supply list detail/print page.
type SupplyListShowData struct {
	PageData
	List              repository.SupplyList
	Quote             repository.Quote
	QuoteID           string
	ProjectName       string
	Items             []repository.SupplyListItem
	Groups            []SupplierGroup
	GrandTotal        float64
	GeneratedAtFmt    string
	EstDeliveryFmt    string
	ActDeliveryFmt    string
	IsAllSuppliers    bool
}

// SupplierGroup is a group of items by supplier for the print view.
type SupplierGroup struct {
	SupplierName string
	Items        []repository.SupplyListItem
	Subtotal     float64
}

// ── Helpers ────────────────────────────────────────────────────

func toNullTime(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func formatDate(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("Jan 2, 2006")
}

func toSupplyListViews(rows []repository.ListSupplyListsByQuoteRow) []SupplyListView {
	views := make([]SupplyListView, len(rows))
	for i, r := range rows {
		views[i] = SupplyListView{
			ID:                r.ID,
			Label:             r.Label,
			ItemCount:         r.ItemCount,
			TotalValue:        r.TotalValue,
			GeneratedAt:       r.GeneratedAt.Format("Jan 2, 2006"),
			EstimatedDelivery: formatDate(r.EstimatedDelivery),
			ActualDelivery:    formatDate(r.ActualDelivery),
			HasActualDelivery: r.ActualDelivery.Valid,
			Notes:             r.Notes.String,
		}
	}
	return views
}

func groupItemsBySupplier(items []repository.SupplyListItem) []SupplierGroup {
	var groups []SupplierGroup
	var current *SupplierGroup

	for _, item := range items {
		name := "No Supplier"
		if item.SupplierName.Valid {
			name = item.SupplierName.String
		}

		if current == nil || current.SupplierName != name {
			if current != nil {
				groups = append(groups, *current)
			}
			current = &SupplierGroup{SupplierName: name}
		}
		current.Items = append(current.Items, item)
		current.Subtotal += item.ExtendedPrice
	}
	if current != nil {
		groups = append(groups, *current)
	}
	return groups
}

// ── Supply List Index ──────────────────────────────────────────

// SupplyListIndex shows all supply lists for a quote.
func (h *Handler) SupplyListIndex(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	project, err := h.queries.GetProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	lists, err := h.queries.ListSupplyListsByQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("listing supply lists", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	suppliers, err := h.queries.ListSuppliersForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing suppliers for project", "error", err)
		suppliers = nil
	}

	data := SupplyListIndexData{
		PageData:  h.pageData(r, "projects"),
		Quote:     quote,
		Lists:     toSupplyListViews(lists),
		Suppliers: suppliers,
		ProjectID: quote.ProjectID,
	}
	data.Project = &ProjectStub{ID: project.ID, Name: project.Name}

	if err := h.renderer.Render(w, "supply_lists.html", data); err != nil {
		h.logger.Error("rendering supply list index", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ── Generate handlers ──────────────────────────────────────────

// GenerateAllSupplyList creates one supply list covering all suppliers on the quote.
func (h *Handler) GenerateAllSupplyList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if quote.Status != "signed" {
		http.Error(w, "Supply lists can only be generated for signed quotes", http.StatusBadRequest)
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		h.logger.Error("beginning transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := h.queries.WithTx(tx)

	// Create the supply list
	list, err := qtx.CreateSupplyList(ctx, repository.CreateSupplyListParams{
		ID:      newID(),
		QuoteID: quoteID,
		Label:   "All Suppliers",
	})
	if err != nil {
		h.logger.Error("creating supply list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Snapshot material line items
	materialItems, err := qtx.ListMaterialLineItemsForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing material line items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for i, item := range materialItems {
		_, err := qtx.CreateSupplyListItem(ctx, repository.CreateSupplyListItemParams{
			ID:            newID(),
			SupplyListID:  list.ID,
			MaterialID:    item.MaterialID,
			SupplierID:    toNullString(item.SupplierID),
			SupplierName:  toNullString(item.SupplierName),
			MaterialName:  item.MaterialName,
			SupplierCode:  item.SupplierCode,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			UnitPrice:     item.UnitPrice,
			ExtendedPrice: item.Quantity * item.UnitPrice,
			SortOrder:     int64(i),
		})
		if err != nil {
			h.logger.Error("creating supply list item", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Snapshot ad hoc items
	adHocItems, err := qtx.ListAdHocLineItemsForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing ad hoc line items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for i, item := range adHocItems {
		_, err := qtx.CreateSupplyListItem(ctx, repository.CreateSupplyListItemParams{
			ID:            newID(),
			SupplyListID:  list.ID,
			MaterialName:  item.ItemName,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			UnitPrice:     item.UnitPrice,
			ExtendedPrice: item.Quantity * item.UnitPrice,
			SortOrder:     int64(len(materialItems) + i),
		})
		if err != nil {
			h.logger.Error("creating ad hoc supply list item", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		h.logger.Error("committing transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/supply-lists", quoteID), http.StatusSeeOther)
}

// GeneratePerSupplierSupplyList creates one supply list per supplier on the quote.
func (h *Handler) GeneratePerSupplierSupplyList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if quote.Status != "signed" {
		http.Error(w, "Supply lists can only be generated for signed quotes", http.StatusBadRequest)
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		h.logger.Error("beginning transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := h.queries.WithTx(tx)

	materialItems, err := qtx.ListMaterialLineItemsForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing material line items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Group items by supplier
	supplierItems := make(map[string][]repository.ListMaterialLineItemsForProjectRow)
	supplierNames := make(map[string]string)
	for _, item := range materialItems {
		supplierItems[item.SupplierID] = append(supplierItems[item.SupplierID], item)
		supplierNames[item.SupplierID] = item.SupplierName
	}

	for supplierID, items := range supplierItems {
		list, err := qtx.CreateSupplyList(ctx, repository.CreateSupplyListParams{
			ID:         newID(),
			QuoteID:    quoteID,
			SupplierID: toNullString(supplierID),
			Label:      supplierNames[supplierID],
		})
		if err != nil {
			h.logger.Error("creating per-supplier supply list", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for i, item := range items {
			_, err := qtx.CreateSupplyListItem(ctx, repository.CreateSupplyListItemParams{
				ID:            newID(),
				SupplyListID:  list.ID,
				MaterialID:    item.MaterialID,
				SupplierID:    toNullString(item.SupplierID),
				SupplierName:  toNullString(item.SupplierName),
				MaterialName:  item.MaterialName,
				SupplierCode:  item.SupplierCode,
				Quantity:      item.Quantity,
				Unit:          item.Unit,
				UnitPrice:     item.UnitPrice,
				ExtendedPrice: item.Quantity * item.UnitPrice,
				SortOrder:     int64(i),
			})
			if err != nil {
				h.logger.Error("creating supply list item", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		h.logger.Error("committing transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/supply-lists", quoteID), http.StatusSeeOther)
}

// GenerateForSupplierSupplyList creates a supply list for a single specified supplier.
func (h *Handler) GenerateForSupplierSupplyList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	supplierID := r.PathValue("sid")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if quote.Status != "signed" {
		http.Error(w, "Supply lists can only be generated for signed quotes", http.StatusBadRequest)
		return
	}

	supplier, err := h.queries.GetSupplier(ctx, supplierID)
	if err != nil {
		h.logger.Error("getting supplier", "error", err)
		http.Error(w, "Supplier not found", http.StatusNotFound)
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		h.logger.Error("beginning transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := h.queries.WithTx(tx)

	list, err := qtx.CreateSupplyList(ctx, repository.CreateSupplyListParams{
		ID:         newID(),
		QuoteID:    quoteID,
		SupplierID: toNullString(supplierID),
		Label:      supplier.Name,
	})
	if err != nil {
		h.logger.Error("creating supply list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	materialItems, err := qtx.ListMaterialLineItemsForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing material line items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sortOrder := 0
	for _, item := range materialItems {
		if item.SupplierID != supplierID {
			continue
		}
		_, err := qtx.CreateSupplyListItem(ctx, repository.CreateSupplyListItemParams{
			ID:            newID(),
			SupplyListID:  list.ID,
			MaterialID:    item.MaterialID,
			SupplierID:    toNullString(item.SupplierID),
			SupplierName:  toNullString(item.SupplierName),
			MaterialName:  item.MaterialName,
			SupplierCode:  item.SupplierCode,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			UnitPrice:     item.UnitPrice,
			ExtendedPrice: item.Quantity * item.UnitPrice,
			SortOrder:     int64(sortOrder),
		})
		if err != nil {
			h.logger.Error("creating supply list item", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		sortOrder++
	}

	if err := tx.Commit(); err != nil {
		h.logger.Error("committing transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/supply-lists/%s", quoteID, list.ID), http.StatusSeeOther)
}

// ── Show / Print ───────────────────────────────────────────────

// SupplyListShow renders a single supply list (view + print).
func (h *Handler) SupplyListShow(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	listID := r.PathValue("lid")
	ctx := r.Context()

	list, err := h.queries.GetSupplyList(ctx, listID)
	if err != nil {
		h.logger.Error("getting supply list", "error", err)
		http.Error(w, "Supply list not found", http.StatusNotFound)
		return
	}

	if list.QuoteID != quoteID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	project, err := h.queries.GetProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	items, err := h.queries.ListSupplyListItemsGrouped(ctx, listID)
	if err != nil {
		h.logger.Error("listing supply list items", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var grandTotal float64
	for _, item := range items {
		grandTotal += item.ExtendedPrice
	}

	data := SupplyListShowData{
		PageData:       h.pageData(r, "projects"),
		List:           list,
		Quote:          quote,
		QuoteID:        quoteID,
		ProjectName:    project.Name,
		Items:          items,
		Groups:         groupItemsBySupplier(items),
		GrandTotal:     grandTotal,
		GeneratedAtFmt: list.GeneratedAt.Format("January 2, 2006"),
		EstDeliveryFmt: formatDate(list.EstimatedDelivery),
		ActDeliveryFmt: formatDate(list.ActualDelivery),
		IsAllSuppliers: !list.SupplierID.Valid,
	}
	data.Project = &ProjectStub{ID: project.ID, Name: project.Name}

	if err := h.renderer.Render(w, "supply_list_show.html", data); err != nil {
		h.logger.Error("rendering supply list show", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ── Update dates ───────────────────────────────────────────────

// UpdateSupplyListDates updates delivery dates on a supply list.
func (h *Handler) UpdateSupplyListDates(w http.ResponseWriter, r *http.Request) {
	listID := r.PathValue("lid")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err := h.queries.UpdateSupplyListDates(ctx, repository.UpdateSupplyListDatesParams{
		ID:                listID,
		EstimatedDelivery: toNullTime(r.FormValue("estimated_delivery")),
		ActualDelivery:    toNullTime(r.FormValue("actual_delivery")),
	})
	if err != nil {
		h.logger.Error("updating supply list dates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Re-render the dates partial
	list, err := h.queries.GetSupplyList(ctx, listID)
	if err != nil {
		h.logger.Error("getting supply list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		List           repository.SupplyList
		EstDeliveryFmt string
		ActDeliveryFmt string
		QuoteID        string
	}{
		List:           list,
		EstDeliveryFmt: formatDate(list.EstimatedDelivery),
		ActDeliveryFmt: formatDate(list.ActualDelivery),
		QuoteID:        list.QuoteID,
	}

	if err := h.renderer.RenderPartial(w, "supply_list_show.html", "dates-form", data); err != nil {
		h.logger.Error("rendering dates partial", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ── Update notes ───────────────────────────────────────────────

// UpdateSupplyListNotes updates notes on a supply list.
func (h *Handler) UpdateSupplyListNotes(w http.ResponseWriter, r *http.Request) {
	listID := r.PathValue("lid")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err := h.queries.UpdateSupplyListNotes(ctx, repository.UpdateSupplyListNotesParams{
		ID:    listID,
		Notes: toNullString(r.FormValue("notes")),
	})
	if err != nil {
		h.logger.Error("updating supply list notes", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Saved")
}

// ── Delete ─────────────────────────────────────────────────────

// DeleteSupplyList deletes a supply list and its items (CASCADE).
func (h *Handler) DeleteSupplyList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	listID := r.PathValue("lid")
	ctx := r.Context()

	// Verify the list belongs to this quote
	list, err := h.queries.GetSupplyList(ctx, listID)
	if err != nil {
		h.logger.Error("getting supply list", "error", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if list.QuoteID != quoteID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteSupplyList(ctx, listID); err != nil {
		h.logger.Error("deleting supply list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/quotes/%s/supply-lists", quoteID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/supply-lists", quoteID), http.StatusSeeOther)
}
