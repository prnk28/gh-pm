package whoami

import (
	"github.com/prnk28/gh-pm/internal/gh"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show who you are",
		Run: func(cmd *cobra.Command, args []string) {
			ghclicmd := []string{"api", "user", "--jq", `"You are @\(.login) (\(.name))"`}
			out, err := gh.Exec(ghclicmd...)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			cmd.Println(out)
		},
	}
	return cmd
}
