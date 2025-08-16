package engine

import (
	"time"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

// LanguageStats contains statistics for a specific programming language
type LanguageStats struct {
	FileCount     int `json:"file_count"`
	LineCount     int `json:"line_count"`
	FunctionCount int `json:"function_count"`
	ClassCount    int `json:"class_count"`
	Complexity    int `json:"complexity"`
}

// ProjectAnalysis contains the complete analysis results for a project
type ProjectAnalysis struct {
	RootPath      string                    `json:"root_path"`
	TotalFiles    int                       `json:"total_files"`
	TotalLines    int                       `json:"total_lines"`
	Languages     map[string]LanguageStats  `json:"languages"`
	FileResults   []*parser.AnalysisResult  `json:"file_results"`
	Summary       string                    `json:"summary,omitempty"`
	GeneratedAt   time.Time                 `json:"generated_at"`
	AnalysisDuration time.Duration          `json:"analysis_duration"`
}

// AnalysisJob represents a single file analysis job
type AnalysisJob struct {
	FilePath string
	Content  []byte
	Parser   parser.Parser
}

// AnalysisJobResult contains the result of processing an analysis job
type AnalysisJobResult struct {
	Result *parser.AnalysisResult
	Error  error
}

// Config contains configuration settings for the analysis engine
type Config struct {
	AIProvider      string   `json:"ai_provider"`
	APIKey          string   `json:"api_key"`
	MaxWorkers      int      `json:"max_workers"`
	OutputFormat    string   `json:"output_format"`
	ExcludePatterns []string `json:"exclude_patterns"`
	IncludePatterns []string `json:"include_patterns"`
	MaxFileSize     int64    `json:"max_file_size"` // in bytes
	Timeout         int      `json:"timeout"`       // in seconds
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		AIProvider:      "anthropic",
		MaxWorkers:      4,
		OutputFormat:    "json",
		ExcludePatterns: []string{"node_modules", ".git", "vendor", "__pycache__", ".venv"},
		IncludePatterns: []string{},
		MaxFileSize:     1024 * 1024, // 1MB
		Timeout:         30,          // 30 seconds
	}
}