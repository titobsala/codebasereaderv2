package parser

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// ParserRegistry manages the registration and lookup of language parsers
type ParserRegistry struct {
	parsers map[string]Parser
	mutex   sync.RWMutex
}

// NewParserRegistry creates a new parser registry
func NewParserRegistry() *ParserRegistry {
	return &ParserRegistry{
		parsers: make(map[string]Parser),
	}
}

// RegisterParser registers a parser for specific file extensions
func (r *ParserRegistry) RegisterParser(parser Parser) error {
	if parser == nil {
		return fmt.Errorf("parser cannot be nil")
	}

	extensions := parser.GetSupportedExtensions()
	if len(extensions) == 0 {
		return fmt.Errorf("parser must support at least one file extension")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, ext := range extensions {
		// Normalize extension (ensure it starts with a dot)
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		ext = strings.ToLower(ext)
		r.parsers[ext] = parser
	}

	return nil
}

// GetParser returns the appropriate parser for a given file path
func (r *ParserRegistry) GetParser(filePath string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	parser, exists := r.parsers[ext]
	if !exists {
		return nil, fmt.Errorf("no parser registered for extension: %s", ext)
	}

	return parser, nil
}

// GetSupportedExtensions returns all supported file extensions
func (r *ParserRegistry) GetSupportedExtensions() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	extensions := make([]string, 0, len(r.parsers))
	for ext := range r.parsers {
		extensions = append(extensions, ext)
	}

	return extensions
}

// IsSupported checks if a file extension is supported
func (r *ParserRegistry) IsSupported(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.parsers[ext]
	return exists
}

// GetRegisteredParsers returns information about all registered parsers
func (r *ParserRegistry) GetRegisteredParsers() map[string]string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]string)
	processed := make(map[string]bool)

	for _, parser := range r.parsers {
		langName := parser.GetLanguageName()
		if !processed[langName] {
			var extensions []string
			for e, p := range r.parsers {
				if p.GetLanguageName() == langName {
					extensions = append(extensions, e)
				}
			}
			result[langName] = strings.Join(extensions, ", ")
			processed[langName] = true
		}
	}

	return result
}
