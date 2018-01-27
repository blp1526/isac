package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/blp1526/isac/lib/keybinding"
)

func main() {
	tmpl := template.Must(template.ParseFiles("_doc/README.tmpl.md"))
	f, err := os.Create("README.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	defer f.Close()

	keybindings := keybinding.Keybindings()

	err = tmpl.Execute(f, keybindings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
