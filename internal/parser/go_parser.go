package parser

import (
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"strings"
	"time"
)

// GoParser implements the Parser interface for Go language files
type GoParser struct{}

// NewGoParser creates a new Go parser instance
func NewGoParser() *GoParser {
	return &GoParser{}
}

// Parse analyzes Go source code and returns structured results
func (g *GoParser) Parse(filePath string, content []byte) (*AnalysisResult, error) {
	result := &AnalysisResult{
		FilePath:     filePath,
		Language:     "Go",
		Functions:    []FunctionInfo{},
		Classes:      []ClassInfo{},
		Imports:      []string{},
		Dependencies: []Dependency{},
		Errors:       []ParseError{},
		AnalyzedAt:   time.Now(),
	}

	// Create a new token file set
	fset := token.NewFileSet()

	// Parse the Go source code
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		// Handle parsing errors
		if list, ok := err.(scanner.ErrorList); ok {
			for _, e := range list {
				result.Errors = append(result.Errors, ParseError{
					Line:    e.Pos.Line,
					Column:  e.Pos.Column,
					Message: e.Msg,
				})
			}
		} else {
			result.Errors = append(result.Errors, ParseError{
				Line:    0,
				Column:  0,
				Message: err.Error(),
			})
		}
		// Return partial results even if there are parsing errors
		return result, nil
	}

	// Count lines
	result.LineCount = strings.Count(string(content), "\n") + 1

	// Extract and categorize imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		result.Imports = append(result.Imports, importPath)

		// Create dependency with proper categorization
		dependency := g.categorizeImport(importPath, filePath)
		result.Dependencies = append(result.Dependencies, dependency)
	}

	result.ImportCount = len(result.Imports)

	// Walk the AST to extract information
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcInfo := g.extractFunctionInfo(fset, x)
			result.Functions = append(result.Functions, funcInfo)
			result.Complexity += funcInfo.Complexity

		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, spec := range x.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							classInfo := g.extractStructInfo(fset, typeSpec, structType)
							result.Classes = append(result.Classes, classInfo)
						} else if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							classInfo := g.extractInterfaceInfo(fset, typeSpec, interfaceType)
							result.Classes = append(result.Classes, classInfo)
						}
					}
				}
			}
		}
		return true
	})

	return result, nil
}

// extractFunctionInfo extracts information about a function declaration
func (g *GoParser) extractFunctionInfo(fset *token.FileSet, funcDecl *ast.FuncDecl) FunctionInfo {
	startPos := fset.Position(funcDecl.Pos())
	endPos := fset.Position(funcDecl.End())

	funcInfo := FunctionInfo{
		Name:                 funcDecl.Name.Name,
		LineStart:            startPos.Line,
		LineEnd:              endPos.Line,
		Parameters:           []string{},
		ReturnType:           "",
		Complexity:           g.calculateComplexity(funcDecl),
		CyclomaticComplexity: g.calculateComplexity(funcDecl),
		LinesOfCode:          endPos.Line - startPos.Line + 1,
		IsPublic:             g.isPublicFunction(funcDecl.Name.Name),
		IsAsync:              false, // Go doesn't have async functions like Python
		HasDocstring:         funcDecl.Doc != nil && len(funcDecl.Doc.List) > 0,
	}

	// Extract parameters
	if funcDecl.Type.Params != nil {
		for _, param := range funcDecl.Type.Params.List {
			paramType := g.typeToString(param.Type)
			if len(param.Names) > 0 {
				for _, name := range param.Names {
					funcInfo.Parameters = append(funcInfo.Parameters, name.Name+" "+paramType)
				}
			} else {
				funcInfo.Parameters = append(funcInfo.Parameters, paramType)
			}
		}
	}

	funcInfo.ParameterCount = len(funcInfo.Parameters)

	// Extract return type
	if funcDecl.Type.Results != nil {
		var returnTypes []string
		for _, result := range funcDecl.Type.Results.List {
			returnTypes = append(returnTypes, g.typeToString(result.Type))
		}
		if len(returnTypes) == 1 {
			funcInfo.ReturnType = returnTypes[0]
		} else if len(returnTypes) > 1 {
			funcInfo.ReturnType = "(" + strings.Join(returnTypes, ", ") + ")"
		}
	}

	return funcInfo
}

// extractStructInfo extracts information about a struct declaration
func (g *GoParser) extractStructInfo(fset *token.FileSet, typeSpec *ast.TypeSpec, structType *ast.StructType) ClassInfo {
	startPos := fset.Position(structType.Pos())
	endPos := fset.Position(structType.End())

	classInfo := ClassInfo{
		Name:         typeSpec.Name.Name,
		LineStart:    startPos.Line,
		LineEnd:      endPos.Line,
		Methods:      []FunctionInfo{},
		Fields:       []string{},
		LinesOfCode:  endPos.Line - startPos.Line + 1,
		IsPublic:     g.isPublicType(typeSpec.Name.Name),
		BaseClasses:  []string{}, // Go doesn't have inheritance
		HasDocstring: typeSpec.Doc != nil && len(typeSpec.Doc.List) > 0,
	}

	// Extract fields
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			fieldType := g.typeToString(field.Type)
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					classInfo.Fields = append(classInfo.Fields, name.Name+" "+fieldType)
				}
			} else {
				// Embedded field
				classInfo.Fields = append(classInfo.Fields, fieldType)
			}
		}
	}

	classInfo.FieldCount = len(classInfo.Fields)
	return classInfo
}

// extractInterfaceInfo extracts information about an interface declaration
func (g *GoParser) extractInterfaceInfo(fset *token.FileSet, typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) ClassInfo {
	startPos := fset.Position(interfaceType.Pos())
	endPos := fset.Position(interfaceType.End())

	classInfo := ClassInfo{
		Name:      typeSpec.Name.Name,
		LineStart: startPos.Line,
		LineEnd:   endPos.Line,
		Methods:   []FunctionInfo{},
		Fields:    []string{},
	}

	// Extract methods from interface
	if interfaceType.Methods != nil {
		for _, method := range interfaceType.Methods.List {
			if len(method.Names) > 0 {
				for _, name := range method.Names {
					if funcType, ok := method.Type.(*ast.FuncType); ok {
						methodInfo := FunctionInfo{
							Name:       name.Name,
							LineStart:  fset.Position(method.Pos()).Line,
							LineEnd:    fset.Position(method.End()).Line,
							Parameters: []string{},
							ReturnType: "",
							Complexity: 1, // Interface methods have minimal complexity
						}

						// Extract parameters
						if funcType.Params != nil {
							for _, param := range funcType.Params.List {
								paramType := g.typeToString(param.Type)
								if len(param.Names) > 0 {
									for _, paramName := range param.Names {
										methodInfo.Parameters = append(methodInfo.Parameters, paramName.Name+" "+paramType)
									}
								} else {
									methodInfo.Parameters = append(methodInfo.Parameters, paramType)
								}
							}
						}

						// Extract return type
						if funcType.Results != nil {
							var returnTypes []string
							for _, result := range funcType.Results.List {
								returnTypes = append(returnTypes, g.typeToString(result.Type))
							}
							if len(returnTypes) == 1 {
								methodInfo.ReturnType = returnTypes[0]
							} else if len(returnTypes) > 1 {
								methodInfo.ReturnType = "(" + strings.Join(returnTypes, ", ") + ")"
							}
						}

						classInfo.Methods = append(classInfo.Methods, methodInfo)
					}
				}
			}
		}
	}

	return classInfo
}

// calculateComplexity calculates cyclomatic complexity for a function
func (g *GoParser) calculateComplexity(funcDecl *ast.FuncDecl) int {
	complexity := 1 // Base complexity

	if funcDecl.Body != nil {
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			switch n.(type) {
			case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
				complexity++
			case *ast.CaseClause:
				complexity++
			}
			return true
		})
	}

	return complexity
}

// typeToString converts an AST type to its string representation
func (g *GoParser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + g.typeToString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + g.typeToString(t.Elt)
		}
		return "[...]" + g.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + g.typeToString(t.Key) + "]" + g.typeToString(t.Value)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + g.typeToString(t.Value)
		case ast.RECV:
			return "<-chan " + g.typeToString(t.Value)
		default:
			return "chan " + g.typeToString(t.Value)
		}
	case *ast.FuncType:
		return "func"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.SelectorExpr:
		return g.typeToString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

// GetSupportedExtensions returns the file extensions supported by this parser
func (g *GoParser) GetSupportedExtensions() []string {
	return []string{".go"}
}

// GetLanguageName returns the human-readable language name
func (g *GoParser) GetLanguageName() string {
	return "Go"
}

// isPublicFunction determines if a function is public (exported) in Go
func (g *GoParser) isPublicFunction(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}

// isPublicType determines if a type is public (exported) in Go
func (g *GoParser) isPublicType(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}

// categorizeImport categorizes an import path into internal, external, or standard library
func (g *GoParser) categorizeImport(importPath, filePath string) Dependency {
	dependency := Dependency{
		Name:        importPath,
		UsageCount:  1, // For now, assume each import is used once
		IsDirectDep: true,
		FilePath:    filePath,
	}

	// Categorize the import
	if g.isStandardLibrary(importPath) {
		dependency.Type = "standard"
	} else if g.isInternalImport(importPath, filePath) {
		dependency.Type = "internal"
	} else {
		dependency.Type = "external"
		// Try to extract version from module path if available
		dependency.Version = g.extractVersion(importPath)
	}

	return dependency
}

// isStandardLibrary checks if an import is from Go's standard library
func (g *GoParser) isStandardLibrary(importPath string) bool {
	// Go standard library packages don't contain dots or slashes in their root
	// and are well-known packages
	standardPackages := map[string]bool{
		// Core packages
		"fmt": true, "os": true, "io": true, "strings": true, "strconv": true,
		"time": true, "context": true, "errors": true, "sync": true, "sort": true,
		"encoding/json": true, "encoding/xml": true, "encoding/base64": true,
		"net/http": true, "net/url": true, "net": true,
		"path": true, "path/filepath": true,
		"regexp": true, "reflect": true, "runtime": true,
		"bufio": true, "bytes": true, "compress/gzip": true,
		"crypto/md5": true, "crypto/sha1": true, "crypto/sha256": true, "crypto/rand": true,
		"database/sql": true, "html/template": true, "text/template": true,
		"log": true, "math": true, "math/rand": true,
		// Testing
		"testing": true, "testing/quick": true,
		// AST and parsing
		"go/ast": true, "go/parser": true, "go/scanner": true, "go/token": true,
		"go/format": true, "go/build": true,
	}

	// Check direct match
	if standardPackages[importPath] {
		return true
	}

	// Check if it's a sub-package of a known standard package
	// Standard library packages don't contain domain names (no dots before first slash)
	if !strings.Contains(importPath, ".") && !strings.Contains(importPath, "/") {
		return true
	}

	// Check common standard library prefixes
	standardPrefixes := []string{
		"archive/", "bufio", "builtin", "bytes", "compress/", "container/",
		"context", "crypto/", "database/", "debug/", "embed", "encoding/",
		"errors", "expvar", "flag", "fmt", "go/", "hash/", "html/", "image/",
		"index/", "io/", "log/", "math/", "mime/", "net/", "os/", "path/",
		"plugin", "reflect", "regexp/", "runtime/", "sort", "strconv",
		"strings", "sync/", "syscall", "testing/", "text/", "time",
		"unicode/", "unsafe",
	}

	for _, prefix := range standardPrefixes {
		if strings.HasPrefix(importPath, prefix) {
			return true
		}
	}

	return false
}

// isInternalImport checks if an import is internal to the current project
func (g *GoParser) isInternalImport(importPath, filePath string) bool {
	// Internal imports typically:
	// 1. Start with a relative path (./  ../)
	// 2. Share the same module prefix as the current file

	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		return true
	}

	// Try to determine if it's from the same module
	// This is a heuristic - we look for common patterns in the import path
	// vs the file path structure

	// If the import path contains the same domain/organization as the project,
	// it's likely internal. For now, we'll use a simple heuristic:
	// if it contains known patterns like github.com/org/project or
	// matches the current project structure

	// Extract potential module name from file path
	// This is simplified - in a full implementation, we'd parse go.mod
	if strings.Contains(importPath, "/internal/") {
		return true
	}

	// For this project specifically, anything starting with our module path
	if strings.HasPrefix(importPath, "github.com/tito-sala/codebasereaderv2") {
		return true
	}

	// Simple heuristic: if import doesn't contain a domain, might be internal
	if !strings.Contains(importPath, ".") && strings.Contains(importPath, "/") {
		return true
	}

	return false
}

// extractVersion attempts to extract version information from import path
func (g *GoParser) extractVersion(importPath string) string {
	// For now, we don't have access to go.mod parsing
	// In a full implementation, we'd parse go.mod or use go list
	// This is a placeholder for future enhancement
	return ""
}
