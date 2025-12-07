package cmd

import (
	"github.com/HenrySchwerdt/tt/db"
	"github.com/HenrySchwerdt/tt/models"
	tterrors "github.com/HenrySchwerdt/tt/tt_errors"
	"github.com/HenrySchwerdt/tt/tui"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show <path>",
	Short: "Shows a project and its sub projects",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.LogAndExitOnError(tterrors.ErrNoProjectPath)
		}

		projectPath := args[0]
		database, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)
		var project *models.Project
		project, err = database.GetProjectByPathRecursive2(projectPath)
		utils.LogAndExitOnError(err)
		tui.RenderProjectTable(project)
	},
}

func init() {
	RootCmd.AddCommand(ShowCmd)
}
