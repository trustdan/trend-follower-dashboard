package domain

import (
	"fmt"
)

// HeatRequest contains input for heat calculation
type HeatRequest struct {
	Equity           float64    `json:"equity"`
	HeatCapPct       float64    `json:"heat_cap_pct"`
	BucketHeatCapPct float64    `json:"bucket_heat_cap_pct"`
	AddRiskDollars   float64    `json:"add_risk_dollars"` // New trade risk
	AddBucket        string     `json:"add_bucket"`        // Bucket for new trade
	OpenPositions    []Position `json:"open_positions"`
}

// Position represents an open trading position
type Position struct {
	Ticker      string  `json:"ticker"`
	Bucket      string  `json:"bucket"`
	RiskDollars float64 `json:"risk_dollars"`
	UnitsOpen   int     `json:"units_open"`
	Status      string  `json:"status"` // "Open", "Closed"
}

// HeatResult contains heat calculation output
type HeatResult struct {
	// Portfolio level
	CurrentPortfolioHeat float64 `json:"current_portfolio_heat"`
	NewPortfolioHeat     float64 `json:"new_portfolio_heat"`
	PortfolioHeatPct     float64 `json:"portfolio_heat_pct"`
	PortfolioCap         float64 `json:"portfolio_cap"`
	PortfolioCapExceeded bool    `json:"portfolio_cap_exceeded"`
	PortfolioOverage     float64 `json:"portfolio_overage"`

	// Bucket level
	CurrentBucketHeat float64 `json:"current_bucket_heat"`
	NewBucketHeat     float64 `json:"new_bucket_heat"`
	BucketHeatPct     float64 `json:"bucket_heat_pct"`
	BucketCap         float64 `json:"bucket_cap"`
	BucketCapExceeded bool    `json:"bucket_cap_exceeded"`
	BucketOverage     float64 `json:"bucket_overage"`

	// Decision
	Allowed         bool   `json:"allowed"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

// CalculateHeat computes portfolio and bucket heat
//
// Heat Management Rules (from CLAUDE.md):
//   - Portfolio heat = sum of risk across all open positions
//   - Portfolio cap = Equity × HeatCapPct (4% = $400 for $10k account)
//   - Bucket heat = sum of risk within one sector bucket
//   - Bucket cap = Equity × BucketHeatCapPct (1.5% = $150 for $10k account)
//   - Any trade exceeding either cap is rejected
//
// Example:
//   Given: Equity=$10,000, PortfolioHeat=$350, TechBucketHeat=$125, NewTrade=$75
//   NewPortfolioHeat = $350 + $75 = $425 > $400 (cap) → REJECT
func CalculateHeat(req HeatRequest) (*HeatResult, error) {
	// Validation
	if req.Equity <= 0 {
		return nil, ErrInvalidEquity
	}
	if req.HeatCapPct <= 0 || req.HeatCapPct > 1 {
		return nil, fmt.Errorf("heat_cap_pct must be between 0 and 1, got %.4f", req.HeatCapPct)
	}
	if req.BucketHeatCapPct <= 0 || req.BucketHeatCapPct > 1 {
		return nil, fmt.Errorf("bucket_heat_cap_pct must be between 0 and 1, got %.4f", req.BucketHeatCapPct)
	}
	if req.AddRiskDollars < 0 {
		return nil, fmt.Errorf("add_risk_dollars must be non-negative, got %.2f", req.AddRiskDollars)
	}

	// Calculate caps
	portfolioCap := req.Equity * req.HeatCapPct
	bucketCap := req.Equity * req.BucketHeatCapPct

	// Sum current portfolio heat (all open positions)
	currentPortfolioHeat := 0.0
	for _, pos := range req.OpenPositions {
		if pos.Status == "Open" {
			currentPortfolioHeat += pos.RiskDollars
		}
	}

	// Sum current bucket heat (positions in target bucket)
	currentBucketHeat := 0.0
	for _, pos := range req.OpenPositions {
		if pos.Status == "Open" && pos.Bucket == req.AddBucket {
			currentBucketHeat += pos.RiskDollars
		}
	}

	// Calculate new heat levels
	newPortfolioHeat := currentPortfolioHeat + req.AddRiskDollars
	newBucketHeat := currentBucketHeat + req.AddRiskDollars

	// Calculate percentages
	portfolioHeatPct := (newPortfolioHeat / portfolioCap) * 100
	bucketHeatPct := (newBucketHeat / bucketCap) * 100

	// Check caps
	portfolioCapExceeded := newPortfolioHeat > portfolioCap
	bucketCapExceeded := newBucketHeat > bucketCap

	portfolioOverage := 0.0
	if portfolioCapExceeded {
		portfolioOverage = newPortfolioHeat - portfolioCap
	}

	bucketOverage := 0.0
	if bucketCapExceeded {
		bucketOverage = newBucketHeat - bucketCap
	}

	// Determine if trade is allowed
	allowed := !portfolioCapExceeded && !bucketCapExceeded
	rejectionReason := ""

	if portfolioCapExceeded {
		rejectionReason = fmt.Sprintf(
			"Portfolio heat ($%.2f) exceeds cap ($%.2f) by $%.2f",
			newPortfolioHeat, portfolioCap, portfolioOverage,
		)
	} else if bucketCapExceeded {
		rejectionReason = fmt.Sprintf(
			"Bucket '%s' heat ($%.2f) exceeds cap ($%.2f) by $%.2f",
			req.AddBucket, newBucketHeat, bucketCap, bucketOverage,
		)
	}

	return &HeatResult{
		CurrentPortfolioHeat: currentPortfolioHeat,
		NewPortfolioHeat:     newPortfolioHeat,
		PortfolioHeatPct:     portfolioHeatPct,
		PortfolioCap:         portfolioCap,
		PortfolioCapExceeded: portfolioCapExceeded,
		PortfolioOverage:     portfolioOverage,
		CurrentBucketHeat:    currentBucketHeat,
		NewBucketHeat:        newBucketHeat,
		BucketHeatPct:        bucketHeatPct,
		BucketCap:            bucketCap,
		BucketCapExceeded:    bucketCapExceeded,
		BucketOverage:        bucketOverage,
		Allowed:              allowed,
		RejectionReason:      rejectionReason,
	}, nil
}
