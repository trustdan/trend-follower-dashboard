package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewCheckHeatCommand creates the check-heat command
func NewCheckHeatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-heat",
		Short: "Check portfolio and bucket heat",
		Long: `Check current portfolio heat and bucket heat status.

Examples:
  # Check current heat with no new trade
  tf-engine check-heat

  # Check heat with a proposed new trade
  tf-engine check-heat --add-risk 75 --add-bucket "Tech/Comm"

  # With JSON output
  tf-engine check-heat --add-risk 75 --format json`,
		RunE: runCheckHeat,
	}

	cmd.Flags().Float64("add-risk", 0, "Risk dollars for proposed new trade")
	cmd.Flags().String("add-bucket", "", "Bucket for proposed new trade")

	return cmd
}

func runCheckHeat(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	format := GetOutputFormat(cmd)
	log := logx.WithCorrelationID(corrID)

	addRisk, _ := cmd.Flags().GetFloat64("add-risk")
	addBucket, _ := cmd.Flags().GetString("add-bucket")

	log.WithFields(map[string]interface{}{
		"add_risk":   addRisk,
		"add_bucket": addBucket,
	}).Info("Checking heat")

	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get settings
	settings, err := db.GetAllSettings()
	if err != nil {
		log.WithError(err).Error("Failed to get settings")
		return fmt.Errorf("failed to get settings: %w", err)
	}

	equity, err := strconv.ParseFloat(settings["Equity_E"], 64)
	if err != nil {
		return fmt.Errorf("invalid Equity_E setting: %w", err)
	}

	heatCapPct, err := strconv.ParseFloat(settings["HeatCap_H_pct"], 64)
	if err != nil {
		return fmt.Errorf("invalid HeatCap_H_pct setting: %w", err)
	}

	bucketHeatCapPct, err := strconv.ParseFloat(settings["BucketHeatCap_pct"], 64)
	if err != nil {
		return fmt.Errorf("invalid BucketHeatCap_pct setting: %w", err)
	}

	// Get open positions from database
	positions, err := db.GetOpenPositions()
	if err != nil {
		log.WithError(err).Error("Failed to get open positions")
		return fmt.Errorf("failed to get open positions: %w", err)
	}

	// Convert storage positions to domain positions
	openPositions := make([]domain.Position, len(positions))
	for i, p := range positions {
		openPositions[i] = domain.Position{
			Ticker:      p.Ticker,
			Bucket:      p.Bucket,
			RiskDollars: p.RiskDollars,
			UnitsOpen:   p.Shares,
			Status:      "Open", // All from GetOpenPositions are open
		}
	}

	// Calculate heat
	heatReq := domain.HeatRequest{
		Equity:           equity,
		HeatCapPct:       heatCapPct,
		BucketHeatCapPct: bucketHeatCapPct,
		AddRiskDollars:   addRisk,
		AddBucket:        addBucket,
		OpenPositions:    openPositions,
	}

	result, err := domain.CalculateHeat(heatReq)
	if err != nil {
		log.WithError(err).Error("Failed to calculate heat")
		return fmt.Errorf("failed to calculate heat: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"current_portfolio_heat": result.CurrentPortfolioHeat,
		"new_portfolio_heat":     result.NewPortfolioHeat,
		"allowed":                result.Allowed,
	}).Info("Heat calculated")

	// Human-readable output
	PrintHuman(format, "Portfolio Heat Status")
	PrintHuman(format, "=====================")
	PrintHumanf(format, "Current Heat:  $%.2f\n", result.CurrentPortfolioHeat)
	if addRisk > 0 {
		PrintHumanf(format, "New Heat:      $%.2f (+$%.2f)\n", result.NewPortfolioHeat, addRisk)
	} else {
		PrintHumanf(format, "New Heat:      $%.2f\n", result.NewPortfolioHeat)
	}
	PrintHumanf(format, "Heat Cap:      $%.2f (%.0f%% of equity)\n", result.PortfolioCap, heatCapPct*100)
	PrintHumanf(format, "Utilization:   %.1f%%\n", result.PortfolioHeatPct)

	if result.PortfolioCapExceeded {
		PrintHumanf(format, "❌ EXCEEDED by $%.2f\n", result.PortfolioOverage)
	} else {
		PrintHuman(format, "✓ Within limits")
	}

	if addBucket != "" {
		PrintHuman(format, "")
		PrintHumanf(format, "Bucket Heat Status: %s\n", addBucket)
		PrintHuman(format, "=====================")
		PrintHumanf(format, "Current Heat:  $%.2f\n", result.CurrentBucketHeat)
		if addRisk > 0 {
			PrintHumanf(format, "New Heat:      $%.2f (+$%.2f)\n", result.NewBucketHeat, addRisk)
		} else {
			PrintHumanf(format, "New Heat:      $%.2f\n", result.NewBucketHeat)
		}
		PrintHumanf(format, "Bucket Cap:    $%.2f (%.1f%% of equity)\n", result.BucketCap, bucketHeatCapPct*100)
		PrintHumanf(format, "Utilization:   %.1f%%\n", result.BucketHeatPct)

		if result.BucketCapExceeded {
			PrintHumanf(format, "❌ EXCEEDED by $%.2f\n", result.BucketOverage)
		} else {
			PrintHuman(format, "✓ Within limits")
		}
	}

	PrintHuman(format, "")
	if result.Allowed {
		PrintHuman(format, "✓ Trade ALLOWED")
	} else {
		PrintHumanf(format, "❌ Trade REJECTED: %s\n", result.RejectionReason)
	}

	// List open positions
	if len(positions) > 0 {
		PrintHuman(format, "")
		PrintHumanf(format, "Open Positions: %d\n", len(positions))
		PrintHuman(format, "-------------------")
		for _, p := range positions {
			PrintHumanf(format, "• %s: $%.2f risk", p.Ticker, p.RiskDollars)
			if p.Bucket != "" {
				PrintHumanf(format, " (%s)", p.Bucket)
			}
			PrintHuman(format, "")
		}
	}

	// JSON output (always)
	PrintJSON(result)

	return nil
}
