package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// Position represents an open or closed trade
type Position struct {
	ID          int       `json:"id"`
	Ticker      string    `json:"ticker"`
	EntryPrice  float64   `json:"entry_price"`
	CurrentStop float64   `json:"current_stop"`
	InitialStop float64   `json:"initial_stop"`
	Shares      int       `json:"shares"`
	RiskDollars float64   `json:"risk_dollars"`
	Bucket      string    `json:"bucket,omitempty"`
	Status      string    `json:"status"` // OPEN or CLOSED
	ExitPrice   float64   `json:"exit_price,omitempty"`
	ExitDate    string    `json:"exit_date,omitempty"`
	Outcome     string    `json:"outcome,omitempty"` // WIN, LOSS, SCRATCH
	PnL         float64   `json:"pnl,omitempty"`
	DecisionID  int       `json:"decision_id"`
	OpenedAt    time.Time `json:"opened_at"`
	ClosedAt    time.Time `json:"closed_at,omitempty"`
}

// OpenPosition creates a new position from a GO decision
func (db *DB) OpenPosition(ticker string) (*Position, error) {
	// Get the GO decision
	decision, err := db.GetDecisionForToday(ticker)
	if err != nil {
		return nil, fmt.Errorf("no decision found for %s: %w", ticker, err)
	}
	if decision == nil {
		return nil, fmt.Errorf("no decision found for %s", ticker)
	}

	if decision.Action != "GO" {
		return nil, fmt.Errorf("cannot open position for NO-GO decision")
	}

	// Get candidate to find bucket
	today := time.Now().Format("2006-01-02")
	candidates, err := db.GetCandidatesForDate(today)
	if err != nil {
		return nil, err
	}

	var bucket string
	for _, c := range candidates {
		if tickerVal, ok := c["ticker"].(string); ok && tickerVal == ticker {
			if bucketVal, ok := c["bucket"].(string); ok {
				bucket = bucketVal
			}
			break
		}
	}

	// Create position
	query := `
		INSERT INTO positions (
			ticker, entry_price, current_stop, initial_stop,
			shares, risk_dollars, bucket, status, decision_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, 'OPEN', ?)
	`

	result, err := db.conn.Exec(query,
		ticker,
		decision.Entry,
		decision.InitialStop,
		decision.InitialStop,
		decision.Shares,
		decision.RiskDollars,
		bucket,
		decision.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to open position: %w", err)
	}

	id, _ := result.LastInsertId()

	position := &Position{
		ID:          int(id),
		Ticker:      ticker,
		EntryPrice:  decision.Entry,
		CurrentStop: decision.InitialStop,
		InitialStop: decision.InitialStop,
		Shares:      decision.Shares,
		RiskDollars: decision.RiskDollars,
		Bucket:      bucket,
		Status:      "OPEN",
		DecisionID:  decision.ID,
		OpenedAt:    time.Now(),
	}

	return position, nil
}

// GetPosition retrieves a position by ID
func (db *DB) GetPosition(id int) (*Position, error) {
	query := `
		SELECT id, ticker, entry_price, current_stop, initial_stop,
		       shares, risk_dollars, bucket, status, exit_price, exit_date,
		       outcome, pnl, decision_id, opened_at, closed_at
		FROM positions
		WHERE id = ?
	`

	var p Position
	var bucket, exitDate, outcome sql.NullString
	var exitPrice, pnl sql.NullFloat64
	var closedAt sql.NullTime

	err := db.conn.QueryRow(query, id).Scan(
		&p.ID, &p.Ticker, &p.EntryPrice, &p.CurrentStop, &p.InitialStop,
		&p.Shares, &p.RiskDollars, &bucket, &p.Status, &exitPrice, &exitDate,
		&outcome, &pnl, &p.DecisionID, &p.OpenedAt, &closedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	if bucket.Valid {
		p.Bucket = bucket.String
	}
	if exitPrice.Valid {
		p.ExitPrice = exitPrice.Float64
	}
	if exitDate.Valid {
		p.ExitDate = exitDate.String
	}
	if outcome.Valid {
		p.Outcome = outcome.String
	}
	if pnl.Valid {
		p.PnL = pnl.Float64
	}
	if closedAt.Valid {
		p.ClosedAt = closedAt.Time
	}

	return &p, nil
}

// GetPositionByTicker retrieves an open position for a ticker
func (db *DB) GetPositionByTicker(ticker string) (*Position, error) {
	query := `
		SELECT id, ticker, entry_price, current_stop, initial_stop,
		       shares, risk_dollars, bucket, status, exit_price, exit_date,
		       outcome, pnl, decision_id, opened_at, closed_at
		FROM positions
		WHERE ticker = ? AND status = 'OPEN'
		ORDER BY opened_at DESC
		LIMIT 1
	`

	var p Position
	var bucket, exitDate, outcome sql.NullString
	var exitPrice, pnl sql.NullFloat64
	var closedAt sql.NullTime

	err := db.conn.QueryRow(query, ticker).Scan(
		&p.ID, &p.Ticker, &p.EntryPrice, &p.CurrentStop, &p.InitialStop,
		&p.Shares, &p.RiskDollars, &bucket, &p.Status, &exitPrice, &exitDate,
		&outcome, &pnl, &p.DecisionID, &p.OpenedAt, &closedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no open position found for %s", ticker)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	if bucket.Valid {
		p.Bucket = bucket.String
	}

	return &p, nil
}

// UpdateStop updates the stop price for a position
func (db *DB) UpdateStop(ticker string, newStop float64) error {
	position, err := db.GetPositionByTicker(ticker)
	if err != nil {
		return err
	}

	// Validate stop movement (cannot move against position)
	// Assume LONG positions for now (can extend later for SHORT)
	if newStop < position.CurrentStop {
		return fmt.Errorf("cannot move stop down for long position (current: %.2f, new: %.2f)",
			position.CurrentStop, newStop)
	}

	// Recalculate risk based on new stop
	newRisk := float64(position.Shares) * (position.EntryPrice - newStop)

	query := `
		UPDATE positions
		SET current_stop = ?, risk_dollars = ?
		WHERE id = ?
	`

	_, err = db.conn.Exec(query, newStop, newRisk, position.ID)
	if err != nil {
		return fmt.Errorf("failed to update stop: %w", err)
	}

	return nil
}

// ClosePosition closes a position and calculates P&L
func (db *DB) ClosePosition(ticker string, exitPrice float64, outcome string) error {
	position, err := db.GetPositionByTicker(ticker)
	if err != nil {
		return err
	}

	// Calculate P&L (assume LONG)
	pnl := float64(position.Shares) * (exitPrice - position.EntryPrice)

	// Update position
	query := `
		UPDATE positions
		SET status = 'CLOSED',
		    exit_price = ?,
		    exit_date = ?,
		    outcome = ?,
		    pnl = ?,
		    closed_at = ?
		WHERE id = ?
	`

	now := time.Now()
	exitDate := now.Format("2006-01-02")

	_, err = db.conn.Exec(query, exitPrice, exitDate, outcome, pnl, now, position.ID)
	if err != nil {
		return fmt.Errorf("failed to close position: %w", err)
	}

	// Trigger cooldown if loss
	if outcome == "LOSS" && position.Bucket != "" {
		reason := fmt.Sprintf("Loss on %s", ticker)
		err = db.TriggerBucketCooldown(position.Bucket, reason)
		if err != nil {
			return fmt.Errorf("failed to trigger cooldown: %w", err)
		}
	}

	return nil
}

// GetAllPositions retrieves all positions, optionally filtered by status
func (db *DB) GetAllPositions(status string) ([]Position, error) {
	query := `
		SELECT id, ticker, entry_price, current_stop, initial_stop,
		       shares, risk_dollars, bucket, status, exit_price, exit_date,
		       outcome, pnl, decision_id, opened_at, closed_at
		FROM positions
	`

	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	query += " ORDER BY opened_at DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query positions: %w", err)
	}
	defer rows.Close()

	positions := []Position{}

	for rows.Next() {
		var p Position
		var bucket, exitDate, outcome sql.NullString
		var exitPrice, pnl sql.NullFloat64
		var closedAt sql.NullTime

		err := rows.Scan(
			&p.ID, &p.Ticker, &p.EntryPrice, &p.CurrentStop, &p.InitialStop,
			&p.Shares, &p.RiskDollars, &bucket, &p.Status, &exitPrice, &exitDate,
			&outcome, &pnl, &p.DecisionID, &p.OpenedAt, &closedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan position: %w", err)
		}

		if bucket.Valid {
			p.Bucket = bucket.String
		}
		if exitPrice.Valid {
			p.ExitPrice = exitPrice.Float64
		}
		if exitDate.Valid {
			p.ExitDate = exitDate.String
		}
		if outcome.Valid {
			p.Outcome = outcome.String
		}
		if pnl.Valid {
			p.PnL = pnl.Float64
		}
		if closedAt.Valid {
			p.ClosedAt = closedAt.Time
		}

		positions = append(positions, p)
	}

	return positions, nil
}

// GetOpenPositions retrieves all open positions
func (db *DB) GetOpenPositions() ([]Position, error) {
	return db.GetAllPositions("OPEN")
}

// CalculatePortfolioHeat calculates total risk from all open positions
func (db *DB) CalculatePortfolioHeat() (float64, error) {
	positions, err := db.GetOpenPositions()
	if err != nil {
		return 0, err
	}

	var totalHeat float64
	for _, p := range positions {
		totalHeat += p.RiskDollars
	}

	return totalHeat, nil
}

// CalculateBucketHeat calculates total risk for a specific bucket from open positions
func (db *DB) CalculateBucketHeat(bucket string) (float64, error) {
	query := `
		SELECT SUM(risk_dollars)
		FROM positions
		WHERE bucket = ? AND status = 'OPEN'
	`

	var heat sql.NullFloat64
	err := db.conn.QueryRow(query, bucket).Scan(&heat)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate bucket heat: %w", err)
	}

	if !heat.Valid {
		return 0, nil
	}

	return heat.Float64, nil
}
