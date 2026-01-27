package handlers

import (
	"net/http"

	"lms-go-be/services"
	"lms-go-be/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest defines register request structure
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=learner instructor admin"`
}

// LoginRequest defines login request structure
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "invalid request", err.Error())
		return
	}

	// Register user
	user, err := h.authService.Register(req.FirstName, req.LastName, req.Email, req.Password, req.Role)
	if err != nil {
		utils.BadRequestError(c, "registration failed", err.Error())
		return
	}

	// Return response without password
	utils.SuccessResponse(c, http.StatusCreated, "user registered successfully", gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
	})
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "invalid request", err.Error())
		return
	}

	// Login user
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.UnauthorizedError(c, err.Error())
		return
	}

	// Return token
	utils.SuccessResponse(c, http.StatusOK, "login successful", gin.H{
		"token": token,
		"type":  "Bearer",
	})
}

// GetProfile handles getting user profile
// GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not found in context")
		return
	}

	// Get user profile
	user, err := h.authService.GetUserProfile(userID.(string))
	if err != nil {
		utils.NotFoundError(c, err.Error())
		return
	}

	// Return profile without password
	utils.SuccessResponse(c, http.StatusOK, "profile retrieved successfully", gin.H{
		"id":             user.ID,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"email":          user.Email,
		"profile_image":  user.ProfileImage,
		"bio":            user.Bio,
		"role":           user.Role,
		"gmfc_coins":     user.GMFCCoins,
		"badge_level":    user.BadgeLevel,
		"current_streak": user.CurrentStreak,
	})
}

// UpdateProfileRequest defines update profile request structure
type UpdateProfileRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfileImage string `json:"profile_image"`
	Bio          string `json:"bio"`
	Department   string `json:"department"`
}

// UpdateProfile handles updating user profile
// PUT /api/auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not found in context")
		return
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "invalid request", err.Error())
		return
	}

	// Get existing user
	user, err := h.authService.GetUserProfile(userID.(string))
	if err != nil {
		utils.NotFoundError(c, err.Error())
		return
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Department != "" {
		user.Department = req.Department
	}

	// Update user
	updatedUser, err := h.authService.UpdateUserProfile(user)
	if err != nil {
		utils.InternalServerError(c, "failed to update profile", err.Error())
		return
	}

	// Return updated profile
	utils.SuccessResponse(c, http.StatusOK, "profile updated successfully", gin.H{
		"id":            updatedUser.ID,
		"first_name":    updatedUser.FirstName,
		"last_name":     updatedUser.LastName,
		"email":         updatedUser.Email,
		"profile_image": updatedUser.ProfileImage,
		"bio":           updatedUser.Bio,
		"role":          updatedUser.Role,
	})
}

// ChangePasswordRequest defines change password request structure
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword handles changing user password
// POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not found in context")
		return
	}

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "invalid request", err.Error())
		return
	}

	// Change password
	err := h.authService.ChangePassword(userID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		utils.BadRequestError(c, "failed to change password", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "password changed successfully", nil)
}

// RefreshTokenRequest defines refresh token request structure
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshToken handles refreshing JWT token
// POST /api/auth/refresh-token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not found in context")
		return
	}

	// Generate new token
	newToken, err := h.authService.RefreshToken(userID.(string))
	if err != nil {
		utils.BadRequestError(c, "failed to refresh token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "token refreshed successfully", gin.H{
		"token": newToken,
		"type":  "Bearer",
	})
}

// GetStats handles getting user statistics
// GET /api/auth/stats
func (h *AuthHandler) GetStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not found in context")
		return
	}

	// Get stats
	stats, err := h.authService.GetUserStats(userID.(string))
	if err != nil {
		utils.NotFoundError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "stats retrieved successfully", stats)
}
