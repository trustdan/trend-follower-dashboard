package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// formatStrategyDisplay converts strategy constant to display text
func formatStrategyDisplay(strategy string) string {
	switch strategy {
	case storage.StrategyLongBreakout:
		return "Long Breakout"
	case storage.StrategyShortBreakout:
		return "Short Breakout"
	case storage.StrategyCustom:
		return "Custom"
	default:
		return strategy
	}
}

// showNoSessionPrompt displays a message prompting the user to create or resume a session
func showNoSessionPrompt(state *AppState, tabName string) fyne.CanvasObject {
	title := canvas.NewText("No Active Trade Session", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel("You must create or resume a trade session before using the " + tabName + " tab.")
	message.Wrapping = fyne.TextWrapWord

	instructions := widget.NewLabel(
		"Trade sessions create cohesion across all tabs by tracking a single trade evaluation from start to finish.\n\n" +
			"Click \"Start New Trade\" to begin evaluating a new setup, or \"Resume Session\" to continue a previous evaluation.",
	)
	instructions.Wrapping = fyne.TextWrapWord

	// We can't call showNewTradeDialog here because it would create a circular dependency
	// Instead, just provide instruction text
	startBtn := widget.NewButton("📝 Instructions: Click \"Start New Trade\" in top bar", func() {
		ShowStyledInformation("How to Start a Session",
			"1. Click the \"Start New Trade\" button in the top bar\n"+
				"2. Select your strategy (Long Breakout, Short Breakout, or Custom)\n"+
				"3. Optionally enter a ticker symbol\n"+
				"4. Click Create Session\n\n"+
				"You'll be taken to the Checklist tab to begin evaluation.",
			state.window)
	})
	startBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
		container.NewPadded(message),
		container.NewPadded(instructions),
		widget.NewSeparator(),
		container.NewCenter(startBtn),
	)

	return container.NewCenter(content)
}

// showPrerequisiteError shows an error when a user tries to access a tab without completing prerequisites
func showPrerequisiteError(state *AppState, requiredTab, currentTab string) fyne.CanvasObject {
	title := canvas.NewText("⚠️ Prerequisite Not Met", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Color = ColorRed()

	sessionInfo := widget.NewLabel(fmt.Sprintf(
		"Session #%d (%s) - %s",
		state.currentSession.SessionNum,
		state.currentSession.Ticker,
		state.currentSession.Strategy,
	))
	sessionInfo.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel(fmt.Sprintf(
		"You must complete the %s tab before accessing %s.",
		requiredTab, currentTab,
	))
	message.Wrapping = fyne.TextWrapWord

	var detailMsg string
	switch requiredTab {
	case "Checklist":
		if state.currentSession.ChecklistBanner == "RED" {
			detailMsg = fmt.Sprintf(
				"The Checklist tab resulted in a RED banner (DO NOT TRADE).\n\n"+
					"Missing required items: %d\n\n"+
					"You must resolve all required checklist items before proceeding to Position Sizing.",
				state.currentSession.ChecklistMissingCount,
			)
		} else if state.currentSession.ChecklistBanner == "YELLOW" {
			detailMsg = "The Checklist tab resulted in a YELLOW banner (CAUTION).\n\n" +
				"You must achieve a GREEN banner before proceeding to Position Sizing."
		} else {
			detailMsg = "The Checklist tab has not been completed.\n\n" +
				"Please complete the checklist evaluation first."
		}
	case "Position Sizing":
		detailMsg = "Position sizing has not been calculated.\n\n" +
			"Please complete position sizing before checking portfolio heat."
	case "Heat Check":
		detailMsg = "Heat check has not been performed.\n\n" +
			"Please verify portfolio and bucket heat caps before proceeding to trade entry."
	}

	details := widget.NewLabel(detailMsg)
	details.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton(fmt.Sprintf("Go to %s Tab", requiredTab), func() {
		// This would need access to the tab navigation function
		ShowStyledInformation("Navigation",
			fmt.Sprintf("Please use the tab bar to navigate to the %s tab.", requiredTab),
			state.window)
	})
	backBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(sessionInfo),
		widget.NewSeparator(),
		container.NewPadded(message),
		container.NewPadded(details),
		widget.NewSeparator(),
		container.NewCenter(backBtn),
	)

	return container.NewCenter(content)
}
