package core

import (
	"github.com/tito-sala/codebasereaderv2/internal/tui/shared"
)

// ViewType is an alias for shared.ViewType
type ViewType = shared.ViewType

const (
	FileTreeView     = shared.FileTreeView
	ContentView      = shared.ContentView
	HelpView         = shared.HelpView
	ConfigView       = shared.ConfigView
	LoadingView      ViewType = iota + 4 // Start after shared constants
	MetricsView
	QualityView
	DependencyView
	ConfirmationView
)

// ConfirmationState holds state for confirmation dialogs
type ConfirmationState struct {
	Message      string
	Action       string
	Data         interface{}
	PreviousView ViewType
}

// ProgressInfo contains information about ongoing analysis progress
type ProgressInfo struct {
	Current  int
	Total    int
	FilePath string
	Message  string
}

// FileTreeItem is an alias for shared.FileTreeItem
type FileTreeItem = shared.FileTreeItem

// KeyBind represents a keyboard shortcut
type KeyBind struct {
	Key         string
	Description string
}

// AnalysisData is an alias for shared.AnalysisData
type AnalysisData = shared.AnalysisData

// These message types are defined in messages.go

// TUIConfig contains configuration for the TUI
type TUIConfig struct {
	ShowHiddenFiles bool
	ColorScheme     string
	KeyBindings     map[string]string
}

// DefaultTUIConfig returns default TUI configuration
func DefaultTUIConfig() *TUIConfig {
	return &TUIConfig{
		ShowHiddenFiles: false,
		ColorScheme:     "default",
		KeyBindings: map[string]string{
			"quit":    "q",
			"help":    "?",
			"up":      "k",
			"down":    "j",
			"select":  "enter",
			"back":    "esc",
			"toggle":  "space",
			"refresh": "r",
			"analyze": "a",
			"export":  "e",
		},
	}
}
