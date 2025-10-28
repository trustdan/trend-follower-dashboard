package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// Decision represents a saved trading decision
type Decision struct {
	ID           int       `json:"id"`
	Date         string    `json:"date"`
	Ticker       string    `json:"ticker"`
	Action       string    `json:"action"` // GO or NO-GO
	Entry        float64   `json:"entry,omitempty"`
	ATR          float64   `json:"atr,omitempty"`
	StopDistance float64   `json:"stop_distance,omitempty"`
	InitialStop  float64   `json:"initial_stop,omitempty"`
	Shares       int       `json:"shares,omitempty"`
	Contracts    int       `json:"contracts,omitempty"`
	RiskDollars  float64   `json:"risk_dollars,omitempty"`
	Banner       string    `json:"banner"`
	Method       string    `json:"method,omitempty"`
	Delta        float64   `json:"delta,omitempty"`
	MaxLoss      float64   `json:"max_loss,omitempty"`
	Bucket       string    `json:"bucket,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	CorrID       string    `json:"corr_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// SaveDecision stores a trading decision
func (db *DB) SaveDecision(d Decision) (int, error) {
	query := `
		INSERT INTO decisions (
			date, ticker, action, entry, atr, stop_distance,
			initial_stop, shares, contracts, risk_dollars, banner,
			method, delta, max_loss, bucket, reason, corr_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		d.Date,
		d.Ticker,
		d.Action,
		d.Entry,
		d.ATR,
		d.StopDistance,
		d.InitialStop,
		d.Shares,
		d.Contracts,
		d.RiskDollars,
		d.Banner,
		d.Method,
		d.Delta,
		d.MaxLoss,
		d.Bucket,
		d.Reason,
		d.CorrID,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to save decision: %w", err)
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

// GetDecisionForToday retrieves today's decision for a ticker
func (db *DB) GetDecisionForToday(ticker string) (*Decision, error) {
	today := time.Now().Format("2006-01-02")
	return db.GetDecisionForDate(ticker, today)
}

// GetDecisionForDate retrieves a decision for a specific ticker and date
func (db *DB) GetDecisionForDate(ticker, date string) (*Decision, error) {
	query := `
		SELECT id, date, ticker, action, entry, atr, stop_distance,
		       initial_stop, shares, contracts, risk_dollars, banner,
		       method, delta, max_loss, bucket, reason, corr_id, created_at
		FROM decisions
		WHERE ticker = ? AND date = ?
		LIMIT 1
	`

	var d Decision

	err := db.conn.QueryRow(query, ticker, date).Scan(
		&d.ID,
		&d.Date,
		&d.Ticker,
		&d.Action,
		&d.Entry,
		&d.ATR,
		&d.StopDistance,
		&d.InitialStop,
		&d.Shares,
		&d.Contracts,
		&d.RiskDollars,
		&d.Banner,
		&d.Method,
		&d.Delta,
		&d.MaxLoss,
		&d.Bucket,
		&d.Reason,
		&d.CorrID,
		&d.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No decision found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get decision: %w", err)
	}

	return &d, nil
}

// CheckForDuplicateDecision checks if a decision already exists for ticker and date
func (db *DB) CheckForDuplicateDecision(ticker, date string) (bool, error) {
	query := `SELECT COUNT(*) FROM decisions WHERE ticker = ? AND date = ?`
	var count int
	err := db.conn.QueryRow(query, ticker, date).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for duplicate: %w", err)
	}
	return count > 0, nil
}
