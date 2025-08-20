package ai

import "github.com/tito-sala/codebasereaderv2/internal/engine"

// AIClient defines the interface for AI service providers
type AIClient interface {
	// GenerateSummary generates a summary for the given project context
	GenerateSummary(request *AIRequest) (*AIResponse, error)

	// GetProviderName returns the name of the AI provider
	GetProviderName() string

	// IsConfigured returns true if the client is properly configured
	IsConfigured() bool
}

// AIRequest contains the data needed for AI analysis
type AIRequest struct {
	Provider string         `json:"provider"`
	Model    string         `json:"model"`
	Prompt   string         `json:"prompt"`
	Context  ProjectContext `json:"context"`
}

// ProjectContext contains structured information about the project for AI analysis
type ProjectContext struct {
	Structure string                  `json:"structure"`
	KeyFiles  []string                `json:"key_files"`
	Languages []string                `json:"languages"`
	Metrics   *engine.ProjectAnalysis `json:"metrics"`
	FileCount int                     `json:"file_count"`
	LineCount int                     `json:"line_count"`
}

// AIResponse contains the AI-generated analysis results
type AIResponse struct {
	Summary     string   `json:"summary"`
	Insights    []string `json:"insights"`
	Suggestions []string `json:"suggestions"`
	Error       error    `json:"error,omitempty"`
}

// PromptTemplate defines templates for different types of AI requests
type PromptTemplate struct {
	Name        string
	Template    string
	Description string
}

// Common prompt templates
var (
	ProjectSummaryPrompt = PromptTemplate{
		Name: "project_summary",
		Template: `Analyze this codebase and provide a comprehensive summary.

Project Structure:
%s

Languages: %s
Total Files: %d
Total Lines: %d

Key Files:
%s

Please provide:
1. A brief overview of what this project does
2. The main technologies and frameworks used
3. Key architectural patterns or design decisions
4. Areas that might need attention or improvement
5. Overall code quality assessment`,
		Description: "Generates a comprehensive project summary",
	}

	CodeQualityPrompt = PromptTemplate{
		Name: "code_quality",
		Template: `Analyze the code quality of this project based on the metrics provided.

Metrics Summary:
- Total Files: %d
- Total Lines: %d
- Average Complexity: %.2f
- Languages: %s

Please assess:
1. Code complexity and maintainability
2. Potential technical debt areas
3. Suggestions for improvement
4. Best practices that are being followed or missing`,
		Description: "Focuses on code quality assessment",
	}
)
