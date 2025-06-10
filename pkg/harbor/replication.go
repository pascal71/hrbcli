package harbor

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pascal71/hrbcli/pkg/api"
)

// ReplicationService handles replication operations
type ReplicationService struct {
	client *api.Client
}

// NewReplicationService creates a new ReplicationService
func NewReplicationService(client *api.Client) *ReplicationService {
	return &ReplicationService{client: client}
}

// ResolvePolicyID resolves a replication policy ID from its name.
// Returns an error if the policy cannot be found.
func (s *ReplicationService) ResolvePolicyID(name string) (int64, error) {
	opts := &api.ListOptions{Query: fmt.Sprintf("name=~%s", name)}
	policies, err := s.ListPolicies(opts)
	if err != nil {
		return 0, err
	}
	for _, p := range policies {
		if p.Name == name {
			return p.ID, nil
		}
	}
	return 0, fmt.Errorf("replication policy '%s' not found", name)
}

// ListPolicies lists replication policies
func (s *ReplicationService) ListPolicies(opts *api.ListOptions) ([]*api.ReplicationPolicy, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.Query != "" {
			params["query"] = opts.Query
		}
		if opts.Sort != "" {
			params["sort"] = opts.Sort
		}
		if opts.Page > 0 {
			params["page"] = fmt.Sprintf("%d", opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = fmt.Sprintf("%d", opts.PageSize)
		}
	}

	resp, err := s.client.Get("/replication/policies", params)
	if err != nil {
		return nil, err
	}

	var policies []*api.ReplicationPolicy
	if err := s.client.DecodeResponse(resp, &policies); err != nil {
		return nil, fmt.Errorf("failed to decode policies: %w", err)
	}

	return policies, nil
}

// GetPolicy retrieves a policy by ID
func (s *ReplicationService) GetPolicy(id int64) (*api.ReplicationPolicy, error) {
	resp, err := s.client.Get(fmt.Sprintf("/replication/policies/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var policy api.ReplicationPolicy
	if err := s.client.DecodeResponse(resp, &policy); err != nil {
		return nil, fmt.Errorf("failed to decode policy: %w", err)
	}
	return &policy, nil
}

// CreatePolicy creates a new replication policy
func (s *ReplicationService) CreatePolicy(req *api.ReplicationPolicy) (*api.ReplicationPolicy, error) {
	resp, err := s.client.Post("/replication/policies", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("no location header in response")
	}

	var id int64
	if _, err := fmt.Sscanf(location, "/api/v2.0/replication/policies/%d", &id); err != nil {
		return nil, fmt.Errorf("failed to parse policy ID from location: %s", location)
	}

	return s.GetPolicy(id)
}

// DeletePolicy deletes a replication policy
func (s *ReplicationService) DeletePolicy(id int64) error {
	resp, err := s.client.Delete(fmt.Sprintf("/replication/policies/%d", id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// StartExecution manually triggers a replication execution
func (s *ReplicationService) StartExecution(policyID int64) (*api.ReplicationExecution, error) {
	req := &api.StartReplicationExecution{PolicyID: policyID}
	resp, err := s.client.Post("/replication/executions", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("no location header in response")
	}

	var id int64
	if _, err := fmt.Sscanf(location, "/api/v2.0/replication/executions/%d", &id); err != nil {
		return nil, fmt.Errorf("failed to parse execution ID from location: %s", location)
	}

	return s.GetExecution(id)
}

// ListExecutions lists executions
func (s *ReplicationService) ListExecutions(policyID int64) ([]*api.ReplicationExecution, error) {
	params := make(map[string]string)
	if policyID > 0 {
		params["policy_id"] = fmt.Sprintf("%d", policyID)
	}

	resp, err := s.client.Get("/replication/executions", params)
	if err != nil {
		return nil, err
	}

	var execs []*api.ReplicationExecution
	if err := s.client.DecodeResponse(resp, &execs); err != nil {
		return nil, fmt.Errorf("failed to decode executions: %w", err)
	}
	return execs, nil
}

// GetExecution retrieves a replication execution
func (s *ReplicationService) GetExecution(id int64) (*api.ReplicationExecution, error) {
	resp, err := s.client.Get(fmt.Sprintf("/replication/executions/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var exec api.ReplicationExecution
	if err := s.client.DecodeResponse(resp, &exec); err != nil {
		return nil, fmt.Errorf("failed to decode execution: %w", err)
	}
	return &exec, nil
}

// ListTasks lists tasks for an execution
func (s *ReplicationService) ListTasks(executionID int64) ([]*api.ReplicationTask, error) {
	resp, err := s.client.Get(fmt.Sprintf("/replication/executions/%d/tasks", executionID), nil)
	if err != nil {
		return nil, err
	}

	var tasks []*api.ReplicationTask
	if err := s.client.DecodeResponse(resp, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}
	return tasks, nil
}

// GetTaskLog retrieves the log of a task
func (s *ReplicationService) GetTaskLog(executionID, taskID int64) (string, error) {
	resp, err := s.client.Get(fmt.Sprintf("/replication/executions/%d/tasks/%d/log", executionID, taskID), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read log: %w", err)
	}
	return string(buf), nil
}
