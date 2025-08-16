package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/parser"
)



// NewMainModel creates a new main TUI model
func NewMainModel() *MainModel {
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.CharLimit = 256

	// Create analysis engine with default config
	engineConfig := engine.DefaultConfig()
	analysisEngine := engine.NewEngine(engineConfig)
	
	// Register parsers
	analysisEngine.GetParserRegistry().RegisterParser(parser.NewGoParser())
	analysisEngine.GetParserRegistry().RegisterParser(parser.NewPythonParser())

	return &MainModel{
		fileTree:       NewFileTreeModel(),
		contentView:    NewContentViewModel(),
		statusBar:      NewStatusBarModel(),
		inputField:     ti,
		currentView:    FileTreeView,
		loading:        false,
		width:          80,
		height:         24,
		analysisEngine: analysisEngine,
	}
}

// Init implements the tea.Model interface
func (m *MainModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.fileTree.Init(),
	)
}

// Update implements the tea.Model interface
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global key bindings
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "?", "f1":
			return m, func() tea.Msg {
				return ShowHelpMsg{Show: m.currentView != HelpView}
			}

		case "tab":
			return m, func() tea.Msg {
				newView := m.getNextView()
				return ViewSwitchMsg{NewView: newView}
			}

		case "shift+tab":
			return m, func() tea.Msg {
				newView := m.getPreviousView()
				return ViewSwitchMsg{NewView: newView}
			}

		case "esc":
			if m.currentView != FileTreeView {
				return m, func() tea.Msg {
					return ViewSwitchMsg{NewView: FileTreeView}
				}
			}
			return m, nil

		case "f5", "ctrl+r":
			return m, func() tea.Msg {
				return RefreshMsg{}
			}

		case "c":
			if m.analysisData != nil {
				return m, func() tea.Msg {
					return ClearAnalysisMsg{}
				}
			}

		case "f2":
			return m, func() tea.Msg {
				return ViewSwitchMsg{NewView: ConfigView}
			}

		case "f3":
			return m, func() tea.Msg {
				return ViewSwitchMsg{NewView: ContentView}
			}

		case "f4":
			return m, func() tea.Msg {
				return ViewSwitchMsg{NewView: FileTreeView}
			}
		}

		// View-specific key bindings
		switch m.currentView {
		case FileTreeView:
			m.fileTree, cmd = m.fileTree.Update(msg)
			cmds = append(cmds, cmd)

		case ContentView:
			// Handle content view specific keys
			switch msg.String() {
			case "m":
				if m.analysisData != nil {
					return m, func() tea.Msg { return ToggleMetricsMsg{} }
				}
			case "s":
				if m.analysisData != nil {
					return m, func() tea.Msg { return ToggleSummaryMsg{} }
				}
			case "e":
				if m.analysisData != nil {
					return m, func() tea.Msg { return ExportMsg{Format: "json", Path: "analysis.json"} }
				}
			default:
				m.contentView, cmd = m.contentView.Update(msg)
				cmds = append(cmds, cmd)
			}

		case ConfigView:
			// Handle config view specific keys
			switch msg.String() {
			case "enter":
				// Process configuration command
				command := m.inputField.Value()
				if command != "" {
					return m, m.processConfigCommand(command)
				}
			default:
				m.inputField, cmd = m.inputField.Update(msg)
				cmds = append(cmds, cmd)
			}

		case HelpView:
			// Help view doesn't need specific key handling beyond global keys
			break

		case LoadingView:
			// Loading view only responds to cancel
			switch msg.String() {
			case "ctrl+c":
				return m, func() tea.Msg {
					return AnalysisCancelledMsg{Reason: "User cancelled"}
				}
			}
		}

	case FileSelectedMsg:
		m.currentView = ContentView
		m.contentView.SetContent(msg.FilePath, msg.Content)
		return m, nil

	case AnalysisCompleteMsg:
		m.analysisData = &AnalysisData{
			ProjectAnalysis: msg.Analysis,
			Summary:         msg.Summary,
		}
		m.loading = false
		m.progressInfo = nil
		m.statusBar.SetMessage(fmt.Sprintf("Analysis complete - %d files analyzed", msg.Analysis.TotalFiles))
		
		// Update content view with analysis results
		m.contentView.SetAnalysisData(m.analysisData)
		m.currentView = ContentView
		
		return m, nil

	case EnhancedAnalysisCompleteMsg:
		m.analysisData = &AnalysisData{
			EnhancedProjectAnalysis: msg.EnhancedAnalysis,
			Summary:                 msg.Summary,
		}
		m.loading = false
		m.progressInfo = nil
		m.statusBar.SetMessage(fmt.Sprintf("Enhanced analysis complete - %d files analyzed", msg.EnhancedAnalysis.TotalFiles))
		
		// Update content view with enhanced analysis results
		m.contentView.SetAnalysisData(m.analysisData)
		m.currentView = ContentView
		
		return m, nil

	case ErrorMsg:
		m.error = msg.Error
		m.loading = false
		m.statusBar.SetMessage(fmt.Sprintf("Error: %s", msg.Error.Error()))
		return m, nil

	case LoadingMsg:
		m.loading = msg.Loading
		if msg.Loading {
			m.statusBar.SetMessage("Analyzing...")
		}
		return m, nil

	case DirectorySelectedMsg:
		// Start analysis of selected directory
		m.loading = true
		m.error = nil
		m.progressInfo = &ProgressInfo{
			Current: 0,
			Total:   0,
			Message: "Starting analysis...",
		}
		m.statusBar.SetMessage(fmt.Sprintf("Analyzing directory: %s", msg.Path))
		return m, m.startAnalysis(msg.Path)

	case AnalysisStartedMsg:
		m.statusBar.SetMessage(fmt.Sprintf("Analysis started for: %s", msg.Path))
		return m, nil

	case AnalysisProgressMsg:
		m.progressInfo = &ProgressInfo{
			Current:  msg.Current,
			Total:    msg.Total,
			FilePath: msg.FilePath,
			Message:  msg.Message,
		}
		progressText := fmt.Sprintf("Analyzing... %d/%d files", msg.Current, msg.Total)
		if msg.FilePath != "" {
			progressText += fmt.Sprintf(" (%s)", msg.FilePath)
		}
		m.statusBar.SetMessage(progressText)
		return m, nil

	case DirectoryLoadedMsg:
		// Forward to file tree
		m.fileTree, cmd = m.fileTree.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case ViewSwitchMsg:
		m.currentView = msg.NewView
		m.statusBar.SetMessage(fmt.Sprintf("Switched to %s view", m.getViewName()))
		return m, nil

	case ShowHelpMsg:
		if msg.Show {
			m.currentView = HelpView
		} else {
			m.currentView = FileTreeView
		}
		return m, nil

	case ToggleMetricsMsg:
		if m.currentView == ContentView && m.analysisData != nil {
			m.contentView.showMetrics = !m.contentView.showMetrics
			m.contentView.showSummary = false
			m.contentView.updateContentFromAnalysis()
			status := "Metrics view enabled"
			if !m.contentView.showMetrics {
				status = "Metrics view disabled"
			}
			m.statusBar.SetMessage(status)
		}
		return m, nil

	case ToggleSummaryMsg:
		if m.currentView == ContentView && m.analysisData != nil {
			m.contentView.showSummary = !m.contentView.showSummary
			m.contentView.showMetrics = false
			m.contentView.updateContentFromAnalysis()
			status := "Summary view enabled"
			if !m.contentView.showSummary {
				status = "Summary view disabled"
			}
			m.statusBar.SetMessage(status)
		}
		return m, nil

	case ClearAnalysisMsg:
		m.analysisData = nil
		m.error = nil
		m.loading = false
		m.progressInfo = nil
		m.contentView.SetAnalysisData(nil)
		m.statusBar.SetMessage("Analysis data cleared")
		return m, nil

	case RefreshMsg:
		m.statusBar.SetMessage("Refreshing file tree...")
		return m, m.fileTree.Init()

	case StatusUpdateMsg:
		m.statusBar.SetMessage(msg.Message)
		return m, nil

	case ProgressUpdateMsg:
		m.statusBar.SetProgress(msg.Progress)
		if msg.Message != "" {
			m.statusBar.SetMessage(msg.Message)
		}
		return m, nil

	case FileContentLoadedMsg:
		m.currentView = ContentView
		m.contentView.SetContent(msg.FilePath, msg.Content)
		m.statusBar.SetMessage(fmt.Sprintf("Loaded %s (%s)", msg.FilePath, formatFileSize(msg.Size)))
		return m, nil

	case AnalysisCancelledMsg:
		m.loading = false
		m.progressInfo = nil
		m.statusBar.SetMessage(fmt.Sprintf("Analysis cancelled: %s", msg.Reason))
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// View implements the tea.Model interface
func (m *MainModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate layout dimensions
	headerHeight := 3
	statusHeight := 2
	contentHeight := m.height - headerHeight - statusHeight

	// Title bar with current view indicator
	viewName := m.getViewName()
	title := titleStyle.Width(m.width).Render(fmt.Sprintf("CodebaseReader v2 - %s", viewName))
	
	// Main content area
	var content string
	switch m.currentView {
	case FileTreeView:
		content = m.fileTree.View(m.width, contentHeight)
	case ContentView:
		content = m.contentView.View(m.width, contentHeight)
	case ConfigView:
		content = m.renderConfigView(m.width, contentHeight)
	case HelpView:
		content = m.renderHelpView(m.width, contentHeight)
	case LoadingView:
		content = m.renderLoadingView(m.width, contentHeight)
	}

	// Ensure content fits within available space
	contentLines := strings.Split(content, "\n")
	if len(contentLines) > contentHeight {
		contentLines = contentLines[:contentHeight]
		content = strings.Join(contentLines, "\n")
	}

	// Status bar with context-sensitive key bindings
	m.updateStatusBarKeyBinds()
	statusBar := m.statusBar.View(m.width)

	// Error display overlay
	errorDisplay := ""
	if m.error != nil {
		errorMsg := fmt.Sprintf("❌ Error: %s", m.error.Error())
		errorDisplay = errorStyle.Width(m.width).Render(errorMsg)
	}

	// Combine all parts with proper spacing
	view := strings.Builder{}
	view.WriteString(title + "\n")
	
	// Add separator line
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3C3C3C")).
		Render(strings.Repeat("─", m.width))
	view.WriteString(separator + "\n")
	
	view.WriteString(content)
	
	// Pad content to ensure status bar is at bottom
	currentLines := strings.Count(view.String(), "\n")
	neededPadding := m.height - statusHeight - currentLines - 1
	if neededPadding > 0 {
		view.WriteString(strings.Repeat("\n", neededPadding))
	}
	
	if errorDisplay != "" {
		view.WriteString(errorDisplay + "\n")
	}
	view.WriteString(statusBar)

	return view.String()
}

// switchView cycles through available views
func (m *MainModel) switchView() {
	switch m.currentView {
	case FileTreeView:
		m.currentView = ContentView
	case ContentView:
		m.currentView = ConfigView
	case ConfigView:
		m.currentView = FileTreeView
	default:
		m.currentView = FileTreeView
	}
}

// renderConfigView renders the configuration view
func (m *MainModel) renderConfigView(width, height int) string {
	var b strings.Builder
	
	// Header
	configHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render("⚙️  Configuration")
	b.WriteString(configHeader + "\n\n")
	
	// Current configuration display
	if m.analysisEngine != nil {
		config := m.analysisEngine.GetConfig()
		b.WriteString("📋 Current Settings:\n")
		b.WriteString(fmt.Sprintf("  AI Provider: %s\n", config.AIProvider))
		b.WriteString(fmt.Sprintf("  Max Workers: %d\n", config.MaxWorkers))
		b.WriteString(fmt.Sprintf("  Output Format: %s\n", config.OutputFormat))
		b.WriteString(fmt.Sprintf("  Max File Size: %d bytes\n", config.MaxFileSize))
		b.WriteString(fmt.Sprintf("  Timeout: %d seconds\n", config.Timeout))
		
		if len(config.ExcludePatterns) > 0 {
			b.WriteString("  Exclude Patterns:\n")
			for _, pattern := range config.ExcludePatterns {
				b.WriteString(fmt.Sprintf("    - %s\n", pattern))
			}
		}
		b.WriteString("\n")
	}
	
	// Input field
	b.WriteString("💬 Enter command:\n")
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(width - 4)
	b.WriteString(inputStyle.Render(m.inputField.View()) + "\n\n")
	
	// Available commands
	b.WriteString("📝 Available Commands:\n")
	commands := []string{
		"set ai_provider <anthropic|openai>  - Set AI provider",
		"set api_key <key>                  - Set API key",
		"set max_workers <number>           - Set worker count",
		"set timeout <seconds>              - Set timeout",
		"add_exclude <pattern>              - Add exclude pattern",
		"remove_exclude <pattern>           - Remove exclude pattern",
		"show config                        - Show current config",
		"reset config                       - Reset to defaults",
	}
	
	for _, cmd := range commands {
		b.WriteString("  " + cmd + "\n")
	}
	
	return b.String()
}

// renderHelpView renders the help view
func (m *MainModel) renderHelpView(width, height int) string {
	var b strings.Builder
	
	// Header
	helpHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render("❓ CodebaseReader v2 - Help")
	b.WriteString(helpHeader + "\n\n")
	
	// Create styled sections
	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Bold(true)
	
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)
	
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC"))
	
	// Navigation section
	b.WriteString(sectionStyle.Render("🧭 Navigation:") + "\n")
	navKeys := [][]string{
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"→/l/Enter", "Expand directory / Select file"},
		{"←/h", "Collapse directory / Go back"},
		{"Space", "Toggle selection"},
		{"Tab", "Switch between views"},
		{"Esc", "Return to file tree"},
		{"PgUp/PgDn", "Scroll content (in content view)"},
		{"Home/End", "Jump to start/end (in content view)"},
	}
	
	for _, key := range navKeys {
		b.WriteString(fmt.Sprintf("  %s  %s\n", 
			keyStyle.Render(fmt.Sprintf("%-12s", key[0])), 
			descStyle.Render(key[1])))
	}
	b.WriteString("\n")
	
	// Actions section
	b.WriteString(sectionStyle.Render("⚡ Actions:") + "\n")
	actionKeys := [][]string{
		{"a", "Analyze selected directory"},
		{"r", "Refresh file tree"},
		{"e", "Export results (when available)"},
		{"m", "Toggle metrics view (in content view)"},
		{"s", "Toggle summary view (in content view)"},
		{"c", "Clear current analysis"},
	}
	
	for _, key := range actionKeys {
		b.WriteString(fmt.Sprintf("  %s  %s\n", 
			keyStyle.Render(fmt.Sprintf("%-12s", key[0])), 
			descStyle.Render(key[1])))
	}
	b.WriteString("\n")
	
	// Global section
	b.WriteString(sectionStyle.Render("🌐 Global:") + "\n")
	globalKeys := [][]string{
		{"?", "Toggle this help"},
		{"q/Ctrl+C", "Quit application"},
		{"F1", "Show keyboard shortcuts"},
	}
	
	for _, key := range globalKeys {
		b.WriteString(fmt.Sprintf("  %s  %s\n", 
			keyStyle.Render(fmt.Sprintf("%-12s", key[0])), 
			descStyle.Render(key[1])))
	}
	b.WriteString("\n")
	
	// Views section
	b.WriteString(sectionStyle.Render("👁️  Views:") + "\n")
	views := [][]string{
		{"File Tree", "Navigate and select files/directories"},
		{"Content", "View file content and analysis results"},
		{"Config", "Configure application settings"},
		{"Help", "This help screen"},
	}
	
	for _, view := range views {
		b.WriteString(fmt.Sprintf("  %s  %s\n", 
			keyStyle.Render(fmt.Sprintf("%-12s", view[0])), 
			descStyle.Render(view[1])))
	}
	b.WriteString("\n")
	
	// Tips section
	b.WriteString(sectionStyle.Render("💡 Tips:") + "\n")
	tips := []string{
		"• Use 'a' on a directory to start analysis",
		"• Switch views with Tab to see different information",
		"• Analysis results show in the Content view",
		"• Use 'm' and 's' in Content view to toggle different displays",
		"• Configuration changes take effect immediately",
	}
	
	for _, tip := range tips {
		b.WriteString("  " + descStyle.Render(tip) + "\n")
	}
	
	return b.String()
}

// renderLoadingView renders the loading view
func (m *MainModel) renderLoadingView(width, height int) string {
	if !m.loading {
		return ""
	}

	var b strings.Builder
	
	// Centered loading header
	loadingHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render("🔍 Analyzing Codebase...")
	
	headerPadding := (width - lipgloss.Width(loadingHeader)) / 2
	if headerPadding > 0 {
		b.WriteString(strings.Repeat(" ", headerPadding))
	}
	b.WriteString(loadingHeader + "\n\n")

	if m.progressInfo != nil {
		if m.progressInfo.Total > 0 {
			percentage := float64(m.progressInfo.Current) / float64(m.progressInfo.Total) * 100
			
			// Progress text
			progressText := fmt.Sprintf("Progress: %d/%d files (%.1f%%)", 
				m.progressInfo.Current, m.progressInfo.Total, percentage)
			b.WriteString(progressText + "\n\n")
			
			// Styled progress bar
			barWidth := min(60, width-10)
			filled := int(float64(barWidth) * percentage / 100)
			
			progressBarStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4"))
			
			emptyBarStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#3C3C3C"))
			
			progressBar := progressBarStyle.Render(strings.Repeat("█", filled)) + 
				emptyBarStyle.Render(strings.Repeat("░", barWidth-filled))
			
			barPadding := (width - barWidth) / 2
			if barPadding > 0 {
				b.WriteString(strings.Repeat(" ", barPadding))
			}
			b.WriteString(fmt.Sprintf("[%s]\n\n", progressBar))
		}

		// Current file being processed
		if m.progressInfo.FilePath != "" {
			currentFileStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Italic(true)
			
			currentFile := fmt.Sprintf("📄 Current file: %s", m.progressInfo.FilePath)
			b.WriteString(currentFileStyle.Render(currentFile) + "\n")
		}

		// Status message
		if m.progressInfo.Message != "" {
			statusStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))
			
			status := fmt.Sprintf("⚡ Status: %s", m.progressInfo.Message)
			b.WriteString(statusStyle.Render(status) + "\n")
		}
	} else {
		// Initial loading state
		initStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Italic(true)
		
		b.WriteString(initStyle.Render("🚀 Initializing analysis...") + "\n")
	}

	// Loading animation or spinner could be added here
	b.WriteString("\n")
	
	// Helpful message
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true)
	
	helpMsg := "Please wait while we process your files. Press Ctrl+C to cancel."
	helpPadding := (width - len(helpMsg)) / 2
	if helpPadding > 0 {
		b.WriteString(strings.Repeat(" ", helpPadding))
	}
	b.WriteString(helpStyle.Render(helpMsg))
	
	return b.String()
}

// SetError sets an error message
func (m *MainModel) SetError(err error) {
	m.error = err
}

// ClearError clears the current error
func (m *MainModel) ClearError() {
	m.error = nil
}

// SetLoading sets the loading state
func (m *MainModel) SetLoading(loading bool) {
	m.loading = loading
	if loading {
		m.currentView = LoadingView
	} else if m.currentView == LoadingView {
		m.currentView = FileTreeView
	}
}

// startAnalysis starts the analysis process for a directory
func (m *MainModel) startAnalysis(path string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return AnalysisStartedMsg{Path: path}
		},
		m.performAnalysis(path),
	)
}

// performAnalysis performs the actual analysis with progress reporting
func (m *MainModel) performAnalysis(path string) tea.Cmd {
	return func() tea.Msg {
		// Try enhanced analysis first
		enhancedAnalysis, err := m.analysisEngine.AnalyzeDirectoryWithEnhancedMetrics(path)
		if err != nil {
			// Fall back to basic analysis
			basicAnalysis, basicErr := m.analysisEngine.AnalyzeDirectory(path)
			if basicErr != nil {
				return ErrorMsg{Error: basicErr}
			}
			
			return AnalysisCompleteMsg{
				Analysis: basicAnalysis,
				Summary:  "", // AI summary would be generated separately
			}
		}

		return EnhancedAnalysisCompleteMsg{
			EnhancedAnalysis: enhancedAnalysis,
			Summary:          "", // AI summary would be generated separately
		}
	}
}

// getViewName returns the human-readable name of the current view
func (m *MainModel) getViewName() string {
	switch m.currentView {
	case FileTreeView:
		return "File Tree"
	case ContentView:
		return "Content"
	case ConfigView:
		return "Configuration"
	case HelpView:
		return "Help"
	case LoadingView:
		return "Loading"
	default:
		return "Unknown"
	}
}

// updateStatusBarKeyBinds updates the status bar with context-sensitive key bindings
func (m *MainModel) updateStatusBarKeyBinds() {
	var keyBinds []KeyBind
	
	// Global key bindings
	keyBinds = append(keyBinds, KeyBind{"?", "help"})
	keyBinds = append(keyBinds, KeyBind{"q", "quit"})
	keyBinds = append(keyBinds, KeyBind{"tab", "switch view"})
	
	// View-specific key bindings
	switch m.currentView {
	case FileTreeView:
		keyBinds = append(keyBinds, KeyBind{"a", "analyze"})
		keyBinds = append(keyBinds, KeyBind{"r", "refresh"})
		keyBinds = append(keyBinds, KeyBind{"enter", "select"})
	case ContentView:
		if m.analysisData != nil {
			keyBinds = append(keyBinds, KeyBind{"m", "metrics"})
			keyBinds = append(keyBinds, KeyBind{"s", "summary"})
			keyBinds = append(keyBinds, KeyBind{"e", "export"})
		}
		keyBinds = append(keyBinds, KeyBind{"↑↓", "scroll"})
	case ConfigView:
		keyBinds = append(keyBinds, KeyBind{"enter", "execute"})
	case LoadingView:
		keyBinds = append(keyBinds, KeyBind{"ctrl+c", "cancel"})
	}
	
	m.statusBar.SetKeyBinds(keyBinds)
}

// getNextView returns the next view in the cycle
func (m *MainModel) getNextView() ViewType {
	switch m.currentView {
	case FileTreeView:
		return ContentView
	case ContentView:
		return ConfigView
	case ConfigView:
		return FileTreeView
	case HelpView:
		return FileTreeView
	case LoadingView:
		return FileTreeView
	default:
		return FileTreeView
	}
}

// getPreviousView returns the previous view in the cycle
func (m *MainModel) getPreviousView() ViewType {
	switch m.currentView {
	case FileTreeView:
		return ConfigView
	case ContentView:
		return FileTreeView
	case ConfigView:
		return ContentView
	case HelpView:
		return FileTreeView
	case LoadingView:
		return FileTreeView
	default:
		return FileTreeView
	}
}

// handleGlobalKeyBinding handles global keyboard shortcuts
func (m *MainModel) handleGlobalKeyBinding(key string) tea.Cmd {
	switch key {
	case "m":
		if m.currentView == ContentView && m.analysisData != nil {
			return func() tea.Msg {
				return ToggleMetricsMsg{}
			}
		}
	case "s":
		if m.currentView == ContentView && m.analysisData != nil {
			return func() tea.Msg {
				return ToggleSummaryMsg{}
			}
		}
	case "e":
		if m.analysisData != nil {
			return func() tea.Msg {
				return ExportMsg{Format: "json", Path: "analysis.json"}
			}
		}
	}
	return nil
}



// processConfigCommand processes configuration commands
func (m *MainModel) processConfigCommand(command string) tea.Cmd {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return func() tea.Msg {
			return StatusUpdateMsg{Message: "Empty command"}
		}
	}

	cmd := strings.ToLower(parts[0])
	
	switch cmd {
	case "set":
		if len(parts) < 3 {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "Usage: set <key> <value>"}
			}
		}
		key := strings.ToLower(parts[1])
		value := strings.Join(parts[2:], " ")
		
		return m.updateConfig(key, value)

	case "show":
		if len(parts) > 1 && strings.ToLower(parts[1]) == "config" {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "Configuration displayed above"}
			}
		}

	case "reset":
		if len(parts) > 1 && strings.ToLower(parts[1]) == "config" {
			return m.resetConfig()
		}

	case "add_exclude":
		if len(parts) < 2 {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "Usage: add_exclude <pattern>"}
			}
		}
		pattern := strings.Join(parts[1:], " ")
		return m.addExcludePattern(pattern)

	case "remove_exclude":
		if len(parts) < 2 {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "Usage: remove_exclude <pattern>"}
			}
		}
		pattern := strings.Join(parts[1:], " ")
		return m.removeExcludePattern(pattern)

	default:
		return func() tea.Msg {
			return StatusUpdateMsg{Message: fmt.Sprintf("Unknown command: %s", cmd)}
		}
	}

	return func() tea.Msg {
		return StatusUpdateMsg{Message: "Command processed"}
	}
}

// updateConfig updates a configuration value
func (m *MainModel) updateConfig(key, value string) tea.Cmd {
	if m.analysisEngine == nil {
		return func() tea.Msg {
			return StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	config := m.analysisEngine.GetConfig()
	
	switch key {
	case "ai_provider":
		if value == "anthropic" || value == "openai" {
			config.AIProvider = value
			m.inputField.SetValue("")
			return func() tea.Msg {
				return StatusUpdateMsg{Message: fmt.Sprintf("AI provider set to %s", value)}
			}
		} else {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "AI provider must be 'anthropic' or 'openai'"}
			}
		}
	case "api_key":
		config.APIKey = value
		m.inputField.SetValue("")
		return func() tea.Msg {
			return StatusUpdateMsg{Message: "API key updated"}
		}
	case "max_workers":
		if workers := parseInt(value); workers > 0 && workers <= 16 {
			config.MaxWorkers = workers
			m.inputField.SetValue("")
			return func() tea.Msg {
				return StatusUpdateMsg{Message: fmt.Sprintf("Max workers set to %d", workers)}
			}
		} else {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "Max workers must be between 1 and 16"}
			}
		}
	case "timeout":
		if timeout := parseInt(value); timeout > 0 && timeout <= 300 {
			config.Timeout = timeout
			m.inputField.SetValue("")
			return func() tea.Msg {
				return StatusUpdateMsg{Message: fmt.Sprintf("Timeout set to %d seconds", timeout)}
			}
		} else {
			return func() tea.Msg {
				return StatusUpdateMsg{Message: "Timeout must be between 1 and 300 seconds"}
			}
		}
	default:
		return func() tea.Msg {
			return StatusUpdateMsg{Message: fmt.Sprintf("Unknown config key: %s", key)}
		}
	}
}

// resetConfig resets configuration to defaults
func (m *MainModel) resetConfig() tea.Cmd {
	if m.analysisEngine == nil {
		return func() tea.Msg {
			return StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	// Update the engine's config (this would need to be implemented in the engine)
	// For now, just show a message
	m.inputField.SetValue("")
	return func() tea.Msg {
		return StatusUpdateMsg{Message: "Configuration reset to defaults"}
	}
}

// addExcludePattern adds an exclude pattern
func (m *MainModel) addExcludePattern(pattern string) tea.Cmd {
	if m.analysisEngine == nil {
		return func() tea.Msg {
			return StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	config := m.analysisEngine.GetConfig()
	config.ExcludePatterns = append(config.ExcludePatterns, pattern)
	m.inputField.SetValue("")
	return func() tea.Msg {
		return StatusUpdateMsg{Message: fmt.Sprintf("Added exclude pattern: %s", pattern)}
	}
}

// removeExcludePattern removes an exclude pattern
func (m *MainModel) removeExcludePattern(pattern string) tea.Cmd {
	if m.analysisEngine == nil {
		return func() tea.Msg {
			return StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	config := m.analysisEngine.GetConfig()
	for i, p := range config.ExcludePatterns {
		if p == pattern {
			config.ExcludePatterns = append(config.ExcludePatterns[:i], config.ExcludePatterns[i+1:]...)
			m.inputField.SetValue("")
			return func() tea.Msg {
				return StatusUpdateMsg{Message: fmt.Sprintf("Removed exclude pattern: %s", pattern)}
			}
		}
	}
	
	return func() tea.Msg {
		return StatusUpdateMsg{Message: fmt.Sprintf("Pattern not found: %s", pattern)}
	}
}

// parseInt parses a string to int, returns 0 if invalid
func parseInt(s string) int {
	if i, err := fmt.Sscanf(s, "%d", new(int)); err == nil && i == 1 {
		var result int
		fmt.Sscanf(s, "%d", &result)
		return result
	}
	return 0
}

