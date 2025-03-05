package app

import (
	"github.com/prnk28/gh-pm/internal/ctx"
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
			c, err := ctx.Get(cmd)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			for _, org := range c.Orgs {
				cmd.Println(org)
			}
			cmd.Println(c.String())
			cmd.Help()
		},
	}
}
