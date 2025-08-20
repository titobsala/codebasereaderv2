package parser

import (
	"testing"
)

func TestParserRegistry(t *testing.T) {
	registry := NewParserRegistry()

	// Test empty registry
	if len(registry.GetSupportedExtensions()) != 0 {
		t.Error("Expected empty registry to have no supported extensions")
	}

	// Create mock parsers
	goParser := &MockParser{
		language:   "Go",
		extensions: []string{".go"},
	}

	pythonParser := &MockParser{
		language:   "Python",
		extensions: []string{".py", ".pyx"},
	}

	// Test parser registration
	err := registry.RegisterParser(goParser)
	if err != nil {
		t.Errorf("Failed to register Go parser: %v", err)
	}

	err = registry.RegisterParser(pythonParser)
	if err != nil {
		t.Errorf("Failed to register Python parser: %v", err)
	}

	// Test GetParser
	parser, err := registry.GetParser("main.go")
	if err != nil {
		t.Errorf("Failed to get parser for .go file: %v", err)
	}
	if parser.GetLanguageName() != "Go" {
		t.Errorf("Expected Go parser, got %s", parser.GetLanguageName())
	}

	parser, err = registry.GetParser("script.py")
	if err != nil {
		t.Errorf("Failed to get parser for .py file: %v", err)
	}
	if parser.GetLanguageName() != "Python" {
		t.Errorf("Expected Python parser, got %s", parser.GetLanguageName())
	}

	// Test unsupported extension
	_, err = registry.GetParser("file.txt")
	if err == nil {
		t.Error("Expected error for unsupported extension")
	}

	// Test IsSupported
	if !registry.IsSupported("main.go") {
		t.Error("Expected .go files to be supported")
	}

	if registry.IsSupported("file.txt") {
		t.Error("Expected .txt files to be unsupported")
	}

	// Test GetSupportedExtensions
	extensions := registry.GetSupportedExtensions()
	if len(extensions) != 3 { // .go, .py, .pyx
		t.Errorf("Expected 3 supported extensions, got %d", len(extensions))
	}

	// Test GetRegisteredParsers
	parsers := registry.GetRegisteredParsers()
	if len(parsers) != 2 { // Go and Python
		t.Errorf("Expected 2 registered parsers, got %d", len(parsers))
	}

	if _, exists := parsers["Go"]; !exists {
		t.Error("Expected Go parser to be registered")
	}

	if _, exists := parsers["Python"]; !exists {
		t.Error("Expected Python parser to be registered")
	}
}

func TestParserRegistryErrors(t *testing.T) {
	registry := NewParserRegistry()

	// Test registering nil parser
	err := registry.RegisterParser(nil)
	if err == nil {
		t.Error("Expected error when registering nil parser")
	}

	// Test parser with no extensions
	emptyParser := &MockParser{
		language:   "Empty",
		extensions: []string{},
	}

	err = registry.RegisterParser(emptyParser)
	if err == nil {
		t.Error("Expected error when registering parser with no extensions")
	}
}
