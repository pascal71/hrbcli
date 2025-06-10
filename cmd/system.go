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
