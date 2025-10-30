package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// StrategyBuilderResult contains the output of a strategy builder
type StrategyBuilderResult struct {
	OptionsStrategy       string
	Legs                  []storage.OptionLeg
	EntryDate             string
	PrimaryExpirationDate string
	DTE                   int
	NetDebit              float64
	MaxProfit             float64
	MaxLoss               float64
	BreakevenLower        float64
	BreakevenUpper        float64
	UnderlyingAtEntry     float64
	RollThresholdDTE      int
	TimeExitMode          string
}

// ShowStrategyBuilder displays the appropriate builder dialog for the given strategy
func ShowStrategyBuilder(strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	switch strategyType {
	case storage.StrategyLongCall:
		showLongCallBuilder(parentWindow, onComplete)
	case storage.StrategyLongPut:
		showLongPutBuilder(parentWindow, onComplete)
	case storage.StrategyBullCallSpread, storage.StrategyBearPutSpread:
		showVerticalSpreadBuilder(strategyType, parentWindow, onComplete)
	case storage.StrategyBullPutSpread, storage.StrategyBearCallSpread:
		showCreditSpreadBuilder(strategyType, parentWindow, onComplete)
	case storage.StrategyIronCondor:
		showIronCondorBuilder(parentWindow, onComplete)
	case storage.StrategyIronButterfly:
		showIronButterflyBuilder(parentWindow, onComplete)
	case storage.StrategyStraddle:
		showStraddleBuilder(parentWindow, onComplete)
	case storage.StrategyStrangle:
		showStrangleBuilder(parentWindow, onComplete)
	case storage.StrategyCalendarCallSpread, storage.StrategyCalendarPutSpread:
		showCalendarSpreadBuilder(strategyType, parentWindow, onComplete)
	case storage.StrategyDiagonalCallSpread, storage.StrategyDiagonalPutSpread:
		showDiagonalSpreadBuilder(strategyType, parentWindow, onComplete)
	case storage.StrategyLongCallButterfly, storage.StrategyLongPutButterfly:
		showButterflyBuilder(strategyType, parentWindow, onComplete)
	default:
		dialog.ShowInformation("Not Implemented",
			fmt.Sprintf("Builder for %s not yet implemented.\nPlease use manual leg entry.",
				storage.GetStrategyDisplayName(strategyType)),
			parentWindow)
	}
}

// =============================================================================
// Directional: Long Call / Long Put
// =============================================================================

func showLongCallBuilder(parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	showSingleLegBuilder("CALL", storage.StrategyLongCall, parentWindow, onComplete)
}

func showLongPutBuilder(parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	showSingleLegBuilder("PUT", storage.StrategyLongPut, parentWindow, onComplete)
}

func showSingleLegBuilder(optionType, strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	// Input fields
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	strikeEntry := widget.NewEntry()
	strikeEntry.SetPlaceHolder("180.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetPlaceHolder("5")
	contractsEntry.SetText("1")

	premiumEntry := widget.NewEntry()
	premiumEntry.SetPlaceHolder("2.50")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetPlaceHolder("21")
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	// Preview labels
	netDebitLabel := widget.NewLabel("Net Debit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: Unlimited")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakeven: $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	// Update preview function
	updatePreview := func() {
		underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
		strike, _ := strconv.ParseFloat(strikeEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		premium, _ := strconv.ParseFloat(premiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		// Build legs
		var legs []storage.OptionLeg
		if optionType == "CALL" {
			legs = storage.BuildLongCall(strike, expiration, contracts, premium)
		} else {
			legs = storage.BuildLongPut(strike, expiration, contracts, premium)
		}

		// Calculate metrics
		netDebit := storage.CalculateNetDebit(legs)
		_, maxLoss, _ := storage.CalculateMaxProfitLoss(strategyType, legs)
		lowerBE, upperBE, _ := storage.CalculateBreakevens(strategyType, legs)

		// Calculate DTE
		dte := calculateDTE(expiration)

		// Update labels
		netDebitLabel.SetText(fmt.Sprintf("Net Debit: $%.2f", netDebit))
		maxProfitLabel.SetText("Max Profit: Unlimited")
		maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))

		if optionType == "CALL" {
			breakevenLabel.SetText(fmt.Sprintf("Breakeven: $%.2f", upperBE))
		} else {
			breakevenLabel.SetText(fmt.Sprintf("Breakeven: $%.2f", lowerBE))
		}

		dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))

		// Update entry date to today if not set
		if underlying > 0 && strike > 0 && premium > 0 {
			// Valid inputs
		}
	}

	// Attach listeners
	underlyingEntry.OnChanged = func(string) { updatePreview() }
	strikeEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	premiumEntry.OnChanged = func(string) { updatePreview() }

	// Form content
	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("=== %s Builder ===", storage.GetStrategyDisplayName(strategyType))),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel(fmt.Sprintf("%s Strike:", optionType)),
		strikeEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Premium per Contract:"),
		premiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	// Create dialog
	d := dialog.NewCustomConfirm(
		fmt.Sprintf("Build %s", storage.GetStrategyDisplayName(strategyType)),
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			// Parse inputs
			underlying, err := strconv.ParseFloat(underlyingEntry.Text, 64)
			if err != nil || underlying <= 0 {
				dialog.ShowError(fmt.Errorf("invalid underlying price"), parentWindow)
				return
			}

			strike, err := strconv.ParseFloat(strikeEntry.Text, 64)
			if err != nil || strike <= 0 {
				dialog.ShowError(fmt.Errorf("invalid strike price"), parentWindow)
				return
			}

			expiration := expirationEntry.Text
			if expiration == "" {
				dialog.ShowError(fmt.Errorf("expiration date required"), parentWindow)
				return
			}

			contracts, err := strconv.Atoi(contractsEntry.Text)
			if err != nil || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid contracts"), parentWindow)
				return
			}

			premium, err := strconv.ParseFloat(premiumEntry.Text, 64)
			if err != nil || premium <= 0 {
				dialog.ShowError(fmt.Errorf("invalid premium"), parentWindow)
				return
			}

			rollThreshold, err := strconv.Atoi(rollThresholdEntry.Text)
			if err != nil || rollThreshold <= 0 {
				rollThreshold = 21
			}

			// Build legs
			var legs []storage.OptionLeg
			if optionType == "CALL" {
				legs = storage.BuildLongCall(strike, expiration, contracts, premium)
			} else {
				legs = storage.BuildLongPut(strike, expiration, contracts, premium)
			}

			// Calculate metrics
			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(strategyType, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(strategyType, legs)
			dte := calculateDTE(expiration)

			// Create result
			result := &StrategyBuilderResult{
				OptionsStrategy:       strategyType,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 700))
	d.Show()
}

// =============================================================================
// Vertical Spreads: Bull Call, Bear Put (Debit Spreads)
// =============================================================================

func showVerticalSpreadBuilder(strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	isBullCall := strategyType == storage.StrategyBullCallSpread
	optionType := "CALL"
	if !isBullCall {
		optionType = "PUT"
	}

	// Input fields
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	lowerStrikeEntry := widget.NewEntry()
	lowerStrikeEntry.SetPlaceHolder("170.00")

	upperStrikeEntry := widget.NewEntry()
	upperStrikeEntry.SetPlaceHolder("180.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetPlaceHolder("5")
	contractsEntry.SetText("1")

	lowerPremiumEntry := widget.NewEntry()
	lowerPremiumEntry.SetPlaceHolder("8.50")

	upperPremiumEntry := widget.NewEntry()
	upperPremiumEntry.SetPlaceHolder("3.50")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	// Preview labels
	netDebitLabel := widget.NewLabel("Net Debit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: $0.00")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakeven: $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	// Update preview
	updatePreview := func() {
		lowerStrike, _ := strconv.ParseFloat(lowerStrikeEntry.Text, 64)
		upperStrike, _ := strconv.ParseFloat(upperStrikeEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		lowerPremium, _ := strconv.ParseFloat(lowerPremiumEntry.Text, 64)
		upperPremium, _ := strconv.ParseFloat(upperPremiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		// Build legs
		var legs []storage.OptionLeg
		if isBullCall {
			legs = storage.BuildBullCallSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
		} else {
			legs = storage.BuildBearPutSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
		}

		// Calculate metrics
		netDebit := storage.CalculateNetDebit(legs)
		maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(strategyType, legs)
		lowerBE, upperBE, _ := storage.CalculateBreakevens(strategyType, legs)
		dte := calculateDTE(expiration)

		// Update labels
		netDebitLabel.SetText(fmt.Sprintf("Net Debit: $%.2f", netDebit))
		maxProfitLabel.SetText(fmt.Sprintf("Max Profit: $%.2f", maxProfit))
		maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))

		if isBullCall {
			breakevenLabel.SetText(fmt.Sprintf("Breakeven: $%.2f", upperBE))
		} else {
			breakevenLabel.SetText(fmt.Sprintf("Breakeven: $%.2f", lowerBE))
		}

		dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))
	}

	// Attach listeners
	lowerStrikeEntry.OnChanged = func(string) { updatePreview() }
	upperStrikeEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	lowerPremiumEntry.OnChanged = func(string) { updatePreview() }
	upperPremiumEntry.OnChanged = func(string) { updatePreview() }

	// Form content
	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("=== %s Builder ===", storage.GetStrategyDisplayName(strategyType))),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel(fmt.Sprintf("Lower %s Strike:", optionType)),
		lowerStrikeEntry,

		widget.NewLabel(fmt.Sprintf("Upper %s Strike:", optionType)),
		upperStrikeEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Lower Strike Premium:"),
		lowerPremiumEntry,

		widget.NewLabel("Upper Strike Premium:"),
		upperPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	// Create dialog
	d := dialog.NewCustomConfirm(
		fmt.Sprintf("Build %s", storage.GetStrategyDisplayName(strategyType)),
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			// Parse and validate inputs
			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			lowerStrike, _ := strconv.ParseFloat(lowerStrikeEntry.Text, 64)
			upperStrike, _ := strconv.ParseFloat(upperStrikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			lowerPremium, _ := strconv.ParseFloat(lowerPremiumEntry.Text, 64)
			upperPremium, _ := strconv.ParseFloat(upperPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || lowerStrike <= 0 || upperStrike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			// Build legs
			var legs []storage.OptionLeg
			if isBullCall {
				legs = storage.BuildBullCallSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
			} else {
				legs = storage.BuildBearPutSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
			}

			// Calculate metrics
			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(strategyType, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(strategyType, legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       strategyType,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 800))
	d.Show()
}

// =============================================================================
// Credit Spreads: Bull Put, Bear Call
// =============================================================================

func showCreditSpreadBuilder(strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	isBullPut := strategyType == storage.StrategyBullPutSpread
	optionType := "PUT"
	if !isBullPut {
		optionType = "CALL"
	}

	// Input fields (similar to debit spreads)
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	lowerStrikeEntry := widget.NewEntry()
	lowerStrikeEntry.SetPlaceHolder("165.00")

	upperStrikeEntry := widget.NewEntry()
	upperStrikeEntry.SetPlaceHolder("170.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	lowerPremiumEntry := widget.NewEntry()
	lowerPremiumEntry.SetPlaceHolder("0.75")

	upperPremiumEntry := widget.NewEntry()
	upperPremiumEntry.SetPlaceHolder("1.50")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	// Preview labels
	netCreditLabel := widget.NewLabel("Net Credit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: $0.00")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakeven: $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	updatePreview := func() {
		lowerStrike, _ := strconv.ParseFloat(lowerStrikeEntry.Text, 64)
		upperStrike, _ := strconv.ParseFloat(upperStrikeEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		lowerPremium, _ := strconv.ParseFloat(lowerPremiumEntry.Text, 64)
		upperPremium, _ := strconv.ParseFloat(upperPremiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		var legs []storage.OptionLeg
		if isBullPut {
			legs = storage.BuildBullPutSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
		} else {
			legs = storage.BuildBearCallSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
		}

		netDebit := storage.CalculateNetDebit(legs)
		maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(strategyType, legs)
		lowerBE, upperBE, _ := storage.CalculateBreakevens(strategyType, legs)
		dte := calculateDTE(expiration)

		netCreditLabel.SetText(fmt.Sprintf("Net Credit: $%.2f", -netDebit))
		maxProfitLabel.SetText(fmt.Sprintf("Max Profit: $%.2f", maxProfit))
		maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))

		if isBullPut {
			breakevenLabel.SetText(fmt.Sprintf("Breakeven: $%.2f", lowerBE))
		} else {
			breakevenLabel.SetText(fmt.Sprintf("Breakeven: $%.2f", upperBE))
		}

		dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))
	}

	lowerStrikeEntry.OnChanged = func(string) { updatePreview() }
	upperStrikeEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	lowerPremiumEntry.OnChanged = func(string) { updatePreview() }
	upperPremiumEntry.OnChanged = func(string) { updatePreview() }

	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("=== %s Builder ===", storage.GetStrategyDisplayName(strategyType))),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel(fmt.Sprintf("Lower %s Strike (BUY):", optionType)),
		lowerStrikeEntry,

		widget.NewLabel(fmt.Sprintf("Upper %s Strike (SELL):", optionType)),
		upperStrikeEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Lower Strike Premium:"),
		lowerPremiumEntry,

		widget.NewLabel("Upper Strike Premium:"),
		upperPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netCreditLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	d := dialog.NewCustomConfirm(
		fmt.Sprintf("Build %s", storage.GetStrategyDisplayName(strategyType)),
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			lowerStrike, _ := strconv.ParseFloat(lowerStrikeEntry.Text, 64)
			upperStrike, _ := strconv.ParseFloat(upperStrikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			lowerPremium, _ := strconv.ParseFloat(lowerPremiumEntry.Text, 64)
			upperPremium, _ := strconv.ParseFloat(upperPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || lowerStrike <= 0 || upperStrike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			var legs []storage.OptionLeg
			if isBullPut {
				legs = storage.BuildBullPutSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
			} else {
				legs = storage.BuildBearCallSpread(lowerStrike, upperStrike, expiration, contracts, lowerPremium, upperPremium)
			}

			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(strategyType, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(strategyType, legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       strategyType,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 800))
	d.Show()
}

// =============================================================================
// Iron Condor Builder
// =============================================================================

func showIronCondorBuilder(parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	// Input fields
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	putSpreadWidthEntry := widget.NewEntry()
	putSpreadWidthEntry.SetPlaceHolder("5.00")

	callSpreadWidthEntry := widget.NewEntry()
	callSpreadWidthEntry.SetPlaceHolder("5.00")

	wingDistanceEntry := widget.NewEntry()
	wingDistanceEntry.SetPlaceHolder("10.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	// Premium entries for all 4 legs
	buyPutPremiumEntry := widget.NewEntry()
	buyPutPremiumEntry.SetPlaceHolder("0.50")

	sellPutPremiumEntry := widget.NewEntry()
	sellPutPremiumEntry.SetPlaceHolder("1.20")

	sellCallPremiumEntry := widget.NewEntry()
	sellCallPremiumEntry.SetPlaceHolder("1.30")

	buyCallPremiumEntry := widget.NewEntry()
	buyCallPremiumEntry.SetPlaceHolder("0.60")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	// Preview labels
	strikePreviewLabel := widget.NewLabel("Strikes: (calculate)")
	netCreditLabel := widget.NewLabel("Net Credit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: $0.00")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakevens: $0.00 / $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	updatePreview := func() {
		underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
		putSpreadWidth, _ := strconv.ParseFloat(putSpreadWidthEntry.Text, 64)
		callSpreadWidth, _ := strconv.ParseFloat(callSpreadWidthEntry.Text, 64)
		wingDistance, _ := strconv.ParseFloat(wingDistanceEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		buyPutPremium, _ := strconv.ParseFloat(buyPutPremiumEntry.Text, 64)
		sellPutPremium, _ := strconv.ParseFloat(sellPutPremiumEntry.Text, 64)
		sellCallPremium, _ := strconv.ParseFloat(sellCallPremiumEntry.Text, 64)
		buyCallPremium, _ := strconv.ParseFloat(buyCallPremiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		if underlying > 0 && putSpreadWidth > 0 && callSpreadWidth > 0 && wingDistance > 0 {
			// Calculate strikes
			sellPutStrike := underlying - wingDistance
			buyPutStrike := sellPutStrike - putSpreadWidth
			sellCallStrike := underlying + wingDistance
			buyCallStrike := sellCallStrike + callSpreadWidth

			strikePreviewLabel.SetText(fmt.Sprintf("Strikes: P %.2f/%.2f | C %.2f/%.2f",
				buyPutStrike, sellPutStrike, sellCallStrike, buyCallStrike))

			// Build legs
			legs := storage.BuildIronCondor(underlying, putSpreadWidth, callSpreadWidth, wingDistance,
				expiration, contracts, buyPutPremium, sellPutPremium, sellCallPremium, buyCallPremium)

			// Calculate metrics
			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyIronCondor, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyIronCondor, legs)
			dte := calculateDTE(expiration)

			netCreditLabel.SetText(fmt.Sprintf("Net Credit: $%.2f", -netDebit))
			maxProfitLabel.SetText(fmt.Sprintf("Max Profit: $%.2f", maxProfit))
			maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))
			breakevenLabel.SetText(fmt.Sprintf("Breakevens: $%.2f / $%.2f", lowerBE, upperBE))
			dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))
		}
	}

	underlyingEntry.OnChanged = func(string) { updatePreview() }
	putSpreadWidthEntry.OnChanged = func(string) { updatePreview() }
	callSpreadWidthEntry.OnChanged = func(string) { updatePreview() }
	wingDistanceEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	buyPutPremiumEntry.OnChanged = func(string) { updatePreview() }
	sellPutPremiumEntry.OnChanged = func(string) { updatePreview() }
	sellCallPremiumEntry.OnChanged = func(string) { updatePreview() }
	buyCallPremiumEntry.OnChanged = func(string) { updatePreview() }

	form := container.NewVBox(
		widget.NewLabel("=== Iron Condor Builder ==="),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel("Put Spread Width:"),
		putSpreadWidthEntry,

		widget.NewLabel("Call Spread Width:"),
		callSpreadWidthEntry,

		widget.NewLabel("Wing Distance from Underlying:"),
		wingDistanceEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewSeparator(),
		widget.NewLabel("Premiums (per contract):"),

		widget.NewLabel("Buy Put Premium:"),
		buyPutPremiumEntry,

		widget.NewLabel("Sell Put Premium:"),
		sellPutPremiumEntry,

		widget.NewLabel("Sell Call Premium:"),
		sellCallPremiumEntry,

		widget.NewLabel("Buy Call Premium:"),
		buyCallPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		strikePreviewLabel,
		netCreditLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	d := dialog.NewCustomConfirm(
		"Build Iron Condor",
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			putSpreadWidth, _ := strconv.ParseFloat(putSpreadWidthEntry.Text, 64)
			callSpreadWidth, _ := strconv.ParseFloat(callSpreadWidthEntry.Text, 64)
			wingDistance, _ := strconv.ParseFloat(wingDistanceEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			buyPutPremium, _ := strconv.ParseFloat(buyPutPremiumEntry.Text, 64)
			sellPutPremium, _ := strconv.ParseFloat(sellPutPremiumEntry.Text, 64)
			sellCallPremium, _ := strconv.ParseFloat(sellCallPremiumEntry.Text, 64)
			buyCallPremium, _ := strconv.ParseFloat(buyCallPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || putSpreadWidth <= 0 || callSpreadWidth <= 0 || wingDistance <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildIronCondor(underlying, putSpreadWidth, callSpreadWidth, wingDistance,
				expiration, contracts, buyPutPremium, sellPutPremium, sellCallPremium, buyCallPremium)

			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyIronCondor, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyIronCondor, legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       storage.StrategyIronCondor,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(550, 900))
	d.Show()
}

// =============================================================================
// Iron Butterfly Builder
// =============================================================================

func showIronButterflyBuilder(parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	atmStrikeEntry := widget.NewEntry()
	atmStrikeEntry.SetPlaceHolder("175.00")

	wingWidthEntry := widget.NewEntry()
	wingWidthEntry.SetPlaceHolder("5.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	buyPutPremiumEntry := widget.NewEntry()
	buyPutPremiumEntry.SetPlaceHolder("0.40")

	sellPutPremiumEntry := widget.NewEntry()
	sellPutPremiumEntry.SetPlaceHolder("2.50")

	sellCallPremiumEntry := widget.NewEntry()
	sellCallPremiumEntry.SetPlaceHolder("2.50")

	buyCallPremiumEntry := widget.NewEntry()
	buyCallPremiumEntry.SetPlaceHolder("0.40")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	strikePreviewLabel := widget.NewLabel("Strikes: (calculate)")
	netCreditLabel := widget.NewLabel("Net Credit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: $0.00")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakevens: $0.00 / $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	updatePreview := func() {
		atmStrike, _ := strconv.ParseFloat(atmStrikeEntry.Text, 64)
		wingWidth, _ := strconv.ParseFloat(wingWidthEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		buyPutPremium, _ := strconv.ParseFloat(buyPutPremiumEntry.Text, 64)
		sellPutPremium, _ := strconv.ParseFloat(sellPutPremiumEntry.Text, 64)
		sellCallPremium, _ := strconv.ParseFloat(sellCallPremiumEntry.Text, 64)
		buyCallPremium, _ := strconv.ParseFloat(buyCallPremiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		if atmStrike > 0 && wingWidth > 0 {
			lowerWing := atmStrike - wingWidth
			upperWing := atmStrike + wingWidth

			strikePreviewLabel.SetText(fmt.Sprintf("Strikes: P %.2f/%.2f | C %.2f/%.2f",
				lowerWing, atmStrike, atmStrike, upperWing))

			legs := storage.BuildIronButterfly(atmStrike, wingWidth, expiration, contracts,
				buyPutPremium, sellPutPremium, sellCallPremium, buyCallPremium)

			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyIronButterfly, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyIronButterfly, legs)
			dte := calculateDTE(expiration)

			netCreditLabel.SetText(fmt.Sprintf("Net Credit: $%.2f", -netDebit))
			maxProfitLabel.SetText(fmt.Sprintf("Max Profit: $%.2f", maxProfit))
			maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))
			breakevenLabel.SetText(fmt.Sprintf("Breakevens: $%.2f / $%.2f", lowerBE, upperBE))
			dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))
		}
	}

	atmStrikeEntry.OnChanged = func(string) { updatePreview() }
	wingWidthEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	buyPutPremiumEntry.OnChanged = func(string) { updatePreview() }
	sellPutPremiumEntry.OnChanged = func(string) { updatePreview() }
	sellCallPremiumEntry.OnChanged = func(string) { updatePreview() }
	buyCallPremiumEntry.OnChanged = func(string) { updatePreview() }

	form := container.NewVBox(
		widget.NewLabel("=== Iron Butterfly Builder ==="),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel("ATM Strike:"),
		atmStrikeEntry,

		widget.NewLabel("Wing Width:"),
		wingWidthEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewSeparator(),
		widget.NewLabel("Premiums:"),

		widget.NewLabel("Buy Put Premium:"),
		buyPutPremiumEntry,

		widget.NewLabel("Sell Put Premium:"),
		sellPutPremiumEntry,

		widget.NewLabel("Sell Call Premium:"),
		sellCallPremiumEntry,

		widget.NewLabel("Buy Call Premium:"),
		buyCallPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		strikePreviewLabel,
		netCreditLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	d := dialog.NewCustomConfirm(
		"Build Iron Butterfly",
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			atmStrike, _ := strconv.ParseFloat(atmStrikeEntry.Text, 64)
			wingWidth, _ := strconv.ParseFloat(wingWidthEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			buyPutPremium, _ := strconv.ParseFloat(buyPutPremiumEntry.Text, 64)
			sellPutPremium, _ := strconv.ParseFloat(sellPutPremiumEntry.Text, 64)
			sellCallPremium, _ := strconv.ParseFloat(sellCallPremiumEntry.Text, 64)
			buyCallPremium, _ := strconv.ParseFloat(buyCallPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || atmStrike <= 0 || wingWidth <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildIronButterfly(atmStrike, wingWidth, expiration, contracts,
				buyPutPremium, sellPutPremium, sellCallPremium, buyCallPremium)

			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyIronButterfly, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyIronButterfly, legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       storage.StrategyIronButterfly,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(550, 900))
	d.Show()
}

// =============================================================================
// Straddle Builder
// =============================================================================

func showStraddleBuilder(parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	atmStrikeEntry := widget.NewEntry()
	atmStrikeEntry.SetPlaceHolder("175.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	callPremiumEntry := widget.NewEntry()
	callPremiumEntry.SetPlaceHolder("8.50")

	putPremiumEntry := widget.NewEntry()
	putPremiumEntry.SetPlaceHolder("8.00")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	netDebitLabel := widget.NewLabel("Net Debit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: Unlimited")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakevens: $0.00 / $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	updatePreview := func() {
		atmStrike, _ := strconv.ParseFloat(atmStrikeEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		callPremium, _ := strconv.ParseFloat(callPremiumEntry.Text, 64)
		putPremium, _ := strconv.ParseFloat(putPremiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		if atmStrike > 0 {
			legs := storage.BuildStraddle(atmStrike, expiration, contracts, callPremium, putPremium)

			netDebit := storage.CalculateNetDebit(legs)
			_, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyStraddle, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyStraddle, legs)
			dte := calculateDTE(expiration)

			netDebitLabel.SetText(fmt.Sprintf("Net Debit: $%.2f", netDebit))
			maxProfitLabel.SetText("Max Profit: Unlimited")
			maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))
			breakevenLabel.SetText(fmt.Sprintf("Breakevens: $%.2f / $%.2f", lowerBE, upperBE))
			dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))
		}
	}

	atmStrikeEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	callPremiumEntry.OnChanged = func(string) { updatePreview() }
	putPremiumEntry.OnChanged = func(string) { updatePreview() }

	form := container.NewVBox(
		widget.NewLabel("=== Straddle Builder ==="),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel("ATM Strike:"),
		atmStrikeEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Call Premium:"),
		callPremiumEntry,

		widget.NewLabel("Put Premium:"),
		putPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	d := dialog.NewCustomConfirm(
		"Build Straddle",
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			atmStrike, _ := strconv.ParseFloat(atmStrikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			callPremium, _ := strconv.ParseFloat(callPremiumEntry.Text, 64)
			putPremium, _ := strconv.ParseFloat(putPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || atmStrike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildStraddle(atmStrike, expiration, contracts, callPremium, putPremium)

			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyStraddle, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyStraddle, legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       storage.StrategyStraddle,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 700))
	d.Show()
}

// =============================================================================
// Strangle Builder
// =============================================================================

func showStrangleBuilder(parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	callStrikeEntry := widget.NewEntry()
	callStrikeEntry.SetPlaceHolder("185.00")

	putStrikeEntry := widget.NewEntry()
	putStrikeEntry.SetPlaceHolder("165.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	callPremiumEntry := widget.NewEntry()
	callPremiumEntry.SetPlaceHolder("4.50")

	putPremiumEntry := widget.NewEntry()
	putPremiumEntry.SetPlaceHolder("4.00")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	netDebitLabel := widget.NewLabel("Net Debit: $0.00")
	maxProfitLabel := widget.NewLabel("Max Profit: Unlimited")
	maxLossLabel := widget.NewLabel("Max Loss: $0.00")
	breakevenLabel := widget.NewLabel("Breakevens: $0.00 / $0.00")
	dteLabel := widget.NewLabel("DTE: 0")

	updatePreview := func() {
		callStrike, _ := strconv.ParseFloat(callStrikeEntry.Text, 64)
		putStrike, _ := strconv.ParseFloat(putStrikeEntry.Text, 64)
		contracts, _ := strconv.Atoi(contractsEntry.Text)
		callPremium, _ := strconv.ParseFloat(callPremiumEntry.Text, 64)
		putPremium, _ := strconv.ParseFloat(putPremiumEntry.Text, 64)
		expiration := expirationEntry.Text

		if contracts == 0 {
			contracts = 1
		}

		if callStrike > 0 && putStrike > 0 {
			legs := storage.BuildStrangle(callStrike, putStrike, expiration, contracts, callPremium, putPremium)

			netDebit := storage.CalculateNetDebit(legs)
			_, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyStrangle, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyStrangle, legs)
			dte := calculateDTE(expiration)

			netDebitLabel.SetText(fmt.Sprintf("Net Debit: $%.2f", netDebit))
			maxProfitLabel.SetText("Max Profit: Unlimited")
			maxLossLabel.SetText(fmt.Sprintf("Max Loss: $%.2f", maxLoss))
			breakevenLabel.SetText(fmt.Sprintf("Breakevens: $%.2f / $%.2f", lowerBE, upperBE))
			dteLabel.SetText(fmt.Sprintf("DTE: %d", dte))
		}
	}

	callStrikeEntry.OnChanged = func(string) { updatePreview() }
	putStrikeEntry.OnChanged = func(string) { updatePreview() }
	expirationEntry.OnChanged = func(string) { updatePreview() }
	contractsEntry.OnChanged = func(string) { updatePreview() }
	callPremiumEntry.OnChanged = func(string) { updatePreview() }
	putPremiumEntry.OnChanged = func(string) { updatePreview() }

	form := container.NewVBox(
		widget.NewLabel("=== Strangle Builder ==="),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel("Call Strike (OTM):"),
		callStrikeEntry,

		widget.NewLabel("Put Strike (OTM):"),
		putStrikeEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Call Premium:"),
		callPremiumEntry,

		widget.NewLabel("Put Premium:"),
		putPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
		maxProfitLabel,
		maxLossLabel,
		breakevenLabel,
		dteLabel,
	)

	d := dialog.NewCustomConfirm(
		"Build Strangle",
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			callStrike, _ := strconv.ParseFloat(callStrikeEntry.Text, 64)
			putStrike, _ := strconv.ParseFloat(putStrikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			callPremium, _ := strconv.ParseFloat(callPremiumEntry.Text, 64)
			putPremium, _ := strconv.ParseFloat(putPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || callStrike <= 0 || putStrike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildStrangle(callStrike, putStrike, expiration, contracts, callPremium, putPremium)

			netDebit := storage.CalculateNetDebit(legs)
			maxProfit, maxLoss, _ := storage.CalculateMaxProfitLoss(storage.StrategyStrangle, legs)
			lowerBE, upperBE, _ := storage.CalculateBreakevens(storage.StrategyStrangle, legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       storage.StrategyStrangle,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             maxProfit,
				MaxLoss:               maxLoss,
				BreakevenLower:        lowerBE,
				BreakevenUpper:        upperBE,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 750))
	d.Show()
}

// =============================================================================
// Calendar Spread Builder
// =============================================================================

func showCalendarSpreadBuilder(strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	isCall := strategyType == storage.StrategyCalendarCallSpread
	optionType := "CALL"
	if !isCall {
		optionType = "PUT"
	}

	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	strikeEntry := widget.NewEntry()
	strikeEntry.SetPlaceHolder("175.00")

	nearExpirationEntry := widget.NewEntry()
	nearExpirationEntry.SetPlaceHolder("2025-11-15")

	farExpirationEntry := widget.NewEntry()
	farExpirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	nearPremiumEntry := widget.NewEntry()
	nearPremiumEntry.SetPlaceHolder("3.50")

	farPremiumEntry := widget.NewEntry()
	farPremiumEntry.SetPlaceHolder("5.00")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	netDebitLabel := widget.NewLabel("Net Debit: $0.00")
	infoLabel := widget.NewLabel("Calendar spreads profit from time decay.")

	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("=== %s Builder ===", storage.GetStrategyDisplayName(strategyType))),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel(fmt.Sprintf("%s Strike (same for both):", optionType)),
		strikeEntry,

		widget.NewLabel("Near Expiration (YYYY-MM-DD):"),
		nearExpirationEntry,

		widget.NewLabel("Far Expiration (YYYY-MM-DD):"),
		farExpirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Near Expiration Premium:"),
		nearPremiumEntry,

		widget.NewLabel("Far Expiration Premium:"),
		farPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
		infoLabel,
	)

	d := dialog.NewCustomConfirm(
		fmt.Sprintf("Build %s", storage.GetStrategyDisplayName(strategyType)),
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			strike, _ := strconv.ParseFloat(strikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			nearPremium, _ := strconv.ParseFloat(nearPremiumEntry.Text, 64)
			farPremium, _ := strconv.ParseFloat(farPremiumEntry.Text, 64)
			nearExpiration := nearExpirationEntry.Text
			farExpiration := farExpirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || strike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildCalendarSpread(optionType, strike, nearExpiration, farExpiration, contracts, nearPremium, farPremium)

			netDebit := storage.CalculateNetDebit(legs)
			dte := calculateDTE(nearExpiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       strategyType,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: nearExpiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             0, // Complex calculation
				MaxLoss:               netDebit,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 700))
	d.Show()
}

// =============================================================================
// Diagonal Spread Builder
// =============================================================================

func showDiagonalSpreadBuilder(strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	isCall := strategyType == storage.StrategyDiagonalCallSpread
	optionType := "CALL"
	if !isCall {
		optionType = "PUT"
	}

	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	nearStrikeEntry := widget.NewEntry()
	nearStrikeEntry.SetPlaceHolder("180.00")

	farStrikeEntry := widget.NewEntry()
	farStrikeEntry.SetPlaceHolder("175.00")

	nearExpirationEntry := widget.NewEntry()
	nearExpirationEntry.SetPlaceHolder("2025-11-15")

	farExpirationEntry := widget.NewEntry()
	farExpirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	nearPremiumEntry := widget.NewEntry()
	nearPremiumEntry.SetPlaceHolder("3.00")

	farPremiumEntry := widget.NewEntry()
	farPremiumEntry.SetPlaceHolder("5.50")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	netDebitLabel := widget.NewLabel("Net Debit: $0.00")

	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("=== %s Builder ===", storage.GetStrategyDisplayName(strategyType))),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel(fmt.Sprintf("Near %s Strike:", optionType)),
		nearStrikeEntry,

		widget.NewLabel(fmt.Sprintf("Far %s Strike:", optionType)),
		farStrikeEntry,

		widget.NewLabel("Near Expiration (YYYY-MM-DD):"),
		nearExpirationEntry,

		widget.NewLabel("Far Expiration (YYYY-MM-DD):"),
		farExpirationEntry,

		widget.NewLabel("Contracts:"),
		contractsEntry,

		widget.NewLabel("Near Expiration Premium:"),
		nearPremiumEntry,

		widget.NewLabel("Far Expiration Premium:"),
		farPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
	)

	d := dialog.NewCustomConfirm(
		fmt.Sprintf("Build %s", storage.GetStrategyDisplayName(strategyType)),
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			nearStrike, _ := strconv.ParseFloat(nearStrikeEntry.Text, 64)
			farStrike, _ := strconv.ParseFloat(farStrikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			nearPremium, _ := strconv.ParseFloat(nearPremiumEntry.Text, 64)
			farPremium, _ := strconv.ParseFloat(farPremiumEntry.Text, 64)
			nearExpiration := nearExpirationEntry.Text
			farExpiration := farExpirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || nearStrike <= 0 || farStrike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildDiagonalSpread(optionType, nearStrike, farStrike, nearExpiration, farExpiration, contracts, nearPremium, farPremium)

			netDebit := storage.CalculateNetDebit(legs)
			dte := calculateDTE(nearExpiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       strategyType,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: nearExpiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             0, // Complex
				MaxLoss:               netDebit,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 750))
	d.Show()
}

// =============================================================================
// Butterfly Builder
// =============================================================================

func showButterflyBuilder(strategyType string, parentWindow fyne.Window, onComplete func(*StrategyBuilderResult)) {
	isCall := strategyType == storage.StrategyLongCallButterfly
	optionType := "CALL"
	if !isCall {
		optionType = "PUT"
	}

	underlyingEntry := widget.NewEntry()
	underlyingEntry.SetPlaceHolder("175.00")

	lowerStrikeEntry := widget.NewEntry()
	lowerStrikeEntry.SetPlaceHolder("170.00")

	middleStrikeEntry := widget.NewEntry()
	middleStrikeEntry.SetPlaceHolder("175.00")

	upperStrikeEntry := widget.NewEntry()
	upperStrikeEntry.SetPlaceHolder("180.00")

	expirationEntry := widget.NewEntry()
	expirationEntry.SetPlaceHolder("2025-12-19")

	contractsEntry := widget.NewEntry()
	contractsEntry.SetText("1")

	lowerPremiumEntry := widget.NewEntry()
	lowerPremiumEntry.SetPlaceHolder("8.00")

	middlePremiumEntry := widget.NewEntry()
	middlePremiumEntry.SetPlaceHolder("5.00")

	upperPremiumEntry := widget.NewEntry()
	upperPremiumEntry.SetPlaceHolder("3.00")

	rollThresholdEntry := widget.NewEntry()
	rollThresholdEntry.SetText("21")

	timeExitModeSelect := widget.NewSelect([]string{"Close", "Roll", "None"}, nil)
	timeExitModeSelect.SetSelected("Close")

	netDebitLabel := widget.NewLabel("Net Debit: $0.00")
	infoLabel := widget.NewLabel("1-2-1 ratio: Buy 1 low, Sell 2 middle, Buy 1 high")

	form := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("=== %s Builder ===", storage.GetStrategyDisplayName(strategyType))),
		widget.NewSeparator(),

		widget.NewLabel("Underlying Price:"),
		underlyingEntry,

		widget.NewLabel(fmt.Sprintf("Lower %s Strike:", optionType)),
		lowerStrikeEntry,

		widget.NewLabel(fmt.Sprintf("Middle %s Strike:", optionType)),
		middleStrikeEntry,

		widget.NewLabel(fmt.Sprintf("Upper %s Strike:", optionType)),
		upperStrikeEntry,

		widget.NewLabel("Expiration Date (YYYY-MM-DD):"),
		expirationEntry,

		widget.NewLabel("Contracts (per wing):"),
		contractsEntry,

		widget.NewLabel("Lower Strike Premium:"),
		lowerPremiumEntry,

		widget.NewLabel("Middle Strike Premium:"),
		middlePremiumEntry,

		widget.NewLabel("Upper Strike Premium:"),
		upperPremiumEntry,

		widget.NewSeparator(),
		widget.NewLabel("Exit Settings:"),

		widget.NewLabel("Roll at DTE:"),
		rollThresholdEntry,

		widget.NewLabel("Time Exit Mode:"),
		timeExitModeSelect,

		widget.NewSeparator(),
		widget.NewLabel("=== Preview ==="),
		netDebitLabel,
		infoLabel,
	)

	d := dialog.NewCustomConfirm(
		fmt.Sprintf("Build %s", storage.GetStrategyDisplayName(strategyType)),
		"Generate",
		"Cancel",
		container.NewVScroll(form),
		func(generate bool) {
			if !generate {
				return
			}

			underlying, _ := strconv.ParseFloat(underlyingEntry.Text, 64)
			lowerStrike, _ := strconv.ParseFloat(lowerStrikeEntry.Text, 64)
			middleStrike, _ := strconv.ParseFloat(middleStrikeEntry.Text, 64)
			upperStrike, _ := strconv.ParseFloat(upperStrikeEntry.Text, 64)
			contracts, _ := strconv.Atoi(contractsEntry.Text)
			lowerPremium, _ := strconv.ParseFloat(lowerPremiumEntry.Text, 64)
			middlePremium, _ := strconv.ParseFloat(middlePremiumEntry.Text, 64)
			upperPremium, _ := strconv.ParseFloat(upperPremiumEntry.Text, 64)
			expiration := expirationEntry.Text
			rollThreshold, _ := strconv.Atoi(rollThresholdEntry.Text)

			if underlying <= 0 || lowerStrike <= 0 || middleStrike <= 0 || upperStrike <= 0 || contracts <= 0 {
				dialog.ShowError(fmt.Errorf("invalid inputs"), parentWindow)
				return
			}

			legs := storage.BuildButterfly(optionType, lowerStrike, middleStrike, upperStrike, expiration, contracts,
				lowerPremium, middlePremium, upperPremium)

			netDebit := storage.CalculateNetDebit(legs)
			dte := calculateDTE(expiration)

			result := &StrategyBuilderResult{
				OptionsStrategy:       strategyType,
				Legs:                  legs,
				EntryDate:             time.Now().Format("2006-01-02"),
				PrimaryExpirationDate: expiration,
				DTE:                   dte,
				NetDebit:              netDebit,
				MaxProfit:             0, // Complex
				MaxLoss:               netDebit,
				UnderlyingAtEntry:     underlying,
				RollThresholdDTE:      rollThreshold,
				TimeExitMode:          timeExitModeSelect.Selected,
			}

			onComplete(result)
		},
		parentWindow,
	)

	d.Resize(fyne.NewSize(500, 800))
	d.Show()
}

// =============================================================================
// Helper Functions
// =============================================================================

// calculateDTE calculates days to expiration from an expiration date string (YYYY-MM-DD)
func calculateDTE(expirationDate string) int {
	if expirationDate == "" {
		return 0
	}

	expDate, err := time.Parse("2006-01-02", expirationDate)
	if err != nil {
		return 0
	}

	now := time.Now()
	duration := expDate.Sub(now)
	days := int(duration.Hours() / 24)

	if days < 0 {
		return 0
	}

	return days
}

// FormatLegsForDisplay formats legs array for display in UI
func FormatLegsForDisplay(legs []storage.OptionLeg) string {
	if len(legs) == 0 {
		return "No legs defined"
	}

	result := ""
	for i, leg := range legs {
		result += fmt.Sprintf("Leg %d: %s %s $%.2f × %d @ $%.2f\n",
			i+1, leg.Action, leg.Type, leg.Strike, leg.Qty, leg.Price)
	}
	return result
}

// SerializeLegs converts legs array to JSON string for database storage
func SerializeLegs(legs []storage.OptionLeg) (string, error) {
	if len(legs) == 0 {
		return "", nil
	}

	bytes, err := json.Marshal(legs)
	if err != nil {
		return "", fmt.Errorf("failed to serialize legs: %w", err)
	}

	return string(bytes), nil
}

// DeserializeLegs converts JSON string from database to legs array
func DeserializeLegs(legsJSON string) ([]storage.OptionLeg, error) {
	if legsJSON == "" {
		return []storage.OptionLeg{}, nil
	}

	var legs []storage.OptionLeg
	err := json.Unmarshal([]byte(legsJSON), &legs)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize legs: %w", err)
	}

	return legs, nil
}
