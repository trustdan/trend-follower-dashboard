-- Migration 002: Add Options Trading Metadata
-- This migration extends the database to support options trading with all 26 strategies
-- while maintaining anti-impulsivity compliance and full backward compatibility.

-- ============================================================================
-- PHASE 1: Extend trade_sessions table with options metadata
-- ============================================================================

-- Core options identification
ALTER TABLE trade_sessions ADD COLUMN instrument_type TEXT DEFAULT 'STOCK' CHECK (instrument_type IN ('STOCK', 'OPTION'));
ALTER TABLE trade_sessions ADD COLUMN options_strategy TEXT; -- See constants in sessions.go
ALTER TABLE trade_sessions ADD COLUMN entry_date TEXT; -- ISO date (YYYY-MM-DD)
ALTER TABLE trade_sessions ADD COLUMN primary_expiration_date TEXT; -- Primary leg expiration
ALTER TABLE trade_sessions ADD COLUMN dte INTEGER; -- Days to expiration at entry
ALTER TABLE trade_sessions ADD COLUMN roll_threshold_dte INTEGER DEFAULT 21; -- DTE to roll/close
ALTER TABLE trade_sessions ADD COLUMN time_exit_mode TEXT DEFAULT 'Close' CHECK (time_exit_mode IN ('None', 'Close', 'Roll'));

-- Multi-leg structure (JSON for flexibility)
-- Example: [{"type":"CALL","strike":180,"exp":"2025-12-19","qty":1,"action":"BUY"}]
ALTER TABLE trade_sessions ADD COLUMN legs_json TEXT;

-- Aggregate pricing (calculated from legs)
ALTER TABLE trade_sessions ADD COLUMN net_debit REAL; -- Total debit paid (negative = credit)
ALTER TABLE trade_sessions ADD COLUMN max_profit REAL; -- Maximum theoretical profit
ALTER TABLE trade_sessions ADD COLUMN max_loss REAL; -- Maximum theoretical loss
ALTER TABLE trade_sessions ADD COLUMN breakeven_lower REAL; -- Lower breakeven price
ALTER TABLE trade_sessions ADD COLUMN breakeven_upper REAL; -- Upper breakeven price
ALTER TABLE trade_sessions ADD COLUMN underlying_at_entry REAL; -- Stock price at entry

-- Pyramiding (Van Tharp method: add every 0.5N up to 4 units)
ALTER TABLE trade_sessions ADD COLUMN max_units INTEGER DEFAULT 4; -- Maximum pyramid units
ALTER TABLE trade_sessions ADD COLUMN add_step_n REAL DEFAULT 0.5; -- Add every X * N
ALTER TABLE trade_sessions ADD COLUMN current_units INTEGER DEFAULT 0; -- Current units (0-4)
ALTER TABLE trade_sessions ADD COLUMN add_price_1 REAL; -- Entry + 0.5N
ALTER TABLE trade_sessions ADD COLUMN add_price_2 REAL; -- Entry + 1.0N
ALTER TABLE trade_sessions ADD COLUMN add_price_3 REAL; -- Entry + 1.5N

-- Breakout system parameters (for documentation/audit)
ALTER TABLE trade_sessions ADD COLUMN entry_lookback INTEGER; -- 20 or 55 for System-1/System-2
ALTER TABLE trade_sessions ADD COLUMN exit_lookback INTEGER DEFAULT 10; -- 10-bar exit

-- ============================================================================
-- PHASE 2: Extend positions table with options metadata
-- ============================================================================

ALTER TABLE positions ADD COLUMN instrument_type TEXT DEFAULT 'STOCK' CHECK (instrument_type IN ('STOCK', 'OPTION'));
ALTER TABLE positions ADD COLUMN options_strategy TEXT;
ALTER TABLE positions ADD COLUMN entry_date TEXT;
ALTER TABLE positions ADD COLUMN primary_expiration_date TEXT;
ALTER TABLE positions ADD COLUMN dte INTEGER;
ALTER TABLE positions ADD COLUMN legs_json TEXT; -- JSON array of legs
ALTER TABLE positions ADD COLUMN net_debit REAL;
ALTER TABLE positions ADD COLUMN max_profit REAL;
ALTER TABLE positions ADD COLUMN max_loss REAL;
ALTER TABLE positions ADD COLUMN breakeven_lower REAL;
ALTER TABLE positions ADD COLUMN breakeven_upper REAL;
ALTER TABLE positions ADD COLUMN underlying_at_entry REAL;
ALTER TABLE positions ADD COLUMN max_units INTEGER DEFAULT 4;
ALTER TABLE positions ADD COLUMN current_units INTEGER DEFAULT 1;
ALTER TABLE positions ADD COLUMN add_step_n REAL DEFAULT 0.5;

-- ============================================================================
-- PHASE 3: Create trade_history table for calendar view
-- ============================================================================

CREATE TABLE IF NOT EXISTS trade_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER,
    ticker TEXT NOT NULL,
    strategy TEXT NOT NULL, -- LONG_BREAKOUT, SHORT_BREAKOUT, CUSTOM
    breakout_system TEXT, -- SYSTEM_1, SYSTEM_2, CUSTOM
    options_strategy TEXT, -- LONG_CALL, IRON_CONDOR, etc.
    instrument_type TEXT DEFAULT 'STOCK',
    sector TEXT, -- Tech/Comm, Finance, etc.
    bucket TEXT, -- For heat tracking
    entry_date TEXT NOT NULL, -- YYYY-MM-DD
    expiration_date TEXT, -- YYYY-MM-DD (NULL for stocks)
    exit_date TEXT, -- YYYY-MM-DD (NULL if still open)
    status TEXT NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED', 'ROLLED')),
    dte INTEGER, -- Days to expiration at entry
    contracts INTEGER, -- Number of option contracts
    shares INTEGER, -- Number of shares
    risk_dollars REAL,
    entry_price REAL,
    exit_price REAL,
    pnl REAL,
    outcome TEXT CHECK (outcome IN ('WIN', 'LOSS', 'SCRATCH')),
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES trade_sessions(id)
);

CREATE INDEX IF NOT EXISTS idx_trade_history_entry_date ON trade_history(entry_date);
CREATE INDEX IF NOT EXISTS idx_trade_history_sector ON trade_history(sector, entry_date);
CREATE INDEX IF NOT EXISTS idx_trade_history_bucket ON trade_history(bucket, entry_date);
CREATE INDEX IF NOT EXISTS idx_trade_history_status ON trade_history(status);
CREATE INDEX IF NOT EXISTS idx_trade_history_ticker ON trade_history(ticker);
CREATE INDEX IF NOT EXISTS idx_trade_history_session ON trade_history(session_id);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS trg_trade_history_updated_at
AFTER UPDATE ON trade_history
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE trade_history
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

-- ============================================================================
-- BACKWARD COMPATIBILITY & MIGRATION NOTES
-- ============================================================================

-- For existing sessions in database:
-- - instrument_type defaults to 'STOCK' (NULL interpreted as stock/ETF)
-- - options_strategy defaults to NULL (stock trades)
-- - entry_date defaults to NULL (can backfill with created_at if needed)
-- - max_units defaults to 4
-- - add_step_n defaults to 0.5
-- - time_exit_mode defaults to 'Close'
-- - roll_threshold_dte defaults to 21

-- No data migration required - all new columns are nullable or have defaults
-- No breaking changes to existing data or queries
