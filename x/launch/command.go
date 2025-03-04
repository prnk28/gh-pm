package launch

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "launch",
		Short: "Manage launches",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(subCommands...)
	return cmd
}
