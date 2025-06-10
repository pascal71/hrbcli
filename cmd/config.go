package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/config"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage hrbcli configuration",
		Long:  `Manage hrbcli configuration including setting Harbor URL, credentials, and other options.`,
	}

	cmd.AddCommand(newConfigInitCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigListCmd())

	return cmd
}

func newConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration interactively",
		Long:  `Initialize Harbor CLI configuration interactively. This will prompt for Harbor URL, credentials, and other settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Welcome to Harbor CLI configuration!")
			output.Info("")

			// Harbor URL
			urlPrompt := promptui.Prompt{
				Label:   "Harbor URL",
				Default: "https://harbor.example.com",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("Harbor URL is required")
					}
					if !strings.HasPrefix(input, "http://") &&
						!strings.HasPrefix(input, "https://") {
						return fmt.Errorf("Harbor URL must start with http:// or https://")
					}
					return nil
				},
			}
			harborURL, err := urlPrompt.Run()
			if err != nil {
				return err
			}

			// Username
			userPrompt := promptui.Prompt{
				Label:   "Username",
				Default: "admin",
			}
			username, err := userPrompt.Run()
			if err != nil {
				return err
			}

			// Password
			fmt.Print("Password: ")
			passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				return err
			}
			password := string(passwordBytes)

			// API Version
			apiPrompt := promptui.Select{
				Label: "API Version",
				Items: []string{"v2.0", "v1.0"},
			}
			_, apiVersion, err := apiPrompt.Run()
			if err != nil {
				return err
			}

			// Output Format
			formatPrompt := promptui.Select{
				Label: "Default Output Format",
				Items: []string{"table", "json", "yaml"},
			}
			_, outputFormat, err := formatPrompt.Run()
			if err != nil {
				return err
			}

			// Insecure
			insecurePrompt := promptui.Select{
				Label: "Skip TLS Verification (only for testing)",
				Items: []string{"false", "true"},
			}
			_, insecureStr, err := insecurePrompt.Run()
			if err != nil {
				return err
			}
			insecure := insecureStr == "true"

			// Create config
			cfg := &config.Config{
				HarborURL:    harborURL,
				Username:     username,
				Password:     password,
				APIVersion:   apiVersion,
				OutputFormat: outputFormat,
				Insecure:     insecure,
			}

			// Verify configuration
			verifyPrompt := promptui.Select{
				Label: "Verify configuration by connecting to Harbor",
				Items: []string{"yes", "no"},
			}
			_, verify, err := verifyPrompt.Run()
			if err != nil {
				return err
			}

			if verify == "yes" {
				output.Info("Testing connection to Harbor...")

				httpClient := &http.Client{Timeout: 10 * time.Second}
				if insecure {
					httpClient.Transport = &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					}
				}

				client := &api.Client{
					BaseURL:    strings.TrimRight(harborURL, "/"),
					Username:   username,
					Password:   password,
					APIVersion: apiVersion,
					HTTPClient: httpClient,
				}

				if err := client.CheckHealth(); err != nil {
					output.Error("Connection failed: %v", err)
				} else {
					output.Success("Successfully connected to Harbor!")
				}
			}

			// Save configuration
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			output.Success("Configuration saved to %s", config.GetConfigPath())
			output.Info("")
			output.Info("You can now use Harbor CLI!")
			output.Info("Try 'hrbcli project list' to list your projects.")

			// Save password to environment variable hint
			if password != "" {
				output.Info("")
				output.Warning(
					"For security, consider using HARBOR_PASSWORD environment variable instead of storing the password.",
				)
			}

			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value. Available keys:
  - harbor_url: Harbor server URL
  - username: Harbor username
  - api_version: Harbor API version (v1.0, v2.0)
  - output_format: Default output format (table, json, yaml)
  - insecure: Skip TLS verification (true, false)
  - default_project: Default project name
  - no_color: Disable colored output (true, false)
  - debug: Enable debug output (true, false)`,
		Args: requireArgs(2, "requires <key> and <value>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// Convert string to appropriate type
			var typedValue interface{}
			switch key {
			case "insecure", "no_color", "debug":
				if value == "true" {
					typedValue = true
				} else if value == "false" {
					typedValue = false
				} else {
					return fmt.Errorf("boolean value must be 'true' or 'false'")
				}
			default:
				typedValue = value
			}

			if err := config.Set(key, typedValue); err != nil {
				return err
			}

			output.Success("Set %s = %v", key, value)
			return nil
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		Long:  `Get a configuration value. If no key is specified, all values are shown.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return listConfig()
			}

			key := args[0]
			value := config.Get(key)
			if value == nil {
				return fmt.Errorf("configuration key '%s' not found", key)
			}

			// Don't show password
			if key == "password" && value != "" {
				value = "********"
			}

			fmt.Printf("%s: %v\n", key, value)
			return nil
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		Long:  `List all configuration values.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listConfig()
		},
	}
}

func listConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Mask password
	if cfg.Password != "" {
		cfg.Password = "********"
	}

	switch output.GetFormat() {
	case "json":
		return output.JSON(cfg)
	case "yaml":
		return output.YAML(cfg)
	default:
		output.Info("Harbor CLI Configuration:")
		output.Info("")
		fmt.Printf("Harbor URL:      %s\n", cfg.HarborURL)
		fmt.Printf("Username:        %s\n", cfg.Username)
		fmt.Printf("Password:        %s\n", cfg.Password)
		fmt.Printf("API Version:     %s\n", cfg.APIVersion)
		fmt.Printf("Insecure:        %v\n", cfg.Insecure)
		fmt.Printf("Output Format:   %s\n", cfg.OutputFormat)
		fmt.Printf("Default Project: %s\n", cfg.DefaultProject)
		fmt.Printf("No Color:        %v\n", cfg.NoColor)
		fmt.Printf("Debug:           %v\n", cfg.Debug)
		output.Info("")
		output.Info("Config file: %s", config.GetConfigPath())

		// Check for environment variables
		if os.Getenv("HARBOR_PASSWORD") != "" {
			output.Info("")
			output.Info("Note: HARBOR_PASSWORD environment variable is set")
		}
	}

	return nil
}
