package tui

import (
	"fmt"

	"github.com/HenrySchwerdt/tt/models"
	"github.com/HenrySchwerdt/tt/utils"
)

// / RenderProjectsTree renders multiple roots as a full tree.
func RenderProjectsTree(roots []*models.Project) {
	for i, root := range roots {
		last := i == len(roots)-1
		printProjectNode(root, "", last)
	}
}

// printProjectNode prints a single node + recursively prints children.
func printProjectNode(p *models.Project, prefix string, isLast bool) {
	branch := "├── "
	nextPrefix := prefix + "│   "
	if isLast {
		branch = "└── "
		nextPrefix = prefix + "    "
	}

	// Print current project
	fmt.Printf(
		"%s%s%s (%s)\n",
		prefix,
		branch,
		p.Name,
		utils.FormatDuration(p.TimeSpend),
	)

	// Recurse over children
	for i, child := range p.Children {
		printProjectNode(child, nextPrefix, i == len(p.Children)-1)
	}
}
