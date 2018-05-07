package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func getQueriesCommand() cli.Command {
	return cli.Command{
		Name:    "queries",
		Aliases: []string{"query", "q"},
		Usage:   "List information about one or more queries",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			name := c.Args().First()

			// if name wasn't provided, list all queries
			if name == "" {
				queries, err := fleet.GetQuerySpecs()
				if err != nil {
					return errors.Wrap(err, "could not list queries")
				}

				if len(queries) == 0 {
					fmt.Println("no queries found")
					return nil
				}

				data := [][]string{}

				for _, query := range queries {
					data = append(data, []string{
						query.Name,
						query.Description,
						query.Query,
					})
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetRowLine(true)
				table.SetHeader([]string{"name", "description", "query"})
				table.AppendBulk(data)
				table.Render()

				return nil
			} else {
				return nil
			}
		},
	}
}
