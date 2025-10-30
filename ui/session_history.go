package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// buildSessionHistoryScreen creates the session history view
func buildSessionHistoryScreen(state *AppState) fyne.CanvasObject {
	// Header
	header := widget.NewLabelWithStyle("Trade Session History", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Filter dropdown
	filterOptions := []string{"All", "COMPLETED", "DRAFT", "ABANDONED"}
	currentFilter := "All"
	filterSelect := widget.NewSelect(filterOptions, func(value string) {
		currentFilter = value
		// Will refresh the list when filter changes
	})
	filterSelect.SetSelected("All")

	// Search entry
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by ticker...")

	// Filter controls
	filterRow := container.NewHBox(
		widget.NewLabel("Filter:"),
		filterSelect,
		layout.NewSpacer(),
		widget.NewLabel("Search:"),
		searchEntry,
		widget.NewButton("🔍 Search", func() {
			// Search functionality
		}),
	)

	// Session list
	var sessions []*storage.TradeSession
	var sessionList *widget.List

	// Load sessions function
	loadSessions := func() {
		allSessions, err := state.db.ListSessionHistory(100) // Get last 100 sessions
		if err != nil {
			log.Printf("Error loading session history: %v", err)
			dialog.ShowError(err, state.window)
			return
		}

		// Filter sessions
		sessions = make([]*storage.TradeSession, 0)
		searchTerm := searchEntry.Text

		for _, session := range allSessions {
			// Apply status filter
			if currentFilter != "All" && session.Status != currentFilter {
				continue
			}

			// Apply search filter
			if searchTerm != "" && session.Ticker != searchTerm {
				continue
			}

			sessions = append(sessions, session)
		}

		if sessionList != nil {
			sessionList.Refresh()
		}
	}

	// Session list widget
	sessionList = widget.NewList(
		func() int {
			return len(sessions)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),      // Session number
				widget.NewLabel(""),      // Ticker
				widget.NewLabel(""),      // Strategy
				widget.NewLabel(""),      // Options badge
				widget.NewLabel(""),      // Status
				widget.NewLabel(""),      // Decision
				widget.NewLabel(""),      // Date
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(sessions) {
				return
			}

			session := sessions[id]
			box := obj.(*fyne.Container)

			// Session number
			box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("#%d", session.SessionNum))

			// Ticker
			box.Objects[1].(*widget.Label).SetText(session.Ticker)

			// Strategy
			box.Objects[2].(*widget.Label).SetText(formatStrategyShort(session.Strategy))

			// Options badge (if applicable)
			optionsBadge := ""
			if session.InstrumentType == "option" && session.OptionsStrategy != "" {
				if session.PrimaryExpirationDate != "" {
					optionsBadge = fmt.Sprintf("[%d DTE]", session.DTE)
				} else {
					optionsBadge = "[OPT]"
				}
			}
			box.Objects[3].(*widget.Label).SetText(optionsBadge)

			// Status
			statusText := session.Status
			if session.Status == storage.StatusCompleted {
				statusText = "✅ " + statusText
			} else if session.Status == storage.StatusDraft {
				statusText = "⏳ " + statusText
			}
			box.Objects[4].(*widget.Label).SetText(statusText)

			// Decision
			decision := "-"
			if session.EntryDecision != "" {
				if session.EntryDecision == "GO" {
					decision = "✅ GO"
				} else {
					decision = "❌ NO-GO"
				}
			}
			box.Objects[5].(*widget.Label).SetText(decision)

			// Date
			dateStr := session.CreatedAt.Format("01/02")
			box.Objects[6].(*widget.Label).SetText(dateStr)
		},
	)

	// Action buttons
	var selectedSessionID int = -1
	sessionList.OnSelected = func(id widget.ListItemID) {
		selectedSessionID = int(id)
	}

	viewBtn := widget.NewButton("👁 View", func() {
		if selectedSessionID < 0 || selectedSessionID >= len(sessions) {
			dialog.ShowInformation("No Selection", "Please select a session to view.", state.window)
			return
		}

		session := sessions[selectedSessionID]
		showSessionDetailsDialog(state, session)
	})
	viewBtn.Importance = widget.HighImportance

	cloneBtn := widget.NewButton("📋 Clone", func() {
		if selectedSessionID < 0 || selectedSessionID >= len(sessions) {
			dialog.ShowInformation("No Selection", "Please select a session to clone.", state.window)
			return
		}

		session := sessions[selectedSessionID]
		cloneSession(state, session, loadSessions)
	})
	cloneBtn.Importance = widget.MediumImportance

	refreshBtn := widget.NewButton("🔄 Refresh", func() {
		loadSessions()
	})
	refreshBtn.Importance = widget.LowImportance

	actionRow := container.NewHBox(
		viewBtn,
		cloneBtn,
		layout.NewSpacer(),
		refreshBtn,
	)

	// Wire up search/filter to reload
	searchEntry.OnSubmitted = func(string) {
		loadSessions()
	}
	filterSelect.OnChanged = func(string) {
		loadSessions()
	}

	// Initial load
	loadSessions()

	// Layout
	content := container.NewBorder(
		container.NewVBox(
			header,
			widget.NewSeparator(),
			filterRow,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			actionRow,
		),
		nil, nil,
		container.NewVScroll(sessionList),
	)

	return content
}

// showSessionDetailsDialog displays full session details in read-only view
func showSessionDetailsDialog(state *AppState, session *storage.TradeSession) {
	// Build session summary
	summary := fmt.Sprintf(`**Session #%d**

**Status:** %s
**Ticker:** %s
**Strategy:** %s
**Created:** %s
**Updated:** %s

**Checklist:**
- Completed: %v
- Banner: %s

**Position Sizing:**
- Completed: %v
- Shares: %d
- Risk: $%.2f

**Heat Check:**
- Completed: %v
- Status: %s

**Trade Entry:**
- Completed: %v
- Decision: %s
`,
		session.SessionNum,
		session.Status,
		session.Ticker,
		formatStrategy(session.Strategy),
		session.CreatedAt.Format(time.RFC822),
		session.UpdatedAt.Format(time.RFC822),
		session.ChecklistCompleted,
		session.ChecklistBanner,
		session.SizingCompleted,
		session.SizingShares,
		session.SizingRiskDollars,
		session.HeatCompleted,
		session.HeatStatus,
		session.EntryCompleted,
		session.EntryDecision,
	)

	summaryWidget := widget.NewRichTextFromMarkdown(summary)

	scrollContainer := container.NewScroll(summaryWidget)
	scrollContainer.SetMinSize(fyne.NewSize(600, 500))

	dialog.ShowCustom(
		fmt.Sprintf("Session #%d Details", session.SessionNum),
		"Close",
		scrollContainer,
		state.window,
	)
}

// cloneSession creates a new session based on an existing one
func cloneSession(state *AppState, original *storage.TradeSession, refreshCallback func()) {
	// Confirm with user
	dialog.ShowConfirm(
		"Clone Session",
		fmt.Sprintf("Create a new draft session based on Session #%d (%s - %s)?\n\nThis will copy the ticker and strategy to a new session.",
			original.SessionNum,
			original.Ticker,
			formatStrategy(original.Strategy),
		),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// Create new session with same ticker and strategy
			newSession, err := state.db.CreateSession(original.Ticker, original.Strategy)
			if err != nil {
				log.Printf("Error cloning session: %v", err)
				dialog.ShowError(err, state.window)
				return
			}

			log.Printf("Session #%d cloned to Session #%d", original.SessionNum, newSession.SessionNum)

			// Refresh history list
			if refreshCallback != nil {
				refreshCallback()
			}

			// Show success
			dialog.ShowInformation(
				"Session Cloned",
				fmt.Sprintf("Session #%d created as a clone of Session #%d.\n\nThe new session is in DRAFT status with the same ticker (%s) and strategy (%s).",
					newSession.SessionNum,
					original.SessionNum,
					newSession.Ticker,
					formatStrategy(newSession.Strategy),
				),
				state.window,
			)

			// Optionally set as current session
			state.SetCurrentSession(newSession)
		},
		state.window,
	)
}
