package main

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// clickableTarget represents any clickable UI element
type clickableTarget struct {
	obj      fyne.CanvasObject
	activate func() // Function to call when activated
	label    string // Display name for logging
}

// VIMMode manages VIM-style keyboard navigation
type VIMMode struct {
	state            *AppState
	enabled          bool
	hintMode         bool
	hints            map[string]*clickableTarget // Maps hint keys to clickable targets
	keyBuffer        string                      // For multi-key commands like 'gg', 'gt'
	hintOverlay      *fyne.Container
	hintLabels       []*canvas.Text // Individual hint labels positioned over buttons
	scrollContainers []*container.Scroll
}

// NewVIMMode creates a new VIM mode handler
func NewVIMMode(state *AppState) *VIMMode {
	return &VIMMode{
		state:     state,
		enabled:   false, // Disabled by default, toggle with button
		hintMode:  false,
		hints:     make(map[string]*clickableTarget),
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
	}
}

// HandleTypedRune processes printable keyboard events for VIM navigation
func (v *VIMMode) HandleTypedRune(r rune) {
	if !v.enabled {
		return // VIM mode disabled, ignore all keys
	}

	key := string(r)
	log.Printf("VIM: Rune typed: %s (buffer: '%s')", key, v.keyBuffer)

	// Check for ? key (help overlay)
	if key == "?" {
		v.showHelpOverlay()
		v.keyBuffer = ""
		return
	}

	// Handle hint mode
	if v.hintMode {
		v.handleHintModeKey(key)
		return
	}

	// Handle multi-key commands
	v.keyBuffer += key
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
		v.scrollVertical(50) // Scroll down 50 pixels
		return true
	case "k":
		log.Println("VIM: k - scroll up")
		v.scrollVertical(-50) // Scroll up 50 pixels
		return true
	case "h":
		log.Println("VIM: h - scroll left")
		v.scrollHorizontal(-50) // Scroll left 50 pixels
		return true
	case "l":
		log.Println("VIM: l - scroll right")
		v.scrollHorizontal(50) // Scroll right 50 pixels
		return true
	case "d":
		log.Println("VIM: d - half-page down")
		v.scrollVertical(300) // Half page down (approx)
		return true
	case "u":
		log.Println("VIM: u - half-page up")
		v.scrollVertical(-300) // Half page up (approx)
		return true

	// Jump to top/bottom
	case "gg":
		log.Println("VIM: gg - jump to top")
		v.scrollToTop()
		return true
	case "G":
		log.Println("VIM: G - jump to bottom")
		v.scrollToBottom()
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
func (v *VIMMode) handleHintModeKey(key string) {
	lowerKey := strings.ToLower(key)
	log.Printf("VIM: Hint mode key: %s", lowerKey)

	if len(lowerKey) == 1 {
		v.activateHint(lowerKey)
	}
}

// enterHintMode activates link hints
func (v *VIMMode) enterHintMode() {
	log.Println("VIM: Entering hint mode...")
	v.hintMode = true
	v.keyBuffer = ""

	// Find all clickable elements (buttons, list items, etc.)
	v.findAllClickableElements()
	log.Printf("VIM: Found %d clickable elements", len(v.hints))

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
	v.hints = make(map[string]*clickableTarget)
	v.hideHintOverlay()
	log.Println("VIM: Hint mode exited")
}

// findAllClickableElements finds all clickable elements and generates hints
func (v *VIMMode) findAllClickableElements() {
	v.hints = make(map[string]*clickableTarget)
	v.scrollContainers = []*container.Scroll{}

	content := v.state.window.Content()
	v.findClickableRecursive(content)

	log.Printf("VIM: Found %d scroll containers", len(v.scrollContainers))

	// Generate hint keys for all targets
	v.generateHintKeys()
}

// findAllButtons recursively finds all buttons in the UI (for scroll operations)
func (v *VIMMode) findAllButtons() []*widget.Button {
	buttons := []*widget.Button{}
	v.scrollContainers = []*container.Scroll{} // Reset scroll containers
	content := v.state.window.Content()
	v.findButtonsRecursive(content, &buttons)
	log.Printf("VIM: Found %d scroll containers", len(v.scrollContainers))
	return buttons
}

// findClickableRecursive recursively finds all clickable elements
func (v *VIMMode) findClickableRecursive(obj fyne.CanvasObject) {
	// Check if this is a button
	if btn, ok := obj.(*widget.Button); ok {
		// Create a target for this button
		target := &clickableTarget{
			obj:   btn,
			label: btn.Text,
			activate: func() {
				if btn.OnTapped != nil {
					btn.OnTapped()
				}
			},
		}
		// We'll assign hint keys later in generateHintKeys
		// For now, just store with temporary key
		tempKey := fmt.Sprintf("btn_%p", btn)
		v.hints[tempKey] = target
		log.Printf("VIM: Found button: %s", btn.Text)
	}

	// Check if this is a list (navigation menu)
	if list, ok := obj.(*widget.List); ok {
		// Get the number of items
		itemCount := list.Length()
		log.Printf("VIM: Found List with %d items", itemCount)

		// Create targets for each list item
		for i := 0; i < itemCount; i++ {
			itemID := i // Capture for closure
			target := &clickableTarget{
				obj:   list, // Use the list as the object
				label: fmt.Sprintf("Nav Item %d", itemID),
				activate: func() {
					log.Printf("VIM: Activating list item %d", itemID)
					if list.OnSelected != nil {
						list.OnSelected(widget.ListItemID(itemID))
					}
				},
			}
			tempKey := fmt.Sprintf("list_%d", itemID)
			v.hints[tempKey] = target
			log.Printf("VIM: Added list item %d as clickable target", itemID)
		}
	}

	// Check all container types
	switch c := obj.(type) {
	case *fyne.Container:
		for _, child := range c.Objects {
			v.findClickableRecursive(child)
		}
	case *container.Scroll:
		// Track this scroll container for j/k/h/l scrolling
		v.scrollContainers = append(v.scrollContainers, c)
		if c.Content != nil {
			v.findClickableRecursive(c.Content)
		}
	case *container.Split:
		if c.Leading != nil {
			v.findClickableRecursive(c.Leading)
		}
		if c.Trailing != nil {
			v.findClickableRecursive(c.Trailing)
		}
	}
}

// findButtonsRecursive recursively searches for buttons and scroll containers
func (v *VIMMode) findButtonsRecursive(obj fyne.CanvasObject, buttons *[]*widget.Button) {
	// Check if this object is a button
	if btn, ok := obj.(*widget.Button); ok {
		*buttons = append(*buttons, btn)
		log.Printf("VIM: Found button: %s", btn.Text)
		// Don't return - button might have children
	}

	// Check all container types
	switch c := obj.(type) {
	case *fyne.Container:
		for _, child := range c.Objects {
			v.findButtonsRecursive(child, buttons)
		}
	case *container.Scroll:
		// Track this scroll container for j/k/h/l scrolling
		v.scrollContainers = append(v.scrollContainers, c)
		if c.Content != nil {
			v.findButtonsRecursive(c.Content, buttons)
		}
	case *container.Split:
		if c.Leading != nil {
			v.findButtonsRecursive(c.Leading, buttons)
		}
		if c.Trailing != nil {
			v.findButtonsRecursive(c.Trailing, buttons)
		}
	}
}

// generateHintKeys assigns letter keys to all clickable targets
func (v *VIMMode) generateHintKeys() {
	letters := "abcdefghijklmnopqrstuvwxyz"
	hintKeys := []string{}

	// Generate single letters
	for _, c := range letters {
		hintKeys = append(hintKeys, string(c))
	}

	// Generate double letters if needed
	numTargets := len(v.hints)
	if len(hintKeys) < numTargets {
		for _, c1 := range letters {
			for _, c2 := range letters {
				hintKeys = append(hintKeys, string(c1)+string(c2))
				if len(hintKeys) >= numTargets {
					break
				}
			}
			if len(hintKeys) >= numTargets {
				break
			}
		}
	}

	// Reassign hints with proper letter keys
	newHints := make(map[string]*clickableTarget)
	idx := 0
	for _, target := range v.hints {
		if idx >= len(hintKeys) {
			break
		}
		newHints[hintKeys[idx]] = target
		log.Printf("VIM: Assigned hint key '%s' to %s", hintKeys[idx], target.label)
		idx++
	}

	v.hints = newHints
}

// showHints displays available hints visually on-screen
func (v *VIMMode) showHints() {
	// Log to file for debugging
	log.Println("===============================")
	log.Println("VIM HINT MODE ACTIVE")
	log.Println("===============================")
	for hint, target := range v.hints {
		log.Printf("  [%s] -> %s", hint, target.label)
	}
	log.Println("Press ESC to exit hint mode")
	log.Println("===============================")

	// Create visual overlay with hints positioned over buttons
	v.createHintOverlay()
}

// createHintOverlay creates hint labels positioned directly over each clickable element
func (v *VIMMode) createHintOverlay() {
	v.hintLabels = []*canvas.Text{}
	overlayObjects := []fyne.CanvasObject{}

	// Get canvas to find absolute positions
	canv := v.state.window.Canvas()

	// Create a hint label for each clickable target
	for hint, target := range v.hints {
		// Get object's absolute position on screen
		objPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(target.obj)
		objSize := target.obj.Size()

		log.Printf("VIM: Target '%s' at position (%v, %v) with size (%v x %v)",
			target.label, objPos.X, objPos.Y, objSize.Width, objSize.Height)

		// Create yellow background box
		bg := canvas.NewRectangle(colorYellow())
		bg.Resize(fyne.NewSize(40, 25))
		bg.Move(objPos)

		// Create hint text (e.g., "A", "B", "C")
		label := canvas.NewText(strings.ToUpper(hint), colorBlack())
		label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		label.TextSize = 16
		label.Alignment = fyne.TextAlignCenter
		label.Resize(fyne.NewSize(40, 25))
		label.Move(objPos)

		v.hintLabels = append(v.hintLabels, label)
		overlayObjects = append(overlayObjects, bg, label)

		log.Printf("VIM: Created hint label '%s' at (%v, %v)", hint, objPos.X, objPos.Y)
	}

	// Create overlay container using NewWithoutLayout for absolute positioning
	v.hintOverlay = container.NewWithoutLayout(overlayObjects...)

	// Add overlay on top of existing content
	currentContent := v.state.window.Content()
	v.state.window.SetContent(container.NewStack(currentContent, v.hintOverlay))
	canv.Refresh(v.state.window.Content())

	log.Printf("VIM: Created %d hint overlays", len(v.hintLabels))
}

// hideHintOverlay removes the visual hint overlay
func (v *VIMMode) hideHintOverlay() {
	if v.hintOverlay == nil {
		return
	}

	// Remove overlay by getting the Stack container and extracting the original content
	currentContent := v.state.window.Content()
	if stack, ok := currentContent.(*fyne.Container); ok {
		if len(stack.Objects) >= 2 {
			// First object is the original content
			v.state.window.SetContent(stack.Objects[0])
			v.state.window.Canvas().Refresh(v.state.window.Content())
		}
	}

	v.hintOverlay = nil
	v.hintLabels = []*canvas.Text{}
	log.Println("VIM: Visual hint overlay hidden")
}

// Helper functions for colors
func colorYellow() color.Color {
	return color.RGBA{R: 255, G: 255, B: 0, A: 230} // Yellow with slight transparency
}

func colorBlack() color.Color {
	return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
}

// activateHint triggers the clickable element associated with a hint
func (v *VIMMode) activateHint(key string) {
	log.Printf("VIM: Looking for hint key: %s", key)
	if target, exists := v.hints[key]; exists {
		log.Printf("VIM: Activating target: %s", target.label)
		// Call the activation function
		target.activate()
		v.exitHintMode()
	} else {
		log.Printf("VIM: No target found for hint: %s", key)
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
	canv.SetOnTypedRune(func(r rune) {
		v.HandleTypedRune(r)
	})
	log.Println("VIM Mode: Handler attached (disabled by default)")
}

// scrollVertical scrolls all scroll containers vertically
func (v *VIMMode) scrollVertical(delta float32) {
	// Find scroll containers if not already found
	if len(v.scrollContainers) == 0 {
		v.findAllButtons() // This also populates scrollContainers
	}

	for _, scroll := range v.scrollContainers {
		currentOffset := scroll.Offset
		newOffset := fyne.NewPos(currentOffset.X, currentOffset.Y+delta)

		// Clamp to valid bounds
		if newOffset.Y < 0 {
			newOffset.Y = 0
		}

		scroll.Offset = newOffset
		scroll.Refresh()
		log.Printf("VIM: Scrolled vertically by %.0f (new Y offset: %.0f)", delta, newOffset.Y)
	}
}

// scrollHorizontal scrolls all scroll containers horizontally
func (v *VIMMode) scrollHorizontal(delta float32) {
	// Find scroll containers if not already found
	if len(v.scrollContainers) == 0 {
		v.findAllButtons() // This also populates scrollContainers
	}

	for _, scroll := range v.scrollContainers {
		currentOffset := scroll.Offset
		newOffset := fyne.NewPos(currentOffset.X+delta, currentOffset.Y)

		// Clamp to valid bounds
		if newOffset.X < 0 {
			newOffset.X = 0
		}

		scroll.Offset = newOffset
		scroll.Refresh()
		log.Printf("VIM: Scrolled horizontally by %.0f (new X offset: %.0f)", delta, newOffset.X)
	}
}

// scrollToTop scrolls all scroll containers to the top
func (v *VIMMode) scrollToTop() {
	// Find scroll containers if not already found
	if len(v.scrollContainers) == 0 {
		v.findAllButtons() // This also populates scrollContainers
	}

	for _, scroll := range v.scrollContainers {
		scroll.Offset = fyne.NewPos(scroll.Offset.X, 0)
		scroll.Refresh()
		log.Println("VIM: Scrolled to top")
	}
}

// scrollToBottom scrolls all scroll containers to the bottom
func (v *VIMMode) scrollToBottom() {
	// Find scroll containers if not already found
	if len(v.scrollContainers) == 0 {
		v.findAllButtons() // This also populates scrollContainers
	}

	for _, scroll := range v.scrollContainers {
		// Get content size to calculate max scroll
		contentSize := scroll.Content.Size()
		scrollSize := scroll.Size()

		maxY := contentSize.Height - scrollSize.Height
		if maxY < 0 {
			maxY = 0
		}

		scroll.Offset = fyne.NewPos(scroll.Offset.X, maxY)
		scroll.Refresh()
		log.Printf("VIM: Scrolled to bottom (Y offset: %.0f)", maxY)
	}
}
