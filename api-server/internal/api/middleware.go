package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tosharewith/postgres-age-operator/api-server/internal/models"
)

// AuthMiddleware provides simple API key authentication
func AuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" {
			// If no API key is configured, allow all requests (development mode)
			c.Next()
			return
		}

		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "MISSING_AUTHORIZATION",
					Message: "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "INVALID_AUTHORIZATION_FORMAT",
					Message: "Authorization header must be 'Bearer <token>'",
				},
			})
			c.Abort()
			return
		}

		// Validate API key
		if parts[1] != apiKey {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "INVALID_API_KEY",
					Message: "Invalid API key",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware provides simple rate limiting
// In production, you'd want to use a more sophisticated solution like Redis
func RateLimitMiddleware() gin.HandlerFunc {
	// Simple in-memory rate limiting (not suitable for production with multiple instances)
	requests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old requests (older than 1 minute)
		if times, exists := requests[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					validTimes = append(validTimes, t)
				}
			}
			requests[clientIP] = validTimes
		}

		// Check rate limit (100 requests per minute per IP)
		if len(requests[clientIP]) >= 100 {
			c.JSON(http.StatusTooManyRequests, models.APIResponse{
				Success: false,
				Error: models.ErrorResponse{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "Too many requests. Limit: 100 requests per minute per IP",
				},
			})
			c.Abort()
			return
		}

		// Add current request
		requests[clientIP] = append(requests[clientIP], now)
		c.Next()
	}
}

// LoggingMiddleware provides request logging
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
	})
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}