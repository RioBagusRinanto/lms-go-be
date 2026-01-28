package handler

import (
	"net/http"

	"lms-go-be/internal/config"
	"lms-go-be/internal/service"
	"lms-go-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
	config      *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Department string `json:"department"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Register user
	user, err := h.authService.Register(req.Email, req.Password, req.FirstName, req.LastName, req.Department)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, "Registration failed", err.Error())
		return
	}

	userDTO := service.ConvertUserToDTO(user)
	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", userDTO)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Authenticate user
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, user.FirstName+" "+user.LastName, h.config.JWT.SecretKey, h.config.JWT.ExpiresIn)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Token generation failed", err.Error())
		return
	}

	userDTO := service.ConvertUserToDTO(user)
	response := service.LoginResponse{
		Token:   token,
		User:    userDTO,
		Message: "Login successful",
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// GetProfile gets current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	user, err := h.authService.GetUser(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	userDTO := service.ConvertUserToDTO(user)
	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", userDTO)
}

// UpdateProfileRequest represents update profile request
type UpdateProfileRequest struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Department string `json:"department"`
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID.(uint), req.FirstName, req.LastName, req.Department)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	userDTO := service.ConvertUserToDTO(user)
	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", userDTO)
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword changes user password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real application, you might invalidate the token here
	// For now, this is a simple endpoint that confirms logout
	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}

// HealthCheck endpoint
func HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Server is running", map[string]string{
		"status": "healthy",
	})
}
