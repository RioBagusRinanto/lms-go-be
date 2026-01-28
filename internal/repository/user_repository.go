package repository

import (
	"fmt"

	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user (soft delete)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetAll gets all users with pagination
func (r *UserRepository) GetAll(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByRole gets users by role
func (r *UserRepository) GetByRole(role string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("role = ?", role).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByDepartment gets users by department
func (r *UserRepository) GetByDepartment(department string) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("department = ?", department).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateCoins updates user's GMFC coins
func (r *UserRepository) UpdateCoins(userID uint, amount int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("gmfc_coins", gorm.Expr("gmfc_coins + ?", amount)).Error
}

// UpdateBadgeLevel updates user's badge level
func (r *UserRepository) UpdateBadgeLevel(userID uint, level string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("current_badge_level", level).Error
}

// IncrementStreak increments user's learning streak
func (r *UserRepository) IncrementStreak(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("current_streak", gorm.Expr("current_streak + ?", 1)).Error
}

// ResetStreak resets user's learning streak
func (r *UserRepository) ResetStreak(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("current_streak", 0).Error
}

// GetLeaderboard gets top users by learning hours or coins
func (r *UserRepository) GetLeaderboard(orderBy string, limit int) ([]models.User, error) {
	var users []models.User
	orderField := "total_learning_hours"
	if orderBy == "coins" {
		orderField = "gmfc_coins"
	}

	if err := r.db.Order(orderField + " DESC").Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserRank gets user's rank on leaderboard
func (r *UserRepository) GetUserRank(userID uint) (int, error) {
	var rank int
	err := r.db.Raw(`
		SELECT COALESCE(COUNT(*), 0) + 1 FROM users 
		WHERE total_learning_hours > (SELECT total_learning_hours FROM users WHERE id = ?)
		AND deleted_at IS NULL
	`, userID).Scan(&rank).Error
	return rank, err
}

// SearchUsers searches users by name or email
func (r *UserRepository) SearchUsers(query string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	searchQuery := fmt.Sprintf("%%%s%%", query)

	if err := r.db.Model(&models.User{}).Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
		searchQuery, searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
		searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin updates user's last login time
func (r *UserRepository) UpdateLastLogin(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}
