package middleware

import (
	"net/http"
	"strings"

	"lms-go-be/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from Authorization header
// It verifies the token and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedError(c, "authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedError(c, "invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			utils.UnauthorizedError(c, err.Error())
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("user_id", claims.ID)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if user has required role(s)
// Parameters:
//   - allowedRoles: list of allowed roles
//
// Returns: gin.HandlerFunc middleware
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.UnauthorizedError(c, "user role not found in context")
			c.Abort()
			return
		}

		// Check if user's role is in allowed roles
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ForbiddenError(c, "user does not have permission for this action")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware configures CORS headers
// This allows requests from frontend applications
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandlingMiddleware handles panics and errors gracefully
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "internal server error",
					"error":   err,
				})
			}
		}()

		c.Next()

		// Handle HTTP errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "request failed",
				"error":   err.Error(),
			})
		}
	}
}

// LoggingMiddleware logs request details
// This is useful for debugging and monitoring
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request information
		method := c.Request.Method
		path := c.Request.RequestURI
		userID, _ := c.Get("user_id")

		c.Next()

		statusCode := c.Writer.Status()
		// In production, use proper logging library
		// For now, we'll skip logging to keep output clean
		_ = statusCode
		_ = method
		_ = path
		_ = userID
	}
}
