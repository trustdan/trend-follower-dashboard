package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	dbPath := "test_new.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err, "New() should not return error")
	require.NotNil(t, db, "DB should not be nil")
	defer db.Close()

	assert.Equal(t, dbPath, db.path, "DB path should match")
	assert.NotNil(t, db.conn, "DB connection should not be nil")
}

func TestInitialize(t *testing.T) {
	dbPath := "test_initialize.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err, "Initialize() should not return error")

	// Verify all tables exist
	tables := []string{"settings", "presets", "candidates", "checklist_evaluations", "decisions", "positions"}
	for _, table := range tables {
		var name string
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		err := db.conn.QueryRow(query, table).Scan(&name)
		require.NoError(t, err, "Table %s should exist", table)
		assert.Equal(t, table, name)
	}

	// Verify default settings
	settings, err := db.GetAllSettings()
	require.NoError(t, err)
	assert.Equal(t, "10000", settings["Equity_E"])
	assert.Equal(t, "0.0075", settings["RiskPct_r"])
	assert.Equal(t, "0.04", settings["HeatCap_H_pct"])
	assert.Equal(t, "0.015", settings["BucketHeatCap_pct"])
	assert.Equal(t, "2", settings["StopMultiple_K"])
}

func TestInitializeIdempotent(t *testing.T) {
	dbPath := "test_idempotent.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	// First initialization
	err = db.Initialize()
	require.NoError(t, err)

	// Modify a setting
	err = db.SetSetting("Equity_E", "20000")
	require.NoError(t, err)

	// Second initialization should preserve the change
	err = db.Initialize()
	require.NoError(t, err)

	value, err := db.GetSetting("Equity_E")
	require.NoError(t, err)
	assert.Equal(t, "20000", value, "Modified setting should be preserved")
}

func TestGetSetting(t *testing.T) {
	dbPath := "test_get_setting.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Test existing setting
	value, err := db.GetSetting("Equity_E")
	require.NoError(t, err)
	assert.Equal(t, "10000", value)

	// Test non-existent setting
	_, err = db.GetSetting("NonExistent")
	assert.Error(t, err, "Should return error for non-existent setting")
	assert.Contains(t, err.Error(), "setting not found")
}

func TestSetSetting(t *testing.T) {
	dbPath := "test_set_setting.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Test updating existing setting
	err = db.SetSetting("Equity_E", "15000")
	require.NoError(t, err)

	value, err := db.GetSetting("Equity_E")
	require.NoError(t, err)
	assert.Equal(t, "15000", value)

	// Test inserting new setting
	err = db.SetSetting("NewSetting", "test_value")
	require.NoError(t, err)

	value, err = db.GetSetting("NewSetting")
	require.NoError(t, err)
	assert.Equal(t, "test_value", value)
}

func TestGetAllSettings(t *testing.T) {
	dbPath := "test_get_all_settings.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	settings, err := db.GetAllSettings()
	require.NoError(t, err)
	assert.Len(t, settings, 5, "Should have 5 default settings")

	// Verify all default settings are present
	expectedKeys := []string{"Equity_E", "RiskPct_r", "HeatCap_H_pct", "BucketHeatCap_pct", "StopMultiple_K"}
	for _, key := range expectedKeys {
		_, ok := settings[key]
		assert.True(t, ok, "Setting %s should exist", key)
	}
}

func TestClose(t *testing.T) {
	dbPath := "test_close.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err, "Close() should not return error")

	// Verify connection is closed by attempting an operation
	_, err = db.GetAllSettings()
	assert.Error(t, err, "Operations should fail after Close()")
}
