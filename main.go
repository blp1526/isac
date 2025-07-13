package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/blp1526/isac/lib/cmd"
)

var version string

func main() {
	configPath := filepath.Join(os.Getenv("HOME"), ".usacloud", "default", "config.json")

	command := cmd.NewCommand(version, configPath)
	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
