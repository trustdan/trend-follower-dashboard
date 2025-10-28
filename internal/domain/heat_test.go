package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateHeat_NoExistingPositions(t *testing.T) {
	// First trade of the day
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,   // 4% = $400
		BucketHeatCapPct: 0.015,  // 1.5% = $150
		AddRiskDollars:   75,
		AddBucket:        "Tech/Comm",
		OpenPositions:    []Position{},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 0.0, result.CurrentPortfolioHeat)
	assert.Equal(t, 75.0, result.NewPortfolioHeat)
	assert.Equal(t, 400.0, result.PortfolioCap)
	assert.False(t, result.PortfolioCapExceeded)
	assert.Equal(t, 0.0, result.PortfolioOverage)

	assert.Equal(t, 0.0, result.CurrentBucketHeat)
	assert.Equal(t, 75.0, result.NewBucketHeat)
	assert.Equal(t, 150.0, result.BucketCap)
	assert.False(t, result.BucketCapExceeded)
	assert.Equal(t, 0.0, result.BucketOverage)

	assert.True(t, result.Allowed)
	assert.Empty(t, result.RejectionReason)
}

func TestCalculateHeat_PortfolioApproachingCap(t *testing.T) {
	// 4 positions at $75 each = $300 heat
	// Adding to Finance bucket (no bucket heat conflict)
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   75,
		AddBucket:        "Finance", // Different bucket to avoid bucket cap
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
			{Ticker: "MSFT", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
			{Ticker: "XOM", Bucket: "Energy", RiskDollars: 75, Status: "Open"},
			{Ticker: "CVX", Bucket: "Energy", RiskDollars: 75, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 300.0, result.CurrentPortfolioHeat)
	assert.Equal(t, 375.0, result.NewPortfolioHeat)
	assert.False(t, result.PortfolioCapExceeded)
	assert.True(t, result.Allowed)
}

func TestCalculateHeat_PortfolioExceedsCap(t *testing.T) {
	// 5 positions at $75 each = $375 heat, adding $75 = $450 > $400
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04, // $400 cap
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   75,
		AddBucket:        "Tech/Comm",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
			{Ticker: "MSFT", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
			{Ticker: "NVDA", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
			{Ticker: "XOM", Bucket: "Energy", RiskDollars: 75, Status: "Open"},
			{Ticker: "CVX", Bucket: "Energy", RiskDollars: 75, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 375.0, result.CurrentPortfolioHeat)
	assert.Equal(t, 450.0, result.NewPortfolioHeat)
	assert.True(t, result.PortfolioCapExceeded)
	assert.Equal(t, 50.0, result.PortfolioOverage)
	assert.False(t, result.Allowed)
	assert.Contains(t, result.RejectionReason, "Portfolio heat")
	assert.Contains(t, result.RejectionReason, "450.00")
	assert.Contains(t, result.RejectionReason, "400.00")
}

func TestCalculateHeat_BucketExceedsCap(t *testing.T) {
	// Tech/Comm bucket has $100, adding $75 = $175 > $150
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015, // $150 cap
		AddRiskDollars:   75,
		AddBucket:        "Tech/Comm",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 100, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 100.0, result.CurrentBucketHeat)
	assert.Equal(t, 175.0, result.NewBucketHeat)
	assert.True(t, result.BucketCapExceeded)
	assert.Equal(t, 25.0, result.BucketOverage)
	assert.False(t, result.Allowed)
	assert.Contains(t, result.RejectionReason, "Bucket")
	assert.Contains(t, result.RejectionReason, "Tech/Comm")
}

func TestCalculateHeat_PortfolioOKBucketExceeds(t *testing.T) {
	// Portfolio at $200 (50% of $400 cap) but bucket at $130
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   75,
		AddBucket:        "Healthcare",
		OpenPositions: []Position{
			{Ticker: "JNJ", Bucket: "Healthcare", RiskDollars: 65, Status: "Open"},
			{Ticker: "PFE", Bucket: "Healthcare", RiskDollars: 65, Status: "Open"},
			{Ticker: "XOM", Bucket: "Energy", RiskDollars: 70, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 200.0, result.CurrentPortfolioHeat)
	assert.False(t, result.PortfolioCapExceeded)
	assert.Equal(t, 130.0, result.CurrentBucketHeat)
	assert.True(t, result.BucketCapExceeded)
	assert.False(t, result.Allowed)
}

func TestCalculateHeat_DifferentBucketNoConflict(t *testing.T) {
	// Tech/Comm has $140, but we're adding to Energy
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   75,
		AddBucket:        "Energy",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 70, Status: "Open"},
			{Ticker: "MSFT", Bucket: "Tech/Comm", RiskDollars: 70, Status: "Open"},
			{Ticker: "NVDA", Bucket: "Tech/Comm", RiskDollars: 60, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 200.0, result.CurrentPortfolioHeat)
	assert.Equal(t, 0.0, result.CurrentBucketHeat) // Energy bucket is empty
	assert.Equal(t, 75.0, result.NewBucketHeat)
	assert.False(t, result.BucketCapExceeded)
	assert.True(t, result.Allowed)
}

func TestCalculateHeat_MultiplePositionsAcrossBuckets(t *testing.T) {
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   75,
		AddBucket:        "Finance",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 100, Status: "Open"},
			{Ticker: "JNJ", Bucket: "Healthcare", RiskDollars: 75, Status: "Open"},
			{Ticker: "XOM", Bucket: "Energy", RiskDollars: 50, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 225.0, result.CurrentPortfolioHeat)
	assert.Equal(t, 300.0, result.NewPortfolioHeat)
	assert.Equal(t, 0.0, result.CurrentBucketHeat) // Finance bucket is empty
	assert.True(t, result.Allowed)
}

func TestCalculateHeat_ZeroRiskQuery(t *testing.T) {
	// Query only - no new trade
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   0,
		AddBucket:        "Tech/Comm",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 100, Status: "Open"},
			{Ticker: "MSFT", Bucket: "Energy", RiskDollars: 200, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 300.0, result.CurrentPortfolioHeat)
	assert.Equal(t, 300.0, result.NewPortfolioHeat)
	assert.True(t, result.Allowed)
}

func TestCalculateHeat_ClosedPositionsIgnored(t *testing.T) {
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,
		BucketHeatCapPct: 0.015,
		AddRiskDollars:   75,
		AddBucket:        "Tech/Comm",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
			{Ticker: "MSFT", Bucket: "Tech/Comm", RiskDollars: 100, Status: "Closed"}, // Ignored
			{Ticker: "NVDA", Bucket: "Tech/Comm", RiskDollars: 50, Status: "Closed"},  // Ignored
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	assert.Equal(t, 75.0, result.CurrentPortfolioHeat)  // Only AAPL
	assert.Equal(t, 75.0, result.CurrentBucketHeat)     // Only AAPL
	assert.Equal(t, 150.0, result.NewPortfolioHeat)
	assert.True(t, result.Allowed)
}

func TestCalculateHeat_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     HeatRequest
		wantErr string
	}{
		{
			name: "Zero equity",
			req: HeatRequest{
				Equity:           0,
				HeatCapPct:       0.04,
				BucketHeatCapPct: 0.015,
				AddRiskDollars:   75,
				AddBucket:        "Tech",
			},
			wantErr: "equity must be greater than zero",
		},
		{
			name: "Negative risk",
			req: HeatRequest{
				Equity:           10000,
				HeatCapPct:       0.04,
				BucketHeatCapPct: 0.015,
				AddRiskDollars:   -75,
				AddBucket:        "Tech",
			},
			wantErr: "add_risk_dollars must be non-negative",
		},
		{
			name: "Invalid heat cap percent",
			req: HeatRequest{
				Equity:           10000,
				HeatCapPct:       1.5,
				BucketHeatCapPct: 0.015,
				AddRiskDollars:   75,
				AddBucket:        "Tech",
			},
			wantErr: "heat_cap_pct must be between 0 and 1",
		},
		{
			name: "Invalid bucket cap percent",
			req: HeatRequest{
				Equity:           10000,
				HeatCapPct:       0.04,
				BucketHeatCapPct: 1.5,
				AddRiskDollars:   75,
				AddBucket:        "Tech",
			},
			wantErr: "bucket_heat_cap_pct must be between 0 and 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateHeat(tt.req)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestCalculateHeat_Percentages(t *testing.T) {
	req := HeatRequest{
		Equity:           10000,
		HeatCapPct:       0.04,   // $400 cap
		BucketHeatCapPct: 0.015,  // $150 cap
		AddRiskDollars:   75,
		AddBucket:        "Tech/Comm",
		OpenPositions: []Position{
			{Ticker: "AAPL", Bucket: "Tech/Comm", RiskDollars: 75, Status: "Open"},
		},
	}

	result, err := CalculateHeat(req)
	require.NoError(t, err)

	// NewPortfolioHeat = 75 + 75 = 150
	// PortfolioHeatPct = (150 / 400) * 100 = 37.5%
	assert.InDelta(t, 37.5, result.PortfolioHeatPct, 0.01)

	// NewBucketHeat = 75 + 75 = 150
	// BucketHeatPct = (150 / 150) * 100 = 100%
	assert.InDelta(t, 100.0, result.BucketHeatPct, 0.01)
}
