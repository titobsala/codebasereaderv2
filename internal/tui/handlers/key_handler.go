package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/codebasereaderv2/internal/tui/core"
)

// KeyHandler handles all keyboard input for the TUI
type KeyHandler struct{}

// NewKeyHandler creates a new key handler
func NewKeyHandler() *KeyHandler {
	return &KeyHandler{}
}

// HandleKeyMsg processes keyboard messages and returns appropriate commands
func (kh *KeyHandler) HandleKeyMsg(msg tea.KeyMsg, m *core.MainModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle tab navigation first
	if m.GetCurrentView() != core.ConfirmationView && m.GetCurrentView() != core.LoadingView {
		oldTab := m.GetTabs().GetActiveTab()
		tabs, _ := m.GetTabs().Update(msg)
		m.SetTabs(tabs)

		// If tab changed, update current view
		if m.GetTabs().GetActiveTab() != oldTab {
			m.SetCurrentView(m.GetTabs().MapTabToViewType())
			return m, nil
		}
	}

	// Global key bindings
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "?", "f1":
		// Switch to help tab
		m.GetTabs().SetActiveTab(3) // Help tab
		m.SetCurrentView(core.HelpView)
		return m, nil

	case "esc":
		if m.GetCurrentView() != core.FileTreeView {
			// Switch to Explorer tab (FileTreeView)
			m.GetTabs().SetActiveTab(0)
			m.SetCurrentView(core.FileTreeView)
		}
		return m, nil

	case "f5", "ctrl+r":
		return m, func() tea.Msg {
			return core.RefreshMsg{}
		}

	case "c":
		if m.GetAnalysisData() != nil {
			return m, func() tea.Msg {
				return core.ClearAnalysisMsg{}
			}
		}
	}

	// View-specific key bindings
	cmd = kh.handleViewSpecificKeys(msg, m)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleViewSpecificKeys processes view-specific keyboard shortcuts
func (kh *KeyHandler) handleViewSpecificKeys(msg tea.KeyMsg, m *core.MainModel) tea.Cmd {
	switch m.GetCurrentView() {
	case core.FileTreeView:
		var cmd tea.Cmd
		fileTree, cmd := m.GetFileTree().Update(msg)
		m.SetFileTree(fileTree)
		return cmd

	case core.ConfirmationView:
		return kh.handleConfirmationKeys(msg, m)

	case core.ContentView:
		return kh.handleContentViewKeys(msg, m)

	case core.ConfigView:
		return kh.handleConfigViewKeys(msg, m)

	case core.HelpView:
		// Handle help navigation
		m.GetHelpView().Update(msg.String())
		return nil

	case core.LoadingView:
		return kh.handleLoadingViewKeys(msg)

	default:
		return nil
	}
}

// handleConfirmationKeys processes confirmation dialog input
func (kh *KeyHandler) handleConfirmationKeys(msg tea.KeyMsg, m *core.MainModel) tea.Cmd {
	switch msg.String() {
	case "y", "Y":
		if m.GetConfirmationState() != nil {
			confirmState := m.GetConfirmationState()
			return func() tea.Msg {
				return core.ConfirmationResponseMsg{
					Confirmed: true,
					Action:    confirmState.Action,
					Data:      confirmState.Data,
				}
			}
		}
	case "n", "N", "esc":
		if m.GetConfirmationState() != nil {
			confirmState := m.GetConfirmationState()
			return func() tea.Msg {
				return core.ConfirmationResponseMsg{
					Confirmed: false,
					Action:    confirmState.Action,
					Data:      confirmState.Data,
				}
			}
		}
	}
	return nil
}

// handleContentViewKeys processes content view specific keys
func (kh *KeyHandler) handleContentViewKeys(msg tea.KeyMsg, m *core.MainModel) tea.Cmd {
	switch msg.String() {
	case "m":
		if m.GetAnalysisData() != nil {
			return func() tea.Msg { return core.ToggleMetricsMsg{} }
		}
	case "s":
		if m.GetAnalysisData() != nil {
			return func() tea.Msg { return core.ToggleSummaryMsg{} }
		}
	case "e":
		if m.GetAnalysisData() != nil {
			return func() tea.Msg { return core.ExportMsg{Format: "json", Path: "analysis.json"} }
		}
	default:
		var cmd tea.Cmd
		contentView, cmd := m.GetContentView().Update(msg)
		m.SetContentView(contentView)
		return cmd
	}
	return nil
}

// handleConfigViewKeys processes config view specific keys
func (kh *KeyHandler) handleConfigViewKeys(msg tea.KeyMsg, m *core.MainModel) tea.Cmd {
	switch msg.String() {
	case "enter":
		// Process configuration command
		command := m.GetInputField().Value()
		if command != "" {
			// This will need to be handled by the config handler
			return func() tea.Msg {
				return core.ProcessConfigCommandMsg{Command: command}
			}
		}
	default:
		var cmd tea.Cmd
		inputField, cmd := m.GetInputField().Update(msg)
		m.SetInputField(inputField)
		return cmd
	}
	return nil
}

// handleLoadingViewKeys processes loading view keys
func (kh *KeyHandler) handleLoadingViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		return func() tea.Msg {
			return core.AnalysisCancelledMsg{Reason: "User cancelled"}
		}
	}
	return nil
}
