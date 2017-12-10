package main

import (
	"fmt"
	"os"
	"path/filepath"

	isac "github.com/blp1526/isac/lib"
	"github.com/blp1526/isac/lib/config"
	"github.com/urfave/cli"
)

const ExitCodeNG = 1

var unanonymize bool
var verbose bool
var version string
var zones string

func main() {
	configPath := filepath.Join(os.Getenv("HOME"), ".usacloud", "default", "config.json")

	app := cli.NewApp()
	app.Name = "isac"
	app.Usage = "interactive SAKURA Cloud"
	app.Version = version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Shingo Kawamura",
			Email: "blp1526@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "unanonymize",
			Usage:       "unanonymize personal data",
			Destination: &unanonymize,
		},

		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "print debug log",
			Destination: &verbose,
		},

		cli.StringFlag{
			Name:        "zones",
			Usage:       "set `ZONES` (separated by \",\", example: \"is1a,is1b,tk1a\")",
			Destination: &zones,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Creates config.json",
			Action: func(c *cli.Context) (err error) {
				err = config.CreateFile(configPath)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("%v", err), ExitCodeNG)
				}
				fmt.Printf("%v has been created\n", configPath)
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		i, err := isac.New(configPath, unanonymize, verbose, zones)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), ExitCodeNG)
		}

		err = i.Run()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), ExitCodeNG)
		}

		return nil
	}

	app.Run(os.Args)
}
