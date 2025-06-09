package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewRepositoryCmd creates the repo command
func NewRepositoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repo",
		Aliases: []string{"repository"},
		Short:   "Manage repositories",
		Long:    `Manage repositories within Harbor projects.`,
	}

	cmd.AddCommand(newRepoListCmd())

	return cmd
}

func newRepoListCmd() *cobra.Command {
	var (
		page     int
		pageSize int
		filter   string
		sortBy   string
		detail   bool
	)

	cmd := &cobra.Command{
		Use:   "list <project>",
		Short: "List repositories",
		Long:  `List repositories in a project.`,
		Args:  cobra.ExactArgs(1),
		Example: `  # List repositories
  hrbcli repo list myproject

  # List with details
  hrbcli repo list myproject --detail

  # Filter by name
  hrbcli repo list myproject --filter "app*"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			repoSvc := harbor.NewRepositoryService(client)

			opts := &api.ListOptions{
				Page:     page,
				PageSize: pageSize,
				Query:    filter,
				Sort:     sortBy,
			}

			repos, err := repoSvc.List(project, opts)
			if err != nil {
				return fmt.Errorf("failed to list repositories: %w", err)
			}

			if len(repos) == 0 {
				output.Info("No repositories found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(repos)
			case "yaml":
				return output.YAML(repos)
			default:
				table := output.Table()
				headers := []string{"NAME", "ARTIFACTS", "PULLS"}
				if detail {
					headers = append(headers, "CREATED", "UPDATED")
				}
				table.Append(headers)

				for _, r := range repos {
					row := []string{
						r.Name,
						strconv.FormatInt(r.ArtifactCount, 10),
						strconv.FormatInt(r.PullCount, 10),
					}
					if detail {
						row = append(row,
							r.CreationTime.Format("2006-01-02"),
							r.UpdateTime.Format("2006-01-02"),
						)
					}
					table.Append(row)
				}

				table.Render()
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter by name")
	cmd.Flags().StringVar(&sortBy, "sort", "", "Sort by field")
	cmd.Flags().BoolVar(&detail, "detail", false, "Show detailed information")

	return cmd
}
