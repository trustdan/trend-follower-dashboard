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

func buildPositionSizingScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("Position Sizing Calculator", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Method selector
	methodLabel := widget.NewLabel("Sizing Method:")
	methodSelect := widget.NewSelect([]string{"Stock/ETF", "Options (Delta-ATR)", "Options (Contracts)"}, nil)
	methodSelect.SetSelected("Stock/ETF")

	// Method explanation (will update based on selection)
	methodExplanation := widget.NewRichTextFromMarkdown(getMethodExplanation("Stock/ETF"))

	// Common inputs
	tickerLabel := widget.NewLabel("Ticker:")
	tickerEntry := widget.NewEntry()
	tickerEntry.SetPlaceHolder("AAPL")

	entryLabel := widget.NewLabel("Entry Price:")
	entryEntry := widget.NewEntry()
	entryEntry.SetPlaceHolder("180.00")

	atrLabel := widget.NewLabel("ATR (N):")
	atrEntry := widget.NewEntry()
	atrEntry.SetPlaceHolder("1.50")

	kLabel := widget.NewLabel("K Multiple (Stop Distance):")
	kEntry := widget.NewEntry()
	kEntry.SetText("2.0")

	// Options-specific inputs (hidden by default)
	deltaLabel := widget.NewLabel("Delta:")
	deltaEntry := widget.NewEntry()
	deltaEntry.SetPlaceHolder("0.70")

	contractPriceLabel := widget.NewLabel("Contract Price:")
	contractPriceEntry := widget.NewEntry()
	contractPriceEntry.SetPlaceHolder("5.00")

	optionsInputs := container.NewVBox(
		deltaLabel,
		deltaEntry,
		contractPriceLabel,
		contractPriceEntry,
	)
	optionsInputs.Hide()

	// Show/hide options inputs based on method
	methodSelect.OnChanged = func(method string) {
		// Update explanation
		methodExplanation.ParseMarkdown(getMethodExplanation(method))
		methodExplanation.Refresh()

		// Show/hide options inputs
		if method == "Stock/ETF" {
			optionsInputs.Hide()
		} else {
			optionsInputs.Show()
		}
	}

	// Load account settings for equity and risk % with defaults
	settings, _ := state.db.GetAllSettings()
	equityStr := getSettingWithDefault(settings, "equity", "100000")
	riskPctStr := getSettingWithDefault(settings, "risk_pct", "0.75")

	settingsLabel := widget.NewLabel(fmt.Sprintf("Account: $%s equity, %s%% risk per trade", equityStr, riskPctStr))

	// Results display
	resultsLabel := widget.NewLabel("")
	resultsLabel.Wrapping = fyne.TextWrapWord

	// Calculate button
	calculateBtn := widget.NewButton("Calculate Position Size", func() {
		ticker := tickerEntry.Text
		if ticker == "" {
			resultsLabel.SetText("❌ Please enter a ticker")
			return
		}

		// Parse inputs
		entry, err := strconv.ParseFloat(entryEntry.Text, 64)
		if err != nil {
			resultsLabel.SetText("❌ Invalid entry price")
			return
		}

		atr, err := strconv.ParseFloat(atrEntry.Text, 64)
		if err != nil {
			resultsLabel.SetText("❌ Invalid ATR")
			return
		}

		k, err := strconv.ParseFloat(kEntry.Text, 64)
		if err != nil {
			resultsLabel.SetText("❌ Invalid K multiple")
			return
		}

		equity, err := strconv.ParseFloat(equityStr, 64)
		if err != nil {
			resultsLabel.SetText("❌ Invalid equity setting")
			return
		}

		riskPct, err := strconv.ParseFloat(riskPctStr, 64)
		if err != nil {
			resultsLabel.SetText("❌ Invalid risk % setting")
			return
		}
		riskPctDecimal := riskPct / 100.0

		// Prepare sizing request
		method := methodSelect.Selected

		sizingReq := domain.SizingRequest{
			Equity:  equity,
			RiskPct: riskPctDecimal,
			Entry:   entry,
			ATR:     atr,
			K:       int(k),
		}

		// Set method and optional params
		switch method {
		case "Stock/ETF":
			sizingReq.Method = "stock"

		case "Options (Delta-ATR)":
			delta, err := strconv.ParseFloat(deltaEntry.Text, 64)
			if err != nil {
				resultsLabel.SetText("❌ Invalid delta")
				return
			}
			sizingReq.Method = "opt-delta-atr"
			sizingReq.Delta = delta

		case "Options (Contracts)":
			contractPrice, err := strconv.ParseFloat(contractPriceEntry.Text, 64)
			if err != nil {
				resultsLabel.SetText("❌ Invalid contract price")
				return
			}
			sizingReq.Method = "opt-maxloss"
			sizingReq.MaxLoss = contractPrice * 100 // Contract price × 100 shares
		}

		// Call backend sizing function
		result, calcErr := domain.CalculatePositionSize(sizingReq)

		if calcErr != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Error: %v", calcErr))
			return
		}

		// Format results
		resultsText := fmt.Sprintf(`✓ Position Size Calculated

Method: %s
Ticker: %s
Risk Dollars: $%.2f
Stop Distance: $%.2f
Initial Stop: $%.2f

`,
			result.Method, ticker, result.RiskDollars, result.StopDistance, result.InitialStop)

		if result.Shares > 0 {
			resultsText += fmt.Sprintf("SHARES TO BUY: %d shares\n\n", result.Shares)
		}

		if result.Contracts > 0 {
			resultsText += fmt.Sprintf("CONTRACTS TO BUY: %d contracts\n\n", result.Contracts)
		}

		resultsText += fmt.Sprintf("Actual Risk: $%.2f (%.2f%% of equity)\n",
			result.ActualRisk, result.ActualRisk/equity*100)

		resultsLabel.SetText(resultsText)
	})
	calculateBtn.Importance = widget.HighImportance

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		settingsLabel,
		widget.NewSeparator(),
		methodLabel,
		methodSelect,
		container.NewPadded(methodExplanation),
		widget.NewSeparator(),
		tickerLabel,
		tickerEntry,
		entryLabel,
		entryEntry,
		atrLabel,
		atrEntry,
		kLabel,
		kEntry,
		optionsInputs,
		widget.NewSeparator(),
		calculateBtn,
		widget.NewSeparator(),
		resultsLabel,
	)

	return container.NewScroll(content)
}

// getMethodExplanation returns the explanation text for each position sizing method
func getMethodExplanation(method string) string {
	switch method {
	case "Stock/ETF":
		return `**Van Tharp Position Sizing Method:**

1. Risk $ = Account × Risk%
2. Stop Distance = K × ATR (typically K=2)
3. Initial Stop = Entry - Stop Distance
4. Shares = floor(Risk$ / Stop Distance)
5. Actual Risk = Shares × Stop Distance

**Example:** $10,000 account, 0.75% risk, AAPL @ $180, ATR=$1.50, K=2
→ Risk=$75, Stop=$3, Shares=25

**Pyramiding:** Add every 0.5×N up to Max Units (typically 4)`

	case "Options (Delta-ATR)":
		return `**Options Sizing Using Delta-Adjusted ATR:**

Per-unit risk adjusted for delta:
**Contracts = Risk$ / (Delta × Stop Distance × 100)**

**Use for:** Single calls/puts
**Best when:** Clear directional bias, high IV
**Delta:** Typically 0.60-0.80 for trending moves

This method accounts for the leverage and directional exposure of options.`

	case "Options (Contracts)":
		return `**Options Sizing Using Max Loss Per Contract:**

**Contracts = Risk$ / (Contract Price × 100)**

**Use for:** Debit spreads, defined-risk strategies
**Best when:** Lower cost, limited downside
**Max Loss:** Debit paid per spread

This method is for strategies where the maximum loss is known upfront (contract price paid).`

	default:
		return "Select a method to see its explanation."
	}
}
