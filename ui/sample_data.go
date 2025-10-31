package main

import (
	"time"

	"github.com/yourusername/trading-engine/internal/storage"
)

// CreateSampleSession creates a sample trade session with realistic data
func CreateSampleSession() *storage.TradeSession {
	now := time.Now()
	return &storage.TradeSession{
		ID:                     999, // Use high ID to avoid conflicts
		SessionNum:             999,
		Ticker:                 "AAPL",
		Strategy:               storage.StrategyLongBreakout,
		CurrentStep:            storage.StepEntry,
		Status:                 "DRAFT",
		InstrumentType:         "stock",
		OptionsStrategy:        "",
		CreatedAt:              now.Add(-30 * time.Minute),
		UpdatedAt:              now,
		ChecklistCompleted:     true,
		ChecklistBanner:        "GREEN",
		ChecklistMissingCount:  0,
		ChecklistQualityScore:  3,
		SizingCompleted:        true,
		SizingMethod:           "stock",
		SizingEntryPrice:       181.50,
		SizingATR:              2.45,
		SizingKMultiple:        2.0,
		SizingShares:           306,
		SizingRiskDollars:      750.00,
		HeatCompleted:          true,
		HeatPortfolioCurrent:   1500.00,
		HeatPortfolioNew:       2250.00,
		HeatBucket:             "Tech/Comm",
		HeatBucketCurrent:      375.00,
		HeatBucketNew:          1125.00,
		EntryCompleted:         false,
		EntryDecision:          "",
		PrimaryExpirationDate:  "",
		DTE:                    0,
	}
}

// CreateSamplePositions creates sample open positions for the dashboard
func CreateSamplePositions() []storage.Position {
	now := time.Now()
	return []storage.Position{
		{
			ID:                1,
			Ticker:            "MSFT",
			InstrumentType:    "stock",
			OptionsStrategy:   "",
			Status:            "OPEN",
			EntryPrice:        380.00,
			CurrentStop:       374.20,
			InitialStop:       374.20,
			Shares:            129,
			RiskDollars:       748.20,
			Bucket:            "Tech/Comm",
			OpenedAt:          now.Add(-7 * 24 * time.Hour),
			DecisionID:        985,
		},
		{
			ID:                2,
			Ticker:            "XLE",
			InstrumentType:    "stock",
			OptionsStrategy:   "",
			Status:            "OPEN",
			EntryPrice:        88.75,
			CurrentStop:       83.85,
			InitialStop:       83.85,
			Shares:            153,
			RiskDollars:       749.70,
			Bucket:            "Energy",
			OpenedAt:          now.Add(-4 * 24 * time.Hour),
			DecisionID:        992,
		},
		{
			ID:                3,
			Ticker:            "GLD",
			InstrumentType:    "stock",
			OptionsStrategy:   "",
			Status:            "OPEN",
			EntryPrice:        184.20,
			CurrentStop:       180.50,
			InitialStop:       180.50,
			Shares:            203,
			RiskDollars:       750.10,
			Bucket:            "Commodities",
			OpenedAt:          now.Add(-2 * 24 * time.Hour),
			DecisionID:        996,
		},
	}
}

// CreateSampleCandidates creates sample candidate tickers for today
func CreateSampleCandidates() []storage.Candidate {
	now := time.Now()
	today := now.Format("2006-01-02")

	return []storage.Candidate{
		{
			ID:     1,
			Ticker: "AAPL",
			Date:   today,
		},
		{
			ID:     2,
			Ticker: "NVDA",
			Date:   today,
		},
		{
			ID:     3,
			Ticker: "TSLA",
			Date:   today,
		},
		{
			ID:     4,
			Ticker: "AMD",
			Date:   today,
		},
		{
			ID:     5,
			Ticker: "META",
			Date:   today,
		},
	}
}

// CreateSampleCalendarTrades creates sample trades for the calendar view
func CreateSampleCalendarTrades() []storage.TradeSession {
	now := time.Now()

	// Get the start of the current week (Sunday)
	weekday := int(now.Weekday())
	weekStart := now.AddDate(0, 0, -weekday)

	trades := []storage.TradeSession{}

	// Week 1 (current week) - 2 trades
	trades = append(trades, storage.TradeSession{
		ID:                 991,
		SessionNum:         991,
		Ticker:             "MSFT",
		Strategy:           storage.StrategyLongBreakout,
		InstrumentType:     "stock",
		Status:             "COMPLETED",
		EntryDecision:      "GO",
		CreatedAt:          weekStart.Add(1 * 24 * time.Hour),
		UpdatedAt:          weekStart.Add(1 * 24 * time.Hour),
		HeatBucket:         "Tech/Comm",
	})

	// Week 2 (1 week ago) - 3 trades
	week2Start := weekStart.AddDate(0, 0, -7)
	trades = append(trades, storage.TradeSession{
		ID:                 988,
		SessionNum:         988,
		Ticker:             "XLE",
		Strategy:           storage.StrategyLongBreakout,
		InstrumentType:     "stock",
		Status:             "COMPLETED",
		EntryDecision:      "GO",
		CreatedAt:          week2Start.Add(2 * 24 * time.Hour),
		UpdatedAt:          week2Start.Add(2 * 24 * time.Hour),
		HeatBucket:         "Energy",
	})
	trades = append(trades, storage.TradeSession{
		ID:                 989,
		SessionNum:         989,
		Ticker:             "GLD",
		Strategy:           storage.StrategyLongBreakout,
		InstrumentType:     "stock",
		Status:             "COMPLETED",
		EntryDecision:      "GO",
		CreatedAt:          week2Start.Add(4 * 24 * time.Hour),
		UpdatedAt:          week2Start.Add(4 * 24 * time.Hour),
		HeatBucket:         "Commodities",
	})

	// Week 3 (2 weeks ago) - 1 trade
	week3Start := weekStart.AddDate(0, 0, -14)
	trades = append(trades, storage.TradeSession{
		ID:                 985,
		SessionNum:         985,
		Ticker:             "GOOGL",
		Strategy:           storage.StrategyLongBreakout,
		InstrumentType:     "stock",
		Status:             "COMPLETED",
		EntryDecision:      "GO",
		CreatedAt:          week3Start.Add(3 * 24 * time.Hour),
		UpdatedAt:          week3Start.Add(3 * 24 * time.Hour),
		HeatBucket:         "Tech/Comm",
	})

	// Week 4 (3 weeks ago) - 2 trades
	week4Start := weekStart.AddDate(0, 0, -21)
	trades = append(trades, storage.TradeSession{
		ID:                 982,
		SessionNum:         982,
		Ticker:             "JPM",
		Strategy:           storage.StrategyLongBreakout,
		InstrumentType:     "stock",
		Status:             "COMPLETED",
		EntryDecision:      "GO",
		CreatedAt:          week4Start.Add(1 * 24 * time.Hour),
		UpdatedAt:          week4Start.Add(1 * 24 * time.Hour),
		HeatBucket:         "Financials",
	})
	trades = append(trades, storage.TradeSession{
		ID:                 983,
		SessionNum:         983,
		Ticker:             "CAT",
		Strategy:           storage.StrategyLongBreakout,
		InstrumentType:     "stock",
		Status:             "COMPLETED",
		EntryDecision:      "GO",
		CreatedAt:          week4Start.Add(5 * 24 * time.Hour),
		UpdatedAt:          week4Start.Add(5 * 24 * time.Hour),
		HeatBucket:         "Industrials",
	})

	return trades
}

// CreateSampleSettings creates sample account settings
func CreateSampleSettings() map[string]string {
	return map[string]string{
		"equity":         "100000.00",
		"risk_pct":       "0.75",
		"portfolio_cap":  "4.0",
		"bucket_cap":     "1.5",
	}
}
