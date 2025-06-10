package api

import "time"

// Schedule represents a job schedule configuration.
type Schedule struct {
	Type string `json:"type"`
	Cron string `json:"cron,omitempty"`
}

// GCHistory represents a garbage collection execution record.
type GCHistory struct {
	ID            int64     `json:"id"`
	JobName       string    `json:"job_name"`
	JobKind       string    `json:"job_kind"`
	JobParameters string    `json:"job_parameters"`
	Schedule      *Schedule `json:"schedule,omitempty"`
	JobStatus     string    `json:"job_status"`
	Deleted       bool      `json:"deleted"`
	CreationTime  time.Time `json:"creation_time"`
	UpdateTime    time.Time `json:"update_time"`
}
