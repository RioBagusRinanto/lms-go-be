package service

import (
	"fmt"
	"time"

	"lms-go-be/internal/models"
	"lms-go-be/internal/repository"
	"lms-go-be/internal/utils"

	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	Token   string   `json:"token"`
	User    *UserDTO `json:"user"`
	Message string   `json:"message"`
}

// UserDTO represents user data transfer object
type UserDTO struct {
	ID                 uint      `json:"id"`
	Email              string    `json:"email"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	FullName           string    `json:"full_name"`
	Department         string    `json:"department"`
	Role               string    `json:"role"`
	GMFCCoins          int64     `json:"gmfc_coins"`
	CurrentBadgeLevel  string    `json:"current_badge_level"`
	TotalLearningHours float64   `json:"total_learning_hours"`
	CurrentStreak      int       `json:"current_streak"`
	CreatedAt          time.Time `json:"created_at"`
}

// Register registers a new user
func (s *AuthService) Register(email, password, firstName, lastName, department string) (*models.User, error) {
	// Validate email
	if !utils.ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate password
	if !utils.ValidatePassword(password) {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Email:      email,
		Password:   hashedPassword,
		FirstName:  firstName,
		LastName:   lastName,
		Department: department,
		Role:       "learner", // Default role
		IsActive:   true,
		GMFCCoins:  0, // Start with 0 coins
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(email, password string) (*models.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Compare password
	if !utils.ComparePassword(user.Password, password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return user, nil
}

// GetUser gets a user by ID
func (s *AuthService) GetUser(userID uint) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(userID uint, firstName, lastName, department string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Department = department

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Validate new password
	if !utils.ValidatePassword(newPassword) {
		return fmt.Errorf("password must be at least 6 characters")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !utils.ComparePassword(user.Password, oldPassword) {
		return fmt.Errorf("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// ConvertUserToDTO converts user model to DTO
func ConvertUserToDTO(user *models.User) *UserDTO {
	return &UserDTO{
		ID:                 user.ID,
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		FullName:           fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Department:         user.Department,
		Role:               user.Role,
		GMFCCoins:          user.GMFCCoins,
		CurrentBadgeLevel:  user.CurrentBadgeLevel,
		TotalLearningHours: user.TotalLearningHours,
		CurrentStreak:      user.CurrentStreak,
		CreatedAt:          user.CreatedAt,
	}
}
