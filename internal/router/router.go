package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/handler/keyboard"
	"github.com/dukerupert/skalkaho/internal/handler/quote"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, h *quote.Handler, kh *keyboard.Handler) {
	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Keyboard UI routes (default)
	mux.HandleFunc("GET /", kh.ListJobs)
	mux.HandleFunc("GET /jobs/{id}", kh.GetJob)
	mux.HandleFunc("GET /categories/{id}", kh.GetCategory)

	// Keyboard UI actions (default)
	mux.HandleFunc("POST /jobs", kh.CreateJob)
	mux.HandleFunc("PUT /jobs/{id}", kh.UpdateJob)
	mux.HandleFunc("DELETE /jobs/{id}", kh.DeleteJob)
	mux.HandleFunc("POST /jobs/{jobID}/categories", kh.CreateCategory)
	mux.HandleFunc("POST /categories/{parentID}/subcategories", kh.CreateSubcategory)
	mux.HandleFunc("DELETE /categories/{id}", kh.DeleteCategory)
	mux.HandleFunc("POST /categories/{categoryID}/items", kh.CreateLineItem)
	mux.HandleFunc("GET /categories/{categoryID}/form", kh.GetInlineForm)
	mux.HandleFunc("GET /category-form", kh.GetCategoryForm)
	mux.HandleFunc("GET /job-form", kh.GetJobForm)
	mux.HandleFunc("GET /items/search", kh.SearchItems)
	mux.HandleFunc("GET /items/{id}/edit", kh.GetEditForm)
	mux.HandleFunc("PUT /items/{id}", kh.UpdateLineItem)
	mux.HandleFunc("DELETE /items/{id}", kh.DeleteLineItem)

	// Markup forms (default)
	mux.HandleFunc("GET /jobs/{id}/markup", kh.GetMarkupForm)
	mux.HandleFunc("PUT /jobs/{id}/markup", kh.UpdateMarkup)
	mux.HandleFunc("GET /categories/{id}/markup", kh.GetCategoryMarkupForm)
	mux.HandleFunc("PUT /categories/{id}/markup", kh.UpdateCategoryMarkup)

	// Mouse UI routes (/m/)
	mux.HandleFunc("GET /m/", h.ListJobs)
	mux.HandleFunc("POST /m/jobs", h.CreateJob)
	mux.HandleFunc("GET /m/jobs/{id}", h.GetJob)
	mux.HandleFunc("PUT /m/jobs/{id}", h.UpdateJob)
	mux.HandleFunc("DELETE /m/jobs/{id}", h.DeleteJob)

	// Mouse UI Categories
	mux.HandleFunc("POST /m/jobs/{jobID}/categories", h.CreateCategory)
	mux.HandleFunc("POST /m/categories/{parentID}/subcategories", h.CreateSubcategory)
	mux.HandleFunc("PUT /m/categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /m/categories/{id}", h.DeleteCategory)

	// Mouse UI Line Items
	mux.HandleFunc("POST /m/categories/{categoryID}/items", h.CreateLineItem)
	mux.HandleFunc("PUT /m/items/{id}", h.UpdateLineItem)
	mux.HandleFunc("DELETE /m/items/{id}", h.DeleteLineItem)

	// Mouse UI Item Templates (autocomplete)
	mux.HandleFunc("GET /m/items/search", h.SearchItemTemplates)

	// Settings (shared)
	mux.HandleFunc("GET /settings", h.GetSettings)
	mux.HandleFunc("PUT /settings", h.UpdateSettings)
}
