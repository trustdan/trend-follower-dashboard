package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// TradeSession represents a complete trade evaluation workflow
type TradeSession struct {
	// Identity
	ID         int    `json:"id"`
	SessionNum int    `json:"session_num"`
	Ticker     string `json:"ticker"`
	Strategy   string `json:"strategy"` // LONG_BREAKOUT, SHORT_BREAKOUT, CUSTOM

	// Provenance from FINVIZ/candidates
	Source       string `json:"source"`        // MANUAL, PRESET, CUSTOM
	CandidateID  *int   `json:"candidate_id"`  // FK to candidates table
	PresetID     *int   `json:"preset_id"`     // FK to presets table
	PresetName   string `json:"preset_name"`   // e.g., "TF_BREAKOUT_LONG"
	ScanDate     string `json:"scan_date"`     // Date of FINVIZ scan

	// Workflow state
	Status      string `json:"status"`       // DRAFT, EVALUATING, COMPLETED, ABANDONED
	CurrentStep string `json:"current_step"` // CHECKLIST, SIZING, HEAT, ENTRY

	// Options Trading Metadata (Phase 1)
	InstrumentType        string  `json:"instrument_type,omitempty"`          // STOCK, OPTION
	OptionsStrategy       string  `json:"options_strategy,omitempty"`         // See constants in options.go
	EntryDate             string  `json:"entry_date,omitempty"`               // YYYY-MM-DD
	PrimaryExpirationDate string  `json:"primary_expiration_date,omitempty"`  // Primary leg expiration
	DTE                   int     `json:"dte,omitempty"`                      // Days to expiration at entry
	RollThresholdDTE      int     `json:"roll_threshold_dte,omitempty"`       // DTE to roll/close (default 21)
	TimeExitMode          string  `json:"time_exit_mode,omitempty"`           // None, Close, Roll
	LegsJSON              string  `json:"legs_json,omitempty"`                // JSON array of OptionLeg
	NetDebit              float64 `json:"net_debit,omitempty"`                // Total debit (negative = credit)
	MaxProfit             float64 `json:"max_profit,omitempty"`               // Maximum theoretical profit
	MaxLoss               float64 `json:"max_loss,omitempty"`                 // Maximum theoretical loss
	BreakevenLower        float64 `json:"breakeven_lower,omitempty"`          // Lower breakeven price
	BreakevenUpper        float64 `json:"breakeven_upper,omitempty"`          // Upper breakeven price
	UnderlyingAtEntry     float64 `json:"underlying_at_entry,omitempty"`      // Stock price at entry

	// Pyramiding Metadata (Van Tharp method)
	MaxUnits     int     `json:"max_units,omitempty"`      // Maximum pyramid units (default 4)
	AddStepN     float64 `json:"add_step_n,omitempty"`     // Add every X * N (default 0.5)
	CurrentUnits int     `json:"current_units,omitempty"`  // Current units (0-4)
	AddPrice1    float64 `json:"add_price_1,omitempty"`    // Entry + 0.5N
	AddPrice2    float64 `json:"add_price_2,omitempty"`    // Entry + 1.0N
	AddPrice3    float64 `json:"add_price_3,omitempty"`    // Entry + 1.5N

	// Breakout System Parameters (for documentation)
	EntryLookback int `json:"entry_lookback,omitempty"` // 20 or 55 for System-1/System-2
	ExitLookback  int `json:"exit_lookback,omitempty"`  // 10-bar exit

	// Gate 1: Checklist
	ChecklistCompleted   bool      `json:"checklist_completed"`
	ChecklistBanner      string    `json:"checklist_banner,omitempty"`      // GREEN, YELLOW, RED
	ChecklistMissingCount int      `json:"checklist_missing_count"`
	ChecklistQualityScore int      `json:"checklist_quality_score"`
	ChecklistCompletedAt *time.Time `json:"checklist_completed_at,omitempty"`

	// Gate 2: Position Sizing
	SizingCompleted    bool       `json:"sizing_completed"`
	SizingMethod       string     `json:"sizing_method,omitempty"`       // stock, opt-delta-atr, opt-contracts
	SizingEntryPrice   float64    `json:"sizing_entry_price,omitempty"`
	SizingATR          float64    `json:"sizing_atr,omitempty"`
	SizingKMultiple    float64    `json:"sizing_k_multiple,omitempty"`
	SizingStopDistance float64    `json:"sizing_stop_distance,omitempty"`
	SizingInitialStop  float64    `json:"sizing_initial_stop,omitempty"`
	SizingShares       int        `json:"sizing_shares,omitempty"`
	SizingContracts    int        `json:"sizing_contracts,omitempty"`
	SizingRiskDollars  float64    `json:"sizing_risk_dollars,omitempty"`
	SizingDelta        float64    `json:"sizing_delta,omitempty"`
	SizingCompletedAt  *time.Time `json:"sizing_completed_at,omitempty"`

	// Gate 3: Heat Check
	HeatCompleted        bool       `json:"heat_completed"`
	HeatStatus           string     `json:"heat_status,omitempty"`           // OK, WARN, REJECT
	HeatPortfolioCurrent float64    `json:"heat_portfolio_current,omitempty"`
	HeatPortfolioNew     float64    `json:"heat_portfolio_new,omitempty"`
	HeatPortfolioCap     float64    `json:"heat_portfolio_cap,omitempty"`
	HeatBucket           string     `json:"heat_bucket,omitempty"`
	HeatBucketCurrent    float64    `json:"heat_bucket_current,omitempty"`
	HeatBucketNew        float64    `json:"heat_bucket_new,omitempty"`
	HeatBucketCap        float64    `json:"heat_bucket_cap,omitempty"`
	HeatCompletedAt      *time.Time `json:"heat_completed_at,omitempty"`

	// Gate 4: Trade Entry
	EntryCompleted   bool       `json:"entry_completed"`
	EntryDecision    string     `json:"entry_decision,omitempty"`    // GO, NO-GO
	EntryDecisionID  *int       `json:"entry_decision_id,omitempty"` // FK to decisions table
	EntryGate1Pass   *bool      `json:"entry_gate1_pass,omitempty"`  // Banner GREEN
	EntryGate2Pass   *bool      `json:"entry_gate2_pass,omitempty"`  // 2-min cooloff
	EntryGate3Pass   *bool      `json:"entry_gate3_pass,omitempty"`  // Ticker not on cooldown
	EntryGate4Pass   *bool      `json:"entry_gate4_pass,omitempty"`  // Heat caps OK
	EntryGate5Pass   *bool      `json:"entry_gate5_pass,omitempty"`  // Sizing complete
	EntryCompletedAt *time.Time `json:"entry_completed_at,omitempty"`

	// Audit trail
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Strategy constants
const (
	StrategyLongBreakout  = "LONG_BREAKOUT"
	StrategyShortBreakout = "SHORT_BREAKOUT"
	StrategyCustom        = "CUSTOM"
)

// Session status constants
const (
	StatusDraft      = "DRAFT"
	StatusEvaluating = "EVALUATING"
	StatusCompleted  = "COMPLETED"
	StatusAbandoned  = "ABANDONED"
)

// Session step constants
const (
	StepChecklist = "CHECKLIST"
	StepSizing    = "SIZING"
	StepHeat      = "HEAT"
	StepEntry     = "ENTRY"
)

// CreateSession creates a new trade session
// Session number is auto-generated from the primary key (id)
func (db *DB) CreateSession(ticker, strategy string) (*TradeSession, error) {
	now := time.Now()

	query := `
		INSERT INTO trade_sessions (
			ticker, strategy, source, status, current_step,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		ticker,
		strategy,
		"MANUAL", // default source
		StatusDraft,
		StepChecklist,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get session id: %w", err)
	}

	// Retrieve the created session (session_num is auto-generated from id)
	return db.GetSession(int(id))
}

// CreateSessionFromPreset creates a new trade session with FINVIZ provenance
func (db *DB) CreateSessionFromPreset(ticker, strategy string, candidateID, presetID int, presetName, scanDate string) (*TradeSession, error) {
	now := time.Now()

	query := `
		INSERT INTO trade_sessions (
			ticker, strategy, source, candidate_id, preset_id, preset_name, scan_date,
			status, current_step, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		ticker,
		strategy,
		"PRESET",
		candidateID,
		presetID,
		presetName,
		scanDate,
		StatusDraft,
		StepChecklist,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session from preset: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get session id: %w", err)
	}

	return db.GetSession(int(id))
}

// CreateSessionWithOptions creates a new trade session with full options metadata
func (db *DB) CreateSessionWithOptions(ticker, strategy, instrumentType, optionsStrategy string,
	entryDate, primaryExpirationDate string, dte, rollThresholdDTE int, timeExitMode, legsJSON string,
	netDebit, maxProfit, maxLoss, breakevenLower, breakevenUpper, underlyingAtEntry float64,
	maxUnits int, addStepN float64, entryLookback, exitLookback int) (*TradeSession, error) {

	now := time.Now()

	query := `
		INSERT INTO trade_sessions (
			ticker, strategy, source, status, current_step,
			instrument_type, options_strategy, entry_date, primary_expiration_date,
			dte, roll_threshold_dte, time_exit_mode, legs_json,
			net_debit, max_profit, max_loss, breakeven_lower, breakeven_upper, underlying_at_entry,
			max_units, add_step_n, current_units,
			entry_lookback, exit_lookback,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		ticker,
		strategy,
		"MANUAL",
		StatusDraft,
		StepChecklist,
		instrumentType,
		optionsStrategy,
		entryDate,
		primaryExpirationDate,
		dte,
		rollThresholdDTE,
		timeExitMode,
		legsJSON,
		netDebit,
		maxProfit,
		maxLoss,
		breakevenLower,
		breakevenUpper,
		underlyingAtEntry,
		maxUnits,
		addStepN,
		0, // current_units starts at 0
		entryLookback,
		exitLookback,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session with options: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get session id: %w", err)
	}

	return db.GetSession(int(id))
}

// GetSession retrieves a session by its ID
func (db *DB) GetSession(id int) (*TradeSession, error) {
	query := `
		SELECT
			id, session_num, ticker, strategy, source,
			candidate_id, preset_id, preset_name, scan_date,
			status, current_step,
			instrument_type, options_strategy, entry_date, primary_expiration_date,
			dte, roll_threshold_dte, time_exit_mode, legs_json,
			net_debit, max_profit, max_loss, breakeven_lower, breakeven_upper, underlying_at_entry,
			max_units, add_step_n, current_units, add_price_1, add_price_2, add_price_3,
			entry_lookback, exit_lookback,
			checklist_completed, checklist_banner, checklist_missing_count,
			checklist_quality_score, checklist_completed_at,
			sizing_completed, sizing_method, sizing_entry_price, sizing_atr,
			sizing_k_multiple, sizing_stop_distance, sizing_initial_stop,
			sizing_shares, sizing_contracts, sizing_risk_dollars, sizing_delta,
			sizing_completed_at,
			heat_completed, heat_status,
			heat_portfolio_current, heat_portfolio_new, heat_portfolio_cap,
			heat_bucket, heat_bucket_current, heat_bucket_new, heat_bucket_cap,
			heat_completed_at,
			entry_completed, entry_decision, entry_decision_id,
			entry_gate1_pass, entry_gate2_pass, entry_gate3_pass,
			entry_gate4_pass, entry_gate5_pass, entry_completed_at,
			created_at, updated_at, completed_at
		FROM trade_sessions
		WHERE id = ?
	`

	session := &TradeSession{}
	var candidateID, presetID, entryDecisionID sql.NullInt64
	var presetName, scanDate, checklistBanner, sizingMethod sql.NullString
	var heatStatus, heatBucket, entryDecision sql.NullString
	var checklistCompletedAt, sizingCompletedAt, heatCompletedAt, entryCompletedAt, completedAt sql.NullTime
	var sizingEntryPrice, sizingATR, sizingKMultiple, sizingStopDistance, sizingInitialStop sql.NullFloat64
	var sizingRiskDollars, sizingDelta sql.NullFloat64
	var sizingShares, sizingContracts sql.NullInt64
	var heatPortfolioCurrent, heatPortfolioNew, heatPortfolioCap sql.NullFloat64
	var heatBucketCurrent, heatBucketNew, heatBucketCap sql.NullFloat64
	var entryGate1, entryGate2, entryGate3, entryGate4, entryGate5 sql.NullInt64
	var checklistCompleted, sizingCompleted, heatCompleted, entryCompleted int

	// Options metadata nullable fields
	var instrumentType, optionsStrategy, entryDate, primaryExpirationDate sql.NullString
	var dte, rollThresholdDTE sql.NullInt64
	var timeExitMode, legsJSON sql.NullString
	var netDebit, maxProfit, maxLoss, breakevenLower, breakevenUpper, underlyingAtEntry sql.NullFloat64
	var maxUnits, currentUnits, entryLookback, exitLookback sql.NullInt64
	var addStepN, addPrice1, addPrice2, addPrice3 sql.NullFloat64

	err := db.conn.QueryRow(query, id).Scan(
		&session.ID, &session.SessionNum, &session.Ticker, &session.Strategy, &session.Source,
		&candidateID, &presetID, &presetName, &scanDate,
		&session.Status, &session.CurrentStep,
		&instrumentType, &optionsStrategy, &entryDate, &primaryExpirationDate,
		&dte, &rollThresholdDTE, &timeExitMode, &legsJSON,
		&netDebit, &maxProfit, &maxLoss, &breakevenLower, &breakevenUpper, &underlyingAtEntry,
		&maxUnits, &addStepN, &currentUnits, &addPrice1, &addPrice2, &addPrice3,
		&entryLookback, &exitLookback,
		&checklistCompleted, &checklistBanner, &session.ChecklistMissingCount,
		&session.ChecklistQualityScore, &checklistCompletedAt,
		&sizingCompleted, &sizingMethod, &sizingEntryPrice, &sizingATR,
		&sizingKMultiple, &sizingStopDistance, &sizingInitialStop,
		&sizingShares, &sizingContracts, &sizingRiskDollars, &sizingDelta,
		&sizingCompletedAt,
		&heatCompleted, &heatStatus,
		&heatPortfolioCurrent, &heatPortfolioNew, &heatPortfolioCap,
		&heatBucket, &heatBucketCurrent, &heatBucketNew, &heatBucketCap,
		&heatCompletedAt,
		&entryCompleted, &entryDecision, &entryDecisionID,
		&entryGate1, &entryGate2, &entryGate3, &entryGate4, &entryGate5,
		&entryCompletedAt,
		&session.CreatedAt, &session.UpdatedAt, &completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Convert nullable fields
	if candidateID.Valid {
		id := int(candidateID.Int64)
		session.CandidateID = &id
	}
	if presetID.Valid {
		id := int(presetID.Int64)
		session.PresetID = &id
	}
	if presetName.Valid {
		session.PresetName = presetName.String
	}
	if scanDate.Valid {
		session.ScanDate = scanDate.String
	}

	// Checklist fields
	session.ChecklistCompleted = checklistCompleted == 1
	if checklistBanner.Valid {
		session.ChecklistBanner = checklistBanner.String
	}
	if checklistCompletedAt.Valid {
		session.ChecklistCompletedAt = &checklistCompletedAt.Time
	}

	// Sizing fields
	session.SizingCompleted = sizingCompleted == 1
	if sizingMethod.Valid {
		session.SizingMethod = sizingMethod.String
	}
	if sizingEntryPrice.Valid {
		session.SizingEntryPrice = sizingEntryPrice.Float64
	}
	if sizingATR.Valid {
		session.SizingATR = sizingATR.Float64
	}
	if sizingKMultiple.Valid {
		session.SizingKMultiple = sizingKMultiple.Float64
	}
	if sizingStopDistance.Valid {
		session.SizingStopDistance = sizingStopDistance.Float64
	}
	if sizingInitialStop.Valid {
		session.SizingInitialStop = sizingInitialStop.Float64
	}
	if sizingShares.Valid {
		session.SizingShares = int(sizingShares.Int64)
	}
	if sizingContracts.Valid {
		session.SizingContracts = int(sizingContracts.Int64)
	}
	if sizingRiskDollars.Valid {
		session.SizingRiskDollars = sizingRiskDollars.Float64
	}
	if sizingDelta.Valid {
		session.SizingDelta = sizingDelta.Float64
	}
	if sizingCompletedAt.Valid {
		session.SizingCompletedAt = &sizingCompletedAt.Time
	}

	// Heat fields
	session.HeatCompleted = heatCompleted == 1
	if heatStatus.Valid {
		session.HeatStatus = heatStatus.String
	}
	if heatPortfolioCurrent.Valid {
		session.HeatPortfolioCurrent = heatPortfolioCurrent.Float64
	}
	if heatPortfolioNew.Valid {
		session.HeatPortfolioNew = heatPortfolioNew.Float64
	}
	if heatPortfolioCap.Valid {
		session.HeatPortfolioCap = heatPortfolioCap.Float64
	}
	if heatBucket.Valid {
		session.HeatBucket = heatBucket.String
	}
	if heatBucketCurrent.Valid {
		session.HeatBucketCurrent = heatBucketCurrent.Float64
	}
	if heatBucketNew.Valid {
		session.HeatBucketNew = heatBucketNew.Float64
	}
	if heatBucketCap.Valid {
		session.HeatBucketCap = heatBucketCap.Float64
	}
	if heatCompletedAt.Valid {
		session.HeatCompletedAt = &heatCompletedAt.Time
	}

	// Entry fields
	session.EntryCompleted = entryCompleted == 1
	if entryDecision.Valid {
		session.EntryDecision = entryDecision.String
	}
	if entryDecisionID.Valid {
		id := int(entryDecisionID.Int64)
		session.EntryDecisionID = &id
	}
	if entryGate1.Valid {
		pass := entryGate1.Int64 == 1
		session.EntryGate1Pass = &pass
	}
	if entryGate2.Valid {
		pass := entryGate2.Int64 == 1
		session.EntryGate2Pass = &pass
	}
	if entryGate3.Valid {
		pass := entryGate3.Int64 == 1
		session.EntryGate3Pass = &pass
	}
	if entryGate4.Valid {
		pass := entryGate4.Int64 == 1
		session.EntryGate4Pass = &pass
	}
	if entryGate5.Valid {
		pass := entryGate5.Int64 == 1
		session.EntryGate5Pass = &pass
	}
	if entryCompletedAt.Valid {
		session.EntryCompletedAt = &entryCompletedAt.Time
	}
	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}

	// Options metadata fields
	if instrumentType.Valid {
		session.InstrumentType = instrumentType.String
	}
	if optionsStrategy.Valid {
		session.OptionsStrategy = optionsStrategy.String
	}
	if entryDate.Valid {
		session.EntryDate = entryDate.String
	}
	if primaryExpirationDate.Valid {
		session.PrimaryExpirationDate = primaryExpirationDate.String
	}
	if dte.Valid {
		session.DTE = int(dte.Int64)
	}
	if rollThresholdDTE.Valid {
		session.RollThresholdDTE = int(rollThresholdDTE.Int64)
	}
	if timeExitMode.Valid {
		session.TimeExitMode = timeExitMode.String
	}
	if legsJSON.Valid {
		session.LegsJSON = legsJSON.String
	}
	if netDebit.Valid {
		session.NetDebit = netDebit.Float64
	}
	if maxProfit.Valid {
		session.MaxProfit = maxProfit.Float64
	}
	if maxLoss.Valid {
		session.MaxLoss = maxLoss.Float64
	}
	if breakevenLower.Valid {
		session.BreakevenLower = breakevenLower.Float64
	}
	if breakevenUpper.Valid {
		session.BreakevenUpper = breakevenUpper.Float64
	}
	if underlyingAtEntry.Valid {
		session.UnderlyingAtEntry = underlyingAtEntry.Float64
	}

	// Pyramid metadata fields
	if maxUnits.Valid {
		session.MaxUnits = int(maxUnits.Int64)
	}
	if addStepN.Valid {
		session.AddStepN = addStepN.Float64
	}
	if currentUnits.Valid {
		session.CurrentUnits = int(currentUnits.Int64)
	}
	if addPrice1.Valid {
		session.AddPrice1 = addPrice1.Float64
	}
	if addPrice2.Valid {
		session.AddPrice2 = addPrice2.Float64
	}
	if addPrice3.Valid {
		session.AddPrice3 = addPrice3.Float64
	}

	// Breakout system metadata fields
	if entryLookback.Valid {
		session.EntryLookback = int(entryLookback.Int64)
	}
	if exitLookback.Valid {
		session.ExitLookback = int(exitLookback.Int64)
	}

	return session, nil
}

// GetSessionByNum retrieves a session by its session_num
func (db *DB) GetSessionByNum(sessionNum int) (*TradeSession, error) {
	query := `SELECT id FROM trade_sessions WHERE session_num = ?`
	var id int
	err := db.conn.QueryRow(query, sessionNum).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session #%d not found", sessionNum)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	return db.GetSession(id)
}

// UpdateSessionChecklist updates the checklist gate completion for a session
func (db *DB) UpdateSessionChecklist(id int, banner string, missingCount, qualityScore int) error {
	now := time.Now()
	completed := 1
	if banner == "RED" || banner == "YELLOW" {
		completed = 0
	}

	query := `
		UPDATE trade_sessions
		SET checklist_completed = ?,
		    checklist_banner = ?,
		    checklist_missing_count = ?,
		    checklist_quality_score = ?,
		    checklist_completed_at = ?,
		    current_step = CASE WHEN ? = 1 THEN 'SIZING' ELSE current_step END
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, completed, banner, missingCount, qualityScore, now, completed, id)
	if err != nil {
		return fmt.Errorf("failed to update session checklist: %w", err)
	}
	return nil
}

// UpdateSessionSizing updates the position sizing gate completion for a session
func (db *DB) UpdateSessionSizing(id int, method string, entryPrice, atr, kMultiple, stopDistance, initialStop float64,
	shares, contracts int, riskDollars, delta float64) error {
	now := time.Now()

	query := `
		UPDATE trade_sessions
		SET sizing_completed = 1,
		    sizing_method = ?,
		    sizing_entry_price = ?,
		    sizing_atr = ?,
		    sizing_k_multiple = ?,
		    sizing_stop_distance = ?,
		    sizing_initial_stop = ?,
		    sizing_shares = ?,
		    sizing_contracts = ?,
		    sizing_risk_dollars = ?,
		    sizing_delta = ?,
		    sizing_completed_at = ?,
		    current_step = 'HEAT'
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, method, entryPrice, atr, kMultiple, stopDistance, initialStop,
		shares, contracts, riskDollars, delta, now, id)
	if err != nil {
		return fmt.Errorf("failed to update session sizing: %w", err)
	}
	return nil
}

// UpdateSessionSizingWithPyramid updates position sizing including pyramid planning
func (db *DB) UpdateSessionSizingWithPyramid(id int, method string, entryPrice, atr, kMultiple, stopDistance, initialStop float64,
	shares, contracts int, riskDollars, delta float64,
	maxUnits int, addStepN, addPrice1, addPrice2, addPrice3 float64) error {
	now := time.Now()

	query := `
		UPDATE trade_sessions
		SET sizing_completed = 1,
		    sizing_method = ?,
		    sizing_entry_price = ?,
		    sizing_atr = ?,
		    sizing_k_multiple = ?,
		    sizing_stop_distance = ?,
		    sizing_initial_stop = ?,
		    sizing_shares = ?,
		    sizing_contracts = ?,
		    sizing_risk_dollars = ?,
		    sizing_delta = ?,
		    max_units = ?,
		    add_step_n = ?,
		    add_price_1 = ?,
		    add_price_2 = ?,
		    add_price_3 = ?,
		    current_units = 1,
		    sizing_completed_at = ?,
		    current_step = 'HEAT'
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, method, entryPrice, atr, kMultiple, stopDistance, initialStop,
		shares, contracts, riskDollars, delta,
		maxUnits, addStepN, addPrice1, addPrice2, addPrice3,
		now, id)
	if err != nil {
		return fmt.Errorf("failed to update session sizing with pyramid: %w", err)
	}
	return nil
}

// UpdateSessionHeat updates the heat check gate completion for a session
func (db *DB) UpdateSessionHeat(id int, status, bucket string,
	portfolioCurrent, portfolioNew, portfolioCap,
	bucketCurrent, bucketNew, bucketCap float64) error {
	now := time.Now()

	query := `
		UPDATE trade_sessions
		SET heat_completed = 1,
		    heat_status = ?,
		    heat_portfolio_current = ?,
		    heat_portfolio_new = ?,
		    heat_portfolio_cap = ?,
		    heat_bucket = ?,
		    heat_bucket_current = ?,
		    heat_bucket_new = ?,
		    heat_bucket_cap = ?,
		    heat_completed_at = ?,
		    current_step = 'ENTRY'
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, status, portfolioCurrent, portfolioNew, portfolioCap,
		bucket, bucketCurrent, bucketNew, bucketCap, now, id)
	if err != nil {
		return fmt.Errorf("failed to update session heat: %w", err)
	}
	return nil
}

// UpdateSessionEntry updates the final trade entry gate and marks session as completed
func (db *DB) UpdateSessionEntry(id int, decision string, decisionID int,
	gate1, gate2, gate3, gate4, gate5 bool) error {
	now := time.Now()

	gate1Int := 0
	if gate1 {
		gate1Int = 1
	}
	gate2Int := 0
	if gate2 {
		gate2Int = 1
	}
	gate3Int := 0
	if gate3 {
		gate3Int = 1
	}
	gate4Int := 0
	if gate4 {
		gate4Int = 1
	}
	gate5Int := 0
	if gate5 {
		gate5Int = 1
	}

	// Handle NULL decision_id (when decisionID is 0 or negative)
	var err error
	if decisionID > 0 {
		query := `
			UPDATE trade_sessions
			SET entry_completed = 1,
			    entry_decision = ?,
			    entry_decision_id = ?,
			    entry_gate1_pass = ?,
			    entry_gate2_pass = ?,
			    entry_gate3_pass = ?,
			    entry_gate4_pass = ?,
			    entry_gate5_pass = ?,
			    entry_completed_at = ?,
			    status = 'COMPLETED',
			    completed_at = ?
			WHERE id = ?
		`
		_, err = db.conn.Exec(query, decision, decisionID,
			gate1Int, gate2Int, gate3Int, gate4Int, gate5Int,
			now, now, id)
	} else {
		query := `
			UPDATE trade_sessions
			SET entry_completed = 1,
			    entry_decision = ?,
			    entry_decision_id = NULL,
			    entry_gate1_pass = ?,
			    entry_gate2_pass = ?,
			    entry_gate3_pass = ?,
			    entry_gate4_pass = ?,
			    entry_gate5_pass = ?,
			    entry_completed_at = ?,
			    status = 'COMPLETED',
			    completed_at = ?
			WHERE id = ?
		`
		_, err = db.conn.Exec(query, decision,
			gate1Int, gate2Int, gate3Int, gate4Int, gate5Int,
			now, now, id)
	}

	if err != nil {
		return fmt.Errorf("failed to update session entry: %w", err)
	}
	return nil
}

// ListActiveSessions returns all DRAFT sessions ordered by most recently updated
func (db *DB) ListActiveSessions() ([]*TradeSession, error) {
	query := `
		SELECT id FROM trade_sessions
		WHERE status = 'DRAFT'
		ORDER BY updated_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*TradeSession
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan session id: %w", err)
		}
		session, err := db.GetSession(id)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// ListSessionHistory returns all sessions (completed and draft) with optional limit
func (db *DB) ListSessionHistory(limit int) ([]*TradeSession, error) {
	query := `
		SELECT id FROM trade_sessions
		ORDER BY created_at DESC
		LIMIT ?
	`

	if limit <= 0 {
		limit = 100 // default limit
	}

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list session history: %w", err)
	}
	defer rows.Close()

	var sessions []*TradeSession
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan session id: %w", err)
		}
		session, err := db.GetSession(id)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// AbandonSession marks a session as abandoned
func (db *DB) AbandonSession(id int) error {
	query := `UPDATE trade_sessions SET status = 'ABANDONED' WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to abandon session: %w", err)
	}
	return nil
}

// CloneSession creates a new session based on an existing one
func (db *DB) CloneSession(sourceID int) (*TradeSession, error) {
	source, err := db.GetSession(sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source session: %w", err)
	}

	// Create new session with same ticker and strategy
	if source.Source == "PRESET" && source.CandidateID != nil && source.PresetID != nil {
		return db.CreateSessionFromPreset(
			source.Ticker,
			source.Strategy,
			*source.CandidateID,
			*source.PresetID,
			source.PresetName,
			source.ScanDate,
		)
	}

	return db.CreateSession(source.Ticker, source.Strategy)
}
