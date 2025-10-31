-- Migration: Add options trading columns to trade_sessions table
-- Date: 2025-10-30
-- Purpose: Support options trading enhancement (Phase 1-7)

-- Add options metadata columns
-- Note: SQLite doesn't support CHECK constraints in ALTER TABLE ADD COLUMN in older versions
-- We add the column first, then add the constraint separately if needed
ALTER TABLE trade_sessions ADD COLUMN instrument_type TEXT DEFAULT 'STOCK';
ALTER TABLE trade_sessions ADD COLUMN options_strategy TEXT;
ALTER TABLE trade_sessions ADD COLUMN entry_date TEXT;
ALTER TABLE trade_sessions ADD COLUMN primary_expiration_date TEXT;
ALTER TABLE trade_sessions ADD COLUMN dte INTEGER;
ALTER TABLE trade_sessions ADD COLUMN roll_threshold_dte INTEGER DEFAULT 21;
ALTER TABLE trade_sessions ADD COLUMN time_exit_mode TEXT DEFAULT 'Close' CHECK (time_exit_mode IN ('None', 'Close', 'Roll'));
ALTER TABLE trade_sessions ADD COLUMN legs_json TEXT;
ALTER TABLE trade_sessions ADD COLUMN net_debit REAL;
ALTER TABLE trade_sessions ADD COLUMN max_profit REAL;
ALTER TABLE trade_sessions ADD COLUMN max_loss REAL;
ALTER TABLE trade_sessions ADD COLUMN breakeven_lower REAL;
ALTER TABLE trade_sessions ADD COLUMN breakeven_upper REAL;
ALTER TABLE trade_sessions ADD COLUMN underlying_at_entry REAL;

-- Add pyramid tracking columns
ALTER TABLE trade_sessions ADD COLUMN max_units INTEGER DEFAULT 4;
ALTER TABLE trade_sessions ADD COLUMN add_step_n REAL DEFAULT 0.5;
ALTER TABLE trade_sessions ADD COLUMN current_units INTEGER DEFAULT 0;
ALTER TABLE trade_sessions ADD COLUMN add_price_1 REAL;
ALTER TABLE trade_sessions ADD COLUMN add_price_2 REAL;
ALTER TABLE trade_sessions ADD COLUMN add_price_3 REAL;

-- Add breakout system columns
ALTER TABLE trade_sessions ADD COLUMN entry_lookback INTEGER;
ALTER TABLE trade_sessions ADD COLUMN exit_lookback INTEGER DEFAULT 10;
