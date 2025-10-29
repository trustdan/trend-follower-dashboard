package storage

import (
	"fmt"
)

// Settings represents the account settings for API responses
type Settings struct {
	Equity       float64 `json:"equity"`
	RiskPct      float64 `json:"riskPct"`
	PortfolioCap float64 `json:"portfolioCap"`
	BucketCap    float64 `json:"bucketCap"`
	MaxUnits     int     `json:"maxUnits"`
}

// Candidate represents a trade candidate for API responses
type Candidate struct {
	ID     int    `json:"id,omitempty"`
	Ticker string `json:"ticker"`
	Date   string `json:"date"`
	Sector string `json:"sector,omitempty"`
	Bucket string `json:"bucket,omitempty"`
}

// GetSettings retrieves settings as a struct for API responses
func (db *DB) GetSettings() (*Settings, error) {
	all, err := db.GetAllSettings()
	if err != nil {
		return nil, err
	}

	settings := &Settings{}

	// Parse each setting with defaults
	if val, ok := all["Equity_E"]; ok {
		fmt.Sscanf(val, "%f", &settings.Equity)
	} else {
		settings.Equity = 100000.0 // default
	}

	if val, ok := all["RiskPct_r"]; ok {
		fmt.Sscanf(val, "%f", &settings.RiskPct)
	} else {
		settings.RiskPct = 0.75 // default
	}

	// HeatCap stored as decimal (0.04) in DB, return as percentage (4.0) for API
	if val, ok := all["HeatCap_H_pct"]; ok {
		var capDecimal float64
		fmt.Sscanf(val, "%f", &capDecimal)
		settings.PortfolioCap = capDecimal * 100 // Convert 0.04 -> 4.0
	} else {
		settings.PortfolioCap = 4.0 // default
	}

	// BucketHeatCap stored as decimal (0.015) in DB, return as percentage (1.5) for API
	if val, ok := all["BucketHeatCap_pct"]; ok {
		var capDecimal float64
		fmt.Sscanf(val, "%f", &capDecimal)
		settings.BucketCap = capDecimal * 100 // Convert 0.015 -> 1.5
	} else {
		settings.BucketCap = 1.5 // default
	}

	if val, ok := all["MaxUnits"]; ok {
		var maxUnits int
		fmt.Sscanf(val, "%d", &maxUnits)
		settings.MaxUnits = maxUnits
	} else {
		settings.MaxUnits = 4 // default
	}

	return settings, nil
}

// GetPositions retrieves all open positions for API responses
func (db *DB) GetPositions() ([]Position, error) {
	return db.GetOpenPositions()
}

// GetCandidates retrieves candidates for a given date
func (db *DB) GetCandidates(date string) ([]Candidate, error) {
	candidates, err := db.GetCandidatesForDate(date)
	if err != nil {
		return nil, err
	}

	result := make([]Candidate, 0, len(candidates))
	for _, c := range candidates {
		candidate := Candidate{}

		if id, ok := c["id"].(int); ok {
			candidate.ID = id
		}
		if ticker, ok := c["ticker"].(string); ok {
			candidate.Ticker = ticker
		}
		if date, ok := c["date"].(string); ok {
			candidate.Date = date
		}
		if sector, ok := c["sector"].(string); ok {
			candidate.Sector = sector
		}
		if bucket, ok := c["bucket"].(string); ok {
			candidate.Bucket = bucket
		}

		result = append(result, candidate)
	}

	return result, nil
}

// AddCandidates adds candidates for a given date (simplified API)
func (db *DB) AddCandidates(tickers []string, date string) error {
	// Use ImportCandidates with nil preset ID and empty sector/bucket
	// Individual tickers can be updated later with sector/bucket info
	return db.ImportCandidates(date, tickers, nil, "", "")
}

// NewDB is an alias for New to match the expected API
func NewDB(dbPath string) (*DB, error) {
	return New(dbPath)
}
