package models

import "time"

// Template management models

// Template represents a template in the system
type Template struct {
	Name         string    `json:"name"`
	Title        string    `json:"title"`
	Version      string    `json:"version"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Features     []string  `json:"features"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	PreviewURL   *string   `json:"preview_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TemplateVersion represents a specific version of a template
type TemplateVersion struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	GitCommit   string    `json:"git_commit"`
	BuildID     string    `json:"build_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Template build models

// TemplateBuildRequest represents a request to build a template
type TemplateBuildRequest struct {
	TemplateName string `json:"template_name" validate:"required"`
	GitCommit    string `json:"git_commit,omitempty"`
	Force        bool   `json:"force,omitempty"`
}

// BuildResponse represents the response from a build request
type BuildResponse struct {
	BuildID   string    `json:"build_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	StartedAt time.Time `json:"started_at"`
}

// Store deployment models

// StoreDeploymentRequest represents a request to deploy a store
type StoreDeploymentRequest struct {
	ShopID       string            `json:"shop_id"`
	Subdomain    string            `json:"subdomain"`
	TemplateName string            `json:"template_name" validate:"required"`
	Version      string            `json:"version,omitempty"`
	DataOverride map[string]string `json:"data_override,omitempty"`
}

// DeploymentResponse represents the response from a deployment request
type DeploymentResponse struct {
	DeploymentID string    `json:"deployment_id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	StartedAt    time.Time `json:"started_at"`
}

// Data update models

// DataUpdateRequest represents a request to update store data
type DataUpdateRequest struct {
	ShopID      string            `json:"shop_id"`
	Subdomain   string            `json:"subdomain"`
	DataType    string            `json:"data_type,omitempty"` // "all", "products", "orders", etc.
	Incremental bool              `json:"incremental,omitempty"`
	Changes     map[string]string `json:"changes,omitempty"`
}

// UpdateResponse represents the response from a data update request
type UpdateResponse struct {
	UpdateID  string    `json:"update_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	StartedAt time.Time `json:"started_at"`
}

// Status models

// DeploymentStatus represents the current deployment status for a store
type DeploymentStatus struct {
	ShopID          string     `json:"shop_id"`
	Subdomain       string     `json:"subdomain"`
	Status          string     `json:"status"` // "deployed", "deploying", "failed", "not_deployed"
	TemplateName    string     `json:"template_name"`
	TemplateVersion string     `json:"template_version"`
	LastDeployedAt  *time.Time `json:"last_deployed_at,omitempty"`
	LastUpdateAt    *time.Time `json:"last_update_at,omitempty"`
	BuildID         string     `json:"build_id,omitempty"`
	DeploymentID    string     `json:"deployment_id,omitempty"`
	Message         string     `json:"message,omitempty"`
	AssetsURL       string     `json:"assets_url,omitempty"`
	PreviewURL      string     `json:"preview_url,omitempty"`
	ProductionURL   string     `json:"production_url,omitempty"`
}

// TemplateDownload represents the response from a template download request
type TemplateDownload struct {
	TemplateName string `json:"template_name"`
	Version      string `json:"version"`
	DownloadURL  string `json:"download_url"`
	AssetSize    int64  `json:"asset_size,omitempty"`
	ContentType  string `json:"content_type,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}

// StoreRedeploymentRequest represents a request to redeploy an existing store
type StoreRedeploymentRequest struct {
	ShopID       string            `json:"shop_id"`
	Subdomain    string            `json:"subdomain"`
	TemplateName string            `json:"template_name,omitempty"`
	Version      string            `json:"version,omitempty"`
	DataOverride map[string]string `json:"data_override,omitempty"`
	Force        bool              `json:"force,omitempty"`
}
