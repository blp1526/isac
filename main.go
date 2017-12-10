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

var configPath string
var showServerID bool
var verbose bool
var version string
var zones string

func main() {
	defaultConfigPath := filepath.Join(os.Getenv("HOME"), ".usacloud", "default", "config.json")

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
		cli.StringFlag{
			Name:        "config, c",
			Value:       defaultConfigPath,
			Usage:       "Set `CONFIG_PATH`.",
			Destination: &configPath,
		},

		cli.BoolFlag{
			Name:        "show-server-id",
			Usage:       "Show server id.",
			Destination: &showServerID,
		},

		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "Print debug log.",
			Destination: &verbose,
		},

		cli.StringFlag{
			Name:        "zones",
			Usage:       "Set `ZONES` (separated by \",\", example: \"is1a,is1b,tk1a\").",
			Destination: &zones,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Create config.json",
			Action: func(c *cli.Context) (err error) {
				cfg := &config.Config{
					AccessToken:       "Write your AccessToken",
					AccessTokenSecret: "Write your AccessTokenSecret",
					Zone:              "Write your default Zone",
				}
				err = cfg.CreateFile(defaultConfigPath)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("%v", err), ExitCodeNG)
				}
				fmt.Printf("%v has been created\n", defaultConfigPath)
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		i, err := isac.New(configPath, showServerID, verbose, zones)
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
