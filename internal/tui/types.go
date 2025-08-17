package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/metrics"
)

// ViewType represents different views in the TUI
type ViewType int

const (
	FileTreeView ViewType = iota
	ContentView
	ConfigView
	HelpView
	LoadingView
	MetricsView
	QualityView
	DependencyView
	ConfirmationView
)

// MainModel represents the main TUI model
type MainModel struct {
	fileTree          *FileTreeModel
	contentView       *ContentViewModel
	statusBar         StatusBarModel
	inputField        textinput.Model
	currentView       ViewType
	analysisData      *AnalysisData
	loading           bool
	error             error
	width             int
	height            int
	analysisEngine    *engine.Engine
	progressInfo      *ProgressInfo
	confirmationState *ConfirmationState
}

// ConfirmationState holds state for confirmation dialogs
type ConfirmationState struct {
	Message    string
	Action     string
	Data       interface{}
	PreviousView ViewType
}

// ProgressInfo contains information about ongoing analysis progress
type ProgressInfo struct {
	Current  int
	Total    int
	FilePath string
	Message  string
}

// FileTreeItem represents an item in the file tree
type FileTreeItem struct {
	Name        string
	Path        string
	IsDirectory bool
	IsSupported bool
	Level       int
	Size        int64
	Children    []FileTreeItem
}

// KeyBind represents a keyboard shortcut
type KeyBind struct {
	Key         string
	Description string
}

// AnalysisData contains the current analysis results
type AnalysisData struct {
	ProjectAnalysis         *engine.ProjectAnalysis
	EnhancedProjectAnalysis *metrics.EnhancedProjectAnalysis
	CurrentFile             string
	Summary                 string
	Error                   error
}

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
			"quit":         "q",
			"help":         "?",
			"up":           "k",
			"down":         "j",
			"select":       "enter",
			"back":         "esc",
			"toggle":       "space",
			"refresh":      "r",
			"analyze":      "a",
			"export":       "e",
		},
	}
}