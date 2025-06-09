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

// NewRegistryCmd creates the registry command
func NewRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "registry",
		Aliases: []string{"reg"},
		Short:   "Manage registry endpoints",
		Long:    `Manage registry endpoints for replication and proxy cache.`,
	}

	cmd.AddCommand(newRegistryListCmd())
	cmd.AddCommand(newRegistryCreateCmd())
	cmd.AddCommand(newRegistryGetCmd())
	cmd.AddCommand(newRegistryUpdateCmd())
	cmd.AddCommand(newRegistryDeleteCmd())
	cmd.AddCommand(newRegistryPingCmd())
	cmd.AddCommand(newRegistryAdaptersCmd())

	return cmd
}

func newRegistryListCmd() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registry endpoints",
		Long:  `List all configured registry endpoints.`,
		Example: `  # List all registries
  hrbcli registry list

  # Search registries by name
  hrbcli registry list --query docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			registrySvc := harbor.NewRegistryService(client)

			opts := &api.ListOptions{
				Query: query,
			}

			registries, err := registrySvc.List(opts)
			if err != nil {
				return fmt.Errorf("failed to list registries: %w", err)
			}

			if len(registries) == 0 {
				output.Info("No registry endpoints found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(registries)
			case "yaml":
				return output.YAML(registries)
			default:
				table := output.Table()
				headers := []string{"ID", "NAME", "TYPE", "URL", "STATUS", "INSECURE", "CREATED"}
				table.Append(headers)

				for _, r := range registries {
					row := []string{
						strconv.FormatInt(r.ID, 10),
						r.Name,
						r.Type,
						r.URL,
						r.Status,
						strconv.FormatBool(r.Insecure),
						r.CreationTime.Format("2006-01-02"),
					}
					table.Append(row)
				}

				table.Render()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Search query")

	return cmd
}

func newRegistryCreateCmd() *cobra.Command {
	var (
		url         string
		description string
		regType     string
		insecure    bool
		username    string
		password    string
		interactive bool
	)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a registry endpoint",
		Long:  `Create a new registry endpoint for replication or proxy cache.`,
		Example: `  # Create Docker Hub endpoint
  hrbcli registry create dockerhub --type docker-hub --url https://hub.docker.com

  # Create private registry with credentials
  hrbcli registry create myregistry --url https://registry.example.com --username user --password pass

  # Create Quay.io endpoint
  hrbcli registry create quay --type quay --url https://quay.io

  # Interactive mode
  hrbcli registry create --interactive`,
		Args: func(cmd *cobra.Command, args []string) error {
			if interactive {
				return nil
			}
			if len(args) != 1 {
				return fmt.Errorf("requires registry name argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			
			if interactive {
				// Interactive mode code here...
				output.Info("Interactive mode not fully implemented yet")
				return fmt.Errorf("please use command line flags")
			} else {
				name = args[0]
				
				// Validate required flags
				if url == "" {
					return fmt.Errorf("--url is required")
				}
				if regType == "" {
					regType = api.RegistryTypeDockerRegistry
				}
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			registrySvc := harbor.NewRegistryService(client)

			// Build request
			req := &api.RegistryReq{
				Name:        name,
				URL:         url,
				Description: description,
				Type:        regType,
				Insecure:    insecure,
			}

			// Add credentials if provided
			if username != "" || password != "" {
				req.Credential = &api.Credential{
					Type:         "basic",
					AccessKey:    username,
					AccessSecret: password,
				}
			}

			// Test connectivity first
			output.Info("Testing connectivity to registry...")
			if err := registrySvc.Ping(req); err != nil {
				output.Warning("Registry ping failed: %v", err)
			} else {
				output.Success("Registry is reachable")
			}

			// Create registry
			registry, err := registrySvc.Create(req)
			if err != nil {
				return fmt.Errorf("failed to create registry: %w", err)
			}

			output.Success("Registry endpoint '%s' created successfully (ID: %d)", name, registry.ID)
			
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Registry URL")
	cmd.Flags().StringVar(&description, "description", "", "Registry description")
	cmd.Flags().StringVar(&regType, "type", api.RegistryTypeDockerRegistry, "Registry type")
	cmd.Flags().BoolVar(&insecure, "insecure", false, "Skip TLS verification")
	cmd.Flags().StringVar(&username, "username", "", "Registry username")
	cmd.Flags().StringVar(&password, "password", "", "Registry password")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")

	return cmd
}

func newRegistryGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get registry details",
		Long:  `Get detailed information about a registry endpoint.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid registry ID: %s", args[0])
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			registrySvc := harbor.NewRegistryService(client)

			registry, err := registrySvc.Get(id)
			if err != nil {
				return fmt.Errorf("failed to get registry: %w", err)
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(registry)
			case "yaml":
				return output.YAML(registry)
			default:
				output.Info("Registry: %s", output.Bold(registry.Name))
				output.Info("")
				fmt.Printf("ID:          %d\n", registry.ID)
				fmt.Printf("Name:        %s\n", registry.Name)
				fmt.Printf("Type:        %s\n", registry.Type)
				fmt.Printf("URL:         %s\n", registry.URL)
				fmt.Printf("Status:      %s\n", registry.Status)
				fmt.Printf("Insecure:    %v\n", registry.Insecure)
				fmt.Printf("Created:     %s\n", registry.CreationTime.Format("2006-01-02 15:04:05"))
				fmt.Printf("Updated:     %s\n", registry.UpdateTime.Format("2006-01-02 15:04:05"))
				
				if registry.Description != "" {
					fmt.Printf("Description: %s\n", registry.Description)
				}
				
				if registry.Credential != nil {
					output.Info("\nCredentials:")
					fmt.Printf("  Type:     %s\n", registry.Credential.Type)
					fmt.Printf("  Username: %s\n", registry.Credential.AccessKey)
				}
			}

			return nil
		},
	}
}

func newRegistryUpdateCmd() *cobra.Command {
	var (
		url         string
		description string
		// insecure    *bool
		username    string
		password    string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update registry endpoint",
		Long:  `Update an existing registry endpoint configuration.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Implementation here...
			output.Info("Update command not fully implemented yet")
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Registry URL")
	cmd.Flags().StringVar(&description, "description", "", "Registry description")
	cmd.Flags().BoolVar(&dummyBool, "insecure", false, "Skip TLS verification")
	cmd.Flags().StringVar(&username, "username", "", "Registry username")
	cmd.Flags().StringVar(&password, "password", "", "Registry password")

	return cmd
}

func newRegistryDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete registry endpoint",
		Long:  `Delete a registry endpoint. The registry must not be used by any replication rules.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid registry ID: %s", args[0])
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			registrySvc := harbor.NewRegistryService(client)

			// Get registry details first
			registry, err := registrySvc.Get(id)
			if err != nil {
				return fmt.Errorf("failed to get registry: %w", err)
			}

			// Confirm deletion if not forced
			if !force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Delete registry '%s'", registry.Name),
					IsConfirm: true,
				}
				result, err := prompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					output.Info("Deletion cancelled")
					return nil
				}
			}

			// Delete registry
			if err := registrySvc.Delete(id); err != nil {
				return fmt.Errorf("failed to delete registry: %w", err)
			}

			output.Success("Registry '%s' deleted successfully", registry.Name)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion without confirmation")

	return cmd
}

func newRegistryPingCmd() *cobra.Command {
	var (
		url      string
		regType  string
		insecure bool
		username string
		password string
	)

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Test registry connectivity",
		Long:  `Test connectivity to a registry endpoint.`,
		Example: `  # Test Docker Hub
  hrbcli registry ping --url https://hub.docker.com --type docker-hub

  # Test private registry with auth
  hrbcli registry ping --url https://registry.example.com --username user --password pass`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				return fmt.Errorf("--url is required")
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}

			registrySvc := harbor.NewRegistryService(client)

			// Build request
			req := &api.RegistryReq{
				Name:     "ping-test",
				URL:      url,
				Type:     regType,
				Insecure: insecure,
			}

			// Add credentials if provided
			if username != "" || password != "" {
				req.Credential = &api.Credential{
					Type:         "basic",
					AccessKey:    username,
					AccessSecret: password,
				}
			}

			// Test connectivity
			output.Info("Testing connectivity to %s...", url)
			if err := registrySvc.Ping(req); err != nil {
				output.Error("Ping failed: %v", err)
				return err
			}

			output.Success("Registry is reachable!")
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Registry URL (required)")
	cmd.Flags().StringVar(&regType, "type", api.RegistryTypeDockerRegistry, "Registry type")
	cmd.Flags().BoolVar(&insecure, "insecure", false, "Skip TLS verification")
	cmd.Flags().StringVar(&username, "username", "", "Registry username")
	cmd.Flags().StringVar(&password, "password", "", "Registry password")

	cmd.MarkFlagRequired("url")

	return cmd
}

func newRegistryAdaptersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "adapters",
		Short: "List available registry adapters",
		Long:  `List all available registry adapters and their capabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			registrySvc := harbor.NewRegistryService(client)

			adapters, err := registrySvc.ListAdapters()
			if err != nil {
				return fmt.Errorf("failed to list adapters: %w", err)
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(adapters)
			case "yaml":
				return output.YAML(adapters)
			default:
				table := output.Table()
				headers := []string{"TYPE", "DESCRIPTION", "SUPPORTED RESOURCES", "SUPPORTED FILTERS"}
				table.Append(headers)

				for name, info := range adapters {
					row := []string{
						name,
						info.Description,
						strings.Join(info.SupportedTypes, ", "),
						strings.Join(info.Filters, ", "),
					}
					table.Append(row)
				}

				table.Render()
			}

			return nil
		},
	}
}

// getDefaultURL returns the default URL for a registry type
func getDefaultURL(regType string) string {
	switch regType {
	case api.RegistryTypeDockerHub:
		return "https://hub.docker.com"
	case api.RegistryTypeQuay:
		return "https://quay.io"
	case api.RegistryTypeGoogleGCR:
		return "https://gcr.io"
	case api.RegistryTypeGitlab:
		return "https://registry.gitlab.com"
	default:
		return "https://"
	}
}
