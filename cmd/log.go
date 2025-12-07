package cmd

import (
	"fmt"

	"github.com/HenrySchwerdt/tt/db"
	tterrors "github.com/HenrySchwerdt/tt/tt_errors"
	"github.com/HenrySchwerdt/tt/tui"
	"github.com/HenrySchwerdt/tt/utils"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "log <path>",
	Short: "Creates a project with the given path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.LogAndExitOnError(tterrors.ErrNoProjectPath)
		}

		projectPath := args[0]

		database, err := db.Init(DefaultDBPath())
		utils.LogAndExitOnError(err)

		feed, err := database.GetTimelineForProject(projectPath)
		utils.LogAndExitOnError(err)

		if len(feed) == 0 {
			fmt.Println("No finished time entries found for project", projectPath)
			return
		}

		tui.RenderTimeline(feed)

	},
}

func init() {
	RootCmd.AddCommand(LogCmd)
}
