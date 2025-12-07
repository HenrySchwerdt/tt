package cmd

import (
	"fmt"

	"github.com/HenrySchwerdt/tt/db"
	tterrors "github.com/HenrySchwerdt/tt/tt_errors"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start <path>",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.LogAndExitOnError(tterrors.ErrNoProjectPath)
		}

		projectName := args[0]

		db, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)

		err = db.StartTimeEntry(projectName)
		utils.LogAndExitOnError(err)
		fmt.Printf("Time-Entry for project '%s' started\n", projectName)
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
}
