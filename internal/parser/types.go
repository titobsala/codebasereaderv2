package parser

import "time"

// ParseError represents an error encountered during parsing
type ParseError struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}

// Dependency represents a code dependency relationship
type Dependency struct {
	Name         string `json:"name"`
	Type         string `json:"type"` // "import", "internal", "external", "standard"
	Version      string `json:"version,omitempty"`
	UsageCount   int    `json:"usage_count"`
	IsDirectDep  bool   `json:"is_direct_dependency"`
	FilePath     string `json:"file_path"`
}

// FunctionInfo contains information about a function or method
type FunctionInfo struct {
	Name              string   `json:"name"`
	LineStart         int      `json:"line_start"`
	LineEnd           int      `json:"line_end"`
	Parameters        []string `json:"parameters"`
	ReturnType        string   `json:"return_type"`
	Complexity        int      `json:"complexity"`
	CyclomaticComplexity int   `json:"cyclomatic_complexity"`
	LinesOfCode       int      `json:"lines_of_code"`
	ParameterCount    int      `json:"parameter_count"`
	IsPublic          bool     `json:"is_public"`
	IsAsync           bool     `json:"is_async"`
	HasDocstring      bool     `json:"has_docstring"`
}

// ClassInfo contains information about a class or struct
type ClassInfo struct {
	Name              string         `json:"name"`
	LineStart         int            `json:"line_start"`
	LineEnd           int            `json:"line_end"`
	Methods           []FunctionInfo `json:"methods"`
	Fields            []string       `json:"fields"`
	LinesOfCode       int            `json:"lines_of_code"`
	MethodCount       int            `json:"method_count"`
	FieldCount        int            `json:"field_count"`
	IsPublic          bool           `json:"is_public"`
	BaseClasses       []string       `json:"base_classes"`
	HasDocstring      bool           `json:"has_docstring"`
	Complexity        int            `json:"complexity"`
}

// AnalysisResult contains the complete analysis results for a single file
type AnalysisResult struct {
	FilePath             string         `json:"file_path"`
	Language             string         `json:"language"`
	LineCount            int            `json:"line_count"`
	Functions            []FunctionInfo `json:"functions"`
	Classes              []ClassInfo    `json:"classes"`
	Imports              []string       `json:"imports"`
	Complexity           int            `json:"complexity"`
	CyclomaticComplexity int            `json:"cyclomatic_complexity"`
	Errors               []ParseError   `json:"errors"`
	AnalyzedAt           time.Time      `json:"analyzed_at"`
	// Quality metrics
	MaintainabilityIndex float64        `json:"maintainability_index"`
	TechnicalDebt        float64        `json:"technical_debt"`
	CodeDuplication      float64        `json:"code_duplication"`
	TestCoverage         float64        `json:"test_coverage"`
	// Dependency metrics
	Dependencies         []Dependency   `json:"dependencies"`
	ImportCount          int            `json:"import_count"`
	ExportCount          int            `json:"export_count"`
	// Code structure metrics
	CommentLines         int            `json:"comment_lines"`
	BlankLines           int            `json:"blank_lines"`
	CodeLines            int            `json:"code_lines"`
	AverageLineLength    float64        `json:"average_line_length"`
	MaxLineLength        int            `json:"max_line_length"`
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