package repositories

import (
	"lms-go-be/models"

	"gorm.io/gorm"
)

// ===== LESSON REPOSITORY =====

// LessonRepository handles lesson data access
type LessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(lesson *models.Lesson) (*models.Lesson, error) {
	if err := r.db.Create(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (r *LessonRepository) GetByID(id string) (*models.Lesson, error) {
	var lesson models.Lesson
	if err := r.db.First(&lesson, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *LessonRepository) GetCourseLesson(courseID, lessonID string) (*models.Lesson, error) {
	var lesson models.Lesson
	if err := r.db.Where("id = ? AND course_id = ?", lessonID, courseID).First(&lesson).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *LessonRepository) GetCourseLessons(courseID string) ([]models.Lesson, error) {
	var lessons []models.Lesson
	if err := r.db.Where("course_id = ?", courseID).Order("order_index ASC").Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *LessonRepository) Update(lesson *models.Lesson) (*models.Lesson, error) {
	if err := r.db.Save(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (r *LessonRepository) Delete(id string) error {
	return r.db.Delete(&models.Lesson{}, "id = ?", id).Error
}

// ===== LESSON PROGRESS REPOSITORY =====

// LessonProgressRepository handles lesson progress data access
type LessonProgressRepository struct {
	db *gorm.DB
}

func NewLessonProgressRepository(db *gorm.DB) *LessonProgressRepository {
	return &LessonProgressRepository{db: db}
}

func (r *LessonProgressRepository) Create(progress *models.LessonProgress) (*models.LessonProgress, error) {
	if err := r.db.Create(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *LessonProgressRepository) GetUserLessonProgress(userID, lessonID string) (*models.LessonProgress, error) {
	var progress models.LessonProgress
	if err := r.db.Where("user_id = ? AND lesson_id = ?", userID, lessonID).First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *LessonProgressRepository) Update(progress *models.LessonProgress) (*models.LessonProgress, error) {
	if err := r.db.Save(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *LessonProgressRepository) GetCompletedLessonCount(userID, courseID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.LessonProgress{}).
		Joins("JOIN lessons ON lesson_progress.lesson_id = lessons.id").
		Where("lesson_progress.user_id = ? AND lessons.course_id = ? AND lesson_progress.is_completed = ?",
			userID, courseID, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ===== QUIZ REPOSITORY =====

// QuizRepository handles quiz data access
type QuizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) Create(quiz *models.Quiz) (*models.Quiz, error) {
	if err := r.db.Create(quiz).Error; err != nil {
		return nil, err
	}
	return quiz, nil
}

func (r *QuizRepository) GetByID(id string) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.Preload("Questions").Preload("Questions.Options").First(&quiz, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *QuizRepository) GetCourseQuizzes(courseID string) ([]models.Quiz, error) {
	var quizzes []models.Quiz
	if err := r.db.Where("course_id = ?", courseID).
		Preload("Questions").
		Preload("Questions.Options").
		Find(&quizzes).Error; err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (r *QuizRepository) Update(quiz *models.Quiz) (*models.Quiz, error) {
	if err := r.db.Save(quiz).Error; err != nil {
		return nil, err
	}
	return quiz, nil
}

func (r *QuizRepository) Delete(id string) error {
	return r.db.Delete(&models.Quiz{}, "id = ?", id).Error
}

// ===== QUIZ ATTEMPT REPOSITORY =====

// QuizAttemptRepository handles quiz attempt data access
type QuizAttemptRepository struct {
	db *gorm.DB
}

func NewQuizAttemptRepository(db *gorm.DB) *QuizAttemptRepository {
	return &QuizAttemptRepository{db: db}
}

func (r *QuizAttemptRepository) Create(attempt *models.QuizAttempt) (*models.QuizAttempt, error) {
	if err := r.db.Create(attempt).Error; err != nil {
		return nil, err
	}
	return attempt, nil
}

func (r *QuizAttemptRepository) GetByID(id string) (*models.QuizAttempt, error) {
	var attempt models.QuizAttempt
	if err := r.db.First(&attempt, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *QuizAttemptRepository) GetUserQuizAttempts(userID, quizID string) ([]models.QuizAttempt, error) {
	var attempts []models.QuizAttempt
	if err := r.db.Where("user_id = ? AND quiz_id = ?", userID, quizID).
		Order("created_at DESC").
		Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *QuizAttemptRepository) Update(attempt *models.QuizAttempt) (*models.QuizAttempt, error) {
	if err := r.db.Save(attempt).Error; err != nil {
		return nil, err
	}
	return attempt, nil
}

// ===== CERTIFICATE REPOSITORY =====

// CertificateRepository handles certificate data access
type CertificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

func (r *CertificateRepository) Create(cert *models.Certificate) (*models.Certificate, error) {
	if err := r.db.Create(cert).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

func (r *CertificateRepository) GetByID(id string) (*models.Certificate, error) {
	var cert models.Certificate
	if err := r.db.First(&cert, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CertificateRepository) GetUserCertificates(userID string) ([]models.Certificate, error) {
	var certs []models.Certificate
	if err := r.db.Where("user_id = ?", userID).Preload("Course").Order("issued_date DESC").Find(&certs).Error; err != nil {
		return nil, err
	}
	return certs, nil
}

func (r *CertificateRepository) GetUserCourseCertificate(userID, courseID string) (*models.Certificate, error) {
	var cert models.Certificate
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// ===== COIN HISTORY REPOSITORY =====

// CoinHistoryRepository handles coin history data access
type CoinHistoryRepository struct {
	db *gorm.DB
}

func NewCoinHistoryRepository(db *gorm.DB) *CoinHistoryRepository {
	return &CoinHistoryRepository{db: db}
}

func (r *CoinHistoryRepository) Create(history *models.CoinHistory) (*models.CoinHistory, error) {
	if err := r.db.Create(history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *CoinHistoryRepository) GetUserHistory(userID string, limit int) ([]models.CoinHistory, error) {
	var history []models.CoinHistory
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

// ===== BADGE HISTORY REPOSITORY =====

// BadgeHistoryRepository handles badge history data access
type BadgeHistoryRepository struct {
	db *gorm.DB
}

func NewBadgeHistoryRepository(db *gorm.DB) *BadgeHistoryRepository {
	return &BadgeHistoryRepository{db: db}
}

func (r *BadgeHistoryRepository) Create(history *models.BadgeHistory) (*models.BadgeHistory, error) {
	if err := r.db.Create(history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *BadgeHistoryRepository) GetUserBadges(userID string) ([]models.BadgeHistory, error) {
	var badges []models.BadgeHistory
	if err := r.db.Where("user_id = ?", userID).Order("achieved_at DESC").Find(&badges).Error; err != nil {
		return nil, err
	}
	return badges, nil
}
