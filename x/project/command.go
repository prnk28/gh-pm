package project

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(subCommands...)
	return cmd
}
