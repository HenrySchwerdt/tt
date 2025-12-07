package cmd

import (
	"log"

	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/tui"

	"github.com/spf13/cobra"
)

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects in a modern TUI",
	Run: func(cmd *cobra.Command, args []string) {

		db, err := db.Init(DefaultDBPath())
		if err != nil {
			log.Fatalln(err)
		}
		projects, err := db.GetAllProjectsRecursive()
		if err != nil {
			log.Fatalln(err)
		}
		tui.RenderProjectsTree(projects)

	},
}

func init() {
	RootCmd.AddCommand(LsCmd)
}
