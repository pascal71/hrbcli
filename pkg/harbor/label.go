package harbor

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pascal71/hrbcli/pkg/api"
)

// LabelService handles label operations
type LabelService struct {
	client *api.Client
}

// NewLabelService creates a new label service
func NewLabelService(client *api.Client) *LabelService {
	return &LabelService{client: client}
}

// List lists labels with optional filters
func (s *LabelService) List(opts *api.LabelListOptions) ([]*api.Label, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.Scope != "" {
			params["scope"] = opts.Scope
		}
		if opts.ProjectID > 0 {
			params["project_id"] = strconv.FormatInt(opts.ProjectID, 10)
		}
	}

	resp, err := s.client.Get("/labels", params)
	if err != nil {
		return nil, err
	}

	var labels []*api.Label
	if err := s.client.DecodeResponse(resp, &labels); err != nil {
		return nil, fmt.Errorf("failed to decode labels: %w", err)
	}

	return labels, nil
}

// Get retrieves a label by ID
func (s *LabelService) Get(id int64) (*api.Label, error) {
	resp, err := s.client.Get(fmt.Sprintf("/labels/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var label api.Label
	if err := s.client.DecodeResponse(resp, &label); err != nil {
		return nil, fmt.Errorf("failed to decode label: %w", err)
	}
	return &label, nil
}

// Create creates a new label
func (s *LabelService) Create(label *api.Label) (*api.Label, error) {
	resp, err := s.client.Post("/labels", label)
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
	if _, err := fmt.Sscanf(location, "/api/v2.0/labels/%d", &id); err != nil {
		return nil, fmt.Errorf("failed to parse label ID from location: %s", location)
	}

	return s.Get(id)
}

// Update updates an existing label
func (s *LabelService) Update(id int64, label *api.Label) error {
	resp, err := s.client.Put(fmt.Sprintf("/labels/%d", id), label)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Delete deletes a label by ID
func (s *LabelService) Delete(id int64) error {
	resp, err := s.client.Delete(fmt.Sprintf("/labels/%d", id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// AddToArtifact attaches a label to an artifact by reference
func (s *LabelService) AddToArtifact(project, repository, reference string, labelID int64) error {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	refEsc := url.PathEscape(reference)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/labels", projectEsc, repoEsc, refEsc)

	body := &api.Label{ID: labelID}
	resp, err := s.client.Post(path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// RemoveFromArtifact removes a label from an artifact
func (s *LabelService) RemoveFromArtifact(project, repository, reference string, labelID int64) error {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	refEsc := url.PathEscape(reference)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/labels/%d", projectEsc, repoEsc, refEsc, labelID)

	resp, err := s.client.Delete(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
