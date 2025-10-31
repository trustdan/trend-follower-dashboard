package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/storage"
)

func buildHeatCheckScreen(state *AppState) fyne.CanvasObject {
	// Get active session (real or sample)
	activeSession := state.GetActiveSession()

	// Session check: require active session
	if activeSession == nil && !state.sampleMode {
		return showNoSessionPrompt(state, "Heat Check")
	}

	// Prerequisite check: sizing must be completed (skip in sample mode)
	if !state.sampleMode && !activeSession.SizingCompleted {
		return showPrerequisiteError(state, "Position Sizing", "Heat Check")
	}

	// Title
	title := canvas.NewText("Heat Check - Portfolio Risk Management", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Session info with sample mode indicator
	sessionInfoText := fmt.Sprintf(
		"Session #%d • %s • %s • Risk: $%.2f",
		activeSession.SessionNum,
		formatStrategyDisplay(activeSession.Strategy),
		activeSession.Ticker,
		activeSession.SizingRiskDollars,
	)
	if state.sampleMode {
		sessionInfoText = "📦 SAMPLE MODE: " + sessionInfoText
	}
	sessionInfo := widget.NewLabel(sessionInfoText)
	sessionInfo.TextStyle = fyne.TextStyle{Bold: true}

	// Explanation panel
	explanationCard := buildHeatCheckExplanation()

	// Load settings with defaults
	settings, _ := state.db.GetAllSettings()
	equityStr := getSettingWithDefault(settings, "equity", "100000")
	portfolioCapStr := getSettingWithDefault(settings, "portfolio_heat_cap", "4.0")
	bucketCapStr := getSettingWithDefault(settings, "bucket_heat_cap", "1.5")

	equity, _ := strconv.ParseFloat(equityStr, 64)
	portfolioCap, _ := strconv.ParseFloat(portfolioCapStr, 64)
	bucketCap, _ := strconv.ParseFloat(bucketCapStr, 64)

	settingsInfo := widget.NewLabel(fmt.Sprintf("Account: $%s | Portfolio Cap: %s%% | Bucket Cap: %s%%",
		equityStr, portfolioCapStr, bucketCapStr))

	// Calculate current heat (use sample data if in sample mode)
	var positions []storage.Position
	if state.sampleMode {
		positions = CreateSamplePositions()
	} else {
		positions, _ = state.db.GetOpenPositions()
	}
	var totalHeat float64
	bucketHeat := make(map[string]float64)

	for _, pos := range positions {
		totalHeat += pos.RiskDollars
		bucketHeat[pos.Bucket] += pos.RiskDollars
	}

	portfolioHeatCap := equity * (portfolioCap / 100.0)
	portfolioHeatPct := (totalHeat / portfolioHeatCap) * 100.0

	// Portfolio heat display
	portfolioLabel := widget.NewLabelWithStyle("Portfolio Heat", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	portfolioHeatLabel := widget.NewLabel(fmt.Sprintf("$%.2f / $%.2f (%.1f%%)", totalHeat, portfolioHeatCap, portfolioHeatPct))

	portfolioProgress := widget.NewProgressBar()
	portfolioProgress.SetValue(totalHeat / portfolioHeatCap)

	// Bucket heat display
	bucketLabel := widget.NewLabelWithStyle("Bucket Heat", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	bucketsList := container.NewVBox()
	bucketHeatCap := equity * (bucketCap / 100.0)

	if len(bucketHeat) == 0 {
		bucketsList.Add(widget.NewLabel("No open positions"))
	} else {
		for bucket, heat := range bucketHeat {
			heatPct := (heat / bucketHeatCap) * 100.0
			bucketRow := container.NewVBox(
				widget.NewLabel(fmt.Sprintf("%s: $%.2f / $%.2f (%.1f%%)", bucket, heat, bucketHeatCap, heatPct)),
			)
			progress := widget.NewProgressBar()
			progress.SetValue(heat / bucketHeatCap)
			bucketRow.Add(progress)
			bucketsList.Add(bucketRow)
		}
	}

	// New trade test
	testLabel := widget.NewLabelWithStyle("Test New Trade", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	riskLabel := widget.NewLabel("Risk Amount ($):")
	riskEntry := widget.NewEntry()
	// Auto-fill from session
	riskEntry.SetText(fmt.Sprintf("%.2f", activeSession.SizingRiskDollars))
	riskEntry.SetPlaceHolder("750.00")

	bucketSelectLabel := widget.NewLabel("Sector Bucket:")
	bucketEntry := widget.NewEntry()
	bucketEntry.SetPlaceHolder("Tech/Comm")
	if state.sampleMode {
		bucketEntry.SetText("Tech/Comm") // Sample data
	}

	testResultLabel := widget.NewLabel("")
	testResultLabel.Wrapping = fyne.TextWrapWord

	// Next button (shown after successful heat check)
	nextBtn := widget.NewButton("Next: Trade Entry →", func() {
		message := "Please use the tab bar to navigate to the Trade Entry tab.\n\n"
		if state.sampleMode {
			message += "📦 Sample mode - explore Trade Entry with sample data"
		} else {
			message += "Your heat check has been saved to Session #" + fmt.Sprintf("%d", activeSession.SessionNum)
		}
		ShowStyledInformation("Next Step", message, state.window)
	})
	nextBtn.Importance = widget.HighImportance
	nextBtn.Hide() // Hidden initially

	testBtn := widget.NewButton("Check Heat", func() {
		riskStr := riskEntry.Text
		if riskStr == "" {
			testResultLabel.SetText("❌ Please enter risk amount")
			return
		}

		risk, err := strconv.ParseFloat(riskStr, 64)
		if err != nil {
			testResultLabel.SetText("❌ Invalid risk amount")
			return
		}

		bucket := bucketEntry.Text
		if bucket == "" {
			bucket = "Unknown"
		}

		// Call backend heat check - prepare request (use sample or real positions)
		var checkPositions []storage.Position
		if state.sampleMode {
			checkPositions = CreateSamplePositions()
		} else {
			checkPositions, _ = state.db.GetOpenPositions()
		}
		openPositions := make([]domain.Position, len(checkPositions))
		for i, p := range checkPositions {
			openPositions[i] = domain.Position{
				Ticker:      p.Ticker,
				Bucket:      p.Bucket,
				RiskDollars: p.RiskDollars,
				UnitsOpen:   p.Shares,
				Status:      "Open",
			}
		}

		heatReq := domain.HeatRequest{
			Equity:           equity,
			HeatCapPct:       portfolioCap / 100.0,
			BucketHeatCapPct: bucketCap / 100.0,
			AddRiskDollars:   risk,
			AddBucket:        bucket,
			OpenPositions:    openPositions,
		}

		result, err := domain.CalculateHeat(heatReq)
		if err != nil {
			testResultLabel.SetText(fmt.Sprintf("❌ Error: %v", err))
			return
		}

		// Determine status
		heatStatus := "OK"
		if !result.Allowed {
			heatStatus = "REJECT"
		}

		// Update session in database (skip in sample mode)
		if !state.sampleMode {
			err = state.db.UpdateSessionHeat(
				activeSession.ID,
				heatStatus,
				bucket,
				result.CurrentPortfolioHeat,
				result.NewPortfolioHeat,
				result.PortfolioCap,
				result.CurrentBucketHeat,
				result.NewBucketHeat,
				result.BucketCap,
			)
			if err != nil {
				testResultLabel.SetText(fmt.Sprintf("❌ Failed to save session: %v", err))
				return
			}

			// Reload session to get updated state
			updatedSession, err := state.db.GetSession(activeSession.ID)
			if err != nil {
				testResultLabel.SetText(fmt.Sprintf("❌ Failed to reload session: %v", err))
				return
			}
			state.SetCurrentSession(updatedSession)
		}

		if result.Allowed {
			message := fmt.Sprintf(`✓ TRADE ALLOWED

New Portfolio Heat: $%.2f / $%.2f (%.1f%%)
New Bucket Heat (%s): $%.2f / $%.2f (%.1f%%)

Both caps OK!

`, result.NewPortfolioHeat, result.PortfolioCap, result.PortfolioHeatPct,
				bucket, result.NewBucketHeat, result.BucketCap, result.BucketHeatPct)

			if state.sampleMode {
				message += "📦 SAMPLE MODE - Heat check shown but not saved"
			} else {
				message += fmt.Sprintf("✓ Session #%d updated - ready for Trade Entry", activeSession.SessionNum)
			}

			testResultLabel.SetText(message)
			nextBtn.Show()
		} else {
			message := fmt.Sprintf(`✗ TRADE REJECTED

Reason: %s

Current Portfolio Heat: $%.2f
New Portfolio Heat: $%.2f
Portfolio Cap: $%.2f
Overage: $%.2f

Current Bucket Heat (%s): $%.2f
New Bucket Heat: $%.2f
Bucket Cap: $%.2f
Overage: $%.2f

`, result.RejectionReason,
				result.CurrentPortfolioHeat, result.NewPortfolioHeat, result.PortfolioCap,
				result.PortfolioOverage,
				bucket, result.CurrentBucketHeat, result.NewBucketHeat, result.BucketCap,
				result.BucketOverage)

			if state.sampleMode {
				message += "📦 SAMPLE MODE - Results shown but not saved"
			} else {
				message += fmt.Sprintf("Session #%d updated - resolve heat issues before proceeding", activeSession.SessionNum)
			}

			testResultLabel.SetText(message)
			nextBtn.Hide()
		}
	})
	testBtn.Importance = widget.HighImportance

	// Disable heat check button if session is completed (but allow in sample mode)
	if activeSession.Status == "COMPLETED" && !state.sampleMode {
		testBtn.Disable()
	}

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(sessionInfo),
		explanationCard,
		widget.NewSeparator(),
		settingsInfo,
		widget.NewSeparator(),
		portfolioLabel,
		portfolioHeatLabel,
		portfolioProgress,
		widget.NewSeparator(),
		bucketLabel,
		bucketsList,
		widget.NewSeparator(),
		testLabel,
		riskLabel,
		riskEntry,
		bucketSelectLabel,
		bucketEntry,
		container.NewHBox(testBtn, nextBtn),
		widget.NewSeparator(),
		testResultLabel,
	)

	return container.NewScroll(content)
}

// buildHeatCheckExplanation creates the educational panel explaining Heat Check
func buildHeatCheckExplanation() fyne.CanvasObject {
	explanationText := widget.NewRichTextFromMarkdown(`### What is Heat Check?

**Heat = Total risk across all open positions**

#### The Rules (from anti-impulsivity system):

1. **Portfolio Heat Cap:** 4% of equity
   - Example: $10,000 account → $400 max total risk

2. **Bucket Heat Cap:** 1.5% of equity per sector
   - Example: $10,000 account → $150 max in Tech/Comm

3. **Purpose:** Prevent concentration risk
   - Forces diversification across sectors
   - Prevents "all-in" on one trade/sector
   - Mechanical enforcement (no override)

#### How It Works:

- Sum risk from all open positions
- Add risk from proposed new trade
- If total > cap → **REJECT** (no exceptions)
- If OK → **ALLOW** trade

**This is discipline enforcement, not flexibility.**`)

	explanationCard := container.NewVBox(
		widget.NewLabelWithStyle("ℹ️ Understanding Heat Check", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		explanationText,
	)

	return container.NewPadded(explanationCard)
}
