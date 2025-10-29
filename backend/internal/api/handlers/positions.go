package handlers

import (
	"log"
	"net/http"

	"github.com/yourusername/trading-engine/internal/api/responses"
	"github.com/yourusername/trading-engine/internal/storage"
)

// PositionsHandler handles positions-related API requests
type PositionsHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewPositionsHandler creates a new positions handler
func NewPositionsHandler(db *storage.DB, logger *log.Logger) *PositionsHandler {
	return &PositionsHandler{
		db:     db,
		logger: logger,
	}
}

// GetPositions handles GET /api/positions
func (h *PositionsHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	positions, err := h.db.GetPositions()
	if err != nil {
		h.logger.Printf("Error getting positions: %v", err)
		responses.InternalError(w, err)
		return
	}

	// Return empty array instead of null if no positions
	if positions == nil {
		positions = []storage.Position{}
	}

	responses.Success(w, positions)
}
