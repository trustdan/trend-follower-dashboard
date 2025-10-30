package storage

import (
	"testing"
)

// TestBuildLongCall tests the long call builder
func TestBuildLongCall(t *testing.T) {
	legs := BuildLongCall(180.0, "2025-12-19", 5, 3.50)

	if len(legs) != 1 {
		t.Errorf("Expected 1 leg, got %d", len(legs))
	}

	leg := legs[0]
	if leg.Type != "CALL" {
		t.Errorf("Expected Type CALL, got %s", leg.Type)
	}
	if leg.Strike != 180.0 {
		t.Errorf("Expected Strike 180.0, got %f", leg.Strike)
	}
	if leg.Action != "BUY" {
		t.Errorf("Expected Action BUY, got %s", leg.Action)
	}
	if leg.Qty != 5 {
		t.Errorf("Expected Qty 5, got %d", leg.Qty)
	}
}

// TestBuildLongPut tests the long put builder
func TestBuildLongPut(t *testing.T) {
	legs := BuildLongPut(175.0, "2025-12-19", 3, 2.75)

	if len(legs) != 1 {
		t.Errorf("Expected 1 leg, got %d", len(legs))
	}

	leg := legs[0]
	if leg.Type != "PUT" {
		t.Errorf("Expected Type PUT, got %s", leg.Type)
	}
	if leg.Strike != 175.0 {
		t.Errorf("Expected Strike 175.0, got %f", leg.Strike)
	}
	if leg.Action != "BUY" {
		t.Errorf("Expected Action BUY, got %s", leg.Action)
	}
}

// TestBuildBullCallSpread tests the bull call spread builder
func TestBuildBullCallSpread(t *testing.T) {
	legs := BuildBullCallSpread(180.0, 185.0, "2025-12-19", 2, 3.50, 1.25)

	if len(legs) != 2 {
		t.Errorf("Expected 2 legs, got %d", len(legs))
	}

	// First leg should be BUY at lower strike
	if legs[0].Action != "BUY" || legs[0].Strike != 180.0 {
		t.Errorf("First leg should be BUY at 180.0")
	}

	// Second leg should be SELL at higher strike
	if legs[1].Action != "SELL" || legs[1].Strike != 185.0 {
		t.Errorf("Second leg should be SELL at 185.0")
	}
}

// TestBuildBearPutSpread tests the bear put spread builder
func TestBuildBearPutSpread(t *testing.T) {
	legs := BuildBearPutSpread(175.0, 180.0, "2025-12-19", 2, 1.00, 3.25)

	if len(legs) != 2 {
		t.Errorf("Expected 2 legs, got %d", len(legs))
	}

	// First leg should be BUY at upper strike
	if legs[0].Action != "BUY" || legs[0].Strike != 180.0 {
		t.Errorf("First leg should be BUY at 180.0")
	}

	// Second leg should be SELL at lower strike
	if legs[1].Action != "SELL" || legs[1].Strike != 175.0 {
		t.Errorf("Second leg should be SELL at 175.0")
	}
}

// TestBuildIronCondor tests the iron condor builder
func TestBuildIronCondor(t *testing.T) {
	// AAPL @ $175, put spread $5 wide, call spread $5 wide, wings $10 from center
	legs := BuildIronCondor(175.0, 5.0, 5.0, 10.0, "2025-12-19", 1,
		0.50, 1.20, 1.30, 0.60)

	if len(legs) != 4 {
		t.Errorf("Expected 4 legs, got %d", len(legs))
	}

	// Leg 0: Buy Put at 160 (175 - 10 - 5)
	if legs[0].Type != "PUT" || legs[0].Action != "BUY" || legs[0].Strike != 160.0 {
		t.Errorf("Leg 0 should be BUY PUT at 160.0, got %s %s at %f",
			legs[0].Action, legs[0].Type, legs[0].Strike)
	}

	// Leg 1: Sell Put at 165 (175 - 10)
	if legs[1].Type != "PUT" || legs[1].Action != "SELL" || legs[1].Strike != 165.0 {
		t.Errorf("Leg 1 should be SELL PUT at 165.0, got %s %s at %f",
			legs[1].Action, legs[1].Type, legs[1].Strike)
	}

	// Leg 2: Sell Call at 185 (175 + 10)
	if legs[2].Type != "CALL" || legs[2].Action != "SELL" || legs[2].Strike != 185.0 {
		t.Errorf("Leg 2 should be SELL CALL at 185.0, got %s %s at %f",
			legs[2].Action, legs[2].Type, legs[2].Strike)
	}

	// Leg 3: Buy Call at 190 (175 + 10 + 5)
	if legs[3].Type != "CALL" || legs[3].Action != "BUY" || legs[3].Strike != 190.0 {
		t.Errorf("Leg 3 should be BUY CALL at 190.0, got %s %s at %f",
			legs[3].Action, legs[3].Type, legs[3].Strike)
	}
}

// TestBuildIronButterfly tests the iron butterfly builder
func TestBuildIronButterfly(t *testing.T) {
	legs := BuildIronButterfly(175.0, 10.0, "2025-12-19", 1,
		0.50, 2.50, 2.60, 0.60)

	if len(legs) != 4 {
		t.Errorf("Expected 4 legs, got %d", len(legs))
	}

	// Leg 0: Buy Put at 165 (175 - 10)
	if legs[0].Type != "PUT" || legs[0].Action != "BUY" || legs[0].Strike != 165.0 {
		t.Errorf("Leg 0 should be BUY PUT at 165.0")
	}

	// Legs 1 and 2 should both be at ATM (175) and SELL
	if legs[1].Strike != 175.0 || legs[1].Action != "SELL" {
		t.Errorf("Leg 1 should be SELL at 175.0")
	}
	if legs[2].Strike != 175.0 || legs[2].Action != "SELL" {
		t.Errorf("Leg 2 should be SELL at 175.0")
	}

	// Leg 3: Buy Call at 185 (175 + 10)
	if legs[3].Type != "CALL" || legs[3].Action != "BUY" || legs[3].Strike != 185.0 {
		t.Errorf("Leg 3 should be BUY CALL at 185.0")
	}
}

// TestBuildStraddle tests the straddle builder
func TestBuildStraddle(t *testing.T) {
	legs := BuildStraddle(175.0, "2025-12-19", 2, 3.50, 3.25)

	if len(legs) != 2 {
		t.Errorf("Expected 2 legs, got %d", len(legs))
	}

	// Both legs should be BUY at ATM strike
	if legs[0].Action != "BUY" || legs[0].Strike != 175.0 {
		t.Errorf("First leg should be BUY at 175.0")
	}
	if legs[1].Action != "BUY" || legs[1].Strike != 175.0 {
		t.Errorf("Second leg should be BUY at 175.0")
	}

	// One CALL, one PUT
	if !(legs[0].Type == "CALL" && legs[1].Type == "PUT") {
		t.Errorf("Should have one CALL and one PUT")
	}
}

// TestBuildStrangle tests the strangle builder
func TestBuildStrangle(t *testing.T) {
	legs := BuildStrangle(180.0, 170.0, "2025-12-19", 2, 2.50, 2.25)

	if len(legs) != 2 {
		t.Errorf("Expected 2 legs, got %d", len(legs))
	}

	// Both should be BUY
	if legs[0].Action != "BUY" || legs[1].Action != "BUY" {
		t.Errorf("Both legs should be BUY")
	}

	// CALL at 180, PUT at 170
	if legs[0].Type != "CALL" || legs[0].Strike != 180.0 {
		t.Errorf("First leg should be CALL at 180.0")
	}
	if legs[1].Type != "PUT" || legs[1].Strike != 170.0 {
		t.Errorf("Second leg should be PUT at 170.0")
	}
}

// TestBuildCalendarSpread tests the calendar spread builder
func TestBuildCalendarSpread(t *testing.T) {
	legs := BuildCalendarSpread("CALL", 175.0, "2025-11-15", "2025-12-19", 1, 2.00, 3.50)

	if len(legs) != 2 {
		t.Errorf("Expected 2 legs, got %d", len(legs))
	}

	// Near expiration should be SELL
	if legs[0].Action != "SELL" || legs[0].Exp != "2025-11-15" {
		t.Errorf("Near leg should be SELL at 2025-11-15")
	}

	// Far expiration should be BUY
	if legs[1].Action != "BUY" || legs[1].Exp != "2025-12-19" {
		t.Errorf("Far leg should be BUY at 2025-12-19")
	}

	// Both should be same strike
	if legs[0].Strike != 175.0 || legs[1].Strike != 175.0 {
		t.Errorf("Both legs should be at strike 175.0")
	}
}

// TestBuildButterfly tests the butterfly builder
func TestBuildButterfly(t *testing.T) {
	legs := BuildButterfly("CALL", 170.0, 175.0, 180.0, "2025-12-19", 1, 1.00, 3.00, 0.50)

	if len(legs) != 3 {
		t.Errorf("Expected 3 legs, got %d", len(legs))
	}

	// Leg 0: BUY 1 at lower strike
	if legs[0].Action != "BUY" || legs[0].Strike != 170.0 || legs[0].Qty != 1 {
		t.Errorf("Leg 0 should be BUY 1 at 170.0")
	}

	// Leg 1: SELL 2 at middle strike
	if legs[1].Action != "SELL" || legs[1].Strike != 175.0 || legs[1].Qty != 2 {
		t.Errorf("Leg 1 should be SELL 2 at 175.0, got %s %d at %f",
			legs[1].Action, legs[1].Qty, legs[1].Strike)
	}

	// Leg 2: BUY 1 at upper strike
	if legs[2].Action != "BUY" || legs[2].Strike != 180.0 || legs[2].Qty != 1 {
		t.Errorf("Leg 2 should be BUY 1 at 180.0")
	}
}

// TestCalculateNetDebit tests the net debit calculator
func TestCalculateNetDebit(t *testing.T) {
	// Test bull call spread: Buy 180 call @ $3.50, Sell 185 call @ $1.25
	// Net debit = (3.50 - 1.25) * 100 = $225
	legs := BuildBullCallSpread(180.0, 185.0, "2025-12-19", 1, 3.50, 1.25)
	netDebit := CalculateNetDebit(legs)

	expected := 225.0 // (3.50 - 1.25) * 100
	if netDebit != expected {
		t.Errorf("Expected net debit $%.2f, got $%.2f", expected, netDebit)
	}

	// Test iron condor (credit spread)
	// Sell Put @ $1.20, Buy Put @ $0.50, Sell Call @ $1.30, Buy Call @ $0.60
	// Credit = (1.20 - 0.50 + 1.30 - 0.60) * 100 = $140
	legs2 := BuildIronCondor(175.0, 5.0, 5.0, 10.0, "2025-12-19", 1,
		0.50, 1.20, 1.30, 0.60)
	netDebit2 := CalculateNetDebit(legs2)

	// Net debit should be negative (it's a credit)
	expectedCredit := -140.0
	if netDebit2 != expectedCredit {
		t.Errorf("Expected net credit $%.2f, got $%.2f", expectedCredit, netDebit2)
	}
}

// TestCalculateMaxProfitLoss tests max profit/loss calculations
func TestCalculateMaxProfitLoss(t *testing.T) {
	// Test bull call spread
	// Spread width = $5, Debit = $2.25, Max profit = $2.75, Max loss = $2.25
	legs := BuildBullCallSpread(180.0, 185.0, "2025-12-19", 1, 3.50, 1.25)
	maxProfit, maxLoss, err := CalculateMaxProfitLoss(StrategyBullCallSpread, legs)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedMaxProfit := 275.0  // ($5 - $2.25) * 100
	expectedMaxLoss := 225.0    // $2.25 * 100

	if maxProfit != expectedMaxProfit {
		t.Errorf("Expected max profit $%.2f, got $%.2f", expectedMaxProfit, maxProfit)
	}
	if maxLoss != expectedMaxLoss {
		t.Errorf("Expected max loss $%.2f, got $%.2f", expectedMaxLoss, maxLoss)
	}

	// Test iron condor
	legs2 := BuildIronCondor(175.0, 5.0, 5.0, 10.0, "2025-12-19", 1,
		0.50, 1.20, 1.30, 0.60)
	maxProfit2, maxLoss2, err2 := CalculateMaxProfitLoss(StrategyIronCondor, legs2)

	if err2 != nil {
		t.Fatalf("Unexpected error: %v", err2)
	}

	expectedMaxProfit2 := 140.0  // Credit received
	expectedMaxLoss2 := 360.0    // $5 spread - $1.40 credit = $3.60 * 100

	if maxProfit2 != expectedMaxProfit2 {
		t.Errorf("Expected max profit $%.2f, got $%.2f", expectedMaxProfit2, maxProfit2)
	}
	if maxLoss2 != expectedMaxLoss2 {
		t.Errorf("Expected max loss $%.2f, got $%.2f", expectedMaxLoss2, maxLoss2)
	}
}

// TestCalculateBreakevens tests breakeven calculations
func TestCalculateBreakevens(t *testing.T) {
	// Test long call
	// Strike $180, Premium $3.50, Breakeven = $183.50
	legs := BuildLongCall(180.0, "2025-12-19", 1, 3.50)
	lower, upper, err := CalculateBreakevens(StrategyLongCall, legs)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedBreakeven := 183.50
	if upper != expectedBreakeven {
		t.Errorf("Expected breakeven $%.2f, got $%.2f", expectedBreakeven, upper)
	}
	if lower != 0 {
		t.Errorf("Long call should have no lower breakeven, got $%.2f", lower)
	}

	// Test iron condor
	// Credit = $1.40, Short Put = $165, Short Call = $185
	// Lower BE = $165 - $1.40 = $163.60, Upper BE = $185 + $1.40 = $186.40
	legs2 := BuildIronCondor(175.0, 5.0, 5.0, 10.0, "2025-12-19", 1,
		0.50, 1.20, 1.30, 0.60)
	lower2, upper2, err2 := CalculateBreakevens(StrategyIronCondor, legs2)

	if err2 != nil {
		t.Fatalf("Unexpected error: %v", err2)
	}

	expectedLowerBE := 163.60
	expectedUpperBE := 186.40

	if lower2 != expectedLowerBE {
		t.Errorf("Expected lower breakeven $%.2f, got $%.2f", expectedLowerBE, lower2)
	}
	if upper2 != expectedUpperBE {
		t.Errorf("Expected upper breakeven $%.2f, got $%.2f", expectedUpperBE, upper2)
	}

	// Test straddle
	// ATM = $175, Total premium = $3.50 + $3.25 = $6.75
	// Lower BE = $175 - $6.75 = $168.25, Upper BE = $175 + $6.75 = $181.75
	legs3 := BuildStraddle(175.0, "2025-12-19", 1, 3.50, 3.25)
	lower3, upper3, err3 := CalculateBreakevens(StrategyStraddle, legs3)

	if err3 != nil {
		t.Fatalf("Unexpected error: %v", err3)
	}

	expectedLowerBE3 := 168.25
	expectedUpperBE3 := 181.75

	if lower3 != expectedLowerBE3 {
		t.Errorf("Expected lower breakeven $%.2f, got $%.2f", expectedLowerBE3, lower3)
	}
	if upper3 != expectedUpperBE3 {
		t.Errorf("Expected upper breakeven $%.2f, got $%.2f", expectedUpperBE3, upper3)
	}
}

// TestGetStrategyDisplayName tests display name lookup
func TestGetStrategyDisplayName(t *testing.T) {
	tests := []struct {
		constant    string
		displayName string
	}{
		{StrategyLongCall, "Long Call"},
		{StrategyIronCondor, "Iron Condor"},
		{StrategyBullCallSpread, "Bull Call Spread"},
		{StrategyStraddle, "Straddle"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		result := GetStrategyDisplayName(tt.constant)
		if result != tt.displayName {
			t.Errorf("GetStrategyDisplayName(%s) = %s, want %s",
				tt.constant, result, tt.displayName)
		}
	}
}

// TestGetStrategyCategory tests category grouping
func TestGetStrategyCategory(t *testing.T) {
	tests := []struct {
		strategy string
		category string
	}{
		{StrategyLongCall, "Directional"},
		{StrategyCoveredCall, "Income"},
		{StrategyBullPutSpread, "Vertical Credit Spreads"},
		{StrategyIronCondor, "Butterflies & Condors"},
		{StrategyCalendarCallSpread, "Time Spreads"},
		{StrategyBullCallSpread, "Vertical Debit Spreads"},
		{StrategyStraddle, "Volatility"},
		{StrategyCallRatioBackspread, "Ratio & Broken Wing"},
		{"UNKNOWN", "Other"},
	}

	for _, tt := range tests {
		result := GetStrategyCategory(tt.strategy)
		if result != tt.category {
			t.Errorf("GetStrategyCategory(%s) = %s, want %s",
				tt.strategy, result, tt.category)
		}
	}
}
