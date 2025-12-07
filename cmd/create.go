package cmd

import (
	"github.com/HenrySchwerdt/tt/db"
	tterrors "github.com/HenrySchwerdt/tt/tt_errors"
	"github.com/HenrySchwerdt/tt/tui"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create <path>",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.LogAndExitOnError(tterrors.ErrNoProjectPath)
		}
		projectName := args[0]
		db, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)
		project, err := db.CreateProject(projectName)
		utils.LogAndExitOnError(err)
		tui.RenderProjectTable(project)
	},
}

func init() {
	RootCmd.AddCommand(CreateCmd)
}
