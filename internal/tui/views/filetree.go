package views

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/tui"
	"github.com/tito-sala/codebasereaderv2/internal/tui/components"
	"github.com/tito-sala/codebasereaderv2/internal/tui/shared"
)

// FileTreeModel handles file tree navigation
type FileTreeModel struct {
	items       []shared.FileTreeItem
	cursor      int
	selected    map[int]bool
	expanded    map[string]bool
	rootPath    string
	currentPath string
	width       int
	height      int
	scrollY     int
	maxScroll   int
}

// NewFileTreeModel creates a new file tree model
func NewFileTreeModel() *FileTreeModel {
	cwd, _ := os.Getwd()
	return &FileTreeModel{
		items:       []shared.FileTreeItem{},
		cursor:      0,
		selected:    make(map[int]bool),
		expanded:    make(map[string]bool),
		rootPath:    cwd,
		currentPath: cwd,
		scrollY:     0,
	}
}

// Init initializes the file tree
func (m *FileTreeModel) Init() tea.Cmd {
	return m.loadDirectory(m.rootPath)
}

// GetCursor returns the current cursor position (for testing)
func (m *FileTreeModel) GetCursor() int {
	return m.cursor
}

// GetSelected returns the selected items map (for testing)
func (m *FileTreeModel) GetSelected() map[int]bool {
	return m.selected
}

// GetExpanded returns the expanded items map (for testing)
func (m *FileTreeModel) GetExpanded() map[string]bool {
	return m.expanded
}

// GetRootPath returns the root path (for testing)
func (m *FileTreeModel) GetRootPath() string {
	return m.rootPath
}

// GetItems returns the items slice (for testing)
func (m *FileTreeModel) GetItems() []shared.FileTreeItem {
	return m.items
}

// SetItems sets the items slice (for testing)
func (m *FileTreeModel) SetItems(items []shared.FileTreeItem) {
	m.items = items
}

// IsFileSupported checks if a file is supported (for testing)
func (m *FileTreeModel) IsFileSupported(filename string) bool {
	return m.isFileSupported(filename)
}

// GetFileIcon returns the icon for a file item (for testing)
func (m *FileTreeModel) GetFileIcon(item shared.FileTreeItem) string {
	return m.getFileIcon(item)
}

// Update handles messages for the file tree
func (m *FileTreeModel) Update(msg tea.Msg) (*FileTreeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.adjustScroll()
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.adjustScroll()
			}
		case "enter":
			// Navigate into directories or open files
			newM, cmd := m.handleSelection()
			*m = newM
			return m, cmd
		case "right", "l":
			// Expand/collapse folders in-place
			newM, cmd := m.handleExpansion()
			*m = newM
			return m, cmd
		case "left", "h":
			newM, cmd := m.handleCollapse()
			*m = newM
			return m, cmd
		case "backspace":
			// Navigate to parent directory
			newM, cmd := m.navigateToParent()
			*m = newM
			return m, cmd
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "r":
			return m, m.loadDirectory(m.rootPath)
		case "a":
			// Simple analysis: selected items vs global analysis
			if m.hasSelectedItems() {
				// Selected items mode: analyze only the selected directories
				return m, m.analyzeSelectedItems()
			} else {
				// Global mode: analyze entire current directory and all subdirectories
				return m, func() tea.Msg {
					return shared.DirectorySelectedMsg{Path: m.currentPath}
				}
			}

		case "d":
			// Show directory details
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if item.IsDirectory {
					return m, m.showDirectoryInfo(item.Path)
				}
			}

		case "home":
			m.cursor = 0
			m.adjustScroll()

		case "end":
			if len(m.items) > 0 {
				m.cursor = len(m.items) - 1
				m.adjustScroll()
			}

		case "pageup":
			m.cursor = max(0, m.cursor-10)
			m.adjustScroll()

		case "pagedown":
			m.cursor = min(len(m.items)-1, m.cursor+10)
			m.adjustScroll()

		case "g":
			// Go to top
			m.cursor = 0
			m.adjustScroll()

		case "G":
			// Go to bottom
			if len(m.items) > 0 {
				m.cursor = len(m.items) - 1
				m.adjustScroll()
			}
		}

	case shared.DirectoryLoadedMsg:
		m.items = msg.Items
		m.maxScroll = max(0, len(m.items)-m.height+2)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.maxScroll = max(0, len(m.items)-m.height+2)

	case shared.RefreshMsg:
		return m, m.loadDirectory(m.rootPath)

	case shared.StatusUpdateMsg:
		// File tree doesn't handle status updates directly
		return m, nil
	}

	return m, nil
}

// View renders the file tree
func (m *FileTreeModel) View(width, height int) string {
	m.width = width
	m.height = height

	var b strings.Builder

	// Header
	header := fmt.Sprintf("Directory: %s", m.currentPath)
	b.WriteString(components.SelectedStyle.Render(header) + "\n\n")

	if len(m.items) == 0 {
		b.WriteString(components.HelpStyle.Render("No files found or loading..."))
		return b.String()
	}

	// Calculate visible range
	contentHeight := height - 4 // Reserve space for header and controls
	startIdx := m.scrollY
	endIdx := min(len(m.items), startIdx+contentHeight)

	// Render visible items
	for i := startIdx; i < endIdx; i++ {
		item := m.items[i]
		line := m.renderFileTreeItem(item, i == m.cursor, m.selected[i])
		b.WriteString(line + "\n")
	}

	// Item counter showing current selection
	if len(m.items) > 0 {
		itemInfo := fmt.Sprintf("Item %d of %d", m.cursor+1, len(m.items))
		b.WriteString("\n" + components.HelpStyle.Render(itemInfo))
	}

	// Controls help - show different help based on selection state
	var controls string
	if m.hasSelectedItems() {
		selectedCount := 0
		for _, selected := range m.selected {
			if selected {
				selectedCount++
			}
		}
		controls = fmt.Sprintf("Controls: Enter=nav, →=expand, ←=parent, Space=select (%d), a=analyze selected", selectedCount)
	} else {
		controls = "Controls: Enter=navigate, →=expand, ←=parent/collapse, Space=select, a=analyze directory (global)"
	}
	b.WriteString("\n" + components.HelpStyle.Render(controls))

	return b.String()
}

// renderFileTreeItem renders a single file tree item
func (m FileTreeModel) renderFileTreeItem(item shared.FileTreeItem, isCursor, isSelected bool) string {
	var b strings.Builder

	// Indentation with tree lines
	indent := m.renderTreeIndent(item.Level)
	b.WriteString(indent)

	// Expansion indicator for directories
	if item.IsDirectory {
		if m.expanded[item.Path] {
			b.WriteString("📂 ")
		} else {
			b.WriteString("📁 ")
		}
	} else {
		b.WriteString("   ")
	}

	// File/directory icon and name
	icon := m.getFileIcon(item)
	name := item.Name

	// Truncate long names
	maxNameLen := m.width - (item.Level * 2) - 15
	if len(name) > maxNameLen && maxNameLen > 0 {
		name = name[:maxNameLen-3] + "..."
	}

	content := fmt.Sprintf("%s %s", icon, name)

	// Apply styling based on state
	var style lipgloss.Style
	if isCursor && isSelected {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#FF5F87")).
			Bold(true)
	} else if isCursor {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true)
	} else if isSelected {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#3C3C3C"))
	} else if !item.IsSupported && !item.IsDirectory {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
	} else if item.IsDirectory {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB")).
			Bold(true)
	} else {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))
	}

	styledContent := style.Render(content)
	b.WriteString(styledContent)

	// File size and metadata for files
	if !item.IsDirectory && item.Size > 0 {
		sizeStr := tui.FormatFileSize(item.Size)

		// Calculate padding to right-align size
		contentWidth := lipgloss.Width(styledContent)
		availableWidth := m.width - len(indent) - 3
		padding := availableWidth - contentWidth - len(sizeStr)

		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
			sizeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Italic(true)
			b.WriteString(sizeStyle.Render(sizeStr))
		}
	}

	return b.String()
}

// renderTreeIndent renders tree-style indentation
func (m FileTreeModel) renderTreeIndent(level int) string {
	// Validate level to catch unexpected values
	if level < 0 {
		level = 0 // Defensive: ensure no negative indentation
	}

	if level == 0 {
		return ""
	}

	indent := strings.Builder{}
	for i := 0; i < level; i++ {
		if i == level-1 {
			indent.WriteString("├─ ")
		} else {
			indent.WriteString("│  ")
		}
	}
	return indent.String()
}

// getFileIcon returns an appropriate icon for the file type
func (m FileTreeModel) getFileIcon(item shared.FileTreeItem) string {
	if item.IsDirectory {
		return "" // Directory icons are handled separately
	}

	ext := strings.ToLower(filepath.Ext(item.Name))
	switch ext {
	case ".go":
		return "🐹"
	case ".py":
		return "🐍"
	case ".js":
		return "🟨"
	case ".ts":
		return "🔷"
	case ".json":
		return "📋"
	case ".md":
		return "📝"
	case ".txt":
		return "📄"
	case ".yml", ".yaml":
		return "⚙️"
	case ".xml":
		return "📰"
	case ".html", ".htm":
		return "🌐"
	case ".css":
		return "🎨"
	case ".java":
		return "☕"
	case ".c", ".h":
		return "🔧"
	case ".cpp", ".hpp", ".cc":
		return "⚡"
	case ".rs":
		return "🦀"
	case ".php":
		return "🐘"
	case ".rb":
		return "💎"
	case ".sh", ".bash":
		return "🐚"
	case ".sql":
		return "🗃️"
	case ".dockerfile":
		return "🐳"
	case ".gitignore", ".gitattributes":
		return "🔀"
	case ".env":
		return "🔐"
	case ".log":
		return "📊"
	default:
		// Check filename patterns
		name := strings.ToLower(item.Name)
		switch {
		case strings.Contains(name, "readme"):
			return "📖"
		case strings.Contains(name, "license"):
			return "📜"
		case strings.Contains(name, "makefile"):
			return "🔨"
		case strings.Contains(name, "dockerfile"):
			return "🐳"
		case item.IsSupported:
			return "📄"
		default:
			return "❓"
		}
	}
}

// handleSelection handles item selection (navigate into directories or select files)
func (m FileTreeModel) handleSelection() (FileTreeModel, tea.Cmd) {
	if m.cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[m.cursor]

	if item.IsDirectory {
		// Navigate into directory (file manager style)
		return m.navigateToDirectory(item.Path)
	} else {
		// Select file
		return m, func() tea.Msg {
			content, err := os.ReadFile(item.Path)
			if err != nil {
				return shared.ErrorMsg{Error: err}
			}
			return shared.FileSelectedMsg{
				FilePath: item.Path,
				Content:  string(content),
			}
		}
	}
}

// handleExpansion handles expanding/collapsing directories in-place
func (m FileTreeModel) handleExpansion() (FileTreeModel, tea.Cmd) {
	if m.cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[m.cursor]

	if item.IsDirectory {
		// Toggle directory expansion
		if m.expanded[item.Path] {
			return m.collapseDirectory(item.Path)
		} else {
			return m.expandDirectory(item.Path)
		}
	}

	// If it's a file, do nothing
	return m, nil
}

// handleCollapse handles collapsing directories or navigating to parent
func (m FileTreeModel) handleCollapse() (FileTreeModel, tea.Cmd) {
	if m.cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[m.cursor]

	if item.IsDirectory && m.expanded[item.Path] {
		// Collapse current directory
		return m.collapseDirectory(item.Path)
	}

	// If directory is not expanded or not a directory, go to parent
	return m.navigateToParent()
}

// expandDirectory expands a directory and loads its contents
func (m FileTreeModel) expandDirectory(path string) (FileTreeModel, tea.Cmd) {
	m.expanded[path] = true
	return m, m.loadDirectory(m.rootPath) // Reload to show expanded items
}

// collapseDirectory collapses a directory
func (m FileTreeModel) collapseDirectory(path string) (FileTreeModel, tea.Cmd) {
	delete(m.expanded, path)
	return m, m.loadDirectory(m.rootPath) // Reload to hide collapsed items
}

// loadDirectory loads directory contents
func (m FileTreeModel) loadDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		items, err := m.buildFileTree(path, 0)
		if err != nil {
			return shared.ErrorMsg{Error: err}
		}
		return shared.DirectoryLoadedMsg{Items: items}
	}
}

// buildFileTree recursively builds the file tree
func (m FileTreeModel) buildFileTree(path string, level int) ([]shared.FileTreeItem, error) {
	var items []shared.FileTreeItem

	// Add parent directory entry (..) when not at filesystem root and at top level
	if level == 0 {
		parentPath := filepath.Dir(path)
		if parentPath != path { // Not at filesystem root
			parentItem := shared.FileTreeItem{
				Name:        "../",
				Path:        parentPath,
				IsDirectory: true,
				IsSupported: false,
				Level:       0,
				Size:        0,
			}
			items = append(items, parentItem)
		}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Sort entries: directories first, then files, both alphabetically
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		// Skip hidden files unless configured to show them
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		var size int64
		if !entry.IsDir() {
			if info, err := entry.Info(); err == nil {
				size = info.Size()
			}
		}

		item := shared.FileTreeItem{
			Name:        entry.Name(),
			Path:        fullPath,
			IsDirectory: entry.IsDir(),
			IsSupported: m.isFileSupported(entry.Name()),
			Level:       level,
			Size:        size,
		}

		items = append(items, item)

		// Recursively add children if directory is expanded
		if entry.IsDir() && m.expanded[fullPath] {
			children, err := m.buildFileTree(fullPath, level+1)
			if err == nil {
				items = append(items, children...)
			}
		}
	}

	return items, nil
}

// isFileSupported checks if a file type is supported for analysis
func (m FileTreeModel) isFileSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supportedExts := []string{".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".h", ".hpp"}

	for _, supported := range supportedExts {
		if ext == supported {
			return true
		}
	}
	return false
}

// adjustScroll adjusts scroll position to keep cursor visible
func (m *FileTreeModel) adjustScroll() {
	if m.cursor < m.scrollY {
		m.scrollY = m.cursor
	} else if m.cursor >= m.scrollY+m.height-4 {
		m.scrollY = m.cursor - m.height + 5
	}

	if m.scrollY < 0 {
		m.scrollY = 0
	}
	if m.scrollY > m.maxScroll {
		m.scrollY = m.maxScroll
	}
}

// showDirectoryInfo shows information about a directory
func (m FileTreeModel) showDirectoryInfo(path string) tea.Cmd {
	return func() tea.Msg {
		info, err := os.Stat(path)
		if err != nil {
			return shared.StatusUpdateMsg{Message: fmt.Sprintf("Error reading directory: %s", err.Error())}
		}

		// Count files in directory
		entries, err := os.ReadDir(path)
		if err != nil {
			return shared.StatusUpdateMsg{Message: fmt.Sprintf("Error reading directory contents: %s", err.Error())}
		}

		fileCount := 0
		dirCount := 0
		for _, entry := range entries {
			if entry.IsDir() {
				dirCount++
			} else {
				fileCount++
			}
		}

		message := fmt.Sprintf("Directory: %d files, %d subdirs, modified %s",
			fileCount, dirCount, info.ModTime().Format("2006-01-02 15:04"))

		return shared.StatusUpdateMsg{Message: message}
	}
}

// hasSelectedItems returns true if any items are currently selected
func (m *FileTreeModel) hasSelectedItems() bool {
	for _, selected := range m.selected {
		if selected {
			return true
		}
	}
	return false
}

// HasSelectedItems returns true if any items are currently selected (public method)
func (m *FileTreeModel) HasSelectedItems() bool {
	return m.hasSelectedItems()
}

// GetCurrentPath returns the current path
func (m *FileTreeModel) GetCurrentPath() string {
	return m.currentPath
}

// SetCurrentPath sets the current path
func (m *FileTreeModel) SetCurrentPath(path string) {
	m.currentPath = path
}

// LoadDirectory loads directory contents (public method)
func (m *FileTreeModel) LoadDirectory(path string) tea.Cmd {
	return m.loadDirectory(path)
}

// analyzeSelectedItems handles analysis of selected items with appropriate feedback
func (m *FileTreeModel) analyzeSelectedItems() tea.Cmd {
	var selectedDirs []string
	var selectedFiles []string
	selectedCount := 0

	// Collect selected items
	for i, selected := range m.selected {
		if selected && i < len(m.items) {
			item := m.items[i]
			selectedCount++
			if item.IsDirectory {
				selectedDirs = append(selectedDirs, item.Path)
			} else {
				selectedFiles = append(selectedFiles, item.Path)
			}
		}
	}

	if selectedCount == 0 {
		// No actual selections, fall back to current directory
		return func() tea.Msg {
			return shared.StatusUpdateMsg{Message: "No items selected, analyzing current directory"}
		}
	}

	// Determine what to analyze
	if len(selectedDirs) > 0 {
		// Analyze selected directories
		if len(selectedDirs) == 1 {
			return func() tea.Msg {
				return shared.DirectorySelectedMsg{Path: selectedDirs[0]}
			}
		} else {
			// Multiple directories - analyze first one and notify user
			firstDir := selectedDirs[0]
			return func() tea.Msg {
				return shared.DirectorySelectedMsg{Path: firstDir}
			}
		}
	} else if len(selectedFiles) > 0 {
		// Only files selected, analyze their parent directory
		parentDir := filepath.Dir(selectedFiles[0])
		return func() tea.Msg {
			return shared.DirectorySelectedMsg{Path: parentDir}
		}
	}

	// Fallback
	return func() tea.Msg {
		return shared.StatusUpdateMsg{Message: "Unable to determine what to analyze from selection"}
	}
}

// navigateToParent navigates to the parent directory (like cd ..)
func (m FileTreeModel) navigateToParent() (FileTreeModel, tea.Cmd) {
	parentPath := filepath.Dir(m.currentPath)

	// Check if we're already at the filesystem root
	if parentPath == m.currentPath {
		return m, func() tea.Msg {
			return shared.StatusUpdateMsg{Message: "Already at filesystem root"}
		}
	}

	// Navigate to parent directory
	m.currentPath = parentPath
	m.rootPath = parentPath
	m.cursor = 0
	m.scrollY = 0
	m.selected = make(map[int]bool)    // Clear selections
	m.expanded = make(map[string]bool) // Clear expansions

	return m, m.loadDirectory(parentPath)
}

// navigateToDirectory navigates into a directory (like cd directory)
func (m FileTreeModel) navigateToDirectory(directoryPath string) (FileTreeModel, tea.Cmd) {
	// Navigate into the directory
	m.currentPath = directoryPath
	m.rootPath = directoryPath
	m.cursor = 0
	m.scrollY = 0
	m.selected = make(map[int]bool)    // Clear selections
	m.expanded = make(map[string]bool) // Clear expansions

	return m, m.loadDirectory(directoryPath)
}
