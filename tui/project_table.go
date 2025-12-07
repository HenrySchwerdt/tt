package tui

import (
	"os"
	"time"

	"github.com/HenrySchwerdt/tt/models"
	"github.com/HenrySchwerdt/tt/utils"
	"github.com/olekukonko/tablewriter"
)

func RenderProjectTable(root *models.Project) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "Path", "Finished", "Created At", "Time Spent"})

	table.Append([]string{
		root.Name,
		root.Path,
		formatBool(root.Finished),
		formatTime(root.CreatedAt),
		utils.FormatDuration(root.TimeSpend),
	})

	var addChildren func(children []*models.Project, prefix string)
	addChildren = func(children []*models.Project, prefix string) {
		for i, child := range children {
			isLast := i == len(children)-1
			branch := "├─ "
			nextPrefix := prefix + "│  "
			if isLast {
				branch = "└─ "
				nextPrefix = prefix + "   "
			}

			name := prefix + branch + child.Name

			table.Append([]string{
				name,
				child.Path,
				formatBool(child.Finished),
				formatTime(child.CreatedAt),
				utils.FormatDuration(child.TimeSpend),
			})

			addChildren(child.Children, nextPrefix)
		}
	}

	addChildren(root.Children, "")
	table.Render()
}

func formatBool(f bool) string {
	if f {
		return "Yes"
	}
	return "No"
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
