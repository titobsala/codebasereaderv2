package metrics

import (
	"testing"
	"time"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

func TestAggregateProjectMetrics(t *testing.T) {
	aggregator := NewAggregator()

	// Create sample analysis results
	results := []*parser.AnalysisResult{
		{
			FilePath:             "main.go",
			Language:             "Go",
			LineCount:            50,
			Functions:            []parser.FunctionInfo{{Name: "main", Complexity: 5, HasDocstring: true}},
			Classes:              []parser.ClassInfo{},
			CyclomaticComplexity: 5,
			MaintainabilityIndex: 85.0,
			TechnicalDebt:        2.5,
			CodeLines:            40,
			CommentLines:         5,
			BlankLines:           5,
			Dependencies: []parser.Dependency{
				{Name: "fmt", Type: "standard"},
				{Name: "os", Type: "standard"},
			},
			AnalyzedAt: time.Now(),
		},
		{
			FilePath:             "utils.go",
			Language:             "Go",
			LineCount:            30,
			Functions:            []parser.FunctionInfo{{Name: "helper", Complexity: 3, HasDocstring: false}},
			Classes:              []parser.ClassInfo{},
			CyclomaticComplexity: 3,
			MaintainabilityIndex: 90.0,
			TechnicalDebt:        1.0,
			CodeLines:            25,
			CommentLines:         2,
			BlankLines:           3,
			Dependencies: []parser.Dependency{
				{Name: "strings", Type: "standard"},
			},
			AnalyzedAt: time.Now(),
		},
		{
			FilePath:             "script.py",
			Language:             "Python",
			LineCount:            40,
			Functions:            []parser.FunctionInfo{{Name: "process", Complexity: 7, HasDocstring: true}},
			Classes:              []parser.ClassInfo{{Name: "DataProcessor", HasDocstring: true}},
			CyclomaticComplexity: 7,
			MaintainabilityIndex: 75.0,
			TechnicalDebt:        4.0,
			CodeLines:            30,
			CommentLines:         5,
			BlankLines:           5,
			Dependencies: []parser.Dependency{
				{Name: "os", Type: "standard"},
				{Name: "requests", Type: "external"},
			},
			AnalyzedAt: time.Now(),
		},
	}

	// Aggregate metrics
	analysis := aggregator.AggregateProjectMetrics(results, "/test/project")

	// Verify project metrics were calculated
	if analysis.ProjectMetrics.TotalComplexity == 0 {
		t.Error("Expected total complexity to be calculated")
	}
	if analysis.ProjectMetrics.AverageComplexity == 0 {
		t.Error("Expected average complexity to be calculated")
	}
	if analysis.ProjectMetrics.MaintainabilityIndex == 0 {
		t.Error("Expected maintainability index to be calculated")
	}
	if analysis.ProjectMetrics.DocumentationRatio == 0 {
		t.Error("Expected documentation ratio to be calculated")
	}

	// Verify directory stats were calculated
	if len(analysis.DirectoryStats) == 0 {
		t.Error("Expected directory stats to be calculated")
	}

	// Verify dependency graph was analyzed
	if len(analysis.DependencyGraph.InternalDependencies) == 0 && len(analysis.DependencyGraph.ExternalDependencies) == 0 {
		// Should have at least some dependencies
		t.Log("Note: No dependencies found in test data")
	}

	// Verify quality score was calculated
	if analysis.QualityScore.Overall == 0 {
		t.Error("Expected overall quality score to be calculated")
	}
	if analysis.QualityScore.Grade == "" {
		t.Error("Expected quality grade to be assigned")
	}

	// Verify specific calculations
	expectedTotalComplexity := 5 + 3 + 7 // Sum of all complexities
	if analysis.ProjectMetrics.TotalComplexity != expectedTotalComplexity {
		t.Errorf("Expected total complexity %d, got %d", expectedTotalComplexity, analysis.ProjectMetrics.TotalComplexity)
	}

	expectedAvgComplexity := float64(expectedTotalComplexity) / float64(len(results))
	if analysis.ProjectMetrics.AverageComplexity != expectedAvgComplexity {
		t.Errorf("Expected average complexity %f, got %f", expectedAvgComplexity, analysis.ProjectMetrics.AverageComplexity)
	}
}

func TestDetectCircularDependencies(t *testing.T) {
	aggregator := NewAggregator()

	// Test case with no circular dependencies
	deps1 := map[string][]string{
		"a.go": {"b.go"},
		"b.go": {"c.go"},
		"c.go": {},
	}

	circular1 := aggregator.detectCircularDependencies(deps1)
	if len(circular1) != 0 {
		t.Errorf("Expected no circular dependencies, got %d", len(circular1))
	}

	// Test case with circular dependencies
	deps2 := map[string][]string{
		"a.go": {"b.go"},
		"b.go": {"c.go"},
		"c.go": {"a.go"},
	}

	circular2 := aggregator.detectCircularDependencies(deps2)
	if len(circular2) == 0 {
		t.Error("Expected circular dependencies to be detected")
	}
}

func TestCalculateDependencyDepth(t *testing.T) {
	aggregator := NewAggregator()

	// Test simple dependency chain
	deps := map[string][]string{
		"a.go": {"b.go"},
		"b.go": {"c.go"},
		"c.go": {"d.go"},
		"d.go": {},
	}

	depth := aggregator.calculateDependencyDepth(deps)
	if depth < 3 {
		t.Errorf("Expected dependency depth >= 3, got %d", depth)
	}

	// Test no dependencies
	emptyDeps := map[string][]string{}
	emptyDepth := aggregator.calculateDependencyDepth(emptyDeps)
	if emptyDepth != 0 {
		t.Errorf("Expected dependency depth 0 for empty deps, got %d", emptyDepth)
	}
}

func TestIsTestFile(t *testing.T) {
	aggregator := NewAggregator()

	testCases := []struct {
		filePath string
		expected bool
	}{
		{"main_test.go", true},
		{"utils_test.go", true},
		{"test_helper.py", true},
		{"helper_test.py", true},
		{"main.test.js", true},
		{"spec.helper.js", true},
		{"main.go", false},
		{"utils.py", false},
		{"helper.js", false},
		{"README.md", false},
	}

	for _, tc := range testCases {
		result := aggregator.isTestFile(tc.filePath)
		if result != tc.expected {
			t.Errorf("isTestFile(%s) = %v, expected %v", tc.filePath, result, tc.expected)
		}
	}
}
