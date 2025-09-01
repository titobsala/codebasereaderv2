package views

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/tui"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
	"github.com/tito-sala/codebasereaderv2/internal/tui/shared"
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
	analysisData   *shared.AnalysisData
	metricsDisplay *components.MetricsDisplay
	currentMode    components.MetricsDisplayMode

	// Performance optimization: cache split lines
	cachedLines   []string
	cachedContent string
}

// NewContentViewModel creates a new content view model
func NewContentViewModel() *ContentViewModel {
	return &ContentViewModel{
		content:        "",
		scrollY:        0,
		maxScroll:      0,
		showMetrics:    false,
		showSummary:    false,
		metricsDisplay: components.NewMetricsDisplay(),
		currentMode:    components.OverviewMode,
	}
}

// Init initializes the content view
func (m *ContentViewModel) Init() tea.Cmd {
	return nil
}

// GetMetricsDisplay returns the metrics display (for testing)
func (m *ContentViewModel) GetMetricsDisplay() *components.MetricsDisplay {
	return m.metricsDisplay
}

// GetAnalysisData returns the analysis data (for testing)
func (m *ContentViewModel) GetAnalysisData() *shared.AnalysisData {
	return m.analysisData
}

// SetAnalysisData sets the analysis data (for testing)
func (m *ContentViewModel) SetAnalysisData(data *shared.AnalysisData) {
	m.analysisData = data
	// When analysis data is set, enable metrics view by default
	if data != nil {
		m.showMetrics = true
		m.UpdateContentFromAnalysis()
	}
}

// ShowMetrics returns whether metrics are shown (for testing)
func (m *ContentViewModel) ShowMetrics() bool {
	return m.showMetrics
}

// SetShowMetrics sets whether metrics are shown (for testing)
func (m *ContentViewModel) SetShowMetrics(show bool) {
	m.showMetrics = show
}

// GetContent returns the current content (for testing)
func (m *ContentViewModel) GetContent() string {
	return m.content
}

// Update handles messages for the content view
func (m *ContentViewModel) Update(msg tea.Msg) (*ContentViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Scrolling commands - unified for both metrics and regular content
		case "up", "k":
			m.scrollUp(1)
		case "down", "j":
			m.scrollDown(1)
		case "pgup":
			m.scrollUp(10)
		case "pgdown":
			m.scrollDown(10)
		case "home", "g":
			m.scrollToTop()
		case "end", "G":
			m.scrollToBottom()
		case "ctrl+u":
			// Scroll up half page
			m.scrollUp(m.height / 2)
		case "ctrl+d":
			// Scroll down half page
			m.scrollDown(m.height / 2)

		// Metrics mode switching (only when metrics are shown) - using 6-9 for terminal compatibility
		case "6":
			if m.showMetrics && m.analysisData != nil {
				m.currentMode = components.OverviewMode
				if m.metricsDisplay != nil {
					m.metricsDisplay.SetMode(components.OverviewMode)
				}
				m.UpdateContentFromAnalysis()
			}
		case "7":
			if m.showMetrics && m.analysisData != nil {
				m.currentMode = components.DetailedMode
				if m.metricsDisplay != nil {
					m.metricsDisplay.SetMode(components.DetailedMode)
				}
				m.UpdateContentFromAnalysis()
			}
		case "8":
			if m.showMetrics && m.analysisData != nil {
				m.currentMode = components.QualityMode
				if m.metricsDisplay != nil {
					m.metricsDisplay.SetMode(components.QualityMode)
				}
				m.UpdateContentFromAnalysis()
			}
		case "9":
			if m.showMetrics && m.analysisData != nil {
				m.currentMode = components.DependencyMode
				if m.metricsDisplay != nil {
					m.metricsDisplay.SetMode(components.DependencyMode)
				}
				m.UpdateContentFromAnalysis()
			}
		}
	}
	return m, nil
}

// ToggleMetrics toggles metrics display
func (m *ContentViewModel) ToggleMetrics() {
	m.showMetrics = !m.showMetrics
	m.showSummary = false
	m.UpdateContentFromAnalysis()
}

// ToggleSummary toggles summary display
func (m *ContentViewModel) ToggleSummary() {
	m.showSummary = !m.showSummary
	m.showMetrics = false
	m.UpdateContentFromAnalysis()
}

// ClearAnalysis clears analysis data
func (m *ContentViewModel) ClearAnalysis() {
	m.analysisData = nil
	m.content = "Select a file to view its content"
	m.filePath = ""
	m.scrollY = 0
	m.showMetrics = false
	m.showSummary = false
}

// ShowSummary returns the current showSummary state
func (m *ContentViewModel) ShowSummary() bool {
	return m.showSummary
}

// View renders the content view
func (m *ContentViewModel) View(width, height int) string {
	m.width = width
	m.height = height

	var b strings.Builder

	// Header
	if m.filePath != "" {
		header := fmt.Sprintf("File: %s", m.filePath)
		b.WriteString(components.SelectedStyle.Render(header) + "\n\n")
	}

	// Content area
	contentHeight := height - 4 // Reserve space for header and controls

	// Handle different content scenarios
	if m.analysisData == nil && m.filePath == "" && !m.showMetrics && !m.showSummary {
		// No file selected and no analysis data - show the analysis dashboard
		dashboardContent := m.renderAnalysisDashboard(width, contentHeight)
		b.WriteString(dashboardContent)
		return b.String()
	} else if m.content == "" {
		b.WriteString(components.HelpStyle.Render("No content to display"))
		return b.String()
	}

	// Get cached lines for better performance
	lines := m.getCachedLines()
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
		b.WriteString("\n" + components.HelpStyle.Render(scrollInfo))
	}

	// Controls help
	var controls string
	if m.showMetrics && m.analysisData != nil {
		// Show metric navigation controls when metrics are active
		controls = "Controls: ↑↓ scroll, 6-9 switch modes (6=overview 7=detailed 8=quality 9=deps), r reset, 1-4 tabs"
	} else if m.analysisData != nil {
		// Show available modes when analysis data exists but metrics not shown
		controls = "Controls: ↑↓ scroll, m metrics, s summary, r reset, 1-4 tabs"
	} else {
		// No analysis data available
		controls = "Controls: ↑↓ scroll, press 'a' in Explorer to analyze directory, 1-4 tabs"
	}
	b.WriteString("\n" + components.HelpStyle.Render(controls))

	return b.String()
}

// getCachedLines returns cached split lines, updating cache if content changed
func (m *ContentViewModel) getCachedLines() []string {
	if m.content != m.cachedContent {
		m.cachedLines = strings.Split(m.content, "\n")
		m.cachedContent = m.content
	}
	return m.cachedLines
}

// SetContent sets the content to display
func (m *ContentViewModel) SetContent(filePath, content string) {
	m.filePath = filePath
	m.content = content
	m.scrollY = 0
	// Clear cache to force update on next access
	m.cachedContent = ""
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

// updateContentFromAnalysis updates the content based on current view mode
func (m *ContentViewModel) UpdateContentFromAnalysis() {
	if m.analysisData == nil {
		// Clear content so the dashboard will be shown
		m.content = ""
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
	// Clear cache to force update on next access
	m.cachedContent = ""
}

// formatAnalysisOverview formats the analysis overview
func (m *ContentViewModel) formatAnalysisOverview() string {
	analysis := m.analysisData.ProjectAnalysis
	var b strings.Builder

	// Header with styling
	header := components.HeaderStyle.Render("📊 Codebase Analysis Results")

	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("=", 50) + "\n\n")

	// Project summary with better formatting
	b.WriteString(components.SectionStyle.Render("📋 Project Summary") + "\n")
	b.WriteString(fmt.Sprintf("📁 Root Path: %s\n", analysis.RootPath))
	b.WriteString(fmt.Sprintf("📄 Total Files: %s\n", tui.FormatNumber(analysis.TotalFiles)))
	b.WriteString(fmt.Sprintf("📝 Total Lines: %s\n", tui.FormatNumber(analysis.TotalLines)))
	b.WriteString(fmt.Sprintf("⏱️  Analysis Duration: %v\n", analysis.AnalysisDuration.Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("🕒 Generated: %s\n\n", analysis.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Language breakdown with visual bars
	if len(analysis.Languages) > 0 {
		b.WriteString(components.SectionStyle.Render("🌐 Language Breakdown") + "\n")

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
			langIcon := tui.GetLangIcon(lang)

			b.WriteString(fmt.Sprintf("  %s %s (%.1f%%):\n", langIcon, lang, percentage))
			b.WriteString(fmt.Sprintf("    📄 Files: %s\n", tui.FormatNumber(stats.FileCount)))
			b.WriteString(fmt.Sprintf("    📝 Lines: %s\n", tui.FormatNumber(stats.LineCount)))
			b.WriteString(fmt.Sprintf("    🔧 Functions: %s\n", tui.FormatNumber(stats.FunctionCount)))
			b.WriteString(fmt.Sprintf("    🏗️  Classes: %s\n", tui.FormatNumber(stats.ClassCount)))

			if stats.Complexity > 0 {
				b.WriteString(fmt.Sprintf("    🧮 Complexity: %s\n", tui.FormatNumber(stats.Complexity)))
			}

			// Enhanced visual percentage bar
			bar := components.CreateProgressBar(percentage, 30, true)
			b.WriteString(fmt.Sprintf("    %s\n", bar))
			b.WriteString("\n")
		}
	}

	// Top files by size/complexity
	if len(analysis.FileResults) > 0 {
		b.WriteString(components.SectionStyle.Render("📋 File Analysis Summary") + "\n")

		// Show top 10 files by line count
		count := len(analysis.FileResults)
		if count > 10 {
			count = 10
		}

		for i := 0; i < count; i++ {
			result := analysis.FileResults[i]
			langIcon := tui.GetLangIcon(result.Language)

			b.WriteString(fmt.Sprintf("  %s %s\n", langIcon, result.FilePath))
			b.WriteString(fmt.Sprintf("    📝 %s lines", tui.FormatNumber(result.LineCount)))

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

// renderAnalysisDashboard renders the main analysis dashboard when no file is selected
func (m *ContentViewModel) renderAnalysisDashboard(width, height int) string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Underline(true)

	b.WriteString(titleStyle.Render("📊 Analysis Dashboard") + "\n\n")

	// Status section
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		MarginLeft(2)

	if m.analysisData == nil {
		b.WriteString(statusStyle.Render("📋 Status: No Analysis Data") + "\n")
		b.WriteString(descStyle.Render("No codebase has been analyzed yet.") + "\n\n")

		b.WriteString(statusStyle.Render("🚀 Getting Started") + "\n")
		b.WriteString(descStyle.Render("1. Navigate to the Explorer tab (Tab 1 or press 1)") + "\n")
		b.WriteString(descStyle.Render("2. Select a directory") + "\n")
		b.WriteString(descStyle.Render("3. Press 'a' to analyze the directory") + "\n")
		b.WriteString(descStyle.Render("4. Return to this tab to view the results") + "\n\n")
	} else {
		b.WriteString(statusStyle.Render("✅ Status: Analysis Complete") + "\n")
		b.WriteString(descStyle.Render("Analysis data is available. Use the controls below to explore.") + "\n\n")
	}

	// Navigation menu
	b.WriteString(m.renderNavigationMenu())

	// Additional help
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		MarginTop(2)

	b.WriteString(helpStyle.Render("💡 Tips:") + "\n")
	b.WriteString(helpStyle.Render("• Use Tab or Shift+Tab to cycle between tabs") + "\n")
	b.WriteString(helpStyle.Render("• Press ? for comprehensive help") + "\n")
	b.WriteString(helpStyle.Render("• Press Esc to return to Explorer") + "\n")

	return b.String()
}

// renderNavigationMenu renders the vertical navigation menu for the analysis tab
func (m *ContentViewModel) renderNavigationMenu() string {
	var b strings.Builder

	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		MarginLeft(2)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	availableStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#66BB6A"))

	unavailableStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Strikethrough(true)

	b.WriteString(menuStyle.Render("📋 Analysis Views") + "\n")

	if m.analysisData != nil {
		// Analysis data available - show all options as available
		b.WriteString(itemStyle.Render(keyStyle.Render("6")+" → "+availableStyle.Render("Overview")+" (Project summary & key metrics)") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("7")+" → "+availableStyle.Render("Detailed")+" (File-level analysis & statistics)") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("8")+" → "+availableStyle.Render("Quality")+" (Code quality scores & insights)") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("9")+" → "+availableStyle.Render("Dependencies")+" (Dependency graph & analysis)") + "\n\n")

		b.WriteString(menuStyle.Render("🎯 Quick Actions") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("m")+" → "+availableStyle.Render("Toggle metrics view")) + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("r")+" → "+availableStyle.Render("Reset to overview")) + "\n")

		// Show current mode if metrics are enabled
		if m.showMetrics {
			currentModeText := ""
			switch m.currentMode {
			case components.OverviewMode:
				currentModeText = "Overview"
			case components.DetailedMode:
				currentModeText = "Detailed"
			case components.QualityMode:
				currentModeText = "Quality"
			case components.DependencyMode:
				currentModeText = "Dependencies"
			}

			currentStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFEB3B")).
				Bold(true)

			b.WriteString(itemStyle.Render("🎯 "+currentStyle.Render("Current: "+currentModeText+" Mode")) + "\n")
		}
	} else {
		// No analysis data - show options as unavailable
		b.WriteString(itemStyle.Render(keyStyle.Render("6")+" → "+unavailableStyle.Render("Overview")+" (requires analysis)") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("7")+" → "+unavailableStyle.Render("Detailed")+" (requires analysis)") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("8")+" → "+unavailableStyle.Render("Quality")+" (requires analysis)") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("9")+" → "+unavailableStyle.Render("Dependencies")+" (requires analysis)") + "\n\n")

		b.WriteString(menuStyle.Render("🎯 Available Actions") + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("1")+" → "+availableStyle.Render("Go to Explorer tab")) + "\n")
		b.WriteString(itemStyle.Render(keyStyle.Render("a")+" → "+availableStyle.Render("Analyze directory")+" (from Explorer)") + "\n")
	}

	b.WriteString(itemStyle.Render(keyStyle.Render("1-5")+" → Switch between main tabs") + "\n\n")

	return b.String()
}
