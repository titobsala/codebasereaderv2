package parser

import (
	"strings"
	"testing"
)

func TestPythonParser_Parse_BasicFunction(t *testing.T) {
	parser := NewPythonParser()
	content := `def hello_world():
    print("Hello, World!")
    return "done"`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Language != "Python" {
		t.Errorf("Expected language 'Python', got '%s'", result.Language)
	}

	if len(result.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(result.Functions))
	}

	fn := result.Functions[0]
	if fn.Name != "hello_world" {
		t.Errorf("Expected function name 'hello_world', got '%s'", fn.Name)
	}

	if fn.LineStart != 1 {
		t.Errorf("Expected function to start at line 1, got %d", fn.LineStart)
	}

	if len(fn.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(fn.Parameters))
	}
}

func TestPythonParser_Parse_FunctionWithParameters(t *testing.T) {
	parser := NewPythonParser()
	content := `def add_numbers(a: int, b: int = 5) -> int:
    return a + b`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(result.Functions))
	}

	fn := result.Functions[0]
	if fn.Name != "add_numbers" {
		t.Errorf("Expected function name 'add_numbers', got '%s'", fn.Name)
	}

	expectedParams := []string{"a: int", "b: int"}
	if len(fn.Parameters) != len(expectedParams) {
		t.Fatalf("Expected %d parameters, got %d", len(expectedParams), len(fn.Parameters))
	}

	for i, expected := range expectedParams {
		if fn.Parameters[i] != expected {
			t.Errorf("Expected parameter '%s', got '%s'", expected, fn.Parameters[i])
		}
	}

	if fn.ReturnType != "int" {
		t.Errorf("Expected return type 'int', got '%s'", fn.ReturnType)
	}
}

func TestPythonParser_Parse_Class(t *testing.T) {
	parser := NewPythonParser()
	content := `class Calculator:
    def __init__(self, name):
        self.name = name
    
    def add(self, a, b):
        return a + b
    
    def multiply(self, a, b):
        if a == 0 or b == 0:
            return 0
        return a * b`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(result.Classes))
	}

	class := result.Classes[0]
	if class.Name != "Calculator" {
		t.Errorf("Expected class name 'Calculator', got '%s'", class.Name)
	}

	if len(class.Methods) != 3 {
		t.Fatalf("Expected 3 methods, got %d", len(class.Methods))
	}

	expectedMethods := []string{"__init__", "add", "multiply"}
	for i, expected := range expectedMethods {
		if class.Methods[i].Name != expected {
			t.Errorf("Expected method '%s', got '%s'", expected, class.Methods[i].Name)
		}
	}

	// Check complexity of multiply method (should be higher due to if statement)
	multiplyMethod := class.Methods[2]
	if multiplyMethod.Complexity < 2 {
		t.Errorf("Expected multiply method complexity >= 2, got %d", multiplyMethod.Complexity)
	}
}

func TestPythonParser_Parse_ClassWithInheritance(t *testing.T) {
	parser := NewPythonParser()
	content := `class AdvancedCalculator(Calculator, Serializable):
    def power(self, base, exponent):
        return base ** exponent`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(result.Classes))
	}

	class := result.Classes[0]
	if class.Name != "AdvancedCalculator" {
		t.Errorf("Expected class name 'AdvancedCalculator', got '%s'", class.Name)
	}

	// Base classes are stored in BaseClasses
	expectedBases := []string{"Calculator", "Serializable"}
	if len(class.BaseClasses) != len(expectedBases) {
		t.Fatalf("Expected %d base classes, got %d", len(expectedBases), len(class.BaseClasses))
	}

	for i, expected := range expectedBases {
		if class.BaseClasses[i] != expected {
			t.Errorf("Expected base class '%s', got '%s'", expected, class.BaseClasses[i])
		}
	}
}

func TestPythonParser_Parse_Imports(t *testing.T) {
	parser := NewPythonParser()
	content := `import os
import sys, json
from collections import defaultdict, Counter
from typing import List, Dict as DictType
import numpy as np`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedImports := []string{
		"os",
		"sys",
		"json",
		"collections.defaultdict",
		"collections.Counter",
		"typing.List",
		"typing.Dict",
		"numpy",
	}

	if len(result.Imports) != len(expectedImports) {
		t.Fatalf("Expected %d imports, got %d", len(expectedImports), len(result.Imports))
	}

	for i, expected := range expectedImports {
		if result.Imports[i] != expected {
			t.Errorf("Expected import '%s', got '%s'", expected, result.Imports[i])
		}
	}
}

func TestPythonParser_Parse_ComplexityCalculation(t *testing.T) {
	parser := NewPythonParser()
	content := `def complex_function(x, y):
    if x > 0:
        if y > 0:
            return x + y
        elif y < 0:
            return x - y
        else:
            return x
    elif x < 0:
        for i in range(abs(x)):
            y += i
        return y
    else:
        while y > 0:
            y -= 1
        return 0`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(result.Functions))
	}

	fn := result.Functions[0]
	// Expected complexity: 1 (base) + 1 (if) + 1 (if) + 1 (elif) + 1 (elif) + 1 (for) + 1 (while) = 7
	expectedComplexity := 7
	if fn.Complexity != expectedComplexity {
		t.Errorf("Expected complexity %d, got %d", expectedComplexity, fn.Complexity)
	}
}

func TestPythonParser_Parse_MixedContent(t *testing.T) {
	parser := NewPythonParser()
	content := `#!/usr/bin/env python3
"""
This is a module docstring.
"""

import os
from typing import Optional

# Global variable
VERSION = "1.0.0"

def standalone_function(name: str) -> str:
    """A standalone function."""
    return f"Hello, {name}!"

class DataProcessor:
    """A class for processing data."""
    
    def __init__(self, config: dict):
        self.config = config
    
    def process(self, data: list) -> Optional[dict]:
        """Process the data."""
        if not data:
            return None
        
        result = {}
        for item in data:
            if item in self.config:
                result[item] = self.config[item]
        
        return result

# Another standalone function
def main():
    processor = DataProcessor({"key": "value"})
    result = processor.process(["key"])
    print(result)`

	result, err := parser.Parse("test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check imports
	expectedImports := []string{"os", "typing.Optional"}
	if len(result.Imports) != len(expectedImports) {
		t.Fatalf("Expected %d imports, got %d", len(expectedImports), len(result.Imports))
	}

	// Check functions
	expectedFunctions := []string{"standalone_function", "main"}
	if len(result.Functions) != len(expectedFunctions) {
		t.Fatalf("Expected %d functions, got %d", len(expectedFunctions), len(result.Functions))
	}

	for i, expected := range expectedFunctions {
		if result.Functions[i].Name != expected {
			t.Errorf("Expected function '%s', got '%s'", expected, result.Functions[i].Name)
		}
	}

	// Check classes
	if len(result.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(result.Classes))
	}

	class := result.Classes[0]
	if class.Name != "DataProcessor" {
		t.Errorf("Expected class name 'DataProcessor', got '%s'", class.Name)
	}

	expectedMethods := []string{"__init__", "process"}
	if len(class.Methods) != len(expectedMethods) {
		t.Fatalf("Expected %d methods, got %d", len(expectedMethods), len(class.Methods))
	}

	// Check line count
	if result.LineCount < 35 || result.LineCount > 45 {
		t.Errorf("Expected line count between 35-45, got %d", result.LineCount)
	}
}

func TestPythonParser_GetSupportedExtensions(t *testing.T) {
	parser := NewPythonParser()
	extensions := parser.GetSupportedExtensions()

	expected := []string{".py", ".pyw"}
	if len(extensions) != len(expected) {
		t.Fatalf("Expected %d extensions, got %d", len(expected), len(extensions))
	}

	for i, exp := range expected {
		if extensions[i] != exp {
			t.Errorf("Expected extension '%s', got '%s'", exp, extensions[i])
		}
	}
}

func TestPythonParser_GetLanguageName(t *testing.T) {
	parser := NewPythonParser()
	name := parser.GetLanguageName()

	if name != "Python" {
		t.Errorf("Expected language name 'Python', got '%s'", name)
	}
}

func TestPythonParser_Parse_AsyncFunctions(t *testing.T) {
	parser := NewPythonParser()
	content := `import asyncio

async def fetch_data(url: str) -> dict:
    async with aiohttp.ClientSession() as session:
        async for attempt in range(3):
            try:
                response = await session.get(url)
                if response.status == 200:
                    return await response.json()
                elif response.status >= 500:
                    continue
                else:
                    break
            except Exception as e:
                if attempt == 2:
                    raise e
                continue
    return {}

class AsyncProcessor:
    async def process(self, items):
        results = []
        async for item in items:
            result = await self.process_item(item)
            results.append(result)
        return results
    
    async def process_item(self, item):
        return item * 2`

	result, err := parser.Parse("async_test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Should find async functions
	if len(result.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(result.Functions))
	}

	asyncFunc := result.Functions[0]
	if !strings.HasPrefix(asyncFunc.Name, "async ") {
		t.Errorf("Expected async function name to start with 'async ', got '%s'", asyncFunc.Name)
	}

	// Should find class with async methods
	if len(result.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(result.Classes))
	}

	class := result.Classes[0]
	if len(class.Methods) != 2 {
		t.Fatalf("Expected 2 methods, got %d", len(class.Methods))
	}

	// Check that complexity is calculated correctly for async constructs
	if asyncFunc.Complexity < 5 {
		t.Errorf("Expected high complexity for async function with loops and conditions, got %d", asyncFunc.Complexity)
	}
}

func TestPythonParser_Parse_Decorators(t *testing.T) {
	parser := NewPythonParser()
	content := `from functools import wraps

@property
def name(self):
    return self._name

@staticmethod
@cache
def utility_function(x, y):
    return x + y

class MyClass:
    @classmethod
    def create(cls, name):
        return cls(name)
    
    @property
    def value(self):
        return self._value
    
    @value.setter
    def value(self, val):
        self._value = val`

	result, err := parser.Parse("decorators_test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Should find functions and methods despite decorators
	if len(result.Functions) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(result.Functions))
	}

	if len(result.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(result.Classes))
	}

	class := result.Classes[0]
	if len(class.Methods) != 3 {
		t.Fatalf("Expected 3 methods, got %d", len(class.Methods))
	}
}

func TestPythonParser_Parse_ComplexExpressions(t *testing.T) {
	parser := NewPythonParser()
	content := `def complex_function(data):
    # List comprehension with condition
    filtered = [x for x in data if x > 0 and x < 100]
    
    # Dictionary comprehension
    mapped = {k: v for k, v in data.items() if v is not None}
    
    # Lambda function
    transform = lambda x: x * 2 if x > 0 else 0
    
    # Complex conditional with multiple operators
    if len(filtered) > 0 and (sum(filtered) > 100 or any(x > 50 for x in filtered)):
        result = [transform(x) for x in filtered if x % 2 == 0]
    elif len(mapped) > 0 and all(v > 0 for v in mapped.values()):
        result = list(mapped.values())
    else:
        result = []
    
    return result`

	result, err := parser.Parse("complex_test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(result.Functions))
	}

	fn := result.Functions[0]
	// Should have high complexity due to multiple conditions, comprehensions, and lambda
	if fn.Complexity < 8 {
		t.Errorf("Expected high complexity for function with comprehensions and complex conditions, got %d", fn.Complexity)
	}
}

func TestPythonParser_Parse_ErrorHandling(t *testing.T) {
	parser := NewPythonParser()

	// Test with empty content
	result, err := parser.Parse("empty.py", []byte(""))
	if err != nil {
		t.Fatalf("Parse failed on empty content: %v", err)
	}

	if result.LineCount != 1 {
		t.Errorf("Expected line count 1 for empty content, got %d", result.LineCount)
	}

	// Test with only comments
	content := `# This is a comment
# Another comment`

	result, err = parser.Parse("comments.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed on comments: %v", err)
	}

	if len(result.Functions) != 0 {
		t.Errorf("Expected 0 functions in comment-only file, got %d", len(result.Functions))
	}

	if len(result.Classes) != 0 {
		t.Errorf("Expected 0 classes in comment-only file, got %d", len(result.Classes))
	}
}

func TestPythonParser_Parse_ExceptionHandling(t *testing.T) {
	parser := NewPythonParser()
	content := `def error_prone_function(data):
    try:
        result = process_data(data)
        if result is None:
            raise ValueError("Invalid result")
        return result
    except ValueError as e:
        logger.error(f"Value error: {e}")
        return None
    except (TypeError, AttributeError) as e:
        logger.error(f"Type/Attribute error: {e}")
        return None
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        raise
    finally:
        cleanup_resources()`

	result, err := parser.Parse("exception_test.py", []byte(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(result.Functions))
	}

	fn := result.Functions[0]
	// Should have complexity for try/except blocks and conditions
	if fn.Complexity < 6 {
		t.Errorf("Expected complexity >= 6 for function with multiple exception handlers, got %d", fn.Complexity)
	}
}

func TestPythonParser_Parse_Docstrings(t *testing.T) {
	parser := NewPythonParser()

	testCases := []struct {
		name          string
		content       string
		hasDocstring  bool
		isClass       bool
		funcOrClsName string
	}{
		{
			name:          "Function with same-line docstring",
			content:       `def my_func(): """This is a docstring."""`,
			hasDocstring:  true,
			isClass:       false,
			funcOrClsName: "my_func",
		},
		{
			name: "Function with next-line docstring",
			content: `def my_func():
	"""This is a docstring."""`,
			hasDocstring:  true,
			isClass:       false,
			funcOrClsName: "my_func",
		},
		{
			name: "Function with docstring after comments and blank lines",
			content: `def my_func():
	# A comment

	"""This is a docstring."""`,
			hasDocstring:  true,
			isClass:       false,
			funcOrClsName: "my_func",
		},
		{
			name:          "Function with no docstring",
			content:       `def my_func():
	pass`,
			hasDocstring:  false,
			isClass:       false,
			funcOrClsName: "my_func",
		},
		{
			name:          "Class with same-line docstring",
			content:       `class MyClass: """This is a docstring."""`,
			hasDocstring:  true,
			isClass:       true,
			funcOrClsName: "MyClass",
		},
		{
			name: "Class with docstring after comments and blank lines",
			content: `class MyClass:
	# A comment

	"""This is a docstring."""`,
			hasDocstring:  true,
			isClass:       true,
			funcOrClsName: "MyClass",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.Parse("test.py", []byte(tc.content))
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if tc.isClass {
				if len(result.Classes) != 1 {
					t.Fatalf("Expected 1 class, got %d", len(result.Classes))
				}
				cls := result.Classes[0]
				if cls.Name != tc.funcOrClsName {
					t.Errorf("Expected class name '%s', got '%s'", tc.funcOrClsName, cls.Name)
				}
				if cls.HasDocstring != tc.hasDocstring {
					t.Errorf("Expected HasDocstring to be %v, but got %v", tc.hasDocstring, cls.HasDocstring)
				}
			} else {
				if len(result.Functions) != 1 {
					t.Fatalf("Expected 1 function, got %d", len(result.Functions))
				}
				fn := result.Functions[0]
				if fn.Name != tc.funcOrClsName {
					t.Errorf("Expected function name '%s', got '%s'", tc.funcOrClsName, fn.Name)
				}
				if fn.HasDocstring != tc.hasDocstring {
					t.Errorf("Expected HasDocstring to be %v, but got %v", tc.hasDocstring, fn.HasDocstring)
				}
			}
		})
	}
}
