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

// SizingHandler handles position sizing API requests
type SizingHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewSizingHandler creates a new sizing handler
func NewSizingHandler(db *storage.DB, logger *log.Logger) *SizingHandler {
	return &SizingHandler{
		db:     db,
		logger: logger,
	}
}

// CalculateSize handles POST /api/sizing/calculate
func (h *SizingHandler) CalculateSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Parse request body
	var req domain.SizingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error parsing sizing request: %v", err)
		responses.BadRequest(w, err)
		return
	}

	h.logger.Printf("Sizing request: %+v", req)

	// Calculate position size
	result, err := domain.CalculatePositionSize(req)
	if err != nil {
		h.logger.Printf("Error calculating position size: %v", err)
		responses.BadRequest(w, err)
		return
	}

	h.logger.Printf("Sizing result: %+v", result)

	// Return result
	responses.Success(w, result)
}
