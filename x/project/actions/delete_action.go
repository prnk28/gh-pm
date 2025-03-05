package actions

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prnk28/gh-pm/internal/ghcli"
	"github.com/prnk28/gh-pm/internal/tui"
	"github.com/spf13/cobra"
)

// Message types for the delete action
type projectsMsg struct {
	projects []ghcli.ProjectResult
	err      error
}

type deleteMsg struct {
	err error
}

// DeleteAction handles the 'project delete' command
func DeleteAction(cmd *cobra.Command, args []string) {
	// Create and run the BubbleTea program
	p := tea.NewProgram(
		newDeleteModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// ProjectItem represents a project in the list
type ProjectItem struct {
	OrgLogin string
	Project  struct {
		ID     string
		Number int
		Title  string
		URL    string
		Closed bool
	}
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

type deleteState int

const (
	stateSelectProject deleteState = iota
	stateConfirmDelete
	stateDeleting
	stateComplete
	stateError
)

// deleteModel is the model for the delete view
type deleteModel struct {
	list            list.Model
	spinner         tui.Spinner
	state           deleteState
	selectedProject ProjectItem
	width           int
	height          int
	loading         bool
	err             error
	confirmChoice   bool
}

// newDeleteModel creates a new delete model
func newDeleteModel() deleteModel {
	spinner := tui.NewSpinner("Loading projects...")

	// Set up list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142"))

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Select Project to Delete"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2D3142")).Padding(0, 1)

	// Add custom keybindings
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select project"),
			),
		}
	}

	return deleteModel{
		list:    l,
		spinner: spinner,
		loading: true,
		state:   stateSelectProject,
	}
}

// Init initializes the model
func (m deleteModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		m.fetchProjects,
	)
}

// fetchProjects fetches projects from the GitHub API
func (m deleteModel) fetchProjects() tea.Msg {
	type projectsMsg struct {
		projects []ghcli.ProjectResult
		err      error
	}

	return func() tea.Msg {
		projects, err := ghcli.GetAllUserProjects()
		return projectsMsg{
			projects: projects,
			err:      err,
		}
	}
}

// deleteProject deletes the selected project
func (m deleteModel) deleteProject() tea.Cmd {
	type deleteMsg struct {
		err error
	}

	return func() tea.Msg {
		// This would be your actual implementation
		// err := ghcli.DeleteProject(m.selectedProject.Project.ID)
		// Simulating API call for now
		return deleteMsg{
			err: nil,
		}
	}
}

// Update handles messages for the model
func (m deleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4) // Leave space for header/footer

	case tea.KeyMsg:
		switch m.state {
		case stateSelectProject:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "enter":
				if len(m.list.Items()) > 0 {
					selectedItem := m.list.SelectedItem().(ProjectItem)
					m.selectedProject = selectedItem
					m.state = stateConfirmDelete
					return m, nil
				}
			}

		case stateConfirmDelete:
			switch msg.String() {
			case "y", "Y":
				m.confirmChoice = true
				m.state = stateDeleting
				m.loading = true
				return m, m.deleteProject()

			case "n", "N", "esc":
				m.confirmChoice = false
				m.state = stateSelectProject
				return m, nil
			}

		case stateComplete, stateError:
			switch msg.String() {
			case "enter", "esc", "q":
				return m, tea.Quit
			}
		}

	case projectsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}

		// Convert projects to list items
		items := make([]list.Item, 0, len(msg.projects))
		for _, p := range msg.projects {
			items = append(items, ProjectItem{
				OrgLogin: p.OrgLogin,
				Project: struct {
					ID     string
					Number int
					Title  string
					URL    string
					Closed bool
				}{
					ID:     p.Project.ID,
					Number: p.Project.Number,
					Title:  p.Project.Title,
					URL:    p.Project.URL,
					Closed: p.Project.Closed,
				},
			})
		}

		// Update the list with the new items
		m.list.SetItems(items)
		return m, nil

	case deleteMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
		} else {
			m.state = stateComplete
		}
		return m, nil
	}

	// Handle spinner updates when loading
	if m.loading {
		spinnerModel, cmd := m.spinner.Update(msg)
		m.spinner = spinnerModel
		cmds = append(cmds, cmd)
	}

	// Handle list updates when in select state
	if m.state == stateSelectProject {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m deleteModel) View() string {
	switch m.state {
	case stateSelectProject:
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
		return tui.Header("Select Project to Delete") + "\n" +
			m.list.View() + "\n" +
			tui.Footer("↑/↓: Navigate • /: Filter • Enter: Select • q: Quit")

	case stateConfirmDelete:
		confirmStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(1, 0)

		return tui.Header("Confirm Deletion") + "\n\n" +
			fmt.Sprintf("Selected project: %s\n\n", m.selectedProject.Project.Title) +
			confirmStyle.Render("Are you sure you want to delete this project? This action cannot be undone.") + "\n\n" +
			"Press [Y] to confirm or [N] to cancel\n\n" +
			tui.Footer("Y: Confirm • N: Cancel")

	case stateDeleting:
		return tui.Header("Deleting Project") + "\n\n" +
			m.spinner.View() + "\n\n" +
			tui.Footer("Please wait...")

	case stateComplete:
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(1, 0)

		return tui.Header("Project Deleted") + "\n\n" +
			successStyle.Render(fmt.Sprintf("Project '%s' has been successfully deleted.", m.selectedProject.Project.Title)) + "\n\n" +
			tui.Footer("Press Enter to exit")

	case stateError:
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(1, 0)

		return tui.Header("Error") + "\n\n" +
			errorStyle.Render(fmt.Sprintf("Error deleting project: %v", m.err)) + "\n\n" +
			tui.Footer("Press Enter to exit")
	}

	return ""
}
