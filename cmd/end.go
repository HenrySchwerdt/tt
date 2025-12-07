package main

import (
	"fmt"
	"log"

	"github.com/HenrySchwerdt/tt/db"

	"github.com/spf13/cobra"
)

var message string

var EndCmd = &cobra.Command{
	Use:   "end",
	Short: "Ends the current time entry",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.Init("./times.db")
		if err != nil {
			log.Fatalln(err)
		}

		err = db.EndTimeEntry(message)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Time entry ended.")
		if message != "" {
			fmt.Printf("Message: %q\n", message)
		}
	},
}

func init() {
	EndCmd.Flags().StringVarP(&message, "message", "m", "", "Optional message for this time entry")
	RootCmd.AddCommand(EndCmd)
}
