package main

import (
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/scrape"
)

func buildScannerScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("FINVIZ Scanner", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Instructions
	instructions := widget.NewLabel("Import ticker candidates from FINVIZ screener URLs")
	instructions.Wrapping = fyne.TextWrapWord

	// Preset selector
	presetLabel := widget.NewLabel("Preset Name:")
	presetEntry := widget.NewEntry()
	presetEntry.SetPlaceHolder("TF_BREAKOUT_LONG")
	presetEntry.Text = "TF_BREAKOUT_LONG"

	// URL entry
	urlLabel := widget.NewLabel("FINVIZ Screener URL:")
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://finviz.com/screener.ashx?v=111&f=...")
	urlEntry.MultiLine = false

	// Common presets - All 5 FINVIZ presets from backend
	presetsLabel := widget.NewLabelWithStyle("Quick Presets:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	breakoutLongBtn := widget.NewButton("TF-Breakout-Long", func() {
		presetEntry.SetText("TF-Breakout-Long")
		urlEntry.SetText("https://finviz.com/screener.ashx?v=111&p=d&s=ta_newhigh&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume")
	})
	breakoutLongBtn.Importance = widget.HighImportance

	momentumUptrendBtn := widget.NewButton("TF-Momentum-Uptrend", func() {
		presetEntry.SetText("TF-Momentum-Uptrend")
		urlEntry.SetText("https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pa,ta_sma50_pa&dr=y1&o=-marketcap")
	})
	momentumUptrendBtn.Importance = widget.HighImportance

	unusualVolumeBtn := widget.NewButton("TF-Unusual-Volume", func() {
		presetEntry.SetText("TF-Unusual-Volume")
		urlEntry.SetText("https://finviz.com/screener.ashx?v=111&p=d&s=ta_unusualvolume&f=cap_largeover,sh_price_o20,ta_sma200_pa,ta_sma50_pa&o=-relativevolume")
	})
	unusualVolumeBtn.Importance = widget.HighImportance

	breakdownShortBtn := widget.NewButton("TF-Breakdown-Short", func() {
		presetEntry.SetText("TF-Breakdown-Short")
		urlEntry.SetText("https://finviz.com/screener.ashx?v=111&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&o=-relativevolume")
	})
	breakdownShortBtn.Importance = widget.HighImportance

	momentumDowntrendBtn := widget.NewButton("TF-Momentum-Downtrend", func() {
		presetEntry.SetText("TF-Momentum-Downtrend")
		urlEntry.SetText("https://finviz.com/screener.ashx?v=111&p=d&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma200_pb,ta_sma50_pb&dr=y1&o=-marketcap")
	})
	momentumDowntrendBtn.Importance = widget.HighImportance

	presetsRow1 := container.NewHBox(
		breakoutLongBtn,
		momentumUptrendBtn,
		unusualVolumeBtn,
	)

	presetsRow2 := container.NewHBox(
		breakdownShortBtn,
		momentumDowntrendBtn,
	)

	// Options
	maxPagesLabel := widget.NewLabel("Max Pages:")
	maxPagesEntry := widget.NewEntry()
	maxPagesEntry.SetText("10")

	rateLimitLabel := widget.NewLabel("Rate Limit (seconds):")
	rateLimitEntry := widget.NewEntry()
	rateLimitEntry.SetText("1")

	optionsRow := container.NewGridWithColumns(4,
		maxPagesLabel, maxPagesEntry,
		rateLimitLabel, rateLimitEntry,
	)

	// Status display
	statusLabel := widget.NewLabel("Status: Ready")
	statusLabel.Wrapping = fyne.TextWrapWord

	// Results display (will be replaced with ticker links)
	resultsContainer := container.NewVBox()
	resultsDateLabel := widget.NewLabel("")

	// Scan button (declare first so it can be referenced in closure)
	var scanBtn *widget.Button
	scanBtn = widget.NewButton("Scan FINVIZ & Import", func() {
		// Validate URL
		finvizURL := urlEntry.Text
		if finvizURL == "" {
			dialog.ShowError(fmt.Errorf("please enter a FINVIZ URL"), state.window)
			return
		}

		if err := scrape.ValidateFinvizURL(finvizURL); err != nil {
			dialog.ShowError(fmt.Errorf("invalid URL: %w", err), state.window)
			return
		}

		// Disable button during scan
		scanBtn.Disable()
		statusLabel.SetText("Status: Scanning FINVIZ...")
		resultsContainer.Objects = []fyne.CanvasObject{}
		resultsDateLabel.SetText("")
		resultsContainer.Refresh()

		// Run scan in goroutine
		go func() {
			// Parse options
			maxPages := 10
			fmt.Sscanf(maxPagesEntry.Text, "%d", &maxPages)

			rateLimit := 1
			fmt.Sscanf(rateLimitEntry.Text, "%d", &rateLimit)

			// Create scraper config
			config := scrape.FinvizConfig{
				MaxPages:       maxPages,
				RateLimit:      time.Duration(rateLimit) * time.Second,
				RequestTimeout: 30 * time.Second,
				MaxRetries:     3,
				UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			}

			scraper := scrape.NewFinvizScraper(config)

			// Scrape
			result, err := scraper.Scrape(finvizURL)
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Status: Error - %v", err))
				scanBtn.Enable()
				return
			}

			statusLabel.SetText(fmt.Sprintf("Status: Found %d tickers from %d pages", result.Count, result.PagesScraped))

			// Import to database
			preset := presetEntry.Text
			if preset == "" {
				preset = "CUSTOM"
			}

			presetID, err := state.db.GetOrCreatePreset(preset, finvizURL)
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Status: Error creating preset - %v", err))
				scanBtn.Enable()
				return
			}

			err = state.db.ImportCandidates(result.Date, result.Tickers, &presetID, "", "")
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Status: Error importing - %v", err))
				scanBtn.Enable()
				return
			}

			// Success
			statusLabel.SetText(fmt.Sprintf("Status: Successfully imported %d tickers", result.Count))
			resultsDateLabel.SetText(fmt.Sprintf("Date: %s", result.Date))

			// Create clickable TradingView links for each ticker
			resultsContainer.Objects = []fyne.CanvasObject{}

			// Add tickers in rows of 5
			tickerRow := container.NewHBox()
			for i, ticker := range result.Tickers {
				// Create clickable hyperlink
				tvURLStr := fmt.Sprintf("https://www.tradingview.com/chart/?symbol=%s", ticker)
				tvURL, _ := url.Parse(tvURLStr)
				link := widget.NewHyperlink(ticker, tvURL)
				tickerRow.Add(link)

				// Start new row after every 5 tickers
				if (i+1)%5 == 0 {
					resultsContainer.Add(tickerRow)
					tickerRow = container.NewHBox()
				}
			}
			// Add any remaining tickers in the last row
			if len(tickerRow.Objects) > 0 {
				resultsContainer.Add(tickerRow)
			}
			resultsContainer.Refresh()

			scanBtn.Enable()
		}()
	})
	scanBtn.Importance = widget.HighImportance

	// Create scroll container for results with minimum height
	resultsScroll := container.NewScroll(resultsContainer)
	resultsScroll.SetMinSize(fyne.NewSize(400, 200))  // Min width 400, height 200 pixels

	// Layout
	form := container.NewVBox(
		container.NewPadded(title),
		instructions,
		widget.NewSeparator(),
		presetsLabel,
		presetsRow1,
		presetsRow2,
		widget.NewSeparator(),
		presetLabel,
		presetEntry,
		urlLabel,
		urlEntry,
		widget.NewSeparator(),
		optionsRow,
		widget.NewSeparator(),
		scanBtn,
		widget.NewSeparator(),
		statusLabel,
		resultsDateLabel,
		widget.NewLabel("Click tickers to view on TradingView:"),
		resultsScroll,
	)

	return container.NewScroll(form)
}
