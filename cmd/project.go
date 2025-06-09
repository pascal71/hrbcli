package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewProjectCmd creates the project command
func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"proj"},
		Short:   "Manage Harbor projects",
		Long:    `Manage Harbor projects including create, list, update, and delete operations.`,
	}

	cmd.AddCommand(newProjectListCmd())
	cmd.AddCommand(newProjectCreateCmd())
	cmd.AddCommand(newProjectGetCmd())
	cmd.AddCommand(newProjectUpdateCmd())
	cmd.AddCommand(newProjectDeleteCmd())
	cmd.AddCommand(newProjectExistsCmd())

	return cmd
}

func newProjectListCmd() *cobra.Command {
	var (
		page     int
		pageSize int
		query    string
		sort     string
		detail   bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Long:  `List all projects accessible to the current user.`,
		Example: `  # List all projects
  hrbcli project list

  # List projects with details
  hrbcli project list --detail

  # Search projects by name
  hrbcli project list --query "name=~prod"

  # List with pagination
  hrbcli project list --page 2 --page-size 20`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			projectSvc := harbor.NewProjectService(client)

			opts := &api.ListOptions{
				Page:     page,
				PageSize: pageSize,
				Query:    query,
				Sort:     sort,
			}

			projects, err := projectSvc.List(opts)
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}

			if len(projects) == 0 {
				output.Info("No projects found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(projects)
			case "yaml":
				return output.YAML(projects)
			default:
				table := output.Table()
				headers := []string{"NAME", "PUBLIC", "REPOS", "OWNER", "CREATED"}
				if detail {
					headers = append(headers, "STORAGE LIMIT", "STORAGE USED")
				}
				// Add headers as first row
				table.Append(headers)

				for _, p := range projects {
					row := []string{
						p.Name,
						strconv.FormatBool(p.Public),
						strconv.FormatInt(p.RepoCount, 10),
						p.OwnerName,
						p.CreationTime.Format("2006-01-02"),
					}

					if detail {
						// Get project summary for quota info
						summary, err := projectSvc.GetSummary(p.Name)
						if err == nil && summary.Quota != nil {
							storageLimit := output.FormatSize(summary.Quota.Hard.Storage)
							storageUsed := output.FormatSize(summary.Quota.Used.Storage)
							row = append(row, storageLimit, storageUsed)
						} else {
							row = append(row, "N/A", "N/A")
						}
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
	cmd.Flags().StringVar(&query, "query", "", "Query string (e.g., 'name=~prod')")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort by field")
	cmd.Flags().BoolVar(&detail, "detail", false, "Show detailed information")

	return cmd
}

func newProjectCreateCmd() *cobra.Command {
	var (
		public       bool
		storageLimit string
		memberLimit  int64
		enableTrust  bool
		preventVul   bool
		severity     string
		autoScan     bool
		reuseSysCVE  bool
		proxyCache   bool
		registryID   int64
		proxySpeedKB int64
		registryName string
	)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new project",
		Long:  `Create a new project in Harbor with the specified name and configuration.`,
		Example: `  # Create a simple private project
  hrbcli project create myproject

  # Create a public project
  hrbcli project create myproject --public

  # Create with storage quota
  hrbcli project create myproject --storage-limit 10G

  # Create with security settings
  hrbcli project create secure-proj --enable-content-trust --prevent-vulnerable --auto-scan`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			// Validate project name
			if err := validateProjectName(projectName); err != nil {
				return err
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			projectSvc := harbor.NewProjectService(client)

			// Check if project already exists
			exists, err := projectSvc.Exists(projectName)
			if err != nil {
				return fmt.Errorf("failed to check project existence: %w", err)
			}
			if exists {
				return fmt.Errorf("project '%s' already exists", projectName)
			}

			// Parse storage limit
			var storageLimitBytes int64 = -1
			if storageLimit != "" {
				storageLimitBytes, err = harbor.ParseStorageLimit(storageLimit)
				if err != nil {
					return err
				}
			}

			// Build project request
			req := &api.ProjectReq{
				ProjectName: projectName,
				Public:      &public,
			}

			// Set metadata
			metadata := &api.ProjectMetadata{}
			if enableTrust {
				metadata.EnableContentTrust = "true"
			}
			if preventVul {
				metadata.PreventVul = "true"
				metadata.Severity = severity
			}
			if autoScan {
				metadata.AutoScan = "true"
			}
			if reuseSysCVE {
				metadata.ReuseSysCVEAllowlist = "true"
			}
			// Look up registry by name if provided
			if proxyCache && registryName != "" && registryID == 0 {
				registrySvc := harbor.NewRegistryService(client)
				registries, err := registrySvc.List(nil)
				if err != nil {
					return fmt.Errorf("failed to list registries: %w", err)
				}
				for _, reg := range registries {
					if reg.Name == registryName {
						registryID = reg.ID
						break
					}
				}
				if registryID == 0 {
					return fmt.Errorf("registry '%s' not found", registryName)
				}
			}

			// Handle proxy cache
			if proxyCache {
				if registryID == 0 {
					return fmt.Errorf("--registry-id is required when creating proxy cache project")
				}
				metadata.ProxySpeedKB = fmt.Sprintf("%d", proxySpeedKB)
			}

			// Set registry ID for proxy cache at project level
			if proxyCache {
				req.RegistryID = &registryID
			}
			if proxyCache {
				if registryID == 0 {
					return fmt.Errorf("--registry-id is required when creating proxy cache project")
				}
				metadata.ProxySpeedKB = fmt.Sprintf("%d", proxySpeedKB)
			}

			// Set limits
			if storageLimitBytes != -1 {
				req.StorageLimit = &storageLimitBytes
			}
			if memberLimit > 0 {
				req.CountLimit = &memberLimit
			}

			// Create project
			if err := projectSvc.Create(req); err != nil {
				return fmt.Errorf("failed to create project: %w", err)
			}

			output.Success("Project '%s' created successfully", projectName)

			// Show created project details
			if output.GetFormat() == "table" {
				output.Info("")
				output.Info("Project Details:")
				fmt.Printf("  Name:          %s\n", projectName)
				fmt.Printf("  Public:        %v\n", public)
				if storageLimit != "" {
					fmt.Printf("  Storage Limit: %s\n", storageLimit)
				}
				if enableTrust || preventVul || autoScan {
					output.Info("  Security:")
					if enableTrust {
						fmt.Printf("    Content Trust: Enabled\n")
					}
					if preventVul {
						fmt.Printf("    Prevent Vulnerable: Yes (Severity: %s)\n", severity)
					}
					if autoScan {
						fmt.Printf("    Auto Scan: Enabled\n")
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&public, "public", false, "Make project public")
	cmd.Flags().
		StringVar(&storageLimit, "storage-limit", "", "Storage quota (e.g., 10G, 500M, -1 for unlimited)")
	cmd.Flags().BoolVar(&proxyCache, "proxy-cache", false, "Create as proxy cache project")
	cmd.Flags().Int64Var(&registryID, "registry-id", 0, "Registry endpoint ID for proxy cache")
	cmd.Flags().StringVar(&registryName, "registry-name", "", "Registry endpoint name for proxy cache (alternative to registry-id)")
	cmd.Flags().Int64Var(&proxySpeedKB, "proxy-speed", -1, "Proxy cache bandwidth limit in KB/s (-1 for unlimited)")
	cmd.Flags().Int64Var(&memberLimit, "member-limit", 0, "Project member limit")
	cmd.Flags().BoolVar(&enableTrust, "enable-content-trust", false, "Enable content trust")
	cmd.Flags().BoolVar(&preventVul, "prevent-vulnerable", false, "Prevent vulnerable images")
	cmd.Flags().
		StringVar(&severity, "severity", "low", "Vulnerability severity threshold (low, medium, high, critical)")
	cmd.Flags().BoolVar(&autoScan, "auto-scan", false, "Automatically scan images on push")
	cmd.Flags().BoolVar(&reuseSysCVE, "reuse-sys-cve", false, "Reuse system CVE allowlist")

	return cmd
}

func newProjectGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get project details",
		Long:  `Get detailed information about a specific project.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			projectSvc := harbor.NewProjectService(client)

			project, err := projectSvc.Get(projectName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Get project summary for additional details
			summary, err := projectSvc.GetSummary(projectName)
			if err != nil {
				output.Warning("Failed to get project summary: %v", err)
			}

			switch output.GetFormat() {
			case "json":
				result := map[string]interface{}{
					"project": project,
					"summary": summary,
				}
				return output.JSON(result)
			case "yaml":
				result := map[string]interface{}{
					"project": project,
					"summary": summary,
				}
				return output.YAML(result)
			default:
				output.Info("Project: %s", output.Bold(project.Name))
				output.Info("")
				fmt.Printf("ID:              %d\n", project.ProjectID)
				fmt.Printf("Public:          %v\n", project.Public)
				fmt.Printf("Owner:           %s\n", project.OwnerName)
				fmt.Printf("Repo Count:      %d\n", project.RepoCount)
				fmt.Printf(
					"Created:         %s\n",
					project.CreationTime.Format("2006-01-02 15:04:05"),
				)
				fmt.Printf(
					"Updated:         %s\n",
					project.UpdateTime.Format("2006-01-02 15:04:05"),
				)

				if project.Metadata != nil {
					output.Info("\nConfiguration:")
					if project.Metadata.EnableContentTrust != "" {
						fmt.Printf("  Content Trust:    %s\n", project.Metadata.EnableContentTrust)
					}
					if project.Metadata.PreventVul != "" {
						fmt.Printf("  Prevent Vulnerable: %s\n", project.Metadata.PreventVul)
						if project.Metadata.Severity != "" {
							fmt.Printf("  Severity:         %s\n", project.Metadata.Severity)
						}
					}
					if project.Metadata.AutoScan != "" {
						fmt.Printf("  Auto Scan:        %s\n", project.Metadata.AutoScan)
					}
				}

				if summary != nil {
					output.Info("\nQuota:")
					if summary.Quota != nil {
						fmt.Printf(
							"  Storage Limit:    %s\n",
							output.FormatSize(summary.Quota.Hard.Storage),
						)
						fmt.Printf(
							"  Storage Used:     %s\n",
							output.FormatSize(summary.Quota.Used.Storage),
						)
						if summary.Quota.Hard.Storage > 0 {
							percentage := float64(
								summary.Quota.Used.Storage,
							) / float64(
								summary.Quota.Hard.Storage,
							) * 100
							fmt.Printf("  Storage Usage:    %.1f%%\n", percentage)
						}
					}

					output.Info("\nMembers:")
					fmt.Printf("  Admins:          %d\n", summary.ProjectAdminCount)
					fmt.Printf("  Developers:      %d\n", summary.DeveloperCount)
					fmt.Printf("  Guests:          %d\n", summary.GuestCount)
					fmt.Printf("  Limited Guests:  %d\n", summary.LimitedGuestCount)
				}
			}

			return nil
		},
	}
}

func newProjectUpdateCmd() *cobra.Command {
	var (
		public       *bool
		storageLimit string
		enableTrust  *bool
		preventVul   *bool
		severity     string
		autoScan     *bool
		reuseSysCVE  *bool
	)

	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Update project settings",
		Long:  `Update settings for an existing project.`,
		Example: `  # Make project public
  hrbcli project update myproject --public=true

  # Update storage quota
  hrbcli project update myproject --storage-limit 20G

  # Enable security features
  hrbcli project update myproject --enable-content-trust=true --auto-scan=true`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			projectSvc := harbor.NewProjectService(client)

			// Get current project to preserve existing settings
			project, err := projectSvc.Get(projectName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Build update request
			req := &api.ProjectReq{}

			// Only update fields that were explicitly set
			if cmd.Flags().Changed("public") {
				req.Public = public
			}

			// Handle storage limit
			if cmd.Flags().Changed("storage-limit") {
				storageLimitBytes, err := harbor.ParseStorageLimit(storageLimit)
				if err != nil {
					return err
				}
				req.StorageLimit = &storageLimitBytes
			}

			// Build metadata with only changed fields
			metadata := &api.ProjectMetadata{}
			hasMetadataChanges := false

			if cmd.Flags().Changed("enable-content-trust") {
				if *enableTrust {
					metadata.EnableContentTrust = "true"
				} else {
					metadata.EnableContentTrust = "false"
				}
				hasMetadataChanges = true
			}

			if cmd.Flags().Changed("prevent-vulnerable") {
				if *preventVul {
					metadata.PreventVul = "true"
				} else {
					metadata.PreventVul = "false"
				}
				hasMetadataChanges = true
			}

			if cmd.Flags().Changed("severity") {
				metadata.Severity = severity
				hasMetadataChanges = true
			}

			if cmd.Flags().Changed("auto-scan") {
				if *autoScan {
					metadata.AutoScan = "true"
				} else {
					metadata.AutoScan = "false"
				}
				hasMetadataChanges = true
			}

			if cmd.Flags().Changed("reuse-sys-cve") {
				if *reuseSysCVE {
					metadata.ReuseSysCVEAllowlist = "true"
				} else {
					metadata.ReuseSysCVEAllowlist = "false"
				}
				hasMetadataChanges = true
			}

			if hasMetadataChanges {
			}

			// Update project
			if err := projectSvc.Update(project.Name, req); err != nil {
				return fmt.Errorf("failed to update project: %w", err)
			}

			output.Success("Project '%s' updated successfully", projectName)
			return nil
		},
	}

	// Use pointers for boolean flags to distinguish between not set and false
	cmd.Flags().BoolVar(&dummyBool, "public", false, "Make project public/private")
	cmd.Flags().
		StringVar(&storageLimit, "storage-limit", "", "Storage quota (e.g., 10G, 500M, -1 for unlimited)")
	cmd.Flags().BoolVar(&dummyBool, "enable-content-trust", false, "Enable/disable content trust")
	cmd.Flags().
		BoolVar(&dummyBool, "prevent-vulnerable", false, "Enable/disable preventing vulnerable images")
	cmd.Flags().StringVar(&severity, "severity", "", "Vulnerability severity threshold")
	cmd.Flags().BoolVar(&dummyBool, "auto-scan", false, "Enable/disable auto scan")
	cmd.Flags().
		BoolVar(&dummyBool, "reuse-sys-cve", false, "Enable/disable reusing system CVE allowlist")

	// Custom parsing for boolean flags
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("public") {
			val, _ := cmd.Flags().GetBool("public")
			public = &val
		}
		if cmd.Flags().Changed("enable-content-trust") {
			val, _ := cmd.Flags().GetBool("enable-content-trust")
			enableTrust = &val
		}
		if cmd.Flags().Changed("prevent-vulnerable") {
			val, _ := cmd.Flags().GetBool("prevent-vulnerable")
			preventVul = &val
		}
		if cmd.Flags().Changed("auto-scan") {
			val, _ := cmd.Flags().GetBool("auto-scan")
			autoScan = &val
		}
		if cmd.Flags().Changed("reuse-sys-cve") {
			val, _ := cmd.Flags().GetBool("reuse-sys-cve")
			reuseSysCVE = &val
		}
		return nil
	}

	return cmd
}

// dummyBool is used for flag parsing
var dummyBool bool

func newProjectDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a project",
		Long:  `Delete a project. The project must be empty (no repositories).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			projectSvc := harbor.NewProjectService(client)

			// Get project details first
			project, err := projectSvc.Get(projectName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Check if project has repositories
			if project.RepoCount > 0 && !force {
				return fmt.Errorf(
					"project '%s' contains %d repositories. Use --force to delete anyway",
					projectName,
					project.RepoCount,
				)
			}

			// Confirm deletion if not forced
			if !force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Delete project '%s'", projectName),
					IsConfirm: true,
				}
				result, err := prompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					output.Info("Deletion cancelled")
					return nil
				}
			}

			// Delete project
			if err := projectSvc.Delete(projectName); err != nil {
				return fmt.Errorf("failed to delete project: %w", err)
			}

			output.Success("Project '%s' deleted successfully", projectName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion without confirmation")

	return cmd
}

func newProjectExistsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exists <name>",
		Short: "Check if a project exists",
		Long:  `Check if a project exists. Returns exit code 0 if exists, 1 if not.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			projectSvc := harbor.NewProjectService(client)

			exists, err := projectSvc.Exists(projectName)
			if err != nil {
				return fmt.Errorf("failed to check project existence: %w", err)
			}

			if exists {
				output.Info("Project '%s' exists", projectName)
				return nil
			} else {
				output.Info("Project '%s' does not exist", projectName)
				os.Exit(1)
			}

			return nil
		},
	}
}

// validateProjectName validates project name according to Harbor rules
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) > 255 {
		return fmt.Errorf("project name cannot exceed 255 characters")
	}

	// Harbor project name rules:
	// - Must be lowercase
	// - Can contain letters, numbers, and special characters ._-
	// - Must start with a letter or number
	for i, c := range name {
		if i == 0 {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
				return fmt.Errorf("project name must start with a lowercase letter or number")
			}
		} else {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '.' || c == '_' || c == '-') {
				return fmt.Errorf("project name can only contain lowercase letters, numbers, and ._- characters")
			}
		}
	}

	return nil
}
