package quote

import (
	"database/sql"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/middleware"
)

// SearchItemTemplates searches for item templates by name.
func (h *Handler) SearchItemTemplates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	query := r.URL.Query().Get("name")
	if query == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	items, err := h.queries.SearchItemTemplates(ctx, sql.NullString{String: query, Valid: true})
	if err != nil {
		logger.Error("failed to search item templates", "error", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	if err := h.renderer.RenderPartial(w, "item_search_results", items); err != nil {
		logger.Error("failed to render search results", "error", err)
	}
}
