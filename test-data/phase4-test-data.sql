-- Phase 4 Integration Test Data
-- Purpose: Pre-populate database for consistent Phase 4 testing
-- Usage: sqlite3 trading.db < phase4-test-data.sql

-- Clear existing test data (optional - uncomment if needed)
-- DELETE FROM candidates;
-- DELETE FROM decisions;
-- DELETE FROM positions;

-- Test Data Set 1: Candidates for Gate 2 testing
-- These are the "approved" tickers for today
INSERT OR REPLACE INTO candidates (ticker, preset, date_imported)
VALUES
    ('AAPL', 'TEST', date('now')),
    ('MSFT', 'TEST', date('now')),
    ('NVDA', 'TEST', date('now')),
    ('SPY', 'TEST', date('now')),
    ('JPM', 'TEST', date('now'));

-- Test Data Set 2: Settings (verify defaults)
-- These should already exist from setup, but verify:
INSERT OR REPLACE INTO settings (key, value)
VALUES
    ('Equity_E', '10000'),
    ('RiskPct_r', '0.0075'),
    ('HeatCap_H_pct', '0.04'),
    ('BucketHeatCap_pct', '0.015'),
    ('StopMultiple_K', '2');

-- Test Data Set 3: Checklist evaluations
-- Pre-create GREEN banner evaluations with timestamps
-- Note: In actual testing, these should be created through Excel UI
-- This is just for reference / quick setup if needed

-- AAPL - GREEN banner (all 6 items)
-- Evaluated 3 minutes ago (past impulse brake)
INSERT OR REPLACE INTO checklist_evaluations (ticker, evaluation_timestamp, banner, missing_count, missing_items, allow_save)
VALUES
    ('AAPL', datetime('now', '-3 minutes'), 'GREEN', 0, '', 1);

-- MSFT - YELLOW banner (2 missing)
INSERT OR REPLACE INTO checklist_evaluations (ticker, evaluation_timestamp, banner, missing_count, missing_items, allow_save)
VALUES
    ('MSFT', datetime('now', '-3 minutes'), 'YELLOW', 2, 'Not overbought,Bucket OK', 0);

-- NVDA - GREEN banner (just evaluated - impulse brake active)
INSERT OR REPLACE INTO checklist_evaluations (ticker, evaluation_timestamp, banner, missing_count, missing_items, allow_save)
VALUES
    ('NVDA', datetime('now', '-30 seconds'), 'GREEN', 0, '', 1);

-- Test Data Set 4: Open position for heat testing
-- AAPL position: $75 risk in Tech/Comm bucket
-- This is used for Test 3.5 and 3.6 (heat with open positions)
-- Uncomment when ready to test cumulative heat:

-- INSERT INTO decisions (
--     ticker, entry_price, atr, method, banner,
--     risk_dollars, shares, contracts, bucket, preset,
--     decision_timestamp, correlation_id
-- ) VALUES (
--     'AAPL', 180.00, 1.5, 'stock', 'GREEN',
--     75.00, 25, 0, 'Tech/Comm', 'TEST',
--     datetime('now'), 'TEST-20251027-120000-HEAT'
-- );

-- INSERT INTO positions (
--     ticker, entry_price, shares, contracts, initial_stop,
--     risk_dollars, bucket, status, opened_at
-- ) VALUES (
--     'AAPL', 180.00, 25, 0, 177.00,
--     75.00, 'Tech/Comm', 'open', datetime('now')
-- );

-- Test Data Set 5: Bucket cooldown testing
-- Create a decision in Tech/Comm bucket to trigger Gate 4
-- Uncomment when ready to test Gate 4 (bucket cooldown):

-- INSERT INTO decisions (
--     ticker, entry_price, atr, method, banner,
--     risk_dollars, shares, contracts, bucket, preset,
--     decision_timestamp, correlation_id
-- ) VALUES (
--     'AAPL', 180.00, 1.5, 'stock', 'GREEN',
--     75.00, 25, 0, 'Tech/Comm', 'TEST',
--     datetime('now', '-1 hour'), 'TEST-20251027-110000-CD'
-- );

-- Verification Queries
-- Run these to confirm data loaded:

-- SELECT '=== CANDIDATES ===' AS section;
-- SELECT ticker, preset, date_imported FROM candidates ORDER BY ticker;

-- SELECT '=== SETTINGS ===' AS section;
-- SELECT key, value FROM settings ORDER BY key;

-- SELECT '=== CHECKLIST EVALUATIONS ===' AS section;
-- SELECT
--     ticker,
--     banner,
--     missing_count,
--     CAST((julianday('now') - julianday(evaluation_timestamp)) * 24 * 60 AS INTEGER) || ' minutes ago' AS age,
--     allow_save
-- FROM checklist_evaluations
-- ORDER BY evaluation_timestamp DESC;

-- SELECT '=== OPEN POSITIONS ===' AS section;
-- SELECT ticker, risk_dollars, bucket, status, opened_at
-- FROM positions
-- WHERE status = 'open'
-- ORDER BY opened_at;

-- Heat Calculation Verification
-- SELECT '=== HEAT SUMMARY ===' AS section;
-- SELECT
--     bucket,
--     SUM(risk_dollars) AS bucket_heat,
--     (SELECT SUM(risk_dollars) FROM positions WHERE status='open') AS portfolio_heat,
--     (SELECT value FROM settings WHERE key='Equity_E') AS equity,
--     (SELECT value FROM settings WHERE key='HeatCap_H_pct') AS portfolio_cap_pct,
--     (SELECT value FROM settings WHERE key='BucketHeatCap_pct') AS bucket_cap_pct
-- FROM positions
-- WHERE status = 'open'
-- GROUP BY bucket;
