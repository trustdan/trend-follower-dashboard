package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/domain"
)

func buildChecklistScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("Checklist - 5 Gates Evaluation", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Ticker entry
	tickerLabel := widget.NewLabel("Ticker Symbol:")
	tickerEntry := widget.NewEntry()
	tickerEntry.SetPlaceHolder("AAPL")

	// Banner (initially gray)
	bannerRect := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	bannerRect.SetMinSize(fyne.NewSize(0, 100))

	bannerText := canvas.NewText("EVALUATE CHECKLIST", color.White)
	bannerText.TextSize = 32
	bannerText.TextStyle = fyne.TextStyle{Bold: true}
	bannerText.Alignment = fyne.TextAlignCenter

	banner := container.NewStack(
		bannerRect,
		container.NewCenter(bannerText),
	)

	// Section 1: Required Gates
	requiredLabel := widget.NewLabelWithStyle("Required Gates (All Must Pass)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	fromPresetCheck := widget.NewCheck("From Preset (SIG_REQ)", nil)
	fromPresetCheck.SetChecked(false)
	fromPresetHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("From Preset",
			"Ticker came from today's FINVIZ preset scan (55-bar breakout filter).\n\n"+
				"This ensures we're trading mechanical breakouts, not random stocks.",
			state.window)
	})
	fromPresetHelp.Importance = widget.HighImportance

	trendCheck := widget.NewCheck("Trend Confirmed (RISK_REQ)", nil)
	trendCheck.SetChecked(false)
	trendHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("Trend Confirmed",
			"Trend confirmed: Long > 55-high OR Short < 55-low. Uses 2×N stop distance.\n\n"+
				"This is the Donchian breakout filter - the core signal.",
			state.window)
	})
	trendHelp.Importance = widget.HighImportance

	liquidityCheck := widget.NewCheck("Liquidity OK (OPT_REQ)", nil)
	liquidityCheck.SetChecked(false)
	liquidityHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("Liquidity OK",
			"Options liquidity: bid-ask < 10% of mid, OI > 100, 60-90 DTE.\n\n"+
				"For options strategies only. Ensures you can enter and exit without slippage.",
			state.window)
	})
	liquidityHelp.Importance = widget.HighImportance

	timeframeCheck := widget.NewCheck("TV Confirm (EXIT_REQ)", nil)
	timeframeCheck.SetChecked(false)
	timeframeHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("TV Confirm",
			"Exit plan confirmed: 10-bar opposite Donchian OR 2×N, whichever closer.\n\n"+
				"Know your exit BEFORE you enter. This is non-negotiable.",
			state.window)
	})
	timeframeHelp.Importance = widget.HighImportance

	earningsCheck := widget.NewCheck("Earnings OK (BEHAV_REQ)", nil)
	earningsCheck.SetChecked(false)
	earningsHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("Earnings OK",
			"2-minute cooloff passed, no intraday overrides, earnings blackout OK.\n\n"+
				"The 2-minute timer prevents impulsive entries. It's intentional friction.",
			state.window)
	})
	earningsHelp.Importance = widget.HighImportance

	requiredChecks := container.NewVBox(
		requiredLabel,
		container.NewBorder(nil, nil, nil, fromPresetHelp, fromPresetCheck),
		container.NewBorder(nil, nil, nil, trendHelp, trendCheck),
		container.NewBorder(nil, nil, nil, liquidityHelp, liquidityCheck),
		container.NewBorder(nil, nil, nil, timeframeHelp, timeframeCheck),
		container.NewBorder(nil, nil, nil, earningsHelp, earningsCheck),
	)

	// Section 2: Optional Quality Items
	optionalLabel := widget.NewLabelWithStyle("Optional Quality Items (Improve Score)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	regimeCheck := widget.NewCheck("Regime OK (e.g., SPY > 200SMA)", nil)
	regimeCheck.SetChecked(false)
	regimeHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("Regime OK",
			"Market regime favorable (e.g., SPY > 200SMA for longs).\n\n"+
				"Optional but improves quality score. Trading with the market tide.",
			state.window)
	})
	regimeHelp.Importance = widget.HighImportance

	chaseCheck := widget.NewCheck("No Chase (< 2N above 20-EMA)", nil)
	chaseCheck.SetChecked(false)
	chaseHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("No Chase",
			"Entry not > 2N above 20-EMA (avoids chasing extended moves).\n\n"+
				"Optional but reduces risk of buying tops. Patience pays.",
			state.window)
	})
	chaseHelp.Importance = widget.HighImportance

	journalCheck := widget.NewCheck("Journal Entry Written", nil)
	journalCheck.SetChecked(false)
	journalHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		dialog.ShowInformation("Journal OK",
			"Trade plan documented: why now, profit targets, risk/reward.\n\n"+
				"Optional but highly recommended. Writing forces clarity.",
			state.window)
	})
	journalHelp.Importance = widget.HighImportance

	optionalChecks := container.NewVBox(
		optionalLabel,
		container.NewBorder(nil, nil, nil, regimeHelp, regimeCheck),
		container.NewBorder(nil, nil, nil, chaseHelp, chaseCheck),
		container.NewBorder(nil, nil, nil, journalHelp, journalCheck),
	)

	// Results display
	resultsLabel := widget.NewLabel("")
	resultsLabel.Wrapping = fyne.TextWrapWord

	// Evaluate button
	evaluateBtn := widget.NewButton("Evaluate Checklist", func() {
		ticker := tickerEntry.Text
		if ticker == "" {
			resultsLabel.SetText("❌ Please enter a ticker symbol")
			return
		}

		// Call backend checklist evaluation
		req := domain.ChecklistRequest{
			Ticker:        ticker,
			FromPreset:    fromPresetCheck.Checked,
			TrendPass:     trendCheck.Checked,
			LiquidityPass: liquidityCheck.Checked,
			TVConfirm:     timeframeCheck.Checked,
			EarningsOK:    earningsCheck.Checked,
			JournalOK:     journalCheck.Checked,
		}
		result, err := domain.EvaluateChecklist(req)

		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Error: %v", err))
			return
		}

		// Update banner based on result
		switch result.Banner {
		case "GREEN":
			bannerRect.FillColor = ColorGreen()
			bannerText.Text = "✓ GREEN - OK TO TRADE"
		case "YELLOW":
			bannerRect.FillColor = ColorYellow()
			bannerText.Text = "⚠ YELLOW - CAUTION"
		case "RED":
			bannerRect.FillColor = ColorRed()
			bannerText.Text = "✗ RED - DO NOT TRADE"
		}
		bannerRect.Refresh()
		bannerText.Refresh()

		// Update results
		resultsText := fmt.Sprintf("Banner: %s\nMissing Required: %d\n",
			result.Banner, result.MissingCount)

		if len(result.MissingItems) > 0 {
			resultsText += fmt.Sprintf("\nMissing Items:\n")
			for _, item := range result.MissingItems {
				resultsText += fmt.Sprintf("  • %s\n", item)
			}
		}

		if result.Banner == "GREEN" {
			resultsText += "\n✓ All gates passed! 2-minute cooloff timer started."
		}

		resultsLabel.SetText(resultsText)
	})
	evaluateBtn.Importance = widget.HighImportance

	// Reset button
	resetBtn := widget.NewButton("Reset", func() {
		tickerEntry.SetText("")
		fromPresetCheck.SetChecked(false)
		trendCheck.SetChecked(false)
		liquidityCheck.SetChecked(false)
		timeframeCheck.SetChecked(false)
		earningsCheck.SetChecked(false)
		regimeCheck.SetChecked(false)
		chaseCheck.SetChecked(false)
		journalCheck.SetChecked(false)
		resultsLabel.SetText("")

		// Reset banner
		bannerRect.FillColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}
		bannerText.Text = "EVALUATE CHECKLIST"
		bannerRect.Refresh()
		bannerText.Refresh()
	})
	resetBtn.Importance = widget.HighImportance

	buttons := container.NewHBox(evaluateBtn, resetBtn)

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		banner,
		widget.NewSeparator(),
		tickerLabel,
		tickerEntry,
		widget.NewSeparator(),
		requiredChecks,
		widget.NewSeparator(),
		optionalChecks,
		widget.NewSeparator(),
		buttons,
		widget.NewSeparator(),
		resultsLabel,
	)

	return container.NewScroll(content)
}
