package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

// buildCalendarScreen creates the 10-week calendar view for trade diversification tracking
func buildCalendarScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("📅 Trade Calendar - 10-Week Sector View", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := widget.NewLabel("Rolling 10-week view (2 weeks back + 8 weeks forward) - Track sector diversification")
	subtitle.Wrapping = fyne.TextWrapWord

	// Filters
	var sectorFilter, strategyFilter, statusFilter string

	// Calculate date range (2 weeks back + 8 weeks forward = 10 weeks total)
	now := time.Now()
	startDate := now.AddDate(0, 0, -14) // 2 weeks back
	endDate := now.AddDate(0, 0, 56)    // 8 weeks forward

	// Container for the grid (will be rebuilt on filter changes)
	gridContainer := container.NewVBox()

	// Function to rebuild the calendar grid
	rebuildCalendar := func() {
		gridContainer.Objects = nil // Clear existing content

		// Query trade_history from database
		startDateStr := startDate.Format("2006-01-02")
		endDateStr := endDate.Format("2006-01-02")

		trades, err := state.db.GetCalendarView(startDateStr, endDateStr, statusFilter)
		if err != nil {
			log.Printf("Error loading calendar view: %v", err)
			gridContainer.Add(widget.NewLabel(fmt.Sprintf("Error loading trades: %v", err)))
			gridContainer.Refresh()
			return
		}

		// Apply additional filters
		filteredTrades := filterTrades(trades, sectorFilter, strategyFilter, statusFilter)

		// Build the calendar grid
		grid := buildCalendarGrid(state, startDate, filteredTrades)
		gridContainer.Add(grid)
		gridContainer.Refresh()
	}

	// Sector filter dropdown
	sectorOptions := []string{
		"All Sectors",
		"Materials/Industrials",
		"Tech/Comm",
		"Financial/Cyclical",
		"Defensive/Utilities",
		"Energy",
		"Healthcare",
		"Real Estate",
		"ETFs",
	}
	sectorSelect := widget.NewSelect(sectorOptions, func(selected string) {
		if selected == "All Sectors" {
			sectorFilter = ""
		} else {
			sectorFilter = selected
		}
		rebuildCalendar()
	})
	sectorSelect.SetSelected("All Sectors")

	// Strategy filter dropdown
	strategyOptions := []string{
		"All Strategies",
		"Long Breakout",
		"Short Breakout",
		"Custom",
	}
	strategySelect := widget.NewSelect(strategyOptions, func(selected string) {
		if selected == "All Strategies" {
			strategyFilter = ""
		} else {
			strategyFilter = selected
		}
		rebuildCalendar()
	})
	strategySelect.SetSelected("All Strategies")

	// Status filter dropdown
	statusOptions := []string{
		"All",
		"Open Only",
		"Closed Only",
		"Rolled Only",
	}
	statusSelect := widget.NewSelect(statusOptions, func(selected string) {
		switch selected {
		case "Open Only":
			statusFilter = "OPEN"
		case "Closed Only":
			statusFilter = "CLOSED"
		case "Rolled Only":
			statusFilter = "ROLLED"
		default:
			statusFilter = ""
		}
		rebuildCalendar()
	})
	statusSelect.SetSelected("All")

	// Filter controls
	filterControls := container.NewHBox(
		widget.NewLabel("Sector:"),
		sectorSelect,
		layout.NewSpacer(),
		widget.NewLabel("Strategy:"),
		strategySelect,
		layout.NewSpacer(),
		widget.NewLabel("Status:"),
		statusSelect,
	)

	// Refresh button
	refreshBtn := widget.NewButton("🔄 Refresh Calendar", func() {
		rebuildCalendar()
	})
	refreshBtn.Importance = widget.HighImportance

	// Legend
	legend := createCalendarLegend()

	// Build initial calendar
	rebuildCalendar()

	// Layout
	content := container.NewBorder(
		container.NewVBox(
			container.NewPadded(title),
			subtitle,
			widget.NewSeparator(),
			filterControls,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			legend,
			container.NewHBox(layout.NewSpacer(), refreshBtn),
		),
		nil,
		nil,
		container.NewScroll(gridContainer),
	)

	return content
}

// buildCalendarGrid creates the actual calendar grid with trades
func buildCalendarGrid(state *AppState, startDate time.Time, trades []storage.TradeHistoryEntry) fyne.CanvasObject {
	// Define sector buckets
	buckets := []string{
		"Materials/Industrials",
		"Tech/Comm",
		"Financial/Cyclical",
		"Defensive/Utilities",
		"Energy",
		"Healthcare",
		"Real Estate",
		"ETFs",
	}

	// Group trades by bucket and week
	tradesByBucketWeek := make(map[string]map[int][]storage.TradeHistoryEntry)

	for _, trade := range trades {
		bucket := trade.Bucket
		if bucket == "" {
			bucket = trade.Sector // Fallback to sector if bucket not set
		}
		if bucket == "" {
			bucket = "Other"
		}

		// Parse entry date and calculate which week it belongs to
		entryDate, err := time.Parse("2006-01-02", trade.EntryDate)
		if err != nil {
			log.Printf("Error parsing entry date %s: %v", trade.EntryDate, err)
			continue
		}

		// Calculate week index (0-9)
		weekIndex := int(entryDate.Sub(startDate).Hours() / 24 / 7)
		if weekIndex < 0 || weekIndex >= 10 {
			continue // Trade outside our 10-week window
		}

		if tradesByBucketWeek[bucket] == nil {
			tradesByBucketWeek[bucket] = make(map[int][]storage.TradeHistoryEntry)
		}
		tradesByBucketWeek[bucket][weekIndex] = append(tradesByBucketWeek[bucket][weekIndex], trade)
	}

	// Build grid container
	grid := container.NewVBox()

	// Header row (week dates)
	headerRow := container.NewHBox(
		widget.NewLabelWithStyle("Sector / Bucket", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	for i := 0; i < 10; i++ {
		weekDate := startDate.AddDate(0, 0, i*7)
		weekStr := weekDate.Format("Jan 2")

		// Highlight current week
		weekLabel := widget.NewLabelWithStyle(weekStr, fyne.TextAlignCenter, fyne.TextStyle{})
		if i == 2 { // Week 2 is the current week (we go 2 weeks back)
			weekLabel.TextStyle = fyne.TextStyle{Bold: true}
		}
		headerRow.Add(weekLabel)
	}
	grid.Add(headerRow)
	grid.Add(widget.NewSeparator())

	// Rows for each bucket
	for _, bucket := range buckets {
		// Skip bucket if no trades
		if tradesByBucketWeek[bucket] == nil {
			continue
		}

		// Find max trades in any week for this bucket (for row spanning)
		maxTradesInWeek := 0
		for weekIdx := 0; weekIdx < 10; weekIdx++ {
			if len(tradesByBucketWeek[bucket][weekIdx]) > maxTradesInWeek {
				maxTradesInWeek = len(tradesByBucketWeek[bucket][weekIdx])
			}
		}

		// Create rows for this bucket (one row per trade position)
		for rowIdx := 0; rowIdx < maxTradesInWeek; rowIdx++ {
			row := container.NewHBox()

			// Bucket label (only on first row)
			if rowIdx == 0 {
				if maxTradesInWeek == 1 {
					row.Add(widget.NewLabel(bucket))
				} else {
					row.Add(widget.NewLabel(fmt.Sprintf("%s (%d)", bucket, maxTradesInWeek)))
				}
			} else {
				row.Add(widget.NewLabel("")) // Empty label for spacing
			}

			// Week cells
			for weekIdx := 0; weekIdx < 10; weekIdx++ {
				trades := tradesByBucketWeek[bucket][weekIdx]

				if rowIdx < len(trades) {
					trade := trades[rowIdx]
					cellWidget := createTradeCell(state, trade)
					row.Add(cellWidget)
				} else {
					// Empty cell
					emptyLabel := widget.NewLabel("-")
					emptyLabel.Alignment = fyne.TextAlignCenter
					row.Add(emptyLabel)
				}
			}

			grid.Add(row)
		}

		// Add separator between buckets
		grid.Add(widget.NewSeparator())
	}

	// If no trades, show message
	if len(tradesByBucketWeek) == 0 {
		grid.Add(widget.NewLabel("No trades found for the selected date range and filters."))
	}

	return grid
}

// createTradeCell creates a clickable, color-coded cell for a trade
func createTradeCell(state *AppState, trade storage.TradeHistoryEntry) fyne.CanvasObject {
	// Create button with ticker
	btn := widget.NewButton(trade.Ticker, func() {
		showTradeDetailsDialog(state, trade)
	})

	// Color-code based on status
	switch trade.Status {
	case "OPEN":
		btn.Importance = widget.SuccessImportance // Green
	case "CLOSED":
		if trade.Outcome == "WIN" {
			btn.Importance = widget.SuccessImportance // Green
		} else if trade.Outcome == "LOSS" {
			btn.Importance = widget.DangerImportance // Red
		} else {
			btn.Importance = widget.MediumImportance // Yellow (scratch)
		}
	case "ROLLED":
		btn.Importance = widget.WarningImportance // Yellow/Orange
	default:
		btn.Importance = widget.LowImportance
	}

	return btn
}

// showTradeDetailsDialog displays detailed information about a trade
func showTradeDetailsDialog(state *AppState, trade storage.TradeHistoryEntry) {
	// Format trade details
	details := formatTradeDetails(trade)

	// Create scrollable content
	content := widget.NewRichTextFromMarkdown(details)
	scrollContent := container.NewScroll(content)
	scrollContent.SetMinSize(fyne.NewSize(600, 400))

	// Show dialog
	dialog.ShowCustom(
		fmt.Sprintf("Trade Details: %s", trade.Ticker),
		"Close",
		scrollContent,
		state.window,
	)
}

// formatTradeDetails creates a markdown-formatted string with trade details
func formatTradeDetails(trade storage.TradeHistoryEntry) string {
	var md string

	md += fmt.Sprintf("# %s - %s\n\n", trade.Ticker, trade.Strategy)

	// Status badge
	statusEmoji := "🔵"
	switch trade.Status {
	case "OPEN":
		statusEmoji = "🟢"
	case "CLOSED":
		if trade.Outcome == "WIN" {
			statusEmoji = "🟢"
		} else if trade.Outcome == "LOSS" {
			statusEmoji = "🔴"
		} else {
			statusEmoji = "🟡"
		}
	case "ROLLED":
		statusEmoji = "🔄"
	}
	md += fmt.Sprintf("**Status:** %s %s\n\n", statusEmoji, trade.Status)

	// Basic info
	md += "## Trade Information\n\n"
	md += fmt.Sprintf("- **Ticker:** %s\n", trade.Ticker)
	md += fmt.Sprintf("- **Strategy:** %s\n", trade.Strategy)
	if trade.BreakoutSystem != "" {
		md += fmt.Sprintf("- **Breakout System:** %s\n", trade.BreakoutSystem)
	}
	if trade.OptionsStrategy != "" {
		md += fmt.Sprintf("- **Options Strategy:** %s\n", trade.OptionsStrategy)
	}
	md += fmt.Sprintf("- **Instrument Type:** %s\n", trade.InstrumentType)
	md += fmt.Sprintf("- **Sector:** %s\n", trade.Sector)
	if trade.Bucket != "" && trade.Bucket != trade.Sector {
		md += fmt.Sprintf("- **Bucket:** %s\n", trade.Bucket)
	}
	md += "\n"

	// Dates
	md += "## Dates\n\n"
	md += fmt.Sprintf("- **Entry Date:** %s\n", trade.EntryDate)
	if trade.ExpirationDate != "" {
		md += fmt.Sprintf("- **Expiration Date:** %s\n", trade.ExpirationDate)
		if trade.DTE > 0 {
			md += fmt.Sprintf("- **DTE at Entry:** %d days\n", trade.DTE)
		}
	}
	if trade.ExitDate != "" {
		md += fmt.Sprintf("- **Exit Date:** %s\n", trade.ExitDate)
	}
	md += "\n"

	// Position sizing
	md += "## Position Details\n\n"
	if trade.Shares > 0 {
		md += fmt.Sprintf("- **Shares:** %d\n", trade.Shares)
	}
	if trade.Contracts > 0 {
		md += fmt.Sprintf("- **Contracts:** %d\n", trade.Contracts)
	}
	if trade.RiskDollars > 0 {
		md += fmt.Sprintf("- **Risk:** $%.2f\n", trade.RiskDollars)
	}
	if trade.EntryPrice > 0 {
		md += fmt.Sprintf("- **Entry Price:** $%.2f\n", trade.EntryPrice)
	}
	if trade.ExitPrice > 0 {
		md += fmt.Sprintf("- **Exit Price:** $%.2f\n", trade.ExitPrice)
	}
	md += "\n"

	// P&L and outcome
	if trade.Status == "CLOSED" {
		md += "## Results\n\n"
		if trade.PnL != 0 {
			if trade.PnL > 0 {
				md += fmt.Sprintf("- **P&L:** +$%.2f ✅\n", trade.PnL)
			} else {
				md += fmt.Sprintf("- **P&L:** -$%.2f ❌\n", -trade.PnL)
			}
		}
		if trade.Outcome != "" {
			md += fmt.Sprintf("- **Outcome:** %s\n", trade.Outcome)
		}
		md += "\n"
	}

	// Notes
	if trade.Notes != "" {
		md += "## Notes\n\n"
		md += trade.Notes + "\n\n"
	}

	// Metadata
	md += "---\n\n"
	md += fmt.Sprintf("*Trade ID: %d*\n", trade.ID)
	if trade.SessionID != nil {
		md += fmt.Sprintf("*Session ID: %d*\n", *trade.SessionID)
	}

	return md
}

// createCalendarLegend creates the legend for calendar status colors
func createCalendarLegend() fyne.CanvasObject {
	legendTitle := widget.NewLabelWithStyle("Legend:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Status indicators
	openBtn := widget.NewButton("OPEN", nil)
	openBtn.Importance = widget.SuccessImportance
	openBtn.Disable()

	winBtn := widget.NewButton("WIN", nil)
	winBtn.Importance = widget.SuccessImportance
	winBtn.Disable()

	lossBtn := widget.NewButton("LOSS", nil)
	lossBtn.Importance = widget.DangerImportance
	lossBtn.Disable()

	scratchBtn := widget.NewButton("SCRATCH", nil)
	scratchBtn.Importance = widget.MediumImportance
	scratchBtn.Disable()

	rolledBtn := widget.NewButton("ROLLED", nil)
	rolledBtn.Importance = widget.WarningImportance
	rolledBtn.Disable()

	legend := container.NewHBox(
		legendTitle,
		openBtn,
		widget.NewLabel("="),
		widget.NewLabel("Open position"),
		widget.NewLabel("•"),
		winBtn,
		widget.NewLabel("="),
		widget.NewLabel("Closed winner"),
		widget.NewLabel("•"),
		lossBtn,
		widget.NewLabel("="),
		widget.NewLabel("Closed loser"),
		widget.NewLabel("•"),
		scratchBtn,
		widget.NewLabel("="),
		widget.NewLabel("Break-even"),
		widget.NewLabel("•"),
		rolledBtn,
		widget.NewLabel("="),
		widget.NewLabel("Rolled position"),
	)

	info := widget.NewLabel("Click any ticker to view full trade details. Calendar shows 2 weeks back + 8 weeks forward for diversification tracking.")
	info.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		legend,
		info,
	)
}

// filterTrades applies additional client-side filters to trades
func filterTrades(trades []storage.TradeHistoryEntry, sectorFilter, strategyFilter, statusFilter string) []storage.TradeHistoryEntry {
	var filtered []storage.TradeHistoryEntry

	for _, trade := range trades {
		// Apply sector filter
		if sectorFilter != "" {
			matchesSector := false
			if trade.Bucket == sectorFilter {
				matchesSector = true
			} else if trade.Sector == sectorFilter {
				matchesSector = true
			}
			if !matchesSector {
				continue
			}
		}

		// Apply strategy filter
		if strategyFilter != "" {
			strategyMatch := false
			switch strategyFilter {
			case "Long Breakout":
				strategyMatch = trade.Strategy == storage.StrategyLongBreakout
			case "Short Breakout":
				strategyMatch = trade.Strategy == storage.StrategyShortBreakout
			case "Custom":
				strategyMatch = trade.Strategy == storage.StrategyCustom
			}
			if !strategyMatch {
				continue
			}
		}

		// Status filter already applied in database query

		filtered = append(filtered, trade)
	}

	return filtered
}
