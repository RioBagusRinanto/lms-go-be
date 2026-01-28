package repository

import (
	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// UserProgressRepository handles user progress database operations
type UserProgressRepository struct {
	db *gorm.DB
}

// NewUserProgressRepository creates a new user progress repository
func NewUserProgressRepository(db *gorm.DB) *UserProgressRepository {
	return &UserProgressRepository{db: db}
}

// Create creates a new user progress record
func (r *UserProgressRepository) Create(progress *models.UserProgress) error {
	return r.db.Create(progress).Error
}

// GetByID gets a progress record by ID
func (r *UserProgressRepository) GetByID(id uint) (*models.UserProgress, error) {
	var progress models.UserProgress
	if err := r.db.First(&progress, id).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetByUserLessonCourse gets progress for specific user, lesson, and course
func (r *UserProgressRepository) GetByUserLessonCourse(userID, lessonID, courseID uint) (*models.UserProgress, error) {
	var progress models.UserProgress
	if err := r.db.Where("user_id = ? AND lesson_id = ? AND course_id = ?",
		userID, lessonID, courseID).First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// Update updates a progress record
func (r *UserProgressRepository) Update(progress *models.UserProgress) error {
	return r.db.Save(progress).Error
}

// Delete deletes a progress record (soft delete)
func (r *UserProgressRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserProgress{}, id).Error
}

// GetUserCourseProgress gets all lessons' progress for a user in a course
func (r *UserProgressRepository) GetUserCourseProgress(userID, courseID uint) ([]models.UserProgress, error) {
	var progresses []models.UserProgress
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		Preload("Lesson").Order("lesson_id").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}

// GetCourseLessonProgress gets progress for specific lesson
func (r *UserProgressRepository) GetCourseLessonProgress(courseID, lessonID uint) ([]models.UserProgress, error) {
	var progresses []models.UserProgress
	if err := r.db.Where("course_id = ? AND lesson_id = ?", courseID, lessonID).
		Preload("User").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}

// MarkLessonComplete marks a lesson as completed
func (r *UserProgressRepository) MarkLessonComplete(userID, lessonID, courseID uint) error {
	return r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND lesson_id = ? AND course_id = ?", userID, lessonID, courseID).
		Updates(map[string]interface{}{
			"is_completed":        true,
			"completed_at":        gorm.Expr("NOW()"),
			"progress_percentage": 100,
		}).Error
}

// UpdateWatchedDuration updates watched duration and calculates percentage
func (r *UserProgressRepository) UpdateWatchedDuration(userID, lessonID, courseID uint, watched int) error {
	return r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND lesson_id = ? AND course_id = ?", userID, lessonID, courseID).
		Update("watched_duration", watched).Error
}

// GetUserTotalLearningHours gets total learning hours for a user
func (r *UserProgressRepository) GetUserTotalLearningHours(userID uint) (float64, error) {
	var totalSeconds int64
	if err := r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND is_completed = ?", userID, true).
		Pluck("SUM(total_duration)", &totalSeconds).Error; err != nil {
		return 0, err
	}
	return float64(totalSeconds) / 3600, nil // Convert to hours
}

// GetCompletedLessonsCount gets count of completed lessons for a user
func (r *UserProgressRepository) GetCompletedLessonsCount(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND is_completed = ?", userID, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetUncompletedLessons gets incomplete lessons for a user in a course
func (r *UserProgressRepository) GetUncompletedLessons(userID, courseID uint) ([]models.UserProgress, error) {
	var progresses []models.UserProgress
	if err := r.db.Where("user_id = ? AND course_id = ? AND is_completed = ?", userID, courseID, false).
		Preload("Lesson").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}

// UpdateLastAccessed updates last accessed timestamp
func (r *UserProgressRepository) UpdateLastAccessed(userID, lessonID, courseID uint) error {
	return r.db.Model(&models.UserProgress{}).
		Where("user_id = ? AND lesson_id = ? AND course_id = ?", userID, lessonID, courseID).
		Update("last_accessed_at", gorm.Expr("NOW()")).Error
}
