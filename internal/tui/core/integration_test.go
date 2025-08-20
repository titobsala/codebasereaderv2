package core

import (
	"testing"

	"github.com/tito-sala/codebasereaderv2/internal/engine"
)

func TestAnalysisIntegration(t *testing.T) {
	// Create main model
	model := NewMainModel()

	// Test directory selection message
	msg := DirectorySelectedMsg{Path: "../testdata"}
	updatedModel, _ := model.Update(msg)

	// Check that loading state is set
	mainModel := updatedModel.(*MainModel)
	if !mainModel.loading {
		t.Error("Expected loading state to be true after directory selection")
	}

	// Check that analysis engine is initialized
	if mainModel.analysisEngine == nil {
		t.Error("Expected analysis engine to be initialized")
	}

	// Test analysis started message
	startedMsg := AnalysisStartedMsg{Path: "../testdata"}
	updatedModel, _ = mainModel.Update(startedMsg)

	// Test analysis complete message
	analysis := &engine.ProjectAnalysis{
		RootPath:   "../testdata",
		TotalFiles: 2,
		TotalLines: 50,
		Languages: map[string]engine.LanguageStats{
			"Go": {
				FileCount:     2,
				LineCount:     50,
				FunctionCount: 5,
				ClassCount:    1,
			},
		},
	}

	completeMsg := AnalysisCompleteMsg{
		Analysis: analysis,
		Summary:  "Test summary",
	}

	updatedModel, _ = mainModel.Update(completeMsg)
	mainModel = updatedModel.(*MainModel)

	// Check that loading state is cleared
	if mainModel.loading {
		t.Error("Expected loading state to be false after analysis completion")
	}

	// Check that analysis data is set
	if mainModel.analysisData == nil {
		t.Error("Expected analysis data to be set")
	}

	if mainModel.analysisData.ProjectAnalysis != analysis {
		t.Error("Expected analysis data to match provided analysis")
	}

	// Check that current view is switched to content view
	if mainModel.currentView != ContentView {
		t.Error("Expected current view to be ContentView after analysis completion")
	}
}

func TestProgressReporting(t *testing.T) {
	model := NewMainModel()

	// Test progress message
	progressMsg := AnalysisProgressMsg{
		Current:  5,
		Total:    10,
		FilePath: "test.go",
		Message:  "Processing...",
	}

	updatedModel, _ := model.Update(progressMsg)
	mainModel := updatedModel.(*MainModel)

	// Check that progress info is updated
	if mainModel.progressInfo == nil {
		t.Error("Expected progress info to be set")
	}

	if mainModel.progressInfo.Current != 5 {
		t.Errorf("Expected current progress to be 5, got %d", mainModel.progressInfo.Current)
	}

	if mainModel.progressInfo.Total != 10 {
		t.Errorf("Expected total progress to be 10, got %d", mainModel.progressInfo.Total)
	}

	if mainModel.progressInfo.FilePath != "test.go" {
		t.Errorf("Expected file path to be 'test.go', got %s", mainModel.progressInfo.FilePath)
	}
}

func TestAnalysisErrorHandling(t *testing.T) {
	model := NewMainModel()

	// Test error message
	errorMsg := ErrorMsg{
		Error: &AnalysisError{message: "Test error"},
	}

	updatedModel, _ := model.Update(errorMsg)
	mainModel := updatedModel.(*MainModel)

	// Check that loading state is cleared
	if mainModel.loading {
		t.Error("Expected loading state to be false after error")
	}

	// Check that error is set
	if mainModel.error == nil {
		t.Error("Expected error to be set")
	}
}

// AnalysisError is a simple error type for testing
type AnalysisError struct {
	message string
}

func (e *AnalysisError) Error() string {
	return e.message
}
