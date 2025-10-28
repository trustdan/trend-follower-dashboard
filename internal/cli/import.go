package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewImportCandidatesCommand creates the import-candidates command
func NewImportCandidatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-candidates",
		Short: "Import daily candidate tickers",
		Long: `Import a list of candidate tickers for trading consideration.
This ensures only screened stocks can pass the hard gates.

Examples:
  # Import candidates with a preset
  tf-engine import-candidates --tickers AAPL,MSFT,NVDA --preset TF_BREAKOUT_LONG

  # Import without preset
  tf-engine import-candidates --tickers TSLA,GOOGL

  # Import with sector and bucket information
  tf-engine import-candidates --tickers AAPL --preset TF_BREAKOUT_LONG --sector Technology --bucket Tech/Comm

  # Import for a specific date
  tf-engine import-candidates --tickers AAPL,MSFT --date 2025-10-26`,
		RunE: runImportCandidates,
	}

	cmd.Flags().String("tickers", "", "Comma-separated list of ticker symbols (required)")
	cmd.Flags().String("preset", "", "Preset name (e.g., TF_BREAKOUT_LONG)")
	cmd.Flags().String("sector", "", "Sector name")
	cmd.Flags().String("bucket", "", "Bucket name")
	cmd.Flags().String("date", "", "Date in YYYY-MM-DD format (defaults to today)")

	cmd.MarkFlagRequired("tickers")

	return cmd
}

func runImportCandidates(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	tickers, _ := cmd.Flags().GetString("tickers")
	preset, _ := cmd.Flags().GetString("preset")
	sector, _ := cmd.Flags().GetString("sector")
	bucket, _ := cmd.Flags().GetString("bucket")
	dateStr, _ := cmd.Flags().GetString("date")

	log.WithField("tickers", tickers).WithField("preset", preset).Info("Importing candidates")

	// Validate and normalize request
	req := domain.ImportCandidatesRequest{
		Tickers: tickers,
		Preset:  preset,
		Sector:  sector,
		Bucket:  bucket,
		Date:    dateStr,
	}

	if err := domain.ValidateImportRequest(req); err != nil {
		log.WithError(err).Error("Import request validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}

	// Normalize tickers
	normalizedTickers, err := domain.NormalizeTickers(req.Tickers)
	if err != nil {
		log.WithError(err).Error("Failed to normalize tickers")
		return fmt.Errorf("failed to normalize tickers: %w", err)
	}

	// Get import date
	importDate := domain.GetImportDate(req.Date)

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get or create preset if provided
	var presetID *int
	if preset != "" {
		id, err := db.GetOrCreatePreset(preset, "")
		if err != nil {
			log.WithError(err).WithField("preset", preset).Error("Failed to get/create preset")
			return fmt.Errorf("failed to get/create preset: %w", err)
		}
		presetID = &id
		log.WithField("preset_id", id).WithField("preset", preset).Info("Using preset")
	}

	// Import candidates
	err = db.ImportCandidates(importDate, normalizedTickers, presetID, sector, bucket)
	if err != nil {
		log.WithError(err).Error("Failed to import candidates")
		return fmt.Errorf("failed to import candidates: %w", err)
	}

	// Build result
	result := domain.ImportCandidatesResult{
		Count:      len(normalizedTickers),
		Date:       importDate,
		Tickers:    normalizedTickers,
		Preset:     preset,
		Sector:     sector,
		Bucket:     bucket,
		Normalized: true,
	}

	// Output JSON
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))

	log.WithField("count", result.Count).WithField("date", importDate).Info("Candidates imported successfully")

	return nil
}

// NewListCandidatesCommand creates the list-candidates command
func NewListCandidatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-candidates",
		Short: "List candidates for a specific date",
		Long: `List all imported candidates for a specific date.

Examples:
  # List today's candidates
  tf-engine list-candidates

  # List candidates for a specific date
  tf-engine list-candidates --date 2025-10-26`,
		RunE: runListCandidates,
	}

	cmd.Flags().String("date", "", "Date in YYYY-MM-DD format (defaults to today)")

	return cmd
}

func runListCandidates(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	dateStr, _ := cmd.Flags().GetString("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	log.WithField("date", dateStr).Info("Listing candidates")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get candidates
	candidates, err := db.GetCandidatesForDate(dateStr)
	if err != nil {
		log.WithError(err).Error("Failed to get candidates")
		return fmt.Errorf("failed to get candidates: %w", err)
	}

	// Build result
	result := map[string]interface{}{
		"date":       dateStr,
		"count":      len(candidates),
		"candidates": candidates,
	}

	// Output JSON
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))

	log.WithField("count", len(candidates)).Info("Candidates retrieved")

	return nil
}

// NewCheckCandidateCommand creates the check-candidate command
func NewCheckCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-candidate",
		Short: "Check if a ticker is in today's candidates",
		Long: `Check if a specific ticker is in the candidates list for today or a specific date.
This is used to enforce the "ticker in today's candidates" hard gate.

Examples:
  # Check if AAPL is in today's candidates
  tf-engine check-candidate --ticker AAPL

  # Check for a specific date
  tf-engine check-candidate --ticker AAPL --date 2025-10-26`,
		RunE: runCheckCandidate,
	}

	cmd.Flags().String("ticker", "", "Ticker symbol to check (required)")
	cmd.Flags().String("date", "", "Date in YYYY-MM-DD format (defaults to today)")

	cmd.MarkFlagRequired("ticker")

	return cmd
}

func runCheckCandidate(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	ticker, _ := cmd.Flags().GetString("ticker")
	dateStr, _ := cmd.Flags().GetString("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	// Normalize ticker
	normalizedTickers, err := domain.NormalizeTickers(ticker)
	if err != nil {
		log.WithError(err).Error("Failed to normalize ticker")
		return fmt.Errorf("failed to normalize ticker: %w", err)
	}
	ticker = normalizedTickers[0]

	log.WithField("ticker", ticker).WithField("date", dateStr).Info("Checking if ticker is in candidates")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Check if ticker is in candidates
	found, err := db.IsTickerInCandidates(dateStr, ticker)
	if err != nil {
		log.WithError(err).Error("Failed to check ticker")
		return fmt.Errorf("failed to check ticker: %w", err)
	}

	// Build result
	result := map[string]interface{}{
		"ticker": ticker,
		"date":   dateStr,
		"found":  found,
	}

	// Output JSON
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))

	log.WithField("ticker", ticker).WithField("found", found).Info("Ticker check complete")

	return nil
}
