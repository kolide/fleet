package main

import (
	"fmt"

	"github.com/kolide/fleet/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func setupCommand() cli.Command {
	var (
		flEmail    string
		flPassword string
		flOrgName  string
	)
	return cli.Command{
		Name:      "setup",
		Usage:     "Setup a Kolide Fleet instance",
		UsageText: `fleetctl config login [options]`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "email",
				EnvVar:      "EMAIL",
				Value:       "",
				Destination: &flEmail,
				Usage:       "The email of the admin user to create",
			},
			cli.StringFlag{
				Name:        "password",
				EnvVar:      "PASSWORD",
				Value:       "",
				Destination: &flPassword,
				Usage:       "The password for the admin user",
			},
			cli.StringFlag{
				Name:        "org-name",
				EnvVar:      "ORG_NAME",
				Value:       "",
				Destination: &flOrgName,
				Usage:       "The name of the organization",
			},
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return errors.Wrap(err, "error creating Fleet API client")
			}

			token, err := fleet.Setup(flEmail, flPassword, flOrgName)
			if err != nil {
				// the Kolide Fleet instance has already been setup
				if setupErr, ok := err.(client.SetupAlreadyErr); ok {
					return setupErr
				}
				return errors.Wrap(err, "error setting up Fleet")
			}

			fmt.Println("Token:", token)

			return nil
		},
	}
}
