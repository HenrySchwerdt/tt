package tui

import (
	"fmt"
	"time"

	"github.com/HenrySchwerdt/tt/models"
	"github.com/HenrySchwerdt/tt/utils"
)

func RenderTimeline(feed []models.FeedItem) {
	var total time.Duration
	var today time.Duration

	todayDate := time.Now().Format("2006-01-02")

	for i, item := range feed {
		isLast := i == len(feed)-1

		total += item.Duration

		if item.Start.Format("2006-01-02") == todayDate {
			today += item.Duration
		}

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

	fmt.Println()
	fmt.Println("━━ Summary ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total: %s\n", utils.FormatDuration2(total))
	fmt.Printf("Today: %s\n", utils.FormatDuration2(today))
}
