package actions

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/prnk28/gh-pm/x/project/views"
)

// ListAction handles the 'project list' command
func ListAction(cmd *cobra.Command, args []string) {
	// Create and run the BubbleTea program
	p := tea.NewProgram(
		views.NewProjectsListViewModel(),
		tea.WithAltScreen(),       // Use the full terminal window
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
