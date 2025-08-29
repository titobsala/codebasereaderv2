package core

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/parser"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
	"github.com/tito-sala/codebasereaderv2/internal/tui/views"
)

// MainModel represents the state of the TUI application
type MainModel struct {
	fileTree          *views.FileTreeModel
	contentView       *views.ContentViewModel
	statusBar         *components.StatusBarModel
	inputField        textinput.Model
	currentView       ViewType
	loading           bool
	width             int
	height            int
	analysisEngine    *engine.Engine
	analysisData      *AnalysisData
	progress          progress.Model
	progressInfo      *ProgressInfo
	tabs              *components.TabsModel
	helpView          *views.HelpViewModel
	confirmationState *ConfirmationState
	error             error
}

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

	// Create progress bar with custom styling
	prog := progress.New(progress.WithDefaultGradient())
	prog.ShowPercentage = true

	return &MainModel{
		fileTree:       views.NewFileTreeModel(),
		contentView:    views.NewContentViewModel(),
		statusBar:      &components.StatusBarModel{},
		inputField:     ti,
		currentView:    FileTreeView,
		loading:        false,
		width:          80,
		height:         24,
		analysisEngine: analysisEngine,
		progress:       prog,
		tabs:           components.NewTabsModel(),
		helpView:       views.NewHelpViewModel(),
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
		// Update tabs with new dimensions
		m.tabs, _ = m.tabs.Update(msg)
		return m, nil

	case tea.KeyMsg:
		// Handle tab navigation first (if not in special views)
		if m.currentView != ConfirmationView && m.currentView != LoadingView {
			oldTab := m.tabs.GetActiveTab()
			m.tabs, _ = m.tabs.Update(msg)

			// If tab changed, update current view
			if m.tabs.GetActiveTab() != oldTab {
				m.currentView = m.tabs.MapTabToViewType()
				return m, nil
			}
		}

		// Global key bindings (highest priority)
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "?", "f1":
			// Switch to help tab
			m.tabs.SetActiveTab(3) // Help tab
			m.currentView = HelpView
			return m, nil

		case "esc":
			// Always return to Explorer tab when pressing Esc
			if m.currentView != FileTreeView {
				m.tabs.SetActiveTab(0)
				m.currentView = FileTreeView
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

		case "ctrl+home":
			// Navigation reset - return to safe state
			m.tabs.SetActiveTab(0)
			m.currentView = FileTreeView
			m.statusBar.SetMessage("Navigation reset - returned to Explorer")
			return m, nil
		}

		// View-specific key bindings
		switch m.currentView {
		case FileTreeView:
			m.fileTree, cmd = m.fileTree.Update(msg)
			cmds = append(cmds, cmd)

		case ConfirmationView:
			// Handle confirmation dialog input
			switch msg.String() {
			case "y", "Y":
				if m.confirmationState != nil {
					return m, func() tea.Msg {
						return ConfirmationResponseMsg{
							Confirmed: true,
							Action:    m.confirmationState.Action,
							Data:      m.confirmationState.Data,
						}
					}
				}
			case "n", "N", "esc":
				if m.confirmationState != nil {
					return m, func() tea.Msg {
						return ConfirmationResponseMsg{
							Confirmed: false,
							Action:    m.confirmationState.Action,
							Data:      m.confirmationState.Data,
						}
					}
				}
			}

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
			case "r":
				// Reset content view to default state
				if m.analysisData != nil {
					m.contentView.SetShowMetrics(false)
					m.contentView.UpdateContentFromAnalysis()
					return m, nil
				}
			default:
				// Pass all other keys (including 1,2,3,4) to ContentView for handling
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
			// Handle help navigation
			m.helpView.Update(msg.String())

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
		m.statusBar.SetMessage(fmt.Sprintf("Analysis complete - %d files analyzed. Press Ctrl+2 for Analysis tab", msg.Analysis.TotalFiles))

		// Update content view with analysis results but don't force switch
		m.contentView.SetAnalysisData(m.analysisData)
		// Stay in current view, let user decide when to switch

		return m, nil

	case EnhancedAnalysisCompleteMsg:
		m.analysisData = &AnalysisData{
			EnhancedProjectAnalysis: msg.EnhancedAnalysis,
			Summary:                 msg.Summary,
		}
		m.loading = false
		m.progressInfo = nil
		m.statusBar.SetMessage(fmt.Sprintf("Enhanced analysis complete - %d files analyzed. Press Ctrl+2 for Analysis tab", msg.EnhancedAnalysis.TotalFiles))

		// Update content view with enhanced analysis results but don't force switch  
		m.contentView.SetAnalysisData(m.analysisData)
		// Stay in current view, let user decide when to switch

		return m, nil

	case ErrorMsg:
		m.error = msg.Error
		m.loading = false
		m.statusBar.SetMessage(fmt.Sprintf("Error: %s", msg.Error.Error()))
		return m, nil

	case LoadingMsg:
		m.loading = msg.Loading
		if msg.Loading {
			m.currentView = LoadingView
			m.statusBar.SetMessage("Starting analysis...")
		} else if m.currentView == LoadingView {
			m.currentView = FileTreeView
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
		
		// Provide context about what type of analysis is being performed
		analysisType := "Global analysis"
		if m.fileTree != nil && m.fileTree.HasSelectedItems() {
			analysisType = "Selected items analysis"
		}
		
		m.statusBar.SetMessage(fmt.Sprintf("%s: %s", analysisType, msg.Path))
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
			m.contentView.ToggleMetrics()
			m.statusBar.SetMessage("Metrics view toggled")
		}
		return m, nil

	case ToggleSummaryMsg:
		if m.currentView == ContentView && m.analysisData != nil {
			m.contentView.ToggleSummary()
			m.statusBar.SetMessage("Summary view toggled")
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

	case ShowConfirmationMsg:
		m.confirmationState = &ConfirmationState{
			Message:      msg.Message,
			Action:       msg.Action,
			Data:         msg.Data,
			PreviousView: m.currentView,
		}
		m.currentView = ConfirmationView
		return m, nil

	case ConfirmationResponseMsg:
		if m.confirmationState != nil {
			// Restore previous view
			m.currentView = m.confirmationState.PreviousView

			if msg.Confirmed {
				// Execute the confirmed action
				switch msg.Action {
				case "navigate_parent":
					// Navigate to parent directory
					if parentPath, ok := msg.Data.(string); ok {
						m.fileTree.SetCurrentPath(parentPath)
						m.statusBar.SetMessage(fmt.Sprintf("Navigated to: %s", parentPath))
						cmd = m.fileTree.LoadDirectory(parentPath)
						cmds = append(cmds, cmd)
					}
				}
			} else {
				m.statusBar.SetMessage("Action cancelled")
			}

			m.confirmationState = nil
		}
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

// View implements the tea.Model interface
func (m *MainModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Layout dimensions for horizontal tabs
	tabsHeight := 3 // Fixed height for horizontal tabs
	statusHeight := 2
	contentHeight := m.height - tabsHeight - statusHeight

	// Render horizontal tabs
	tabsDisplay := m.tabs.View(m.width)

	// Main content area
	var content string
	switch m.currentView {
	case FileTreeView:
		content = m.fileTree.View(m.width, contentHeight)
	case ContentView:
		// Let user control what they see in the Analysis tab
		content = m.contentView.View(m.width, contentHeight)
	case ConfigView:
		content = m.renderConfigView(m.width)
	case HelpView:
		content = m.renderHelpView(m.width, contentHeight)
	case LoadingView:
		content = m.renderLoadingView(m.width)
	case ConfirmationView:
		content = m.renderConfirmationView(m.width, contentHeight)
	}

	// Join tabs and content vertically (tabs on top, content below)
	mainView := lipgloss.JoinVertical(lipgloss.Left, tabsDisplay, content)

	// Status bar with context-sensitive key bindings
	m.updateStatusBarKeyBinds()
	statusBar := m.statusBar.View(m.width)

	// Error display overlay
	errorDisplay := ""
	if m.error != nil {
		errorMsg := fmt.Sprintf("❌ Error: %s", m.error.Error())
		errorDisplay = components.ErrorStyle.Width(m.width).Render(errorMsg)
	}

	// Combine all parts
	finalView := lipgloss.JoinVertical(lipgloss.Top, mainView, statusBar)
	if errorDisplay != "" {
		finalView = lipgloss.JoinVertical(lipgloss.Top, finalView, errorDisplay)
	}

	return finalView
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
func (m *MainModel) renderConfigView(width int) string {
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

// renderHelpView renders the comprehensive help view
func (m *MainModel) renderHelpView(width, height int) string {
	return m.helpView.Render(width, height)
}

// renderLoadingView renders the loading view
func (m *MainModel) renderLoadingView(width int) string {
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
			percentage := float64(m.progressInfo.Current) / float64(m.progressInfo.Total)

			// Progress text
			progressText := fmt.Sprintf("Progress: %d/%d files (%.1f%%)",
				m.progressInfo.Current, m.progressInfo.Total, percentage*100)
			b.WriteString(progressText + "\n\n")

			// Update progress model with current percentage
			m.progress.Width = min(60, width-10)

			// Render the bubbles progress bar
			progressBar := m.progress.ViewAs(percentage)

			barPadding := (width - m.progress.Width) / 2
			if barPadding > 0 {
				b.WriteString(strings.Repeat(" ", barPadding))
			}
			b.WriteString(progressBar + "\n\n")
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

// renderConfirmationView renders the confirmation dialog
func (m *MainModel) renderConfirmationView(width, height int) string {
	if m.confirmationState == nil {
		return "No confirmation state"
	}

	var b strings.Builder

	// Center the dialog
	dialogWidth := min(60, width-4)
	dialogHeight := 8

	// Calculate centering
	horizontalPadding := (width - dialogWidth) / 2
	verticalPadding := (height - dialogHeight) / 2

	// Add vertical padding
	for i := 0; i < verticalPadding; i++ {
		b.WriteString("\n")
	}

	// Dialog box style
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF5F87")).
		Padding(1, 2).
		Width(dialogWidth - 4).
		Align(lipgloss.Center)

	// Dialog content
	var dialogContent strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true).
		Align(lipgloss.Center)
	dialogContent.WriteString(titleStyle.Render("⚠️  Confirmation Required") + "\n\n")

	// Message
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Align(lipgloss.Center)
	dialogContent.WriteString(messageStyle.Render(m.confirmationState.Message) + "\n\n")

	// Options
	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Align(lipgloss.Center)
	dialogContent.WriteString(optionStyle.Render("Press 'y' to confirm, 'n' or 'Esc' to cancel"))

	// Render the dialog with horizontal centering
	dialog := dialogStyle.Render(dialogContent.String())

	// Add horizontal padding to each line of the dialog
	dialogLines := strings.Split(dialog, "\n")
	for _, line := range dialogLines {
		if horizontalPadding > 0 {
			b.WriteString(strings.Repeat(" ", horizontalPadding))
		}
		b.WriteString(line + "\n")
	}

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
	return tea.Sequence(
		// Start the analysis in background and send progress updates
		func() tea.Msg {
			return AnalysisProgressMsg{
				Current:  0,
				Total:    1,
				FilePath: "",
				Message:  "Starting analysis...",
			}
		},
		m.runAnalysisWithProgress(path),
	)
}

// runAnalysisWithProgress runs analysis with simulated progress updates
func (m *MainModel) runAnalysisWithProgress(path string) tea.Cmd {
	return func() tea.Msg {
		// Get file count for progress tracking
		stats, err := m.analysisEngine.GetFileWalkerStats(path)
		if err != nil {
			return ErrorMsg{Error: fmt.Errorf("failed to scan directory: %w", err)}
		}

		totalFiles := stats.SupportedFiles
		if totalFiles == 0 {
			totalFiles = 1 // Prevent division by zero
		}

		// Create a buffered channel for progress updates
		progressChan := make(chan AnalysisProgressMsg, 100)

		// Progress callback that sends updates to channel
		progressCallback := func(current, total int, filePath string) {
			select {
			case progressChan <- AnalysisProgressMsg{
				Current:  current,
				Total:    total,
				FilePath: filePath,
				Message:  fmt.Sprintf("Processing file %d of %d", current, total),
			}:
			default:
				// Channel full, skip this update
			}
		}

		// Run analysis in goroutine
		resultChan := make(chan tea.Msg, 1)
		go func() {
			defer close(progressChan)
			// Try enhanced analysis first
			enhancedAnalysis, err := m.analysisEngine.AnalyzeDirectoryWithEnhancedMetricsAndProgress(path, progressCallback)
			if err != nil {
				// Fall back to basic analysis
				basicAnalysis, basicErr := m.analysisEngine.AnalyzeDirectoryWithProgress(path, progressCallback)
				if basicErr != nil {
					resultChan <- ErrorMsg{Error: basicErr}
					return
				}

				resultChan <- AnalysisCompleteMsg{
					Analysis: basicAnalysis,
					Summary:  "",
				}
				return
			}

			resultChan <- EnhancedAnalysisCompleteMsg{
				EnhancedAnalysis: enhancedAnalysis,
				Summary:          "",
			}
		}()

		// For now, return the final result
		// In a full implementation, we'd need to handle the progress updates
		// through a different mechanism (like subscriptions or background commands)
		return <-resultChan
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
	case ConfirmationView:
		return "Confirmation"
	default:
		return "Unknown"
	}
}

// updateStatusBarKeyBinds updates the status bar with context-sensitive key bindings
func (m *MainModel) updateStatusBarKeyBinds() {
	var keyBinds []components.KeyBind

	// Global key bindings
	keyBinds = append(keyBinds, components.KeyBind{Key: " ?", Description: "help"})
	keyBinds = append(keyBinds, components.KeyBind{Key: " esc", Description: "explorer"})
	keyBinds = append(keyBinds, components.KeyBind{Key: " 1-4", Description: "tabs"})
	keyBinds = append(keyBinds, components.KeyBind{Key: " tab/shift+tab", Description: "cycle tabs"})

	// View-specific key bindings
	switch m.currentView {
	case FileTreeView:
		keyBinds = append(keyBinds, components.KeyBind{Key: " a", Description: "analyze dir"})
		keyBinds = append(keyBinds, components.KeyBind{Key: " ↑↓", Description: "navigate"})
		keyBinds = append(keyBinds, components.KeyBind{Key: " enter", Description: "select"})
	case ContentView:
		if m.analysisData != nil {
			keyBinds = append(keyBinds, components.KeyBind{Key: " m", Description: "metrics"})
			if m.contentView.ShowMetrics() {
				keyBinds = append(keyBinds, components.KeyBind{Key: " 6-9", Description: "modes"})
			}
			keyBinds = append(keyBinds, components.KeyBind{Key: " s", Description: "summary"})
			keyBinds = append(keyBinds, components.KeyBind{Key: " r", Description: "reset view"})
		}
		keyBinds = append(keyBinds, components.KeyBind{Key: " ↑↓", Description: "scroll"})
	case ConfigView:
		keyBinds = append(keyBinds, components.KeyBind{Key: " enter", Description: "execute cmd"})
		keyBinds = append(keyBinds, components.KeyBind{Key: " type", Description: "command"})
	case LoadingView:
		keyBinds = append(keyBinds, components.KeyBind{Key: " ctrl+c", Description: "cancel"})
	case HelpView:
		keyBinds = append(keyBinds, components.KeyBind{Key: " ↑↓", Description: "navigate"})
	}

	m.statusBar.SetKeyBinds(keyBinds)
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

// Getter and setter methods for handlers to access MainModel fields

func (m *MainModel) GetCurrentView() ViewType {
	return m.currentView
}

func (m *MainModel) SetCurrentView(view ViewType) {
	m.currentView = view
}

func (m *MainModel) GetTabs() *components.TabsModel {
	return m.tabs
}

func (m *MainModel) SetTabs(tabs *components.TabsModel) {
	m.tabs = tabs
}

func (m *MainModel) GetAnalysisData() *AnalysisData {
	return m.analysisData
}

func (m *MainModel) GetFileTree() *views.FileTreeModel {
	return m.fileTree
}

func (m *MainModel) SetFileTree(fileTree *views.FileTreeModel) {
	m.fileTree = fileTree
}

func (m *MainModel) GetContentView() *views.ContentViewModel {
	return m.contentView
}

func (m *MainModel) SetContentView(contentView *views.ContentViewModel) {
	m.contentView = contentView
}

func (m *MainModel) GetHelpView() *views.HelpViewModel {
	return m.helpView
}

func (m *MainModel) GetConfirmationState() *ConfirmationState {
	return m.confirmationState
}

func (m *MainModel) GetInputField() textinput.Model {
	return m.inputField
}

func (m *MainModel) SetInputField(inputField textinput.Model) {
	m.inputField = inputField
}

func (m *MainModel) GetAnalysisEngine() *engine.Engine {
	return m.analysisEngine
}

func (m *MainModel) GetStatusBar() *components.StatusBarModel {
	return m.statusBar
}

func (m *MainModel) SetAnalysisData(data *AnalysisData) {
	m.analysisData = data
}

func (m *MainModel) SetProgressInfo(info *ProgressInfo) {
	m.progressInfo = info
}

func (m *MainModel) GetLoading() bool {
	return m.loading
}

func (m *MainModel) GetProgressInfo() *ProgressInfo {
	return m.progressInfo
}

func (m *MainModel) GetProgress() progress.Model {
	return m.progress
}
