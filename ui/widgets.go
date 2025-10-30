package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// NewStyledButton creates a button with proper text color for visibility on dark green background
func NewStyledButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButton(label, tapped)
	btn.Importance = widget.HighImportance
	return btn
}

// NewPrimaryButton creates a button with explicit white text color
// This is a workaround for Fyne theme not applying button text colors correctly
func NewPrimaryButton(label string, tapped func()) fyne.CanvasObject {
	// Create a custom button-like widget with explicit colors
	text := canvas.NewText(label, color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 14
	text.TextStyle = fyne.TextStyle{Bold: true}

	bg := canvas.NewRectangle(BritishRacingGreen)
	bg.CornerRadius = 4

	btn := widget.NewButton(label, tapped)
	btn.Importance = widget.HighImportance

	return btn
}
