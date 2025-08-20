package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.AIProvider != "anthropic" {
		t.Errorf("Expected AI provider 'anthropic', got '%s'", config.AIProvider)
	}

	if config.MaxWorkers != 4 {
		t.Errorf("Expected max workers 4, got %d", config.MaxWorkers)
	}

	if config.OutputFormat != "json" {
		t.Errorf("Expected output format 'json', got '%s'", config.OutputFormat)
	}

	if config.MaxFileSize != 1024*1024 {
		t.Errorf("Expected max file size 1MB, got %d", config.MaxFileSize)
	}

	if config.Timeout != 30 {
		t.Errorf("Expected timeout 30s, got %d", config.Timeout)
	}

	// Check default exclude patterns
	expectedPatterns := []string{"node_modules", ".git", "vendor", "__pycache__", ".venv"}
	if len(config.ExcludePatterns) != len(expectedPatterns) {
		t.Errorf("Expected %d exclude patterns, got %d", len(expectedPatterns), len(config.ExcludePatterns))
	}
}

func TestNewEngine(t *testing.T) {
	// Test with nil config (should use defaults)
	engine := NewEngine(nil)
	if engine == nil {
		t.Error("Expected engine to be created with nil config")
	}

	config := engine.GetConfig()
	if config.MaxWorkers != 4 {
		t.Errorf("Expected default max workers 4, got %d", config.MaxWorkers)
	}

	// Test with custom config
	customConfig := &Config{
		MaxWorkers: 8,
		AIProvider: "openai",
	}

	engine = NewEngine(customConfig)
	if engine.GetConfig().MaxWorkers != 8 {
		t.Errorf("Expected max workers 8, got %d", engine.GetConfig().MaxWorkers)
	}

	if engine.GetConfig().AIProvider != "openai" {
		t.Errorf("Expected AI provider 'openai', got '%s'", engine.GetConfig().AIProvider)
	}
}

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(2)

	if pool.maxWorkers != 2 {
		t.Errorf("Expected 2 workers, got %d", pool.maxWorkers)
	}

	// Test starting and stopping
	pool.Start()
	if !pool.running {
		t.Error("Expected worker pool to be running after Start()")
	}

	// Give workers a moment to start
	time.Sleep(10 * time.Millisecond)

	pool.Stop()
	if pool.running {
		t.Error("Expected worker pool to be stopped after Stop()")
	}
}

func TestProjectAnalysis(t *testing.T) {
	analysis := &ProjectAnalysis{
		RootPath:   "/test/project",
		TotalFiles: 10,
		TotalLines: 1000,
		Languages: map[string]LanguageStats{
			"Go": {
				FileCount:     5,
				LineCount:     600,
				FunctionCount: 20,
				ClassCount:    0,
				Complexity:    15,
			},
			"Python": {
				FileCount:     5,
				LineCount:     400,
				FunctionCount: 15,
				ClassCount:    5,
				Complexity:    12,
			},
		},
		FileResults:      []*parser.AnalysisResult{},
		Summary:          "Test project analysis",
		GeneratedAt:      time.Now(),
		AnalysisDuration: time.Second * 5,
	}

	if analysis.RootPath != "/test/project" {
		t.Errorf("Expected root path '/test/project', got '%s'", analysis.RootPath)
	}

	if len(analysis.Languages) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(analysis.Languages))
	}

	goStats, exists := analysis.Languages["Go"]
	if !exists {
		t.Error("Expected Go language stats to exist")
	}

	if goStats.FileCount != 5 {
		t.Errorf("Expected Go file count 5, got %d", goStats.FileCount)
	}
}

// MockParser for testing
type MockParser struct {
	name       string
	extensions []string
	parseFunc  func(filePath string, content []byte) (*parser.AnalysisResult, error)
}

func (m *MockParser) Parse(filePath string, content []byte) (*parser.AnalysisResult, error) {
	if m.parseFunc != nil {
		return m.parseFunc(filePath, content)
	}

	return &parser.AnalysisResult{
		FilePath:  filePath,
		Language:  m.name,
		LineCount: len(strings.Split(string(content), "\n")),
		Functions: []parser.FunctionInfo{
			{Name: "testFunc", LineStart: 1, LineEnd: 5},
		},
		Classes:    []parser.ClassInfo{},
		Imports:    []string{"fmt"},
		Complexity: 1,
		Errors:     []parser.ParseError{},
		AnalyzedAt: time.Now(),
	}, nil
}

func (m *MockParser) GetSupportedExtensions() []string {
	return m.extensions
}

func (m *MockParser) GetLanguageName() string {
	return m.name
}

func setupTestProject(t *testing.T) (string, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "engine_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test file structure
	testFiles := map[string]string{
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

func helper() {
	// Helper function
}`,
		"src/utils.go": `package src

func Utils() {
	// Utility function
}`,
		"src/parser.py": `def parse():
    """Parse function"""
    pass

class Parser:
    def __init__(self):
        pass
    
    def run(self):
        pass`,
		"README.md":   "# Test Project\n\nThis is a test project.",
		"config.json": `{"setting": "value"}`,
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestEngine_AnalyzeDirectory(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine with test parsers
	config := DefaultConfig()
	config.MaxWorkers = 2
	engine := NewEngine(config)

	// Register mock parsers
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	pyParser := &MockParser{name: "Python", extensions: []string{"py"}}

	engine.GetParserRegistry().RegisterParser(goParser)
	engine.GetParserRegistry().RegisterParser(pyParser)

	// Analyze directory
	analysis, err := engine.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeDirectory failed: %v", err)
	}

	// Verify results
	if analysis.RootPath != tempDir {
		t.Errorf("Expected root path %s, got %s", tempDir, analysis.RootPath)
	}

	if analysis.TotalFiles != 3 { // main.go, src/utils.go, src/parser.py
		t.Errorf("Expected 3 files, got %d", analysis.TotalFiles)
	}

	if len(analysis.Languages) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(analysis.Languages))
	}

	// Check Go language stats
	goStats, exists := analysis.Languages["Go"]
	if !exists {
		t.Error("Expected Go language stats")
	} else {
		if goStats.FileCount != 2 {
			t.Errorf("Expected 2 Go files, got %d", goStats.FileCount)
		}
		if goStats.FunctionCount != 2 { // 2 files * 1 function each from mock
			t.Errorf("Expected 2 Go functions, got %d", goStats.FunctionCount)
		}
	}

	// Check Python language stats
	pyStats, exists := analysis.Languages["Python"]
	if !exists {
		t.Error("Expected Python language stats")
	} else {
		if pyStats.FileCount != 1 {
			t.Errorf("Expected 1 Python file, got %d", pyStats.FileCount)
		}
	}
}

func TestEngine_AnalyzeDirectoryWithProgress(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine with test parsers
	config := DefaultConfig()
	engine := NewEngine(config)

	// Register mock parsers
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	engine.GetParserRegistry().RegisterParser(goParser)

	// Track progress
	var progressCalls []struct {
		current  int
		total    int
		filePath string
	}

	progressCallback := func(current, total int, filePath string) {
		progressCalls = append(progressCalls, struct {
			current  int
			total    int
			filePath string
		}{current, total, filePath})
	}

	// Analyze directory with progress
	analysis, err := engine.AnalyzeDirectoryWithProgress(tempDir, progressCallback)
	if err != nil {
		t.Fatalf("AnalyzeDirectoryWithProgress failed: %v", err)
	}

	// Verify progress was reported
	if len(progressCalls) == 0 {
		t.Error("Expected progress callbacks")
	}

	// Verify final progress
	if len(progressCalls) > 0 {
		lastCall := progressCalls[len(progressCalls)-1]
		if lastCall.current != analysis.TotalFiles {
			t.Errorf("Expected final progress current=%d, got %d", analysis.TotalFiles, lastCall.current)
		}
		if lastCall.total != analysis.TotalFiles {
			t.Errorf("Expected final progress total=%d, got %d", analysis.TotalFiles, lastCall.total)
		}
	}
}

func TestEngine_AnalyzeFile(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine with test parser
	engine := NewEngine(nil)
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	engine.GetParserRegistry().RegisterParser(goParser)

	// Analyze single file
	filePath := filepath.Join(tempDir, "main.go")
	result, err := engine.AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Verify result
	if result.FilePath != filePath {
		t.Errorf("Expected file path %s, got %s", filePath, result.FilePath)
	}

	if result.Language != "Go" {
		t.Errorf("Expected language Go, got %s", result.Language)
	}

	if len(result.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result.Functions))
	}
}

func TestEngine_AnalyzeFile_UnsupportedFile(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine without registering parsers for .md files
	engine := NewEngine(nil)

	// Try to analyze unsupported file
	filePath := filepath.Join(tempDir, "README.md")
	_, err := engine.AnalyzeFile(filePath)
	if err == nil {
		t.Error("Expected error for unsupported file type")
	}
}

func TestEngine_GetFileWalkerStats(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine with test parsers
	engine := NewEngine(nil)
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	pyParser := &MockParser{name: "Python", extensions: []string{"py"}}

	engine.GetParserRegistry().RegisterParser(goParser)
	engine.GetParserRegistry().RegisterParser(pyParser)

	// Get stats
	stats, err := engine.GetFileWalkerStats(tempDir)
	if err != nil {
		t.Fatalf("GetFileWalkerStats failed: %v", err)
	}

	// Verify stats
	if stats.TotalFiles == 0 {
		t.Error("Expected some total files")
	}

	if stats.SupportedFiles == 0 {
		t.Error("Expected some supported files")
	}

	if stats.FilesByExtension["go"] == 0 {
		t.Error("Expected some .go files")
	}

	if stats.FilesByExtension["py"] == 0 {
		t.Error("Expected some .py files")
	}
}

func TestEngine_ConcurrentAnalysis(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine with multiple workers
	config := DefaultConfig()
	config.MaxWorkers = 4
	engine := NewEngine(config)

	// Register parser that simulates work
	slowParser := &MockParser{
		name:       "Go",
		extensions: []string{"go"},
		parseFunc: func(filePath string, content []byte) (*parser.AnalysisResult, error) {
			// Simulate some processing time
			time.Sleep(10 * time.Millisecond)
			return &parser.AnalysisResult{
				FilePath:   filePath,
				Language:   "Go",
				LineCount:  len(strings.Split(string(content), "\n")),
				Functions:  []parser.FunctionInfo{{Name: "test", LineStart: 1, LineEnd: 2}},
				Classes:    []parser.ClassInfo{},
				Imports:    []string{},
				Complexity: 1,
				Errors:     []parser.ParseError{},
				AnalyzedAt: time.Now(),
			}, nil
		},
	}

	engine.GetParserRegistry().RegisterParser(slowParser)

	// Measure analysis time
	start := time.Now()
	analysis, err := engine.AnalyzeDirectory(tempDir)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Concurrent analysis failed: %v", err)
	}

	// Verify results
	if analysis.TotalFiles != 2 { // Only .go files
		t.Errorf("Expected 2 files, got %d", analysis.TotalFiles)
	}

	// With concurrent processing, it should be faster than sequential
	// (This is a rough check - in practice, the overhead might make it not much faster for small files)
	t.Logf("Analysis took %v for %d files", duration, analysis.TotalFiles)
}

func TestEngine_ErrorHandling(t *testing.T) {
	tempDir, cleanup := setupTestProject(t)
	defer cleanup()

	// Create engine with parser that returns errors
	engine := NewEngine(nil)
	errorParser := &MockParser{
		name:       "Go",
		extensions: []string{"go"},
		parseFunc: func(filePath string, content []byte) (*parser.AnalysisResult, error) {
			return nil, fmt.Errorf("mock parsing error for %s", filePath)
		},
	}

	engine.GetParserRegistry().RegisterParser(errorParser)

	// Analyze directory - should handle errors gracefully
	analysis, err := engine.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("Expected analysis to complete despite parsing errors: %v", err)
	}

	// Should have no successful results due to parsing errors
	if len(analysis.FileResults) != 0 {
		t.Errorf("Expected 0 successful results due to parsing errors, got %d", len(analysis.FileResults))
	}
}
