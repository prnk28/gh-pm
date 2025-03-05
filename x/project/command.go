package project

import (
	"github.com/prnk28/gh-pm/x/project/actions"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	// Define the subcommands
	subCommands := []*cobra.Command{
		{
			Use:   "create",
			Short: "Create a project",
			Run:   actions.CreateAction,
		},
		{
			Use:   "list",
			Short: "List all projects",
			Run:   actions.ListAction,
		},
	}

	// Create the root command
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
		Run: func(cmd *cobra.Command, args []string) {
			actions.ListAction(cmd, args)
		},
	}

	// Add the subcommands to the root command
	cmd.AddCommand(subCommands...)
	return cmd
}
