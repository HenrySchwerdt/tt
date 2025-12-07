package cmd

import (
	"log"

	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/tui"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create <path>",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("you must provide a project name")
		}

		projectName := args[0]

		db, err := db.Init(DefaultDBPath())
		if err != nil {
			log.Fatalln(err)
		}

		project, err := db.CreateProject(projectName)
		if err != nil {
			log.Fatalln(err)
		}
		tui.RenderProjectTable(project)
	},
}

func init() {
	RootCmd.AddCommand(CreateCmd)
}
