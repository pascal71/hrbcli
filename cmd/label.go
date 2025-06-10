package cmd

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewLabelCmd creates the label command
func NewLabelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels",
		Long:  `Manage Harbor labels.`,
	}

	cmd.AddCommand(newLabelListCmd())
	cmd.AddCommand(newLabelCreateCmd())
	cmd.AddCommand(newLabelGetCmd())
	cmd.AddCommand(newLabelUpdateCmd())
	cmd.AddCommand(newLabelDeleteCmd())

	return cmd
}

func newLabelListCmd() *cobra.Command {
	var (
		page      int
		pageSize  int
		name      string
		scope     string
		projectID int64
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}

			labelSvc := harbor.NewLabelService(client)

			opts := &api.LabelListOptions{
				Page:      page,
				PageSize:  pageSize,
				Name:      name,
				Scope:     scope,
				ProjectID: projectID,
			}

			labels, err := labelSvc.List(opts)
			if err != nil {
				return fmt.Errorf("failed to list labels: %w", err)
			}

			if len(labels) == 0 {
				output.Info("No labels found")
				return nil
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(labels)
			case "yaml":
				return output.YAML(labels)
			default:
				table := output.Table()
				table.Append([]string{"ID", "NAME", "SCOPE", "PROJECT"})
				for _, l := range labels {
					pid := ""
					if l.ProjectID != 0 {
						pid = strconv.FormatInt(l.ProjectID, 10)
					}
					table.Append([]string{
						strconv.FormatInt(l.ID, 10),
						l.Name,
						l.Scope,
						pid,
					})
				}
				table.Render()
				return nil
			}
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().StringVar(&name, "name", "", "Filter by name")
	cmd.Flags().StringVar(&scope, "scope", "", "Label scope (g or p)")
	cmd.Flags().Int64Var(&projectID, "project-id", 0, "Project ID for project labels")
	return cmd
}

func newLabelCreateCmd() *cobra.Command {
	var (
		description string
		color       string
		scope       string
		projectID   int64
	)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a label",
		Args:  requireArgs(1, "requires <name>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			client, err := api.NewClient()
			if err != nil {
				return err
			}
			labelSvc := harbor.NewLabelService(client)

			label := &api.Label{
				Name:        name,
				Description: description,
				Color:       color,
				Scope:       scope,
			}
			if projectID > 0 {
				label.ProjectID = projectID
			}

			created, err := labelSvc.Create(label)
			if err != nil {
				if apiErr, ok := err.(*api.APIError); ok && apiErr.IsConflict() {
					return fmt.Errorf("label '%s' already exists", name)
				}
				return fmt.Errorf("failed to create label: %w", err)
			}

			output.Success("Label '%s' created (ID: %d)", created.Name, created.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Label description")
	cmd.Flags().StringVar(&color, "color", "", "Label color")
	cmd.Flags().StringVar(&scope, "scope", "g", "Label scope (g or p)")
	cmd.Flags().Int64Var(&projectID, "project-id", 0, "Project ID when scope is p")
	return cmd
}

func newLabelGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get label details",
		Args:  requireArgs(1, "requires <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid label ID: %s", args[0])
			}

			client, err := api.NewClient()
			if err != nil {
				return err
			}
			labelSvc := harbor.NewLabelService(client)

			label, err := labelSvc.Get(id)
			if err != nil {
				return fmt.Errorf("failed to get label: %w", err)
			}

			switch output.GetFormat() {
			case "json":
				return output.JSON(label)
			case "yaml":
				return output.YAML(label)
			default:
				output.Info("Label: %s", output.Bold(label.Name))
				fmt.Printf("ID:          %d\n", label.ID)
				fmt.Printf("Scope:       %s\n", label.Scope)
				if label.ProjectID != 0 {
					fmt.Printf("Project ID:  %d\n", label.ProjectID)
				}
				if label.Description != "" {
					fmt.Printf("Description: %s\n", label.Description)
				}
				if label.Color != "" {
					fmt.Printf("Color:       %s\n", label.Color)
				}
				return nil
			}
		},
	}
	return cmd
}

func newLabelUpdateCmd() *cobra.Command {
	var (
		name        string
		description string
		color       string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a label",
		Args:  requireArgs(1, "requires <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid label ID: %s", args[0])
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			labelSvc := harbor.NewLabelService(client)

			label := &api.Label{}
			if name != "" {
				label.Name = name
			}
			if description != "" {
				label.Description = description
			}
			if color != "" {
				label.Color = color
			}

			if err := labelSvc.Update(id, label); err != nil {
				return fmt.Errorf("failed to update label: %w", err)
			}
			output.Success("Label %d updated", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New label name")
	cmd.Flags().StringVar(&description, "description", "", "New description")
	cmd.Flags().StringVar(&color, "color", "", "New color")
	return cmd
}

func newLabelDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a label",
		Args:  requireArgs(1, "requires <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid label ID: %s", args[0])
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			labelSvc := harbor.NewLabelService(client)

			if !force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Delete label %d", id),
					IsConfirm: true,
				}
				result, err := prompt.Run()
				if err != nil || result != "y" {
					output.Info("Deletion cancelled")
					return nil
				}
			}

			if err := labelSvc.Delete(id); err != nil {
				return fmt.Errorf("failed to delete label: %w", err)
			}
			output.Success("Label %d deleted", id)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion without confirmation")
	return cmd
}
