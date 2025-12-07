package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tt",
	Short: "Terminal Time Tracker",
}

func DefaultDBPath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		configDir = os.Getenv("APPDATA") // typically C:\Users\<User>\AppData\Roaming
	} else {
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}

	dbDir := filepath.Join(configDir, "tt")
	_ = os.MkdirAll(dbDir, os.ModePerm) // ensure folder exists

	return filepath.Join(dbDir, "tt.db")
}
