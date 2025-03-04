package milestone

import (
	"github.com/spf13/cobra"
)

var subCommands = []*cobra.Command{
	{
		Use:   "create",
		Short: "Create a milestone",
		Run:   createAction,
	},
	{
		Use:   "view",
		Short: "View a milestone",
		Run:   viewAction,
	},
	{
		Use:   "delete",
		Short: "Delete a milestone",
		Run:   deleteAction,
	},
	{
		Use:   "remove",
		Short: "Remove a milestone",
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
