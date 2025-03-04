package project

import (
	"github.com/spf13/cobra"
)

var subCommands = []*cobra.Command{
	{
		Use:   "create",
		Short: "Create a project",
		Run:   createAction,
	},
	{
		Use:   "view",
		Short: "View a project",
		Run:   viewAction,
	},
	{
		Use:   "delete",
		Short: "Delete a project",
		Run:   deleteAction,
	},
	{
		Use:   "remove",
		Short: "Remove a project",
		Run:   removeAction,
	},
}

func createAction(cmd *cobra.Command, args []string) {
}

func viewAction(cmd *cobra.Command, args []string) {
}

func deleteAction(cmd *cobra.Command, args []string) {
}

func removeAction(cmd *cobra.Command, args []string) {
}
