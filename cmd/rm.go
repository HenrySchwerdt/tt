package cmd

import (
	"fmt"

	"github.com/HenrySchwerdt/tt/db"
	tterrors "github.com/HenrySchwerdt/tt/tt_errors"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "rm <path>",
	Short: "Removes a project with the given path and all sub projects",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.LogAndExitOnError(tterrors.ErrNoProjectPath)
		}
		projectName := args[0]

		db, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)

		err = db.RemoveProject(projectName)
		utils.LogAndExitOnError(err)
		fmt.Printf("All projects from '%s' downwards are removed\n", projectName)
	},
}

func init() {
	RootCmd.AddCommand(RemoveCmd)
}
