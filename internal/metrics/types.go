package metrics

import "time"

// ProjectMetrics contains overall project metrics
type ProjectMetrics struct {
	TotalComplexity      int     `json:"total_complexity"`
	AverageComplexity    float64 `json:"average_complexity"`
	MaxComplexity        int     `json:"max_complexity"`
	MaintainabilityIndex float64 `json:"maintainability_index"`
	TechnicalDebt        float64 `json:"technical_debt"`
	CodeDuplication      float64 `json:"code_duplication"`
	TestCoverage         float64 `json:"test_coverage"`
	DocumentationRatio   float64 `json:"documentation_ratio"`
	CodeToCommentRatio   float64 `json:"code_to_comment_ratio"`
}

// DirectoryStats contains statistics for a specific directory
type DirectoryStats struct {
	Path                 string                   `json:"path"`
	FileCount            int                      `json:"file_count"`
	LineCount            int                      `json:"line_count"`
	Languages            map[string]LanguageStats `json:"languages"`
	Complexity           int                      `json:"complexity"`
	MaintainabilityIndex float64                  `json:"maintainability_index"`
	SubDirectories       []string                 `json:"sub_directories"`
}

// LanguageStats contains statistics for a specific programming language
type LanguageStats struct {
	FileCount            int     `json:"file_count"`
	LineCount            int     `json:"line_count"`
	FunctionCount        int     `json:"function_count"`
	ClassCount           int     `json:"class_count"`
	Complexity           int     `json:"complexity"`
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	AverageComplexity    float64 `json:"average_complexity"`
	MaxComplexity        int     `json:"max_complexity"`
	MaintainabilityIndex float64 `json:"maintainability_index"`
	TechnicalDebt        float64 `json:"technical_debt"`
	CodeLines            int     `json:"code_lines"`
	CommentLines         int     `json:"comment_lines"`
	BlankLines           int     `json:"blank_lines"`
	TestFiles            int     `json:"test_files"`
	TestCoverage         float64 `json:"test_coverage"`
}

// DependencyGraph represents the project's dependency relationships
type DependencyGraph struct {
	InternalDependencies map[string][]string `json:"internal_dependencies"`
	ExternalDependencies map[string][]string `json:"external_dependencies"`
	StandardDependencies map[string][]string `json:"standard_dependencies"`
	CircularDependencies [][]string          `json:"circular_dependencies"`
	DependencyDepth      int                 `json:"dependency_depth"`
	UnusedDependencies   []string            `json:"unused_dependencies"`
}

// QualityScore represents overall code quality metrics
type QualityScore struct {
	Overall         float64 `json:"overall"`
	Maintainability float64 `json:"maintainability"`
	Complexity      float64 `json:"complexity"`
	Documentation   float64 `json:"documentation"`
	TestCoverage    float64 `json:"test_coverage"`
	CodeDuplication float64 `json:"code_duplication"`
	Grade           string  `json:"grade"` // A, B, C, D, F
}

// EnhancedProjectAnalysis contains the complete analysis results for a project with enhanced metrics
type EnhancedProjectAnalysis struct {
	RootPath         string                   `json:"root_path"`
	TotalFiles       int                      `json:"total_files"`
	TotalLines       int                      `json:"total_lines"`
	Languages        map[string]LanguageStats `json:"languages"`
	FileResults      interface{}              `json:"file_results"` // Will be []*parser.AnalysisResult
	Summary          string                   `json:"summary,omitempty"`
	GeneratedAt      time.Time                `json:"generated_at"`
	AnalysisDuration time.Duration            `json:"analysis_duration"`
	// Enhanced metrics
	ProjectMetrics  ProjectMetrics            `json:"project_metrics"`
	DirectoryStats  map[string]DirectoryStats `json:"directory_stats"`
	DependencyGraph DependencyGraph           `json:"dependency_graph"`
	QualityScore    QualityScore              `json:"quality_score"`
}
