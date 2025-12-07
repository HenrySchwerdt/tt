package cmd

import (
	"fmt"
	"log"

	"github.com/HenrySchwerdt/tt/db"

	"github.com/spf13/cobra"
)

var EndCmd = &cobra.Command{
	Use:   "end",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.Init("./times.db")
		if err != nil {
			log.Fatalln(err)
		}

		err = db.EndTimeEntry()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Time-Entry endend\n") // log logged time
	},
}

func init() {
	RootCmd.AddCommand(EndCmd)
}
