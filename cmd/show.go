package cmd

import (
	"log"

	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/models"
	"github.com/HenrySchwerdt/tt/tui"

	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show <path>",
	Short: "Shows a project and its sub projects",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("you must provide a project path")
		}

		projectPath := args[0]

		database, err := db.Init("./times.db")
		if err != nil {
			log.Fatalln(err)
		}

		var project *models.Project

		project, err = database.GetProjectByPathRecursive2(projectPath)

		if err != nil {
			log.Fatalln(err)
		}

		tui.RenderProjectTable(project)
	},
}

func init() {
	RootCmd.AddCommand(ShowCmd)
}
