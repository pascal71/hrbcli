package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "hrbcli",
	Short: "Harbor CLI - Manage Harbor from the command line",
	Long:  `Harbor CLI (hrbcli) is a command-line interface for Harbor registry.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	// Add commands here
}

func initConfig() {
	viper.SetEnvPrefix("HARBOR")
	viper.AutomaticEnv()
}
