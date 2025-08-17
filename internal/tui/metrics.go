package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/metrics"
)

// MetricsDisplayMode represents different ways to display metrics
type MetricsDisplayMode int

const (
	OverviewMode MetricsDisplayMode = iota
	DetailedMode
	QualityMode
	DependencyMode
)

// MetricsDisplay handles the display of enhanced metrics
type MetricsDisplay struct {
	mode         MetricsDisplayMode
	scrollY      int
	maxScroll    int
	width        int
	height       int
	sortBy       string
	filterLang   string
}

// NewMetricsDisplay creates a new metrics display
func NewMetricsDisplay() *MetricsDisplay {
	return &MetricsDisplay{
		mode:   OverviewMode,
		sortBy: "complexity",
	}
}

// SetMode sets the display mode
func (m *MetricsDisplay) SetMode(mode MetricsDisplayMode) {
	m.mode = mode
	m.scrollY = 0
}

// SetFilter sets the language filter
func (m *MetricsDisplay) SetFilter(lang string) {
	m.filterLang = lang
	m.scrollY = 0
}

// SetSort sets the sort criteria
func (m *MetricsDisplay) SetSort(sortBy string) {
	m.sortBy = sortBy
	m.scrollY = 0
}

// Scroll handles scrolling
func (m *MetricsDisplay) Scroll(delta int) {
	if m.maxScroll <= 0 || delta == 0 {
		return // No scrolling needed
	}
	
	oldScroll := m.scrollY
	newScroll := m.scrollY + delta
	if newScroll < 0 {
		newScroll = 0
	}
	if newScroll > m.maxScroll {
		newScroll = m.maxScroll
	}
	
	// Only update if position actually changed
	if newScroll != oldScroll {
		m.scrollY = newScroll
	}
}

// Render renders the metrics display
func (m *MetricsDisplay) Render(analysis *metrics.EnhancedProjectAnalysis, width, height int) string {
	m.width = width
	m.height = height

	if analysis == nil {
		return m.renderNoData()
	}

	switch m.mode {
	case OverviewMode:
		return m.renderOverview(analysis)
	case DetailedMode:
		return m.renderDetailed(analysis)
	case QualityMode:
		return m.renderQuality(analysis)
	case DependencyMode:
		return m.renderDependencies(analysis)
	default:
		return m.renderOverview(analysis)
	}
}

// renderNoData renders when no analysis data is available
func (m *MetricsDisplay) renderNoData() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)
	
	return style.Render("📊 No metrics data available\n\nRun analysis on a directory to see comprehensive metrics")
}

// renderOverview renders the overview mode
func (m *MetricsDisplay) renderOverview(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	
	b.WriteString(headerStyle.Render("📊 Project Metrics Overview") + "\n")
	b.WriteString(strings.Repeat("═", 60) + "\n\n")

	// Project summary
	b.WriteString(m.renderProjectSummary(analysis))
	b.WriteString("\n")

	// Quality score
	b.WriteString(m.renderQualityScore(analysis.QualityScore))
	b.WriteString("\n")

	// Language breakdown
	b.WriteString(m.renderLanguageBreakdown(analysis.Languages))
	b.WriteString("\n")

	// Top metrics
	b.WriteString(m.renderTopMetrics(analysis))

	return m.applyScrolling(b.String())
}

// renderDetailed renders the detailed mode
func (m *MetricsDisplay) renderDetailed(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	
	b.WriteString(headerStyle.Render("📈 Detailed Metrics Analysis") + "\n")
	b.WriteString(strings.Repeat("═", 60) + "\n\n")

	// Project metrics
	b.WriteString(m.renderProjectMetrics(analysis.ProjectMetrics))
	b.WriteString("\n")

	// Directory breakdown
	b.WriteString(m.renderDirectoryStats(analysis.DirectoryStats))
	b.WriteString("\n")

	// Language details
	b.WriteString(m.renderLanguageDetails(analysis.Languages))

	return m.applyScrolling(b.String())
}

// renderQuality renders the quality mode
func (m *MetricsDisplay) renderQuality(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	
	b.WriteString(headerStyle.Render("🏆 Code Quality Analysis") + "\n")
	b.WriteString(strings.Repeat("═", 60) + "\n\n")

	// Quality score breakdown
	b.WriteString(m.renderQualityBreakdown(analysis.QualityScore))
	b.WriteString("\n")

	// Technical debt analysis
	b.WriteString(m.renderTechnicalDebt(analysis.ProjectMetrics))
	b.WriteString("\n")

	// Maintainability insights
	b.WriteString(m.renderMaintainabilityInsights(analysis))

	return m.applyScrolling(b.String())
}

// renderDependencies renders the dependency mode
func (m *MetricsDisplay) renderDependencies(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	
	b.WriteString(headerStyle.Render("🔗 Dependency Analysis") + "\n")
	b.WriteString(strings.Repeat("═", 60) + "\n\n")

	// Dependency overview
	b.WriteString(m.renderDependencyOverview(analysis.DependencyGraph))
	b.WriteString("\n")

	// Internal dependencies
	if len(analysis.DependencyGraph.InternalDependencies) > 0 {
		b.WriteString(m.renderInternalDependencies(analysis.DependencyGraph))
		b.WriteString("\n")
	}

	// External dependencies
	if len(analysis.DependencyGraph.ExternalDependencies) > 0 {
		b.WriteString(m.renderExternalDependencies(analysis.DependencyGraph))
		b.WriteString("\n")
	}

	// Circular dependencies
	if len(analysis.DependencyGraph.CircularDependencies) > 0 {
		b.WriteString(m.renderCircularDependencies(analysis.DependencyGraph))
	}

	return m.applyScrolling(b.String())
}

// renderProjectSummary renders the project summary
func (m *MetricsDisplay) renderProjectSummary(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("📋 Project Summary") + "\n")
	b.WriteString(fmt.Sprintf("📁 Root Path: %s\n", analysis.RootPath))
	b.WriteString(fmt.Sprintf("📄 Total Files: %s\n", formatNumber(analysis.TotalFiles)))
	b.WriteString(fmt.Sprintf("📝 Total Lines: %s\n", formatNumber(analysis.TotalLines)))
	b.WriteString(fmt.Sprintf("🌐 Languages: %d\n", len(analysis.Languages)))

	return b.String()
}

// renderQualityScore renders the quality score
func (m *MetricsDisplay) renderQualityScore(score metrics.QualityScore) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🏆 Quality Score") + "\n")
	
	// Grade with color
	gradeColor := m.getGradeColor(score.Grade)
	gradeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(gradeColor)).
		Bold(true)
	
	b.WriteString(fmt.Sprintf("📊 Overall: %.1f%% (%s)\n", 
		score.Overall, gradeStyle.Render(score.Grade)))
	b.WriteString(fmt.Sprintf("🔧 Maintainability: %.1f%%\n", score.Maintainability))
	b.WriteString(fmt.Sprintf("🧮 Complexity: %.1f%%\n", score.Complexity))
	b.WriteString(fmt.Sprintf("📚 Documentation: %.1f%%\n", score.Documentation))
	b.WriteString(fmt.Sprintf("🧪 Test Coverage: %.1f%%\n", score.TestCoverage))
	b.WriteString(fmt.Sprintf("📋 Code Duplication: %.1f%%\n", score.CodeDuplication))

	return b.String()
}

// renderLanguageBreakdown renders language statistics
func (m *MetricsDisplay) renderLanguageBreakdown(languages map[string]metrics.LanguageStats) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🌐 Language Breakdown") + "\n")

	// Sort languages by line count
	type langStat struct {
		name  string
		stats metrics.LanguageStats
	}
	var sortedLangs []langStat
	totalLines := 0
	
	for lang, stats := range languages {
		sortedLangs = append(sortedLangs, langStat{lang, stats})
		totalLines += stats.LineCount
	}
	
	sort.Slice(sortedLangs, func(i, j int) bool {
		return sortedLangs[i].stats.LineCount > sortedLangs[j].stats.LineCount
	})

	for _, langStat := range sortedLangs {
		lang := langStat.name
		stats := langStat.stats
		
		percentage := float64(stats.LineCount) / float64(totalLines) * 100
		langIcon := getLangIcon(lang)
		
		b.WriteString(fmt.Sprintf("  %s %s (%.1f%%):\n", langIcon, lang, percentage))
		b.WriteString(fmt.Sprintf("    📄 Files: %s • 📝 Lines: %s • 🔧 Functions: %s\n", 
			formatNumber(stats.FileCount), formatNumber(stats.LineCount), formatNumber(stats.FunctionCount)))
		
		if stats.AverageComplexity > 0 {
			b.WriteString(fmt.Sprintf("    🧮 Avg Complexity: %.1f • 🏆 Maintainability: %.1f%%\n", 
				stats.AverageComplexity, stats.MaintainabilityIndex))
		}
		
		// Visual bar
		barWidth := 30
		filled := int(float64(barWidth) * percentage / 100)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		b.WriteString(fmt.Sprintf("    [%s]\n", bar))
		b.WriteString("\n")
	}

	return b.String()
}

// renderTopMetrics renders top metrics
func (m *MetricsDisplay) renderTopMetrics(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🔝 Key Metrics") + "\n")
	b.WriteString(fmt.Sprintf("🧮 Total Complexity: %s\n", formatNumber(analysis.ProjectMetrics.TotalComplexity)))
	b.WriteString(fmt.Sprintf("📊 Average Complexity: %.1f\n", analysis.ProjectMetrics.AverageComplexity))
	b.WriteString(fmt.Sprintf("⚠️  Max Complexity: %s\n", formatNumber(analysis.ProjectMetrics.MaxComplexity)))
	b.WriteString(fmt.Sprintf("🏗️  Technical Debt: %.1f\n", analysis.ProjectMetrics.TechnicalDebt))
	b.WriteString(fmt.Sprintf("📚 Documentation Ratio: %.1f%%\n", analysis.ProjectMetrics.DocumentationRatio))

	return b.String()
}

// renderProjectMetrics renders detailed project metrics
func (m *MetricsDisplay) renderProjectMetrics(metrics metrics.ProjectMetrics) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("📊 Project Metrics") + "\n")
	b.WriteString(fmt.Sprintf("🧮 Total Complexity: %s\n", formatNumber(metrics.TotalComplexity)))
	b.WriteString(fmt.Sprintf("📊 Average Complexity: %.2f\n", metrics.AverageComplexity))
	b.WriteString(fmt.Sprintf("⚠️  Maximum Complexity: %s\n", formatNumber(metrics.MaxComplexity)))
	b.WriteString(fmt.Sprintf("🏗️  Technical Debt Score: %.2f\n", metrics.TechnicalDebt))
	b.WriteString(fmt.Sprintf("🔧 Maintainability Index: %.2f%%\n", metrics.MaintainabilityIndex))
	b.WriteString(fmt.Sprintf("📚 Documentation Ratio: %.2f%%\n", metrics.DocumentationRatio))
	b.WriteString(fmt.Sprintf("💬 Code to Comment Ratio: %.2f:1\n", metrics.CodeToCommentRatio))
	b.WriteString(fmt.Sprintf("🧪 Test Coverage: %.2f%%\n", metrics.TestCoverage))
	b.WriteString(fmt.Sprintf("📋 Code Duplication: %.2f%%\n", metrics.CodeDuplication))

	return b.String()
}

// renderDirectoryStats renders directory statistics
func (m *MetricsDisplay) renderDirectoryStats(dirStats map[string]metrics.DirectoryStats) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("📁 Directory Analysis") + "\n")

	// Sort directories by complexity
	type dirStat struct {
		path  string
		stats metrics.DirectoryStats
	}
	var sortedDirs []dirStat
	
	for path, stats := range dirStats {
		sortedDirs = append(sortedDirs, dirStat{path, stats})
	}
	
	sort.Slice(sortedDirs, func(i, j int) bool {
		return sortedDirs[i].stats.Complexity > sortedDirs[j].stats.Complexity
	})

	for _, dirStat := range sortedDirs {
		path := dirStat.path
		stats := dirStat.stats
		
		b.WriteString(fmt.Sprintf("📁 %s:\n", path))
		b.WriteString(fmt.Sprintf("  📄 Files: %s • 📝 Lines: %s • 🧮 Complexity: %s\n", 
			formatNumber(stats.FileCount), formatNumber(stats.LineCount), formatNumber(stats.Complexity)))
		b.WriteString(fmt.Sprintf("  🏗️  Maintainability: %.1f%%\n", stats.MaintainabilityIndex))
		b.WriteString("\n")
	}

	return b.String()
}

// renderLanguageDetails renders detailed language statistics
func (m *MetricsDisplay) renderLanguageDetails(languages map[string]metrics.LanguageStats) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🌐 Language Details") + "\n")

	for lang, stats := range languages {
		langIcon := getLangIcon(lang)
		b.WriteString(fmt.Sprintf("%s %s:\n", langIcon, lang))
		b.WriteString(fmt.Sprintf("  📄 Files: %s\n", formatNumber(stats.FileCount)))
		b.WriteString(fmt.Sprintf("  📝 Lines: %s (Code: %s, Comments: %s, Blank: %s)\n", 
			formatNumber(stats.LineCount), formatNumber(stats.CodeLines), 
			formatNumber(stats.CommentLines), formatNumber(stats.BlankLines)))
		b.WriteString(fmt.Sprintf("  🔧 Functions: %s • 🏗️ Classes: %s\n", 
			formatNumber(stats.FunctionCount), formatNumber(stats.ClassCount)))
		b.WriteString(fmt.Sprintf("  🧮 Complexity: %s (Avg: %.1f, Max: %s)\n", 
			formatNumber(stats.Complexity), stats.AverageComplexity, formatNumber(stats.MaxComplexity)))
		b.WriteString(fmt.Sprintf("  🏗️  Maintainability: %.1f%% • 🏗️ Technical Debt: %.1f\n", 
			stats.MaintainabilityIndex, stats.TechnicalDebt))
		if stats.TestFiles > 0 {
			b.WriteString(fmt.Sprintf("  🧪 Test Files: %s • Coverage: %.1f%%\n", 
				formatNumber(stats.TestFiles), stats.TestCoverage))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderQualityBreakdown renders quality score breakdown
func (m *MetricsDisplay) renderQualityBreakdown(score metrics.QualityScore) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🏆 Quality Score Breakdown") + "\n")

	// Overall score with visual bar
	gradeColor := m.getGradeColor(score.Grade)
	gradeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(gradeColor)).
		Bold(true)
	
	b.WriteString(fmt.Sprintf("📊 Overall Score: %.1f%% (%s)\n", 
		score.Overall, gradeStyle.Render(score.Grade)))
	
	// Individual metrics with bars
	metrics := []struct {
		name  string
		value float64
		icon  string
	}{
		{"Maintainability", score.Maintainability, "🔧"},
		{"Complexity", score.Complexity, "🧮"},
		{"Documentation", score.Documentation, "📚"},
		{"Test Coverage", score.TestCoverage, "🧪"},
		{"Code Duplication", score.CodeDuplication, "📋"},
	}

	for _, metric := range metrics {
		bar := m.renderProgressBar(metric.value, 30)
		b.WriteString(fmt.Sprintf("%s %s: %.1f%% %s\n", 
			metric.icon, metric.name, metric.value, bar))
	}

	return b.String()
}

// renderTechnicalDebt renders technical debt analysis
func (m *MetricsDisplay) renderTechnicalDebt(projectMetrics metrics.ProjectMetrics) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🏗️ Technical Debt Analysis") + "\n")
	b.WriteString(fmt.Sprintf("💰 Total Technical Debt: %.2f\n", projectMetrics.TechnicalDebt))
	
	// Debt level assessment
	debtLevel := "Low"
	debtColor := "#00FF00"
	if projectMetrics.TechnicalDebt > 50 {
		debtLevel = "High"
		debtColor = "#FF0000"
	} else if projectMetrics.TechnicalDebt > 20 {
		debtLevel = "Medium"
		debtColor = "#FFA500"
	}
	
	debtStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(debtColor)).
		Bold(true)
	
	b.WriteString(fmt.Sprintf("📊 Debt Level: %s\n", debtStyle.Render(debtLevel)))
	
	// Recommendations
	b.WriteString("\n💡 Recommendations:\n")
	if projectMetrics.TechnicalDebt > 50 {
		b.WriteString("  • High technical debt detected - consider refactoring\n")
		b.WriteString("  • Focus on reducing complexity in high-complexity functions\n")
		b.WriteString("  • Improve code documentation\n")
	} else if projectMetrics.TechnicalDebt > 20 {
		b.WriteString("  • Moderate technical debt - monitor and improve gradually\n")
		b.WriteString("  • Consider adding more tests\n")
	} else {
		b.WriteString("  • Low technical debt - good code quality!\n")
		b.WriteString("  • Maintain current practices\n")
	}

	return b.String()
}

// renderMaintainabilityInsights renders maintainability insights
func (m *MetricsDisplay) renderMaintainabilityInsights(analysis *metrics.EnhancedProjectAnalysis) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🔧 Maintainability Insights") + "\n")
	b.WriteString(fmt.Sprintf("📊 Project Maintainability: %.1f%%\n", analysis.ProjectMetrics.MaintainabilityIndex))
	
	// Language maintainability comparison
	b.WriteString("\n🌐 By Language:\n")
	for lang, stats := range analysis.Languages {
		langIcon := getLangIcon(lang)
		b.WriteString(fmt.Sprintf("  %s %s: %.1f%%\n", langIcon, lang, stats.MaintainabilityIndex))
	}

	return b.String()
}

// renderDependencyOverview renders dependency overview
func (m *MetricsDisplay) renderDependencyOverview(graph metrics.DependencyGraph) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🔗 Dependency Overview") + "\n")
	b.WriteString(fmt.Sprintf("🏠 Internal Dependencies: %d\n", len(graph.InternalDependencies)))
	b.WriteString(fmt.Sprintf("🌐 External Dependencies: %d\n", len(graph.ExternalDependencies)))
	b.WriteString(fmt.Sprintf("🔄 Circular Dependencies: %d\n", len(graph.CircularDependencies)))
	b.WriteString(fmt.Sprintf("📊 Dependency Depth: %d\n", graph.DependencyDepth))
	b.WriteString(fmt.Sprintf("🗑️  Unused Dependencies: %d\n", len(graph.UnusedDependencies)))

	return b.String()
}

// renderInternalDependencies renders internal dependencies
func (m *MetricsDisplay) renderInternalDependencies(graph metrics.DependencyGraph) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🏠 Internal Dependencies") + "\n")
	
	count := 0
	for file, deps := range graph.InternalDependencies {
		if count >= 10 { // Limit display
			b.WriteString(fmt.Sprintf("  ... and %d more files\n", len(graph.InternalDependencies)-10))
			break
		}
		b.WriteString(fmt.Sprintf("📄 %s (%d deps)\n", file, len(deps)))
		count++
	}

	return b.String()
}

// renderExternalDependencies renders external dependencies
func (m *MetricsDisplay) renderExternalDependencies(graph metrics.DependencyGraph) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🌐 External Dependencies") + "\n")
	
	count := 0
	for file, deps := range graph.ExternalDependencies {
		if count >= 10 { // Limit display
			b.WriteString(fmt.Sprintf("  ... and %d more files\n", len(graph.ExternalDependencies)-10))
			break
		}
		b.WriteString(fmt.Sprintf("📄 %s (%d deps)\n", file, len(deps)))
		count++
	}

	return b.String()
}

// renderCircularDependencies renders circular dependencies
func (m *MetricsDisplay) renderCircularDependencies(graph metrics.DependencyGraph) string {
	var b strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)
	
	b.WriteString(sectionStyle.Render("🔄 Circular Dependencies (Issues)") + "\n")
	
	for i, cycle := range graph.CircularDependencies {
		if i >= 5 { // Limit display
			b.WriteString(fmt.Sprintf("  ... and %d more cycles\n", len(graph.CircularDependencies)-5))
			break
		}
		b.WriteString(fmt.Sprintf("⚠️  Cycle %d: %s\n", i+1, strings.Join(cycle, " → ")))
	}

	return b.String()
}

// Helper functions

// getGradeColor returns color for quality grade
func (m *MetricsDisplay) getGradeColor(grade string) string {
	switch grade {
	case "A":
		return "#00FF00"
	case "B":
		return "#7FFF00"
	case "C":
		return "#FFFF00"
	case "D":
		return "#FFA500"
	case "F":
		return "#FF0000"
	default:
		return "#CCCCCC"
	}
}

// renderProgressBar renders a progress bar
func (m *MetricsDisplay) renderProgressBar(value float64, width int) string {
	filled := int(float64(width) * value / 100)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s]", bar)
}

// applyScrolling applies scrolling to content
func (m *MetricsDisplay) applyScrolling(content string) string {
	if m.height <= 4 {
		return content // Not enough space for scrolling
	}
	
	lines := strings.Split(content, "\n")
	availableHeight := m.height - 4 // Reserve space for controls
	
	// Only update maxScroll if it's not set or content changed significantly
	newMaxScroll := max(0, len(lines)-availableHeight)
	if m.maxScroll != newMaxScroll {
		m.maxScroll = newMaxScroll
		// Ensure scroll position is still valid
		if m.scrollY > m.maxScroll {
			m.scrollY = m.maxScroll
		}
	}
	
	if m.maxScroll == 0 {
		return content
	}
	
	startLine := m.scrollY
	if startLine >= len(lines) {
		startLine = max(0, len(lines)-availableHeight)
		m.scrollY = startLine
	}
	if startLine < 0 {
		startLine = 0
		m.scrollY = 0
	}
	
	endLine := min(len(lines), startLine+availableHeight)
	
	if startLine >= endLine || startLine >= len(lines) {
		return content // Return original content if scroll position is invalid
	}
	
	visibleLines := lines[startLine:endLine]
	result := strings.Join(visibleLines, "\n")
	
	// Add scroll indicator only if we have scrollable content
	if m.maxScroll > 0 && len(lines) > availableHeight {
		scrollInfo := fmt.Sprintf("\n\n📊 Line %d-%d of %d (↑↓ to scroll)", 
			startLine+1, min(endLine, len(lines)), len(lines))
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true).
			Render(scrollInfo)
	}
	
	return result
}

// Helper functions are defined in other files