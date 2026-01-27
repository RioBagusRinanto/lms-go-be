package repository

import (
	"fmt"

	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// CourseRepository handles course database operations
type CourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new course repository
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course
func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// GetByID gets a course by ID with relations
func (r *CourseRepository) GetByID(id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.Preload("Instructor").Preload("Lessons").Preload("Quizzes").
		First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

// Update updates a course
func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Delete deletes a course (soft delete)
func (r *CourseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}

// GetAll gets all published courses with pagination
func (r *CourseRepository) GetAll(page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	if err := r.db.Where("is_published = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("is_published = ?", true).Preload("Instructor").
		Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// GetByCategory gets courses by category
func (r *CourseRepository) GetByCategory(category string, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	if err := r.db.Where("category = ? AND is_published = ?", category, true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("category = ? AND is_published = ?", category, true).
		Preload("Instructor").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// GetMandatoryCourses gets all mandatory courses
func (r *CourseRepository) GetMandatoryCourses() ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("is_mandatory = ? AND is_published = ?", true, true).
		Preload("Instructor").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// GetByInstructor gets courses by instructor ID
func (r *CourseRepository) GetByInstructor(instructorID uint, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	if err := r.db.Where("instructor_id = ?", instructorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("instructor_id = ?", instructorID).
		Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// SearchCourses searches courses by title or description
func (r *CourseRepository) SearchCourses(query string, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	searchQuery := fmt.Sprintf("%%%s%%", query)

	if err := r.db.Where("(title ILIKE ? OR description ILIKE ?) AND is_published = ?",
		searchQuery, searchQuery, true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("(title ILIKE ? OR description ILIKE ?) AND is_published = ?",
		searchQuery, searchQuery, true).Preload("Instructor").
		Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// GetTopRatedCourses gets top rated courses
func (r *CourseRepository) GetTopRatedCourses(limit int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("is_published = ?", true).
		Order("average_rating DESC").Limit(limit).
		Preload("Instructor").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// GetPopularCourses gets most enrolled courses
func (r *CourseRepository) GetPopularCourses(limit int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Where("is_published = ?", true).
		Order("enrollment_count DESC").Limit(limit).
		Preload("Instructor").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// UpdateEnrollmentCount increments enrollment count
func (r *CourseRepository) UpdateEnrollmentCount(courseID uint, increment int) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("enrollment_count", gorm.Expr("enrollment_count + ?", increment)).Error
}

// UpdateCompletionCount increments completion count
func (r *CourseRepository) UpdateCompletionCount(courseID uint, increment int) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		Update("completion_count", gorm.Expr("completion_count + ?", increment)).Error
}

// GetCategoryList gets all distinct course categories
func (r *CourseRepository) GetCategoryList() ([]string, error) {
	var categories []string
	if err := r.db.Model(&models.Course{}).Where("is_published = ?", true).
		Distinct("category").Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCourseStats gets statistics for a course
func (r *CourseRepository) GetCourseStats(courseID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var enrollmentCount int64
	var completionCount int64
	var avgRating float64

	r.db.Model(&models.Enrollment{}).Where("course_id = ?", courseID).Count(&enrollmentCount)
	r.db.Model(&models.Enrollment{}).Where("course_id = ? AND is_passed = ?", courseID, true).Count(&completionCount)
	r.db.Model(&models.CourseReview{}).Where("course_id = ?", courseID).Select("AVG(rating)").Row().Scan(&avgRating)

	stats["enrollment_count"] = enrollmentCount
	stats["completion_count"] = completionCount
	stats["completion_rate"] = 0.0
	if enrollmentCount > 0 {
		stats["completion_rate"] = float64(completionCount) / float64(enrollmentCount) * 100
	}
	stats["average_rating"] = avgRating

	return stats, nil
}
