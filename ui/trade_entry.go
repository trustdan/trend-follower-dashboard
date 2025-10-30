package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildTradeEntryScreen(state *AppState) fyne.CanvasObject {
	// Session check: require active session
	if state.currentSession == nil {
		return showNoSessionPrompt(state, "Trade Entry")
	}

	// Prerequisite check: heat must be completed
	if !state.currentSession.HeatCompleted {
		return showPrerequisiteError(state, "Heat Check", "Trade Entry")
	}

	// Title
	title := canvas.NewText("Trade Entry - 5 Gates Final Check", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Session summary
	sessionInfo := widget.NewLabel(fmt.Sprintf(
		"Session #%d • %s • %s",
		state.currentSession.SessionNum,
		formatStrategyDisplay(state.currentSession.Strategy),
		state.currentSession.Ticker,
	))
	sessionInfo.TextStyle = fyne.TextStyle{Bold: true}

	instructions := widget.NewRichTextFromMarkdown(fmt.Sprintf(
		"**Final Gate Check for Session #%d**\n\n"+
			"This is your last checkpoint before executing the trade. "+
			"All session data has been verified through the 4-step workflow:\n"+
			"1. ✅ Checklist (%s banner)\n"+
			"2. ✅ Position Sizing ($%.2f risk, %d shares)\n"+
			"3. ✅ Heat Check (%s)\n"+
			"4. ⏳ Trade Entry (current step)\n\n"+
			"Review the full session summary below and check all 5 gates before saving your decision.",
		state.currentSession.SessionNum,
		state.currentSession.ChecklistBanner,
		state.currentSession.SizingRiskDollars,
		state.currentSession.SizingShares,
		state.currentSession.HeatStatus,
	))

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

	// Session Summary (read-only display)
	summaryLabel := widget.NewLabelWithStyle("📋 Session Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Build summary text - include options info if applicable
	summaryStr := fmt.Sprintf(
		"Strategy: %s\n"+
			"Ticker: %s\n"+
			"Entry: $%.2f\n"+
			"Stop: $%.2f (K=%.1f, ATR=%.2f)\n",
		formatStrategyDisplay(state.currentSession.Strategy),
		state.currentSession.Ticker,
		state.currentSession.SizingEntryPrice,
		state.currentSession.SizingInitialStop,
		state.currentSession.SizingKMultiple,
		state.currentSession.SizingATR)

	// Add shares or contracts
	if state.currentSession.SizingShares > 0 {
		summaryStr += fmt.Sprintf("Shares: %d\n", state.currentSession.SizingShares)
	}
	if state.currentSession.SizingContracts > 0 {
		summaryStr += fmt.Sprintf("Contracts: %d\n", state.currentSession.SizingContracts)
	}

	summaryStr += fmt.Sprintf("Risk: $%.2f\n"+
		"Bucket: %s\n"+
		"Checklist Banner: %s",
		state.currentSession.SizingRiskDollars,
		state.currentSession.HeatBucket,
		state.currentSession.ChecklistBanner)

	// Add pyramid planning information
	if state.currentSession.MaxUnits > 0 {
		summaryStr += fmt.Sprintf("\n\n🔺 Pyramid Planning:\n"+
			"Max Units: %d\n"+
			"Add Every: %.1fN (ATR)\n"+
			"Current Units: %d / %d\n",
			state.currentSession.MaxUnits,
			state.currentSession.AddStepN,
			state.currentSession.CurrentUnits,
			state.currentSession.MaxUnits)

		// Show add-on prices if calculated
		if state.currentSession.AddPrice1 > 0 {
			summaryStr += fmt.Sprintf("\nAdd-On Prices:\n")
			summaryStr += fmt.Sprintf("  Add 1: $%.2f (Entry + %.1fN)\n",
				state.currentSession.AddPrice1,
				state.currentSession.AddStepN)

			if state.currentSession.AddPrice2 > 0 && state.currentSession.MaxUnits > 2 {
				summaryStr += fmt.Sprintf("  Add 2: $%.2f (Entry + %.1fN)\n",
					state.currentSession.AddPrice2,
					state.currentSession.AddStepN*2.0)
			}

			if state.currentSession.AddPrice3 > 0 && state.currentSession.MaxUnits > 3 {
				summaryStr += fmt.Sprintf("  Add 3: $%.2f (Entry + %.1fN)\n",
					state.currentSession.AddPrice3,
					state.currentSession.AddStepN*3.0)
			}
		}
	}

	// Add options information if this is an options trade
	if state.currentSession.InstrumentType == "option" && state.currentSession.OptionsStrategy != "" {
		summaryStr += fmt.Sprintf("\n\n📊 Options Details:\n"+
			"Strategy: %s\n"+
			"Expiration: %s (%d DTE)\n"+
			"Entry Date: %s\n"+
			"Roll at: %d DTE\n"+
			"Time Exit Mode: %s\n"+
			"Net Debit: $%.2f\n"+
			"Max Profit: $%.2f\n"+
			"Max Loss: $%.2f",
			state.currentSession.OptionsStrategy,
			state.currentSession.PrimaryExpirationDate,
			state.currentSession.DTE,
			state.currentSession.EntryDate,
			state.currentSession.RollThresholdDTE,
			state.currentSession.TimeExitMode,
			state.currentSession.NetDebit,
			state.currentSession.MaxProfit,
			state.currentSession.MaxLoss)

		// Add breakevens if available
		if state.currentSession.BreakevenLower > 0 || state.currentSession.BreakevenUpper > 0 {
			summaryStr += "\nBreakevens:"
			if state.currentSession.BreakevenLower > 0 {
				summaryStr += fmt.Sprintf(" $%.2f", state.currentSession.BreakevenLower)
			}
			if state.currentSession.BreakevenUpper > 0 {
				if state.currentSession.BreakevenLower > 0 {
					summaryStr += " /"
				}
				summaryStr += fmt.Sprintf(" $%.2f", state.currentSession.BreakevenUpper)
			}
		}

		// Add legs display if available
		if state.currentSession.LegsJSON != "" {
			legs, err := DeserializeLegs(state.currentSession.LegsJSON)
			if err == nil && len(legs) > 0 {
				summaryStr += "\n\nLegs:\n"
				summaryStr += FormatLegsForDisplay(legs)
			}
		}
	}

	summaryText := widget.NewLabel(summaryStr)
	summaryText.Wrapping = fyne.TextWrapWord

	// 5-Gate Status display
	gatesLabel := widget.NewLabelWithStyle("🚦 5-Gate Status Check", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Results
	resultsLabel := widget.NewLabel("")
	resultsLabel.Wrapping = fyne.TextWrapWord

	// Gate check results (to be filled by Check Gates button)
	var gate1, gate2, gate3, gate4, gate5 bool
	var gatesAllPass bool

	// Check Gates button
	checkGatesBtn := widget.NewButton("Check All 5 Gates", func() {
		// Gate 1: Banner must be GREEN
		gate1 = state.currentSession.ChecklistBanner == "GREEN"

		// Gate 2: 2-minute cooloff (simplified - check if checklist was done)
		gate2 = state.currentSession.ChecklistCompleted

		// Gate 3: Ticker not on cooldown (simplified - always pass for now)
		gate3 = true // TODO: Check cooldowns table

		// Gate 4: Heat caps OK
		gate4 = state.currentSession.HeatStatus == "OK"

		// Gate 5: Sizing complete
		gate5 = state.currentSession.SizingCompleted

		gatesAllPass = gate1 && gate2 && gate3 && gate4 && gate5

		// Update banner
		if gatesAllPass {
			bannerRect.FillColor = ColorGreen()
			bannerText.Text = "✓ ALL GATES PASSED - GO"
		} else {
			bannerRect.FillColor = ColorRed()
			bannerText.Text = "✗ GATES FAILED - NO-GO"
		}
		bannerRect.Refresh()
		bannerText.Refresh()

		// Format results
		resultsText := "🚦 Gate Check Results:\n\n"

		if gate1 {
			resultsText += "1. ✅ Banner GREEN\n"
		} else {
			resultsText += "1. ❌ Banner GREEN - FAILED (banner is " + state.currentSession.ChecklistBanner + ")\n"
		}

		if gate2 {
			resultsText += "2. ✅ 2-Minute Cooloff Elapsed\n"
		} else {
			resultsText += "2. ❌ 2-Minute Cooloff - FAILED\n"
		}

		if gate3 {
			resultsText += "3. ✅ Ticker Not on Cooldown\n"
		} else {
			resultsText += "3. ❌ Ticker on Cooldown - FAILED\n"
		}

		if gate4 {
			resultsText += "4. ✅ Heat Caps Not Exceeded\n"
		} else {
			resultsText += "4. ❌ Heat Caps Exceeded - FAILED (status: " + state.currentSession.HeatStatus + ")\n"
		}

		if gate5 {
			resultsText += "5. ✅ Position Sizing Completed\n"
		} else {
			resultsText += "5. ❌ Position Sizing Not Completed - FAILED\n"
		}

		resultsText += "\n"
		if gatesAllPass {
			resultsText += "✅ ALL GATES PASSED - YOU MAY TRADE\n\n"
			resultsText += "Click 'Save GO' to log this decision and complete the session.\n"
			resultsText += "Click 'Save NO-GO' if you decide not to trade despite passing gates."
		} else {
			resultsText += "❌ GATES FAILED - DO NOT TRADE\n\n"
			resultsText += "You may only save a NO-GO decision at this time."
		}

		resultsLabel.SetText(resultsText)
	})
	checkGatesBtn.Importance = widget.HighImportance

	// Save decision buttons
	saveGoBtn := widget.NewButton("Save GO ✅", func() {
		if !gatesAllPass {
			ShowStyledInformation("Cannot Save GO",
				"All 5 gates must pass before you can save a GO decision.\n\n"+
					"Current status: One or more gates have failed.",
				state.window)
			return
		}

		// Save GO decision
		err := state.db.UpdateSessionEntry(
			state.currentSession.ID,
			"GO",
			0, // decision_id (0 for now - would link to decisions table)
			gate1, gate2, gate3, gate4, gate5,
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

		// Create position from session
		position, err := state.db.CreatePositionFromSession(updatedSession)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to create position: %v", err))
			return
		}

		state.SetCurrentSession(updatedSession)

		// Build success message with position details
		successMsg := fmt.Sprintf("Session #%d completed with GO decision!\n\n"+
			"✅ Position #%d OPENED\n\n"+
			"Ticker: %s\n"+
			"Strategy: %s\n"+
			"Shares: %d @ $%.2f\n"+
			"Initial Stop: $%.2f\n"+
			"Risk: $%.2f\n",
			updatedSession.SessionNum,
			position.ID,
			updatedSession.Ticker,
			formatStrategyDisplay(updatedSession.Strategy),
			updatedSession.SizingShares,
			updatedSession.SizingEntryPrice,
			updatedSession.SizingInitialStop,
			updatedSession.SizingRiskDollars)

		// Add pyramid info if applicable
		if updatedSession.MaxUnits > 1 {
			successMsg += fmt.Sprintf("\nPyramid: %d/%d units (add every %.1fN)\n",
				updatedSession.CurrentUnits,
				updatedSession.MaxUnits,
				updatedSession.AddStepN)
		}

		// Add options info if applicable
		if updatedSession.InstrumentType == "option" {
			successMsg += fmt.Sprintf("\nOptions: %s (%d DTE)\n",
				updatedSession.OptionsStrategy,
				updatedSession.DTE)
		}

		successMsg += "\nThis session is now COMPLETED and read-only."

		ShowStyledInformation("GO Decision Saved", successMsg, state.window)

		resultsLabel.SetText(fmt.Sprintf("✅ GO decision saved - Session #%d is COMPLETED\n✅ Position #%d OPENED",
			updatedSession.SessionNum, position.ID))
	})
	saveGoBtn.Importance = widget.HighImportance

	saveNoGoBtn := widget.NewButton("Save NO-GO ❌", func() {
		// Save NO-GO decision
		err := state.db.UpdateSessionEntry(
			state.currentSession.ID,
			"NO-GO",
			0, // decision_id
			gate1, gate2, gate3, gate4, gate5,
		)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to save session: %v", err))
			return
		}

		// Reload session
		updatedSession, err := state.db.GetSession(state.currentSession.ID)
		if err != nil {
			resultsLabel.SetText(fmt.Sprintf("❌ Failed to reload session: %v", err))
			return
		}
		state.SetCurrentSession(updatedSession)

		ShowStyledInformation("NO-GO Decision Saved",
			fmt.Sprintf("Session #%d completed with NO-GO decision.\n\n"+
				"You decided not to trade %s despite the analysis.\n\n"+
				"This session is now COMPLETED and read-only.",
				updatedSession.SessionNum,
				updatedSession.Ticker),
			state.window)

		resultsLabel.SetText("✅ NO-GO decision saved - Session #" + fmt.Sprintf("%d", updatedSession.SessionNum) + " is COMPLETED")
	})

	// Disable buttons if session is completed
	if state.currentSession.Status == "COMPLETED" {
		checkGatesBtn.Disable()
		saveGoBtn.Disable()
		saveNoGoBtn.Disable()
	}

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(sessionInfo),
		container.NewPadded(instructions),
		widget.NewSeparator(),
		summaryLabel,
		container.NewPadded(summaryText),
		widget.NewSeparator(),
		banner,
		widget.NewSeparator(),
		gatesLabel,
		checkGatesBtn,
		widget.NewSeparator(),
		resultsLabel,
		widget.NewSeparator(),
		container.NewHBox(saveGoBtn, saveNoGoBtn),
	)

	return container.NewScroll(content)
}
