package harbor

import (
	"fmt"

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
