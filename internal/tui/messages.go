package tui

import (
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/metrics"
)

// Message types for the TUI

// FileSelectedMsg is sent when a file is selected
type FileSelectedMsg struct {
	FilePath string
	Content  string
}

// DirectorySelectedMsg is sent when a directory is selected for analysis
type DirectorySelectedMsg struct {
	Path string
}

// AnalysisCompleteMsg is sent when analysis is complete
type AnalysisCompleteMsg struct {
	Analysis *engine.ProjectAnalysis
	Summary  string
}

// EnhancedAnalysisCompleteMsg is sent when enhanced analysis is complete
type EnhancedAnalysisCompleteMsg struct {
	EnhancedAnalysis *metrics.EnhancedProjectAnalysis
	Summary          string
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Error error
}

// LoadingMsg is sent to update loading state
type LoadingMsg struct {
	Loading bool
}

// AnalysisStartedMsg is sent when analysis starts
type AnalysisStartedMsg struct {
	Path string
}

// AnalysisProgressMsg is sent to update analysis progress
type AnalysisProgressMsg struct {
	Current  int
	Total    int
	FilePath string
	Message  string
}

// DirectoryLoadedMsg is already defined in filetree.go

// ViewSwitchMsg is sent to switch views
type ViewSwitchMsg struct {
	NewView ViewType
}

// ShowHelpMsg is sent to show/hide help
type ShowHelpMsg struct {
	Show bool
}

// ToggleMetricsMsg is sent to toggle metrics view
type ToggleMetricsMsg struct{}

// ToggleSummaryMsg is sent to toggle summary view
type ToggleSummaryMsg struct{}

// ClearAnalysisMsg is sent to clear analysis data
type ClearAnalysisMsg struct{}

// RefreshMsg is sent to refresh the file tree
type RefreshMsg struct{}

// StatusUpdateMsg is sent to update status bar
type StatusUpdateMsg struct {
	Message string
}

// ProgressUpdateMsg is sent to update progress
type ProgressUpdateMsg struct {
	Progress float64
	Message  string
}

// FileContentLoadedMsg is sent when file content is loaded
type FileContentLoadedMsg struct {
	FilePath string
	Content  string
	Size     int64
}

// AnalysisCancelledMsg is sent when analysis is cancelled
type AnalysisCancelledMsg struct {
	Reason string
}

// ExportMsg is sent to export analysis results
type ExportMsg struct {
	Format string
	Path   string
}

// MetricsModeChangeMsg is sent to change metrics display mode
type MetricsModeChangeMsg struct {
	Mode MetricsDisplayMode
}

// MetricsFilterMsg is sent to filter metrics by language
type MetricsFilterMsg struct {
	Language string
}

// MetricsSortMsg is sent to sort metrics
type MetricsSortMsg struct {
	SortBy string
}

// ShowConfirmationMsg shows a confirmation dialog
type ShowConfirmationMsg struct {
	Message string
	Action  string
	Data    interface{}
}

// ConfirmationResponseMsg handles confirmation responses
type ConfirmationResponseMsg struct {
	Confirmed bool
	Action    string
	Data      interface{}
}