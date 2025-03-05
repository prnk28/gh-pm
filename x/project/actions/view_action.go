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

	// TODO: Fix ProjectSelector
	//
	// // Check if project ID was provided
	// if len(args) > 0 {
	// 	projectID = args[0]
	// } else {
	// 	// If no project ID provided, first show the list of projects
	// 	p := tea.NewProgram(
	// 		newProjectSelectorModel(),
	// 		tea.WithAltScreen(),
	// 		tea.WithMouseCellMotion(),
	// 	)
	//
	// 	model, err := p.Run()
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
	// 		os.Exit(1)
	// 	}
	//
	// 	// Check if a project was selected
	// 	if selector, ok := model.(projectSelectorModel); ok && selector.selectedProject.Project.ID != "" {
	// 		projectID = selector.selectedProject.Project.ID
	// 	} else {
	// 		// No project selected, exit
	// 		return
	// 	}
	// }

	// Now show the project view
	p := tea.NewProgram(
		views.NewProjectsView(projectID, true),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// projectSelectorModel is a simplified version of deleteModel for selecting a project
type projectSelectorModel struct {
	// Same as deleteModel but with only the project selection functionality
	// We can reuse most of the code from deleteModel
	// For brevity, I'm not including the full implementation here
	selectedProject ProjectItem
}

func newProjectSelectorModel() projectSelectorModel {
	// Similar to newDeleteModel but without the deletion functionality
	return projectSelectorModel{}
}
