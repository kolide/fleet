package main

import (
	"github.com/kolide/fleet/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func setupCommand() cli.Command {
	var (
		flAddress            string
		flEmail              string
		flPassword           string
		flOrgName            string
		flInsecureSkipVerify bool
		flDebug              bool
	)
	return cli.Command{
		Name:      "setup",
		Usage:     "Setup a Kolide Fleet instance",
		UsageText: `fleetctl config login [options]`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "address",
				EnvVar:      "ADDRESS",
				Value:       "",
				Destination: &flAddress,
				Usage:       "The address of the Kolide Fleet instance",
			},
			cli.BoolFlag{
				Name:        "insecure-skip-verify",
				EnvVar:      "INSECURE_SKIP_VERIFY",
				Destination: &flInsecureSkipVerify,
				Usage:       "Whether or not to validate the remote TLS certificate",
			},
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
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
		Action: func(cliCtx *cli.Context) error {
			fleet, err := client.New(flAddress, flInsecureSkipVerify)
			if err != nil {
				return errors.Wrap(err, "error creating Fleet API client")
			}

			return fleet.Setup(flEmail, flPassword, flOrgName)
		},
	}
}
