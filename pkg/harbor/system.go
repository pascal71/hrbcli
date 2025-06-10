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

// GetInfo retrieves Harbor system information.
func (s *SystemService) GetInfo(withStorage bool) (*api.SystemInfo, error) {
	params := make(map[string]string)
	if withStorage {
		params["with_storage"] = "true"
	}

	resp, err := s.client.Get("/systeminfo", params)
	if err != nil {
		return nil, err
	}

	var info api.SystemInfo
	if err := s.client.DecodeResponse(resp, &info); err != nil {
		return nil, fmt.Errorf("failed to decode system info: %w", err)
	}

	return &info, nil
}

// GetConfig retrieves Harbor configuration settings.
func (s *SystemService) GetConfig() (map[string]interface{}, error) {
	resp, err := s.client.Get("/configurations", nil)
	if err != nil {
		return nil, err
	}

	// API returns map of objects with `value` and `editable` fields.
	raw := make(map[string]struct {
		Value    interface{} `json:"value"`
		Editable bool        `json:"editable"`
	})
	if err := s.client.DecodeResponse(resp, &raw); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	cfg := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		cfg[k] = v.Value
	}
	return cfg, nil
}

// UpdateConfig updates Harbor configuration settings.
func (s *SystemService) UpdateConfig(cfg map[string]interface{}) error {
	resp, err := s.client.Put("/configurations", cfg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
