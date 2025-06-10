package harbor

import (
	"fmt"
	"net/http"

	"github.com/pascal71/hrbcli/pkg/api"
)

// ConfigService handles Harbor system configuration operations
// such as retrieving and updating settings.
type ConfigService struct {
	client *api.Client
}

// NewConfigService creates a new ConfigService
func NewConfigService(client *api.Client) *ConfigService {
	return &ConfigService{client: client}
}

// Get retrieves the Harbor system configuration as a key/value map
func (s *ConfigService) Get() (map[string]interface{}, error) {
	resp, err := s.client.Get("/configurations", nil)
	if err != nil {
		return nil, err
	}

	// Harbor API returns configuration in the form
	// {"setting": {"value": <any>, "editable": <bool>}}
	// Convert it to a simple key/value map for easier consumption.
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

// Update updates Harbor system configuration.
// The cfg map may contain a subset of settings to modify.
func (s *ConfigService) Update(cfg map[string]interface{}) error {
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
