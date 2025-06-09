package api

import "time"

// Registry represents a registry endpoint in Harbor
type Registry struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	URL          string      `json:"url"`
	Description  string      `json:"description,omitempty"`
	Type         string      `json:"type"`
	Insecure     bool        `json:"insecure"`
	Credential   *Credential `json:"credential,omitempty"`
	Status       string      `json:"status,omitempty"`
	CreationTime time.Time   `json:"creation_time,omitempty"`
	UpdateTime   time.Time   `json:"update_time,omitempty"`
}

// RegistryReq represents a registry creation/update request
type RegistryReq struct {
	Name        string      `json:"name"`
	URL         string      `json:"url"`
	Description string      `json:"description,omitempty"`
	Type        string      `json:"type"`
	Insecure    bool        `json:"insecure"`
	Credential  *Credential `json:"credential,omitempty"`
}

// Credential represents registry credentials
type Credential struct {
	Type         string `json:"type"` // basic, oauth
	AccessKey    string `json:"access_key,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
}

// RegistryPing represents registry ping response
type RegistryPing struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// RegistryInfo represents registry information
type RegistryInfo struct {
	Type           string   `json:"type"`
	Description    string   `json:"description"`
	SupportedTypes []string `json:"supported_resource_types"`
	Filters        []string `json:"supported_resource_filters"`
	Triggers       []string `json:"supported_triggers"`
}

// ProxyCache represents proxy cache configuration
type ProxyCache struct {
	ID         int64  `json:"id"`
	ProjectID  int64  `json:"project_id"`
	RegistryID int64  `json:"registry_id"`
	SpeedKB    int64  `json:"speed_kb,omitempty"`
	Status     string `json:"status,omitempty"`
}

// Common registry types
const (
	RegistryTypeDockerHub      = "docker-hub"
	RegistryTypeDockerRegistry = "docker-registry"
	RegistryTypeHarbor         = "harbor"
	RegistryTypeAzureACR       = "azure-acr"
	RegistryTypeAWS            = "aws-ecr"
	RegistryTypeGoogleGCR      = "google-gcr"
	RegistryTypeQuay           = "quay"
	RegistryTypeGitlab         = "gitlab"
)
