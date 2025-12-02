package tui

import (
	"fmt"
	"sort"
	"strings"

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

// func formatDuration(d time.Duration) string {
// 	h := int(d.Hours())
// 	m := int(d.Minutes()) % 60
// 	return fmt.Sprintf("%dh %dm", h, m)
// }

func renderNode(p *models.Project, prefix string, last bool) {
	// tree branch characters
	var connector string
	if last {
		connector = "└─ "
	} else {
		connector = "├─ "
	}

	fmt.Println(prefix + connector + p.Name + " ") // TODO: set total time

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
			Name:     "PhD",
			Finished: false,
			Children: []*models.Project{
				{Name: "Research"},
				{Name: "Writing"},
			},
		},
		{
			Name:     "Side Project",
			Finished: true,
			Children: []*models.Project{
				{Name: "Experiment"},
			},
		},
	}
}
