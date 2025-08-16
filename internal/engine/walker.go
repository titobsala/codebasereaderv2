package engine

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tito-sala/codebasereaderv2/internal/parser"
)

// FileWalker handles concurrent directory traversal and file discovery
type FileWalker struct {
	parserRegistry *parser.ParserRegistry
	config         *Config
	gitignoreRules []string
	mutex          sync.RWMutex
}

// NewFileWalker creates a new file walker with the given configuration
func NewFileWalker(parserRegistry *parser.ParserRegistry, config *Config) *FileWalker {
	if config == nil {
		config = DefaultConfig()
	}

	return &FileWalker{
		parserRegistry: parserRegistry,
		config:         config,
		gitignoreRules: make([]string, 0),
	}
}

// WalkResult contains information about a discovered file
type WalkResult struct {
	FilePath string
	Parser   parser.Parser
	Error    error
}

// Walk traverses the directory tree and returns a channel of supported files
func (fw *FileWalker) Walk(rootPath string) (<-chan WalkResult, error) {
	resultChan := make(chan WalkResult, 100)

	// Load .gitignore rules from the root directory
	if err := fw.loadGitignoreRules(rootPath); err != nil {
		// Log warning but continue - .gitignore is optional
		fmt.Printf("Warning: failed to load .gitignore rules: %v\n", err)
	}

	go func() {
		defer close(resultChan)
		
		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				resultChan <- WalkResult{
					FilePath: path,
					Error:    fmt.Errorf("error accessing path %s: %w", path, err),
				}
				return nil // Continue walking
			}

			// Skip directories
			if d.IsDir() {
				// Check if directory should be excluded
				if fw.shouldExcludeDirectory(path, rootPath) {
					return filepath.SkipDir
				}
				return nil
			}

			// Check if file should be excluded
			if fw.shouldExcludeFile(path, rootPath) {
				return nil
			}

			// Check if we have a parser for this file
			parser := fw.getParserForFile(path)
			if parser == nil {
				return nil // Skip unsupported files
			}

			// Check file size limit
			if fw.config.MaxFileSize > 0 {
				info, err := d.Info()
				if err != nil {
					resultChan <- WalkResult{
						FilePath: path,
						Error:    fmt.Errorf("error getting file info for %s: %w", path, err),
					}
					return nil
				}
				
				if info.Size() > fw.config.MaxFileSize {
					return nil // Skip files that are too large
				}
			}

			resultChan <- WalkResult{
				FilePath: path,
				Parser:   parser,
			}

			return nil
		})

		if err != nil {
			resultChan <- WalkResult{
				Error: fmt.Errorf("error walking directory tree: %w", err),
			}
		}
	}()

	return resultChan, nil
}

// loadGitignoreRules loads .gitignore patterns from the root directory
func (fw *FileWalker) loadGitignoreRules(rootPath string) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	gitignorePath := filepath.Join(rootPath, ".gitignore")
	
	file, err := os.Open(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // .gitignore doesn't exist, which is fine
		}
		return err
	}
	defer file.Close()

	fw.gitignoreRules = make([]string, 0)
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		fw.gitignoreRules = append(fw.gitignoreRules, line)
	}

	return scanner.Err()
}

// shouldExcludeDirectory checks if a directory should be excluded from traversal
func (fw *FileWalker) shouldExcludeDirectory(dirPath, rootPath string) bool {
	// Get relative path from root
	relPath, err := filepath.Rel(rootPath, dirPath)
	if err != nil {
		return false
	}

	// Normalize path separators
	relPath = filepath.ToSlash(relPath)
	dirName := filepath.Base(dirPath)

	// Check exclude patterns from config
	for _, pattern := range fw.config.ExcludePatterns {
		if fw.matchesPattern(pattern, relPath) || fw.matchesPattern(pattern, dirName) {
			return true
		}
	}

	// Check .gitignore rules
	fw.mutex.RLock()
	gitignoreRules := fw.gitignoreRules
	fw.mutex.RUnlock()
	
	for _, rule := range gitignoreRules {
		if fw.matchesGitignoreRule(rule, relPath, true) {
			return true
		}
	}

	return false
}

// shouldExcludeFile checks if a file should be excluded from analysis
func (fw *FileWalker) shouldExcludeFile(filePath, rootPath string) bool {
	// Get relative path from root
	relPath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return false
	}

	// Normalize path separators
	relPath = filepath.ToSlash(relPath)
	fileName := filepath.Base(filePath)

	// Check exclude patterns from config
	for _, pattern := range fw.config.ExcludePatterns {
		if fw.matchesPattern(pattern, relPath) || fw.matchesPattern(pattern, fileName) {
			return true
		}
	}

	// Check include patterns if specified
	if len(fw.config.IncludePatterns) > 0 {
		included := false
		for _, pattern := range fw.config.IncludePatterns {
			if fw.matchesPattern(pattern, relPath) || fw.matchesPattern(pattern, fileName) {
				included = true
				break
			}
		}
		if !included {
			return true
		}
	}

	// Check .gitignore rules
	fw.mutex.RLock()
	gitignoreRules := fw.gitignoreRules
	fw.mutex.RUnlock()
	
	for _, rule := range gitignoreRules {
		if fw.matchesGitignoreRule(rule, relPath, false) {
			return true
		}
	}

	return false
}

// getParserForFile returns the appropriate parser for a file, or nil if unsupported
func (fw *FileWalker) getParserForFile(filePath string) parser.Parser {
	parser, err := fw.parserRegistry.GetParser(filePath)
	if err != nil {
		return nil
	}
	return parser
}

// matchesPattern checks if a path matches a glob-like pattern
func (fw *FileWalker) matchesPattern(pattern, path string) bool {
	// Simple pattern matching - can be enhanced with proper glob matching
	pattern = strings.ToLower(pattern)
	path = strings.ToLower(path)

	// Exact match
	if pattern == path {
		return true
	}

	// Wildcard matching
	if strings.Contains(pattern, "*") {
		return fw.matchesWildcard(pattern, path)
	}

	// Check if pattern matches any part of the path (for directory names)
	pathParts := strings.Split(path, "/")
	for _, part := range pathParts {
		if pattern == part {
			return true
		}
	}

	return false
}

// matchesWildcard performs simple wildcard matching
func (fw *FileWalker) matchesWildcard(pattern, str string) bool {
	// Simple implementation - can be enhanced with proper glob matching
	if pattern == "*" {
		return true
	}

	if strings.HasPrefix(pattern, "*.") {
		// File extension pattern
		ext := pattern[2:]
		return strings.HasSuffix(str, "."+ext)
	}

	if strings.HasSuffix(pattern, "*") && strings.HasPrefix(pattern, "*") {
		// Contains pattern (*substring*)
		if len(pattern) > 2 {
			substring := pattern[1 : len(pattern)-1]
			return strings.Contains(str, substring)
		}
		return true
	}

	if strings.HasSuffix(pattern, "*") {
		// Prefix pattern
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(str, prefix)
	}

	if strings.HasPrefix(pattern, "*") {
		// Suffix pattern
		suffix := pattern[1:]
		return strings.HasSuffix(str, suffix)
	}

	return false
}

// matchesGitignoreRule checks if a path matches a .gitignore rule
func (fw *FileWalker) matchesGitignoreRule(rule, path string, isDir bool) bool {
	// Handle negation rules (starting with !)
	if strings.HasPrefix(rule, "!") {
		return false // Negation rules are complex, skip for now
	}

	originalRule := rule
	
	// Handle directory-only rules (ending with /)
	dirOnly := strings.HasSuffix(rule, "/")
	if dirOnly {
		rule = rule[:len(rule)-1]
	}

	// Handle absolute paths (starting with /)
	if strings.HasPrefix(originalRule, "/") {
		rule = rule[1:] // Remove leading slash
		// For absolute paths, match from root
		if dirOnly {
			// For directory rules like /build/, check if path starts with the directory
			return strings.HasPrefix(path, rule+"/") || path == rule
		}
		return fw.matchesPattern(rule, path)
	}

	// For relative patterns
	if dirOnly {
		// Directory rules like dist/ should match dist/* but not files named dist
		pathParts := strings.Split(path, "/")
		for i, part := range pathParts {
			if fw.matchesPattern(rule, part) {
				// Check if this is a directory (has more parts after it)
				if i < len(pathParts)-1 {
					return true
				}
				// If it's the last part, it should only match if isDir is true
				if isDir {
					return true
				}
			}
		}
		return false
	}

	// Handle patterns that should match anywhere in the path
	pathParts := strings.Split(path, "/")
	
	// Check if rule matches any part of the path
	for _, part := range pathParts {
		if fw.matchesPattern(rule, part) {
			return true
		}
	}

	// Check if the rule matches the full path
	if fw.matchesPattern(rule, path) {
		return true
	}

	// Check if rule matches from any directory level
	for i := 0; i < len(pathParts); i++ {
		subPath := strings.Join(pathParts[i:], "/")
		if fw.matchesPattern(rule, subPath) {
			return true
		}
	}

	return false
}

// GetSupportedFiles returns a list of all supported file extensions
func (fw *FileWalker) GetSupportedFiles() []string {
	return fw.parserRegistry.GetSupportedExtensions()
}

// GetStats returns statistics about the file discovery process
func (fw *FileWalker) GetStats(rootPath string) (*WalkStats, error) {
	// Load .gitignore rules first
	if err := fw.loadGitignoreRules(rootPath); err != nil {
		// Log warning but continue - .gitignore is optional
		fmt.Printf("Warning: failed to load .gitignore rules: %v\n", err)
	}

	stats := &WalkStats{
		TotalFiles:      0,
		SupportedFiles:  0,
		ExcludedFiles:   0,
		DirectoriesSkipped: 0,
		FilesByExtension: make(map[string]int),
	}

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Continue walking
		}

		if d.IsDir() {
			if fw.shouldExcludeDirectory(path, rootPath) {
				stats.DirectoriesSkipped++
				return filepath.SkipDir
			}
			return nil
		}

		stats.TotalFiles++

		if fw.shouldExcludeFile(path, rootPath) {
			stats.ExcludedFiles++
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != "" {
			ext = ext[1:] // Remove the dot
			stats.FilesByExtension[ext]++
		}

		if fw.getParserForFile(path) != nil {
			stats.SupportedFiles++
		}

		return nil
	})

	return stats, err
}

// WalkStats contains statistics about the file discovery process
type WalkStats struct {
	TotalFiles         int            `json:"total_files"`
	SupportedFiles     int            `json:"supported_files"`
	ExcludedFiles      int            `json:"excluded_files"`
	DirectoriesSkipped int            `json:"directories_skipped"`
	FilesByExtension   map[string]int `json:"files_by_extension"`
}