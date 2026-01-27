package service

import (
	"fmt"
	"time"

	"lms-go-be/internal/models"
	"lms-go-be/internal/repository"
)

// CourseService handles course business logic
type CourseService struct {
	courseRepo     *repository.CourseRepository
	enrollmentRepo *repository.EnrollmentRepository
	reviewRepo     *repository.CourseReviewRepository
}

// NewCourseService creates a new course service
func NewCourseService(
	courseRepo *repository.CourseRepository,
	enrollmentRepo *repository.EnrollmentRepository,
	reviewRepo *repository.CourseReviewRepository,
) *CourseService {
	return &CourseService{
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
		reviewRepo:     reviewRepo,
	}
}

// CreateCourseRequest represents create course request
type CreateCourseRequest struct {
	Title            string     `json:"title" binding:"required"`
	Description      string     `json:"description"`
	Category         string     `json:"category" binding:"required"`
	DurationMinutes  int        `json:"duration_minutes" binding:"required,min=1"`
	DifficultyLevel  string     `json:"difficulty_level"`
	PassingScore     int        `json:"passing_score"`
	IsMandatory      bool       `json:"is_mandatory"`
	MandatoryDueDate *time.Time `json:"mandatory_due_date"`
	CoinsReward      int        `json:"coins_reward"`
}

// CourseDTO represents course data transfer object
type CourseDTO struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	DurationMinutes int       `json:"duration_minutes"`
	DifficultyLevel string    `json:"difficulty_level"`
	PassingScore    int       `json:"passing_score"`
	IsMandatory     bool      `json:"is_mandatory"`
	IsPublished     bool      `json:"is_published"`
	EnrollmentCount int       `json:"enrollment_count"`
	CompletionCount int       `json:"completion_count"`
	AverageRating   float64   `json:"average_rating"`
	CoinsReward     int       `json:"coins_reward"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(instructorID uint, req CreateCourseRequest) (*models.Course, error) {
	course := &models.Course{
		Title:            req.Title,
		Description:      req.Description,
		Category:         req.Category,
		InstructorID:     instructorID,
		DurationMinutes:  req.DurationMinutes,
		PassingScore:     req.PassingScore,
		IsMandatory:      req.IsMandatory,
		MandatoryDueDate: req.MandatoryDueDate,
		CoinsReward:      req.CoinsReward,
		IsPublished:      false,
	}

	if course.PassingScore == 0 {
		course.PassingScore = 70
	}

	if course.CoinsReward == 0 {
		course.CoinsReward = 100
	}

	if req.DifficultyLevel != "" {
		course.DifficultyLevel = req.DifficultyLevel
	}

	if err := s.courseRepo.Create(course); err != nil {
		return nil, err
	}

	return course, nil
}

// GetCourse gets a course by ID
func (s *CourseService) GetCourse(courseID uint) (*models.Course, error) {
	return s.courseRepo.GetByID(courseID)
}

// GetAllCourses gets all published courses
func (s *CourseService) GetAllCourses(page, pageSize int) ([]models.Course, int64, error) {
	return s.courseRepo.GetAll(page, pageSize)
}

// GetCoursesByCategory gets courses by category
func (s *CourseService) GetCoursesByCategory(category string, page, pageSize int) ([]models.Course, int64, error) {
	return s.courseRepo.GetByCategory(category, page, pageSize)
}

// SearchCourses searches courses
func (s *CourseService) SearchCourses(query string, page, pageSize int) ([]models.Course, int64, error) {
	return s.courseRepo.SearchCourses(query, page, pageSize)
}

// GetMandatoryCourses gets all mandatory courses
func (s *CourseService) GetMandatoryCourses() ([]models.Course, error) {
	return s.courseRepo.GetMandatoryCourses()
}

// GetTopRatedCourses gets top rated courses
func (s *CourseService) GetTopRatedCourses(limit int) ([]models.Course, error) {
	return s.courseRepo.GetTopRatedCourses(limit)
}

// GetPopularCourses gets popular courses
func (s *CourseService) GetPopularCourses(limit int) ([]models.Course, error) {
	return s.courseRepo.GetPopularCourses(limit)
}

// PublishCourse publishes a course
func (s *CourseService) PublishCourse(courseID uint) (*models.Course, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, err
	}

	course.IsPublished = true
	if err := s.courseRepo.Update(course); err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourse updates a course
func (s *CourseService) UpdateCourse(courseID uint, req CreateCourseRequest) (*models.Course, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, err
	}

	course.Title = req.Title
	course.Description = req.Description
	course.Category = req.Category
	course.DurationMinutes = req.DurationMinutes
	course.PassingScore = req.PassingScore
	course.IsMandatory = req.IsMandatory
	course.MandatoryDueDate = req.MandatoryDueDate
	course.CoinsReward = req.CoinsReward

	if req.DifficultyLevel != "" {
		course.DifficultyLevel = req.DifficultyLevel
	}

	if err := s.courseRepo.Update(course); err != nil {
		return nil, err
	}

	return course, nil
}

// AddReview adds a review to a course
func (s *CourseService) AddReview(userID, courseID uint, rating int, reviewText string) (*models.CourseReview, error) {
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Check if user already reviewed
	existingReview, err := s.reviewRepo.GetUserReview(userID, courseID)
	if err == nil && existingReview != nil {
		return nil, fmt.Errorf("user already reviewed this course")
	}

	review := &models.CourseReview{
		UserID:     userID,
		CourseID:   courseID,
		Rating:     rating,
		ReviewText: reviewText,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	// Update course average rating
	avgRating, _ := s.reviewRepo.GetCourseAverageRating(courseID)
	course, _ := s.courseRepo.GetByID(courseID)
	if course != nil {
		course.AverageRating = avgRating
		_ = s.courseRepo.Update(course)
	}

	return review, nil
}

// GetCourseReviews gets reviews for a course
func (s *CourseService) GetCourseReviews(courseID uint, page, pageSize int) ([]models.CourseReview, int64, error) {
	return s.reviewRepo.GetCourseReviews(courseID, page, pageSize)
}

// ConvertCourseToDTO converts course model to DTO
func ConvertCourseToDTO(course *models.Course) *CourseDTO {
	return &CourseDTO{
		ID:              course.ID,
		Title:           course.Title,
		Description:     course.Description,
		Category:        course.Category,
		DurationMinutes: course.DurationMinutes,
		DifficultyLevel: course.DifficultyLevel,
		PassingScore:    course.PassingScore,
		IsMandatory:     course.IsMandatory,
		IsPublished:     course.IsPublished,
		EnrollmentCount: course.EnrollmentCount,
		CompletionCount: course.CompletionCount,
		AverageRating:   course.AverageRating,
		CoinsReward:     course.CoinsReward,
		CreatedAt:       course.CreatedAt,
	}
}
