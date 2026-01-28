package repository

import (
	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// EnrollmentRepository handles enrollment database operations
type EnrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create creates a new enrollment
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

// GetByID gets an enrollment by ID
func (r *EnrollmentRepository) GetByID(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Preload("User").Preload("Course").First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// GetByUserAndCourse gets enrollment by user and course
func (r *EnrollmentRepository) GetByUserAndCourse(userID, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		Preload("User").Preload("Course").First(&enrollment).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// Update updates an enrollment
func (r *EnrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

// Delete deletes an enrollment (soft delete)
func (r *EnrollmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Enrollment{}, id).Error
}

// GetUserEnrollments gets all enrollments for a user
func (r *EnrollmentRepository) GetUserEnrollments(userID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	if err := r.db.Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Preload("Course").Preload("User").
		Offset(offset).Limit(pageSize).Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

// GetUserInProgressCourses gets in-progress enrollments for a user
func (r *EnrollmentRepository) GetUserInProgressCourses(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Where("user_id = ? AND completion_status = ?", userID, "in_progress").
		Preload("Course").Order("last_accessed_at DESC").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// GetUserCompletedCourses gets completed enrollments for a user
func (r *EnrollmentRepository) GetUserCompletedCourses(userID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	if err := r.db.Where("user_id = ? AND completion_status = ?", userID, "completed").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ? AND completion_status = ?", userID, "completed").
		Preload("Course").Order("completed_at DESC").
		Offset(offset).Limit(pageSize).Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

// GetUserMandatoryCourses gets mandatory course enrollments for a user
func (r *EnrollmentRepository) GetUserMandatoryCourses(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("enrollments.user_id = ? AND courses.is_mandatory = ?", userID, true).
		Preload("Course").Order("courses.mandatory_due_date ASC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// GetCourseEnrollments gets all enrollments for a course
func (r *EnrollmentRepository) GetCourseEnrollments(courseID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	if err := r.db.Model(&models.Enrollment{}).Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("course_id = ?", courseID).Preload("User").
		Offset(offset).Limit(pageSize).Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

// IsEnrolled checks if user is enrolled in a course
func (r *EnrollmentRepository) IsEnrolled(userID, courseID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetOverdueEnrollments gets overdue mandatory course enrollments
func (r *EnrollmentRepository) GetOverdueEnrollments() ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.is_mandatory = ? AND courses.mandatory_due_date < NOW() AND enrollments.completion_status != ?",
			true, "completed").
		Preload("User").Preload("Course").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// MarkAsStarted updates enrollment status to in_progress
func (r *EnrollmentRepository) MarkAsStarted(userID, courseID uint) error {
	return r.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Updates(map[string]interface{}{
			"completion_status": "in_progress",
			"last_accessed_at":  gorm.Expr("NOW()"),
		}).Error
}

// MarkAsCompleted updates enrollment status to completed
func (r *EnrollmentRepository) MarkAsCompleted(userID, courseID uint, finalScore int) error {
	return r.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Updates(map[string]interface{}{
			"completion_status": "completed",
			"completed_at":      gorm.Expr("NOW()"),
			"final_score":       finalScore,
			"is_passed":         finalScore >= 70,
			"overall_progress":  100,
		}).Error
}

// UpdateProgress updates enrollment progress
func (r *EnrollmentRepository) UpdateProgress(userID, courseID uint, progress int) error {
	return r.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Update("overall_progress", progress).Error
}
