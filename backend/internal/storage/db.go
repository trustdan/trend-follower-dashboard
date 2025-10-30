package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection
type DB struct {
	conn  *sql.DB
	path  string
	cache *Cache
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Performance optimizations for SQLite
	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	if _, err := conn.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}

	if _, err := conn.Exec("PRAGMA cache_size = -64000"); err != nil { // 64MB cache
		conn.Close()
		return nil, fmt.Errorf("failed to set cache size: %w", err)
	}

	if _, err := conn.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set temp store: %w", err)
	}

	return &DB{
		conn:  conn,
		path:  dbPath,
		cache: NewCache(),
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Initialize creates all tables and bootstraps default settings
func (db *DB) Initialize() error {
	// Create tables
	if _, err := db.conn.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Bootstrap default settings
	if _, err := db.conn.Exec(defaultSettings); err != nil {
		return fmt.Errorf("failed to bootstrap settings: %w", err)
	}

	return nil
}

// GetSetting retrieves a configuration value by key
func (db *DB) GetSetting(key string) (string, error) {
	var value string
	query := `SELECT value FROM settings WHERE key = ?`
	err := db.conn.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("setting not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get setting: %w", err)
	}
	return value, nil
}

// SetSetting updates or inserts a configuration value
func (db *DB) SetSetting(key, value string) error {
	query := `
		INSERT INTO settings (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.conn.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}

	// Invalidate cache
	db.cache.Delete("all_settings")
	db.cache.Delete("setting:" + key)

	return nil
}

// GetAllSettings retrieves all configuration key-value pairs
func (db *DB) GetAllSettings() (map[string]string, error) {
	// Try cache first (5 minute TTL)
	if cached, ok := db.cache.Get("all_settings"); ok {
		return cached.(map[string]string), nil
	}

	query := `SELECT key, value FROM settings ORDER BY key`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	// Cache for 5 minutes
	db.cache.Set("all_settings", settings, 5*time.Minute)

	return settings, nil
}

// GetOrCreatePreset gets an existing preset by name or creates it if it doesn't exist
// Returns the preset ID
func (db *DB) GetOrCreatePreset(name, queryString string) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("preset name cannot be empty")
	}

	// Try to get existing preset
	var presetID int
	query := `SELECT id FROM presets WHERE name = ?`
	err := db.conn.QueryRow(query, name).Scan(&presetID)
	if err == nil {
		return presetID, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query preset: %w", err)
	}

	// Create new preset
	insertQuery := `INSERT INTO presets (name, query_string) VALUES (?, ?)`
	result, err := db.conn.Exec(insertQuery, name, queryString)
	if err != nil {
		return 0, fmt.Errorf("failed to create preset: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get preset ID: %w", err)
	}

	return int(id), nil
}

// ImportCandidates imports a list of candidates for a specific date
// This replaces any existing candidates for the same date and preset
func (db *DB) ImportCandidates(date string, tickers []string, presetID *int, sector, bucket string) error {
	if len(tickers) == 0 {
		return fmt.Errorf("at least one ticker required")
	}

	// Start transaction
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing candidates for this date and preset
	deleteQuery := `DELETE FROM candidates WHERE date = ? AND preset_id IS ?`
	_, err = tx.Exec(deleteQuery, date, presetID)
	if err != nil {
		return fmt.Errorf("failed to delete existing candidates: %w", err)
	}

	// Insert new candidates
	insertQuery := `
		INSERT INTO candidates (date, ticker, preset_id, sector, bucket)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(date, ticker, preset_id) DO UPDATE SET
			sector = excluded.sector,
			bucket = excluded.bucket
	`

	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, ticker := range tickers {
		_, err := stmt.Exec(date, ticker, presetID, sector, bucket)
		if err != nil {
			return fmt.Errorf("failed to insert candidate %s: %w", ticker, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetCandidatesForDate retrieves all candidates for a specific date
func (db *DB) GetCandidatesForDate(date string) ([]map[string]interface{}, error) {
	query := `
		SELECT c.id, c.date, c.ticker, c.preset_id, p.name as preset_name, c.sector, c.bucket
		FROM candidates c
		LEFT JOIN presets p ON c.preset_id = p.id
		WHERE c.date = ?
		ORDER BY c.ticker
	`

	rows, err := db.conn.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query candidates: %w", err)
	}
	defer rows.Close()

	candidates := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var date, ticker, sector, bucket string
		var presetID sql.NullInt64
		var presetName sql.NullString

		err := rows.Scan(&id, &date, &ticker, &presetID, &presetName, &sector, &bucket)
		if err != nil {
			return nil, fmt.Errorf("failed to scan candidate: %w", err)
		}

		candidate := map[string]interface{}{
			"id":     id,
			"date":   date,
			"ticker": ticker,
			"sector": sector,
			"bucket": bucket,
		}

		if presetID.Valid {
			candidate["preset_id"] = int(presetID.Int64)
		}
		if presetName.Valid {
			candidate["preset"] = presetName.String
		}

		candidates = append(candidates, candidate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating candidates: %w", err)
	}

	return candidates, nil
}

// IsTickerInCandidates checks if a ticker is in the candidates list for a specific date
func (db *DB) IsTickerInCandidates(date, ticker string) (bool, error) {
	query := `SELECT COUNT(*) FROM candidates WHERE date = ? AND ticker = ?`
	var count int
	err := db.conn.QueryRow(query, date, ticker).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ticker: %w", err)
	}
	return count > 0, nil
}

// GetCandidatesCount returns the count of candidates for a specific date
func (db *DB) GetCandidatesCount(date string) (int, error) {
	query := `SELECT COUNT(*) FROM candidates WHERE date = ?`
	var count int
	err := db.conn.QueryRow(query, date).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count candidates: %w", err)
	}
	return count, nil
}
