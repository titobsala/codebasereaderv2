package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Min returns the minimum of two integers
func Min(a, b int) int {
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

// FormatFileSize formats file size in human readable format
func FormatFileSize(size int64) string {
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

// FormatNumber formats numbers with thousand separators
func FormatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}

// GetLangIcon returns an icon for the programming language
func GetLangIcon(lang string) string {
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

// GetGradeStyle returns the appropriate pre-cached style for a quality grade
func GetGradeStyle(grade string) lipgloss.Style {
	switch grade {
	case "A":
		return gradeStyleA
	case "B":
		return gradeStyleB
	case "C":
		return gradeStyleC
	case "D":
		return gradeStyleD
	case "F":
		return gradeStyleF
	default:
		return gradeStyleDefault
	}
}
