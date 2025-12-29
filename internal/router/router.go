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
	mux.HandleFunc("GET /jobs/{id}/rename", h.GetJobRenameForm)
	mux.HandleFunc("PUT /jobs/{id}/name", h.UpdateJobName)
	mux.HandleFunc("GET /jobs/{id}/order-list", h.GetOrderList)
	mux.HandleFunc("GET /jobs/{id}/site-materials", h.GetSiteMaterials)

	// Categories
	mux.HandleFunc("GET /categories/{id}", h.GetCategory)
	mux.HandleFunc("POST /jobs/{jobID}/categories", h.CreateCategory)
	mux.HandleFunc("POST /categories/{parentID}/subcategories", h.CreateSubcategory)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)
	mux.HandleFunc("GET /category-form", h.GetCategoryForm)
	mux.HandleFunc("GET /categories/{id}/markup", h.GetCategoryMarkupForm)
	mux.HandleFunc("PUT /categories/{id}/markup", h.UpdateCategoryMarkup)
	mux.HandleFunc("GET /categories/{id}/rename", h.GetCategoryRenameForm)
	mux.HandleFunc("PUT /categories/{id}/name", h.UpdateCategoryName)

	// Line Items
	mux.HandleFunc("POST /categories/{categoryID}/items", h.CreateLineItem)
	mux.HandleFunc("GET /categories/{categoryID}/form", h.GetInlineForm)
	mux.HandleFunc("GET /items/search", h.SearchItems)
	mux.HandleFunc("GET /items/{id}/edit", h.GetEditForm)
	mux.HandleFunc("PUT /items/{id}", h.UpdateLineItem)
	mux.HandleFunc("DELETE /items/{id}", h.DeleteLineItem)

	// Item Templates
	mux.HandleFunc("GET /items", h.ListItemTemplates)
	mux.HandleFunc("POST /items", h.CreateItemTemplate)
	mux.HandleFunc("GET /items/new", h.GetItemTemplateForm)
	mux.HandleFunc("GET /item-templates/{id}/edit", h.GetItemTemplateEditForm)
	mux.HandleFunc("PUT /item-templates/{id}", h.UpdateItemTemplate)
	mux.HandleFunc("DELETE /item-templates/{id}", h.DeleteItemTemplate)

	// Settings
	mux.HandleFunc("GET /settings", h.GetSettings)
	mux.HandleFunc("PUT /settings", h.UpdateSettings)
}
