package tui

import (
	"fmt"

	"github.com/HenrySchwerdt/tt/models"
	"github.com/HenrySchwerdt/tt/utils"
)

func RenderTimeline(feed []models.FeedItem) {
	for i, item := range feed {
		isLast := i == len(feed)-1

		// Dot + timestamp + duration
		fmt.Printf("● %s (%s)\n",
			item.Start.Format("2006-01-02 15:04"),
			utils.FormatDuration2(item.Duration),
		)
		fmt.Printf("│   %s\n", item.ProjectPath)

		if item.Message != "" {
			fmt.Printf("│   “%s”\n", item.Message)
		}

		if !isLast {
			fmt.Println("│")
		}
	}
}
