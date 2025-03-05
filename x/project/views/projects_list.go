package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prnk28/gh-pm/internal/ghc"
	"github.com/prnk28/gh-pm/internal/models"
	"github.com/prnk28/gh-pm/internal/tui"
)

// Message types for the projects list view
type projectsMsg struct {
	projects []models.ProjectsJson
	err      error
}

// ProjectItem represents a project in the list
type ProjectItem struct {
	OrgLogin string
	Project  models.ProjectsJson
}

// Title returns the title for the list item
func (i ProjectItem) Title() string {
	return fmt.Sprintf("#%d: %s", i.Project.Number, i.Project.Title)
}

// Description returns the description for the list item
func (i ProjectItem) Description() string {
	status := "Open"
	if i.Project.Closed {
		status = "Closed"
	}
	return fmt.Sprintf("Organization: %s • Status: %s", i.OrgLogin, status)
}

// FilterValue returns the value to use for filtering
func (i ProjectItem) FilterValue() string {
	return fmt.Sprintf("%s %s", i.OrgLogin, i.Project.Title)
}

// ProjectsListViewModel is the model for the projects list view
type ProjectsListViewModel struct {
	list    list.Model
	spinner tui.Spinner
	loading bool
	err     error
	width   int
	height  int
}

// NewProjectsListViewModel creates a new projects list view model
func NewProjectsListViewModel() ProjectsListViewModel {
	spinner := tui.NewSpinner("Loading projects...")

	// Set up list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142"))

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "GitHub Projects"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142")).Padding(0, 1)

	return ProjectsListViewModel{
		list:    l,
		spinner: spinner,
		loading: true,
	}
}

// Init initializes the model
func (m ProjectsListViewModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		m.fetchProjects,
	)
}

// fetchProjects fetches projects from the GitHub API
func (m ProjectsListViewModel) fetchProjects() tea.Msg {
	return func() tea.Msg {
		projects, err := ghc.GetProjects()
		return projectsMsg{
			projects: projects,
			err:      err,
		}
	}
}

// Update handles messages for the model
func (m ProjectsListViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4) // Leave space for header/footer

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case projectsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		// Convert projects to list items
		items := make([]list.Item, 0, len(msg.projects))
		for _, p := range msg.projects {
			items = append(items, ProjectItem{
				OrgLogin: p.Owner.Login,
				Project:  p,
			})
		}

		// Update the list with the new items
		m.list.SetItems(items)
		return m, nil
	}

	// Handle spinner updates when loading
	if m.loading {
		spinnerModel, cmd := m.spinner.Update(msg)
		m.spinner = spinnerModel
		cmds = append(cmds, cmd)
	}

	// Handle list updates
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m ProjectsListViewModel) View() string {
	if m.err != nil {
		return tui.Header("Error") + "\n\n" +
			fmt.Sprintf("Error fetching projects: %v", m.err) + "\n\n" +
			tui.Footer("Press q to quit")
	}

	if m.loading {
		return tui.Header("GitHub Projects") + "\n\n" +
			m.spinner.View() + "\n\n" +
			tui.Footer("Press q to quit")
	}

	// Render the list with header and footer
	return strings.Join([]string{
		tui.Header("GitHub Projects"),
		m.list.View(),
		tui.Footer("↑/↓: Navigate • /: Filter • Enter: Select • q: Quit"),
	}, "\n")
}
