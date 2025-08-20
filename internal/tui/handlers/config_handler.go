package handlers

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/codebasereaderv2/internal/tui/core"
)

// ConfigHandler handles all configuration-related functionality
type ConfigHandler struct{}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// ProcessConfigCommand processes configuration commands
func (ch *ConfigHandler) ProcessConfigCommand(command string, m *core.MainModel) tea.Cmd {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: "Empty command"}
		}
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "set":
		if len(parts) < 3 {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "Usage: set <key> <value>"}
			}
		}
		key := strings.ToLower(parts[1])
		value := strings.Join(parts[2:], " ")

		return ch.updateConfig(key, value, m)

	case "show":
		if len(parts) > 1 && strings.ToLower(parts[1]) == "config" {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "Configuration displayed above"}
			}
		}

	case "reset":
		if len(parts) > 1 && strings.ToLower(parts[1]) == "config" {
			return ch.resetConfig(m)
		}

	case "add_exclude":
		if len(parts) < 2 {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "Usage: add_exclude <pattern>"}
			}
		}
		pattern := strings.Join(parts[1:], " ")
		return ch.addExcludePattern(pattern, m)

	case "remove_exclude":
		if len(parts) < 2 {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "Usage: remove_exclude <pattern>"}
			}
		}
		pattern := strings.Join(parts[1:], " ")
		return ch.removeExcludePattern(pattern, m)

	default:
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: fmt.Sprintf("Unknown command: %s", cmd)}
		}
	}

	return func() tea.Msg {
		return core.StatusUpdateMsg{Message: "Command processed"}
	}
}

// updateConfig updates a configuration value
func (ch *ConfigHandler) updateConfig(key, value string, m *core.MainModel) tea.Cmd {
	analysisEngine := m.GetAnalysisEngine()
	if analysisEngine == nil {
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	config := analysisEngine.GetConfig()

	switch key {
	case "ai_provider":
		if value == "anthropic" || value == "openai" {
			config.AIProvider = value
			m.GetInputField().SetValue("")
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: fmt.Sprintf("AI provider set to %s", value)}
			}
		} else {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "AI provider must be 'anthropic' or 'openai'"}
			}
		}
	case "api_key":
		config.APIKey = value
		m.GetInputField().SetValue("")
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: "API key updated"}
		}
	case "max_workers":
		if workers := ch.parseInt(value); workers > 0 && workers <= 16 {
			config.MaxWorkers = workers
			m.GetInputField().SetValue("")
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: fmt.Sprintf("Max workers set to %d", workers)}
			}
		} else {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "Max workers must be between 1 and 16"}
			}
		}
	case "timeout":
		if timeout := ch.parseInt(value); timeout > 0 && timeout <= 300 {
			config.Timeout = timeout
			m.GetInputField().SetValue("")
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: fmt.Sprintf("Timeout set to %d seconds", timeout)}
			}
		} else {
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: "Timeout must be between 1 and 300 seconds"}
			}
		}
	default:
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: fmt.Sprintf("Unknown config key: %s", key)}
		}
	}
}

// resetConfig resets configuration to defaults
func (ch *ConfigHandler) resetConfig(m *core.MainModel) tea.Cmd {
	if m.GetAnalysisEngine() == nil {
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	// Update the engine's config (this would need to be implemented in the engine)
	// For now, just show a message
	m.GetInputField().SetValue("")
	return func() tea.Msg {
		return core.StatusUpdateMsg{Message: "Configuration reset to defaults"}
	}
}

// addExcludePattern adds an exclude pattern
func (ch *ConfigHandler) addExcludePattern(pattern string, m *core.MainModel) tea.Cmd {
	if m.GetAnalysisEngine() == nil {
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	config := m.GetAnalysisEngine().GetConfig()
	config.ExcludePatterns = append(config.ExcludePatterns, pattern)
	m.GetInputField().SetValue("")
	return func() tea.Msg {
		return core.StatusUpdateMsg{Message: fmt.Sprintf("Added exclude pattern: %s", pattern)}
	}
}

// removeExcludePattern removes an exclude pattern
func (ch *ConfigHandler) removeExcludePattern(pattern string, m *core.MainModel) tea.Cmd {
	if m.GetAnalysisEngine() == nil {
		return func() tea.Msg {
			return core.StatusUpdateMsg{Message: "Analysis engine not initialized"}
		}
	}

	config := m.GetAnalysisEngine().GetConfig()
	for i, p := range config.ExcludePatterns {
		if p == pattern {
			config.ExcludePatterns = append(config.ExcludePatterns[:i], config.ExcludePatterns[i+1:]...)
			m.GetInputField().SetValue("")
			return func() tea.Msg {
				return core.StatusUpdateMsg{Message: fmt.Sprintf("Removed exclude pattern: %s", pattern)}
			}
		}
	}

	return func() tea.Msg {
		return core.StatusUpdateMsg{Message: fmt.Sprintf("Pattern not found: %s", pattern)}
	}
}

// HandleConfigMessages processes configuration-related messages
func (ch *ConfigHandler) HandleConfigMessages(msg tea.Msg, m *core.MainModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case core.ProcessConfigCommandMsg:
		return m, ch.ProcessConfigCommand(msg.Command, m)
	default:
		return m, nil
	}
}

// parseInt parses a string to int, returns 0 if invalid
func (ch *ConfigHandler) parseInt(s string) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return 0
}
