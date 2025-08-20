package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	primaryPurple = lipgloss.Color("#7D56F4")
	primaryBlue   = lipgloss.Color("#4FC3F7")
	primaryCyan   = lipgloss.Color("#26C6DA")
	primaryGreen  = lipgloss.Color("#66BB6A")
	primaryYellow = lipgloss.Color("#FFEB3B")
	primaryOrange = lipgloss.Color("#FF9800")
	primaryRed    = lipgloss.Color("#EF5350")

	// Accent colors
	accentPink   = lipgloss.Color("#E91E63")
	accentTeal   = lipgloss.Color("#009688")
	accentIndigo = lipgloss.Color("#3F51B5")

	// Neutral colors
	neutralWhite  = lipgloss.Color("#FAFAFA")
	neutralLight  = lipgloss.Color("#E0E0E0")
	neutralMedium = lipgloss.Color("#9E9E9E")
	neutralDark   = lipgloss.Color("#424242")
	neutralBlack  = lipgloss.Color("#212121")

	// Semantic colors
	successGreen  = lipgloss.Color("#4CAF50")
	warningOrange = lipgloss.Color("#FF9800")
	errorRed      = lipgloss.Color("#F44336")
	infoBlue      = lipgloss.Color("#2196F3")
)

// Enhanced styles for the TUI
var (
	// Main title with gradient effect
	titleStyle = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(primaryPurple).
			Padding(0, 2).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryBlue)

	// Enhanced status bar with gradient-like effect
	statusStyle = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(neutralDark).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.Border{
			Top:    "▔",
			Bottom: "▁",
			Left:   "▏",
			Right:  "▕",
		}).
		BorderForeground(primaryCyan)

	// Subtle help text with improved readability
	helpStyle = lipgloss.NewStyle().
			Foreground(neutralMedium).
			Italic(true).
			MarginLeft(1)

	// Enhanced error styling with background
	errorStyle = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(errorRed).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.ThickBorder()).
			BorderForeground(errorRed)

	// Selected items with enhanced visual feedback
	selectedStyle = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(primaryPurple).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryBlue)

	// Normal text with better contrast
	normalStyle = lipgloss.NewStyle().
			Foreground(neutralLight)

	// Main headers with sophisticated styling
	headerStyle = lipgloss.NewStyle().
			Foreground(primaryPurple).
			Bold(true).
			Underline(true).
			MarginBottom(1).
			Padding(0, 1).
			Border(lipgloss.Border{
			Bottom: "═",
		}).
		BorderForeground(primaryBlue)

	// Section headers with accent styling
	sectionStyle = lipgloss.NewStyle().
			Foreground(primaryCyan).
			Bold(true).
			MarginTop(1).
			MarginBottom(1).
			Padding(0, 1).
			Border(lipgloss.Border{
			Left: "▌",
		}).
		BorderForeground(primaryBlue)

	// Summary style with enhanced visual appeal
	summaryStyle = lipgloss.NewStyle().
			Foreground(primaryBlue).
			Bold(true).
			Italic(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryCyan).
			Background(lipgloss.Color("#1A1A2E"))

	// Enhanced grade styles with backgrounds
	gradeStyleA = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(successGreen).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successGreen)

	gradeStyleB = lipgloss.NewStyle().
			Foreground(neutralBlack).
			Background(primaryGreen).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryGreen)

	gradeStyleC = lipgloss.NewStyle().
			Foreground(neutralBlack).
			Background(primaryYellow).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryYellow)

	gradeStyleD = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(warningOrange).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningOrange)

	gradeStyleF = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(errorRed).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorRed)

	gradeStyleDefault = lipgloss.NewStyle().
				Foreground(neutralMedium).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(neutralMedium)

	// Enhanced separator with visual flair
	separatorStyle = lipgloss.NewStyle().
			Foreground(primaryBlue).
			Bold(true)

	// Improved scroll info with better visibility
	scrollInfoStyle = lipgloss.NewStyle().
			Foreground(neutralMedium).
			Italic(true).
			Align(lipgloss.Center).
			MarginTop(1).
			Padding(0, 1).
			Border(lipgloss.Border{
			Top: "─",
		}).
		BorderForeground(neutralDark)

	// New advanced styles for enhanced visual appeal

	// Card-like containers for content sections
	cardStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryBlue).
			Background(lipgloss.Color("#1E1E2E"))

	// Highlight boxes for important information
	highlightStyle = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(accentTeal).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryCyan)

	// Progress bar container
	progressContainerStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryBlue).
				Background(neutralDark)

	// Metric value styling
	metricValueStyle = lipgloss.NewStyle().
				Foreground(primaryGreen).
				Bold(true)

	// Warning text styling
	warningStyle = lipgloss.NewStyle().
			Foreground(warningOrange).
			Bold(true).
			Background(lipgloss.Color("#2D1B00")).
			Padding(0, 1).
			Border(lipgloss.Border{
			Left: "▌",
		}).
		BorderForeground(warningOrange)

	// Info box styling
	infoStyle = lipgloss.NewStyle().
			Foreground(infoBlue).
			Background(lipgloss.Color("#0A1929")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(infoBlue)

	// Language tag styling
	languageTagStyle = lipgloss.NewStyle().
				Foreground(neutralWhite).
				Background(accentIndigo).
				Padding(0, 1).
				Bold(true).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryPurple)

	// Table header styling
	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(neutralWhite).
				Background(primaryPurple).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.Border{
			Bottom: "═",
		}).
		BorderForeground(primaryBlue)

	// Table row styling
	tableRowStyle = lipgloss.NewStyle().
			Foreground(neutralLight).
			Padding(0, 1).
			Border(lipgloss.Border{
			Bottom: "─",
		}).
		BorderForeground(neutralDark)

	// Badge styling for counts and numbers
	badgeStyle = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(accentPink).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentPink)
)

// Helper functions for advanced styling

// FormatNumberStyled formats numbers with styling applied
func FormatNumberStyled(n int) string {
	formatted := FormatNumber(n) // Use existing function from utils.go
	return metricValueStyle.Render(formatted)
}

// createProgressBar creates an enhanced progress bar with gradient colors
func createProgressBar(value float64, width int, withPercentage bool) string {
	if width <= 0 {
		width = 20
	}

	filled := int(float64(width) * value / 100.0)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	// Create gradient effect based on progress
	var fillChar, emptyChar string
	var fillColor lipgloss.Color

	if value >= 80 {
		fillColor = successGreen
		fillChar = "█"
	} else if value >= 60 {
		fillColor = primaryGreen
		fillChar = "█"
	} else if value >= 40 {
		fillColor = primaryYellow
		fillChar = "▓"
	} else if value >= 20 {
		fillColor = warningOrange
		fillChar = "▒"
	} else {
		fillColor = errorRed
		fillChar = "░"
	}

	emptyChar = "░"

	filledPart := lipgloss.NewStyle().
		Foreground(fillColor).
		Render(strings.Repeat(fillChar, filled))

	emptyPart := lipgloss.NewStyle().
		Foreground(neutralDark).
		Render(strings.Repeat(emptyChar, width-filled))

	bar := lipgloss.JoinHorizontal(lipgloss.Left, filledPart, emptyPart)

	if withPercentage {
		percentage := lipgloss.NewStyle().
			Foreground(neutralMedium).
			Render(fmt.Sprintf(" %.1f%%", value))
		return lipgloss.JoinHorizontal(lipgloss.Left, "[", bar, "]", percentage)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, "[", bar, "]")
}

// createCard wraps content in a styled card container
func createCard(title, content string) string {
	header := titleStyle.
		Width(50).
		Align(lipgloss.Center).
		Render(title)

	body := cardStyle.
		Width(50).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

// createBadge creates a styled badge for counts or labels
func createBadge(text string, badgeType string) string {
	var style lipgloss.Style

	switch badgeType {
	case "success":
		style = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(successGreen).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successGreen)
	case "warning":
		style = lipgloss.NewStyle().
			Foreground(neutralBlack).
			Background(warningOrange).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningOrange)
	case "error":
		style = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(errorRed).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorRed)
	case "info":
		style = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(infoBlue).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(infoBlue)
	default:
		style = badgeStyle
	}

	return style.Render(text)
}

// createSeparator creates a styled separator line
func createSeparator(width int, char string) string {
	if char == "" {
		char = "─"
	}

	return separatorStyle.
		Width(width).
		Align(lipgloss.Center).
		Render(strings.Repeat(char, width))
}

// createHighlight creates a highlighted text box
func createHighlight(text string, highlightType string) string {
	var style lipgloss.Style

	switch highlightType {
	case "success":
		style = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(successGreen).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.Border{Left: "▌"}).
			BorderForeground(successGreen)
	case "warning":
		style = warningStyle
	case "error":
		style = lipgloss.NewStyle().
			Foreground(neutralWhite).
			Background(errorRed).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.Border{Left: "▌"}).
			BorderForeground(errorRed)
	case "info":
		style = infoStyle
	default:
		style = highlightStyle
	}

	return style.Render(text)
}
