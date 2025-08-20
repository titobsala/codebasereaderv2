package handlers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/codebasereaderv2/internal/tui/core"
)

// AnalysisHandler handles all analysis-related functionality
type AnalysisHandler struct{}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler() *AnalysisHandler {
	return &AnalysisHandler{}
}

// StartAnalysis starts the analysis process for a directory
func (ah *AnalysisHandler) StartAnalysis(path string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return core.AnalysisStartedMsg{Path: path}
		},
		ah.performAnalysis(path),
	)
}

// performAnalysis performs the actual analysis with progress reporting
func (ah *AnalysisHandler) performAnalysis(path string) tea.Cmd {
	return tea.Sequence(
		// Start the analysis in background and send progress updates
		func() tea.Msg {
			return core.AnalysisProgressMsg{
				Current:  0,
				Total:    1,
				FilePath: "",
				Message:  "Starting analysis...",
			}
		},
		ah.runAnalysisWithProgress(path),
	)
}

// runAnalysisWithProgress runs analysis with simulated progress updates
// TODO: This method needs to be refactored to accept the analysis engine as a parameter
func (ah *AnalysisHandler) runAnalysisWithProgress(path string) tea.Cmd {
	return func() tea.Msg {
		// For now, return a placeholder message indicating analysis would run here
		// The full implementation would need access to the analysis engine
		return core.ErrorMsg{Error: fmt.Errorf("analysis handler refactoring incomplete - needs engine access")}
	}
}

// HandleAnalysisMessages processes analysis-related messages
func (ah *AnalysisHandler) HandleAnalysisMessages(msg tea.Msg, m *core.MainModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case core.DirectorySelectedMsg:
		// Start analysis of selected directory
		m.SetLoading(true)
		m.ClearError()
		m.SetProgressInfo(&core.ProgressInfo{
			Current: 0,
			Total:   0,
			Message: "Starting analysis...",
		})
		m.GetStatusBar().SetMessage(fmt.Sprintf("Analyzing directory: %s", msg.Path))
		return m, ah.StartAnalysis(msg.Path)

	case core.AnalysisStartedMsg:
		m.GetStatusBar().SetMessage(fmt.Sprintf("Analysis started for: %s", msg.Path))
		return m, nil

	case core.AnalysisProgressMsg:
		progressInfo := &core.ProgressInfo{
			Current:  msg.Current,
			Total:    msg.Total,
			FilePath: msg.FilePath,
			Message:  msg.Message,
		}
		m.SetProgressInfo(progressInfo)

		progressText := fmt.Sprintf("Analyzing... %d/%d files", msg.Current, msg.Total)
		if msg.FilePath != "" {
			progressText += fmt.Sprintf(" (%s)", msg.FilePath)
		}
		m.GetStatusBar().SetMessage(progressText)
		return m, nil

	case core.AnalysisCompleteMsg:
		analysisData := &core.AnalysisData{
			ProjectAnalysis: msg.Analysis,
			Summary:         msg.Summary,
		}
		m.SetAnalysisData(analysisData)
		m.SetLoading(false)
		m.SetProgressInfo(nil)
		m.GetStatusBar().SetMessage(fmt.Sprintf("Analysis complete - %d files analyzed", msg.Analysis.TotalFiles))

		// Update content view with analysis results
		m.GetContentView().SetAnalysisData(analysisData)
		m.SetCurrentView(core.ContentView)

		return m, nil

	case core.EnhancedAnalysisCompleteMsg:
		analysisData := &core.AnalysisData{
			EnhancedProjectAnalysis: msg.EnhancedAnalysis,
			Summary:                 msg.Summary,
		}
		m.SetAnalysisData(analysisData)
		m.SetLoading(false)
		m.SetProgressInfo(nil)
		m.GetStatusBar().SetMessage(fmt.Sprintf("Enhanced analysis complete - %d files analyzed", msg.EnhancedAnalysis.TotalFiles))

		// Update content view with enhanced analysis results
		m.GetContentView().SetAnalysisData(analysisData)
		m.SetCurrentView(core.ContentView)

		return m, nil

	case core.AnalysisCancelledMsg:
		m.SetLoading(false)
		m.SetProgressInfo(nil)
		m.GetStatusBar().SetMessage(fmt.Sprintf("Analysis cancelled: %s", msg.Reason))
		return m, nil

	case core.ClearAnalysisMsg:
		m.SetAnalysisData(nil)
		m.ClearError()
		m.SetLoading(false)
		m.SetProgressInfo(nil)
		m.GetContentView().SetAnalysisData(nil)
		m.GetStatusBar().SetMessage("Analysis data cleared")
		return m, nil

	case core.ToggleMetricsMsg:
		if m.GetCurrentView() == core.ContentView && m.GetAnalysisData() != nil {
			// These operations need to be handled through methods since the fields are private
			// For now, we'll send a status message indicating the toggle was attempted
			// The actual implementation would require proper getter/setter methods on ContentViewModel
			m.GetStatusBar().SetMessage("Metrics view toggled (implementation needed)")
		}
		return m, nil

	case core.ToggleSummaryMsg:
		if m.GetCurrentView() == core.ContentView && m.GetAnalysisData() != nil {
			// These operations need to be handled through methods since the fields are private
			// For now, we'll send a status message indicating the toggle was attempted
			// The actual implementation would require proper getter/setter methods on ContentViewModel
			m.GetStatusBar().SetMessage("Summary view toggled (implementation needed)")
		}
		return m, nil

	case core.ExportMsg:
		if m.GetAnalysisData() != nil {
			// TODO: Implement export functionality
			m.GetStatusBar().SetMessage(fmt.Sprintf("Exporting analysis to %s (%s format)", msg.Path, msg.Format))
		}
		return m, nil

	default:
		return m, nil
	}
}
