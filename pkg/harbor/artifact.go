package harbor

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pascal71/hrbcli/pkg/api"
)

// ArtifactService handles artifact related operations
type ArtifactService struct {
	client *api.Client
}

// NewArtifactService creates a new artifact service
func NewArtifactService(client *api.Client) *ArtifactService {
	return &ArtifactService{client: client}
}

// List lists artifacts in a repository
func (s *ArtifactService) List(project, repository string, opts *api.ArtifactListOptions) ([]*api.Artifact, error) {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts", projectEsc, repoEsc)

	params := make(map[string]string)
	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.WithTag {
			params["with_tag"] = "true"
		}
		if opts.WithLabel {
			params["with_label"] = "true"
		}
		if opts.WithSignature {
			params["with_signature"] = "true"
		}
		if opts.WithScanOverview {
			params["with_scan_overview"] = "true"
		}
	}

	resp, err := s.client.Get(path, params)
	if err != nil {
		return nil, err
	}

	var artifacts []*api.Artifact
	if err := s.client.DecodeResponse(resp, &artifacts); err != nil {
		return nil, fmt.Errorf("failed to decode artifacts: %w", err)
	}

	return artifacts, nil
}

// Get retrieves details of a specific artifact
func (s *ArtifactService) Get(project, repository, reference string) (*api.Artifact, error) {
	opts := &api.ArtifactGetOptions{
		WithTag:       true,
		WithLabel:     true,
		WithSignature: true,
	}
	return s.GetWithOptions(project, repository, reference, opts)
}

// GetWithOptions retrieves details of a specific artifact with optional parameters
func (s *ArtifactService) GetWithOptions(project, repository, reference string, opts *api.ArtifactGetOptions) (*api.Artifact, error) {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	refEsc := url.PathEscape(reference)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s", projectEsc, repoEsc, refEsc)

	params := make(map[string]string)
	if opts != nil {
		if opts.WithTag {
			params["with_tag"] = "true"
		}
		if opts.WithLabel {
			params["with_label"] = "true"
		}
		if opts.WithSignature {
			params["with_signature"] = "true"
		}
		if opts.WithScanOverview {
			params["with_scan_overview"] = "true"
		}
	}

	resp, err := s.client.Get(path, params)
	if err != nil {
		return nil, err
	}

	var art api.Artifact
	if err := s.client.DecodeResponse(resp, &art); err != nil {
		return nil, fmt.Errorf("failed to decode artifact: %w", err)
	}

	return &art, nil
}

// Scan triggers vulnerability scan for the specified artifact
func (s *ArtifactService) Scan(project, repository, reference string, scanType string) error {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	refEsc := url.PathEscape(reference)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/scan", projectEsc, repoEsc, refEsc)

	var body interface{}
	if scanType != "" {
		body = map[string]string{"scan_type": scanType}
	}

	resp, err := s.client.Post(path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Vulnerabilities retrieves the vulnerability report for the specified artifact
func (s *ArtifactService) Vulnerabilities(project, repository, reference string) (*api.VulnerabilityReport, error) {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	refEsc := url.PathEscape(reference)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/additions/vulnerabilities", projectEsc, repoEsc, refEsc)

	resp, err := s.client.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var report api.VulnerabilityReport
	if err := s.client.DecodeResponse(resp, &report); err != nil {
		return nil, fmt.Errorf("failed to decode vulnerability report: %w", err)
	}

	return &report, nil
}

// SBOM retrieves the SBOM report for the specified artifact
func (s *ArtifactService) SBOM(project, repository, reference string) (map[string]interface{}, error) {
	projectEsc := url.PathEscape(project)
	repoEsc := url.PathEscape(repository)
	refEsc := url.PathEscape(reference)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/additions/sbom", projectEsc, repoEsc, refEsc)

	resp, err := s.client.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var report map[string]interface{}
	if err := s.client.DecodeResponse(resp, &report); err != nil {
		return nil, fmt.Errorf("failed to decode SBOM: %w", err)
	}

	return report, nil
}
