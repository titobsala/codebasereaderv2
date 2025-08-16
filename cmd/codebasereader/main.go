package main

import (
	"fmt"
	"log"

	"github.com/tito-sala/codebasereaderv2/internal/core"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
)

func main() {
	fmt.Println("CodebaseReader v2 – Initializing...")

	// Create application with default configuration
	config := engine.DefaultConfig()
	app := core.NewApplication(config)

	// Validate setup
	if err := app.ValidateSetup(); err != nil {
		log.Printf("Setup validation failed: %v", err)
		log.Println("Note: No parsers registered yet. This is expected at this stage.")
	}

	// Display configuration
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Max Workers: %d\n", config.MaxWorkers)
	fmt.Printf("  Max File Size: %d bytes\n", config.MaxFileSize)
	fmt.Printf("  AI Provider: %s\n", config.AIProvider)
	fmt.Printf("  Exclude Patterns: %v\n", config.ExcludePatterns)

	// Display supported languages (will be empty until parsers are registered)
	languages := app.GetSupportedLanguages()
	fmt.Printf("Supported Languages: %d registered\n", len(languages))
	for lang, exts := range languages {
		fmt.Printf("  %s: %s\n", lang, exts)
	}

	fmt.Println("Core interfaces and project structure initialized successfully!")
}
