package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/tui"
	"github.com/tito-sala/codebasereaderv2/internal/tui/core"
)

// ViewHandler handles view rendering and management
type ViewHandler struct{}

// NewViewHandler creates a new view handler
func NewViewHandler() *ViewHandler {
	return &ViewHandler{}
}

// RenderConfigView renders the configuration view
func (vh *ViewHandler) RenderConfigView(width, height int, m *core.MainModel) string {
	var b strings.Builder

	// Header
	configHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render("⚙️  Configuration")
	b.WriteString(configHeader + "\n\n")

	// Current configuration display
	analysisEngine := m.GetAnalysisEngine()
	if analysisEngine != nil {
		config := analysisEngine.GetConfig()
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
	b.WriteString(inputStyle.Render(m.GetInputField().View()) + "\n\n")

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

// RenderHelpView renders the comprehensive help view
func (vh *ViewHandler) RenderHelpView(width, height int, m *core.MainModel) string {
	helpView := m.GetHelpView()
	return helpView.Render(width, height)
}

// RenderLoadingView renders the loading view
func (vh *ViewHandler) RenderLoadingView(width, height int, m *core.MainModel) string {
	if !m.GetLoading() {
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

	progressInfo := m.GetProgressInfo()
	if progressInfo != nil {
		if progressInfo.Total > 0 {
			percentage := float64(progressInfo.Current) / float64(progressInfo.Total)

			// Progress text
			progressText := fmt.Sprintf("Progress: %d/%d files (%.1f%%)",
				progressInfo.Current, progressInfo.Total, percentage*100)
			b.WriteString(progressText + "\n\n")

			// Update progress model with current percentage
			progress := m.GetProgress()
			progress.Width = tui.Min(60, width-10)

			// Render the bubbles progress bar
			progressBar := progress.ViewAs(percentage)

			barPadding := (width - progress.Width) / 2
			if barPadding > 0 {
				b.WriteString(strings.Repeat(" ", barPadding))
			}
			b.WriteString(progressBar + "\n\n")
		}

		// Current file being processed
		if progressInfo.FilePath != "" {
			currentFileStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Italic(true)

			currentFile := fmt.Sprintf("📄 Current file: %s", progressInfo.FilePath)
			b.WriteString(currentFileStyle.Render(currentFile) + "\n")
		}

		// Status message
		if progressInfo.Message != "" {
			statusStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

			status := fmt.Sprintf("⚡ Status: %s", progressInfo.Message)
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

// RenderConfirmationView renders the confirmation dialog
func (vh *ViewHandler) RenderConfirmationView(width, height int, m *core.MainModel) string {
	confirmationState := m.GetConfirmationState()
	if confirmationState == nil {
		return "No confirmation state"
	}

	var b strings.Builder

	// Center the dialog
	dialogWidth := tui.Min(60, width-4)
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
	dialogContent.WriteString(messageStyle.Render(confirmationState.Message) + "\n\n")

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
