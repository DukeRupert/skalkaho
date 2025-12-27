package keyboard

import (
	"net/http"
	"strconv"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// GetSettings shows the settings page.
func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		logger.Error("failed to get settings", "error", err)
		http.Error(w, "Failed to load settings", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Settings": settings,
	}

	if err := h.renderer.Render(w, "settings", data); err != nil {
		logger.Error("failed to render settings", "error", err)
	}
}

// UpdateSettings updates the application settings.
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	surchargePercent, _ := strconv.ParseFloat(r.FormValue("default_surcharge_percent"), 64)

	_, err := h.queries.UpdateSettings(ctx, repository.UpdateSettingsParams{
		DefaultSurchargeMode:    r.FormValue("default_surcharge_mode"),
		DefaultSurchargePercent: surchargePercent,
	})
	if err != nil {
		logger.Error("failed to update settings", "error", err)
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	// For HTMX, trigger a toast notification
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Settings saved", "type": "success"}}`)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}
