package cmd

import (
	"fmt"
	"log"

	"github.com/HenrySchwerdt/tt/db"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start <path>",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("you must provide a project name")
		}

		projectName := args[0]

		db, err := db.Init("./times.db")
		if err != nil {
			log.Fatalln(err)
		}

		err = db.StartTimeEntry(projectName)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Time-Entry for project '%s' started\n", projectName)
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
}
