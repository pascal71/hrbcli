package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pascal71/hrbcli/pkg/output"
)

var (
	cfgFile      string
	outputFormat string
	debug        bool
	noColor      bool
)

var rootCmd = &cobra.Command{
	Use:   "hrbcli",
	Short: "Harbor CLI - Manage Harbor from the command line",
	Long: `Harbor CLI (hrbcli) is a command-line interface for Harbor registry.

It provides access to Harbor functionality through the Harbor API,
allowing you to manage projects, repositories, users, and more.

Examples:
  # List all projects
  hrbcli project list

  # Create a new project
  hrbcli project create myproject --public

  # List repositories in a project
  hrbcli repo list myproject

  # Get system health status
  hrbcli system health`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip validation for config and completion commands
		if cmd.Name() == "config" || cmd.Name() == "completion" || cmd.Parent().Name() == "config" {
			return nil
		}

		// Validate required configuration
		if viper.GetString("harbor_url") == "" {
			return fmt.Errorf(
				"Harbor URL not configured. Run 'hrbcli config init' or set HARBOR_URL",
			)
		}

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hrbcli.yaml)")
	rootCmd.PersistentFlags().String("harbor-url", "", "Harbor server URL")
	rootCmd.PersistentFlags().String("username", "", "Harbor username")
	rootCmd.PersistentFlags().String("password", "", "Harbor password")
	rootCmd.PersistentFlags().String("api-version", "v2.0", "Harbor API version")
	rootCmd.PersistentFlags().Bool("insecure", false, "Skip TLS certificate verification")
	rootCmd.PersistentFlags().
		StringVarP(&outputFormat, "output", "o", "table", "Output format (table|json|yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	// Bind flags to viper
	viper.BindPFlag("harbor_url", rootCmd.PersistentFlags().Lookup("harbor-url"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("api_version", rootCmd.PersistentFlags().Lookup("api-version"))
	viper.BindPFlag("insecure", rootCmd.PersistentFlags().Lookup("insecure"))
	viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("no_color", rootCmd.PersistentFlags().Lookup("no-color"))

	// Add commands - we'll implement these next
	rootCmd.AddCommand(NewProjectCmd())
	rootCmd.AddCommand(NewRegistryCmd())
	rootCmd.AddCommand(NewArtifactCmd())
	rootCmd.AddCommand(NewUserCmd())
	rootCmd.AddCommand(NewRepositoryCmd())
	rootCmd.AddCommand(NewReplicationCmd())
	rootCmd.AddCommand(NewSystemCmd())
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewCompletionCmd())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search for config in home directory
		configPath := filepath.Join(home, ".hrbcli.yaml")
		viper.SetConfigFile(configPath)
	}

	// Set env prefix
	viper.SetEnvPrefix("HARBOR")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Set output formatter settings
	output.SetFormat(viper.GetString("output_format"))
	output.SetNoColor(viper.GetBool("no_color"))
	output.SetDebug(viper.GetBool("debug"))
}
