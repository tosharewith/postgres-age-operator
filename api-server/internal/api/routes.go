package api

import (
	"github.com/gin-gonic/gin"

	"github.com/tosharewith/postgres-age-operator/api-server/internal/k8s"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, k8sClient *k8s.Client, apiKey string) {
	handler := NewHandler(k8sClient)

	// Add middleware
	router.Use(CORSMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(ErrorHandlingMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(RateLimitMiddleware())

	// Health check (no auth required)
	router.GET("/health", handler.HealthCheck)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service":     "PostgreSQL AGE Operator API",
			"version":     "1.0.0",
			"description": "REST API for managing PostgreSQL AGE Operator customer instances",
			"endpoints": gin.H{
				"health":     "GET /health",
				"customers":  "GET /api/v1/customers",
				"create":     "POST /api/v1/customers",
				"get":        "GET /api/v1/customers/:name",
				"update":     "PUT /api/v1/customers/:name",
				"delete":     "DELETE /api/v1/customers/:name",
				"status":     "GET /api/v1/customers/:name/status",
				"docs":       "GET /api/v1/docs",
			},
		})
	})

	// API v1 routes with authentication
	v1 := router.Group("/api/v1")
	v1.Use(AuthMiddleware(apiKey))
	{
		// Customer management
		customers := v1.Group("/customers")
		{
			customers.POST("", handler.CreateCustomer)           // Create customer instance
			customers.GET("", handler.ListCustomers)             // List all customer instances
			customers.GET("/:name", handler.GetCustomer)         // Get specific customer instance
			customers.PUT("/:name", handler.UpdateCustomer)      // Update customer instance
			customers.DELETE("/:name", handler.DeleteCustomer)   // Delete customer instance
			customers.GET("/:name/status", handler.GetCustomerStatus) // Get customer status
		}

		// API documentation
		v1.GET("/docs", func(c *gin.Context) {
			c.JSON(200, getAPIDocumentation())
		})
	}
}

// getAPIDocumentation returns API documentation
func getAPIDocumentation() gin.H {
	return gin.H{
		"title":       "PostgreSQL AGE Operator API",
		"version":     "1.0.0",
		"description": "REST API for managing PostgreSQL AGE Operator customer instances",
		"baseURL":     "/api/v1",
		"authentication": gin.H{
			"type":        "Bearer Token",
			"description": "Include 'Authorization: Bearer <api-key>' header in requests",
		},
		"endpoints": []gin.H{
			{
				"method":      "POST",
				"path":        "/customers",
				"description": "Create a new customer instance",
				"body": gin.H{
					"name":        "string (required, max 20 chars, lowercase alphanumeric with hyphens)",
					"displayName": "string (optional)",
					"imageTag":    "string (optional, default: 'latest')",
					"config": gin.H{
						"resources": gin.H{
							"requests": gin.H{"cpu": "string", "memory": "string"},
							"limits":   gin.H{"cpu": "string", "memory": "string"},
						},
						"highAvailability": "boolean",
						"backupEnabled":    "boolean",
						"monitoringEnabled": "boolean",
					},
					"labels": "object (optional key-value pairs)",
				},
			},
			{
				"method":      "GET",
				"path":        "/customers",
				"description": "List all customer instances",
				"parameters": gin.H{
					"page":     "int (optional, default: 1)",
					"pageSize": "int (optional, default: 10, max: 100)",
				},
			},
			{
				"method":      "GET",
				"path":        "/customers/:name",
				"description": "Get a specific customer instance",
			},
			{
				"method":      "PUT",
				"path":        "/customers/:name",
				"description": "Update a customer instance",
				"body": gin.H{
					"displayName": "string (optional)",
					"imageTag":    "string (optional)",
					"config":      "object (optional)",
					"labels":      "object (optional)",
				},
			},
			{
				"method":      "DELETE",
				"path":        "/customers/:name",
				"description": "Delete a customer instance",
			},
			{
				"method":      "GET",
				"path":        "/customers/:name/status",
				"description": "Get customer instance status",
			},
		},
		"examples": gin.H{
			"createCustomer": gin.H{
				"url":    "POST /api/v1/customers",
				"headers": gin.H{
					"Authorization": "Bearer your-api-key",
					"Content-Type":  "application/json",
				},
				"body": gin.H{
					"name":        "acme-corp",
					"displayName": "ACME Corporation",
					"imageTag":    "latest",
					"config": gin.H{
						"resources": gin.H{
							"requests": gin.H{"cpu": "200m", "memory": "256Mi"},
							"limits":   gin.H{"cpu": "1000m", "memory": "1Gi"},
						},
						"highAvailability":  true,
						"backupEnabled":     true,
						"monitoringEnabled": true,
					},
					"labels": gin.H{
						"environment": "production",
						"team":        "platform",
					},
				},
			},
			"listCustomers": gin.H{
				"url": "GET /api/v1/customers?page=1&pageSize=10",
				"headers": gin.H{
					"Authorization": "Bearer your-api-key",
				},
			},
			"getCustomer": gin.H{
				"url": "GET /api/v1/customers/acme-corp",
				"headers": gin.H{
					"Authorization": "Bearer your-api-key",
				},
			},
		},
		"responseFormat": gin.H{
			"success": gin.H{
				"success": true,
				"message": "optional success message",
				"data":    "response data",
			},
			"error": gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERROR_CODE",
					"message": "Error description",
					"details": "[]ValidationError (optional)",
				},
			},
		},
	}
}