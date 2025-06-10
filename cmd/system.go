package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewSystemCmd creates the system command
func NewSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System administration commands",
		Long:  `Manage Harbor system operations such as statistics.`,
	}

	cmd.AddCommand(newSystemStatisticsCmd())
	cmd.AddCommand(newSystemConfigCmd())

	return cmd
}

func newSystemStatisticsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Show Harbor statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			sysSvc := harbor.NewSystemService(client)

			stats, err := sysSvc.GetStatistics()
			if err != nil {
				return fmt.Errorf("failed to get statistics: %w", err)
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(stats)
			case "yaml":
				return output.YAML(stats)
			default:
				table := output.Table()
				table.Append(
					[]string{
						"PRIVATE PROJECTS",
						"PUBLIC PROJECTS",
						"TOTAL PROJECTS",
						"PRIVATE REPOS",
						"PUBLIC REPOS",
						"TOTAL REPOS",
						"STORAGE",
					},
				)
				table.Append([]string{
					strconv.FormatInt(stats.PrivateProjectCount, 10),
					strconv.FormatInt(stats.PublicProjectCount, 10),
					strconv.FormatInt(stats.TotalProjectCount, 10),
					strconv.FormatInt(stats.PrivateRepoCount, 10),
					strconv.FormatInt(stats.PublicRepoCount, 10),
					strconv.FormatInt(stats.TotalRepoCount, 10),
					harbor.FormatStorageSize(stats.TotalStorageConsumption),
				})
				table.Render()
				return nil
			}
		},
	}

	return cmd
}

func newSystemConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Harbor configuration",
	}

	cmd.AddCommand(newSystemConfigGetCmd())
	cmd.AddCommand(newSystemConfigSetCmd())

	return cmd
}

func newSystemConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get Harbor configuration",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			sysSvc := harbor.NewSystemService(client)

			cfg, err := sysSvc.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to get configuration: %w", err)
			}

			if len(args) == 1 {
				key := args[0]
				val, ok := cfg[key]
				if !ok {
					return fmt.Errorf("configuration key '%s' not found", key)
				}
				switch output.GetFormat() {
				case "json":
					return output.JSON(val)
				case "yaml":
					return output.YAML(val)
				default:
					fmt.Printf("%s: %v\n", key, val)
					return nil
				}
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(cfg)
			case "yaml":
				return output.YAML(cfg)
			default:
				table := output.Table()
				table.Append([]string{"KEY", "VALUE"})
				keys := make([]string, 0, len(cfg))
				for k := range cfg {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					table.Append([]string{k, fmt.Sprintf("%v", cfg[k])})
				}
				table.Render()
				return nil
			}
		},
	}
}

func newSystemConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set Harbor configuration value",
		Args:  requireArgs(2, "requires <key> <value>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			var typed interface{} = value
			if value == "true" {
				typed = true
			} else if value == "false" {
				typed = false
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			sysSvc := harbor.NewSystemService(client)
			update := map[string]interface{}{key: typed}
			if err := sysSvc.UpdateConfig(update); err != nil {
				return fmt.Errorf("failed to update configuration: %w", err)
			}

			output.Success("Updated %s", key)
			return nil
		},
	}
}
