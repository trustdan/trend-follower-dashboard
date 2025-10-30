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
)

func buildDashboardScreen(state *AppState) fyne.CanvasObject {
	log.Println("buildDashboardScreen: Starting...")

	// Title
	title := canvas.NewText("Dashboard", nil)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	log.Println("buildDashboardScreen: Title created")

	// Account Settings Card
	log.Println("buildDashboardScreen: Building settings card...")
	settingsCard := buildSettingsCard(state)
	log.Println("buildDashboardScreen: Settings card built")

	// Open Positions Card
	log.Println("buildDashboardScreen: Building positions card...")
	positionsCard := buildPositionsCard(state)
	log.Println("buildDashboardScreen: Positions card built")

	// Heat Status Card
	log.Println("buildDashboardScreen: Building heat card...")
	heatCard := buildHeatCard(state)
	log.Println("buildDashboardScreen: Heat card built")

	// Today's Candidates Card
	log.Println("buildDashboardScreen: Building candidates card...")
	candidatesCard := buildCandidatesCard(state)
	log.Println("buildDashboardScreen: Candidates card built")

	// Layout: 2x2 grid of cards
	topRow := container.NewGridWithColumns(2,
		settingsCard,
		heatCard,
	)

	bottomRow := container.NewGridWithColumns(2,
		positionsCard,
		candidatesCard,
	)

	content := container.NewVBox(
		container.NewPadded(title),
		topRow,
		bottomRow,
		layout.NewSpacer(),
	)

	return container.NewScroll(content)
}

func buildSettingsCard(state *AppState) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Account Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Load settings from database
	settings, err := state.db.GetAllSettings()
	if err != nil {
		return container.NewVBox(
			title,
			widget.NewLabel("Error loading settings"),
		)
	}

	// Extract settings with defaults for nil values
	equity := getSettingWithDefault(settings, "equity", "0")
	riskPct := getSettingWithDefault(settings, "risk_pct", "0")
	portfolioCap := getSettingWithDefault(settings, "portfolio_heat_cap", "0")
	bucketCap := getSettingWithDefault(settings, "bucket_heat_cap", "0")

	equityLabel := widget.NewLabel(fmt.Sprintf("Equity: $%s", equity))
	riskLabel := widget.NewLabel(fmt.Sprintf("Risk per Trade: %s%%", riskPct))
	portfolioCapLabel := widget.NewLabel(fmt.Sprintf("Portfolio Heat Cap: %s%%", portfolioCap))
	bucketCapLabel := widget.NewLabel(fmt.Sprintf("Bucket Heat Cap: %s%%", bucketCap))

	editBtn := widget.NewButton("Edit Settings", func() {
		showSettingsDialog(state)
	})
	editBtn.Importance = widget.HighImportance

	card := container.NewVBox(
		title,
		widget.NewSeparator(),
		equityLabel,
		riskLabel,
		portfolioCapLabel,
		bucketCapLabel,
		editBtn,
	)

	return container.NewPadded(card)
}

func buildPositionsCard(state *AppState) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Open Positions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Load open positions from database
	positions, err := state.db.GetOpenPositions()
	if err != nil {
		return container.NewVBox(
			title,
			widget.NewLabel("Error loading positions"),
		)
	}

	if len(positions) == 0 {
		return container.NewPadded(
			container.NewVBox(
				title,
				widget.NewSeparator(),
				widget.NewLabel("No open positions"),
			),
		)
	}

	// Create table of positions
	positionsList := container.NewVBox()
	for _, pos := range positions {
		posLabel := widget.NewLabel(fmt.Sprintf("%s - $%.2f risk - %s", pos.Ticker, pos.RiskDollars, pos.Bucket))
		positionsList.Add(posLabel)
	}

	card := container.NewVBox(
		title,
		widget.NewSeparator(),
		positionsList,
	)

	return container.NewPadded(card)
}

func buildHeatCard(state *AppState) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Heat Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Get settings for caps
	settings, err := state.db.GetAllSettings()
	if err != nil {
		return container.NewVBox(
			title,
			widget.NewLabel("Error loading heat data"),
		)
	}

	// Calculate total portfolio heat
	positions, err := state.db.GetOpenPositions()
	if err != nil {
		return container.NewVBox(
			title,
			widget.NewLabel("Error loading positions"),
		)
	}

	var totalHeat float64
	for _, pos := range positions {
		totalHeat += pos.RiskDollars
	}

	// Parse equity and portfolio cap with defaults
	equity := getSettingWithDefault(settings, "equity", "0")
	portfolioCap := getSettingWithDefault(settings, "portfolio_heat_cap", "0")

	// Calculate heat percentage
	// TODO: Proper parsing and calculation
	heatLabel := widget.NewLabel(fmt.Sprintf("Portfolio Heat: $%.2f", totalHeat))
	capLabel := widget.NewLabel(fmt.Sprintf("Cap: %s%% of $%s", portfolioCap, equity))

	// Create progress bar
	heatProgress := widget.NewProgressBar()
	heatProgress.SetValue(0.5) // TODO: Calculate actual percentage

	card := container.NewVBox(
		title,
		widget.NewSeparator(),
		heatLabel,
		capLabel,
		heatProgress,
	)

	return container.NewPadded(card)
}

func buildCandidatesCard(state *AppState) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Today's Candidates", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Load today's candidates (use today's date)
	today := time.Now().Format("2006-01-02")
	candidates, err := state.db.GetCandidatesForDate(today)
	if err != nil {
		return container.NewVBox(
			title,
			widget.NewLabel("Error loading candidates"),
		)
	}

	if len(candidates) == 0 {
		return container.NewPadded(
			container.NewVBox(
				title,
				widget.NewSeparator(),
				widget.NewLabel("No candidates found"),
				widget.NewLabel("Use Scanner to import from FINVIZ"),
			),
		)
	}

	// Create list of candidates
	candidatesList := container.NewVBox()
	for _, cand := range candidates {
		// Safe type assertion with defaults
		ticker, ok := cand["ticker"].(string)
		if !ok || ticker == "" {
			ticker = "UNKNOWN"
		}
		scanDate, ok := cand["scan_date"].(string)
		if !ok || scanDate == "" {
			scanDate = time.Now().Format("2006-01-02")
		}
		candLabel := widget.NewLabel(fmt.Sprintf("%s - %s", ticker, scanDate))
		candidatesList.Add(candLabel)
	}

	refreshBtn := widget.NewButton("Refresh", func() {
		// TODO: Refresh candidates
	})
	refreshBtn.Importance = widget.HighImportance

	// Create scroll container with minimum height for better visibility
	scroll := container.NewScroll(candidatesList)
	scroll.SetMinSize(fyne.NewSize(200, 300))  // Min width 200, height 300 pixels

	card := container.NewVBox(
		title,
		widget.NewSeparator(),
		scroll,
		refreshBtn,
	)

	return container.NewPadded(card)
}

// showSettingsDialog displays a dialog for editing account settings
func showSettingsDialog(state *AppState) {
	// Load current settings
	settings, err := state.db.GetAllSettings()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load settings: %w", err), state.window)
		return
	}

	// Create form entries with defaults
	equityEntry := widget.NewEntry()
	equityEntry.SetText(getSettingWithDefault(settings, "equity", "100000"))
	equityEntry.SetPlaceHolder("100000")

	riskPctEntry := widget.NewEntry()
	riskPctEntry.SetText(getSettingWithDefault(settings, "risk_pct", "0.75"))
	riskPctEntry.SetPlaceHolder("0.75")

	portfolioCapEntry := widget.NewEntry()
	portfolioCapEntry.SetText(getSettingWithDefault(settings, "portfolio_heat_cap", "4.0"))
	portfolioCapEntry.SetPlaceHolder("4.0")

	bucketCapEntry := widget.NewEntry()
	bucketCapEntry.SetText(getSettingWithDefault(settings, "bucket_heat_cap", "1.5"))
	bucketCapEntry.SetPlaceHolder("1.5")

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Account Equity ($)", Widget: equityEntry, HintText: "Total account value in dollars"},
			{Text: "Risk per Trade (%)", Widget: riskPctEntry, HintText: "Percentage of equity to risk per trade (e.g., 0.75)"},
			{Text: "Portfolio Heat Cap (%)", Widget: portfolioCapEntry, HintText: "Maximum total risk across all positions (e.g., 4.0)"},
			{Text: "Bucket Heat Cap (%)", Widget: bucketCapEntry, HintText: "Maximum risk per sector bucket (e.g., 1.5)"},
		},
		OnSubmit: func() {
			// Save to database
			if err := state.db.SetSetting("equity", equityEntry.Text); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save equity: %w", err), state.window)
				return
			}
			if err := state.db.SetSetting("risk_pct", riskPctEntry.Text); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save risk_pct: %w", err), state.window)
				return
			}
			if err := state.db.SetSetting("portfolio_heat_cap", portfolioCapEntry.Text); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save portfolio_heat_cap: %w", err), state.window)
				return
			}
			if err := state.db.SetSetting("bucket_heat_cap", bucketCapEntry.Text); err != nil {
				dialog.ShowError(fmt.Errorf("failed to save bucket_heat_cap: %w", err), state.window)
				return
			}

			// Show success message
			dialog.ShowInformation("Success", "Settings saved successfully. Refresh dashboard to see changes.", state.window)
		},
	}

	// Show form dialog
	dialog.ShowForm("Edit Settings", "Save", "Cancel", form.Items, func(submitted bool) {
		if submitted {
			form.OnSubmit()
		}
	}, state.window)
}
