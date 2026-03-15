package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/auth"
	apphandler "github.com/dukerupert/skalkaho/internal/handler/app"
	authhandler "github.com/dukerupert/skalkaho/internal/handler/auth"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, authH *authhandler.Handler, appH *apphandler.Handler, sm *auth.SessionManager) {
	// Health check (public)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static files (public)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Auth routes (public — session middleware added for redirect-if-logged-in check)
	mux.Handle("GET /login", sm.SessionMiddleware(http.HandlerFunc(authH.GetLogin)))
	mux.HandleFunc("POST /login", authH.PostLogin)
	mux.HandleFunc("GET /logout", authH.Logout)
	mux.HandleFunc("POST /logout", authH.Logout)

	// Protected app routes
	mux.Handle("GET /{$}", protect(sm, http.HandlerFunc(appH.ListProjects)))
	mux.Handle("GET /projects/new-modal", protect(sm, http.HandlerFunc(appH.NewProjectModal)))
	mux.Handle("POST /projects", protect(sm, http.HandlerFunc(appH.CreateProject)))
	mux.Handle("DELETE /projects/{id}", protect(sm, http.HandlerFunc(appH.DeleteProject)))
	mux.Handle("PATCH /projects/{id}/status", protect(sm, http.HandlerFunc(appH.UpdateProjectStatus)))
	mux.Handle("PATCH /projects/{id}/client", protect(sm, http.HandlerFunc(appH.UpdateProjectClient)))
	mux.Handle("GET /clients", protect(sm, http.HandlerFunc(appH.ListClients)))
	mux.Handle("GET /clients/new-modal", protect(sm, http.HandlerFunc(appH.NewClientModal)))
	mux.Handle("POST /clients", protect(sm, http.HandlerFunc(appH.CreateClient)))
	mux.Handle("GET /clients/{id}/edit", protect(sm, http.HandlerFunc(appH.GetClientEditForm)))
	mux.Handle("POST /clients/{id}", protect(sm, http.HandlerFunc(appH.UpdateClient)))
	mux.Handle("DELETE /clients/{id}", protect(sm, http.HandlerFunc(appH.DeleteClient)))
	mux.Handle("GET /materials", protect(sm, http.HandlerFunc(appH.ListMaterials)))
	mux.Handle("POST /materials", protect(sm, http.HandlerFunc(appH.CreateMaterial)))
	mux.Handle("GET /materials/{id}/edit", protect(sm, http.HandlerFunc(appH.GetMaterialEditForm)))
	mux.Handle("POST /materials/{id}", protect(sm, http.HandlerFunc(appH.UpdateMaterial)))
	mux.Handle("DELETE /materials/{id}", protect(sm, http.HandlerFunc(appH.DeleteMaterial)))
	mux.Handle("POST /suppliers", protect(sm, http.HandlerFunc(appH.CreateSupplier)))
	mux.Handle("DELETE /suppliers/{id}", protect(sm, http.HandlerFunc(appH.DeleteSupplier)))
	mux.Handle("GET /rates", protect(sm, http.HandlerFunc(appH.ListRates)))
	mux.Handle("GET /rates/new-modal", protect(sm, http.HandlerFunc(appH.NewRateModal)))
	mux.Handle("POST /rates", protect(sm, http.HandlerFunc(appH.CreateRate)))
	mux.Handle("GET /rates/{id}/edit", protect(sm, http.HandlerFunc(appH.GetRateEditForm)))
	mux.Handle("POST /rates/{id}", protect(sm, http.HandlerFunc(appH.UpdateRate)))
	mux.Handle("DELETE /rates/{id}", protect(sm, http.HandlerFunc(appH.DeleteRate)))
	mux.Handle("POST /rate-categories", protect(sm, http.HandlerFunc(appH.CreateRateCategory)))
	mux.Handle("DELETE /rate-categories/{id}", protect(sm, http.HandlerFunc(appH.DeleteRateCategory)))
	mux.Handle("GET /subcontractors", protect(sm, http.HandlerFunc(appH.ListSubcontractors)))
	mux.Handle("GET /subcontractors/new-modal", protect(sm, http.HandlerFunc(appH.NewSubcontractorModal)))
	mux.Handle("POST /subcontractors", protect(sm, http.HandlerFunc(appH.CreateSubcontractor)))
	mux.Handle("GET /subcontractors/{id}/edit", protect(sm, http.HandlerFunc(appH.GetSubcontractorEditForm)))
	mux.Handle("POST /subcontractors/{id}", protect(sm, http.HandlerFunc(appH.UpdateSubcontractor)))
	mux.Handle("DELETE /subcontractors/{id}", protect(sm, http.HandlerFunc(appH.DeleteSubcontractor)))
	mux.Handle("PATCH /subcontractors/{id}/favorite", protect(sm, http.HandlerFunc(appH.ToggleSubcontractorFavorite)))
	// Template management
	mux.Handle("GET /templates", protect(sm, http.HandlerFunc(appH.ListTemplates)))
	mux.Handle("GET /templates/new-modal", protect(sm, http.HandlerFunc(appH.NewTemplateModal)))
	mux.Handle("POST /templates", protect(sm, http.HandlerFunc(appH.CreateTemplate)))
	mux.Handle("GET /templates/{id}", protect(sm, http.HandlerFunc(appH.GetTemplate)))
	mux.Handle("GET /templates/{id}/edit-modal", protect(sm, http.HandlerFunc(appH.GetTemplateEditModal)))
	mux.Handle("POST /templates/{id}", protect(sm, http.HandlerFunc(appH.UpdateTemplate)))
	mux.Handle("DELETE /templates/{id}", protect(sm, http.HandlerFunc(appH.DeleteTemplate)))
	mux.Handle("POST /templates/{id}/sections", protect(sm, http.HandlerFunc(appH.CreateTemplateSection)))
	mux.Handle("POST /templates/{id}/sections/{sid}", protect(sm, http.HandlerFunc(appH.UpdateTemplateSection)))
	mux.Handle("DELETE /templates/{id}/sections/{sid}", protect(sm, http.HandlerFunc(appH.DeleteTemplateSection)))
	mux.Handle("POST /templates/{id}/sections/{sid}/subcategories", protect(sm, http.HandlerFunc(appH.CreateTemplateSubcategory)))
	mux.Handle("POST /templates/{id}/subcategories/{scid}", protect(sm, http.HandlerFunc(appH.UpdateTemplateSubcategory)))
	mux.Handle("DELETE /templates/{id}/subcategories/{scid}", protect(sm, http.HandlerFunc(appH.DeleteTemplateSubcategory)))
	mux.Handle("POST /templates/{id}/subcategories/{scid}/groups", protect(sm, http.HandlerFunc(appH.CreateTemplateComponentGroup)))
	mux.Handle("POST /templates/{id}/groups/{gid}", protect(sm, http.HandlerFunc(appH.UpdateTemplateComponentGroup)))
	mux.Handle("DELETE /templates/{id}/groups/{gid}", protect(sm, http.HandlerFunc(appH.DeleteTemplateComponentGroup)))

	mux.Handle("GET /projects/{id}", protect(sm, http.HandlerFunc(appH.GetProjectOverview)))
	mux.Handle("GET /projects/{id}/status-modal", protect(sm, http.HandlerFunc(appH.GetStatusModal)))
	mux.Handle("GET /projects/{id}/estimate", protect(sm, http.HandlerFunc(appH.EstimateBuilder)))

	// Estimate builder API (JSON)
	mux.Handle("GET /api/estimate/{id}", protect(sm, http.HandlerFunc(appH.GetEstimate)))
	mux.Handle("POST /api/estimate/{id}", protect(sm, http.HandlerFunc(appH.SaveEstimate)))

	// Quote management (protected)
	mux.Handle("GET /projects/{id}/quotes", protect(sm, http.HandlerFunc(appH.ListQuotes)))
	mux.Handle("POST /projects/{id}/quotes", protect(sm, http.HandlerFunc(appH.CreateQuote)))
	mux.Handle("POST /quotes/{id}/send", protect(sm, http.HandlerFunc(appH.SendQuote)))
	mux.Handle("POST /quotes/{id}/resend", protect(sm, http.HandlerFunc(appH.ResendQuote)))
	mux.Handle("GET /quotes/{id}/preview", protect(sm, http.HandlerFunc(appH.PreviewQuote)))
	mux.Handle("GET /quotes/{id}/notes-modal", protect(sm, http.HandlerFunc(appH.GetNotesModal)))
	mux.Handle("POST /quotes/{id}/notes", protect(sm, http.HandlerFunc(appH.UpdateQuoteNotes)))
	mux.Handle("GET /quotes/{id}/send-modal", protect(sm, http.HandlerFunc(appH.GetSendModal)))

	// Public quote page (no auth)
	mux.HandleFunc("GET /q/{token}", appH.GetQuotePage)
	mux.HandleFunc("POST /q/{token}", appH.SubmitSignature)
}

// protect wraps a handler with session loading and authentication requirement.
func protect(sm *auth.SessionManager, handler http.Handler) http.Handler {
	return sm.SessionMiddleware(auth.RequireAuth(handler))
}
