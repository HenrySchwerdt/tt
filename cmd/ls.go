package cmd

import (
	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/tui"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects and project structure",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)
		projects, err := db.GetAllProjectsRecursive()
		utils.LogAndExitOnError(err)
		tui.RenderProjectsTree(projects)
	},
}

func init() {
	RootCmd.AddCommand(LsCmd)
}
