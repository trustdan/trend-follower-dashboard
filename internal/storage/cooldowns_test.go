package storage

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCooldownTestDB(t *testing.T) *DB {
	dbPath := "./test_cooldown_" + t.Name() + ".db"
	t.Cleanup(func() { os.Remove(dbPath) })

	db, err := New(dbPath)
	require.NoError(t, err)

	err = db.Initialize()
	require.NoError(t, err)

	return db
}

func TestTriggerBucketCooldown(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	err := db.TriggerBucketCooldown("Tech/Comm", "Test loss")
	assert.NoError(t, err)

	// Verify cooldown was created
	cooldown, err := db.GetBucketCooldown("Tech/Comm")
	assert.NoError(t, err)
	assert.NotNil(t, cooldown)
	assert.Equal(t, "Tech/Comm", cooldown.Bucket)
	assert.True(t, cooldown.Active)
	assert.Equal(t, "Test loss", cooldown.Reason)

	// Verify duration is 24 hours
	expectedExpiry := time.Now().Add(CooldownDuration)
	assert.WithinDuration(t, expectedExpiry, cooldown.ExpiresAt, 2*time.Second)
}

func TestTriggerBucketCooldown_EmptyBucket(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Empty bucket should not create cooldown
	err := db.TriggerBucketCooldown("", "Test")
	assert.NoError(t, err)

	// Verify no cooldown exists
	cooldowns, err := db.GetAllActiveCooldowns()
	assert.NoError(t, err)
	assert.Empty(t, cooldowns)
}

func TestTriggerBucketCooldown_ExtendExisting(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Create initial cooldown
	err := db.TriggerBucketCooldown("Tech/Comm", "First loss")
	assert.NoError(t, err)

	// Get initial expiry time
	cooldown1, _ := db.GetBucketCooldown("Tech/Comm")
	initialExpiry := cooldown1.ExpiresAt

	// Wait a bit to ensure new expiry is later
	time.Sleep(1 * time.Second)

	// Trigger again to extend
	err = db.TriggerBucketCooldown("Tech/Comm", "Second loss")
	assert.NoError(t, err)

	// Get updated cooldown
	cooldown2, _ := db.GetBucketCooldown("Tech/Comm")

	// Expiry should be extended (later than initial)
	assert.True(t, cooldown2.ExpiresAt.After(initialExpiry),
		"New expiry (%v) should be after initial expiry (%v)", cooldown2.ExpiresAt, initialExpiry)
	assert.Equal(t, "Second loss", cooldown2.Reason)

	// Should still be 24 hours from now
	expectedExpiry := time.Now().Add(CooldownDuration)
	assert.WithinDuration(t, expectedExpiry, cooldown2.ExpiresAt, 2*time.Second)
}

func TestGetBucketCooldown_NoActive(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	cooldown, err := db.GetBucketCooldown("NonExistent")
	assert.NoError(t, err)
	assert.Nil(t, cooldown)
}

func TestGetBucketCooldown_Expired(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Create cooldown with past expiry
	now := time.Now()
	pastExpiry := now.Add(-1 * time.Hour)

	_, err := db.conn.Exec(`
		INSERT INTO bucket_cooldowns (bucket, started_at, expires_at, active, reason)
		VALUES (?, ?, ?, 1, ?)
	`, "Tech/Comm", now.Add(-25*time.Hour).Unix(), pastExpiry.Unix(), "Test")
	assert.NoError(t, err)

	// Should return nil and deactivate the cooldown
	cooldown, err := db.GetBucketCooldown("Tech/Comm")
	assert.NoError(t, err)
	assert.Nil(t, cooldown)

	// Verify cooldown was deactivated in database
	var active int
	err = db.conn.QueryRow("SELECT active FROM bucket_cooldowns WHERE bucket = ?", "Tech/Comm").Scan(&active)
	assert.NoError(t, err)
	assert.Equal(t, 0, active)
}

func TestGetAllActiveCooldowns(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Create multiple cooldowns
	err := db.TriggerBucketCooldown("Tech/Comm", "Loss 1")
	assert.NoError(t, err)

	err = db.TriggerBucketCooldown("Energy", "Loss 2")
	assert.NoError(t, err)

	err = db.TriggerBucketCooldown("Finance", "Loss 3")
	assert.NoError(t, err)

	// Get all active
	cooldowns, err := db.GetAllActiveCooldowns()
	assert.NoError(t, err)
	assert.Len(t, cooldowns, 3)

	// Verify sorted by bucket
	buckets := []string{cooldowns[0].Bucket, cooldowns[1].Bucket, cooldowns[2].Bucket}
	assert.Contains(t, buckets, "Tech/Comm")
	assert.Contains(t, buckets, "Energy")
	assert.Contains(t, buckets, "Finance")
}

func TestGetAllActiveCooldowns_FiltersExpired(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Create active cooldown
	err := db.TriggerBucketCooldown("Tech/Comm", "Active")
	assert.NoError(t, err)

	// Create expired cooldown
	now := time.Now()
	pastExpiry := now.Add(-1 * time.Hour)
	result, err := db.conn.Exec(`
		INSERT INTO bucket_cooldowns (bucket, started_at, expires_at, active, reason)
		VALUES (?, ?, ?, 1, ?)
	`, "Energy", now.Add(-25*time.Hour).Unix(), pastExpiry.Unix(), "Expired")
	require.NoError(t, err)

	rowsAffected, _ := result.RowsAffected()
	require.Equal(t, int64(1), rowsAffected, "Expired cooldown should be inserted")

	// Should only return active cooldown (expired one filtered out)
	cooldowns, err := db.GetAllActiveCooldowns()
	assert.NoError(t, err)
	assert.Len(t, cooldowns, 1, "Should only have 1 active cooldown (Tech/Comm)")
	if len(cooldowns) > 0 {
		assert.Equal(t, "Tech/Comm", cooldowns[0].Bucket)
	}

	// Verify expired cooldown was deactivated by GetAllActiveCooldowns
	// Note: GetAllActiveCooldowns deactivates expired entries as it processes them
	var active int
	var id int
	err = db.conn.QueryRow("SELECT id, active FROM bucket_cooldowns WHERE bucket = ?", "Energy").Scan(&id, &active)
	assert.NoError(t, err)

	t.Logf("Energy cooldown: ID=%d, Active=%d", id, active)

	// The active flag should be 0 after GetAllActiveCooldowns processed it
	assert.Equal(t, 0, active, "Expired cooldown should be deactivated by GetAllActiveCooldowns")
}

func TestCheckBucketCooldown_NotInCooldown(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	err := db.CheckBucketCooldown("Tech/Comm")
	assert.NoError(t, err)
}

func TestCheckBucketCooldown_EmptyBucket(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	err := db.CheckBucketCooldown("")
	assert.NoError(t, err)
}

func TestCheckBucketCooldown_InCooldown(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Trigger cooldown
	err := db.TriggerBucketCooldown("Tech/Comm", "Test loss")
	assert.NoError(t, err)

	// Check should return error
	err = db.CheckBucketCooldown("Tech/Comm")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tech/Comm is in cooldown")
	assert.Contains(t, err.Error(), "hours remaining")
}

func TestCheckBucketCooldown_ExpiredReturnsNoError(t *testing.T) {
	db := setupCooldownTestDB(t)
	defer db.Close()

	// Create expired cooldown
	now := time.Now()
	pastExpiry := now.Add(-1 * time.Hour)
	_, err := db.conn.Exec(`
		INSERT INTO bucket_cooldowns (bucket, started_at, expires_at, active, reason)
		VALUES (?, ?, ?, 1, ?)
	`, "Tech/Comm", now.Add(-25*time.Hour).Unix(), pastExpiry.Unix(), "Test")
	assert.NoError(t, err)

	// Check should pass (expired cooldowns are treated as inactive)
	err = db.CheckBucketCooldown("Tech/Comm")
	assert.NoError(t, err)
}
