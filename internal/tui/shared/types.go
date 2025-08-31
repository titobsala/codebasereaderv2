package shared

import (
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/metrics"
)

// ViewType represents different views in the TUI
type ViewType int

const (
	FileTreeView ViewType = iota
	ContentView
	HelpView
	ConfigView
)

// FileTreeItem represents a file or directory in the tree
type FileTreeItem struct {
	Name        string
	Path        string
	IsDirectory bool
	IsSupported bool
	Level       int
	Size        int64
	Children    []FileTreeItem
}

// Message types used by views

// DirectorySelectedMsg is sent when a directory is selected
type DirectorySelectedMsg struct {
	Path string
}

// StatusUpdateMsg is sent to update status message
type StatusUpdateMsg struct {
	Message string
}

// RefreshMsg is sent to refresh the view
type RefreshMsg struct{}

// ErrorMsg represents an error message
type ErrorMsg struct {
	Error error
}

// FileSelectedMsg is sent when a file is selected
type FileSelectedMsg struct {
	FilePath string
	Content  string
}

// ShowConfirmationMsg is sent to show a confirmation dialog
type ShowConfirmationMsg struct {
	Message string
	Action  string
	Data    interface{}
}

// DirectoryLoadedMsg is sent when directory loading is complete
type DirectoryLoadedMsg struct {
	Items []FileTreeItem
}

// AnalysisData contains the current analysis results
type AnalysisData struct {
	ProjectAnalysis         *engine.ProjectAnalysis
	EnhancedProjectAnalysis *metrics.EnhancedProjectAnalysis
	Summary                 string
	Timestamp               string
}
