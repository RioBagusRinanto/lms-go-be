package service

import (
	"fmt"
	"time"

	"lms-go-be/internal/models"
	"lms-go-be/internal/repository"
	"lms-go-be/internal/utils"
)

// QuizService handles quiz operations
type QuizService struct {
	quizRepo        *repository.QuizRepository
	quizAttemptRepo *repository.QuizAttemptRepository
	enrollmentRepo  *repository.EnrollmentRepository
	gamificationSvc *GamificationService
}

// NewQuizService creates a new quiz service
func NewQuizService(
	quizRepo *repository.QuizRepository,
	quizAttemptRepo *repository.QuizAttemptRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	gamificationSvc *GamificationService,
) *QuizService {
	return &QuizService{
		quizRepo:        quizRepo,
		quizAttemptRepo: quizAttemptRepo,
		enrollmentRepo:  enrollmentRepo,
		gamificationSvc: gamificationSvc,
	}
}

// StartQuizAttemptRequest represents start quiz attempt request
type StartQuizAttemptRequest struct {
	QuizID uint `json:"quiz_id" binding:"required"`
}

// SubmitQuizRequest represents submit quiz request
type SubmitQuizRequest struct {
	Answers map[uint]string `json:"answers"` // question_id: answer
}

// QuizAttemptDTO represents quiz attempt DTO
type QuizAttemptDTO struct {
	ID               uint       `json:"id"`
	QuizID           uint       `json:"quiz_id"`
	UserID           uint       `json:"user_id"`
	AttemptNumber    int        `json:"attempt_number"`
	Score            int        `json:"score"`
	MaxScore         int        `json:"max_score"`
	Percentage       int        `json:"percentage"`
	IsPassed         bool       `json:"is_passed"`
	TimeSpentSeconds int        `json:"time_spent_seconds"`
	StartedAt        time.Time  `json:"started_at"`
	SubmittedAt      *time.Time `json:"submitted_at"`
}

// GetQuiz gets a quiz by ID
func (s *QuizService) GetQuiz(quizID uint) (*models.Quiz, error) {
	return s.quizRepo.GetByID(quizID)
}

// GetQuizByLesson gets quiz for a lesson
func (s *QuizService) GetQuizByLesson(lessonID uint) (*models.Quiz, error) {
	return s.quizRepo.GetByLesson(lessonID)
}

// StartAttempt starts a new quiz attempt
func (s *QuizService) StartAttempt(userID, quizID uint) (*models.QuizAttempt, error) {
	// Check attempt limit
	quiz, err := s.quizRepo.GetByID(quizID)
	if err != nil {
		return nil, err
	}

	attemptCount, err := s.quizAttemptRepo.GetQuizAttemptCount(userID, quizID)
	if err != nil {
		return nil, err
	}

	if attemptCount >= int64(quiz.Attempts) {
		return nil, fmt.Errorf("maximum quiz attempts exceeded")
	}

	// Create new attempt
	attempt := &models.QuizAttempt{
		QuizID:        quizID,
		UserID:        userID,
		AttemptNumber: int(attemptCount) + 1,
		StartedAt:     time.Now(),
	}

	if err := s.quizAttemptRepo.Create(attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

// SubmitAttempt submits quiz answers
func (s *QuizService) SubmitAttempt(userID, quizID, quizAttemptID uint, answers map[uint]string, timeSpent int) (*models.QuizAttempt, error) {
	attempt, err := s.quizAttemptRepo.GetByID(quizAttemptID)
	if err != nil {
		return nil, err
	}

	// Get quiz
	quiz, err := s.quizRepo.GetByID(quizID)
	if err != nil {
		return nil, err
	}

	// Calculate score (simplified - in production, implement proper grading)
	score := 0
	maxScore := len(quiz.Questions)

	for _, question := range quiz.Questions {
		if userAnswer, exists := answers[question.ID]; exists {
			if isCorrectAnswer(&question, userAnswer) {
				score++
			}
		}
	}

	// Update attempt with results
	percentage := 0
	if maxScore > 0 {
		percentage = (score * 100) / maxScore
	}

	isPassed := percentage >= quiz.PassingScore

	attempt.SubmittedAt = utils.TimePtr(time.Now())
	attempt.Score = score
	attempt.MaxScore = maxScore
	attempt.Percentage = percentage
	attempt.IsPassed = isPassed
	attempt.TimeSpentSeconds = timeSpent

	if err := s.quizAttemptRepo.Update(attempt); err != nil {
		return nil, err
	}

	// Award coins if passed
	if isPassed {
		coinReward := int64(quiz.PassingScore * 2) // Simplified coin calculation
		_ = s.gamificationSvc.AwardCoins(userID, coinReward, fmt.Sprintf("Quiz Passed: %s", quiz.Title), "quiz", &quiz.ID)
	}

	return attempt, nil
}

// GetUserAttempts gets all attempts by user for a quiz
func (s *QuizService) GetUserAttempts(userID, quizID uint) ([]models.QuizAttempt, error) {
	return s.quizAttemptRepo.GetUserQuizAttempts(userID, quizID)
}

// GetAttemptCount gets attempt count for user
func (s *QuizService) GetAttemptCount(userID, quizID uint) (int64, error) {
	return s.quizAttemptRepo.GetQuizAttemptCount(userID, quizID)
}

// isCorrectAnswer checks if answer is correct (simplified)
func isCorrectAnswer(question *models.Question, answer string) bool {
	// This is simplified - in production, implement proper answer checking
	for _, option := range question.Options {
		if fmt.Sprintf("%d", option.ID) == answer && option.IsCorrect {
			return true
		}
	}
	return false
}

// DashboardService provides dashboard data
type DashboardService struct {
	enrollmentRepo      *repository.EnrollmentRepository
	progressRepo        *repository.UserProgressRepository
	certificateRepo     *repository.CertificateRepository
	coinTransactionRepo *repository.CoinTransactionRepository
	badgeProgressRepo   *repository.BadgeProgressRepository
	userRepo            *repository.UserRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	enrollmentRepo *repository.EnrollmentRepository,
	progressRepo *repository.UserProgressRepository,
	certificateRepo *repository.CertificateRepository,
	coinTransactionRepo *repository.CoinTransactionRepository,
	badgeProgressRepo *repository.BadgeProgressRepository,
	userRepo *repository.UserRepository,
) *DashboardService {
	return &DashboardService{
		enrollmentRepo:      enrollmentRepo,
		progressRepo:        progressRepo,
		certificateRepo:     certificateRepo,
		coinTransactionRepo: coinTransactionRepo,
		badgeProgressRepo:   badgeProgressRepo,
		userRepo:            userRepo,
	}
}

// DashboardData represents dashboard data
type DashboardData struct {
	MandatoryCourses   []models.Enrollment      `json:"mandatory_courses"`
	InProgressCourses  []models.Enrollment      `json:"in_progress_courses"`
	CompletedCourses   int64                    `json:"completed_courses"`
	Certificates       int64                    `json:"certificates"`
	GMFCCoins          int64                    `json:"gmfc_coins"`
	CurrentBadgeLevel  string                   `json:"current_badge_level"`
	EarnedBadges       int64                    `json:"earned_badges"`
	TotalLearningHours float64                  `json:"total_learning_hours"`
	CurrentStreak      int                      `json:"current_streak"`
	LeaderboardRank    int                      `json:"leaderboard_rank"`
	RecentTransactions []models.CoinTransaction `json:"recent_transactions"`
}

// GetUserDashboard gets complete dashboard data for a user
func (s *DashboardService) GetUserDashboard(userID uint) (*DashboardData, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get mandatory courses
	mandatoryCourses, _ := s.enrollmentRepo.GetUserMandatoryCourses(userID)

	// Get in-progress courses
	inProgressCourses, _ := s.enrollmentRepo.GetUserInProgressCourses(userID)

	// Get completed courses count
	completedEnrollments, totalCompleted, _ := s.enrollmentRepo.GetUserCompletedCourses(userID, 1, 1000)
	_ = completedEnrollments // Use if needed

	// Get certificates
	_, totalCertificates, _ := s.certificateRepo.GetUserCertificates(userID, 1, 1000)

	// Get earned badges count
	earnedBadgeCount, _ := s.badgeProgressRepo.GetUserEarnedBadgeCount(userID)

	// Get leaderboard rank
	leaderboardRank, _ := s.userRepo.GetUserRank(userID)

	// Get recent coin transactions
	recentTransactions, _ := s.coinTransactionRepo.GetUserRecentTransactions(userID, 5)

	dashboard := &DashboardData{
		MandatoryCourses:   mandatoryCourses,
		InProgressCourses:  inProgressCourses,
		CompletedCourses:   totalCompleted,
		Certificates:       totalCertificates,
		GMFCCoins:          user.GMFCCoins,
		CurrentBadgeLevel:  user.CurrentBadgeLevel,
		EarnedBadges:       earnedBadgeCount,
		TotalLearningHours: user.TotalLearningHours,
		CurrentStreak:      user.CurrentStreak,
		LeaderboardRank:    leaderboardRank,
		RecentTransactions: recentTransactions,
	}

	return dashboard, nil
}
