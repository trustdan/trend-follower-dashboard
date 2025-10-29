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
