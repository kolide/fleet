package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kolide/kit/env"
	"github.com/kolide/kit/version"
	"github.com/urfave/cli"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "fleetctl"
	app.Usage = "The CLI for operating Kolide Fleet"
	app.Version = version.Version().Version
	cli.VersionPrinter = func(c *cli.Context) {
		version.PrintFull()
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  fmt.Sprintf("%s/.fleet/config", env.String("HOME", "~/")),
			EnvVar: "CONFIG",
			Usage:  "The path to the Fleet config file",
		},
		cli.StringFlag{
			Name:   "context",
			Value:  "default",
			EnvVar: "CONTEXT",
			Usage:  "The Fleet config context",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "query",
			Usage:       "run a query across your fleet",
			Subcommands: []cli.Command{},
		},
		cli.Command{
			Name:        "apply",
			Usage:       "apply a set of osquery configurations",
			Subcommands: []cli.Command{},
		},
		cli.Command{
			Name:        "edit",
			Usage:       "edit your complete configuration in an ephemeral editor",
			Subcommands: []cli.Command{},
		},
		setupCommand(),
		loginCommand(),
		cli.Command{
			Name:  "config",
			Usage: "modify how and which Fleet server to connect to",
			Subcommands: []cli.Command{
				configSetCommand(),
				configGetCommand(),
			},
		},
	}

	app.RunAndExitOnError()
}
