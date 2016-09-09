package cli

import (
	"fmt"

	"github.com/kolide/kolide-ose/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: This should probably be killed once we are done developing the config
// patterns (or turned into something more user friendly)

func createTestCmd(confManager config.ConfigManager) *cobra.Command {

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Test",
		Long:  `Subcommand for debug testing`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Config",
		Long:  `Spit out the config values`,
		Run: func(cmd *cobra.Command, args []string) {
			viper.Debug()
			fmt.Println(viper.AllSettings())
			fmt.Println(viper.AllKeys())
			fmt.Println(confManager.LoadConfig())
		},
	}

	testCmd.AddCommand(configCmd)

	return testCmd
}
