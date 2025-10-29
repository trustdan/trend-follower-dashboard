package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateStockPosition_BasicExample(t *testing.T) {
	// Example from CLAUDE.md:
	// Equity=$10,000, Risk=0.75%, Entry=$180, ATR=$1.50, K=2
	// Expected: R=$75, StopDist=$3, InitStop=$177, Shares=25, ActualRisk=$75
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Entry:   180,
		ATR:     1.5,
		K:       2,
		Method:  "stock",
	}

	result, err := CalculateStockPosition(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 75.0, result.RiskDollars)
	assert.Equal(t, 3.0, result.StopDistance)
	assert.Equal(t, 177.0, result.InitialStop)
	assert.Equal(t, 25, result.Shares)
	assert.Equal(t, 0, result.Contracts)
	assert.Equal(t, 75.0, result.ActualRisk)
	assert.Equal(t, "stock", result.Method)
}

func TestCalculateStockPosition_PriceRanges(t *testing.T) {
	tests := []struct {
		name   string
		entry  float64
		atr    float64
		shares int
	}{
		{"Low price", 10, 0.25, 150},
		{"Medium price", 50, 0.50, 75},
		{"High price 1", 100, 1.00, 37},
		{"High price 2", 180, 1.50, 25},
		{"Very high price", 500, 5.00, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   tt.entry,
				ATR:     tt.atr,
				K:       2,
				Method:  "stock",
			}

			result, err := CalculateStockPosition(req)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.shares, result.Shares, "Shares mismatch for entry=%v, atr=%v", tt.entry, tt.atr)
			assert.InDelta(t, 75.0, result.RiskDollars, 0.01, "Risk dollars should be ~$75")
		})
	}
}

func TestCalculateStockPosition_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     SizingRequest
		wantErr error
	}{
		{
			name: "Zero ATR",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     0,
				K:       2,
			},
			wantErr: ErrInvalidATR,
		},
		{
			name: "Negative ATR",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     -1.5,
				K:       2,
			},
			wantErr: ErrInvalidATR,
		},
		{
			name: "Negative entry price",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   -180,
				ATR:     1.5,
				K:       2,
			},
			wantErr: ErrInvalidEntry,
		},
		{
			name: "Zero entry price",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   0,
				ATR:     1.5,
				K:       2,
			},
			wantErr: ErrInvalidEntry,
		},
		{
			name: "Zero equity",
			req: SizingRequest{
				Equity:  0,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       2,
			},
			wantErr: ErrInvalidEquity,
		},
		{
			name: "Negative equity",
			req: SizingRequest{
				Equity:  -10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       2,
			},
			wantErr: ErrInvalidEquity,
		},
		{
			name: "Zero K multiple",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       0,
			},
			wantErr: ErrInvalidK,
		},
		{
			name: "Negative K multiple",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   180,
				ATR:     1.5,
				K:       -2,
			},
			wantErr: ErrInvalidK,
		},
		{
			name: "Zero risk percent",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 0,
				Entry:   180,
				ATR:     1.5,
				K:       2,
			},
			wantErr: ErrInvalidRiskPct,
		},
		{
			name: "Risk percent > 1",
			req: SizingRequest{
				Equity:  10000,
				RiskPct: 1.5,
				Entry:   180,
				ATR:     1.5,
				K:       2,
			},
			wantErr: ErrInvalidRiskPct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateStockPosition(tt.req)
			assert.Error(t, err, "Should return error")
			assert.Nil(t, result, "Result should be nil on error")
			assert.ErrorIs(t, err, tt.wantErr, "Should return expected error type")
		})
	}
}

func TestCalculateStockPosition_ActualRiskNeverExceedsTarget(t *testing.T) {
	// Test various combinations to ensure actual risk never exceeds target risk
	entries := []float64{10, 25, 50, 100, 180, 500, 1000}
	atrs := []float64{0.10, 0.25, 0.50, 1.00, 2.50, 5.00, 10.00}

	for _, entry := range entries {
		for _, atr := range atrs {
			req := SizingRequest{
				Equity:  10000,
				RiskPct: 0.0075,
				Entry:   entry,
				ATR:     atr,
				K:       2,
				Method:  "stock",
			}

			result, err := CalculateStockPosition(req)
			require.NoError(t, err, "entry=%v, atr=%v", entry, atr)

			// Actual risk should never exceed target risk
			assert.LessOrEqual(t, result.ActualRisk, result.RiskDollars,
				"ActualRisk ($%.2f) should not exceed RiskDollars ($%.2f) for entry=%v, atr=%v",
				result.ActualRisk, result.RiskDollars, entry, atr)

			// Verify calculation consistency
			expectedStopDistance := float64(req.K) * req.ATR
			assert.Equal(t, expectedStopDistance, result.StopDistance)

			expectedInitialStop := req.Entry - expectedStopDistance
			assert.Equal(t, expectedInitialStop, result.InitialStop)

			expectedActualRisk := float64(result.Shares) * result.StopDistance
			assert.Equal(t, expectedActualRisk, result.ActualRisk)
		}
	}
}

func TestCalculatePositionSize_RouterStock(t *testing.T) {
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Entry:   180,
		ATR:     1.5,
		K:       2,
		Method:  "stock",
	}

	result, err := CalculatePositionSize(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 25, result.Shares)
	assert.Equal(t, "stock", result.Method)
}

func TestCalculatePositionSize_InvalidMethod(t *testing.T) {
	req := SizingRequest{
		Equity:  10000,
		RiskPct: 0.0075,
		Entry:   180,
		ATR:     1.5,
		K:       2,
		Method:  "invalid-method",
	}

	result, err := CalculatePositionSize(req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidMethod)
}
