package api

import "time"

// Label represents a Harbor label
// Only fields needed by the CLI are included.
type Label struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Color        string    `json:"color,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	ProjectID    int64     `json:"project_id,omitempty"`
	CreationTime time.Time `json:"creation_time,omitempty"`
	UpdateTime   time.Time `json:"update_time,omitempty"`
}

// LabelListOptions represents options when listing labels
// such as pagination or filtering by name and scope.
type LabelListOptions struct {
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	Name      string `json:"name,omitempty"`
	Scope     string `json:"scope,omitempty"`
	ProjectID int64  `json:"project_id,omitempty"`
}
