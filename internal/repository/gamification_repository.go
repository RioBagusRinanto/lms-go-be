package repository

import (
	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// CoinTransactionRepository handles coin transaction database operations
type CoinTransactionRepository struct {
	db *gorm.DB
}

// NewCoinTransactionRepository creates a new coin transaction repository
func NewCoinTransactionRepository(db *gorm.DB) *CoinTransactionRepository {
	return &CoinTransactionRepository{db: db}
}

// Create creates a new coin transaction
func (r *CoinTransactionRepository) Create(transaction *models.CoinTransaction) error {
	return r.db.Create(transaction).Error
}

// GetByID gets a coin transaction by ID
func (r *CoinTransactionRepository) GetByID(id uint) (*models.CoinTransaction, error) {
	var transaction models.CoinTransaction
	if err := r.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetUserTransactions gets all transactions for a user
func (r *CoinTransactionRepository) GetUserTransactions(userID uint, page, pageSize int) ([]models.CoinTransaction, int64, error) {
	var transactions []models.CoinTransaction
	var total int64

	if err := r.db.Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetUserRecentTransactions gets recent transactions for a user
func (r *CoinTransactionRepository) GetUserRecentTransactions(userID uint, limit int) ([]models.CoinTransaction, error) {
	var transactions []models.CoinTransaction
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetUserTotalEarned gets total coins earned by a user
func (r *CoinTransactionRepository) GetUserTotalEarned(userID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&models.CoinTransaction{}).
		Where("user_id = ? AND transaction_type IN ?", userID, []string{"earned", "admin_adjustment"}).
		Where("amount > 0").
		Pluck("COALESCE(SUM(amount), 0)", &total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetUserTotalSpent gets total coins spent by a user
func (r *CoinTransactionRepository) GetUserTotalSpent(userID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&models.CoinTransaction{}).
		Where("user_id = ? AND transaction_type IN ?", userID, []string{"spent", "redeemed"}).
		Where("amount < 0").
		Pluck("COALESCE(SUM(amount), 0)", &total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// BadgeRepository handles badge database operations
type BadgeRepository struct {
	db *gorm.DB
}

// NewBadgeRepository creates a new badge repository
func NewBadgeRepository(db *gorm.DB) *BadgeRepository {
	return &BadgeRepository{db: db}
}

// Create creates a new badge
func (r *BadgeRepository) Create(badge *models.Badge) error {
	return r.db.Create(badge).Error
}

// GetByID gets a badge by ID
func (r *BadgeRepository) GetByID(id uint) (*models.Badge, error) {
	var badge models.Badge
	if err := r.db.First(&badge, id).Error; err != nil {
		return nil, err
	}
	return &badge, nil
}

// GetByName gets a badge by name
func (r *BadgeRepository) GetByName(name string) (*models.Badge, error) {
	var badge models.Badge
	if err := r.db.Where("name = ?", name).First(&badge).Error; err != nil {
		return nil, err
	}
	return &badge, nil
}

// GetAll gets all badges
func (r *BadgeRepository) GetAll() ([]models.Badge, error) {
	var badges []models.Badge
	if err := r.db.Find(&badges).Error; err != nil {
		return nil, err
	}
	return badges, nil
}

// GetByLevel gets badges by level
func (r *BadgeRepository) GetByLevel(level string) ([]models.Badge, error) {
	var badges []models.Badge
	if err := r.db.Where("level = ?", level).Find(&badges).Error; err != nil {
		return nil, err
	}
	return badges, nil
}

// BadgeProgressRepository handles badge progress database operations
type BadgeProgressRepository struct {
	db *gorm.DB
}

// NewBadgeProgressRepository creates a new badge progress repository
func NewBadgeProgressRepository(db *gorm.DB) *BadgeProgressRepository {
	return &BadgeProgressRepository{db: db}
}

// Create creates a new badge progress record
func (r *BadgeProgressRepository) Create(progress *models.BadgeProgress) error {
	return r.db.Create(progress).Error
}

// GetByID gets a badge progress record by ID
func (r *BadgeProgressRepository) GetByID(id uint) (*models.BadgeProgress, error) {
	var progress models.BadgeProgress
	if err := r.db.Preload("Badge").First(&progress, id).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetUserBadgeProgress gets user's progress for a specific badge
func (r *BadgeProgressRepository) GetUserBadgeProgress(userID, badgeID uint) (*models.BadgeProgress, error) {
	var progress models.BadgeProgress
	if err := r.db.Where("user_id = ? AND badge_id = ?", userID, badgeID).
		Preload("Badge").First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// Update updates a badge progress record
func (r *BadgeProgressRepository) Update(progress *models.BadgeProgress) error {
	return r.db.Save(progress).Error
}

// GetUserBadges gets all badges for a user
func (r *BadgeProgressRepository) GetUserBadges(userID uint) ([]models.BadgeProgress, error) {
	var progresses []models.BadgeProgress
	if err := r.db.Where("user_id = ?", userID).Preload("Badge").
		Order("earned_at DESC NULLS LAST").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}

// GetUserEarnedBadges gets earned badges for a user
func (r *BadgeProgressRepository) GetUserEarnedBadges(userID uint) ([]models.BadgeProgress, error) {
	var progresses []models.BadgeProgress
	if err := r.db.Where("user_id = ? AND is_earned = ?", userID, true).
		Preload("Badge").Order("earned_at DESC").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}

// GetUserEarnedBadgeCount gets count of earned badges for a user
func (r *BadgeProgressRepository) GetUserEarnedBadgeCount(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.BadgeProgress{}).
		Where("user_id = ? AND is_earned = ?", userID, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateProgress updates badge progress percentage
func (r *BadgeProgressRepository) UpdateProgress(userID, badgeID uint, progress int) error {
	return r.db.Model(&models.BadgeProgress{}).
		Where("user_id = ? AND badge_id = ?", userID, badgeID).
		Update("progress", progress).Error
}

// MarkBadgeEarned marks a badge as earned
func (r *BadgeProgressRepository) MarkBadgeEarned(userID, badgeID uint) error {
	return r.db.Model(&models.BadgeProgress{}).
		Where("user_id = ? AND badge_id = ?", userID, badgeID).
		Updates(map[string]interface{}{
			"is_earned": true,
			"earned_at": gorm.Expr("NOW()"),
			"progress":  100,
		}).Error
}
