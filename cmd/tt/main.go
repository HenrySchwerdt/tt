package main

import (
	"log"

	"github.com/HenrySchwerdt/tt/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
