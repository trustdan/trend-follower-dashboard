package storage

import (
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name     string
		ticker   string
		strategy string
		wantErr  bool
	}{
		{
			name:     "Create long breakout session",
			ticker:   "AAPL",
			strategy: StrategyLongBreakout,
			wantErr:  false,
		},
		{
			name:     "Create short breakout session",
			ticker:   "TSLA",
			strategy: StrategyShortBreakout,
			wantErr:  false,
		},
		{
			name:     "Create custom session",
			ticker:   "NVDA",
			strategy: StrategyCustom,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := db.CreateSession(tt.ticker, tt.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Verify session created correctly
			if session.Ticker != tt.ticker {
				t.Errorf("Expected ticker %s, got %s", tt.ticker, session.Ticker)
			}
			if session.Strategy != tt.strategy {
				t.Errorf("Expected strategy %s, got %s", tt.strategy, session.Strategy)
			}
			if session.Status != StatusDraft {
				t.Errorf("Expected status DRAFT, got %s", session.Status)
			}
			if session.CurrentStep != StepChecklist {
				t.Errorf("Expected current_step CHECKLIST, got %s", session.CurrentStep)
			}
			if session.ChecklistCompleted {
				t.Error("Expected checklist_completed to be false")
			}
			if session.SessionNum != session.ID {
				t.Errorf("Expected session_num (%d) to equal id (%d)", session.SessionNum, session.ID)
			}
		})
	}
}

func TestGetSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a session first
	created, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Retrieve it
	retrieved, err := db.GetSession(created.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}
	if retrieved.Ticker != created.Ticker {
		t.Errorf("Expected ticker %s, got %s", created.Ticker, retrieved.Ticker)
	}
	if retrieved.Strategy != created.Strategy {
		t.Errorf("Expected strategy %s, got %s", created.Strategy, retrieved.Strategy)
	}
}

func TestGetSessionByNum(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	created, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	retrieved, err := db.GetSessionByNum(created.SessionNum)
	if err != nil {
		t.Fatalf("Failed to get session by num: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}

	// Test non-existent session
	_, err = db.GetSessionByNum(99999)
	if err == nil {
		t.Error("Expected error for non-existent session, got nil")
	}
}

func TestUpdateSessionChecklist(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name              string
		banner            string
		missingCount      int
		qualityScore      int
		expectCompleted   bool
		expectCurrentStep string
	}{
		{
			name:              "GREEN banner completes checklist",
			banner:            "GREEN",
			missingCount:      0,
			qualityScore:      5,
			expectCompleted:   true,
			expectCurrentStep: StepSizing,
		},
		{
			name:              "YELLOW banner does not complete",
			banner:            "YELLOW",
			missingCount:      1,
			qualityScore:      3,
			expectCompleted:   false,
			expectCurrentStep: StepChecklist, // stays on checklist
		},
		{
			name:              "RED banner does not complete",
			banner:            "RED",
			missingCount:      2,
			qualityScore:      1,
			expectCompleted:   false,
			expectCurrentStep: StepChecklist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh session for each test
			s, err := db.CreateSession("AAPL", StrategyLongBreakout)
			if err != nil {
				t.Fatalf("Failed to create session: %v", err)
			}

			err = db.UpdateSessionChecklist(s.ID, tt.banner, tt.missingCount, tt.qualityScore)
			if err != nil {
				t.Errorf("UpdateSessionChecklist() error = %v", err)
				return
			}

			// Retrieve updated session
			updated, err := db.GetSession(s.ID)
			if err != nil {
				t.Fatalf("Failed to get updated session: %v", err)
			}

			if updated.ChecklistCompleted != tt.expectCompleted {
				t.Errorf("Expected checklist_completed %v, got %v", tt.expectCompleted, updated.ChecklistCompleted)
			}
			if updated.ChecklistBanner != tt.banner {
				t.Errorf("Expected banner %s, got %s", tt.banner, updated.ChecklistBanner)
			}
			if updated.ChecklistMissingCount != tt.missingCount {
				t.Errorf("Expected missing_count %d, got %d", tt.missingCount, updated.ChecklistMissingCount)
			}
			if updated.ChecklistQualityScore != tt.qualityScore {
				t.Errorf("Expected quality_score %d, got %d", tt.qualityScore, updated.ChecklistQualityScore)
			}
			if updated.CurrentStep != tt.expectCurrentStep {
				t.Errorf("Expected current_step %s, got %s", tt.expectCurrentStep, updated.CurrentStep)
			}
			if updated.ChecklistCompletedAt == nil {
				t.Error("Expected checklist_completed_at to be set")
			}
		})
	}
}

func TestUpdateSessionSizing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	session, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Complete checklist first (prerequisite)
	err = db.UpdateSessionChecklist(session.ID, "GREEN", 0, 5)
	if err != nil {
		t.Fatalf("Failed to update checklist: %v", err)
	}

	// Update sizing
	err = db.UpdateSessionSizing(
		session.ID,
		"stock",   // method
		180.0,     // entry price
		1.5,       // atr
		2.0,       // k multiple
		3.0,       // stop distance
		177.0,     // initial stop
		25,        // shares
		0,         // contracts
		75.0,      // risk dollars
		0.0,       // delta
	)
	if err != nil {
		t.Errorf("UpdateSessionSizing() error = %v", err)
		return
	}

	// Retrieve updated session
	updated, err := db.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}

	if !updated.SizingCompleted {
		t.Error("Expected sizing_completed to be true")
	}
	if updated.SizingMethod != "stock" {
		t.Errorf("Expected method 'stock', got %s", updated.SizingMethod)
	}
	if updated.SizingShares != 25 {
		t.Errorf("Expected 25 shares, got %d", updated.SizingShares)
	}
	if updated.SizingRiskDollars != 75.0 {
		t.Errorf("Expected risk $75, got $%.2f", updated.SizingRiskDollars)
	}
	if updated.CurrentStep != StepHeat {
		t.Errorf("Expected current_step HEAT, got %s", updated.CurrentStep)
	}
	if updated.SizingCompletedAt == nil {
		t.Error("Expected sizing_completed_at to be set")
	}
}

func TestUpdateSessionHeat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	session, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Complete checklist and sizing first (prerequisites)
	err = db.UpdateSessionChecklist(session.ID, "GREEN", 0, 5)
	if err != nil {
		t.Fatalf("Failed to update checklist: %v", err)
	}
	err = db.UpdateSessionSizing(session.ID, "stock", 180.0, 1.5, 2.0, 3.0, 177.0, 25, 0, 75.0, 0.0)
	if err != nil {
		t.Fatalf("Failed to update sizing: %v", err)
	}

	// Update heat
	err = db.UpdateSessionHeat(
		session.ID,
		"OK",       // status
		"Tech/Comm", // bucket
		2100.0,     // portfolio current
		2175.0,     // portfolio new
		4000.0,     // portfolio cap
		1400.0,     // bucket current
		1475.0,     // bucket new
		1500.0,     // bucket cap
	)
	if err != nil {
		t.Errorf("UpdateSessionHeat() error = %v", err)
		return
	}

	// Retrieve updated session
	updated, err := db.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}

	if !updated.HeatCompleted {
		t.Error("Expected heat_completed to be true")
	}
	if updated.HeatStatus != "OK" {
		t.Errorf("Expected status 'OK', got %s", updated.HeatStatus)
	}
	if updated.HeatBucket != "Tech/Comm" {
		t.Errorf("Expected bucket 'Tech/Comm', got %s", updated.HeatBucket)
	}
	if updated.HeatPortfolioNew != 2175.0 {
		t.Errorf("Expected portfolio_new $2175, got $%.2f", updated.HeatPortfolioNew)
	}
	if updated.CurrentStep != StepEntry {
		t.Errorf("Expected current_step ENTRY, got %s", updated.CurrentStep)
	}
	if updated.HeatCompletedAt == nil {
		t.Error("Expected heat_completed_at to be set")
	}
}

func TestUpdateSessionEntry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	session, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Complete all prerequisites
	err = db.UpdateSessionChecklist(session.ID, "GREEN", 0, 5)
	if err != nil {
		t.Fatalf("Failed to update checklist: %v", err)
	}
	err = db.UpdateSessionSizing(session.ID, "stock", 180.0, 1.5, 2.0, 3.0, 177.0, 25, 0, 75.0, 0.0)
	if err != nil {
		t.Fatalf("Failed to update sizing: %v", err)
	}
	err = db.UpdateSessionHeat(session.ID, "OK", "Tech/Comm", 2100.0, 2175.0, 4000.0, 1400.0, 1475.0, 1500.0)
	if err != nil {
		t.Fatalf("Failed to update heat: %v", err)
	}

	// Update entry (all gates pass)
	// Note: using decision_id=0 since we don't have a real decision in test
	err = db.UpdateSessionEntry(
		session.ID,
		"GO",  // decision
		0,     // decision_id (will be NULL in DB)
		true,  // gate1
		true,  // gate2
		true,  // gate3
		true,  // gate4
		true,  // gate5
	)
	if err != nil {
		t.Errorf("UpdateSessionEntry() error = %v", err)
		return
	}

	// Retrieve updated session
	updated, err := db.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}

	if !updated.EntryCompleted {
		t.Error("Expected entry_completed to be true")
	}
	if updated.EntryDecision != "GO" {
		t.Errorf("Expected decision 'GO', got %s", updated.EntryDecision)
	}
	// decision_id is NULL in test
	if updated.EntryDecisionID != nil {
		t.Errorf("Expected decision_id to be NULL, got %v", *updated.EntryDecisionID)
	}
	if updated.Status != StatusCompleted {
		t.Errorf("Expected status COMPLETED, got %s", updated.Status)
	}
	if updated.CompletedAt == nil {
		t.Error("Expected completed_at to be set")
	}
	if updated.EntryCompletedAt == nil {
		t.Error("Expected entry_completed_at to be set")
	}

	// Verify all gates
	if updated.EntryGate1Pass == nil || !*updated.EntryGate1Pass {
		t.Error("Expected gate1 to pass")
	}
	if updated.EntryGate2Pass == nil || !*updated.EntryGate2Pass {
		t.Error("Expected gate2 to pass")
	}
	if updated.EntryGate3Pass == nil || !*updated.EntryGate3Pass {
		t.Error("Expected gate3 to pass")
	}
	if updated.EntryGate4Pass == nil || !*updated.EntryGate4Pass {
		t.Error("Expected gate4 to pass")
	}
	if updated.EntryGate5Pass == nil || !*updated.EntryGate5Pass {
		t.Error("Expected gate5 to pass")
	}
}

func TestListActiveSessions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create multiple sessions with different statuses
	s1, _ := db.CreateSession("AAPL", StrategyLongBreakout)
	s2, _ := db.CreateSession("TSLA", StrategyShortBreakout)
	s3, _ := db.CreateSession("NVDA", StrategyLongBreakout)

	// Complete session 3
	db.UpdateSessionChecklist(s3.ID, "GREEN", 0, 5)
	db.UpdateSessionSizing(s3.ID, "stock", 180.0, 1.5, 2.0, 3.0, 177.0, 25, 0, 75.0, 0.0)
	db.UpdateSessionHeat(s3.ID, "OK", "Tech/Comm", 2100.0, 2175.0, 4000.0, 1400.0, 1475.0, 1500.0)
	db.UpdateSessionEntry(s3.ID, "GO", 0, true, true, true, true, true)

	// Sleep briefly to ensure different updated_at times
	time.Sleep(10 * time.Millisecond)

	// Update s2 to make it most recent
	db.UpdateSessionChecklist(s2.ID, "GREEN", 0, 5)

	// List active sessions
	sessions, err := db.ListActiveSessions()
	if err != nil {
		t.Fatalf("ListActiveSessions() error = %v", err)
	}

	// Should only return s1 and s2 (not s3 which is COMPLETED)
	if len(sessions) != 2 {
		t.Errorf("Expected 2 active sessions, got %d", len(sessions))
	}

	// Verify ordering (most recently updated first)
	if sessions[0].ID != s2.ID {
		t.Errorf("Expected first session to be s2 (most recent), got ID %d", sessions[0].ID)
	}
	if sessions[1].ID != s1.ID {
		t.Errorf("Expected second session to be s1, got ID %d", sessions[1].ID)
	}
}

func TestListSessionHistory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create multiple sessions
	s1, _ := db.CreateSession("AAPL", StrategyLongBreakout)
	time.Sleep(10 * time.Millisecond)
	s2, _ := db.CreateSession("TSLA", StrategyShortBreakout)
	time.Sleep(10 * time.Millisecond)
	s3, _ := db.CreateSession("NVDA", StrategyCustom)

	// Complete one session
	db.UpdateSessionChecklist(s2.ID, "GREEN", 0, 5)
	db.UpdateSessionSizing(s2.ID, "stock", 180.0, 1.5, 2.0, 3.0, 177.0, 25, 0, 75.0, 0.0)
	db.UpdateSessionHeat(s2.ID, "OK", "Tech/Comm", 2100.0, 2175.0, 4000.0, 1400.0, 1475.0, 1500.0)
	db.UpdateSessionEntry(s2.ID, "GO", 0, true, true, true, true, true)

	// List all sessions
	sessions, err := db.ListSessionHistory(10)
	if err != nil {
		t.Fatalf("ListSessionHistory() error = %v", err)
	}

	// Should return all 3 sessions
	if len(sessions) != 3 {
		t.Errorf("Expected 3 sessions in history, got %d", len(sessions))
	}

	// Verify ordering (newest first)
	if sessions[0].ID != s3.ID {
		t.Errorf("Expected first session to be s3 (newest), got ID %d", sessions[0].ID)
	}
	if sessions[1].ID != s2.ID {
		t.Errorf("Expected second session to be s2, got ID %d", sessions[1].ID)
	}
	if sessions[2].ID != s1.ID {
		t.Errorf("Expected third session to be s1 (oldest), got ID %d", sessions[2].ID)
	}

	// Verify statuses
	foundCompleted := false
	for _, s := range sessions {
		if s.Status == StatusCompleted {
			foundCompleted = true
		}
	}
	if !foundCompleted {
		t.Error("Expected to find at least one COMPLETED session")
	}
}

func TestAbandonSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	session, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	err = db.AbandonSession(session.ID)
	if err != nil {
		t.Errorf("AbandonSession() error = %v", err)
		return
	}

	// Verify status changed
	updated, err := db.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if updated.Status != StatusAbandoned {
		t.Errorf("Expected status ABANDONED, got %s", updated.Status)
	}
}

func TestCloneSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create and complete a session
	original, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Complete all gates
	db.UpdateSessionChecklist(original.ID, "GREEN", 0, 5)
	db.UpdateSessionSizing(original.ID, "stock", 180.0, 1.5, 2.0, 3.0, 177.0, 25, 0, 75.0, 0.0)
	db.UpdateSessionHeat(original.ID, "OK", "Tech/Comm", 2100.0, 2175.0, 4000.0, 1400.0, 1475.0, 1500.0)
	db.UpdateSessionEntry(original.ID, "GO", 0, true, true, true, true, true)

	// Clone the session
	cloned, err := db.CloneSession(original.ID)
	if err != nil {
		t.Errorf("CloneSession() error = %v", err)
		return
	}

	// Verify clone has same ticker and strategy but is a new DRAFT session
	if cloned.Ticker != original.Ticker {
		t.Errorf("Expected ticker %s, got %s", original.Ticker, cloned.Ticker)
	}
	if cloned.Strategy != original.Strategy {
		t.Errorf("Expected strategy %s, got %s", original.Strategy, cloned.Strategy)
	}
	if cloned.Status != StatusDraft {
		t.Errorf("Expected status DRAFT, got %s", cloned.Status)
	}
	if cloned.ID == original.ID {
		t.Error("Cloned session should have different ID")
	}
	if cloned.SessionNum == original.SessionNum {
		t.Error("Cloned session should have different session_num")
	}
	if cloned.ChecklistCompleted {
		t.Error("Cloned session should not have completed checklist")
	}
	if cloned.SizingCompleted {
		t.Error("Cloned session should not have completed sizing")
	}
}

func TestSessionWorkflowProgression(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create session
	session, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Verify initial state
	if session.CurrentStep != StepChecklist {
		t.Errorf("Expected initial step CHECKLIST, got %s", session.CurrentStep)
	}

	// Complete checklist with GREEN
	err = db.UpdateSessionChecklist(session.ID, "GREEN", 0, 5)
	if err != nil {
		t.Fatalf("Failed to update checklist: %v", err)
	}

	// Verify progressed to SIZING
	session, _ = db.GetSession(session.ID)
	if session.CurrentStep != StepSizing {
		t.Errorf("Expected step SIZING after GREEN checklist, got %s", session.CurrentStep)
	}

	// Complete sizing
	err = db.UpdateSessionSizing(session.ID, "stock", 180.0, 1.5, 2.0, 3.0, 177.0, 25, 0, 75.0, 0.0)
	if err != nil {
		t.Fatalf("Failed to update sizing: %v", err)
	}

	// Verify progressed to HEAT
	session, _ = db.GetSession(session.ID)
	if session.CurrentStep != StepHeat {
		t.Errorf("Expected step HEAT after sizing, got %s", session.CurrentStep)
	}

	// Complete heat
	err = db.UpdateSessionHeat(session.ID, "OK", "Tech/Comm", 2100.0, 2175.0, 4000.0, 1400.0, 1475.0, 1500.0)
	if err != nil {
		t.Fatalf("Failed to update heat: %v", err)
	}

	// Verify progressed to ENTRY
	session, _ = db.GetSession(session.ID)
	if session.CurrentStep != StepEntry {
		t.Errorf("Expected step ENTRY after heat, got %s", session.CurrentStep)
	}

	// Complete entry
	err = db.UpdateSessionEntry(session.ID, "GO", 0, true, true, true, true, true)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	// Verify completed
	session, _ = db.GetSession(session.ID)
	if session.Status != StatusCompleted {
		t.Errorf("Expected status COMPLETED after entry, got %s", session.Status)
	}
}

// TestSessionsWithOptionsMetadata tests creating and retrieving sessions with options data
func TestSessionsWithOptionsMetadata(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Run migration to add options columns (this should be idempotent)
	// Note: In production, migrations are run via separate migration tool
	// For tests, we're adding columns inline (if they don't exist, SQLite will error - that's OK)

	t.Run("Create Session With Options Metadata", func(t *testing.T) {
		session, err := db.CreateSession("MSFT", StrategyLongBreakout)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Update session with options metadata (Long Call strategy)
		query := `
			UPDATE trade_sessions
			SET instrument_type = ?,
			    options_strategy = ?,
			    entry_date = ?,
			    primary_expiration_date = ?,
			    dte = ?,
			    roll_threshold_dte = ?,
			    time_exit_mode = ?,
			    net_debit = ?,
			    max_profit = ?,
			    max_loss = ?,
			    underlying_at_entry = ?,
			    max_units = ?,
			    add_step_n = ?,
			    current_units = ?,
			    entry_lookback = ?,
			    exit_lookback = ?
			WHERE id = ?
		`

		_, err = db.conn.Exec(query,
			InstrumentOption,
			StrategyLongCall,
			"2025-10-30",
			"2025-12-19",
			60,
			21,
			TimeExitRoll,
			1250.0,  // net debit (5 contracts × $2.50 × 100)
			99999.0, // unlimited max profit for long call
			1250.0,  // max loss = debit paid
			175.0,   // underlying at entry
			4,       // max units
			0.5,     // add every 0.5N
			1,       // currently 1 unit
			55,      // System-2 (55-bar entry)
			10,      // 10-bar exit
			session.ID,
		)

		if err != nil {
			t.Fatalf("Failed to update session with options metadata: %v", err)
		}

		// Retrieve the session and verify
		retrieved, err := db.GetSession(session.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve session: %v", err)
		}

		if retrieved.InstrumentType != InstrumentOption {
			t.Errorf("Expected instrument type %s, got %s", InstrumentOption, retrieved.InstrumentType)
		}

		if retrieved.OptionsStrategy != StrategyLongCall {
			t.Errorf("Expected options strategy %s, got %s", StrategyLongCall, retrieved.OptionsStrategy)
		}

		if retrieved.EntryDate != "2025-10-30" {
			t.Errorf("Expected entry date 2025-10-30, got %s", retrieved.EntryDate)
		}

		if retrieved.PrimaryExpirationDate != "2025-12-19" {
			t.Errorf("Expected expiration 2025-12-19, got %s", retrieved.PrimaryExpirationDate)
		}

		if retrieved.DTE != 60 {
			t.Errorf("Expected DTE 60, got %d", retrieved.DTE)
		}

		if retrieved.RollThresholdDTE != 21 {
			t.Errorf("Expected roll threshold 21, got %d", retrieved.RollThresholdDTE)
		}

		if retrieved.TimeExitMode != TimeExitRoll {
			t.Errorf("Expected time exit mode %s, got %s", TimeExitRoll, retrieved.TimeExitMode)
		}

		if retrieved.NetDebit != 1250.0 {
			t.Errorf("Expected net debit 1250.0, got %.2f", retrieved.NetDebit)
		}

		if retrieved.MaxUnits != 4 {
			t.Errorf("Expected max units 4, got %d", retrieved.MaxUnits)
		}

		if retrieved.AddStepN != 0.5 {
			t.Errorf("Expected add step 0.5, got %.2f", retrieved.AddStepN)
		}

		if retrieved.CurrentUnits != 1 {
			t.Errorf("Expected current units 1, got %d", retrieved.CurrentUnits)
		}

		if retrieved.EntryLookback != 55 {
			t.Errorf("Expected entry lookback 55, got %d", retrieved.EntryLookback)
		}

		if retrieved.ExitLookback != 10 {
			t.Errorf("Expected exit lookback 10, got %d", retrieved.ExitLookback)
		}
	})

	t.Run("Stock Trade Backward Compatibility", func(t *testing.T) {
		session, err := db.CreateSession("TSLA", StrategyLongBreakout)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Stock trades should work without options metadata
		retrieved, err := db.GetSession(session.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve session: %v", err)
		}

		// These should be zero/empty for stock trades
		if retrieved.OptionsStrategy != "" {
			t.Errorf("Expected empty options strategy for stock, got %s", retrieved.OptionsStrategy)
		}

		if retrieved.DTE != 0 {
			t.Errorf("Expected DTE 0 for stock, got %d", retrieved.DTE)
		}
	})
}

// TestPyramidPricing tests pyramid add-on price calculation
func TestPyramidPricing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	session, err := db.CreateSession("AAPL", StrategyLongBreakout)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Calculate pyramid prices
	entry := 180.0
	atr := 1.5
	addStepN := 0.5

	addPrice1 := entry + (addStepN * atr)      // 180 + 0.75 = 180.75
	addPrice2 := entry + (2 * addStepN * atr)  // 180 + 1.50 = 181.50
	addPrice3 := entry + (3 * addStepN * atr)  // 180 + 2.25 = 182.25

	query := `
		UPDATE trade_sessions
		SET max_units = ?,
		    add_step_n = ?,
		    current_units = ?,
		    add_price_1 = ?,
		    add_price_2 = ?,
		    add_price_3 = ?
		WHERE id = ?
	`

	_, err = db.conn.Exec(query, 4, addStepN, 1, addPrice1, addPrice2, addPrice3, session.ID)
	if err != nil {
		t.Fatalf("Failed to update pyramid pricing: %v", err)
	}

	retrieved, err := db.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve session: %v", err)
	}

	if retrieved.AddPrice1 != addPrice1 {
		t.Errorf("Expected add price 1 %.2f, got %.2f", addPrice1, retrieved.AddPrice1)
	}
	if retrieved.AddPrice2 != addPrice2 {
		t.Errorf("Expected add price 2 %.2f, got %.2f", addPrice2, retrieved.AddPrice2)
	}
	if retrieved.AddPrice3 != addPrice3 {
		t.Errorf("Expected add price 3 %.2f, got %.2f", addPrice3, retrieved.AddPrice3)
	}
	if retrieved.MaxUnits != 4 {
		t.Errorf("Expected max units 4, got %d", retrieved.MaxUnits)
	}
	if retrieved.CurrentUnits != 1 {
		t.Errorf("Expected current units 1, got %d", retrieved.CurrentUnits)
	}
}
