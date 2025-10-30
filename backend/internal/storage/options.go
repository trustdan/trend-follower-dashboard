package storage

import (
	"fmt"
	"math"
)

// OptionLeg represents a single leg in a multi-leg options strategy
type OptionLeg struct {
	Type   string  `json:"type"`   // CALL or PUT
	Strike float64 `json:"strike"` // Strike price
	Exp    string  `json:"exp"`    // Expiration date (YYYY-MM-DD)
	Qty    int     `json:"qty"`    // Number of contracts
	Action string  `json:"action"` // BUY or SELL
	Price  float64 `json:"price"`  // Price per contract (optional, for record keeping)
}

// Options strategy type constants (26 total strategies)
const (
	// Directional (Simple)
	StrategyLongCall = "LONG_CALL"
	StrategyLongPut  = "LONG_PUT"

	// Income Generation
	StrategyCoveredCall      = "COVERED_CALL"
	StrategyCashSecuredPut   = "CASH_SECURED_PUT"

	// Vertical Credit Spreads
	StrategyBullPutSpread  = "BULL_PUT_SPREAD"
	StrategyBearCallSpread = "BEAR_CALL_SPREAD"

	// Butterflies & Condors (Neutral)
	StrategyIronButterfly        = "IRON_BUTTERFLY"
	StrategyIronCondor           = "IRON_CONDOR"
	StrategyLongPutButterfly     = "LONG_PUT_BUTTERFLY"
	StrategyLongCallButterfly    = "LONG_CALL_BUTTERFLY"
	StrategyInverseIronButterfly = "INVERSE_IRON_BUTTERFLY"
	StrategyInverseIronCondor    = "INVERSE_IRON_CONDOR"
	StrategyShortPutButterfly    = "SHORT_PUT_BUTTERFLY"
	StrategyShortCallButterfly   = "SHORT_CALL_BUTTERFLY"

	// Time Spreads (Calendar & Diagonal)
	StrategyCalendarCallSpread = "CALENDAR_CALL_SPREAD"
	StrategyCalendarPutSpread  = "CALENDAR_PUT_SPREAD"
	StrategyDiagonalCallSpread = "DIAGONAL_CALL_SPREAD"
	StrategyDiagonalPutSpread  = "DIAGONAL_PUT_SPREAD"

	// Vertical Debit Spreads
	StrategyBullCallSpread = "BULL_CALL_SPREAD"
	StrategyBearPutSpread  = "BEAR_PUT_SPREAD"

	// Volatility Plays
	StrategyStraddle = "STRADDLE"
	StrategyStrangle = "STRANGLE"

	// Ratio & Broken Wing
	StrategyCallRatioBackspread = "CALL_RATIO_BACKSPREAD"
	StrategyPutBrokenWing       = "PUT_BROKEN_WING"
	StrategyPutRatioBackspread  = "PUT_RATIO_BACKSPREAD"
	StrategyCallBrokenWing      = "CALL_BROKEN_WING"
)

// Breakout system constants
const (
	SystemOne    = "SYSTEM_1" // 20-bar breakout, 10-bar exit
	SystemTwo    = "SYSTEM_2" // 55-bar breakout, 10-bar exit (DEFAULT)
	SystemCustom = "CUSTOM"   // Manual parameters
)

// Time exit mode constants
const (
	TimeExitNone  = "None"  // Do not exit on time
	TimeExitClose = "Close" // Close position at roll threshold
	TimeExitRoll  = "Roll"  // Roll to next expiration
)

// Instrument type constants
const (
	InstrumentStock  = "STOCK"
	InstrumentOption = "OPTION"
)

// OptionsMetadata groups all options-related fields for a session or position
type OptionsMetadata struct {
	InstrumentType        string       `json:"instrument_type"`
	OptionsStrategy       string       `json:"options_strategy,omitempty"`
	EntryDate             string       `json:"entry_date,omitempty"`              // YYYY-MM-DD
	PrimaryExpirationDate string       `json:"primary_expiration_date,omitempty"` // YYYY-MM-DD
	DTE                   int          `json:"dte,omitempty"`                     // Days to expiration at entry
	RollThresholdDTE      int          `json:"roll_threshold_dte,omitempty"`      // DTE to roll/close (default 21)
	TimeExitMode          string       `json:"time_exit_mode,omitempty"`          // None, Close, Roll
	Legs                  []OptionLeg  `json:"legs,omitempty"`
	NetDebit              float64      `json:"net_debit,omitempty"`           // Total debit paid (negative = credit)
	MaxProfit             float64      `json:"max_profit,omitempty"`          // Maximum theoretical profit
	MaxLoss               float64      `json:"max_loss,omitempty"`            // Maximum theoretical loss
	BreakevenLower        float64      `json:"breakeven_lower,omitempty"`     // Lower breakeven price
	BreakevenUpper        float64      `json:"breakeven_upper,omitempty"`     // Upper breakeven price
	UnderlyingAtEntry     float64      `json:"underlying_at_entry,omitempty"` // Stock price at entry
}

// PyramidMetadata groups all pyramiding-related fields
type PyramidMetadata struct {
	MaxUnits     int     `json:"max_units"`               // Maximum pyramid units (default 4)
	AddStepN     float64 `json:"add_step_n"`              // Add every X * N (default 0.5)
	CurrentUnits int     `json:"current_units"`           // Current units in position (0-4)
	AddPrice1    float64 `json:"add_price_1,omitempty"`   // Entry + 0.5N
	AddPrice2    float64 `json:"add_price_2,omitempty"`   // Entry + 1.0N
	AddPrice3    float64 `json:"add_price_3,omitempty"`   // Entry + 1.5N
}

// BreakoutSystemMetadata groups breakout system parameters
type BreakoutSystemMetadata struct {
	EntryLookback int `json:"entry_lookback,omitempty"` // 20 or 55 for System-1/System-2
	ExitLookback  int `json:"exit_lookback,omitempty"`  // 10-bar exit (default)
}

// GetStrategyDisplayName returns a human-readable name for an options strategy constant
func GetStrategyDisplayName(strategy string) string {
	names := map[string]string{
		StrategyLongCall:             "Long Call",
		StrategyLongPut:              "Long Put",
		StrategyCoveredCall:          "Covered Call",
		StrategyCashSecuredPut:       "Cash-Secured Put",
		StrategyBullPutSpread:        "Bull Put Credit Spread",
		StrategyBearCallSpread:       "Bear Call Credit Spread",
		StrategyIronButterfly:        "Iron Butterfly",
		StrategyIronCondor:           "Iron Condor",
		StrategyLongPutButterfly:     "Long Put Butterfly",
		StrategyLongCallButterfly:    "Long Call Butterfly",
		StrategyInverseIronButterfly: "Inverse Iron Butterfly",
		StrategyInverseIronCondor:    "Inverse Iron Condor",
		StrategyShortPutButterfly:    "Short Put Butterfly",
		StrategyShortCallButterfly:   "Short Call Butterfly",
		StrategyCalendarCallSpread:   "Calendar Call Spread",
		StrategyCalendarPutSpread:    "Calendar Put Spread",
		StrategyDiagonalCallSpread:   "Diagonal Call Spread",
		StrategyDiagonalPutSpread:    "Diagonal Put Spread",
		StrategyBullCallSpread:       "Bull Call Spread",
		StrategyBearPutSpread:        "Bear Put Spread",
		StrategyStraddle:             "Straddle",
		StrategyStrangle:             "Strangle",
		StrategyCallRatioBackspread:  "Call Ratio Backspread",
		StrategyPutBrokenWing:        "Put Broken Wing",
		StrategyPutRatioBackspread:   "Put Ratio Backspread",
		StrategyCallBrokenWing:       "Call Broken Wing",
	}

	if name, ok := names[strategy]; ok {
		return name
	}
	return strategy
}

// GetStrategyCategory returns the category for grouping strategies in UI
func GetStrategyCategory(strategy string) string {
	switch strategy {
	case StrategyLongCall, StrategyLongPut:
		return "Directional"
	case StrategyCoveredCall, StrategyCashSecuredPut:
		return "Income"
	case StrategyBullPutSpread, StrategyBearCallSpread:
		return "Vertical Credit Spreads"
	case StrategyIronButterfly, StrategyIronCondor, StrategyLongPutButterfly,
		StrategyLongCallButterfly, StrategyInverseIronButterfly,
		StrategyInverseIronCondor, StrategyShortPutButterfly, StrategyShortCallButterfly:
		return "Butterflies & Condors"
	case StrategyCalendarCallSpread, StrategyCalendarPutSpread,
		StrategyDiagonalCallSpread, StrategyDiagonalPutSpread:
		return "Time Spreads"
	case StrategyBullCallSpread, StrategyBearPutSpread:
		return "Vertical Debit Spreads"
	case StrategyStraddle, StrategyStrangle:
		return "Volatility"
	case StrategyCallRatioBackspread, StrategyPutBrokenWing,
		StrategyPutRatioBackspread, StrategyCallBrokenWing:
		return "Ratio & Broken Wing"
	default:
		return "Other"
	}
}

// ========================================
// Strategy Builder Functions
// ========================================

// BuildLongCall creates a single-leg long call structure
func BuildLongCall(strike float64, expiration string, contracts int, premium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "CALL",
			Strike: strike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  premium,
		},
	}
}

// BuildLongPut creates a single-leg long put structure
func BuildLongPut(strike float64, expiration string, contracts int, premium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "PUT",
			Strike: strike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  premium,
		},
	}
}

// BuildBullCallSpread creates a 2-leg bull call spread (buy lower, sell higher)
func BuildBullCallSpread(lowerStrike, upperStrike float64, expiration string, contracts int, lowerPremium, upperPremium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "CALL",
			Strike: lowerStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  lowerPremium,
		},
		{
			Type:   "CALL",
			Strike: upperStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  upperPremium,
		},
	}
}

// BuildBearPutSpread creates a 2-leg bear put spread (buy higher, sell lower)
func BuildBearPutSpread(lowerStrike, upperStrike float64, expiration string, contracts int, lowerPremium, upperPremium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "PUT",
			Strike: upperStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  upperPremium,
		},
		{
			Type:   "PUT",
			Strike: lowerStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  lowerPremium,
		},
	}
}

// BuildBullPutSpread creates a 2-leg bull put credit spread (sell higher, buy lower)
func BuildBullPutSpread(lowerStrike, upperStrike float64, expiration string, contracts int, lowerPremium, upperPremium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "PUT",
			Strike: lowerStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  lowerPremium,
		},
		{
			Type:   "PUT",
			Strike: upperStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  upperPremium,
		},
	}
}

// BuildBearCallSpread creates a 2-leg bear call credit spread (sell lower, buy higher)
func BuildBearCallSpread(lowerStrike, upperStrike float64, expiration string, contracts int, lowerPremium, upperPremium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "CALL",
			Strike: lowerStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  lowerPremium,
		},
		{
			Type:   "CALL",
			Strike: upperStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  upperPremium,
		},
	}
}

// BuildIronCondor creates a 4-leg iron condor structure
// underlying: current stock price
// putSpreadWidth: width of put spread (e.g., $5)
// callSpreadWidth: width of call spread (e.g., $5)
// wingDistance: distance from underlying to short strikes (e.g., $10)
func BuildIronCondor(underlying, putSpreadWidth, callSpreadWidth, wingDistance float64, expiration string, contracts int,
	buyPutPremium, sellPutPremium, sellCallPremium, buyCallPremium float64) []OptionLeg {

	// Calculate strikes
	sellPutStrike := underlying - wingDistance
	buyPutStrike := sellPutStrike - putSpreadWidth
	sellCallStrike := underlying + wingDistance
	buyCallStrike := sellCallStrike + callSpreadWidth

	return []OptionLeg{
		{
			Type:   "PUT",
			Strike: buyPutStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  buyPutPremium,
		},
		{
			Type:   "PUT",
			Strike: sellPutStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  sellPutPremium,
		},
		{
			Type:   "CALL",
			Strike: sellCallStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  sellCallPremium,
		},
		{
			Type:   "CALL",
			Strike: buyCallStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  buyCallPremium,
		},
	}
}

// BuildIronButterfly creates a 4-leg iron butterfly structure
// atmStrike: at-the-money strike
// wingWidth: distance from ATM to wings
func BuildIronButterfly(atmStrike, wingWidth float64, expiration string, contracts int,
	buyPutPremium, sellPutPremium, sellCallPremium, buyCallPremium float64) []OptionLeg {

	lowerWing := atmStrike - wingWidth
	upperWing := atmStrike + wingWidth

	return []OptionLeg{
		{
			Type:   "PUT",
			Strike: lowerWing,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  buyPutPremium,
		},
		{
			Type:   "PUT",
			Strike: atmStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  sellPutPremium,
		},
		{
			Type:   "CALL",
			Strike: atmStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  sellCallPremium,
		},
		{
			Type:   "CALL",
			Strike: upperWing,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  buyCallPremium,
		},
	}
}

// BuildStraddle creates a 2-leg straddle (buy ATM call + ATM put)
func BuildStraddle(atmStrike float64, expiration string, contracts int, callPremium, putPremium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "CALL",
			Strike: atmStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  callPremium,
		},
		{
			Type:   "PUT",
			Strike: atmStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  putPremium,
		},
	}
}

// BuildStrangle creates a 2-leg strangle (buy OTM call + OTM put)
func BuildStrangle(callStrike, putStrike float64, expiration string, contracts int, callPremium, putPremium float64) []OptionLeg {
	return []OptionLeg{
		{
			Type:   "CALL",
			Strike: callStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  callPremium,
		},
		{
			Type:   "PUT",
			Strike: putStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  putPremium,
		},
	}
}

// BuildCalendarSpread creates a 2-leg calendar spread (same strike, different expirations)
func BuildCalendarSpread(optionType string, strike float64, nearExpiration, farExpiration string, contracts int,
	nearPremium, farPremium float64) []OptionLeg {

	return []OptionLeg{
		{
			Type:   optionType,
			Strike: strike,
			Exp:    nearExpiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  nearPremium,
		},
		{
			Type:   optionType,
			Strike: strike,
			Exp:    farExpiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  farPremium,
		},
	}
}

// BuildDiagonalSpread creates a 2-leg diagonal spread (different strikes AND expirations)
func BuildDiagonalSpread(optionType string, nearStrike, farStrike float64, nearExpiration, farExpiration string, contracts int,
	nearPremium, farPremium float64) []OptionLeg {

	return []OptionLeg{
		{
			Type:   optionType,
			Strike: nearStrike,
			Exp:    nearExpiration,
			Qty:    contracts,
			Action: "SELL",
			Price:  nearPremium,
		},
		{
			Type:   optionType,
			Strike: farStrike,
			Exp:    farExpiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  farPremium,
		},
	}
}

// BuildButterfly creates a 3-leg butterfly (buy 1 lower, sell 2 middle, buy 1 upper)
func BuildButterfly(optionType string, lowerStrike, middleStrike, upperStrike float64, expiration string, contracts int,
	lowerPremium, middlePremium, upperPremium float64) []OptionLeg {

	return []OptionLeg{
		{
			Type:   optionType,
			Strike: lowerStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  lowerPremium,
		},
		{
			Type:   optionType,
			Strike: middleStrike,
			Exp:    expiration,
			Qty:    contracts * 2, // Sell 2x middle
			Action: "SELL",
			Price:  middlePremium,
		},
		{
			Type:   optionType,
			Strike: upperStrike,
			Exp:    expiration,
			Qty:    contracts,
			Action: "BUY",
			Price:  upperPremium,
		},
	}
}

// ========================================
// Calculation Functions
// ========================================

// CalculateNetDebit computes net debit/credit from legs (negative = credit received)
func CalculateNetDebit(legs []OptionLeg) float64 {
	netDebit := 0.0
	for _, leg := range legs {
		if leg.Action == "BUY" {
			netDebit += leg.Price * float64(leg.Qty) * 100 // Options are per 100 shares
		} else if leg.Action == "SELL" {
			netDebit -= leg.Price * float64(leg.Qty) * 100
		}
	}
	return netDebit
}

// CalculateMaxProfitLoss computes max profit and max loss for common strategies
func CalculateMaxProfitLoss(strategyType string, legs []OptionLeg) (maxProfit, maxLoss float64, err error) {
	netDebit := CalculateNetDebit(legs)

	switch strategyType {
	case StrategyLongCall, StrategyLongPut:
		// Long single option: max loss = debit paid, max profit = unlimited (use large number)
		return 999999.99, netDebit, nil

	case StrategyBullCallSpread, StrategyBearPutSpread:
		// Vertical debit spread: max profit = spread width - debit, max loss = debit
		if len(legs) != 2 {
			return 0, 0, fmt.Errorf("vertical spread requires exactly 2 legs, got %d", len(legs))
		}
		spreadWidth := math.Abs(legs[0].Strike-legs[1].Strike) * 100
		maxProfit := spreadWidth - netDebit
		maxLoss := netDebit
		return maxProfit, maxLoss, nil

	case StrategyBullPutSpread, StrategyBearCallSpread:
		// Vertical credit spread: max profit = credit received, max loss = spread width - credit
		if len(legs) != 2 {
			return 0, 0, fmt.Errorf("vertical spread requires exactly 2 legs, got %d", len(legs))
		}
		spreadWidth := math.Abs(legs[0].Strike-legs[1].Strike) * 100
		maxProfit := -netDebit // Credit received (netDebit is negative)
		maxLoss := spreadWidth - maxProfit
		return maxProfit, maxLoss, nil

	case StrategyIronCondor:
		// Iron Condor: max profit = credit, max loss = wider spread width - credit
		if len(legs) != 4 {
			return 0, 0, fmt.Errorf("iron condor requires exactly 4 legs, got %d", len(legs))
		}
		// Find the width of the wider spread
		putSpreadWidth := math.Abs(legs[0].Strike-legs[1].Strike) * 100
		callSpreadWidth := math.Abs(legs[2].Strike-legs[3].Strike) * 100
		widerSpread := math.Max(putSpreadWidth, callSpreadWidth)
		maxProfit := -netDebit
		maxLoss := widerSpread - maxProfit
		return maxProfit, maxLoss, nil

	case StrategyIronButterfly:
		// Iron Butterfly: max profit = credit, max loss = wing width - credit
		if len(legs) != 4 {
			return 0, 0, fmt.Errorf("iron butterfly requires exactly 4 legs, got %d", len(legs))
		}
		wingWidth := math.Abs(legs[0].Strike-legs[1].Strike) * 100
		maxProfit := -netDebit
		maxLoss := wingWidth - maxProfit
		return maxProfit, maxLoss, nil

	case StrategyStraddle, StrategyStrangle:
		// Straddle/Strangle: max loss = debit, max profit = unlimited
		return 999999.99, netDebit, nil

	default:
		// For other strategies, return conservative estimates
		return 0, netDebit, fmt.Errorf("max profit/loss calculation not implemented for strategy %s", strategyType)
	}
}

// CalculateBreakevens computes breakeven prices for common strategies
func CalculateBreakevens(strategyType string, legs []OptionLeg) (lower, upper float64, err error) {
	netDebit := CalculateNetDebit(legs)
	debitPerContract := netDebit / 100 // Convert to per-share basis

	switch strategyType {
	case StrategyLongCall:
		// Breakeven = strike + debit paid
		be := legs[0].Strike + debitPerContract
		return 0, be, nil

	case StrategyLongPut:
		// Breakeven = strike - debit paid
		be := legs[0].Strike - debitPerContract
		return be, 0, nil

	case StrategyBullCallSpread:
		// Breakeven = lower strike + debit paid
		lowerStrike := math.Min(legs[0].Strike, legs[1].Strike)
		be := lowerStrike + debitPerContract
		return 0, be, nil

	case StrategyBearPutSpread:
		// Breakeven = upper strike - debit paid
		upperStrike := math.Max(legs[0].Strike, legs[1].Strike)
		be := upperStrike - debitPerContract
		return be, 0, nil

	case StrategyBullPutSpread:
		// Breakeven = higher short strike - credit
		upperStrike := math.Max(legs[0].Strike, legs[1].Strike)
		be := upperStrike + debitPerContract // debitPerContract is negative (credit)
		return be, 0, nil

	case StrategyBearCallSpread:
		// Breakeven = lower short strike + credit
		lowerStrike := math.Min(legs[0].Strike, legs[1].Strike)
		be := lowerStrike + debitPerContract // debitPerContract is negative (credit)
		return 0, be, nil

	case StrategyIronCondor:
		// Two breakevens: put side and call side
		if len(legs) != 4 {
			return 0, 0, fmt.Errorf("iron condor requires exactly 4 legs")
		}
		creditPerShare := -debitPerContract
		// Find the short strikes
		var shortPutStrike, shortCallStrike float64
		for _, leg := range legs {
			if leg.Action == "SELL" && leg.Type == "PUT" {
				shortPutStrike = leg.Strike
			}
			if leg.Action == "SELL" && leg.Type == "CALL" {
				shortCallStrike = leg.Strike
			}
		}
		lowerBE := shortPutStrike - creditPerShare
		upperBE := shortCallStrike + creditPerShare
		return lowerBE, upperBE, nil

	case StrategyIronButterfly:
		// Two breakevens around ATM
		if len(legs) != 4 {
			return 0, 0, fmt.Errorf("iron butterfly requires exactly 4 legs")
		}
		creditPerShare := -debitPerContract
		// Find ATM strike (should be same for both short options)
		var atmStrike float64
		for _, leg := range legs {
			if leg.Action == "SELL" {
				atmStrike = leg.Strike
				break
			}
		}
		lowerBE := atmStrike - creditPerShare
		upperBE := atmStrike + creditPerShare
		return lowerBE, upperBE, nil

	case StrategyStraddle:
		// Two breakevens: strike ± debit
		strike := legs[0].Strike // ATM strike
		lowerBE := strike - debitPerContract
		upperBE := strike + debitPerContract
		return lowerBE, upperBE, nil

	case StrategyStrangle:
		// Two breakevens: put strike - debit, call strike + debit
		var putStrike, callStrike float64
		for _, leg := range legs {
			if leg.Type == "PUT" {
				putStrike = leg.Strike
			} else if leg.Type == "CALL" {
				callStrike = leg.Strike
			}
		}
		lowerBE := putStrike - debitPerContract
		upperBE := callStrike + debitPerContract
		return lowerBE, upperBE, nil

	default:
		return 0, 0, fmt.Errorf("breakeven calculation not implemented for strategy %s", strategyType)
	}
}
