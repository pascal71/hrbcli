package harbor

import (
	"fmt"

	"github.com/pascal71/hrbcli/pkg/api"
)

// PreheatService handles distribution/preheat operations
// such as listing providers and policies.
type PreheatService struct {
	client *api.Client
}

// NewPreheatService creates a new PreheatService
func NewPreheatService(client *api.Client) *PreheatService {
	return &PreheatService{client: client}
}

// ListProviders lists distribution providers under a project
func (s *PreheatService) ListProviders(project string) ([]*api.PreheatProvider, error) {
	path := fmt.Sprintf("/projects/%s/preheat/providers", project)
	resp, err := s.client.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var providers []*api.PreheatProvider
	if err := s.client.DecodeResponse(resp, &providers); err != nil {
		return nil, fmt.Errorf("failed to decode providers: %w", err)
	}
	return providers, nil
}

// ListPolicies lists preheat policies under a project
func (s *PreheatService) ListPolicies(project string, opts *api.ListOptions) ([]*api.PreheatPolicy, error) {
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
	path := fmt.Sprintf("/projects/%s/preheat/policies", project)
	resp, err := s.client.Get(path, params)
	if err != nil {
		return nil, err
	}
	var policies []*api.PreheatPolicy
	if err := s.client.DecodeResponse(resp, &policies); err != nil {
		return nil, fmt.Errorf("failed to decode policies: %w", err)
	}
	return policies, nil
}

// GetPolicy retrieves a preheat policy
func (s *PreheatService) GetPolicy(project, name string) (*api.PreheatPolicy, error) {
	path := fmt.Sprintf("/projects/%s/preheat/policies/%s", project, name)
	resp, err := s.client.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var policy api.PreheatPolicy
	if err := s.client.DecodeResponse(resp, &policy); err != nil {
		return nil, fmt.Errorf("failed to decode policy: %w", err)
	}
	return &policy, nil
}
