package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger logs request information.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get request ID from context
			requestID := RequestIDFromContext(r.Context())

			// Create request-scoped logger
			reqLogger := logger.With(
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			// Add logger to context
			ctx := WithLogger(r.Context(), reqLogger)

			// Wrap response writer to capture status
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Process request
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Log completion
			reqLogger.Info("request completed",
				"status", wrapped.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
