package parser

import (
	"regexp"
	"strings"
	"time"
)

// PythonParser implements the Parser interface for Python language files
type PythonParser struct {
	// Regex patterns for Python constructs
	classPattern         *regexp.Regexp
	functionPattern      *regexp.Regexp
	asyncFunctionPattern *regexp.Regexp
	methodPattern        *regexp.Regexp
	importPattern        *regexp.Regexp
	fromImportPattern    *regexp.Regexp
	decoratorPattern     *regexp.Regexp
}

// NewPythonParser creates a new Python parser instance
func NewPythonParser() *PythonParser {
	return &PythonParser{
		classPattern:         regexp.MustCompile(`^(\s*)class\s+(\w+)(?:\(([^)]*)\))?:`),
		functionPattern:      regexp.MustCompile(`^(\s*)def\s+(\w+)\s*\(([^)]*)\)(?:\s*->\s*([^:]+))?:`),
		asyncFunctionPattern: regexp.MustCompile(`^(\s*)async\s+def\s+(\w+)\s*\(([^)]*)\)(?:\s*->\s*([^:]+))?:`),
		methodPattern:        regexp.MustCompile(`^(\s+)def\s+(\w+)\s*\(([^)]*)\)(?:\s*->\s*([^:]+))?:`),
		importPattern:        regexp.MustCompile(`^import\s+(.+)$`),
		fromImportPattern:    regexp.MustCompile(`^from\s+(\S+)\s+import\s+(.+)$`),
		decoratorPattern:     regexp.MustCompile(`^(\s*)@(\w+(?:\.\w+)*)`),
	}
}

// Parse analyzes Python source code and returns structured results
func (p *PythonParser) Parse(filePath string, content []byte) (*AnalysisResult, error) {
	result := &AnalysisResult{
		FilePath:     filePath,
		Language:     "Python",
		Functions:    []FunctionInfo{},
		Classes:      []ClassInfo{},
		Imports:      []string{},
		Dependencies: []Dependency{},
		Errors:       []ParseError{},
		AnalyzedAt:   time.Now(),
	}

	lines := strings.Split(string(content), "\n")
	result.LineCount = len(lines)

	// Track current context for nested structures
	var currentClass *ClassInfo
	var decorators []string

	for lineNum, line := range lines {
		lineNum++ // Convert to 1-based indexing
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// Handle decorators
		if matches := p.decoratorPattern.FindStringSubmatch(line); matches != nil {
			decorators = append(decorators, matches[2])
			continue
		}

		// Handle imports
		if matches := p.importPattern.FindStringSubmatch(trimmedLine); matches != nil {
			imports := p.parseImports(matches[1])
			result.Imports = append(result.Imports, imports...)
			continue
		}

		if matches := p.fromImportPattern.FindStringSubmatch(trimmedLine); matches != nil {
			module := matches[1]
			imports := p.parseImports(matches[2])
			for _, imp := range imports {
				result.Imports = append(result.Imports, module+"."+imp)
			}
			continue
		}

		// Calculate indentation level
		indent := p.getIndentLevel(line)

		// Handle class definitions
		if matches := p.classPattern.FindStringSubmatch(line); matches != nil {
			// Close previous class if we're at a lower or equal indentation level
			if currentClass != nil {
				prevClassIndent := 0
				if currentClass.LineStart > 1 {
					prevClassIndent = p.getIndentLevel(lines[currentClass.LineStart-2])
				}
				if indent <= prevClassIndent {
					currentClass.LineEnd = lineNum - 1
					result.Classes = append(result.Classes, *currentClass)
					currentClass = nil
				}
			}

			className := matches[2]
			baseClasses := p.parseBaseClasses(matches[3])

			currentClass = &ClassInfo{
				Name:         className,
				LineStart:    lineNum,
				LineEnd:      lineNum, // Will be updated when class ends
				Methods:      []FunctionInfo{},
				Fields:       []string{}, // Will be populated separately
				BaseClasses:  baseClasses,
				IsPublic:     p.isPublicClass(className),
				HasDocstring: p.hasDocstring(lines, lineNum),
			}

			// Clear decorators after use
			decorators = []string{}
			continue
		}

		// Handle function/method definitions (regular and async)
		var funcMatches []string
		var isAsync bool

		if matches := p.functionPattern.FindStringSubmatch(line); matches != nil {
			funcMatches = matches
			isAsync = false
		} else if matches := p.asyncFunctionPattern.FindStringSubmatch(line); matches != nil {
			funcMatches = matches
			isAsync = true
		}

		if funcMatches != nil {
			funcName := funcMatches[2]
			if isAsync {
				funcName = "async " + funcName
			}
			params := p.parseParameters(funcMatches[3])
			returnType := strings.TrimSpace(funcMatches[4])

			funcInfo := FunctionInfo{
				Name:                 funcName,
				LineStart:            lineNum,
				LineEnd:              lineNum, // Will be updated when function ends
				Parameters:           params,
				ReturnType:           returnType,
				Complexity:           p.calculateComplexity(lines, lineNum-1, indent),
				CyclomaticComplexity: p.calculateComplexity(lines, lineNum-1, indent),
				ParameterCount:       len(params),
				IsPublic:             p.isPublicFunction(funcMatches[2]),
				IsAsync:              isAsync,
				HasDocstring:         p.hasDocstring(lines, lineNum),
			}

			// Determine if this is a method or standalone function
			if currentClass != nil {
				classIndent := 0
				if currentClass.LineStart > 1 {
					classIndent = p.getIndentLevel(lines[currentClass.LineStart-2])
				}
				if indent > classIndent {
					// This is a method
					currentClass.Methods = append(currentClass.Methods, funcInfo)
				} else {
					// This is a standalone function
					result.Functions = append(result.Functions, funcInfo)
				}
			} else {
				// This is a standalone function
				result.Functions = append(result.Functions, funcInfo)
			}

			// Clear decorators after use
			decorators = []string{}
		}
	}

	// Close any remaining open class
	if currentClass != nil {
		currentClass.LineEnd = len(lines)
		currentClass.LinesOfCode = currentClass.LineEnd - currentClass.LineStart + 1
		currentClass.MethodCount = len(currentClass.Methods)
		currentClass.FieldCount = len(currentClass.Fields)

		// Calculate class complexity as sum of method complexities
		classComplexity := 0
		for _, method := range currentClass.Methods {
			classComplexity += method.Complexity
		}
		currentClass.Complexity = classComplexity

		result.Classes = append(result.Classes, *currentClass)
	}

	// Calculate total complexity
	for _, fn := range result.Functions {
		result.Complexity += fn.Complexity
	}
	for _, class := range result.Classes {
		for _, method := range class.Methods {
			result.Complexity += method.Complexity
		}
	}

	return result, nil
}

// getIndentLevel calculates the indentation level of a line
func (p *PythonParser) getIndentLevel(line string) int {
	indent := 0
	for _, char := range line {
		if char == ' ' {
			indent++
		} else if char == '\t' {
			indent += 4 // Treat tab as 4 spaces
		} else {
			break
		}
	}
	return indent
}

// parseImports parses import statements and returns a list of imported modules
func (p *PythonParser) parseImports(importStr string) []string {
	var imports []string
	parts := strings.Split(importStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Handle "as" aliases
		if strings.Contains(part, " as ") {
			parts := strings.Split(part, " as ")
			imports = append(imports, strings.TrimSpace(parts[0]))
		} else {
			imports = append(imports, part)
		}
	}

	return imports
}

// parseBaseClasses parses the base classes from a class definition
func (p *PythonParser) parseBaseClasses(baseStr string) []string {
	if baseStr == "" {
		return []string{}
	}

	var bases []string
	parts := strings.Split(baseStr, ",")

	for _, part := range parts {
		base := strings.TrimSpace(part)
		if base != "" {
			bases = append(bases, base)
		}
	}

	return bases
}

// parseParameters parses function parameters
func (p *PythonParser) parseParameters(paramStr string) []string {
	if paramStr == "" {
		return []string{}
	}

	var params []string
	parts := strings.Split(paramStr, ",")

	for _, part := range parts {
		param := strings.TrimSpace(part)
		if param != "" {
			// Handle default values
			if strings.Contains(param, "=") {
				paramParts := strings.Split(param, "=")
				param = strings.TrimSpace(paramParts[0])
			}
			params = append(params, param)
		}
	}

	return params
}

// calculateComplexity calculates cyclomatic complexity for a Python function
func (p *PythonParser) calculateComplexity(lines []string, startLine, baseIndent int) int {
	complexity := 1 // Base complexity

	// Look ahead to find the end of the function
	for i := startLine + 1; i < len(lines); i++ {
		line := lines[i]
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// If we've reached a line with equal or less indentation, we're done
		currentIndent := p.getIndentLevel(line)
		if currentIndent <= baseIndent && trimmedLine != "" {
			break
		}

		// Count complexity-increasing constructs
		if strings.HasPrefix(trimmedLine, "if ") ||
			strings.HasPrefix(trimmedLine, "elif ") ||
			strings.HasPrefix(trimmedLine, "for ") ||
			strings.HasPrefix(trimmedLine, "while ") ||
			strings.HasPrefix(trimmedLine, "except ") ||
			strings.HasPrefix(trimmedLine, "with ") ||
			strings.HasPrefix(trimmedLine, "try:") ||
			strings.HasPrefix(trimmedLine, "finally:") ||
			strings.HasPrefix(trimmedLine, "async for ") ||
			strings.HasPrefix(trimmedLine, "async with ") {
			complexity++
		}

		// Count logical operators that increase complexity
		if strings.Contains(trimmedLine, " and ") ||
			strings.Contains(trimmedLine, " or ") {
			// Count the number of logical operators
			andCount := strings.Count(trimmedLine, " and ")
			orCount := strings.Count(trimmedLine, " or ")
			complexity += andCount + orCount
		}

		// Count list/dict comprehensions
		if strings.Contains(trimmedLine, " for ") &&
			(strings.Contains(trimmedLine, "[") || strings.Contains(trimmedLine, "{")) {
			complexity++
		}

		// Count lambda functions
		if strings.Contains(trimmedLine, "lambda ") {
			complexity++
		}
	}

	return complexity
}

// GetSupportedExtensions returns the file extensions supported by this parser
func (p *PythonParser) GetSupportedExtensions() []string {
	return []string{".py", ".pyw"}
}

// GetLanguageName returns the human-readable language name
func (p *PythonParser) GetLanguageName() string {
	return "Python"
}

// isPublicFunction determines if a function is public in Python
func (p *PythonParser) isPublicFunction(name string) bool {
	return !strings.HasPrefix(name, "_")
}

// isPublicClass determines if a class is public in Python
func (p *PythonParser) isPublicClass(name string) bool {
	return !strings.HasPrefix(name, "_")
}

// hasDocstring checks if a function or class has a docstring.
// It checks the definition line itself, and subsequent lines.
func (p *PythonParser) hasDocstring(lines []string, lineNum int) bool {
	// lineNum is 1-based, so the definition line is at index lineNum - 1
	defLineIdx := lineNum - 1

	if defLineIdx >= len(lines) {
		return false
	}

	// First, check if the docstring is on the same line as the definition
	defLine := lines[defLineIdx]
	if strings.Contains(defLine, `"""`) || strings.Contains(defLine, `'''`) {
		// More robust check to ensure it's not just a string literal in code
		trimmedDefLine := strings.TrimSpace(defLine)
		if strings.HasSuffix(trimmedDefLine, `"""`) || strings.HasSuffix(trimmedDefLine, `'''`) {
			return true
		}
	}

	// If not on the same line, look for docstring in the next lines
	for i := lineNum; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// If we find a non-empty, non-comment line that is not a docstring, stop.
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, `"""`) && !strings.HasPrefix(line, `'''`) {
			return false
		}

		if strings.HasPrefix(line, `"""`) || strings.HasPrefix(line, `'''`) {
			return true
		}
	}

	return false
}
