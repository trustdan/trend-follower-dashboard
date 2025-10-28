package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/scrape"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewScrapeFinvizCommand creates the scrape-finviz command
func NewScrapeFinvizCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scrape-finviz",
		Short: "Scrape tickers from FINVIZ screener",
		Long: `Scrape ticker symbols from a FINVIZ screener URL.
Supports pagination, rate limiting, and automatic retry on failures.

The scraper will:
  - Extract tickers from all pages (up to --max-pages)
  - Normalize ticker symbols (uppercase, BRK.B -> BRK-B)
  - Remove duplicates
  - Optionally import tickers as candidates

Examples:
  # Scrape a FINVIZ screener URL
  tf-engine scrape-finviz --query "https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa,ta_sma50_pa"

  # Scrape and automatically import as candidates
  tf-engine scrape-finviz --query "<url>" --preset TF_BREAKOUT_LONG --import

  # Limit to first 5 pages
  tf-engine scrape-finviz --query "<url>" --max-pages 5

  # Custom rate limit (2 seconds between pages)
  tf-engine scrape-finviz --query "<url>" --rate-limit 2s`,
		RunE: runScrapeFinviz,
	}

	cmd.Flags().String("query", "", "FINVIZ screener URL (required)")
	cmd.Flags().String("preset", "", "Preset name for auto-import")
	cmd.Flags().Bool("import", false, "Auto-import scraped tickers as candidates")
	cmd.Flags().Int("max-pages", 10, "Maximum pages to scrape (0 = unlimited)")
	cmd.Flags().Duration("rate-limit", 1*time.Second, "Delay between page requests")
	cmd.Flags().Duration("timeout", 30*time.Second, "HTTP request timeout")
	cmd.Flags().Int("max-retries", 3, "Maximum retry attempts per page")

	cmd.MarkFlagRequired("query")

	return cmd
}

func runScrapeFinviz(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	queryURL, _ := cmd.Flags().GetString("query")
	preset, _ := cmd.Flags().GetString("preset")
	autoImport, _ := cmd.Flags().GetBool("import")
	maxPages, _ := cmd.Flags().GetInt("max-pages")
	rateLimit, _ := cmd.Flags().GetDuration("rate-limit")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	maxRetries, _ := cmd.Flags().GetInt("max-retries")

	log.WithField("query", queryURL).WithField("auto_import", autoImport).Info("Scraping FINVIZ")

	// Validate URL
	if err := scrape.ValidateFinvizURL(queryURL); err != nil {
		log.WithError(err).Error("Invalid FINVIZ URL")
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create scraper with config
	config := scrape.FinvizConfig{
		MaxPages:       maxPages,
		RateLimit:      rateLimit,
		RequestTimeout: timeout,
		MaxRetries:     maxRetries,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	scraper := scrape.NewFinvizScraper(config)

	// Scrape FINVIZ
	log.Info("Starting FINVIZ scrape")
	result, err := scraper.Scrape(queryURL)
	if err != nil {
		log.WithError(err).Error("Failed to scrape FINVIZ")
		return fmt.Errorf("scrape failed: %w", err)
	}

	log.WithField("count", result.Count).
		WithField("pages", result.PagesScraped).
		WithField("more_available", result.MoreAvailable).
		Info("FINVIZ scrape completed")

	// Auto-import if requested
	if autoImport {
		if len(result.Tickers) == 0 {
			log.Warn("No tickers to import (empty result)")
			fmt.Println("⚠ No tickers found to import")
		} else {
			log.WithField("count", result.Count).Info("Auto-importing tickers as candidates")

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
				id, err := db.GetOrCreatePreset(preset, queryURL)
				if err != nil {
					log.WithError(err).WithField("preset", preset).Error("Failed to get/create preset")
					return fmt.Errorf("failed to get/create preset: %w", err)
				}
				presetID = &id
				log.WithField("preset_id", id).WithField("preset", preset).Info("Using preset")
			}

			// Import candidates
			err = db.ImportCandidates(result.Date, result.Tickers, presetID, "", "")
			if err != nil {
				log.WithError(err).Error("Failed to import candidates")
				return fmt.Errorf("failed to import candidates: %w", err)
			}

			log.WithField("count", result.Count).Info("Tickers imported as candidates")
			fmt.Printf("✓ %d tickers imported as candidates for %s\n", result.Count, result.Date)
		}
	}

	// Build output result
	outputResult := map[string]interface{}{
		"tickers":        result.Tickers,
		"count":          result.Count,
		"date":           result.Date,
		"pages_scraped":  result.PagesScraped,
		"more_available": result.MoreAvailable,
		"normalized":     result.Normalized,
	}

	if autoImport {
		outputResult["imported"] = true
		outputResult["preset"] = preset
	}

	// Output JSON
	jsonResult, _ := json.MarshalIndent(outputResult, "", "  ")
	fmt.Println(string(jsonResult))

	return nil
}
