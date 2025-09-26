package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/tosharewith/postgres-age-operator/api-server/internal/api"
	"github.com/tosharewith/postgres-age-operator/api-server/internal/k8s"
)

func main() {
	// Command line flags
	var (
		port   = flag.String("port", getEnvOrDefault("PORT", "8080"), "Port to run the server on")
		apiKey = flag.String("api-key", os.Getenv("API_KEY"), "API key for authentication (optional)")
		debug  = flag.Bool("debug", getEnvOrDefault("DEBUG", "false") == "true", "Enable debug mode")
	)
	flag.Parse()

	// Set gin mode
	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Create router
	router := gin.New()

	// Setup routes
	api.SetupRoutes(router, k8sClient, *apiKey)

	// Log startup information
	log.Printf("Starting PostgreSQL AGE Operator API server...")
	log.Printf("Port: %s", *port)
	if *apiKey != "" {
		log.Printf("Authentication: Enabled (API Key required)")
	} else {
		log.Printf("Authentication: Disabled (Development mode)")
	}
	log.Printf("Debug mode: %v", *debug)

	// Start server
	if err := router.Run(":" + *port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}