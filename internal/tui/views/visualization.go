package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/metrics"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
)

// VisualizationMode represents different visualization types
type VisualizationMode int

const (
	DependencyTreeMode VisualizationMode = iota
	ComplexityHeatmapMode
	LanguageCompositionMode
	QualityGaugesMode
	TechnicalDebtMode
	FunctionUsageMode
)

// VisualizationViewModel handles the visualization system
type VisualizationViewModel struct {
	currentMode   VisualizationMode
	modes         []VisualizationModeInfo
	scrollY       int
	maxScroll     int
	width         int
	height        int
	analysisData  *AnalysisData
	filterOptions FilterOptions
}

// VisualizationModeInfo contains information about a visualization mode
type VisualizationModeInfo struct {
	Name        string
	Icon        string
	Description string
	ShortKey    string
}

// FilterOptions contains filtering options for visualizations
type FilterOptions struct {
	Language      string
	MinComplexity int
	MaxComplexity int
	ShowTestFiles bool
	FilePattern   string
}

// AnalysisData wraps the analysis data for the visualization
type AnalysisData struct {
	EnhancedProjectAnalysis *metrics.EnhancedProjectAnalysis
	Summary                 string
}

// NewVisualizationViewModel creates a new visualization view model
func NewVisualizationViewModel() *VisualizationViewModel {
	return &VisualizationViewModel{
		currentMode: DependencyTreeMode,
		modes:       createVisualizationModes(),
		scrollY:     0,
		maxScroll:   0,
		filterOptions: FilterOptions{
			Language:      "",
			MinComplexity: 0,
			MaxComplexity: 100,
			ShowTestFiles: true,
			FilePattern:   "*",
		},
	}
}

// createVisualizationModes creates all available visualization modes
func createVisualizationModes() []VisualizationModeInfo {
	return []VisualizationModeInfo{
		{
			Name:        "Dependency Tree",
			Icon:        "🌳",
			Description: "Function call trees and import relationships",
			ShortKey:    "1",
		},
		{
			Name:        "Complexity Heatmap",
			Icon:        "🔥",
			Description: "Visual complexity distribution across files",
			ShortKey:    "2",
		},
		{
			Name:        "Language Composition",
			Icon:        "🎨",
			Description: "Programming language breakdown charts",
			ShortKey:    "3",
		},
		{
			Name:        "Quality Gauges",
			Icon:        "🏆",
			Description: "Code quality metrics visualization",
			ShortKey:    "4",
		},
		{
			Name:        "Technical Debt",
			Icon:        "⚠️",
			Description: "Technical debt distribution and hotspots",
			ShortKey:    "5",
		},
		{
			Name:        "Function Usage",
			Icon:        "📊",
			Description: "Function call frequency and usage patterns",
			ShortKey:    "6",
		},
	}
}

// SetAnalysisData sets the analysis data for visualization
func (v *VisualizationViewModel) SetAnalysisData(data *AnalysisData) {
	v.analysisData = data
	v.scrollY = 0
}

// GetCurrentMode returns the current visualization mode
func (v *VisualizationViewModel) GetCurrentMode() VisualizationMode {
	return v.currentMode
}

// SetMode sets the current visualization mode
func (v *VisualizationViewModel) SetMode(mode VisualizationMode) {
	if mode >= 0 && int(mode) < len(v.modes) {
		v.currentMode = mode
		v.scrollY = 0
	}
}

// Update handles navigation within the visualization system
func (v *VisualizationViewModel) Update(key string) {
	switch key {
	case "left", "h":
		if v.currentMode > 0 {
			v.currentMode--
			v.scrollY = 0
		}
	case "right", "l":
		if int(v.currentMode) < len(v.modes)-1 {
			v.currentMode++
			v.scrollY = 0
		}
	case "1":
		v.SetMode(DependencyTreeMode)
	case "2":
		v.SetMode(ComplexityHeatmapMode)
	case "3":
		v.SetMode(LanguageCompositionMode)
	case "4":
		v.SetMode(QualityGaugesMode)
	case "5":
		v.SetMode(TechnicalDebtMode)
	case "6":
		v.SetMode(FunctionUsageMode)
	case "up", "k":
		if v.scrollY > 0 {
			v.scrollY--
		}
	case "down", "j":
		if v.scrollY < v.maxScroll {
			v.scrollY++
		}
	case "pgup":
		v.scrollY = max(0, v.scrollY-10)
	case "pgdown":
		v.scrollY = min(v.maxScroll, v.scrollY+10)
	case "home", "g":
		v.scrollY = 0
	case "end", "G":
		v.scrollY = v.maxScroll
	}
}

// Render renders the visualization view
func (v *VisualizationViewModel) Render(width, height int) string {
	v.width = width
	v.height = height

	var b strings.Builder

	// Header with current mode indicator
	currentMode := v.modes[v.currentMode]
	headerText := fmt.Sprintf("%s %s (%s)", currentMode.Icon, currentMode.Name, currentMode.ShortKey)
	header := components.HeaderStyle.Render(headerText)
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("=", width-4) + "\n\n")

	// Mode navigation tabs
	b.WriteString(v.renderModeTabs() + "\n\n")

	// Main visualization content
	visualizationContent := v.renderCurrentMode()
	scrolledContent := v.applyScrolling(visualizationContent)
	b.WriteString(scrolledContent)

	// Footer with navigation hints
	footer := v.renderFooter()
	b.WriteString("\n" + footer)

	return b.String()
}

// renderModeTabs renders the mode selection tabs
func (v *VisualizationViewModel) renderModeTabs() string {
	var tabs []string

	for i, mode := range v.modes {
		var style lipgloss.Style
		if i == int(v.currentMode) {
			style = lipgloss.NewStyle().
				Foreground(components.NeutralWhite).
				Background(components.PrimaryPurple).
				Padding(0, 1).
				Bold(true)
		} else {
			style = lipgloss.NewStyle().
				Foreground(components.NeutralMedium).
				Padding(0, 1)
		}

		tabText := fmt.Sprintf("%s %s", mode.ShortKey, mode.Icon)
		tabs = append(tabs, style.Render(tabText))
	}

	return strings.Join(tabs, " ")
}

// renderCurrentMode renders the currently selected visualization mode
func (v *VisualizationViewModel) renderCurrentMode() string {
	if v.analysisData == nil || v.analysisData.EnhancedProjectAnalysis == nil {
		return v.renderNoData()
	}

	switch v.currentMode {
	case DependencyTreeMode:
		return v.renderDependencyTree()
	case ComplexityHeatmapMode:
		return v.renderComplexityHeatmap()
	case LanguageCompositionMode:
		return v.renderLanguageComposition()
	case QualityGaugesMode:
		return v.renderQualityGauges()
	case TechnicalDebtMode:
		return v.renderTechnicalDebt()
	case FunctionUsageMode:
		return v.renderFunctionUsage()
	default:
		return v.renderNoData()
	}
}

// renderNoData renders when no analysis data is available
func (v *VisualizationViewModel) renderNoData() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Width(v.width).
		Height(v.height - 10)

	return style.Render("📊 No analysis data available\n\nRun analysis on a directory to see visualizations")
}

// renderDependencyTree renders dependency tree visualization
func (v *VisualizationViewModel) renderDependencyTree() string {
	var b strings.Builder

	analysis := v.analysisData.EnhancedProjectAnalysis
	deps := analysis.DependencyGraph

	b.WriteString("🌳 Dependency Tree Visualization\n\n")

	// Internal dependencies tree
	if len(deps.InternalDependencies) > 0 {
		b.WriteString("📦 Internal Dependencies:\n")
		count := 0
		for file, fileDeps := range deps.InternalDependencies {
			if count >= 15 { // Limit display for performance
				b.WriteString(fmt.Sprintf("   ... and %d more files\n", len(deps.InternalDependencies)-15))
				break
			}
			b.WriteString(fmt.Sprintf("├── %s\n", file))
			for i, dep := range fileDeps {
				if i >= 5 { // Limit deps per file
					b.WriteString(fmt.Sprintf("│   └── ... and %d more\n", len(fileDeps)-5))
					break
				}
				if i == len(fileDeps)-1 && len(fileDeps) <= 5 {
					b.WriteString(fmt.Sprintf("│   └── %s\n", dep))
				} else {
					b.WriteString(fmt.Sprintf("│   ├── %s\n", dep))
				}
			}
			count++
		}
		b.WriteString("\n")
	}

	// Circular dependencies (critical issues)
	if len(deps.CircularDependencies) > 0 {
		b.WriteString("⚠️  Circular Dependencies (Critical Issues):\n")
		for i, cycle := range deps.CircularDependencies {
			if i >= 5 {
				b.WriteString(fmt.Sprintf("   ... and %d more cycles\n", len(deps.CircularDependencies)-5))
				break
			}
			b.WriteString(fmt.Sprintf("🔄 %s\n", strings.Join(cycle, " → ")))
		}
		b.WriteString("\n")
	}

	// Dependency statistics
	b.WriteString("📊 Dependency Statistics:\n")
	b.WriteString(fmt.Sprintf("• Internal Dependencies: %d files\n", len(deps.InternalDependencies)))
	b.WriteString(fmt.Sprintf("• External Dependencies: %d files\n", len(deps.ExternalDependencies)))
	b.WriteString(fmt.Sprintf("• Standard Library: %d files\n", len(deps.StandardDependencies)))
	b.WriteString(fmt.Sprintf("• Circular Dependencies: %d cycles\n", len(deps.CircularDependencies)))
	b.WriteString(fmt.Sprintf("• Dependency Depth: %d levels\n", deps.DependencyDepth))

	return v.applyScrolling(b.String())
}

// renderComplexityHeatmap renders complexity heatmap visualization
func (v *VisualizationViewModel) renderComplexityHeatmap() string {
	var b strings.Builder

	analysis := v.analysisData.EnhancedProjectAnalysis

	b.WriteString("🔥 Complexity Heatmap\n\n")

	// Create a visual complexity distribution
	b.WriteString("📊 Complexity Distribution by Directory:\n")

	// Sort directories by complexity for better visualization
	type dirComplexity struct {
		path       string
		complexity int
		files      int
	}

	var dirs []dirComplexity
	for path, stats := range analysis.DirectoryStats {
		dirs = append(dirs, dirComplexity{
			path:       path,
			complexity: stats.Complexity,
			files:      stats.FileCount,
		})
	}

	// Sort by complexity descending
	for i := 0; i < len(dirs)-1; i++ {
		for j := i + 1; j < len(dirs); j++ {
			if dirs[i].complexity < dirs[j].complexity {
				dirs[i], dirs[j] = dirs[j], dirs[i]
			}
		}
	}

	// Render complexity bars
	maxComplexity := 0
	if len(dirs) > 0 {
		maxComplexity = dirs[0].complexity
	}

	for i, dir := range dirs {
		if i >= 20 { // Limit display
			b.WriteString(fmt.Sprintf("... and %d more directories\n", len(dirs)-20))
			break
		}

		// Calculate bar length (max 40 characters)
		barLength := 40
		if maxComplexity > 0 {
			barLength = int(float64(dir.complexity) / float64(maxComplexity) * 40)
		}

		// Create color-coded bar based on complexity level
		var barChar string
		var barColor lipgloss.Color

		complexityRatio := float64(dir.complexity) / float64(maxComplexity)
		if complexityRatio > 0.8 {
			barChar = "█"
			barColor = lipgloss.Color("#FF0000") // Red for high complexity
		} else if complexityRatio > 0.6 {
			barChar = "▓"
			barColor = lipgloss.Color("#FF8800") // Orange for medium-high
		} else if complexityRatio > 0.4 {
			barChar = "▒"
			barColor = lipgloss.Color("#FFFF00") // Yellow for medium
		} else {
			barChar = "░"
			barColor = lipgloss.Color("#00FF00") // Green for low
		}

		bar := strings.Repeat(barChar, barLength)
		coloredBar := lipgloss.NewStyle().Foreground(barColor).Render(bar)

		b.WriteString(fmt.Sprintf("%-30s %s %d (%d files)\n",
			dir.path, coloredBar, dir.complexity, dir.files))
	}

	b.WriteString("\n📈 Complexity Legend:\n")
	b.WriteString("🟥 High (>80%):    Requires immediate attention\n")
	b.WriteString("🟧 Medium-High (60-80%): Consider refactoring\n")
	b.WriteString("🟨 Medium (40-60%): Monitor for growth\n")
	b.WriteString("🟩 Low (<40%):     Well-structured code\n")

	return v.applyScrolling(b.String())
}

// renderLanguageComposition renders language composition charts
func (v *VisualizationViewModel) renderLanguageComposition() string {
	var b strings.Builder

	analysis := v.analysisData.EnhancedProjectAnalysis

	b.WriteString("🎨 Language Composition\n\n")

	// Calculate total lines for percentages
	totalLines := 0
	for _, stats := range analysis.Languages {
		totalLines += stats.LineCount
	}

	b.WriteString("📊 Lines of Code by Language:\n")

	// Sort languages by line count
	type langStat struct {
		name  string
		stats metrics.LanguageStats
	}
	var sortedLangs []langStat
	for lang, stats := range analysis.Languages {
		sortedLangs = append(sortedLangs, langStat{lang, stats})
	}

	// Sort by line count descending
	for i := 0; i < len(sortedLangs)-1; i++ {
		for j := i + 1; j < len(sortedLangs); j++ {
			if sortedLangs[i].stats.LineCount < sortedLangs[j].stats.LineCount {
				sortedLangs[i], sortedLangs[j] = sortedLangs[j], sortedLangs[i]
			}
		}
	}

	// Render language composition bars
	for _, langStat := range sortedLangs {
		lang := langStat.name
		stats := langStat.stats

		percentage := float64(stats.LineCount) / float64(totalLines) * 100
		barLength := int(percentage * 50 / 100) // Max 50 chars

		// Language-specific colors/icons
		var langIcon string
		var barColor lipgloss.Color

		switch strings.ToLower(lang) {
		case "go":
			langIcon = "🐹"
			barColor = lipgloss.Color("#00ADD8")
		case "python":
			langIcon = "🐍"
			barColor = lipgloss.Color("#3776AB")
		case "javascript":
			langIcon = "🟨"
			barColor = lipgloss.Color("#F7DF1E")
		case "typescript":
			langIcon = "🔷"
			barColor = lipgloss.Color("#3178C6")
		case "java":
			langIcon = "☕"
			barColor = lipgloss.Color("#ED8B00")
		default:
			langIcon = "📄"
			barColor = lipgloss.Color("#888888")
		}

		bar := strings.Repeat("█", barLength)
		if barLength == 0 {
			bar = "▏" // Show something even for very small percentages
		}

		coloredBar := lipgloss.NewStyle().Foreground(barColor).Render(bar)

		b.WriteString(fmt.Sprintf("%s %-12s %s %.1f%% (%d lines, %d files)\n",
			langIcon, lang, coloredBar, percentage, stats.LineCount, stats.FileCount))
	}

	b.WriteString("\n🔢 Detailed Statistics:\n")
	for _, langStat := range sortedLangs {
		lang := langStat.name
		stats := langStat.stats

		b.WriteString(fmt.Sprintf("• %s:\n", lang))
		b.WriteString(fmt.Sprintf("  - Functions: %d (avg complexity: %.1f)\n", stats.FunctionCount, stats.AverageComplexity))
		b.WriteString(fmt.Sprintf("  - Classes: %d\n", stats.ClassCount))
		b.WriteString(fmt.Sprintf("  - Code/Comment/Blank: %d/%d/%d lines\n",
			stats.CodeLines, stats.CommentLines, stats.BlankLines))
		if stats.TestFiles > 0 {
			b.WriteString(fmt.Sprintf("  - Test files: %d (coverage: %.1f%%)\n", stats.TestFiles, stats.TestCoverage))
		}
		b.WriteString("\n")
	}

	return v.applyScrolling(b.String())
}

// renderQualityGauges renders quality metrics as gauges
func (v *VisualizationViewModel) renderQualityGauges() string {
	var b strings.Builder

	analysis := v.analysisData.EnhancedProjectAnalysis
	quality := analysis.QualityScore

	b.WriteString("🏆 Code Quality Gauges\n\n")

	// Overall quality score with large gauge
	b.WriteString("📊 Overall Quality Score:\n")
	overallGauge := v.createGauge(quality.Overall, 60)
	gradeColor := v.getGradeColor(quality.Grade)
	gradeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(gradeColor)).Bold(true)

	b.WriteString(fmt.Sprintf("%s\n", overallGauge))
	b.WriteString(fmt.Sprintf("Score: %.1f%% (Grade: %s)\n\n", quality.Overall, gradeStyle.Render(quality.Grade)))

	// Individual quality metrics
	qualityMetrics := []struct {
		name  string
		value float64
		icon  string
		desc  string
	}{
		{"Maintainability", quality.Maintainability, "🔧", "How easy it is to modify and extend"},
		{"Complexity", quality.Complexity, "🧮", "Code complexity and readability"},
		{"Documentation", quality.Documentation, "📚", "Code documentation coverage"},
		{"Test Coverage", quality.TestCoverage, "🧪", "Automated test coverage"},
		{"Code Duplication", quality.CodeDuplication, "📋", "Amount of duplicated code"},
	}

	for _, metric := range qualityMetrics {
		gauge := v.createGauge(metric.value, 40)
		b.WriteString(fmt.Sprintf("%s %s: %.1f%%\n", metric.icon, metric.name, metric.value))
		b.WriteString(fmt.Sprintf("%s\n", gauge))
		b.WriteString(fmt.Sprintf("   %s\n\n", metric.desc))
	}

	// Quality insights and recommendations
	b.WriteString("💡 Quality Insights:\n")
	if quality.Overall >= 90 {
		b.WriteString("✅ Excellent code quality! Keep up the great work.\n")
	} else if quality.Overall >= 75 {
		b.WriteString("👍 Good code quality with room for improvement.\n")
		if quality.Complexity < 70 {
			b.WriteString("• Consider reducing complexity in high-complexity functions\n")
		}
		if quality.Documentation < 70 {
			b.WriteString("• Improve code documentation coverage\n")
		}
		if quality.TestCoverage < 70 {
			b.WriteString("• Add more comprehensive test coverage\n")
		}
	} else {
		b.WriteString("⚠️  Code quality needs attention. Focus on:\n")
		if quality.Maintainability < 60 {
			b.WriteString("• Improving maintainability through refactoring\n")
		}
		if quality.Complexity < 60 {
			b.WriteString("• Reducing complexity in critical functions\n")
		}
		if quality.Documentation < 50 {
			b.WriteString("• Adding comprehensive documentation\n")
		}
		if quality.TestCoverage < 50 {
			b.WriteString("• Implementing thorough test coverage\n")
		}
	}

	return v.applyScrolling(b.String())
}

// renderTechnicalDebt renders technical debt visualization
func (v *VisualizationViewModel) renderTechnicalDebt() string {
	var b strings.Builder

	analysis := v.analysisData.EnhancedProjectAnalysis
	projectMetrics := analysis.ProjectMetrics

	b.WriteString("⚠️  Technical Debt Analysis\n\n")

	// Overall debt score
	b.WriteString("💰 Technical Debt Overview:\n")
	debtScore := projectMetrics.TechnicalDebt

	// Create debt level visualization
	debtGauge := v.createDebtGauge(debtScore)
	b.WriteString(fmt.Sprintf("%s\n", debtGauge))

	debtLevel, debtColor := v.getDebtLevel(debtScore)
	debtStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(debtColor)).Bold(true)
	b.WriteString(fmt.Sprintf("Debt Score: %.1f (%s)\n\n", debtScore, debtStyle.Render(debtLevel)))

	// Debt breakdown by language
	b.WriteString("🌐 Technical Debt by Language:\n")
	for lang, stats := range analysis.Languages {
		langIcon := v.getLangIcon(lang)
		debtBar := v.createMiniGauge(stats.TechnicalDebt, 20)
		b.WriteString(fmt.Sprintf("%s %-12s %s %.1f\n", langIcon, lang, debtBar, stats.TechnicalDebt))
	}
	b.WriteString("\n")

	// Debt sources analysis
	b.WriteString("🔍 Debt Sources Analysis:\n")
	avgComplexity := projectMetrics.AverageComplexity
	maxComplexity := projectMetrics.MaxComplexity
	maintainabilityIndex := projectMetrics.MaintainabilityIndex

	if avgComplexity > 10 {
		b.WriteString("🔥 High Average Complexity: Functions are too complex\n")
		b.WriteString("   → Consider breaking down large functions\n")
	}

	if maxComplexity > 25 {
		b.WriteString("⚡ Extremely Complex Functions: Some functions are very complex\n")
		b.WriteString("   → Identify and refactor the most complex functions\n")
	}

	if maintainabilityIndex < 60 {
		b.WriteString("🔧 Low Maintainability: Code is hard to maintain\n")
		b.WriteString("   → Focus on improving code structure and readability\n")
	}

	if projectMetrics.CodeDuplication > 15 {
		b.WriteString("📋 High Code Duplication: Significant code repetition\n")
		b.WriteString("   → Extract common functionality into reusable components\n")
	}

	b.WriteString("\n💡 Debt Reduction Recommendations:\n")
	if debtScore > 50 {
		b.WriteString("🚨 Critical: Immediate action required\n")
		b.WriteString("• Allocate 30-40% of development time to debt reduction\n")
		b.WriteString("• Focus on highest complexity functions first\n")
		b.WriteString("• Implement code review processes\n")
	} else if debtScore > 20 {
		b.WriteString("⚠️  Moderate: Plan debt reduction activities\n")
		b.WriteString("• Allocate 15-20% of development time to refactoring\n")
		b.WriteString("• Address debt incrementally during feature development\n")
		b.WriteString("• Improve test coverage to prevent regression\n")
	} else {
		b.WriteString("✅ Low: Maintain current quality practices\n")
		b.WriteString("• Continue current development practices\n")
		b.WriteString("• Monitor for debt accumulation\n")
		b.WriteString("• Regular code quality assessments\n")
	}

	return v.applyScrolling(b.String())
}

// renderFunctionUsage renders function usage patterns
func (v *VisualizationViewModel) renderFunctionUsage() string {
	var b strings.Builder

	analysis := v.analysisData.EnhancedProjectAnalysis

	b.WriteString("📊 Function Usage Analysis\n\n")

	// Function statistics summary
	totalFunctions := 0
	totalClasses := 0
	for _, stats := range analysis.Languages {
		totalFunctions += stats.FunctionCount
		totalClasses += stats.ClassCount
	}

	b.WriteString("📈 Function Statistics:\n")
	b.WriteString(fmt.Sprintf("• Total Functions: %d\n", totalFunctions))
	b.WriteString(fmt.Sprintf("• Total Classes: %d\n", totalClasses))
	if totalFunctions > 0 {
		avgComplexityPerFunc := float64(analysis.ProjectMetrics.TotalComplexity) / float64(totalFunctions)
		b.WriteString(fmt.Sprintf("• Average Complexity per Function: %.1f\n", avgComplexityPerFunc))
	}
	b.WriteString("\n")

	// Function complexity distribution
	b.WriteString("🧮 Function Complexity Distribution by Language:\n")
	for lang, stats := range analysis.Languages {
		if stats.FunctionCount == 0 {
			continue
		}

		langIcon := v.getLangIcon(lang)

		// Create complexity distribution visualization
		complexityRatio := stats.AverageComplexity / 20.0 // Normalize to reasonable scale
		if complexityRatio > 1.0 {
			complexityRatio = 1.0
		}

		complexityBar := v.createComplexityBar(complexityRatio, 25)

		b.WriteString(fmt.Sprintf("%s %-12s %s Functions: %d, Avg Complexity: %.1f\n",
			langIcon, lang, complexityBar, stats.FunctionCount, stats.AverageComplexity))
	}
	b.WriteString("\n")

	// Function size analysis (lines of code per function)
	b.WriteString("📏 Function Size Analysis:\n")
	for lang, stats := range analysis.Languages {
		if stats.FunctionCount == 0 {
			continue
		}

		avgLinesPerFunction := float64(stats.CodeLines) / float64(stats.FunctionCount)
		sizeCategory := "Small"
		sizeIcon := "🟢"

		if avgLinesPerFunction > 50 {
			sizeCategory = "Large"
			sizeIcon = "🔴"
		} else if avgLinesPerFunction > 25 {
			sizeCategory = "Medium"
			sizeIcon = "🟡"
		}

		b.WriteString(fmt.Sprintf("%s %s: %.1f lines/function (%s)\n",
			sizeIcon, lang, avgLinesPerFunction, sizeCategory))
	}
	b.WriteString("\n")

	// Quality recommendations based on function patterns
	b.WriteString("💡 Function Quality Recommendations:\n")

	highComplexityLangs := 0
	for _, stats := range analysis.Languages {
		if stats.AverageComplexity > 10 {
			highComplexityLangs++
		}
	}

	if highComplexityLangs > 0 {
		b.WriteString("⚠️  High Complexity Detected:\n")
		b.WriteString("• Break down complex functions into smaller, focused functions\n")
		b.WriteString("• Consider extracting common logic into utility functions\n")
		b.WriteString("• Use early returns to reduce nesting levels\n")
	}

	if totalFunctions > 500 {
		b.WriteString("📊 Large Codebase:\n")
		b.WriteString("• Consider organizing functions into modules/packages\n")
		b.WriteString("• Implement consistent naming conventions\n")
		b.WriteString("• Use static analysis tools for function dependency tracking\n")
	}

	// Mock function usage tree (since we don't have actual call graph data yet)
	b.WriteString("\n🌳 Sample Function Call Tree:\n")
	b.WriteString("├── main()\n")
	b.WriteString("│   ├── initializeApp()\n")
	b.WriteString("│   ├── processFiles()\n")
	b.WriteString("│   │   ├── walkDirectory()\n")
	b.WriteString("│   │   └── analyzeFile()\n")
	b.WriteString("│   │       ├── parseCode()\n")
	b.WriteString("│   │       └── calculateMetrics()\n")
	b.WriteString("│   └── generateReport()\n")
	b.WriteString("\n")
	b.WriteString("📝 Note: Detailed function call graphs require additional\n")
	b.WriteString("    static analysis and will be available in future updates.\n")

	return v.applyScrolling(b.String())
}

// Helper methods for visualization rendering

// createGauge creates a visual gauge for percentage values
func (v *VisualizationViewModel) createGauge(value float64, width int) string {
	filled := int(value * float64(width) / 100)
	empty := width - filled

	var fillChar, emptyChar string
	var fillColor lipgloss.Color

	if value >= 80 {
		fillChar = "█"
		fillColor = lipgloss.Color("#00FF00") // Green
	} else if value >= 60 {
		fillChar = "▓"
		fillColor = lipgloss.Color("#FFFF00") // Yellow
	} else if value >= 40 {
		fillChar = "▒"
		fillColor = lipgloss.Color("#FF8800") // Orange
	} else {
		fillChar = "░"
		fillColor = lipgloss.Color("#FF0000") // Red
	}

	emptyChar = "░"

	filledPart := lipgloss.NewStyle().Foreground(fillColor).Render(strings.Repeat(fillChar, filled))
	emptyPart := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render(strings.Repeat(emptyChar, empty))

	return fmt.Sprintf("[%s%s] %.1f%%", filledPart, emptyPart, value)
}

// createDebtGauge creates a specialized gauge for technical debt
func (v *VisualizationViewModel) createDebtGauge(debt float64) string {
	// Normalize debt to 0-100 scale (assuming 100 is very high debt)
	normalizedDebt := debt
	if normalizedDebt > 100 {
		normalizedDebt = 100
	}

	return v.createGauge(normalizedDebt, 50)
}

// createMiniGauge creates a small gauge for compact display
func (v *VisualizationViewModel) createMiniGauge(value float64, width int) string {
	filled := int(value * float64(width) / 100)
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s]", bar)
}

// createComplexityBar creates a complexity-colored bar
func (v *VisualizationViewModel) createComplexityBar(ratio float64, width int) string {
	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}

	var fillColor lipgloss.Color
	if ratio > 0.8 {
		fillColor = lipgloss.Color("#FF0000") // Red for high complexity
	} else if ratio > 0.6 {
		fillColor = lipgloss.Color("#FF8800") // Orange for medium-high
	} else if ratio > 0.4 {
		fillColor = lipgloss.Color("#FFFF00") // Yellow for medium
	} else {
		fillColor = lipgloss.Color("#00FF00") // Green for low
	}

	filledPart := lipgloss.NewStyle().Foreground(fillColor).Render(strings.Repeat("█", filled))
	emptyPart := strings.Repeat("░", width-filled)

	return fmt.Sprintf("[%s%s]", filledPart, emptyPart)
}

// getGradeColor returns color for quality grade
func (v *VisualizationViewModel) getGradeColor(grade string) string {
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

// getDebtLevel returns debt level and color
func (v *VisualizationViewModel) getDebtLevel(debt float64) (string, string) {
	if debt > 50 {
		return "Critical", "#FF0000"
	} else if debt > 20 {
		return "High", "#FF8800"
	} else if debt > 10 {
		return "Medium", "#FFFF00"
	} else {
		return "Low", "#00FF00"
	}
}

// getLangIcon returns an icon for a programming language
func (v *VisualizationViewModel) getLangIcon(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		return "🐹"
	case "python":
		return "🐍"
	case "javascript":
		return "🟨"
	case "typescript":
		return "🔷"
	case "java":
		return "☕"
	case "c":
		return "🔧"
	case "cpp", "c++":
		return "⚡"
	case "rust":
		return "🦀"
	case "php":
		return "🐘"
	case "ruby":
		return "💎"
	default:
		return "📄"
	}
}

// renderFooter renders the visualization footer with navigation hints
func (v *VisualizationViewModel) renderFooter() string {
	navigation := "Navigate: ←→/hl (modes) • ↑↓/kj (scroll) • 1-6 (jump) • f (filter)"
	return components.HelpStyle.
		Align(lipgloss.Center).
		Render(navigation)
}

// applyScrolling applies scrolling to visualization content
func (v *VisualizationViewModel) applyScrolling(content string) string {
	if v.height <= 8 {
		return content
	}

	lines := strings.Split(content, "\n")
	availableHeight := v.height - 8 // Reserve space for header, tabs, footer

	newMaxScroll := max(0, len(lines)-availableHeight)
	if v.maxScroll != newMaxScroll {
		v.maxScroll = newMaxScroll
		if v.scrollY > v.maxScroll {
			v.scrollY = v.maxScroll
		}
	}

	if v.maxScroll == 0 {
		return content
	}

	startLine := v.scrollY
	if startLine >= len(lines) {
		startLine = max(0, len(lines)-availableHeight)
		v.scrollY = startLine
	}

	endLine := min(len(lines), startLine+availableHeight)
	visibleLines := lines[startLine:endLine]
	result := strings.Join(visibleLines, "\n")

	// Add scroll indicator
	if v.maxScroll > 0 {
		scrollInfo := fmt.Sprintf("\n%s",
			components.ScrollInfoStyle.Render(fmt.Sprintf("📊 Showing lines %d-%d of %d (↑↓ to scroll)",
				startLine+1, min(endLine, len(lines)), len(lines))))
		result += scrollInfo
	}

	return result
}
