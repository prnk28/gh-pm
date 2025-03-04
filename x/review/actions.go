package review

import (
	"github.com/spf13/cobra"
)

var subCommands = []*cobra.Command{
	{
		Use:   "create",
		Short: "Create a PR",
		Run:   createAction,
	},
	{
		Use:   "view",
		Short: "View a PR",
		Run:   viewAction,
	},
	{
		Use:   "delete",
		Short: "Delete a PR",
		Run:   deleteAction,
	},
	{
		Use:   "remove",
		Short: "Remove a PR",
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
