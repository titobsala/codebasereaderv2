package core

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMainModel(t *testing.T) {
	model := NewMainModel()

	if model == nil {
		t.Fatal("NewMainModel() returned nil")
	}

	if model.currentView != FileTreeView {
		t.Errorf("Expected initial view to be FileTreeView, got %v", model.currentView)
	}

	if model.loading {
		t.Error("Expected initial loading state to be false")
	}

	if model.width != 80 || model.height != 24 {
		t.Errorf("Expected default dimensions 80x24, got %dx%d", model.width, model.height)
	}
}

func TestMainModelUpdate(t *testing.T) {
	model := NewMainModel()

	// Test quit command
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Error("Expected quit command to return a command")
	}

	// Test help toggle - now uses messages
	updatedModel, cmd = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	if cmd == nil {
		t.Error("Expected help command to return a command")
	}

	// Simulate the message being processed
	updatedModel, _ = updatedModel.Update(ShowHelpMsg{Show: true})
	mainModel := updatedModel.(*MainModel)
	if mainModel.currentView != HelpView {
		t.Error("Expected help message to switch to HelpView")
	}

	// Test help toggle again
	updatedModel, _ = updatedModel.Update(ShowHelpMsg{Show: false})
	mainModel = updatedModel.(*MainModel)
	if mainModel.currentView != FileTreeView {
		t.Error("Expected help message to switch back to FileTreeView")
	}
}

func TestMainModelView(t *testing.T) {
	model := NewMainModel()

	// Test that View() doesn't panic and returns a string
	view := model.View()
	if view == "" {
		t.Error("Expected View() to return non-empty string")
	}

	// Test with proper dimensions
	model.width = 80
	model.height = 24
	view = model.View()
	if view == "" {
		t.Error("Expected View() with dimensions to return non-empty string")
	}
}

func TestViewSwitching(t *testing.T) {
	model := NewMainModel()

	originalView := model.currentView
	model.switchView()

	if model.currentView == originalView {
		t.Error("Expected switchView() to change the current view")
	}
}

func TestErrorHandling(t *testing.T) {
	model := NewMainModel()

	// Test setting error
	testErr := ErrorMsg{Error: &testError{"test error"}}
	updatedModel, _ := model.Update(testErr)
	mainModel := updatedModel.(*MainModel)

	if mainModel.error == nil {
		t.Error("Expected error to be set")
	}

	// Test clearing error
	mainModel.ClearError()
	if mainModel.error != nil {
		t.Error("Expected error to be cleared")
	}
}

func TestLoadingState(t *testing.T) {
	model := NewMainModel()

	// Test setting loading state
	loadingMsg := LoadingMsg{Loading: true}
	updatedModel, _ := model.Update(loadingMsg)
	mainModel := updatedModel.(*MainModel)

	if !mainModel.loading {
		t.Error("Expected loading state to be true")
	}

	// Test clearing loading state
	loadingMsg = LoadingMsg{Loading: false}
	updatedModel, _ = updatedModel.Update(loadingMsg)
	mainModel = updatedModel.(*MainModel)

	if mainModel.loading {
		t.Error("Expected loading state to be false")
	}
}

// testError is a simple error implementation for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
