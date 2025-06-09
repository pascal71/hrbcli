package api

import (
	"time"
)

// Project represents a Harbor project
type Project struct {
	ProjectID    int64            `json:"project_id"`
	Name         string           `json:"name"`
	Public       bool             `json:"public"`
	OwnerID      int              `json:"owner_id"`
	OwnerName    string           `json:"owner_name"`
	CreationTime time.Time        `json:"creation_time"`
	UpdateTime   time.Time        `json:"update_time"`
	Deleted      bool             `json:"deleted"`
	RepoCount    int64            `json:"repo_count"`
	Metadata     *ProjectMetadata `json:"metadata,omitempty"`
	CVEAllowlist *CVEAllowlist    `json:"cve_allowlist,omitempty"`
	StorageLimit int64            `json:"storage_limit,omitempty"`
}

// ProjectReq represents a project creation/update request
type ProjectReq struct {
	ProjectName  string           `json:"project_name"`
	Public       *bool            `json:"public,omitempty"`
	Metadata     *ProjectMetadata `json:"metadata,omitempty"`
	CVEAllowlist *CVEAllowlist    `json:"cve_allowlist,omitempty"`
	StorageLimit *int64           `json:"storage_limit,omitempty"`
	CountLimit   *int64           `json:"count_limit,omitempty"`
	RegistryID   *int64           `json:"registry_id,omitempty"`
}

// ProjectMetadata represents project metadata
type ProjectMetadata struct {
	Public               string `json:"public,omitempty"`
	EnableContentTrust   string `json:"enable_content_trust,omitempty"`
	PreventVul           string `json:"prevent_vul,omitempty"`
	Severity             string `json:"severity,omitempty"`
	AutoScan             string `json:"auto_scan,omitempty"`
	ReuseSysCVEAllowlist string `json:"reuse_sys_cve_allowlist,omitempty"`
	RetentionID          string `json:"retention_id,omitempty"`
	ProxySpeedKB         string `json:"proxy_speed_kb,omitempty"`
}

// CVEAllowlist represents a CVE allowlist
type CVEAllowlist struct {
	ID           int64              `json:"id,omitempty"`
	ProjectID    int64              `json:"project_id,omitempty"`
	Items        []CVEAllowlistItem `json:"items,omitempty"`
	CreationTime time.Time          `json:"creation_time,omitempty"`
	UpdateTime   time.Time          `json:"update_time,omitempty"`
}

// CVEAllowlistItem represents a CVE allowlist item
type CVEAllowlistItem struct {
	CVEID string `json:"cve_id"`
}

// ProjectSummary represents project summary information
type ProjectSummary struct {
	RepoCount         int64         `json:"repo_count"`
	ProjectAdminCount int64         `json:"project_admin_count"`
	MasterCount       int64         `json:"master_count"`
	DeveloperCount    int64         `json:"developer_count"`
	GuestCount        int64         `json:"guest_count"`
	LimitedGuestCount int64         `json:"limited_guest_count"`
	Quota             *ProjectQuota `json:"quota,omitempty"`
}

// ProjectQuota represents project quota information
type ProjectQuota struct {
	Hard QuotaHard `json:"hard"`
	Used QuotaUsed `json:"used"`
}

// QuotaHard represents hard quota limits
type QuotaHard struct {
	Storage int64 `json:"storage"`
	Count   int64 `json:"count"`
}

// QuotaUsed represents used quota
type QuotaUsed struct {
	Storage int64 `json:"storage"`
	Count   int64 `json:"count"`
}

// User represents a Harbor user
type User struct {
	UserID          int       `json:"user_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"password,omitempty"`
	Realname        string    `json:"realname"`
	Comment         string    `json:"comment"`
	Deleted         bool      `json:"deleted"`
	AdminRoleInAuth bool      `json:"admin_role_in_auth"`
	SysadminFlag    bool      `json:"sysadmin_flag"`
	CreationTime    time.Time `json:"creation_time"`
	UpdateTime      time.Time `json:"update_time"`

}

// UserReq represents a request to create or update a user
type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Realname string `json:"realname,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// UserProfile represents a user profile update
type UserProfile struct {
	Email    string `json:"email"`
	Realname string `json:"realname"`
	Comment  string `json:"comment"`
}

// SysAdminFlag is used to set or unset Harbor admin privilege
type SysAdminFlag struct {
	SysadminFlag bool `json:"sysadmin_flag"`

}

// SystemInfo represents Harbor system information
type SystemInfo struct {
	HarborVersion               string `json:"harbor_version"`
	RegistryURL                 string `json:"registry_url"`
	ExternalURL                 string `json:"external_url"`
	AuthMode                    string `json:"auth_mode"`
	ProjectCreationRestriction  string `json:"project_creation_restriction"`
	SelfRegistration            bool   `json:"self_registration"`
	HasCARoot                   bool   `json:"has_ca_root"`
	WithNotary                  bool   `json:"with_notary"`
	WithChartmuseum             bool   `json:"with_chartmuseum"`
	RegistryStorageProviderName string `json:"registry_storage_provider_name"`
}

// Pagination represents pagination information
type Pagination struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}

// ListOptions represents common list options
type ListOptions struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Query    string `json:"q,omitempty"`
	Sort     string `json:"sort,omitempty"`
}

// Repository represents a Harbor repository
type Repository struct {
	ID            int64     `json:"id"`
	ProjectID     int64     `json:"project_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ArtifactCount int64     `json:"artifact_count"`
	PullCount     int64     `json:"pull_count"`
	CreationTime  time.Time `json:"creation_time"`
	UpdateTime    time.Time `json:"update_time"`
}
