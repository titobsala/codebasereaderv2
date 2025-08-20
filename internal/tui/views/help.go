package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
)

// HelpViewModel handles the comprehensive help system
type HelpViewModel struct {
	currentSection int
	sections       []HelpSection
	scrollY        int
	maxScroll      int
	width          int
	height         int
}

// HelpSection represents a section in the help system
type HelpSection struct {
	Title       string
	Icon        string
	Content     string
	KeyBindings []KeyBinding
	Tips        []string
}

// KeyBinding represents a keyboard shortcut with description
type KeyBinding struct {
	Keys        []string
	Description string
	Context     string
}

// NewHelpViewModel creates a new help view model
func NewHelpViewModel() *HelpViewModel {
	return &HelpViewModel{
		currentSection: 0,
		sections:       createHelpSections(),
		scrollY:        0,
		maxScroll:      0,
	}
}

// createHelpSections creates all help sections with comprehensive documentation
func createHelpSections() []HelpSection {
	return []HelpSection{
		{
			Title: "Overview",
			Icon:  "📖",
			Content: `Welcome to CodebaseReader v2 - A powerful TUI for codebase analysis!

This tool provides comprehensive static analysis of your code projects with:
• Language-specific metrics and complexity analysis
• Dependency graph visualization and analysis  
• Code quality scoring and maintainability insights
• Interactive file tree navigation and content viewing
• Multiple visualization modes for different analysis aspects
• Export functionality for analysis results

Navigate through this help using the arrow keys or number shortcuts (1-8).`,
			KeyBindings: []KeyBinding{
				{[]string{"←", "→", "h", "l"}, "Navigate help sections", "Help View"},
				{[]string{"1-8"}, "Jump to specific help section", "Help View"},
				{[]string{"?"}, "Toggle help on/off", "Global"},
				{[]string{"Esc"}, "Return to previous view", "Help View"},
			},
		},
		{
			Title: "Navigation",
			Icon:  "🧭",
			Content: `Master the navigation system to efficiently browse your codebase:

The interface is organized into tabs that you can navigate between:
• Explorer Tab: File tree and directory navigation
• Analysis Tab: Metrics, quality scores, and analysis results  
• Configuration Tab: Settings and preferences
• Help Tab: This comprehensive help system

Use the tabbed interface or keyboard shortcuts to switch between views seamlessly.`,
			KeyBindings: []KeyBinding{
				{[]string{"↑", "↓", "k", "j"}, "Move up/down in lists and content", "All Views"},
				{[]string{"←", "→", "h", "l"}, "Navigate horizontally, collapse/expand", "File Tree"},
				{[]string{"Enter", "Space"}, "Select item or toggle expansion", "File Tree"},
				{[]string{"Tab", "Shift+Tab"}, "Switch between tabs", "Global"},
				{[]string{"1", "2", "3", "4"}, "Jump to specific tabs", "Global"},
				{[]string{"PgUp", "PgDn"}, "Scroll content by page", "Content Areas"},
				{[]string{"Home", "End", "g", "G"}, "Jump to start/end", "Content Areas"},
				{[]string{"Ctrl+u", "Ctrl+d"}, "Scroll half page up/down", "Content Areas"},
			},
		},
		{
			Title: "Analysis Features",
			Icon:  "📊",
			Content: `Comprehensive codebase analysis with multiple visualization modes:

DEPENDENCY ANALYSIS:
• Internal Dependencies: Project-internal imports and relationships
• External Dependencies: Third-party libraries and frameworks  
• Standard Library: Language standard library usage
• Circular Dependencies: Detected circular import issues

METRICS ANALYSIS:
• Lines of Code: Total, code lines, comments, blank lines
• Complexity Metrics: Cyclomatic complexity, maintainability index
• Function/Class Counts: Detailed breakdown by language
• Technical Debt: Calculated debt score and recommendations

QUALITY SCORING:
• Overall Quality Grade: A-F grading system
• Maintainability Score: Code maintainability percentage
• Documentation Ratio: Comment-to-code ratio analysis
• Test Coverage: Test file detection and coverage metrics`,
			KeyBindings: []KeyBinding{
				{[]string{"a"}, "Analyze selected directory", "File Tree"},
				{[]string{"r"}, "Refresh analysis", "Global"},
				{[]string{"c"}, "Clear current analysis", "Global"},
				{[]string{"m"}, "Toggle metrics view", "Content View"},
				{[]string{"s"}, "Toggle summary view", "Content View"},
			},
		},
		{
			Title: "Metrics Modes",
			Icon:  "📈",
			Content: `Four specialized metrics visualization modes (use 1-4 keys in Analysis tab):

MODE 1 - OVERVIEW:
• Project summary with key statistics
• Quality score with visual grade indicator
• Language breakdown with percentage bars
• Top-level metrics at a glance

MODE 2 - DETAILED ANALYSIS:
• Comprehensive project metrics breakdown
• Directory-level statistics and complexity
• Per-language detailed analysis
• Average metrics and ratios

MODE 3 - QUALITY ASSESSMENT:
• Quality score breakdown by category
• Technical debt analysis with recommendations
• Maintainability insights by language
• Code quality best practices suggestions

MODE 4 - DEPENDENCY GRAPH:
• Internal dependency relationships
• External library dependencies  
• Standard library usage patterns
• Circular dependency detection and warnings

Each mode provides scrollable content with detailed insights.`,
			KeyBindings: []KeyBinding{
				{[]string{"1"}, "Overview mode - project summary", "Analysis Tab"},
				{[]string{"2"}, "Detailed mode - comprehensive metrics", "Analysis Tab"},
				{[]string{"3"}, "Quality mode - quality assessment", "Analysis Tab"},
				{[]string{"4"}, "Dependency mode - dependency analysis", "Analysis Tab"},
				{[]string{"↑", "↓"}, "Scroll within metrics view", "Metrics Active"},
			},
		},
		{
			Title: "File Operations",
			Icon:  "📁",
			Content: `Efficient file and directory operations:

FILE TREE NAVIGATION:
• Hierarchical display of project structure
• Language-specific file icons and indicators
• Size indicators for files and directories
• Support for hidden files (configurable)

FILE CONTENT VIEWING:
• Syntax-aware content display
• Line numbers and content formatting
• Scrollable content with navigation indicators
• File metadata and statistics

DIRECTORY OPERATIONS:
• Expand/collapse directory trees
• Recursive analysis of subdirectories
• Selective analysis of specific directories
• Project root detection and management`,
			KeyBindings: []KeyBinding{
				{[]string{"Enter"}, "Open file or expand directory", "File Tree"},
				{[]string{"Space"}, "Toggle directory expansion", "File Tree"},
				{[]string{"←", "h"}, "Collapse directory or go up", "File Tree"},
				{[]string{"→", "l"}, "Expand directory or go down", "File Tree"},
				{[]string{"Ctrl+a"}, "Select all in current directory", "File Tree"},
			},
		},
		{
			Title: "Configuration",
			Icon:  "⚙️",
			Content: `Customize the application to your preferences:

DISPLAY SETTINGS:
• Color scheme selection (default, dark, light themes)
• Hidden file visibility toggle
• Language-specific syntax highlighting
• Progress bar and animation preferences

ANALYSIS SETTINGS:
• Worker thread configuration for performance
• File type inclusion/exclusion patterns
• Analysis depth and recursion limits
• Timeout settings for large projects

EXPORT SETTINGS:
• Output format selection (JSON, CSV, XML)
• Report detail level configuration
• Custom template support
• Automatic export location settings

All configuration changes take effect immediately without restart.`,
			KeyBindings: []KeyBinding{
				{[]string{"Enter"}, "Modify configuration option", "Config View"},
				{[]string{"Space"}, "Toggle boolean settings", "Config View"},
				{[]string{"→", "←"}, "Adjust numeric values", "Config View"},
				{[]string{"r"}, "Reset to defaults", "Config View"},
			},
		},
		{
			Title: "Export & Reports",
			Icon:  "📤",
			Content: `Generate and export comprehensive analysis reports:

EXPORT FORMATS:
• JSON: Machine-readable structured data
• CSV: Spreadsheet-compatible tabular data
• XML: Structured markup for integration
• HTML: Web-viewable formatted reports

REPORT CONTENTS:
• Executive summary with key metrics
• Detailed analysis by language and directory
• Quality assessment with recommendations
• Dependency graph with relationships
• Historical comparison (when available)

EXPORT OPTIONS:
• Full report with all metrics
• Summary report with key insights only
• Custom report with selected sections
• Raw data export for further analysis

Reports are automatically timestamped and saved to your specified output directory.`,
			KeyBindings: []KeyBinding{
				{[]string{"e"}, "Export current analysis", "Global"},
				{[]string{"Ctrl+e"}, "Export with custom options", "Global"},
				{[]string{"Shift+e"}, "Quick export summary", "Global"},
			},
		},
		{
			Title: "Tips & Shortcuts",
			Icon:  "💡",
			Content: `Advanced tips and hidden features to boost your productivity:

PRODUCTIVITY TIPS:
• Use number keys (1-4) to quickly switch between analysis modes
• Hold Shift while navigating to select multiple items
• Use Ctrl+C at any time to safely quit the application
• Press 'r' to refresh when file system changes are detected

ANALYSIS BEST PRACTICES:
• Start analysis from the project root for complete dependency graphs
• Use different metrics modes to understand different aspects of code quality
• Pay attention to circular dependencies as they indicate design issues
• Monitor technical debt scores to guide refactoring efforts

PERFORMANCE OPTIMIZATION:
• Exclude test directories for faster analysis when not needed
• Use file patterns to focus analysis on specific file types
• Adjust worker threads based on your system capabilities
• Consider using summary exports for large projects

KEYBOARD MASTERY:
• Learn the Tab navigation for seamless view switching
• Use PgUp/PgDn for efficient content scrolling
• Master the vim-style navigation (hjkl) for speed
• Combine Ctrl with navigation keys for faster movement`,
			KeyBindings: []KeyBinding{
				{[]string{"F1"}, "Show context-sensitive quick help", "Global"},
				{[]string{"Ctrl+?"}, "Show all keyboard shortcuts", "Global"},
				{[]string{"F12"}, "Toggle debug information", "Global"},
				{[]string{"Ctrl+r"}, "Force refresh all views", "Global"},
			},
		},
	}
}

// Update handles navigation within the help system
func (h *HelpViewModel) Update(key string) {
	switch key {
	case "left", "h":
		if h.currentSection > 0 {
			h.currentSection--
			h.scrollY = 0
		}
	case "right", "l":
		if h.currentSection < len(h.sections)-1 {
			h.currentSection++
			h.scrollY = 0
		}
	case "1":
		h.currentSection = 0
		h.scrollY = 0
	case "2":
		if len(h.sections) > 1 {
			h.currentSection = 1
			h.scrollY = 0
		}
	case "3":
		if len(h.sections) > 2 {
			h.currentSection = 2
			h.scrollY = 0
		}
	case "4":
		if len(h.sections) > 3 {
			h.currentSection = 3
			h.scrollY = 0
		}
	case "5":
		if len(h.sections) > 4 {
			h.currentSection = 4
			h.scrollY = 0
		}
	case "6":
		if len(h.sections) > 5 {
			h.currentSection = 5
			h.scrollY = 0
		}
	case "7":
		if len(h.sections) > 6 {
			h.currentSection = 6
			h.scrollY = 0
		}
	case "8":
		if len(h.sections) > 7 {
			h.currentSection = 7
			h.scrollY = 0
		}
	case "up", "k":
		if h.scrollY > 0 {
			h.scrollY--
		}
	case "down", "j":
		if h.scrollY < h.maxScroll {
			h.scrollY++
		}
	case "pgup":
		h.scrollY = max(0, h.scrollY-10)
	case "pgdown":
		h.scrollY = min(h.maxScroll, h.scrollY+10)
	case "home", "g":
		h.scrollY = 0
	case "end", "G":
		h.scrollY = h.maxScroll
	}
}

// Render renders the enhanced help view
func (h *HelpViewModel) Render(width, height int) string {
	h.width = width
	h.height = height

	var b strings.Builder

	// Header with current section indicator
	currentSection := h.sections[h.currentSection]
	header := components.CreateCard(
		fmt.Sprintf("%s %s (%d/%d)", currentSection.Icon, currentSection.Title, h.currentSection+1, len(h.sections)),
		"",
	)
	b.WriteString(header + "\n")

	// Section navigation tabs
	b.WriteString(h.renderSectionTabs() + "\n")

	// Section content
	contentArea := h.renderSectionContent(currentSection)
	b.WriteString(contentArea)

	// Footer with navigation hints
	footer := h.renderFooter()
	b.WriteString("\n" + footer)

	return b.String()
}

// renderSectionTabs renders navigation tabs for help sections
func (h *HelpViewModel) renderSectionTabs() string {
	var tabs []string

	for i, _ := range h.sections {
		var style lipgloss.Style
		if i == h.currentSection {
			style = lipgloss.NewStyle().
				Foreground(components.NeutralWhite).
				Background(components.PrimaryPurple).
				Padding(0, 1).
				Bold(true).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(components.PrimaryBlue)
		} else {
			style = lipgloss.NewStyle().
				Foreground(components.NeutralMedium).
				Background(components.NeutralDark).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(components.NeutralDark)
		}

		tabText := fmt.Sprintf("%d", i+1)
		tabs = append(tabs, style.Render(tabText))
	}

	return strings.Join(tabs, " ")
}

// renderSectionContent renders the content of the current section
func (h *HelpViewModel) renderSectionContent(section HelpSection) string {
	var b strings.Builder

	// Content description
	if section.Content != "" {
		contentCard := components.CardStyle.Render(section.Content)
		b.WriteString(contentCard + "\n")
	}

	// Key bindings
	if len(section.KeyBindings) > 0 {
		b.WriteString(components.SectionStyle.Render("⌨️  Keyboard Shortcuts") + "\n")
		for _, kb := range section.KeyBindings {
			keyText := strings.Join(kb.Keys, ", ")
			b.WriteString(fmt.Sprintf("  %s  %s\n",
				components.CreateBadge(keyText, "info"),
				kb.Description))
			if kb.Context != "" {
				b.WriteString(fmt.Sprintf("    %s\n",
					components.HelpStyle.Render("Context: "+kb.Context)))
			}
		}
		b.WriteString("\n")
	}

	// Tips
	if len(section.Tips) > 0 {
		b.WriteString(components.SectionStyle.Render("💡 Pro Tips") + "\n")
		for _, tip := range section.Tips {
			b.WriteString(components.CreateHighlight(tip, "info") + "\n")
		}
	}

	return h.applyScrolling(b.String())
}

// renderFooter renders the help footer with navigation hints
func (h *HelpViewModel) renderFooter() string {
	navigation := "Navigate: ←→/hl (sections) • ↑↓/kj (scroll) • 1-8 (jump) • ? (close help)"
	return components.HelpStyle.
		Align(lipgloss.Center).
		Render(navigation)
}

// applyScrolling applies scrolling to help content
func (h *HelpViewModel) applyScrolling(content string) string {
	if h.height <= 8 {
		return content
	}

	lines := strings.Split(content, "\n")
	availableHeight := h.height - 8 // Reserve space for header, tabs, footer

	newMaxScroll := max(0, len(lines)-availableHeight)
	if h.maxScroll != newMaxScroll {
		h.maxScroll = newMaxScroll
		if h.scrollY > h.maxScroll {
			h.scrollY = h.maxScroll
		}
	}

	if h.maxScroll == 0 {
		return content
	}

	startLine := h.scrollY
	if startLine >= len(lines) {
		startLine = max(0, len(lines)-availableHeight)
		h.scrollY = startLine
	}

	endLine := min(len(lines), startLine+availableHeight)
	visibleLines := lines[startLine:endLine]
	result := strings.Join(visibleLines, "\n")

	// Add scroll indicator
	if h.maxScroll > 0 {
		scrollInfo := fmt.Sprintf("\n%s",
			components.ScrollInfoStyle.Render(fmt.Sprintf("📊 Showing lines %d-%d of %d (↑↓ to scroll)",
				startLine+1, min(endLine, len(lines)), len(lines))))
		result += scrollInfo
	}

	return result
}
