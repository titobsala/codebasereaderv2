package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewFileTreeModel(t *testing.T) {
	model := NewFileTreeModel()
	
	if model.cursor != 0 {
		t.Error("Expected initial cursor position to be 0")
	}
	
	if model.selected == nil {
		t.Error("Expected selected map to be initialized")
	}
	
	if model.expanded == nil {
		t.Error("Expected expanded map to be initialized")
	}
	
	if model.rootPath == "" {
		t.Error("Expected rootPath to be set")
	}
}

func TestFileTreeNavigation(t *testing.T) {
	model := NewFileTreeModel()
	
	// Add some test items
	model.items = []FileTreeItem{
		{Name: "file1.go", Path: "/test/file1.go", IsDirectory: false, IsSupported: true},
		{Name: "file2.py", Path: "/test/file2.py", IsDirectory: false, IsSupported: true},
		{Name: "dir1", Path: "/test/dir1", IsDirectory: true, IsSupported: false},
	}
	
	// Test down navigation
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if updatedModel.cursor != 1 {
		t.Errorf("Expected cursor to be 1 after down, got %d", updatedModel.cursor)
	}
	
	// Test up navigation
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if updatedModel.cursor != 0 {
		t.Errorf("Expected cursor to be 0 after up, got %d", updatedModel.cursor)
	}
}

func TestFileSelection(t *testing.T) {
	model := NewFileTreeModel()
	
	// Add a test file
	model.items = []FileTreeItem{
		{Name: "test.go", Path: "/test/test.go", IsDirectory: false, IsSupported: true},
	}
	
	// Test space key for selection - create a proper KeyMsg that will return "space" from String()
	keyMsg := tea.KeyMsg{Type: tea.KeySpace, Runes: []rune(" ")}
	updatedModel, _ := model.Update(keyMsg)
	
	if !updatedModel.selected[0] {
		t.Error("Expected file to be selected after space key")
	}
	
	// Test space key again to deselect
	updatedModel, _ = updatedModel.Update(keyMsg)
	
	if updatedModel.selected[0] {
		t.Error("Expected file to be deselected after second space key")
	}
}

func TestIsFileSupported(t *testing.T) {
	model := NewFileTreeModel()
	
	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.go", true},
		{"test.py", true},
		{"test.js", true},
		{"test.ts", true},
		{"test.txt", false},
		{"test.md", false},
		{"test", false},
	}
	
	for _, tc := range testCases {
		result := model.isFileSupported(tc.filename)
		if result != tc.expected {
			t.Errorf("isFileSupported(%s) = %v, expected %v", tc.filename, result, tc.expected)
		}
	}
}

func TestFormatFileSize(t *testing.T) {
	testCases := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}
	
	for _, tc := range testCases {
		result := formatFileSize(tc.size)
		if result != tc.expected {
			t.Errorf("formatFileSize(%d) = %s, expected %s", tc.size, result, tc.expected)
		}
	}
}

func TestFileTreeView(t *testing.T) {
	model := NewFileTreeModel()
	
	// Add some test items
	model.items = []FileTreeItem{
		{Name: "test.go", Path: "/test/test.go", IsDirectory: false, IsSupported: true, Size: 1024},
	}
	
	// Test that View doesn't panic
	view := model.View(80, 24)
	if view == "" {
		t.Error("Expected View() to return non-empty string")
	}
}

func TestGetFileIcon(t *testing.T) {
	model := NewFileTreeModel()
	
	testCases := []struct {
		item     FileTreeItem
		expected string
	}{
		{FileTreeItem{Name: "dir", IsDirectory: true}, ""}, // Directory icons are handled separately
		{FileTreeItem{Name: "test.go", IsDirectory: false, IsSupported: true}, "🐹"},
		{FileTreeItem{Name: "test.py", IsDirectory: false, IsSupported: true}, "🐍"},
		{FileTreeItem{Name: "test.js", IsDirectory: false, IsSupported: true}, "🟨"}, // Updated icon
		{FileTreeItem{Name: "test.md", IsDirectory: false, IsSupported: false}, "📝"},
	}
	
	for _, tc := range testCases {
		result := model.getFileIcon(tc.item)
		if result != tc.expected {
			t.Errorf("getFileIcon(%s) = %s, expected %s", tc.item.Name, result, tc.expected)
		}
	}
}