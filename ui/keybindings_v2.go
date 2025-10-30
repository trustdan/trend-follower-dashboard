package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// VIMMode manages VIM-style keyboard navigation
type VIMMode struct {
	state          *AppState
	enabled        bool
	hintMode       bool
	hints          map[string]*widget.Button
	currentButtons []*widget.Button
	keyBuffer      string // For multi-key commands like 'gg', 'gt'
}

// NewVIMMode creates a new VIM mode handler
func NewVIMMode(state *AppState) *VIMMode {
	return &VIMMode{
		state:     state,
		enabled:   false, // Disabled by default, toggle with button
		hintMode:  false,
		hints:     make(map[string]*widget.Button),
		keyBuffer: "",
	}
}

// Toggle enables/disables VIM mode
func (v *VIMMode) Toggle() {
	v.enabled = !v.enabled
	if v.enabled {
		log.Println("VIM Mode: ENABLED")
	} else {
		log.Println("VIM Mode: DISABLED")
		v.exitHintMode()
	}
}

// IsEnabled returns whether VIM mode is active
func (v *VIMMode) IsEnabled() bool {
	return v.enabled
}

// HandleKeyDown processes keyboard events for VIM navigation
func (v *VIMMode) HandleKeyDown(ev *fyne.KeyEvent) {
	if !v.enabled {
		return // VIM mode disabled, ignore all keys
	}

	log.Printf("VIM: Key pressed: %v (buffer: '%s')", ev.Name, v.keyBuffer)

	// Special keys that work regardless of mode
	switch ev.Name {
	case fyne.KeyEscape:
		v.exitHintMode()
		v.keyBuffer = ""
		log.Println("VIM: ESC - cleared buffer and exited hint mode")
		return
	}

	// Check for ? key (help overlay)
	if string(ev.Name) == "?" {
		v.showHelpOverlay()
		return
	}

	// Handle hint mode
	if v.hintMode {
		v.handleHintModeKey(ev)
		return
	}

	// Handle multi-key commands
	v.keyBuffer += string(ev.Name)
	if v.handleCommand(v.keyBuffer) {
		v.keyBuffer = "" // Command executed, clear buffer
	} else if len(v.keyBuffer) > 2 {
		// Buffer too long, clear it
		log.Printf("VIM: Buffer too long, clearing: %s", v.keyBuffer)
		v.keyBuffer = ""
	}
}

// handleCommand processes VIM commands
func (v *VIMMode) handleCommand(cmd string) bool {
	log.Printf("VIM: Checking command: '%s'", cmd)

	switch cmd {
	// Link hints
	case "f", "F":
		log.Println("VIM: f/F - activating link hints")
		v.enterHintMode()
		return true

	// Scrolling
	case "j":
		log.Println("VIM: j - scroll down")
		// TODO: Implement scroll down
		return true
	case "k":
		log.Println("VIM: k - scroll up")
		// TODO: Implement scroll up
		return true
	case "h":
		log.Println("VIM: h - scroll left")
		// TODO: Implement scroll left
		return true
	case "l":
		log.Println("VIM: l - scroll right")
		// TODO: Implement scroll right
		return true
	case "d":
		log.Println("VIM: d - half-page down")
		// TODO: Implement half-page down
		return true
	case "u":
		log.Println("VIM: u - half-page up")
		// TODO: Implement half-page up
		return true

	// Jump to top/bottom
	case "gg":
		log.Println("VIM: gg - jump to top")
		// TODO: Implement jump to top
		return true
	case "G":
		log.Println("VIM: G - jump to bottom")
		// TODO: Implement jump to bottom
		return true

	// Find
	case "/":
		log.Println("VIM: / - open find")
		v.openFind()
		return true
	case "n":
		log.Println("VIM: n - next search result")
		// TODO: Implement next search
		return true
	case "N":
		log.Println("VIM: N - previous search result")
		// TODO: Implement previous search
		return true

	// Refresh
	case "r":
		log.Println("VIM: r - refresh view")
		v.refreshCurrentView()
		return true

	// Focus
	case "gi":
		log.Println("VIM: gi - focus first input")
		// TODO: Implement focus first input
		return true

	// History
	case "H":
		log.Println("VIM: H - navigate back")
		// TODO: Implement back
		return true
	case "L":
		log.Println("VIM: L - navigate forward")
		// TODO: Implement forward
		return true

	// Command palette
	case "o":
		log.Println("VIM: o - open command palette")
		v.openCommandPalette()
		return true
	case "T":
		log.Println("VIM: T - view switcher")
		v.openViewSwitcher()
		return true

	// Tabs (if AppTabs used)
	case "gt":
		log.Println("VIM: gt - next tab")
		// TODO: Implement next tab
		return true
	case "gT":
		log.Println("VIM: gT - previous tab")
		// TODO: Implement previous tab
		return true
	case "t":
		log.Println("VIM: t - new tab")
		// TODO: Implement new tab
		return true
	case "x":
		log.Println("VIM: x - close tab")
		// TODO: Implement close tab
		return true
	case "X":
		log.Println("VIM: X - reopen last tab")
		// TODO: Implement reopen tab
		return true
	case "g0":
		log.Println("VIM: g0 - first tab")
		// TODO: Implement first tab
		return true
	case "g$":
		log.Println("VIM: g$ - last tab")
		// TODO: Implement last tab
		return true
	}

	// Command not recognized yet (might be multi-key)
	return false
}

// handleHintModeKey processes keys in hint mode
func (v *VIMMode) handleHintModeKey(ev *fyne.KeyEvent) {
	key := strings.ToLower(string(ev.Name))
	log.Printf("VIM: Hint mode key: %s", key)

	if len(key) == 1 {
		v.activateHint(key)
	}
}

// enterHintMode activates link hints
func (v *VIMMode) enterHintMode() {
	log.Println("VIM: Entering hint mode...")
	v.hintMode = true
	v.keyBuffer = ""
	v.currentButtons = v.findAllButtons()
	log.Printf("VIM: Found %d buttons", len(v.currentButtons))
	v.generateHints()
	v.showHints()
	log.Println("VIM: Hint mode activated")
}

// exitHintMode deactivates link hints
func (v *VIMMode) exitHintMode() {
	if !v.hintMode {
		return
	}
	log.Println("VIM: Exiting hint mode...")
	v.hintMode = false
	v.hints = make(map[string]*widget.Button)
	v.currentButtons = []*widget.Button{}
	log.Println("VIM: Hint mode exited")
}

// findAllButtons recursively finds all buttons in the UI
func (v *VIMMode) findAllButtons() []*widget.Button {
	buttons := []*widget.Button{}
	content := v.state.window.Content()
	v.findButtonsRecursive(content, &buttons)
	return buttons
}

// findButtonsRecursive recursively searches for buttons
func (v *VIMMode) findButtonsRecursive(obj fyne.CanvasObject, buttons *[]*widget.Button) {
	if btn, ok := obj.(*widget.Button); ok {
		*buttons = append(*buttons, btn)
		return
	}

	if c, ok := obj.(*fyne.Container); ok {
		for _, child := range c.Objects {
			v.findButtonsRecursive(child, buttons)
		}
	}
}

// generateHints creates keyboard shortcuts for all buttons
func (v *VIMMode) generateHints() {
	letters := "abcdefghijklmnopqrstuvwxyz"
	hints := []string{}

	// Single letters
	for _, c := range letters {
		hints = append(hints, string(c))
		if len(hints) >= len(v.currentButtons) {
			break
		}
	}

	// Double letters if needed
	if len(hints) < len(v.currentButtons) {
		for _, c1 := range letters {
			for _, c2 := range letters {
				hints = append(hints, string(c1)+string(c2))
				if len(hints) >= len(v.currentButtons) {
					break
				}
			}
			if len(hints) >= len(v.currentButtons) {
				break
			}
		}
	}

	// Assign hints to buttons
	v.hints = make(map[string]*widget.Button)
	for i, btn := range v.currentButtons {
		if i >= len(hints) {
			break
		}
		v.hints[hints[i]] = btn
	}
}

// showHints displays available hints in log
func (v *VIMMode) showHints() {
	log.Println("===============================")
	log.Println("VIM HINT MODE ACTIVE")
	log.Println("===============================")
	for hint, btn := range v.hints {
		log.Printf("  [%s] -> %s", hint, btn.Text)
	}
	log.Println("Press ESC to exit hint mode")
	log.Println("===============================")
}

// activateHint triggers the button associated with a hint
func (v *VIMMode) activateHint(key string) {
	log.Printf("VIM: Looking for hint key: %s", key)
	if btn, exists := v.hints[key]; exists {
		log.Printf("VIM: Activating button: %s", btn.Text)
		btn.OnTapped()
		v.exitHintMode()
	} else {
		log.Printf("VIM: No button found for hint: %s", key)
	}
}

// showHelpOverlay shows the VIM keybindings help
func (v *VIMMode) showHelpOverlay() {
	helpText := `# VIM Mode Keybindings

## Link Hints
- **f / F** - Show link hints (press letter to activate)

## Scrolling
- **j / k** - Scroll down / up
- **h / l** - Scroll left / right
- **d / u** - Half-page down / up
- **gg / G** - Jump to top / bottom

## Find
- **/** - Open find bar
- **n / N** - Next / previous match

## Actions
- **r** - Refresh current view
- **gi** - Focus first input
- **?** - Show this help

## Navigation
- **H / L** - Back / forward
- **o** - Command palette
- **T** - View switcher

## Tabs (if available)
- **gt / gT** - Next / previous tab
- **t** - New tab
- **x / X** - Close / reopen tab
- **g0 / g$** - First / last tab

## General
- **ESC** - Exit overlays / hint mode`

	helpContent := widget.NewRichTextFromMarkdown(helpText)
	scrollContainer := container.NewScroll(helpContent)
	scrollContainer.SetMinSize(fyne.NewSize(600, 500))

	dialog.ShowCustom(
		"VIM Mode Keybindings",
		"Close",
		scrollContainer,
		v.state.window,
	)
}

// openFind opens a find dialog
func (v *VIMMode) openFind() {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search...")

	d := dialog.NewCustom(
		"Find",
		"Close",
		container.NewVBox(
			widget.NewLabel("Search for:"),
			searchEntry,
		),
		v.state.window,
	)

	searchEntry.OnSubmitted = func(query string) {
		log.Printf("VIM: Searching for: %s", query)
		// TODO: Implement actual search
		d.Hide()
	}

	d.Show()
	v.state.window.Canvas().Focus(searchEntry)
}

// refreshCurrentView refreshes the current view
func (v *VIMMode) refreshCurrentView() {
	// TODO: Implement refresh logic
	log.Println("VIM: Refreshing current view (not yet implemented)")
}

// openCommandPalette opens a command palette
func (v *VIMMode) openCommandPalette() {
	commands := []string{
		"Dashboard",
		"Scanner",
		"Checklist",
		"Position Sizing",
		"Heat Check",
		"Trade Entry",
		"Calendar",
	}

	list := widget.NewList(
		func() int { return len(commands) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(commands[id])
		},
	)

	d := dialog.NewCustom(
		"Command Palette",
		"Cancel",
		container.NewVBox(
			widget.NewLabel("Select a view:"),
			list,
		),
		v.state.window,
	)

	list.OnSelected = func(id widget.ListItemID) {
		log.Printf("VIM: Selected command: %s", commands[id])
		// TODO: Navigate to selected view
		d.Hide()
	}

	d.Show()
}

// openViewSwitcher opens a view switcher
func (v *VIMMode) openViewSwitcher() {
	v.openCommandPalette() // Same as command palette for now
}

// AttachToWindow attaches the VIM handler to the window
func (v *VIMMode) AttachToWindow() {
	canv := v.state.window.Canvas()
	canv.SetOnTypedKey(func(ev *fyne.KeyEvent) {
		v.HandleKeyDown(ev)
	})
	log.Println("VIM Mode: Handler attached (disabled by default)")
}
