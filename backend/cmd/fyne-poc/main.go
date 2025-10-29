package main

import (
	"fmt"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	// Import backend packages
	"github.com/yourusername/trading-engine/internal/storage"
)

func main() {
	// Initialize database
	dbPath := filepath.Join("../../", "trading.db")
	db, err := storage.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize database schema if needed
	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Fyne application
	myApp := app.New()
	myWindow := myApp.NewWindow("TF-Engine Settings POC")
	myWindow.Resize(fyne.NewSize(600, 400))

	// Create UI components
	title := widget.NewLabel("TF-Engine Settings (Backend Integration)")
	title.TextStyle.Bold = true

	equityLabel := widget.NewLabel("Equity: $0")
	riskLabel := widget.NewLabel("Risk %: 0.00%")
	portfolioCapLabel := widget.NewLabel("Portfolio Cap: 0.00%")
	bucketCapLabel := widget.NewLabel("Bucket Cap: 0.00%")

	refreshBtn := widget.NewButton("Load from Backend", func() {
		// Call backend to get settings
		settings, err := db.GetAllSettings()
		if err != nil {
			log.Printf("Error loading settings: %v", err)
			equityLabel.SetText("Equity: ERROR loading")
			return
		}

		// Update UI with real data from backend
		if equity, ok := settings["equity"]; ok {
			equityLabel.SetText(fmt.Sprintf("Equity: $%s", equity))
		} else {
			equityLabel.SetText("Equity: Not set")
		}
		if riskPct, ok := settings["risk_pct"]; ok {
			riskLabel.SetText(fmt.Sprintf("Risk %%: %s%%", riskPct))
		} else {
			riskLabel.SetText("Risk %: Not set")
		}
		if portfolioCap, ok := settings["portfolio_heat_cap"]; ok {
			portfolioCapLabel.SetText(fmt.Sprintf("Portfolio Cap: %s%%", portfolioCap))
		} else {
			portfolioCapLabel.SetText("Portfolio Cap: Not set")
		}
		if bucketCap, ok := settings["bucket_heat_cap"]; ok {
			bucketCapLabel.SetText(fmt.Sprintf("Bucket Cap: %s%%", bucketCap))
		} else {
			bucketCapLabel.SetText("Bucket Cap: Not set")
		}
		log.Println("✓ Data loaded from backend successfully")
	})

	updateBtn := widget.NewButton("Update", func() {
		log.Println("Update clicked (not implemented in POC)")
	})

	statusLabel := widget.NewLabel("Status: Ready (Click 'Load from Backend' to test)")

	// Layout
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		equityLabel,
		riskLabel,
		portfolioCapLabel,
		bucketCapLabel,
		widget.NewSeparator(),
		container.NewHBox(refreshBtn, updateBtn),
		widget.NewSeparator(),
		statusLabel,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
