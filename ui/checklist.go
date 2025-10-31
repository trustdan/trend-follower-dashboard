package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/domain"
)

func buildChecklistScreen(state *AppState) fyne.CanvasObject {
	// Get active session (real or sample)
	activeSession := state.GetActiveSession()

	// Session check: require active session (real or sample)
	if activeSession == nil && !state.sampleMode {
		return showNoSessionPrompt(state, "Checklist")
	}

	// Title
	title := canvas.NewText("Checklist - 5 Gates Evaluation", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Session info with read-only indicator
	sessionInfoText := fmt.Sprintf(
		"Session #%d • %s • %s",
		activeSession.SessionNum,
		formatStrategyDisplay(activeSession.Strategy),
		activeSession.Ticker,
	)
	if state.sampleMode {
		sessionInfoText = "📦 SAMPLE MODE: " + sessionInfoText
	} else if activeSession.Status == "COMPLETED" {
		sessionInfoText = "🔒 READ-ONLY: " + sessionInfoText
	}
	sessionInfo := widget.NewLabel(sessionInfoText)
	sessionInfo.TextStyle = fyne.TextStyle{Bold: true}

	// Ticker entry (auto-filled from session)
	tickerLabel := widget.NewLabel("Ticker Symbol:")
	tickerEntry := widget.NewEntry()
	tickerEntry.SetText(activeSession.Ticker)
	tickerEntry.SetPlaceHolder("AAPL")

	// Disable ticker entry if session is completed or in sample mode
	if activeSession.Status == "COMPLETED" || state.sampleMode {
		tickerEntry.Disable()
	}

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
		ShowStyledInformation("From Preset - Signal Required",
			"What it means:\n"+
				"The stock came from today's FINVIZ screener results, not from a random idea or tip.\n\n"+
				"Why it matters:\n"+
				"This prevents you from trading random stocks based on emotions or hunches. "+
				"You're only considering stocks that meet specific technical criteria.\n\n"+
				"Example:\n"+
				"✓ GOOD: AAPL showed up in your FINVIZ scan for \"new 55-day highs\"\n"+
				"✗ BAD: Your friend texted you about AAPL looking good\n\n"+
				"Think of it like:\n"+
				"Only dating people who meet your criteria, not random people you bump into.",
			state.window)
	})
	fromPresetHelp.Importance = widget.HighImportance

	trendCheck := widget.NewCheck("Trend Confirmed (RISK_REQ)", nil)
	trendCheck.SetChecked(false)
	trendHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ShowStyledInformation("Trend Confirmed - Risk/Sizing Required",
			"What it means:\n"+
				"For longs: Stock price just broke above its highest point in 55 days\n"+
				"For shorts: Stock price just broke below its lowest point in 55 days\n\n"+
				"Why it matters:\n"+
				"You're catching a strong momentum move. The trend is clearly established.\n"+
				"You're not trying to predict a reversal or catch a falling knife.\n\n"+
				"Example:\n"+
				"AAPL was trading between $150-$180 for 2 months, then broke above $180\n"+
				"✓ GOOD: Enter long when it breaks $180 (new 55-day high)\n"+
				"✗ BAD: Try to guess if $175 is \"good enough\" to enter\n\n"+
				"Think of it like:\n"+
				"Joining a race car that's already speeding up, not trying to time when it will start.",
			state.window)
	})
	trendHelp.Importance = widget.HighImportance

	liquidityCheck := widget.NewCheck("Liquidity OK (OPT_REQ)", nil)
	liquidityCheck.SetChecked(false)
	liquidityHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ShowStyledInformation("Liquidity OK - Options Required",
			"What it means:\n"+
				"If trading options: There's enough trading volume to get in and out easily.\n"+
				"If trading stocks: Average daily volume is high enough (usually 500K+ shares).\n\n"+
				"Why it matters:\n"+
				"Low liquidity = you might not be able to exit when you want to.\n"+
				"The bid-ask spread could eat up your profits.\n\n"+
				"Example (Options):\n"+
				"Option shows: Bid $4.80, Ask $5.20, Open Interest 250 contracts\n"+
				"✓ GOOD: Spread is $0.40 (8% of $5.00) and 250 contracts available\n"+
				"✗ BAD: Spread is $0.80 (16%) or only 20 contracts available\n\n"+
				"Example (Stocks):\n"+
				"✓ GOOD: AAPL trades 50M shares per day - easy to buy/sell\n"+
				"✗ BAD: Tiny company trades 10K shares per day - might get stuck\n\n"+
				"Think of it like:\n"+
				"Trading at a busy market vs. trying to sell something on a deserted street.",
			state.window)
	})
	liquidityHelp.Importance = widget.HighImportance

	timeframeCheck := widget.NewCheck("TV Confirm (EXIT_REQ)", nil)
	timeframeCheck.SetChecked(false)
	timeframeHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ShowStyledInformation("TV Confirm - Exit Plan Required",
			"What it means:\n"+
				"You've confirmed your exit plan BEFORE entering the trade.\n"+
				"You know exactly when you'll get out, whether you win or lose.\n\n"+
				"Why it matters:\n"+
				"No exit plan = holding losers too long and selling winners too early.\n"+
				"Emotional decisions happen when you don't know your exit in advance.\n\n"+
				"The Rule:\n"+
				"Exit when price breaks the 10-day low (for longs) OR hits your stop loss,\n"+
				"whichever happens first.\n\n"+
				"Example:\n"+
				"You bought AAPL at $180 with a stop at $177\n"+
				"✓ Exit if: Price hits $177 (stop loss) OR breaks below 10-day low\n"+
				"✗ Don't: \"Let me see what happens\" or \"I'll decide later\"\n\n"+
				"Think of it like:\n"+
				"Setting your GPS destination before driving, not figuring it out as you go.",
			state.window)
	})
	timeframeHelp.Importance = widget.HighImportance

	earningsCheck := widget.NewCheck("Earnings OK (BEHAV_REQ)", nil)
	earningsCheck.SetChecked(false)
	earningsHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ShowStyledInformation("Earnings OK - Behavior Required",
			"What it means:\n"+
				"1. You've waited at least 2 minutes since evaluating this trade\n"+
				"2. The stock doesn't have earnings coming up in the next 5 days\n"+
				"3. You're not changing your mind multiple times today\n\n"+
				"Why it matters:\n"+
				"The 2-minute wait prevents impulsive \"I gotta get in NOW!\" feelings.\n"+
				"Earnings announcements cause wild price swings that break trend systems.\n\n"+
				"Example:\n"+
				"✓ GOOD: Evaluated at 10:00am, waited until 10:02am to enter\n"+
				"✓ GOOD: Checked - no earnings until next month\n"+
				"✗ BAD: Immediately hitting buy after seeing the breakout\n"+
				"✗ BAD: Company reports earnings tomorrow morning\n\n"+
				"The 2-Minute Rule:\n"+
				"If you can't wait 2 minutes, you're probably being impulsive.\n"+
				"Good trades will still be good trades in 2 minutes.\n\n"+
				"Think of it like:\n"+
				"Waiting 2 minutes before sending an angry email - often saves you from mistakes.",
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
		ShowStyledInformation("Regime OK (Optional)",
			"What it means:\n"+
				"The overall market is moving in your direction.\n"+
				"For longs: The S&P 500 (SPY) is above its 200-day average\n"+
				"For shorts: The S&P 500 is below its 200-day average\n\n"+
				"Why it matters:\n"+
				"\"A rising tide lifts all boats\" - it's easier to make money when the\n"+
				"whole market is moving with you, not against you.\n\n"+
				"Example:\n"+
				"SPY at $450, 200-day average at $420\n"+
				"✓ GOOD for longs: Market is in uptrend (+7% above average)\n"+
				"✗ RISKY for longs: Market at $390 (below average) - fighting the tide\n\n"+
				"Not Required But:\n"+
				"If this is checked, your quality score goes up. Think of it as bonus points.\n\n"+
				"Think of it like:\n"+
				"Swimming downstream vs. upstream - both work, but one is easier.",
			state.window)
	})
	regimeHelp.Importance = widget.HighImportance

	chaseCheck := widget.NewCheck("No Chase (< 2N above 20-EMA)", nil)
	chaseCheck.SetChecked(false)
	chaseHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ShowStyledInformation("No Chase (Optional)",
			"What it means:\n"+
				"The stock hasn't run up too far, too fast from its recent average price.\n"+
				"Specifically: Entry price isn't more than 2× ATR above the 20-day average.\n\n"+
				"Why it matters:\n"+
				"Chasing stocks that have already run too far often means buying at the top.\n"+
				"You want to catch the move early, not late.\n\n"+
				"Example:\n"+
				"Stock's 20-day average: $100\n"+
				"ATR (daily volatility): $3\n"+
				"Current price: $104\n"+
				"✓ GOOD: Only $4 above average (< 2×$3 = $6 limit)\n"+
				"✗ RISKY: Price at $108 (too extended - likely near a pullback)\n\n"+
				"Not Required But:\n"+
				"Checking this improves quality score. It's okay to trade extended stocks,\n"+
				"but you're taking on more risk of a near-term pullback.\n\n"+
				"Think of it like:\n"+
				"Joining a party that just started vs. showing up at 2am when it's winding down.",
			state.window)
	})
	chaseHelp.Importance = widget.HighImportance

	journalCheck := widget.NewCheck("Journal Entry Written", nil)
	journalCheck.SetChecked(false)
	journalHelp := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		ShowStyledInformation("Journal Entry Written (Optional)",
			"What it means:\n"+
				"You wrote down your trade plan BEFORE entering:\n"+
				"• Why are you taking this trade right now?\n"+
				"• What's your profit target?\n"+
				"• What's your stop loss?\n"+
				"• What could go wrong?\n\n"+
				"Why it matters:\n"+
				"Writing forces you to think clearly. If you can't explain the trade\n"+
				"in writing, you probably don't understand it well enough to risk money.\n\n"+
				"What to Write:\n"+
				"\"AAPL broke above 55-day high at $180. Entry at $181, stop at $177.\n"+
				"Target: hold for trend exit at 10-day low. Risk $300 for $900+ potential.\n"+
				"Could fail if market sells off or sector rotation happens.\"\n\n"+
				"Not Required But:\n"+
				"Highly recommended. Your future self will thank you when reviewing trades.\n"+
				"Improves quality score.\n\n"+
				"Think of it like:\n"+
				"Planning a road trip vs. just getting in the car and driving randomly.",
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

	// Next button (shown after GREEN banner)
	nextBtn := widget.NewButton("Next: Position Sizing →", func() {
		// Navigate to Position Sizing tab (index 3)
		if state.navigateToTab != nil {
			state.navigateToTab(3)
		}
	})
	nextBtn.Importance = widget.HighImportance
	nextBtn.Hide() // Hidden initially

	// Evaluate button
	evaluateBtn := widget.NewButton("Evaluate Checklist", func() {
		// In sample mode, just show sample results without database updates
		if state.sampleMode {
			// Pre-check all boxes for sample mode
			fromPresetCheck.SetChecked(true)
			trendCheck.SetChecked(true)
			liquidityCheck.SetChecked(true)
			timeframeCheck.SetChecked(true)
			earningsCheck.SetChecked(true)
			regimeCheck.SetChecked(true)
			chaseCheck.SetChecked(true)
			journalCheck.SetChecked(true)

			// Show GREEN banner
			bannerRect.FillColor = ColorGreen()
			bannerText.Text = "✓ GREEN - OK TO TRADE"
			bannerRect.Refresh()
			bannerText.Refresh()

			resultsLabel.SetText("📦 SAMPLE MODE\n\nBanner: GREEN\nMissing Required: 0\n\n✓ All gates passed! (Sample data)\n✓ Sample session ready for Position Sizing")
			nextBtn.Show()
			return
		}

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

		// Calculate quality score from optional checks
		qualityScore := 0
		if regimeCheck.Checked {
			qualityScore++
		}
		if chaseCheck.Checked {
			qualityScore++
		}
		if journalCheck.Checked {
			qualityScore++
		}

		// Update session in database
		err = state.db.UpdateSessionChecklist(
			activeSession.ID,
			result.Banner,
			result.MissingCount,
			qualityScore,
		)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to save session: %v", err))
			return
		}

		// Reload session to get updated state
		updatedSession, err := state.db.GetSession(activeSession.ID)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to reload session: %v", err))
			return
		}
		state.SetCurrentSession(updatedSession)

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
			resultsText += fmt.Sprintf("\n✓ Session #%d updated - ready for Position Sizing", activeSession.SessionNum)
			nextBtn.Show()
		} else {
			nextBtn.Hide()
		}

		resultsLabel.SetText(resultsText)
	})
	evaluateBtn.Importance = widget.HighImportance

	// Disable evaluate button if session is completed (but allow in sample mode)
	if activeSession.Status == "COMPLETED" && !state.sampleMode {
		evaluateBtn.Disable()
	}

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

	buttons := container.NewHBox(evaluateBtn, resetBtn, nextBtn)

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(sessionInfo),
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
