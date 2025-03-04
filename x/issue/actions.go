package issue

import (
	"github.com/spf13/cobra"
)

var subCommands = []*cobra.Command{
	{
		Use:   "create",
		Short: "Create a todo",
		Run:   createAction,
	},
	{
		Use:   "view",
		Short: "View a todo",
		Run:   viewAction,
	},
	{
		Use:   "delete",
		Short: "Delete a todo",
		Run:   deleteAction,
	},
	{
		Use:   "complete",
		Short: "Complete a todo",
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
