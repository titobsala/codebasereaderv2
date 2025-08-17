package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
)

// ContentViewModel handles the content display area
type ContentViewModel struct {
	content        string
	filePath       string
	scrollY        int
	maxScroll      int
	showMetrics    bool
	showSummary    bool
	width          int
	height         int
	analysisData   *AnalysisData
	metricsDisplay *MetricsDisplay
	currentMode    MetricsDisplayMode
}

// NewContentViewModel creates a new content view model
func NewContentViewModel() *ContentViewModel {
	return &ContentViewModel{
		content:        "Select a file to view its content",
		scrollY:        0,
		maxScroll:      0,
		showMetrics:    false,
		showSummary:    false,
		metricsDisplay: NewMetricsDisplay(),
		currentMode:    OverviewMode,
	}
}

// Init initializes the content view
func (m ContentViewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the content view
func (m *ContentViewModel) Update(msg tea.Msg) (*ContentViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(-1)
			} else if m.scrollY > 0 {
				m.scrollY--
			}
		case "down", "j":
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(1)
			} else if m.scrollY < m.maxScroll {
				m.scrollY++
			}
		case "pgup":
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(-10)
			} else {
				m.scrollY = max(0, m.scrollY-10)
			}
		case "pgdown":
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(10)
			} else {
				m.scrollY = min(m.maxScroll, m.scrollY+10)
			}
		case "home", "g":
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(-1000) // Reset to top
			} else {
				m.scrollY = 0
			}
		case "end", "G":
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(1000) // Go to bottom
			} else {
				m.scrollY = m.maxScroll
			}
		case "ctrl+u":
			// Scroll up half page
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(-m.height/2)
			} else {
				m.scrollY = max(0, m.scrollY-m.height/2)
			}
		case "ctrl+d":
			// Scroll down half page
			if m.showMetrics && m.metricsDisplay != nil {
				m.metricsDisplay.Scroll(m.height/2)
			} else {
				m.scrollY = min(m.maxScroll, m.scrollY+m.height/2)
			}
		case "1":
			// Switch to overview mode
			if m.showMetrics {
				m.currentMode = OverviewMode
				m.metricsDisplay.SetMode(OverviewMode)
			}
		case "2":
			// Switch to detailed mode
			if m.showMetrics {
				m.currentMode = DetailedMode
				m.metricsDisplay.SetMode(DetailedMode)
			}
		case "3":
			// Switch to quality mode
			if m.showMetrics {
				m.currentMode = QualityMode
				m.metricsDisplay.SetMode(QualityMode)
			}
		case "4":
			// Switch to dependency mode
			if m.showMetrics {
				m.currentMode = DependencyMode
				m.metricsDisplay.SetMode(DependencyMode)
			}
		}

	case ToggleMetricsMsg:
		m.showMetrics = !m.showMetrics
		m.showSummary = false
		m.updateContentFromAnalysis()

	case ToggleSummaryMsg:
		m.showSummary = !m.showSummary
		m.showMetrics = false
		m.updateContentFromAnalysis()

	case FileContentLoadedMsg:
		m.SetContent(msg.FilePath, msg.Content)

	case ClearAnalysisMsg:
		m.analysisData = nil
		m.content = "Select a file to view its content"
		m.filePath = ""
		m.scrollY = 0
		m.showMetrics = false
		m.showSummary = false
	}
	return m, nil
}

// View renders the content view
func (m *ContentViewModel) View(width, height int) string {
	m.width = width
	m.height = height

	var b strings.Builder

	// Header
	if m.filePath != "" {
		header := fmt.Sprintf("File: %s", m.filePath)
		b.WriteString(selectedStyle.Render(header) + "\n\n")
	}

	// Content area
	contentHeight := height - 4 // Reserve space for header and controls
	
	if m.content == "" {
		b.WriteString(helpStyle.Render("No content to display"))
		return b.String()
	}

	// Split content into lines
	lines := strings.Split(m.content, "\n")
	newMaxScroll := max(0, len(lines)-contentHeight)
	
	// Only update maxScroll if content changed significantly
	if m.maxScroll != newMaxScroll {
		m.maxScroll = newMaxScroll
		// Ensure scroll position is still valid
		if m.scrollY > m.maxScroll {
			m.scrollY = m.maxScroll
		}
	}

	// Display visible lines
	startLine := m.scrollY
	if startLine < 0 {
		startLine = 0
		m.scrollY = 0
	}
	if startLine >= len(lines) && len(lines) > 0 {
		startLine = max(0, len(lines)-contentHeight)
		m.scrollY = startLine
	}
	
	endLine := min(len(lines), startLine+contentHeight)

	for i := startLine; i < endLine; i++ {
		if i < len(lines) {
			line := lines[i]
			if len(line) > width-2 {
				line = line[:width-5] + "..."
			}
			b.WriteString(line + "\n")
		}
	}

	// Scroll indicator
	if m.maxScroll > 0 {
		scrollInfo := fmt.Sprintf("Line %d-%d of %d", startLine+1, endLine, len(lines))
		b.WriteString("\n" + helpStyle.Render(scrollInfo))
	}

	// Controls help
	var controls string
	if m.showMetrics && m.analysisData != nil && m.analysisData.EnhancedProjectAnalysis != nil {
		controls = "Controls: ↑↓ scroll, 1 overview, 2 detailed, 3 quality, 4 deps, m toggle, Esc back"
	} else {
		controls = "Controls: ↑↓/kj scroll, m metrics, s summary, Esc back"
	}
	b.WriteString("\n" + helpStyle.Render(controls))

	return b.String()
}

// SetContent sets the content to display
func (m *ContentViewModel) SetContent(filePath, content string) {
	m.filePath = filePath
	m.content = content
	m.scrollY = 0
}

// SetMetrics sets metrics content
func (m *ContentViewModel) SetMetrics(metrics string) {
	if m.showMetrics {
		m.content = metrics
	}
}

// SetSummary sets AI summary content
func (m *ContentViewModel) SetSummary(summary string) {
	if m.showSummary {
		m.content = summary
	}
}

// SetAnalysisData sets the analysis data for display
func (m *ContentViewModel) SetAnalysisData(data *AnalysisData) {
	m.analysisData = data
	if data != nil && data.ProjectAnalysis != nil {
		m.updateContentFromAnalysis()
	}
}

// updateContentFromAnalysis updates the content based on current view mode
func (m *ContentViewModel) updateContentFromAnalysis() {
	if m.analysisData == nil {
		return
	}

	if m.showMetrics {
		// Use enhanced metrics if available, otherwise fall back to basic
		if m.analysisData.EnhancedProjectAnalysis != nil {
			m.content = m.metricsDisplay.Render(m.analysisData.EnhancedProjectAnalysis, m.width, m.height)
		} else if m.analysisData.ProjectAnalysis != nil {
			m.content = m.formatAnalysisMetrics()
		}
	} else if m.showSummary && m.analysisData.Summary != "" {
		m.content = m.analysisData.Summary
	} else {
		// Show basic overview
		if m.analysisData.ProjectAnalysis != nil {
			m.content = m.formatAnalysisOverview()
		}
	}
	m.scrollY = 0 // Reset scroll when content changes
}

// formatAnalysisOverview formats the analysis overview
func (m *ContentViewModel) formatAnalysisOverview() string {
	analysis := m.analysisData.ProjectAnalysis
	var b strings.Builder

	// Header with styling
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render("📊 Codebase Analysis Results")
	
	b.WriteString(headerStyle + "\n")
	b.WriteString(strings.Repeat("═", 50) + "\n\n")

	// Project summary with better formatting
	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(summaryStyle.Render("📋 Project Summary") + "\n")
	b.WriteString(fmt.Sprintf("📁 Root Path: %s\n", analysis.RootPath))
	b.WriteString(fmt.Sprintf("📄 Total Files: %s\n", formatNumber(analysis.TotalFiles)))
	b.WriteString(fmt.Sprintf("📝 Total Lines: %s\n", formatNumber(analysis.TotalLines)))
	b.WriteString(fmt.Sprintf("⏱️  Analysis Duration: %v\n", analysis.AnalysisDuration.Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("🕒 Generated: %s\n\n", analysis.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Language breakdown with visual bars
	if len(analysis.Languages) > 0 {
		b.WriteString(summaryStyle.Render("🌐 Language Breakdown") + "\n")
		
		// Sort languages by line count for better display
		type langStat struct {
			name  string
			stats engine.LanguageStats
		}
		var sortedLangs []langStat
		for lang, stats := range analysis.Languages {
			sortedLangs = append(sortedLangs, langStat{lang, stats})
		}
		
		// Simple sort by line count (descending)
		for i := 0; i < len(sortedLangs)-1; i++ {
			for j := i + 1; j < len(sortedLangs); j++ {
				if sortedLangs[i].stats.LineCount < sortedLangs[j].stats.LineCount {
					sortedLangs[i], sortedLangs[j] = sortedLangs[j], sortedLangs[i]
				}
			}
		}
		
		for _, langStat := range sortedLangs {
			lang := langStat.name
			stats := langStat.stats
			
			percentage := float64(stats.LineCount) / float64(analysis.TotalLines) * 100
			
			// Language icon
			langIcon := getLangIcon(lang)
			
			b.WriteString(fmt.Sprintf("  %s %s (%.1f%%):\n", langIcon, lang, percentage))
			b.WriteString(fmt.Sprintf("    📄 Files: %s\n", formatNumber(stats.FileCount)))
			b.WriteString(fmt.Sprintf("    📝 Lines: %s\n", formatNumber(stats.LineCount)))
			b.WriteString(fmt.Sprintf("    🔧 Functions: %s\n", formatNumber(stats.FunctionCount)))
			b.WriteString(fmt.Sprintf("    🏗️  Classes: %s\n", formatNumber(stats.ClassCount)))
			
			if stats.Complexity > 0 {
				b.WriteString(fmt.Sprintf("    🧮 Complexity: %s\n", formatNumber(stats.Complexity)))
			}
			
			// Visual percentage bar
			barWidth := 30
			filled := int(float64(barWidth) * percentage / 100)
			bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
			b.WriteString(fmt.Sprintf("    [%s] %.1f%%\n", bar, percentage))
			b.WriteString("\n")
		}
	}

	// Top files by size/complexity
	if len(analysis.FileResults) > 0 {
		b.WriteString(summaryStyle.Render("📋 File Analysis Summary") + "\n")
		
		// Show top 10 files by line count
		count := len(analysis.FileResults)
		if count > 10 {
			count = 10
		}
		
		for i := 0; i < count; i++ {
			result := analysis.FileResults[i]
			langIcon := getLangIcon(result.Language)
			
			b.WriteString(fmt.Sprintf("  %s %s\n", langIcon, result.FilePath))
			b.WriteString(fmt.Sprintf("    📝 %s lines", formatNumber(result.LineCount)))
			
			if len(result.Functions) > 0 {
				b.WriteString(fmt.Sprintf(" • 🔧 %d functions", len(result.Functions)))
			}
			if len(result.Classes) > 0 {
				b.WriteString(fmt.Sprintf(" • 🏗️ %d classes", len(result.Classes)))
			}
			if result.Complexity > 0 {
				b.WriteString(fmt.Sprintf(" • 🧮 %d complexity", result.Complexity))
			}
			if len(result.Errors) > 0 {
				b.WriteString(fmt.Sprintf(" • ❌ %d errors", len(result.Errors)))
			}
			b.WriteString("\n")
		}
		
		if len(analysis.FileResults) > 10 {
			b.WriteString(fmt.Sprintf("  ... and %d more files\n", len(analysis.FileResults)-10))
		}
	}

	return b.String()
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

// formatAnalysisMetrics formats detailed metrics
func (m *ContentViewModel) formatAnalysisMetrics() string {
	analysis := m.analysisData.ProjectAnalysis
	var b strings.Builder

	b.WriteString("📈 Detailed Metrics\n")
	b.WriteString(strings.Repeat("=", 40) + "\n\n")

	// Overall statistics
	b.WriteString("📊 Overall Statistics:\n")
	b.WriteString(fmt.Sprintf("  Total Files: %d\n", analysis.TotalFiles))
	b.WriteString(fmt.Sprintf("  Total Lines: %d\n", analysis.TotalLines))

	totalFunctions := 0
	totalClasses := 0
	totalComplexity := 0
	for _, stats := range analysis.Languages {
		totalFunctions += stats.FunctionCount
		totalClasses += stats.ClassCount
		totalComplexity += stats.Complexity
	}

	b.WriteString(fmt.Sprintf("  Total Functions: %d\n", totalFunctions))
	b.WriteString(fmt.Sprintf("  Total Classes: %d\n", totalClasses))
	if totalComplexity > 0 {
		b.WriteString(fmt.Sprintf("  Total Complexity: %d\n", totalComplexity))
	}
	b.WriteString("\n")

	// Language-specific metrics
	for lang, stats := range analysis.Languages {
		b.WriteString(fmt.Sprintf("🔧 %s Metrics:\n", lang))
		b.WriteString(fmt.Sprintf("  Files: %d (%.1f%%)\n", 
			stats.FileCount, float64(stats.FileCount)/float64(analysis.TotalFiles)*100))
		b.WriteString(fmt.Sprintf("  Lines: %d (%.1f%%)\n", 
			stats.LineCount, float64(stats.LineCount)/float64(analysis.TotalLines)*100))
		b.WriteString(fmt.Sprintf("  Functions: %d\n", stats.FunctionCount))
		b.WriteString(fmt.Sprintf("  Classes: %d\n", stats.ClassCount))
		
		if stats.LineCount > 0 && stats.FunctionCount > 0 {
			avgLinesPerFunction := float64(stats.LineCount) / float64(stats.FunctionCount)
			b.WriteString(fmt.Sprintf("  Avg Lines/Function: %.1f\n", avgLinesPerFunction))
		}
		
		if stats.Complexity > 0 {
			b.WriteString(fmt.Sprintf("  Complexity: %d\n", stats.Complexity))
			if stats.FunctionCount > 0 {
				avgComplexity := float64(stats.Complexity) / float64(stats.FunctionCount)
				b.WriteString(fmt.Sprintf("  Avg Complexity/Function: %.1f\n", avgComplexity))
			}
		}
		b.WriteString("\n")
	}

	// File details
	if len(analysis.FileResults) > 0 {
		b.WriteString("📄 File Details:\n")
		for _, result := range analysis.FileResults {
			b.WriteString(fmt.Sprintf("  %s:\n", result.FilePath))
			b.WriteString(fmt.Sprintf("    Language: %s\n", result.Language))
			b.WriteString(fmt.Sprintf("    Lines: %d\n", result.LineCount))
			b.WriteString(fmt.Sprintf("    Functions: %d\n", len(result.Functions)))
			b.WriteString(fmt.Sprintf("    Classes: %d\n", len(result.Classes)))
			if result.Complexity > 0 {
				b.WriteString(fmt.Sprintf("    Complexity: %d\n", result.Complexity))
			}
			if len(result.Errors) > 0 {
				b.WriteString(fmt.Sprintf("    Errors: %d\n", len(result.Errors)))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

