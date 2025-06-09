package api

import "time"

// ReplicationPolicy represents a replication policy
// Only fields used by the CLI are included
type ReplicationPolicy struct {
	ID                int64               `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description,omitempty"`
	SrcRegistry       *Registry           `json:"src_registry,omitempty"`
	DestRegistry      *Registry           `json:"dest_registry,omitempty"`
	DestNamespace     string              `json:"dest_namespace,omitempty"`
	Trigger           *ReplicationTrigger `json:"trigger,omitempty"`
	Filters           []ReplicationFilter `json:"filters,omitempty"`
	ReplicateDeletion bool                `json:"replicate_deletion,omitempty"`
	Override          bool                `json:"override"`
	Enabled           bool                `json:"enabled"`
	CreationTime      time.Time           `json:"creation_time,omitempty"`
	UpdateTime        time.Time           `json:"update_time,omitempty"`
}

// ReplicationTrigger represents policy trigger
type ReplicationTrigger struct {
	Type            string                      `json:"type"`
	TriggerSettings *ReplicationTriggerSettings `json:"trigger_settings,omitempty"`
}

// ReplicationTriggerSettings represents trigger settings
type ReplicationTriggerSettings struct {
	Cron string `json:"cron,omitempty"`
}

// ReplicationFilter represents replication filter
type ReplicationFilter struct {
	Type       string      `json:"type"`
	Value      interface{} `json:"value"`
	Decoration string      `json:"decoration,omitempty"`
}

// ReplicationExecution represents a replication execution
type ReplicationExecution struct {
	ID         int64     `json:"id"`
	PolicyID   int64     `json:"policy_id"`
	Status     string    `json:"status"`
	Trigger    string    `json:"trigger"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	StatusText string    `json:"status_text"`
	Total      int       `json:"total"`
	Failed     int       `json:"failed"`
	Succeed    int       `json:"succeed"`
	InProgress int       `json:"in_progress"`
	Stopped    int       `json:"stopped"`
}

// StartReplicationExecution represents start execution request
type StartReplicationExecution struct {
	PolicyID int64 `json:"policy_id"`
}

// ReplicationTask represents a replication task
type ReplicationTask struct {
	ID           int64     `json:"id"`
	ExecutionID  int64     `json:"execution_id"`
	Status       string    `json:"status"`
	JobID        string    `json:"job_id"`
	Operation    string    `json:"operation"`
	ResourceType string    `json:"resource_type"`
	SrcResource  string    `json:"src_resource"`
	DstResource  string    `json:"dst_resource"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}
