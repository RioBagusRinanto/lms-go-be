package repository

import (
	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// LearningReportRepository handles learning report database operations
type LearningReportRepository struct {
	db *gorm.DB
}

// NewLearningReportRepository creates a new learning report repository
func NewLearningReportRepository(db *gorm.DB) *LearningReportRepository {
	return &LearningReportRepository{db: db}
}

// Create creates a new learning report
func (r *LearningReportRepository) Create(report *models.LearningReport) error {
	return r.db.Create(report).Error
}

// GetByID gets a report by ID
func (r *LearningReportRepository) GetByID(id uint) (*models.LearningReport, error) {
	var report models.LearningReport
	if err := r.db.Preload("User").First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// GetUserReport gets the latest report for a user
func (r *LearningReportRepository) GetUserReport(userID uint) (*models.LearningReport, error) {
	var report models.LearningReport
	if err := r.db.Where("user_id = ? AND report_type = ?", userID, "user").
		Order("generated_at DESC").First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// GetOrganizationReport gets organization-wide report
func (r *LearningReportRepository) GetOrganizationReport() (*models.LearningReport, error) {
	var report models.LearningReport
	if err := r.db.Where("report_type = ?", "organization").
		Order("generated_at DESC").First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// GetCourseReport gets report for a specific course
func (r *LearningReportRepository) GetCourseReport(courseID uint) (*models.LearningReport, error) {
	var report models.LearningReport
	if err := r.db.Where("report_type = ? AND (SELECT id FROM courses WHERE id = ? LIMIT 1) IS NOT NULL", "course", courseID).
		Order("generated_at DESC").First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// SystemAuditLogRepository handles system audit log database operations
type SystemAuditLogRepository struct {
	db *gorm.DB
}

// NewSystemAuditLogRepository creates a new system audit log repository
func NewSystemAuditLogRepository(db *gorm.DB) *SystemAuditLogRepository {
	return &SystemAuditLogRepository{db: db}
}

// Create creates a new audit log
func (r *SystemAuditLogRepository) Create(log *models.SystemAuditLog) error {
	return r.db.Create(log).Error
}

// GetByID gets an audit log by ID
func (r *SystemAuditLogRepository) GetByID(id uint) (*models.SystemAuditLog, error) {
	var log models.SystemAuditLog
	if err := r.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// GetAll gets all audit logs with pagination
func (r *SystemAuditLogRepository) GetAll(page, pageSize int) ([]models.SystemAuditLog, int64, error) {
	var logs []models.SystemAuditLog
	var total int64

	if err := r.db.Model(&models.SystemAuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetUserLogs gets all logs for a specific user
func (r *SystemAuditLogRepository) GetUserLogs(userID uint, page, pageSize int) ([]models.SystemAuditLog, int64, error) {
	var logs []models.SystemAuditLog
	var total int64

	if err := r.db.Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogsByAction gets logs by action type
func (r *SystemAuditLogRepository) GetLogsByAction(action string, page, pageSize int) ([]models.SystemAuditLog, int64, error) {
	var logs []models.SystemAuditLog
	var total int64

	if err := r.db.Where("action = ?", action).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("action = ?", action).Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DownloadLogRepository handles download log database operations
type DownloadLogRepository struct {
	db *gorm.DB
}

// NewDownloadLogRepository creates a new download log repository
func NewDownloadLogRepository(db *gorm.DB) *DownloadLogRepository {
	return &DownloadLogRepository{db: db}
}

// Create creates a new download log
func (r *DownloadLogRepository) Create(log *models.DownloadLog) error {
	return r.db.Create(log).Error
}

// GetUserDownloads gets all downloads by a user
func (r *DownloadLogRepository) GetUserDownloads(userID uint, page, pageSize int) ([]models.DownloadLog, int64, error) {
	var logs []models.DownloadLog
	var total int64

	if err := r.db.Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Preload("Material").
		Order("downloaded_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetMaterialDownloadCount gets download count for a material
func (r *DownloadLogRepository) GetMaterialDownloadCount(materialID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.DownloadLog{}).
		Where("material_id = ?", materialID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CourseReviewRepository handles course review database operations
type CourseReviewRepository struct {
	db *gorm.DB
}

// NewCourseReviewRepository creates a new course review repository
func NewCourseReviewRepository(db *gorm.DB) *CourseReviewRepository {
	return &CourseReviewRepository{db: db}
}

// Create creates a new course review
func (r *CourseReviewRepository) Create(review *models.CourseReview) error {
	return r.db.Create(review).Error
}

// GetByID gets a review by ID
func (r *CourseReviewRepository) GetByID(id uint) (*models.CourseReview, error) {
	var review models.CourseReview
	if err := r.db.Preload("User").Preload("Course").First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// GetCourseReviews gets all reviews for a course
func (r *CourseReviewRepository) GetCourseReviews(courseID uint, page, pageSize int) ([]models.CourseReview, int64, error) {
	var reviews []models.CourseReview
	var total int64

	if err := r.db.Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("course_id = ?", courseID).Preload("User").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// GetUserReview gets review by user for a course
func (r *CourseReviewRepository) GetUserReview(userID, courseID uint) (*models.CourseReview, error) {
	var review models.CourseReview
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// Update updates a review
func (r *CourseReviewRepository) Update(review *models.CourseReview) error {
	return r.db.Save(review).Error
}

// Delete deletes a review (soft delete)
func (r *CourseReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.CourseReview{}, id).Error
}

// GetCourseAverageRating gets average rating for a course
func (r *CourseReviewRepository) GetCourseAverageRating(courseID uint) (float64, error) {
	var avgRating float64
	if err := r.db.Model(&models.CourseReview{}).
		Where("course_id = ?", courseID).
		Select("COALESCE(AVG(rating), 0)").Row().Scan(&avgRating).Error; err != nil {
		return 0, err
	}
	return avgRating, nil
}
