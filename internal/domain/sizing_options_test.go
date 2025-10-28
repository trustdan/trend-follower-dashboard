package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateOptionDeltaATRPosition_LowDelta(t *testing.T) {
	// Low delta (0.30) yields 0 contracts
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Entry:   180,
		ATR:     1.5,
		K:       2,
		Method:  "opt-delta-atr",
		Delta:   0.30,
	}

	result, err := CalculateOptionDeltaATRPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 3.0, result.StopDistance)
	assert.Equal(t, 177.0, result.InitialStop)
	assert.Equal(t, 0, result.Shares)
	assert.Equal(t, 0, result.Contracts)
	assert.Equal(t, 0.0, result.ActualRisk)
	assert.Equal(t, "opt-delta-atr", result.Method)
}

func TestCalculateOptionDeltaATRPosition_MediumDelta(t *testing.T) {
	// Medium delta (0.70) still yields 0 contracts
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Entry:   50,
		ATR:     1.0,
		K:       2,
		Method:  "opt-delta-atr",
		Delta:   0.70,
	}

	result, err := CalculateOptionDeltaATRPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 2.0, result.StopDistance)
	assert.Equal(t, 48.0, result.InitialStop)
	assert.Equal(t, 0, result.Shares)
	assert.Equal(t, 0, result.Contracts)
	assert.Equal(t, 0.0, result.ActualRisk)
}

func TestCalculateOptionDeltaATRPosition_HighDelta80(t *testing.T) {
	// High delta (0.80) with larger ATR yields contracts
	// StockShares = floor(75 / 4) = 18
	// DeltaShares = 18 × 0.80 = 14.4
	// Contracts = floor(14.4 / 100) = 0
	// Wait, we need more shares to get contracts
	// Let's use entry=100, atr=2.0, K=2
	// StockShares = floor(75 / 4) = 18
	// DeltaShares = 18 × 0.80 = 14.4
	// Contracts = 0

	// Actually to get 2 contracts we need:
	// Contracts = 2
	// DeltaShares = 200+
	// StockShares = 200 / 0.80 = 250
	// StopDistance = 75 / 250 = 0.3
	// With K=2, ATR = 0.15

	// Let's recalculate: Entry=100, ATR=2.0, K=2, Delta=0.80
	// StopDistance = 4.0
	// StockShares = floor(75 / 4) = 18
	// DeltaShares = 18 × 0.80 = 14.4
	// Contracts = 0

	// We need bigger account to get contracts. Let's try equity=100000
	req := SizingRequest{
		Equity:  100000,
		RiskPct: 0.0075,
		Entry:   100,
		ATR:     2.0,
		K:       2,
		Method:  "opt-delta-atr",
		Delta:   0.80,
	}

	result, err := CalculateOptionDeltaATRPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	// R = 100000 × 0.0075 = 750
	// StopDistance = 2 × 2 = 4
	// StockShares = floor(750 / 4) = 187
	// DeltaShares = 187 × 0.80 = 149.6
	// Contracts = floor(149.6 / 100) = 1
	// ActualRisk = 1 × 100 × 4 × 0.80 = 320

	assert.Equal(t, 750.0, result.RiskDollars)
	assert.Equal(t, 4.0, result.StopDistance)
	assert.Equal(t, 96.0, result.InitialStop)
	assert.Equal(t, 0, result.Shares)
	assert.Equal(t, 1, result.Contracts)
	assert.Equal(t, 320.0, result.ActualRisk)
}

func TestCalculateOptionDeltaATRPosition_VeryHighDelta90(t *testing.T) {
	// Very high delta (0.90) with large account
	req := SizingRequest{
		Equity:  100000,
		RiskPct: 0.0075,
		Entry:   100,
		ATR:     1.0,
		K:       2,
		Method:  "opt-delta-atr",
		Delta:   0.90,
	}

	result, err := CalculateOptionDeltaATRPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	// R = 100000 × 0.0075 = 750
	// StopDistance = 2 × 1 = 2
	// StockShares = floor(750 / 2) = 375
	// DeltaShares = 375 × 0.90 = 337.5
	// Contracts = floor(337.5 / 100) = 3
	// ActualRisk = 3 × 100 × 2 × 0.90 = 540

	assert.Equal(t, 750.0, result.RiskDollars)
	assert.Equal(t, 2.0, result.StopDistance)
	assert.Equal(t, 98.0, result.InitialStop)
	assert.Equal(t, 0, result.Shares)
	assert.Equal(t, 3, result.Contracts)
	assert.Equal(t, 540.0, result.ActualRisk)
}

func TestCalculateOptionDeltaATRPosition_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     SizingRequest
		wantErr error
	}{
		{
			name: "Negative delta",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       2,
				Delta:   -0.3,
			},
			wantErr: ErrInvalidDelta,
		},
		{
			name: "Delta > 1",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       2,
				Delta:   1.5,
			},
			wantErr: ErrInvalidDelta,
		},
		{
			name: "Zero delta",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       2,
				Delta:   0,
			},
			wantErr: ErrInvalidDelta,
		},
		{
			name: "Missing delta (defaults to 0)",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       2,
				// Delta not specified
			},
			wantErr: ErrInvalidDelta,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateOptionDeltaATRPosition(tt.req)
			assert.Error(t, err, "Should return error")
			assert.Nil(t, result, "Result should be nil on error")
			assert.ErrorIs(t, err, tt.wantErr, "Should return expected error type")
		})
	}
}

func TestCalculateOptionDeltaATRPosition_ActualRiskNeverExceedsTarget(t *testing.T) {
	// Test various combinations to ensure actual risk never exceeds target risk
	entries := []float64{50, 100, 180, 500}
	atrs := []float64{0.50, 1.00, 2.00, 5.00}
	deltas := []float64{0.30, 0.50, 0.70, 0.80, 0.90}

	for _, entry := range entries {
		for _, atr := range atrs {
			for _, delta := range deltas {
				req := SizingRequest{
					Equity:  10000,
					RiskPct: 0.0075,
					Entry:   entry,
					ATR:     atr,
					K:       2,
					Method:  "opt-delta-atr",
					Delta:   delta,
				}

				result, err := CalculateOptionDeltaATRPosition(req)
				require.NoError(t, err, "entry=%v, atr=%v, delta=%v", entry, atr, delta)

				// Actual risk should never exceed target risk
				assert.LessOrEqual(t, result.ActualRisk, result.RiskDollars,
					"ActualRisk ($%.2f) should not exceed RiskDollars ($%.2f) for entry=%v, atr=%v, delta=%v",
					result.ActualRisk, result.RiskDollars, entry, atr, delta)
			}
		}
	}
}

func TestCalculatePositionSize_RouterOptionDeltaATR(t *testing.T) {
	req := SizingRequest{
		Equity:  100000,
		RiskPct: 0.0075,
		Entry:   100,
		ATR:     1.0,
		K:       2,
		Method:  "opt-delta-atr",
		Delta:   0.90,
	}

	result, err := CalculatePositionSize(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 3, result.Contracts)
	assert.Equal(t, "opt-delta-atr", result.Method)
}

// Max-Loss Method Tests

func TestCalculateOptionMaxLossPosition_BasicExample(t *testing.T) {
	// Example from documentation:
	// Equity=$10,000, RiskPct=0.75%, MaxLoss=$50
	// R = $75, Contracts = 1, ActualRisk = $50
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 50,
	}

	result, err := CalculateOptionMaxLossPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 1, result.Contracts)
	assert.Equal(t, 50.0, result.ActualRisk)
	assert.Equal(t, 0, result.Shares)
	assert.Equal(t, 0.0, result.StopDistance)
	assert.Equal(t, 0.0, result.InitialStop)
	assert.Equal(t, "opt-maxloss", result.Method)
}

func TestCalculateOptionMaxLossPosition_HigherMaxLoss(t *testing.T) {
	// MaxLoss > RiskBudget = 0 contracts
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 100,
	}

	result, err := CalculateOptionMaxLossPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 0, result.Contracts)
	assert.Equal(t, 0.0, result.ActualRisk)
}

func TestCalculateOptionMaxLossPosition_LowerMaxLoss(t *testing.T) {
	// Lower MaxLoss = more contracts
	// R=$75, MaxLoss=$25, Contracts=3, ActualRisk=$75
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 25,
	}

	result, err := CalculateOptionMaxLossPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 3, result.Contracts)
	assert.Equal(t, 75.0, result.ActualRisk)
}

func TestCalculateOptionMaxLossPosition_FractionalRoundDown(t *testing.T) {
	// R=$75, MaxLoss=$30
	// Contracts = floor(75/30) = floor(2.5) = 2
	// ActualRisk = 2 × 30 = 60
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 30,
	}

	result, err := CalculateOptionMaxLossPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 2, result.Contracts)
	assert.Equal(t, 60.0, result.ActualRisk)
}

func TestCalculateOptionMaxLossPosition_VerySmallMaxLoss(t *testing.T) {
	// R=$75, MaxLoss=$10
	// Contracts = floor(75/10) = 7
	// ActualRisk = 7 × 10 = 70
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 10,
	}

	result, err := CalculateOptionMaxLossPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 7, result.Contracts)
	assert.Equal(t, 70.0, result.ActualRisk)
}

func TestCalculateOptionMaxLossPosition_ExactDivisor(t *testing.T) {
	// R=$75, MaxLoss=$15
	// Contracts = floor(75/15) = 5
	// ActualRisk = 5 × 15 = 75
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 15,
	}

	result, err := CalculateOptionMaxLossPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 5, result.Contracts)
	assert.Equal(t, 75.0, result.ActualRisk)
}

func TestCalculateOptionMaxLossPosition_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     SizingRequest
		wantErr error
	}{
		{
			name: "Zero max loss",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				MaxLoss: 0,
			},
			wantErr: ErrInvalidMaxLoss,
		},
		{
			name: "Negative max loss",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				MaxLoss: -50,
			},
			wantErr: ErrInvalidMaxLoss,
		},
		{
			name: "Zero equity",
			req: SizingRequest{
				Equity:  0,
				RiskPct: 0.0075,
				MaxLoss: 50,
			},
			wantErr: ErrInvalidEquity,
		},
		{
			name: "Invalid risk percent",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 1.5,
				MaxLoss: 50,
			},
			wantErr: ErrInvalidRiskPct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateOptionMaxLossPosition(tt.req)
			assert.Error(t, err, "Should return error")
			assert.Nil(t, result, "Result should be nil on error")
			assert.ErrorIs(t, err, tt.wantErr, "Should return expected error type")
		})
	}
}

func TestCalculateOptionMaxLossPosition_ActualRiskNeverExceedsTarget(t *testing.T) {
	// Test various max-loss values to ensure actual risk never exceeds target
	maxLosses := []float64{5, 10, 15, 20, 25, 30, 40, 50, 60, 75, 100, 150}

	for _, maxLoss := range maxLosses {
		req := SizingRequest{
			Equity:  10000,
			RiskPct: 0.0075,
			Method:  "opt-maxloss",
			MaxLoss: maxLoss,
		}

		result, err := CalculateOptionMaxLossPosition(req)
		require.NoError(t, err, "maxLoss=%v", maxLoss)

		// Actual risk should never exceed target risk
		assert.LessOrEqual(t, result.ActualRisk, result.RiskDollars,
			"ActualRisk ($%.2f) should not exceed RiskDollars ($%.2f) for maxLoss=%v",
			result.ActualRisk, result.RiskDollars, maxLoss)
	}
}

func TestCalculatePositionSize_RouterOptionMaxLoss(t *testing.T) {
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Method:  "opt-maxloss",
		MaxLoss: 50,
	}

	result, err := CalculatePositionSize(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.Contracts)
	assert.Equal(t, "opt-maxloss", result.Method)
}
