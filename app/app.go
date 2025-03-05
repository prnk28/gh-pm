package app

import (
	"github.com/prnk28/gh-pm/internal/ghcli"
	"github.com/spf13/cobra"
)

const (
	Name  = "pm"
	Long  = "A Github CLI Extension for managing projects"
	Short = `gh pm [command]`
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pm",
		Short: "gh pm [command]",
		Long:  "A Github CLI Extension for managing projects",
		Run: func(cmd *cobra.Command, args []string) {
			out, err := ghcli.Whoami()
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			cmd.Println(out)
			cmd.Help()
		},
	}
}
