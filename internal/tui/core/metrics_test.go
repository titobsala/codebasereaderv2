package core

import (
	"testing"

	"github.com/tito-sala/codebasereaderv2/internal/metrics"
)

func TestMetricsDisplay(t *testing.T) {
	display := NewMetricsDisplay()

	// Test initial state
	if display.mode != OverviewMode {
		t.Errorf("Expected initial mode to be OverviewMode, got %v", display.mode)
	}

	// Test mode switching
	display.SetMode(DetailedMode)
	if display.mode != DetailedMode {
		t.Errorf("Expected mode to be DetailedMode, got %v", display.mode)
	}

	// Set up dimensions and max scroll for testing
	display.height = 20
	display.maxScroll = 10

	// Test scrolling
	display.Scroll(5)
	if display.scrollY != 5 {
		t.Errorf("Expected scrollY to be 5, got %d", display.scrollY)
	}

	// Test scroll bounds
	display.Scroll(-10)
	if display.scrollY != 0 {
		t.Errorf("Expected scrollY to be 0 (bounded), got %d", display.scrollY)
	}
}

func TestMetricsDisplayRender(t *testing.T) {
	display := NewMetricsDisplay()

	// Test with nil analysis
	result := display.Render(nil, 80, 24)
	if result == "" {
		t.Error("Expected non-empty result for nil analysis")
	}

	// Test with sample analysis
	analysis := &metrics.EnhancedProjectAnalysis{
		RootPath:   "/test/project",
		TotalFiles: 10,
		TotalLines: 1000,
		Languages: map[string]metrics.LanguageStats{
			"Go": {
				FileCount:            5,
				LineCount:            600,
				FunctionCount:        20,
				ClassCount:           5,
				Complexity:           50,
				AverageComplexity:    2.5,
				MaintainabilityIndex: 85.0,
			},
			"Python": {
				FileCount:            5,
				LineCount:            400,
				FunctionCount:        15,
				ClassCount:           3,
				Complexity:           30,
				AverageComplexity:    2.0,
				MaintainabilityIndex: 90.0,
			},
		},
		ProjectMetrics: metrics.ProjectMetrics{
			TotalComplexity:      80,
			AverageComplexity:    2.3,
			MaxComplexity:        10,
			MaintainabilityIndex: 87.5,
			TechnicalDebt:        15.5,
			DocumentationRatio:   75.0,
		},
		QualityScore: metrics.QualityScore{
			Overall:         85.0,
			Maintainability: 87.5,
			Complexity:      75.0,
			Documentation:   75.0,
			TestCoverage:    60.0,
			CodeDuplication: 5.0,
			Grade:           "B",
		},
		DirectoryStats: map[string]metrics.DirectoryStats{
			"/test/project/src": {
				Path:                 "/test/project/src",
				FileCount:            8,
				LineCount:            800,
				Complexity:           60,
				MaintainabilityIndex: 85.0,
			},
		},
		DependencyGraph: metrics.DependencyGraph{
			InternalDependencies: map[string][]string{
				"main.go": {"utils.go", "config.go"},
			},
			ExternalDependencies: map[string][]string{
				"main.go": {"fmt", "os"},
			},
			CircularDependencies: [][]string{},
			DependencyDepth:      3,
			UnusedDependencies:   []string{},
		},
	}

	// Test overview mode
	display.SetMode(OverviewMode)
	result = display.Render(analysis, 80, 24)
	if result == "" {
		t.Error("Expected non-empty result for overview mode")
	}

	// Test detailed mode
	display.SetMode(DetailedMode)
	result = display.Render(analysis, 80, 24)
	if result == "" {
		t.Error("Expected non-empty result for detailed mode")
	}

	// Test quality mode
	display.SetMode(QualityMode)
	result = display.Render(analysis, 80, 24)
	if result == "" {
		t.Error("Expected non-empty result for quality mode")
	}

	// Test dependency mode
	display.SetMode(DependencyMode)
	result = display.Render(analysis, 80, 24)
	if result == "" {
		t.Error("Expected non-empty result for dependency mode")
	}
}

func TestContentViewModelWithEnhancedMetrics(t *testing.T) {
	contentView := NewContentViewModel()

	// Test initial state
	if contentView.metricsDisplay == nil {
		t.Error("Expected metricsDisplay to be initialized")
	}

	// Test setting analysis data
	analysisData := &AnalysisData{
		EnhancedProjectAnalysis: &metrics.EnhancedProjectAnalysis{
			RootPath:   "/test",
			TotalFiles: 5,
			TotalLines: 500,
			Languages: map[string]metrics.LanguageStats{
				"Go": {
					FileCount:     5,
					LineCount:     500,
					FunctionCount: 10,
				},
			},
			QualityScore: metrics.QualityScore{
				Overall: 85.0,
				Grade:   "B",
			},
		},
	}

	contentView.SetAnalysisData(analysisData)
	if contentView.analysisData != analysisData {
		t.Error("Expected analysis data to be set")
	}

	// Test metrics toggle
	contentView.showMetrics = true
	contentView.updateContentFromAnalysis()

	// Content should be updated with metrics
	if contentView.content == "" {
		t.Error("Expected content to be updated with metrics")
	}
}

func TestEnhancedAnalysisIntegration(t *testing.T) {
	model := NewMainModel()

	// Test enhanced analysis complete message
	enhancedAnalysis := &metrics.EnhancedProjectAnalysis{
		RootPath:   "/test/project",
		TotalFiles: 10,
		TotalLines: 1000,
		Languages: map[string]metrics.LanguageStats{
			"Go": {
				FileCount:     10,
				LineCount:     1000,
				FunctionCount: 50,
			},
		},
		QualityScore: metrics.QualityScore{
			Overall: 90.0,
			Grade:   "A",
		},
	}

	msg := EnhancedAnalysisCompleteMsg{
		EnhancedAnalysis: enhancedAnalysis,
		Summary:          "Test summary",
	}

	updatedModel, _ := model.Update(msg)
	mainModel := updatedModel.(*MainModel)

	// Check that analysis data was set
	if mainModel.analysisData == nil {
		t.Error("Expected analysis data to be set")
	}

	if mainModel.analysisData.EnhancedProjectAnalysis != enhancedAnalysis {
		t.Error("Expected enhanced analysis to be set correctly")
	}

	// Check that loading state was cleared
	if mainModel.loading {
		t.Error("Expected loading to be false after analysis complete")
	}

	// Check that view switched to content
	if mainModel.currentView != ContentView {
		t.Errorf("Expected current view to be ContentView, got %v", mainModel.currentView)
	}
}
