package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pascal71/hrbcli/pkg/config"
	"github.com/pascal71/hrbcli/pkg/output"
)

// Client represents a Harbor API client
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	APIVersion string
	HTTPClient *http.Client
	Debug      bool
}

// NewClient creates a new Harbor API client
func NewClient() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.HarborURL == "" {
		return nil, fmt.Errorf("Harbor URL not configured")
	}

	// Ensure URL has no trailing slash
	baseURL := strings.TrimRight(cfg.HarborURL, "/")

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Configure TLS
	if cfg.Insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := &Client{
		BaseURL:    baseURL,
		Username:   cfg.Username,
		Password:   cfg.Password,
		APIVersion: cfg.APIVersion,
		HTTPClient: httpClient,
		Debug:      cfg.Debug,
	}

	// Ensure global debug matches configuration
	output.SetDebug(cfg.Debug)

	return client, nil
}

// Request makes an HTTP request to the Harbor API
func (c *Client) Request(method, path string, body interface{}) (*http.Response, error) {
	// Build URL
	fullURL := fmt.Sprintf("%s/api/%s%s", c.BaseURL, c.APIVersion, path)

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)

		if c.Debug {
			output.Debug("Request body: %s", string(jsonBody))
		}
	}

	// Create request
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set basic auth
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	if c.Debug {
		output.Debug("%s %s", method, fullURL)
	}

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Provide clearer error for TLS verification failures
		if urlErr, ok := err.(*url.Error); ok {
			switch urlErr.Err.(type) {
			case x509.UnknownAuthorityError, x509.HostnameError, x509.CertificateInvalidError:
				return nil, fmt.Errorf("TLS certificate verification failed: %w", urlErr.Err)
			}
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if c.Debug {
		output.Debug("Response status: %s", resp.Status)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)

		var apiErr APIError
		if err := json.Unmarshal(bodyBytes, &apiErr); err == nil && apiErr.Message != "" {
			return nil, &apiErr
		}

		return nil, &APIError{
			Code:    resp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	return resp, nil
}

// Get makes a GET request
func (c *Client) Get(path string, params map[string]string) (*http.Response, error) {
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		path = path + "?" + values.Encode()
	}
	return c.Request("GET", path, nil)
}

// Post makes a POST request
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.Request("POST", path, body)
}

// Put makes a PUT request
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.Request("PUT", path, body)
}

// Delete makes a DELETE request
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Request("DELETE", path, nil)
}

// Head makes a HEAD request
func (c *Client) Head(path string) (*http.Response, error) {
	return c.Request("HEAD", path, nil)
}

// DecodeResponse decodes a JSON response
func (c *Client) DecodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if c.Debug {
		// Read body for debugging
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		output.Debug("Response body: %s", string(bodyBytes))

		// Decode from bytes
		return json.Unmarshal(bodyBytes, v)
	}

	// Direct decode
	return json.NewDecoder(resp.Body).Decode(v)
}

// CheckHealth checks if Harbor is healthy
func (c *Client) CheckHealth() error {
	resp, err := c.Get("/health", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %s", resp.Status)
	}

	return nil
}
