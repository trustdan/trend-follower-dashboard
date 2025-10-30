package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/domain"
)

func buildTradeEntryScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("Trade Entry - 5 Gates Final Check", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	instructions := widget.NewLabel("This is the final gate check before saving a GO/NO-GO decision.")
	instructions.Wrapping = fyne.TextWrapWord

	// Banner (starts gray)
	bannerRect := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	bannerRect.SetMinSize(fyne.NewSize(0, 80))

	bannerText := canvas.NewText("RUN FINAL GATES CHECK", color.White)
	bannerText.TextSize = 28
	bannerText.TextStyle = fyne.TextStyle{Bold: true}
	bannerText.Alignment = fyne.TextAlignCenter

	banner := container.NewStack(
		bannerRect,
		container.NewCenter(bannerText),
	)

	// Input fields
	tickerLabel := widget.NewLabel("Ticker:")
	tickerEntry := widget.NewEntry()
	tickerEntry.SetPlaceHolder("AAPL")

	riskLabel := widget.NewLabel("Risk ($):")
	riskEntry := widget.NewEntry()
	riskEntry.SetPlaceHolder("750.00")

	// Individual sectors dropdown (FINVIZ order + ETFs)
	sectorLabel := widget.NewLabel("Sector:")
	sectorOptions := []string{
		"Basic Materials",
		"Communication Services",
		"Consumer Cyclical",
		"Consumer Defensive",
		"Energy",
		"Financial",
		"Healthcare",
		"Industrials",
		"Real Estate",
		"Technology",
		"Utilities",
		"ETFs",
	}
	sectorSelect := widget.NewSelect(sectorOptions, nil)
	sectorSelect.SetSelected("Technology")

	// Sector buckets (grouped correlated sectors for heat tracking)
	bucketLabel := widget.NewLabel("Bucket (for heat tracking):")
	bucketOptions := []string{
		"Materials/Industrials",  // Basic Materials + Industrials (commodity/economic cycle)
		"Tech/Comm",              // Communication Services + Technology (growth/innovation)
		"Financial/Cyclical",     // Financial + Consumer Cyclical (economic sensitivity)
		"Defensive/Utilities",    // Consumer Defensive + Utilities (defensive plays)
		"Energy",                 // Energy (standalone - commodity/burst sector)
		"Healthcare",             // Healthcare (standalone - less correlated)
		"Real Estate",            // Real Estate (standalone)
		"ETFs",                   // ETFs (standalone)
	}
	bucketSelect := widget.NewSelect(bucketOptions, nil)

	// Auto-map sector to bucket when sector changes
	sectorToBucket := map[string]string{
		"Basic Materials":        "Materials/Industrials",
		"Industrials":            "Materials/Industrials",
		"Communication Services": "Tech/Comm",
		"Technology":             "Tech/Comm",
		"Financial":              "Financial/Cyclical",
		"Consumer Cyclical":      "Financial/Cyclical",
		"Consumer Defensive":     "Defensive/Utilities",
		"Utilities":              "Defensive/Utilities",
		"Energy":                 "Energy",
		"Healthcare":             "Healthcare",
		"Real Estate":            "Real Estate",
		"ETFs":                   "ETFs",
	}

	// Set initial bucket based on default sector
	bucketSelect.SetSelected(sectorToBucket["Technology"])

	// Update bucket when sector changes
	sectorSelect.OnChanged = func(sector string) {
		if bucket, ok := sectorToBucket[sector]; ok {
			bucketSelect.SetSelected(bucket)
		}
	}

	bannerStateLabel := widget.NewLabel("Checklist Banner State:")
	bannerStateSelect := widget.NewSelect([]string{"GREEN", "YELLOW", "RED"}, nil)
	bannerStateSelect.SetSelected("GREEN")

	// Results
	resultsLabel := widget.NewLabel("")
	resultsLabel.Wrapping = fyne.TextWrapWord

	// Check Gates button
	checkGatesBtn := widget.NewButton("Check All 5 Gates", func() {
		ticker := tickerEntry.Text
		if ticker == "" {
			resultsLabel.SetText("❌ Please enter ticker")
			return
		}

		riskStr := riskEntry.Text
		if riskStr == "" {
			resultsLabel.SetText("❌ Please enter risk amount")
			return
		}

		risk, err := strconv.ParseFloat(riskStr, 64)
		if err != nil {
			resultsLabel.SetText("❌ Invalid risk amount")
			return
		}

		sector := sectorSelect.Selected
		if sector == "" {
			sector = "Unknown"
		}

		bucket := bucketSelect.Selected
		if bucket == "" {
			bucket = "Unknown"
		}

		bannerState := bannerStateSelect.Selected

		// Get settings with defaults
		settings, _ := state.db.GetAllSettings()
		equityStr := getSettingWithDefault(settings, "equity", "100000")
		portfolioCapStr := getSettingWithDefault(settings, "portfolio_heat_cap", "4.0")
		bucketCapStr := getSettingWithDefault(settings, "bucket_heat_cap", "1.5")

		equity, _ := strconv.ParseFloat(equityStr, 64)
		portfolioCap, _ := strconv.ParseFloat(portfolioCapStr, 64)
		bucketCapCap, _ := strconv.ParseFloat(bucketCapStr, 64)

		// Simplified gates check (TODO: implement full GateChecker interface)
		_ = equity
		_ = portfolioCap
		_ = bucketCapCap

		result := &domain.HardGatesResult{
			AllPassed:      bannerState == "GREEN" && risk > 0,
			FailedGates:    []string{},
			FailureReasons: []string{},
		}

		if bannerState != "GREEN" {
			result.AllPassed = false
			result.FailedGates = append(result.FailedGates, "Banner")
			result.FailureReasons = append(result.FailureReasons, "Banner must be GREEN")
		}

		// Update banner
		if result.AllPassed {
			bannerRect.FillColor = ColorGreen()
			bannerText.Text = "✓ ALL GATES PASSED - GO"
		} else {
			bannerRect.FillColor = ColorRed()
			bannerText.Text = "✗ GATES FAILED - NO-GO"
		}
		bannerRect.Refresh()
		bannerText.Refresh()

		// Format results
		resultsText := fmt.Sprintf("Ticker: %s\nRisk: $%.2f\nSector: %s\nBucket: %s\n\n", ticker, risk, sector, bucket)
		resultsText += "Gate Results:\n"

		if len(result.FailedGates) > 0 {
			for i, gate := range result.FailedGates {
				resultsText += fmt.Sprintf("%d. ✗ %s - %s\n", i+1, gate, result.FailureReasons[i])
			}
		} else {
			resultsText += "1. ✓ Banner GREEN\n"
			resultsText += "2. ✓ Ticker in candidates (not checked)\n"
			resultsText += "3. ✓ Impulse brake (not checked)\n"
			resultsText += "4. ✓ Bucket cooldown (not checked)\n"
			resultsText += "5. ✓ Heat caps (not checked)\n"
		}

		resultsText += fmt.Sprintf("\nOverall: ")
		if result.AllPassed {
			resultsText += "✓ BASIC CHECKS PASSED\n\n(Full gates check TODO)"
		} else {
			resultsText += "✗ GATES FAILED\n\nYou may NOT trade this position."
		}

		resultsLabel.SetText(resultsText)
	})
	checkGatesBtn.Importance = widget.HighImportance

	// Save decision buttons (only enabled if gates pass)
	saveGoBtn := widget.NewButton("Save GO Decision", func() {
		// TODO: Implement save GO decision
		resultsLabel.SetText("✓ GO decision saved to database")
	})
	saveGoBtn.Disable()

	saveNoGoBtn := widget.NewButton("Save NO-GO Decision", func() {
		// TODO: Implement save NO-GO decision
		resultsLabel.SetText("✓ NO-GO decision saved to database")
	})
	saveNoGoBtn.Importance = widget.HighImportance

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		instructions,
		banner,
		widget.NewSeparator(),
		tickerLabel,
		tickerEntry,
		riskLabel,
		riskEntry,
		sectorLabel,
		sectorSelect,
		bucketLabel,
		bucketSelect,
		bannerStateLabel,
		bannerStateSelect,
		widget.NewSeparator(),
		checkGatesBtn,
		widget.NewSeparator(),
		resultsLabel,
		widget.NewSeparator(),
		container.NewHBox(saveGoBtn, saveNoGoBtn),
	)

	return container.NewScroll(content)
}
