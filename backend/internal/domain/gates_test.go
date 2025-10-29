package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSaveDecisionRequest_ValidGO(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		Entry:  180.0,
		ATR:    1.5,
		Method: "stock",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.NoError(t, err)
}

func TestValidateSaveDecisionRequest_ValidNOGO(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "NO-GO",
		Reason: "Bad setup",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.NoError(t, err)
}

func TestValidateSaveDecisionRequest_MissingTicker(t *testing.T) {
	req := SaveDecisionRequest{
		Action: "GO",
		Entry:  180.0,
		ATR:    1.5,
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ticker is required")
}

func TestValidateSaveDecisionRequest_InvalidAction(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "MAYBE",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "action must be GO or NO-GO")
}

func TestValidateSaveDecisionRequest_GOWithoutEntry(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		ATR:    1.5,
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "entry price must be positive")
}

func TestValidateSaveDecisionRequest_GOWithoutATR(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		Entry:  180.0,
		Method: "stock",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ATR must be positive")
}

func TestValidateSaveDecisionRequest_NOGOWithoutReason(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "NO-GO",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reason is required")
}

func TestValidateSaveDecisionRequest_OptionDeltaATRValid(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		Entry:  5.5,
		ATR:    1.5,
		Method: "opt-delta-atr",
		Delta:  0.30,
	}
	err := ValidateSaveDecisionRequest(req)
	assert.NoError(t, err)
}

func TestValidateSaveDecisionRequest_OptionDeltaATRInvalidDelta(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		Entry:  5.5,
		ATR:    1.5,
		Method: "opt-delta-atr",
		Delta:  1.5,
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delta must be between 0 and 1")
}

func TestValidateSaveDecisionRequest_OptionMaxLossValid(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker:  "AAPL",
		Action:  "GO",
		Entry:   5.5,
		Method:  "opt-maxloss",
		MaxLoss: 0.75,
	}
	err := ValidateSaveDecisionRequest(req)
	assert.NoError(t, err)
}

func TestValidateSaveDecisionRequest_OptionMaxLossInvalid(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		Entry:  5.5,
		Method: "opt-maxloss",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max-loss must be positive")
}

func TestValidateSaveDecisionRequest_InvalidMethod(t *testing.T) {
	req := SaveDecisionRequest{
		Ticker: "AAPL",
		Action: "GO",
		Entry:  180.0,
		ATR:    1.5,
		Method: "invalid",
	}
	err := ValidateSaveDecisionRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid method")
}

// Mock gate checker for testing
type MockGateChecker struct {
	BannerError        error
	CandidatesError    error
	ImpulseError       error
	CooldownError      error
	HeatError          error
}

func (m *MockGateChecker) CheckBannerGreen(ticker string) error {
	return m.BannerError
}

func (m *MockGateChecker) CheckTickerInCandidates(ticker, date string) error {
	return m.CandidatesError
}

func (m *MockGateChecker) CheckImpulseBrake(ticker string) error {
	return m.ImpulseError
}

func (m *MockGateChecker) CheckBucketCooldown(bucket string) error {
	return m.CooldownError
}

func (m *MockGateChecker) CheckHeatCaps(addRisk float64, bucket string) error {
	return m.HeatError
}

func TestValidateHardGates_AllPass(t *testing.T) {
	checker := &MockGateChecker{}
	result, err := ValidateHardGates(checker, "AAPL", "Tech/Comm", 75.0, "2025-10-27")

	assert.NoError(t, err)
	assert.True(t, result.AllPassed)
	assert.Empty(t, result.FailedGates)
	assert.Empty(t, result.FailureReasons)
}

func TestValidateHardGates_BannerFails(t *testing.T) {
	checker := &MockGateChecker{
		BannerError: assert.AnError,
	}
	result, err := ValidateHardGates(checker, "AAPL", "Tech/Comm", 75.0, "2025-10-27")

	assert.NoError(t, err)
	assert.False(t, result.AllPassed)
	assert.Contains(t, result.FailedGates, "Banner")
	assert.Len(t, result.FailureReasons, 1)
}

func TestValidateHardGates_MultipleGatesFail(t *testing.T) {
	checker := &MockGateChecker{
		BannerError:     assert.AnError,
		CandidatesError: assert.AnError,
		ImpulseError:    assert.AnError,
	}
	result, err := ValidateHardGates(checker, "AAPL", "Tech/Comm", 75.0, "2025-10-27")

	assert.NoError(t, err)
	assert.False(t, result.AllPassed)
	assert.Len(t, result.FailedGates, 3)
	assert.Contains(t, result.FailedGates, "Banner")
	assert.Contains(t, result.FailedGates, "Candidates")
	assert.Contains(t, result.FailedGates, "ImpulseBrake")
}

func TestValidateHardGates_HeatCapsFail(t *testing.T) {
	checker := &MockGateChecker{
		HeatError: assert.AnError,
	}
	result, err := ValidateHardGates(checker, "AAPL", "Tech/Comm", 75.0, "2025-10-27")

	assert.NoError(t, err)
	assert.False(t, result.AllPassed)
	assert.Contains(t, result.FailedGates, "HeatCaps")
}
