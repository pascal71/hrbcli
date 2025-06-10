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
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Schedule garbage collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			sysSvc := harbor.NewSystemService(client)
			id, err := sysSvc.StartGC()
			if err != nil {
				return fmt.Errorf("failed to schedule GC: %w", err)
			}
			if id > 0 {
				output.Success("GC job %d scheduled", id)
			} else {
				output.Success("GC scheduled")
			}
			return nil
		},
	}
	return cmd
}

func newSystemGCHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show garbage collection history",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			sysSvc := harbor.NewSystemService(client)
			history, err := sysSvc.GCHistory(nil)
			if err != nil {
				return fmt.Errorf("failed to get GC history: %w", err)
			}

			if len(history) == 0 {
				output.Info("No GC history found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(history)
			case "yaml":
				return output.YAML(history)
			default:
				table := output.Table()
				table.Append([]string{"ID", "STATUS", "JOB", "CREATED"})
				for _, h := range history {
					table.Append([]string{
						strconv.FormatInt(h.ID, 10),
						h.JobStatus,
						h.JobName,
						h.CreationTime.Format("2006-01-02 15:04:05"),
					})
				}
				table.Render()
				return nil
			}
		},
	}
	return cmd
}

func newSystemGCStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <id>",
		Short: "Show garbage collection job status",
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
			sysSvc := harbor.NewSystemService(client)
			gc, err := sysSvc.GCStatus(id)
			if err != nil {
				return fmt.Errorf("failed to get GC status: %w", err)
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
				table.Append([]string{"JOB", gc.JobName})
				table.Append([]string{"CREATED", gc.CreationTime.Format("2006-01-02 15:04:05")})
				table.Render()
				return nil
			}
		},
	}
	return cmd
}
