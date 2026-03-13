package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/auth"
	authhandler "github.com/dukerupert/skalkaho/internal/handler/auth"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, authH *authhandler.Handler, sm *auth.SessionManager) {
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

	// Protected routes
	// Placeholder: root redirects to login if not authenticated
	mux.Handle("GET /", protect(sm, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><head><title>Skalkaho</title></head><body><h1>Welcome to Skalkaho</h1><p>You are logged in.</p><a href="/logout">Sign out</a></body></html>`))
	})))
}

// protect wraps a handler with session loading and authentication requirement.
func protect(sm *auth.SessionManager, handler http.Handler) http.Handler {
	return sm.SessionMiddleware(auth.RequireAuth(handler))
}
