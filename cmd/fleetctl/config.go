package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/kolide/kit/env"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type configFile struct {
	Contexts map[string]Context `json:"contexts"`
}

type Context struct {
	Address   string `json:"address"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	IgnoreTLS bool   `json:"ignore_tls"`
}

func makeConfigIfNotExists(fp string) error {
	if _, err := os.Stat(filepath.Dir(fp)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Dir(fp), 0700); err != nil {
			return err
		}
	}

	_, err := os.OpenFile(fp, os.O_RDONLY|os.O_CREATE, 0600)
	return err
}

func readConfig(fp string) (c configFile, err error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &c)

	if c.Contexts == nil {
		c.Contexts = map[string]Context{
			"default": Context{},
		}
	}
	return
}

func writeConfig(fp string, c configFile) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fp, b, 0400)
}

func configSetCommand() cli.Command {
	var (
		flContext string
		flConfig  string
	)

	return cli.Command{
		Name:      "set",
		Usage:     "Set a config option",
		UsageText: `fleetctl config set [options]`,
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 2 {
				return cli.ShowCommandHelp(c, "set")
			}

			flConfig = c.String("config")
			flContext = c.String("context")

			key, value := c.Args()[0], c.Args()[1]

			// validate key
			switch key {
			case "address", "email", "token", "ignore_tls":
			default:
				return cli.ShowCommandHelp(c, "set")
			}

			if err := makeConfigIfNotExists(flConfig); err != nil {
				return errors.Wrapf(err, "error verifying that config exists at %s", flConfig)
			}

			config, err := readConfig(flConfig)
			if err != nil {
				return errors.Wrapf(err, "error reading config at %s", flConfig)
			}

			currentContext, ok := config.Contexts[flContext]
			if !ok {
				fmt.Printf("[+] Context %q not found, creating it with default values\n", flContext)
				currentContext = Context{}
			}
			switch key {
			case "address":
				currentContext.Address = value
			case "email":
				currentContext.Email = value
			case "token":
				currentContext.Token = value
			case "ignore_tls":
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					return errors.Wrapf(err, "error parsing %q as bool", value)
				}
				currentContext.IgnoreTLS = boolValue
			default:
				return cli.ShowCommandHelp(c, "set")
			}

			config.Contexts[flContext] = currentContext

			if err := writeConfig(flConfig, config); err != nil {
				return errors.Wrap(err, "error saving config file")
			}

			fmt.Printf("[+] Set the %q config key to %q in the %q context\n", key, value, flContext)
			return nil
		},
	}
}

func configGetCommand() cli.Command {
	var (
		flContext string
		flConfig  string
	)

	return cli.Command{
		Name:      "get",
		Usage:     "Get a config option",
		UsageText: `fleetctl config get [options]`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "context",
				Value:       "default",
				EnvVar:      "CONTEXT",
				Usage:       "The Fleet config context",
				Destination: &flContext,
			},
			cli.StringFlag{
				Name:        "config",
				Value:       fmt.Sprintf("%s/.fleet/config", env.String("HOME", "~/")),
				EnvVar:      "CONFIG",
				Usage:       "The path to the Fleet config file",
				Destination: &flConfig,
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.ShowCommandHelp(c, "get")
			}

			key := c.Args()[0]

			// validate key
			switch key {
			case "address", "email", "token", "ignore_tls":
			default:
				return cli.ShowCommandHelp(c, "get")
			}

			if err := makeConfigIfNotExists(flConfig); err != nil {
				return errors.Wrapf(err, "error verifying that config exists at %s", flConfig)
			}

			config, err := readConfig(flConfig)
			if err != nil {
				return errors.Wrapf(err, "error reading config at %s", flConfig)
			}

			currentContext, ok := config.Contexts[flContext]
			if !ok {
				fmt.Printf("[+] Context %q not found, creating it with default values\n", flContext)
				currentContext = Context{}
			}

			switch key {
			case "address":
				fmt.Printf("  %s.%s => %s\n", flContext, key, currentContext.Address)
			case "email":
				fmt.Printf("  %s.%s => %s\n", flContext, key, currentContext.Email)
			case "token":
				fmt.Printf("  %s.%s => %s\n", flContext, key, currentContext.Token)
			case "ignore_tls":
				fmt.Printf("  %s.%s => %s\n", flContext, key, currentContext.IgnoreTLS)
			default:
				return cli.ShowCommandHelp(c, "set")
			}

			return nil
		},
	}
}
