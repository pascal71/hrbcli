package harbor

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pascal71/hrbcli/pkg/api"
)

// UserService handles user-related operations
type UserService struct {
	client *api.Client
}

// NewUserService creates a new user service
func NewUserService(client *api.Client) *UserService {
	return &UserService{client: client}
}

// List lists users
func (s *UserService) List(opts *api.ListOptions) ([]*api.User, error) {
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
	resp, err := s.client.Get("/users", params)
	if err != nil {
		return nil, err
	}
	var users []*api.User
	if err := s.client.DecodeResponse(resp, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}
	return users, nil
}

// Search searches users by username
func (s *UserService) Search(username string, opts *api.ListOptions) ([]*api.User, error) {
	params := map[string]string{"username": username}
	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
	}
	resp, err := s.client.Get("/users/search", params)
	if err != nil {
		return nil, err
	}
	var users []*api.User
	if err := s.client.DecodeResponse(resp, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}
	return users, nil
}

// Get retrieves a user by ID
func (s *UserService) Get(id int64) (*api.User, error) {
	resp, err := s.client.Get(fmt.Sprintf("/users/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var user api.User
	if err := s.client.DecodeResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	return &user, nil
}

// GetByUsername fetches a user by username
func (s *UserService) GetByUsername(username string) (*api.User, error) {
	users, err := s.Search(username, nil)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username {
			return s.Get(int64(u.UserID))
		}
	}
	return nil, fmt.Errorf("user '%s' not found", username)
}

// Create creates a new user
func (s *UserService) Create(req *api.UserReq) (*api.User, error) {
	resp, err := s.client.Post("/users", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	if location == "" {
		return nil, nil
	}
	var id int64
	fmt.Sscanf(location, "/api/%*s/users/%d", &id)
	return s.Get(id)
}

// Delete deletes a user by ID
func (s *UserService) Delete(id int64) error {
	resp, err := s.client.Delete(fmt.Sprintf("/users/%d", id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// SetAdmin toggles admin role for a user
func (s *UserService) SetAdmin(id int64, admin bool) error {
	flag := &api.SysAdminFlag{SysadminFlag: admin}
	resp, err := s.client.Put(fmt.Sprintf("/users/%d/sysadmin", id), flag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// UpdateProfile updates user's profile fields
func (s *UserService) UpdateProfile(id int64, profile *api.UserProfile) error {
	resp, err := s.client.Put(fmt.Sprintf("/users/%d", id), profile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
