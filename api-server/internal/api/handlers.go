package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tosharewith/postgres-age-operator/api-server/internal/k8s"
	"github.com/tosharewith/postgres-age-operator/api-server/internal/models"
)

// Handler handles HTTP requests for customer management
type Handler struct {
	k8sClient *k8s.Client
}

// NewHandler creates a new API handler
func NewHandler(k8sClient *k8s.Client) *Handler {
	return &Handler{
		k8sClient: k8sClient,
	}
}

// CreateCustomer handles POST /api/v1/customers
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
				Details: []models.ValidationError{{Field: "body", Message: err.Error()}},
			},
		})
		return
	}

	// Validate request
	if err := h.validateCreateRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err,
		})
		return
	}

	// Create customer instance
	customer := &models.CustomerInstance{
		ID:          uuid.New().String(),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		ImageTag:    req.ImageTag,
		Config:      req.Config,
		Labels:      req.Labels,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status: models.InstanceStatus{
			Phase:       models.InstancePhaseCreating,
			Message:     "Creating customer instance",
			Ready:       false,
			LastUpdated: time.Now(),
		},
	}

	// Set default image tag if not provided
	if customer.ImageTag == "" {
		customer.ImageTag = "latest"
	}

	// Create in Kubernetes
	err := h.k8sClient.CreateCustomerInstance(c.Request.Context(), customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "CREATION_FAILED",
				Message: fmt.Sprintf("Failed to create customer instance: %s", err.Error()),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Customer instance created successfully",
		Data:    customer,
	})
}

// GetCustomer handles GET /api/v1/customers/:name
func (h *Handler) GetCustomer(c *gin.Context) {
	customerName := c.Param("name")

	customer, err := h.k8sClient.GetCustomerInstance(c.Request.Context(), customerName)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "CUSTOMER_NOT_FOUND",
					Message: fmt.Sprintf("Customer instance '%s' not found", customerName),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "RETRIEVAL_FAILED",
				Message: fmt.Sprintf("Failed to retrieve customer instance: %s", err.Error()),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    customer,
	})
}

// UpdateCustomer handles PUT /api/v1/customers/:name
func (h *Handler) UpdateCustomer(c *gin.Context) {
	customerName := c.Param("name")

	var req models.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
				Details: []models.ValidationError{{Field: "body", Message: err.Error()}},
			},
		})
		return
	}

	// Update in Kubernetes
	err := h.k8sClient.UpdateCustomerInstance(c.Request.Context(), customerName, &req)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "CUSTOMER_NOT_FOUND",
					Message: fmt.Sprintf("Customer instance '%s' not found", customerName),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "UPDATE_FAILED",
				Message: fmt.Sprintf("Failed to update customer instance: %s", err.Error()),
			},
		})
		return
	}

	// Get updated customer
	customer, err := h.k8sClient.GetCustomerInstance(c.Request.Context(), customerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "RETRIEVAL_FAILED",
				Message: "Customer updated but failed to retrieve updated data",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Customer instance updated successfully",
		Data:    customer,
	})
}

// DeleteCustomer handles DELETE /api/v1/customers/:name
func (h *Handler) DeleteCustomer(c *gin.Context) {
	customerName := c.Param("name")

	err := h.k8sClient.DeleteCustomerInstance(c.Request.Context(), customerName)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "CUSTOMER_NOT_FOUND",
					Message: fmt.Sprintf("Customer instance '%s' not found", customerName),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "DELETION_FAILED",
				Message: fmt.Sprintf("Failed to delete customer instance: %s", err.Error()),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Customer instance deleted successfully",
	})
}

// ListCustomers handles GET /api/v1/customers
func (h *Handler) ListCustomers(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get all customers
	customers, err := h.k8sClient.ListCustomerInstances(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "LIST_FAILED",
				Message: fmt.Sprintf("Failed to list customer instances: %s", err.Error()),
			},
		})
		return
	}

	// Apply pagination
	total := len(customers)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		customers = []models.CustomerInstance{}
	} else {
		if end > total {
			end = total
		}
		customers = customers[start:end]
	}

	response := models.CustomerListResponse{
		Customers: customers,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		HasNext:   end < total,
		HasPrev:   page > 1,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetCustomerStatus handles GET /api/v1/customers/:name/status
func (h *Handler) GetCustomerStatus(c *gin.Context) {
	customerName := c.Param("name")

	customer, err := h.k8sClient.GetCustomerInstance(c.Request.Context(), customerName)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "CUSTOMER_NOT_FOUND",
					Message: fmt.Sprintf("Customer instance '%s' not found", customerName),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: models.ErrorResponse{
				Code:    "STATUS_RETRIEVAL_FAILED",
				Message: fmt.Sprintf("Failed to retrieve customer status: %s", err.Error()),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    customer.Status,
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// validateCreateRequest validates the create customer request
func (h *Handler) validateCreateRequest(req *models.CreateCustomerRequest) *models.ErrorResponse {
	var errors []models.ValidationError

	if req.Name == "" {
		errors = append(errors, models.ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(req.Name) > 20 {
		errors = append(errors, models.ValidationError{
			Field:   "name",
			Message: "Name must be 20 characters or less",
		})
	} else if !isValidKubernetesName(req.Name) {
		errors = append(errors, models.ValidationError{
			Field:   "name",
			Message: "Name must be a valid Kubernetes name (lowercase alphanumeric with optional hyphens)",
		})
	}

	if len(errors) > 0 {
		return &models.ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Request validation failed",
			Details: errors,
		}
	}

	return nil
}

// isNotFoundError checks if an error indicates a resource was not found
func isNotFoundError(err error) bool {
	return err != nil && (
		err.Error() == "customer instance not found" ||
		err.Error() == "Customer instance not found")
}

// isValidKubernetesName validates if a string is a valid Kubernetes name
func isValidKubernetesName(name string) bool {
	if name == "" {
		return false
	}

	for i, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
		if i == 0 && r == '-' {
			return false
		}
		if i == len(name)-1 && r == '-' {
			return false
		}
	}
	return true
}