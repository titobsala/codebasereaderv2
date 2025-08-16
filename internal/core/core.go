package core

import (
	"fmt"

	"github.com/tito-sala/codebasereaderv2/internal/ai"
	"github.com/tito-sala/codebasereaderv2/internal/engine"
	"github.com/tito-sala/codebasereaderv2/internal/parser"
	"github.com/tito-sala/codebasereaderv2/internal/tui"
)

// Application represents the main application structure
type Application struct {
	engine    *engine.Engine
	aiClient  ai.AIClient
	tuiConfig *tui.TUIConfig
}

// NewApplication creates a new application instance
func NewApplication(config *engine.Config) *Application {
	if config == nil {
		config = engine.DefaultConfig()
	}

	app := &Application{
		engine:    engine.NewEngine(config),
		tuiConfig: tui.DefaultTUIConfig(),
	}

	return app
}

// RegisterParser registers a new language parser with the application
func (app *Application) RegisterParser(parser parser.Parser) error {
	return app.engine.GetParserRegistry().RegisterParser(parser)
}

// SetAIClient sets the AI client for the application
func (app *Application) SetAIClient(client ai.AIClient) {
	app.aiClient = client
}

// GetEngine returns the analysis engine
func (app *Application) GetEngine() *engine.Engine {
	return app.engine
}

// GetAIClient returns the AI client
func (app *Application) GetAIClient() ai.AIClient {
	return app.aiClient
}

// GetTUIConfig returns the TUI configuration
func (app *Application) GetTUIConfig() *tui.TUIConfig {
	return app.tuiConfig
}

// ValidateSetup validates that the application is properly configured
func (app *Application) ValidateSetup() error {
	// Check if at least one parser is registered
	extensions := app.engine.GetParserRegistry().GetSupportedExtensions()
	if len(extensions) == 0 {
		return fmt.Errorf("no parsers registered")
	}

	// Validate configuration
	config := app.engine.GetConfig()
	if config.MaxWorkers <= 0 {
		return fmt.Errorf("max workers must be greater than 0")
	}

	if config.MaxFileSize <= 0 {
		return fmt.Errorf("max file size must be greater than 0")
	}

	return nil
}

// GetSupportedLanguages returns information about supported languages
func (app *Application) GetSupportedLanguages() map[string]string {
	return app.engine.GetParserRegistry().GetRegisteredParsers()
}