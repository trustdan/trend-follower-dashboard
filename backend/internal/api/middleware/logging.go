package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// CorrelationIDKey is the context key for correlation IDs
	CorrelationIDKey contextKey = "correlationID"
)

// Logging logs HTTP requests with correlation IDs and performance metrics
func Logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate or extract correlation ID
			correlationID := r.Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			// Add correlation ID to response header
			w.Header().Set("X-Correlation-ID", correlationID)

			// Create response writer wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Log request
			logger.Printf("[%s] --> %s %s %s",
				correlationID,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			// Process request
			next.ServeHTTP(rw, r)

			// Log response with duration
			duration := time.Since(start)
			logger.Printf("[%s] <-- %s %s %d %s",
				correlationID,
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
			)

			// Log performance warning if slow
			if duration > 500*time.Millisecond {
				logger.Printf("[%s] SLOW REQUEST: %s %s took %s",
					correlationID,
					r.Method,
					r.URL.Path,
					duration,
				)
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
