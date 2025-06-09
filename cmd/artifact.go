package cmd

import (
	"fmt"
	"strings"

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

	cmd.AddCommand(newArtifactScanCmd())

	return cmd
}

func newArtifactScanCmd() *cobra.Command {
	var scanType string

	cmd := &cobra.Command{
		Use:   "scan <project>/<repository>[:tag|@digest]",
		Short: "Scan an image",
		Long:  `Trigger vulnerability scan for a specific image in Harbor.`,
		Args:  cobra.ExactArgs(1),
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
			return nil
		},
	}

	cmd.Flags().StringVar(&scanType, "scan-type", "", "Scan type (vulnerability|sbom)")

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
