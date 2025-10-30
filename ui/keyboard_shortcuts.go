package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// setupKeyboardShortcuts registers global keyboard shortcuts for the application
func setupKeyboardShortcuts(state *AppState, window fyne.Window, onNewTrade, onResume, onHistory func()) {
	// Get canvas
	canvas := window.Canvas()

	// Ctrl+N: New Trade
	canvas.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		log.Println("Keyboard shortcut: Ctrl+N (New Trade)")
		if onNewTrade != nil {
			onNewTrade()
		}
	})

	// Ctrl+R: Resume Session
	canvas.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyR,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		log.Println("Keyboard shortcut: Ctrl+R (Resume Session)")
		if onResume != nil {
			onResume()
		}
	})

	// Ctrl+H: Session History
	canvas.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyH,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		log.Println("Keyboard shortcut: Ctrl+H (Session History)")
		if onHistory != nil {
			onHistory()
		}
	})

	log.Println("Keyboard shortcuts registered: Ctrl+N (New), Ctrl+R (Resume), Ctrl+H (History)")
}
