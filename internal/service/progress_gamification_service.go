package service

import (
	"encoding/json"
	"fmt"
	"time"

	"lms-go-be/internal/models"
	"lms-go-be/internal/repository"
	"lms-go-be/internal/utils"
)

// ProgressService handles user progress tracking
type ProgressService struct {
	userProgressRepo *repository.UserProgressRepository
	enrollmentRepo   *repository.EnrollmentRepository
	userRepo         *repository.UserRepository
}

// NewProgressService creates a new progress service
func NewProgressService(
	userProgressRepo *repository.UserProgressRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	userRepo *repository.UserRepository,
) *ProgressService {
	return &ProgressService{
		userProgressRepo: userProgressRepo,
		enrollmentRepo:   enrollmentRepo,
		userRepo:         userRepo,
	}
}

// UpdateVideoProgressRequest represents update video progress request
type UpdateVideoProgressRequest struct {
	WatchedDuration int `json:"watched_duration" binding:"min=0"`
	TotalDuration   int `json:"total_duration" binding:"min=0"`
}

// ProgressDTO represents progress data transfer object
type ProgressDTO struct {
	ID                 uint `json:"id"`
	UserID             uint `json:"user_id"`
	LessonID           uint `json:"lesson_id"`
	CourseID           uint `json:"course_id"`
	WatchedDuration    int  `json:"watched_duration_seconds"`
	TotalDuration      int  `json:"total_duration_seconds"`
	ProgressPercentage int  `json:"progress_percentage"`
	IsCompleted        bool `json:"is_completed"`
}

// TrackProgress tracks user's lesson progress
func (s *ProgressService) TrackProgress(userID, courseID, lessonID uint, watchedSeconds, totalSeconds int) (*models.UserProgress, error) {
	// Get or create progress record
	progress, err := s.userProgressRepo.GetByUserLessonCourse(userID, lessonID, courseID)
	if err != nil {
		// Create new progress record
		progress = &models.UserProgress{
			UserID:        userID,
			CourseID:      courseID,
			LessonID:      lessonID,
			TotalDuration: totalSeconds,
		}

		if err := s.userProgressRepo.Create(progress); err != nil {
			return nil, err
		}
	}

	// Update watched duration
	progress.WatchedDuration = watchedSeconds
	progress.TotalDuration = totalSeconds
	progress.ProgressPercentage = utils.CalculateProgressPercentage(watchedSeconds, totalSeconds)
	progress.LastAccessedAt = utils.TimePtr(time.Now())

	// Mark as completed if fully watched (90% or more)
	if progress.ProgressPercentage >= 90 {
		progress.IsCompleted = true
		progress.CompletedAt = utils.TimePtr(time.Now())
	}

	if err := s.userProgressRepo.Update(progress); err != nil {
		return nil, err
	}

	return progress, nil
}

// GetLessonProgress gets progress for a lesson
func (s *ProgressService) GetLessonProgress(userID, courseID, lessonID uint) (*models.UserProgress, error) {
	return s.userProgressRepo.GetByUserLessonCourse(userID, lessonID, courseID)
}

// GetCourseProgress gets all lesson progress for a course
func (s *ProgressService) GetCourseProgress(userID, courseID uint) ([]models.UserProgress, error) {
	return s.userProgressRepo.GetUserCourseProgress(userID, courseID)
}

// CalculateCourseProgress calculates overall course progress
func (s *ProgressService) CalculateCourseProgress(userID, courseID uint) (int, error) {
	progresses, err := s.userProgressRepo.GetUserCourseProgress(userID, courseID)
	if err != nil {
		return 0, err
	}

	if len(progresses) == 0 {
		return 0, nil
	}

	totalProgress := 0
	for _, p := range progresses {
		totalProgress += p.ProgressPercentage
	}

	averageProgress := totalProgress / len(progresses)
	return averageProgress, nil
}

// GetUserTotalLearningHours gets total learning hours for a user
func (s *ProgressService) GetUserTotalLearningHours(userID uint) (float64, error) {
	return s.userProgressRepo.GetUserTotalLearningHours(userID)
}

// GamificationService handles gamification operations
type GamificationService struct {
	coinTransactionRepo *repository.CoinTransactionRepository
	badgeRepo           *repository.BadgeRepository
	badgeProgressRepo   *repository.BadgeProgressRepository
	userRepo            *repository.UserRepository
	certificateRepo     *repository.CertificateRepository
}

// NewGamificationService creates a new gamification service
func NewGamificationService(
	coinTransactionRepo *repository.CoinTransactionRepository,
	badgeRepo *repository.BadgeRepository,
	badgeProgressRepo *repository.BadgeProgressRepository,
	userRepo *repository.UserRepository,
	certificateRepo *repository.CertificateRepository,
) *GamificationService {
	return &GamificationService{
		coinTransactionRepo: coinTransactionRepo,
		badgeRepo:           badgeRepo,
		badgeProgressRepo:   badgeProgressRepo,
		userRepo:            userRepo,
		certificateRepo:     certificateRepo,
	}
}

// CoinTransactionDTO represents coin transaction DTO
type CoinTransactionDTO struct {
	ID              uint      `json:"id"`
	Amount          int64     `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Reason          string    `json:"reason"`
	CreatedAt       time.Time `json:"created_at"`
}

// AwardCoins awards coins to a user
func (s *GamificationService) AwardCoins(userID uint, amount int64, reason, referenceType string, referenceID *uint) error {
	if amount <= 0 {
		return fmt.Errorf("coin amount must be positive")
	}

	// Update user coins
	if err := s.userRepo.UpdateCoins(userID, amount); err != nil {
		return err
	}

	// Log transaction
	transaction := &models.CoinTransaction{
		UserID:          userID,
		Amount:          amount,
		TransactionType: "earned",
		Reason:          reason,
		ReferenceType:   referenceType,
		ReferenceID:     referenceID,
	}

	return s.coinTransactionRepo.Create(transaction)
}

// SpendCoins spends coins from a user
func (s *GamificationService) SpendCoins(userID uint, amount int64, reason string) error {
	if amount <= 0 {
		return fmt.Errorf("coin amount must be positive")
	}

	// Check if user has enough coins
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user.GMFCCoins < amount {
		return fmt.Errorf("insufficient coins")
	}

	// Deduct coins
	if err := s.userRepo.UpdateCoins(userID, -amount); err != nil {
		return err
	}

	// Log transaction
	transaction := &models.CoinTransaction{
		UserID:          userID,
		Amount:          -amount,
		TransactionType: "spent",
		Reason:          reason,
	}

	return s.coinTransactionRepo.Create(transaction)
}

// GetUserCoins gets user's coin balance
func (s *GamificationService) GetUserCoins(userID uint) (int64, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, err
	}
	return user.GMFCCoins, nil
}

// GetCoinTransactions gets coin transactions for a user
func (s *GamificationService) GetCoinTransactions(userID uint, page, pageSize int) ([]models.CoinTransaction, int64, error) {
	return s.coinTransactionRepo.GetUserTransactions(userID, page, pageSize)
}

// CheckAndAwardBadges checks and awards badges to a user
func (s *GamificationService) CheckAndAwardBadges(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Get all badges
	badges, err := s.badgeRepo.GetAll()
	if err != nil {
		return err
	}

	for _, badge := range badges {
		// Check if badge is already earned
		progress, _ := s.badgeProgressRepo.GetUserBadgeProgress(userID, badge.ID)
		if progress != nil && progress.IsEarned {
			continue
		}

		// Parse badge criteria
		var criteria map[string]interface{}
		_ = json.Unmarshal([]byte(badge.Criteria), &criteria)

		// Check criteria (simplified logic)
		if shouldEarnBadge(user, criteria) {
			// Create or update badge progress
			if progress == nil {
				progress = &models.BadgeProgress{
					UserID:  userID,
					BadgeID: badge.ID,
				}
				_ = s.badgeProgressRepo.Create(progress)
			}

			// Mark badge as earned
			_ = s.badgeProgressRepo.MarkBadgeEarned(userID, badge.ID)

			// Update user's badge level
			_ = s.userRepo.UpdateBadgeLevel(userID, badge.Level)
		}
	}

	return nil
}

// shouldEarnBadge checks if user meets badge criteria
func shouldEarnBadge(user *models.User, criteria map[string]interface{}) bool {
	// This is simplified logic - in production, implement proper criteria checking
	badgeType, ok := criteria["type"].(string)
	if !ok {
		return false
	}

	switch badgeType {
	case "coins_earned":
		value := int64(criteria["value"].(float64))
		return user.GMFCCoins >= value
	case "learning_hours":
		value := criteria["hours"].(float64)
		return user.TotalLearningHours >= value
	default:
		return false
	}
}

// GetUserBadges gets all badges for a user
func (s *GamificationService) GetUserBadges(userID uint) ([]models.BadgeProgress, error) {
	return s.badgeProgressRepo.GetUserBadges(userID)
}

// GetUserEarnedBadges gets earned badges for a user
func (s *GamificationService) GetUserEarnedBadges(userID uint) ([]models.BadgeProgress, error) {
	return s.badgeProgressRepo.GetUserEarnedBadges(userID)
}
