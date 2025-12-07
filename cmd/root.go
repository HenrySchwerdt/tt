package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tt",
	Short: "Terminal Tracker",
}

const MAJOR = 0
const MINOR = 0
const PATCH = 2

func DefaultDBPath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		configDir = os.Getenv("APPDATA")
	} else {
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}

	dbDir := filepath.Join(configDir, "tt")
	_ = os.MkdirAll(dbDir, os.ModePerm)

	return filepath.Join(dbDir, "tt.db")
}
