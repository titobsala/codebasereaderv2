package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Common styling helper functions following Go's "a little copying is better than a little dependency" principle

// NewHeaderStyle creates a consistent header style
func NewHeaderStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true)
}

// NewSectionStyle creates a consistent section style
func NewSectionStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true)
}

// NewHighlightStyle creates a highlight style with background
func NewHighlightStyle(fg, bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Padding(0, 1)
}

// NewBorderedStyle creates a style with border
func NewBorderedStyle(borderColor lipgloss.Color, padding ...int) lipgloss.Style {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	if len(padding) >= 2 {
		style = style.Padding(padding[0], padding[1])
	} else if len(padding) == 1 {
		style = style.Padding(padding[0])
	}

	return style
}

// NewCenteredStyle creates a centered style
func NewCenteredStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(color).
		Align(lipgloss.Center)
}

// NewDialogStyle creates a dialog box style
func NewDialogStyle(borderColor lipgloss.Color, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(width - 4).
		Align(lipgloss.Center)
}
