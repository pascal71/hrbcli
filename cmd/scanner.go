package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	cmd.AddCommand(newScannerReportsCmd())

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

func newScannerReportsCmd() *cobra.Command {
	var reportType string
	var summary bool
	var outputDir string
	var sortBy string
	var reverse bool

	cmd := &cobra.Command{
		Use:   "reports <project>[/<repository>]",
		Short: "Get reports for artifacts",
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

			if outputDir != "" {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			}

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
				Repository string      `json:"repository"`
				Reference  string      `json:"reference"`
				Report     interface{} `json:"report"`
				Count      int         `json:"count,omitempty"`
				Critical   int         `json:"critical,omitempty"`
				High       int         `json:"high,omitempty"`
				Medium     int         `json:"medium,omitempty"`
				Low        int         `json:"low,omitempty"`
				Total      int         `json:"total,omitempty"`
			}
			var reports []entry

			for _, r := range repos {
				arts, err := artSvc.List(project, r, &api.ArtifactListOptions{WithTag: true, WithScanOverview: summary})
				if err != nil {
					return fmt.Errorf("failed to list artifacts for %s: %w", r, err)
				}
				for _, a := range arts {
					ref := a.Digest
					if len(a.Tags) > 0 {
						ref = a.Tags[0].Name
					}

					if summary {
						if outputDir != "" {
							ext := "json"
							if output.GetFormat() == "yaml" {
								ext = "yaml"
							}
							name := fmt.Sprintf("%s_%s_summary.%s", strings.ReplaceAll(r, "/", "_"), strings.ReplaceAll(strings.ReplaceAll(ref, ":", "_"), "/", "_"), ext)
							path := filepath.Join(outputDir, name)
							if err := output.WriteFile(path, ext, a.ScanOverview); err != nil {
								return fmt.Errorf("failed to write report: %w", err)
							}
							output.Success("Saved report to %s", path)
						} else {
							c := 0
							crit := 0
							high := 0
							med := 0
							low := 0
							for _, ov := range a.ScanOverview {
								c += ov.Summary.Total
								crit += ov.Summary.Summary["Critical"]
								high += ov.Summary.Summary["High"]
								med += ov.Summary.Summary["Medium"]
								low += ov.Summary.Summary["Low"]
							}
							reports = append(reports, entry{Repository: r, Reference: ref, Report: a.ScanOverview, Count: c, Critical: crit, High: high, Medium: med, Low: low, Total: c})
						}
						continue
					}

					if strings.ToLower(reportType) == "sbom" {
						report, err := artSvc.SBOM(project, r, a.Digest)
						if err != nil {
							output.Warning("Failed to get SBOM for %s/%s@%s: %v", project, r, output.Truncate(a.Digest, 13), err)
							continue
						}
						if outputDir != "" {
							ext := "json"
							if output.GetFormat() == "yaml" {
								ext = "yaml"
							}
							name := fmt.Sprintf("%s_%s_sbom.%s", strings.ReplaceAll(r, "/", "_"), strings.ReplaceAll(strings.ReplaceAll(ref, ":", "_"), "/", "_"), ext)
							path := filepath.Join(outputDir, name)
							if err := output.WriteFile(path, ext, report); err != nil {
								return fmt.Errorf("failed to write report: %w", err)
							}
							output.Success("Saved report to %s", path)
						} else {
							reports = append(reports, entry{Repository: r, Reference: ref, Report: report})
						}
					} else {
						report, err := artSvc.Vulnerabilities(project, r, a.Digest)
						if err != nil {
							output.Warning("Failed to get vulnerabilities for %s/%s@%s: %v", project, r, output.Truncate(a.Digest, 13), err)
							continue
						}
						if outputDir != "" {
							ext := "json"
							if output.GetFormat() == "yaml" {
								ext = "yaml"
							}
							name := fmt.Sprintf("%s_%s_vuln.%s", strings.ReplaceAll(r, "/", "_"), strings.ReplaceAll(strings.ReplaceAll(ref, ":", "_"), "/", "_"), ext)
							path := filepath.Join(outputDir, name)
							if err := output.WriteFile(path, ext, report); err != nil {
								return fmt.Errorf("failed to write report: %w", err)
							}
							output.Success("Saved report to %s", path)
						} else {
							count := len(report.Vulnerabilities)
							if count == 0 && report.Summary.Total > 0 {
								count = report.Summary.Total
							}

							crit := report.Summary.Summary["Critical"]
							high := report.Summary.Summary["High"]
							med := report.Summary.Summary["Medium"]
							low := report.Summary.Summary["Low"]

							reports = append(reports, entry{
								Repository: r,
								Reference:  ref,
								Report:     report,
								Count:      count,
								Critical:   crit,
								High:       high,
								Medium:     med,
								Low:        low,
								Total:      report.Summary.Total,
							})
						}
					}
				}
			}

			if outputDir != "" {
				return nil
			}

			if len(reports) == 0 {
				output.Info("No reports found")
				return nil
			}

			sort.SliceStable(reports, func(i, j int) bool {
				field := strings.ToLower(sortBy)
				switch field {
				case "repo", "repository":
					if reverse {
						return reports[i].Repository > reports[j].Repository
					}
					return reports[i].Repository < reports[j].Repository
				case "ref", "reference":
					if reverse {
						return reports[i].Reference > reports[j].Reference
					}
					return reports[i].Reference < reports[j].Reference
				case "vuln":
					if reverse {
						return reports[i].Count < reports[j].Count
					}
					return reports[i].Count > reports[j].Count
				case "crit", "critical":
					if reverse {
						return reports[i].Critical < reports[j].Critical
					}
					return reports[i].Critical > reports[j].Critical
				case "high":
					if reverse {
						return reports[i].High < reports[j].High
					}
					return reports[i].High > reports[j].High
				case "med", "medium":
					if reverse {
						return reports[i].Medium < reports[j].Medium
					}
					return reports[i].Medium > reports[j].Medium
				case "low":
					if reverse {
						return reports[i].Low < reports[j].Low
					}
					return reports[i].Low > reports[j].Low
				case "total":
					if reverse {
						return reports[i].Total < reports[j].Total
					}
					return reports[i].Total > reports[j].Total
				default:
					// default severity sort
					if reports[i].Critical != reports[j].Critical {
						if reverse {
							return reports[i].Critical < reports[j].Critical
						}
						return reports[i].Critical > reports[j].Critical
					}
					if reports[i].High != reports[j].High {
						if reverse {
							return reports[i].High < reports[j].High
						}
						return reports[i].High > reports[j].High
					}
					if reports[i].Medium != reports[j].Medium {
						if reverse {
							return reports[i].Medium < reports[j].Medium
						}
						return reports[i].Medium > reports[j].Medium
					}
					if reports[i].Low != reports[j].Low {
						if reverse {
							return reports[i].Low < reports[j].Low
						}
						return reports[i].Low > reports[j].Low
					}
					if reports[i].Total != reports[j].Total {
						if reverse {
							return reports[i].Total < reports[j].Total
						}
						return reports[i].Total > reports[j].Total
					}
					if reverse {
						return reports[i].Repository > reports[j].Repository
					}
					return reports[i].Repository < reports[j].Repository
				}
			})

			switch output.GetFormat() {
			case "json":
				return output.JSON(reports)
			case "yaml":
				return output.YAML(reports)
			default:
				table := output.Table()
				if summary {

					table.Append([]string{"REPOSITORY", "REFERENCE", "SCANNER", "STATUS", "TOTAL", "CRITICAL", "HIGH", "MEDIUM", "LOW"})

					for _, e := range reports {
						overview, ok := e.Report.(map[string]api.NativeReportSummary)
						if !ok {
							continue
						}
						for name, ov := range overview {

							sum := ov.Summary.Summary
							table.Append([]string{
								e.Repository,
								e.Reference,
								name,
								ov.ScanStatus,
								fmt.Sprintf("%d", ov.Summary.Total),
								fmt.Sprintf("%d", sum["Critical"]),
								fmt.Sprintf("%d", sum["High"]),
								fmt.Sprintf("%d", sum["Medium"]),
								fmt.Sprintf("%d", sum["Low"]),
							})

						}
					}
				} else if strings.ToLower(reportType) == "sbom" {
					table.Append([]string{"REPOSITORY", "REFERENCE", "SBOM"})
					for _, e := range reports {
						table.Append([]string{e.Repository, e.Reference, "available"})
					}
				} else {
					table.Append([]string{"REPOSITORY", "REFERENCE", "VULNERABILITIES"})
					for _, e := range reports {
						rep, ok := e.Report.(*api.VulnerabilityReport)
						if !ok || rep == nil {
							table.Append([]string{e.Repository, e.Reference, ""})
							continue
						}

						count := len(rep.Vulnerabilities)
						if count == 0 && rep.Summary.Total > 0 {
							count = rep.Summary.Total
						}
						table.Append([]string{e.Repository, e.Reference, fmt.Sprintf("%d", count)})

					}
				}
				table.Render()
				return nil
			}
		},
	}

	cmd.Flags().StringVar(&reportType, "type", "vulnerability", "Report type (vulnerability|sbom)")
	cmd.Flags().BoolVar(&summary, "summary", false, "Show summary instead of full report")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory to save reports")
	cmd.Flags().StringVar(&sortBy, "sort", "severity", "Sort by field (severity|crit|high|medium|low|total|vuln|repo|ref)")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "Reverse sort order")

	return cmd
}
