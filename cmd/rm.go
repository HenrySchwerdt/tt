package cmd

import (
	"fmt"
	"log"

	"github.com/HenrySchwerdt/tt/db"

	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "rm <path>",
	Short: "Removes a project with the given path and all sub projects",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("you must provide a project name")
		}

		projectName := args[0]

		db, err := db.Init("./times.db")
		if err != nil {
			log.Fatalln(err)
		}

		err = db.RemoveProject(projectName)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("All projects from '%s' downwards are removed\n", projectName)
	},
}

func init() {
	RootCmd.AddCommand(RemoveCmd)
}
