package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/trading-engine/internal/storage"
)

func main() {
	// Setup logging to file
	logFile, err := os.OpenFile("tf-gui.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Recover from panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %v", r)
			fmt.Fprintf(os.Stderr, "Application crashed: %v\nCheck tf-gui.log for details\n", r)
		}
	}()

	log.Println("========== TF-Engine GUI Starting ==========")
	log.Printf("Working directory: %s", getWorkingDir())

	// Initialize database
	dbPath := filepath.Join(".", "trading.db")
	log.Printf("Database path: %s", dbPath)

	db, err := storage.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	log.Println("Database opened successfully")

	// Initialize database schema if needed
	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Create application
	log.Println("Creating Fyne application...")
	myApp := app.NewWithID("com.tfengine.dashboard")
	log.Println("Setting theme...")
	myApp.Settings().SetTheme(&tfTheme{variant: ThemeLight})

	// Create main window
	log.Println("Creating main window...")
	mainWindow := myApp.NewWindow("TF-Engine - Trend Following Dashboard")
	mainWindow.Resize(fyne.NewSize(1200, 800))

	// Create app state
	log.Println("Creating app state...")
	appState := &AppState{
		db:         db,
		window:     mainWindow,
		isDarkMode: false,
		myApp:      myApp,
	}

	// Build UI
	log.Println("Building UI...")
	content := buildMainUI(appState)
	mainWindow.SetContent(content)
	log.Println("UI built successfully")

	// Show window first
	log.Println("Showing window...")
	mainWindow.Show()

	// Attach VIM mode keyboard handlers AFTER window is shown
	// (Canvas must be initialized first for handlers to register properly)
	if appState.vimMode != nil {
		log.Println("Attaching VIM mode handlers to window canvas...")
		appState.vimMode.AttachToWindow()
	}

	// Always show welcome dialog on first run
	// (User can disable with "don't show again" checkbox)
	log.Println("Checking first run...")
	if isFirstRun(db) {
		log.Println("First run detected, showing welcome dialog...")
		showWelcomeDialog(appState)
		setNotFirstRun(db)
	} else {
		log.Println("Not first run (use Welcome button to show dialog manually)")
	}

	// Run app
	log.Println("Starting application event loop...")
	myApp.Run()
	log.Println("Application exited normally")
}

// AppState holds the application state
type AppState struct {
	db                     *storage.DB
	window                 fyne.Window
	isDarkMode             bool
	myApp                  fyne.App
	currentSession         *storage.TradeSession
	sessionChangeCallbacks []func(*storage.TradeSession)
	vimMode                *VIMMode // VIM mode handler (attached after window shown)
	navigateToTab          func(int) // Navigation callback for switching tabs
	sampleMode             bool      // Whether sample data mode is enabled
	sampleSession          *storage.TradeSession // Sample session when in sample mode
}

// SetCurrentSession updates the current session and triggers all callbacks
func (s *AppState) SetCurrentSession(session *storage.TradeSession) {
	s.currentSession = session
	for _, callback := range s.sessionChangeCallbacks {
		callback(session)
	}
}

// RegisterSessionChangeCallback adds a callback to be called when session changes
func (s *AppState) RegisterSessionChangeCallback(callback func(*storage.TradeSession)) {
	s.sessionChangeCallbacks = append(s.sessionChangeCallbacks, callback)
}

// GetActiveSession returns the current session (real or sample depending on mode)
func (s *AppState) GetActiveSession() *storage.TradeSession {
	if s.sampleMode && s.sampleSession != nil {
		return s.sampleSession
	}
	return s.currentSession
}

// EnableSampleMode activates sample data mode
func (s *AppState) EnableSampleMode() {
	s.sampleMode = true
	if s.sampleSession == nil {
		s.sampleSession = CreateSampleSession()
	}
	// Trigger callbacks to refresh UI
	for _, callback := range s.sessionChangeCallbacks {
		callback(s.sampleSession)
	}
}

// DisableSampleMode deactivates sample data mode
func (s *AppState) DisableSampleMode() {
	s.sampleMode = false
	// Trigger callbacks to refresh UI with real session (or nil)
	for _, callback := range s.sessionChangeCallbacks {
		callback(s.currentSession)
	}
}

// buildMainUI constructs the main application UI with navigation
func buildMainUI(state *AppState) fyne.CanvasObject {
	// Create theme toggle button (declare var first for closure)
	var themeToggleBtn *widget.Button
	themeToggleBtn = widget.NewButton("🌙 Dark Mode", func() {
		state.isDarkMode = !state.isDarkMode
		if state.isDarkMode {
			state.myApp.Settings().SetTheme(&tfTheme{variant: ThemeDark})
			themeToggleBtn.SetText("☀️ Light Mode")
		} else {
			state.myApp.Settings().SetTheme(&tfTheme{variant: ThemeLight})
			themeToggleBtn.SetText("🌙 Dark Mode")
		}
	})
	themeToggleBtn.Importance = widget.HighImportance

	// Create help button
	helpBtn := widget.NewButton("❓ Help", func() {
		showHelpDialog(state)
	})
	helpBtn.Importance = widget.HighImportance

	// Create VIM mode toggle button
	vimBtn := widget.NewButton("VIM: On", func() {
		// Placeholder; actual handler assigned after vimMode is created
	})
	vimBtn.Importance = widget.HighImportance

	// Create "Show Welcome" button
	welcomeBtn := widget.NewButton("👋 Welcome", func() {
		showWelcomeDialog(state)
	})
	welcomeBtn.Importance = widget.HighImportance

	// Create Sample Mode toggle button
	var sampleBtn *widget.Button
	sampleBtn = widget.NewButton("📦 Sample Data: Off", func() {
		if state.sampleMode {
			state.DisableSampleMode()
			sampleBtn.SetText("📦 Sample Data: Off")
		} else {
			state.EnableSampleMode()
			sampleBtn.SetText("📦 Sample Data: On")
		}
		sampleBtn.Refresh()
	})
	sampleBtn.Importance = widget.HighImportance

	// Navigation buttons
	navItems := []struct {
		title string
		build func(*AppState) fyne.CanvasObject
	}{
		{"📊 Dashboard", buildDashboardScreen},
		{"🔍 Scanner", buildScannerScreen},
		{"✅ Checklist", buildChecklistScreen},
		{"📐 Position Sizing", buildPositionSizingScreen},
		{"🔥 Heat Check", buildHeatCheckScreen},
		{"💰 Trade Entry", buildTradeEntryScreen},
		{"📅 Calendar", buildCalendarScreen},
		{"📜 Session History", buildSessionHistoryScreen},
	}

	navButtons := make([]*widget.Button, len(navItems))
	contentArea := container.NewStack()

	updateNavSelection := func(idx int) {
		if idx < 0 || idx >= len(navItems) {
			return
		}

		contentArea.Objects = []fyne.CanvasObject{navItems[idx].build(state)}
		contentArea.Refresh()

		for i, btn := range navButtons {
			if btn == nil {
				continue
			}
			if i == idx {
				btn.Importance = widget.HighImportance
			} else {
				btn.Importance = widget.LowImportance
			}
			btn.Refresh()
		}
	}

	for i, item := range navItems {
		idx := i
		btn := widget.NewButton(item.title, func() {
			updateNavSelection(idx)
		})
		btn.Alignment = widget.ButtonAlignLeading
		btn.Importance = widget.LowImportance
		navButtons[i] = btn
	}

	// Convert navButtons to []fyne.CanvasObject for VBox
	navObjects := make([]fyne.CanvasObject, len(navButtons))
	for i, btn := range navButtons {
		navObjects[i] = btn
	}
	navContainer := container.NewVBox(navObjects...)
	updateNavSelection(0)

	// Store navigation function in AppState so screens can navigate programmatically
	state.navigateToTab = updateNavSelection

	// Initialize VIM mode and wire up the toggle button
	vimMode := NewVIMMode(state)
	state.vimMode = vimMode // Store in AppState for later attachment
	updateVIMButton := func() {
		if vimMode.IsEnabled() {
			vimBtn.SetText("VIM: On")
		} else {
			vimBtn.SetText("VIM: Off")
		}
	}
	vimBtn.OnTapped = func() {
		vimMode.Toggle()
		updateVIMButton()
	}
	updateVIMButton()
	// NOTE: AttachToWindow() will be called in main() AFTER window.Show()

	// Setup keyboard shortcuts
	setupKeyboardShortcuts(state, state.window, func() {
		showNewTradeDialog(state, updateNavSelection)
	}, func() {
		showResumeSessionMenu(state, updateNavSelection)
	}, func() {
		updateNavSelection(7) // Session History is tab index 7
	})

	// Create session bar
	sessionBar := NewSessionBar(state)

	// Create session control buttons
	newTradeBtn := widget.NewButton("Start New Trade", func() {
		showNewTradeDialog(state, updateNavSelection)
	})
	newTradeBtn.Importance = widget.HighImportance

	resumeSessionBtn := createResumeSessionButton(state, updateNavSelection)

	// Session controls row
	sessionControls := container.NewHBox(
		newTradeBtn,
		resumeSessionBtn,
		layout.NewSpacer(),
		sampleBtn,
		themeToggleBtn,
		helpBtn,
		vimBtn,
		welcomeBtn,
	)

	// Create split container with navigation on left, content on right
	split := container.NewHSplit(
		container.NewBorder(
			container.NewVBox(
				widget.NewLabel("TF-Engine"),
			),
			nil, nil, nil,
			container.NewVScroll(navContainer),
		),
		contentArea,
	)
	split.SetOffset(0.2) // 20% for navigation, 80% for content

	// Main layout with session bar at top
	mainLayout := container.NewBorder(
		container.NewVBox(
			sessionBar,
			widget.NewSeparator(),
			sessionControls,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		split,
	)

	return mainLayout
}

// isFirstRun checks if this is the first time the app is being run
func isFirstRun(db *storage.DB) bool {
	settings, err := db.GetAllSettings()
	if err != nil {
		return true
	}

	// Check if user disabled the welcome dialog
	showWelcome, exists := settings["show_welcome_dialog"]
	if exists && showWelcome == "false" {
		return false
	}

	// Check if this is truly the first run
	firstRun, exists := settings["first_run"]
	return !exists || firstRun != "false"
}

// setNotFirstRun marks the app as having been run
func setNotFirstRun(db *storage.DB) {
	db.SetSetting("first_run", "false")
}

// showWelcomeDialog displays the first-run welcome message
func showWelcomeDialog(state *AppState) {
	welcomeContent := widget.NewRichTextFromMarkdown(`# Welcome to TF-Engine

This is a **discipline enforcement system** for trend-following trading.

## Quick Start:

1. **Dashboard** - Set your account size and risk parameters
2. **Scanner** - Import FINVIZ candidates daily
3. **Checklist** - Evaluate trades (RED/YELLOW/GREEN banner)
4. **Position Sizing** - Calculate shares/contracts
5. **Heat Check** - Verify portfolio limits
6. **Trade Entry** - Final gates check before GO/NO-GO

## Philosophy:

This system **prevents impulsive trading** through:
- 5 Hard Gates (cannot bypass)
- Heat caps (no concentration)
- 2-minute cooloff
- Mechanical exits

**The value is in what it prevents, not what it allows.**

## Keyboard Shortcuts:

Press **F** to show link hints for keyboard navigation (VIM-style).
Use **J/K** for down/up, **D/U** for page down/up.

See docs/anti-impulsivity.md for full details.`)

	// Create "don't show again" checkbox
	dontShowAgain := widget.NewCheck("Don't show this message again", nil)

	// Create custom dialog with checkbox
	scrollWelcome := container.NewScroll(welcomeContent)
	scrollWelcome.SetMinSize(fyne.NewSize(700, 500))

	content := container.NewVBox(
		scrollWelcome,
		widget.NewSeparator(),
		dontShowAgain,
	)

	d := dialog.NewCustom(
		"Welcome to TF-Engine!",
		"Get Started",
		content,
		state.window,
	)

	// Override the dismiss callback to save the preference
	d.SetOnClosed(func() {
		if dontShowAgain.Checked {
			state.db.SetSetting("show_welcome_dialog", "false")
		}
	})

	d.Show()
}

// getWorkingDir returns the current working directory
func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}

// showHelpDialog displays the help/FAQ dialog
func showHelpDialog(state *AppState) {
	helpContent := widget.NewRichTextFromMarkdown(`# TF-Engine Help & FAQ

## Quick Navigation

Use the menu on the left to navigate between screens:
- **Dashboard** - View positions and configure settings
- **Scanner** - Import FINVIZ candidates
- **Checklist** - Evaluate trades (5 gates + quality items)
- **Position Sizing** - Calculate shares/contracts
- **Heat Check** - Verify portfolio and bucket heat limits
- **Trade Entry** - Final GO/NO-GO decision
- **Calendar** - 10-week sector diversification grid

## Keyboard Shortcuts

- **F** - Show link hints for keyboard navigation (VIM-style)
- **J** - Navigate down
- **K** - Navigate up
- **D** - Page down
- **U** - Page up
- **ESC** - Cancel/dismiss hints

## Core Philosophy

This is a **discipline enforcement system**, not a flexible trading platform.

### The 5 Hard Gates (Cannot Bypass)

1. **Signal** - 55-bar Donchian breakout (long > 55-high / short < 55-low)
2. **Risk/Size** - ATR-based position sizing (2×N stop distance)
3. **Options/Liquidity** - 60-90 DTE, liquid options only
4. **Exits** - Exit by 10-bar opposite Donchian OR 2×N (whichever closer)
5. **Behavior** - 2-minute cooloff + no intraday overrides

### Heat Management

- **Portfolio Cap** - Maximum 4% of equity at risk across all positions
- **Bucket Cap** - Maximum 1.5% of equity at risk in one sector
- Prevents concentration risk and forced diversification

### Banner States

- **GREEN** - All required gates passed, quality score good → OK to trade
- **YELLOW** - Required gates passed, quality score low → Caution
- **RED** - Missing required gates → DO NOT TRADE

## Common Questions

**Q: Can I bypass the 2-minute cooloff?**
A: No. This is intentional friction to prevent impulsive trades.

**Q: Can I increase the heat caps?**
A: You can adjust settings, but the recommended caps are 4% portfolio / 1.5% bucket.

**Q: What if I want to add more than the system allows?**
A: That's the point. The system prevents overconcentration and impulsive sizing.

**Q: Can I trade without running the checklist?**
A: No. The banner must be GREEN to proceed to Trade Entry.

**Q: How do I re-show the welcome message?**
A: Delete the "show_welcome_dialog" setting from the database settings table.

## Position Sizing Methods

1. **Stock/ETF** - Direct stock purchase with ATR-based stops
2. **Options (Delta-ATR)** - Delta-adjusted ATR for options
3. **Options (Contracts)** - Contract-based risk (premium × 100)

## File Locations

- **Database**: trading.db (same directory as executable)
- **Logs**: tf-gui.log (same directory as executable)

## For More Information

See the documentation files:
- **docs/anti-impulsivity.md** - Core design philosophy
- **docs/PROJECT_STATUS.md** - Current implementation status
- **CLAUDE.md** - Development instructions
- **README.md** - Project overview

---

**Remember:** The value is in what this system prevents, not what it allows.`)

	// Create a sized container for the help content
	scrollContainer := container.NewScroll(helpContent)
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))

	dialog.ShowCustom(
		"Help & FAQ",
		"Close",
		scrollContainer,
		state.window,
	)
}
