package cli

import (
	"context"
	"fmt"

	"github.com/kolide/kolide-ose/config"
	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
	"github.com/kolide/kolide-ose/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createPrepareCmd(configManager config.Manager) *cobra.Command {

	var prepareCmd = &cobra.Command{
		Use:   "prepare",
		Short: "Subcommands for initializing kolide infrastructure",
		Long: `
Subcommands for initializing kolide infrastructure

To setup kolide infrastructure, use one of the available commands.
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Given correct database configurations, prepare the databases for use",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			config := configManager.LoadConfig()
			connString := fmt.Sprintf(
				"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
				config.Mysql.Username,
				config.Mysql.Password,
				config.Mysql.Address,
				config.Mysql.Database,
			)
			ds, err := datastore.New("gorm-mysql", connString)
			if err != nil {
				initFatal(err, "creating db connection")
			}
			if err := ds.Drop(); err != nil {
				initFatal(err, "dropping db tables")
			}

			if err := ds.Migrate(); err != nil {
				initFatal(err, "migrating db schema")
			}
		},
	}

	prepareCmd.AddCommand(dbCmd)

	var testDataCmd = &cobra.Command{
		Use:   "test-data",
		Short: "Generate test data",
		Long:  ``,
		Run: func(cmd *cobra.Command, arg []string) {
			connString := fmt.Sprintf(
				"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
				viper.GetString("mysql.username"),
				viper.GetString("mysql.password"),
				viper.GetString("mysql.address"),
				viper.GetString("mysql.database"),
			)
			ds, err := datastore.New("gorm-mysql", connString)
			if err != nil {
				initFatal(err, "creating db connection")
			}
			svc, err := server.NewService(server.ServiceConfig{Datastore: ds})
			if err != nil {
				initFatal(err, "creating new service")
			}
			var (
				name     = "admin"
				username = "admin"
				password = "secret"
				email    = "admin@kolide.co"
				enabled  = true
				isAdmin  = true
			)
			admin := kolide.UserPayload{
				Name:     &name,
				Username: &username,
				Password: &password,
				Email:    &email,
				Enabled:  &enabled,
				Admin:    &isAdmin,
			}
			_, err = svc.NewUser(context.Background(), admin)
			if err != nil {
				initFatal(err, "saving new user")
			}
		},
	}

	prepareCmd.AddCommand(testDataCmd)

	return prepareCmd

}
