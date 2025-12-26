package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/handler/keyboard"
	"github.com/dukerupert/skalkaho/internal/handler/quote"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, h *quote.Handler, kh *keyboard.Handler) {
	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Jobs
	mux.HandleFunc("GET /", h.ListJobs)
	mux.HandleFunc("POST /jobs", h.CreateJob)
	mux.HandleFunc("GET /jobs/{id}", h.GetJob)
	mux.HandleFunc("PUT /jobs/{id}", h.UpdateJob)
	mux.HandleFunc("DELETE /jobs/{id}", h.DeleteJob)

	// Categories
	mux.HandleFunc("POST /jobs/{jobID}/categories", h.CreateCategory)
	mux.HandleFunc("POST /categories/{parentID}/subcategories", h.CreateSubcategory)
	mux.HandleFunc("PUT /categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)

	// Line Items
	mux.HandleFunc("POST /categories/{categoryID}/items", h.CreateLineItem)
	mux.HandleFunc("PUT /items/{id}", h.UpdateLineItem)
	mux.HandleFunc("DELETE /items/{id}", h.DeleteLineItem)

	// Item Templates (autocomplete)
	mux.HandleFunc("GET /items/search", h.SearchItemTemplates)

	// Settings
	mux.HandleFunc("GET /settings", h.GetSettings)
	mux.HandleFunc("PUT /settings", h.UpdateSettings)

	// Keyboard UI routes
	mux.HandleFunc("GET /k/", kh.ListJobs)
	mux.HandleFunc("GET /k/jobs/{id}", kh.GetJob)
	mux.HandleFunc("GET /k/categories/{id}", kh.GetCategory)

	// Keyboard UI actions
	mux.HandleFunc("POST /k/jobs", kh.CreateJob)
	mux.HandleFunc("PUT /k/jobs/{id}", kh.UpdateJob)
	mux.HandleFunc("DELETE /k/jobs/{id}", kh.DeleteJob)
	mux.HandleFunc("POST /k/jobs/{jobID}/categories", kh.CreateCategory)
	mux.HandleFunc("POST /k/categories/{parentID}/subcategories", kh.CreateSubcategory)
	mux.HandleFunc("DELETE /k/categories/{id}", kh.DeleteCategory)
	mux.HandleFunc("POST /k/categories/{categoryID}/items", kh.CreateLineItem)
	mux.HandleFunc("GET /k/categories/{categoryID}/form", kh.GetInlineForm)
	mux.HandleFunc("GET /k/category-form", kh.GetCategoryForm)
	mux.HandleFunc("GET /k/job-form", kh.GetJobForm)
	mux.HandleFunc("GET /k/items/search", kh.SearchItems)
	mux.HandleFunc("DELETE /k/items/{id}", kh.DeleteLineItem)
}
