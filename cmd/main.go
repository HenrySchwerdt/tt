package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

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

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
