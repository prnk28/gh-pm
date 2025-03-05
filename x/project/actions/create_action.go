package actions

import (
	"fmt"
	"os"

	"github.com/prnk28/gh-pm/internal/ctx"
	"github.com/prnk28/gh-pm/x/project/views"
	"github.com/spf13/cobra"
)

// CreateAction handles the 'project create' command
func CreateAction(cmd *cobra.Command, args []string) {
	c, err := ctx.Get(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	// Create and run the form
	form, err := views.NewProjectForm(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if the form was submitted or canceled
	if !form.Submitted {
		fmt.Println("Project creation canceled")
		return
	}

	// Create the project using the GitHub API
	// This would be your actual implementation
	// err = ghcli.CreateProject(form.Title, form.Organization, form.Description)
	// if err != nil {
	//     fmt.Fprintf(os.Stderr, "Error creating project: %v\n", err)
	//     os.Exit(1)
	// }

	// For now, just simulate success
	fmt.Println(form.FormatSummary())
}
