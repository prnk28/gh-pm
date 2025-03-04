package app

import (
	"github.com/prnk28/gh-pm/pkg/exc"
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
			ghclicmd := []string{"api", "user", "--jq", `"You are @\(.login) (\(.name))"`}
			out, err := exc.Gh(ghclicmd...)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			cmd.Println(out)
			cmd.Help()
		},
	}
}

