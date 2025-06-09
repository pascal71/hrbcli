package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewReplicationCmd creates the replication command
func NewReplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replication",
		Short: "Manage replication policies",
	}

	cmd.AddCommand(newReplicationListCmd())
	cmd.AddCommand(newReplicationCreateCmd())
	cmd.AddCommand(newReplicationGetCmd())
	cmd.AddCommand(newReplicationDeleteCmd())
	cmd.AddCommand(newReplicationExecuteCmd())
	cmd.AddCommand(newReplicationExecutionsCmd())
	cmd.AddCommand(newReplicationExecutionCmd())
	cmd.AddCommand(newReplicationLogsCmd())
	cmd.AddCommand(newReplicationStatsCmd())

	return cmd
}

func newReplicationListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List replication policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			policies, err := svc.ListPolicies(nil)
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
						strconv.FormatInt(p.ID, 10),
						p.Name,
						strconv.FormatBool(p.Enabled),
					})
				}
				table.Render()
			}
			return nil
		},
	}
	return cmd
}

func newReplicationGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Show replication policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			policy, err := svc.GetPolicy(id)
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
				table.Append([]string{"ID", strconv.FormatInt(policy.ID, 10)})
				table.Append([]string{"NAME", policy.Name})
				table.Append([]string{"ENABLED", strconv.FormatBool(policy.Enabled)})
				table.Render()
			}
			return nil
		},
	}
	return cmd
}

func newReplicationCreateCmd() *cobra.Command {
	var src, dst, name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create replication policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" || src == "" || dst == "" {
				return fmt.Errorf("--name, --source and --destination are required")
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			req := &api.ReplicationPolicy{Name: name, SrcRegistry: &api.Registry{Name: src}, DestRegistry: &api.Registry{Name: dst}, Enabled: true}
			policy, err := svc.CreatePolicy(req)
			if err != nil {
				return err
			}
			output.Success("Created replication policy %s (ID %d)", policy.Name, policy.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&src, "source", "", "Source registry")
	cmd.Flags().StringVar(&dst, "destination", "", "Destination registry")
	return cmd
}

func newReplicationDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete replication policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			if err := svc.DeletePolicy(id); err != nil {
				return err
			}
			output.Success("Deleted replication policy %d", id)
			return nil
		},
	}
	return cmd
}

func newReplicationExecuteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute <policy-id>",
		Short: "Trigger replication execution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			exec, err := svc.StartExecution(id)
			if err != nil {
				return err
			}
			output.Success("Started execution %d", exec.ID)
			return nil
		},
	}
	return cmd
}

func newReplicationExecutionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "executions [policy-id]",
		Short: "List replication executions",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id int64
			var err error
			if len(args) == 1 {
				id, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			execs, err := svc.ListExecutions(id)
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(execs)
			case "yaml":
				return output.YAML(execs)
			default:
				table := output.Table()
				table.Append([]string{"ID", "POLICY", "STATUS", "START", "END"})
				for _, e := range execs {
					table.Append([]string{
						strconv.FormatInt(e.ID, 10),
						strconv.FormatInt(e.PolicyID, 10),
						e.Status,
						e.StartTime.Format("2006-01-02 15:04"),
						e.EndTime.Format("2006-01-02 15:04"),
					})
				}
				table.Render()
			}
			return nil
		},
	}
	return cmd
}

func newReplicationExecutionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution <id>",
		Short: "Show execution details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			exec, err := svc.GetExecution(id)
			if err != nil {
				return err
			}
			switch output.GetFormat() {
			case "json":
				return output.JSON(exec)
			case "yaml":
				return output.YAML(exec)
			default:
				table := output.Table()
				table.Append([]string{"FIELD", "VALUE"})
				table.Append([]string{"STATUS", exec.Status})
				table.Append([]string{"SUCCEED", strconv.Itoa(exec.Succeed)})
				table.Append([]string{"FAILED", strconv.Itoa(exec.Failed)})
				table.Append([]string{"IN_PROGRESS", strconv.Itoa(exec.InProgress)})
				table.Append([]string{"STOPPED", strconv.Itoa(exec.Stopped)})
				table.Render()
			}
			return nil
		},
	}
	return cmd
}

func newReplicationLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <execution-id>",
		Short: "Show execution logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			tasks, err := svc.ListTasks(id)
			if err != nil {
				return err
			}
			for _, t := range tasks {
				log, err := svc.GetTaskLog(id, t.ID)
				if err != nil {
					return err
				}
				output.Info("Task %d", t.ID)
				fmt.Println(log)
			}
			return nil
		},
	}
	return cmd
}

func newReplicationStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Show replication statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewReplicationService(client)
			execs, err := svc.ListExecutions(0)
			if err != nil {
				return err
			}
			var succeed, failed, running int
			for _, e := range execs {
				succeed += e.Succeed
				failed += e.Failed
				running += e.InProgress
			}
			fmt.Printf("Succeeded: %d\nFailed: %d\nRunning: %d\n", succeed, failed, running)
			return nil
		},
	}
	return cmd
}
