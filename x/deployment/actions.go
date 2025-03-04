package deployment

import (
	"github.com/spf13/cobra"
)

var subCommands = []*cobra.Command{
	{
		Use:   "create",
		Short: "Create a deployment",
		Run:   createAction,
	},
	{
		Use:   "view",
		Short: "View a deployment",
		Run:   viewAction,
	},
	{
		Use:   "delete",
		Short: "Delete a deployment",
		Run:   deleteAction,
	},
	{
		Use:   "remove",
		Short: "Remove a deployment",
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
