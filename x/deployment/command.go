package deployment

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "Manage deployments",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(subCommands...)
	return cmd
}
