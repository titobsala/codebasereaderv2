package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/codebasereaderv2/internal/tui/shared"
)

// TabItem represents a single tab
type TabItem struct {
	Title   string
	Content string
	Icon    string
}

// TabsModel handles tab navigation and display
type TabsModel struct {
	Tabs      []TabItem
	activeTab int
	width     int
	height    int
}

// NewTabsModel creates a new tabs model with default tabs
func NewTabsModel() *TabsModel {
	tabs := []TabItem{
		{Title: "Explorer", Icon: "📁", Content: "File tree and content view"},
		{Title: "Analysis", Icon: "📊", Content: "Metrics and analysis results"},
		{Title: "Configuration", Icon: "⚙️", Content: "Settings and preferences"},
		{Title: "Help", Icon: "❓", Content: "Help and documentation"},
	}

	return &TabsModel{
		Tabs:      tabs,
		activeTab: 0,
	}
}

// Init initializes the tabs model
func (m *TabsModel) Init() tea.Cmd {
	return nil
}

// Update handles tab navigation
func (m *TabsModel) Update(msg tea.Msg) (*TabsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "l", "tab":
			m.activeTab = (m.activeTab + 1) % len(m.Tabs)
		case "left", "h", "shift+tab":
			m.activeTab = (m.activeTab - 1 + len(m.Tabs)) % len(m.Tabs)
		case "1":
			m.activeTab = 0
		case "2":
			if len(m.Tabs) > 1 {
				m.activeTab = 1
			}
		case "3":
			if len(m.Tabs) > 2 {
				m.activeTab = 2
			}
		case "4":
			if len(m.Tabs) > 3 {
				m.activeTab = 3
			}
		}
	case tea.WindowSizeMsg:
		// We don't need to handle window size here anymore for the main width,
		// as the parent model will pass the correct width to the View function.
		m.height = msg.Height
	}
	return m, nil
}

// View renders the tabs, now accepting a width for the tab column
func (m *TabsModel) View(width int) string {
	m.width = width
	return m.renderTabs()
}

// GetActiveTab returns the currently active tab index
func (m *TabsModel) GetActiveTab() int {
	return m.activeTab
}

// SetActiveTab sets the active tab
func (m *TabsModel) SetActiveTab(index int) {
	if index >= 0 && index < len(m.Tabs) {
		m.activeTab = index
	}
}

// renderTabs renders the tab bar as a horizontal row with advanced styling
func (m *TabsModel) renderTabs() string {
	var renderedTabs []string

	tabWidth := (m.width - len(m.Tabs) - 1) / len(m.Tabs) // Account for spacing between tabs
	if tabWidth < 12 {
		tabWidth = 12 // Minimum tab width for better readability
	}

	for i, tab := range m.Tabs {
		var style lipgloss.Style
		isActive := i == m.activeTab

		if isActive {
			// Enhanced active tab styling with gradient-like effect and borders
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(NeutralWhite).
				Background(PrimaryPurple).
				Padding(1, 2).
				Align(lipgloss.Center).
				Border(lipgloss.Border{
					Bottom: "▂",
					Left:   "▏",
					Right:  "▕",
				}).
				BorderForeground(PrimaryBlue).
				BorderBackground(PrimaryPurple)
		} else {
			// Enhanced inactive tab styling with subtle hover effects
			style = lipgloss.NewStyle().
				Foreground(NeutralMedium).
				Padding(1, 2).
				Align(lipgloss.Center).
				Border(lipgloss.Border{
					Bottom: "─",
				}).
				BorderForeground(NeutralDark).
				Underline(false)

			// Add subtle background for better visual separation
			if i == m.activeTab-1 || i == m.activeTab+1 {
				// Adjacent tabs get slightly different styling for better UX
				style = style.Background(lipgloss.Color("#2A2A3A"))
			}
		}

		// Enhanced tab text with better spacing and icons
		tabText := tab.Icon + " " + tab.Title

		// Apply additional styling for active tab to make it stand out
		if isActive {
			tabText = "🔥 " + tabText + " 🔥" // Add fire emojis for active tab
		}

		renderedTabs = append(renderedTabs, style.Width(tabWidth).Render(tabText))
	}

	// Join tabs with better spacing and add subtle background
	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Add a subtle background container for the entire tab bar
	tabBarContainer := lipgloss.NewStyle().
		Background(lipgloss.Color("#1A1A2A")).
		Padding(0, 1).
		Render(tabsRow)

	return tabBarContainer
}

// GetTabTitle returns the title of the active tab
func (m *TabsModel) GetTabTitle() string {
	if m.activeTab >= 0 && m.activeTab < len(m.Tabs) {
		return m.Tabs[m.activeTab].Title
	}
	return "Unknown"
}

// MapTabToViewType maps tab index to ViewType
func (m *TabsModel) MapTabToViewType() shared.ViewType {
	switch m.activeTab {
	case 0:
		return shared.FileTreeView // Explorer
	case 1:
		return shared.ContentView // Analysis (will show metrics when analysis is available)
	case 2:
		return shared.ConfigView // Configuration
	case 3:
		return shared.HelpView // Help
	default:
		return shared.FileTreeView
	}
}
