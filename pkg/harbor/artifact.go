package harbor

import (
	"fmt"
	"net/http"
	"net/url"

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
