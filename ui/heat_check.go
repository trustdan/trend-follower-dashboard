package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/domain"
)

func buildHeatCheckScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("Heat Check - Portfolio Risk Management", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

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

	// Calculate current heat
	positions, _ := state.db.GetOpenPositions()
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
	riskEntry.SetPlaceHolder("750.00")

	bucketSelectLabel := widget.NewLabel("Sector Bucket:")
	bucketEntry := widget.NewEntry()
	bucketEntry.SetPlaceHolder("Tech/Comm")

	testResultLabel := widget.NewLabel("")
	testResultLabel.Wrapping = fyne.TextWrapWord

	testBtn := widget.NewButton("Check If Trade Allowed", func() {
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

		// Call backend heat check - prepare request
		positions, _ := state.db.GetOpenPositions()
		openPositions := make([]domain.Position, len(positions))
		for i, p := range positions {
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

		if result.Allowed {
			testResultLabel.SetText(fmt.Sprintf(`✓ TRADE ALLOWED

New Portfolio Heat: $%.2f / $%.2f (%.1f%%)
New Bucket Heat (%s): $%.2f / $%.2f (%.1f%%)

Both caps OK!
`, result.NewPortfolioHeat, result.PortfolioCap, result.PortfolioHeatPct,
				bucket, result.NewBucketHeat, result.BucketCap, result.BucketHeatPct))
		} else {
			testResultLabel.SetText(fmt.Sprintf(`✗ TRADE REJECTED

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
				result.BucketOverage))
		}
	})
	testBtn.Importance = widget.HighImportance

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
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
		testBtn,
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
