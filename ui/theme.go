package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ThemeVariantType represents the custom theme variant
type ThemeVariantType int

const (
	ThemeLight ThemeVariantType = iota
	ThemeDark
)

// tfTheme is the custom theme for TF-Engine with British Racing Green
type tfTheme struct {
	variant ThemeVariantType
}

var _ fyne.Theme = (*tfTheme)(nil)

// Color scheme constants
var (
	// British Racing Green palette
	BritishRacingGreen = color.RGBA{R: 0x00, G: 0x42, B: 0x25, A: 0xFF} // #004225
	ForestGreen        = color.RGBA{R: 0x22, G: 0x8B, B: 0x22, A: 0xFF} // #228B22
	LightGreen         = color.RGBA{R: 0x90, G: 0xEE, B: 0x90, A: 0xFF} // #90EE90

	// Complementary pink/rose palette
	DustyRose = color.RGBA{R: 0xD8, G: 0xA7, B: 0xB1, A: 0xFF} // #D8A7B1
	SoftPink  = color.RGBA{R: 0xFF, G: 0xB3, B: 0xC1, A: 0xFF} // #FFB3C1

	// Banner colors (discipline enforcement)
	BannerGreen  = color.RGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0xFF} // #2E7D32 - GO
	BannerYellow = color.RGBA{R: 0xFF, G: 0xC1, B: 0x07, A: 0xFF} // #FFC107 - CAUTION
	BannerRed    = color.RGBA{R: 0xD3, G: 0x2F, B: 0x2F, A: 0xFF} // #D32F2F - NO-GO

	// Light mode backgrounds
	LightBackground = color.RGBA{R: 0xF5, G: 0xF5, B: 0xF0, A: 0xFF} // #F5F5F0 - Warm off-white
	LightCard       = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // #FFFFFF - White

	// Dark mode backgrounds
	DarkBackground = color.RGBA{R: 0x1E, G: 0x1E, B: 0x1E, A: 0xFF} // #1E1E1E - Dark charcoal
	DarkCard       = color.RGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xFF} // #282828 - Dark card

	// Text colors
	LightText = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Black
	DarkText  = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // White
)

func (t *tfTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	isLight := t.variant == ThemeLight

	switch name {
	case theme.ColorNamePrimary:
		return BritishRacingGreen

	case theme.ColorNameButton:
		// For light mode, we want buttons to still be British Racing Green
		// The button text color will be handled by ForegroundOnPrimary
		return BritishRacingGreen

	// Force button text to be white for all button types
	case theme.ColorNameDisabled:
		return color.RGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}

	case theme.ColorNameHover:
		return ForestGreen

	case theme.ColorNameFocus:
		return ForestGreen

	case theme.ColorNameSuccess:
		return BannerGreen

	case theme.ColorNameWarning:
		if isLight {
			return DustyRose
		}
		return SoftPink

	case theme.ColorNameError:
		return BannerRed

	case theme.ColorNameBackground:
		if isLight {
			return LightBackground
		}
		return DarkBackground

	// Foreground color for general text (labels, etc.)
	case theme.ColorNameForeground:
		if isLight {
			return LightText  // Black text for light backgrounds
		}
		return DarkText  // White text for dark backgrounds

	// CRITICAL: Ensure ALL button text is white for good contrast on green backgrounds
	// These apply to text on colored buttons
	case theme.ColorNameForegroundOnPrimary:
		// Always white for good contrast on British Racing Green
		return color.White

	case theme.ColorNameForegroundOnError:
		// White text on red danger buttons
		return color.White

	case theme.ColorNameForegroundOnSuccess:
		// White text on green success buttons
		return color.White

	case theme.ColorNameForegroundOnWarning:
		// White text on yellow/orange warning buttons (changed from black)
		return color.White

	case theme.ColorNameInputBackground:
		if isLight {
			return LightCard
		}
		return DarkCard

	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}

	case theme.ColorNamePlaceHolder:
		if isLight {
			return color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
		}
		return color.RGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF}

	case theme.ColorNameScrollBar:
		if isLight {
			return ForestGreen
		}
		return LightGreen

	case theme.ColorNameShadow:
		if isLight {
			return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33}
		}
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x66}
	}

	// Fall back to Fyne's default theme with appropriate variant
	fyneVariant := theme.VariantLight
	if !isLight {
		fyneVariant = theme.VariantDark
	}
	return theme.DefaultTheme().Color(name, fyneVariant)
}

func (t *tfTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *tfTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *tfTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 20
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	}
	return theme.DefaultTheme().Size(name)
}

// Helper functions for banner colors (discipline enforcement)
func ColorGreen() color.Color {
	return BannerGreen
}

func ColorYellow() color.Color {
	return BannerYellow
}

func ColorRed() color.Color {
	return BannerRed
}
