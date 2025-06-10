package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewDistributionCmd creates the distribution command
func NewDistributionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distribution",
		Short: "Manage distribution providers and policies",
	}

	cmd.AddCommand(newDistributionProvidersCmd())
	cmd.AddCommand(newDistributionPoliciesCmd())
	cmd.AddCommand(newDistributionPolicyGetCmd())

	return cmd
}

func newDistributionProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers <project>",
		Short: "List distribution providers",
		Args:  requireArgs(1, "requires <project>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := args[0]
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewPreheatService(client)
			providers, err := svc.ListProviders(project)
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(providers)
			case "yaml":
				return output.YAML(providers)
			default:
				table := output.Table()
				table.Append([]string{"ID", "PROVIDER", "ENABLED", "DEFAULT"})
				for _, p := range providers {
					table.Append([]string{
						fmt.Sprintf("%d", p.ID),
						p.Provider,
						fmt.Sprintf("%v", p.Enabled),
						fmt.Sprintf("%v", p.Default),
					})
				}
				table.Render()
				return nil
			}
		},
	}
	return cmd
}

func newDistributionPoliciesCmd() *cobra.Command {
	var page, pageSize int
	var query, sort string

	cmd := &cobra.Command{
		Use:   "policies <project>",
		Short: "List distribution policies",
		Args:  requireArgs(1, "requires <project>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := args[0]
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewPreheatService(client)
			opts := &api.ListOptions{Page: page, PageSize: pageSize, Query: query, Sort: sort}
			policies, err := svc.ListPolicies(project, opts)
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(policies)
			case "yaml":
				return output.YAML(policies)
			default:
				table := output.Table()
				table.Append([]string{"ID", "NAME", "ENABLED"})
				for _, p := range policies {
					table.Append([]string{
						fmt.Sprintf("%d", p.ID),
						p.Name,
						fmt.Sprintf("%v", p.Enabled),
					})
				}
				table.Render()
				return nil
			}
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	cmd.Flags().StringVar(&query, "query", "", "Query filter")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort order")
	return cmd
}

func newDistributionPolicyGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy <project> <name>",
		Short: "Show distribution policy",
		Args:  requireArgs(2, "requires <project> <name>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := args[0]
			name := args[1]
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewPreheatService(client)
			policy, err := svc.GetPolicy(project, name)
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(policy)
			case "yaml":
				return output.YAML(policy)
			default:
				table := output.Table()
				table.Append([]string{"FIELD", "VALUE"})
				table.Append([]string{"ID", fmt.Sprintf("%d", policy.ID)})
				table.Append([]string{"NAME", policy.Name})
				table.Append([]string{"ENABLED", fmt.Sprintf("%v", policy.Enabled)})
				table.Render()
				return nil
			}
		},
	}
	return cmd
}
