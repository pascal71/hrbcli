package harbor

import (
	"fmt"
	"net/http"

	"github.com/pascal71/hrbcli/pkg/api"
)

// SystemService handles system related operations
type SystemService struct {
	client *api.Client
}

// NewSystemService creates a new SystemService
func NewSystemService(client *api.Client) *SystemService {
	return &SystemService{client: client}
}

// GetStatistics retrieves Harbor statistics
func (s *SystemService) GetStatistics() (*api.Statistic, error) {
	resp, err := s.client.Get("/statistics", nil)
	if err != nil {
		return nil, err
	}

	var stats api.Statistic
	if err := s.client.DecodeResponse(resp, &stats); err != nil {
		return nil, fmt.Errorf("failed to decode statistics: %w", err)
	}

	return &stats, nil
}

// StartGC triggers a manual garbage collection job.
// It returns the ID of the GC job if available.
func (s *SystemService) StartGC() (int64, error) {
	body := map[string]interface{}{
		"schedule": api.Schedule{Type: "Manual"},
	}
	resp, err := s.client.Post("/system/gc/schedule", body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return 0, nil
	}
	var id int64
	fmt.Sscanf(location, "/api/%*s/system/gc/%d", &id)
	return id, nil
}

// GCHistory retrieves past garbage collection executions
func (s *SystemService) GCHistory(opts *api.ListOptions) ([]*api.GCHistory, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.Page > 0 {
			params["page"] = fmt.Sprintf("%d", opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = fmt.Sprintf("%d", opts.PageSize)
		}
	}
	resp, err := s.client.Get("/system/gc", params)
	if err != nil {
		return nil, err
	}

	var history []*api.GCHistory
	if err := s.client.DecodeResponse(resp, &history); err != nil {
		return nil, fmt.Errorf("failed to decode gc history: %w", err)
	}

	return history, nil
}

// GCStatus retrieves a garbage collection job by ID
func (s *SystemService) GCStatus(id int64) (*api.GCHistory, error) {
	resp, err := s.client.Get(fmt.Sprintf("/system/gc/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var gc api.GCHistory
	if err := s.client.DecodeResponse(resp, &gc); err != nil {
		return nil, fmt.Errorf("failed to decode gc status: %w", err)
	}

	return &gc, nil
}
