package harbor

import (
	"fmt"
	"net/http"

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
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s/scan", project, repository, reference)

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
