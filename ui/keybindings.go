package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// VIMKeyHandler manages VIM-style keyboard navigation
type VIMKeyHandler struct {
	state          *AppState
	hintMode       bool
	hints          map[string]*widget.Button
	hintOverlay    *fyne.Container
	currentButtons []*widget.Button
}

// NewVIMKeyHandler creates a new VIM keyboard handler
func NewVIMKeyHandler(state *AppState) *VIMKeyHandler {
	return &VIMKeyHandler{
		state:       state,
		hintMode:    false,
		hints:       make(map[string]*widget.Button),
		hintOverlay: container.NewWithoutLayout(),
	}
}

// HandleKeyDown processes keyboard events for VIM navigation
func (v *VIMKeyHandler) HandleKeyDown(ev *fyne.KeyEvent) {
	log.Printf("VIM: Key pressed: %v", ev.Name)

	// F key - Toggle hint mode
	if ev.Name == fyne.KeyF {
		log.Println("VIM: F key detected - toggling hint mode")
		v.toggleHintMode()
		return
	}

	// ESC key - Exit hint mode
	if ev.Name == fyne.KeyEscape {
		if v.hintMode {
			log.Println("VIM: ESC key detected - exiting hint mode")
			v.exitHintMode()
		}
		return
	}

	// In hint mode, handle letter keys
	if v.hintMode {
		if len(ev.Name) == 1 {
			log.Printf("VIM: Letter key in hint mode: %s", ev.Name)
			v.activateHint(string(ev.Name))
		}
		return
	}

	// VIM navigation keys (when not in hint mode)
	switch ev.Name {
	case fyne.KeyJ:
		log.Println("VIM: J key - navigate down")
		// TODO: Implement scroll down or focus next item
	case fyne.KeyK:
		log.Println("VIM: K key - navigate up")
		// TODO: Implement scroll up or focus previous item
	case fyne.KeyD:
		log.Println("VIM: D key - page down")
		// TODO: Implement page down
	case fyne.KeyU:
		log.Println("VIM: U key - page up")
		// TODO: Implement page up
	}
}

// toggleHintMode enables/disables hint mode
func (v *VIMKeyHandler) toggleHintMode() {
	if v.hintMode {
		v.exitHintMode()
	} else {
		v.enterHintMode()
	}
}

// enterHintMode activates link hints
func (v *VIMKeyHandler) enterHintMode() {
	log.Println("VIM: Entering hint mode...")
	v.hintMode = true
	v.currentButtons = v.findAllButtons()
	log.Printf("VIM: Found %d buttons", len(v.currentButtons))

	// Generate hints
	v.generateHints()
	log.Println("VIM: Hint mode activated")
}

// exitHintMode deactivates link hints
func (v *VIMKeyHandler) exitHintMode() {
	log.Println("VIM: Exiting hint mode...")
	v.hintMode = false
	v.hints = make(map[string]*widget.Button)
	v.hintOverlay.Objects = []fyne.CanvasObject{}
	v.hintOverlay.Refresh()
	v.currentButtons = []*widget.Button{}
	log.Println("VIM: Hint mode exited")
}

// findAllButtons recursively finds all buttons in the UI
func (v *VIMKeyHandler) findAllButtons() []*widget.Button {
	buttons := []*widget.Button{}

	// Get the main window content
	content := v.state.window.Content()

	// Recursively search for buttons
	v.findButtonsRecursive(content, &buttons)

	return buttons
}

// findButtonsRecursive recursively searches for buttons
func (v *VIMKeyHandler) findButtonsRecursive(obj fyne.CanvasObject, buttons *[]*widget.Button) {
	if btn, ok := obj.(*widget.Button); ok {
		*buttons = append(*buttons, btn)
		return
	}

	// Check if it's a container with multiple objects
	if c, ok := obj.(*fyne.Container); ok {
		for _, child := range c.Objects {
			v.findButtonsRecursive(child, buttons)
		}
	}

	// Check container types
	if c, ok := obj.(interface{ Objects() []fyne.CanvasObject }); ok {
		for _, child := range c.Objects() {
			v.findButtonsRecursive(child, buttons)
		}
	}
}

// generateHints creates keyboard hints for all buttons
func (v *VIMKeyHandler) generateHints() {
	// Generate letter combinations (a-z, then aa-zz, etc.)
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
	for i, btn := range v.currentButtons {
		if i >= len(hints) {
			break
		}

		hint := hints[i]
		v.hints[hint] = btn

		// Create visual hint label
		// Note: In a real implementation, we'd need to position these
		// labels over the buttons. This is challenging in Fyne as we
		// don't have direct control over absolute positioning.
		// For now, we'll just show a notification about available hints.
	}

	// Show hint notification
	v.showHintNotification()
}

// showHintNotification displays available hints
func (v *VIMKeyHandler) showHintNotification() {
	// Log all available hints
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
func (v *VIMKeyHandler) activateHint(key string) {
	log.Printf("VIM: Looking for hint key: %s", key)
	if btn, exists := v.hints[key]; exists {
		log.Printf("VIM: Activating button for hint: %s", key)
		btn.OnTapped()
		v.exitHintMode()
	} else {
		log.Printf("VIM: No button found for hint: %s", key)
	}
}

// AttachToWindow attaches the VIM key handler to the window
func (v *VIMKeyHandler) AttachToWindow() {
	// Add keyboard event handler
	canv := v.state.window.Canvas()

	// Set typed key handler
	canv.SetOnTypedKey(func(ev *fyne.KeyEvent) {
		v.HandleKeyDown(ev)
	})
}
