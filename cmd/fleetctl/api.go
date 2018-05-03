package main

import (
	"fmt"

	"github.com/kolide/fleet/client"
	"github.com/urfave/cli"
)

func clientFromCLI(c *cli.Context) (*client.Client, error) {
	config, err := readConfig(c.String("config"))
	if err != nil {
		return nil, err
	}

	cc, ok := config.Contexts[c.String("context")]
	if !ok {
		return nil, fmt.Errorf("context %q is not found", c.String("context"))
	}

	return client.New(cc.Address, cc.IgnoreTLS)
}
