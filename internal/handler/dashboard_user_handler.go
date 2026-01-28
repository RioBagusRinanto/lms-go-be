package handler

import (
	"net/http"
	"strconv"

	"lms-go-be/internal/repository"
	"lms-go-be/internal/service"
	"lms-go-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard gets user dashboard data
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	dashboard, err := h.dashboardService.GetUserDashboard(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve dashboard", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard retrieved successfully", dashboard)
}

// UserHandler handles user-related endpoints
type UserHandler struct {
	userRepo            *repository.UserRepository
	gamificationService *service.GamificationService
	badgeProgressRepo   *repository.BadgeProgressRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userRepo *repository.UserRepository,
	gamificationService *service.GamificationService,
	badgeProgressRepo *repository.BadgeProgressRepository,
) *UserHandler {
	return &UserHandler{
		userRepo:            userRepo,
		gamificationService: gamificationService,
		badgeProgressRepo:   badgeProgressRepo,
	}
}

// GetUserProfile gets a user's profile
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userRepo.GetByID(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	userDTO := service.ConvertUserToDTO(user)
	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", userDTO)
}

// GetLeaderboard gets the leaderboard
func (h *UserHandler) GetLeaderboard(c *gin.Context) {
	orderBy := c.Query("order_by")
	if orderBy == "" {
		orderBy = "hours"
	}

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	users, err := h.userRepo.GetLeaderboard(orderBy, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve leaderboard", err.Error())
		return
	}

	dtos := make([]interface{}, len(users))
	for i, user := range users {
		dtos[i] = service.ConvertUserToDTO(&user)
	}

	utils.SuccessResponse(c, http.StatusOK, "Leaderboard retrieved successfully", dtos)
}

// GetCoins gets user's coin balance
func (h *UserHandler) GetCoins(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	coins, err := h.gamificationService.GetUserCoins(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve coins", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Coins retrieved successfully", map[string]interface{}{
		"balance": coins,
	})
}

// GetCoinTransactions gets coin transactions for a user
func (h *UserHandler) GetCoinTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	transactions, total, err := h.gamificationService.GetCoinTransactions(userID.(uint), page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions, page, 10, total)
}

// GetBadges gets all badges for a user
func (h *UserHandler) GetBadges(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	badges, err := h.gamificationService.GetUserBadges(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve badges", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Badges retrieved successfully", badges)
}

// GetEarnedBadges gets earned badges for a user
func (h *UserHandler) GetEarnedBadges(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	badges, err := h.gamificationService.GetUserEarnedBadges(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve badges", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Earned badges retrieved successfully", badges)
}

// ListUsers gets all users (admin only)
func (h *UserHandler) ListUsers(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	users, total, err := h.userRepo.GetAll(page, 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", err.Error())
		return
	}

	dtos := make([]interface{}, len(users))
	for i, user := range users {
		dtos[i] = service.ConvertUserToDTO(&user)
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Users retrieved successfully", dtos, page, 10, total)
}

// AdjustCoins adjusts user coins (admin only)
func (h *UserHandler) AdjustCoins(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req struct {
		Amount int64  `json:"amount" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.userRepo.UpdateCoins(uint(userID), req.Amount); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to adjust coins", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Coins adjusted successfully", nil)
}
