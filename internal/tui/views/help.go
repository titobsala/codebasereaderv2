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
			Title: "Getting Started",
			Icon:  "🚀",
			Content: `Welcome to CodebaseReader v2 - A powerful TUI for codebase analysis!

QUICK START:
• Navigate files with ↑↓ arrow keys or j/k
• Press 'a' on any directory to analyze it
• Use Tab to switch between main views
• Press ? to toggle this help system

MAIN INTERFACE:
The interface has 4 main tabs accessible via Tab/Shift+Tab:
• Explorer (1): File tree and directory navigation
• Analysis (2): Metrics, quality scores, and analysis results  
• Configuration (3): Settings and preferences
• Help (4): This help system

SUPPORTED LANGUAGES:
• Go - Full AST analysis with complexity metrics
• Python - Function/class detection and imports
• More languages coming soon!`,
			KeyBindings: []KeyBinding{
				{[]string{"Tab", "Shift+Tab"}, "Switch between main tabs", "Global"},
				{[]string{"1", "2", "3", "4"}, "Jump directly to specific tab", "Global"},
				{[]string{"?"}, "Toggle help on/off", "Global"},
				{[]string{"q", "Ctrl+c"}, "Quit application", "Global"},
			},
		},
		{
			Title: "Navigation & File Operations", 
			Icon:  "🧭",
			Content: `Navigate efficiently through your codebase:

FILE TREE NAVIGATION:
• Use ↑↓ or j/k to move between files and directories
• Press Enter to open files or expand/collapse directories  
• Use ←→ or h/l to collapse/expand directories
• Press Space to toggle directory expansion

FILE OPERATIONS:
• View file contents in the content area
• Files show syntax highlighting and line numbers
• Navigate large files with PgUp/PgDn or scroll arrows
• Press 'a' on directories to analyze them

DIRECTORY ANALYSIS:
• Select any directory and press 'a' to analyze
• Analysis processes all supported files recursively
• Results appear in the Analysis tab automatically
• Use 'c' to clear analysis results when done`,
			KeyBindings: []KeyBinding{
				{[]string{"↑", "↓", "j", "k"}, "Navigate up/down in file tree", "Explorer"},
				{[]string{"←", "→", "h", "l"}, "Collapse/expand directories", "Explorer"}, 
				{[]string{"Enter"}, "Open file or toggle directory", "Explorer"},
				{[]string{"Space"}, "Toggle directory expansion", "Explorer"},
				{[]string{"a"}, "Analyze selected directory", "Explorer"},
				{[]string{"r"}, "Refresh file tree", "Explorer"},
				{[]string{"PgUp", "PgDn"}, "Scroll content by page", "Content View"},
			},
		},
		{
			Title: "Analysis & Metrics",
			Icon:  "📊", 
			Content: `Comprehensive codebase analysis and metrics:

HOW TO ANALYZE:
• Navigate to any directory in the Explorer tab
• Press 'a' to start analysis of that directory
• Analysis runs in background with progress indicator
• Results automatically appear in the Analysis tab

WHAT YOU GET:
• Lines of Code: Total lines, code lines, comments, blank lines
• Complexity Metrics: Cyclomatic complexity for functions
• Function/Class Counts: Detailed breakdown by language  
• Quality Scoring: Overall project quality assessment
• Dependency Analysis: Import relationships and circular dependencies

VIEWING RESULTS:
• Switch to Analysis tab (Tab or press '2') to see results
• Use 'm' to toggle between detailed metrics and overview
• Use 's' to toggle summary view with key insights
• Scroll with ↑↓ to navigate through large result sets

ANALYSIS TIPS:
• Start from project root for complete dependency analysis
• Large projects may take longer - watch the progress indicator
• Use 'c' to clear analysis when switching between projects`,
			KeyBindings: []KeyBinding{
				{[]string{"a"}, "Analyze selected directory", "Explorer"},
				{[]string{"2"}, "Switch to Analysis tab", "Global"},
				{[]string{"m"}, "Toggle metrics/overview view", "Analysis"},
				{[]string{"s"}, "Toggle summary view", "Analysis"}, 
				{[]string{"c"}, "Clear current analysis", "Global"},
				{[]string{"r"}, "Refresh/re-run analysis", "Global"},
			},
		},
		{
			Title: "Configuration & Tips",
			Icon:  "⚙️",
			Content: `Configuration and productivity tips:

CONFIGURATION TAB:
• Switch to Configuration tab (press '3') for settings
• Adjust worker thread count for better performance
• Set file exclusion patterns for faster analysis
• Configure analysis timeout for large projects

PRODUCTIVITY SHORTCUTS:
• Use number keys 1-4 to jump between tabs quickly
• Press 'r' to refresh when files change
• Use Ctrl+C to safely quit at any time
• Navigate with vim-style keys (hjkl) for speed

ANALYSIS BEST PRACTICES:
• Start analysis from project root for complete results
• Watch for circular dependencies - they indicate design issues  
• Use different view modes (m/s) to understand different aspects
• Clear analysis (c) when switching between projects

PERFORMANCE TIPS:
• Exclude test directories if not needed for faster analysis
• Larger projects will take longer - be patient with progress
• The tool handles concurrency automatically
• Results are cached until you clear them`,
			KeyBindings: []KeyBinding{
				{[]string{"3"}, "Switch to Configuration tab", "Global"},
				{[]string{"1", "2", "3", "4"}, "Jump between tabs quickly", "Global"},
				{[]string{"r"}, "Refresh current view", "Global"},
				{[]string{"Ctrl+r", "F5"}, "Force refresh all views", "Global"},
				{[]string{"←", "→", "h", "l"}, "Navigate help sections", "Help"},
				{[]string{"Esc"}, "Exit help and return to Explorer", "Help"},
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
	headerText := fmt.Sprintf("%s %s (%d/%d)", currentSection.Icon, currentSection.Title, h.currentSection+1, len(h.sections))
	header := components.HeaderStyle.Render(headerText)
	b.WriteString(header + "\n\n")

	// Section navigation tabs
	b.WriteString(h.renderSectionTabs() + "\n\n")

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

	for i := range h.sections {
		var style lipgloss.Style
		if i == h.currentSection {
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
		b.WriteString(section.Content + "\n\n")
	}

	// Key bindings
	if len(section.KeyBindings) > 0 {
		b.WriteString("\n⌨️  Keyboard Shortcuts:\n")
		for _, kb := range section.KeyBindings {
			keyText := strings.Join(kb.Keys, ", ")
			b.WriteString(fmt.Sprintf("  %s - %s", keyText, kb.Description))
			if kb.Context != "" {
				b.WriteString(fmt.Sprintf(" (%s)", kb.Context))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Tips
	if len(section.Tips) > 0 {
		b.WriteString("💡 Pro Tips:\n")
		for _, tip := range section.Tips {
			b.WriteString("  • " + tip + "\n")
		}
	}

	return h.applyScrolling(b.String())
}

// renderFooter renders the help footer with navigation hints
func (h *HelpViewModel) renderFooter() string {
	navigation := "Navigate: ←→/hl (sections) • ↑↓/kj (scroll) • 1-4 (jump) • ? (close help)"
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
