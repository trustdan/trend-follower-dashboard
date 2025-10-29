package domain

import (
	"fmt"
	"time"
)

// HardGatesResult represents the result of checking all 5 hard gates
type HardGatesResult struct {
	AllPassed      bool     `json:"all_passed"`
	FailedGates    []string `json:"failed_gates,omitempty"`
	FailureReasons []string `json:"failure_reasons,omitempty"`
}

// GateChecker defines the interface for checking each hard gate
type GateChecker interface {
	CheckBannerGreen(ticker string) error
	CheckTickerInCandidates(ticker, date string) error
	CheckImpulseBrake(ticker string) error
	CheckBucketCooldown(bucket string) error
	CheckHeatCaps(addRisk float64, bucket string) error
}

// ValidateHardGates checks all 5 hard gates for a GO decision
// This is the core discipline enforcement mechanism
//
// The 5 Hard Gates:
//  1. Banner GREEN (all 6 checklist items satisfied)
//  2. Ticker in today's candidates (from FINVIZ screen)
//  3. 2-minute impulse brake expired
//  4. Bucket not in cooldown (24hr after loss)
//  5. Heat caps not exceeded (4% portfolio, 1.5% bucket)
//
// All gates must pass for a GO decision to be saved.
func ValidateHardGates(checker GateChecker, ticker, bucket string, riskDollars float64, date string) (*HardGatesResult, error) {
	result := &HardGatesResult{
		AllPassed:      true,
		FailedGates:    make([]string, 0),
		FailureReasons: make([]string, 0),
	}

	// Gate 1: Banner GREEN
	if err := checker.CheckBannerGreen(ticker); err != nil {
		result.AllPassed = false
		result.FailedGates = append(result.FailedGates, "Banner")
		result.FailureReasons = append(result.FailureReasons, err.Error())
	}

	// Gate 2: Ticker in today's candidates
	if err := checker.CheckTickerInCandidates(ticker, date); err != nil {
		result.AllPassed = false
		result.FailedGates = append(result.FailedGates, "Candidates")
		result.FailureReasons = append(result.FailureReasons, err.Error())
	}

	// Gate 3: 2-minute impulse brake
	if err := checker.CheckImpulseBrake(ticker); err != nil {
		result.AllPassed = false
		result.FailedGates = append(result.FailedGates, "ImpulseBrake")
		result.FailureReasons = append(result.FailureReasons, err.Error())
	}

	// Gate 4: Bucket cooldown
	if bucket != "" {
		if err := checker.CheckBucketCooldown(bucket); err != nil {
			result.AllPassed = false
			result.FailedGates = append(result.FailedGates, "BucketCooldown")
			result.FailureReasons = append(result.FailureReasons, err.Error())
		}
	}

	// Gate 5: Heat caps
	if err := checker.CheckHeatCaps(riskDollars, bucket); err != nil {
		result.AllPassed = false
		result.FailedGates = append(result.FailedGates, "HeatCaps")
		result.FailureReasons = append(result.FailureReasons, err.Error())
	}

	return result, nil
}

// SaveDecisionRequest represents a request to save a trading decision
type SaveDecisionRequest struct {
	Ticker   string
	Action   string  // GO or NO-GO
	Entry    float64 // Required for GO
	ATR      float64 // Required for GO (stock or opt-delta-atr)
	Method   string  // "stock", "opt-delta-atr", "opt-maxloss"
	Delta    float64 // For opt-delta-atr
	MaxLoss  float64 // For opt-maxloss
	Bucket   string
	Reason   string // Required for NO-GO
	Date     string // Defaults to today
	CorrID   string
}

// SaveDecisionResult represents the result of saving a decision
type SaveDecisionResult struct {
	DecisionID   int                `json:"decision_id"`
	Ticker       string             `json:"ticker"`
	Action       string             `json:"action"`
	Shares       int                `json:"shares,omitempty"`
	Contracts    int                `json:"contracts,omitempty"`
	Entry        float64            `json:"entry,omitempty"`
	InitialStop  float64            `json:"initial_stop,omitempty"`
	RiskDollars  float64            `json:"risk_dollars,omitempty"`
	Banner       string             `json:"banner"`
	GatesChecked bool               `json:"gates_checked"`
	GatesResult  *HardGatesResult   `json:"gates_result,omitempty"`
	Message      string             `json:"message"`
}

// ValidateSaveDecisionRequest validates the save decision request
func ValidateSaveDecisionRequest(req SaveDecisionRequest) error {
	if req.Ticker == "" {
		return fmt.Errorf("ticker is required")
	}

	if req.Action != "GO" && req.Action != "NO-GO" {
		return fmt.Errorf("action must be GO or NO-GO, got: %s", req.Action)
	}

	if req.Action == "GO" {
		if req.Entry <= 0 {
			return fmt.Errorf("entry price must be positive for GO decision")
		}

		// Method defaults to stock
		if req.Method == "" {
			req.Method = "stock"
		}

		// Validate method-specific requirements
		switch req.Method {
		case "stock", "opt-delta-atr":
			if req.ATR <= 0 {
				return fmt.Errorf("ATR must be positive for %s method", req.Method)
			}
			if req.Method == "opt-delta-atr" && (req.Delta <= 0 || req.Delta > 1) {
				return fmt.Errorf("delta must be between 0 and 1 for opt-delta-atr method")
			}
		case "opt-maxloss":
			if req.MaxLoss <= 0 {
				return fmt.Errorf("max-loss must be positive for opt-maxloss method")
			}
		default:
			return fmt.Errorf("invalid method: %s (must be stock, opt-delta-atr, or opt-maxloss)", req.Method)
		}
	}

	if req.Action == "NO-GO" && req.Reason == "" {
		return fmt.Errorf("reason is required for NO-GO decision")
	}

	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	return nil
}
