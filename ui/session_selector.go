package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// createResumeSessionButton creates a dropdown button to resume sessions
func createResumeSessionButton(state *AppState, navigateToTab func(int)) *widget.Button {
	btn := widget.NewButton("Resume Session ▼", func() {
		showResumeSessionMenu(state, navigateToTab)
	})
	btn.Importance = widget.MediumImportance
	return btn
}

// showResumeSessionMenu displays a menu of active sessions to resume
func showResumeSessionMenu(state *AppState, navigateToTab func(int)) {
	// Get all active DRAFT sessions
	sessions, err := state.db.ListActiveSessions()
	if err != nil {
		log.Printf("Error listing active sessions: %v", err)
		dialog.ShowError(err, state.window)
		return
	}

	if len(sessions) == 0 {
		dialog.ShowInformation(
			"No Active Sessions",
			"You have no active sessions to resume.\n\nClick 'Start New Trade' to begin a new session.",
			state.window,
		)
		return
	}

	// Build list of session options
	options := make([]string, len(sessions))
	sessionMap := make(map[string]*storage.TradeSession)

	for i, session := range sessions {
		// Calculate time since last update
		timeSince := time.Since(session.UpdatedAt)
		timeStr := formatTimeSince(timeSince)

		// Build progress string
		progress := buildShortProgress(session)

		// Build options badge if applicable
		optionsBadge := ""
		if session.InstrumentType == "option" && session.OptionsStrategy != "" {
			optionsBadge = " [OPT]"
			// Add expiration info if available
			if session.PrimaryExpirationDate != "" {
				optionsBadge = fmt.Sprintf(" [OPT: %d DTE]", session.DTE)
			}
		}

		// Format: "#47 (AAPL - Long) [OPT: 60 DTE] | ✅⏳○○ | 2h ago"
		optionText := fmt.Sprintf("#%d (%s - %s)%s | %s | %s",
			session.SessionNum,
			session.Ticker,
			formatStrategyShort(session.Strategy),
			optionsBadge,
			progress,
			timeStr,
		)

		options[i] = optionText
		sessionMap[optionText] = session
	}

	// Show selection dialog
	var selectedSession *storage.TradeSession
	dialog.ShowCustomConfirm(
		"Resume Session",
		"Resume",
		"Cancel",
		widget.NewList(
			func() int { return len(options) },
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(options[i])
			},
		),
		func(resume bool) {
			if !resume || selectedSession == nil {
				return
			}

			// Set as current session
			state.SetCurrentSession(selectedSession)

			// Navigate to current step
			tabIndex := getTabIndexForStep(selectedSession.CurrentStep)
			if navigateToTab != nil {
				navigateToTab(tabIndex)
			}

			dialog.ShowInformation(
				"Session Resumed",
				fmt.Sprintf("Resumed Session #%d (%s - %s).\n\nNavigated to %s tab.",
					selectedSession.SessionNum,
					selectedSession.Ticker,
					formatStrategy(selectedSession.Strategy),
					selectedSession.CurrentStep,
				),
				state.window,
			)
		},
		state.window,
	)
}

// buildShortProgress builds a compact progress indicator
func buildShortProgress(session *storage.TradeSession) string {
	checklist := "○"
	if session.ChecklistCompleted {
		checklist = "✅"
	} else if session.CurrentStep == storage.StepChecklist {
		checklist = "⏳"
	}

	sizing := "○"
	if session.SizingCompleted {
		sizing = "✅"
	} else if session.CurrentStep == storage.StepSizing {
		sizing = "⏳"
	}

	heat := "○"
	if session.HeatCompleted {
		heat = "✅"
	} else if session.CurrentStep == storage.StepHeat {
		heat = "⏳"
	}

	entry := "○"
	if session.EntryCompleted {
		entry = "✅"
	} else if session.CurrentStep == storage.StepEntry {
		entry = "⏳"
	}

	return fmt.Sprintf("%s%s%s%s", checklist, sizing, heat, entry)
}

// formatStrategyShort converts strategy to short format
func formatStrategyShort(strategy string) string {
	switch strategy {
	case storage.StrategyLongBreakout:
		return "Long"
	case storage.StrategyShortBreakout:
		return "Short"
	case storage.StrategyCustom:
		return "Custom"
	default:
		return strategy
	}
}

// formatTimeSince formats a duration as a human-readable string
func formatTimeSince(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		mins := int(d.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}

// getTabIndexForStep returns the tab index for a given step
func getTabIndexForStep(step string) int {
	switch step {
	case storage.StepChecklist:
		return 2 // Checklist tab
	case storage.StepSizing:
		return 3 // Position Sizing tab
	case storage.StepHeat:
		return 4 // Heat Check tab
	case storage.StepEntry:
		return 5 // Trade Entry tab
	default:
		return 2 // Default to Checklist
	}
}
