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
	// Session check: require active session
	if state.currentSession == nil {
		return showNoSessionPrompt(state, "Position Sizing")
	}

	// Prerequisite check: checklist must be completed with GREEN banner
	if !state.currentSession.ChecklistCompleted {
		return showPrerequisiteError(state, "Checklist", "Position Sizing")
	}

	if state.currentSession.ChecklistBanner == "RED" || state.currentSession.ChecklistBanner == "YELLOW" {
		return showPrerequisiteError(state, "Checklist", "Position Sizing")
	}

	// Title
	title := canvas.NewText("Position Sizing Calculator", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Session info
	sessionInfo := widget.NewLabel(fmt.Sprintf(
		"Session #%d • %s • %s • Banner: %s",
		state.currentSession.SessionNum,
		formatStrategyDisplay(state.currentSession.Strategy),
		state.currentSession.Ticker,
		state.currentSession.ChecklistBanner,
	))
	sessionInfo.TextStyle = fyne.TextStyle{Bold: true}

	// Method selector
	methodLabel := widget.NewLabel("Sizing Method:")
	methodSelect := widget.NewSelect([]string{"Stock/ETF", "Options (Delta-ATR)", "Options (Contracts)"}, nil)
	methodSelect.SetSelected("Stock/ETF")

	// Method explanation (will update based on selection)
	methodExplanation := widget.NewRichTextFromMarkdown(getMethodExplanation("Stock/ETF"))

	// Common inputs
	tickerLabel := widget.NewLabel("Ticker:")
	tickerEntry := widget.NewEntry()
	tickerEntry.SetText(state.currentSession.Ticker) // Auto-fill from session
	tickerEntry.SetPlaceHolder("AAPL")
	tickerEntry.Disable() // Lock ticker to session ticker

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

	// Pyramid planning inputs
	maxUnitsLabel := widget.NewLabel("Max Units (Pyramid):")
	maxUnitsEntry := widget.NewEntry()
	maxUnitsEntry.SetText("4")
	maxUnitsEntry.SetPlaceHolder("4")

	addStepLabel := widget.NewLabel("Add Every (× N):")
	addStepEntry := widget.NewEntry()
	addStepEntry.SetText("0.5")
	addStepEntry.SetPlaceHolder("0.5")

	currentUnitsLabel := widget.NewLabel("Current Units:")
	currentUnitsEntry := widget.NewEntry()
	currentUnitsEntry.SetText("0")
	currentUnitsEntry.SetPlaceHolder("0")

	// Add-on prices display (calculated after sizing)
	addOnPricesLabel := widget.NewLabel("")
	addOnPricesLabel.Wrapping = fyne.TextWrapWord

	pyramidSection := container.NewVBox(
		widget.NewLabelWithStyle("=== Pyramid Planning ===", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		maxUnitsLabel,
		maxUnitsEntry,
		addStepLabel,
		addStepEntry,
		currentUnitsLabel,
		currentUnitsEntry,
		widget.NewSeparator(),
		addOnPricesLabel,
	)

	// Options info display (if session has options)
	optionsInfoLabel := widget.NewLabel("")
	optionsInfoLabel.Wrapping = fyne.TextWrapWord
	optionsInfoLabel.Hide()

	if state.currentSession.InstrumentType == "option" && state.currentSession.OptionsStrategy != "" {
		// Show options information
		optionsInfo := fmt.Sprintf("📊 Options Strategy: %s\n"+
			"Expiration: %s (%d DTE)\n"+
			"Entry Date: %s\n"+
			"Roll at: %d DTE\n"+
			"Time Exit Mode: %s",
			state.currentSession.OptionsStrategy,
			state.currentSession.PrimaryExpirationDate,
			state.currentSession.DTE,
			state.currentSession.EntryDate,
			state.currentSession.RollThresholdDTE,
			state.currentSession.TimeExitMode)

		optionsInfoLabel.SetText(optionsInfo)
		optionsInfoLabel.Show()
	}

	// Next button (shown after successful calculation)
	nextBtn := widget.NewButton("Next: Heat Check →", func() {
		ShowStyledInformation("Next Step",
			"Please use the tab bar to navigate to the Heat Check tab.\n\n"+
				"Your position sizing has been saved to Session #"+fmt.Sprintf("%d", state.currentSession.SessionNum),
			state.window)
	})
	nextBtn.Importance = widget.HighImportance
	nextBtn.Hide() // Hidden initially

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

		// Parse pyramid inputs with validation
		maxUnits, err := strconv.Atoi(maxUnitsEntry.Text)
		if err != nil || maxUnits < 1 || maxUnits > 10 {
			resultsLabel.SetText("❌ Max units must be between 1 and 10")
			return
		}

		addStep, err := strconv.ParseFloat(addStepEntry.Text, 64)
		if err != nil || addStep <= 0 {
			resultsLabel.SetText("❌ Add step must be greater than 0")
			return
		}

		currentUnits, err := strconv.Atoi(currentUnitsEntry.Text)
		if err != nil || currentUnits < 0 || currentUnits > maxUnits {
			resultsLabel.SetText(fmt.Sprintf("❌ Current units must be between 0 and %d", maxUnits))
			return
		}

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

		// Calculate add-on prices (up to 3 additional units)
		addPrice1 := entry + (addStep * atr)
		addPrice2 := entry + (addStep * 2.0 * atr)
		addPrice3 := entry + (addStep * 3.0 * atr)

		addOnText := fmt.Sprintf("📊 Add-On Prices (Every %.1f × N):\n", addStep)
		if maxUnits > 1 {
			addOnText += fmt.Sprintf("  Add 1: $%.2f (Entry + %.1fN)\n", addPrice1, addStep)
		}
		if maxUnits > 2 {
			addOnText += fmt.Sprintf("  Add 2: $%.2f (Entry + %.1fN)\n", addPrice2, addStep*2.0)
		}
		if maxUnits > 3 {
			addOnText += fmt.Sprintf("  Add 3: $%.2f (Entry + %.1fN)\n", addPrice3, addStep*3.0)
		}
		addOnText += fmt.Sprintf("\nCurrent Units: %d / %d", currentUnits, maxUnits)
		addOnPricesLabel.SetText(addOnText)

		// Update session in database with pyramid data
		methodStr := result.Method
		deltaValue := 0.0
		if method == "Options (Delta-ATR)" && sizingReq.Delta > 0 {
			deltaValue = sizingReq.Delta
		}
		err = state.db.UpdateSessionSizingWithPyramid(
			state.currentSession.ID,
			methodStr,
			entry,
			atr,
			k,
			result.StopDistance,
			result.InitialStop,
			result.Shares,
			result.Contracts,
			result.RiskDollars,
			deltaValue,
			maxUnits,
			addStep,
			addPrice1,
			addPrice2,
			addPrice3,
		)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to save session: %v", err))
			return
		}

		// Reload session to get updated state
		updatedSession, err := state.db.GetSession(state.currentSession.ID)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to reload session: %v", err))
			return
		}
		state.SetCurrentSession(updatedSession)

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

		resultsText += fmt.Sprintf("\n✓ Session #%d updated - ready for Heat Check", state.currentSession.SessionNum)

		resultsLabel.SetText(resultsText)
		nextBtn.Show()
	})
	calculateBtn.Importance = widget.HighImportance

	// Disable calculate button if session is completed
	if state.currentSession.Status == "COMPLETED" {
		calculateBtn.Disable()
	}

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(sessionInfo),
		settingsLabel,
		widget.NewSeparator(),
		optionsInfoLabel, // Show options info if applicable
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
		pyramidSection, // Pyramid planning section
		widget.NewSeparator(),
		container.NewHBox(calculateBtn, nextBtn),
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
