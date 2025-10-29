package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fresh-start-trading-platform/svelte-poc-server/webui"
)

// Settings represents the trading system configuration
type Settings struct {
	Equity       float64 `json:"equity"`
	RiskPct      float64 `json:"riskPct"`
	PortfolioCap float64 `json:"portfolioCap"`
	BucketCap    float64 `json:"bucketCap"`
}

func main() {
	// Get embedded Svelte files
	sfs, err := webui.Sub()
	if err != nil {
		log.Fatalf("Failed to get embedded files: %v", err)
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// API endpoint: GET /api/settings
	// For POC: Using mock data instead of database
	mux.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Mock settings data (matches backend defaults)
		settings := Settings{
			Equity:       100000.0, // $100,000
			RiskPct:      0.75,     // 0.75%
			PortfolioCap: 4.0,      // 4.0%
			BucketCap:    1.5,      // 1.5%
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(settings); err != nil {
			log.Printf("Error encoding settings: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("Served settings: equity=$%.2f, risk=%.2f%%", settings.Equity, settings.RiskPct)
	})

	// Serve static files (SPA fallback)
	mux.Handle("/", http.FileServer(http.FS(sfs)))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting TF-Engine Svelte POC server on http://localhost:%s", port)
	log.Printf("Demonstrating: Svelte UI + Go backend + embedded files")
	log.Printf("Press Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
