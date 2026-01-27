package services

import (
	"errors"
	"lms-go-be/models"
	"lms-go-be/repositories"
	"lms-go-be/utils"
	"time"

	"gorm.io/gorm"
)

// CourseService handles course business logic
type CourseService struct {
	courseRepo     *repositories.CourseRepository
	enrollmentRepo *repositories.EnrollmentRepository
	lessonRepo     *repositories.LessonRepository
}

// NewCourseService creates a new course service instance
func NewCourseService(
	courseRepo *repositories.CourseRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
	lessonRepo *repositories.LessonRepository,
) *CourseService {
	return &CourseService{
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
		lessonRepo:     lessonRepo,
	}
}

// CreateCourse creates a new course
// Parameters:
//   - course: course model to create
//
// Returns: created course and error if any
func (s *CourseService) CreateCourse(course *models.Course) (*models.Course, error) {
	if course.Title == "" {
		return nil, errors.New("course title is required")
	}

	course.Status = "draft"
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()

	return s.courseRepo.Create(course)
}

// GetCourseByID retrieves a course with all details
// Parameters:
//   - courseID: course ID
//
// Returns: course model and error if any
func (s *CourseService) GetCourseByID(courseID string) (*models.Course, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}
	return course, nil
}

// UpdateCourse updates an existing course
// Parameters:
//   - course: course model with updated data
//
// Returns: updated course and error if any
func (s *CourseService) UpdateCourse(course *models.Course) (*models.Course, error) {
	course.UpdatedAt = time.Now()
	return s.courseRepo.Update(course)
}

// PublishCourse publishes a course
// Parameters:
//   - courseID: course ID
//
// Returns: error if any
func (s *CourseService) PublishCourse(courseID string) error {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return err
	}

	course.Status = "published"
	course.UpdatedAt = time.Now()
	_, err = s.courseRepo.Update(course)
	return err
}

// GetAllCourses retrieves all published courses with pagination
// Parameters:
//   - page: page number
//   - pageSize: number of records per page
//   - filters: optional filters
//
// Returns: slice of courses, total count, and error if any
func (s *CourseService) GetAllCourses(page, pageSize int, filters map[string]interface{}) ([]models.Course, int64, error) {
	return s.courseRepo.GetAll(page, pageSize, filters)
}

// GetMandatoryCourses retrieves all mandatory courses
// Returns: slice of mandatory courses and error if any
func (s *CourseService) GetMandatoryCourses() ([]models.Course, error) {
	return s.courseRepo.GetMandatoryCourses()
}

// SearchCourses searches courses by keyword
// Parameters:
//   - keyword: search keyword
//   - page: page number
//   - pageSize: number of records per page
//
// Returns: slice of courses, total count, and error if any
func (s *CourseService) SearchCourses(keyword string, page, pageSize int) ([]models.Course, int64, error) {
	if keyword == "" {
		return []models.Course{}, 0, errors.New("search keyword is required")
	}
	return s.courseRepo.SearchCourses(keyword, page, pageSize)
}

// DeleteCourse deletes a course
// Parameters:
//   - courseID: course ID
//
// Returns: error if any
func (s *CourseService) DeleteCourse(courseID string) error {
	return s.courseRepo.Delete(courseID)
}

// GetCourseStats retrieves course statistics
// Parameters:
//   - courseID: course ID
//
// Returns: map with stats and error if any
func (s *CourseService) GetCourseStats(courseID string) (map[string]interface{}, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, err
	}

	enrollmentStats, err := s.enrollmentRepo.GetEnrollmentStats(courseID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"id":                course.ID,
		"title":             course.Title,
		"total_duration":    course.TotalDuration,
		"total_enrollments": enrollmentStats["total"],
		"completed":         enrollmentStats["completed"],
		"in_progress":       enrollmentStats["in_progress"],
		"completion_rate":   float64(0),
	}

	if enrollmentStats["total"] > 0 {
		completionRate := float64(enrollmentStats["completed"]) / float64(enrollmentStats["total"]) * 100
		stats["completion_rate"] = utils.RoundFloat(completionRate, 2)
	}

	return stats, nil
}

// ===== ENROLLMENT SERVICE =====

// EnrollmentService handles enrollment business logic
type EnrollmentService struct {
	enrollmentRepo     *repositories.EnrollmentRepository
	courseRepo         *repositories.CourseRepository
	lessonRepo         *repositories.LessonRepository
	lessonProgressRepo *repositories.LessonProgressRepository
	userRepo           *repositories.UserRepository
}

// NewEnrollmentService creates a new enrollment service instance
func NewEnrollmentService(
	enrollmentRepo *repositories.EnrollmentRepository,
	courseRepo *repositories.CourseRepository,
	lessonRepo *repositories.LessonRepository,
	lessonProgressRepo *repositories.LessonProgressRepository,
	userRepo *repositories.UserRepository,
) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:     enrollmentRepo,
		courseRepo:         courseRepo,
		lessonRepo:         lessonRepo,
		lessonProgressRepo: lessonProgressRepo,
		userRepo:           userRepo,
	}
}

// EnrollUser enrolls a user in a course
// Parameters:
//   - userID: user ID
//   - courseID: course ID
//
// Returns: enrollment model and error if any
func (s *EnrollmentService) EnrollUser(userID, courseID string) (*models.Enrollment, error) {
	// Check if user already enrolled
	existing, _ := s.enrollmentRepo.GetUserCourseEnrollment(userID, courseID)
	if existing != nil {
		return nil, errors.New("user is already enrolled in this course")
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		Status:     "enrolled",
		Progress:   0,
		EnrolledAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.enrollmentRepo.Create(enrollment)
}

// GetUserCourses retrieves all courses enrolled by a user
// Parameters:
//   - userID: user ID
//
// Returns: slice of enrollments and error if any
func (s *EnrollmentService) GetUserCourses(userID string) ([]models.Enrollment, error) {
	return s.enrollmentRepo.GetUserEnrollments(userID)
}

// GetMandatoryCoursesForUser retrieves mandatory courses not yet completed by user
// Parameters:
//   - userID: user ID
//
// Returns: slice of enrollments and error if any
func (s *EnrollmentService) GetMandatoryCoursesForUser(userID string) ([]models.Enrollment, error) {
	mandatoryCourses, err := s.courseRepo.GetMandatoryCourses()
	if err != nil {
		return nil, err
	}

	var enrollments []models.Enrollment
	for _, course := range mandatoryCourses {
		enrollment, err := s.enrollmentRepo.GetUserCourseEnrollment(userID, course.ID)
		if err == nil && enrollment != nil && enrollment.Status != "completed" {
			enrollment.Course = &course
			enrollments = append(enrollments, *enrollment)
		}
	}

	return enrollments, nil
}

// GetInProgressCourses retrieves in-progress courses for a user
// Parameters:
//   - userID: user ID
//
// Returns: slice of enrollments and error if any
func (s *EnrollmentService) GetInProgressCourses(userID string) ([]models.Enrollment, error) {
	return s.enrollmentRepo.GetInProgressEnrollments(userID)
}

// UpdateEnrollmentProgress updates enrollment progress
// Parameters:
//   - enrollmentID: enrollment ID
//
// Returns: error if any
func (s *EnrollmentService) UpdateEnrollmentProgress(enrollmentID string) error {
	enrollment, err := s.enrollmentRepo.GetByID(enrollmentID)
	if err != nil {
		return err
	}

	// Get course lessons
	lessons, err := s.lessonRepo.GetCourseLessons(enrollment.CourseID)
	if err != nil || len(lessons) == 0 {
		return errors.New("no lessons found for course")
	}

	// Count completed lessons
	completedCount, err := s.lessonProgressRepo.GetCompletedLessonCount(enrollment.UserID, enrollment.CourseID)
	if err != nil {
		return err
	}

	// Calculate progress percentage
	progress := utils.CalculateCompletionPercentage(int(completedCount), len(lessons))
	enrollment.Progress = progress

	// Update status
	if progress == 100 {
		now := time.Now()
		enrollment.CompletionDate = &now
		enrollment.Status = "completed"
	} else if progress > 0 {
		enrollment.Status = "in_progress"
	}

	enrollment.UpdatedAt = time.Now()
	_, err = s.enrollmentRepo.Update(enrollment)
	return err
}

// CompleteEnrollment marks an enrollment as completed
// Parameters:
//   - enrollmentID: enrollment ID
//   - score: final score (if applicable)
//
// Returns: error if any
func (s *EnrollmentService) CompleteEnrollment(enrollmentID string, score *int) error {
	enrollment, err := s.enrollmentRepo.GetByID(enrollmentID)
	if err != nil {
		return err
	}

	now := time.Now()
	enrollment.Status = "completed"
	enrollment.CompletionDate = &now
	enrollment.FinalScore = score
	enrollment.Progress = 100
	enrollment.UpdatedAt = now

	_, err = s.enrollmentRepo.Update(enrollment)
	return err
}

// GetCompletedCourses retrieves completed courses for a user
// Parameters:
//   - userID: user ID
//
// Returns: slice of enrollments and error if any
func (s *EnrollmentService) GetCompletedCourses(userID string) ([]models.Enrollment, error) {
	return s.enrollmentRepo.GetCompletedEnrollments(userID)
}

// GetUserStats retrieves enrollment statistics for a user
// Parameters:
//   - userID: user ID
//
// Returns: map with stats and error if any
func (s *EnrollmentService) GetUserStats(userID string) (map[string]interface{}, error) {
	enrollments, err := s.enrollmentRepo.GetUserEnrollments(userID)
	if err != nil {
		return nil, err
	}

	completedCount, _ := s.enrollmentRepo.GetCompletionCount(userID)

	var totalHours int
	var inProgressCount int
	for _, e := range enrollments {
		if e.Status == "in_progress" || e.Status == "enrolled" {
			inProgressCount++
		}
		if e.Course != nil {
			totalHours += e.Course.TotalDuration
		}
	}

	stats := map[string]interface{}{
		"total_courses_enrolled": len(enrollments),
		"completed_courses":      completedCount,
		"in_progress_courses":    inProgressCount,
		"total_learning_hours":   totalHours,
	}

	return stats, nil
}
