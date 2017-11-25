package main

import (
	"fmt"
	"os"

	isac "github.com/blp1526/isac/lib"
	"github.com/urfave/cli"
)

var verbose bool
var version string
var zones string

func main() {

	app := cli.NewApp()
	app.Name = "isac"
	app.Usage = "interactive SAKURA Cloud"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Shingo Kawamura",
			Email: "blp1526@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "Print debug log.",
			Destination: &verbose,
		},
		cli.StringFlag{
			Name:        "zones",
			Value:       "is1a,is1b,tk1a",
			Usage:       "Set `ZONES` (separated by \",\")",
			Destination: &zones,
		},
	}

	app.Action = func(c *cli.Context) error {
		i := isac.New(verbose, zones)

		err := i.Run()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err), 1)
		}
		return nil
	}
	app.Run(os.Args)
}
