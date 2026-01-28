package service

import (
	"fmt"
	"time"

	"lms-go-be/internal/models"
	"lms-go-be/internal/repository"
)

// EnrollmentService handles enrollment business logic
type EnrollmentService struct {
	enrollmentRepo      *repository.EnrollmentRepository
	courseRepo          *repository.CourseRepository
	userProgressRepo    *repository.UserProgressRepository
	userRepo            *repository.UserRepository
	coinTransactionRepo *repository.CoinTransactionRepository
	certificateRepo     *repository.CertificateRepository
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(
	enrollmentRepo *repository.EnrollmentRepository,
	courseRepo *repository.CourseRepository,
	userProgressRepo *repository.UserProgressRepository,
	userRepo *repository.UserRepository,
	coinTransactionRepo *repository.CoinTransactionRepository,
	certificateRepo *repository.CertificateRepository,
) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:      enrollmentRepo,
		courseRepo:          courseRepo,
		userProgressRepo:    userProgressRepo,
		userRepo:            userRepo,
		coinTransactionRepo: coinTransactionRepo,
		certificateRepo:     certificateRepo,
	}
}

// EnrollRequest represents enrollment request
type EnrollRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}

// EnrollmentDTO represents enrollment data transfer object
type EnrollmentDTO struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	CourseID         uint       `json:"course_id"`
	CompletionStatus string     `json:"completion_status"`
	OverallProgress  int        `json:"overall_progress"`
	FinalScore       int        `json:"final_score"`
	IsPassed         bool       `json:"is_passed"`
	EnrolledAt       time.Time  `json:"enrolled_at"`
	CompletedAt      *time.Time `json:"completed_at"`
}

// EnrollUser enrolls a user in a course
func (s *EnrollmentService) EnrollUser(userID, courseID uint) (*models.Enrollment, error) {
	// Check if user already enrolled
	isEnrolled, err := s.enrollmentRepo.IsEnrolled(userID, courseID)
	if err != nil {
		return nil, err
	}

	if isEnrolled {
		return nil, fmt.Errorf("user is already enrolled in this course")
	}

	// Check if course exists
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found")
	}

	// Check enrollment limit
	if course.MaxEnrollments > 0 && course.EnrollmentCount >= course.MaxEnrollments {
		return nil, fmt.Errorf("course is full, cannot enroll")
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserID:           userID,
		CourseID:         courseID,
		CompletionStatus: "not_started",
		OverallProgress:  0,
	}

	if err := s.enrollmentRepo.Create(enrollment); err != nil {
		return nil, err
	}

	// Increment enrollment count
	_ = s.courseRepo.UpdateEnrollmentCount(courseID, 1)

	return enrollment, nil
}

// GetUserEnrollments gets all enrollments for a user
func (s *EnrollmentService) GetUserEnrollments(userID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	return s.enrollmentRepo.GetUserEnrollments(userID, page, pageSize)
}

// GetInProgressCourses gets in-progress courses for a user
func (s *EnrollmentService) GetInProgressCourses(userID uint) ([]models.Enrollment, error) {
	return s.enrollmentRepo.GetUserInProgressCourses(userID)
}

// GetCompletedCourses gets completed courses for a user
func (s *EnrollmentService) GetCompletedCourses(userID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	return s.enrollmentRepo.GetUserCompletedCourses(userID, page, pageSize)
}

// GetMandatoryCourses gets mandatory course enrollments for a user
func (s *EnrollmentService) GetMandatoryCourses(userID uint) ([]models.Enrollment, error) {
	return s.enrollmentRepo.GetUserMandatoryCourses(userID)
}

// StartCourse marks enrollment as started
func (s *EnrollmentService) StartCourse(userID, courseID uint) error {
	return s.enrollmentRepo.MarkAsStarted(userID, courseID)
}

// CompleteCourse marks course as completed and awards coins
func (s *EnrollmentService) CompleteCourse(userID, courseID uint, finalScore int) (*models.Enrollment, error) {
	// Mark enrollment as completed
	if err := s.enrollmentRepo.MarkAsCompleted(userID, courseID, finalScore); err != nil {
		return nil, err
	}

	// Get enrollment
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(userID, courseID)
	if err != nil {
		return nil, err
	}

	// Get course
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, err
	}

	// Award coins
	if finalScore >= course.PassingScore {
		coinAmount := int64(course.CoinsReward)
		_ = s.userRepo.UpdateCoins(userID, coinAmount)

		// Log coin transaction
		_ = s.coinTransactionRepo.Create(&models.CoinTransaction{
			UserID:          userID,
			Amount:          coinAmount,
			TransactionType: "earned",
			Reason:          fmt.Sprintf("Course Completion: %s", course.Title),
			ReferenceID:     &course.ID,
			ReferenceType:   "course",
		})

		// Generate certificate
		certificateNumber := fmt.Sprintf("CERT-%d-%d-%d", userID, courseID, time.Now().Unix())
		certificate := &models.Certificate{
			UserID:            userID,
			CourseID:          courseID,
			CertificateNumber: certificateNumber,
			Score:             finalScore,
		}
		_ = s.certificateRepo.Create(certificate)
	}

	// Increment course completion count
	_ = s.courseRepo.UpdateCompletionCount(courseID, 1)

	return enrollment, nil
}

// UpdateProgress updates course progress
func (s *EnrollmentService) UpdateProgress(userID, courseID uint, progress int) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	return s.enrollmentRepo.UpdateProgress(userID, courseID, progress)
}

// ConvertEnrollmentToDTO converts enrollment model to DTO
func ConvertEnrollmentToDTO(enrollment *models.Enrollment) *EnrollmentDTO {
	return &EnrollmentDTO{
		ID:               enrollment.ID,
		UserID:           enrollment.UserID,
		CourseID:         enrollment.CourseID,
		CompletionStatus: enrollment.CompletionStatus,
		OverallProgress:  enrollment.OverallProgress,
		FinalScore:       enrollment.FinalScore,
		IsPassed:         enrollment.IsPassed,
		EnrolledAt:       enrollment.EnrolledAt,
		CompletedAt:      enrollment.CompletedAt,
	}
}
