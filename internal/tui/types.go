package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
)

// ViewType represents different views in the TUI
type ViewType int

const (
	FileTreeView ViewType = iota
	ContentView
	ConfigView
	HelpView
	LoadingView
)

// MainModel represents the main TUI model
type MainModel struct {
	fileTree     FileTreeModel
	contentView  ContentViewModel
	statusBar    StatusBarModel
	inputField   textinput.Model
	currentView  ViewType
	analysisData *AnalysisData
	loading      bool
	error        error
	width        int
	height       int
}

// FileTreeModel represents the file tree navigation component
type FileTreeModel struct {
	items       []FileTreeItem
	cursor      int
	selected    map[int]bool
	expanded    map[string]bool
	rootPath    string
	currentPath string
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

// ContentViewModel represents the content display area
type ContentViewModel struct {
	content     string
	scrollY     int
	maxScroll   int
	showMetrics bool
	showSummary bool
}

// StatusBarModel represents the status bar at the bottom
type StatusBarModel struct {
	message    string
	progress   float64
	showHelp   bool
	keybinds   []KeyBind
}

// KeyBind represents a keyboard shortcut
type KeyBind struct {
	Key         string
	Description string
}

// AnalysisData contains the current analysis results
type AnalysisData struct {
	ProjectAnalysis *engine.ProjectAnalysis
	CurrentFile     string
	Summary         string
	Error           error
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