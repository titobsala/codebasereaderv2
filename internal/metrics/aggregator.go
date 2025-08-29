package metrics

import (
	"math"
	"path/filepath"
	"strings"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

// Aggregator handles aggregation of metrics across multiple files and directories
type Aggregator struct {
	calculator *Calculator
}

// NewAggregator creates a new metrics aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		calculator: NewCalculator(),
	}
}

// AggregateProjectMetrics calculates comprehensive project-level metrics
func (a *Aggregator) AggregateProjectMetrics(results []*parser.AnalysisResult, rootPath string) *EnhancedProjectAnalysis {
	analysis := &EnhancedProjectAnalysis{
		RootPath:        rootPath,
		TotalFiles:      len(results),
		Languages:       make(map[string]LanguageStats),
		FileResults:     results,
		ProjectMetrics:  ProjectMetrics{},
		DirectoryStats:  make(map[string]DirectoryStats),
		DependencyGraph: DependencyGraph{},
		QualityScore:    QualityScore{},
	}
	
	// Calculate project-level metrics
	a.calculateProjectMetrics(analysis)

	// Calculate directory-level statistics
	a.calculateDirectoryStats(analysis)

	// Analyze dependency relationships
	a.analyzeDependencyGraph(analysis)

	// Calculate overall quality score
	a.calculateOverallQualityScore(analysis)

	return analysis
}

// calculateProjectMetrics calculates overall project metrics
func (a *Aggregator) calculateProjectMetrics(analysis *EnhancedProjectAnalysis) {
	var totalComplexity, maxComplexity int
	var totalMaintainability, totalTechnicalDebt float64
	var totalCodeLines, totalCommentLines int
	var functionsWithDocs, totalFunctions int

	// Safe type assertion with check
	fileResults, ok := analysis.FileResults.([]*parser.AnalysisResult)
	if !ok {
		// Handle the case where FileResults is not the expected type
		return
	}
	for _, result := range fileResults {
		totalComplexity += result.CyclomaticComplexity
		totalMaintainability += result.MaintainabilityIndex
		totalTechnicalDebt += result.TechnicalDebt
		totalCodeLines += result.CodeLines
		totalCommentLines += result.CommentLines

		// Track maximum complexity
		if result.CyclomaticComplexity > maxComplexity {
			maxComplexity = result.CyclomaticComplexity
		}

		// Count documented functions
		for _, fn := range result.Functions {
			totalFunctions++
			if fn.HasDocstring {
				functionsWithDocs++
			}
		}

		// Count documented classes and their methods
		for _, class := range result.Classes {
			if class.HasDocstring {
				functionsWithDocs++ // Count class as documented entity
			}
			for _, method := range class.Methods {
				totalFunctions++
				if method.HasDocstring {
					functionsWithDocs++
				}
			}
		}
	}

	fileCount := len(fileResults)
	if fileCount == 0 {
		return
	}

	// Calculate averages and ratios
	avgComplexity := float64(totalComplexity) / float64(fileCount)
	avgMaintainability := totalMaintainability / float64(fileCount)

	var documentationRatio float64
	if totalFunctions > 0 {
		documentationRatio = float64(functionsWithDocs) / float64(totalFunctions) * 100
	}

	var codeToCommentRatio float64
	if totalCommentLines > 0 {
		codeToCommentRatio = float64(totalCodeLines) / float64(totalCommentLines)
	}

	analysis.ProjectMetrics = ProjectMetrics{
		TotalComplexity:      totalComplexity,
		AverageComplexity:    avgComplexity,
		MaxComplexity:        maxComplexity,
		MaintainabilityIndex: avgMaintainability,
		TechnicalDebt:        totalTechnicalDebt,
		DocumentationRatio:   documentationRatio,
		CodeToCommentRatio:   codeToCommentRatio,
		// CodeDuplication and TestCoverage would require more sophisticated analysis
		CodeDuplication: 0.0, // Placeholder
		TestCoverage:    0.0, // Placeholder
	}
}

// calculateDirectoryStats calculates statistics for each directory
func (a *Aggregator) calculateDirectoryStats(analysis *EnhancedProjectAnalysis) {
	dirStats := make(map[string]*DirectoryStats)

	// Safe type assertion with check
	fileResults, ok := analysis.FileResults.([]*parser.AnalysisResult)
	if !ok {
		return
	}
	for _, result := range fileResults {
		dir := filepath.Dir(result.FilePath)

		// Initialize directory stats if not exists
		if _, exists := dirStats[dir]; !exists {
			dirStats[dir] = &DirectoryStats{
				Path:      dir,
				Languages: make(map[string]LanguageStats),
			}
		}

		stats := dirStats[dir]
		stats.FileCount++
		stats.LineCount += result.LineCount
		stats.Complexity += result.CyclomaticComplexity

		// Update language-specific stats for this directory
		langStats := stats.Languages[result.Language]
		langStats.FileCount++
		langStats.LineCount += result.LineCount
		langStats.FunctionCount += len(result.Functions)
		langStats.ClassCount += len(result.Classes)
		langStats.Complexity += result.Complexity
		langStats.CyclomaticComplexity += result.CyclomaticComplexity
		langStats.CodeLines += result.CodeLines
		langStats.CommentLines += result.CommentLines
		langStats.BlankLines += result.BlankLines
		langStats.MaintainabilityIndex += result.MaintainabilityIndex
		langStats.TechnicalDebt += result.TechnicalDebt

		// Check if this is a test file
		if a.isTestFile(result.FilePath) {
			langStats.TestFiles++
		}

		stats.Languages[result.Language] = langStats
	}

	// Calculate averages and finalize stats
	finalStats := make(map[string]DirectoryStats)
	for path, stats := range dirStats {
		if stats.FileCount > 0 {
			stats.MaintainabilityIndex = 0
			totalMaintainability := 0.0

			for lang, langStats := range stats.Languages {
				if langStats.FileCount > 0 {
					langStats.AverageComplexity = float64(langStats.Complexity) / float64(langStats.FileCount)
					langStats.MaintainabilityIndex = langStats.MaintainabilityIndex / float64(langStats.FileCount)
					totalMaintainability += langStats.MaintainabilityIndex
				}
				stats.Languages[lang] = langStats
			}

			stats.MaintainabilityIndex = totalMaintainability / float64(len(stats.Languages))
		}

		finalStats[path] = *stats
	}

	analysis.DirectoryStats = finalStats
}

// analyzeDependencyGraph analyzes project dependency relationships
func (a *Aggregator) analyzeDependencyGraph(analysis *EnhancedProjectAnalysis) {
	internalDeps := make(map[string][]string)
	externalDeps := make(map[string][]string)
	standardDeps := make(map[string][]string)
	allDeps := make(map[string]map[string]bool)

	// Collect all dependencies
	// Safe type assertion with check
	fileResults, ok := analysis.FileResults.([]*parser.AnalysisResult)
	if !ok {
		return
	}
	for _, result := range fileResults {
		fileDeps := make(map[string]bool)

		for _, dep := range result.Dependencies {
			switch dep.Type {
			case "internal":
				if _, exists := internalDeps[result.FilePath]; !exists {
					internalDeps[result.FilePath] = []string{}
				}
				internalDeps[result.FilePath] = append(internalDeps[result.FilePath], dep.Name)
				fileDeps[dep.Name] = true
			case "external":
				if _, exists := externalDeps[result.FilePath]; !exists {
					externalDeps[result.FilePath] = []string{}
				}
				externalDeps[result.FilePath] = append(externalDeps[result.FilePath], dep.Name)
				fileDeps[dep.Name] = true
			case "standard":
				if _, exists := standardDeps[result.FilePath]; !exists {
					standardDeps[result.FilePath] = []string{}
				}
				standardDeps[result.FilePath] = append(standardDeps[result.FilePath], dep.Name)
				fileDeps[dep.Name] = true
			}
		}

		allDeps[result.FilePath] = fileDeps
	}

	// Detect circular dependencies (simplified detection)
	circularDeps := a.detectCircularDependencies(internalDeps)

	// Calculate dependency depth
	maxDepth := a.calculateDependencyDepth(internalDeps)

	analysis.DependencyGraph = DependencyGraph{
		InternalDependencies: internalDeps,
		ExternalDependencies: externalDeps,
		StandardDependencies: standardDeps,
		CircularDependencies: circularDeps,
		DependencyDepth:      maxDepth,
		UnusedDependencies:   []string{}, // Would require more sophisticated analysis
	}
}

// detectCircularDependencies detects circular dependency chains
func (a *Aggregator) detectCircularDependencies(deps map[string][]string) [][]string {
	// Simplified circular dependency detection with safety limits
	var circular [][]string
	
	// Safety limit to prevent infinite loops
	const maxDepth = 100
	
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	var dfs func(string, []string, int) bool
	dfs = func(node string, path []string, depth int) bool {
		// Safety check for maximum depth
		if depth > maxDepth {
			return false
		}
		
		if recursionStack[node] {
			// Found a cycle, extract the cycle
			cycleStart := -1
			for i, p := range path {
				if p == node {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 && cycleStart < len(path) {
				cycle := make([]string, len(path)-cycleStart+1)
				copy(cycle, path[cycleStart:])
				cycle[len(cycle)-1] = node
				circular = append(circular, cycle)
			}
			return true
		}

		if visited[node] {
			return false
		}

		visited[node] = true
		recursionStack[node] = true
		newPath := make([]string, len(path)+1)
		copy(newPath, path)
		newPath[len(path)] = node

		// Only process dependencies that exist in our graph
		if nodeDeps, exists := deps[node]; exists {
			for _, dep := range nodeDeps {
				if dfs(dep, newPath, depth+1) {
					break // Stop after finding first cycle to prevent excessive processing
				}
			}
		}

		recursionStack[node] = false
		return false
	}

	// Process each node with safety limits
	processedCount := 0
	const maxNodes = 1000
	
	for node := range deps {
		if processedCount >= maxNodes {
			break // Safety limit on number of nodes processed
		}
		if !visited[node] {
			dfs(node, []string{}, 0)
		}
		processedCount++
	}

	return circular
}

// calculateDependencyDepth calculates the maximum dependency depth
func (a *Aggregator) calculateDependencyDepth(deps map[string][]string) int {
	maxDepth := 0
	const maxAllowedDepth = 50 // Safety limit

	var calculateDepth func(string, map[string]bool, int) int
	calculateDepth = func(node string, visited map[string]bool, currentDepth int) int {
		// Safety checks
		if currentDepth > maxAllowedDepth {
			return currentDepth
		}
		
		if visited[node] {
			return currentDepth // Avoid infinite recursion
		}

		visited[node] = true
		depth := currentDepth

		// Only process if node exists in deps
		if nodeDeps, exists := deps[node]; exists {
			for _, dep := range nodeDeps {
				depDepth := calculateDepth(dep, visited, currentDepth+1)
				if depDepth > depth {
					depth = depDepth
				}
			}
		}

		delete(visited, node)
		return depth
	}

	// Process with safety limits
	processedCount := 0
	const maxNodes = 1000
	
	for node := range deps {
		if processedCount >= maxNodes {
			break
		}
		depth := calculateDepth(node, make(map[string]bool), 0)
		if depth > maxDepth {
			maxDepth = depth
		}
		processedCount++
	}

	return maxDepth
}

// calculateOverallQualityScore calculates the overall project quality score
func (a *Aggregator) calculateOverallQualityScore(analysis *EnhancedProjectAnalysis) {
	metrics := analysis.ProjectMetrics

	// Normalize complexity score (0-100, lower complexity is better)
	complexityScore := math.Max(0, 100-metrics.AverageComplexity*5)
	if complexityScore > 100 {
		complexityScore = 100
	}

	overall, grade := a.calculator.CalculateQualityScore(
		metrics.MaintainabilityIndex,
		complexityScore,
		metrics.DocumentationRatio,
		metrics.TestCoverage,
		metrics.CodeDuplication,
	)

	analysis.QualityScore = QualityScore{
		Overall:         overall,
		Maintainability: metrics.MaintainabilityIndex,
		Complexity:      complexityScore,
		Documentation:   metrics.DocumentationRatio,
		TestCoverage:    metrics.TestCoverage,
		CodeDuplication: metrics.CodeDuplication,
		Grade:           grade,
	}
}

// isTestFile determines if a file is a test file based on naming conventions
func (a *Aggregator) isTestFile(filePath string) bool {
	fileName := filepath.Base(filePath)
	lowerName := strings.ToLower(fileName)

	// Common test file patterns
	testPatterns := []string{
		"_test.go",
		"test_",
		"_test.py",
		".test.",
		"spec.",
		"_spec.",
	}

	for _, pattern := range testPatterns {
		if strings.Contains(lowerName, pattern) {
			return true
		}
	}

	return false
}
