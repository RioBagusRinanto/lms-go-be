package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is a standardized API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination struct {
		Total      int64 `json:"total"`
		Page       int   `json:"page"`
		PageSize   int   `json:"page_size"`
		TotalPages int   `json:"total_pages"`
	} `json:"pagination"`
}

// SuccessResponse sends a success response to the client
// Parameters:
//   - c: Gin context
//   - statusCode: HTTP status code
//   - message: success message
//   - data: response data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response to the client
// Parameters:
//   - c: Gin context
//   - statusCode: HTTP status code
//   - message: error message
//   - errorDetail: detailed error information
func ErrorResponse(c *gin.Context, statusCode int, message string, errorDetail string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   errorDetail,
	})
}

// PaginationSuccessResponse sends a paginated success response
// Parameters:
//   - c: Gin context
//   - statusCode: HTTP status code
//   - message: success message
//   - data: response data
//   - total: total number of records
//   - page: current page number
//   - pageSize: number of records per page
func PaginationSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	response.Pagination.Total = total
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalPages = totalPages

	c.JSON(statusCode, response)
}

// BadRequestError sends a 400 Bad Request error
// Parameters:
//   - c: Gin context
//   - message: error message
//   - errorDetail: detailed error information
func BadRequestError(c *gin.Context, message string, errorDetail string) {
	ErrorResponse(c, http.StatusBadRequest, message, errorDetail)
}

// UnauthorizedError sends a 401 Unauthorized error
// Parameters:
//   - c: Gin context
//   - message: error message
func UnauthorizedError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, "authentication required")
}

// ForbiddenError sends a 403 Forbidden error
// Parameters:
//   - c: Gin context
//   - message: error message
func ForbiddenError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, "access denied")
}

// NotFoundError sends a 404 Not Found error
// Parameters:
//   - c: Gin context
//   - message: error message
func NotFoundError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, "resource not found")
}

// InternalServerError sends a 500 Internal Server Error
// Parameters:
//   - c: Gin context
//   - message: error message
//   - errorDetail: detailed error information
func InternalServerError(c *gin.Context, message string, errorDetail string) {
	ErrorResponse(c, http.StatusInternalServerError, message, errorDetail)
}

// ConflictError sends a 409 Conflict error
// Parameters:
//   - c: Gin context
//   - message: error message
//   - errorDetail: detailed error information
func ConflictError(c *gin.Context, message string, errorDetail string) {
	ErrorResponse(c, http.StatusConflict, message, errorDetail)
}
