package parser

import (
	"os"
	"testing"
)

func TestGoParser_IntegrationWithRealFile(t *testing.T) {
	parser := NewGoParser()

	// Test with the actual config.go file
	content, err := os.ReadFile("../config/config.go")
	if err != nil {
		t.Skipf("Could not read config.go file: %v", err)
	}

	result, err := parser.Parse("../config/config.go", content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Basic validation
	if result.Language != "Go" {
		t.Errorf("Expected language 'Go', got '%s'", result.Language)
	}

	if result.LineCount == 0 {
		t.Error("Expected line count > 0")
	}

	if len(result.Functions) == 0 {
		t.Error("Expected to find functions in config.go")
	}

	if len(result.Imports) == 0 {
		t.Error("Expected to find imports in config.go")
	}

	// Should have no parse errors for valid Go code
	if len(result.Errors) > 0 {
		t.Errorf("Expected no parse errors, got %d errors", len(result.Errors))
		for _, err := range result.Errors {
			t.Logf("Parse error: Line %d, Column %d: %s", err.Line, err.Column, err.Message)
		}
	}

	t.Logf("Successfully parsed config.go: %d lines, %d functions, %d imports",
		result.LineCount, len(result.Functions), len(result.Imports))
}

func TestPythonParser_IntegrationWithSampleCode(t *testing.T) {
	parser := NewPythonParser()

	// Test with sample Python code
	content := `#!/usr/bin/env python3
"""Sample Python module for testing."""

import os
from typing import List, Optional

class Calculator:
    """A simple calculator class."""
    
    def __init__(self, name: str):
        self.name = name
    
    def add(self, a: int, b: int) -> int:
        """Add two numbers."""
        return a + b
    
    def divide(self, a: int, b: int) -> Optional[float]:
        """Divide two numbers."""
        if b == 0:
            return None
        return a / b

def main():
    """Main function."""
    calc = Calculator("test")
    result = calc.add(5, 3)
    print(f"Result: {result}")

if __name__ == "__main__":
    main()
`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Basic validation
	if result.Language != "Python" {
		t.Errorf("Expected language 'Python', got '%s'", result.Language)
	}

	if result.LineCount == 0 {
		t.Error("Expected line count > 0")
	}

	// Should find the main function
	if len(result.Functions) == 0 {
		t.Error("Expected to find functions")
	}

	// Should find the Calculator class
	if len(result.Classes) == 0 {
		t.Error("Expected to find classes")
	}

	// Should find imports
	if len(result.Imports) == 0 {
		t.Error("Expected to find imports")
	}

	t.Logf("Successfully parsed Python code: %d lines, %d functions, %d classes, %d imports",
		result.LineCount, len(result.Functions), len(result.Classes), len(result.Imports))
}

func TestGoParser_RegistryIntegration(t *testing.T) {
	registry := NewParserRegistry()
	goParser := NewGoParser()

	// Register the Go parser
	err := registry.RegisterParser(goParser)
	if err != nil {
		t.Fatalf("Failed to register Go parser: %v", err)
	}

	// Test that we can get the parser back
	parser, err := registry.GetParser("test.go")
	if err != nil {
		t.Fatalf("Failed to get parser for .go file: %v", err)
	}

	if parser.GetLanguageName() != "Go" {
		t.Errorf("Expected Go parser, got %s", parser.GetLanguageName())
	}

	// Test that the parser works through the registry
	code := `package main

func hello() string {
	return "Hello, World!"
}`

	result, err := parser.Parse("test.go", []byte(code))
	if err != nil {
		t.Fatalf("Parse through registry failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result.Functions))
	}

	if result.Functions[0].Name != "hello" {
		t.Errorf("Expected function name 'hello', got '%s'", result.Functions[0].Name)
	}
}

func TestPythonParser_RegistryIntegration(t *testing.T) {
	registry := NewParserRegistry()
	pythonParser := NewPythonParser()

	// Register the Python parser
	err := registry.RegisterParser(pythonParser)
	if err != nil {
		t.Fatalf("Failed to register Python parser: %v", err)
	}

	// Test that we can get the parser back
	parser, err := registry.GetParser("test.py")
	if err != nil {
		t.Fatalf("Failed to get parser for .py file: %v", err)
	}

	if parser.GetLanguageName() != "Python" {
		t.Errorf("Expected Python parser, got %s", parser.GetLanguageName())
	}

	// Test that the parser works through the registry
	code := `def greet(name):
    return f"Hello, {name}!"

class Person:
    def __init__(self, name):
        self.name = name`

	result, err := parser.Parse("test.py", []byte(code))
	if err != nil {
		t.Fatalf("Parse through registry failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result.Functions))
	}

	if result.Functions[0].Name != "greet" {
		t.Errorf("Expected function name 'greet', got '%s'", result.Functions[0].Name)
	}

	if len(result.Classes) != 1 {
		t.Errorf("Expected 1 class, got %d", len(result.Classes))
	}

	if result.Classes[0].Name != "Person" {
		t.Errorf("Expected class name 'Person', got '%s'", result.Classes[0].Name)
	}
}
