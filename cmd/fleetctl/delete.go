package main

import (
	"io/ioutil"

	"github.com/kolide/fleet/server/service"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func deleteCommand() cli.Command {
	var (
		flFilename string
		flDebug    bool
	)
	return cli.Command{
		Name:      "delete",
		Usage:     "Specify files to declaratively batch delete osquery configurations",
		UsageText: `fleetctl delete [options]`,
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			cli.StringFlag{
				Name:        "f",
				EnvVar:      "FILENAME",
				Value:       "",
				Destination: &flFilename,
				Usage:       "A file to apply",
			},
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(c *cli.Context) error {
			if flFilename == "" {
				return errors.New("-f must be specified")
			}

			b, err := ioutil.ReadFile(flFilename)
			if err != nil {
				return err
			}

			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			specs, err := specGroupFromBytes(b)
			if err != nil {
				return err
			}

			for _, query := range specs.Queries {
				if err := fleet.DeleteQuery(query.Name); err != nil {
					switch err.(type) {
					case service.NotFoundErr:
						continue
					}
					return err
				}
			}

			for _, pack := range specs.Packs {
				if err := fleet.DeletePack(pack.Name); err != nil {
					switch err.(type) {
					case service.NotFoundErr:
						continue
					}
					return err
				}
			}

			for _, label := range specs.Labels {
				if err := fleet.DeleteLabel(label.Name); err != nil {
					switch err.(type) {
					case service.NotFoundErr:
						continue
					}
					return err
				}
			}

			return nil
		},
	}
}
