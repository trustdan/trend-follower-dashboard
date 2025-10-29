package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewChecklistCommand creates the checklist command
func NewChecklistCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checklist",
		Short: "Evaluate entry checklist for a trade setup",
		Long: `Evaluate the 6-item entry checklist and determine banner color.

The system enforces a 6-item checklist before any trade:
  1. FromPreset    - Ticker came from today's FINVIZ preset
  2. TrendPass     - Trend alignment confirmed
  3. LiquidityPass - Adequate volume and spread
  4. TVConfirm     - TradingView setup confirmation
  5. EarningsOK    - No earnings in next 7 days
  6. JournalOK     - Trade thesis documented

Banner Logic:
  - 0 missing → GREEN  (proceed to impulse timer)
  - 1 missing → YELLOW (caution, do not proceed)
  - 2+ missing → RED   (no-go, do not proceed)

Only GREEN banner allows you to proceed with the trade.

Examples:
  # Perfect setup
  tf-engine checklist --ticker AAPL --from-preset --trend-pass --liquidity-pass --tv-confirm --earnings-ok --journal-ok

  # One item missing (YELLOW)
  tf-engine checklist --ticker AAPL --from-preset --trend-pass --liquidity-pass --tv-confirm --earnings-ok

  # Two items missing (RED)
  tf-engine checklist --ticker AAPL --from-preset --trend-pass --liquidity-pass`,
		RunE: runChecklist,
	}

	// Required flags
	cmd.Flags().String("ticker", "", "Stock ticker symbol (required)")
	cmd.MarkFlagRequired("ticker")

	// Checklist item flags (all boolean)
	cmd.Flags().Bool("from-preset", false, "Ticker is from today's FINVIZ preset")
	cmd.Flags().Bool("trend-pass", false, "Trend alignment confirmed")
	cmd.Flags().Bool("liquidity-pass", false, "Adequate volume and spread")
	cmd.Flags().Bool("tv-confirm", false, "TradingView setup confirmation")
	cmd.Flags().Bool("earnings-ok", false, "No earnings in next 7 days")
	cmd.Flags().Bool("journal-ok", false, "Trade thesis documented in journal")

	return cmd
}

func runChecklist(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	format := GetOutputFormat(cmd)
	log := logx.WithCorrelationID(corrID)

	log.Info("Starting checklist evaluation")

	// Get flags
	ticker, _ := cmd.Flags().GetString("ticker")
	fromPreset, _ := cmd.Flags().GetBool("from-preset")
	trendPass, _ := cmd.Flags().GetBool("trend-pass")
	liquidityPass, _ := cmd.Flags().GetBool("liquidity-pass")
	tvConfirm, _ := cmd.Flags().GetBool("tv-confirm")
	earningsOK, _ := cmd.Flags().GetBool("earnings-ok")
	journalOK, _ := cmd.Flags().GetBool("journal-ok")

	// Build request
	req := domain.ChecklistRequest{
		Ticker:        ticker,
		FromPreset:    fromPreset,
		TrendPass:     trendPass,
		LiquidityPass: liquidityPass,
		TVConfirm:     tvConfirm,
		EarningsOK:    earningsOK,
		JournalOK:     journalOK,
	}

	log.WithField("request", req).Info("Evaluating checklist")

	// Evaluate checklist
	result, err := domain.EvaluateChecklist(req)
	if err != nil {
		log.WithError(err).Error("Checklist evaluation failed")
		return fmt.Errorf("checklist evaluation failed: %w", err)
	}

	log.WithField("result", result).Info("Checklist evaluation completed")

	// Start impulse timer if banner is GREEN
	if result.Banner == "GREEN" {
		// Open database
		db, err := storage.New(dbPath)
		if err != nil {
			log.WithError(err).Error("Failed to open database")
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		err = db.StartImpulseTimer(ticker)
		if err != nil {
			log.WithError(err).Error("Failed to start impulse timer")
			return fmt.Errorf("failed to start impulse timer: %w", err)
		}

		log.WithField("ticker", ticker).Info("Impulse timer started (2 minutes)")

		// Output human-readable message (only if format is human)
		PrintHuman(format, "\n⏱️  Impulse brake timer started")
		PrintHuman(format, "   Wait 2 minutes before saving decision")
	}

	// Output JSON result (always)
	if err := PrintJSON(result); err != nil {
		log.WithError(err).Error("Failed to marshal result to JSON")
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	return nil
}
