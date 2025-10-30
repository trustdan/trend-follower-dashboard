package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// TradeHistoryEntry represents a trade in the calendar view
type TradeHistoryEntry struct {
	ID                int       `json:"id"`
	SessionID         *int      `json:"session_id,omitempty"`
	Ticker            string    `json:"ticker"`
	Strategy          string    `json:"strategy"`           // LONG_BREAKOUT, SHORT_BREAKOUT, CUSTOM
	BreakoutSystem    string    `json:"breakout_system"`    // SYSTEM_1, SYSTEM_2, CUSTOM
	OptionsStrategy   string    `json:"options_strategy"`   // LONG_CALL, IRON_CONDOR, etc.
	InstrumentType    string    `json:"instrument_type"`    // STOCK, OPTION
	Sector            string    `json:"sector"`             // Tech/Comm, Finance, etc.
	Bucket            string    `json:"bucket"`             // For heat tracking
	EntryDate         string    `json:"entry_date"`         // YYYY-MM-DD
	ExpirationDate    string    `json:"expiration_date"`    // YYYY-MM-DD (NULL for stocks)
	ExitDate          string    `json:"exit_date"`          // YYYY-MM-DD (NULL if still open)
	Status            string    `json:"status"`             // OPEN, CLOSED, ROLLED
	DTE               int       `json:"dte"`                // Days to expiration at entry
	Contracts         int       `json:"contracts"`          // Number of option contracts
	Shares            int       `json:"shares"`             // Number of shares
	RiskDollars       float64   `json:"risk_dollars"`
	EntryPrice        float64   `json:"entry_price"`
	ExitPrice         float64   `json:"exit_price"`
	PnL               float64   `json:"pnl"`
	Outcome           string    `json:"outcome"`            // WIN, LOSS, SCRATCH
	Notes             string    `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AddTradeToHistory creates a new entry in the trade_history table
func (db *DB) AddTradeToHistory(entry *TradeHistoryEntry) error {
	query := `
		INSERT INTO trade_history (
			session_id, ticker, strategy, breakout_system, options_strategy,
			instrument_type, sector, bucket, entry_date, expiration_date,
			exit_date, status, dte, contracts, shares, risk_dollars,
			entry_price, exit_price, pnl, outcome, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var sessionID interface{}
	if entry.SessionID != nil {
		sessionID = *entry.SessionID
	}

	_, err := db.conn.Exec(query,
		sessionID,
		entry.Ticker,
		entry.Strategy,
		entry.BreakoutSystem,
		entry.OptionsStrategy,
		entry.InstrumentType,
		entry.Sector,
		entry.Bucket,
		entry.EntryDate,
		nullString(entry.ExpirationDate),
		nullString(entry.ExitDate),
		entry.Status,
		nullInt(entry.DTE),
		nullInt(entry.Contracts),
		nullInt(entry.Shares),
		nullFloat(entry.RiskDollars),
		nullFloat(entry.EntryPrice),
		nullFloat(entry.ExitPrice),
		nullFloat(entry.PnL),
		nullString(entry.Outcome),
		nullString(entry.Notes),
	)

	if err != nil {
		return fmt.Errorf("failed to add trade to history: %w", err)
	}

	return nil
}

// GetCalendarView retrieves trades for a date range grouped by sector and week
// startDate and endDate should be in YYYY-MM-DD format
func (db *DB) GetCalendarView(startDate, endDate string, statusFilter string) ([]TradeHistoryEntry, error) {
	query := `
		SELECT
			id, session_id, ticker, strategy, breakout_system, options_strategy,
			instrument_type, sector, bucket, entry_date, expiration_date,
			exit_date, status, dte, contracts, shares, risk_dollars,
			entry_price, exit_price, pnl, outcome, notes,
			created_at, updated_at
		FROM trade_history
		WHERE entry_date >= ? AND entry_date <= ?
	`

	args := []interface{}{startDate, endDate}

	if statusFilter != "" {
		query += " AND status = ?"
		args = append(args, statusFilter)
	}

	query += " ORDER BY entry_date ASC, sector ASC, ticker ASC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query calendar view: %w", err)
	}
	defer rows.Close()

	var entries []TradeHistoryEntry

	for rows.Next() {
		var entry TradeHistoryEntry
		var sessionID sql.NullInt64
		var breakoutSystem, optionsStrategy, instrumentType sql.NullString
		var sector, bucket, expirationDate, exitDate sql.NullString
		var dte, contracts, shares sql.NullInt64
		var riskDollars, entryPrice, exitPrice, pnl sql.NullFloat64
		var outcome, notes sql.NullString

		err := rows.Scan(
			&entry.ID, &sessionID, &entry.Ticker, &entry.Strategy,
			&breakoutSystem, &optionsStrategy, &instrumentType,
			&sector, &bucket, &entry.EntryDate, &expirationDate,
			&exitDate, &entry.Status, &dte, &contracts, &shares,
			&riskDollars, &entryPrice, &exitPrice, &pnl,
			&outcome, &notes, &entry.CreatedAt, &entry.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade history entry: %w", err)
		}

		// Convert nullable fields
		if sessionID.Valid {
			id := int(sessionID.Int64)
			entry.SessionID = &id
		}
		if breakoutSystem.Valid {
			entry.BreakoutSystem = breakoutSystem.String
		}
		if optionsStrategy.Valid {
			entry.OptionsStrategy = optionsStrategy.String
		}
		if instrumentType.Valid {
			entry.InstrumentType = instrumentType.String
		}
		if sector.Valid {
			entry.Sector = sector.String
		}
		if bucket.Valid {
			entry.Bucket = bucket.String
		}
		if expirationDate.Valid {
			entry.ExpirationDate = expirationDate.String
		}
		if exitDate.Valid {
			entry.ExitDate = exitDate.String
		}
		if dte.Valid {
			entry.DTE = int(dte.Int64)
		}
		if contracts.Valid {
			entry.Contracts = int(contracts.Int64)
		}
		if shares.Valid {
			entry.Shares = int(shares.Int64)
		}
		if riskDollars.Valid {
			entry.RiskDollars = riskDollars.Float64
		}
		if entryPrice.Valid {
			entry.EntryPrice = entryPrice.Float64
		}
		if exitPrice.Valid {
			entry.ExitPrice = exitPrice.Float64
		}
		if pnl.Valid {
			entry.PnL = pnl.Float64
		}
		if outcome.Valid {
			entry.Outcome = outcome.String
		}
		if notes.Valid {
			entry.Notes = notes.String
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trade history: %w", err)
	}

	return entries, nil
}

// GetTradeHistoryBySector retrieves trades for a specific sector and date range
func (db *DB) GetTradeHistoryBySector(sector, startDate, endDate string) ([]TradeHistoryEntry, error) {
	query := `
		SELECT
			id, session_id, ticker, strategy, breakout_system, options_strategy,
			instrument_type, sector, bucket, entry_date, expiration_date,
			exit_date, status, dte, contracts, shares, risk_dollars,
			entry_price, exit_price, pnl, outcome, notes,
			created_at, updated_at
		FROM trade_history
		WHERE sector = ? AND entry_date >= ? AND entry_date <= ?
		ORDER BY entry_date ASC, ticker ASC
	`

	rows, err := db.conn.Query(query, sector, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query sector trades: %w", err)
	}
	defer rows.Close()

	var entries []TradeHistoryEntry

	for rows.Next() {
		var entry TradeHistoryEntry
		var sessionID sql.NullInt64
		var breakoutSystem, optionsStrategy, instrumentType sql.NullString
		var sectorVal, bucket, expirationDate, exitDate sql.NullString
		var dte, contracts, shares sql.NullInt64
		var riskDollars, entryPrice, exitPrice, pnl sql.NullFloat64
		var outcome, notes sql.NullString

		err := rows.Scan(
			&entry.ID, &sessionID, &entry.Ticker, &entry.Strategy,
			&breakoutSystem, &optionsStrategy, &instrumentType,
			&sectorVal, &bucket, &entry.EntryDate, &expirationDate,
			&exitDate, &entry.Status, &dte, &contracts, &shares,
			&riskDollars, &entryPrice, &exitPrice, &pnl,
			&outcome, &notes, &entry.CreatedAt, &entry.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade history entry: %w", err)
		}

		// Convert nullable fields (same as GetCalendarView)
		if sessionID.Valid {
			id := int(sessionID.Int64)
			entry.SessionID = &id
		}
		if breakoutSystem.Valid {
			entry.BreakoutSystem = breakoutSystem.String
		}
		if optionsStrategy.Valid {
			entry.OptionsStrategy = optionsStrategy.String
		}
		if instrumentType.Valid {
			entry.InstrumentType = instrumentType.String
		}
		if sectorVal.Valid {
			entry.Sector = sectorVal.String
		}
		if bucket.Valid {
			entry.Bucket = bucket.String
		}
		if expirationDate.Valid {
			entry.ExpirationDate = expirationDate.String
		}
		if exitDate.Valid {
			entry.ExitDate = exitDate.String
		}
		if dte.Valid {
			entry.DTE = int(dte.Int64)
		}
		if contracts.Valid {
			entry.Contracts = int(contracts.Int64)
		}
		if shares.Valid {
			entry.Shares = int(shares.Int64)
		}
		if riskDollars.Valid {
			entry.RiskDollars = riskDollars.Float64
		}
		if entryPrice.Valid {
			entry.EntryPrice = entryPrice.Float64
		}
		if exitPrice.Valid {
			entry.ExitPrice = exitPrice.Float64
		}
		if pnl.Valid {
			entry.PnL = pnl.Float64
		}
		if outcome.Valid {
			entry.Outcome = outcome.String
		}
		if notes.Valid {
			entry.Notes = notes.String
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// UpdateTradeHistory updates an existing trade history entry
func (db *DB) UpdateTradeHistory(id int, exitDate, exitPrice, pnl, outcome *string) error {
	query := `
		UPDATE trade_history
		SET exit_date = ?, exit_price = ?, pnl = ?, outcome = ?, status = ?
		WHERE id = ?
	`

	status := "CLOSED"
	if exitDate == nil {
		status = "OPEN"
	}

	_, err := db.conn.Exec(query, exitDate, exitPrice, pnl, outcome, status, id)
	if err != nil {
		return fmt.Errorf("failed to update trade history: %w", err)
	}

	return nil
}

// Helper functions for nullable values
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

func nullFloat(f float64) interface{} {
	if f == 0 {
		return nil
	}
	return f
}
