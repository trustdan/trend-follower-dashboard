package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/trading-engine/internal/storage"
)

// Server represents the HTTP server
type Server struct {
	db     *storage.DB
	addr   string
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(db *storage.DB, addr string) *Server {
	return &Server{
		db:   db,
		addr: addr,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", s.healthHandler)

	// API endpoints
	mux.HandleFunc("/api/size", corsMiddleware(s.sizeHandler))
	mux.HandleFunc("/api/checklist", corsMiddleware(s.checklistHandler))
	mux.HandleFunc("/api/decision", corsMiddleware(s.decisionHandler))
	mux.HandleFunc("/api/candidates", corsMiddleware(s.candidatesHandler))
	mux.HandleFunc("/api/heat", corsMiddleware(s.heatHandler))
	mux.HandleFunc("/api/timer", corsMiddleware(s.timerHandler))
	mux.HandleFunc("/api/cooldown", corsMiddleware(s.cooldownHandler))
	mux.HandleFunc("/api/positions", corsMiddleware(s.positionsHandler))
	mux.HandleFunc("/api/settings", corsMiddleware(s.settingsHandler))

	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Starting HTTP server on %s\n", s.addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Correlation-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response
func respondError(w http.ResponseWriter, status int, message string, corrID string) {
	respondJSON(w, status, map[string]string{
		"error":          message,
		"correlation_id": corrID,
	})
}

// healthHandler handles health checks
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "3.0.0-dev",
		"time":    time.Now().Format(time.RFC3339),
	})
}
