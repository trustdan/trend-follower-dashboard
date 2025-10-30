package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/yourusername/trading-engine/internal/api/handlers"
	"github.com/yourusername/trading-engine/internal/api/middleware"
	"github.com/yourusername/trading-engine/internal/storage"
	"github.com/yourusername/trading-engine/internal/webui"
)

// ServerCommand runs the HTTP server
func ServerCommand() {
	// Parse flags
	listen := flag.String("listen", "127.0.0.1:8080", "Address to listen on")
	dbPath := flag.String("db", getDefaultDBPath(), "Path to database file")
	flag.Parse()

	// Initialize logger (simple stdout logger for now)
	logger := log.New(os.Stdout, "[TF-Engine] ", log.LstdFlags)

	logger.Println("Starting TF-Engine HTTP Server...")

	// Initialize database
	absDBPath, _ := filepath.Abs(*dbPath)
	logger.Printf("Opening database: %s", absDBPath)

	db, err := storage.New(*dbPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	settingsHandler := handlers.NewSettingsHandler(db, logger)
	positionsHandler := handlers.NewPositionsHandler(db, logger)
	candidatesHandler := handlers.NewCandidatesHandler(db, logger)
	sizingHandler := handlers.NewSizingHandler(db, logger)
	heatHandler := handlers.NewHeatHandler(db, logger)
	decisionsHandler := handlers.NewDecisionHandler(db, logger)
	calendarHandler := handlers.NewCalendarHandler(db, logger)

	// Create router
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/settings", settingsHandler.GetSettings)
	mux.HandleFunc("/api/positions", positionsHandler.GetPositions)
	mux.HandleFunc("/api/candidates", candidatesHandler.GetCandidates)
	mux.HandleFunc("/api/candidates/scan", candidatesHandler.ScanCandidates)
	mux.HandleFunc("/api/candidates/import", candidatesHandler.ImportCandidates)
	mux.HandleFunc("/api/sizing/calculate", sizingHandler.CalculateSize)
	mux.HandleFunc("/api/heat/check", heatHandler.CheckHeat)
	mux.HandleFunc("/api/decisions/save", decisionsHandler.SaveDecision)
	mux.HandleFunc("/api/calendar", calendarHandler.GetCalendar)

	// Serve embedded Svelte UI
	sfs, err := webui.Sub()
	if err != nil {
		logger.Printf("Warning: Could not load embedded UI: %v", err)
		logger.Println("API endpoints will still work")
	} else {
		mux.Handle("/", http.FileServer(http.FS(sfs)))
		logger.Println("Embedded UI loaded successfully")
	}

	// Apply middleware
	handler := middleware.Recovery(logger)(
		middleware.Logging(logger)(
			middleware.CORS(mux),
		),
	)

	// Create server
	srv := &http.Server{
		Addr:         *listen,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Printf("Server listening on http://%s", *listen)
		logger.Println("Press Ctrl+C to stop")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped")
}
