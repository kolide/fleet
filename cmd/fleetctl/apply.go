package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type genericSpec struct {
	Kind    string      `json:"kind"`
	Version string      `json:"apiVersion"`
	Spec    interface{} `json:"spec"`
}

type specContainer struct {
	Queries []*kolide.QuerySpec
	Packs   []*kolide.PackSpec
	Labels  []*kolide.LabelSpec
	Options *kolide.OptionsSpec
}

func applyCommand() cli.Command {
	var (
		flFilename string
		flDebug    bool
	)
	return cli.Command{
		Name:      "apply",
		Usage:     "Apply files to declaratively manage osquery configurations",
		UsageText: `fleetctl apply [options]`,
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
				return errors.Wrap(err, "error reading file")
			}

			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			specs := &specContainer{
				Queries: []*kolide.QuerySpec{},
				Packs:   []*kolide.PackSpec{},
				Labels:  []*kolide.LabelSpec{},
			}

			for _, specYaml := range strings.Split(string(b), "---") {
				if strings.TrimSpace(specYaml) == "" {
					if flDebug {
						fmt.Println("[!] Skipping empty spec in file")
					}
					continue
				}

				var s genericSpec
				if err := yaml.Unmarshal([]byte(specYaml), &s); err != nil {
					return errors.Wrap(err, "error unmarshaling spec")
				}

				if s.Spec == nil {
					return errors.Errorf("no spec field on %q document", s.Kind)
				}

				specBytes, err := yaml.Marshal(s.Spec)
				if err != nil {
					return errors.Errorf("error marshaling spec for %q kind", s.Kind)
				}

				switch strings.ToLower(s.Kind) {
				case "query":
					var querySpec *kolide.QuerySpec
					if err := yaml.Unmarshal(specBytes, &querySpec); err != nil {
						return errors.Wrap(err, "error unmarshaling query spec")
					}
					specs.Queries = append(specs.Queries, querySpec)

				case "pack":
					var packSpec *kolide.PackSpec
					if err := yaml.Unmarshal(specBytes, &packSpec); err != nil {
						return errors.Wrap(err, "error unmarshaling pack spec")
					}
					specs.Packs = append(specs.Packs, packSpec)

				case "label":
					var labelSpec *kolide.LabelSpec
					if err := yaml.Unmarshal(specBytes, &labelSpec); err != nil {
						return errors.Wrap(err, "error unmarshaling label spec")
					}
					specs.Labels = append(specs.Labels, labelSpec)

				case "options":
					if specs.Options != nil {
						return errors.New("options defined twice in the same file")
					}

					var optionSpec *kolide.OptionsSpec
					if err := yaml.Unmarshal(specBytes, &optionSpec); err != nil {
						return errors.Wrap(err, "error unmarshaling option spec")
					}
					specs.Options = optionSpec

				default:
					return errors.Errorf("unknown kind %q", s.Kind)
				}
			}

			if err := fleet.ApplyQuerySpecs(specs.Queries); err != nil {
				return errors.Wrap(err, "error applying queries")
			}
			if flDebug {
				fmt.Printf("[+] applied %d queries\n", len(specs.Queries))
			}

			if err := fleet.ApplyLabelSpecs(specs.Labels); err != nil {
				return errors.Wrap(err, "error applying labels")
			}
			if flDebug {
				fmt.Printf("[+] applied %d labels\n", len(specs.Labels))
			}

			if err := fleet.ApplyPackSpecs(specs.Packs); err != nil {
				return errors.Wrap(err, "error applying packs")
			}
			if flDebug {
				fmt.Printf("[+] applied %d packs\n", len(specs.Packs))
			}

			return nil
		},
	}
}
