package storage

// schema defines the complete database structure for Trading Engine v3
const schema = `
-- Settings table: key-value configuration
CREATE TABLE IF NOT EXISTS settings (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Presets table: FINVIZ screen configurations
CREATE TABLE IF NOT EXISTS presets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	query_string TEXT NOT NULL,
	active INTEGER DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Candidates table: Daily screening results
CREATE TABLE IF NOT EXISTS candidates (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL,
	ticker TEXT NOT NULL,
	preset_id INTEGER,
	sector TEXT,
	bucket TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (preset_id) REFERENCES presets(id),
	UNIQUE(date, ticker, preset_id)
);

CREATE INDEX IF NOT EXISTS idx_candidates_date ON candidates(date);
CREATE INDEX IF NOT EXISTS idx_candidates_ticker ON candidates(ticker);

-- Checklist evaluations table: Banner determination history
CREATE TABLE IF NOT EXISTS checklist_evaluations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ticker TEXT NOT NULL,
	banner TEXT NOT NULL,
	missing_count INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_checklist_ticker ON checklist_evaluations(ticker);

-- Decisions table: Trade evaluation history
CREATE TABLE IF NOT EXISTS decisions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL,
	ticker TEXT NOT NULL,
	action TEXT NOT NULL,
	entry REAL,
	atr REAL,
	stop_distance REAL,
	initial_stop REAL,
	shares INTEGER DEFAULT 0,
	contracts INTEGER DEFAULT 0,
	risk_dollars REAL,
	banner TEXT NOT NULL,
	method TEXT,
	delta REAL,
	max_loss REAL,
	bucket TEXT,
	reason TEXT,
	corr_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(date, ticker)
);

CREATE INDEX IF NOT EXISTS idx_decisions_date ON decisions(date);
CREATE INDEX IF NOT EXISTS idx_decisions_ticker ON decisions(ticker);
CREATE INDEX IF NOT EXISTS idx_decisions_created_at ON decisions(created_at DESC);

-- Positions table: Open and closed trades
CREATE TABLE IF NOT EXISTS positions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ticker TEXT NOT NULL,
	entry_price REAL NOT NULL,
	current_stop REAL NOT NULL,
	initial_stop REAL NOT NULL,
	shares INTEGER NOT NULL,
	risk_dollars REAL NOT NULL,
	bucket TEXT,
	status TEXT NOT NULL DEFAULT 'OPEN',
	exit_price REAL,
	exit_date TEXT,
	outcome TEXT,
	pnl REAL,
	decision_id INTEGER,
	opened_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	closed_at DATETIME,
	FOREIGN KEY (decision_id) REFERENCES decisions(id)
);

CREATE INDEX IF NOT EXISTS idx_positions_ticker ON positions(ticker);
CREATE INDEX IF NOT EXISTS idx_positions_status ON positions(status);
CREATE INDEX IF NOT EXISTS idx_positions_bucket ON positions(bucket);
CREATE INDEX IF NOT EXISTS idx_positions_status_opened ON positions(status, opened_at DESC);

-- Impulse timers table: 2-minute brake enforcement
CREATE TABLE IF NOT EXISTS impulse_timers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ticker TEXT NOT NULL,
	started_at INTEGER NOT NULL,
	expires_at INTEGER NOT NULL,
	active INTEGER NOT NULL DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_impulse_timers_ticker ON impulse_timers(ticker, active);

-- Bucket cooldowns table: 24-hour sector lockout after losses
CREATE TABLE IF NOT EXISTS bucket_cooldowns (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	bucket TEXT NOT NULL,
	started_at INTEGER NOT NULL,
	expires_at INTEGER NOT NULL,
	active INTEGER NOT NULL DEFAULT 1,
	reason TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bucket_cooldowns_bucket ON bucket_cooldowns(bucket, active);

-- Trade sessions table: cohesive workflow and provenance tracking
CREATE TABLE IF NOT EXISTS trade_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_num INTEGER GENERATED ALWAYS AS (id) STORED UNIQUE,
    ticker TEXT NOT NULL,
    strategy TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'MANUAL' CHECK (source IN ('MANUAL', 'PRESET', 'CUSTOM')),
    candidate_id INTEGER,
    preset_id INTEGER,
    preset_name TEXT,
    scan_date TEXT,
    status TEXT NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'EVALUATING', 'COMPLETED', 'ABANDONED')),
    current_step TEXT NOT NULL DEFAULT 'CHECKLIST' CHECK (current_step IN ('CHECKLIST', 'SIZING', 'HEAT', 'ENTRY')),
    checklist_completed INTEGER NOT NULL DEFAULT 0 CHECK (checklist_completed IN (0,1)),
    checklist_banner TEXT,
    checklist_missing_count INTEGER DEFAULT 0,
    checklist_quality_score INTEGER DEFAULT 0,
    checklist_completed_at DATETIME,
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
    entry_completed INTEGER NOT NULL DEFAULT 0 CHECK (entry_completed IN (0,1)),
    entry_decision TEXT,
    entry_decision_id INTEGER,
    entry_gate1_pass INTEGER,
    entry_gate2_pass INTEGER,
    entry_gate3_pass INTEGER,
    entry_gate4_pass INTEGER,
    entry_gate5_pass INTEGER,
    entry_completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (candidate_id) REFERENCES candidates(id),
    FOREIGN KEY (preset_id) REFERENCES presets(id),
    FOREIGN KEY (entry_decision_id) REFERENCES decisions(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_session_num ON trade_sessions(session_num);
CREATE INDEX IF NOT EXISTS idx_sessions_ticker ON trade_sessions(ticker);
CREATE INDEX IF NOT EXISTS idx_sessions_strategy ON trade_sessions(strategy);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON trade_sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_created ON trade_sessions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_updated ON trade_sessions(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_active ON trade_sessions(status, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_candidate ON trade_sessions(candidate_id);

CREATE TRIGGER IF NOT EXISTS trg_sessions_updated_at
AFTER UPDATE ON trade_sessions
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE trade_sessions
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;
`

// defaultSettings contains the bootstrap configuration values
const defaultSettings = `
INSERT OR IGNORE INTO settings (key, value) VALUES
	('Equity_E', '10000'),
	('RiskPct_r', '0.0075'),
	('HeatCap_H_pct', '0.04'),
	('BucketHeatCap_pct', '0.015'),
	('StopMultiple_K', '2');
`
