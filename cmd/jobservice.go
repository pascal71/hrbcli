package cmd

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/pascal71/hrbcli/pkg/harbor"
	"github.com/pascal71/hrbcli/pkg/output"
)

// NewJobServiceCmd creates the jobservice command
func NewJobServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobservice",
		Short: "Manage Harbor job service",
	}

	cmd.AddCommand(newJobServiceDashboardCmd())

	return cmd
}

func newJobServiceDashboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Show job service dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return err
			}
			svc := harbor.NewJobService(client)

			pools, err := svc.GetWorkerPools()
			if err != nil {
				return err
			}
			workers, err := svc.GetWorkers("all")
			if err != nil {
				return err
			}
			queues, err := svc.ListJobQueues()
			if err != nil {
				return err
			}

			switch output.GetFormat() {
			case "json":
				data := struct {
					Pools   []*api.WorkerPool `json:"pools"`
					Workers []*api.Worker     `json:"workers"`
					Queues  []*api.JobQueue   `json:"queues"`
				}{pools, workers, queues}
				return output.JSON(data)
			case "yaml":
				data := struct {
					Pools   []*api.WorkerPool `yaml:"pools"`
					Workers []*api.Worker     `yaml:"workers"`
					Queues  []*api.JobQueue   `yaml:"queues"`
				}{pools, workers, queues}
				return output.YAML(data)
			default:
				table := output.Table()
				table.Append([]string{"POOL ID", "PID", "HOST", "CONC", "START", "HEARTBEAT"})
				for _, p := range pools {
					table.Append([]string{
						p.WorkerPoolID,
						strconv.FormatInt(p.PID, 10),
						p.Host,
						strconv.Itoa(p.Concurrency),
						p.StartAt.Format("2006-01-02 15:04:05"),
						p.HeartbeatAt.Format("2006-01-02 15:04:05"),
					})
				}
				table.Render()

				table = output.Table()
				table.Append([]string{"WORKER ID", "POOL", "JOB", "JOB ID", "START", "CHECKIN"})
				for _, w := range workers {
					start := ""
					if !w.StartAt.IsZero() {
						start = w.StartAt.Format("2006-01-02 15:04:05")
					}
					check := w.CheckIn
					if w.CheckInAt.After(w.StartAt) {
						check = w.CheckIn + " @ " + w.CheckInAt.Format("2006-01-02 15:04:05")
					}
					table.Append([]string{
						w.ID,
						w.PoolID,
						w.JobName,
						w.JobID,
						start,
						check,
					})
				}
				table.Render()

				table = output.Table()
				table.Append([]string{"JOB TYPE", "COUNT", "LATENCY", "PAUSED"})
				for _, q := range queues {
					table.Append([]string{
						q.JobType,
						strconv.Itoa(q.Count),
						strconv.Itoa(q.Latency),
						strconv.FormatBool(q.Paused),
					})
				}
				table.Render()
				return nil
			}
		},
	}
	return cmd
}
