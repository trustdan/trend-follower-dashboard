package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewOpenPositionCommand creates the open-position command
func NewOpenPositionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-position",
		Short: "Open position from GO decision",
		Long: `Create an open position from today's GO decision.

Examples:
  # Open position for AAPL
  tf-engine open-position --ticker AAPL

  # With JSON output
  tf-engine open-position --ticker AAPL --json`,
		RunE: runOpenPosition,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol (required)")
	cmd.Flags().Bool("json", false, "Output in JSON format")

	cmd.MarkFlagRequired("ticker")

	return cmd
}

func runOpenPosition(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	ticker, _ := cmd.Flags().GetString("ticker")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	log.WithField("ticker", ticker).Info("Opening position")

	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	position, err := db.OpenPosition(ticker)
	if err != nil {
		log.WithError(err).Error("Failed to open position")
		return fmt.Errorf("failed to open position: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"ticker": ticker,
		"shares": position.Shares,
		"risk":   position.RiskDollars,
	}).Info("Position opened")

	if jsonOutput {
		output, _ := json.MarshalIndent(position, "", "  ")
		fmt.Println(string(output))
	} else {
		fmt.Printf("✓ Position opened: %s\n", ticker)
		fmt.Printf("  Entry:  $%.2f\n", position.EntryPrice)
		fmt.Printf("  Stop:   $%.2f\n", position.CurrentStop)
		fmt.Printf("  Shares: %d\n", position.Shares)
		fmt.Printf("  Risk:   $%.2f\n", position.RiskDollars)
		if position.Bucket != "" {
			fmt.Printf("  Bucket: %s\n", position.Bucket)
		}
	}

	return nil
}

// NewUpdateStopCommand creates the update-stop command
func NewUpdateStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-stop",
		Short: "Update position stop price",
		Long: `Update the stop loss price for an open position.

Examples:
  # Update stop for AAPL to $179
  tf-engine update-stop --ticker AAPL --new-stop 179

  # With JSON output
  tf-engine update-stop --ticker AAPL --new-stop 179 --json`,
		RunE: runUpdateStop,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol (required)")
	cmd.Flags().Float64("new-stop", 0, "New stop price (required)")
	cmd.Flags().Bool("json", false, "Output in JSON format")

	cmd.MarkFlagRequired("ticker")
	cmd.MarkFlagRequired("new-stop")

	return cmd
}

func runUpdateStop(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	ticker, _ := cmd.Flags().GetString("ticker")
	newStop, _ := cmd.Flags().GetFloat64("new-stop")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	log.WithFields(map[string]interface{}{
		"ticker":   ticker,
		"new_stop": newStop,
	}).Info("Updating stop")

	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = db.UpdateStop(ticker, newStop)
	if err != nil {
		log.WithError(err).Error("Failed to update stop")
		return fmt.Errorf("failed to update stop: %w", err)
	}

	log.Info("Stop updated")

	if jsonOutput {
		result := map[string]interface{}{
			"ticker":   ticker,
			"new_stop": newStop,
			"message":  "Stop updated successfully",
		}
		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("✓ Stop updated for %s to $%.2f\n", ticker, newStop)
	}

	return nil
}

// NewClosePositionCommand creates the close-position command
func NewClosePositionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-position",
		Short: "Close an open position",
		Long: `Close a position and record outcome (WIN/LOSS/SCRATCH).

Examples:
  # Close AAPL position with exit at $185, outcome WIN
  tf-engine close-position --ticker AAPL --exit 185 --outcome WIN

  # Close with LOSS
  tf-engine close-position --ticker AAPL --exit 176 --outcome LOSS

  # With JSON output
  tf-engine close-position --ticker AAPL --exit 180 --outcome SCRATCH --json`,
		RunE: runClosePosition,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol (required)")
	cmd.Flags().Float64("exit", 0, "Exit price (required)")
	cmd.Flags().String("outcome", "", "Outcome: WIN/LOSS/SCRATCH (required)")
	cmd.Flags().Bool("json", false, "Output in JSON format")

	cmd.MarkFlagRequired("ticker")
	cmd.MarkFlagRequired("exit")
	cmd.MarkFlagRequired("outcome")

	return cmd
}

func runClosePosition(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	ticker, _ := cmd.Flags().GetString("ticker")
	exit, _ := cmd.Flags().GetFloat64("exit")
	outcome, _ := cmd.Flags().GetString("outcome")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// Validate outcome
	if outcome != "WIN" && outcome != "LOSS" && outcome != "SCRATCH" {
		return fmt.Errorf("outcome must be WIN, LOSS, or SCRATCH, got: %s", outcome)
	}

	log.WithFields(map[string]interface{}{
		"ticker":  ticker,
		"exit":    exit,
		"outcome": outcome,
	}).Info("Closing position")

	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get position before closing (for display)
	position, err := db.GetPositionByTicker(ticker)
	if err != nil {
		log.WithError(err).Error("Failed to get position")
		return fmt.Errorf("failed to get position: %w", err)
	}

	// Close position
	err = db.ClosePosition(ticker, exit, outcome)
	if err != nil {
		log.WithError(err).Error("Failed to close position")
		return fmt.Errorf("failed to close position: %w", err)
	}

	// Calculate P&L for display
	pnl := float64(position.Shares) * (exit - position.EntryPrice)

	log.WithFields(map[string]interface{}{
		"ticker":  ticker,
		"outcome": outcome,
		"pnl":     pnl,
	}).Info("Position closed")

	if jsonOutput {
		result := map[string]interface{}{
			"ticker":      ticker,
			"entry":       position.EntryPrice,
			"exit":        exit,
			"shares":      position.Shares,
			"pnl":         pnl,
			"outcome":     outcome,
			"bucket":      position.Bucket,
			"cooldown":    outcome == "LOSS" && position.Bucket != "",
		}
		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("✓ Position closed: %s\n", ticker)
		fmt.Printf("  Entry:   $%.2f\n", position.EntryPrice)
		fmt.Printf("  Exit:    $%.2f\n", exit)
		fmt.Printf("  Shares:  %d\n", position.Shares)
		fmt.Printf("  P&L:     $%.2f\n", pnl)
		fmt.Printf("  Outcome: %s\n", outcome)

		if outcome == "LOSS" && position.Bucket != "" {
			fmt.Printf("\n⚠️  Bucket %s entered cooldown (24 hours)\n", position.Bucket)
		}
	}

	return nil
}

// NewListPositionsCommand creates the list-positions command
func NewListPositionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-positions",
		Short: "List positions",
		Long: `List all positions, optionally filtered by status.

Examples:
  # List all positions
  tf-engine list-positions

  # List only open positions
  tf-engine list-positions --status OPEN

  # List only closed positions
  tf-engine list-positions --status CLOSED

  # With JSON output
  tf-engine list-positions --json`,
		RunE: runListPositions,
	}

	cmd.Flags().String("status", "", "Filter by status (OPEN or CLOSED)")
	cmd.Flags().Bool("json", false, "Output in JSON format")

	return cmd
}

func runListPositions(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	status, _ := cmd.Flags().GetString("status")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	log.WithField("status", status).Info("Listing positions")

	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	positions, err := db.GetAllPositions(status)
	if err != nil {
		log.WithError(err).Error("Failed to get positions")
		return fmt.Errorf("failed to get positions: %w", err)
	}

	log.WithField("count", len(positions)).Info("Positions retrieved")

	if jsonOutput {
		output, _ := json.MarshalIndent(positions, "", "  ")
		fmt.Println(string(output))
		return nil
	}

	if len(positions) == 0 {
		fmt.Println("No positions found")
		return nil
	}

	fmt.Printf("Positions: %d\n\n", len(positions))

	for _, p := range positions {
		fmt.Printf("• %s (%s)\n", p.Ticker, p.Status)
		fmt.Printf("  Entry: $%.2f | Stop: $%.2f | Shares: %d\n",
			p.EntryPrice, p.CurrentStop, p.Shares)
		fmt.Printf("  Risk: $%.2f", p.RiskDollars)
		if p.Bucket != "" {
			fmt.Printf(" | Bucket: %s", p.Bucket)
		}
		fmt.Println()

		if p.Status == "CLOSED" {
			fmt.Printf("  Exit: $%.2f | P&L: $%.2f | Outcome: %s\n",
				p.ExitPrice, p.PnL, p.Outcome)
		}
		fmt.Println()
	}

	return nil
}

// NewGetPositionCommand creates the get-position command
func NewGetPositionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-position",
		Short: "Get position details",
		Long: `Retrieve detailed information about a position.

Examples:
  # Get details for AAPL position
  tf-engine get-position --ticker AAPL

  # With JSON output
  tf-engine get-position --ticker AAPL --json`,
		RunE: runGetPosition,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol (required)")
	cmd.Flags().Bool("json", false, "Output in JSON format")

	cmd.MarkFlagRequired("ticker")

	return cmd
}

func runGetPosition(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	ticker, _ := cmd.Flags().GetString("ticker")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	log.WithField("ticker", ticker).Info("Getting position details")

	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	position, err := db.GetPositionByTicker(ticker)
	if err != nil {
		log.WithError(err).Error("Failed to get position")
		return fmt.Errorf("failed to get position: %w", err)
	}

	log.Info("Position retrieved")

	if jsonOutput {
		output, _ := json.MarshalIndent(position, "", "  ")
		fmt.Println(string(output))
		return nil
	}

	fmt.Printf("Position: %s\n\n", ticker)
	fmt.Printf("Status:        %s\n", position.Status)
	fmt.Printf("Entry:         $%.2f\n", position.EntryPrice)
	fmt.Printf("Current Stop:  $%.2f\n", position.CurrentStop)
	fmt.Printf("Initial Stop:  $%.2f\n", position.InitialStop)
	fmt.Printf("Shares:        %d\n", position.Shares)
	fmt.Printf("Risk:          $%.2f\n", position.RiskDollars)
	if position.Bucket != "" {
		fmt.Printf("Bucket:        %s\n", position.Bucket)
	}
	fmt.Printf("Opened:        %s\n", position.OpenedAt.Format("2006-01-02 15:04"))

	if position.Status == "CLOSED" {
		fmt.Printf("\nExit:          $%.2f\n", position.ExitPrice)
		fmt.Printf("Exit Date:     %s\n", position.ExitDate)
		fmt.Printf("P&L:           $%.2f\n", position.PnL)
		fmt.Printf("Outcome:       %s\n", position.Outcome)
	}

	return nil
}
