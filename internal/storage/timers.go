package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// ImpulseTimer represents an active brake timer
type ImpulseTimer struct {
	ID        int       `json:"id"`
	Ticker    string    `json:"ticker"`
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Active    bool      `json:"active"`
}

// ImpulseBrakeDuration is the mandatory 2-minute wait period
const ImpulseBrakeDuration = 2 * time.Minute

// StartImpulseTimer creates a new timer for a ticker
// This is called automatically when checklist evaluation returns GREEN
func (db *DB) StartImpulseTimer(ticker string) error {
	now := time.Now()
	expiresAt := now.Add(ImpulseBrakeDuration)

	// Deactivate any existing timers for this ticker
	_, err := db.conn.Exec(`UPDATE impulse_timers SET active = 0 WHERE ticker = ? AND active = 1`, ticker)
	if err != nil {
		return fmt.Errorf("failed to deactivate old timers: %w", err)
	}

	// Insert new timer
	query := `
		INSERT INTO impulse_timers (ticker, started_at, expires_at, active)
		VALUES (?, ?, ?, 1)
	`

	_, err = db.conn.Exec(query, ticker, now.Unix(), expiresAt.Unix())
	if err != nil {
		return fmt.Errorf("failed to start impulse timer: %w", err)
	}

	return nil
}

// GetActiveTimer retrieves the active timer for a ticker
// Returns nil if no active timer exists
func (db *DB) GetActiveTimer(ticker string) (*ImpulseTimer, error) {
	query := `
		SELECT id, ticker, started_at, expires_at, active
		FROM impulse_timers
		WHERE ticker = ? AND active = 1
		ORDER BY started_at DESC
		LIMIT 1
	`

	var timer ImpulseTimer
	var startedUnix, expiresUnix int64

	err := db.conn.QueryRow(query, ticker).Scan(
		&timer.ID,
		&timer.Ticker,
		&startedUnix,
		&expiresUnix,
		&timer.Active,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No active timer
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get timer: %w", err)
	}

	timer.StartedAt = time.Unix(startedUnix, 0)
	timer.ExpiresAt = time.Unix(expiresUnix, 0)

	return &timer, nil
}

// CheckImpulseBrake validates timer status before allowing save
// Returns error if:
//   - No active timer exists (must evaluate checklist first)
//   - Timer has not expired yet (must wait full 2 minutes)
func (db *DB) CheckImpulseBrake(ticker string) error {
	timer, err := db.GetActiveTimer(ticker)
	if err != nil {
		return fmt.Errorf("failed to check timer: %w", err)
	}

	if timer == nil {
		return fmt.Errorf("no impulse timer active for %s (evaluate checklist first)", ticker)
	}

	now := time.Now()
	if now.Before(timer.ExpiresAt) {
		remaining := timer.ExpiresAt.Sub(now)
		return fmt.Errorf("impulse brake active, wait %.0f more seconds", remaining.Seconds())
	}

	return nil
}
