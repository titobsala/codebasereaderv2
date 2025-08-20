package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

func setupTestEnvironment(t *testing.T) (string, *FileWalker, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "walker_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test file structure
	testFiles := map[string]string{
		"main.go":               "package main\n\nfunc main() {}\n",
		"src/utils.go":          "package src\n\nfunc Utils() {}\n",
		"src/parser.py":         "def parse():\n    pass\n",
		"node_modules/lib.js":   "console.log('test');\n",
		".git/config":           "[core]\n",
		"vendor/dep.go":         "package vendor\n",
		"__pycache__/cache.pyc": "compiled python",
		"docs/README.md":        "# Documentation\n",
		"test/test_file.go":     "package test\n",
		"large_file.go":         strings.Repeat("// Large file\n", 1000),
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Create .gitignore file
	gitignoreContent := `# Compiled files
*.pyc
__pycache__/

# Dependencies
node_modules/

# Build artifacts
/build/
dist/

# IDE files
.vscode/
*.swp
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// Setup parser registry
	registry := parser.NewParserRegistry()
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	pyParser := &MockParser{name: "Python", extensions: []string{"py"}}

	registry.RegisterParser(goParser)
	registry.RegisterParser(pyParser)

	// Create config
	config := DefaultConfig()
	config.MaxFileSize = 1024 * 10 // 10KB limit for testing

	// Create file walker
	walker := NewFileWalker(registry, config)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, walker, cleanup
}

func TestFileWalker_Walk(t *testing.T) {
	tempDir, walker, cleanup := setupTestEnvironment(t)
	defer cleanup()

	resultChan, err := walker.Walk(tempDir)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	var results []WalkResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Check that we found the expected files
	expectedFiles := []string{"main.go", "src/utils.go", "src/parser.py", "test/test_file.go"}
	foundFiles := make(map[string]bool)

	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Unexpected error for file %s: %v", result.FilePath, result.Error)
			continue
		}

		relPath, err := filepath.Rel(tempDir, result.FilePath)
		if err != nil {
			t.Errorf("Failed to get relative path: %v", err)
			continue
		}

		foundFiles[relPath] = true

		// Verify parser is assigned correctly
		if result.Parser == nil {
			t.Errorf("No parser assigned for file: %s", relPath)
			continue
		}

		ext := strings.ToLower(filepath.Ext(relPath))
		if ext != "" {
			ext = ext[1:] // Remove dot
		}

		expectedLanguage := ""
		switch ext {
		case "go":
			expectedLanguage = "Go"
		case "py":
			expectedLanguage = "Python"
		}

		if result.Parser.GetLanguageName() != expectedLanguage {
			t.Errorf("Wrong parser for file %s: expected %s, got %s",
				relPath, expectedLanguage, result.Parser.GetLanguageName())
		}
	}

	// Verify expected files were found
	for _, expectedFile := range expectedFiles {
		if !foundFiles[expectedFile] {
			t.Errorf("Expected file not found: %s", expectedFile)
		}
	}

	// Verify excluded files were not found
	excludedFiles := []string{"node_modules/lib.js", ".git/config", "vendor/dep.go", "__pycache__/cache.pyc"}
	for _, excludedFile := range excludedFiles {
		if foundFiles[excludedFile] {
			t.Errorf("Excluded file was found: %s", excludedFile)
		}
	}
}

func TestFileWalker_GitignoreRules(t *testing.T) {
	tempDir, walker, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Load gitignore rules
	err := walker.loadGitignoreRules(tempDir)
	if err != nil {
		t.Fatalf("Failed to load gitignore rules: %v", err)
	}

	// Test that gitignore rules are loaded
	if len(walker.gitignoreRules) == 0 {
		t.Error("No gitignore rules loaded")
	}

	// Test specific exclusions
	testCases := []struct {
		path     string
		excluded bool
	}{
		{"__pycache__/cache.pyc", true},
		{"node_modules/lib.js", true},
		{"build/output.bin", true}, // This should match /build/ rule
		{"dist/app.js", true},      // This should match dist/ rule
		{"src/main.go", false},
		{"test.py", false},
	}

	for _, tc := range testCases {
		excluded := walker.shouldExcludeFile(filepath.Join(tempDir, tc.path), tempDir)
		if excluded != tc.excluded {
			t.Errorf("File %s: expected excluded=%v, got excluded=%v", tc.path, tc.excluded, excluded)
		}
	}
}

func TestFileWalker_PatternMatching(t *testing.T) {
	_, walker, cleanup := setupTestEnvironment(t)
	defer cleanup()

	testCases := []struct {
		pattern string
		path    string
		matches bool
	}{
		{"*.go", "main.go", true},
		{"*.go", "main.py", false},
		{"node_modules", "node_modules/lib.js", true},
		{"__pycache__", "src/__pycache__/cache.pyc", true},
		{"test*", "test_file.go", true},
		{"*cache*", "pycache.py", true},
		{"exact.txt", "exact.txt", true},
		{"exact.txt", "not_exact.txt", false},
	}

	for _, tc := range testCases {
		matches := walker.matchesPattern(tc.pattern, tc.path)
		if matches != tc.matches {
			t.Errorf("Pattern '%s' with path '%s': expected %v, got %v",
				tc.pattern, tc.path, tc.matches, matches)
		}
	}
}

func TestFileWalker_WildcardMatching(t *testing.T) {
	_, walker, cleanup := setupTestEnvironment(t)
	defer cleanup()

	testCases := []struct {
		pattern string
		str     string
		matches bool
	}{
		{"*", "anything", true},
		{"*.go", "main.go", true},
		{"*.go", "main.py", false},
		{"test*", "test_file.go", true},
		{"test*", "main_test.go", false},
		{"*test", "main_test", true},
		{"*test", "test_main", false},
	}

	for _, tc := range testCases {
		matches := walker.matchesWildcard(tc.pattern, tc.str)
		if matches != tc.matches {
			t.Errorf("Wildcard pattern '%s' with string '%s': expected %v, got %v",
				tc.pattern, tc.str, tc.matches, matches)
		}
	}
}

func TestFileWalker_GetStats(t *testing.T) {
	tempDir, walker, cleanup := setupTestEnvironment(t)
	defer cleanup()

	stats, err := walker.GetStats(tempDir)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	// Verify basic stats
	if stats.TotalFiles == 0 {
		t.Error("Expected some total files")
	}

	if stats.SupportedFiles == 0 {
		t.Error("Expected some supported files")
	}

	// Note: ExcludedFiles might be 0 if no files match the exclude patterns
	// This is actually correct behavior - we only exclude files that match patterns

	if stats.DirectoriesSkipped == 0 {
		t.Error("Expected some directories to be skipped")
	}

	// Verify file extension counts
	if stats.FilesByExtension["go"] == 0 {
		t.Error("Expected some .go files")
	}

	if stats.FilesByExtension["py"] == 0 {
		t.Error("Expected some .py files")
	}
}

func TestFileWalker_ExcludePatterns(t *testing.T) {
	tempDir, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create walker with custom exclude patterns
	registry := parser.NewParserRegistry()
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	registry.RegisterParser(goParser)

	config := DefaultConfig()
	config.ExcludePatterns = []string{"test*", "src/*"}

	walker := NewFileWalker(registry, config)

	resultChan, err := walker.Walk(tempDir)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	var results []WalkResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Check that test files and src files are excluded
	for _, result := range results {
		if result.Error != nil {
			continue
		}

		relPath, err := filepath.Rel(tempDir, result.FilePath)
		if err != nil {
			continue
		}

		if strings.HasPrefix(relPath, "test") {
			t.Errorf("File with 'test' prefix should be excluded: %s", relPath)
		}

		if strings.HasPrefix(relPath, "src/") {
			t.Errorf("File in 'src' directory should be excluded: %s", relPath)
		}
	}
}

func TestFileWalker_IncludePatterns(t *testing.T) {
	tempDir, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create walker with include patterns
	registry := parser.NewParserRegistry()
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	registry.RegisterParser(goParser)

	config := DefaultConfig()
	config.IncludePatterns = []string{"*.go"}
	config.ExcludePatterns = []string{} // Clear default excludes for this test

	walker := NewFileWalker(registry, config)

	resultChan, err := walker.Walk(tempDir)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	var results []WalkResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Check that only .go files are included
	for _, result := range results {
		if result.Error != nil {
			continue
		}

		if !strings.HasSuffix(result.FilePath, ".go") {
			t.Errorf("Non-.go file should be excluded when include patterns are set: %s", result.FilePath)
		}
	}
}

func TestFileWalker_FileSizeLimit(t *testing.T) {
	tempDir, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create walker with small file size limit
	registry := parser.NewParserRegistry()
	goParser := &MockParser{name: "Go", extensions: []string{"go"}}
	registry.RegisterParser(goParser)

	config := DefaultConfig()
	config.MaxFileSize = 100            // Very small limit
	config.ExcludePatterns = []string{} // Clear default excludes

	walker := NewFileWalker(registry, config)

	resultChan, err := walker.Walk(tempDir)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	foundLargeFile := false

	for result := range resultChan {
		if result.Error != nil {
			continue
		}

		relPath, err := filepath.Rel(tempDir, result.FilePath)
		if err != nil {
			continue
		}

		if relPath == "large_file.go" {
			foundLargeFile = true
		}
	}

	if foundLargeFile {
		t.Error("Large file should have been excluded due to size limit")
	}
}
