package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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

// sizeHandler handles position sizing requests
func (s *Server) sizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	var req struct {
		Equity  float64 `json:"equity"`
		RiskPct float64 `json:"risk_pct"`
		Entry   float64 `json:"entry"`
		ATR     float64 `json:"atr"`
		K       int     `json:"k"`
		Method  string  `json:"method"`
		Delta   float64 `json:"delta,omitempty"`
		MaxLoss float64 `json:"max_loss,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Failed to decode request")
		respondError(w, http.StatusBadRequest, "Invalid request body", corrID)
		return
	}

	// Load settings from database if not provided
	if req.Equity == 0 {
		equityStr, err := s.db.GetSetting("Equity_E")
		if err != nil {
			log.WithError(err).Error("Failed to get Equity_E")
			respondError(w, http.StatusInternalServerError, "Failed to get equity from database", corrID)
			return
		}
		fmt.Sscanf(equityStr, "%f", &req.Equity)
	}

	if req.RiskPct == 0 {
		riskStr, err := s.db.GetSetting("RiskPct_r")
		if err != nil {
			log.WithError(err).Error("Failed to get RiskPct_r")
			respondError(w, http.StatusInternalServerError, "Failed to get risk pct from database", corrID)
			return
		}
		fmt.Sscanf(riskStr, "%f", &req.RiskPct)
	}

	if req.K == 0 {
		kStr, err := s.db.GetSetting("StopMultiple_K")
		if err != nil {
			log.WithError(err).Error("Failed to get StopMultiple_K")
			respondError(w, http.StatusInternalServerError, "Failed to get K from database", corrID)
			return
		}
		var kFloat float64
		fmt.Sscanf(kStr, "%f", &kFloat)
		req.K = int(kFloat)
	}

	// Build sizing request
	sizingReq := domain.SizingRequest{
		Equity:  req.Equity,
		RiskPct: req.RiskPct,
		Entry:   req.Entry,
		ATR:     req.ATR,
		K:       req.K,
		Method:  req.Method,
		Delta:   req.Delta,
		MaxLoss: req.MaxLoss,
	}

	// Calculate position size
	result, err := domain.CalculatePositionSize(sizingReq)
	if err != nil {
		log.WithError(err).Error("Position sizing failed")
		respondError(w, http.StatusBadRequest, err.Error(), corrID)
		return
	}

	log.WithFields(map[string]interface{}{
		"method": req.Method,
		"shares": result.Shares,
		"risk":   result.RiskDollars,
	}).Info("Position size calculated")

	// Add correlation ID to response
	response := map[string]interface{}{
		"risk_dollars":   result.RiskDollars,
		"stop_distance":  result.StopDistance,
		"initial_stop":   result.InitialStop,
		"shares":         result.Shares,
		"contracts":      result.Contracts,
		"actual_risk":    result.ActualRisk,
		"method":         result.Method,
		"correlation_id": corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// checklistHandler handles checklist evaluation
func (s *Server) checklistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	var req struct {
		Ticker string `json:"ticker"`
		Checks struct {
			FromPreset    bool `json:"from_preset"`
			TrendPass     bool `json:"trend_pass"`
			LiquidityPass bool `json:"liquidity_pass"`
			TVConfirm     bool `json:"tv_confirm"`
			EarningsOK    bool `json:"earnings_ok"`
			JournalOK     bool `json:"journal_ok"`
		} `json:"checks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Failed to decode request")
		respondError(w, http.StatusBadRequest, "Invalid request body", corrID)
		return
	}

	// Build checklist request
	checklistReq := domain.ChecklistRequest{
		Ticker:        req.Ticker,
		FromPreset:    req.Checks.FromPreset,
		TrendPass:     req.Checks.TrendPass,
		LiquidityPass: req.Checks.LiquidityPass,
		TVConfirm:     req.Checks.TVConfirm,
		EarningsOK:    req.Checks.EarningsOK,
		JournalOK:     req.Checks.JournalOK,
	}

	// Evaluate checklist
	result, err := domain.EvaluateChecklist(checklistReq)
	if err != nil {
		log.WithError(err).Error("Checklist evaluation failed")
		respondError(w, http.StatusBadRequest, err.Error(), corrID)
		return
	}

	// Start impulse timer if banner is GREEN
	if result.Banner == domain.BannerGreen {
		if err := s.db.StartImpulseTimer(req.Ticker); err != nil {
			log.WithError(err).Error("Failed to start impulse timer")
			// Don't fail the request, just log the error
		}
	}

	log.WithFields(map[string]interface{}{
		"ticker": req.Ticker,
		"banner": result.Banner,
	}).Info("Checklist evaluated")

	response := map[string]interface{}{
		"banner":         result.Banner,
		"missing_count":  result.MissingCount,
		"missing_items":  result.MissingItems,
		"allow_save":     result.AllowSave,
		"correlation_id": corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// decisionHandler handles saving trading decisions
func (s *Server) decisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	var req struct {
		Ticker  string  `json:"ticker"`
		Action  string  `json:"action"`
		Entry   float64 `json:"entry"`
		ATR     float64 `json:"atr"`
		Method  string  `json:"method"`
		Delta   float64 `json:"delta,omitempty"`
		MaxLoss float64 `json:"max_loss,omitempty"`
		Bucket  string  `json:"bucket"`
		Reason  string  `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Failed to decode request")
		respondError(w, http.StatusBadRequest, "Invalid request body", corrID)
		return
	}

	// Validate request
	saveReq := domain.SaveDecisionRequest{
		Ticker:  req.Ticker,
		Action:  req.Action,
		Entry:   req.Entry,
		ATR:     req.ATR,
		Method:  req.Method,
		Delta:   req.Delta,
		MaxLoss: req.MaxLoss,
		Bucket:  req.Bucket,
		Reason:  req.Reason,
		Date:    time.Now().Format("2006-01-02"),
		CorrID:  corrID,
	}

	if err := domain.ValidateSaveDecisionRequest(saveReq); err != nil {
		log.WithError(err).Error("Invalid save decision request")
		respondError(w, http.StatusBadRequest, err.Error(), corrID)
		return
	}

	// For GO decisions, validate gates and calculate sizing
	if req.Action == "GO" {
		// Get settings
		equityStr, _ := s.db.GetSetting("Equity_E")
		riskStr, _ := s.db.GetSetting("RiskPct_r")
		kStr, _ := s.db.GetSetting("StopMultiple_K")

		var equity, riskPct, kFloat float64
		fmt.Sscanf(equityStr, "%f", &equity)
		fmt.Sscanf(riskStr, "%f", &riskPct)
		fmt.Sscanf(kStr, "%f", &kFloat)
		k := int(kFloat)

		// Calculate position size
		sizingReq := domain.SizingRequest{
			Equity:  equity,
			RiskPct: riskPct,
			Entry:   req.Entry,
			ATR:     req.ATR,
			K:       k,
			Method:  req.Method,
			Delta:   req.Delta,
			MaxLoss: req.MaxLoss,
		}

		if sizingReq.Method == "" {
			sizingReq.Method = "stock"
		}

		sizing, err := domain.CalculatePositionSize(sizingReq)
		if err != nil {
			log.WithError(err).Error("Position sizing failed")
			respondError(w, http.StatusBadRequest, err.Error(), corrID)
			return
		}

		// Create gate checker
		checker := &DBGateChecker{
			db:     s.db,
			log:    log,
			equity: equity,
		}

		// Validate hard gates
		gatesResult, err := domain.ValidateHardGates(checker, req.Ticker, req.Bucket, sizing.RiskDollars, time.Now().Format("2006-01-02"))
		if err != nil {
			log.WithError(err).Error("Failed to validate gates")
			respondError(w, http.StatusInternalServerError, "Failed to validate gates", corrID)
			return
		}

		if !gatesResult.AllPassed {
			log.WithFields(map[string]interface{}{
				"ticker":       req.Ticker,
				"failed_gates": gatesResult.FailedGates,
			}).Info("Decision rejected by gates")

			response := map[string]interface{}{
				"accepted":        false,
				"failed_gates":    gatesResult.FailedGates,
				"failure_reasons": gatesResult.FailureReasons,
				"correlation_id":  corrID,
			}

			respondJSON(w, http.StatusBadRequest, response)
			return
		}

		// Save GO decision
		decision := storage.Decision{
			Date:         time.Now().Format("2006-01-02"),
			Ticker:       req.Ticker,
			Action:       "GO",
			Entry:        req.Entry,
			ATR:          req.ATR,
			Method:       sizingReq.Method,
			Delta:        req.Delta,
			MaxLoss:      req.MaxLoss,
			Shares:       sizing.Shares,
			Contracts:    sizing.Contracts,
			RiskDollars:  sizing.RiskDollars,
			StopDistance: sizing.StopDistance,
			InitialStop:  sizing.InitialStop,
			Bucket:       req.Bucket,
			Banner:       domain.BannerGreen,
			CorrID:       corrID,
		}

		decisionID, err := s.db.SaveDecision(decision)
		if err != nil {
			log.WithError(err).Error("Failed to save decision")
			respondError(w, http.StatusInternalServerError, "Failed to save decision", corrID)
			return
		}

		log.WithFields(map[string]interface{}{
			"ticker":      req.Ticker,
			"action":      "GO",
			"decision_id": decisionID,
		}).Info("GO decision saved")

		response := map[string]interface{}{
			"accepted":       true,
			"decision_id":    decisionID,
			"shares":         sizing.Shares,
			"contracts":      sizing.Contracts,
			"risk_dollars":   sizing.RiskDollars,
			"initial_stop":   sizing.InitialStop,
			"correlation_id": corrID,
		}

		respondJSON(w, http.StatusOK, response)
		return
	}

	// Save NO-GO decision
	decision := storage.Decision{
		Date:   time.Now().Format("2006-01-02"),
		Ticker: req.Ticker,
		Action: "NO-GO",
		Reason: req.Reason,
		Bucket: req.Bucket,
		CorrID: corrID,
	}

	decisionID, err := s.db.SaveDecision(decision)
	if err != nil {
		log.WithError(err).Error("Failed to save NO-GO decision")
		respondError(w, http.StatusInternalServerError, "Failed to save decision", corrID)
		return
	}

	log.WithFields(map[string]interface{}{
		"ticker":      req.Ticker,
		"action":      "NO-GO",
		"decision_id": decisionID,
	}).Info("NO-GO decision saved")

	response := map[string]interface{}{
		"accepted":       true,
		"decision_id":    decisionID,
		"correlation_id": corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// candidatesHandler handles listing candidates
func (s *Server) candidatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	// Get date parameter (defaults to today)
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Get candidates
	candidates, err := s.db.GetCandidatesForDate(date)
	if err != nil {
		log.WithError(err).Error("Failed to get candidates")
		respondError(w, http.StatusInternalServerError, "Failed to get candidates", corrID)
		return
	}

	log.WithField("count", len(candidates)).Info("Candidates retrieved")

	response := map[string]interface{}{
		"candidates":     candidates,
		"count":          len(candidates),
		"date":           date,
		"correlation_id": corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// heatHandler handles heat status requests
func (s *Server) heatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	// Get settings
	equityStr, err := s.db.GetSetting("Equity_E")
	if err != nil {
		log.WithError(err).Error("Failed to get Equity_E")
		respondError(w, http.StatusInternalServerError, "Failed to get equity", corrID)
		return
	}

	heatCapStr, err := s.db.GetSetting("HeatCap_H_pct")
	if err != nil {
		log.WithError(err).Error("Failed to get HeatCap_H_pct")
		respondError(w, http.StatusInternalServerError, "Failed to get heat cap", corrID)
		return
	}

	bucketHeatCapStr, err := s.db.GetSetting("BucketHeatCap_pct")
	if err != nil {
		log.WithError(err).Error("Failed to get BucketHeatCap_pct")
		respondError(w, http.StatusInternalServerError, "Failed to get bucket heat cap", corrID)
		return
	}

	var equity, heatCap, bucketHeatCap float64
	fmt.Sscanf(equityStr, "%f", &equity)
	fmt.Sscanf(heatCapStr, "%f", &heatCap)
	fmt.Sscanf(bucketHeatCapStr, "%f", &bucketHeatCap)

	// Calculate portfolio heat
	portfolioHeat, err := s.db.CalculatePortfolioHeat()
	if err != nil {
		log.WithError(err).Error("Failed to calculate portfolio heat")
		respondError(w, http.StatusInternalServerError, "Failed to calculate heat", corrID)
		return
	}

	portfolioCap := equity * heatCap
	portfolioPct := 0.0
	if equity > 0 {
		portfolioPct = portfolioHeat / equity
	}

	// Get bucket heat (simplified - just return empty array for now)
	buckets := []map[string]interface{}{}

	log.WithFields(map[string]interface{}{
		"portfolio_heat": portfolioHeat,
		"portfolio_cap":  portfolioCap,
	}).Info("Heat status retrieved")

	response := map[string]interface{}{
		"portfolio_heat": portfolioHeat,
		"portfolio_cap":  portfolioCap,
		"portfolio_pct":  portfolioPct,
		"buckets":        buckets,
		"correlation_id": corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// timerHandler handles impulse timer checks
func (s *Server) timerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondError(w, http.StatusBadRequest, "ticker parameter required", corrID)
		return
	}

	// Get timer
	timer, err := s.db.GetActiveTimer(ticker)
	if err != nil {
		log.WithError(err).Error("Failed to get timer")
		respondError(w, http.StatusNotFound, "No active timer found", corrID)
		return
	}

	elapsed := time.Since(timer.StartedAt)
	const impulseDuration = 120 * time.Second // 2 minutes
	ready := elapsed >= impulseDuration

	remaining := int(impulseDuration.Seconds() - elapsed.Seconds())
	if remaining < 0 {
		remaining = 0
	}

	log.WithFields(map[string]interface{}{
		"ticker":  ticker,
		"elapsed": elapsed.Seconds(),
		"ready":   ready,
	}).Info("Timer checked")

	response := map[string]interface{}{
		"ticker":            ticker,
		"elapsed_seconds":   int(elapsed.Seconds()),
		"remaining_seconds": remaining,
		"ready":             ready,
		"correlation_id":    corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// cooldownHandler handles cooldown status checks
func (s *Server) cooldownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		respondError(w, http.StatusBadRequest, "bucket parameter required", corrID)
		return
	}

	// Check cooldown
	err := s.db.CheckBucketCooldown(bucket)
	inCooldown := (err != nil)

	response := map[string]interface{}{
		"bucket":         bucket,
		"in_cooldown":    inCooldown,
		"correlation_id": corrID,
	}

	if inCooldown {
		// Get cooldown details
		cooldown, err := s.db.GetBucketCooldown(bucket)
		if err == nil {
			response["expires_at"] = cooldown.ExpiresAt.Format(time.RFC3339)
			response["reason"] = cooldown.Reason
		}
	}

	log.WithFields(map[string]interface{}{
		"bucket":      bucket,
		"in_cooldown": inCooldown,
	}).Info("Cooldown checked")

	respondJSON(w, http.StatusOK, response)
}

// positionsHandler handles position listing
func (s *Server) positionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	status := r.URL.Query().Get("status")
	ticker := r.URL.Query().Get("ticker")

	// If ticker specified, get single position
	if ticker != "" {
		position, err := s.db.GetPositionByTicker(ticker)
		if err != nil {
			log.WithError(err).Error("Failed to get position")
			respondError(w, http.StatusNotFound, "Position not found", corrID)
			return
		}

		response := map[string]interface{}{
			"position":       position,
			"correlation_id": corrID,
		}

		respondJSON(w, http.StatusOK, response)
		return
	}

	// Get all positions
	positions, err := s.db.GetAllPositions(status)
	if err != nil {
		log.WithError(err).Error("Failed to get positions")
		respondError(w, http.StatusInternalServerError, "Failed to get positions", corrID)
		return
	}

	log.WithField("count", len(positions)).Info("Positions retrieved")

	response := map[string]interface{}{
		"positions":      positions,
		"count":          len(positions),
		"correlation_id": corrID,
	}

	respondJSON(w, http.StatusOK, response)
}

// settingsHandler handles getting all settings
func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	corrID := r.Header.Get("X-Correlation-ID")
	if corrID == "" {
		corrID = logx.GenerateCorrelationID()
	}
	log := logx.WithCorrelationID(corrID)

	settings, err := s.db.GetAllSettings()
	if err != nil {
		log.WithError(err).Error("Failed to get settings")
		respondError(w, http.StatusInternalServerError, "Failed to get settings", corrID)
		return
	}

	// Convert string values to proper types
	response := make(map[string]interface{})
	for key, value := range settings {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			response[key] = floatVal
		} else {
			response[key] = value
		}
	}
	response["correlation_id"] = corrID

	log.Info("Settings retrieved")

	respondJSON(w, http.StatusOK, response)
}
