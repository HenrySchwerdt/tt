package cmd

import (
	"github.com/HenrySchwerdt/tt/tui"

	"github.com/spf13/cobra"
)

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects in a modern TUI",
	Run: func(cmd *cobra.Command, args []string) {

		// For now use mock data
		projects := tui.MockProjects()

		tui.RenderProjectTree(projects)

	},
}

func init() {
	RootCmd.AddCommand(LsCmd)
}
