package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/handler/keyboard"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, h *keyboard.Handler) {
	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Jobs
	mux.HandleFunc("GET /", h.ListJobs)
	mux.HandleFunc("GET /jobs/{id}", h.GetJob)
	mux.HandleFunc("POST /jobs", h.CreateJob)
	mux.HandleFunc("PUT /jobs/{id}", h.UpdateJob)
	mux.HandleFunc("DELETE /jobs/{id}", h.DeleteJob)
	mux.HandleFunc("GET /job-form", h.GetJobForm)
	mux.HandleFunc("GET /jobs/{id}/markup", h.GetMarkupForm)
	mux.HandleFunc("PUT /jobs/{id}/markup", h.UpdateMarkup)

	// Categories
	mux.HandleFunc("GET /categories/{id}", h.GetCategory)
	mux.HandleFunc("POST /jobs/{jobID}/categories", h.CreateCategory)
	mux.HandleFunc("POST /categories/{parentID}/subcategories", h.CreateSubcategory)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)
	mux.HandleFunc("GET /category-form", h.GetCategoryForm)
	mux.HandleFunc("GET /categories/{id}/markup", h.GetCategoryMarkupForm)
	mux.HandleFunc("PUT /categories/{id}/markup", h.UpdateCategoryMarkup)

	// Line Items
	mux.HandleFunc("POST /categories/{categoryID}/items", h.CreateLineItem)
	mux.HandleFunc("GET /categories/{categoryID}/form", h.GetInlineForm)
	mux.HandleFunc("GET /items/search", h.SearchItems)
	mux.HandleFunc("GET /items/{id}/edit", h.GetEditForm)
	mux.HandleFunc("PUT /items/{id}", h.UpdateLineItem)
	mux.HandleFunc("DELETE /items/{id}", h.DeleteLineItem)

	// Settings
	mux.HandleFunc("GET /settings", h.GetSettings)
	mux.HandleFunc("PUT /settings", h.UpdateSettings)
}
