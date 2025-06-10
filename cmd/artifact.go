package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewArtifactCmd creates the artifact command
func NewArtifactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Manage artifacts",
		Long:  `Manage artifacts in Harbor.`,
	}

	cmd.AddCommand(newArtifactListCmd())
	cmd.AddCommand(newArtifactGetCmd())
	cmd.AddCommand(newArtifactScanCmd())
	cmd.AddCommand(newArtifactVulnCmd())
	cmd.AddCommand(newArtifactSbomCmd())

	return cmd
}

func newArtifactScanCmd() *cobra.Command {
	var (
		scanType string
		wait     bool
	)

	cmd := &cobra.Command{
		Use:   "scan <project>/<repository>[:tag|@digest]",
		Short: "Scan an image",
		Long:  `Trigger vulnerability scan for a specific image in Harbor.`,
		Args:  requireArgs(1, "requires <project>/<repository>[:tag|@digest]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, ref, err := parseArtifactRef(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			artSvc := harbor.NewArtifactService(client)
			if err := artSvc.Scan(project, repo, ref, scanType); err != nil {
				return fmt.Errorf("failed to scan artifact: %w", err)
			}

			output.Success("Scan triggered for %s/%s:%s", project, repo, ref)

			if wait {
				ctx := cmd.Context()
				for {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}

					art, err := artSvc.GetWithOptions(project, repo, ref, &api.ArtifactGetOptions{WithScanOverview: true})
					if err != nil {
						return fmt.Errorf("failed to get scan status: %w", err)
					}

					done := true
					for _, ov := range art.ScanOverview {
						status := strings.ToLower(ov.ScanStatus)
						if status != "success" && status != "finished" {
							done = false
							break
						}
					}

					if done {
						output.Success("Scan completed for %s/%s:%s", project, repo, ref)
						return nil
					}

					time.Sleep(2 * time.Second)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&scanType, "scan-type", "", "Scan type (vulnerability|sbom)")
	cmd.Flags().BoolVar(&wait, "wait", false, "Wait for scan to complete")

	return cmd
}

func parseArtifactRef(input string) (string, string, string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid reference format")
	}
	project := parts[0]
	remainder := parts[1]

	repo := remainder
	ref := "latest"

	if idx := strings.Index(remainder, "@"); idx != -1 {
		repo = remainder[:idx]
		ref = remainder[idx+1:]
	} else if idx := strings.LastIndex(remainder, ":"); idx != -1 {
		repo = remainder[:idx]
		ref = remainder[idx+1:]
	}

	if repo == "" || ref == "" {
		return "", "", "", fmt.Errorf("invalid reference format")
	}

	return project, repo, ref, nil
}

func parseProjectRepo(input string) (string, string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "", "", fmt.Errorf("invalid reference format")
	}
	project := parts[0]
	repo := ""
	if len(parts) == 2 {
		repo = parts[1]
		if repo == "" {
			return "", "", fmt.Errorf("invalid repository reference")
		}
	}
	return project, repo, nil
}

func newArtifactListCmd() *cobra.Command {
	var (
		page             int
		pageSize         int
		withLabel        bool
		withScanOverview bool
		detail           bool
	)

	cmd := &cobra.Command{
		Use:   "list <project>[/<repository>]",
		Short: "List artifacts",
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

			artSvc := harbor.NewArtifactService(client)

			opts := &api.ArtifactListOptions{
				Page:             page,
				PageSize:         pageSize,
				WithTag:          true,
				WithLabel:        withLabel,
				WithSignature:    true,
				WithScanOverview: withScanOverview,
			}

			printArtifacts := func(repoName string, arts []*api.Artifact) error {
				if len(arts) == 0 {
					return nil
				}

				switch output.GetFormat() {
				case "json":
					return output.JSON(arts)
				case "yaml":
					return output.YAML(arts)
				default:
					table := output.Table()
					headers := []string{"REPOSITORY", "DIGEST", "TAGS"}
					if detail {
						headers = append(headers, "SIZE", "ARCH", "SIGNED")
					}
					table.Append(headers)
					for _, a := range arts {
						if detail && (a.ExtraAttrs == nil || a.ExtraAttrs.Architecture == "") {
							if art, err := artSvc.Get(project, repoName, a.Digest); err == nil && art.ExtraAttrs != nil {
								a.ExtraAttrs = art.ExtraAttrs
							}
						}
						tags := make([]string, len(a.Tags))
						for i, t := range a.Tags {
							tags[i] = t.Name
						}
						row := []string{repoName, output.Truncate(a.Digest, 13), strings.Join(tags, ",")}
						if detail {
							arch := ""
							if a.ExtraAttrs != nil {
								arch = a.ExtraAttrs.Architecture
							}
							signed := "no"
							if len(a.Signatures) > 0 {
								signed = "yes"
							}
							row = append(row,
								harbor.FormatStorageSize(a.Size),
								arch,
								signed,
							)
						}
						table.Append(row)
					}
					table.Render()
					return nil
				}
			}

			if repo != "" {
				arts, err := artSvc.List(project, repo, opts)
				if err != nil {
					return fmt.Errorf("failed to list artifacts: %w", err)
				}
				return printArtifacts(repo, arts)
			}

			repoSvc := harbor.NewRepositoryService(client)
			repos, err := repoSvc.List(project, nil)
			if err != nil {
				return fmt.Errorf("failed to list repositories: %w", err)
			}

			for _, r := range repos {

				repoName := strings.TrimPrefix(r.Name, project+"/")
				arts, err := artSvc.List(project, repoName, opts)

				if err != nil {
					return fmt.Errorf("failed to list artifacts for %s: %w", r.Name, err)
				}
				if err := printArtifacts(r.Name, arts); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().BoolVar(&withLabel, "with-label", false, "Include labels")
	cmd.Flags().BoolVar(&withScanOverview, "with-scan-overview", false, "Include scan overview")
	cmd.Flags().BoolVar(&detail, "detail", false, "Show detailed information")
	return cmd
}

func newArtifactVulnCmd() *cobra.Command {
	var (
		severity string
		summary  bool
	)

	cmd := &cobra.Command{
		Use:   "vulnerabilities <project>/<repository>[:tag|@digest]",
		Short: "Show vulnerability report",
		Args:  requireArgs(1, "requires <project>/<repository>[:tag|@digest]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, ref, err := parseArtifactRef(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			artSvc := harbor.NewArtifactService(client)

			if summary {
				art, err := artSvc.GetWithOptions(project, repo, ref, &api.ArtifactGetOptions{WithScanOverview: true})
				if err != nil {
					return fmt.Errorf("failed to get summary: %w", err)
				}
				if len(art.ScanOverview) == 0 {
					output.Info("No vulnerability summary available")
					return nil
				}
				switch output.GetFormat() {
				case "json":
					return output.JSON(art.ScanOverview)
				case "yaml":
					return output.YAML(art.ScanOverview)
				default:
					table := output.Table()

					table.Append([]string{"SCANNER", "STATUS", "SEVERITY", "TOTAL", "CRITICAL", "HIGH", "MEDIUM", "LOW"})
					for name, ov := range art.ScanOverview {
						sum := ov.Summary.Summary
						table.Append([]string{
							name,
							ov.ScanStatus,
							ov.Severity,
							fmt.Sprintf("%d", ov.Summary.Total),
							fmt.Sprintf("%d", sum["Critical"]),
							fmt.Sprintf("%d", sum["High"]),
							fmt.Sprintf("%d", sum["Medium"]),
							fmt.Sprintf("%d", sum["Low"]),
						})

					}
					table.Render()
					return nil
				}
			}

			report, err := artSvc.Vulnerabilities(project, repo, ref)
			if err != nil {
				return fmt.Errorf("failed to get vulnerabilities: %w", err)
			}

			if report == nil {
				output.Info("No vulnerabilities found")
				return nil
			}

			if len(report.Vulnerabilities) == 0 {
				if report.Summary.Total == 0 {
					output.Info("No vulnerabilities found")
					return nil
				}
				output.Info("%d vulnerabilities found (summary only)", report.Summary.Total)
				return nil
			}

			sevOrder := map[string]int{"none": 0, "negligible": 1, "low": 2, "medium": 3, "high": 4, "critical": 5}
			var vulns []api.VulnerabilityItem
			normalized := strings.ToLower(severity)
			for _, v := range report.Vulnerabilities {
				if normalized != "" {
					if sevOrder[strings.ToLower(v.Severity)] < sevOrder[normalized] {
						continue
					}
				}
				vulns = append(vulns, v)
			}

			if len(vulns) == 0 {
				output.Info("No vulnerabilities found")
			} else {
				switch output.GetFormat() {
				case "json":
					return output.JSON(vulns)
				case "yaml":
					return output.YAML(vulns)
				default:
					table := output.Table()
					table.Append([]string{"SEVERITY", "CVE", "PACKAGE", "VERSION", "FIXED VERSION"})
					for _, v := range vulns {
						table.Append([]string{v.Severity, v.CVEID, v.Package, v.Version, v.FixedVersion})
					}
					table.Render()
				}
			}

			if normalized != "" && len(vulns) > 0 {
				return fmt.Errorf("vulnerabilities with severity >= %s found", severity)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&severity, "severity", "", "Fail if vulnerabilities of this severity or higher are found")
	cmd.Flags().BoolVar(&summary, "summary", false, "Show vulnerability summary instead of detailed report")

	return cmd
}

func newArtifactSbomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sbom <project>/<repository>[:tag|@digest]",
		Short: "Show SBOM report",
		Args:  requireArgs(1, "requires <project>/<repository>[:tag|@digest]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, ref, err := parseArtifactRef(args[0])
			if err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			artSvc := harbor.NewArtifactService(client)
			report, err := artSvc.SBOM(project, repo, ref)
			if err != nil {
				return fmt.Errorf("failed to get SBOM: %w", err)
			}

			if len(report) == 0 {
				output.Info("No SBOM data found")
				return nil
			}

			switch output.GetFormat() {
			case "yaml":
				return output.YAML(report)
			default:
				return output.JSON(report)
			}
		},
	}

	return cmd
}

func newArtifactGetCmd() *cobra.Command {
	var showTags bool
	cmd := &cobra.Command{
		Use:   "get <project>/<repository>[:tag|@digest]",
		Short: "Get artifact details",
		Args:  requireArgs(1, "requires <project>/<repository>[:tag|@digest]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, repo, ref, err := parseArtifactRef(args[0])
			if err != nil {
				return err
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			artSvc := harbor.NewArtifactService(client)
			artifact, err := artSvc.Get(project, repo, ref)
			if err != nil {
				return fmt.Errorf("failed to get artifact: %w", err)
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(artifact)
			case "yaml":
				return output.YAML(artifact)
			default:
				output.Info("Artifact: %s", output.Bold(artifact.Digest))
				fmt.Printf("Size:         %s\n", harbor.FormatStorageSize(artifact.Size))
				if artifact.ExtraAttrs != nil {
					if artifact.ExtraAttrs.Architecture != "" {
						fmt.Printf("Architecture: %s\n", artifact.ExtraAttrs.Architecture)
					}
					if artifact.ExtraAttrs.OS != "" {
						fmt.Printf("OS:           %s\n", artifact.ExtraAttrs.OS)
					}
				}
				signed := "no"
				if len(artifact.Signatures) > 0 {
					signed = "yes"
				}
				fmt.Printf("Signed:       %s\n", signed)
				if showTags && len(artifact.Tags) > 0 {
					names := make([]string, len(artifact.Tags))
					for i, t := range artifact.Tags {
						names[i] = t.Name
					}
					fmt.Printf("Tags:         %s\n", strings.Join(names, ", "))
				}
				return nil
			}
		},
	}
	cmd.Flags().BoolVar(&showTags, "tags", false, "Display tag names")
	return cmd
}
