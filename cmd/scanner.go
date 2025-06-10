package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewScannerCmd creates the scanner command
func NewScannerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scanner",
		Short: "Scanner related commands",
	}

	cmd.AddCommand(newScannerRunningCmd())

	cmd.AddCommand(newScannerScanCmd())

	return cmd
}

func newScannerRunningCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "running <project>[/<repository>]",
		Short: "Show running scans",
		Args:  requireArgs(1, "requires <project>[/<repository>]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, err := parseProjectRepo(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			repoSvc := harbor.NewRepositoryService(client)
			artSvc := harbor.NewArtifactService(client)

			var repos []string
			if repo != "" {
				repos = []string{repo}
			} else {
				list, err := repoSvc.List(project, nil)
				if err != nil {
					return fmt.Errorf("failed to list repositories: %w", err)
				}
				for _, r := range list {
					repos = append(repos, strings.TrimPrefix(r.Name, project+"/"))
				}
			}

			type entry struct {
				Repository string `json:"repository"`
				Digest     string `json:"digest"`
				Tags       string `json:"tags"`
				Status     string `json:"status"`
			}
			var running []entry

			for _, r := range repos {
				arts, err := artSvc.List(project, r, &api.ArtifactListOptions{WithTag: true, WithScanOverview: true})
				if err != nil {
					return fmt.Errorf("failed to list artifacts for %s: %w", r, err)
				}
				for _, a := range arts {
					status := ""
					for _, ov := range a.ScanOverview {
						status = ov.ScanStatus
						break
					}
					if status != "" && strings.ToLower(status) != "success" && strings.ToLower(status) != "finished" {
						tags := make([]string, len(a.Tags))
						for i, t := range a.Tags {
							tags[i] = t.Name
						}
						running = append(running, entry{
							Repository: r,
							Digest:     output.Truncate(a.Digest, 13),
							Tags:       strings.Join(tags, ","),
							Status:     status,
						})
					}
				}
			}

			if len(running) == 0 {
				output.Info("No running scans")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(running)
			case "yaml":
				return output.YAML(running)
			default:
				table := output.Table()
				table.Append([]string{"REPOSITORY", "DIGEST", "TAGS", "STATUS"})
				for _, e := range running {
					table.Append([]string{e.Repository, e.Digest, e.Tags, e.Status})
				}
				table.Render()
				return nil
			}
		},
	}

	return cmd
}

func newScannerScanCmd() *cobra.Command {
	var scanType string

	cmd := &cobra.Command{
		Use:   "scan <project>[/<repository>]",
		Short: "Trigger scan for artifacts",
		Args:  requireArgs(1, "requires <project>[/<repository>]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, err := parseProjectRepo(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			repoSvc := harbor.NewRepositoryService(client)
			artSvc := harbor.NewArtifactService(client)

			var repos []string
			if repo != "" {
				repos = []string{repo}
			} else {
				list, err := repoSvc.List(project, nil)
				if err != nil {
					return fmt.Errorf("failed to list repositories: %w", err)
				}
				for _, r := range list {
					repos = append(repos, strings.TrimPrefix(r.Name, project+"/"))
				}
			}

			for _, r := range repos {
				arts, err := artSvc.List(project, r, nil)
				if err != nil {
					return fmt.Errorf("failed to list artifacts for %s: %w", r, err)
				}
				for _, a := range arts {
					if err := artSvc.Scan(project, r, a.Digest, scanType); err != nil {
						output.Warning("Failed to scan %s/%s@%s: %v", project, r, output.Truncate(a.Digest, 13), err)
					} else {
						output.Success("Scan triggered for %s/%s@%s", project, r, output.Truncate(a.Digest, 13))
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&scanType, "scan-type", "", "Scan type (vulnerability|sbom)")

	return cmd
}
