package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createButtonWithWhiteText creates a button with guaranteed white text on dark green background
// This fixes contrast issues where Fyne's theme doesn't always apply white text correctly
func createButtonWithWhiteText(label string, importance widget.ButtonImportance, tapped func()) fyne.CanvasObject {
	// Create a container that overlays white text on top of a standard button
	btn := widget.NewButton(label, tapped)
	btn.Importance = importance

	// Create white text overlay
	text := canvas.NewText(label, color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 14
	text.TextStyle = fyne.TextStyle{Bold: true}

	// Stack the white text on top of the button
	// This ensures white text even if theme doesn't apply correctly
	return container.NewStack(btn, container.NewCenter(text))
}

// NewPrimaryButtonWithWhiteText creates a high-importance button with white text
func NewPrimaryButtonWithWhiteText(label string, tapped func()) fyne.CanvasObject {
	return createButtonWithWhiteText(label, widget.HighImportance, tapped)
}

// NewDangerButtonWithWhiteText creates a danger button with white text
func NewDangerButtonWithWhiteText(label string, tapped func()) fyne.CanvasObject {
	return createButtonWithWhiteText(label, widget.DangerImportance, tapped)
}

// NewSuccessButtonWithWhiteText creates a success button with white text
func NewSuccessButtonWithWhiteText(label string, tapped func()) fyne.CanvasObject {
	return createButtonWithWhiteText(label, widget.SuccessImportance, tapped)
}

// NewMediumButtonWithWhiteText creates a medium-importance button with white text
func NewMediumButtonWithWhiteText(label string, tapped func()) fyne.CanvasObject {
	return createButtonWithWhiteText(label, widget.MediumImportance, tapped)
}

// FixButtonTextColor wraps an existing button to ensure white text
// Use this for buttons you can't easily replace
func FixButtonTextColor(btn *widget.Button, label string) fyne.CanvasObject {
	text := canvas.NewText(label, color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 14
	text.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewStack(btn, container.NewCenter(text))
}
