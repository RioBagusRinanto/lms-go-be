package repositories

import (
	"lms-go-be/models"

	"gorm.io/gorm"
)

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
// Parameters:
//   - user: user model to create
//
// Returns: created user and error if any
func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
// Parameters:
//   - id: user ID
//
// Returns: user model and error if not found
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
// Parameters:
//   - email: user email
//
// Returns: user model and error if not found
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
// Parameters:
//   - user: user model with updated data
//
// Returns: updated user and error if any
func (r *UserRepository) Update(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll retrieves all users with pagination
// Parameters:
//   - page: page number
//   - pageSize: number of records per page
//   - role: optional role filter
//
// Returns: slice of users, total count, and error if any
func (r *UserRepository) GetAll(page, pageSize int, role string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Delete deletes a user by ID
// Parameters:
//   - id: user ID
//
// Returns: error if any
func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// GetActiveUsers retrieves active users count
// Returns: count of active users and error if any
func (r *UserRepository) GetActiveUsers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("status = ?", "active").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateCoins updates user's GMFC coins
// Parameters:
//   - userID: user ID
//   - coins: amount to add (can be negative)
//
// Returns: error if any
func (r *UserRepository) UpdateCoins(userID string, coins int64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("gmfc_coins", gorm.Expr("gmfc_coins + ?", coins)).Error
}

// UpdateBadgeLevel updates user's badge level
// Parameters:
//   - userID: user ID
//   - badgeLevel: new badge level
//
// Returns: error if any
func (r *UserRepository) UpdateBadgeLevel(userID, badgeLevel string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("badge_level", badgeLevel).Error
}

// UpdateStreak updates user's current streak
// Parameters:
//   - userID: user ID
//   - streak: new streak value
//
// Returns: error if any
func (r *UserRepository) UpdateStreak(userID string, streak int) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("current_streak", streak).Error
}
