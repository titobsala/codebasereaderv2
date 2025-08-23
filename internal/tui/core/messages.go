package core

import (
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/metrics"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
	"github.com/tito-sala/codebasereaderv2/internal/tui/shared"
)

// Message types for the TUI

// FileSelectedMsg is an alias for shared.FileSelectedMsg
type FileSelectedMsg = shared.FileSelectedMsg

// DirectorySelectedMsg is an alias for shared.DirectorySelectedMsg
type DirectorySelectedMsg = shared.DirectorySelectedMsg

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

// ErrorMsg is an alias for shared.ErrorMsg
type ErrorMsg = shared.ErrorMsg

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

// DirectoryLoadedMsg is an alias for shared.DirectoryLoadedMsg
type DirectoryLoadedMsg = shared.DirectoryLoadedMsg

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

// RefreshMsg is an alias for shared.RefreshMsg  
type RefreshMsg = shared.RefreshMsg

// StatusUpdateMsg is an alias for shared.StatusUpdateMsg
type StatusUpdateMsg = shared.StatusUpdateMsg

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
	Mode components.MetricsDisplayMode
}

// MetricsFilterMsg is sent to filter metrics by language
type MetricsFilterMsg struct {
	Language string
}

// MetricsSortMsg is sent to sort metrics
type MetricsSortMsg struct {
	SortBy string
}

// ShowConfirmationMsg is an alias for shared.ShowConfirmationMsg
type ShowConfirmationMsg = shared.ShowConfirmationMsg

// ConfirmationResponseMsg handles confirmation responses
type ConfirmationResponseMsg struct {
	Confirmed bool
	Action    string
	Data      interface{}
}

// ProcessConfigCommandMsg is sent to process configuration commands
type ProcessConfigCommandMsg struct {
	Command string
}
