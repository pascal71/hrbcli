package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/output"
)

// NewSystemCmd creates the system command
func NewSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System administration tasks",
		Long:  `Manage Harbor system level tasks like backups and maintenance.`,
	}

	cmd.AddCommand(newSystemBackupCmd())

	return cmd
}

func newSystemBackupCmd() *cobra.Command {
	var outputDir string
	var dbOnly bool

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup Harbor data",
		Long: `Create a backup of Harbor's database and storage using the community backup script. 
Docker must be installed and the command should be executed on the Harbor host.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = "."
			}
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			scriptURL := "https://raw.githubusercontent.com/spagno/harbor-backup/master/scripts/67558_harbor-backup.sh"
			resp, err := http.Get(scriptURL)
			if err != nil {
				return fmt.Errorf("failed to download backup script: %w", err)
			}
			defer resp.Body.Close()

			scriptPath := filepath.Join(outputDir, "harbor-backup.sh")
			f, err := os.Create(scriptPath)
			if err != nil {
				return fmt.Errorf("failed to create script file: %w", err)
			}
			if _, err := io.Copy(f, resp.Body); err != nil {
				f.Close()
				return fmt.Errorf("failed to save script: %w", err)
			}
			f.Close()

			if err := os.Chmod(scriptPath, 0755); err != nil {
				return fmt.Errorf("failed to make script executable: %w", err)
			}

			cmdArgs := []string{scriptPath}
			if dbOnly {
				cmdArgs = append(cmdArgs, "--dbonly")
			}

			backupCmd := exec.Command("bash", cmdArgs...)
			backupCmd.Stdout = os.Stdout
			backupCmd.Stderr = os.Stderr
			backupCmd.Dir = outputDir

			if err := backupCmd.Run(); err != nil {
				return fmt.Errorf("backup execution failed: %w", err)
			}

			output.Success("Backup files stored in %s", outputDir)
			return nil
		},
	}

	cmd.Flags().StringVar(&outputDir, "dir", ".", "Directory to store the backup archive")
	cmd.Flags().BoolVar(&dbOnly, "db-only", false, "Backup database only")

	return cmd
}
