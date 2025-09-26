package models

import (
	"time"
)

// CustomerInstance represents a customer's PostgreSQL AGE operator instance
type CustomerInstance struct {
	ID          string            `json:"id"`
	Name        string            `json:"name" binding:"required"`
	DisplayName string            `json:"displayName,omitempty"`
	ImageTag    string            `json:"imageTag,omitempty"`
	Namespace   string            `json:"namespace"`
	Status      InstanceStatus    `json:"status"`
	Config      CustomerConfig    `json:"config,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// CustomerConfig holds configuration options for a customer instance
type CustomerConfig struct {
	Resources     ResourceConfig    `json:"resources,omitempty"`
	Storage       StorageConfig     `json:"storage,omitempty"`
	HighAvailability bool           `json:"highAvailability,omitempty"`
	BackupEnabled bool             `json:"backupEnabled,omitempty"`
	MonitoringEnabled bool          `json:"monitoringEnabled,omitempty"`
	CustomSettings map[string]interface{} `json:"customSettings,omitempty"`
}

// ResourceConfig defines resource limits and requests
type ResourceConfig struct {
	Requests ResourceRequirements `json:"requests,omitempty"`
	Limits   ResourceRequirements `json:"limits,omitempty"`
}

// ResourceRequirements defines CPU and memory requirements
type ResourceRequirements struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// StorageConfig defines storage requirements
type StorageConfig struct {
	Size         string `json:"size,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
	BackupSize   string `json:"backupSize,omitempty"`
}

// InstanceStatus represents the current status of a customer instance
type InstanceStatus struct {
	Phase        InstancePhase `json:"phase"`
	Message      string        `json:"message,omitempty"`
	Ready        bool          `json:"ready"`
	Replicas     int32         `json:"replicas,omitempty"`
	ReadyReplicas int32        `json:"readyReplicas,omitempty"`
	LastUpdated  time.Time     `json:"lastUpdated"`
}

// InstancePhase represents the phase of an instance
type InstancePhase string

const (
	InstancePhaseCreating   InstancePhase = "Creating"
	InstancePhaseRunning    InstancePhase = "Running"
	InstancePhaseFailed     InstancePhase = "Failed"
	InstancePhaseDeleting   InstancePhase = "Deleting"
	InstancePhaseDeleted    InstancePhase = "Deleted"
	InstancePhaseUpdating   InstancePhase = "Updating"
)

// CreateCustomerRequest represents a request to create a new customer instance
type CreateCustomerRequest struct {
	Name        string         `json:"name" binding:"required"`
	DisplayName string         `json:"displayName,omitempty"`
	ImageTag    string         `json:"imageTag,omitempty"`
	Config      CustomerConfig `json:"config,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// UpdateCustomerRequest represents a request to update a customer instance
type UpdateCustomerRequest struct {
	DisplayName string         `json:"displayName,omitempty"`
	ImageTag    string         `json:"imageTag,omitempty"`
	Config      CustomerConfig `json:"config,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// CustomerListResponse represents a paginated list of customer instances
type CustomerListResponse struct {
	Customers  []CustomerInstance `json:"customers"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	HasNext    bool               `json:"hasNext"`
	HasPrev    bool               `json:"hasPrev"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Details []ValidationError  `json:"details,omitempty"`
}