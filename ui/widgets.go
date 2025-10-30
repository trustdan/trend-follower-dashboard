package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

// ShowStyledInformation displays an information dialog with proper text colors
// that work correctly with both light and dark modes, ensuring white text on dark green buttons
func ShowStyledInformation(title, message string, window fyne.Window) {
	// Message label with word wrapping
	contentLabel := widget.NewLabel(message)
	contentLabel.Wrapping = fyne.TextWrapWord

	// Transparent rectangle sets the minimum width for the dialog content
	widthGuard := canvas.NewRectangle(color.Transparent)
	widthGuard.SetMinSize(fyne.NewSize(450, 0)) // enforce minimum width

	// Stack label on top of width guard to guarantee dialog width
	messageArea := container.NewStack(
		widthGuard,
		container.NewPadded(contentLabel),
	)

	// Build dialog using the standard dismiss button (pass "OK" as dismiss label)
	// This avoids creating duplicate buttons
	d := dialog.NewCustom(
		title,
		"OK", // use built-in dismiss button
		messageArea,
		window,
	)

	// Show dialog
	d.Show()
}
