package repository

import (
	"lms-go-be/internal/models"

	"gorm.io/gorm"
)

// QuizRepository handles quiz database operations
type QuizRepository struct {
	db *gorm.DB
}

// NewQuizRepository creates a new quiz repository
func NewQuizRepository(db *gorm.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

// Create creates a new quiz
func (r *QuizRepository) Create(quiz *models.Quiz) error {
	return r.db.Create(quiz).Error
}

// GetByID gets a quiz by ID with questions
func (r *QuizRepository) GetByID(id uint) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_number")
	}).Preload("Questions.Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_number")
	}).First(&quiz, id).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

// Update updates a quiz
func (r *QuizRepository) Update(quiz *models.Quiz) error {
	return r.db.Save(quiz).Error
}

// Delete deletes a quiz (soft delete)
func (r *QuizRepository) Delete(id uint) error {
	return r.db.Delete(&models.Quiz{}, id).Error
}

// GetByLesson gets quiz for a lesson
func (r *QuizRepository) GetByLesson(lessonID uint) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.Where("lesson_id = ?", lessonID).
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number")
		}).
		Preload("Questions.Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number")
		}).First(&quiz).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

// GetByCourse gets quizzes for a course
func (r *QuizRepository) GetByCourse(courseID uint) ([]models.Quiz, error) {
	var quizzes []models.Quiz
	if err := r.db.Where("course_id = ?", courseID).
		Preload("Questions").Order("title").Find(&quizzes).Error; err != nil {
		return nil, err
	}
	return quizzes, nil
}

// QuizAttemptRepository handles quiz attempt database operations
type QuizAttemptRepository struct {
	db *gorm.DB
}

// NewQuizAttemptRepository creates a new quiz attempt repository
func NewQuizAttemptRepository(db *gorm.DB) *QuizAttemptRepository {
	return &QuizAttemptRepository{db: db}
}

// Create creates a new quiz attempt
func (r *QuizAttemptRepository) Create(attempt *models.QuizAttempt) error {
	return r.db.Create(attempt).Error
}

// GetByID gets a quiz attempt by ID
func (r *QuizAttemptRepository) GetByID(id uint) (*models.QuizAttempt, error) {
	var attempt models.QuizAttempt
	if err := r.db.Preload("Answers").First(&attempt, id).Error; err != nil {
		return nil, err
	}
	return &attempt, nil
}

// Update updates a quiz attempt
func (r *QuizAttemptRepository) Update(attempt *models.QuizAttempt) error {
	return r.db.Save(attempt).Error
}

// Delete deletes a quiz attempt (soft delete)
func (r *QuizAttemptRepository) Delete(id uint) error {
	return r.db.Delete(&models.QuizAttempt{}, id).Error
}

// GetUserQuizAttempts gets all attempts by a user for a quiz
func (r *QuizAttemptRepository) GetUserQuizAttempts(userID, quizID uint) ([]models.QuizAttempt, error) {
	var attempts []models.QuizAttempt
	if err := r.db.Where("user_id = ? AND quiz_id = ?", userID, quizID).
		Order("started_at DESC").Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

// GetUserQuizLastAttempt gets the last attempt by a user for a quiz
func (r *QuizAttemptRepository) GetUserQuizLastAttempt(userID, quizID uint) (*models.QuizAttempt, error) {
	var attempt models.QuizAttempt
	if err := r.db.Where("user_id = ? AND quiz_id = ?", userID, quizID).
		Order("started_at DESC").First(&attempt).Error; err != nil {
		return nil, err
	}
	return &attempt, nil
}

// GetQuizAttemptCount gets number of attempts user has made for a quiz
func (r *QuizAttemptRepository) GetQuizAttemptCount(userID, quizID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.QuizAttempt{}).
		Where("user_id = ? AND quiz_id = ?", userID, quizID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetQuizStats gets statistics for a quiz
func (r *QuizAttemptRepository) GetQuizStats(quizID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalAttempts, passedAttempts int64
	var avgScore float64

	r.db.Model(&models.QuizAttempt{}).Where("quiz_id = ?", quizID).Count(&totalAttempts)
	r.db.Model(&models.QuizAttempt{}).Where("quiz_id = ? AND is_passed = ?", quizID, true).Count(&passedAttempts)
	r.db.Model(&models.QuizAttempt{}).Where("quiz_id = ?", quizID).Select("AVG(percentage)").Row().Scan(&avgScore)

	stats["total_attempts"] = totalAttempts
	stats["passed_attempts"] = passedAttempts
	stats["pass_rate"] = 0.0
	if totalAttempts > 0 {
		stats["pass_rate"] = float64(passedAttempts) / float64(totalAttempts) * 100
	}
	stats["average_score"] = avgScore

	return stats, nil
}

// CertificateRepository handles certificate database operations
type CertificateRepository struct {
	db *gorm.DB
}

// NewCertificateRepository creates a new certificate repository
func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

// Create creates a new certificate
func (r *CertificateRepository) Create(certificate *models.Certificate) error {
	return r.db.Create(certificate).Error
}

// GetByID gets a certificate by ID
func (r *CertificateRepository) GetByID(id uint) (*models.Certificate, error) {
	var certificate models.Certificate
	if err := r.db.Preload("User").Preload("Course").First(&certificate, id).Error; err != nil {
		return nil, err
	}
	return &certificate, nil
}

// GetByUserAndCourse gets certificate for user and course
func (r *CertificateRepository) GetByUserAndCourse(userID, courseID uint) (*models.Certificate, error) {
	var certificate models.Certificate
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		Preload("User").Preload("Course").First(&certificate).Error; err != nil {
		return nil, err
	}
	return &certificate, nil
}

// GetUserCertificates gets all certificates for a user
func (r *CertificateRepository) GetUserCertificates(userID uint, page, pageSize int) ([]models.Certificate, int64, error) {
	var certificates []models.Certificate
	var total int64

	if err := r.db.Model(&models.Certificate{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Preload("Course").
		Order("issued_at DESC").Offset(offset).Limit(pageSize).Find(&certificates).Error; err != nil {
		return nil, 0, err
	}

	return certificates, total, nil
}

// GetCourseCertificates gets all certificates issued for a course
func (r *CertificateRepository) GetCourseCertificates(courseID uint) ([]models.Certificate, error) {
	var certificates []models.Certificate
	if err := r.db.Where("course_id = ?", courseID).
		Preload("User").Order("issued_at DESC").Find(&certificates).Error; err != nil {
		return nil, err
	}
	return certificates, nil
}

// GetCertificateByCertificateNumber gets certificate by certificate number
func (r *CertificateRepository) GetCertificateByCertificateNumber(certificateNumber string) (*models.Certificate, error) {
	var certificate models.Certificate
	if err := r.db.Where("certificate_number = ?", certificateNumber).
		Preload("User").Preload("Course").First(&certificate).Error; err != nil {
		return nil, err
	}
	return &certificate, nil
}

// CountUserCertificates counts certificates for a user
func (r *CertificateRepository) CountUserCertificates(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Certificate{}).
		Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
