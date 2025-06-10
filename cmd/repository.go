package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
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
	cmd.AddCommand(newRepoGetCmd())
	cmd.AddCommand(newRepoDeleteCmd())
	cmd.AddCommand(newRepoTagsCmd())

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
		Args:  requireArgs(1, "requires <project>"),
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

func parseRepoRef(input string) (string, string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository reference")
	}
	project := parts[0]
	repo := parts[1]
	if project == "" || repo == "" {
		return "", "", fmt.Errorf("invalid repository reference")
	}
	return project, repo, nil
}

func newRepoGetCmd() *cobra.Command {
	var detail bool

	cmd := &cobra.Command{
		Use:   "get <project>/<repository>",
		Short: "Get repository details",
		Long:  `Get detailed information about a repository.`,
		Args:  requireArgs(1, "requires <project>/<repository>"),
		Example: `  # Get repository details
  hrbcli repo get myproject/myrepo

  # Show timestamps
  hrbcli repo get myproject/myrepo --detail`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, err := parseRepoRef(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			repoSvc := harbor.NewRepositoryService(client)

			repository, err := repoSvc.Get(project, repo)
			if err != nil {
				return fmt.Errorf("failed to get repository: %w", err)
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(repository)
			case "yaml":
				return output.YAML(repository)
			default:
				output.Info("Repository: %s", output.Bold(repository.Name))
				output.Info("")
				fmt.Printf("Artifacts:  %d\n", repository.ArtifactCount)
				fmt.Printf("Pulls:      %d\n", repository.PullCount)
				if detail {
					fmt.Printf("Created:    %s\n", repository.CreationTime.Format("2006-01-02 15:04:05"))
					fmt.Printf("Updated:    %s\n", repository.UpdateTime.Format("2006-01-02 15:04:05"))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&detail, "detail", false, "Show creation and update times")
	return cmd
}

func newRepoDeleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <project>/<repository>[<:tag>|@digest]",
		Short: "Delete repository or tag",
		Args:  requireArgs(1, "requires <project>/<repository>[:tag]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]
			project, repo, ref := "", "", ""
			if strings.ContainsAny(input, ":@") {
				p, r, rr, err := parseArtifactRef(input)
				if err != nil {
					return err
				}
				project, repo, ref = p, r, rr
			} else {
				var err error
				project, repo, err = parseRepoRef(input)
				if err != nil {
					return err
				}
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			repoSvc := harbor.NewRepositoryService(client)
			artSvc := harbor.NewArtifactService(client)

			if ref != "" {
				if err := artSvc.Delete(project, repo, ref); err != nil {
					return fmt.Errorf("failed to delete artifact: %w", err)
				}
				output.Success("Deleted %s/%s:%s", project, repo, ref)
				return nil
			}

			if !force {
				prompt := promptui.Prompt{Label: fmt.Sprintf("Delete repository '%s/%s'", project, repo), IsConfirm: true}
				result, err := prompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					output.Info("Deletion cancelled")
					return nil
				}
			}

			if err := repoSvc.Delete(project, repo); err != nil {
				return fmt.Errorf("failed to delete repository: %w", err)
			}
			output.Success("Deleted repository %s/%s", project, repo)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Force deletion without confirmation")
	return cmd
}

func newRepoTagsCmd() *cobra.Command {
	var detail bool
	var filter string
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "tags <project>/<repository>",
		Short: "List tags for repository",
		Args:  requireArgs(1, "requires <project>/<repository>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, err := parseRepoRef(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			repoSvc := harbor.NewRepositoryService(client)
			opts := &api.ListOptions{Query: filter, Page: page, PageSize: pageSize}
			tags, err := repoSvc.ListTags(project, repo, opts)
			if err != nil {
				return fmt.Errorf("failed to list tags: %w", err)
			}

			if len(tags) == 0 {
				output.Info("No tags found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(tags)
			case "yaml":
				return output.YAML(tags)
			default:
				table := output.Table()
				headers := []string{"NAME"}
				if detail {
					headers = append(headers, "IMMUTABLE")
				}
				table.Append(headers)
				for _, t := range tags {
					row := []string{t.Name}
					if detail {
						row = append(row, strconv.FormatBool(t.Immutable))
					}
					table.Append(row)
				}
				table.Render()
				return nil
			}
		},
	}
	cmd.Flags().BoolVar(&detail, "detail", false, "Show immutable flag")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter by tag name")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	return cmd
}
