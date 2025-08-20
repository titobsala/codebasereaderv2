package core

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// formatFileSize formats file size in human readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// formatNumber formats numbers with thousand separators
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}

// getLangIcon returns an icon for the programming language
func getLangIcon(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		return "⚡" // Lightning for Go (fast)
	case "python":
		return "🐍" // Snake for Python
	case "javascript":
		return "JS"
	case "typescript":
		return "TS"
	case "java":
		return "☕"
	case "c":
		return "C"
	case "c++", "cpp":
		return "C++"
	case "rust":
		return "🦀"
	case "php":
		return "PHP"
	case "ruby":
		return "💎"
	default:
		return "📄"
	}
}

// getGradeStyle returns the appropriate pre-cached style for a quality grade
func getGradeStyle(grade string) lipgloss.Style {
	switch grade {
	case "A":
		return components.GradeStyleA
	case "B":
		return components.GradeStyleB
	case "C":
		return components.GradeStyleC
	case "D":
		return components.GradeStyleD
	case "F":
		return components.GradeStyleF
	default:
		return components.GradeStyleDefault
	}
}
