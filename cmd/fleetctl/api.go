package main

import (
	"fmt"

	"github.com/kolide/fleet/server/service"
	"github.com/urfave/cli"
)

func clientFromCLI(c *cli.Context) (*service.Client, error) {
	config, err := readConfig(c.String("config"))
	if err != nil {
		return nil, err
	}

	cc, ok := config.Contexts[c.String("context")]
	if !ok {
		return nil, fmt.Errorf("context %q is not found", c.String("context"))
	}

	return service.NewClient(cc.Address, cc.IgnoreTLS)
}
