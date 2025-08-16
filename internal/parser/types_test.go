package parser

import (
	"testing"
	"time"
)

// MockParser implements the Parser interface for testing
type MockParser struct {
	language   string
	extensions []string
}

func (m *MockParser) Parse(filePath string, content []byte) (*AnalysisResult, error) {
	return &AnalysisResult{
		FilePath:   filePath,
		Language:   m.language,
		LineCount:  len(content),
		Functions:  []FunctionInfo{},
		Classes:    []ClassInfo{},
		Imports:    []string{},
		Complexity: 1,
		Errors:     []ParseError{},
		AnalyzedAt: time.Now(),
	}, nil
}

func (m *MockParser) GetSupportedExtensions() []string {
	return m.extensions
}

func (m *MockParser) GetLanguageName() string {
	return m.language
}

func TestParserInterface(t *testing.T) {
	parser := &MockParser{
		language:   "Test",
		extensions: []string{".test"},
	}

	// Test GetLanguageName
	if parser.GetLanguageName() != "Test" {
		t.Errorf("Expected language name 'Test', got '%s'", parser.GetLanguageName())
	}

	// Test GetSupportedExtensions
	extensions := parser.GetSupportedExtensions()
	if len(extensions) != 1 || extensions[0] != ".test" {
		t.Errorf("Expected extensions ['.test'], got %v", extensions)
	}

	// Test Parse
	result, err := parser.Parse("test.test", []byte("test content"))
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}

	if result.FilePath != "test.test" {
		t.Errorf("Expected file path 'test.test', got '%s'", result.FilePath)
	}

	if result.Language != "Test" {
		t.Errorf("Expected language 'Test', got '%s'", result.Language)
	}
}

func TestAnalysisResult(t *testing.T) {
	result := &AnalysisResult{
		FilePath:  "test.go",
		Language:  "Go",
		LineCount: 100,
		Functions: []FunctionInfo{
			{
				Name:       "TestFunction",
				LineStart:  10,
				LineEnd:    20,
				Parameters: []string{"param1", "param2"},
				ReturnType: "string",
				Complexity: 2,
			},
		},
		Classes: []ClassInfo{
			{
				Name:      "TestClass",
				LineStart: 30,
				LineEnd:   50,
				Methods:   []FunctionInfo{},
				Fields:    []string{"field1", "field2"},
			},
		},
		Imports:    []string{"fmt", "testing"},
		Complexity: 5,
		Errors:     []ParseError{},
		AnalyzedAt: time.Now(),
	}

	// Verify all fields are properly set
	if result.FilePath != "test.go" {
		t.Errorf("Expected file path 'test.go', got '%s'", result.FilePath)
	}

	if len(result.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result.Functions))
	}

	if len(result.Classes) != 1 {
		t.Errorf("Expected 1 class, got %d", len(result.Classes))
	}

	if len(result.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(result.Imports))
	}
}