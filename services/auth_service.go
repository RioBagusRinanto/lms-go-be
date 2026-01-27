package services

import (
	"errors"
	"lms-go-be/models"
	"lms-go-be/repositories"
	"lms-go-be/utils"
	"time"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService creates a new auth service instance
// Parameters:
//   - userRepo: user repository instance
//
// Returns: AuthService instance
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register registers a new user
// Parameters:
//   - firstName: user's first name
//   - lastName: user's last name
//   - email: user's email
//   - password: plain text password
//   - role: user role (learner, instructor, admin)
//
// Returns: created user and error if any
func (s *AuthService) Register(firstName, lastName, email, password, role string) (*models.User, error) {
	// Validate email format
	if !utils.ValidateEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      hashedPassword,
		Role:          role,
		Status:        "active",
		GMFCCoins:     0,
		BadgeLevel:    "Bronze",
		CurrentStreak: 0,
		IsVerified:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// Login authenticates a user and returns JWT token
// Parameters:
//   - email: user's email
//   - password: plain text password
//
// Returns: JWT token string and error if any
func (s *AuthService) Login(email, password string) (string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Verify password
	if !utils.ComparePassword(user.Password, password) {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserProfile retrieves user profile by ID
// Parameters:
//   - userID: user ID
//
// Returns: user model and error if any
func (s *AuthService) GetUserProfile(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// UpdateUserProfile updates user profile information
// Parameters:
//   - user: user model with updated data
//
// Returns: updated user and error if any
func (s *AuthService) UpdateUserProfile(user *models.User) (*models.User, error) {
	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

// ChangePassword changes user password
// Parameters:
//   - userID: user ID
//   - oldPassword: current password
//   - newPassword: new password
//
// Returns: error if any
func (s *AuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !utils.ComparePassword(user.Password, oldPassword) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	_, err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken generates a new JWT token for an existing user
// Parameters:
//   - userID: user ID
//
// Returns: new JWT token and error if any
func (s *AuthService) RefreshToken(userID string) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate new token
	token, err := utils.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserStats retrieves user statistics
// Parameters:
//   - userID: user ID
//
// Returns: map with user stats and error if any
func (s *AuthService) GetUserStats(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	stats := map[string]interface{}{
		"id":             user.ID,
		"email":          user.Email,
		"gmfc_coins":     user.GMFCCoins,
		"badge_level":    user.BadgeLevel,
		"current_streak": user.CurrentStreak,
	}

	return stats, nil
}
