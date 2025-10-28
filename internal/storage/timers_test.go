package storage

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartImpulseTimer(t *testing.T) {
	// Create temp database
	dbPath := "./test_impulse_timer.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Initialize schema
	err = db.Initialize()
	require.NoError(t, err)

	// Start timer
	err = db.StartImpulseTimer("AAPL")
	assert.NoError(t, err)

	// Verify timer exists
	timer, err := db.GetActiveTimer("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, timer)
	assert.Equal(t, "AAPL", timer.Ticker)
	assert.True(t, timer.Active)

	// Verify expiration is 2 minutes from now
	expectedExpiry := time.Now().Add(ImpulseBrakeDuration)
	assert.WithinDuration(t, expectedExpiry, timer.ExpiresAt, 2*time.Second)
}

func TestStartImpulseTimer_ReplacesOldTimer(t *testing.T) {
	dbPath := "./test_impulse_replace.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Start first timer
	err = db.StartImpulseTimer("AAPL")
	require.NoError(t, err)

	timer1, err := db.GetActiveTimer("AAPL")
	require.NoError(t, err)
	require.NotNil(t, timer1)

	// Wait a moment to ensure different timestamps
	time.Sleep(1 * time.Second)

	// Start second timer (should replace first)
	err = db.StartImpulseTimer("AAPL")
	assert.NoError(t, err)

	timer2, err := db.GetActiveTimer("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, timer2)

	// Should be a different timer (different ID)
	assert.NotEqual(t, timer1.ID, timer2.ID)

	// New timer should have later timestamps
	assert.True(t, timer2.StartedAt.After(timer1.StartedAt))
}

func TestGetActiveTimer_NoTimer(t *testing.T) {
	dbPath := "./test_no_timer.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// No timer exists
	timer, err := db.GetActiveTimer("MSFT")
	assert.NoError(t, err)
	assert.Nil(t, timer)
}

func TestGetActiveTimer_MultipleTickers(t *testing.T) {
	dbPath := "./test_multiple_timers.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Start timers for different tickers
	err = db.StartImpulseTimer("AAPL")
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	err = db.StartImpulseTimer("MSFT")
	require.NoError(t, err)

	// Verify both timers exist independently
	timerAAPL, err := db.GetActiveTimer("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, timerAAPL)
	assert.Equal(t, "AAPL", timerAAPL.Ticker)

	timerMSFT, err := db.GetActiveTimer("MSFT")
	assert.NoError(t, err)
	assert.NotNil(t, timerMSFT)
	assert.Equal(t, "MSFT", timerMSFT.Ticker)

	// MSFT timer should have later start time
	assert.True(t, timerMSFT.StartedAt.After(timerAAPL.StartedAt))
}

func TestCheckImpulseBrake_NoTimer(t *testing.T) {
	dbPath := "./test_brake_no_timer.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Check brake without starting timer
	err = db.CheckImpulseBrake("AAPL")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no impulse timer active")
	assert.Contains(t, err.Error(), "evaluate checklist first")
}

func TestCheckImpulseBrake_TimerActive(t *testing.T) {
	dbPath := "./test_brake_active.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Start timer
	err = db.StartImpulseTimer("AAPL")
	require.NoError(t, err)

	// Immediately check - should fail (timer still active)
	err = db.CheckImpulseBrake("AAPL")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impulse brake active")
	assert.Contains(t, err.Error(), "more seconds")
}

func TestCheckImpulseBrake_TimerExpired(t *testing.T) {
	dbPath := "./test_brake_expired.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	require.NoError(t, err)
	defer db.Close()

	err = db.Initialize()
	require.NoError(t, err)

	// Start timer and manually set expiry to past
	err = db.StartImpulseTimer("AAPL")
	require.NoError(t, err)

	// Manually update expires_at to be in the past
	pastTime := time.Now().Add(-1 * time.Minute).Unix()
	_, err = db.conn.Exec("UPDATE impulse_timers SET expires_at = ? WHERE ticker = ?", pastTime, "AAPL")
	require.NoError(t, err)

	// Check brake - should pass (timer expired)
	err = db.CheckImpulseBrake("AAPL")
	assert.NoError(t, err)
}

func TestImpulseBrakeDuration(t *testing.T) {
	// Verify the constant is set correctly
	assert.Equal(t, 2*time.Minute, ImpulseBrakeDuration)
	assert.Equal(t, 120*time.Second, ImpulseBrakeDuration)
}
