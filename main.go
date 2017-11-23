package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

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
		cli.StringFlag{
			Name:        "zones",
			Value:       "is1a,is1b,tk1a",
			Usage:       `Request zones (separated by ",")`,
			Destination: &zones,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Printf("%v\n", zones)
		return nil
	}
	app.Run(os.Args)
}
