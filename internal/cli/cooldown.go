package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewCheckCooldownCommand creates the check-cooldown command
func NewCheckCooldownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-cooldown",
		Short: "Check bucket cooldown status",
		Long: `Check if a bucket is in cooldown and show remaining time.

Examples:
  # Check if Tech/Comm bucket is in cooldown (human-readable)
  tf-engine check-cooldown --bucket "Tech/Comm"

  # Check with JSON output
  tf-engine check-cooldown --bucket "Tech/Comm" --format json`,
		RunE: runCheckCooldown,
	}

	cmd.Flags().String("bucket", "", "Bucket name (required)")
	cmd.MarkFlagRequired("bucket")

	return cmd
}

func runCheckCooldown(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	format := GetOutputFormat(cmd)
	log := logx.WithCorrelationID(corrID)

	bucket, _ := cmd.Flags().GetString("bucket")

	log.WithField("bucket", bucket).Info("Checking bucket cooldown")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get cooldown
	cooldown, err := db.GetBucketCooldown(bucket)
	if err != nil {
		log.WithError(err).Error("Failed to get cooldown")
		return fmt.Errorf("failed to get cooldown: %w", err)
	}

	if cooldown == nil {
		// Not in cooldown
		result := map[string]interface{}{
			"bucket":      bucket,
			"in_cooldown": false,
		}

		PrintHumanf(format, "✓ Bucket %s is NOT in cooldown\n", bucket)
		PrintJSON(result)

		log.Info("Bucket not in cooldown")
		return nil
	}

	// Cooldown is active
	remaining := cooldown.ExpiresAt.Sub(time.Now())
	hoursRemaining := remaining.Hours()

	result := map[string]interface{}{
		"bucket":           bucket,
		"in_cooldown":      true,
		"started_at":       cooldown.StartedAt.Format(time.RFC3339),
		"expires_at":       cooldown.ExpiresAt.Format(time.RFC3339),
		"remaining_hours":  hoursRemaining,
		"reason":           cooldown.Reason,
	}

	PrintHumanf(format, "⏱️  Bucket %s is in cooldown\n", bucket)
	PrintHumanf(format, "   Started: %s\n", cooldown.StartedAt.Format("2006-01-02 15:04"))
	PrintHumanf(format, "   Expires: %s\n", cooldown.ExpiresAt.Format("2006-01-02 15:04"))
	PrintHumanf(format, "   Remaining: %.1f hours\n", hoursRemaining)
	if cooldown.Reason != "" {
		PrintHumanf(format, "   Reason: %s\n", cooldown.Reason)
	}

	PrintJSON(result)

	log.WithField("hours_remaining", hoursRemaining).Info("Bucket in cooldown")
	return nil
}

// NewListCooldownsCommand creates the list-cooldowns command
func NewListCooldownsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-cooldowns",
		Short: "List all active bucket cooldowns",
		Long: `Show all buckets currently in cooldown with remaining time.

Examples:
  # List all active cooldowns (human-readable)
  tf-engine list-cooldowns

  # List with JSON output
  tf-engine list-cooldowns --format json`,
		RunE: runListCooldowns,
	}

	return cmd
}

func runListCooldowns(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	format := GetOutputFormat(cmd)
	log := logx.WithCorrelationID(corrID)

	log.Info("Listing all active cooldowns")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get all active cooldowns
	cooldowns, err := db.GetAllActiveCooldowns()
	if err != nil {
		log.WithError(err).Error("Failed to get cooldowns")
		return fmt.Errorf("failed to get cooldowns: %w", err)
	}

	if len(cooldowns) == 0 {
		result := map[string]interface{}{
			"cooldowns": []interface{}{},
			"count":     0,
		}

		PrintHuman(format, "✓ No active cooldowns")
		PrintJSON(result)

		log.Info("No active cooldowns")
		return nil
	}

	// Build JSON response
	type cooldownInfo struct {
		Bucket         string  `json:"bucket"`
		StartedAt      string  `json:"started_at"`
		ExpiresAt      string  `json:"expires_at"`
		RemainingHours float64 `json:"remaining_hours"`
		Reason         string  `json:"reason,omitempty"`
	}

	list := make([]cooldownInfo, len(cooldowns))
	for i, c := range cooldowns {
		remaining := c.ExpiresAt.Sub(time.Now())
		list[i] = cooldownInfo{
			Bucket:         c.Bucket,
			StartedAt:      c.StartedAt.Format(time.RFC3339),
			ExpiresAt:      c.ExpiresAt.Format(time.RFC3339),
			RemainingHours: remaining.Hours(),
			Reason:         c.Reason,
		}
	}

	result := map[string]interface{}{
		"cooldowns": list,
		"count":     len(cooldowns),
	}

	// Human-readable output
	PrintHumanf(format, "⏱️  Active cooldowns: %d\n", len(cooldowns))
	for _, c := range cooldowns {
		remaining := c.ExpiresAt.Sub(time.Now())
		PrintHumanf(format, "Bucket: %s\n", c.Bucket)
		PrintHumanf(format, "  Expires: %s (%.1f hours remaining)\n",
			c.ExpiresAt.Format("2006-01-02 15:04"), remaining.Hours())
		if c.Reason != "" {
			PrintHumanf(format, "  Reason: %s\n", c.Reason)
		}
		PrintHuman(format, "")
	}

	// JSON output
	PrintJSON(result)

	log.WithField("count", len(cooldowns)).Info("Listed active cooldowns")
	return nil
}

// NewTriggerCooldownCommand creates the trigger-cooldown command (for testing)
func NewTriggerCooldownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trigger-cooldown",
		Short: "Manually trigger a bucket cooldown (for testing)",
		Long: `Manually create or extend a bucket cooldown period.

This is primarily for testing and manual intervention. In normal operation,
cooldowns are triggered automatically when positions are closed at a loss.

Examples:
  # Trigger cooldown for Tech/Comm bucket
  tf-engine trigger-cooldown --bucket "Tech/Comm" --reason "Manual test"

  # Trigger with JSON output
  tf-engine trigger-cooldown --bucket "Tech/Comm" --format json`,
		RunE: runTriggerCooldown,
	}

	cmd.Flags().String("bucket", "", "Bucket name (required)")
	cmd.Flags().String("reason", "Manual trigger", "Reason for cooldown")
	cmd.MarkFlagRequired("bucket")

	return cmd
}

func runTriggerCooldown(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	format := GetOutputFormat(cmd)
	log := logx.WithCorrelationID(corrID)

	bucket, _ := cmd.Flags().GetString("bucket")
	reason, _ := cmd.Flags().GetString("reason")

	log.WithField("bucket", bucket).Info("Triggering bucket cooldown")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Trigger cooldown
	err = db.TriggerBucketCooldown(bucket, reason)
	if err != nil {
		log.WithError(err).Error("Failed to trigger cooldown")
		return fmt.Errorf("failed to trigger cooldown: %w", err)
	}

	// Get the cooldown to show details
	cooldown, _ := db.GetBucketCooldown(bucket)

	if cooldown != nil {
		result := map[string]interface{}{
			"bucket":     bucket,
			"triggered":  true,
			"expires_at": cooldown.ExpiresAt.Format(time.RFC3339),
			"reason":     reason,
		}

		PrintHumanf(format, "✓ Cooldown triggered for bucket: %s\n", bucket)
		PrintHumanf(format, "  Expires: %s (24 hours)\n", cooldown.ExpiresAt.Format("2006-01-02 15:04"))
		PrintHumanf(format, "  Reason: %s\n", reason)

		PrintJSON(result)
	}

	log.Info("Bucket cooldown triggered successfully")
	return nil
}
