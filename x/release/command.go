package release

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Manage releases",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(subCommands...)
	return cmd
}
