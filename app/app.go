package app

import "github.com/spf13/cobra"

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
			cmd.Help()
		},
	}
}

func DepStatus() map[string]bool {
	installStatus := map[string]bool{
		"gh":   hasCmdInstalled("gh"),
		"git":  hasCmdInstalled("git"),
		"hub":  hasCmdInstalled("fzf"),
		"jq":   hasCmdInstalled("jq"),
		"sed":  hasCmdInstalled("sed"),
		"curl": hasCmdInstalled("curl"),
		"gum":  hasCmdInstalled("gum"),
	}
	return installStatus
}
