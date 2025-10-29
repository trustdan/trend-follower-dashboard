package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/trading-engine/internal/api/responses"
	"github.com/yourusername/trading-engine/internal/storage"
)

// DecisionHandler handles trade decision-related API requests
type DecisionHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewDecisionHandler creates a new DecisionHandler
func NewDecisionHandler(db *storage.DB, logger *log.Logger) *DecisionHandler {
	return &DecisionHandler{
		db:     db,
		logger: logger,
	}
}

// SaveDecisionRequest represents the request body for saving a trade decision
type SaveDecisionRequest struct {
	Ticker       string  `json:"ticker"`
	Entry        float64 `json:"entry"`
	ATR          float64 `json:"atr"`
	Method       string  `json:"method"`
	BannerStatus string  `json:"banner_status"`
	Shares       int     `json:"shares"`
	Contracts    int     `json:"contracts"`
	Sector       string  `json:"sector"`
	Strategy     string  `json:"strategy"`
	RiskDollars  float64 `json:"risk_dollars"`
	Decision     string  `json:"decision"` // "GO" or "NO-GO"
	Notes        string  `json:"notes"`

	// Gate states
	BannerGreen    bool `json:"banner_green"`
	TimerComplete  bool `json:"timer_complete"`
	NotOnCooldown  bool `json:"not_on_cooldown"`
	HeatPassed     bool `json:"heat_passed"`
	SizingComplete bool `json:"sizing_complete"`
}

// SaveDecisionResponse represents the response from saving a decision
type SaveDecisionResponse struct {
	ID        int64     `json:"id"`
	Ticker    string    `json:"ticker"`
	Decision  string    `json:"decision"`
	Timestamp time.Time `json:"timestamp"`
}

// SaveDecision handles POST /api/decisions/save
func (h *DecisionHandler) SaveDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	var req SaveDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error parsing save decision request: %v", err)
		responses.BadRequest(w, err)
		return
	}

	h.logger.Printf("Saving trade decision: ticker=%s decision=%s", req.Ticker, req.Decision)

	// Validate request
	if req.Ticker == "" {
		responses.BadRequest(w, fmt.Errorf("ticker is required"))
		return
	}

	if req.Decision != "GO" && req.Decision != "NO-GO" {
		responses.BadRequest(w, fmt.Errorf("decision must be 'GO' or 'NO-GO'"))
		return
	}

	if req.Notes == "" {
		responses.BadRequest(w, fmt.Errorf("notes are required"))
		return
	}

	// For GO decisions, validate all gates passed
	if req.Decision == "GO" {
		if !req.BannerGreen {
			responses.BadRequest(w, fmt.Errorf("cannot save GO decision: banner is not GREEN"))
			return
		}
		if !req.TimerComplete {
			responses.BadRequest(w, fmt.Errorf("cannot save GO decision: timer not complete"))
			return
		}
		if !req.NotOnCooldown {
			responses.BadRequest(w, fmt.Errorf("cannot save GO decision: ticker on cooldown"))
			return
		}
		if !req.HeatPassed {
			responses.BadRequest(w, fmt.Errorf("cannot save GO decision: heat check not passed"))
			return
		}
		if !req.SizingComplete {
			responses.BadRequest(w, fmt.Errorf("cannot save GO decision: position sizing not complete"))
			return
		}
	}

	// Save decision to database
	timestamp := time.Now()
	id, err := h.db.SaveDecision(storage.Decision{
		Date:        timestamp.Format("2006-01-02"),
		Ticker:      req.Ticker,
		Action:      req.Decision,
		Entry:       req.Entry,
		ATR:         req.ATR,
		Shares:      req.Shares,
		Contracts:   req.Contracts,
		RiskDollars: req.RiskDollars,
		Banner:      req.BannerStatus,
		Method:      req.Method,
		Bucket:      req.Sector,
		Reason:      req.Notes,
		CreatedAt:   timestamp,
	})

	if err != nil {
		h.logger.Printf("Error saving decision to database: %v", err)
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Printf("Decision saved successfully: id=%d ticker=%s decision=%s", id, req.Ticker, req.Decision)

	// Return success response
	resp := SaveDecisionResponse{
		ID:        int64(id),
		Ticker:    req.Ticker,
		Decision:  req.Decision,
		Timestamp: timestamp,
	}

	responses.Success(w, resp)
}
