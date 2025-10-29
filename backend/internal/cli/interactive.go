package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/scrape"
	"github.com/yourusername/trading-engine/internal/storage"
)

// ASCII Art and Visual Elements
const banner = `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   ████████╗███████╗      ███████╗███╗   ██╗ ██████╗      ║
║   ╚══██╔══╝██╔════╝      ██╔════╝████╗  ██║██╔════╝      ║
║      ██║   █████╗  █████╗█████╗  ██╔██╗ ██║██║  ███╗     ║
║      ██║   ██╔══╝  ╚════╝██╔══╝  ██║╚██╗██║██║   ██║     ║
║      ██║   ██║           ███████╗██║ ╚████║╚██████╔╝     ║
║      ╚═╝   ╚═╝           ╚══════╝╚═╝  ╚═══╝ ╚═════╝      ║
║                                                           ║
║        📊 FINVIZ SCRAPER 9000 📈 Market Data Extraction  ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`

const divider = "════════════════════════════════════════════════════════════"

// Loading animations
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var progressFrames = []string{
	"[          ]",
	"[■         ]",
	"[■■        ]",
	"[■■■       ]",
	"[■■■■      ]",
	"[■■■■■     ]",
	"[■■■■■■    ]",
	"[■■■■■■■   ]",
	"[■■■■■■■■  ]",
	"[■■■■■■■■■ ]",
	"[■■■■■■■■■■]",
}

// showBanner displays the ASCII art banner
func showBanner() {
	// Clear screen for Windows
	fmt.Print("\033[2J\033[H")
	fmt.Println(banner)
	time.Sleep(500 * time.Millisecond)
}

// showProgress displays an animated progress bar
func showProgress(message string, duration time.Duration) {
	steps := len(progressFrames)
	stepDuration := duration / time.Duration(steps)

	for i, frame := range progressFrames {
		percent := (i + 1) * 100 / steps
		fmt.Printf("\r%s %s %d%% ", message, frame, percent)
		time.Sleep(stepDuration)
	}
	fmt.Print("\r")
}

// showSpinner displays a spinning animation
func showSpinner(message string, done chan bool) {
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\r%s %s ", spinnerFrames[i%len(spinnerFrames)], message)
			time.Sleep(80 * time.Millisecond)
			i++
		}
	}
}

// printSuccess prints a success message with styling
func printSuccess(message string) {
	fmt.Printf("\n✅ %s\n", message)
}

// printInfo prints an info message with styling
func printInfo(message string) {
	fmt.Printf("\n💡 %s\n", message)
}

// printWarning prints a warning message with styling
func printWarning(message string) {
	fmt.Printf("\n⚠️  %s\n", message)
}

// printError prints an error message with styling
func printError(message string) {
	fmt.Printf("\n❌ %s\n", message)
}

// printHeader prints a section header
func printHeader(title string) {
	fmt.Printf("\n%s\n", divider)
	fmt.Printf("  📋 %s\n", title)
	fmt.Printf("%s\n\n", divider)
}

// printTickerBox displays tickers in a fancy box
func printTickerBox(tickers []string, count int) {
	fmt.Println("\n┌────────────────────────────────────────────────────┐")
	fmt.Printf("│  📈 Tickers Found: %-32d │\n", count)
	fmt.Println("├────────────────────────────────────────────────────┤")

	limit := 10
	if count < limit {
		limit = count
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("│  💰 %-46s │\n", tickers[i])
	}

	if count > 10 {
		fmt.Printf("│  ... and %d more                                   │\n", count-10)
	}

	fmt.Println("└────────────────────────────────────────────────────┘")
}

// printFinalStats displays final statistics with flair
func printFinalStats(count int, pages int, date string, preset string) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                          ║")
	fmt.Println("║              ✅ IMPORT COMPLETE! ✅                      ║")
	fmt.Println("║                                                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  📊 Tickers Imported:  %d\n", count)
	fmt.Printf("  📅 Date:              %s\n", date)
	fmt.Printf("  📋 Preset:            %s\n", preset)
	fmt.Printf("  📄 Pages Scraped:     %d\n", pages)
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║  🎯 NEXT STEPS                                           ║")
	fmt.Println("╠══════════════════════════════════════════════════════════╣")
	fmt.Println("║  1. Open TradingPlatform.xlsm                            ║")
	fmt.Println("║  2. Check Dashboard for today's candidates               ║")
	fmt.Println("║  3. Evaluate trades using Trade Entry sheet              ║")
	fmt.Println("║                                                          ║")
	fmt.Println("║  💪 The 5 Hard Gates will enforce discipline!            ║")
	fmt.Println("║  🚫 Only trade tickers from today's candidates!          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// Predefined FINVIZ screener presets
var finvizPresets = map[string]string{
	"TF-Breakout-Long":        "https://finviz.com/screener.ashx?v=111&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume",
	"TF-Momentum-Uptrend":     "https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&dr=y1&o=-marketcap",
	"TF-Unusual-Volume":       "https://finviz.com/screener.ashx?v=111&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume",
	"TF-Breakdown-Short":      "https://finviz.com/screener.ashx?v=111&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&o=-relativevolume",
	"TF-Momentum-Downtrend":   "https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&dr=y1&o=-marketcap",
	"Enter Custom URL":        "",
}

// NewInteractiveCommand creates the interactive command
func NewInteractiveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode for candidate import",
		Long: `Launch interactive mode to scrape FINVIZ and import candidates.

This guides you through:
  1. Selecting a FINVIZ screener preset (or entering custom URL)
  2. Configuring scraper options (pages, rate limit, etc.)
  3. Scraping the screener
  4. Reviewing results
  5. Importing as today's candidates

No need to remember command-line flags - just follow the prompts!

Examples:
  # Launch interactive mode
  tf-engine interactive

  # Launch with auto-defaults (skip confirmations)
  tf-engine interactive --auto`,
		RunE: runInteractive,
	}

	cmd.Flags().Bool("auto", false, "Use default options without prompts")

	return cmd
}

func runInteractive(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)
	auto, _ := cmd.Flags().GetBool("auto")

	// Show epic banner
	showBanner()

	printInfo("Launching interactive candidate import wizard...")
	showProgress("Initializing", 1*time.Second)

	// Step 1: Select preset
	printHeader("STEP 1/5: Select FINVIZ Screener 🎯")

	presetNames := []string{
		"TF-Breakout-Long",
		"TF-Momentum-Uptrend",
		"TF-Unusual-Volume",
		"TF-Breakdown-Short",
		"TF-Momentum-Downtrend",
		"Enter Custom URL",
	}

	var selectedPreset string
	var queryURL string

	if auto {
		selectedPreset = presetNames[0]
		queryURL = finvizPresets[selectedPreset]
		fmt.Printf("Using default: %s\n", selectedPreset)
	} else {
		promptSelect := promptui.Select{
			Label: "Choose a screener preset",
			Items: presetNames,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U0001F449 {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "\U00002705 {{ . | green }}",
			},
		}

		idx, result, err := promptSelect.Run()
		if err != nil {
			return fmt.Errorf("selection cancelled: %w", err)
		}

		selectedPreset = result
		queryURL = finvizPresets[selectedPreset]

		// If custom URL selected, prompt for it
		if idx == len(presetNames)-1 {
			promptURL := promptui.Prompt{
				Label:   "Enter FINVIZ screener URL",
				Default: "https://finviz.com/screener.ashx?v=111&f=",
				Validate: func(input string) error {
					if !strings.HasPrefix(input, "https://finviz.com/") {
						return fmt.Errorf("must be a valid FINVIZ URL")
					}
					return nil
				},
			}

			queryURL, err = promptURL.Run()
			if err != nil {
				return fmt.Errorf("URL entry cancelled: %w", err)
			}
		}
	}

	printSuccess(fmt.Sprintf("Selected: %s", selectedPreset))
	printInfo(fmt.Sprintf("URL: %s", queryURL))

	// Step 2: Configure scraper options
	printHeader("STEP 2/5: Configure Scraper Options ⚙️")

	maxPages := 10
	rateLimit := 1 * time.Second
	maxRetries := 3

	if !auto {
		// Max pages
		promptPages := promptui.Prompt{
			Label:   "Max pages to scrape (0 = unlimited)",
			Default: "10",
			Validate: func(input string) error {
				var n int
				_, err := fmt.Sscan(input, &n)
				if err != nil || n < 0 {
					return fmt.Errorf("must be a non-negative number")
				}
				return nil
			},
		}

		pagesStr, err := promptPages.Run()
		if err == nil {
			fmt.Sscan(pagesStr, &maxPages)
		}

		// Rate limit
		promptRate := promptui.Select{
			Label: "Rate limit between pages",
			Items: []string{"0.5 seconds", "1 second (recommended)", "2 seconds", "3 seconds"},
			Templates: &promptui.SelectTemplates{
				Active:   "\U0001F449 {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "\U00002705 {{ . | green }}",
			},
		}

		_, rateStr, err := promptRate.Run()
		if err == nil {
			switch rateStr {
			case "0.5 seconds":
				rateLimit = 500 * time.Millisecond
			case "1 second (recommended)":
				rateLimit = 1 * time.Second
			case "2 seconds":
				rateLimit = 2 * time.Second
			case "3 seconds":
				rateLimit = 3 * time.Second
			}
		}
	}

	fmt.Println()
	printSuccess("Configuration Complete!")
	fmt.Printf("  ⚡ Max Pages: %d\n", maxPages)
	fmt.Printf("  ⏱️  Rate Limit: %v\n", rateLimit)
	fmt.Printf("  🔄 Max Retries: %d\n", maxRetries)

	// Step 3: Confirm and scrape
	printHeader("STEP 3/5: Scrape FINVIZ 🚀")

	if !auto {
		promptConfirm := promptui.Prompt{
			Label:     "Start scraping",
			IsConfirm: true,
		}

		_, err := promptConfirm.Run()
		if err != nil {
			fmt.Println("\nScraping cancelled.")
			return nil
		}
	}

	// Validate URL
	if err := scrape.ValidateFinvizURL(queryURL); err != nil {
		log.WithError(err).Error("Invalid FINVIZ URL")
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create scraper
	config := scrape.FinvizConfig{
		MaxPages:       maxPages,
		RateLimit:      rateLimit,
		RequestTimeout: 30 * time.Second,
		MaxRetries:     maxRetries,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	scraper := scrape.NewFinvizScraper(config)

	// Scrape with animation
	fmt.Println()
	printInfo("🌐 Connecting to FINVIZ...")
	showProgress("Establishing connection", 800*time.Millisecond)

	printInfo("📡 Fetching screener data...")

	// Start spinner in background
	done := make(chan bool)
	go showSpinner("Scraping pages", done)

	result, err := scraper.Scrape(queryURL)
	done <- true

	if err != nil {
		log.WithError(err).Error("Failed to scrape FINVIZ")
		printError("Scraping failed!")
		return fmt.Errorf("scrape failed: %w", err)
	}

	printSuccess("Scraping complete!")

	// Step 4: Review results
	printHeader("STEP 4/5: Review Results 📊")

	fmt.Printf("  📈 Tickers Found:   %d\n", result.Count)
	fmt.Printf("  📄 Pages Scraped:   %d\n", result.PagesScraped)
	fmt.Printf("  📅 Date:            %s\n", result.Date)

	if result.MoreAvailable {
		printWarning("More tickers available (reached max pages limit)")
	}

	if result.Count > 0 {
		printTickerBox(result.Tickers, result.Count)
	} else {
		printWarning("No tickers found!")
		return nil
	}

	// Step 5: Import as candidates
	printHeader("STEP 5/5: Import Candidates 💾")

	shouldImport := auto

	if !auto {
		promptImport := promptui.Prompt{
			Label:     fmt.Sprintf("💾 Import %d tickers as today's candidates", result.Count),
			IsConfirm: true,
		}

		_, err := promptImport.Run()
		shouldImport = (err == nil)
	}

	if !shouldImport {
		printWarning("Import cancelled. Tickers not saved.")
		return nil
	}

	// Get preset name for database
	var presetName string
	if !auto {
		promptPresetName := promptui.Prompt{
			Label:   "Enter preset name (for tracking)",
			Default: strings.ReplaceAll(selectedPreset, " ", "_"),
		}

		presetName, err = promptPresetName.Run()
		if err != nil {
			presetName = "manual"
		}
	} else {
		presetName = "trend_following"
	}

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get or create preset
	presetID, err := db.GetOrCreatePreset(presetName, queryURL)
	if err != nil {
		log.WithError(err).WithField("preset", presetName).Error("Failed to get/create preset")
		return fmt.Errorf("failed to get/create preset: %w", err)
	}

	// Import candidates with progress animation
	fmt.Println()
	printInfo("💾 Saving to database...")
	showProgress("Importing tickers", 1500*time.Millisecond)

	err = db.ImportCandidates(result.Date, result.Tickers, &presetID, "", "")
	if err != nil {
		log.WithError(err).Error("Failed to import candidates")
		printError("Database import failed!")
		return fmt.Errorf("failed to import candidates: %w", err)
	}

	printSuccess("Database updated successfully!")

	// Show final stats with epic flair
	printFinalStats(result.Count, result.PagesScraped, result.Date, presetName)

	return nil
}
