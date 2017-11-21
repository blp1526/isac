package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "isac"
	app.Usage = "interactive SAKURA Cloud"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello isac!")
		return nil
	}
	app.Run(os.Args)
}
