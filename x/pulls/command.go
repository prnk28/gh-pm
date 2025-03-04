package pulls

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pulls",
		Aliases: []string{
			"pr",
		},
		Short: "Manage reviews",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(subCommands...)
	return cmd
}
