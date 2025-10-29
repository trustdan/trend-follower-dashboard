package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewCheckTimerCommand creates the check-timer command
func NewCheckTimerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-timer",
		Short: "Check impulse brake timer status",
		Long: `Check the status of the impulse brake timer for a ticker.

The impulse brake timer is a 2-minute mandatory delay that starts
when checklist evaluation returns GREEN. This prevents impulsive
trading decisions by enforcing a cooling-off period.

Examples:
  # Check timer status for AAPL
  tf-engine check-timer --ticker AAPL

  # Check with JSON output
  tf-engine check-timer --ticker AAPL --format json`,
		RunE: runCheckTimer,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol (required)")
	cmd.MarkFlagRequired("ticker")

	return cmd
}

func runCheckTimer(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	format := GetOutputFormat(cmd)
	log := logx.WithCorrelationID(corrID)

	ticker, _ := cmd.Flags().GetString("ticker")

	log.WithField("ticker", ticker).Info("Checking impulse timer")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get active timer
	timer, err := db.GetActiveTimer(ticker)
	if err != nil {
		log.WithError(err).Error("Failed to get timer")
		return fmt.Errorf("failed to get timer: %w", err)
	}

	// Build result
	now := time.Now()
	result := make(map[string]interface{})
	result["ticker"] = ticker
	result["timer_active"] = timer != nil

	if timer == nil {
		result["elapsed_seconds"] = 0
		result["remaining_seconds"] = 0
		result["brake_cleared"] = false

		PrintHumanf(format, "No active timer for %s\n", ticker)
		PrintHuman(format, "You must evaluate checklist with GREEN banner first")

		PrintJSON(result)

		log.WithField("ticker", ticker).Info("No active timer found")
		return nil
	}

	// Timer exists - check if expired
	elapsed := now.Sub(timer.StartedAt)
	elapsedSeconds := int(elapsed.Seconds())

	result["started_at"] = timer.StartedAt.Format(time.RFC3339)
	result["elapsed_seconds"] = elapsedSeconds

	if now.Before(timer.ExpiresAt) {
		// Timer still active
		remaining := timer.ExpiresAt.Sub(now)
		remainingSeconds := int(remaining.Seconds())

		result["remaining_seconds"] = remainingSeconds
		result["brake_cleared"] = false

		PrintHumanf(format, "⏱️  Timer active for %s\n", ticker)
		PrintHumanf(format, "   Started: %s\n", timer.StartedAt.Format("15:04:05"))
		PrintHumanf(format, "   Expires: %s\n", timer.ExpiresAt.Format("15:04:05"))
		PrintHumanf(format, "   Remaining: %d seconds\n", remainingSeconds)

		log.WithField("ticker", ticker).WithField("remaining", remainingSeconds).Info("Timer active")
	} else {
		// Timer expired
		result["remaining_seconds"] = 0
		result["brake_cleared"] = true

		PrintHumanf(format, "✓ Timer expired for %s\n", ticker)
		PrintHuman(format, "  You may proceed with save-decision")

		log.WithField("ticker", ticker).Info("Timer expired")
	}

	// JSON output (always)
	PrintJSON(result)

	return nil
}
