package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prnk28/gh-pm/internal/ghcli"
	"github.com/prnk28/gh-pm/internal/tui"
)

// Message types for the project view
type projectDetailsMsg struct {
	project ghcli.ProjectResult
	columns []Column
	err     error
}

// ProjectsView represents the view for a project and its cards
type ProjectsView struct {
	project         ghcli.ProjectResult
	columns         []Column
	activeColumn    int
	spinner         tui.Spinner
	loading         bool
	err             error
	width           int
	height          int
	ready           bool
	columnViewports []viewport.Model
	table           table.Model
	tabs            []string
	activeTab       int

	// When true, bypasses loading and uses sample data
	UseSampleData bool
}

// Column represents a column in the kanban board
type Column struct {
	Name  string
	Cards []Card
}

// Card represents a card in the kanban board
type Card struct {
	ID     string
	Title  string
	Status string
	URL    string
}

// NewProjectsView creates a new projects view,
// with UseSampleData set to true if you want to bypass loading.
func NewProjectsView(projectID string, useSample bool) ProjectsView {
	spinner := tui.NewSpinner("Loading project details...")

	// Set up table for project details
	columns := []table.Column{
		{Title: "Property", Width: 20},
		{Title: "Value", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#2D3142")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2D3142")).
		Bold(true)
	t.SetStyles(s)

	// Create the model with UseSampleData flag
	return ProjectsView{
		spinner:       spinner,
		loading:       true,
		table:         t,
		tabs:          []string{"Details", "Kanban Board"},
		activeTab:     0,
		UseSampleData: useSample,
	}
}

// Init initializes the model. If UseSampleData is true, it returns sample data right away.
func (m ProjectsView) Init() tea.Cmd {
	if m.UseSampleData {
		// Immediately return sample data without waiting for spinner
		return m.fetchSampleProjectDetails
	}

	// Otherwise, run the spinner and then fetch actual details (or sample data in real API call)
	return tea.Batch(
		m.spinner.Init(),
		m.fetchProjectDetails,
	)
}

// fetchSampleProjectDetails returns sample data for testing
func (m ProjectsView) fetchSampleProjectDetails() tea.Msg {
	// Sample project data
	project := ghcli.ProjectResult{
		OrgLogin: "example-org",
		Project: struct {
			ID        string
			Number    int
			Title     string
			URL       string
			Closed    bool
			CreatedAt string
		}{
			ID:        "proj123",
			Number:    42,
			Title:     "Example Project",
			URL:       "https://github.com/orgs/example-org/projects/42",
			Closed:    false,
			CreatedAt: "2023-01-01T00:00:00Z",
		},
	}

	columns := []Column{
		{
			Name: "To Do",
			Cards: []Card{
				{ID: "card1", Title: "Implement feature A", Status: "Open", URL: "https://github.com/example-org/repo/issues/1"},
				{ID: "card2", Title: "Fix bug in module B", Status: "Open", URL: "https://github.com/example-org/repo/issues/2"},
			},
		},
		{
			Name: "In Progress",
			Cards: []Card{
				{ID: "card3", Title: "Refactor authentication", Status: "In Progress", URL: "https://github.com/example-org/repo/issues/3"},
			},
		},
		{
			Name: "Done",
			Cards: []Card{
				{ID: "card4", Title: "Update documentation", Status: "Closed", URL: "https://github.com/example-org/repo/issues/4"},
				{ID: "card5", Title: "Release version 1.0", Status: "Closed", URL: "https://github.com/example-org/repo/issues/5"},
			},
		},
	}

	return projectDetailsMsg{
		project: project,
		columns: columns,
		err:     nil,
	}
}

// fetchProjectDetails fetches project details from the GitHub API.
// When UseSampleData is true, you can simply call m.fetchSampleProjectDetails.
func (m ProjectsView) fetchProjectDetails() tea.Msg {
	// In a real implementation, you'd fetch data from GitHub here.
	// For testing, you can use the sample data function:
	return m.fetchSampleProjectDetails()
}

// Update handles messages for the model
func (m ProjectsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize column viewports
			m.columnViewports = make([]viewport.Model, len(m.columns))
			columnWidth := (m.width - 4) / max(1, len(m.columns))

			for i := range m.columns {
				m.columnViewports[i] = viewport.New(columnWidth, m.height-10)
				m.columnViewports[i].SetContent(m.renderColumn(m.columns[i]))
			}

			m.table.SetWidth(m.width - 4)
			m.table.SetHeight(10)

			m.ready = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			// Switch between details and kanban views
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			return m, nil

		case "left", "h":
			if m.activeTab == 1 { // Only in kanban view
				m.activeColumn = max(0, m.activeColumn-1)
			}

		case "right", "l":
			if m.activeTab == 1 { // Only in kanban view
				m.activeColumn = min(len(m.columns)-1, m.activeColumn+1)
			}

		case "up", "k":
			if m.activeTab == 0 { // Details view
				var cmd tea.Cmd
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.activeTab == 1 && len(m.columnViewports) > 0 { // Kanban view
				var cmd tea.Cmd
				m.columnViewports[m.activeColumn], cmd = m.columnViewports[m.activeColumn].Update(msg)
				cmds = append(cmds, cmd)
			}

		case "down", "j":
			if m.activeTab == 0 { // Details view
				var cmd tea.Cmd
				m.table, cmd = m.table.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.activeTab == 1 && len(m.columnViewports) > 0 { // Kanban view
				var cmd tea.Cmd
				m.columnViewports[m.activeColumn], cmd = m.columnViewports[m.activeColumn].Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case projectDetailsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.project = msg.project
		m.columns = msg.columns

		// Update table with project details
		rows := []table.Row{
			{"ID", m.project.Project.ID},
			{"Number", fmt.Sprintf("%d", m.project.Project.Number)},
			{"Title", m.project.Project.Title},
			{"Organization", m.project.OrgLogin},
			{"Status", m.getStatusText()},
			{"URL", m.project.Project.URL},
			{"Created At", m.project.Project.CreatedAt},
		}
		m.table.SetRows(rows)

		// Initialize column viewports if we have window dimensions
		if m.width > 0 && m.height > 0 {
			m.columnViewports = make([]viewport.Model, len(m.columns))
			columnWidth := (m.width - 4) / max(1, len(m.columns))
			for i := range m.columns {
				m.columnViewports[i] = viewport.New(columnWidth, m.height-10)
				m.columnViewports[i].SetContent(m.renderColumn(m.columns[i]))
			}
			m.ready = true
		}

		return m, nil
	}

	// Handle spinner updates when loading
	if m.loading {
		spinnerModel, cmd := m.spinner.Update(msg)
		m.spinner = spinnerModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m ProjectsView) View() string {
	if m.err != nil {
		return tui.Header("Error") + "\n\n" +
			fmt.Sprintf("Error fetching project details: %v", m.err) + "\n\n" +
			tui.Footer("Press q to quit")
	}

	if m.loading {
		return tui.Header("Loading Project") + "\n\n" +
			m.spinner.View() + "\n\n" +
			tui.Footer("Press q to quit")
	}

	// Render tabs
	tabsView := m.renderTabs()

	var content string
	if m.activeTab == 0 {
		// Details view
		content = m.table.View()
	} else {
		// Kanban view
		content = m.renderKanban()
	}

	return tui.Header(fmt.Sprintf("Project: %s", m.project.Project.Title)) + "\n" +
		tabsView + "\n" +
		content + "\n" +
		m.renderFooter()
}

// renderTabs renders the tab navigation
func (m ProjectsView) renderTabs() string {
	var b strings.Builder

	tabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#2D3142"))

	activeTabStyle := tabStyle.Copy().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2D3142"))

	for i, tab := range m.tabs {
		if i == m.activeTab {
			b.WriteString(activeTabStyle.Render(tab))
		} else {
			b.WriteString(tabStyle.Render(tab))
		}
	}

	return b.String()
}

// renderKanban renders the kanban board view
func (m ProjectsView) renderKanban() string {
	if !m.ready || len(m.columnViewports) == 0 {
		return "Loading kanban board..."
	}

	var b strings.Builder

	// Render column headers
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#2D3142")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width((m.width - 4) / len(m.columns))

	activeHeaderStyle := headerStyle.Copy().
		Background(lipgloss.Color("#4F5D75"))

	for i, column := range m.columns {
		if i == m.activeColumn {
			b.WriteString(activeHeaderStyle.Render(fmt.Sprintf("%s (%d)", column.Name, len(column.Cards))))
		} else {
			b.WriteString(headerStyle.Render(fmt.Sprintf("%s (%d)", column.Name, len(column.Cards))))
		}
	}

	b.WriteString("\n")

	// Render column contents
	for i, vp := range m.columnViewports {
		if i == m.activeColumn {
			b.WriteString(lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#4F5D75")).
				Render(vp.View()))
		} else {
			b.WriteString(lipgloss.NewStyle().
				BorderStyle(lipgloss.HiddenBorder()).
				Render(vp.View()))
		}
	}

	return b.String()
}

// renderColumn renders a single column of the kanban board
func (m ProjectsView) renderColumn(column Column) string {
	var b strings.Builder

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2D3142")).
		Padding(1, 2).
		Margin(0, 0, 1, 0).
		Width((m.width - 10) / len(m.columns))

	for _, card := range column.Cards {
		statusStyle := lipgloss.NewStyle().
			Foreground(m.getStatusColor(card.Status))

		b.WriteString(cardStyle.Render(
			fmt.Sprintf("%s\n%s\n",
				card.Title,
				statusStyle.Render(card.Status),
			),
		))
		b.WriteString("\n")
	}

	return b.String()
}

// renderFooter renders the footer with help text
func (m ProjectsView) renderFooter() string {
	var helpText string

	if m.activeTab == 0 {
		helpText = "↑/↓: Navigate • Tab: Switch View • q: Quit"
	} else {
		helpText = "←/→: Change Column • ↑/↓: Scroll Cards • Tab: Switch View • q: Quit"
	}

	return tui.Footer(helpText)
}

// getStatusText returns a formatted status text
func (m ProjectsView) getStatusText() string {
	if m.project.Project.Closed {
		return "Closed"
	}
	return "Open"
}

// getStatusColor returns a color based on status
func (m ProjectsView) getStatusColor(status string) lipgloss.Color {
	switch status {
	case "Open":
		return lipgloss.Color("#00FF00") // Green
	case "In Progress":
		return lipgloss.Color("#FFFF00") // Yellow
	case "Closed":
		return lipgloss.Color("#FF0000") // Red
	default:
		return lipgloss.Color("#FFFFFF") // White
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
