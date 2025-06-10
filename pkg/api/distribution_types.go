package api

import "time"

// PreheatProvider represents a distribution provider under a project
// Only the fields used by the CLI are included.
type PreheatProvider struct {
	ID       int64  `json:"id"`
	Provider string `json:"provider"`
	Enabled  bool   `json:"enabled"`
	Default  bool   `json:"default"`
}

// PreheatPolicy represents a preheat policy
// Only fields needed for listing and getting are included.
type PreheatPolicy struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	ProviderID   int64     `json:"provider_id,omitempty"`
	ProviderName string    `json:"provider_name,omitempty"`
	Filters      string    `json:"filters,omitempty"`
	Trigger      string    `json:"trigger,omitempty"`
	Enabled      bool      `json:"enabled"`
	CreationTime time.Time `json:"creation_time,omitempty"`
	UpdateTime   time.Time `json:"update_time,omitempty"`
}
