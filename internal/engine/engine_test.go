package engine

import (
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