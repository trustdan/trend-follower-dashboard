package domain

import (
	"fmt"
	"math"
)

// CalculateOptionDeltaATRPosition implements position sizing for options using delta-adjusted ATR
//
// The Delta-ATR method sizes options positions to risk the same dollar amount
// as an equivalent stock position, adjusted for the option's delta.
//
// Formula:
//   1. Stock equivalent: shares = floor(R ÷ (K × ATR))
//   2. Delta adjustment: effectiveShares = shares × delta
//   3. Contracts = floor(effectiveShares ÷ 100)
//   4. Actual risk = contracts × 100 × (K × ATR × delta)
//
// Why: Options with 0.30 delta move ~$0.30 for every $1.00 the stock moves.
// So a $3.00 stock stop becomes a $0.90 option stop.
//
// Example:
//   Given: R=$75, Entry=$180, ATR=$1.50, K=2, Delta=0.30
//   StopDistance = 2 × $1.50 = $3.00
//   StockShares = floor($75 ÷ $3.00) = 25
//   DeltaShares = 25 × 0.30 = 7.5
//   Contracts = floor(7.5 ÷ 100) = 0
//   Result: 0 contracts (delta too low, position too small)
func CalculateOptionDeltaATRPosition(req SizingRequest) (*SizingResult, error) {
	// Validation
	if err := validateCommonInputs(req); err != nil {
		return nil, err
	}

	// Validate delta
	if req.Delta <= 0 || req.Delta > 1 {
		return nil, fmt.Errorf("%w: got %.2f", ErrInvalidDelta, req.Delta)
	}

	// Step 1: Calculate risk dollars (R)
	riskDollars := req.Equity * req.RiskPct

	// Step 2: Calculate stop distance
	stopDistance := float64(req.K) * req.ATR

	// Step 3: Calculate initial stop (for underlying)
	initialStop := req.Entry - stopDistance

	// Step 4: Calculate stock-equivalent shares
	stockShares := math.Floor(riskDollars / stopDistance)

	// Step 5: Adjust for delta
	deltaShares := stockShares * req.Delta

	// Step 6: Convert to contracts (100 shares per contract)
	contracts := int(math.Floor(deltaShares / 100))

	// Step 7: Calculate actual risk
	// ActualRisk = contracts × 100 × (K × ATR × delta)
	actualRisk := float64(contracts) * 100 * stopDistance * req.Delta

	return &SizingResult{
		RiskDollars:  riskDollars,
		StopDistance: stopDistance,
		InitialStop:  initialStop,
		Shares:       0,        // Options use contracts, not shares
		Contracts:    contracts,
		ActualRisk:   actualRisk,
		Method:       "opt-delta-atr",
	}, nil
}

// CalculateOptionMaxLossPosition implements position sizing for options using max loss method
//
// The Max-Loss method sizes options positions based on the maximum loss per contract.
// This is ideal for defined-risk strategies where you know the maximum loss upfront.
//
// Formula:
//   1. Calculate risk budget: R = Equity × RiskPct
//   2. Calculate contracts: contracts = floor(R ÷ MaxLoss)
//   3. Calculate actual risk: actualRisk = contracts × MaxLoss
//
// Example:
//   Given: Equity=$10,000, RiskPct=0.75%, MaxLoss=$50
//   R = $10,000 × 0.0075 = $75
//   Contracts = floor($75 ÷ $50) = 1
//   ActualRisk = 1 × $50 = $50
//
// When to use:
//   - Trading defined-risk strategies (spreads, buying options)
//   - You know maximum loss per contract upfront
//   - Trading small accounts where delta-ATR yields 0 contracts
//
// How to determine MaxLoss:
//   - Buying calls/puts: Premium paid per contract (e.g., $3.50 × 100 = $350)
//   - Credit spreads: Max loss = (width of spread - credit) × 100
//   - Debit spreads: Max loss = premium paid
func CalculateOptionMaxLossPosition(req SizingRequest) (*SizingResult, error) {
	// Validate equity and risk percent
	if req.Equity <= 0 {
		return nil, ErrInvalidEquity
	}
	if req.RiskPct <= 0 || req.RiskPct > 1 {
		return nil, fmt.Errorf("%w: got %.4f", ErrInvalidRiskPct, req.RiskPct)
	}

	// Validate max loss
	if req.MaxLoss <= 0 {
		return nil, fmt.Errorf("%w: got %.2f", ErrInvalidMaxLoss, req.MaxLoss)
	}

	// Step 1: Calculate risk budget
	// R = Equity × RiskPct
	riskDollars := req.Equity * req.RiskPct

	// Step 2: Calculate contracts
	// contracts = floor(R ÷ MaxLoss)
	// Always round down to avoid exceeding risk budget
	contracts := int(math.Floor(riskDollars / req.MaxLoss))

	// Step 3: Calculate actual risk
	actualRisk := float64(contracts) * req.MaxLoss

	// Sanity check - this should never happen due to floor division
	if actualRisk > riskDollars {
		return nil, fmt.Errorf(
			"internal error: actual risk $%.2f exceeds specified risk $%.2f",
			actualRisk, riskDollars,
		)
	}

	return &SizingResult{
		RiskDollars:  riskDollars,
		StopDistance: 0, // Not applicable for max-loss method
		InitialStop:  0, // Not applicable for max-loss method
		Shares:       0, // Options use contracts, not shares
		Contracts:    contracts,
		ActualRisk:   actualRisk,
		Method:       "opt-maxloss",
	}, nil
}
