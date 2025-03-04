package main

import (
	"fmt"
	"os"

	"github.com/prnk28/gh-pm/x/deployment"
	"github.com/prnk28/gh-pm/x/release"
	"github.com/prnk28/gh-pm/x/milestone"
	"github.com/prnk28/gh-pm/x/project"
	"github.com/prnk28/gh-pm/x/pulls"
	"github.com/prnk28/gh-pm/x/issue"

	"github.com/prnk28/gh-pm/app"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	deployment.Command(),
	release.Command(),
	milestone.Command(),
	project.Command(),
	pulls.Command(),
	issue.Command(),
}

func main() {
	rootCmd := app.RootCmd()
	rootCmd.AddCommand(commands...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
