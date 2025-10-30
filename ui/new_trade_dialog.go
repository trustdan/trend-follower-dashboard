package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// showNewTradeDialog displays the enhanced "Start New Trade" dialog with options support
func showNewTradeDialog(state *AppState, navigateToTab func(int)) {
	// Ticker entry
	tickerEntry := widget.NewEntry()
	tickerEntry.SetPlaceHolder("AAPL")

	// Breakout system selection
	var selectedSystem string = storage.SystemTwo
	breakoutSystemRadio := widget.NewRadioGroup([]string{
		"System-2 (55/10) - 55-bar breakout, 10-bar exit [DEFAULT]",
		"System-1 (20/10) - 20-bar breakout, 10-bar exit",
		"Custom - Manual parameters",
	}, func(value string) {
		switch {
		case value == "System-2 (55/10) - 55-bar breakout, 10-bar exit [DEFAULT]":
			selectedSystem = storage.SystemTwo
		case value == "System-1 (20/10) - 20-bar breakout, 10-bar exit":
			selectedSystem = storage.SystemOne
		default:
			selectedSystem = storage.SystemCustom
		}
	})
	breakoutSystemRadio.SetSelected("System-2 (55/10) - 55-bar breakout, 10-bar exit [DEFAULT]")
	breakoutSystemRadio.Horizontal = false

	// Direction selection
	var selectedDirection string = storage.StrategyLongBreakout
	directionRadio := widget.NewRadioGroup([]string{
		"Long (bullish)",
		"Short (bearish)",
	}, func(value string) {
		if value == "Long (bullish)" {
			selectedDirection = storage.StrategyLongBreakout
		} else {
			selectedDirection = storage.StrategyShortBreakout
		}
	})
	directionRadio.SetSelected("Long (bullish)")
	directionRadio.Horizontal = false

	// Instrument type selection
	var selectedInstrumentType string = storage.InstrumentStock
	var selectedOptionsStrategy string = ""

	instrumentTypeRadio := widget.NewRadioGroup([]string{
		"Stock/ETF (no options)",
		"Options (select strategy below)",
	}, func(value string) {
		if value == "Stock/ETF (no options)" {
			selectedInstrumentType = storage.InstrumentStock
			selectedOptionsStrategy = ""
		} else {
			selectedInstrumentType = storage.InstrumentOption
		}
	})
	instrumentTypeRadio.SetSelected("Stock/ETF (no options)")
	instrumentTypeRadio.Horizontal = false

	// Options strategy dropdown (categorized)
	optionsStrategies := []string{
		"--- Directional ---",
		"Long Call",
		"Long Put",
		"--- Income ---",
		"Covered Call",
		"Cash-Secured Put",
		"--- Vertical Credit Spreads ---",
		"Bull Put Credit Spread",
		"Bear Call Credit Spread",
		"--- Butterflies & Condors ---",
		"Iron Butterfly",
		"Iron Condor",
		"Long Put Butterfly",
		"Long Call Butterfly",
		"Inverse Iron Butterfly",
		"Inverse Iron Condor",
		"Short Put Butterfly",
		"Short Call Butterfly",
		"--- Time Spreads ---",
		"Calendar Call Spread",
		"Calendar Put Spread",
		"Diagonal Call Spread",
		"Diagonal Put Spread",
		"--- Vertical Debit Spreads ---",
		"Bull Call Spread",
		"Bear Put Spread",
		"--- Volatility ---",
		"Straddle",
		"Strangle",
		"--- Ratio & Broken Wing ---",
		"Call Ratio Backspread",
		"Put Broken Wing",
		"Put Ratio Backspread",
		"Call Broken Wing",
	}

	optionsStrategySelect := widget.NewSelect(optionsStrategies, func(value string) {
		// Map display name to constant
		selectedOptionsStrategy = mapDisplayNameToConstant(value)
	})
	optionsStrategySelect.PlaceHolder = "Select an options strategy..."
	optionsStrategySelect.Disable() // Start disabled

	// Enable/disable options dropdown based on instrument type
	instrumentTypeRadio.OnChanged = func(value string) {
		if value == "Options (select strategy below)" {
			optionsStrategySelect.Enable()
			selectedInstrumentType = storage.InstrumentOption
		} else {
			optionsStrategySelect.Disable()
			optionsStrategySelect.ClearSelected()
			selectedInstrumentType = storage.InstrumentStock
			selectedOptionsStrategy = ""
		}
	}

	// Info text
	infoText := widget.NewLabel(
		"Select your breakout system, direction, and instrument type.\n" +
			"For options trades, you'll configure legs in the next step.\n" +
			"The session will track your progress through all 5 gates.",
	)
	infoText.Wrapping = fyne.TextWrapWord

	// Create form
	formContent := container.NewVBox(
		widget.NewLabel("Ticker:"),
		tickerEntry,
		widget.NewSeparator(),

		widget.NewLabel("Breakout System:"),
		breakoutSystemRadio,
		widget.NewSeparator(),

		widget.NewLabel("Direction:"),
		directionRadio,
		widget.NewSeparator(),

		widget.NewLabel("Instrument Type:"),
		instrumentTypeRadio,
		widget.NewSeparator(),

		widget.NewLabel("Options Strategy:"),
		optionsStrategySelect,
		widget.NewSeparator(),

		infoText,
	)

	// Create custom dialog
	d := dialog.NewCustomConfirm(
		"Start New Trade Session",
		"Next",
		"Cancel",
		container.NewVScroll(formContent),
		func(proceed bool) {
			if !proceed {
				return
			}

			// Validate ticker
			ticker := tickerEntry.Text
			if ticker == "" {
				dialog.ShowError(fmt.Errorf("ticker is required"), state.window)
				return
			}

			// If options selected, validate strategy selection
			if selectedInstrumentType == storage.InstrumentOption && selectedOptionsStrategy == "" {
				dialog.ShowError(fmt.Errorf("please select an options strategy"), state.window)
				return
			}

			// If options selected, show strategy builder
			if selectedInstrumentType == storage.InstrumentOption {
				ShowStrategyBuilder(selectedOptionsStrategy, state.window, func(result *StrategyBuilderResult) {
					// Create session with options metadata
					createOptionsSession(state, ticker, selectedDirection, selectedSystem, result, navigateToTab)
				})
			} else {
				// Create stock/ETF session (simple)
				createStockSession(state, ticker, selectedDirection, selectedSystem, navigateToTab)
			}
		},
		state.window,
	)

	d.Resize(fyne.NewSize(600, 800))
	d.Show()
}

// createStockSession creates a simple stock/ETF session
func createStockSession(state *AppState, ticker, direction, system string, navigateToTab func(int)) {
	log.Printf("Creating stock session: ticker=%s, direction=%s, system=%s", ticker, direction, system)

	session, err := state.db.CreateSession(ticker, direction)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		dialog.ShowError(err, state.window)
		return
	}

	log.Printf("Session created successfully: ID=%d, Num=%d", session.ID, session.SessionNum)

	// Set as current session
	state.SetCurrentSession(session)

	// Navigate to Checklist tab (index 2)
	if navigateToTab != nil {
		navigateToTab(2) // Dashboard=0, Scanner=1, Checklist=2
	}

	// Show success message
	dialog.ShowInformation(
		"Session Created",
		fmt.Sprintf("Session #%d created for %s (%s).\n\nNavigated to Checklist tab to begin evaluation.",
			session.SessionNum,
			ticker,
			formatStrategy(direction),
		),
		state.window,
	)
}

// createOptionsSession creates an options session with full metadata
func createOptionsSession(state *AppState, ticker, direction, system string, result *StrategyBuilderResult, navigateToTab func(int)) {
	log.Printf("Creating options session: ticker=%s, strategy=%s", ticker, result.OptionsStrategy)

	// Serialize legs to JSON
	legsJSON, err := SerializeLegs(result.Legs)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to serialize legs: %w", err), state.window)
		return
	}

	// Determine entry/exit lookback based on system
	entryLookback := 55
	exitLookback := 10
	if system == storage.SystemOne {
		entryLookback = 20
	}

	// Create session with options metadata
	session, err := state.db.CreateSessionWithOptions(
		ticker,
		direction,
		storage.InstrumentOption,
		result.OptionsStrategy,
		result.EntryDate,
		result.PrimaryExpirationDate,
		result.DTE,
		result.RollThresholdDTE,
		result.TimeExitMode,
		legsJSON,
		result.NetDebit,
		result.MaxProfit,
		result.MaxLoss,
		result.BreakevenLower,
		result.BreakevenUpper,
		result.UnderlyingAtEntry,
		4,   // max_units default
		0.5, // add_step_n default
		entryLookback,
		exitLookback,
	)

	if err != nil {
		log.Printf("Error creating options session: %v", err)
		dialog.ShowError(err, state.window)
		return
	}

	log.Printf("Options session created successfully: ID=%d, Num=%d", session.ID, session.SessionNum)

	// Set as current session
	state.SetCurrentSession(session)

	// Navigate to Checklist tab (index 2)
	if navigateToTab != nil {
		navigateToTab(2)
	}

	// Show detailed success message
	dialog.ShowInformation(
		"Options Session Created",
		fmt.Sprintf(
			"Session #%d created for %s (%s).\n\n"+
				"Strategy: %s\n"+
				"Expiration: %s (%d DTE)\n"+
				"Net Debit: $%.2f\n"+
				"Max Profit: $%.2f\n"+
				"Max Loss: $%.2f\n\n"+
				"Navigated to Checklist tab to begin evaluation.",
			session.SessionNum,
			ticker,
			storage.GetStrategyDisplayName(result.OptionsStrategy),
			storage.GetStrategyDisplayName(result.OptionsStrategy),
			result.PrimaryExpirationDate,
			result.DTE,
			result.NetDebit,
			result.MaxProfit,
			result.MaxLoss,
		),
		state.window,
	)
}

// mapDisplayNameToConstant maps display name to storage constant
func mapDisplayNameToConstant(displayName string) string {
	mapping := map[string]string{
		"Long Call":                  storage.StrategyLongCall,
		"Long Put":                   storage.StrategyLongPut,
		"Covered Call":               storage.StrategyCoveredCall,
		"Cash-Secured Put":           storage.StrategyCashSecuredPut,
		"Bull Put Credit Spread":     storage.StrategyBullPutSpread,
		"Bear Call Credit Spread":    storage.StrategyBearCallSpread,
		"Iron Butterfly":             storage.StrategyIronButterfly,
		"Iron Condor":                storage.StrategyIronCondor,
		"Long Put Butterfly":         storage.StrategyLongPutButterfly,
		"Long Call Butterfly":        storage.StrategyLongCallButterfly,
		"Inverse Iron Butterfly":     storage.StrategyInverseIronButterfly,
		"Inverse Iron Condor":        storage.StrategyInverseIronCondor,
		"Short Put Butterfly":        storage.StrategyShortPutButterfly,
		"Short Call Butterfly":       storage.StrategyShortCallButterfly,
		"Calendar Call Spread":       storage.StrategyCalendarCallSpread,
		"Calendar Put Spread":        storage.StrategyCalendarPutSpread,
		"Diagonal Call Spread":       storage.StrategyDiagonalCallSpread,
		"Diagonal Put Spread":        storage.StrategyDiagonalPutSpread,
		"Bull Call Spread":           storage.StrategyBullCallSpread,
		"Bear Put Spread":            storage.StrategyBearPutSpread,
		"Straddle":                   storage.StrategyStraddle,
		"Strangle":                   storage.StrategyStrangle,
		"Call Ratio Backspread":      storage.StrategyCallRatioBackspread,
		"Put Broken Wing":            storage.StrategyPutBrokenWing,
		"Put Ratio Backspread":       storage.StrategyPutRatioBackspread,
		"Call Broken Wing":           storage.StrategyCallBrokenWing,
	}

	if constant, ok := mapping[displayName]; ok {
		return constant
	}

	return ""
}
