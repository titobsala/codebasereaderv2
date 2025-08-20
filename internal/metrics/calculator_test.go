package metrics

import (
	"testing"
	"time"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

func TestCalculateFileMetrics(t *testing.T) {
	calculator := NewCalculator()

	// Test Go file metrics
	goContent := []byte(`package main

import (
	"fmt"
	"os"
)

// PublicFunction is a public function
func PublicFunction(param1 string, param2 int) error {
	if param1 == "" {
		return fmt.Errorf("param1 cannot be empty")
	}
	
	for i := 0; i < param2; i++ {
		fmt.Println(param1)
	}
	
	return nil
}

func privateFunction() {
	// This is a private function
	fmt.Println("private")
}

type PublicStruct struct {
	Field1 string
	Field2 int
}
`)

	result := &parser.AnalysisResult{
		FilePath:     "test.go",
		Language:     "Go",
		LineCount:    25,
		Functions:    []parser.FunctionInfo{},
		Classes:      []parser.ClassInfo{},
		Imports:      []string{"fmt", "os"},
		Dependencies: []parser.Dependency{},
		AnalyzedAt:   time.Now(),
	}

	calculator.CalculateFileMetrics(result, goContent)

	// Verify basic line metrics were calculated
	if result.CodeLines == 0 {
		t.Error("Expected code lines to be calculated")
	}
	if result.CommentLines == 0 {
		t.Error("Expected comment lines to be calculated")
	}
	if result.BlankLines == 0 {
		t.Error("Expected blank lines to be calculated")
	}
	if result.AverageLineLength == 0 {
		t.Error("Expected average line length to be calculated")
	}
	if result.MaxLineLength == 0 {
		t.Error("Expected max line length to be calculated")
	}

	// Verify maintainability index was calculated
	if result.MaintainabilityIndex == 0 {
		t.Error("Expected maintainability index to be calculated")
	}

	// Verify technical debt was calculated
	if result.TechnicalDebt < 0 {
		t.Error("Expected technical debt to be non-negative")
	}

	// Verify dependencies were analyzed
	if result.ImportCount != 2 {
		t.Errorf("Expected import count to be 2, got %d", result.ImportCount)
	}
	if len(result.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(result.Dependencies))
	}
}

func TestCalculateQualityScore(t *testing.T) {
	calculator := NewCalculator()

	// Test quality score calculation
	score, grade := calculator.CalculateQualityScore(85.0, 15.0, 75.0, 80.0, 5.0)

	if score < 0 || score > 100 {
		t.Errorf("Expected score to be between 0 and 100, got %f", score)
	}

	validGrades := map[string]bool{"A": true, "B": true, "C": true, "D": true, "F": true}
	if !validGrades[grade] {
		t.Errorf("Expected valid grade (A-F), got %s", grade)
	}

	// Test high quality score
	highScore, highGrade := calculator.CalculateQualityScore(95.0, 5.0, 90.0, 95.0, 2.0)
	if highGrade != "A" {
		t.Errorf("Expected grade A for high quality, got %s", highGrade)
	}
	if highScore < 90 {
		t.Errorf("Expected high score (>90) for high quality, got %f", highScore)
	}

	// Test low quality score
	lowScore, lowGrade := calculator.CalculateQualityScore(30.0, 50.0, 20.0, 10.0, 40.0)
	if lowGrade == "A" {
		t.Errorf("Expected low grade for low quality, got %s", lowGrade)
	}
	if lowScore > 60 {
		t.Errorf("Expected low score (<60) for low quality, got %f", lowScore)
	}
}

func TestClassifyDependency(t *testing.T) {
	calculator := NewCalculator()

	tests := []struct {
		importPath string
		language   string
		expected   string
	}{
		{"fmt", "go", "standard"},
		{"os", "go", "standard"},
		{"github.com/user/repo", "go", "external"},
		{"golang.org/x/tools", "go", "external"},
		{"myproject/internal/pkg", "go", "internal"},
		{"os", "python", "standard"},
		{"sys", "python", "standard"},
		{"requests", "python", "external"},
		{"mymodule.submodule", "python", "internal"},
		{"numpy", "python", "external"},
	}

	for _, test := range tests {
		result := calculator.classifyDependency(test.importPath, test.language)
		if result != test.expected {
			t.Errorf("classifyDependency(%s, %s) = %s, expected %s",
				test.importPath, test.language, result, test.expected)
		}
	}
}
