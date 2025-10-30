package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildCalendarScreen(state *AppState) fyne.CanvasObject {
	// Title
	title := canvas.NewText("Calendar - 10-Week Sector View", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := widget.NewLabel("Rolling 10-week view (2 weeks back + 8 weeks forward)")
	subtitle.Wrapping = fyne.TextWrapWord

	// Calculate weeks
	now := time.Now()
	startWeek := now.AddDate(0, 0, -14) // 2 weeks back

	// Sector buckets (grouped correlated sectors for diversification tracking)
	buckets := []string{
		"Materials/Industrials",  // Basic Materials + Industrials
		"Tech/Comm",              // Communication Services + Technology
		"Financial/Cyclical",     // Financial + Consumer Cyclical
		"Defensive/Utilities",    // Consumer Defensive + Utilities
		"Energy",                 // Energy (standalone)
		"Healthcare",             // Healthcare (standalone)
		"Real Estate",            // Real Estate (standalone)
		"ETFs",                   // ETFs (standalone)
	}

	// Build calendar grid
	grid := container.NewVBox()

	// Header row (week dates)
	headerRow := container.NewHBox(widget.NewLabel("Sector"))
	for i := 0; i < 10; i++ {
		weekDate := startWeek.AddDate(0, 0, i*7)
		weekLabel := widget.NewLabel(weekDate.Format("Jan 2"))
		headerRow.Add(weekLabel)
	}
	grid.Add(headerRow)
	grid.Add(widget.NewSeparator())

	// Get positions grouped by week and bucket
	positions, _ := state.db.GetOpenPositions()
	positionMap := make(map[string]map[string][]string) // bucket -> week -> tickers

	// TODO: Load actual position entry dates from database
	// For now, show current positions in current week
	currentWeekStr := now.Format("Jan 2")
	for _, pos := range positions {
		if positionMap[pos.Bucket] == nil {
			positionMap[pos.Bucket] = make(map[string][]string)
		}
		positionMap[pos.Bucket][currentWeekStr] = append(positionMap[pos.Bucket][currentWeekStr], pos.Ticker)
	}

	// Rows for each bucket (with expandable rows if multiple positions per week)
	for _, bucket := range buckets {
		// Find max positions in any single week for this bucket
		maxPositionsInWeek := 0
		if positionMap[bucket] != nil {
			for weekStr := range positionMap[bucket] {
				if len(positionMap[bucket][weekStr]) > maxPositionsInWeek {
					maxPositionsInWeek = len(positionMap[bucket][weekStr])
				}
			}
		}

		// If no positions, show one empty row
		if maxPositionsInWeek == 0 {
			maxPositionsInWeek = 1
		}

		// Create a row for each position slot
		for rowIdx := 0; rowIdx < maxPositionsInWeek; rowIdx++ {
			var bucketLabel string
			if rowIdx == 0 {
				if maxPositionsInWeek == 1 {
					bucketLabel = bucket
				} else {
					bucketLabel = fmt.Sprintf("%s (1/%d)", bucket, maxPositionsInWeek)
				}
			} else {
				bucketLabel = fmt.Sprintf("%s (%d/%d)", bucket, rowIdx+1, maxPositionsInWeek)
			}

			row := container.NewHBox(widget.NewLabel(bucketLabel))

			for i := 0; i < 10; i++ {
				weekDate := startWeek.AddDate(0, 0, i*7)
				weekStr := weekDate.Format("Jan 2")

				// Check if this bucket has positions in this week
				cell := widget.NewLabel("")
				if positionMap[bucket] != nil && len(positionMap[bucket][weekStr]) > rowIdx {
					ticker := positionMap[bucket][weekStr][rowIdx]
					cell.SetText(ticker)
				} else {
					cell.SetText("-")
				}

				row.Add(cell)
			}

			grid.Add(row)
		}
	}

	// Legend
	legend := widget.NewLabel(`
Legend:
- Ticker symbols show positions entered that week in that sector
- Multiple rows per bucket = multiple positions in same sector (e.g., Tech/Comm 1/3, 2/3, 3/3)
- Goal: Avoid clustering trades in same sector/week
- Diversify across time AND sectors
`)
	legend.Wrapping = fyne.TextWrapWord

	refreshBtn := widget.NewButton("Refresh Calendar", func() {
		// TODO: Refresh calendar data
	})
	refreshBtn.Importance = widget.HighImportance

	// Layout
	content := container.NewVBox(
		container.NewPadded(title),
		subtitle,
		widget.NewSeparator(),
		grid,
		widget.NewSeparator(),
		legend,
		refreshBtn,
	)

	return container.NewScroll(content)
}
