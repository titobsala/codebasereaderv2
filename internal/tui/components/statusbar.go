package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyBind represents a keyboard shortcut
type KeyBind struct {
	Key         string
	Description string
}

// StatusBarModel handles the status bar display
type StatusBarModel struct {
	message  string
	progress float64
	showHelp bool
	keybinds []KeyBind
	width    int
}

// NewStatusBarModel creates a new status bar model
func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		message:  "Ready",
		progress: 0.0,
		showHelp: true,
		keybinds: []KeyBind{
			{"?", "help"},
			{"q", "quit"},
			{"tab", "switch view"},
			{"a", "analyze"},
		},
	}
}

// Update handles messages for the status bar
func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	// Status bar is mostly passive, updated by other components
	return m, nil
}

// View renders the status bar
func (m StatusBarModel) View(width int) string {
	m.width = width

	// Create styled status bar background
	statusBarStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(width).
		Padding(0, 1)

	var content strings.Builder

	// Main status message with progress indicator
	statusMsg := m.message
	if m.progress > 0 && m.progress < 1 {
		// Create a mini progress bar
		progressWidth := 20
		filled := int(float64(progressWidth) * m.progress)
		progressBar := strings.Repeat("█", filled) + strings.Repeat("░", progressWidth-filled)
		statusMsg = fmt.Sprintf("%s [%s] %.0f%%", statusMsg, progressBar, m.progress*100)
	}

	// Key bindings help
	var helpText string
	if m.showHelp && len(m.keybinds) > 0 {
		var bindings []string
		for _, kb := range m.keybinds {
			// Style key bindings
			keyStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC"))

			binding := fmt.Sprintf("%s:%s",
				keyStyle.Render(kb.Key),
				descStyle.Render(kb.Description))
			bindings = append(bindings, binding)
		}
		helpText = strings.Join(bindings, " • ")
	}

	// Calculate available space (accounting for padding)
	availableWidth := width - 2
	statusWidth := lipgloss.Width(statusMsg)
	helpWidth := lipgloss.Width(helpText)

	if statusWidth+helpWidth+3 <= availableWidth {
		// Both fit on one line
		content.WriteString(statusMsg)
		padding := availableWidth - statusWidth - helpWidth - 3
		if padding > 0 {
			content.WriteString(strings.Repeat(" ", padding))
		}
		if helpText != "" {
			content.WriteString(" • " + helpText)
		}
	} else {
		// Status message takes priority
		if statusWidth <= availableWidth {
			content.WriteString(statusMsg)
			remaining := availableWidth - statusWidth - 3
			if remaining > 10 && helpText != "" {
				// Truncate help text to fit
				if helpWidth > remaining {
					// Find a good truncation point
					truncated := helpText
					for lipgloss.Width(truncated) > remaining-3 {
						if len(truncated) <= 10 {
							break
						}
						truncated = truncated[:len(truncated)-10]
					}
					helpText = truncated + "..."
				}
				content.WriteString(" • " + helpText)
			}
		} else {
			// Truncate status message
			truncated := statusMsg
			for lipgloss.Width(truncated) > availableWidth-3 {
				if len(truncated) <= 10 {
					break
				}
				truncated = truncated[:len(truncated)-10]
			}
			content.WriteString(truncated + "...")
		}
	}

	return statusBarStyle.Render(content.String())
}

// SetMessage sets the status message
func (m *StatusBarModel) SetMessage(message string) {
	m.message = message
}

// SetProgress sets the progress value (0.0 to 1.0)
func (m *StatusBarModel) SetProgress(progress float64) {
	m.progress = progress
}

// SetKeyBinds sets the key bindings to display
func (m *StatusBarModel) SetKeyBinds(keybinds []KeyBind) {
	m.keybinds = keybinds
}

// ShowHelp toggles help display
func (m *StatusBarModel) ShowHelp(show bool) {
	m.showHelp = show
}
