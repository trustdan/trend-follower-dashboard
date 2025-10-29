package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// BucketCooldown represents an active cooldown period for a sector bucket
type BucketCooldown struct {
	ID        int       `json:"id"`
	Bucket    string    `json:"bucket"`
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Active    bool      `json:"active"`
	Reason    string    `json:"reason"`
}

// CooldownDuration is the mandatory 24-hour wait period after a loss
const CooldownDuration = 24 * time.Hour

// TriggerBucketCooldown creates or extends cooldown for a bucket
// Called when a loss is recorded in a bucket
func (db *DB) TriggerBucketCooldown(bucket, reason string) error {
	if bucket == "" {
		return nil // No bucket, no cooldown
	}

	now := time.Now()
	expiresAt := now.Add(CooldownDuration)

	// Check if cooldown already exists
	existing, err := db.GetBucketCooldown(bucket)
	if err != nil {
		return fmt.Errorf("failed to check existing cooldown: %w", err)
	}

	if existing != nil && existing.Active {
		// Extend existing cooldown (reset to 24 hours from now)
		query := `
			UPDATE bucket_cooldowns
			SET expires_at = ?, reason = ?
			WHERE bucket = ? AND active = 1
		`
		_, err := db.conn.Exec(query, expiresAt.Unix(), reason, bucket)
		if err != nil {
			return fmt.Errorf("failed to extend cooldown: %w", err)
		}
	} else {
		// Create new cooldown
		query := `
			INSERT INTO bucket_cooldowns (bucket, started_at, expires_at, active, reason)
			VALUES (?, ?, ?, 1, ?)
		`
		_, err := db.conn.Exec(query, bucket, now.Unix(), expiresAt.Unix(), reason)
		if err != nil {
			return fmt.Errorf("failed to create cooldown: %w", err)
		}
	}

	return nil
}

// GetBucketCooldown retrieves the active cooldown for a bucket
// Returns nil if no active cooldown exists or if cooldown has expired
func (db *DB) GetBucketCooldown(bucket string) (*BucketCooldown, error) {
	query := `
		SELECT id, bucket, started_at, expires_at, active, reason
		FROM bucket_cooldowns
		WHERE bucket = ? AND active = 1
		ORDER BY started_at DESC
		LIMIT 1
	`

	var cooldown BucketCooldown
	var startedUnix, expiresUnix int64

	err := db.conn.QueryRow(query, bucket).Scan(
		&cooldown.ID,
		&cooldown.Bucket,
		&startedUnix,
		&expiresUnix,
		&cooldown.Active,
		&cooldown.Reason,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No active cooldown
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cooldown: %w", err)
	}

	cooldown.StartedAt = time.Unix(startedUnix, 0)
	cooldown.ExpiresAt = time.Unix(expiresUnix, 0)

	// Check if expired
	if time.Now().After(cooldown.ExpiresAt) {
		// Deactivate expired cooldown
		_, _ = db.conn.Exec(`UPDATE bucket_cooldowns SET active = 0 WHERE id = ?`, cooldown.ID)
		return nil, nil // Expired, return nil
	}

	return &cooldown, nil
}

// GetAllActiveCooldowns retrieves all currently active cooldowns
// Automatically deactivates any that have expired
func (db *DB) GetAllActiveCooldowns() ([]BucketCooldown, error) {
	query := `
		SELECT id, bucket, started_at, expires_at, active, reason
		FROM bucket_cooldowns
		WHERE active = 1
		ORDER BY bucket
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query cooldowns: %w", err)
	}
	defer rows.Close()

	cooldowns := []BucketCooldown{}
	expiredIDs := []int{}
	now := time.Now()

	for rows.Next() {
		var c BucketCooldown
		var startedUnix, expiresUnix int64

		err := rows.Scan(&c.ID, &c.Bucket, &startedUnix, &expiresUnix, &c.Active, &c.Reason)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cooldown: %w", err)
		}

		c.StartedAt = time.Unix(startedUnix, 0)
		c.ExpiresAt = time.Unix(expiresUnix, 0)

		// Only include if not expired
		if now.Before(c.ExpiresAt) {
			cooldowns = append(cooldowns, c)
		} else {
			// Mark for deactivation
			expiredIDs = append(expiredIDs, c.ID)
		}
	}

	// Deactivate expired cooldowns after closing rows
	for _, id := range expiredIDs {
		_, _ = db.conn.Exec(`UPDATE bucket_cooldowns SET active = 0 WHERE id = ?`, id)
	}

	return cooldowns, nil
}

// CheckBucketCooldown validates cooldown status before allowing save
// Returns error if bucket is in active cooldown
func (db *DB) CheckBucketCooldown(bucket string) error {
	if bucket == "" {
		return nil // No bucket specified, skip cooldown check
	}

	cooldown, err := db.GetBucketCooldown(bucket)
	if err != nil {
		return fmt.Errorf("failed to check cooldown: %w", err)
	}

	if cooldown == nil {
		return nil // No active cooldown
	}

	remaining := cooldown.ExpiresAt.Sub(time.Now())
	return fmt.Errorf("bucket %s is in cooldown (%.1f hours remaining)",
		bucket, remaining.Hours())
}
