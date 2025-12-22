package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/handler/quote"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, h *quote.Handler) {
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
	mux.HandleFunc("PUT /categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)

	// Line Items
	mux.HandleFunc("POST /categories/{categoryID}/items", h.CreateLineItem)
	mux.HandleFunc("PUT /items/{id}", h.UpdateLineItem)
	mux.HandleFunc("DELETE /items/{id}", h.DeleteLineItem)

	// Settings
	mux.HandleFunc("GET /settings", h.GetSettings)
	mux.HandleFunc("PUT /settings", h.UpdateSettings)
}
