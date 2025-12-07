package cmd

import (
	"fmt"

	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var message string

var EndCmd = &cobra.Command{
	Use:   "end",
	Short: "Ends the current time entry",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)

		err = db.EndTimeEntry(message)
		utils.LogAndExitOnError(err)

		fmt.Println("Time entry ended.")
	},
}

func init() {
	EndCmd.Flags().StringVarP(&message, "message", "m", "", "Optional message for this time entry")
	RootCmd.AddCommand(EndCmd)
}
