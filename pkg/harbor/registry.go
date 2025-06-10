package harbor

import (
	"fmt"
	"net/http"

	"github.com/pascal71/hrbcli/pkg/api"
)

// RegistryService handles registry endpoint operations
type RegistryService struct {
	client *api.Client
}

// NewRegistryService creates a new registry service
func NewRegistryService(client *api.Client) *RegistryService {
	return &RegistryService{client: client}
}

// List lists all registry endpoints
func (s *RegistryService) List(opts *api.ListOptions) ([]*api.Registry, error) {
	params := make(map[string]string)

	if opts != nil {
		if opts.Query != "" {
			params["q"] = opts.Query
		}
		if opts.Page > 0 {
			params["page"] = fmt.Sprintf("%d", opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = fmt.Sprintf("%d", opts.PageSize)
		}
	}

	resp, err := s.client.Get("/registries", params)
	if err != nil {
		return nil, err
	}

	var registries []*api.Registry
	if err := s.client.DecodeResponse(resp, &registries); err != nil {
		return nil, fmt.Errorf("failed to decode registries: %w", err)
	}

	return registries, nil
}

// Get gets a registry by ID
func (s *RegistryService) Get(id int64) (*api.Registry, error) {
	resp, err := s.client.Get(fmt.Sprintf("/registries/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var registry api.Registry
	if err := s.client.DecodeResponse(resp, &registry); err != nil {
		return nil, fmt.Errorf("failed to decode registry: %w", err)
	}

	return &registry, nil
}

// Create creates a new registry endpoint
func (s *RegistryService) Create(req *api.RegistryReq) (*api.Registry, error) {
	resp, err := s.client.Post("/registries", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Get location header to find the created registry ID
	location := resp.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("no location header in response")
	}

	// Extract ID from location (format: /api/v2.0/registries/{id})
	var id int64
	if _, err := fmt.Sscanf(location, "/api/v2.0/registries/%d", &id); err != nil {
		return nil, fmt.Errorf("failed to parse registry ID from location: %s", location)
	}

	// Get the created registry
	return s.Get(id)
}

// Update updates a registry endpoint
func (s *RegistryService) Update(id int64, req *api.RegistryReq) error {
	resp, err := s.client.Put(fmt.Sprintf("/registries/%d", id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Delete deletes a registry endpoint
func (s *RegistryService) Delete(id int64) error {
	resp, err := s.client.Delete(fmt.Sprintf("/registries/%d", id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Ping tests connectivity to a registry
func (s *RegistryService) Ping(req *api.RegistryReq) error {
	resp, err := s.client.Post("/registries/ping", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var pingResp api.RegistryPing
		if err := s.client.DecodeResponse(resp, &pingResp); err == nil && pingResp.Reason != "" {
			return fmt.Errorf("ping failed: %s", pingResp.Reason)
		}
		return fmt.Errorf("ping failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetInfo returns adapter information for the specified registry type.
// Harbor 2.6+ exposes adapter details via `/replication/adapterinfos`.
func (s *RegistryService) GetInfo(registryType string) (*api.RegistryInfo, error) {
	resp, err := s.client.Get("/replication/adapterinfos", nil)
	if err != nil {
		return nil, err
	}

	infos := make(map[string]*api.RegistryInfo)
	if err := s.client.DecodeResponse(resp, &infos); err != nil {
		return nil, fmt.Errorf("failed to decode adapter infos: %w", err)
	}

	info, ok := infos[registryType]
	if !ok {
		return nil, fmt.Errorf("adapter info %s not found", registryType)
	}

	return info, nil
}

// ListAdapters lists available registry adapters

func (s *RegistryService) ListAdapters() (map[string]*api.RegistryInfo, error) {
	resp, err := s.client.Get("/replication/adapterinfos", nil)
	if err != nil {
		return nil, err
	}

	var adapters map[string]*api.RegistryInfo
	if err := s.client.DecodeResponse(resp, &adapters); err != nil {
		return nil, fmt.Errorf("failed to decode adapters: %w", err)
	}

	return adapters, nil
}
