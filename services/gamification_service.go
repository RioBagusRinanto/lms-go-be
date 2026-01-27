package services

import (
	"errors"
	"lms-go-be/models"
	"lms-go-be/repositories"
	"lms-go-be/utils"
	"time"
)

// GamificationService handles gamification business logic (coins, badges, streaks)
type GamificationService struct {
	userRepo       *repositories.UserRepository
	coinRepo       *repositories.CoinHistoryRepository
	badgeRepo      *repositories.BadgeHistoryRepository
	enrollmentRepo *repositories.EnrollmentRepository
}

// NewGamificationService creates a new gamification service instance
func NewGamificationService(
	userRepo *repositories.UserRepository,
	coinRepo *repositories.CoinHistoryRepository,
	badgeRepo *repositories.BadgeHistoryRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
) *GamificationService {
	return &GamificationService{
		userRepo:       userRepo,
		coinRepo:       coinRepo,
		badgeRepo:      badgeRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// AwardCoinsForCourseCompletion awards coins when user completes a course
// Parameters:
//   - userID: user ID
//   - score: course completion score
//   - duration: course duration in minutes
//
// Returns: error if any
func (s *GamificationService) AwardCoinsForCourseCompletion(userID string, score int, duration int) error {
	// Calculate coins
	coins := utils.CalculateCoinsEarned(score, duration)

	// Update user coins
	if err := s.userRepo.UpdateCoins(userID, coins); err != nil {
		return err
	}

	// Record coin history
	history := &models.CoinHistory{
		UserID:      userID,
		Amount:      coins,
		Reason:      "course_completion",
		ReferenceID: "",
		CreatedAt:   time.Now(),
	}

	_, err := s.coinRepo.Create(history)
	return err
}

// AwardCoinsForQuizPass awards coins when user passes a quiz
// Parameters:
//   - userID: user ID
//   - score: quiz score
//
// Returns: error if any
func (s *GamificationService) AwardCoinsForQuizPass(userID string, score int) error {
	coins := int64(0)

	if score >= 90 {
		coins = 50
	} else if score >= 80 {
		coins = 25
	} else if score >= 70 {
		coins = 10
	}

	if coins == 0 {
		return nil // No coins for low scores
	}

	// Update user coins
	if err := s.userRepo.UpdateCoins(userID, coins); err != nil {
		return err
	}

	// Record coin history
	history := &models.CoinHistory{
		UserID:    userID,
		Amount:    coins,
		Reason:    "quiz_passed",
		CreatedAt: time.Now(),
	}

	_, err := s.coinRepo.Create(history)
	return err
}

// AwardStreakBonus awards coins for maintaining a streak
// Parameters:
//   - userID: user ID
//   - streak: current streak
//
// Returns: error if any
func (s *GamificationService) AwardStreakBonus(userID string, streak int) error {
	// Bonus coins based on streak
	coins := int64(streak * 10)

	// Update user coins
	if err := s.userRepo.UpdateCoins(userID, coins); err != nil {
		return err
	}

	// Record coin history
	history := &models.CoinHistory{
		UserID:    userID,
		Amount:    coins,
		Reason:    "streak_bonus",
		CreatedAt: time.Now(),
	}

	_, err := s.coinRepo.Create(history)
	return err
}

// UpdateBadgeLevel updates user's badge level based on completed courses
// Parameters:
//   - userID: user ID
//
// Returns: error if any
func (s *GamificationService) UpdateBadgeLevel(userID string) error {
	// Count completed courses
	completedCount, err := s.enrollmentRepo.GetCompletionCount(userID)
	if err != nil {
		return err
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Calculate new badge level
	newBadgeLevel := utils.CalculateBadgeLevel(int(completedCount))

	// Only update if badge level changed
	if user.BadgeLevel != newBadgeLevel {
		// Update user badge level
		if err := s.userRepo.UpdateBadgeLevel(userID, newBadgeLevel); err != nil {
			return err
		}

		// Record badge achievement
		badgeHistory := &models.BadgeHistory{
			UserID:     userID,
			BadgeLevel: newBadgeLevel,
			AchievedAt: time.Now(),
			CreatedAt:  time.Now(),
		}

		_, err := s.badgeRepo.Create(badgeHistory)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateStreak updates user's learning streak
// Parameters:
//   - userID: user ID
//   - increment: whether to increment (true) or reset (false) streak
//
// Returns: error if any
func (s *GamificationService) UpdateStreak(userID string, increment bool) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	newStreak := user.CurrentStreak
	if increment {
		newStreak++
	} else {
		newStreak = 0
	}

	return s.userRepo.UpdateStreak(userID, newStreak)
}

// GetUserCoins retrieves user's coin balance and recent transactions
// Parameters:
//   - userID: user ID
//
// Returns: map with coin info and error if any
func (s *GamificationService) GetUserCoins(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get recent transactions
	history, err := s.coinRepo.GetUserHistory(userID, 5)
	if err != nil {
		return nil, err
	}

	coinInfo := map[string]interface{}{
		"balance":             user.GMFCCoins,
		"recent_transactions": history,
	}

	return coinInfo, nil
}

// GetUserBadges retrieves user's badge information
// Parameters:
//   - userID: user ID
//
// Returns: map with badge info and error if any
func (s *GamificationService) GetUserBadges(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get badge history
	badges, err := s.badgeRepo.GetUserBadges(userID)
	if err != nil {
		return nil, err
	}

	// Get completion count for progress to next badge
	completedCount, err := s.enrollmentRepo.GetCompletionCount(userID)
	if err != nil {
		return nil, err
	}

	// Calculate progress to next badge
	currentBadgeReq := utils.GetBadgeRequirements(user.BadgeLevel)
	var nextBadgeLevel string
	var nextBadgeReq int

	switch user.BadgeLevel {
	case "Bronze":
		nextBadgeLevel = "Silver"
		nextBadgeReq = utils.GetBadgeRequirements("Silver")
	case "Silver":
		nextBadgeLevel = "Gold"
		nextBadgeReq = utils.GetBadgeRequirements("Gold")
	case "Gold":
		nextBadgeLevel = "Platinum"
		nextBadgeReq = utils.GetBadgeRequirements("Platinum")
	case "Platinum":
		nextBadgeLevel = "Diamond"
		nextBadgeReq = utils.GetBadgeRequirements("Diamond")
	default:
		nextBadgeLevel = "Diamond"
		nextBadgeReq = utils.GetBadgeRequirements("Diamond")
	}

	progressToNext := int(completedCount) - currentBadgeReq
	if progressToNext < 0 {
		progressToNext = 0
	}

	badgeInfo := map[string]interface{}{
		"current_badge":          user.BadgeLevel,
		"next_badge":             nextBadgeLevel,
		"completed_courses":      completedCount,
		"progress_to_next_badge": progressToNext,
		"courses_needed_next":    nextBadgeReq - int(completedCount),
		"badge_history":          badges,
	}

	return badgeInfo, nil
}

// GetLeaderboardPosition retrieves user's position on the leaderboard
// Parameters:
//   - userID: user ID
//
// Returns: position information and error if any
func (s *GamificationService) GetLeaderboardPosition(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// This would typically involve a database query to rank users by coins
	// For now, we'll return basic info
	leaderboardInfo := map[string]interface{}{
		"user_id":     user.ID,
		"gmfc_coins":  user.GMFCCoins,
		"badge_level": user.BadgeLevel,
	}

	return leaderboardInfo, nil
}

// RedeemCoins allows user to redeem coins
// Parameters:
//   - userID: user ID
//   - coinsToRedeem: number of coins to redeem
//   - reason: reason for redemption
//
// Returns: error if any
func (s *GamificationService) RedeemCoins(userID string, coinsToRedeem int64, reason string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user.GMFCCoins < coinsToRedeem {
		return errors.New("insufficient coins for redemption")
	}

	// Deduct coins
	if err := s.userRepo.UpdateCoins(userID, -coinsToRedeem); err != nil {
		return err
	}

	// Record coin history
	history := &models.CoinHistory{
		UserID:    userID,
		Amount:    -coinsToRedeem,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	_, err = s.coinRepo.Create(history)
	return err
}
