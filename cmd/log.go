package main

import (
	"log"

	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/tui"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "log <path>",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("you must provide a project path")
		}

		projectPath := args[0]

		database, err := db.Init(DefaultDBPath())
		if err != nil {
			log.Fatalln(err)
		}

		feed, err := database.GetTimelineForProject(projectPath)
		if err != nil {
			log.Fatalln(err)
		}

		if len(feed) == 0 {
			log.Println("No finished time entries found for project", projectPath)
			return
		}

		tui.RenderTimeline(feed)

	},
}

func init() {
	RootCmd.AddCommand(LogCmd)
}
