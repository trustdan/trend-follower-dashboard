package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// sessionBarWidget displays the current trade session status
type sessionBarWidget struct {
	widget.BaseWidget
	state         *AppState
	sessionLabel  *widget.Label
	progressLabel *widget.Label
	background    *canvas.Rectangle
	container     *fyne.Container
}

// NewSessionBar creates a new session bar widget
func NewSessionBar(state *AppState) *sessionBarWidget {
	bar := &sessionBarWidget{
		state:         state,
		sessionLabel:  widget.NewLabel("No Active Session"),
		progressLabel: widget.NewLabel(""),
		background:    canvas.NewRectangle(color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
	}

	bar.sessionLabel.TextStyle = fyne.TextStyle{Bold: true}
	bar.ExtendBaseWidget(bar)

	// Register callback to update when session changes
	state.RegisterSessionChangeCallback(func(session *storage.TradeSession) {
		bar.updateFromSession(session)
	})

	// Initial update
	bar.updateFromSession(state.currentSession)

	return bar
}

// CreateRenderer implements fyne.Widget
func (s *sessionBarWidget) CreateRenderer() fyne.WidgetRenderer {
	s.container = container.NewVBox(
		container.NewBorder(nil, nil, nil, nil, s.background),
		container.NewHBox(
			s.sessionLabel,
			layout.NewSpacer(),
			s.progressLabel,
		),
	)
	return widget.NewSimpleRenderer(s.container)
}

// updateFromSession updates the widget to show the current session
func (s *sessionBarWidget) updateFromSession(session *storage.TradeSession) {
	if session == nil {
		s.sessionLabel.SetText("No Active Session")
		s.progressLabel.SetText("")
		s.background.FillColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
		s.background.Refresh()
		s.Refresh()
		return
	}

	// Build session label: "Session #47 • LONG_BREAKOUT • AAPL [OPT: 60 DTE]"
	sessionText := fmt.Sprintf("Session #%d • %s • %s",
		session.SessionNum,
		formatStrategy(session.Strategy),
		session.Ticker,
	)

	// Add options badge if applicable
	if session.InstrumentType == "option" && session.OptionsStrategy != "" {
		if session.PrimaryExpirationDate != "" {
			sessionText += fmt.Sprintf(" [%s: %d DTE]",
				storage.GetStrategyDisplayName(session.OptionsStrategy),
				session.DTE)
		} else {
			sessionText += " [OPTIONS]"
		}
	}

	// Build progress indicator: "✅ Checklist | ⏳ Sizing | ○ Heat | ○ Entry"
	progress := buildProgressIndicator(session)

	s.sessionLabel.SetText(sessionText)
	s.progressLabel.SetText(progress)

	// Color based on status
	if session.Status == storage.StatusCompleted {
		s.background.FillColor = ColorGreen() // Green for completed
	} else if session.ChecklistBanner == "RED" {
		s.background.FillColor = ColorRed() // Red if banner is red
	} else if session.ChecklistBanner == "YELLOW" {
		s.background.FillColor = ColorYellow() // Yellow if banner is yellow
	} else if session.ChecklistBanner == "GREEN" {
		s.background.FillColor = theme.PrimaryColor() // Blue for active green session
	} else {
		s.background.FillColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255} // Default gray
	}

	s.background.Refresh()
	s.Refresh()
}

// formatStrategy converts strategy constant to readable format
func formatStrategy(strategy string) string {
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

// buildProgressIndicator builds the progress string showing gate completion
func buildProgressIndicator(session *storage.TradeSession) string {
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

	return fmt.Sprintf("%s Checklist | %s Sizing | %s Heat | %s Entry",
		checklist, sizing, heat, entry)
}
