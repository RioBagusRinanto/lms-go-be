package repositories

import (
	"lms-go-be/models"

	"gorm.io/gorm"
)

// EnrollmentRepository handles enrollment data access
type EnrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new enrollment repository instance
func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create creates a new enrollment
// Parameters:
//   - enrollment: enrollment model to create
//
// Returns: created enrollment and error if any
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) (*models.Enrollment, error) {
	if err := r.db.Create(enrollment).Error; err != nil {
		return nil, err
	}
	return enrollment, nil
}

// GetByID retrieves an enrollment by ID
// Parameters:
//   - id: enrollment ID
//
// Returns: enrollment model and error if not found
func (r *EnrollmentRepository) GetByID(id string) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Preload("User").Preload("Course").First(&enrollment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// GetUserEnrollments retrieves all enrollments for a user
// Parameters:
//   - userID: user ID
//
// Returns: slice of enrollments and error if any
func (r *EnrollmentRepository) GetUserEnrollments(userID string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Where("user_id = ?", userID).
		Preload("Course").
		Preload("Course.Lessons").
		Order("updated_at DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// GetCourseEnrollments retrieves all enrollments for a course
// Parameters:
//   - courseID: course ID
//
// Returns: slice of enrollments and error if any
func (r *EnrollmentRepository) GetCourseEnrollments(courseID string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Where("course_id = ?", courseID).
		Preload("User").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// GetUserCourseEnrollment retrieves a specific user's enrollment in a course
// Parameters:
//   - userID: user ID
//   - courseID: course ID
//
// Returns: enrollment model and error if not found
func (r *EnrollmentRepository) GetUserCourseEnrollment(userID, courseID string) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&enrollment).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// Update updates an existing enrollment
// Parameters:
//   - enrollment: enrollment model with updated data
//
// Returns: updated enrollment and error if any
func (r *EnrollmentRepository) Update(enrollment *models.Enrollment) (*models.Enrollment, error) {
	if err := r.db.Save(enrollment).Error; err != nil {
		return nil, err
	}
	return enrollment, nil
}

// UpdateProgress updates enrollment progress
// Parameters:
//   - enrollmentID: enrollment ID
//   - progress: progress percentage (0-100)
//
// Returns: error if any
func (r *EnrollmentRepository) UpdateProgress(enrollmentID string, progress int) error {
	return r.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentID).
		Update("progress", progress).Error
}

// UpdateStatus updates enrollment status
// Parameters:
//   - enrollmentID: enrollment ID
//   - status: new status
//
// Returns: error if any
func (r *EnrollmentRepository) UpdateStatus(enrollmentID, status string) error {
	return r.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentID).
		Update("status", status).Error
}

// GetCompletedEnrollments retrieves completed enrollments for a user
// Parameters:
//   - userID: user ID
//
// Returns: slice of completed enrollments and error if any
func (r *EnrollmentRepository) GetCompletedEnrollments(userID string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Where("user_id = ? AND status = ?", userID, "completed").
		Preload("Course").
		Order("completion_date DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// GetInProgressEnrollments retrieves in-progress enrollments for a user
// Parameters:
//   - userID: user ID
//
// Returns: slice of in-progress enrollments and error if any
func (r *EnrollmentRepository) GetInProgressEnrollments(userID string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Where("user_id = ? AND status IN ?", userID, []string{"enrolled", "in_progress"}).
		Preload("Course").
		Order("updated_at DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// GetCompletionCount returns number of completed courses for a user
// Parameters:
//   - userID: user ID
//
// Returns: count and error if any
func (r *EnrollmentRepository) GetCompletionCount(userID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetEnrollmentStats returns enrollment statistics for a course
// Parameters:
//   - courseID: course ID
//
// Returns: map with stats and error if any
func (r *EnrollmentRepository) GetEnrollmentStats(courseID string) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Total enrollments
	var total int64
	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Completed
	var completed int64
	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ? AND status = ?", courseID, "completed").Count(&completed).Error; err != nil {
		return nil, err
	}
	stats["completed"] = completed

	// In progress
	var inProgress int64
	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ? AND status IN ?", courseID, []string{"enrolled", "in_progress"}).Count(&inProgress).Error; err != nil {
		return nil, err
	}
	stats["in_progress"] = inProgress

	return stats, nil
}

// Delete deletes an enrollment
// Parameters:
//   - id: enrollment ID
//
// Returns: error if any
func (r *EnrollmentRepository) Delete(id string) error {
	return r.db.Delete(&models.Enrollment{}, "id = ?", id).Error
}
