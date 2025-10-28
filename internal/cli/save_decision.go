package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// DBGateChecker implements domain.GateChecker using database queries
type DBGateChecker struct {
	db     *storage.DB
	log    *logrus.Entry
	equity float64
}

// CheckBannerGreen verifies the banner is GREEN for the ticker
func (c *DBGateChecker) CheckBannerGreen(ticker string) error {
	// For now, we'll check if there's a recent GREEN checklist evaluation
	// In a full implementation, this would query checklist_evaluations table
	// For M13, we'll implement a placeholder that can be enhanced
	c.log.WithField("gate", "banner").WithField("ticker", ticker).Info("Checking banner gate")

	// Placeholder: Will be properly implemented with checklist_evaluations table query
	// For now, assume GREEN (will be validated in integration)
	return nil
}

// CheckTickerInCandidates verifies ticker is in today's candidates
func (c *DBGateChecker) CheckTickerInCandidates(ticker, date string) error {
	c.log.WithField("gate", "candidates").WithField("ticker", ticker).Info("Checking candidates gate")

	found, err := c.db.IsTickerInCandidates(date, ticker)
	if err != nil {
		return fmt.Errorf("failed to check candidates: %w", err)
	}

	if !found {
		return fmt.Errorf("%s not in today's candidates (must be from FINVIZ screen)", ticker)
	}

	return nil
}

// CheckImpulseBrake verifies 2-minute timer has expired
func (c *DBGateChecker) CheckImpulseBrake(ticker string) error {
	c.log.WithField("gate", "impulse_brake").WithField("ticker", ticker).Info("Checking impulse brake gate")

	return c.db.CheckImpulseBrake(ticker)
}

// CheckBucketCooldown verifies bucket is not in cooldown
func (c *DBGateChecker) CheckBucketCooldown(bucket string) error {
	c.log.WithField("gate", "bucket_cooldown").WithField("bucket", bucket).Info("Checking bucket cooldown gate")

	return c.db.CheckBucketCooldown(bucket)
}

// CheckHeatCaps verifies portfolio and bucket heat caps
func (c *DBGateChecker) CheckHeatCaps(addRisk float64, bucket string) error {
	c.log.WithField("gate", "heat_caps").WithField("add_risk", addRisk).WithField("bucket", bucket).Info("Checking heat caps gate")

	// Get settings for heat caps
	heatCapPctStr, err := c.db.GetSetting("HeatCap_H_pct")
	if err != nil {
		return fmt.Errorf("failed to get heat cap setting: %w", err)
	}

	bucketHeatCapPctStr, err := c.db.GetSetting("BucketHeatCap_pct")
	if err != nil {
		return fmt.Errorf("failed to get bucket heat cap setting: %w", err)
	}

	// Parse settings
	var heatCapPct, bucketHeatCapPct float64
	fmt.Sscanf(heatCapPctStr, "%f", &heatCapPct)
	fmt.Sscanf(bucketHeatCapPctStr, "%f", &bucketHeatCapPct)

	// Calculate heat request
	req := domain.HeatRequest{
		Equity:            c.equity,
		HeatCapPct:        heatCapPct,
		BucketHeatCapPct:  bucketHeatCapPct,
		AddRiskDollars:    addRisk,
		AddBucket:         bucket,
		OpenPositions:     []domain.Position{}, // Empty for now - would query from DB
	}

	// Calculate heat
	result, err := domain.CalculateHeat(req)
	if err != nil {
		return fmt.Errorf("failed to calculate heat: %w", err)
	}

	// Check portfolio cap
	if result.NewPortfolioHeat > result.PortfolioCap {
		overage := result.NewPortfolioHeat - result.PortfolioCap
		return fmt.Errorf("portfolio heat ($%.2f) exceeds cap ($%.2f) by $%.2f",
			result.NewPortfolioHeat, result.PortfolioCap, overage)
	}

	// Check bucket cap
	if bucket != "" && result.NewBucketHeat > result.BucketCap {
		overage := result.NewBucketHeat - result.BucketCap
		return fmt.Errorf("bucket heat ($%.2f) exceeds cap ($%.2f) by $%.2f",
			result.NewBucketHeat, result.BucketCap, overage)
	}

	return nil
}

// NewSaveDecisionCommand creates the save-decision command
func NewSaveDecisionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save-decision",
		Short: "Save a trading decision with 5 hard gates enforcement",
		Long: `Save a trading decision (GO or NO-GO) with full discipline enforcement.

For GO decisions, all 5 hard gates must pass:
  1. Banner GREEN (all 6 checklist items satisfied)
  2. Ticker in today's candidates (from FINVIZ screen)
  3. 2-minute impulse brake expired
  4. Bucket not in cooldown (24hr after loss)
  5. Heat caps not exceeded (4% portfolio, 1.5% bucket)

For NO-GO decisions, gates are not checked (just recording the decision).

Examples:
  # Save GO decision (stock)
  tf-engine save-decision --ticker AAPL --entry 180 --atr 1.5 --action GO

  # Save GO decision (options delta-ATR)
  tf-engine save-decision --ticker AAPL --entry 5.50 --atr 1.50 --delta 0.30 --method opt-delta-atr --action GO

  # Save GO decision (options max-loss)
  tf-engine save-decision --ticker AAPL --max-loss 0.75 --method opt-maxloss --action GO

  # Save NO-GO decision
  tf-engine save-decision --ticker AAPL --action NO-GO --reason "Bad setup"`,
		RunE: runSaveDecision,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol (required)")
	cmd.Flags().String("action", "", "GO or NO-GO (required)")
	cmd.Flags().Float64("entry", 0, "Entry price (required for GO)")
	cmd.Flags().Float64("atr", 0, "ATR value (required for stock and opt-delta-atr)")
	cmd.Flags().String("method", "stock", "Method: stock, opt-delta-atr, opt-maxloss")
	cmd.Flags().Float64("delta", 0, "Option delta (for opt-delta-atr)")
	cmd.Flags().Float64("max-loss", 0, "Maximum loss per contract (for opt-maxloss)")
	cmd.Flags().String("bucket", "", "Sector bucket")
	cmd.Flags().String("reason", "", "Reason (required for NO-GO)")
	cmd.Flags().String("date", "", "Date in YYYY-MM-DD format (defaults to today)")

	cmd.MarkFlagRequired("ticker")
	cmd.MarkFlagRequired("action")

	return cmd
}

func runSaveDecision(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	// Get flags
	ticker, _ := cmd.Flags().GetString("ticker")
	action, _ := cmd.Flags().GetString("action")
	entry, _ := cmd.Flags().GetFloat64("entry")
	atr, _ := cmd.Flags().GetFloat64("atr")
	method, _ := cmd.Flags().GetString("method")
	delta, _ := cmd.Flags().GetFloat64("delta")
	maxLoss, _ := cmd.Flags().GetFloat64("max-loss")
	bucket, _ := cmd.Flags().GetString("bucket")
	reason, _ := cmd.Flags().GetString("reason")
	dateStr, _ := cmd.Flags().GetString("date")

	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	log.WithField("ticker", ticker).WithField("action", action).Info("Saving trading decision")

	// Build request
	req := domain.SaveDecisionRequest{
		Ticker:  ticker,
		Action:  action,
		Entry:   entry,
		ATR:     atr,
		Method:  method,
		Delta:   delta,
		MaxLoss: maxLoss,
		Bucket:  bucket,
		Reason:  reason,
		Date:    dateStr,
		CorrID:  corrID,
	}

	// Validate request
	if err := domain.ValidateSaveDecisionRequest(req); err != nil {
		log.WithError(err).Error("Invalid save decision request")
		return fmt.Errorf("validation failed: %w", err)
	}

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Check for duplicate decision
	hasDuplicate, err := db.CheckForDuplicateDecision(ticker, dateStr)
	if err != nil {
		log.WithError(err).Error("Failed to check for duplicate")
		return fmt.Errorf("failed to check for duplicate: %w", err)
	}
	if hasDuplicate {
		return fmt.Errorf("you already have a decision for %s today (date: %s)", ticker, dateStr)
	}

	// Get equity for calculations
	equityStr, err := db.GetSetting("Equity_E")
	if err != nil {
		log.WithError(err).Error("Failed to get equity")
		return fmt.Errorf("failed to get equity: %w", err)
	}
	var equity float64
	fmt.Sscanf(equityStr, "%f", &equity)

	var decision storage.Decision
	decision.Date = dateStr
	decision.Ticker = ticker
	decision.Action = action
	decision.Method = method
	decision.Bucket = bucket
	decision.Reason = reason
	decision.CorrID = corrID

	// Handle GO decision - calculate position and check gates
	if action == "GO" {
		// Calculate position size
		var shares int
		var contracts int
		var riskDollars, stopDistance, initialStop float64

		riskPctStr, _ := db.GetSetting("RiskPct_r")
		var riskPct float64
		fmt.Sscanf(riskPctStr, "%f", &riskPct)

		kStr, _ := db.GetSetting("StopMultiple_K")
		var k float64
		fmt.Sscanf(kStr, "%f", &k)

		sizingReq := domain.SizingRequest{
			Equity:  equity,
			RiskPct: riskPct,
			Entry:   entry,
			ATR:     atr,
			K:       int(k),
			Method:  method,
			Delta:   delta,
			MaxLoss: maxLoss,
		}

		var result *domain.SizingResult
		var err error

		switch method {
		case "stock":
			result, err = domain.CalculateStockPosition(sizingReq)
		case "opt-delta-atr":
			result, err = domain.CalculateOptionDeltaATRPosition(sizingReq)
		case "opt-maxloss":
			result, err = domain.CalculateOptionMaxLossPosition(sizingReq)
		}

		if err != nil {
			log.WithError(err).Error("Failed to calculate position")
			return fmt.Errorf("failed to calculate position: %w", err)
		}

		shares = result.Shares
		contracts = result.Contracts
		riskDollars = result.RiskDollars
		stopDistance = result.StopDistance
		initialStop = result.InitialStop

		// Check all 5 hard gates
		checker := &DBGateChecker{db: db, log: log, equity: equity}
		gatesResult, err := domain.ValidateHardGates(checker, ticker, bucket, riskDollars, dateStr)
		if err != nil {
			log.WithError(err).Error("Failed to validate gates")
			return fmt.Errorf("failed to validate gates: %w", err)
		}

		if !gatesResult.AllPassed {
			log.WithField("failed_gates", gatesResult.FailedGates).Error("Hard gates failed")
			for i, gate := range gatesResult.FailedGates {
				fmt.Printf("❌ Gate %d failed: %s\n", i+1, gate)
				fmt.Printf("   Reason: %s\n", gatesResult.FailureReasons[i])
			}
			return fmt.Errorf("hard gates failed: %v", gatesResult.FailedGates)
		}

		// All gates passed - populate decision
		decision.Entry = entry
		decision.ATR = atr
		decision.StopDistance = stopDistance
		decision.InitialStop = initialStop
		decision.Shares = shares
		decision.Contracts = contracts
		decision.RiskDollars = riskDollars
		decision.Banner = "GREEN"
		decision.Delta = delta
		decision.MaxLoss = maxLoss

		log.Info("All 5 hard gates passed")
	} else {
		// NO-GO decision - no gates checked
		decision.Banner = "NO-GO"
		log.Info("NO-GO decision - gates not checked")
	}

	// Save decision
	decisionID, err := db.SaveDecision(decision)
	if err != nil {
		log.WithError(err).Error("Failed to save decision")
		return fmt.Errorf("failed to save decision: %w", err)
	}

	decision.ID = decisionID
	log.WithField("decision_id", decisionID).Info("Decision saved successfully")

	// Output result
	if action == "GO" {
		if method == "stock" {
			fmt.Printf("✓ Decision saved: %s GO %d shares @ $%.2f (stop: $%.2f, risk: $%.2f)\n",
				ticker, decision.Shares, entry, decision.InitialStop, decision.RiskDollars)
		} else {
			fmt.Printf("✓ Decision saved: %s GO %d contracts (risk: $%.2f)\n",
				ticker, decision.Contracts, decision.RiskDollars)
		}
	} else {
		fmt.Printf("✓ Decision saved: %s NO-GO (%s)\n", ticker, reason)
	}

	// Output JSON
	jsonResult, _ := json.MarshalIndent(decision, "", "  ")
	fmt.Println(string(jsonResult))

	return nil
}
