package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.AIProvider != "anthropic" {
		t.Errorf("Expected default AI provider to be 'anthropic', got '%s'", config.AIProvider)
	}

	if config.APIKey != "" {
		t.Errorf("Expected default API key to be empty, got '%s'", config.APIKey)
	}

	if config.MaxWorkers != runtime.NumCPU() {
		t.Errorf("Expected default max workers to be %d, got %d", runtime.NumCPU(), config.MaxWorkers)
	}

	if config.OutputFormat != "json" {
		t.Errorf("Expected default output format to be 'json', got '%s'", config.OutputFormat)
	}

	expectedPatterns := []string{"node_modules", ".git", "vendor", "__pycache__", ".pytest_cache", "*.pyc"}
	if len(config.ExcludePatterns) != len(expectedPatterns) {
		t.Errorf("Expected %d exclude patterns, got %d", len(expectedPatterns), len(config.ExcludePatterns))
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testConfig := &Config{
		AIProvider:      "openai",
		APIKey:          "test-key",
		MaxWorkers:      4,
		OutputFormat:    "yaml",
		ExcludePatterns: []string{"test", "build"},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.AIProvider != "openai" {
		t.Errorf("Expected AI provider 'openai', got '%s'", config.AIProvider)
	}

	if config.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", config.APIKey)
	}

	if config.MaxWorkers != 4 {
		t.Errorf("Expected max workers 4, got %d", config.MaxWorkers)
	}

	if config.OutputFormat != "yaml" {
		t.Errorf("Expected output format 'yaml', got '%s'", config.OutputFormat)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("CODEBASE_AI_PROVIDER", "openai-compatible")
	os.Setenv("CODEBASE_API_KEY", "env-key")
	os.Setenv("CODEBASE_MAX_WORKERS", "8")
	os.Setenv("CODEBASE_OUTPUT_FORMAT", "text")
	os.Setenv("CODEBASE_EXCLUDE_PATTERNS", "dist, build, tmp")

	defer func() {
		os.Unsetenv("CODEBASE_AI_PROVIDER")
		os.Unsetenv("CODEBASE_API_KEY")
		os.Unsetenv("CODEBASE_MAX_WORKERS")
		os.Unsetenv("CODEBASE_OUTPUT_FORMAT")
		os.Unsetenv("CODEBASE_EXCLUDE_PATTERNS")
	}()

	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load config from env: %v", err)
	}

	if config.AIProvider != "openai-compatible" {
		t.Errorf("Expected AI provider 'openai-compatible', got '%s'", config.AIProvider)
	}

	if config.APIKey != "env-key" {
		t.Errorf("Expected API key 'env-key', got '%s'", config.APIKey)
	}

	if config.MaxWorkers != 8 {
		t.Errorf("Expected max workers 8, got %d", config.MaxWorkers)
	}

	if config.OutputFormat != "text" {
		t.Errorf("Expected output format 'text', got '%s'", config.OutputFormat)
	}

	expectedPatterns := []string{"dist", "build", "tmp"}
	if len(config.ExcludePatterns) != len(expectedPatterns) {
		t.Errorf("Expected %d exclude patterns, got %d", len(expectedPatterns), len(config.ExcludePatterns))
	}

	for i, pattern := range expectedPatterns {
		if config.ExcludePatterns[i] != pattern {
			t.Errorf("Expected exclude pattern '%s', got '%s'", pattern, config.ExcludePatterns[i])
		}
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid config",
			config:      DefaultConfig(),
			expectError: false,
		},
		{
			name: "Invalid AI provider",
			config: &Config{
				AIProvider:   "invalid",
				MaxWorkers:   4,
				OutputFormat: "json",
			},
			expectError: true,
			errorMsg:    "invalid ai_provider",
		},
		{
			name: "Zero max workers",
			config: &Config{
				AIProvider:   "anthropic",
				MaxWorkers:   0,
				OutputFormat: "json",
			},
			expectError: true,
			errorMsg:    "max_workers must be greater than 0",
		},
		{
			name: "Too many max workers",
			config: &Config{
				AIProvider:   "anthropic",
				MaxWorkers:   150,
				OutputFormat: "json",
			},
			expectError: true,
			errorMsg:    "max_workers cannot exceed 100",
		},
		{
			name: "Invalid output format",
			config: &Config{
				AIProvider:   "anthropic",
				MaxWorkers:   4,
				OutputFormat: "invalid",
			},
			expectError: true,
			errorMsg:    "invalid output_format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestConfigSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.json")

	config := &Config{
		AIProvider:      "openai",
		APIKey:          "test-key",
		MaxWorkers:      4,
		OutputFormat:    "yaml",
		ExcludePatterns: []string{"test", "build"},
	}

	err := config.Save(configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Load and verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config file: %v", err)
	}

	var savedConfig Config
	if err := json.Unmarshal(data, &savedConfig); err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if savedConfig.AIProvider != config.AIProvider {
		t.Errorf("Expected AI provider '%s', got '%s'", config.AIProvider, savedConfig.AIProvider)
	}

	if savedConfig.APIKey != config.APIKey {
		t.Errorf("Expected API key '%s', got '%s'", config.APIKey, savedConfig.APIKey)
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	config, err := LoadConfig("/non/existent/path/config.json")
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got: %v", err)
	}

	// Should return default config
	defaultConfig := DefaultConfig()
	if config.AIProvider != defaultConfig.AIProvider {
		t.Errorf("Expected default AI provider, got '%s'", config.AIProvider)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.json")

	// Write invalid JSON
	if err := os.WriteFile(configPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got none")
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()
	if path == "" {
		t.Errorf("Expected non-empty default config path")
	}

	// Should contain .codebasereader directory
	if !strings.Contains(path, ".codebasereader") {
		t.Errorf("Expected path to contain '.codebasereader', got '%s'", path)
	}
}

func TestLoadConfigEnvOverridesFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	fileConfig := &Config{
		AIProvider:   "anthropic",
		APIKey:       "file-key",
		MaxWorkers:   2,
		OutputFormat: "json",
	}

	data, err := json.MarshalIndent(fileConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set environment variables that should override file values
	os.Setenv("CODEBASE_AI_PROVIDER", "openai")
	os.Setenv("CODEBASE_API_KEY", "env-key")
	os.Setenv("CODEBASE_MAX_WORKERS", "8")

	defer func() {
		os.Unsetenv("CODEBASE_AI_PROVIDER")
		os.Unsetenv("CODEBASE_API_KEY")
		os.Unsetenv("CODEBASE_MAX_WORKERS")
	}()

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Environment variables should override file values
	if config.AIProvider != "openai" {
		t.Errorf("Expected env AI provider 'openai', got '%s'", config.AIProvider)
	}

	if config.APIKey != "env-key" {
		t.Errorf("Expected env API key 'env-key', got '%s'", config.APIKey)
	}

	if config.MaxWorkers != 8 {
		t.Errorf("Expected env max workers 8, got %d", config.MaxWorkers)
	}

	// File value should be used where no env override exists
	if config.OutputFormat != "json" {
		t.Errorf("Expected file output format 'json', got '%s'", config.OutputFormat)
	}
}
