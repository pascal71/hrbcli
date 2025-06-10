package cmd

import (
	"fmt"

	"github.com/pascal71/hrbcli/internal/version"
	"github.com/pascal71/hrbcli/pkg/output"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version information of Harbor CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			ver := struct {
				Version   string `json:"version" yaml:"version"`
				BuildTime string `json:"buildTime" yaml:"buildTime"`
				GoVersion string `json:"goVersion" yaml:"goVersion"`
				Platform  string `json:"platform" yaml:"platform"`
			}{
				Version:   version.Version,
				BuildTime: version.BuildTime,
				GoVersion: version.GoVersion,
				Platform:  version.Platform,
			}

			switch output.GetFormat() {
			case "json":
				output.JSON(ver)
			case "yaml":
				output.YAML(ver)
			default:
				fmt.Printf("Harbor CLI (hrbcli) version %s\n", version.Version)
				fmt.Printf("Build Time: %s\n", version.BuildTime)
				fmt.Printf("Go Version: %s\n", version.GoVersion)
				fmt.Printf("Platform: %s\n", version.Platform)
			}
		},
	}

	return cmd
}
