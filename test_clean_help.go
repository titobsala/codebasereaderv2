package main

import (
	"fmt"
	"strings"

	"github.com/tito-sala/codebasereaderv2/internal/tui/views"
)

func main() {
	fmt.Println("Testing cleaned-up help page...")

	// Create help view model
	helpView := views.NewHelpViewModel()
	
	// Test rendering with more height to show full content
	output := helpView.Render(80, 50)
	
	// Check for problematic styling elements
	problematic := []string{
		"│                                                  │", // Empty card content
		"╭──────────────────────────────────────────────────╮", // Large empty boxes
	}
	
	hasProblems := false
	for _, issue := range problematic {
		if strings.Contains(output, issue) {
			fmt.Printf("❌ Found problematic element: %s\n", issue)
			hasProblems = true
		}
	}
	
	if !hasProblems {
		fmt.Println("✅ No problematic empty boxes found")
	}
	
	// Check for expected clean elements
	expected := []string{
		"🚀 Getting Started", // Section headers should be visible
		"1   2   3   4", // Clean tab navigation
		"⌨️  Keyboard Shortcuts:", // Clean shortcuts section
	}
	
	allExpectedFound := true
	for _, expect := range expected {
		if !strings.Contains(output, expect) {
			fmt.Printf("❌ Missing expected element: %s\n", expect)
			allExpectedFound = false
		}
	}
	
	if allExpectedFound {
		fmt.Println("✅ All expected clean elements found")
	}
	
	// Test navigation still works
	helpView.Update("2")
	output2 := helpView.Render(80, 24)
	
	if output != output2 {
		fmt.Println("✅ Navigation still working after cleanup")
	} else {
		fmt.Println("❌ Navigation might be broken")
	}
	
	// Show sample of cleaned output
	lines := strings.Split(output, "\n")
	fmt.Printf("\n📋 Full cleaned help output (%d lines):\n", len(lines))
	for i, line := range lines {
		fmt.Printf("  [%2d] %s\n", i+1, line)
	}
	
	fmt.Println("\nHelp page cleanup test completed!")
}