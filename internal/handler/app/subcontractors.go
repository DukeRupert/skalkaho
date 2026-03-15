package app

import (
	"database/sql"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// SubcontractorsPageData extends PageData for the subcontractors page.
type SubcontractorsPageData struct {
	PageData
	Subcontractors    []repository.ListSubcontractorsRow
	Trades            []repository.Trade
	TotalSubs         int64
	TotalFavorites    int64
	Search            string
	TradeFilter       string
	EditSubcontractor *repository.Subcontractor
	EditTrades        []repository.GetSubcontractorTradesRow
}

// ListSubcontractors renders the subcontractors list with optional search/filter.
func (h *Handler) ListSubcontractors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	search := r.URL.Query().Get("search")
	tradeFilter := r.URL.Query().Get("trade")

	var subs []repository.ListSubcontractorsRow
	var err error

	switch {
	case search != "":
		searchRows, searchErr := h.queries.SearchSubcontractors(ctx, search)
		if searchErr != nil {
			h.logger.Error("searching subcontractors", "error", searchErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Convert SearchSubcontractorsRow to ListSubcontractorsRow (same shape)
		for _, sr := range searchRows {
			subs = append(subs, repository.ListSubcontractorsRow{
				ID: sr.ID, Name: sr.Name, Company: sr.Company,
				Phone: sr.Phone, Email: sr.Email, Address: sr.Address,
				Notes: sr.Notes, IsFavorite: sr.IsFavorite,
				CreatedAt: sr.CreatedAt, UpdatedAt: sr.UpdatedAt,
				PrimaryTrade: sr.PrimaryTrade,
			})
		}
	case tradeFilter != "":
		tradeRows, tradeErr := h.queries.ListSubcontractorsByTrade(ctx, tradeFilter)
		if tradeErr != nil {
			h.logger.Error("listing subcontractors by trade", "error", tradeErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		for _, tr := range tradeRows {
			subs = append(subs, repository.ListSubcontractorsRow{
				ID: tr.ID, Name: tr.Name, Company: tr.Company,
				Phone: tr.Phone, Email: tr.Email, Address: tr.Address,
				Notes: tr.Notes, IsFavorite: tr.IsFavorite,
				CreatedAt: tr.CreatedAt, UpdatedAt: tr.UpdatedAt,
				PrimaryTrade: sql.NullString{String: tr.PrimaryTrade, Valid: true},
			})
		}
	default:
		subs, err = h.queries.ListSubcontractors(ctx)
		if err != nil {
			h.logger.Error("listing subcontractors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	trades, err := h.queries.ListTrades(ctx)
	if err != nil {
		h.logger.Error("listing trades", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalSubs, err := h.queries.CountSubcontractors(ctx)
	if err != nil {
		h.logger.Error("counting subcontractors", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalFavs, err := h.queries.CountFavoriteSubcontractors(ctx)
	if err != nil {
		h.logger.Error("counting favorites", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if subs == nil {
		subs = []repository.ListSubcontractorsRow{}
	}

	data := SubcontractorsPageData{
		PageData:       h.pageData(r, "subcontractors"),
		Subcontractors: subs,
		Trades:         trades,
		TotalSubs:      totalSubs,
		TotalFavorites: totalFavs,
		Search:         search,
		TradeFilter:    tradeFilter,
	}

	if isHTMXPartial(r) {
		if err := h.renderer.RenderPartial(w, "subcontractors.html", "sub-rows", data); err != nil {
			h.logger.Error("rendering sub rows", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.renderer.Render(w, "subcontractors.html", data); err != nil {
		h.logger.Error("rendering subcontractors", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewSubcontractorModal returns the new subcontractor form modal.
func (h *Handler) NewSubcontractorModal(w http.ResponseWriter, r *http.Request) {
	trades, err := h.queries.ListTrades(r.Context())
	if err != nil {
		h.logger.Error("listing trades", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := SubcontractorsPageData{
		PageData: h.pageData(r, "subcontractors"),
		Trades:   trades,
	}

	if err := h.renderer.RenderPartial(w, "subcontractors.html", "new-sub-modal", data); err != nil {
		h.logger.Error("rendering new sub modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateSubcontractor handles POST /subcontractors.
func (h *Handler) CreateSubcontractor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateSubcontractor(r.Context(), repository.CreateSubcontractorParams{
		ID:         id,
		Name:       name,
		Company:    toNullString(r.FormValue("company")),
		Phone:      toNullString(r.FormValue("phone")),
		Email:      toNullString(r.FormValue("email")),
		Address:    toNullString(r.FormValue("address")),
		Notes:      toNullString(r.FormValue("notes")),
		IsFavorite: r.FormValue("is_favorite") == "on",
	})
	if err != nil {
		h.logger.Error("creating subcontractor", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set trades
	h.setSubcontractorTrades(r, id)

	http.Redirect(w, r, "/subcontractors", http.StatusSeeOther)
}

// GetSubcontractorEditForm returns the edit form partial.
func (h *Handler) GetSubcontractorEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()

	sub, err := h.queries.GetSubcontractor(ctx, id)
	if err != nil {
		h.logger.Error("getting subcontractor", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	subTrades, err := h.queries.GetSubcontractorTrades(ctx, id)
	if err != nil {
		h.logger.Error("getting subcontractor trades", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	trades, err := h.queries.ListTrades(ctx)
	if err != nil {
		h.logger.Error("listing trades", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := SubcontractorsPageData{
		PageData:          h.pageData(r, "subcontractors"),
		Trades:            trades,
		EditSubcontractor: &sub,
		EditTrades:        subTrades,
	}

	if err := h.renderer.RenderPartial(w, "subcontractors.html", "edit-sub-modal", data); err != nil {
		h.logger.Error("rendering edit sub modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateSubcontractor handles POST /subcontractors/{id}.
func (h *Handler) UpdateSubcontractor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	_, err := h.queries.UpdateSubcontractor(r.Context(), repository.UpdateSubcontractorParams{
		ID:         id,
		Name:       name,
		Company:    toNullString(r.FormValue("company")),
		Phone:      toNullString(r.FormValue("phone")),
		Email:      toNullString(r.FormValue("email")),
		Address:    toNullString(r.FormValue("address")),
		Notes:      toNullString(r.FormValue("notes")),
		IsFavorite: r.FormValue("is_favorite") == "on",
	})
	if err != nil {
		h.logger.Error("updating subcontractor", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update trades
	h.setSubcontractorTrades(r, id)

	http.Redirect(w, r, "/subcontractors", http.StatusSeeOther)
}

// ToggleSubcontractorFavorite handles PATCH /subcontractors/{id}/favorite.
func (h *Handler) ToggleSubcontractorFavorite(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, err := h.queries.ToggleFavorite(r.Context(), id)
	if err != nil {
		h.logger.Error("toggling favorite", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return just the updated star button
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	starClass := "favorite-btn"
	if sub.IsFavorite {
		starClass += " favorite-btn--active"
	}
	w.Write([]byte(`<button class="` + starClass + `" hx-patch="/subcontractors/` + id + `/favorite" hx-swap="outerHTML">` +
		`<svg fill="` + boolToFill(sub.IsFavorite) + `" stroke="currentColor" viewBox="0 0 24 24" width="18" height="18">` +
		`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>` +
		`</svg></button>`))
}

// DeleteSubcontractor handles DELETE /subcontractors/{id}.
func (h *Handler) DeleteSubcontractor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteSubcontractor(r.Context(), id); err != nil {
		h.logger.Error("deleting subcontractor", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", "/subcontractors")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/subcontractors", http.StatusSeeOther)
}

// setSubcontractorTrades clears and re-sets trade associations.
func (h *Handler) setSubcontractorTrades(r *http.Request, subID string) {
	ctx := r.Context()
	if err := h.queries.ClearSubcontractorTrades(ctx, subID); err != nil {
		h.logger.Error("clearing subcontractor trades", "error", err)
		return
	}

	tradeIDs := r.Form["trade_ids"]
	primaryTradeID := r.FormValue("primary_trade_id")

	for _, tid := range tradeIDs {
		position := int64(1)
		if tid == primaryTradeID {
			position = 0
		}
		if err := h.queries.SetSubcontractorTrade(ctx, repository.SetSubcontractorTradeParams{
			SubcontractorID: subID,
			TradeID:         tid,
			Position:        position,
		}); err != nil {
			h.logger.Error("setting subcontractor trade", "error", err, "trade_id", tid)
		}
	}

	// If trades were selected but no primary was specified, make the first one primary
	if len(tradeIDs) > 0 && primaryTradeID == "" {
		if err := h.queries.ClearSubcontractorTrades(ctx, subID); err != nil {
			h.logger.Error("clearing trades for re-set", "error", err)
			return
		}
		for i, tid := range tradeIDs {
			position := int64(1)
			if i == 0 {
				position = 0
			}
			_ = h.queries.SetSubcontractorTrade(ctx, repository.SetSubcontractorTradeParams{
				SubcontractorID: subID,
				TradeID:         tid,
				Position:        position,
			})
		}
	}
}

func boolToFill(b bool) string {
	if b {
		return "currentColor"
	}
	return "none"
}
