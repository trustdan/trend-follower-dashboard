-- Migration: Add Trade Sessions Table
-- Version: 001
-- Created: 2025-10-30
-- Description: Adds trade_sessions table for cohesive trade evaluation workflow

-- ============================================================================
-- Trade Sessions Table
-- ============================================================================
-- Stores the complete lifecycle of a trade evaluation session, linking together
-- Checklist → Sizing → Heat → Entry steps with full audit trail.

CREATE TABLE IF NOT EXISTS trade_sessions (
    -- Primary key
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    -- Human-facing identifiers
    session_num INTEGER GENERATED ALWAYS AS (id) STORED UNIQUE,
    ticker TEXT NOT NULL,
    strategy TEXT NOT NULL,

    -- Provenance from FINVIZ/candidates
    source TEXT NOT NULL DEFAULT 'MANUAL' CHECK (source IN ('MANUAL', 'PRESET', 'CUSTOM')),
    candidate_id INTEGER,
    preset_id INTEGER,
    preset_name TEXT,
    scan_date TEXT,

    -- Workflow state
    status TEXT NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'EVALUATING', 'COMPLETED', 'ABANDONED')),
    current_step TEXT NOT NULL DEFAULT 'CHECKLIST' CHECK (current_step IN ('CHECKLIST', 'SIZING', 'HEAT', 'ENTRY')),

    -- Gate 1: Checklist completion
    checklist_completed INTEGER NOT NULL DEFAULT 0 CHECK (checklist_completed IN (0,1)),
    checklist_banner TEXT,
    checklist_missing_count INTEGER DEFAULT 0,
    checklist_quality_score INTEGER DEFAULT 0,
    checklist_completed_at DATETIME,

    -- Gate 2: Position Sizing completion
    sizing_completed INTEGER NOT NULL DEFAULT 0 CHECK (sizing_completed IN (0,1)),
    sizing_method TEXT,
    sizing_entry_price REAL,
    sizing_atr REAL,
    sizing_k_multiple REAL,
    sizing_stop_distance REAL,
    sizing_initial_stop REAL,
    sizing_shares INTEGER,
    sizing_contracts INTEGER,
    sizing_risk_dollars REAL,
    sizing_delta REAL,
    sizing_completed_at DATETIME,

    -- Gate 3: Heat Check completion
    heat_completed INTEGER NOT NULL DEFAULT 0 CHECK (heat_completed IN (0,1)),
    heat_status TEXT,
    heat_portfolio_current REAL,
    heat_portfolio_new REAL,
    heat_portfolio_cap REAL,
    heat_bucket TEXT,
    heat_bucket_current REAL,
    heat_bucket_new REAL,
    heat_bucket_cap REAL,
    heat_completed_at DATETIME,

    -- Gate 4: Trade Entry completion (Final decision)
    entry_completed INTEGER NOT NULL DEFAULT 0 CHECK (entry_completed IN (0,1)),
    entry_decision TEXT,
    entry_decision_id INTEGER,
    entry_gate1_pass INTEGER,
    entry_gate2_pass INTEGER,
    entry_gate3_pass INTEGER,
    entry_gate4_pass INTEGER,
    entry_gate5_pass INTEGER,
    entry_completed_at DATETIME,

    -- Audit trail
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,

    -- Foreign keys
    FOREIGN KEY (candidate_id) REFERENCES candidates(id),
    FOREIGN KEY (preset_id) REFERENCES presets(id),
    FOREIGN KEY (entry_decision_id) REFERENCES decisions(id)
);

-- ============================================================================
-- Indexes for Performance
-- ============================================================================


-- Index on session_num (unique, for fast lookup by user-facing ID)
CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_session_num ON trade_sessions(session_num);

-- Index on ticker (for filtering by ticker)
CREATE INDEX IF NOT EXISTS idx_sessions_ticker ON trade_sessions(ticker);

-- Index on strategy (for filtering by strategy type)
CREATE INDEX IF NOT EXISTS idx_sessions_strategy ON trade_sessions(strategy);

-- Index on status (for filtering DRAFT vs COMPLETED sessions)
CREATE INDEX IF NOT EXISTS idx_sessions_status ON trade_sessions(status);

-- Index on created_at (for sorting by date, descending)
CREATE INDEX IF NOT EXISTS idx_sessions_created ON trade_sessions(created_at DESC);

-- Index on updated_at (for "recent sessions" queries)
CREATE INDEX IF NOT EXISTS idx_sessions_updated ON trade_sessions(updated_at DESC);

-- Composite index for active session queries (status + updated_at)
CREATE INDEX IF NOT EXISTS idx_sessions_active ON trade_sessions(status, updated_at DESC);

-- Index for provenance lookups
CREATE INDEX IF NOT EXISTS idx_sessions_candidate ON trade_sessions(candidate_id);

-- ============================================================================
-- Triggers for Auto-Update
-- ============================================================================

-- Trigger: Auto-update `updated_at` timestamp on any UPDATE
CREATE TRIGGER IF NOT EXISTS trg_sessions_updated_at
AFTER UPDATE ON trade_sessions
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE trade_sessions
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

-- ============================================================================
-- Migration Verification
-- ============================================================================

-- Verify table created successfully
SELECT 'Migration 001: trade_sessions table created' AS status;

-- ============================================================================
-- Rollback Script (for reference only, not executed)
-- ============================================================================
-- To rollback this migration, run:
--
--   DROP TABLE IF EXISTS trade_sessions;
--   DROP INDEX IF EXISTS idx_sessions_session_num;
--   DROP INDEX IF EXISTS idx_sessions_ticker;
--   DROP INDEX IF EXISTS idx_sessions_strategy;
--   DROP INDEX IF EXISTS idx_sessions_status;
--   DROP INDEX IF EXISTS idx_sessions_created;
--   DROP INDEX IF EXISTS idx_sessions_updated;
--   DROP INDEX IF EXISTS idx_sessions_active;
--   DROP TRIGGER IF EXISTS trg_sessions_updated_at;
--
-- ============================================================================
