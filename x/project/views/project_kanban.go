package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prnk28/gh-pm/internal/ghcli"
	"github.com/prnk28/gh-pm/internal/models"
	"github.com/prnk28/gh-pm/internal/tui"
)

// Message types for the project view
type fetchCardsMsg struct {
	columns []models.ProjectColumn
	err     error
}

// KanbanView represents the view for a project's kanban board
type KanbanView struct {
	projectID       string
	columns         []models.ProjectColumn
	activeColumn    int
	spinner         tui.Spinner
	loading         bool
	err             error
	width           int
	height          int
	ready           bool
	columnViewports []viewport.Model
}

// NewKanbanView creates a new kanban view for a project
func NewKanbanView(projectID string) KanbanView {
	spinner := tui.NewSpinner("Loading project cards...")

	return KanbanView{
		projectID:    projectID,
		spinner:      spinner,
		loading:      true,
		activeColumn: 0,
	}
}

// Init initializes the model
func (m KanbanView) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		m.fetchCards,
	)
}

// fetchCards fetches cards for the project from the GitHub API
func (m KanbanView) fetchCards() tea.Msg {
	columns, err := ghcli.GetProjectCards(m.projectID)
	return fetchCardsMsg{
		columns: columns,
		err:     err,
	}
}

// Update handles messages for the model
func (m KanbanView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready && len(m.columns) > 0 {
			// Initialize column viewports
			m.columnViewports = make([]viewport.Model, len(m.columns))
			columnWidth := (m.width - 4) / max(1, len(m.columns))

			for i := range m.columns {
				m.columnViewports[i] = viewport.New(columnWidth, m.height-10)
				m.columnViewports[i].SetContent(m.renderColumn(m.columns[i]))
			}

			m.ready = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			m.activeColumn = max(0, m.activeColumn-1)

		case "right", "l":
			m.activeColumn = min(len(m.columns)-1, m.activeColumn+1)

		case "up", "k":
			if len(m.columnViewports) > 0 {
				var cmd tea.Cmd
				m.columnViewports[m.activeColumn], cmd = m.columnViewports[m.activeColumn].Update(msg)
				cmds = append(cmds, cmd)
			}

		case "down", "j":
			if len(m.columnViewports) > 0 {
				var cmd tea.Cmd
				m.columnViewports[m.activeColumn], cmd = m.columnViewports[m.activeColumn].Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case fetchCardsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.columns = msg.columns

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
func (m KanbanView) View() string {
	if m.err != nil {
		return tui.Header("Error") + "\n\n" +
			fmt.Sprintf("Error fetching project cards: %v", m.err) + "\n\n" +
			tui.Footer("Press q to quit")
	}

	if m.loading {
		return tui.Header("Loading Project Cards") + "\n\n" +
			m.spinner.View() + "\n\n" +
			tui.Footer("Press q to quit")
	}

	if !m.ready || len(m.columnViewports) == 0 {
		return "Initializing..."
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

	return tui.Header("Project Kanban Board") + "\n" +
		b.String() + "\n" +
		tui.Footer("←/→: Change Column • ↑/↓: Scroll Cards • q: Quit")
}

// renderColumn renders a single column of the kanban board
func (m KanbanView) renderColumn(column models.ProjectColumn) string {
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

		title := card.Title
		if len(title) > cardStyle.GetWidth()-10 {
			title = title[:cardStyle.GetWidth()-10] + "..."
		}

		b.WriteString(cardStyle.Render(
			fmt.Sprintf("%s\n%s\n",
				title,
				statusStyle.Render(card.Status),
			),
		))
		b.WriteString("\n")
	}

	return b.String()
}

// getStatusColor returns a color based on status
func (m KanbanView) getStatusColor(status string) lipgloss.Color {
	switch status {
	case "To Do":
		return lipgloss.Color("#7997E5") // Blue
	case "In Progress":
		return lipgloss.Color("#FFBE3C") // Yellow/Orange
	case "Done":
		return lipgloss.Color("#5CBF40") // Green
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
