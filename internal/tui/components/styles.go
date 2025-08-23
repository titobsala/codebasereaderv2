package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/tui"
)

// Color palette
var (
	// Primary colors
	PrimaryPurple = lipgloss.Color("#7D56F4")
	PrimaryBlue   = lipgloss.Color("#4FC3F7")
	PrimaryCyan   = lipgloss.Color("#26C6DA")
	PrimaryGreen  = lipgloss.Color("#66BB6A")
	PrimaryYellow = lipgloss.Color("#FFEB3B")
	PrimaryOrange = lipgloss.Color("#FF9800")
	PrimaryRed    = lipgloss.Color("#EF5350")

	// Accent colors
	AccentPink   = lipgloss.Color("#E91E63")
	AccentTeal   = lipgloss.Color("#009688")
	AccentIndigo = lipgloss.Color("#3F51B5")

	// Neutral colors
	NeutralWhite  = lipgloss.Color("#FAFAFA")
	NeutralLight  = lipgloss.Color("#E0E0E0")
	NeutralMedium = lipgloss.Color("#9E9E9E")
	NeutralDark   = lipgloss.Color("#424242")
	NeutralBlack  = lipgloss.Color("#212121")

	// Semantic colors
	SuccessGreen  = lipgloss.Color("#4CAF50")
	WarningOrange = lipgloss.Color("#FF9800")
	ErrorRed      = lipgloss.Color("#F44336")
	InfoBlue      = lipgloss.Color("#2196F3")
)

// Enhanced styles for the TUI
var (
	// Main title with gradient effect
	TitleStyle = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(PrimaryPurple).
			Padding(0, 2).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryBlue)

	// Enhanced status bar with gradient-like effect
	StatusStyle = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(NeutralDark).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.Border{
			Top:    "▔",
			Bottom: "▁",
			Left:   "▏",
			Right:  "▕",
		}).
		BorderForeground(PrimaryCyan)

	// Subtle help text with improved readability
	HelpStyle = lipgloss.NewStyle().
			Foreground(NeutralMedium).
			Italic(true).
			MarginLeft(1)

	// Enhanced error styling with background
	ErrorStyle = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(ErrorRed).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.ThickBorder()).
			BorderForeground(ErrorRed)

	// Selected items with enhanced visual feedback
	SelectedStyle = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(PrimaryPurple).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(PrimaryBlue)

	// Main headers with sophisticated styling
	HeaderStyle = lipgloss.NewStyle().
			Foreground(PrimaryPurple).
			Bold(true).
			Underline(true).
			MarginBottom(1).
			Padding(0, 1).
			Border(lipgloss.Border{
			Bottom: "=",
		}).
		BorderForeground(PrimaryBlue)

	// Section headers with accent styling
	SectionStyle = lipgloss.NewStyle().
			Foreground(PrimaryCyan).
			Bold(true).
			MarginTop(1).
			MarginBottom(1).
			Padding(0, 1).
			Border(lipgloss.Border{
			Left: "|",
		}).
		BorderForeground(PrimaryBlue)

	// Summary style with enhanced visual appeal
	SummaryStyle = lipgloss.NewStyle().
			Foreground(PrimaryBlue).
			Bold(true).
			Italic(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryCyan).
			Background(lipgloss.Color("#1A1A2E"))

	// Enhanced grade styles with backgrounds
	GradeStyleA = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(SuccessGreen).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessGreen)

	GradeStyleB = lipgloss.NewStyle().
			Foreground(NeutralBlack).
			Background(PrimaryGreen).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryGreen)

	GradeStyleC = lipgloss.NewStyle().
			Foreground(NeutralBlack).
			Background(PrimaryYellow).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryYellow)

	GradeStyleD = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(WarningOrange).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(WarningOrange)

	GradeStyleF = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(ErrorRed).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ErrorRed)

	GradeStyleDefault = lipgloss.NewStyle().
				Foreground(NeutralMedium).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(NeutralMedium)

	// Improved scroll info with better visibility
	ScrollInfoStyle = lipgloss.NewStyle().
			Foreground(NeutralMedium).
			Italic(true).
			Align(lipgloss.Center).
			MarginTop(1).
			Padding(0, 1).
			Border(lipgloss.Border{
			Top: "-",
		}).
		BorderForeground(NeutralDark)

	// New advanced styles for enhanced visual appeal

	// Card-like containers for content sections
	CardStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryBlue).
			Background(lipgloss.Color("#1E1E2E"))

	// Highlight boxes for important information
	highlightStyle = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(AccentTeal).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(PrimaryCyan)

	// Metric value styling
	MetricValueStyle = lipgloss.NewStyle().
				Foreground(PrimaryGreen).
				Bold(true)

	// Warning text styling
	warningStyle = lipgloss.NewStyle().
			Foreground(WarningOrange).
			Bold(true).
			Background(lipgloss.Color("#2D1B00")).
			Padding(0, 1).
			Border(lipgloss.Border{
			Left: "|",
		}).
		BorderForeground(WarningOrange)

	// Info box styling
	infoStyle = lipgloss.NewStyle().
			Foreground(InfoBlue).
			Background(lipgloss.Color("#0A1929")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(InfoBlue)

	// Badge styling for counts and numbers
	badgeStyle = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(AccentPink).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentPink)
)

// Helper functions for advanced styling

// FormatNumberStyled formats numbers with styling applied
func FormatNumberStyled(n int) string {
	formatted := tui.FormatNumber(n) // Use existing function from utils.go
	return MetricValueStyle.Render(formatted)
}

// CreateProgressBar creates an enhanced progress bar with gradient colors
func CreateProgressBar(value float64, width int, withPercentage bool) string {
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
		fillColor = SuccessGreen
		fillChar = "█"
	} else if value >= 60 {
		fillColor = PrimaryGreen
		fillChar = "█"
	} else if value >= 40 {
		fillColor = PrimaryYellow
		fillChar = "▓"
	} else if value >= 20 {
		fillColor = WarningOrange
		fillChar = "▒"
	} else {
		fillColor = ErrorRed
		fillChar = "░"
	}

	emptyChar = "░"

	filledPart := lipgloss.NewStyle().
		Foreground(fillColor).
		Render(strings.Repeat(fillChar, filled))

	emptyPart := lipgloss.NewStyle().
		Foreground(NeutralDark).
		Render(strings.Repeat(emptyChar, width-filled))

	bar := lipgloss.JoinHorizontal(lipgloss.Left, filledPart, emptyPart)

	if withPercentage {
		percentage := lipgloss.NewStyle().
			Foreground(NeutralMedium).
			Render(fmt.Sprintf(" %.1f%%", value))
		return lipgloss.JoinHorizontal(lipgloss.Left, "[", bar, "]", percentage)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, "[", bar, "]")
}

// CreateCard wraps content in a styled card container
func CreateCard(title, content string) string {
	header := TitleStyle.
		Width(50).
		Align(lipgloss.Center).
		Render(title)

	body := CardStyle.
		Width(50).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

// CreateBadge creates a styled badge for counts or labels
func CreateBadge(text string, badgeType string) string {
	var style lipgloss.Style

	switch badgeType {
	case "success":
		style = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(SuccessGreen).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessGreen)
	case "warning":
		style = lipgloss.NewStyle().
			Foreground(NeutralBlack).
			Background(WarningOrange).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(WarningOrange)
	case "error":
		style = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(ErrorRed).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ErrorRed)
	case "info":
		style = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(InfoBlue).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(InfoBlue)
	default:
		style = badgeStyle
	}

	return style.Render(text)
}

// CreateHighlight creates a highlighted text box
func CreateHighlight(text string, highlightType string) string {
	var style lipgloss.Style

	switch highlightType {
	case "success":
		style = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(SuccessGreen).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.Border{Left: "▌"}).
			BorderForeground(SuccessGreen)
	case "warning":
		style = warningStyle
	case "error":
		style = lipgloss.NewStyle().
			Foreground(NeutralWhite).
			Background(ErrorRed).
			Padding(0, 1).
			Bold(true).
			Border(lipgloss.Border{Left: "▌"}).
			BorderForeground(ErrorRed)
	case "info":
		style = infoStyle
	default:
		style = highlightStyle
	}

	return style.Render(text)
}
