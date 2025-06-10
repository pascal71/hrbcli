package harbor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pascal71/hrbcli/pkg/api"
)

// ProjectService handles project-related operations
type ProjectService struct {
	client *api.Client
}

// NewProjectService creates a new project service
func NewProjectService(client *api.Client) *ProjectService {
	return &ProjectService{client: client}
}

// List lists projects
func (s *ProjectService) List(opts *api.ListOptions) ([]*api.Project, error) {
	params := make(map[string]string)

	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.Query != "" {
			params["q"] = opts.Query
		}
		if opts.Sort != "" {
			params["sort"] = opts.Sort
		}
	}

	resp, err := s.client.Get("/projects", params)
	if err != nil {
		return nil, err
	}

	var projects []*api.Project
	if err := s.client.DecodeResponse(resp, &projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

// Get gets a project by name or ID
func (s *ProjectService) Get(nameOrID string) (*api.Project, error) {
	resp, err := s.client.Get(fmt.Sprintf("/projects/%s", nameOrID), nil)
	if err != nil {
		return nil, err
	}

	var project api.Project
	if err := s.client.DecodeResponse(resp, &project); err != nil {
		return nil, fmt.Errorf("failed to decode project: %w", err)
	}

	return &project, nil
}

// Create creates a new project
func (s *ProjectService) Create(req *api.ProjectReq) error {
	resp, err := s.client.Post("/projects", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Update updates a project
func (s *ProjectService) Update(nameOrID string, req *api.ProjectReq) error {
	resp, err := s.client.Put(fmt.Sprintf("/projects/%s", nameOrID), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Delete deletes a project
func (s *ProjectService) Delete(nameOrID string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/projects/%s", nameOrID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSummary gets project summary
func (s *ProjectService) GetSummary(nameOrID string) (*api.ProjectSummary, error) {
	resp, err := s.client.Get(fmt.Sprintf("/projects/%s/summary", nameOrID), nil)
	if err != nil {
		return nil, err
	}

	var summary api.ProjectSummary
	if err := s.client.DecodeResponse(resp, &summary); err != nil {
		return nil, fmt.Errorf("failed to decode project summary: %w", err)
	}

	return &summary, nil
}

// Exists checks if a project exists
func (s *ProjectService) Exists(projectName string) (bool, error) {
	resp, err := s.client.Head(fmt.Sprintf("/projects?project_name=%s", projectName))
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok && apiErr.IsNotFound() {
			return false, nil
		}
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// SearchByName searches projects by name
func (s *ProjectService) SearchByName(name string) ([]*api.Project, error) {
	opts := &api.ListOptions{
		Query: fmt.Sprintf("name=~%s", name),
	}
	return s.List(opts)
}

// ParseStorageLimit parses storage limit string (e.g., "10G", "500M") to bytes
func ParseStorageLimit(limit string) (int64, error) {
	if limit == "" || limit == "-1" {
		return -1, nil
	}

	limit = strings.ToUpper(strings.TrimSpace(limit))

	var multiplier int64 = 1
	var numStr string

	switch {
	case strings.HasSuffix(limit, "T"):
		multiplier = 1024 * 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(limit, "T")
	case strings.HasSuffix(limit, "G"):
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(limit, "G")
	case strings.HasSuffix(limit, "M"):
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(limit, "M")
	case strings.HasSuffix(limit, "K"):
		multiplier = 1024
		numStr = strings.TrimSuffix(limit, "K")
	default:
		numStr = limit
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid storage limit format: %s", limit)
	}

	return int64(num * float64(multiplier)), nil
}

// FormatStorageSize formats bytes to human readable format
func FormatStorageSize(bytes int64) string {
	if bytes < 0 {
		return "unlimited"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
