package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/HenrySchwerdt/tt/models"
)

// RenderProjectTree prints the project tree using ASCII characters.
func RenderProjectTree(projects []*models.Project) {
	fmt.Println("Project:")
	for i, p := range projects {
		isLast := i == len(projects)-1
		renderNode(p, "", isLast)
	}
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", h, m)
}

func renderNode(p *models.Project, prefix string, last bool) {
	// tree branch characters
	var connector string
	if last {
		connector = "└─ "
	} else {
		connector = "├─ "
	}

	fmt.Println(prefix + connector + p.Name + " " + formatDuration(p.TotalTime))

	// next prefix
	var newPrefix string
	if last {
		newPrefix = prefix + "   "
	} else {
		newPrefix = prefix + "│  "
	}

	// ensure stable ordering (optional)
	sort.Slice(p.Children, func(i, j int) bool {
		return strings.Compare(p.Children[i].Name, p.Children[j].Name) < 0
	})

	for i, child := range p.Children {
		renderNode(child, newPrefix, i == len(p.Children)-1)
	}
}

func MockProjects() []*models.Project {
	return []*models.Project{
		{
			Name:      "PhD",
			TotalTime: 12*time.Hour + 35*time.Minute,
			Finished:  false,
			Children: []*models.Project{
				{Name: "Research", TotalTime: 5*time.Hour + 20*time.Minute},
				{Name: "Writing", TotalTime: 7*time.Hour + 15*time.Minute},
			},
		},
		{
			Name:      "Side Project",
			TotalTime: 3*time.Hour + 10*time.Minute,
			Finished:  true,
			Children: []*models.Project{
				{Name: "Experiment", TotalTime: 3*time.Hour + 10*time.Minute, Finished: true},
			},
		},
	}
}
