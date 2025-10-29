// Package handlers provides HTTP handlers for API endpoints
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yourusername/trading-engine/internal/api/responses"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/storage"
)

// HeatHandler handles heat check API requests
type HeatHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewHeatHandler creates a new heat handler
func NewHeatHandler(db *storage.DB, logger *log.Logger) *HeatHandler {
	return &HeatHandler{
		db:     db,
		logger: logger,
	}
}

// HeatCheckRequest contains the API request for heat checking
type HeatCheckRequest struct {
	AddRiskDollars float64 `json:"add_risk_dollars"` // Proposed trade risk
	AddBucket      string  `json:"add_bucket"`        // Bucket for proposed trade
}

// CheckHeat handles POST /api/heat/check
func (h *HeatHandler) CheckHeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Parse request body
	var req HeatCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error parsing heat check request: %v", err)
		responses.BadRequest(w, err)
		return
	}

	h.logger.Printf("Heat check request: risk=$%.2f, bucket=%s", req.AddRiskDollars, req.AddBucket)

	// Get settings
	settings, err := h.db.GetSettings()
	if err != nil {
		h.logger.Printf("Error getting settings: %v", err)
		responses.InternalError(w, err)
		return
	}

	// Get open positions
	positions, err := h.db.GetOpenPositions()
	if err != nil {
		h.logger.Printf("Error getting open positions: %v", err)
		responses.InternalError(w, err)
		return
	}

	h.logger.Printf("Found %d open positions", len(positions))

	// Convert storage positions to domain positions
	domainPositions := make([]domain.Position, len(positions))
	for i, p := range positions {
		domainPositions[i] = domain.Position{
			Ticker:      p.Ticker,
			Bucket:      p.Bucket,
			RiskDollars: p.RiskDollars,
			UnitsOpen:   1, // Assuming 1 unit per position for now
			Status:      p.Status,
		}
	}

	// Build domain heat request
	heatReq := domain.HeatRequest{
		Equity:           settings.Equity,
		HeatCapPct:       settings.PortfolioCap / 100, // Convert from % to decimal (4.0 -> 0.04)
		BucketHeatCapPct: settings.BucketCap / 100,    // Convert from % to decimal (1.5 -> 0.015)
		AddRiskDollars:   req.AddRiskDollars,
		AddBucket:        req.AddBucket,
		OpenPositions:    domainPositions,
	}

	h.logger.Printf("Heat request: equity=$%.2f, portfolio_cap=%.2f%%, bucket_cap=%.2f%%",
		heatReq.Equity, heatReq.HeatCapPct*100, heatReq.BucketHeatCapPct*100)

	// Calculate heat
	result, err := domain.CalculateHeat(heatReq)
	if err != nil {
		h.logger.Printf("Error calculating heat: %v", err)
		responses.BadRequest(w, err)
		return
	}

	h.logger.Printf("Heat result: portfolio=%.2f/%.2f (%.1f%%), bucket=%.2f/%.2f (%.1f%%), allowed=%v",
		result.NewPortfolioHeat, result.PortfolioCap, result.PortfolioHeatPct,
		result.NewBucketHeat, result.BucketCap, result.BucketHeatPct,
		result.Allowed)

	// Return result
	responses.Success(w, result)
}
