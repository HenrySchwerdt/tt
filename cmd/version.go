package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version of Terminal Tracker",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tt - v%d.%d.%d\n", MAJOR, MINOR, PATCH)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
