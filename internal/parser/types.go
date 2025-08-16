package parser

import "time"

// ParseError represents an error encountered during parsing
type ParseError struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}

// FunctionInfo contains information about a function or method
type FunctionInfo struct {
	Name       string   `json:"name"`
	LineStart  int      `json:"line_start"`
	LineEnd    int      `json:"line_end"`
	Parameters []string `json:"parameters"`
	ReturnType string   `json:"return_type"`
	Complexity int      `json:"complexity"`
}

// ClassInfo contains information about a class or struct
type ClassInfo struct {
	Name      string         `json:"name"`
	LineStart int            `json:"line_start"`
	LineEnd   int            `json:"line_end"`
	Methods   []FunctionInfo `json:"methods"`
	Fields    []string       `json:"fields"`
}

// AnalysisResult contains the complete analysis results for a single file
type AnalysisResult struct {
	FilePath     string         `json:"file_path"`
	Language     string         `json:"language"`
	LineCount    int            `json:"line_count"`
	Functions    []FunctionInfo `json:"functions"`
	Classes      []ClassInfo    `json:"classes"`
	Imports      []string       `json:"imports"`
	Complexity   int            `json:"complexity"`
	Errors       []ParseError   `json:"errors"`
	AnalyzedAt   time.Time      `json:"analyzed_at"`
}

// Parser defines the interface that all language parsers must implement
type Parser interface {
	// Parse analyzes file content and returns structured results
	Parse(filePath string, content []byte) (*AnalysisResult, error)
	
	// GetSupportedExtensions returns file extensions this parser handles
	GetSupportedExtensions() []string
	
	// GetLanguageName returns the human-readable language name
	GetLanguageName() string
}