package parser

import (
	"testing"
)

func TestGoParser_GetLanguageName(t *testing.T) {
	parser := NewGoParser()
	if parser.GetLanguageName() != "Go" {
		t.Errorf("Expected language name 'Go', got '%s'", parser.GetLanguageName())
	}
}

func TestGoParser_GetSupportedExtensions(t *testing.T) {
	parser := NewGoParser()
	extensions := parser.GetSupportedExtensions()
	
	if len(extensions) != 1 {
		t.Errorf("Expected 1 supported extension, got %d", len(extensions))
	}
	
	if extensions[0] != ".go" {
		t.Errorf("Expected extension '.go', got '%s'", extensions[0])
	}
}

func TestGoParser_ParseSimpleFunction(t *testing.T) {
	parser := NewGoParser()
	
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

func add(a, b int) int {
	if a > 0 {
		return a + b
	}
	return b
}`

	result, err := parser.Parse("test.go", []byte(code))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check basic properties
	if result.Language != "Go" {
		t.Errorf("Expected language 'Go', got '%s'", result.Language)
	}

	if result.FilePath != "test.go" {
		t.Errorf("Expected file path 'test.go', got '%s'", result.FilePath)
	}

	if result.LineCount == 0 {
		t.Error("Expected line count > 0")
	}

	// Check imports
	if len(result.Imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(result.Imports))
	}

	if result.Imports[0] != "fmt" {
		t.Errorf("Expected import 'fmt', got '%s'", result.Imports[0])
	}

	// Check functions
	if len(result.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(result.Functions))
	}

	// Check main function
	mainFunc := findFunction(result.Functions, "main")
	if mainFunc == nil {
		t.Error("Expected to find 'main' function")
	} else {
		if len(mainFunc.Parameters) != 0 {
			t.Errorf("Expected main function to have 0 parameters, got %d", len(mainFunc.Parameters))
		}
		if mainFunc.ReturnType != "" {
			t.Errorf("Expected main function to have no return type, got '%s'", mainFunc.ReturnType)
		}
	}

	// Check add function
	addFunc := findFunction(result.Functions, "add")
	if addFunc == nil {
		t.Error("Expected to find 'add' function")
	} else {
		if len(addFunc.Parameters) != 2 {
			t.Errorf("Expected add function to have 2 parameters, got %d", len(addFunc.Parameters))
		}
		if addFunc.ReturnType != "int" {
			t.Errorf("Expected add function return type 'int', got '%s'", addFunc.ReturnType)
		}
		if addFunc.Complexity < 2 {
			t.Errorf("Expected add function complexity >= 2 (has if statement), got %d", addFunc.Complexity)
		}
	}
}

func TestGoParser_ParseStruct(t *testing.T) {
	parser := NewGoParser()
	
	code := `package main

type Person struct {
	Name string
	Age  int
	Email string
}

type Address struct {
	Street string
	City   string
	Country string
}

func (p Person) GetName() string {
	return p.Name
}

func (p *Person) SetAge(age int) {
	p.Age = age
}`

	result, err := parser.Parse("test.go", []byte(code))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check structs
	if len(result.Classes) != 2 {
		t.Errorf("Expected 2 structs, got %d", len(result.Classes))
	}

	// Check Person struct
	personStruct := findClass(result.Classes, "Person")
	if personStruct == nil {
		t.Error("Expected to find 'Person' struct")
	} else {
		if len(personStruct.Fields) != 3 {
			t.Errorf("Expected Person struct to have 3 fields, got %d", len(personStruct.Fields))
		}
		
		expectedFields := []string{"Name string", "Age int", "Email string"}
		for i, expected := range expectedFields {
			if i < len(personStruct.Fields) && personStruct.Fields[i] != expected {
				t.Errorf("Expected field '%s', got '%s'", expected, personStruct.Fields[i])
			}
		}
	}

	// Check functions (methods are also captured as functions)
	if len(result.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(result.Functions))
	}

	// Check GetName method
	getNameFunc := findFunction(result.Functions, "GetName")
	if getNameFunc == nil {
		t.Error("Expected to find 'GetName' method")
	} else {
		// Method receivers are not included in parameters in go/ast
		if len(getNameFunc.Parameters) != 0 {
			t.Errorf("Expected GetName to have 0 parameters, got %d", len(getNameFunc.Parameters))
		}
		if getNameFunc.ReturnType != "string" {
			t.Errorf("Expected GetName return type 'string', got '%s'", getNameFunc.ReturnType)
		}
	}
}

func TestGoParser_ParseInterface(t *testing.T) {
	parser := NewGoParser()
	
	code := `package main

type Writer interface {
	Write([]byte) (int, error)
}

type ReadWriter interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}`

	result, err := parser.Parse("test.go", []byte(code))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check interfaces
	if len(result.Classes) != 2 {
		t.Errorf("Expected 2 interfaces, got %d", len(result.Classes))
	}

	// Check Writer interface
	writerInterface := findClass(result.Classes, "Writer")
	if writerInterface == nil {
		t.Error("Expected to find 'Writer' interface")
	} else {
		if len(writerInterface.Methods) != 1 {
			t.Errorf("Expected Writer interface to have 1 method, got %d", len(writerInterface.Methods))
		}
		
		if len(writerInterface.Methods) > 0 {
			writeMethod := writerInterface.Methods[0]
			if writeMethod.Name != "Write" {
				t.Errorf("Expected method name 'Write', got '%s'", writeMethod.Name)
			}
			if writeMethod.ReturnType != "(int, error)" {
				t.Errorf("Expected Write method return type '(int, error)', got '%s'", writeMethod.ReturnType)
			}
		}
	}

	// Check ReadWriter interface
	readWriterInterface := findClass(result.Classes, "ReadWriter")
	if readWriterInterface == nil {
		t.Error("Expected to find 'ReadWriter' interface")
	} else {
		if len(readWriterInterface.Methods) != 3 {
			t.Errorf("Expected ReadWriter interface to have 3 methods, got %d", len(readWriterInterface.Methods))
		}
	}
}

func TestGoParser_ParseComplexFunction(t *testing.T) {
	parser := NewGoParser()
	
	code := `package main

func complexFunction(items []string, threshold int) ([]string, error) {
	var result []string
	
	for i, item := range items {
		if len(item) > threshold {
			switch item[0] {
			case 'A':
				result = append(result, "Type A: " + item)
			case 'B':
				result = append(result, "Type B: " + item)
			default:
				result = append(result, "Other: " + item)
			}
		} else if i%2 == 0 {
			result = append(result, "Even: " + item)
		}
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("no items processed")
	}
	
	return result, nil
}`

	result, err := parser.Parse("test.go", []byte(code))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result.Functions))
	}

	complexFunc := result.Functions[0]
	if complexFunc.Name != "complexFunction" {
		t.Errorf("Expected function name 'complexFunction', got '%s'", complexFunc.Name)
	}

	// Check parameters
	if len(complexFunc.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(complexFunc.Parameters))
	}

	// Check return type
	if complexFunc.ReturnType != "([]string, error)" {
		t.Errorf("Expected return type '([]string, error)', got '%s'", complexFunc.ReturnType)
	}

	// Check complexity (should be high due to multiple control structures)
	if complexFunc.Complexity < 5 {
		t.Errorf("Expected complexity >= 5, got %d", complexFunc.Complexity)
	}
}

func TestGoParser_ParseErrors(t *testing.T) {
	parser := NewGoParser()
	
	// Test with invalid Go code
	invalidCode := `package main

func main() {
	fmt.Println("Hello, World!"
	// Missing closing parenthesis
}`

	result, err := parser.Parse("test.go", []byte(invalidCode))
	// Parser should not return an error, but should capture parse errors
	if err != nil {
		t.Fatalf("Parse should not return error for invalid code: %v", err)
	}

	// Should have parse errors recorded
	if len(result.Errors) == 0 {
		t.Error("Expected parse errors to be recorded")
	}
}

func TestGoParser_ParseMultipleImports(t *testing.T) {
	parser := NewGoParser()
	
	code := `package main

import (
	"fmt"
	"strings"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Hello")
}`

	result, err := parser.Parse("test.go", []byte(code))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedImports := []string{"fmt", "strings", "os", "path/filepath"}
	if len(result.Imports) != len(expectedImports) {
		t.Errorf("Expected %d imports, got %d", len(expectedImports), len(result.Imports))
	}

	for i, expected := range expectedImports {
		if i < len(result.Imports) && result.Imports[i] != expected {
			t.Errorf("Expected import '%s', got '%s'", expected, result.Imports[i])
		}
	}
}

// Helper functions for tests

func findFunction(functions []FunctionInfo, name string) *FunctionInfo {
	for i := range functions {
		if functions[i].Name == name {
			return &functions[i]
		}
	}
	return nil
}

func findClass(classes []ClassInfo, name string) *ClassInfo {
	for i := range classes {
		if classes[i].Name == name {
			return &classes[i]
		}
	}
	return nil
}