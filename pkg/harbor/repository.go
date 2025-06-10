package harbor

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pascal71/hrbcli/pkg/api"
)

// RepositoryService handles repository operations
type RepositoryService struct {
	client *api.Client
}

// NewRepositoryService creates a new RepositoryService
func NewRepositoryService(client *api.Client) *RepositoryService {
	return &RepositoryService{client: client}
}

// List lists repositories within a project
func (s *RepositoryService) List(projectName string, opts *api.ListOptions) ([]*api.Repository, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.Query != "" {
			params["q"] = opts.Query
		}
		if opts.Sort != "" {
			params["sort"] = opts.Sort
		}
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
	}

	projectEsc := url.PathEscape(projectName)
	path := fmt.Sprintf("/projects/%s/repositories", projectEsc)
	resp, err := s.client.Get(path, params)
	if err != nil {
		return nil, err
	}

	var repos []*api.Repository
	if err := s.client.DecodeResponse(resp, &repos); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}

	return repos, nil
}

// Get retrieves a repository by name within a project
func (s *RepositoryService) Get(projectName, repositoryName string) (*api.Repository, error) {
	projectEsc := url.PathEscape(projectName)
	repoEsc := url.PathEscape(repositoryName)
	path := fmt.Sprintf("/projects/%s/repositories/%s", projectEsc, repoEsc)

	resp, err := s.client.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var repo api.Repository
	if err := s.client.DecodeResponse(resp, &repo); err != nil {
		return nil, fmt.Errorf("failed to decode repository: %w", err)
	}

	return &repo, nil
}

// Delete removes a repository from a project
func (s *RepositoryService) Delete(projectName, repositoryName string) error {
	projectEsc := url.PathEscape(projectName)
	repoEsc := url.PathEscape(repositoryName)
	path := fmt.Sprintf("/projects/%s/repositories/%s", projectEsc, repoEsc)

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

// ListTags lists all tags for a repository
func (s *RepositoryService) ListTags(projectName, repositoryName string, opts *api.ListOptions) ([]*api.ArtifactTag, error) {
	projectEsc := url.PathEscape(projectName)
	repoEsc := url.PathEscape(repositoryName)
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts", projectEsc, repoEsc)

	params := map[string]string{"with_tag": "true"}
	if opts != nil {
		if opts.Query != "" {
			params["q"] = opts.Query
		}
		if opts.Sort != "" {
			params["sort"] = opts.Sort
		}
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
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

	var tags []*api.ArtifactTag
	for _, a := range artifacts {
		for _, t := range a.Tags {
			tag := t
			tags = append(tags, &api.ArtifactTag{Name: tag.Name, Immutable: tag.Immutable})
		}
	}
	return tags, nil
}
