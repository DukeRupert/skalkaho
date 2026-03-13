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
	mux.Handle("POST /projects", protect(sm, http.HandlerFunc(appH.CreateProject)))
	mux.Handle("DELETE /projects/{id}", protect(sm, http.HandlerFunc(appH.DeleteProject)))
	mux.Handle("PATCH /projects/{id}/status", protect(sm, http.HandlerFunc(appH.UpdateProjectStatus)))
	mux.Handle("GET /clients", protect(sm, http.HandlerFunc(appH.Clients)))
	mux.Handle("GET /materials", protect(sm, http.HandlerFunc(appH.Materials)))
	mux.Handle("GET /rates", protect(sm, http.HandlerFunc(appH.Rates)))
	mux.Handle("GET /projects/{id}", protect(sm, http.HandlerFunc(appH.ProjectOverview)))
	mux.Handle("GET /projects/{id}/estimate", protect(sm, http.HandlerFunc(appH.EstimateBuilder)))
}

// protect wraps a handler with session loading and authentication requirement.
func protect(sm *auth.SessionManager, handler http.Handler) http.Handler {
	return sm.SessionMiddleware(auth.RequireAuth(handler))
}
