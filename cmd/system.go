package cmd

import (
	"fmt"
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
	cmd.AddCommand(newSystemInfoCmd())
	cmd.AddCommand(newSystemHealthCmd())
	cmd.AddCommand(newSystemConfigCmd())
	cmd.AddCommand(newSystemGCCmd())

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

func newSystemInfoCmd() *cobra.Command {
	var withStorage bool
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show system information",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			svc := harbor.NewSystemService(client)
			info, err := svc.GetInfo(withStorage)
			if err != nil {
				return fmt.Errorf("failed to get system info: %w", err)
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(info)
			case "yaml":
				return output.YAML(info)
			default:
				table := output.Table()
				table.Append([]string{"FIELD", "VALUE"})
				table.Append([]string{"Harbor Version", info.HarborVersion})
				table.Append([]string{"Registry", info.RegistryURL})
				table.Render()
				return nil
			}
		},
	}
	cmd.Flags().BoolVar(&withStorage, "with-storage", false, "Include storage information")
	return cmd
}

func newSystemHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check system health",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			if err := client.CheckHealth(); err != nil {
				return err
			}
			output.Success("Harbor is healthy")
			return nil
		},
	}
}

func newSystemConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",

		Short: "Manage Harbor system configuration",
	}
	cmd.AddCommand(newSystemConfigGetCmd())
	cmd.AddCommand(newSystemConfigSetCmd())

	return cmd
}

func newSystemConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get [key]",

		Short: "Get Harbor system configuration",

		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			svc := harbor.NewConfigService(client)
			cfg, err := svc.Get()
			if err != nil {
				return fmt.Errorf("failed to get configuration: %w", err)
			}
			if len(args) == 1 {
				val, ok := cfg[args[0]]
				if !ok {
					return fmt.Errorf("configuration key '%s' not found", args[0])
				}
				fmt.Printf("%s: %v\n", args[0], val)
				return nil
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(cfg)
			default:
				return output.YAML(cfg)
			}
		},
	}
}

func parseConfigValue(val string) interface{} {
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	return val
}

func newSystemConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set Harbor system configuration",
		Args:  requireArgs(2, "requires <key> and <value>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			svc := harbor.NewConfigService(client)
			cfg := map[string]interface{}{args[0]: parseConfigValue(args[1])}
			if err := svc.Update(cfg); err != nil {
				return fmt.Errorf("failed to update configuration: %w", err)
			}
			output.Success("Updated %s", args[0])

			return nil
		},
	}
}

func newSystemGCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gc",
		Short: "Manage garbage collection",
	}
	cmd.AddCommand(newSystemGCScheduleCmd())
	cmd.AddCommand(newSystemGCHistoryCmd())
	cmd.AddCommand(newSystemGCStatusCmd())
	return cmd
}

func newSystemGCScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule",
		Short: "Schedule garbage collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewSystemService(client)
			if err := svc.ScheduleGC(); err != nil {
				return err
			}
			output.Success("Garbage collection scheduled")
			return nil
		},
	}
}

func newSystemGCHistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Show garbage collection history",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewSystemService(client)
			history, err := svc.GetGCHistory()
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(history)
			case "yaml":
				return output.YAML(history)
			default:
				table := output.Table()
				table.Append([]string{"ID", "STATUS", "START", "END"})
				for _, h := range history {
					table.Append([]string{
						strconv.FormatInt(h.ID, 10),
						h.JobStatus,
						h.CreationTime.Format("2006-01-02 15:04:05"),
						h.UpdateTime.Format("2006-01-02 15:04:05"),
					})
				}
				table.Render()
				return nil
			}
		},
	}
}

func newSystemGCStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <id>",
		Short: "Get garbage collection job status",
		Args:  requireArgs(1, "requires <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewSystemService(client)
			gc, err := svc.GetGC(id)
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(gc)
			case "yaml":
				return output.YAML(gc)
			default:
				table := output.Table()
				table.Append([]string{"FIELD", "VALUE"})
				table.Append([]string{"ID", strconv.FormatInt(gc.ID, 10)})
				table.Append([]string{"STATUS", gc.JobStatus})
				table.Append([]string{"START", gc.CreationTime.Format("2006-01-02 15:04:05")})
				table.Append([]string{"END", gc.UpdateTime.Format("2006-01-02 15:04:05")})
				table.Render()
				return nil
			}
		},
	}
}
