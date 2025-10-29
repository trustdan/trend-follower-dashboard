package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/trading-engine/internal/api/responses"
	"github.com/yourusername/trading-engine/internal/scrape"
	"github.com/yourusername/trading-engine/internal/storage"
)

// CandidatesHandler handles candidates-related API requests
type CandidatesHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewCandidatesHandler creates a new candidates handler
func NewCandidatesHandler(db *storage.DB, logger *log.Logger) *CandidatesHandler {
	return &CandidatesHandler{
		db:     db,
		logger: logger,
	}
}

// GetCandidates handles GET /api/candidates?date=YYYY-MM-DD
func (h *CandidatesHandler) GetCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Get date from query params (default to today)
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	candidates, err := h.db.GetCandidates(dateStr)
	if err != nil {
		h.logger.Printf("Error getting candidates for %s: %v", dateStr, err)
		responses.InternalError(w, err)
		return
	}

	// Return empty array instead of null
	if candidates == nil {
		candidates = []storage.Candidate{}
	}

	responses.Success(w, candidates)
}

// ScanRequest represents the request body for scanning FINVIZ
type ScanRequest struct {
	Preset string `json:"preset"`
}

// ScanResponse represents the response for a FINVIZ scan
type ScanResponse struct {
	Count   int      `json:"count"`
	Tickers []string `json:"tickers"`
	Date    string   `json:"date"`
}

// ScanCandidates handles POST /api/candidates/scan
func (h *CandidatesHandler) ScanCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Parse request body
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding scan request: %v", err)
		responses.BadRequest(w, err)
		return
	}

	// Validate preset
	if req.Preset == "" {
		req.Preset = "TF_BREAKOUT_LONG" // Default preset
	}

	h.logger.Printf("Starting FINVIZ scan with preset: %s", req.Preset)

	// Get FINVIZ URL for preset (this should be configurable)
	// For now, use a hardcoded map
	presetURLs := map[string]string{
		"TF_BREAKOUT_LONG": "https://finviz.com/screener.ashx?v=111&f=ta_pattern_channelup,ta_perf_1w10o",
		// Add more presets as needed
	}

	url, ok := presetURLs[req.Preset]
	if !ok {
		h.logger.Printf("Unknown preset: %s", req.Preset)
		responses.BadRequest(w, nil)
		return
	}

	// Scrape FINVIZ
	tickers, err := scrape.ScrapeFinviz(url)
	if err != nil {
		h.logger.Printf("Error scraping FINVIZ: %v", err)
		responses.InternalError(w, err)
		return
	}

	h.logger.Printf("FINVIZ scan found %d tickers", len(tickers))

	// Return results (don't save yet, user will review and import)
	result := ScanResponse{
		Count:   len(tickers),
		Tickers: tickers,
		Date:    time.Now().Format("2006-01-02"),
	}

	responses.Success(w, result)
}

// ImportRequest represents the request body for importing candidates
type ImportRequest struct {
	Tickers []string `json:"tickers"`
	Date    string   `json:"date"`
}

// ImportCandidates handles POST /api/candidates/import
func (h *CandidatesHandler) ImportCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Parse request body
	var req ImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding import request: %v", err)
		responses.BadRequest(w, err)
		return
	}

	// Validate
	if len(req.Tickers) == 0 {
		h.logger.Printf("No tickers provided for import")
		responses.BadRequest(w, nil)
		return
	}

	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	h.logger.Printf("Importing %d candidates for %s", len(req.Tickers), req.Date)

	// Save to database
	if err := h.db.AddCandidates(req.Tickers, req.Date); err != nil {
		h.logger.Printf("Error importing candidates: %v", err)
		responses.InternalError(w, err)
		return
	}

	responses.Success(w, map[string]interface{}{
		"imported": len(req.Tickers),
		"date":     req.Date,
	})
}

// DeleteCandidate handles DELETE /api/candidates/:ticker
func (h *CandidatesHandler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// TODO: Implement candidate deletion (Phase 1 Step 9)
	responses.Error(w, http.StatusNotImplemented, nil)
}
