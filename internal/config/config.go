package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Config represents the application configuration
type Config struct {
	AIProvider      string   `json:"ai_provider"`
	APIKey          string   `json:"api_key"`
	MaxWorkers      int      `json:"max_workers"`
	OutputFormat    string   `json:"output_format"`
	ExcludePatterns []string `json:"exclude_patterns"`
	ConfigPath      string   `json:"-"` // Not serialized, used internally
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		AIProvider:      "anthropic",
		APIKey:          "",
		MaxWorkers:      runtime.NumCPU(),
		OutputFormat:    "json",
		ExcludePatterns: []string{"node_modules", ".git", "vendor", "__pycache__", ".pytest_cache", "*.pyc"},
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()
	config.ConfigPath = configPath

	// Try to load from file if it exists
	if configPath != "" {
		if err := config.loadFromFile(configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	config.loadFromEnv()

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from a JSON file
func (c *Config) loadFromFile(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File doesn't exist, use defaults
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func (c *Config) loadFromEnv() {
	if provider := os.Getenv("CODEBASE_AI_PROVIDER"); provider != "" {
		c.AIProvider = provider
	}

	if apiKey := os.Getenv("CODEBASE_API_KEY"); apiKey != "" {
		c.APIKey = apiKey
	}

	if workersStr := os.Getenv("CODEBASE_MAX_WORKERS"); workersStr != "" {
		if workers, err := strconv.Atoi(workersStr); err == nil && workers > 0 {
			c.MaxWorkers = workers
		}
	}

	if format := os.Getenv("CODEBASE_OUTPUT_FORMAT"); format != "" {
		c.OutputFormat = format
	}

	if excludeStr := os.Getenv("CODEBASE_EXCLUDE_PATTERNS"); excludeStr != "" {
		patterns := strings.Split(excludeStr, ",")
		for i, pattern := range patterns {
			patterns[i] = strings.TrimSpace(pattern)
		}
		c.ExcludePatterns = patterns
	}
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate AI provider
	validProviders := []string{"anthropic", "openai", "openai-compatible"}
	if !contains(validProviders, c.AIProvider) {
		return fmt.Errorf("invalid ai_provider '%s', must be one of: %s", 
			c.AIProvider, strings.Join(validProviders, ", "))
	}

	// Validate max workers
	if c.MaxWorkers <= 0 {
		return fmt.Errorf("max_workers must be greater than 0, got %d", c.MaxWorkers)
	}

	if c.MaxWorkers > 100 {
		return fmt.Errorf("max_workers cannot exceed 100, got %d", c.MaxWorkers)
	}

	// Validate output format
	validFormats := []string{"json", "yaml", "text"}
	if !contains(validFormats, c.OutputFormat) {
		return fmt.Errorf("invalid output_format '%s', must be one of: %s", 
			c.OutputFormat, strings.Join(validFormats, ", "))
	}

	return nil
}

// Save saves the configuration to a file
func (c *Config) Save(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(homeDir, ".codebasereader", "config.json")
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}