package main

import (
	"os"
	"path/filepath"

	"github.com/blp1526/isac/lib/cmd"
)

var version string

func main() {
	var unanonymize bool
	var zones string
	configPath := filepath.Join(os.Getenv("HOME"), ".usacloud", "default", "config.json")

	app := cmd.NewApp(version, unanonymize, zones, configPath)
	app.Run(os.Args)
}
