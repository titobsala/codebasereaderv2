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
		FilePath:   filePath,
		Language:   "Go",
		Functions:  []FunctionInfo{},
		Classes:    []ClassInfo{},
		Imports:    []string{},
		Errors:     []ParseError{},
		AnalyzedAt: time.Now(),
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

	// Extract imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		result.Imports = append(result.Imports, importPath)
	}

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
		Name:       funcDecl.Name.Name,
		LineStart:  startPos.Line,
		LineEnd:    endPos.Line,
		Parameters: []string{},
		ReturnType: "",
		Complexity: g.calculateComplexity(funcDecl),
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
		Name:      typeSpec.Name.Name,
		LineStart: startPos.Line,
		LineEnd:   endPos.Line,
		Methods:   []FunctionInfo{},
		Fields:    []string{},
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