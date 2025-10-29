package domain

import (
	"math"
)

// CalculateStockPosition implements Van Tharp position sizing for stocks
//
// The Van Tharp method calculates position size based on risk per trade:
// 1. Calculate risk dollars: R = Equity × RiskPct (0.75%)
// 2. Calculate stop distance: StopDist = K × ATR (K=2)
// 3. Calculate initial stop: InitStop = Entry - StopDist
// 4. Calculate shares: Shares = floor(R ÷ StopDist)
// 5. Verify: ActualRisk = Shares × StopDist ≤ R
//
// Example (from CLAUDE.md):
//   Given: Equity=$10,000, Risk=0.75%, Entry=$180, ATR=$1.50, K=2
//   R = $10,000 × 0.0075 = $75.00
//   StopDistance = 2 × $1.50 = $3.00
//   InitialStop = $180 - $3.00 = $177.00
//   Shares = floor($75 ÷ $3.00) = 25 shares
//   ActualRisk = 25 × $3.00 = $75.00 ✓
func CalculateStockPosition(req SizingRequest) (*SizingResult, error) {
	// Validation
	if err := validateCommonInputs(req); err != nil {
		return nil, err
	}

	// Step 1: Calculate risk dollars (R)
	// R = Equity × RiskPct
	riskDollars := req.Equity * req.RiskPct

	// Step 2: Calculate stop distance
	// StopDistance = K × ATR
	stopDistance := float64(req.K) * req.ATR

	// Step 3: Calculate initial stop
	// InitialStop = Entry - StopDistance
	initialStop := req.Entry - stopDistance

	// Step 4: Calculate shares (round down)
	// Shares = floor(R ÷ StopDistance)
	shares := int(math.Floor(riskDollars / stopDistance))

	// Step 5: Calculate actual risk
	// ActualRisk = Shares × StopDistance
	actualRisk := float64(shares) * stopDistance

	return &SizingResult{
		RiskDollars:  riskDollars,
		StopDistance: stopDistance,
		InitialStop:  initialStop,
		Shares:       shares,
		Contracts:    0, // Stocks don't use contracts
		ActualRisk:   actualRisk,
		Method:       "stock",
	}, nil
}
