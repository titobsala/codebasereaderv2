package metrics

import (
	"math"
	"regexp"
	"strings"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

// Calculator provides methods for calculating various code metrics
type Calculator struct {
	commentPatterns map[string]*regexp.Regexp
}

// NewCalculator creates a new metrics calculator
func NewCalculator() *Calculator {
	return &Calculator{
		commentPatterns: map[string]*regexp.Regexp{
			"go":     regexp.MustCompile(`^\s*//|/\*[\s\S]*?\*/`),
			"python": regexp.MustCompile(`^\s*#|'''[\s\S]*?'''|"""[\s\S]*?"""`),
		},
	}
}

// CalculateFileMetrics calculates comprehensive metrics for a single file
func (c *Calculator) CalculateFileMetrics(result *parser.AnalysisResult, content []byte) {
	lines := strings.Split(string(content), "\n")
	
	// Calculate line-based metrics
	c.calculateLineMetrics(result, lines)
	
	// Calculate complexity metrics
	c.calculateComplexityMetrics(result)
	
	// Calculate maintainability index
	c.calculateMaintainabilityIndex(result)
	
	// Calculate technical debt
	c.calculateTechnicalDebt(result)
	
	// Calculate dependency metrics
	c.calculateDependencyMetrics(result)
}

// calculateLineMetrics calculates line-based metrics
func (c *Calculator) calculateLineMetrics(result *parser.AnalysisResult, lines []string) {
	var codeLines, commentLines, blankLines int
	var totalLineLength int
	var maxLineLength int
	
	commentPattern := c.commentPatterns[strings.ToLower(result.Language)]
	
	for _, line := range lines {
		lineLength := len(line)
		totalLineLength += lineLength
		
		if lineLength > maxLineLength {
			maxLineLength = lineLength
		}
		
		trimmedLine := strings.TrimSpace(line)
		
		if trimmedLine == "" {
			blankLines++
		} else if commentPattern != nil && commentPattern.MatchString(trimmedLine) {
			commentLines++
		} else {
			codeLines++
		}
	}
	
	result.CodeLines = codeLines
	result.CommentLines = commentLines
	result.BlankLines = blankLines
	result.MaxLineLength = maxLineLength
	
	if len(lines) > 0 {
		result.AverageLineLength = float64(totalLineLength) / float64(len(lines))
	}
}

// calculateComplexityMetrics calculates various complexity metrics
func (c *Calculator) calculateComplexityMetrics(result *parser.AnalysisResult) {
	totalCyclomatic := 0
	
	// Calculate cyclomatic complexity for functions
	for i, fn := range result.Functions {
		result.Functions[i].CyclomaticComplexity = fn.Complexity
		result.Functions[i].LinesOfCode = fn.LineEnd - fn.LineStart + 1
		result.Functions[i].ParameterCount = len(fn.Parameters)
		totalCyclomatic += fn.Complexity
	}
	
	// Calculate cyclomatic complexity for classes and their methods
	for i, class := range result.Classes {
		classComplexity := 0
		for j, method := range class.Methods {
			result.Classes[i].Methods[j].CyclomaticComplexity = method.Complexity
			result.Classes[i].Methods[j].LinesOfCode = method.LineEnd - method.LineStart + 1
			result.Classes[i].Methods[j].ParameterCount = len(method.Parameters)
			classComplexity += method.Complexity
			totalCyclomatic += method.Complexity
		}
		result.Classes[i].Complexity = classComplexity
		result.Classes[i].LinesOfCode = class.LineEnd - class.LineStart + 1
		result.Classes[i].MethodCount = len(class.Methods)
		result.Classes[i].FieldCount = len(class.Fields)
	}
	
	result.CyclomaticComplexity = totalCyclomatic
}

// calculateMaintainabilityIndex calculates the maintainability index
// Based on the formula: MI = 171 - 5.2 * ln(HV) - 0.23 * CC - 16.2 * ln(LOC)
// Where HV = Halstead Volume, CC = Cyclomatic Complexity, LOC = Lines of Code
func (c *Calculator) calculateMaintainabilityIndex(result *parser.AnalysisResult) {
	if result.CodeLines == 0 {
		result.MaintainabilityIndex = 100.0
		return
	}
	
	// Simplified calculation without full Halstead metrics
	// Using approximations based on available data
	halsteadVolume := float64(result.CodeLines) * 2.0 // Approximation
	cyclomaticComplexity := float64(result.CyclomaticComplexity)
	linesOfCode := float64(result.CodeLines)
	
	if halsteadVolume <= 0 {
		halsteadVolume = 1.0
	}
	if linesOfCode <= 0 {
		linesOfCode = 1.0
	}
	
	mi := 171.0 - 5.2*math.Log(halsteadVolume) - 0.23*cyclomaticComplexity - 16.2*math.Log(linesOfCode)
	
	// Normalize to 0-100 scale
	if mi < 0 {
		mi = 0
	} else if mi > 100 {
		mi = 100
	}
	
	result.MaintainabilityIndex = mi
}

// calculateTechnicalDebt calculates technical debt score
func (c *Calculator) calculateTechnicalDebt(result *parser.AnalysisResult) {
	debt := 0.0
	
	// High complexity functions contribute to technical debt
	for _, fn := range result.Functions {
		if fn.CyclomaticComplexity > 10 {
			debt += float64(fn.CyclomaticComplexity-10) * 0.5
		}
		if fn.ParameterCount > 5 {
			debt += float64(fn.ParameterCount-5) * 0.3
		}
		if fn.LinesOfCode > 50 {
			debt += float64(fn.LinesOfCode-50) * 0.1
		}
	}
	
	// Large classes contribute to technical debt
	for _, class := range result.Classes {
		if class.MethodCount > 20 {
			debt += float64(class.MethodCount-20) * 0.2
		}
		if class.LinesOfCode > 200 {
			debt += float64(class.LinesOfCode-200) * 0.05
		}
	}
	
	// Long lines contribute to technical debt
	if result.MaxLineLength > 120 {
		debt += float64(result.MaxLineLength-120) * 0.01
	}
	
	// Low comment ratio contributes to technical debt
	if result.CodeLines > 0 {
		commentRatio := float64(result.CommentLines) / float64(result.CodeLines)
		if commentRatio < 0.1 {
			debt += (0.1 - commentRatio) * 10.0
		}
	}
	
	result.TechnicalDebt = debt
}

// calculateDependencyMetrics calculates dependency-related metrics
func (c *Calculator) calculateDependencyMetrics(result *parser.AnalysisResult) {
	result.ImportCount = len(result.Imports)
	
	// Analyze dependencies
	dependencies := make([]parser.Dependency, 0)
	
	for _, imp := range result.Imports {
		dep := parser.Dependency{
			Name:        imp,
			Type:        c.classifyDependency(imp, result.Language),
			UsageCount:  1, // Simplified - would need more analysis for actual usage
			IsDirectDep: true,
			FilePath:    result.FilePath,
		}
		dependencies = append(dependencies, dep)
	}
	
	result.Dependencies = dependencies
}

// classifyDependency classifies a dependency as standard, internal, or external
func (c *Calculator) classifyDependency(importPath, language string) string {
	switch strings.ToLower(language) {
	case "go":
		if strings.Contains(importPath, ".") {
			if strings.HasPrefix(importPath, "github.com") || 
			   strings.HasPrefix(importPath, "golang.org") ||
			   strings.HasPrefix(importPath, "google.golang.org") {
				return "external"
			}
			return "internal"
		}
		// Check for relative imports or internal packages
		if strings.Contains(importPath, "/") {
			return "internal"
		}
		return "standard"
	case "python":
		standardLibs := map[string]bool{
			"os": true, "sys": true, "json": true, "re": true, "time": true,
			"datetime": true, "collections": true, "itertools": true, "functools": true,
			"math": true, "random": true, "string": true, "urllib": true, "http": true,
		}
		
		if standardLibs[importPath] {
			return "standard"
		}
		
		if strings.Contains(importPath, ".") {
			return "internal"
		}
		
		return "external"
	default:
		return "unknown"
	}
}

// CalculateQualityScore calculates an overall quality score
func (c *Calculator) CalculateQualityScore(maintainability, complexity, documentation, testCoverage, duplication float64) (float64, string) {
	// Weighted average of different quality factors
	weights := map[string]float64{
		"maintainability": 0.3,
		"complexity":      0.25,
		"documentation":   0.2,
		"testCoverage":    0.15,
		"duplication":     0.1,
	}
	
	// Normalize complexity score (lower is better)
	normalizedComplexity := math.Max(0, 100-complexity)
	
	// Normalize duplication score (lower is better)
	normalizedDuplication := math.Max(0, 100-duplication)
	
	score := maintainability*weights["maintainability"] +
		normalizedComplexity*weights["complexity"] +
		documentation*weights["documentation"] +
		testCoverage*weights["testCoverage"] +
		normalizedDuplication*weights["duplication"]
	
	// Determine grade
	var grade string
	switch {
	case score >= 90:
		grade = "A"
	case score >= 80:
		grade = "B"
	case score >= 70:
		grade = "C"
	case score >= 60:
		grade = "D"
	default:
		grade = "F"
	}
	
	return score, grade
}