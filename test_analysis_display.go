package main

import (
	"fmt"
	"strings"

	"github.com/tito-sala/codebasereaderv2/internal/metrics"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
)

func main() {
	fmt.Println("Testing analysis display rendering...")

	// Create a test analysis
	testAnalysis := &metrics.EnhancedProjectAnalysis{
		RootPath:    "/home/tito-sala/Code/Personal/codebasereaderv2/internal",
		TotalFiles:  47,
		TotalLines:  12500,
		Languages:   make(map[string]metrics.LanguageStats),
		QualityScore: metrics.QualityScore{
			Overall: 85.5,
		},
		ProjectMetrics: metrics.ProjectMetrics{
			TotalComplexity:      150,
			AverageComplexity:    3.2,
			MaxComplexity:        15,
			MaintainabilityIndex: 78.9,
		},
	}

	// Add a language
	testAnalysis.Languages["Go"] = metrics.LanguageStats{
		FileCount:            47,
		LineCount:            12500,
		FunctionCount:        180,
		ClassCount:           25,
		Complexity:           150,
		CyclomaticComplexity: 150,
		AverageComplexity:    3.2,
		MaxComplexity:        15,
		MaintainabilityIndex: 78.9,
		TechnicalDebt:        15.2,
	}

	// Create metrics display
	display := components.NewMetricsDisplay()

	// Test different rendering modes
	fmt.Println("\n=== OVERVIEW MODE ===")
	result := display.Render(testAnalysis, 80, 40)
	fmt.Printf("Result (first 500 chars): %q\n", result[:min(500, len(result))])
	
	// Check for truncation issues
	lines := strings.Split(result, "\n")
	fmt.Printf("Total lines: %d\n", len(lines))
	fmt.Println("First 10 lines:")
	for i, line := range lines {
		if i >= 10 {
			break
		}
		fmt.Printf("  [%d] %q\n", i+1, line)
	}
}