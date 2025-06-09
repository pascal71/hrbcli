package harbor

import (
	"fmt"
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

	path := fmt.Sprintf("/projects/%s/repositories", projectName)
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
