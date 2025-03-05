package actions

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prnk28/gh-pm/x/project/views"
	"github.com/spf13/cobra"
)

// ViewAction handles the 'project view' command
func ViewAction(cmd *cobra.Command, args []string) {
	var projectID string

	// Check if project ID was provided
	if len(args) == 0 {
		fmt.Println("Error: Project ID is required.")
		fmt.Println("Usage: gh pm project cards <project-id>")
		os.Exit(1)
	}

	projectID = args[0]

	// Create and run the BubbleTea program with the kanban view
	p := tea.NewProgram(
		views.NewKanbanView(projectID),
		tea.WithAltScreen(),       // Use the full terminal window
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
