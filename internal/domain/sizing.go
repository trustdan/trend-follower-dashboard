package domain

import (
	"errors"
	"fmt"
)

// SizingRequest contains input parameters for position sizing
type SizingRequest struct {
	Equity    float64 `json:"equity"`
	RiskPct   float64 `json:"risk_pct"`
	Entry     float64 `json:"entry"`
	ATR       float64 `json:"atr_n"`
	K         int     `json:"k"`
	Method    string  `json:"method"` // "stock", "opt-delta-atr", "opt-maxloss"
	Delta     float64 `json:"delta,omitempty"`
	MaxLoss   float64 `json:"max_loss,omitempty"`
}

// SizingResult contains calculated position sizing
type SizingResult struct {
	RiskDollars  float64 `json:"risk_dollars"`
	StopDistance float64 `json:"stop_distance"`
	InitialStop  float64 `json:"initial_stop"`
	Shares       int     `json:"shares"`
	Contracts    int     `json:"contracts"`
	ActualRisk   float64 `json:"actual_risk"`
	Method       string  `json:"method"`
}

// Validation errors
var (
	ErrInvalidEquity  = errors.New("equity must be greater than zero")
	ErrInvalidRiskPct = errors.New("risk percent must be between 0 and 1")
	ErrInvalidEntry   = errors.New("entry price must be positive")
	ErrInvalidATR     = errors.New("ATR must be greater than zero")
	ErrInvalidK       = errors.New("K multiple must be positive")
	ErrInvalidMethod  = errors.New("method must be 'stock', 'opt-delta-atr', or 'opt-maxloss'")
	ErrInvalidDelta   = errors.New("delta must be between 0 and 1")
	ErrInvalidMaxLoss = errors.New("max loss must be positive")
)

// validateCommonInputs validates inputs common to all sizing methods
func validateCommonInputs(req SizingRequest) error {
	if req.Equity <= 0 {
		return ErrInvalidEquity
	}
	if req.RiskPct <= 0 || req.RiskPct > 1 {
		return fmt.Errorf("%w: got %.4f", ErrInvalidRiskPct, req.RiskPct)
	}
	if req.Entry <= 0 {
		return fmt.Errorf("%w: got %.2f", ErrInvalidEntry, req.Entry)
	}
	if req.ATR <= 0 {
		return fmt.Errorf("%w: got %.2f", ErrInvalidATR, req.ATR)
	}
	if req.K <= 0 {
		return fmt.Errorf("%w: got %d", ErrInvalidK, req.K)
	}
	return nil
}

// CalculatePositionSize routes to the appropriate sizing method
func CalculatePositionSize(req SizingRequest) (*SizingResult, error) {
	switch req.Method {
	case "stock":
		return CalculateStockPosition(req)
	case "opt-delta-atr":
		return CalculateOptionDeltaATRPosition(req)
	case "opt-maxloss":
		return CalculateOptionMaxLossPosition(req)
	default:
		return nil, fmt.Errorf("%w: got '%s'", ErrInvalidMethod, req.Method)
	}
}
