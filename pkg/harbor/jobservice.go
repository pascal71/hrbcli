package harbor

import (
	"fmt"

	"github.com/pascal71/hrbcli/pkg/api"
)

// JobService handles job service operations
// such as retrieving worker pools and job queues.
type JobService struct {
	client *api.Client
}

// NewJobService creates a new JobService
func NewJobService(client *api.Client) *JobService {
	return &JobService{client: client}
}

// GetWorkerPools retrieves worker pools
func (s *JobService) GetWorkerPools() ([]*api.WorkerPool, error) {
	resp, err := s.client.Get("/jobservice/pools", nil)
	if err != nil {
		return nil, err
	}
	var pools []*api.WorkerPool
	if err := s.client.DecodeResponse(resp, &pools); err != nil {
		return nil, fmt.Errorf("failed to decode worker pools: %w", err)
	}
	return pools, nil
}

// GetWorkers retrieves workers in a pool. Use poolID="all" to get all workers.
func (s *JobService) GetWorkers(poolID string) ([]*api.Worker, error) {
	resp, err := s.client.Get(fmt.Sprintf("/jobservice/pools/%s/workers", poolID), nil)
	if err != nil {
		return nil, err
	}
	var workers []*api.Worker
	if err := s.client.DecodeResponse(resp, &workers); err != nil {
		return nil, fmt.Errorf("failed to decode workers: %w", err)
	}
	return workers, nil
}

// ListJobQueues lists job queues
func (s *JobService) ListJobQueues() ([]*api.JobQueue, error) {
	resp, err := s.client.Get("/jobservice/queues", nil)
	if err != nil {
		return nil, err
	}
	var queues []*api.JobQueue
	if err := s.client.DecodeResponse(resp, &queues); err != nil {
		return nil, fmt.Errorf("failed to decode job queues: %w", err)
	}
	return queues, nil
}
