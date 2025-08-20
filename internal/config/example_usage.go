package config

import (
	"fmt"
	"log"
)

// ExampleUsage demonstrates how to use the configuration system
func ExampleUsage() {
	// Load configuration from default path or environment
	config, err := LoadConfig(GetDefaultConfigPath())
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Configuration loaded successfully:\n")
	fmt.Printf("  AI Provider: %s\n", config.AIProvider)
	fmt.Printf("  Max Workers: %d\n", config.MaxWorkers)
	fmt.Printf("  Output Format: %s\n", config.OutputFormat)
	fmt.Printf("  Exclude Patterns: %v\n", config.ExcludePatterns)

	// API key should not be logged for security
	if config.APIKey != "" {
		fmt.Printf("  API Key: [CONFIGURED]\n")
	} else {
		fmt.Printf("  API Key: [NOT SET]\n")
	}

	// Save configuration to a specific path
	if err := config.Save("./my-config.json"); err != nil {
		log.Printf("Failed to save config: %v", err)
	} else {
		fmt.Printf("Configuration saved to ./my-config.json\n")
	}
}
