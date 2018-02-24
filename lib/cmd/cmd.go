package cmd

import (
	"fmt"

	isac "github.com/blp1526/isac/lib"
	"github.com/blp1526/isac/lib/config"
	"github.com/urfave/cli"
)

// NewApp initializes *cli.App.
func NewApp(version string, unanonymize bool, zones string, configPath string) *cli.App {
	const exitCodeNG = 1

	app := cli.NewApp()
	app.Name = "isac"
	app.Version = version
	app.Usage = "Interactive SAKURA Cloud"
	app.Version = version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Shingo Kawamura",
			Email: "blp1526@gmail.com",
		},
	}
	app.Copyright = "(c) 2017 Shingo Kawamura"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "unanonymize",
			Usage:       "unanonymize personal data",
			Destination: &unanonymize,
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
					return cli.NewExitError(fmt.Sprintf("%v", err), exitCodeNG)
				}
				fmt.Printf("%v has been created\n", configPath)
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		i, err := isac.New(configPath, unanonymize, zones)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), exitCodeNG)
		}

		err = i.Run()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), exitCodeNG)
		}

		return nil
	}

	return app
}
