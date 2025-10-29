// Package handlers provides HTTP handlers for API endpoints
package handlers

import (
	"log"
	"net/http"

	"github.com/yourusername/trading-engine/internal/api/responses"
	"github.com/yourusername/trading-engine/internal/storage"
)

// SettingsHandler handles settings-related API requests
type SettingsHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(db *storage.DB, logger *log.Logger) *SettingsHandler {
	return &SettingsHandler{
		db:     db,
		logger: logger,
	}
}

// GetSettings handles GET /api/settings
func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	settings, err := h.db.GetSettings()
	if err != nil {
		h.logger.Printf("Error getting settings: %v", err)
		responses.InternalError(w, err)
		return
	}

	responses.Success(w, settings)
}

// UpdateSettings handles PUT /api/settings
func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// TODO: Implement settings update (Phase 2+)
	responses.Error(w, http.StatusNotImplemented, nil)
}
