package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// RatesPageData extends PageData with rates-specific data.
type RatesPageData struct {
	PageData
	Rates      []repository.Rate
	Categories []repository.RateCategory
	// CategoryCounts maps category ID to rate count for badge display.
	CategoryCounts map[string]int64
	// CategoryNames maps category ID to name for table display.
	CategoryNames map[string]string
	CategoryID    string
	Search        string
	EditRate      *repository.Rate
}

// ListRates renders the rates list with category tabs and search.
func (h *Handler) ListRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	categories, err := h.queries.ListRateCategories(ctx)
	if err != nil {
		h.logger.Error("listing rate categories", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Build count and name maps
	counts := make(map[string]int64, len(categories))
	names := make(map[string]string, len(categories))
	for _, cat := range categories {
		names[cat.ID] = cat.Name
		count, err := h.queries.CountRatesByCategory(ctx, cat.ID)
		if err != nil {
			h.logger.Error("counting rates", "error", err, "category", cat.ID)
			continue
		}
		counts[cat.ID] = count
	}

	var rates []repository.Rate
	switch {
	case search != "" && categoryID != "":
		rates, err = h.queries.SearchRatesByCategory(ctx, repository.SearchRatesByCategoryParams{
			CategoryID: categoryID,
			SearchTerm: search,
		})
	case search != "":
		rates, err = h.queries.SearchRates(ctx, search)
	case categoryID != "":
		rates, err = h.queries.ListRatesByCategory(ctx, categoryID)
	default:
		rates, err = h.queries.ListRates(ctx)
	}
	if err != nil {
		h.logger.Error("listing rates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := RatesPageData{
		PageData:       h.pageData(r, "rates"),
		Rates:          rates,
		Categories:     categories,
		CategoryCounts: counts,
		CategoryNames:  names,
		CategoryID:     categoryID,
		Search:         search,
	}

	if isHTMXPartial(r) {
		if err := h.renderer.RenderPartial(w, "rates.html", "rate-rows", data); err != nil {
			h.logger.Error("rendering rate rows", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.renderer.Render(w, "rates.html", data); err != nil {
		h.logger.Error("rendering rates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateRate handles POST /rates.
func (h *Handler) CreateRate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	categoryID := r.FormValue("category_id")
	if name == "" || categoryID == "" {
		http.Error(w, "Name and category are required", http.StatusBadRequest)
		return
	}

	rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateRate(r.Context(), repository.CreateRateParams{
		ID:         id,
		Name:       name,
		CategoryID: categoryID,
		Supplier:   toNullString(r.FormValue("supplier")),
		Rate:       rate,
		Unit:       r.FormValue("unit"),
		Notes:      toNullString(r.FormValue("notes")),
	})
	if err != nil {
		h.logger.Error("creating rate", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/rates?category=%s", categoryID), http.StatusSeeOther)
}

// GetRateEditForm returns the edit form partial for a rate.
func (h *Handler) GetRateEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rate, err := h.queries.GetRate(r.Context(), id)
	if err != nil {
		h.logger.Error("getting rate", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	categories, err := h.queries.ListRateCategories(r.Context())
	if err != nil {
		h.logger.Error("listing rate categories", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := RatesPageData{
		PageData:   h.pageData(r, "rates"),
		Categories: categories,
		EditRate:   &rate,
	}

	if err := h.renderer.RenderPartial(w, "rates.html", "edit-modal", data); err != nil {
		h.logger.Error("rendering edit form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateRate handles POST /rates/{id}.
func (h *Handler) UpdateRate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	categoryID := r.FormValue("category_id")
	if name == "" || categoryID == "" {
		http.Error(w, "Name and category are required", http.StatusBadRequest)
		return
	}

	rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)

	_, err := h.queries.UpdateRate(r.Context(), repository.UpdateRateParams{
		ID:         id,
		Name:       name,
		CategoryID: categoryID,
		Supplier:   toNullString(r.FormValue("supplier")),
		Rate:       rate,
		Unit:       r.FormValue("unit"),
		Notes:      toNullString(r.FormValue("notes")),
	})
	if err != nil {
		h.logger.Error("updating rate", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/rates?category=%s", categoryID), http.StatusSeeOther)
}

// DeleteRate handles DELETE /rates/{id}.
func (h *Handler) DeleteRate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteRate(r.Context(), id); err != nil {
		h.logger.Error("deleting rate", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", "/rates")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/rates", http.StatusSeeOther)
}

// CreateRateCategory handles POST /rate-categories.
func (h *Handler) CreateRateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	categories, _ := h.queries.ListRateCategories(r.Context())
	sortOrder := int64(len(categories))

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateRateCategory(r.Context(), repository.CreateRateCategoryParams{
		ID:        id,
		Name:      name,
		SortOrder: sortOrder,
	})
	if err != nil {
		h.logger.Error("creating rate category", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/rates?category=%s", id), http.StatusSeeOther)
}

// DeleteRateCategory handles DELETE /rate-categories/{id}.
func (h *Handler) DeleteRateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteRateCategory(r.Context(), id); err != nil {
		h.logger.Error("deleting rate category", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", "/rates")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/rates", http.StatusSeeOther)
}
