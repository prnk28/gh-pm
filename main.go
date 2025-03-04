package main

import (
	"fmt"
	"os"

	"github.com/prnk28/gh-pm/x/deployment"
	"github.com/prnk28/gh-pm/x/launch"
	"github.com/prnk28/gh-pm/x/milestone"
	"github.com/prnk28/gh-pm/x/project"
	"github.com/prnk28/gh-pm/x/review"
	"github.com/prnk28/gh-pm/x/todo"
	"github.com/prnk28/gh-pm/x/whoami"

	"github.com/prnk28/gh-pm/app"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	deployment.Command(),
	launch.Command(),
	milestone.Command(),
	project.Command(),
	review.Command(),
	todo.Command(),
	whoami.Command(),
}

func main() {
	rootCmd := app.RootCmd()
	rootCmd.AddCommand(commands...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
