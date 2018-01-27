package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/blp1526/isac/lib/cmd"
	"github.com/blp1526/isac/lib/keybinding"
	"github.com/urfave/cli"
)

func main() {
	tmpl := template.Must(template.ParseFiles("_generator/README.tmpl.md"))
	f, err := os.Create("README.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	defer f.Close()

	data := &struct {
		Keybindings []keybinding.Keybinding
		Options     []cli.Flag
	}{
		Keybindings: keybinding.Keybindings(),
		Options:     cmd.NewApp("", false, "", "").Flags,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
