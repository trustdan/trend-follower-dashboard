package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewSizeCommand creates the size command
func NewSizeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "size",
		Short: "Calculate position size for a trade",
		Long: `Calculate position size using the Van Tharp method.

This command calculates the number of shares or contracts to trade based on:
- Account equity
- Risk per trade (as % of equity)
- Entry price
- ATR (Average True Range)
- K multiple (stop distance multiplier)

The system enforces Van Tharp's position sizing rules:
  Risk dollars (R) = Equity × RiskPct
  Stop distance = K × ATR
  Initial stop = Entry - Stop distance
  Shares = floor(R ÷ Stop distance)

Examples:
  # Stock position sizing
  tf-engine size --entry 180 --atr 1.5 --k 2 --method stock

  # Use custom equity and risk
  tf-engine size --entry 180 --atr 1.5 --equity 20000 --risk 0.01 --method stock

  # Option sizing (delta-ATR method)
  tf-engine size --entry 5.00 --atr 0.50 --delta 0.50 --method opt-delta-atr`,
		RunE: runSize,
	}

	// Required flags
	cmd.Flags().Float64("entry", 0, "Entry price (required)")
	cmd.Flags().Float64("atr", 0, "Average True Range / ATR (required)")
	cmd.Flags().String("method", "stock", "Sizing method: stock, opt-delta-atr, opt-maxloss")

	// Optional flags (will use settings from DB if not provided)
	cmd.Flags().Float64("equity", 0, "Account equity (default: from database)")
	cmd.Flags().Float64("risk", 0, "Risk per trade as decimal (default: from database)")
	cmd.Flags().Int("k", 0, "Stop multiple (default: from database)")

	// Option-specific flags
	cmd.Flags().Float64("delta", 0, "Option delta (required for opt-delta-atr)")
	cmd.Flags().Float64("maxloss", 0, "Max loss per contract (required for opt-maxloss)")

	cmd.MarkFlagRequired("entry")
	cmd.MarkFlagRequired("atr")

	return cmd
}

func runSize(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	log.Info("Starting position sizing calculation")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get flags
	entry, _ := cmd.Flags().GetFloat64("entry")
	atr, _ := cmd.Flags().GetFloat64("atr")
	method, _ := cmd.Flags().GetString("method")
	equity, _ := cmd.Flags().GetFloat64("equity")
	risk, _ := cmd.Flags().GetFloat64("risk")
	k, _ := cmd.Flags().GetInt("k")
	delta, _ := cmd.Flags().GetFloat64("delta")
	maxloss, _ := cmd.Flags().GetFloat64("maxloss")

	// Load settings from database if not provided
	if equity == 0 {
		equityStr, err := db.GetSetting("Equity_E")
		if err != nil {
			log.WithError(err).Error("Failed to get Equity_E from database")
			return fmt.Errorf("failed to get equity from database: %w", err)
		}
		fmt.Sscanf(equityStr, "%f", &equity)
		log.WithField("equity", equity).Debug("Loaded equity from database")
	}

	if risk == 0 {
		riskStr, err := db.GetSetting("RiskPct_r")
		if err != nil {
			log.WithError(err).Error("Failed to get RiskPct_r from database")
			return fmt.Errorf("failed to get risk percent from database: %w", err)
		}
		fmt.Sscanf(riskStr, "%f", &risk)
		log.WithField("risk_pct", risk).Debug("Loaded risk percent from database")
	}

	if k == 0 {
		kStr, err := db.GetSetting("StopMultiple_K")
		if err != nil {
			log.WithError(err).Error("Failed to get StopMultiple_K from database")
			return fmt.Errorf("failed to get K multiple from database: %w", err)
		}
		fmt.Sscanf(kStr, "%d", &k)
		log.WithField("k", k).Debug("Loaded K multiple from database")
	}

	// Build request
	req := domain.SizingRequest{
		Equity:  equity,
		RiskPct: risk,
		Entry:   entry,
		ATR:     atr,
		K:       k,
		Method:  method,
		Delta:   delta,
		MaxLoss: maxloss,
	}

	log.WithField("request", req).Info("Calculating position size")

	// Calculate position size
	result, err := domain.CalculatePositionSize(req)
	if err != nil {
		log.WithError(err).Error("Position sizing calculation failed")
		return fmt.Errorf("position sizing failed: %w", err)
	}

	log.WithField("result", result).Info("Position sizing calculation completed")

	// Output JSON result
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.WithError(err).Error("Failed to marshal result to JSON")
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	fmt.Println(string(jsonResult))

	return nil
}
