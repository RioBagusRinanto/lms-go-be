package middleware

import (
	"net/http"
	"strings"

	"lms-go-be/internal/config"
	"lms-go-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is the JWT authentication middleware
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "missing authorization header")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify token
		claims, err := utils.VerifyToken(tokenString, cfg.JWT.SecretKey)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "invalid token")
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("full_name", claims.FullName)

		c.Next()
	}
}

// RoleMiddleware checks if user has the required role
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", "user role not found")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", "invalid user role")
			c.Abort()
			return
		}

		// Check if user role is in required roles
		allowed := false
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware adds a request ID to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("X-Request-ID")
		if requestID == "" {
			requestID = c.GetHeader("X-Request-ID")
		}
		if requestID != "" {
			c.Set("request_id", requestID)
		}
		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", "an unexpected error occurred")
			}
		}()
		c.Next()
	}
}
