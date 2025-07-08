package cmd

import (
	"context"
	"fmt"

	isac "github.com/blp1526/isac/lib"
	"github.com/blp1526/isac/lib/config"
	"github.com/urfave/cli/v3"
)

// NewCommand initializes *cli.Command for v3.
func NewCommand(version string, configPath string) *cli.Command {
	const exitCodeNG = 1

	var unanonymize bool
	var zones string

	cmd := &cli.Command{
		Name:    "isac",
		Version: version,
		Usage:   "Interactive SAKURA Cloud",
		Authors: []any{
			"Shingo Kawamura <blp1526@gmail.com>",
		},
		Copyright: "(c) 2017 Shingo Kawamura",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "unanonymize",
				Usage:       "unanonymize personal data",
				Destination: &unanonymize,
			},
			&cli.StringFlag{
				Name:        "zones",
				Usage:       "set `ZONES` (separated by \",\", example: \"is1a,is1b,tk1a\")",
				Destination: &zones,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Creates config.json",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := config.CreateFile(configPath)
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), exitCodeNG)
					}
					fmt.Printf("%v has been created\n", configPath)
					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			i, err := isac.New(configPath, unanonymize, zones)
			if err != nil {
				return cli.Exit(fmt.Sprintf("%v", err), exitCodeNG)
			}

			err = i.Run()
			if err != nil {
				return cli.Exit(fmt.Sprintf("%v", err), exitCodeNG)
			}

			return nil
		},
	}

	return cmd
}
