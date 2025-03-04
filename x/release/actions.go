package release

import (
	"github.com/spf13/cobra"
)

var subCommands = []*cobra.Command{
	{
		Use:   "create",
		Short: "Create a release",
		Run:   createAction,
	},
	{
		Use:   "view",
		Short: "View a release",
		Run:   viewAction,
	},
	{
		Use:   "delete",
		Short: "Delete a release",
		Run:   deleteAction,
	},
	{
		Use:   "remove",
		Short: "Remove a release",
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
